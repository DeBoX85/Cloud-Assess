# Cloud Assess Implementation Plan

Status: working implementation plan

Reference implementation: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Current implementation milestone: the generic core scan path and deterministic semantic source-versus-target equivalence harness are implemented. Live Azure regression is in progress: default-stage, optional-stage, resource-group-scoped, two-subscription, leaf-management-group, scanner selection, RG/tag include/exclude, observed recommendation exclusion and exact-resource exclusion passes are equivalent for the data compared. Several Diagnostics batch warnings remain uncorrelated. A local check classified the missing cross-run Advisor row as having unknown recorded tag scope in the target inventory, while historical Azure timing remains unresolved. Parent-to-child management-group traversal and non-empty evidence for currently absent datasets remain next.

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

Steps 1-27 have the generic-path implementation described below; step 2 now has the accepted bounded immutable five-field customization workflow (PR #75). Bounded YAML/KQL Graph extensions are implemented with offline evidence; internal table plugins and remaining live extension validation remain gaps. Step 28 is in progress, with equivalent live baselines recorded for default stages, optional Policy/Defender Recommendations/Cost execution, resource-group and two-subscription scopes, a leaf management group, separate scanner, RG and tag include/exclude, same-RG include/exclude precedence, recommendation-exclude and exact-resource-exclude filters. Several passes retain Diagnostics warning limits. The missing Advisor row has unknown recorded tag scope in the tag-only target report, but the separate runs cannot establish unchanged Azure state. Step 29 has enforced Linux and Windows CI jobs plus development-candidate ZIP packaging and isolated extraction checks for their amd64 targets. Release publication, fresh-OS installation and release approval remain incomplete. These numbered steps are not equal-sized progress units; see the dated estimates and limitations in [ROADMAP.md](ROADMAP.md#progress-estimate-and-next-checkpoint).

The ordered remaining work, evidence gates, and release-scope decisions are maintained in [ROADMAP.md](ROADMAP.md).

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

The CI foundation from step 29 is already partly implemented ahead of sequence. It checks pinned source-data provenance, formatting, module graph cleanliness, branding boundaries, PowerShell validation-helper behavior, the actual CLI build, root/scan help and version smoke tests, race-enabled tests, a minimum statement-coverage floor, `go vet`, and reachable vulnerabilities. External actions are commit-pinned. Required CI additionally executes built-CLI offline preflight/report-preservation cases and compares compiled dependencies with the generated inventory under CGO-disabled amd64 builds; see [BUILT_CLI_VALIDATION.md](BUILT_CLI_VALIDATION.md). Development candidate packaging adds notices/checksums/build metadata and isolated extraction/runtime checks; see [PACKAGE_BUILDS.md](PACKAGE_BUILDS.md). Successful live installed scans, release publication and approval remain future work.

## Current known gaps

The generic core `scan` path is runnable and semantically equivalent for the live behaviors exercised so far, but core v1 is not yet declared release-complete.

Outstanding work includes:

- investigation of two Diagnostics HTTP 400 subrequests in the two-subscription pass
- parent-to-child management-group live regression against the pinned reference; the leaf pass is complete
- non-empty live or targeted-fixture evidence for Policy and Defender Recommendations
- resolution of the live Arc SQL `vcores` response shape
- live YAML/KQL extension equivalence on a suitable environment
- internal plugin migration/parity
- `plugins list/info` CLI surface
- final release-specific dependency/license review and packaged notices
- packaging/release artifacts
- security and operational review at distributable-product level

Until internal table-plugin execution is implemented, explicitly enabling the plugin stage in the core-v1 CLI returns a clear configuration error before Azure authentication.

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

Current state and next boundary:

```text
Completed: default-stage subscription baseline
Completed: optional Policy/Defender Recommendations/Cost pass
Completed: resource-group-scoped and two-subscription passes (the latter with Diagnostics warnings and non-empty Defender plan status)
Completed: AdvisoryDev leaf management group, Storage/VM selection, RG include and tag-only include passes
Next: parent-to-child management-group traversal, remaining filter combinations and Advisor cross-run scope check
Also open: historical Diagnostics batch request correlation and targeted non-empty fixtures for absent scenarios
```

Primary finding comparison key:

```text
Recommendation ID
Resource ID
Category
Impact
Source
```

The existing non-production Azure environment supplied the first three live passes. Targeted Terraform fixtures should be added only for important behaviors not represented there.

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

A [generated dependency/license evidence inventory](DEPENDENCIES.md) now covers the selected module graph, Linux/Windows CLI membership, standard-library notices and bundled source pins, with CI freshness checks. Regeneration on the exact release tree, package-specific license review and inclusion of notices with artifacts remain part of the release process.

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

## Access and operations evidence boundary

[ACCESS_MODEL.md](ACCESS_MODEL.md) and [OPERATIONS.md](OPERATIONS.md) map current request semantics, scope/credential selection and recovery. Execution completeness does not establish intended estate visibility: discovery intersects requested IDs with visible subscriptions and ARG may silently omit inaccessible objects. Canonical JSON now reports requested/resolved scope and fails before resource queries for unresolved explicit subscriptions; see [SCOPE_REPORTING.md](SCOPE_REPORTING.md). Gate 004 still requires independent membership/access and representative live evidence; no production or minimum-role certification is claimed. Cost remains subscription-wide under RG/tag/resource filters.

The eight-task offline QA checkpoint adds a requirement-to-test map, combined scope fixtures, request/credential boundary safeguards, default retry-count checks, documentation parser/help validation and incomplete package rejection; see [OFFLINE_QA_REVIEW.md](OFFLINE_QA_REVIEW.md). Advisor origin/cycle checks and default shared-client redirect refusal are deliberate security corrections for unsafe response paths, not new feature parity or release approval.

SDK middleware review also disables automatic provider registration in shared ARM/scope configuration. Synthetic registration-required errors must preserve failed retrieval with no write attempt; this correction and its regressions are logged as FN-019.

The [alignment sanity review](ALIGNMENT_REVIEW.md) compares these implementation milestones with eventual full AZQR parity and records SDK origin/token-audience corrections. Core implementation is not a claim that every source command or plugin is complete.

The [pagination/lifecycle checkpoint](PAGINATION_LIFECYCLE.md) adds SDK cycle errors, context checks and an opt-in cooperative assessment budget without declaring universal scan/load bounds. Zero timeout preserves existing behavior; report persistence remains outside the budget.

The [ARG completeness checkpoint](ARG_COMPLETENESS.md) implements two documented target corrections: service-contract page size 1000 and explicit tokenless-truncation failure. Normal query text, subscription batching, valid continuation behavior and empty-result handling are preserved; metadata/load and live large-result limits remain open.

The [branding design review](BRANDING_REVIEW.md) supplies the implementation batch and acceptance checks. The first two steps now have a [validated profile and native development builder](BRANDING_PROFILES.md), with real default/custom executable checks. Five fields have presentation consumers; four remain unused. The subsequent [paired report checks](BRANDING_REPORT_QA.md) implement custom report/data acceptance using production renderers and synthetic input. [Profile-aware development candidate packages](BRANDED_PACKAGES.md) passed final native acceptance in PR #75 for the bounded immutable five-field development workflow. Release distribution and approval remain separate.

The [source-feature characterization](FEATURE_PARITY_CHARACTERIZATION.md) now maps separate YAML Graph versus internal table execution, missing commands and six plugin dependencies. Next implement offline rules inspection. Production plugin/command parity and live validation remain open; source failure swallowing and unforwarded region history flags require explicit target decisions and regressions, not silent copying.

[Offline rules inspection](RULES_INSPECTION.md) now implements the source-selected catalog/Diagnostics command with an independent 380-row reference executable capture, Markdown safety, invalid-argument/writer-error tests and actual installed CLI tripwires. Next implement scanner-specific subcommands; external/internal plugin execution and release/live evidence remain open. Exact-head native acceptance remains mandatory before merge.

The [development execution plan](DEVELOPMENT_EXECUTION_PLAN.md) now defines bounded batches, scope controls, exact-head QA, failure feedback, durable proposal checkpoints and reviewed rollback. [Scanner commands](SCANNER_COMMANDS.md) implement B1 through existing orchestration with independent 87-key source capture and offline selection/preflight checks. After native acceptance, resume B2 bounded YAML Graph integration; no Azure/laptop tasks are required now. Earlier next-task statements are historical checkpoints.

B2 [bounded YAML Graph plugins](YAML_GRAPH_PLUGINS.md) now implements source-object parsing, deterministic discovery, independent catalog overlays and normal Graph execution. Exact-head native acceptance remains mandatory. Next B3 integrates a real zone-mapping plugin with canonical tables/health/reports, then verifies the shared path against a second plugin. Routine handoffs are not needed under the ongoing authorization recorded in the execution plan. Laptop/Azure acceptance stays deferred.

B2 is VERIFIED OFFLINE through PR #79 (merge `7cddf9368272dd425315bdf3635b21fb974a9fe4`, native run `36963458729`). B3a [zone adapter](ZONE_MAPPING.md) is the first bounded internal-plugin slice, with source/API-grounded pagination, partial-health inputs and body/request limits. Native acceptance remains mandatory before merge. B3 table/schema/privacy/report and CLI integration follows; an unexposed adapter does not establish completed plugin parity. No Azure/laptop task is needed for the next slice.

B3a is VERIFIED OFFLINE through PR #80 (merge `0b76ad6b2870432357e5a38c1e8fe73d5214f11b`, run `36966029765`). B3b [canonical plugin tables](PLUGIN_TABLES.md) adds owned versioned schema/health and validated report/privacy infrastructure; ordinary schema/fields/tables stay compatible. This requires final native acceptance and does not expose plugin CLI execution. Next integrate zone operations/selection/health/exit and actual commands, then verify a second real adapter.
