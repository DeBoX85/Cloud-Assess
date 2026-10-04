#!/usr/bin/env python3
"""Require compiling latency-calculation faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/latency.go'
MUTATIONS = (
    ('selected-scope', 'if !selected || name != c.SubscriptionName {', 'if (!selected && false) || name != c.SubscriptionName {', 'TestLatencyRuntimeAdmission', 1, 'internal/plugins/region/primary.go'),
    ('pinned-bytes', 'hex.EncodeToString(hash[:]) != latencyDataHash', '(hex.EncodeToString(hash[:]) != latencyDataHash && false)', 'TestLatencyRuntimeDataRejection', 1, FILE),
    ('input-work', 'n > MaxLatencyEntries-entries', '(n > MaxLatencyEntries-entries && false)', 'TestLatencyRuntimeBudgets', 1, FILE),
    ('data-label', 'regionID.MatchString(s)', '(regionID.MatchString(s) || true)', 'TestLatencyRuntimeDataRejection', 1, FILE),
    ('negative-data', '|| ms < 0 ||', '|| (ms < 0 && false) ||', 'TestLatencyRuntimeDataRejection', 1, FILE),
    ('forward-preference', 'data.Matrix[c.SourceRegion][c.TargetRegion]', 'data.Matrix[c.TargetRegion][c.SourceRegion]', 'TestLatencyRuntimeCapturedComparisons', 1, FILE),
    ('reverse-measurement', '} else if ms, ok := data.Matrix[c.TargetRegion][c.SourceRegion]; ok {\n\t\t\t\tc.AvgLatencyMs = ms', '} else if ms, ok := data.Matrix[c.TargetRegion][c.SourceRegion]; ok {\n\t\t\t\tc.AvgLatencyMs = ms * 2', 'TestLatencyRuntimeCapturedComparisons', 1, FILE),
    ('cluster-direction', 'totals[a+":"+b]', 'totals[b+":"+a]', 'TestLatencyRuntimeCapturedComparisons', 1, FILE),
    ('cluster-average', 'value.sum / float64(value.count)', 'value.sum', 'TestLatencyRuntimeCapturedComparisons', 1, FILE),
    ('unknown-health', 'warn("latency_unknown")', '_ = exists', 'TestLatencyRuntimeCapturedComparisons', 1, FILE),
    ('estimated-health', 'warn("latency_estimated")', '_ = value', 'TestLatencyRuntimeCapturedComparisons', 1, FILE),
    ('zone-ownership', 'maps.Clone(before.TargetZoneMappings)', 'before.TargetZoneMappings', 'TestLatencyRuntimeOwnershipAndCancellation', 1, FILE),
    ('terminal-cancellation', 'if ctx.Err() != nil {', 'if ctx.Err() != nil && false {', 'TestLatencyRuntimeOwnershipAndCancellation', 1, FILE),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-latency-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count, file in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Latency baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / file
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Latency mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout or 'panic:' in outcome.stdout + outcome.stderr or '[build failed]' in outcome.stdout:
            raise SystemExit('Latency mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored latency baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region latency mutation rejected:', name)
