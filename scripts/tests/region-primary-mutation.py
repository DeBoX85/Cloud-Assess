#!/usr/bin/env python3
"""Require region primary guards to reject compiling faults, with restored baselines."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (
    ('selected-identity', 'internal/plugins/region/primary.go', 'if !selected || name != c.SubscriptionName {',
     'if !selected && false || name != c.SubscriptionName && false {', './internal/plugins/region', 'TestPrimaryRejectsForeignScope'),
    ('unknown-sku-denominator', 'internal/plugins/region/primary.go', 'if confirmed := c.TotalSKUsChecked - c.UnknownSKUs; confirmed > 0 {\n\t\tsku =',
     'if confirmed := c.TotalSKUsChecked; confirmed > 0 {\n\t\tsku =', './internal/plugins/region', 'TestPrimarySourceScoresAndEveryCell'),
    ('aggregate-detail-work', 'internal/plugins/region/primary.go', 'if detailEntries > MaxDetailEntries {',
     'if detailEntries > MaxDetailEntries && false {', './internal/plugins/region', 'TestPrimaryAggregateDetailWorkRejectedBeforeProjection'),
    ('joined-separator-bytes', 'internal/plugins/region/primary.go', 'textBytes += max(0, len(details)-1) * 2',
     'textBytes += 0', './internal/plugins/region', 'TestPrimaryJoinedSeparatorsCountBeforeProjection'),
    ('comparison-work-limit', 'internal/plugins/region/primary.go', 'if len(input) > MaxComparisons || len(subscriptions) > MaxSubscriptions {',
     'if len(input) > MaxComparisons && false || len(subscriptions) > MaxSubscriptions {', './internal/plugins/region', 'TestPrimaryWorkLimits'),
)



def run(directory, package, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$', package],
                          cwd=directory, capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-primary-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, filename, before, after, package, test in MUTATIONS:
        baseline = run(directory, package, test)
        if baseline.returncode:
            raise SystemExit('Region primary baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / filename
        original = path.read_text(encoding='utf-8')
        if original.count(before) != (2 if name == 'aggregate-detail-work' else 1):
            raise SystemExit('Region primary anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, package, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout:
            raise SystemExit('Region primary mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, package, test)
        if restored.returncode:
            raise SystemExit('Restored Region primary failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region primary mutation rejected:', name)
