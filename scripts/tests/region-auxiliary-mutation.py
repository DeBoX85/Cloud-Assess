#!/usr/bin/env python3
"""Require region auxiliary guards to reject compiling faults, with restored baselines."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ('selected-identity', 'internal/plugins/region/auxiliary.go', 'return selected && selectedName == name',
     'return selected && selectedName == name || true', './internal/plugins/region', 'TestAuxiliarySelectedScope'),
    ('row-work-limit', 'internal/plugins/region/auxiliary.go', 'if len(input) > MaxAuxRows {',
     'if len(input) > MaxAuxRows && false {', './internal/plugins/region', 'TestAuxiliaryWorkLimits'),
    ('aggregate-text-limit', 'internal/plugins/region/auxiliary.go', 'return *total <= auxTextBudget',
     'return *total <= auxTextBudget || true', './internal/plugins/region', 'TestAuxiliaryWorkLimits'),
)



def run(directory, package, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$', package],
                          cwd=directory, capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-auxiliary-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, filename, before, after, package, test in MUTATIONS:
        baseline = run(directory, package, test)
        if baseline.returncode:
            raise SystemExit('Region auxiliary baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / filename
        original = path.read_text(encoding='utf-8')
        if original.count(before) != (2 if name == 'row-work-limit' else 1):
            raise SystemExit('Region auxiliary anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, package, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout:
            raise SystemExit('Region auxiliary mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, package, test)
        if restored.returncode:
            raise SystemExit('Restored Region auxiliary failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region auxiliary mutation rejected:', name)
