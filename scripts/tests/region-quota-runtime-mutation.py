#!/usr/bin/env python3
"""Require compiling quota-calculation faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/quota_runtime.go'
MUTATIONS = (
    ('exact-threshold', 'available*100 < usage.Limit*15', 'available*100 <= usage.Limit*15', 'TestQuotaRuntimeArithmeticAndDisplay', FILE),
    ('provider-filter', 'return strings.HasSuffix(resource, "PerServer") || strings.HasSuffix(resource, "PerDatabase")', 'return false', 'TestQuotaRuntimeProviderFilters', FILE),
    ('selected-scope', 'return request, selected', 'return request, selected || true', 'TestQuotaRuntimeScopeRequestsAndCancellation', FILE),
    ('duplicate-query', '!valid || expected[key]', '!valid || (expected[key] && false)', 'TestQuotaRuntimeScopeRequestsAndCancellation', FILE),
    ('duplicate-raw-identity', 'if seen[usage.ResourceName] {', 'if seen[usage.ResourceName] && false {', 'TestQuotaRuntimeAdmissionCorrectionsAndBounds', FILE),
    ('negative-count', 'usage.Current < 0 ||', '(usage.Current < 0 && false) ||', 'TestQuotaRuntimeAdmissionCorrectionsAndBounds', FILE),
    ('filtered-work-budget', 'len(response.Usages) > MaxAuxRows-count', '(len(response.Usages) > MaxAuxRows-count && false)', 'TestQuotaRuntimeAdmissionCorrectionsAndBounds', FILE),
    ('filtered-label-budget', '!auxiliaryText(usage.ResourceName, &textBytes) || !auxiliaryText(usage.LocalizedName, &textBytes)', '(!auxiliaryText(usage.ResourceName, &textBytes) || !auxiliaryText(usage.LocalizedName, &textBytes)) && false', 'TestQuotaRuntimeAdmissionCorrectionsAndBounds', FILE),
    ('missing-health', 'warn("quota_evidence_missing")', '_ = i', 'TestQuotaRuntimeHealthAndOwnership', FILE),
    ('health-ownership', 'table.Health.Warnings = slices.Clone(result.Health.Warnings)', 'table.Health.Warnings = result.Health.Warnings', 'TestQuotaRuntimeHealthAndOwnership', FILE),
    ('terminal-cancellation', 'if err := ctx.Err(); err != nil {\n\t\treturn nil, err\n\t}\n\treturn result, nil', 'if err := ctx.Err(); err != nil && false {\n\t\treturn nil, err\n\t}\n\treturn result, nil', 'TestQuotaRuntimeScopeRequestsAndCancellation', FILE),
    ('display-identity', 'key := [4]string{strings.ToLower(row.SubscriptionID), row.Region, row.QuotaType, row.ResourceName}', 'key := [4]string{strings.ToLower(row.SubscriptionID), row.Region, row.QuotaType, row.DisplayName}', 'TestQuotaRuntimeArithmeticAndDisplay', 'internal/plugins/region/auxiliary.go'),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-quota-runtime-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, file in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Quota baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / file
        original = path.read_text(encoding='utf-8')
        if original.count(before) != 1:
            raise SystemExit('Quota mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        output = outcome.stdout + outcome.stderr
        if outcome.returncode == 0 or '--- FAIL: ' + test not in output or 'panic:' in output or '[build failed]' in output:
            raise SystemExit('Quota mutation survived or failed outside named assertions: ' + name + '\n' + output)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored quota baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Region quota runtime mutation rejected:', name)
