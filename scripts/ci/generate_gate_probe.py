#!/usr/bin/env python3
"""Copy production gate wiring into branch-isolated hosted contract exercises."""
import argparse
import importlib.util
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUTPUT = ROOT / '.github/workflows/gate-contract-probe.yml'


def production_jobs():
    spec = importlib.util.spec_from_file_location('quality_generator', ROOT / 'scripts/ci/generate_quality_workflow.py')
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    _, jobs = module.split_jobs((ROOT / '.github/workflows/test.yml').read_text(encoding='utf-8'))
    return jobs


def generate():
    jobs = production_jobs()
    out = '''name: Hosted gate contract probes

on:
  push:
    branches:
      - 'qa/gate-contract-*-success'
      - 'qa/gate-contract-*-failure'
      - 'qa/gate-contract-*-skipped'

permissions:
  contents: read

concurrency:
  group: hosted-gate-probe-${{ github.ref }}
  cancel-in-progress: true

jobs:
'''
    for platform, original, runner in [('linux', 'quality', 'ubuntu-24.04'), ('windows', 'windows-validation', 'windows-2022')]:
        py = 'python3' if platform == 'linux' else 'python'
        workers = [platform + '-native'] + [platform + '-controls-' + str(i) for i in (1, 2, 3)]
        for worker in workers:
            out += '  ' + worker + ':\n'
            controlled = worker == platform + '-controls-1'
            if controlled:
                out += "    if: ${{ ! endsWith(github.ref_name, '-skipped') }}\n"
            out += '    runs-on: ' + runner + '\n    timeout-minutes: 5\n    steps:\n'
            out += '      - name: Controlled worker outcome\n'
            if controlled:
                out += '        env:\n          PROBE_BRANCH: ${{ github.ref_name }}\n'
                out += '        run: ' + py + ' -c "import os,sys; sys.exit(1 if os.environ[\'PROBE_BRANCH\'].endswith(\'-failure\') else 0)"\n\n'
            else:
                out += '        run: ' + py + ' -c "print(\'Controlled worker success\')"\n\n'
        # Body, needs, always(), checkout and validator command are copied verbatim.
        out += '  ' + platform + '-probe-gate:\n' + jobs[original].rstrip() + '\n\n'
    return out.rstrip() + '\n'


def verify(text):
    if text.replace('\r\n', '\n') != generate():
        raise ValueError('hosted probe no longer matches production gate wiring')


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--write', action='store_true')
    args = parser.parse_args()
    if args.write:
        OUTPUT.write_bytes(generate().encode('utf-8'))
    else:
        verify(OUTPUT.read_text(encoding='utf-8'))
        print('Hosted gate probe wiring matches production; isolated push branches only')


if __name__ == '__main__':
    main()
