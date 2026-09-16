# Characterization Baseline

Reference repository: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

This document tracks source behavior that Cloud Assess intentionally preserves or changes.

## Implemented deterministic contracts

### Filtering and scope

- Filter identifiers are compared case-insensitively where the reference does so.
- Tag keys are case-insensitive; tag values are exact-match.
- All include tags must match.
- Any matching exclude tag removes a resource.
- Explicit subscription/resource-group include scope takes precedence over matching exclusion, matching the reference implementation.
- `include.resourceTypes` preserves the reference meaning: values are scanner/service keys such as `aks`, `ca`, and `st`, not literal ARM type strings.
- A normal scan with multiple available scanner keys lets `include.resourceTypes` choose scanner keys.
- A scanner-specific command with one scanner key takes precedence over `include.resourceTypes`.
- Selected scanner keys are expanded to an allowed ARM resource-type set.
- Scanner/subscription/resource-group/individual-resource structural scope is applied during inventory discovery and reapplied to downstream findings.
- Child resource/service results inherit a recorded parent tag-scope decision.
- Unknown tag scope is excluded for include-tag filters and retained for exclude-only tag filters.
- Recommendation exclusions are case-insensitive.

### Findings and gates

- Recommendation severity ordering: High > Medium > Low.
- Finding summary deduplicates by normalized recommendation ID plus resource ID.
- SLA-category definitions and findings are excluded from the ordinary findings summary.
- Findings without a known recommendation definition are excluded from the ordinary findings summary.
- Summary ordering is deterministic by resource type and then recommendation ID.
- The quality gate counts findings at or above the configured High/Medium/Low threshold.

### Stage configuration

- Default stages are graph, diagnostics, advisor, and defender.
- Graph remains mandatory for a normal assessment.
- Invalid stage names return errors rather than terminating the process. This is an intentional reliability change.

### Azure/resource helpers

- Resource ID subscription, resource-group, resource-group-ID, resource-type, and resource-name parsing preserves the pinned reference behavior.
- Subscription-ID redaction preserves the reference mask shape.
- Cost reporting uses the previous completed UTC calendar month.

### Scanner and recommendation catalog

- The scanner registry contains 87 keys and preserves special AVD, AVS, HPC, SAP, VNet, and dual-Redis mappings from the pinned source.
- APRL is pinned to `60eaddda76541f6adbc1c5ffa686829807e55e29`.
- AOR and custom-rule imported tree revisions are recorded in `internal/rules/provenance.go`.
- Embedded catalog precedence is APRL, then AOR, then custom rules; later definitions replace an identical resource-type + recommendation-ID pair.
- Embedded disabled/development/non-ARG rules preserve reference skip semantics.
- External YAML plugin rules preserve the reference's less restrictive execution filtering semantics.

### Azure Resource Graph

- Up to 300 subscriptions are sent per ARG request.
- ARG result pages request up to 5,000 rows and follow skip tokens.
- Management-group-aware authorization scope is opt-in for queries that require it.
- Individual malformed rows are skipped rather than discarding valid rows.
- Malformed-row counts are surfaced as target warnings; this is an observability improvement over source log-only behavior.
- JSON strings are unquoted while object/array/primitive projected parameters preserve their JSON text representation.
- Missing `id` stops processing the remaining rows for that recommendation, preserving source behavior.
- Unsupported/disallowed Resource Graph logical-table errors become warnings/skips.
- Other recommendation-query failures are returned as errors instead of terminating from a worker goroutine. This is an intentional reliability change.
- Recommendation execution uses bounded concurrency with a default of 10 workers.
- Findings and rule warnings are sorted deterministically after concurrent execution. This is an intentional determinism improvement.

### Discovery

- Subscription discovery excludes Disabled and Deleted subscriptions.
- Explicit subscription selection and subscription filters preserve reference behavior.
- Management-group discovery recursively resolves descendant management groups and subscriptions.
- Management-group traversal uses a visited-group guard to avoid duplicate traversal without changing the resolved subscription set.
- Scope IDs are deterministic SHA-256 hashes of sorted subscription IDs.
- Resource inventory uses the pinned source ARG projection for identity, tags, SKU fields, capacity, and kind.
- Resource inventory is partitioned into in-scope and out-of-scope collections.
- Malformed inventory rows are skipped and counted.

## Known intentional differences

- Product/CLI/configuration branding is neutralized.
- The configuration root is `assessment:` rather than the legacy product name.
- The legacy `exclude.services` configuration concept is named `exclude.resources` in Cloud Assess.
- Lower-level errors are propagated rather than calling `log.Fatal`.
- Deterministic ordering is added where source map/concurrency ordering was unstable.
- Cloud Assess exposes explicit warning/completeness signals for malformed ARG rows.

## Next characterization targets

1. Stage parameter parsing (`--stage-param`) and stage-specific option semantics.
2. Diagnostics recommendation definitions and ARM batch behavior.
3. Advisor result retrieval and normalization.
4. Defender status and Defender recommendation behavior.
5. Azure Policy noncompliance behavior.
6. Arc SQL behavior.
7. Cost API result retrieval and normalization beyond the already-characterized date range.
8. Canonical assessment-result assembly and completeness calculation.
9. JSON, Excel, CSV, SARIF, and stdout rendering equivalence.
10. End-to-end CLI exit-code and severity-gate behavior.
