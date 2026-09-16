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

### Diagnostics

- The diagnostics subsystem preserves the pinned reference diagnostic-settings support table.
- 41 supported resource types have dedicated missing-diagnostics recommendation definitions; additional legacy/system types are queried but intentionally do not create a dedicated recommendation.
- Diagnostic settings are validated directly through Azure Resource Manager, not Azure Resource Graph.
- The ARM batch endpoint is `/batch?api-version=2020-06-01`.
- Each subrequest queries `<resource-id>/providers/microsoft.insights/diagnosticSettings?api-version=2021-05-01-preview`.
- Only resource types in the diagnostics support table are included in ARM batch requests.
- Batch size remains 20 resources and concurrency is capped at 30 workers.
- A resource with at least one returned diagnostic setting is treated as compliant for the dedicated diagnostic recommendation.
- A supported resource without settings receives its dedicated recommendation when one exists.
- Downstream structural/tag filtering is reapplied before a diagnostic finding is emitted.
- The source storage-account behavior is characterized: a storage account without settings receives `st-001`; a storage account with settings does not.
- A non-success ARM batch subrequest preserves the reference output behavior, which can still result in a missing-diagnostics finding, but Cloud Assess adds a `diagnostics_subrequest_non_success` warning so uncertainty is visible.
- A malformed diagnostic-setting ID is converted to a `diagnostics_malformed_setting_id` warning rather than risking a panic.
- Batch/transport/decode failures are returned to the caller rather than terminating the process.
- Diagnostics findings use source `DIAGNOSTICS` and validation mechanism `Azure Resource Manager`. This intentionally corrects the reference report's historical tendency to describe the canonical finding stream as ARG-validated even when a finding came from direct ARM validation.

### Advisor

- Advisor remains an auxiliary assessment dataset rather than being folded into the primary `Finding` stream.
- Recommendation instances are retrieved from the `AdvisorResources` Azure Resource Graph table.
- The pinned query only considers `microsoft.advisor/recommendations` records with a non-empty assessed resource ID.
- Suppressed recommendations are excluded when `properties.suppressionIds` is present and non-empty.
- Results are deduplicated by resource ID plus recommendation type ID using `summarize take_any(*)`.
- Human-readable recommendation descriptions are resolved through the Microsoft Advisor metadata contract at `/providers/Microsoft.Advisor/metadata?api-version=2020-01-01`.
- Advisor metadata pagination follows `nextLink` until exhausted.
- Only the metadata entity named `recommendationType` contributes recommendation ID to display-name mappings.
- Subscription exclusion and downstream resource/service exclusion are reapplied before an Advisor record is emitted.
- Advisor records preserve recommendation ID, subscription ID/name, resource type, resource name, resource ID, category, impact, and metadata-derived description.
- Missing metadata for a recommendation ID preserves the source behavior of an empty description rather than dropping the Advisor record.
- Malformed ARG rows are skipped and surfaced as `advisor_malformed_arg_rows` warnings.
- Advisor output ordering is deterministic in Cloud Assess. This is an intentional determinism improvement over source ordering.
- Metadata or ARG query failures are returned to the caller rather than being logged and silently converted to a nil dataset.
- Cloud Assess implements the Advisor metadata wire contract through the existing authenticated/retrying ARM HTTP layer instead of adding the `armadvisor` SDK dependency. This is an implementation change, not a behavior change.

### Defender

- Defender is represented as two separate auxiliary datasets: plan/tier status and security recommendations.
- Defender status comes from `SecurityResources` records of type `microsoft.security/pricings`, joined to subscription containers to obtain the subscription display name.
- Defender status preserves subscription ID, subscription name, plan name, and pricing tier.
- The Defender status path reapplies subscription exclusion but does not apply resource-level filtering because the records are subscription/plan scoped.
- Defender recommendations come from `microsoft.security/assessments` with `properties.status.code == 'Unhealthy'`.
- Recommendation metadata categories are expanded with `mvexpand`, so one unhealthy assessment can yield category-specific rows.
- The pinned recommendation query derives resource group, resource type, and resource name by splitting `properties.resourceDetails.Id`.
- The source query projects `ResourceType = tostring(ResourceIdsplit[6])`; Cloud Assess preserves the projected value rather than silently replacing it with a full normalized ARM type.
- Source-side `distinct` is preserved in the query and Cloud Assess also preserves the source post-query deduplication key `(resourceID, category, recommendationName)`.
- Defender recommendation filtering reapplies downstream resource/service exclusion using the assessed resource ID.
- Subscription names for Defender recommendations are resolved from the already-discovered subscription map; missing names remain empty.
- Defender portal links preserve the source normalization behavior by prefixing `https://` to the projected portal-link value.
- Malformed Defender status rows are skipped and surfaced as `defender_status_malformed_arg_rows` warnings.
- Malformed Defender recommendation rows are skipped and surfaced as `defender_recommendations_malformed_arg_rows` warnings.
- Defender status and recommendation outputs are sorted deterministically in Cloud Assess.
- ARG query failures and nil ARG results are returned to the caller rather than being logged and silently converted to nil datasets.

## Known intentional differences

- Product/CLI/configuration branding is neutralized.
- The configuration root is `assessment:` rather than the legacy product name.
- The legacy `exclude.services` configuration concept is named `exclude.resources` in Cloud Assess.
- Lower-level errors are propagated rather than calling `log.Fatal` or silently returning nil datasets.
- Deterministic ordering is added where source map/concurrency ordering was unstable.
- Cloud Assess exposes explicit warning/completeness signals for malformed ARG rows.
- Diagnostics carries its real validation mechanism (`Azure Resource Manager`) rather than labeling every primary finding as Azure Resource Graph validated.
- Non-success diagnostics subrequests preserve reference finding semantics but additionally produce explicit uncertainty warnings.
- Malformed diagnostic-setting IDs become warnings instead of panics.
- Advisor metadata uses the existing authenticated ARM HTTP layer rather than the `armadvisor` SDK while preserving the same API contract.
- Advisor malformed ARG rows are surfaced as warnings instead of being log-only behavior.
- Defender malformed ARG rows are surfaced as warnings rather than being log-only behavior.
- Defender datasets are sorted deterministically after normalization.

## Next characterization targets

1. Stage parameter parsing (`--stage-param`) and stage-specific option semantics.
2. Azure Policy noncompliance behavior.
3. Arc SQL behavior.
4. Cost API result retrieval and normalization beyond the already-characterized date range.
5. Canonical assessment-result assembly and completeness calculation.
6. JSON, Excel, CSV, SARIF, and stdout rendering equivalence.
7. End-to-end CLI exit-code and severity-gate behavior.
