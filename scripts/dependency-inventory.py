#!/usr/bin/env python3
"""Collect source evidence, not inferred license classifications or release approval."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[1]
TARGETS = ('linux/amd64', 'windows/amd64')


def run(args, cwd=ROOT, env=None):
    completed = subprocess.run(args, cwd=cwd, env=env, capture_output=True, text=True, timeout=300)
    if completed.returncode:
        raise RuntimeError(f'command failed ({completed.returncode}): {args[0]} {args[1]}\n{completed.stderr[-4000:]}')
    return completed.stdout


def objects(text):
    decoder = json.JSONDecoder()
    while text.strip():
        text = text.lstrip()
        value, end = decoder.raw_decode(text)
        yield value
        text = text[end:]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def evidence(directory):
    rows = []
    for path in sorted(directory.iterdir(), key=lambda p: p.name):
        if path.is_file() and re.match(r'^(license|licence|copying|notice)(?:[._-].*)?$', path.name, re.I):
            if path.is_symlink():
                raise ValueError(f'symlink notice rejected: {path.name}')
            data = path.read_bytes()
            if not data.strip():
                raise ValueError(f'empty notice: {path.name}')
            rows.append({'file': path.name, 'sha256': digest(data), 'text': data.decode('utf-8')})
    if not any(re.match(r'^(license|licence|copying)(?:[._-].*)?$', row['file'], re.I) for row in rows):
        raise ValueError(f'no license evidence in {directory}')
    return rows


def collect(go):
    env = dict(os.environ, GOWORK='off', GOFLAGS='-mod=readonly')
    modules = [m for m in objects(run([go, 'list', '-m', '-json', 'all'], env=env)) if not m.get('Main')]
    if any(m.get('Replace') or not m.get('Version') for m in modules):
        raise ValueError('replaced/unversioned modules require an explicit provenance policy')
    usage = {m['Path']: [] for m in modules}
    for target in TARGETS:
        goos, goarch = target.split('/')
        target_env = dict(env, GOOS=goos, GOARCH=goarch, CGO_ENABLED='0')
        for package in objects(run([go, 'list', '-buildvcs=false', '-deps', '-json', './cmd/cloud-assess'], env=target_env)):
            module = package.get('Module', {})
            if module.get('Path') in usage and target not in usage[module['Path']]:
                usage[module['Path']].append(target)
    # Download in a disposable empty module so additional graph sums cannot mutate go.sum.
    with tempfile.TemporaryDirectory(prefix='dependency-inventory-') as temp:
        Path(temp, 'go.mod').write_text('module inventory.invalid/evidence\n\ngo 1.26.8\n')
        downloads = list(objects(run([go, 'mod', 'download', '-json'] +
                                    [m['Path'] + '@' + m['Version'] for m in modules], cwd=temp, env=env)))
    downloaded = {(m['Path'], m['Version']): m for m in downloads}
    notices = []
    entries = []
    for module in sorted(modules, key=lambda m: m['Path']):
        item = downloaded[(module['Path'], module['Version'])]
        if item.get('Error') or not item.get('Sum') or not item.get('GoModSum'):
            raise ValueError(f'incomplete module download: {module["Path"]}')
        files = evidence(Path(item['Dir']))
        label = module['Path'] + '@' + module['Version']
        notices.append((label, files))
        entries.append({'path': module['Path'], 'version': module['Version'],
                        'direct': not module.get('Indirect', False), 'cliTargets': usage[module['Path']],
                        'sum': item['Sum'], 'goModSum': item['GoModSum'],
                        'notices': [{k: f[k] for k in ('file', 'sha256')} for f in files]})
    goroot = Path(run([go, 'env', 'GOROOT'], env=env).strip())
    go_version = run([go, 'env', 'GOVERSION'], env=env).strip()
    if go_version != 'go1.26.8':
        raise ValueError(f'inventory requires pinned go1.26.8; got {go_version}')
    std = evidence(goroot)
    patents = goroot / 'PATENTS'
    if patents.is_file():
        data = patents.read_bytes()
        std.append({'file': 'PATENTS', 'sha256': digest(data), 'text': data.decode('utf-8')})
    notices.append(('Go standard library/toolchain ' + go_version, std))
    pins = {}
    source = (ROOT / 'internal/rules/provenance.go').read_text()
    for key in ('ReferenceRepositoryCommit', 'APRLCommit', 'AORSnapshotTree', 'CustomRulesSnapshotTree'):
        pins[key] = re.search(r'\b' + key + r'\s*=\s*"([0-9a-f]{40})"', source)[1]
    expected = [('internal/rules/upstream/aprl', pins['APRLCommit']),
                ('internal/rules/upstream/orphan-resources', pins['AORSnapshotTree']),
                ('internal/rules/custom', pins['CustomRulesSnapshotTree'])]
    for path, pin in expected:
        if run(['git', 'rev-parse', 'HEAD:' + path]).strip() != pin:
            raise ValueError(f'committed source pin mismatch: {path}')
        if run(['git', 'status', '--porcelain', '--', path]).strip():
            raise ValueError(f'dirty bundled content: {path}')
    if run(['git', '-C', 'internal/rules/upstream/aprl', 'rev-parse', 'HEAD']).strip() != pins['APRLCommit']:
        raise ValueError('APRL checkout does not match pin')
    bundled = [{'name': 'APRL normal catalog', 'path': 'internal/rules/upstream/aprl/azure-resources', 'revision': pins['APRLCommit']},
               {'name': 'AOR snapshot', 'path': 'internal/rules/upstream/orphan-resources', 'revision': pins['AORSnapshotTree']},
               {'name': 'Custom reference rules', 'path': 'internal/rules/custom/azure-resources', 'revision': pins['CustomRulesSnapshotTree']},
               {'name': 'Known SKU snapshot', 'path': 'internal/skus/known_skus.yaml',
                'sha256': digest((ROOT / 'internal/skus/known_skus.yaml').read_bytes())}]
    own = evidence(ROOT)
    # Attribution file contains full source-family licenses, not merely a guessed SPDX label.
    for filename in ('THIRD_PARTY_LICENSES.md',):
        data = (ROOT / filename).read_bytes()
        own.append({'file': filename, 'sha256': digest(data), 'text': data.decode('utf-8')})
    notices.append(('Cloud Assess and incorporated source/rule attribution', own))
    aprl_license = evidence(ROOT / 'internal/rules/upstream/aprl')
    notices.append(('Pinned APRL source license', aprl_license))
    result = {'schemaVersion': 1, 'scope': 'selected module graph, CLI targets, standard library and bundled data',
              'goVersion': go_version, 'targets': list(TARGETS), 'referencePins': pins,
              'inputs': {p: digest((ROOT / p).read_bytes()) for p in ('go.mod', 'go.sum', 'internal/rules/embedded.go', 'internal/skus/skus.go')},
              'modules': entries, 'bundledContent': bundled,
              'standardLibraryNotices': [{k: f[k] for k in ('file', 'sha256')} for f in std]}
    sections = ['# Generated dependency and incorporated-content notices\n\nGenerated by scripts/dependency-inventory.py. Preserve alongside distribution artifacts.\n']
    for label, files in notices:
        for item in files:
            sections.append('\n## ' + label + ' / ' + item['file'] + '\n\nSHA-256: ' + item['sha256'] + '\n\n' + item['text'] + '\n')
    return json.dumps(result, indent=2, ensure_ascii=False) + '\n', ''.join(sections)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--check', action='store_true')
    parser.add_argument('--go', default=os.environ.get('QA_GO', 'go'))
    args = parser.parse_args()
    manifest, notices = collect(args.go)
    outputs = {'docs/dependencies/inventory.json': manifest,
               'docs/dependencies/NOTICES.md': notices}
    if args.check:
        stale = [p for p, content in outputs.items() if not (ROOT / p).is_file() or (ROOT / p).read_bytes() != content.encode('utf-8')]
        if stale:
            raise SystemExit('stale inventory: ' + ', '.join(stale) + '; regenerate and review')
    else:
        # Build all evidence before touching either output.
        for path, content in outputs.items():
            (ROOT / path).parent.mkdir(parents=True, exist_ok=True)
            (ROOT / path).write_bytes(content.encode('utf-8'))
    print(f'Dependency inventory {"verified" if args.check else "generated"}.')


if __name__ == '__main__':
    main()
