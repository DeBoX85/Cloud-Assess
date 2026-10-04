#!/usr/bin/env python3
"""Reject compiling service-sheet faults, with named failures and restored baselines."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/service.go'
MUTATIONS = (
    ('comparison-selected-identity', 'if !auxiliarySelected(scope, c.SubscriptionID, c.SubscriptionName) || inventory != nil && !contributors[strings.ToLower(c.SubscriptionID)] {',
     'if false {', 'TestServiceSelectedScope'),
    ('replicated-row-limit', 'if len(inventory.ResourceTypes) > MaxServiceRows/len(targets) {',
     'if len(inventory.ResourceTypes) > MaxServiceRows/len(targets) && false {', 'TestServiceWorkBounds'),
    ('global-entry-limit', 'if n > MaxServiceEntries-entries {',
     'if n > MaxServiceEntries-entries && false {', 'TestServiceWorkBounds'),
    ('joined-cell-delimiters', 'units := max(0, len(group)-1) * 2',
     'units := 0', 'TestServiceJoinedAndReplicatedTextBounds'),
    ('replicated-text-limit', 'if *total > serviceTextBudget {',
     'if *total > serviceTextBudget && false {', 'TestServiceJoinedAndReplicatedTextBounds'),
    ('sheet-collision', 'if sheets[name] {',
     'if sheets[name] && false {', 'TestServiceMalformedAndCollision'),
    ('pinned-disk-provider', '"microsoft.compute/disks"',
     '"microsoft.compute/unsupported"', 'TestServiceRegistryAndDetailSemantics'),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-service-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Service baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / FILE
        original = path.read_text(encoding='utf-8')
        if original.count(before) != 1:
            raise SystemExit('Service mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout:
            raise SystemExit('Service mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored service baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region service mutation rejected:', name)
