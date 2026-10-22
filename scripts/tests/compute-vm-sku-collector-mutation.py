#!/usr/bin/env python3
"""Compile isolated VM collector reversals, require named failures and restore bytes."""
from pathlib import Path
import os
import subprocess

ROOT = Path(__file__).resolve().parents[2]
SOURCE = ROOT / 'internal/plugins/region/compute_vm_sku_collector.go'
GO = os.environ.get('QA_GO', 'go')


def run(select):
    return subprocess.run([GO, 'test', './internal/plugins/region', '-run', select, '-count=1', '-v'],
                          cwd=ROOT, capture_output=True, text=True, encoding='utf-8',
                          errors='backslashreplace', timeout=180)


def check(select):
    result = run(select)
    if result.returncode:
        raise SystemExit('VM restored baseline failed\n' + result.stdout + result.stderr)


original = SOURCE.read_bytes()
# Match logical lines on both native checkouts; restore exact original bytes.
text = original.decode('utf-8').replace('\r\n', '\n')
manifest = 'TestComputeVMSKUFrozenManifest'
controls = (
    ('request-version', 'const vmSKUVersion = "2021-07-01"', 'const vmSKUVersion = "2023-03-01"', manifest),
    ('getter-bound', 'target.String(), MaxQuotaPageBytes)', 'target.String(), 0)', manifest),
    ('pre-auth-scope', 'scope[id] == ""', 'false', manifest),
    ('required-filter', 'len(q["$filter"]) != 1 || q.Get("$filter") != filter || ', '', manifest),
    ('json-admission', 'if !quotaJSON(ctx, body) {', 'if false {', manifest),
    ('type-filter', 'if !strings.EqualFold(kind, "virtualMachines") {', 'if false {', manifest),
    ('optional-metadata', '!vmSKUMetadata(ctx, row)', 'false', manifest),
    ('nonnegative-capacity', 'diskScalar(raw, &n) && n >= 0', 'diskScalar(raw, &n)', 'TestComputeVMSKUPageCollisionAndMetadata'),
    ('missing-restrictions', '} else {\n\t\tout.Known = false\n\t}\n\tif raw, has := row["capabilities"]; has {', '} else {\n\t\tout.Known = true\n\t}\n\tif raw, has := row["capabilities"]; has {', manifest),
    ('location-membership', 'if !listed || hard {', 'if hard {', manifest),
    ('restriction-union', 'range append(values, locations...)', 'range values', manifest),
    ('subscription-reason', 'reason != "QuotaId" && reason != "NotAvailableForSubscription"', 'reason != "QuotaId"', manifest),
    ('region-zone-precedence', 'if restricted {\n\t\tout.State = SKURestricted', 'if restricted && false {\n\t\tout.State = SKURestricted', manifest),
    ('capability-conflict', 'prior != "" && prior != value', 'false && prior != value', manifest),
    ('page-collision', 'if _, exists := page.groups[name]; exists {', 'if _, exists := page.groups[name]; exists && false {', manifest),
    ('cross-page-collision', 'if _, exists := prior[name]; exists {', 'if _, exists := prior[name]; exists && false {', 'TestComputeVMSKUPageCollisionAndMetadata'),
    ('group-count', 'len(prior)+len(page.groups) >= MaxComputeVMSKUGroups', 'len(prior)+len(page.groups) > MaxComputeVMSKUGroups', 'TestComputeVMSKUGroupsAndScopeBounds'),
    ('atomic-page', 'if code != "" {\n\t\t\treturn fail(code)\n\t\t}', 'if code != "" && false {\n\t\t\treturn fail(code)\n\t\t}', manifest),
    ('partial-projection', 'if c == nil || c.Status != "complete" {', 'if c == nil || c.Status != "complete" && c.Status != "partial" {', manifest),
    ('unknown-required', 'if found && !g.Known {', 'if found && !g.Known && false {', manifest),
    ('unrequired-isolation', 'if found && !g.Known {', 'if len(c.Groups) > len(names) || found && !g.Known {', manifest),
    ('vmss-unsupported', 'if c.ResourceType == vmSSSKUType {', 'if c.ResourceType == vmSSSKUType && false {', manifest),
    ('owned-zone-projection', 'zones = append(zones, z)', 'zones = g.BlockedZones', 'TestComputeVMSKUOwnershipAndLateContext'),
    ('status-error-pair', 'if response != nil && response.StatusCode != http.StatusOK {', 'if e == nil && response != nil && response.StatusCode != http.StatusOK {', 'TestComputeVMSKUActualHTTPAndCancellation'),
    ('context-priority', 'return nil, cause', 'return result, nil', 'TestComputeVMSKUActualHTTPAndCancellation'),
    ('body-bytes', 'if len(body) > MaxQuotaPageBytes {', 'if false {', 'TestComputeVMSKUPageAndAggregateBounds'),
    ('aggregate-bytes', 'if len(body) > MaxQuotaTotalBytes-total {', 'if false {', 'TestComputeVMSKUPageAndAggregateBounds'),
    ('page-count', 'if pages >= MaxQuotaPages {', 'if pages > MaxQuotaPages {', 'TestComputeVMSKUPageAndAggregateBounds'),
    ('raw-work', 'b.work >= MaxComputeVMSKUWork', 'b.work > MaxComputeVMSKUWork', 'TestComputeVMSKUJSONAndRawWorkBounds'),
    ('projection-work', 'b.work >= MaxComputeVMSKUWork', 'b.work > MaxComputeVMSKUWork', 'TestComputeVMSKUProjectionBoundsAndCorruption'),
    ('projection-text', 'len(s) > auxTextBudget-b.text', 'len(s) > auxTextBudget', 'TestComputeVMSKUProjectionBoundsAndCorruption'),
    ('late-unknown-unsupported', 'empty := func(status string) (SKUEvidence, error) {\n\t\tif e := ctx.Err(); e != nil {\n\t\t\treturn SKUEvidence{}, e\n\t\t}', 'empty := func(status string) (SKUEvidence, error) {', 'TestComputeVMSKUOwnershipAndLateContext'),
)
check('^TestComputeVMSKU')
try:
    for name, old, replacement, expected in controls:
        if text.count(old) != 1:
            raise SystemExit('VM production anchor changed: ' + name)
        SOURCE.write_bytes(text.replace(old, replacement, 1).encode('utf-8'))
        compiled = run('^$')
        if compiled.returncode:
            raise SystemExit('VM mutation did not compile: ' + name + '\n' + compiled.stdout + compiled.stderr)
        failed = run('^' + expected + '$')
        if failed.returncode == 0 or ('--- FAIL: ' + expected) not in failed.stdout or '[build failed]' in failed.stdout + failed.stderr:
            raise SystemExit('VM mutation missed named assertion: ' + name + '\n' + failed.stdout + failed.stderr)
        SOURCE.write_bytes(original)
        check('^' + expected + '$')
        print('VM SKU compiling mutation rejected: ' + name + '; restored named baseline passed', flush=True)
finally:
    SOURCE.write_bytes(original)
check('^TestComputeVMSKU')
if SOURCE.read_bytes() != original:
    raise SystemExit('VM production source was not restored')
