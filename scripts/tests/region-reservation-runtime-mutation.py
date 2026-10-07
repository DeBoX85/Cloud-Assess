#!/usr/bin/env python3
"""Require compiling reservation faults and restored independent named assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = 'internal/plugins/region/reservation_runtime.go'
MUTATIONS = (
    ('status-precedence', 'case usage.Allocated == 0:', 'case available == 0:', 'TestReservationRuntimeLiteralStatusAndCells', 1),
    ('structural-provider', 'if len(parts[index]) != len(fixed) || !strings.EqualFold(parts[index], fixed) {', 'if (len(parts[index]) != len(fixed) || !strings.EqualFold(parts[index], fixed)) && false {', 'TestReservationRuntimeIdentityAdmission', 1),
    ('selected-scope', 'return request, selected', 'return request, selected || true', 'TestReservationRuntimeRequestsAndCancellation', 1),
    ('record-subscription', 'sub != key.SubscriptionID', '(sub != key.SubscriptionID && false)', 'TestReservationRuntimeIdentityAdmission', 1),
    ('record-region', 'usage.Region != key.Region', '(usage.Region != key.Region && false)', 'TestReservationRuntimeIdentityAdmission', 1),
    ('response-name', '!strings.EqualFold(usage.ResponseName, name)', '(!strings.EqualFold(usage.ResponseName, name) && false)', 'TestReservationRuntimeIdentityAdmission', 1),
    ('response-region', 'usage.ResponseRegion != usage.Region', '(usage.ResponseRegion != usage.Region && false)', 'TestReservationRuntimeIdentityAdmission', 1),
    ('duplicate-identity', 'if seen[canonical] {', 'if seen[canonical] && false {', 'TestReservationRuntimeIdentityAdmission', 1),
    ('negative-count', 'usage.Reserved < 0 ||', '(usage.Reserved < 0 && false) ||', 'TestReservationRuntimeIdentityAdmission', 1),
    ('presence-contradiction', '!usage.ReservedKnown && usage.Reserved != 0', '(!usage.ReservedKnown && usage.Reserved != 0 && false)', 'TestReservationRuntimeIdentityAdmission', 1),
    ('unknown-idle', 'if !usage.ReservedKnown || !usage.AllocatedKnown {', 'if (!usage.ReservedKnown || !usage.AllocatedKnown) && false {', 'TestReservationRuntimeHealthAndOwnership', 1),
    ('skipped-work', 'len(response.Reservations) > MaxAuxRows-count', '(len(response.Reservations) > MaxAuxRows-count && false)', 'TestReservationRuntimeWorkAndDeclaredOrder', 1),
    ('raw-id-accounting', 'textBytes += len(usage.ResourceID)', 'textBytes += 0', 'TestReservationRuntimeRawAndOutputTextBudgets', 1),
    ('skipped-label', 'if !auxiliaryText(label, &textBytes) {', 'if !auxiliaryText(label, &textBytes) && false {', 'TestReservationRuntimeIdentityAdmission', 2),
    ('output-accounting', 'outputText += 128', 'outputText += 0', 'TestReservationRuntimeRawAndOutputTextBudgets', 1),
    ('missing-health', 'warn("reservation_evidence_missing")', '_ = i', 'TestReservationRuntimeHealthAndOwnership', 1),
    ('warning-ownership', 'table.Health.Warnings = slices.Clone(result.Health.Warnings)', '_ = slices.Clone(result.Health.Warnings); table.Health.Warnings = result.Health.Warnings', 'TestReservationRuntimeHealthAndOwnership', 1),
    ('terminal-cancellation', 'if err := ctx.Err(); err != nil {\n\t\treturn nil, err\n\t}\n\treturn result, nil', 'if err := ctx.Err(); err != nil && false {\n\t\treturn nil, err\n\t}\n\treturn result, nil', 'TestReservationRuntimeRequestsAndCancellation', 1),
)


def run(directory, test):
    return subprocess.run(['go', 'test', '-count=1', '-run', '^'+test+'$',
                           './internal/plugins/region'], cwd=directory,
                          capture_output=True, text=True, timeout=120)


with tempfile.TemporaryDirectory(prefix='cloud-assess-reservation-runtime-mutations-') as temporary:
    directory = Path(temporary)
    for name in ('go.mod', 'go.sum'):
        shutil.copyfile(ROOT/name, directory/name)
    for name in ('cmd', 'internal'):
        shutil.copytree(ROOT/name, directory/name, ignore=shutil.ignore_patterns('.git'))
    for name, before, after, test, count in MUTATIONS:
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit('Reservation baseline failed: '+name+'\n'+baseline.stdout+baseline.stderr)
        path = directory/FILE
        original = path.read_text(encoding='utf-8')
        if original.count(before) != count:
            raise SystemExit('Reservation mutation anchor changed: '+name)
        try:
            path.write_text(original.replace(before, after), encoding='utf-8')
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding='utf-8')
        output = outcome.stdout+outcome.stderr
        if outcome.returncode == 0 or '--- FAIL: '+test not in output or 'panic:' in output or '[build failed]' in output:
            raise SystemExit('Reservation fault survived or failed outside assertions: '+name+'\n'+output)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit('Restored reservation baseline failed: '+name+'\n'+restored.stdout+restored.stderr)
        print('Region reservation runtime mutation rejected:', name)
