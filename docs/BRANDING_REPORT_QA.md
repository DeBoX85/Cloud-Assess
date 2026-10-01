# Paired immutable-profile report acceptance

Date: 2026-10-02 (Europe/Oslo). Baseline: merged PR #73, `9cb3cfe0ff639d80e2cde8bbb8996208ed757adb`. This implements step 3 of [BRANDING_REVIEW.md](BRANDING_REVIEW.md), complementing the [profile and actual CLI builder checks](BRANDING_PROFILES.md). No production code, defaults, source/dependency pins or comparator normalization changes are required by the passing report checks. AR-03 remains partial until profile-aware native packaging and installation pass.

## What is exercised

The native checker `scripts/tests/branding-reports.py` compiles the real application test binary twice, once with built-in identity and once with an encoded custom profile through the same dedicated Go linker variable as the supported builder. `TestBrandingReportEvidence` supplies a deterministic synthetic canonical assessment to the production application Runner and all production renderers. There is no runtime test hook in the CLI, no successful Azure scan and no assertion that this proves live assessment or catalog parity. Synthetic APRL/AOR/CUSTOM recommendation labels, IDs and guidance links must remain unchanged; actual pinned-library provenance remains protected by existing native checks.

Each profile produces default timestamped outputs, operator-selected paths with spaces, and raw/unredacted outputs. Every set has twelve non-empty XLSX sheets, twelve CSV tables, JSON, stdout and SARIF. All auxiliary datasets have at least one synthetic row. The current profile's reportFilePrefix is independently checked through the actual default-path artifacts. Explicit output paths must win. Private Unix report mode 0600 is checked; Windows privacy remains governed by directory ACLs, not FileMode.

The checker independently reads the generated XLSX ZIP/XML, rather than deriving expectations solely from the Go renderer or its table projection. All twelve A1 cells must contain the expected title with no formula. The custom title deliberately starts with `=1+1` and contains quotes/ampersand, establishing literal text handling rather than formula execution. Worksheet cell values and formulas below/around the title must agree between profiles. XLSX ZIP bytes are not compared because container/metadata byte equality is not the report-data contract.

JSON, stdout and every CSV are compared byte-for-byte. Raw subscription/resource identities must survive unredacted output and raw subscription IDs must be absent from redacted canonical output. Source labels and finding recommendation IDs have independent expected sets. SARIF is compared structurally after normalizing exactly driver.name and driver.informationUri. All other fields, including version, rules, results, locations, guidance links, fingerprints and automationDetails, must agree. Three non-empty findings/fingerprints prevent empty-result equivalence from masquerading as identity coverage. The fixed emitted namespaces are retained; this does not prove cross-tool alert continuity when a consumer sees a different tool name. See GitHub's [SARIF tool and fingerprint contract](https://docs.github.com/en/code-security/reference/code-scanning/sarif-files/sarif-support).

Four concurrent runners per profile use distinct output paths and shared immutable profile/result inputs. Their JSON, SARIF and decoded workbook data must match the sequential result. The Go fixture also verifies that the canonical in-memory assessment was not mutated. Linux paired binaries use the race detector; Windows checks execute native builds without claiming race instrumentation. An HTTP tripwire and isolated Azure configuration detect unexpected authentication/network attempts by the synthetic report execution. Build commands use the normal trusted Go environment and may resolve dependencies; the tripwire applies to execution of the generated test binaries.

## Negative controls and execution

The same paired evidence comparison must reject five independently altered evidence objects: JSON, stdout, a CSV row, an XLSX recommendation cell and SARIF automation identity. These are checker controls after permitted presentation normalization, not production mutation or proof of all possible defects.

```text
python scripts/tests/branding-reports.py
```

Both required quality/windows-validation jobs execute this command alongside the existing default/custom/malformed CLI builder checks and all fourteen package cases. Local paired checks, focused application/renderer race checks and vet passed before publication. Exact-head native CI/head/tree/merge evidence remains blocking and is retained in the PR, then indexed in the next ledger checkpoint. FN-029 records fixture preparation mistakes without counting their failed runs as QA.

Next implement step 4: canonical profile/hash-aware candidate producer/verifier, safe executable/archive names, matching installation instructions, exact identity/version checks and native custom-package/tampering controls. Preserve clean stamped checkout, dependency inventory, notices, checksums, payload allowlist and private-report safeguards. The subsequent [profile-aware package implementation](BRANDED_PACKAGES.md) adds this workflow, pending final native acceptance; no release approval is provided.
