#!/usr/bin/env python3
"""Require compiling inventory-calculation faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/inventory_aggregation.go'
MUTATIONS = (
    ('empty-selected-name', 'if name == "" {', 'if name == "" && false {',
     'TestAggregationSelectedScopeAndMalformed', 1),
    ('selected-resource-preflight', '|| !selected {', '|| (!selected && false) {',
     'TestAggregationSelectedScopeAndMalformed', 1),
    ('embedded-subscription', 'if !inventoryIdentity(resource.ID, resource.SubscriptionID) {',
     'if !inventoryIdentity(resource.ID, resource.SubscriptionID) && false {',
     'TestAggregationSelectedScopeAndMalformed', 1),
    ('row-work-limit', 'if len(resources) > MaxInventoryRows {',
     'if len(resources) > MaxInventoryRows && false {', 'TestAggregationWorkAndEntryLimits', 1),
    ('global-entry-limit', 'entries >= MaxInventoryEntries',
     'entries >= MaxInventoryEntries && false', 'TestAggregationWorkAndEntryLimits', 1),
    ('decoded-text-limit', 'if !auxiliaryText(label, &decodedBytes) {',
     'if !auxiliaryText(label, &decodedBytes) && false {', 'TestAggregationDecodedAndProjectedText', 2),
    ('projected-key-text', 'len(key) > auxTextBudget-outputBytes',
     'len(key) > auxTextBudget-outputBytes && false', 'TestAggregationDecodedAndProjectedText', 1),
    ('ascii-space-normalization', 'strings.ReplaceAll(region, " ", "")',
     'region', 'TestAggregationSourceNormalization', 1),
    ('raw-sku-keys', 'resourceType, location, resource.SKUName, insert)',
     'resourceType, location, strings.TrimSpace(resource.SKUName), insert)',
     'TestAggregationLiteralCountsAndIgnoredFields', 2),
    ('resource-not-capacity', 'values[key]++', 'values[key] += 2',
     'TestAggregationCapturedMapsAndIdentityCorrection', 1),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-aggregation-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Aggregation baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / FILE
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Aggregation mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout or 'panic:' in outcome.stdout + outcome.stderr or '[build failed]' in outcome.stdout:
            raise SystemExit('Aggregation mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored aggregation baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region aggregation mutation rejected:', name)
