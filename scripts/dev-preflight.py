#!/usr/bin/env python3
"""Fast local feedback; never substitutes for full final-revision native/source QA."""
import argparse
import os
from pathlib import Path
import subprocess
import sys

ROOT = Path(__file__).resolve().parents[1]


def checked(command, capture=False):
    result = subprocess.run(command, cwd=ROOT, capture_output=capture, text=True)
    if result.returncode:
        if capture:
            print(result.stdout, end='')
            print(result.stderr, end='', file=sys.stderr)
        raise SystemExit(result.returncode)
    return result.stdout if capture else ''


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--package', action='append', help='Go package beginning ./; repeat for focused feedback')
    parser.add_argument('--static-only', action='store_true', help='Workflow and whitespace checks only')
    args = parser.parse_args()
    packages = args.package or ['./...']
    if any(not p.startswith('./') or any(c.isspace() for c in p) for p in packages):
        parser.error('packages must be local ./ paths without whitespace')
    checked(['git', 'diff', '--check', 'HEAD'])
    checked([sys.executable, str(ROOT / 'scripts/ci/generate_quality_workflow.py')])
    checked([sys.executable, str(ROOT / 'scripts/tests/ci-workflow.py')])
    checked([sys.executable, str(ROOT / 'scripts/tests/gate-probe.py')])
    checked([sys.executable, str(ROOT / 'scripts/tests/handoff.py')])
    if not args.static_only:
        tracked = checked(['git', 'ls-files', '-z', '--', '*.go'], capture=True).split('\0')
        extra = checked(['git', 'ls-files', '--others', '--exclude-standard', '-z', '--', '*.go'], capture=True).split('\0')
        files = sorted(set(filter(None, tracked + extra)))
        go = os.environ.get('QA_GO', 'go')
        gofmt = str(Path(go).with_name('gofmt' + Path(go).suffix))
        for i in range(0, len(files), 80):
            output = checked([gofmt, '-l'] + files[i:i + 80], capture=True)
            if output.strip():
                raise SystemExit('Unformatted Go files:\n' + output)
        checked([go, 'test', '-count=1'] + packages)
        checked([go, 'vet'] + packages)
    print('Developer preflight passed. Full final-revision native/source gates remain required.')


if __name__ == '__main__':
    main()
