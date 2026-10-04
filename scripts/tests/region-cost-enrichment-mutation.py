#!/usr/bin/env python3
"""Require compiling cost-enrichment-calculation faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/cost_enrichment.go'
MUTATIONS = (
    ('tiny-history', 'meter.HistoricalCost > 0 && meter.HistoricalCost < MinCostWeight', '(meter.HistoricalCost > 0 && meter.HistoricalCost < MinCostWeight && false)', 'TestCostRuntimeHistoricalWeightFloor', 1, FILE),
    ('entry-budget', 'n > MaxCostEntries-entries', '(n > MaxCostEntries-entries && false)', 'TestCostRuntimeBudgets', 1, FILE),
    ('region-budget', 'len(regions) > MaxCostRegions', '(len(regions) > MaxCostRegions && false)', 'TestCostRuntimeBudgets', 1, FILE),
    ('history-completeness', '!history.Complete', '(!history.Complete && false)', 'TestCostRuntimeHealthAndSubscriptionWeights', 1, FILE),
    ('free-roundoff', 'math.Max(-100, weighted/total)', 'weighted/total', 'TestCostRuntimeFreeTargetRoundoff', 1, FILE),
    ('selected-scope', 'if !selected || name != c.SubscriptionName {', 'if (!selected && false) || name != c.SubscriptionName {', 'TestCostRuntimeAdmission', 1, 'internal/plugins/region/primary.go'),
    ('evidence-scope', '!selected || id != strings.ToLower(id)', '(!selected && false) || id != strings.ToLower(id)', 'TestCostRuntimeAdmission', 1, FILE),
    ('duplicate-history', 'if seen[meter.MeterID] {', 'if seen[meter.MeterID] && false {', 'TestCostRuntimeAdmission', 1, FILE),
    ('negative-values', 'n >= 0', '(n >= 0 || true)', 'TestCostRuntimeAdmission', 1, FILE),
    ('work-budget', 'n > MaxCostWork-work', '(n > MaxCostWork-work && false)', 'TestCostRuntimeWorkBudget', 1, FILE),
    ('weighted-average', 'weighted/total', 'weighted', 'TestCostRuntimeCapturedComparisons', 1, FILE),
    ('zero-history', 'weight = 1', 'weight = 0', 'TestCostRuntimeCapturedComparisons', 1, FILE),
    ('free-denominator', 'total += weight', 'if diff != 0 { total += weight }', 'TestCostRuntimeCapturedComparisons', 1, FILE),
    ('exact-threshold', 'source < 0.0001', 'source <= 0.0001', 'TestCostRuntimeCapturedComparisons', 2, FILE),
    ('stale-state', 'c.AvgCostDifference, c.HasCostData = 0, false', 'c.AvgCostDifference, c.HasCostData = before.AvgCostDifference, before.HasCostData', 'TestCostRuntimeCapturedComparisons', 1, FILE),
    ('partial-health', 'warn("cost_pricing_partial")', '_ = history', 'TestCostRuntimeHealthAndSubscriptionWeights', 1, FILE),
    ('zone-ownership', 'c.TargetZoneMappings = maps.Clone(before.TargetZoneMappings)', 'c.TargetZoneMappings = before.TargetZoneMappings; _ = maps.Clone(before.TargetZoneMappings)', 'TestCostRuntimeOwnershipAndCancellation', 1, FILE),
    ('terminal-cancellation', 'if ctx.Err() != nil {\n\t\treturn nil, ctx.Err()\n\t}\n\treturn result, nil', 'if ctx.Err() != nil && false {\n\t\treturn nil, ctx.Err()\n\t}\n\treturn result, nil', 'TestCostRuntimeOwnershipAndCancellation', 1, FILE),
    ('deterministic-sum', 'ordered[id][i].MeterID < ordered[id][j].MeterID', 'ordered[id][i].MeterID > ordered[id][j].MeterID', 'TestCostRuntimeDeterministicSum', 1, FILE),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-cost-enrichment-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count, file in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Cost enrichment baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / file
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Cost enrichment mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout or 'panic:' in outcome.stdout + outcome.stderr or '[build failed]' in outcome.stdout:
            raise SystemExit('Cost enrichment mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored cost-enrichment baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region cost-enrichment mutation rejected:', name)
