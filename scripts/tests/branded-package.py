#!/usr/bin/env python3
"""Build a real custom candidate and run all installed package/preflight checks."""
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[2]
GO = os.environ.get('QA_GO', 'go')


def run(args):
    result = subprocess.run(args, cwd=ROOT, capture_output=True, text=True, timeout=180)
    if result.returncode:
        raise AssertionError(result.stdout + result.stderr)
    sys.stderr.write(result.stderr)
    return result.stdout


with tempfile.TemporaryDirectory(prefix='cloud-assess-custom-package-') as directory:
    directory = Path(directory)
    profile = {'schemaVersion': 1, 'productName': 'Example "<&>" Æ Toolkit', 'cliName': 'example-cloud',
               'reportTitle': 'Example Cloud Assessment', 'reportFilePrefix': 'example_report',
               'websiteURL': 'https://example.test/tool?a=1&b=2'}
    path = directory / 'public profile.json'
    path.write_text(json.dumps(profile), encoding='utf-8')
    binary = run([GO, 'run', './tools/brand-build', '--profile', str(path), '--output', str(directory / 'build with spaces')]).strip()
    print(run([sys.executable, 'scripts/tests/package-candidate.py', '--binary', binary, '--go', GO]))
    print('Custom native package: all package cases and installed offline CLI checks passed')
