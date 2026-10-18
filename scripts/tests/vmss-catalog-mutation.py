#!/usr/bin/env python3
"""Compile isolated catalog reversals, require named failures, restore exact bytes."""
from pathlib import Path
import os
import subprocess

ROOT = Path(__file__).resolve().parents[2]
CATALOG = ROOT / 'internal/plugins/region/vmss_catalog.go'
LEGACY = ROOT / 'internal/plugins/region/compute_vm_sku_collector.go'
GO = os.environ.get('QA_GO', 'go')


def run(select):
    return subprocess.run([GO, 'test', './internal/plugins/region', '-run', select,
                           '-count=1', '-v'], cwd=ROOT, capture_output=True,
                          text=True, encoding='utf-8', errors='backslashreplace', timeout=180)


def check(select):
    result = run(select)
    if result.returncode:
        raise SystemExit('VMSS restored baseline failed\n' + result.stdout + result.stderr)


manifest = 'TestVMSSCatalogFrozenManifest'
originals = {p: p.read_bytes() for p in (CATALOG, LEGACY)}
texts = {p: b.decode('utf-8').replace('\r\n', '\n') for p, b in originals.items()}
controls = (
    ('expected-target', CATALOG, 'expected.ResourceType != vmSSSKUType', 'false', manifest),
    ('selected-binding', CATALOG, 'scope[id] == ""', 'false', 'TestVMSSCatalogScopeIsolation'),
    ('empty-label', CATALOG, 'if label == "" {', 'if false {', 'TestVMSSCatalogScopeIsolation'),
    ('collection-binding', CATALOG, 'c.SubscriptionID != id || c.Region != region || c.ResourceType != vmSSSKUType', 'false', manifest),
    ('failure-enum', CATALOG, '!vmSKUFailure(c.FailureCode)', 'false', manifest),
    ('unrequired-validation', CATALOG, 'for name, g := range c.Groups {', 'for name, g := range c.Groups {\n\t\t\tif len(required) == 0 { continue }', 'TestVMSSCatalogValidationAndSanitization'),
    ('canonical-zones', CATALOG, 'zone <= prior', 'zone == "" && prior != ""', manifest),
    ('unknown-required', CATALOG, 'found && !g.Known', 'found && !g.Known && false', manifest),
    ('partial-isolation', CATALOG, 'c == nil || c.Status != "complete"', 'c == nil || c.Status != "complete" && c.Status != "partial"', manifest),
    ('restriction-state', CATALOG, 'State: g.State, BlockedZones: zones', 'State: SKUAvailable, BlockedZones: zones', manifest),
    ('missing-unavailable', CATALOG, 'g = ComputeVMSKUGroup{Known: true, State: SKUUnavailable}', 'g = ComputeVMSKUGroup{Known: true, State: SKUAvailable}', manifest),
    ('owned-zones', CATALOG, 'zones = append(zones, zone)', 'zones = g.BlockedZones', 'TestVMSSCatalogOwnedConcurrentOutput'),
    ('eligibility-limitation', CATALOG, 'DeploymentEligibility: "not_assessed"', 'DeploymentEligibility: "eligible"', manifest),
    ('fixed-warning', CATALOG, 'Message: vmssCatalogWarning', 'Message: ""', manifest),
    ('sanitized-error', CATALOG, 'fmt.Errorf("VMSS catalog projection rejected [%s]", code)', 'fmt.Errorf("VMSS catalog projection rejected [%s]", code+expected.SubscriptionID)', 'TestVMSSCatalogValidationAndSanitization'),
    ('count-bound', CATALOG, 'len(selected) > MaxSubscriptions', 'len(selected) > MaxSubscriptions+1', 'TestVMSSCatalogCountBoundaries'),
    ('required-count', CATALOG, 'len(required) > MaxComputeVMSKUGroups', 'len(required) > MaxComputeVMSKUGroups+1', 'TestVMSSCatalogCountBoundaries'),
    ('group-count', CATALOG, 'len(c.Groups) > MaxComputeVMSKUGroups', 'len(c.Groups) > MaxComputeVMSKUGroups+1', 'TestVMSSCatalogCountBoundaries'),
    ('logical-work', LEGACY, 'b.work >= MaxComputeVMSKUWork', 'b.work > MaxComputeVMSKUWork', 'TestVMSSCatalogWorkBoundary'),
    ('input-text', LEGACY, 'len(s) > auxTextBudget-b.text', 'len(s) > auxTextBudget', 'TestVMSSCatalogTextBoundaries'),
    ('label-bound', LEGACY, 'len(s) > limit ||', 'false ||', 'TestVMSSCatalogLabelBoundaries'),
    ('output-text', CATALOG, 'if !output.label(zone, MaxLabelBytes) {', 'if false {', 'TestVMSSCatalogTextBoundaries'),
    ('late-return-context', CATALOG, 'finish := func(out VMSSCatalogEvidence) (VMSSCatalogEvidence, error) {\n\t\tif err := ctx.Err(); err != nil {\n\t\t\treturn VMSSCatalogEvidence{}, err\n\t\t}', 'finish := func(out VMSSCatalogEvidence) (VMSSCatalogEvidence, error) {', 'TestVMSSCatalogFinalReturnContext'),
    ('unsupported-scoring', LEGACY, 'if c.ResourceType == vmSSSKUType {', 'if c.ResourceType == vmSSSKUType && false {', 'TestVMSSCatalogScoringInvariant'),
)
check('^TestVMSSCatalog')
try:
    for name, source, old, replacement, expected in controls:
        text = texts[source]
        if text.count(old) != 1:
            raise SystemExit('VMSS production anchor changed: ' + name)
        source.write_bytes(text.replace(old, replacement, 1).encode('utf-8'))
        compiled = run('^$')
        if compiled.returncode:
            raise SystemExit('VMSS reversal did not compile: ' + name + '\n' + compiled.stdout + compiled.stderr)
        failed = run('^' + expected + '$')
        if failed.returncode == 0 or ('--- FAIL: ' + expected) not in failed.stdout or '[build failed]' in failed.stdout + failed.stderr:
            raise SystemExit('VMSS reversal missed named assertion: ' + name + '\n' + failed.stdout + failed.stderr)
        source.write_bytes(originals[source])
        check('^' + expected + '$')
        print('VMSS catalog compiling mutation rejected: ' + name + '; restored named baseline passed', flush=True)
finally:
    for source, original in originals.items():
        source.write_bytes(original)
check('^TestVMSSCatalog')
if any(p.read_bytes() != b for p, b in originals.items()):
    raise SystemExit('VMSS source was not restored exactly')
