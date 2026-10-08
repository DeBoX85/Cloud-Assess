#!/usr/bin/env python3
"""Require compiling REST quota collection faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/quota_collector.go'
MUTATIONS = (
    ('vm-version', 'return "Microsoft.Compute", "2024-11-01"', 'return "Microsoft.Compute", "2023-01-01"', 'TestVMQuotaCollectorRetainedEvidence', FILE),
    ('vm-unknown-status', 'request.QuotaType != "VM" && (status == 404 || status == 405)', 'status == 404 || status == 405', 'TestVMQuotaCollectorRetainedEvidence', FILE),
    ('origin', 'len(endpoint) > 2048 || u.Scheme != "https"', 'len(endpoint) > 2048 || (u.Scheme != "https" && false)', 'TestRESTQuotaCollectorRequestAdmission', 'internal/plugins/region/quota_collector.go'),
    ('scope', 'scope[request.SubscriptionID] == ""', '(scope[request.SubscriptionID] == "" && false)', 'TestRESTQuotaCollectorRequestAdmission', 'internal/plugins/region/quota_collector.go'),
    ('version', 'return "Microsoft.Sql", "2021-11-01"', 'return "Microsoft.Sql", "2022-07-01"', 'TestRESTQuotaCollectorLiteralProviders', 'internal/plugins/region/quota_collector.go'),
    ('continuation-origin', '!strings.EqualFold(u.Host, c.origin.Host)', '(!strings.EqualFold(u.Host, c.origin.Host) && false)', 'TestRESTQuotaCollectorContinuation', 'internal/plugins/region/quota_collector.go'),
    ('continuation-path', 'u.Path != first.Path', '(u.Path != first.Path && false)', 'TestRESTQuotaCollectorContinuation', 'internal/plugins/region/quota_collector.go'),
    ('continuation-version', 'q.Get("api-version") != version', '(q.Get("api-version") != version && false)', 'TestRESTQuotaCollectorContinuation', 'internal/plugins/region/quota_collector.go'),
    ('canonical-authority', 'u.Host = c.origin.Host', '_ = c.origin.Host', 'TestRESTQuotaCollectorContinuation', 'internal/plugins/region/quota_collector.go'),
    ('noncharacters', 'return r >= 0xfdd0 && r <= 0xfdef || r&0xffff == 0xfffe || r&0xffff == 0xffff', 'return r==0xfffe || r==0xffff', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/primary.go'),
    ('cycle', 'if seenLinks[u.String()] {', 'if seenLinks[u.String()] && false {', 'TestRESTQuotaCollectorContinuation', 'internal/plugins/region/quota_collector.go'),
    ('status', 'response.StatusCode != http.StatusOK', '(response.StatusCode != http.StatusOK && false)', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/quota_collector.go'),
    ('unsupported', 'status == 404 || status == 405', 'status == 404 || status == 405 || status == 400', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/quota_collector.go'),
    ('partial', 'if pages > 0 {', 'if pages > 0 && false {', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/quota_collector.go'),
    ('pages', 'pages >= MaxQuotaPages', 'pages > MaxQuotaPages', 'TestRESTQuotaCollectorBudgets', 'internal/plugins/region/quota_collector.go'),
    ('page-bytes', 'len(body) > MaxQuotaPageBytes', '(len(body) > MaxQuotaPageBytes && false)', 'TestRESTQuotaCollectorBudgets', 'internal/plugins/region/quota_collector.go'),
    ('aggregate-bytes', 'len(body) > MaxQuotaTotalBytes-total', '(len(body) > MaxQuotaTotalBytes-total && false)', 'TestRESTQuotaCollectorBudgets', 'internal/plugins/region/quota_collector.go'),
    ('aggregate-rows', 'len(rows) > MaxAuxRows-len(result.Evidence.Usages)', '(len(rows) > MaxAuxRows-len(result.Evidence.Usages) && false)', 'TestRESTQuotaCollectorBudgets', 'internal/plugins/region/quota_collector.go'),
    ('duplicate-keys', 'keys[strings.ToLower(s)] {', '(keys[strings.ToLower(s)] && false) {', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/quota_collector.go'),
    ('tokens', 'tokens >= 65536', 'tokens > 65536', 'TestRESTQuotaCollectorBudgets', 'internal/plugins/region/quota_collector.go'),
    ('depth', 'depth > 32', 'depth > 33', 'TestRESTQuotaCollectorBudgets', 'internal/plugins/region/quota_collector.go'),
    ('count', 'row.Current < 0', '(row.Current < 0 && false)', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/quota_collector.go'),
    ('presence', 'row.CurrentKnown = true', 'row.CurrentKnown = false', 'TestRESTQuotaCollectorLiteralProviders', 'internal/plugins/region/quota_collector.go'),
    ('crosspage-identity', 'seenNames[row.ResourceName] || pageNames[row.ResourceName]', '(seenNames[row.ResourceName] && false) || pageNames[row.ResourceName]', 'TestRESTQuotaCollectorFailureHealth', 'internal/plugins/region/quota_collector.go'),
    ('terminal-cancellation', 'if err := ctx.Err(); err != nil {\n\t\treturn nil, err\n\t}\n\treturn result, nil', 'if err := ctx.Err(); err != nil && false {\n\t\treturn nil, err\n\t}\n\treturn result, nil', 'TestRESTQuotaCollectorIsolationCancellation', 'internal/plugins/region/quota_collector.go'),
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
        print('REST quota collector mutation rejected:', name)
