#!/usr/bin/env python3
"""Require compiling CostComparison faults to fail named tests and restored baselines."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/cost.go'
MUTATIONS = (
    ('empty-selected-name', 'if name == "" {', 'if name == "" && false {',
     'TestCostSelectedScopeAndMalformed', 1),
    ('selected-contributor', '|| !selected ||', '|| false ||',
     'TestCostSelectedScopeAndMalformed', 1),
    ('meter-work-limit', 'len(input.MeterInputs) > MaxCostMeters',
     'len(input.MeterInputs) > MaxCostMeters && false', 'TestCostWorkLimits', 1),
    ('global-entry-limit', 'if n > MaxCostEntries-entries {',
     'if n > MaxCostEntries-entries && false {', 'TestCostWorkLimits', 1),
    ('decoded-text-limit', 'if !auxiliaryText(s, &decodedBytes) {',
     'if !auxiliaryText(s, &decodedBytes) && false {', 'TestCostTextLimitsAndUnicode', 2),
    ('first-received-meter', 'if _, exists := firstMeter[key]; !exists {',
     'if _, exists := firstMeter[key]; !exists || true {', 'TestCostMetadataReceivedOrder', 1),
    ('first-price-metadata', 'if _, exists := metadata[id]; !exists {',
     'if _, exists := metadata[id]; !exists || true {', 'TestCostMetadataReceivedOrder', 1),
    ('zero-price-is-absent', 'if price > 0 {', 'if price >= 0 {',
     'TestCostPriceSemantics', 1),
    ('nonfinite-price', 'math.IsNaN(price)', 'false',
     'TestCostSelectedScopeAndMalformed', 1),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-cost-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Cost baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / FILE
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Cost mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout:
            raise SystemExit('Cost mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored cost baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region cost mutation rejected:', name)
