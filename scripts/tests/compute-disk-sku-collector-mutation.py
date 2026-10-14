#!/usr/bin/env python3
"""Compiling disk SKU production controls with restored named independent cases."""
from pathlib import Path
import os
import subprocess

ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / 'internal/plugins/region/compute_disk_sku_collector.go'
GO = os.environ.get('QA_GO', 'go')


def run(select):
    return subprocess.run([GO, 'test', './internal/plugins/region', '-run', select, '-count=1', '-v'],
                          cwd=ROOT, capture_output=True, text=True, timeout=180)


def check(select):
    result = run(select)
    if result.returncode:
        raise SystemExit('Disk SKU restored baseline failed\n' + result.stdout + result.stderr)


original = SOURCE.read_bytes()
text = original.decode('utf-8')
controls = (
    ("request-version", "computeDiskSKUVersion = \"2021-07-01\"", "computeDiskSKUVersion = \"2023-03-01\"", "TestComputeDiskSKULiteralStates"),
    ("continuation-origin", "!strings.EqualFold(target.Host, c.transport.origin.Host)", "false", "TestComputeDiskSKURequestsAndContinuation"),
    ("continuation-path", "target.Path != first.Path", "false", "TestComputeDiskSKURequestsAndContinuation"),
    ("continuation-version", "query.Get(\"api-version\") != computeDiskSKUVersion", "false", "TestComputeDiskSKURequestsAndContinuation"),
    ("json-admission", "if !quotaJSON(ctx, body) {", "if false {", "TestComputeDiskSKURejectsMalformedPages"),
    ("type-filter", "if strings.ToLower(kind) != \"disks\" {", "if false {", "TestComputeDiskSKULiteralStates"),
    ("restriction-values", "range append(values, scope...)", "range scope", "TestComputeDiskSKULiteralStates"),
    ("restriction-precedence", "} else if restricted {", "} else if false {", "TestComputeDiskSKULiteralStates"),
    ("zone-order", "sort.Strings(value.BlockedZones)", "sort.Sort(sort.Reverse(sort.StringSlice(value.BlockedZones)))", "TestComputeDiskSKULiteralStates"),
    ("prior-collision", "if _, exists := prior[key]; exists {", "if false {", "TestComputeDiskSKUAtomicFailureAndSnapshot"),
    ("atomic-page", "page.values[key] = value", "page.values[key] = value\n\t\tprior[key] = value", "TestComputeDiskSKUAtomicFailureAndSnapshot"),
    ("partial-snapshot", "if c == nil || c.Status != \"complete\" {", "if c == nil {", "TestComputeDiskSKUAtomicFailureAndSnapshot"),
    ("snapshot-slice", "value.BlockedZones = append([]string{}, value.BlockedZones...)", "// omitted slice copy", "TestComputeDiskSKUAtomicFailureAndSnapshot"),
    ("page-bytes", "len(body) > MaxQuotaPageBytes", "false", "TestComputeDiskSKUBounds"),
    ("aggregate-bytes", "len(body) > MaxQuotaTotalBytes-total", "false", "TestComputeDiskSKUBounds"),
    ("page-count", "admitted >= MaxQuotaPages", "admitted > MaxQuotaPages", "TestComputeDiskSKUBounds"),
    ("raw-work", "if b.work >= MaxComputeDiskSKUWork {", "if b.work > MaxComputeDiskSKUWork {", "TestComputeDiskSKUBounds"),
    ("pre-auth-scope", "scope[id] == \"\"", "false", "TestComputeDiskSKURequestAndCancellation"),
    ("optional-metadata", "!diskMetadata(ctx, row)", "false", "TestComputeDiskSKURejectsMalformedPages"),
    ("cancellation-no-result", "return nil, cause", "return result, nil", "TestComputeDiskSKURequestAndCancellation"),
    ("getter-limit", "target.String(), MaxQuotaPageBytes)", "target.String(), 0)", "TestComputeDiskSKULiteralStates"),
)
check('^TestComputeDiskSKU')
try:
    for name, old, replacement, expected in controls:
        if text.count(old) != 1:
            raise SystemExit('Disk SKU production anchor changed: ' + name)
        SOURCE.write_bytes(text.replace(old, replacement, 1).encode('utf-8'))
        compiled = run('^$')
        if compiled.returncode:
            raise SystemExit('Disk SKU mutation did not compile: ' + name + '\n' + compiled.stdout + compiled.stderr)
        failed = run('^' + expected + '$')
        if failed.returncode == 0 or ('--- FAIL: ' + expected) not in failed.stdout or '[build failed]' in failed.stdout + failed.stderr:
            raise SystemExit('Disk SKU mutation missed named assertion: ' + name + '\n' + failed.stdout + failed.stderr)
        SOURCE.write_bytes(original)
        check('^' + expected + '$')
        print('Compute disk SKU production compiling mutation rejected: ' + name + '; restored named baseline passed', flush=True)
finally:
    SOURCE.write_bytes(original)
check('^TestComputeDiskSKU')
if SOURCE.read_bytes() != original:
    raise SystemExit('Disk SKU production source was not restored')
