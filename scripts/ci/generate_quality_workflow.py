#!/usr/bin/env python3
"""Generate a conserved parallel workflow without a new YAML dependency."""
import argparse
import hashlib
import json
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
BASELINE = ROOT / 'scripts/ci/testdata/quality-before-20261017.yml.txt'
PLAN = ROOT / 'scripts/ci/quality-plan.json'
WORKFLOW = ROOT / '.github/workflows/test.yml'
BASELINE_SHA256 = '1c9e08848af4aabc9b29402a0d0f264efbd679c019466c628d14048e8032e5ee'
ACTIVE_GO_VERSION = '1.26.10'


def split_jobs(text):
    """Parse only the fixed, generated indentation grammar, not arbitrary YAML."""
    header, body = text.replace('\r\n', '\n').split('jobs:\n', 1)
    matches = list(re.finditer(r'^  ([a-z][a-z0-9-]*):\n', body, re.M))
    jobs = {}
    for i, match in enumerate(matches):
        key = match.group(1)
        if key in jobs:
            raise ValueError('duplicate job identity')
        jobs[key] = body[match.end():matches[i + 1].start() if i + 1 < len(matches) else len(body)]
    return header, jobs


def split_steps(body):
    preamble, body = body.split('    steps:\n', 1)
    matches = list(re.finditer(r'^      - name: (.+)\n', body, re.M))
    steps = {}
    for i, match in enumerate(matches):
        name = match.group(1)
        if name in steps:
            raise ValueError('duplicate step identity')
        steps[name] = body[match.start():matches[i + 1].start() if i + 1 < len(matches) else len(body)].rstrip() + '\n'
    return preamble, steps


def read_baseline():
    # A native Windows checkout may use CRLF; the authority is canonical Git text.
    raw = BASELINE.read_bytes().replace(b'\r\n', b'\n')
    plan = json.loads(PLAN.read_text(encoding='utf-8'))
    if hashlib.sha256(raw).hexdigest() != BASELINE_SHA256 or plan['baseline_sha256'] != BASELINE_SHA256:
        raise ValueError('frozen original workflow changed')
    header, jobs = split_jobs(raw.decode('utf-8'))
    if set(jobs) != {'quality', 'windows-validation'}:
        raise ValueError('unexpected baseline jobs')
    return header, {key: split_steps(body) for key, body in jobs.items()}, plan


def generate():
    header, baseline, plan = read_baseline()
    out = header + 'jobs:\n'
    for platform, original in (('linux', 'quality'), ('windows', 'windows-validation')):
        preamble, steps = baseline[original]
        groups = plan['groups'][platform]
        controls = [name for name in steps if re.search(r'compiling|mutation|retained source', name)]
        assigned = [name for group in groups for name in group]
        if len(groups) != 3 or any(not group for group in groups) or len(assigned) != len(set(assigned)) or set(assigned) != set(controls):
            raise ValueError('control suites must occur exactly once on their original platform')
        setup = ['Checkout', 'Set up Go']
        if list(steps)[:2] != setup:
            raise ValueError('baseline setup changed')
        additions = plan.get('additions', {})
        if set(additions) != {'linux', 'windows'} or not isinstance(additions[platform], list):
            raise ValueError('explicit per-platform additive inventory required')
        extras = {}
        for addition in additions[platform]:
            if not isinstance(addition, dict) or set(addition) != {'name', 'script', 'worker'}:
                raise ValueError('additive checks have only name, script and worker fields')
            name, script, worker = addition['name'], addition['script'], addition['worker']
            if (not isinstance(name, str) or not re.fullmatch(r'[A-Za-z0-9][A-Za-z0-9 ()/-]{0,119}', name) or
                    name in steps or any(name == prior['name'] for group in extras.values() for prior in group) or
                    worker not in ('native', 'controls-1', 'controls-2', 'controls-3') or
                    not isinstance(script, str) or not re.fullmatch(r'scripts/tests/[a-z0-9][a-z0-9-]*\.py', script) or
                    not (ROOT / script).is_file()):
                raise ValueError('invalid or duplicate additive required check')
            extras.setdefault(platform + '-' + worker, []).append(addition)
        native = [name for name in steps if name not in controls]
        for worker, names in [(platform + '-native', native)] + [(platform + '-controls-' + str(i + 1), setup + group) for i, group in enumerate(groups)]:
            out += '  ' + worker + ':\n' + preamble + '    steps:\n'
            py = 'python3' if platform == 'linux' else 'python'
            # Run the separate bounded Build/Admission proofs before the full
            # native suite. They remain additive and never replace that suite.
            early = [e for e in extras.get(worker, []) if e['script'] in ('scripts/tests/required-inventory-build.py', 'scripts/tests/required-inventory-admission.py', 'scripts/tests/mutation-assertions.py')]
            late = [e for e in extras.get(worker, []) if e not in early]
            early_blocks = ''.join('      - name: ' + e['name'] + '\n        run: ' + py + ' ' + e['script'] + '\n\n' for e in early)
            extra_blocks = ''.join('      - name: ' + e['name'] + '\n        run: ' + py + ' ' + e['script'] + '\n\n' for e in late)
            inserted = False
            for index, name in enumerate(names):
                # Fail quickly on workflow drift, before expensive suites begin.
                if index == 2:
                    out += '      - name: Verify parallel workflow conservation and gate guards\n        run: ' + py + ' scripts/tests/ci-workflow.py\n\n'
                    out += early_blocks
                if name == 'Retain synthetic QA failure evidence':
                    out += extra_blocks
                    inserted = True
                block = steps[name]
                # Keep the historical baseline immutable; apply the approved
                # security toolchain override to setup blocks only.
                if name == 'Set up Go':
                    if block.count("go-version: '1.26.9'") != 1:
                        raise ValueError('baseline Go setup changed')
                    block = block.replace("go-version: '1.26.9'",
                                          "go-version: '" + ACTIVE_GO_VERSION + "'")
                # Keep build caching, but require this final native run to execute tests.
                if name == 'Test Go packages on Windows':
                    block = block.replace('run: go test ./...', 'run: go test -count=1 ./...')
                if name == 'Verify coverage floor':
                    block = block.replace('go test -covermode=atomic', 'go test -count=1 -covermode=atomic')
                out += block + '\n'
            if not inserted:
                out += extra_blocks
        gate = 'quality' if platform == 'linux' else 'windows-validation'
        workers = [platform + '-native'] + [platform + '-controls-' + str(i) for i in (1, 2, 3)]
        out += '  ' + gate + ':\n'
        out += '    if: ${{ always() }}\n    needs: [' + ', '.join(workers) + ']\n'
        out += '    runs-on: ubuntu-24.04\n    timeout-minutes: 5\n    steps:\n'
        out += '      - name: Checkout gate validator\n        uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1\n        with:\n          persist-credentials: false\n\n'
        out += '      - name: Require every native worker to succeed\n        env:\n          NEEDS_JSON: ${{ toJSON(needs) }}\n        run: python3 scripts/ci/required_gate.py --platform ' + platform + '\n\n'
    return out.rstrip() + '\n'


def verify(text):
    if text.replace('\r\n', '\n') != generate():
        raise ValueError('workflow differs from the conserved generated plan')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--write', action='store_true')
    args = parser.parse_args()
    if args.write:
        WORKFLOW.write_bytes(generate().encode('utf-8'))
    else:
        verify(WORKFLOW.read_text(encoding='utf-8'))
        print('Parallel workflow: all original platform checks conserved; protected gates strict')


if __name__ == '__main__':
    main()
