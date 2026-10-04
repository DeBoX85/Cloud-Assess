# Cloud Assess Implementation Plan

Status: working implementation plan

Reference implementation: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Current implementation milestone: the generic core scan path and deterministic semantic source-versus-target equivalence harness are implemented. Live Azure regression is in progress: default-stage, optional-stage, resource-group-scoped, two-subscription, leaf-management-group, scanner selection, RG/tag include/exclude, observed recommendation exclusion and exact-resource exclusion passes are equivalent for the data compared. Several Diagnostics batch warnings remain uncorrelated. A local check classified the missing cross-run Advisor row as having unknown recorded tag scope in the target inventory, while historical Azure timing remains unresolved. Parent-to-child management-group traversal is explicitly deferred as DV-001; non-empty evidence for currently absent datasets remains open. Current offline milestones also include scanner commands, YAML Graph execution, zone/service-health/SQL EOL execution and carbon pure/request/public execution through PR90/93/94. Executor/postmerge recovery is complete through PR95 and AI governance source captures/pure core are accepted offline through PR96. The bounded AI request library is now accepted offline through PR97; source-specific Graph discovery is accepted through PR101 and public/report integration through PR103; region primary scoring/projection is accepted offline through PR104; the auxiliary source runner is accepted offline through PR105, with pure quota/reservation projections accepted offline through PR106 and service-availability pure sheets accepted offline through PR108, with pure CostComparison accepted offline through PR109 and pure Inventory accepted offline through PR110 and separate source-only inventory calculations characterization IN PROGRESS/unaccepted; live tasks remain deferred while the operator is unavailable.

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


Current B4 advance (2026-10-03): service-health accepted through PR87, SQL EOL library through PR88. SQL execution candidate preserves source subscription-only/string projection across standalone/normal/scanner paths and three-plugin dispatch. Restored local race/vet and independent mutation checks passed; clean built/native acceptance precedes B5 carbon, AI governance and region selection migrations. See SQL_EOL and SESSION_HANDOVER for exact current checks and live deferrals.


B4 SQL execution accepted offline through PR89 after exact native/postmerge verification. B5 current carbon core includes actual source scanner/calculation captures and pure bounded aggregation; see CARBON_EMISSIONS and SESSION_HANDOVER. Finish its full/built/native gates, then separately implement bounded authenticated request/access/pagination and public execution/report integration. This is not live emissions/API permission validation.


Carbon pure aggregation core accepted offline through PR90; exact next request-adapter contract/acceptance matrix in CARBON_EMISSIONS. Implement authenticated date/report library, then separate public command/registry/report integration. Verified pure calculations do not establish service authorization/completeness or a finished carbon plugin.


Current B5 checkpoint (2026-10-03): carbon request library accepted through PR93 with exact native/postmerge verification. Public carbon standalone/mixed/scanner/registry/report integration is the separate unaccepted candidate; focused tests pass, final acceptance remains pending. CARBON_EMISSIONS and SESSION_HANDOVER hold the current contract/resume point. Earlier request-library-pending paragraphs are historical; AI governance then region selection follow separately.


Historical PR94 outage snapshot (superseded by PR95/96 recovery and the latest checkpoint below): carbon public execution passed final exact-head native checks and protected remote merge verification; local postmerge fetch/switch completion remains unknown after executor outage. SESSION_HANDOVER holds exact proof/recovery sequence. AI governance remains unimplemented and follows after recovery, then region selection separately; all live/Gate004/release boundaries remain open.


Historical pre-PR97-acceptance B5 checkpoint (superseded below): AI pure core is VERIFIED OFFLINE through PR96 merge9c6d0e79d05d06559fb4b543406b928013aed835. [AI_GOVERNANCE_REQUESTS.md](AI_GOVERNANCE_REQUESTS.md) / [PR97](https://github.com/DeBoX85/Cloud-Assess/pull/97) define the current unaccepted strict decoder/bounded metrics/deployment library. Source captures, full local race/vet, strict stamped build, eleven CLI cases, docs and three compiling request guard mutations pass. Final exact-head native steps/logs, current rules/base/head/preview and protected postmerge verification remain mandatory. Next: source-specific Graph account discovery/tag-scope correlation and real standalone/mixed/scanner/report integration; public ai-gov remains unavailable until that separate acceptance. Region selection/B6/B7 and earlier live/load/advisory/Gate004/release boundaries remain open. No laptop or Azure access required for current engineering.


Latest B5 checkpoint (2026-10-03): bounded AI requests VERIFIED OFFLINE through PR97 merge `aadd0d03b13a093a9ec3b9e69e3f70208d1e1714`, final head `f58895fcb9cd51817d49b83f26f71b7073443739`/run37092703311 with both mandatory jobs/full logs and remote tree/ordered parents/identity verified. Four compiling mutations include the fixed partial-cloud audience guard. Public AI remains unavailable. [AI_GOVERNANCE_EXECUTION.md](AI_GOVERNANCE_EXECUTION.md) defines bounded source-specific Graph discovery, honest recorded tag scope and separate actual standalone/mixed/scanner/report integration. Workspace recovery and fresh local offline QA are now accepted through PR100; bounded discovery is VERIFIED OFFLINE through PR101 merge3c4019c; next is separately accepted actual public/registry/report integration under the current handover. No live equivalence claimed. Existing user laptop/Azure, Gate004/release/module-only-advisory deferrals remain open.


## AI public execution WIP and interruption revalidation (2026-10-03)

Live accepted baseline032740d/treebf5880/ordered parents3c4019c and9edf28a, PR100-102 merges, human identity, unchanged rules23890737 and no open proposals verified. PR102 accepted-merge push run37106303179 now passed Linux111155398115/Windows111155398010 with all required steps/logs inspected; coverage81.9%, eight controls and broad existing QA, only conditional failure upload skipped. This completes the previous distinct pending push proof.

Previous scratch workspace/tools absent; accepted bootstrap recreated cloud-assess-recovery-20261003 with exact verified tools/clones/modules/source/APRL/provenance. Fresh pristine full Linux race/vet/stamped build, eleven actual CLI cases/docs307/9/1/eight controls/restored checks/5000x discovery fuzz passed on032740d. A later interruption retained all candidate files; actual inventory/status and current ref were checked before another edit.

Candidate on feat/ai-public-execution adds real AI standalone/mixed/scanner/registry/report integration, owned filter decisions and sanitized16-cell tables, public-cloud pre-auth/factory/scope guard, enclosing5-minute lifecycle and truthful retained discovery failure. Focused tests and corrected58 authenticated in-memory complete-cell all-format cases pass; new ownership/control/broad/native acceptance remains pending. FN060 records actual initial fixture failures. Publish/read back code/docs WIP, then finish tests and final exact-head native/protected QA. Source/dependency/capture pins and Azure/laptop/Gate004/release deferrals unchanged; reviewed protected revert after checking dependents remains rollback.


### Public AI candidate QA completed locally; final native acceptance next

Verified WIP82fbda557bf2e3d6990e846bb9b6bcf2b7d4be26/tree d0f425ca3d4f79058c2e889746e5faaa6c12269c, direct parent032740d, mandatory human identity and20 exact changed-file readbacks. Fetched object/tree matches; single-branch tracking failures were resolved without reset/force (FN061), preserving untracked owned tests. Post-interruption focused race,58 actual authenticated Cobra/factory all-format cases, earlier-deadline/empty/concurrent library tests, repeated/eight concurrent recorded-tag/scope/result snapshots, missing-operation/cloud-before-scope and critical pending-header cases passed. Added all four compiling public execution controls/restored checks to both native jobs; full Linux race/vet passed. Mixed/partial reports preserve independent normal findings and raw SARIF identity even with default JSON/CSV/XLSX/stdout masking. Final clean stamped twelve-CLI/default+branded packages/docs/security and both final-head native jobs/protected merge remain required. Source/capture/module pins and live/release deferrals unchanged.


## AI public acceptance and region primary core (2026-10-03)

PR103 AI public execution is VERIFIED OFFLINE at merge `8c8f8b9f3eb3e53e9d1829b9295d1bb170d7a4da`, tree `cbeaa60a893664d9cec3a4752b3422f7f0402a74`, ordered parents032740d / finald78be7. Final [PR103](https://github.com/DeBoX85/Cloud-Assess/pull/103) retains exact29-file readback, human identity, local full race/vet/12 compiled CLI/default+branded19-case packages/docs308/9/1 and required [run37127332387](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37127332387), quality111215198140/Windows111215198379, both tested previewdcedd209 with exact candidate tree/parents. Distinct accepted-push [run37127893618](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37127893618) also passed quality111216865978/Windows111216865886 on exact8c8f8b9; all required steps/full logs verified, only conditional failure upload skipped. Both native runs passed all eight request/four execution compiling controls,12 built CLI and19 default/19 branded cases per OS, Linux race/vet/fuzz, coverage81.4%, no reachable/imported vulnerabilities; existing one module-only advisory remains open. Clean fetched accepted checkout and unchanged actual AZQR/APRL verified. Self-review is not independent-person approval; no live/Gate004/release closure.

Region selection primary scoring/projection is IN PROGRESS under [REGION_SELECTION.md](REGION_SELECTION.md). A separate byte-identical pinned source/APRL copy called only unchanged calculateScores/generateOutputTable, producing42 synthetic literal comparison captures with exact25-column rows and source metadata. No Scan, credential/SDK construction or HTTP; production source and dependencies unchanged. Retained inputs/output/harness are byte protected. Target pure Project preserves source weights/unknown/restricted/zone/cost/latency/quality/bands and owns canonical scoped rows; rejects malformed/nonfinite/inconsistent/oversized inputs before reporting, bounds joined UTF-16 cells/aggregate text and deterministically orders equal score/name ties. Focused source-cell/guard/ownership/repeated42-concurrent/work-limit race tests pass. Final broad/native acceptance remains required; public region-selection stays unavailable. Next characterize/migrate the remaining source sheets/calculations, then bounded ARM/retail adapters and separate actual public/report execution. Both source empty sentinel branches, full-score input completeness, flag forwarding, singleton state, quotas/prices/latency and debug JSON remain explicit future obligations.


Primary region local checkpoint: complete42 independent source score/cell cases, foreign/malformed/bounded input,8192-row exact boundary, joined/aggregate text, ownership/order/empty/cancellation and repeated42-concurrent race tests pass. All three compiling primary controls and restored named baselines pass; workflow failure propagation6 cases and full Linux race/vet pass. Final clean compiled CLI/docs/package/security and final-head native/protected acceptance remain required. No new public plugin is advertised.


### PR104 aggregate-work review correction and executor recovery

Automated [Code Review](https://github.com/DeBoX85/Cloud-Assess/pull/104#discussion_r4173516981) identified missing joined-separator bytes in the pre-projection aggregate budget at original PR104 head596f102. Native run37129101097 passed on that head, but it remains UNACCEPTED and its green gates do not certify the correction. A new named regression failed on the original code. Corrected pre-projection accounting includes detail separators and mapping arrows/separators, and an independent65536 global detail/mapping-entry cap bounds CPU even for empty labels. New aggregate-empty-entry and near16MiB separator-specific regressions distinguish early rejection from final output validation; compiling removal of each guard fails its named assertion and restored code passes. Both native jobs now run five region controls. FN063 records the finding; new final exact-head native/review/merge proof is required. Source captures/cells/pins unchanged; no faulty head accepted.

After correction, focused/full Linux race, vet and all five compiling controls/restored baselines completed successfully. The later commit/build/package checkpoint command stopped returning; a fresh pwd request also did not return. Its local commit and build/package outcome are UNKNOWN. Do not infer completion or replay commits blindly. The correction is reconstructed from verified remote596f102 plus the reviewed bounded changes and published through the authenticated GitHub connector, with exact file/commit/tree/parent/identity readback. Fresh mandatory native Linux/Windows gates must validate this published correction; earlier local/native/package evidence cannot substitute. No alternative executor is exposed; available GitHub code publication/native Actions permit independent authorized work without credentials or user laptop/Azure tests. Recover actual local state only after executor requests return.


## Region primary acceptance and reproducible source runner (2026-10-03)

PR104 primary region scoring/projection is VERIFIED OFFLINE at merge `c270ea22b036054382f6a34ca072f9a57bcac333`, tree `b7a49c874f7bad2fb5f6b292beaec026ee9188e8`, ordered parents accepted8c8f8b9 / corrected16e572a. [PR104](https://github.com/DeBoX85/Cloud-Assess/pull/104) retains exact19-file/commit/tree/identity readback, completed corrected local focused/full race/vet/five compiling controls, the automated review finding/reproduction/fix and unknown local checkpoint boundary. Required corrected [run37130189380](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37130189380), quality111223592562/Windows111223592666, and distinct accepted-push [run37130591651](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37130591651), quality111224747914/Windows111224747853, passed every required step with full logs inspected; only conditional Linux failure upload skipped. Both OSes passed complete42 source score/25-cell/hash cases, all five compiling region/eight AI request/four execution controls,12 compiled CLI and19 default/19 branded package cases, docs316/9/1; Linux race/vet/fuzz and81.8% coverage passed. No reachable/imported vulnerability findings; module-only advisory stays open. Remote accepted commit/tree/parents/identity/ref verified. Local postmerge fetch cannot be claimed while executor is unresponsive; no source or dependency pin advance, live/Gate004/release closure or public region availability.

Current source characterization task is IN PROGRESS on feat/region-source-characterization under [REGION_SOURCE_CHARACTERIZATION.md](REGION_SOURCE_CHARACTERIZATION.md). Added a reproducible isolated pinned-source Python runner, retained pure-source Go harness and a read-only Ubuntu24.04/Go1.26.8 Actions workflow, using already supported platforms and pinned actions. It captures fourteen helper branches for actual literal Svc Avail <region>, CostComparison, Quota, Capacity Reservations and raw/masked/empty Inventory tables, retaining literal inputs and complete source outputs with original hashes through bounded JSON log chunks. It invokes no Scan, credentials or service client; original source/APRL/module files must remain unchanged. No auxiliary target runtime or public command is implemented. Fresh source-runner output/provenance, required quality/windows logs/review and protected acceptance remain pending. Existing local executor still does not return even pwd; unknown local files/commits/processes must be inspected on recovery, never reset or replayed. Use verified remote checkpoints and native Actions to continue independent authorized work without laptop/Azure requests.


Inventory source characterization PR111 is VERIFIED OFFLINE. The bounded owned aggregation helper is IN PROGRESS/unaccepted under [REGION_INVENTORY_AGGREGATION.md](REGION_INVENTORY_AGGREGATION.md). Public region-selection remains unavailable. Finish fresh exact-head native/source QA, review and protected acceptance, then availability/latency/cost/quota/reservation arithmetic, bounded request adapters and coordinator/public/all-format execution. Azure/laptop-dependent validation remains explicitly deferred.


PR108 service exact final and separate accepted-push proof are indexed in REGION_SERVICE_AVAILABILITY/current handover. CostComparison is independent; prior service QA does not certify its new code.


Current Inventory checkpoint2026-10-04: CostComparison is VERIFIED OFFLINE through [PR109](https://github.com/DeBoX85/Cloud-Assess/pull/109), merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents. 

Final candidate1cf8e4cf6adc38ffe377894824b51f3b26431373/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, parentca16157 on accepted1957f6d7. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165402 quality111338402524/Windows111338402644 passed every mandatory step; only conditional failure upload skipped. Both tested preview77c2b59d40f89e1966402668a98efcdfad2e4d5d with identical final tree/ordered base/head parents. Full logs inspected: all nine Cost controls/restored baselines, ten service/three auxiliary/five primary/eight request/four execution controls, actual CLI/default+branded package/docs/provenance/module/inventory/branding guards; Linux race/vet/fuzz and82.5% coverage; zero reachable/imported vulnerability findings, existing module-only advisory remains open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165342 /111338403371 completed; full14 branches, indexed3-input/4-output chunks, exact retained bytes/hashes and source/tree/APRL provenance verified. Complete14-file scope/deletions/remote bytes/original scratch blobs, parent/human identity/pins/base/head/preview/current no-bypass rules23890737 verified. Semantic self-review checked actual unchanged source first-tuple/item break behavior, all source topology/cells and safe bounded identity/nonfinite/text/cancellation behavior. No callable final manual Code Review tool or independent-person approval is claimed; automatic comments currently empty. FN072 compile failure was not accepted. Expected-head protected merge and distinct accepted-push evidence follow; no fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure.

Protected acceptance verified: PR109 merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents, Denis author/GitHub committer/livecore ref and no open proposals. Distinct accepted-push native/source QA remains pending; no dependent live proof claimed.

Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513599 quality111339415060/Windows111339415184 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513582 /111339414976 passed on exact6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f. All mandatory steps/full logs inspected; only conditional Linux failure upload skipped. Both hosts passed9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages/docs/branding/provenance; Linux race/vet/fuzz82.5%, zero reachable/imported findings. Source14 complete branches, indexed chunks/both retained exact bytes/hashes and pinned source/tree/APRL provenance verified. Module-only advisory and all live/public/laptop/Azure/Gate004/release obligations stay open. No fresh local Go/source or independent-person approval claimed.

ProjectInventory/new tests/eight compiling controls and explicit source-sheet composition correction are unaccepted; existing source/captures/SKU/APRL/dependency pins unchanged. Native exact-head full proof remains required; no public or live closure.


2026-10-04 historical PR111 preparation baseline, superseded below: Inventory is VERIFIED OFFLINE through [PR110](https://github.com/DeBoX85/Cloud-Assess/pull/110), merge2dcd8dae924ce3fad5402890cdeca0ecbef6e01e/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, ordered6dfcd1be/76ac3f8d parents. 

Final76ac3f8d65acf9970305c58186e13bb19cf78ffe/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, parent accepted6dfcd1be. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985397 quality111340777403/Windows111340777269 passed every mandatory step; only conditional failure upload skipped, full logs inspected. Both tested preview4e8f21f81155a55383272fe7ade895896fb72db6 with identical tree/orderedbase+head parents. Eight Inventory controls/restored baselines and all9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution controls passed each host;12 actual CLI/19 default/19 branded packages/docs/branding/provenance/module/inventory checks. Linux race/vet/fuzz passed; 82.7% coverage, zero reachable/imported findings; existing module-only advisory stays open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985428 /111340777428 passed full14 branches/indexed3-input4-output chunks/both retained exact bytes/hashes/source-tree-APRL provenance. Original scratch Go/test/script blobs/complete16 remote files/scope/no deleted or pinned files/commit identity/tree/parent/base/preview/rules verified. Semantic self-review checked actual source helper/mask/SKU/duplicate-skip plus explicit sheet/scope/budget/product corrections. Automated comments currently empty; no callable final manual Code Review tool or independent-person approval is claimed. No fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure. 
 Protected acceptance/ref/tree/parents/Denis author/GitHub committer verified. Distinct accepted-push proof follows below.
New inventory_capture_test.go.txt/extended capture-region-source.py invoke only unchanged pure inventory/merge/normalization/predicate functions in isolated pinned source/APRL.13 structural branches/30 literal normalization cases; actual output retention/topology/hash/final source/native/remote acceptance is pending. No target calculation code, pin/public/request/live/platform change. See REGION_INVENTORY_CALCULATIONS for scope/contracts/acceptance/rollback/next action.


PR110 distinct accepted-push verified: 
Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170292930 quality111341722998/Windows111341722896 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170292817 /111341722273 passed on exact2dcd8dae924ce3fad5402890cdeca0ecbef6e01e. All mandatory steps/full logs inspected, only conditional failure upload skipped. Both hosts passed8 Inventory/9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines;12 CLI/19 default/19 branded packages/docs/branding/provenance. Linux race/vet/fuzz82.7%; zero reachable/imported findings, existing module-only advisory open. Source14 full branches/complete indexed chunks/retained hashes/bytes/source-tree-APRL provenance exact. No fresh local Go/source/public/live/laptop/Azure/Gate004/release/advisory closure.



Observed source characterization preparation is accepted through PR111; full final-head and distinct accepted-push evidence is indexed in SESSION_HANDOVER and REGION_INVENTORY_CALCULATIONS. The original observed fixture hashes remain unchanged.

## PR111 accepted and owned aggregation active (2026-10-04)

Inventory source characterization PR111 is VERIFIED OFFLINE. The bounded owned aggregation helper is IN PROGRESS/unaccepted under [REGION_INVENTORY_AGGREGATION.md](REGION_INVENTORY_AGGREGATION.md). Public region-selection remains unavailable. Finish fresh exact-head native/source QA, review and protected acceptance, then availability/latency/cost/quota/reservation arithmetic, bounded request adapters and coordinator/public/all-format execution. Azure/laptop-dependent validation remains explicitly deferred.

Accepted PR111 final native37170957575/source37170957525 and distinct accepted-push native37171306079/source37171306081 passed all mandatory steps/full logs on the documented exact heads; unchanged four captures/13 inventory branches/30 normalization cases and all prior controls. See SESSION_HANDOVER/REGION_INVENTORY_CALCULATIONS and [PR111](https://github.com/DeBoX85/Cloud-Assess/pull/111) for commit/tree/ordered parents/job/coverage/security proof. Current aggregation contract separately defines corrections, independent budgets, actual complete five-map oracles, valid literal subset, ownership/cancellation and ten compiling controls; no new Go/native/public/live acceptance yet.
