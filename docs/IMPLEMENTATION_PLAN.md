# Cloud Assess Implementation Plan

Status: working implementation plan

Reference implementation: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Current implementation milestone: the generic core scan path and deterministic semantic source-versus-target equivalence harness are implemented. Live Azure source-versus-target regression is the next major phase.

## Principle

Cloud Assess is rebuilt behavior-first, not by blindly copying the source tree.

For each major subsystem:

```text
Observe source behavior
  -> Create characterization tests
  -> Define target interface
  -> Implement target subsystem
  -> Run equivalence tests
  -> Record intentional differences
```

## Core v1 development order

1. Repository bootstrap
2. Branding abstraction
3. Neutral domain model
4. Characterization-test harness
5. Configuration and filters
6. Azure authentication/environment
7. Scope discovery
8. Resource discovery
9. Scanner registry
10. Rule catalog
11. Azure Resource Graph client
12. Core recommendation execution
13. Diagnostics
14. Advisor
15. Defender
16. Policy
17. Arc SQL
18. Cost
19. Stage health and assessment completeness
20. Canonical assessment result and findings summary
21. JSON renderer
22. CSV renderer
23. Excel renderer
24. SARIF renderer
25. CLI integration and production orchestration
26. Severity and exit-code gates
27. Equivalence harness
28. Live Azure regression suite
29. CI and packaging
30. Security/license review
31. Plugin migration

## Current completion boundary

Steps 1-27 are implemented: the generic core scan path is quality-gated and the development equivalence harness can normalize the pinned reference table-JSON and Cloud Assess canonical JSON into comparable semantic datasets.

The current executable path is:

```text
cloud-assess scan
  -> validate filters/stages/gate
  -> DefaultAzureCredential
  -> resolve subscription or management-group scope
  -> discover and filter inventory
  -> build/prune pinned recommendation catalog
  -> execute Graph recommendations
  -> execute enabled auxiliary stages
  -> build canonical assessment result
  -> render requested reports
  -> apply completeness / severity exit semantics
```

Implemented Azure-data subsystems include Diagnostics, Advisor, Defender status/recommendations, Azure Policy noncompliance, Arc-enabled SQL inventory/status, and previous-month Cost Management data.

Implemented output formats are XLSX, canonical JSON, CSV, SARIF 2.1.0, and canonical JSON stdout. Excel remains enabled by default.

Exit semantics are finalized:

```text
0 = complete successful assessment and severity gate passed
1 = execution, configuration, authentication, or rendering failure
2 = quality/severity gate failed
3 = partial assessment because a requested noncritical stage failed
```

Reports are rendered before exit 2 or 3 is returned, preserving evidence for CI and troubleshooting. When a critical stage returns a partial result plus an error, requested reports are also persisted when possible before exit 1.

The CI foundation from step 29 is already partly implemented ahead of sequence. It currently checks pinned source-data provenance, formatting, module graph cleanliness, branding boundaries, the actual CLI build, root/scan help and version smoke tests, race-enabled tests, and `go vet`. Packaging/release automation remains future work.

## Current known gaps

The generic core `scan` path is runnable, but core v1 is not yet declared equivalent or release-complete.

Outstanding work includes:

- live Azure regression against the pinned reference
- resolution of the live Arc SQL `vcores` response shape
- external/YAML plugin execution in production orchestration
- internal plugin migration/parity
- scanner-specific CLI commands
- `rules` CLI command
- `plugins list/info` CLI surface
- final dependency/license inventory
- packaging/release artifacts
- security and operational review at distributable-product level

Until plugin execution is implemented, explicitly enabling the plugin stage in the core-v1 CLI returns a clear configuration error before Azure authentication.

## Characterization levels

### 1. Deterministic unit behavior

No Azure connection required.

Covered priority cases include:

- filter precedence
- tag matching
- resource-group validation
- stage defaults and validation
- stage parameter validation
- severity thresholds
- previous-calendar-month cost period
- subscription-ID redaction
- finding deduplication
- summary calculation
- recommendation applicability
- scanner registry mappings
- recommendation normalization
- report field mappings

### 2. Fixture-based components

Sanitized Azure/API fixtures cover mappings such as:

- ARG row -> Finding
- Advisor response -> Advisor record
- Defender pricing row -> Defender plan status
- Defender assessment -> Defender recommendation
- Policy state -> Policy record
- Arc SQL row -> Arc SQL record
- Cost Management row -> Cost record
- resource row -> Resource
- diagnostic-settings batch response -> diagnostic finding

### 3. Cross-package and report characterization

The target now includes a cross-package path using fake Azure operations but the real:

```text
Coordinator
  -> Application Runner
  -> Canonical Result
  -> JSON Renderer
  -> Excel Renderer
```

This verifies orchestration-to-report contracts without requiring live Azure.

Report characterization focuses on semantic output for:

- JSON
- CSV
- XLSX
- SARIF

Excel comparison focuses on worksheet names, headers, rows, ordering, counts, redaction and relevant formatting contracts rather than raw XLSX ZIP bytes.

### 4. Live Azure equivalence

Next major phase:

```text
Select stable Azure test scope
Run pinned reference
Run Cloud Assess
Normalize both outputs
Compare findings and auxiliary datasets
Classify every delta
Repeat with targeted fixtures for missing scenarios
```

Primary finding comparison key:

```text
Recommendation ID
Resource ID
Category
Impact
Source
```

A pre-existing non-production Azure test environment is suitable for the first pass. Targeted Terraform fixtures should be added only for behaviors not represented there.

## Existing reference fixtures to reuse

The pinned source already includes useful integration scenarios for:

- authentication
- Redis Enterprise
- VNet DNS
- generic resource findings
- Storage diagnostic settings
- Storage HTTPS
- Storage TLS
- Storage immutable versioning

These form the initial targeted live-equivalence supplement.

## Additional live characterization scenarios

Validate:

- subscription include/exclude
- resource-group include/exclude
- scanner/resource-type selection
- resource exclusion
- recommendation exclusion
- include tags
- exclude tags
- include/exclude tag precedence
- multiple subscriptions
- management-group recursion
- recommendation not applicable
- recommendation compliant
- recommendation noncompliant
- duplicate findings
- diagnostic setting present/absent
- Advisor data present
- Defender pricing status
- Defender unhealthy recommendation
- Policy noncompliance
- Arc SQL presence and numeric-vCore response shape
- Cost available
- Cost unauthorized/unavailable

## Intentional differences from source

The following differences are deliberate and should not fail equivalence tests:

- new product and CLI identity
- no legacy CLI compatibility requirement
- no legacy configuration compatibility requirement
- `GraphResult` replaced by `Finding`
- `GraphRecommendation` replaced by `RecommendationDefinition`
- lower-level packages return errors instead of calling `log.Fatal` or silently returning nil datasets
- stage status and assessment completeness are explicit
- subscription masking is named subscription-ID redaction
- deterministic ordering is added where source map/concurrency ordering was unstable
- custom rule source/path is neutral rather than legacy-branded
- Diagnostics findings identify `Azure Resource Manager` as their validation mechanism instead of inheriting a generic Azure Resource Graph label
- non-success Diagnostics subrequests retain source-compatible finding semantics while producing explicit uncertainty warnings
- malformed diagnostic-setting IDs become warnings rather than panic-prone parsing
- Advisor metadata is retrieved through the shared authenticated ARM HTTP layer rather than adding the `armadvisor` SDK dependency
- malformed Advisor/Defender/Policy/Arc SQL rows become explicit warnings
- Cost Management uses the shared authenticated ARM HTTP layer rather than adding the `armcostmanagement` SDK dependency
- Cost records populate subscription display names from the already-discovered subscription map, correcting a pinned source stage bug
- Arc SQL preserves the pinned source `vcores` string-decoder contract pending live-equivalence evidence
- SARIF uses Cloud Assess branding and a Cloud Assess fingerprint namespace

## Output and redaction strategy

JSON is canonical machine-readable assessment state. Excel remains the default human-facing report. Both are generated from the same canonical assessment result.

When subscription-ID redaction is enabled, XLSX, CSV, JSON, and JSON stdout mask subscription IDs, including subscription IDs embedded in serialized ARM/resource strings.

SARIF intentionally retains stable Azure resource identities because its results and fingerprints are intended for automation and baselining, matching the identity-bearing behavior of the pinned reference. SARIF should therefore be treated as sensitive/identity-bearing output.

Exact source-versus-target equivalence runs should disable redaction where stable raw resource identity is required for comparison.

## Licensing and attribution

Cloud Assess is currently licensed under Apache 2.0 at the repository level.

Reused or derived MIT-licensed source and recommendation material retains applicable copyright and license notices, including Microsoft-originated source, APRL material, Azure Orphan Resources material, and other third-party dependencies as required.

A generated/maintained dependency-license inventory remains part of the release process.

## Definition of done for core v1

Core v1 is complete when:

- repository builds independently
- no legacy product-facing branding remains
- required legal attribution remains
- scope/filter behavior is equivalent
- rule loading and scanner pruning are equivalent
- normalized core findings are materially equivalent for representative Azure inputs
- Diagnostics, Advisor, Defender, Policy, Arc SQL and Cost are reproduced
- stage failures are explicitly visible
- assessment completeness is explicit
- Excel and JSON are both first-class outputs
- CSV and SARIF work
- redaction and severity gates work
- unit, cross-package, race, vet and executable smoke tests pass
- selected live Azure equivalence tests pass
- remaining CLI/plugin gaps are either implemented or explicitly deferred from the release target
- security and license checks pass
