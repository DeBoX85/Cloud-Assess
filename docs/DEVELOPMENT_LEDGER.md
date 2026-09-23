# Cloud Assess Development Ledger

Status: active development-process audit index

Purpose: provide a durable, repository-backed record of how Cloud Assess was designed, implemented, reviewed, corrected, and validated.

This ledger is intended for later troubleshooting, forensic review, handover, and development audit. It supplements rather than replaces Git history.

## Audit model

Cloud Assess development evidence is intentionally split into several layers:

| Evidence | Purpose | Authority |
|---|---|---|
| Git commit history | Exact chronological code/document changes | Primary per-change record |
| GitHub Actions runs | Executable validation of a specific commit SHA | Primary validation record |
| This development ledger | Chronological index connecting milestones, decisions, defects, commits, and validation | Primary process index |
| `IMPLEMENTATION_PLAN.md` | Planned phases, completion boundary, outstanding work | Current roadmap |
| `CHARACTERIZATION.md` | Source behavior preserved or intentionally changed | Behavioral contract |
| `QUALITY_GATE_001.md` | Formal audit of the core foundation | Audit snapshot |
| `QUALITY_GATE_002.md` | Formal audit of the runnable integration path | Audit snapshot |
| `QUALITY_GATE_003.md` | Formal post-live-validation repository and quality-gate audit | Audit snapshot |
| `EQUIVALENCE.md` | Procedure and normalization rules for source-vs-target comparison | Validation runbook |
| `ROADMAP.md` | Ordered remaining work, evidence gates, and release-scope decisions | Current execution roadmap |
| `TARGET_SPECIFICATION.md` | Target architecture/product requirements | Design contract |
| `NOTICE.md` / `THIRD_PARTY_LICENSES.md` | Attribution and incorporated-license evidence | Legal/provenance record |

If this ledger and Git disagree about the exact contents of a change, Git is authoritative. If documentation and the pinned reference disagree about source behavior, the pinned source implementation and characterization tests are authoritative.

## Reference baseline

Reference repository:

```text
DeBoX85/azqr
```

Pinned reference commit:

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

Pinned source tree:

```text
17d93b20c303f90f7843036be82f0dc32f3260f1
```

Pinned APRL revision:

```text
60eaddda76541f6adbc1c5ffa686829807e55e29
```

Pinned AOR imported tree:

```text
a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5
```

Pinned custom-rule imported tree:

```text
674b9b3dcb443ce6dc445b48e1db47e4a0ca7082
```

Pinned known-SKU blob:

```text
a2d97a60ec445ce4023a1a5a9b6c4dca76e31f5a
```

Active development branch:

```text
bootstrap/core-v1
```

## Standing development rules

These rules have governed the reconstruction effort:

1. Code/runtime behavior wins over stale or ambiguous documentation.
2. Source behavior is characterized before target behavior is declared equivalent.
3. Functional behavior is preserved unless an intentional change is explicitly recorded.
4. Lower-level packages return errors instead of terminating the process.
5. Target improvements must not be silently normalized as source-equivalent; they must be documented.
6. Exact source-data provenance is pinned and checked in CI.
7. Every milestone head must pass the permanent quality gate before it is treated as complete.
8. Partial or failed live assessments cannot be used as valid equivalence baselines.
9. New normalization rules in the equivalence harness require an explicit compatibility rationale.
10. Public/release readiness is not inferred from unit-test success alone; live Azure equivalence remains required.

## Development chronology

### Phase A: Repository and target architecture

**Objective**

Create an independent Cloud Assess repository and define the target architecture before reproducing assessment behavior.

**Key decisions**

- Working product name: Cloud Assess
- CLI name: `cloud-assess`
- Go retained as implementation language
- Apache-2.0 repository-level license
- Required MIT attribution retained for incorporated/derived material
- Canonical neutral domain model instead of source `GraphResult` coupling
- Excel remains the default human-facing report
- JSON becomes the canonical machine-facing result contract
- Core implementation may be restructured as long as behavior remains materially equivalent
- Legacy CLI/config compatibility is not required

**Primary records**

- `docs/TARGET_SPECIFICATION.md`
- `docs/IMPLEMENTATION_PLAN.md`
- `NOTICE.md`
- `THIRD_PARTY_LICENSES.md`

### Phase B: Foundation and characterization

**Objective**

Reproduce and characterize deterministic core behavior before adding the Azure auxiliary stages.

**Implemented/characterized**

- neutral branding/domain model
- filter semantics
- scanner/service selection
- stage configuration
- severity gate
- subscription-ID redaction
- previous-completed-month Cost period
- resource-ID helpers
- scanner registry
- APRL/AOR/custom recommendation corpus
- Azure cloud/environment selection
- DefaultAzureCredential construction
- authenticated retry/throttling HTTP stack
- Resource Graph transport, pagination, batching, row decoding and error handling
- resource inventory discovery
- subscription discovery abstractions
- management-group traversal
- findings deduplication and summary
- recommendation applicability

### Quality Gate 001

**Purpose**

Loop back over the foundation before beginning the auxiliary Azure stages.

**Validated code baseline**

```text
dc82bc81755acae7ad84759483b59509f0889e2d
```

**Workflow**

```text
35095316848
```

**Material defect found**

The initial target interpreted `include.resourceTypes` as literal ARM resource types. The reference uses scanner/service keys and then expands them into ARM resource-type scope. Downstream findings were also not reapplying the complete structural scope.

**Remediation**

- restored scanner-key semantics
- applied subscription/RG/resource-type/resource scope to inventory and downstream findings
- added regression tests
- corrected specification drift
- strengthened CI
- added complete third-party MIT license texts

**Full record**

`docs/QUALITY_GATE_001.md`

### Phase C: Azure assessment subsystems

**Objective**

Reproduce the non-core-Graph Azure datasets independently before orchestration.

**Diagnostics**

- direct ARM batch behavior
- 20-resource batch size
- 30-worker ceiling
- 41 dedicated recommendation definitions
- explicit Azure Resource Manager validation mechanism
- non-success subrequests surfaced as warnings
- lower-level failures returned instead of process termination

**Advisor**

- ARG recommendation instances
- Advisor metadata API for descriptions
- no additional Advisor SDK dependency; shared ARM HTTP layer used
- malformed rows surfaced as warnings
- deterministic output

**Defender**

- separate plan/tier status dataset
- separate unhealthy security recommendations dataset
- source-compatible deduplication
- preserved source ResourceType projection behavior

**Azure Policy**

- noncompliant policy-state dataset
- management-group-aware ARG scope
- source-compatible deduplication and metadata joins

**Arc SQL**

- pinned multi-join ARG query
- source license/DPS/telemetry derivation
- pinned `vcores` decoder mismatch retained as a live-equivalence review item

**Cost**

- previous completed UTC month
- ActualCost grouped by ServiceName
- source-compatible two-worker ceiling
- direct authenticated Cost Management REST call
- source skip/error codes preserved
- source bug corrected: subscription display name is populated

### Phase D: Stage health and canonical result

**Objective**

Make assessment health explicit and separate from process success.

**Implemented**

- per-stage status
- `complete`
- `complete_with_warnings`
- `partial`
- `failed`
- critical stage stop behavior
- optional-stage continuation
- stage parameter validation
- canonical immutable-style assessment result
- deterministic dataset ordering

**Exit semantics**

```text
0 = success
1 = execution/config/auth/render failure
2 = quality/severity gate failure
3 = partial assessment
```

### Phase E: Rendering

**Implemented**

- canonical JSON
- CSV
- Excel
- SARIF 2.1.0
- JSON stdout
- shared CSV/Excel table projection
- Assessment Status as first Excel worksheet
- source-compatible SKU capacity behavior using exact pinned known-SKU data
- centralized Cloud Assess branding
- private report file permissions

**Intentional renderer differences**

- canonical JSON is domain-oriented rather than source table-oriented
- SARIF uses Cloud Assess branding/fingerprint namespace
- Diagnostics reports the real ARM validation mechanism
- Assessment Status exposes stage health directly

### Phase F: Production orchestration and CLI

**Objective**

Turn the independently tested subsystems into a runnable internal product path.

**Implemented**

```text
cloud-assess scan
  -> validate configuration/stages/gate
  -> Azure credential
  -> scope discovery
  -> inventory discovery
  -> two-phase Graph catalog/pruning/execution
  -> enabled auxiliary stages
  -> canonical result
  -> report rendering
  -> completeness/severity exit semantics
```

**Key behavior**

- CLI subscription/RG scope is folded into include filters as in the reference
- Graph/resource inventory/scope are critical
- optional stage failures preserve later results
- reports are written before exit 2/3 when possible
- SIGINT/SIGTERM cancellation is propagated through command context
- default report timestamp is selected at scan start
- explicit plugin-stage execution fails early while production plugin execution is deferred

### Quality Gate 002

**Purpose**

Audit the recent integration work after an interrupted development session.

**Validated code/CI baseline**

```text
ec24ddb0eaab7d4f8902a4eb962e9f392499ee9f
```

**Workflow**

```text
35295328518
```

**Material findings**

1. JSON/stdout redaction was not honoring default subscription-ID redaction.
2. Deferred plugin stage could be requested misleadingly.
3. Default report timestamp was selected too late.
4. Cross-package integration coverage was insufficient.
5. CI did not explicitly build/invoke the final binary.
6. GitHub checkout action used a deprecated Node runtime.
7. Repository status documentation was stale.

**Remediation**

- JSON/stdout redaction fixed, including embedded/warning-only subscription IDs
- plugin stage now fails before Azure authentication while deferred
- timestamp moved to scan start
- authoritative coordinator -> app -> JSON/XLSX smoke path added
- CI builds the binary and runs root/scan/version smoke checks
- checkout action updated
- status docs brought current

**Full record**

`docs/QUALITY_GATE_002.md`

### Phase G: Semantic equivalence harness

**Objective**

Create a deterministic comparison boundary before live Azure validation.

**Current validated head**

```text
f357a8b4e149b46540b1687779899abc1cc76138
```

**Workflow**

```text
35297058652
```

**Result**

PASS: provenance, formatting, module graph, branding boundary, executable build/smoke checks, full race-enabled tests, and `go vet`.

**Implemented**

- development-only `tools/equivalence`
- reference table-JSON projection
- Cloud Assess canonical JSON projection
- semantic comparison by stable identities
- dataset coverage mismatch detection
- missing/extra/changed record details
- per-dataset reference/target/missing/extra/changed counts
- target `partial`/`failed` rejection
- redacted-input rejection
- documented intentional normalization only

**Intentional normalizations**

- reference `AZQR` -> `DIAGNOSTICS` for known Diagnostics recommendation IDs
- reference `AZQR` -> `CUSTOM` for legacy embedded custom-rule recommendations
- APRL/AOR provenance remains strict
- Cost subscription display-name source bug ignored
- Azure Policy timestamp ignored
- target-only timing/schema metadata ignored
- numeric formatting-only differences normalized

**Runbook**

`docs/EQUIVALENCE.md`

### Phase H: Reproducible live-equivalence runner

**Objective**

Turn the documented live-comparison procedure into an auditable, repeatable Windows workflow before touching the Azure test environment.

**Validated baseline**

```text
b7eca668b67a3a3515ed4d2fe6bc6e3a8469e706
```

**Workflow**

```text
35359172689
```

**Result**

PASS: pinned-data provenance, Go formatting/module consistency, branding boundaries, PowerShell parsing, actual CLI build/help/version smoke tests, full race-enabled tests, and `go vet`.

**Implemented**

- `scripts/live-equivalence.ps1`
- subscription, resource-group, and management-group scopes
- strict pinned-reference commit verification
- clean-working-tree requirement for both source and target
- exact APRL submodule verification for both repositories
- automatic unredacted AZQR reference run
- automatic unredacted Cloud Assess target run
- automatic semantic comparator execution
- timestamped local evidence bundles under `artifacts/equivalence/`
- Git ignore protection for sensitive live evidence
- per-command stdout/stderr capture
- exact commit/branch/scope/stage/filter/environment metadata
- SHA-256 hashes for separate source/target filter files
- explicit warning that the evidence bundle contains sensitive unredacted Azure identifiers
- CI parsing of PowerShell validation scripts

**Operational decision**

The live runner does not clone, checkout, or mutate the source repositories automatically. It fails closed when the reference is not at the pinned commit, a working tree is dirty, or a required submodule is not initialized at the pinned revision. This preserves forensic reproducibility and avoids silently changing the developer's environment.

**Runbook**

`docs/EQUIVALENCE.md`

### Phase I: First live Azure equivalence pass

**Date**

```text
2026-09-22
```

**Reference commit**

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

**Target assessment commit**

```text
27346cab18873402ab6f74468924ca3f72ce6539
```

**Scope**

One non-production Azure subscription, default stage set, no filter files.

**Execution health**

- pinned reference scan exit code: 0
- Cloud Assess scan exit code: 0
- semantic comparator exit code: 1 on the initial comparison

**Observed semantic counts**

- recommendations: 314 reference / 314 target, no missing/extra/changed records
- primary findings: 71 reference / 71 target, no missing/extra records, 11 changed records
- resource types: 7 / 7, exact
- in-scope inventory: 21 / 21, exact
- out-of-scope inventory: 3 / 3, exact
- Advisor: 12 / 12, exact
- Defender status: enabled on both sides, 0 / 0
- Policy, Arc SQL, Defender recommendations, and Cost were not enabled in this default-stage pass

**Only observed delta**

All 11 changed primary findings differed only in `subscriptionName`.

Pinned-source inspection confirmed that these records were Diagnostics findings. The AZQR Diagnostics scanner populates `SubscriptionID` but does not assign `GraphResult.SubscriptionName`. Cloud Assess intentionally fills the discovered subscription display name.

**Classification**

```text
Pinned-source data omission corrected by target
```

This is not treated as an assessment-semantic defect.

**Remediation**

- equivalence normalization narrowed to known Diagnostics recommendation IDs only
- ordinary Graph findings continue to compare subscription display name strictly
- regression fixture changed to reproduce the actual live AZQR omission
- dedicated test added to prove the normalization does not apply to non-Diagnostics findings

**Validation state**

The original Azure scans do not need to be repeated for this correction because the raw reference and target JSON artifacts already exist. Re-running the comparator against those same artifacts after pulling the normalization fix is sufficient to validate the reclassification.

### Phase J: First successful live Azure equivalence baseline

**Date**

```text
2026-09-22
```

**Evidence**

The original live Azure reference and target reports from Phase I were re-compared after adding the narrowly scoped Diagnostics subscription-name normalization.

**Result**

```text
equivalent: true
```

**Semantic coverage**

- recommendations: 314 reference / 314 target, 0 missing, 0 extra, 0 changed
- primary findings: 71 / 71, 0 missing, 0 extra, 0 changed
- resource types: 7 / 7, exact
- in-scope inventory: 21 / 21, exact
- out-of-scope inventory: 3 / 3, exact
- Advisor: 12 / 12, exact
- Defender plan status: enabled on both sides, 0 / 0
- Azure Policy: not enabled in this pass
- Arc SQL: not enabled in this pass
- Defender recommendations: not enabled in this pass
- Cost: not enabled in this pass

**Interpretation**

For the default-stage behavior actually exercised by this subscription, Cloud Assess is semantically equivalent to the pinned AZQR reference after accounting for the confirmed pinned-source Diagnostics subscription-name omission.

This milestone does not yet establish equivalence for optional stages that were not enabled or for Azure behaviors absent from the selected subscription.

**Next validation boundary**

Expand live coverage deliberately rather than repeating the same baseline:

1. enable optional stages that can return data in the existing test subscription;
2. validate filtered/resource-group scope;
3. validate multi-subscription or management-group scope;
4. add targeted fixtures only for important behaviors absent from the existing environment.

### Phase K: Optional-stage live Azure equivalence pass

**Date**

```text
2026-09-22
```

**Reference commit**

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

**Target commit**

```text
498511488401247a7c8cecaba6438580ba6fc09a
```

**Scope**

Same non-production Azure subscription used for the first baseline, no filter files.

**Additional stages enabled**

```text
policy
defender-recommendations
cost
```

The default graph/diagnostics/advisor/defender stages remained enabled on both source and target.

**Execution health**

- pinned reference scan exit code: 0
- Cloud Assess scan exit code: 0
- semantic comparator exit code: 0
- semantic result: equivalent = true

**Semantic coverage**

- recommendations: 314 / 314, exact
- primary findings: 71 / 71, exact
- resource types: 7 / 7, exact
- in-scope inventory: 21 / 21, exact
- out-of-scope inventory: 3 / 3, exact
- Advisor: 12 / 12, exact
- Azure Policy: enabled on both sides, 0 / 0
- Defender recommendations: enabled on both sides, 0 / 0
- Defender plan status: enabled on both sides, 0 / 0
- Cost: 5 / 5, exact
- Arc SQL: not enabled in this pass

**Interpretation**

The Cost implementation now has live semantic-equivalence evidence with non-empty data. Policy and Defender Recommendations have live execution/coverage equivalence evidence for an empty-result environment, but still require a future environment or targeted fixture containing non-empty results before their row-level projections can be considered live-validated.

**Next validation boundary**

Validate scope semantics next, beginning with resource-group scope on the same subscription. After that, validate multi-subscription or management-group traversal if an appropriate test scope is available.

### Phase L: Resource-group-scoped live Azure equivalence pass

**Date**

```text
2026-09-22
```

**Reference commit**

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

**Target commit**

```text
3f869e7cdce26b6bf11fa4306337efd45018fa3a
```

**Pinned APRL commit**

```text
60eaddda76541f6adbc1c5ffa686829807e55e29
```

**Scope**

The same non-production Azure subscription used for the preceding live passes, restricted to resource group:

```text
haz-rg-aitestbed-wus
```

No filter files were used.

**Stages**

The default graph, diagnostics, Advisor, and Defender plan-status stages were enabled on both source and target. Policy, Arc SQL, Defender Recommendations, and Cost were not enabled.

**Execution health**

- pinned reference scan exit code: 0
- Cloud Assess scan exit code: 0
- semantic comparator exit code: 0
- semantic result: equivalent = true

**Semantic coverage**

- recommendations: 314 / 314, exact
- primary findings: 47 / 47, exact
- resource types: 7 / 7, exact
- in-scope inventory: 11 / 11, exact
- out-of-scope inventory: 13 / 13, exact
- Advisor: 6 / 6, exact
- Defender plan status: enabled on both sides, 0 / 0
- Azure Policy: not enabled in this pass
- Arc SQL: not enabled in this pass
- Defender recommendations: not enabled in this pass
- Cost: not enabled in this pass

**Interpretation**

Cloud Assess is semantically equivalent to the pinned AZQR reference for the resource-group-scoped behavior exercised by this environment. The non-empty findings, inventory, out-of-scope, and Advisor datasets provide live evidence that the resource-group boundary and downstream scope filtering behave equivalently.

The empty Defender plan-status result establishes execution and empty-result equivalence only. It does not validate non-empty Defender row projection.

**Evidence note**

The run metadata captured the full source and target commands, scope, commits, dataset enablement, and exit codes. Its `stages` property was null because the runner used the implicit default stage set. The executed commands and dataset enablement make the effective coverage reconstructable, but the runner should record the resolved default stages explicitly in future evidence.

**Next validation boundary**

Validate broader scope traversal next, using multi-subscription or management-group scope if an appropriate test scope and permissions are available. Important behaviors still absent from live non-empty evidence should subsequently be covered through a suitable Azure environment or targeted fixtures.

### Phase M: Post-live-validation repository audit and quality-gate hardening

**Date**

```text
2026-09-23
```

**Audited baseline**

```text
5cc911ea853cbcdf402a43cb0c96d359f8ece81d
```

**Validated remediation commit**

```text
954dd75584f96dff1f5017b952c0fd11def40a53
```

**Workflow evidence**

```text
Go quality gate run 35847062551
```

**Reason for the gate**

Development paused before broader live scope validation to perform a complete repository QA and check for implementation, documentation, evidence, CI, or security drift.

**Material findings**

- Go `1.26.0` produced 19 reachable standard-library findings under `govulncheck`
- implicit live-runner stages executed correctly but were serialized as null in evidence metadata
- the live runner could not express multiple explicit subscriptions even though both executables support them
- PowerShell evidence logic was parser-checked but not behavior-tested
- README, implementation-plan, characterization, and current-boundary text lagged behind the completed Phase J/K/L live passes
- pull requests targeting the active development branch were not quality-gated
- external GitHub Actions used mutable major-version tags
- the permanent gate had no coverage floor or reachable-vulnerability scan

**Remediation**

- raised the repository and CI toolchain to Go `1.26.8`
- added pinned `govulncheck v1.8.0`; remediated result is zero reachable vulnerabilities
- added a 75% total statement-coverage floor; validated total is 77.2%
- pinned external actions to full commit SHAs
- added `bootstrap/core-v1` as a pull-request gate target
- added tested scope/stage helper logic for the live runner
- runner metadata schema `1.1` records requested controls, effective stages, and implicit-default selection
- runner now supports multiple subscriptions, multiple management groups, and multiple resource groups within one subscription
- reconciled current-state and next-boundary documentation
- created `QUALITY_GATE_003.md` as the detailed audit record

**Validation state**

- independent repeated shuffled tests: pass
- PowerShell parse and behavior tests: pass
- executable build/help/version checks: pass
- full race-enabled tests: pass
- coverage floor: pass at 77.2%
- `go vet`: pass
- `actionlint`: pass
- `govulncheck`: pass with zero reachable vulnerabilities
- hosted quality gate: pass

**Residual boundary**

The Phase J/K/L raw evidence bundles remain intentionally untracked because they contain unredacted Azure identifiers. This repository gate validates the runner, metadata contract, code, documentation, and hosted workflow; it does not replace independent replay of those sensitive bundles.

Broader live equivalence remains next, beginning with multi-subscription or management-group scope. Arc SQL and non-empty Policy/Defender datasets still require suitable live data or targeted fixtures.

### Phase N: Two-subscription live Azure equivalence pass

**Date**

```text
2026-09-23
```

**Reference commit**

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

**Target commit**

```text
612fc765fe74677fc05cd4f5e5574595030c4c4f
```

**Pinned APRL commit**

```text
60eaddda76541f6adbc1c5ffa686829807e55e29
```

**Scope and evidence**

Two distinct, accessible non-production Azure subscriptions were passed explicitly. No resource-group scope or filter files were used. The effective stages were the implicit defaults: graph, diagnostics, Advisor, and Defender plan status. Policy, Defender Recommendations, Arc SQL, Cost, and plugins were not enabled.

The local, untracked evidence bundle is identified by UTC run stamp `20260923_112151Z`. The reviewer examined its `run-metadata.json` and `equivalence.json`, plus stage health and per-subscription counts extracted locally from the unredacted target JSON. The unredacted source and target reports and execution logs were not transferred to this repository or independently replayed by the reviewer.

**Execution health**

- pinned reference scan exit code: 0
- Cloud Assess scan exit code: 0
- semantic comparator exit code: 0
- semantic result: `equivalent = true`
- target assessment completeness: `complete_with_warnings`
- scope discovery: completed, 2 subscriptions resolved

**Semantic coverage**

- recommendations: 314 / 314, exact
- primary findings: 98 / 98, exact
- resource types: 23 / 23, exact
- in-scope inventory: 38 / 38, exact
- out-of-scope inventory: 21 / 21, exact
- Advisor: 46 / 46, exact
- Defender plan status: 18 / 18, exact and non-empty
- Policy, Defender Recommendations, Arc SQL, and Cost: not enabled

The user's local per-subscription target summary showed the following, with subscription IDs omitted from this ledger:

| Selected subscription | In-scope inventory | Out-of-scope inventory | Raw target findings | Advisor | Defender plan status |
|---|---:|---:|---:|---:|---:|
| 1 | 12 | 10 | 50 | 11 | 0 |
| 2 | 26 | 11 | 74 | 35 | 18 |

Both subscriptions contributed non-empty inventory, findings, and Advisor records. Raw target findings include SLA-category findings; the semantic primary-findings dataset excludes those and compares the associated SLA values through inventory instead. Thus the raw per-subscription finding counts must not be added and compared directly with the 98 primary-finding records.

**Warning classification and interpretation**

The target Diagnostics stage completed with two `diagnostics_subrequest_non_success` warnings. Both were HTTP 400 responses from ARM batch subrequests for diagnostic settings. For a non-success subrequest, the target preserves the pinned reference's missing-diagnostics finding behavior while explicitly reporting that the diagnostic-setting state of the affected resource is uncertain. The available aggregate evidence does not identify the two affected resources or establish why Azure returned HTTP 400. These warnings are an unresolved coverage limitation, not an observed source-versus-target delta.

Cloud Assess and the pinned reference are semantically equivalent for the compared two-subscription datasets in this run, and the non-empty Defender plan-status result supplies live row-level equivalence evidence. This does not prove the correctness of diagnostic-setting conclusions for the two unsuccessful subrequests, and it does not establish management-group traversal or the optional stages absent from this pass.

**Next validation boundary**

Investigate the two Diagnostics HTTP 400 subrequests using the locally retained evidence and read-only Azure requests if the affected resources can be identified. Then validate management-group traversal when an appropriate scope is available. Obtain non-empty Policy and Defender Recommendations evidence and resolve the Arc SQL response shape separately.

### Planning checkpoint: Remaining-work roadmap

**Date**

```text
2026-09-23
```

The [execution roadmap](ROADMAP.md) now orders the remaining live validation, missing stage evidence, first-release feature decisions, packaging, and final security/operational review. It records explicit completion evidence and required user inputs, including separate authorization before Azure fixture provisioning. This planning checkpoint does not claim that any of those open items are complete or expand the agreed core-v1 scope.

### Phase O: Diagnostics warning investigation with individual GETs

**Date**

```text
2026-09-23
```

**Evidence and method**

The user ran the read-only `tools/diagnostics-probe` from target commit `f8fc77900ba0803232dbaba4e14053edca2454c0` against the locally retained, unredacted Phase N target report (`20260923_112151Z`) in an authenticated Azure environment. The probe made individual diagnostic-settings GET requests using the scanner's ARM endpoint, supported-resource list, and API version. The user supplied its sanitized console summary; the raw report and resource IDs remain outside Git. The identity used by the probe was not independently verified against the original scan metadata.

- original target report: 2 Diagnostics batch subrequest warnings with HTTP 400
- probe: 37 eligible resources, 35 successful GETs, 2 failed GETs
- both failed GETs: `microsoft.network/networkwatchers`, HTTP 400, Azure error code `ResourceTypeNotSupported`

**Classification and impact**

The pinned AZQR reference (`8e4f0577f3615e6c9014c031bcad079f235369cc`) and Cloud Assess both include Network Watchers in the Diagnostics request list. Both have a `nil` recommendation entry for that resource type. Consequently, these two failed requests cannot themselves produce dedicated missing-diagnostics findings in either implementation. This is a pinned-source eligibility mismatch with the API behavior observed in this Azure environment, not evidence of a target-only request-construction defect. The target preserved reference output and reported the non-success status as warnings; do not suppress these warnings to make the stage appear complete.

The probe ran after the original scan and issued individual GETs, while Phase N used ARM batch. The matching number and status of failures strongly suggest the Network Watchers caused the two batch warnings, but the original batch responses did not retain a resource correlation or error code. Their exact identity and the absence of other transient batch failures remain unproved. The original target report remains `complete_with_warnings`; this investigation does not turn it into a fully successful Diagnostics pass. A batch response capture with request correlation would be needed to confirm the historical mapping if required for release evidence.

**Next validation boundary**

Proceed with management-group traversal on a suitable read-only test scope. Keep the Phase N Diagnostics caveat in the evidence matrix, and correlate batch failures with request IDs in a future controlled run if that limitation needs to be closed.

### Phase P: AdvisoryDev leaf management-group live equivalence pass

**Date and provenance**

```text
2026-09-23
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: bc69fc004f67df1094305fb3f91ba41d625152cf
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260923_145252Z
```

**Scope and evidence**

The requested management-group ID was `AdvisoryDev`, the accessible non-production leaf under `Advisory`. The parent `Advisory` also contains a production child, so it was not scanned in this pass. The user supplied the generated `run-metadata.json`, `equivalence.json`, and a local summary of the unredacted target report. Raw source/target JSON and logs remain local and were not independently inspected by the reviewer. No filter files were supplied; implicit default stages were graph, diagnostics, Advisor, and Defender plan status on both sides.

The target scope stage completed with one resolved subscription. The user's local check found one distinct subscription ID across inventory, out-of-scope inventory, findings, Advisor, and Defender records, and confirmed it was the expected Dev subscription. This pass exercises leaf management-group subscription discovery and subsequent assessment. It does not exercise recursion from a parent management group into child groups or assess the sibling production subscription.

**Execution health and semantic coverage**

- pinned reference scan exit code: 0
- Cloud Assess scan exit code: 0
- comparator exit code: 0; `equivalent = true`
- target completeness: `complete_with_warnings`
- recommendations: 314 / 314, exact
- primary findings: 44 / 44, exact
- resource types: 11 / 11, exact
- in-scope inventory: 12 / 12, exact
- out-of-scope inventory: 10 / 10, exact
- Advisor: 11 / 11, exact
- Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL, and Cost: not enabled

The local target stage summary reported inventory 12, graph 44, diagnostics 6, Advisor 11, and Defender plan status 0. Scope, inventory, graph, Advisor, and Defender stages completed; Diagnostics completed with warnings. It contained one `diagnostics_subrequest_non_success` warning. A follow-up local read of that warning showed `diagnostic settings batch subrequest returned HTTP 400`. The batch warning does not contain the affected resource or Azure error code; the comparator independently noted target warnings. It cannot be assumed to be the Network Watcher response observed by the later individual GET probe in Phase O without request correlation. The target report remains `complete_with_warnings`, so diagnostic-setting conclusions for the unsuccessful batch subrequest remain uncertain.

**Read-only probe follow-up**

The user ran `tools/diagnostics-probe` against the retained Phase P target report. It selected 11 supported inventory resources, returned 10 successful individual diagnostic-settings GETs and one HTTP 400 `ResourceTypeNotSupported` response for `microsoft.network/networkwatchers`. The resource-ID hash for that failure exactly matched one of the two failed Network Watchers in the Phase O probe of the broader two-subscription report. This confirms that the same Network Watcher produces the same error on subsequent individual GETs and strengthens the explanation for the Phase N and P warnings.

Neither original ARM batch response retained request-to-subresponse correlation, so the matching count and status still do not prove which resource generated each historical warning. In the pinned source and target, the Network Watcher Diagnostics support entry has no dedicated recommendation, so this observed unsupported-resource response does not itself generate a missing-diagnostics finding. Preserve both historical reports as `complete_with_warnings`; an exact batch mapping remains open if required for release evidence.

**Next validation boundary**

Keep parent-to-child management-group traversal open until a suitable non-production parent with accessible descendants is available or a separate scope decision is made. Continue with representative filter combinations and missing non-empty optional-stage evidence. Preserve the raw reports and logs under the local evidence stamp for later warning classification.

### Planning checkpoint: Paired Storage/VM filter fixture

**Date**

```text
2026-09-23
```

The Dev inventory contains two Storage Accounts and one Virtual Machine. A paired, subscription-agnostic fixture now selects the `st` and `vm` scanner keys under the pinned source's `azqr:` root and the target's `assessment:` root. The [equivalence guide](EQUIVALENCE.md#filter-file-note) records the exact read-only runner command on `AdvisoryDev`. This fixture has been syntax-checked locally; a filtered live pass and its warning/scope review remain pending. It does not establish tag, resource-group, resource, or recommendation filter equivalence.

## Current boundary

The generic core scan, deterministic equivalence tooling, reproducible live runner, default/optional/resource-group/two-subscription/leaf-management-group live passes, and post-live repository remediation are complete and quality-gated for the behavior exercised so far. The two-subscription and leaf-management-group passes have explicit Diagnostics warning boundaries.

The next live validation boundary is:

```text
Nested management-group traversal and representative filtered-scope equivalence; retain Diagnostics warning caveats
```

Required evidence set for each live pass:

- exact pinned reference commit
- exact Cloud Assess commit
- Azure scope/stage/filter parameters
- unredacted reference JSON
- unredacted Cloud Assess JSON
- generated equivalence JSON
- classification notes for every non-empty delta

## Known open items

The following are not forgotten; they remain intentionally open:

- Diagnostics HTTP 400 subrequest root cause and affected-resource coverage
- nested management-group traversal equivalence
- non-empty Policy and Defender Recommendations evidence
- Arc SQL numeric `vcores` response-shape resolution
- production external/YAML plugin execution
- internal plugin migration/parity
- scanner-specific CLI commands
- `rules` CLI command
- `plugins list/info` CLI
- packaging/distribution
- generated dependency/license inventory
- final release-level security and operational review
- deferred compare/alternative-VM-SKU/MCP/public website/distribution work

## Troubleshooting and forensic reconstruction

For a later investigation, use this order:

1. Identify the affected behavior/milestone in this ledger.
2. Open the linked quality gate or characterization section.
3. Find the relevant target commit SHA.
4. Inspect the exact change with:
   ```bash
   git show <commit-sha>
   ```
5. Inspect chronological history around it:
   ```bash
   git log --date=iso --decorate --oneline bootstrap/core-v1
   ```
6. Re-run or inspect the associated GitHub Actions workflow using its recorded run ID.
7. Compare the behavior to the pinned source commit.
8. Run the relevant characterization/unit test.
9. For live mismatches, reproduce both unredacted JSON reports and run `tools/equivalence`.
10. Classify the delta before changing code.

## Maintenance rule for this ledger

From this point onward, update this ledger whenever any of the following occurs:

- a milestone becomes complete
- a material defect is found
- a behavior is intentionally changed from the source
- a source defect is deliberately corrected
- a new quality gate is performed
- the pinned reference/provenance changes
- live equivalence identifies a new class of delta
- a deferred feature enters or leaves scope
- a release baseline is created

Routine commits do not need individual prose entries because the Git log already records them. Significant commits and every validated milestone must be indexed here.

This makes the development process reconstructable even if the conversational history that produced it is unavailable.
