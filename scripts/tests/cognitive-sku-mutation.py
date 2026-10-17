#!/usr/bin/env python3
"""Compile bounded Cognitive reversals; require named failure and exact restoration."""
from pathlib import Path
import os
import subprocess
ROOT = Path(__file__).resolve().parents[2]
CATALOG = ROOT / 'internal/plugins/region/cognitive_sku_catalog.go'
COLLECTOR = ROOT / 'internal/plugins/region/cognitive_sku_collector.go'
BUDGET = ROOT / 'internal/plugins/region/compute_vm_sku_collector.go'
TRANSPORT = ROOT / 'internal/plugins/region/quota_collector.go'
HTTP = ROOT / 'internal/azure/bounded_get.go'
GO = os.environ.get('QA_GO', 'go')
def run(select):
    return subprocess.run([GO, 'test', './internal/plugins/region', '-run', select,
                           '-count=1', '-v'], cwd=ROOT, capture_output=True,
                          text=True, encoding='utf-8', errors='backslashreplace', timeout=180)


def check(select):
    result = run(select)
    if result.returncode:
        raise SystemExit('Cognitive restored baseline failed\n' + result.stdout + result.stderr)


manifest = 'TestCognitiveSKUFrozenManifest'
originals = {p: p.read_bytes() for p in (CATALOG, COLLECTOR, BUDGET, TRANSPORT, HTTP)}
texts = {p: b.decode('utf-8').replace('\r\n', '\n') for p, b in originals.items()}
controls = (
    ('http-body-closure', HTTP, 'func (b *boundedResponseBody) Read(p []byte) (int, error) {', 'func (b *boundedResponseBody) Close() error { return nil }; func (b *boundedResponseBody) Read(p []byte) (int, error) {', 'TestCognitiveSKUActualHTTPClosure'),
    ('request-kind', COLLECTOR, '&& storageIdentifier(q.Kind)', '&& true', manifest),
    ('request-scope', COLLECTOR, 'scope[id] != ""', 'true', manifest),
    ('alias-provenance', COLLECTOR, 'strings.ReplaceAll(key, " ", "") != value', 'false', 'TestCognitiveSKUAliasProvenance'),
    ('kind-filter', COLLECTOR, '!strings.EqualFold(kind, q.Kind)', 'false', manifest),
    ('resource-filter', COLLECTOR, '!strings.EqualFold(resourceType, "accounts")', 'false', manifest),
    ('known-before-filter', COLLECTOR, 'f, ok := cognitiveRow(ctx, row, q.Region, q.LocationAliases)', 'if !strings.EqualFold(kind, q.Kind) { continue }; f, ok := cognitiveRow(ctx, row, q.Region, q.LocationAliases)', manifest),
    ('overlapping-scope', COLLECTOR, 'if old.scopes[loc] {', 'if old.scopes[loc] && false {', manifest),
    ('disjoint-reduction', COLLECTOR, 'old.listed = old.listed || f.listed', 'old.listed = f.listed', manifest),
    ('known-uncertainty', COLLECTOR, 'Known: !f.unknown', 'Known: true', manifest),
    ('regional-restriction', COLLECTOR, 'g.State = SKURestricted', 'g.State = SKUAvailable', manifest),
    ('zone-restriction', COLLECTOR, 'g.State = SKUZoneRestricted', 'g.State = SKUAvailable', manifest),
    ('atomic-scope-state', COLLECTOR, 'copyFlags.scopes = map[string]bool{}', 'copyFlags.scopes = f.scopes', 'TestCognitiveSKUAtomicRegionalState'),
    ('null-next', COLLECTOR, 'has && !diskScalar(raw, &page.next)', 'has && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) && !diskScalar(raw, &page.next)', manifest),
    ('status-precedence', COLLECTOR, 'response != nil && response.StatusCode != http.StatusOK', 'response != nil && response.StatusCode != http.StatusOK && e == nil', manifest),
    ('wrapped-cancel', COLLECTOR, 'errors.Is(e, context.Canceled)', 'e == context.Canceled', manifest),
    ('wrapped-deadline', COLLECTOR, 'errors.Is(e, context.DeadlineExceeded)', 'e == context.DeadlineExceeded', manifest),
    ('page-byte-bound', COLLECTOR, 'len(body) > MaxQuotaPageBytes', 'len(body) > MaxQuotaPageBytes+1', 'TestCognitiveSKUByteAndPageBounds'),
    ('aggregate-byte-bound', COLLECTOR, 'len(body) > MaxQuotaTotalBytes-total', 'len(body) > MaxQuotaTotalBytes-total+1', 'TestCognitiveSKUByteAndPageBounds'),
    ('page-count-bound', COLLECTOR, 'pages >= MaxQuotaPages', 'pages > MaxQuotaPages', 'TestCognitiveSKUByteAndPageBounds'),
    ('raw-row-bound', COLLECTOR, 'page.rows > MaxCognitiveSKURows', 'page.rows > MaxCognitiveSKURows+1', 'TestCognitiveSKUFilteredRawRowBound'),
    ('collector-final-context', COLLECTOR, 'if e := ctx.Err(); e != nil {\n\t\treturn nil, e\n\t}\n\treturn result, nil', 'return result, nil', 'TestCognitiveSKUFinalReturnContext'),
    ('selected-empty-label', CATALOG, 'if label == "" {', 'if false {', 'TestCognitiveCatalogScopeIsolation'),
    ('collection-binding', CATALOG, 'c.SubscriptionID != id || c.Region != region || c.Kind != kind || c.ResourceType != cognitiveSKUType', 'false', manifest),
    ('failure-enum', CATALOG, '!cognitiveSKUFailure(c.FailureCode)', 'false', manifest),
    ('unrequired-validation', CATALOG, 'for name, g := range c.Groups {', 'for name, g := range c.Groups {\n\t\t\tif len(required) == 0 { continue }', 'TestCognitiveCatalogValidationAndSanitization'),
    ('canonical-zones', CATALOG, 'zone <= prior', 'zone == "" && prior != ""', manifest),
    ('unknown-required', CATALOG, 'found && !g.Known', 'found && !g.Known && false', manifest),
    ('partial-isolation', CATALOG, 'c == nil || c.Status != "complete"', 'c == nil || c.Status != "complete" && c.Status != "partial"', manifest),
    ('missing-unavailable', CATALOG, 'g = CognitiveSKUGroup{Known: true, State: SKUUnavailable}', 'g = CognitiveSKUGroup{Known: true, State: SKUAvailable}', manifest),
    ('owned-zones', CATALOG, 'zones = append(zones, zone)', 'zones = g.BlockedZones', 'TestCognitiveCatalogOwnedConcurrentOutput'),
    ('eligibility-limitation', CATALOG, 'DeploymentEligibility: "not_assessed"', 'DeploymentEligibility: "eligible"', manifest),
    ('fixed-warning', CATALOG, 'Message: cognitiveCatalogWarning', 'Message: ""', manifest),
    ('sanitized-error', CATALOG, 'fmt.Errorf("Cognitive SKU catalog projection rejected [%s]", code)', 'fmt.Errorf("Cognitive SKU catalog projection rejected [%s]", code+expected.SubscriptionID)', 'TestCognitiveCatalogValidationAndSanitization'),
    ('scope-count', CATALOG, 'len(selected) > MaxSubscriptions', 'len(selected) > MaxSubscriptions+1', 'TestCognitiveCatalogCountBoundaries'),
    ('alias-count', CATALOG, 'len(expected.LocationAliases) > MaxCognitiveSKUAliases', 'len(expected.LocationAliases) > MaxCognitiveSKUAliases+1', 'TestCognitiveSKUScopeAliasPreauth'),
    ('required-count', CATALOG, 'len(required) > MaxCognitiveSKUGroups', 'len(required) > MaxCognitiveSKUGroups+1', 'TestCognitiveCatalogCountBoundaries'),
    ('group-count', CATALOG, 'len(c.Groups) > MaxCognitiveSKUGroups', 'len(c.Groups) > MaxCognitiveSKUGroups+1', 'TestCognitiveCatalogCountBoundaries'),
    ('logical-work', BUDGET, 'b.work >= MaxComputeVMSKUWork', 'b.work > MaxComputeVMSKUWork', 'TestCognitiveCatalogWorkBoundary'),
    ('raw-logical-work', BUDGET, 'b.work >= MaxComputeVMSKUWork', 'b.work > MaxComputeVMSKUWork', 'TestCognitiveSKURowAndWorkBounds'),
    ('input-text', BUDGET, 'len(s) > auxTextBudget-b.text', 'len(s) > auxTextBudget', 'TestCognitiveCatalogTextBoundaries'),
    ('label-bound', BUDGET, 'len(s) > limit ||', 'false ||', 'TestCognitiveCatalogLabelBoundaries'),
    ('output-text', CATALOG, 'if !output.label(zone, MaxLabelBytes) {', 'if false {', 'TestCognitiveCatalogTextBoundaries'),
    ('catalog-final-context', CATALOG, 'finish := func(out CognitiveCatalogEvidence) (CognitiveCatalogEvidence, error) {\n\t\tif err := ctx.Err(); err != nil {\n\t\t\treturn CognitiveCatalogEvidence{}, err\n\t\t}', 'finish := func(out CognitiveCatalogEvidence) (CognitiveCatalogEvidence, error) {', 'TestCognitiveCatalogFinalReturnContext'),
    ('no-scoring-conversion', COLLECTOR, 'func cognitiveSKUFailure(code string) bool {', 'func(c *CognitiveSKUCollection) Evidence() SKUEvidence { return SKUEvidence{} }; func cognitiveSKUFailure(code string) bool {', manifest),
    ('transport-raw-URL', TRANSPORT, 'len(link) > 8192', 'len(link) > 8193', 'TestCognitiveSKUTransportURLBound'),
)
check('^TestCognitive')
try:
    for name, source, old, replacement, expected in controls:
        text = texts[source]
        if text.count(old) != 1:
            raise SystemExit('Cognitive production anchor changed: ' + name)
        source.write_bytes(text.replace(old, replacement, 1).encode('utf-8'))
        compiled = run('^$')
        if compiled.returncode:
            raise SystemExit('Cognitive reversal did not compile: ' + name + '\n' + compiled.stdout + compiled.stderr)
        failed = run('^' + expected + '$')
        if failed.returncode == 0 or ('--- FAIL: ' + expected) not in failed.stdout or '[build failed]' in failed.stdout + failed.stderr:
            raise SystemExit('Cognitive reversal missed named assertion: ' + name + '\n' + failed.stdout + failed.stderr)
        source.write_bytes(originals[source])
        check('^' + expected + '$')
        print('Cognitive catalog compiling mutation rejected: ' + name + '; restored named baseline passed', flush=True)
finally:
    for source, original in originals.items():
        source.write_bytes(original)
check('^TestCognitive')
if any(p.read_bytes() != b for p, b in originals.items()):
    raise SystemExit('Cognitive source was not restored exactly')
