#!/usr/bin/env python3
"""Require compiling availability-calculation faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/availability.go'
MUTATIONS = (
    ('selected-scope', '|| !selected ||', '|| (!selected && false) ||', 'TestAvailabilityRuntimeScopeAndEvidence', 1),
    ('bound-evidence', 'evidence.SubscriptionID != id', '(evidence.SubscriptionID != id && false)', 'TestAvailabilityRuntimeScopeAndEvidence', 1),
    ('input-work', 'n > MaxAvailabilityEntries-entries', '(n > MaxAvailabilityEntries-entries && false)', 'TestAvailabilityRuntimeWorkAndText', 1),
    ('decoded-text', 'return auxiliaryText(label, &decodedBytes)', 'return len(label) <= MaxLabelBytes', 'TestAvailabilityRuntimeWorkAndText', 1),
    ('output-list', 'len(*list) >= MaxDetails', '(len(*list) >= MaxDetails && false)', 'TestAvailabilityRuntimeOutputBounds', 1),
    ('output-label', '!auxiliaryText(detail, &outputBytes)', '(!auxiliaryText(detail, &outputBytes) && false)', 'TestAvailabilityRuntimeOutputBounds', 1),
    ('raw-sku', 'identifier := key + ":" + raw', 'identifier := key + ":" + strings.ToLower(strings.TrimSpace(raw))', 'TestAvailabilityRuntimeCapturedComparisons', 1),
    ('confirmed-denominator', 'c.TotalSKUsChecked - c.UnknownSKUs', 'c.TotalSKUsChecked', 'TestAvailabilityRuntimeCapturedComparisons', 1),
    ('unknown-health', 'warn("availability_sku_unknown")', '_ = skuEvidence', 'TestAvailabilityRuntimeUnknownAndCompleteness', 1),
    ('cancellation', 'if ctx.Err() != nil {', 'if ctx.Err() != nil && false {', 'TestAvailabilityRuntimeOwnershipAndCancellation', 5),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-region-availability-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Availability baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / FILE
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Availability mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        if outcome.returncode == 0 or '--- FAIL: ' + test not in outcome.stdout or 'panic:' in outcome.stdout + outcome.stderr or '[build failed]' in outcome.stdout:
            raise SystemExit('Availability mutation survived or failed outside named assertions: ' + name + '\n' + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored availability baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region availability mutation rejected:', name)
