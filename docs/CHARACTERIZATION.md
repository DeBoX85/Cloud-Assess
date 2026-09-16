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
- At the pinned source commit, `--stage-param` only registers `plugin.target-regions` as a string option.
- Stage parameters use `stage.key=value`; malformed parameters, unknown stages/options, and type mismatches are errors.

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
- Diagnostics findings use source `DIAGNOSTICS` and validation mechanism `Azure Resource Manager`.

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
- Missing metadata preserves an empty description rather than dropping the record.
- Malformed ARG rows are skipped and surfaced as `advisor_malformed_arg_rows` warnings.
- Advisor output ordering is deterministic in Cloud Assess.
- Metadata or ARG query failures are returned to the caller rather than being logged and silently converted to a nil dataset.
- Cloud Assess implements the Advisor metadata wire contract through the existing authenticated/retrying ARM HTTP layer instead of adding the `armadvisor` SDK dependency.

### Defender

- Defender is represented as two separate auxiliary datasets: plan/tier status and security recommendations.
- Defender status comes from `microsoft.security/pricings` joined to subscription containers.
- Defender status preserves subscription ID/name, plan name, and pricing tier and reapplies subscription exclusion.
- Defender recommendations come from unhealthy `microsoft.security/assessments` and expand metadata categories with `mvexpand`.
- The source query projects `ResourceType = tostring(ResourceIdsplit[6])`; Cloud Assess preserves that projected value.
- Source-side `distinct` plus the post-query `(resourceID, category, recommendationName)` dedup key are preserved.
- Defender recommendation filtering reapplies downstream resource/service exclusion.
- Defender portal links preserve source normalization by prefixing `https://`.
- Malformed status/recommendation rows become explicit warnings.
- Defender outputs are sorted deterministically and ARG failures are returned.

### Azure Policy

- Azure Policy remains an auxiliary dataset of non-compliant policy states.
- The pinned query reads `microsoft.policyinsights/policystates` and filters to `complianceState == 'NonCompliant'`.
- Policy definition metadata is left-joined to obtain display name and description.
- Subscription containers are left-joined to obtain subscription display names.
- The query executes with management-group-aware ARG authorization scope enabled.
- Subscription and downstream resource/service filters are reapplied after query execution.
- Source post-query deduplication by `(resourceID, policyDefinitionID)` is preserved.
- Resource group, full resource type, and resource name are derived from the policy state's resource ID.
- Policy display/definition/assignment identities, timestamp, description, and compliance state are preserved.
- Malformed rows become `policy_malformed_arg_rows` warnings and failures are returned.
- Azure Policy output is sorted deterministically.

### Arc SQL

- Arc SQL uses the pinned multi-join Resource Graph query for `Microsoft.AzureArcData/sqlServerInstances`.
- The query joins the `WindowsAgent.SqlServer` Hybrid Compute extension and Hybrid Compute machine status.
- Source derivation of license (`Paid` -> `SA`, `PAYG` -> `PAYG`, otherwise `unset`) is preserved.
- DPS and telemetry status parsing/fallback logic is preserved in the KQL.
- Subscription and downstream resource/service filters are reapplied; the SQL instance resource ID is the resource-scope key.
- Subscription display names are resolved from the discovered subscription map.
- Malformed rows become `arcsql_malformed_arg_rows` warnings and failures are returned.
- Arc SQL output is sorted deterministically.
- The pinned query projects `vcores = toint(properties.vCore)` while the pinned row decoder expects a string. Cloud Assess currently preserves that decoder contract and explicitly characterizes numeric `vcores` as a malformed row. This remains a live-equivalence review item rather than an untracked source ambiguity.

### Cost Management

- Cost uses the exact previous completed UTC calendar month.
- The pinned query is `ActualCost` with `Custom` timeframe, `TotalCost = Sum(Cost)`, grouped by the `ServiceName` dimension.
- The wire endpoint is `/subscriptions/{id}/providers/Microsoft.CostManagement/query?api-version=2021-10-01`.
- The source Cost Management row contract is `[value, serviceName, currency]` and values are converted with generic `%v` formatting.
- The source stage runs at most two subscription workers concurrently; Cloud Assess preserves that ceiling.
- Cloud Assess queries subscriptions independently and then sorts the combined records deterministically.
- Azure response codes `MissingRegistrationForResourceProvider`, `MissingSubscriptionRegistration`, `DisallowedOperation`, and `NotFound` preserve source skip semantics but produce `cost_subscription_skipped` warnings.
- Non-skippable Cost Management failures are returned to the caller.
- Short/malformed rows become `cost_malformed_row` warnings instead of risking an index panic.
- Empty/204 Cost responses produce an empty dataset without failure.
- Cloud Assess uses the shared authenticated ARM HTTP layer instead of importing `armcostmanagement`, preserving the same API contract.
- Cloud Assess intentionally populates `SubscriptionName` from the discovered subscription map. The pinned source `CostStage` fails to pass the name into `ScannerConfig`, leaving that report column empty despite modeling it; this is treated as a target correctness fix.

## Known intentional differences

- Product/CLI/configuration branding is neutralized.
- The configuration root is `assessment:` rather than the legacy product name.
- The legacy `exclude.services` configuration concept is named `exclude.resources` in Cloud Assess.
- Lower-level errors are propagated rather than calling `log.Fatal` or silently returning nil datasets.
- Deterministic ordering is added where source map/concurrency ordering was unstable.
- Cloud Assess exposes explicit warning/completeness signals for malformed rows.
- Diagnostics carries its real validation mechanism (`Azure Resource Manager`) rather than a generic ARG label.
- Non-success diagnostics subrequests preserve reference finding semantics but additionally produce explicit uncertainty warnings.
- Malformed diagnostic-setting IDs become warnings instead of panics.
- Advisor metadata uses the existing authenticated ARM HTTP layer rather than the `armadvisor` SDK.
- Cost Management uses the existing authenticated ARM HTTP layer rather than the `armcostmanagement` SDK.
- Cost subscription names are populated, correcting the pinned source stage bug.
- Arc SQL retains the pinned `vcores` decoder mismatch until live equivalence establishes the correct current Azure response behavior.

## Next characterization targets

1. Stage health and overall assessment completeness behavior.
2. Canonical assessment-result assembly.
3. JSON, Excel, CSV, SARIF, and stdout rendering equivalence.
4. End-to-end CLI exit-code and severity-gate behavior.
5. Live-equivalence resolution of the Arc SQL `vcores` response shape.
