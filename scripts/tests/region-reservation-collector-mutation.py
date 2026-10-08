#!/usr/bin/env python3
"""Require compiling reservation collection faults and restored named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
MUTATIONS = (('version',
  'u.RawQuery = "api-version=2024-11-01"',
  'u.RawQuery = "api-version=2023-01-01"',
  'TestReservationCollectorLiteral',
  'internal/plugins/region/reservation_collector.go'),
 ('expand',
  'u.RawQuery = "%24expand=instanceView&api-version=2024-11-01"',
  'u.RawQuery = "api-version=2024-11-01"',
  'TestReservationCollectorLiteral',
  'internal/plugins/region/reservation_collector.go'),
 ('scope',
  'scope[request.SubscriptionID] == ""',
  '(scope[request.SubscriptionID] == "" && false)',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/reservation_collector.go'),
 ('group-scope',
  'ok && sub == selected',
  'ok && (sub == selected || true)',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/reservation_collector.go'),
 ('summary-identity',
  '(r.ID != "" && !strings.EqualFold(r.ID, id))',
  '(r.ID != "" && !strings.EqualFold(r.ID, id) && false)',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/reservation_collector.go'),
 ('get-identity',
  '(raw.ID != "" && !strings.EqualFold(raw.ID, id))',
  '(raw.ID != "" && !strings.EqualFold(raw.ID, id) && false)',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/reservation_collector.go'),
 ('get-region',
  '(raw.Location != "" && strings.ToLower(raw.Location) != region)',
  '(raw.Location != "" && strings.ToLower(raw.Location) != region && false)',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/reservation_collector.go'),
 ('allocated-presence',
  'r.AllocatedKnown = true',
  'r.AllocatedKnown = false',
  'TestReservationCollectorLiteral',
  'internal/plugins/region/reservation_collector.go'),
 ('allocated-type',
  '!reservationVM(ref.ID, subscription)',
  '(!reservationVM(ref.ID, subscription) && false)',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/reservation_collector.go'),
 ('partial',
  'if accepted {',
  'if accepted && false {',
  'TestReservationCollectorFailurePagination',
  'internal/plugins/region/reservation_collector.go'),
 ('calls',
  'calls >= MaxReservationCalls',
  'calls > MaxReservationCalls',
  'TestReservationCollectorBudgets',
  'internal/plugins/region/reservation_collector.go'),
 ('pages',
  'pages >= MaxQuotaPages',
  'pages > MaxQuotaPages',
  'TestReservationCollectorBudgets',
  'internal/plugins/region/reservation_collector.go'),
 ('page-bytes',
  'len(b) > MaxQuotaPageBytes',
  '(len(b) > MaxQuotaPageBytes && false)',
  'TestReservationCollectorBudgets',
  'internal/plugins/region/reservation_collector.go'),
 ('total-bytes',
  'len(b) > MaxQuotaTotalBytes-total',
  '(len(b) > MaxQuotaTotalBytes-total && false)',
  'TestReservationCollectorBudgets',
  'internal/plugins/region/reservation_collector.go'),
 ('work',
  'n > MaxAuxRows-work',
  '(n > MaxAuxRows-work && false)',
  'TestReservationCollectorBudgets',
  'internal/plugins/region/reservation_collector.go'),
 ('id-cap',
  'stringLimit = MaxReservationIDBytes',
  'stringLimit = 8192',
  'TestReservationCollectorIdentity',
  'internal/plugins/region/quota_collector.go'),
 ('terminal-cancellation',
  'if e := ctx.Err(); e != nil {\n\t\treturn nil, e\n\t}\n\treturn result, nil',
  'if e := ctx.Err(); e != nil && false {\n\t\treturn nil, e\n\t}\n\treturn result, nil',
  'TestReservationCollectorCancellationOwnership',
  'internal/plugins/region/reservation_collector.go'))


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^' + test + '$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-reservation-runtime-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, file in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Reservation baseline failed: ' + name + '\n' + baseline.stdout + baseline.stderr)
        path = directory / file
        original = path.read_text(encoding='utf-8')
        if original.count(before) != 1:
            raise SystemExit('Reservation mutation anchor changed: ' + name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        output = outcome.stdout + outcome.stderr
        if outcome.returncode == 0 or '--- FAIL: ' + test not in output or 'panic:' in output or '[build failed]' in output:
            raise SystemExit('Reservation mutation survived or failed outside named assertions: ' + name + '\n' + output)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored reservation baseline failed: ' + name + '\n' + restored.stdout + restored.stderr)
        print('Reservation collector mutation rejected:', name)
