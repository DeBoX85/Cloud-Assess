#!/usr/bin/env python3
"""Require compiling named provider collector faults and restored assertions."""
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT = Path(__file__).resolve().parents[2]
FILE = "internal/plugins/region/provider_locations_collector.go"
MUTATIONS = (
    ("version", 'first.RawQuery = "api-version=2021-04-01"', 'first.RawQuery = "api-version=2020-01-01"', "TestProviderCollectorLiteral"),
    ("selected-scope", 'scope[id] == ""', '(scope[id] == "" && false)', "TestProviderCollectorAdmission"),
    ("returned-id-structure", 'len(parts[1]) == len("subscriptions") && strings.EqualFold(parts[1], "subscriptions")', 'len(parts[1]) == len("subscriptions")', "TestProviderCollectorIDs"),
    ("returned-id-binding", 'strings.EqualFold(parts[4], namespace)', '(strings.EqualFold(parts[4], namespace) || true)', "TestProviderCollectorIDs"),
    ("namespace-ascii", '!providerIdentifier.MatchString(namespace)', '(!providerIdentifier.MatchString(namespace) && false)', "TestProviderCollectorShapes"),
    ("location-ascii", 'if r > 127 {', 'if r > 127 && false {', "TestProviderCollectorShapes"),
    ("logical-location", '(normalized == "global" || costPhysical(normalized))', '(normalized == "global" || costPhysical(normalized) || true)', "TestProviderCollectorShapes"),
    ("mandatory-array", 'len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null"))', '(len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null"))) && false', "TestProviderCollectorShapes"),
    ("json-whole-guard", '!quotaJSON(ctx, body)', '(!quotaJSON(ctx, body) && false)', "TestProviderCollectorShapes"),
    ("same-page-collision", 'if _, exists := page.locations[key]; exists {', 'if _, exists := page.locations[key]; exists && false {', "TestProviderCollectorAtomicHealth"),
    ("prior-page-collision", 'if _, exists := prior[key]; exists {', 'if _, exists := prior[key]; exists && false {', "TestProviderCollectorAtomicHealth"),
    ("partial-prefix", 'if admitted > 0 {', 'if admitted > 0 && false {', "TestProviderCollectorAtomicHealth"),
    ("request-status", 'response.StatusCode != http.StatusOK', '(response.StatusCode != http.StatusOK && false)', "TestProviderCollectorAtomicHealth"),
    ("unsafe-continuation-health", 'return fail("provider_unsafe_continuation")', 'return result, nil', "TestProviderCollectorContinuation"),
    ("canonical-cycle", 'if seen[u.String()] {', 'if seen[u.String()] && false {', "TestProviderCollectorContinuation"),
    ("page-limit", 'if admitted >= MaxQuotaPages {', 'if admitted >= MaxQuotaPages && false {', "TestProviderCollectorBounds"),
    ("page-bytes", 'len(body) > MaxQuotaPageBytes ||', '(len(body) > MaxQuotaPageBytes && false) ||', "TestProviderCollectorBounds"),
    ("aggregate-bytes", 'len(body) > MaxQuotaTotalBytes-total', '(len(body) > MaxQuotaTotalBytes-total && false)', "TestProviderCollectorBounds"),
    ("raw-work", 'if page.work >= MaxProviderLocationWork {', 'if page.work >= MaxProviderLocationWork && false {', "TestProviderCollectorBounds"),
    ("deadline-health", 'if errors.Is(err, context.DeadlineExceeded) {', 'if errors.Is(err, context.DeadlineExceeded) && false {', "TestProviderCollectorCancellationOwnership"),
    ("cancel-sentinel", 'return nil, context.Canceled', 'return nil, err', "TestProviderCollectorWrappedCancellation"),
    ("deadline-sentinel", 'return nil, context.DeadlineExceeded', 'return nil, err', "TestProviderCollectorWrappedCancellation"),
    ("raw-duplicate-locations", 'if !charge() {\n\t\t\t\t\treturn providerPage{}, "provider_work_limit"\n\t\t\t\t}\n\t\t\t\tvar value string', 'if !charge() && false {\n\t\t\t\t\treturn providerPage{}, "provider_work_limit"\n\t\t\t\t}\n\t\t\t\tvar value string', "TestProviderCollectorDuplicateLocationWork"),
    ("canonical-properties", 'if strings.EqualFold(key, name) && key != name {', 'if strings.EqualFold(key, name) && key != name && false {', "TestProviderCollectorCanonicalProperties"),
    ("decoded-text", 'return auxiliaryText(value, &page.text)', 'return auxiliaryText(value, new(int))', "TestProviderCollectorDecoderBoundaries"),
)

def run(directory, test):
    return subprocess.run(["go", "test", "-count=1", "-run", "^" + test + "$",
                           "./internal/plugins/region"], cwd=directory,
                          capture_output=True, text=True, timeout=120)

with tempfile.TemporaryDirectory(prefix="cloud-assess-provider-collector-control-") as temporary:
    directory = Path(temporary)
    for name in ("go.mod", "go.sum"):
        shutil.copyfile(ROOT / name, directory / name)
    for name in ("cmd", "internal"):
        shutil.copytree(ROOT / name, directory / name, ignore=shutil.ignore_patterns(".git"))
    path = directory / FILE
    original = path.read_text(encoding="utf-8")
    for name, before, after, test in MUTATIONS:
        if original.count(before) != 1:
            raise SystemExit("Provider collector mutation ambiguous or absent: " + name)
        baseline = run(directory, test)
        if baseline.returncode:
            raise SystemExit("Provider collector baseline failed: " + name + "\n" + baseline.stdout + baseline.stderr)
        try:
            path.write_text(original.replace(before, after), encoding="utf-8")
            outcome = run(directory, test)
        finally:
            path.write_text(original, encoding="utf-8")
        if outcome.returncode == 0 or "--- FAIL: " + test not in outcome.stdout or "[build failed]" in outcome.stdout or "panic:" in outcome.stdout + outcome.stderr:
            raise SystemExit("Provider collector mutation survived or failed outside named assertion: " + name + "\n" + outcome.stdout + outcome.stderr)
        restored = run(directory, test)
        if restored.returncode:
            raise SystemExit("Provider collector restored baseline failed: " + name + "\n" + restored.stdout + restored.stderr)
        print("Provider collector mutation rejected:", name)
