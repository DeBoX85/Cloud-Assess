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
| `FAILURE_NOTES.md` | Confirmed mistakes, recurrence checks, corrections and prevention steps | Troubleshooting index; underlying Git/CI/ledger evidence remains authoritative |
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

### Phase Q: AdvisoryDev Storage/VM filtered semantic comparison

**Date and provenance**

```text
2026-09-23
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 229555a6ecea67642909f002ec0ec2e986743256
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260923_153403Z
```

The user supplied the generated `run-metadata.json` and `equivalence.json`. The runner requested management group `AdvisoryDev`, with graph, diagnostics, Advisor, and Defender stages. The reference used `examples/filters/azqr-storage-vm.yml` (SHA-256 `d18089f17792644ccb51828c55ea6272362701139c193023a4af44c06daf5422`); the target used `examples/filters/cloud-assess-storage-vm.yml` (SHA-256 `231db54615fea8a91265387303f66f6eecedfcf42b8bc1b8c6e983458be6bf51`). Each command in the metadata includes its respective filter path. The raw source and target reports and logs remain in the user's local evidence directory and were not independently inspected.

Both filter hashes match the checked-in fixture contents with CRLF checkout line endings on Windows; Git's LF blob hashes naturally differ. This is a line-ending difference, not evidence of changed filter selection.

**Execution and semantic result**

- reference, target, and comparator exit codes: 0 / 0 / 0
- comparator result: `equivalent = true`; no missing, extra, or changed records in any enabled dataset
- recommendations: 61 / 61
- primary findings: 18 / 18
- resource types: 3 / 3
- in-scope inventory: 4 / 4
- out-of-scope inventory: 18 / 18
- Advisor: 9 / 9
- Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL, and Cost: not enabled

**Target stage and scope follow-up**

The user read the retained unredacted target report locally and supplied its summary. Target completeness was `complete`. Its scope stage completed with one subscription; across inventory, out-of-scope inventory, findings, Advisor, and Defender records, the only populated subscription ID was the expected non-production Dev subscription `c09f96df-19de-4c49-80ff-0c1a94d93ab2`. Scope, inventory (4 records), graph (19), diagnostics (2), Advisor (9), and Defender (0) stages all completed with no warning codes. Defender Recommendations, Policy, Arc SQL, Cost, and plugins were skipped. The graph and diagnostics stage record counts are stage counters and should not be equated with the comparator's primary-finding count of 18.

This completes target stage-health and scope review for this filtered pass. The supplied comparison establishes exact semantic agreement for the selected scanner keys and requested scope, with non-empty findings, inventory, out-of-scope inventory, and Advisor data. The Network Watcher HTTP 400 batch warning seen in the earlier unfiltered Dev run is absent from this filtered target report; this does not resolve the historical warning or identify its exact batch request. Raw reports and logs were retained locally and were not independently inspected by the reviewer. The unfiltered Phase N and P Diagnostics warning boundaries remain open, as do tag, resource-group, individual-resource, recommendation, and nested management-group filter/traversal checks.

### Planning checkpoint: Separate Dev resource-group and tag filters

**Date**

```text
2026-09-23
```

The user identified `rg-fos-FinopsHub-dev-euw` and `Environment: dev`. Their local read of the unfiltered Phase P target report counted three in-scope resources in that group, all three with the tag on the resource records. This supports a non-empty RG include pass, but a combined RG-and-tag pass cannot isolate the tag condition because it retains the same three resources as an RG-only pass. The resource types for these three records were not established by the supplied output.

Paired AZQR and Cloud Assess fixtures now express the full Dev resource-group ID for an RG-only pass, and another pair expresses a tag-only include. The [equivalence guide](EQUIVALENCE.md#filter-file-note) records the RG runner command and the condition for using the tag-only pair. Before running the tag pass, check whether the unfiltered Dev inventory contains any nonmatching resources; if none exist, a passing comparison would not prove tag exclusion. Neither new fixture pair has live equivalence evidence yet, and the earlier unfiltered Diagnostics warning caveats remain open.

### Phase R: AdvisoryDev resource-group filter semantic comparison

**Date and provenance**

```text
2026-09-23
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 1cd44a4904dd1ea1272c2eff89782c4ab3e0f63d
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260923_160228Z
```

The user supplied `run-metadata.json` and `equivalence.json` for the separate RG-only fixture pair on `AdvisoryDev`. Both commands requested the same management group and implicit default graph, diagnostics, Advisor, and Defender stages. Neither command used a CLI resource-group scope flag. The reference filter was `examples/filters/azqr-dev-rg.yml` (SHA-256 `f287e8c45e9febf1216d6d688242765fa545856f7049ed7337343cc0a635c022`); the target filter was `examples/filters/cloud-assess-dev-rg.yml` (SHA-256 `09d6664d56a3c91a01d060e2ca162dd8e6ae6a46fa99917af9dce5551f7f09e5`). These hashes match the checked-in fixture content with Windows CRLF line endings. The raw reports and logs remain in the user's local evidence directory and were not independently inspected.

**Execution and semantic result**

- reference, target, and comparator exit codes: 0 / 0 / 0
- comparator result: `equivalent = true`; zero missing, extra, or changed records in every enabled dataset
- recommendations: 314 / 314
- primary findings: 11 / 11
- resource types: 3 / 3
- in-scope inventory: 3 / 3
- out-of-scope inventory: 19 / 19
- Advisor: 3 / 3
- Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL, and Cost: not enabled

**Target stage and scope follow-up**

The user read the retained target report locally and reported `complete` status. The scope stage completed with one subscription, and the only populated subscription ID across inventory, out-of-scope, findings, Advisor, and Defender data was the expected Dev subscription `c09f96df-19de-4c49-80ff-0c1a94d93ab2`. Scope, inventory (3 records), graph (11), diagnostics (3), Advisor (3), and Defender (0) stages completed with no warning codes. Defender Recommendations, Policy, Arc SQL, Cost, and plugins were skipped. This closes target stage-health and scope review for the RG-only pass; its filtered result does not resolve the earlier unfiltered Diagnostics warnings. The raw reports and logs remain local and were not independently inspected by the reviewer.

The same local check of the unfiltered Phase P target inventory found 12 in-scope resources: 3 with resource tag `Environment: dev` and 9 without that exact tag match. The three matching resources are the same three already identified in the selected RG. The separate tag-only fixture can therefore test exclusion of nonmatching resources without an RG filter; it has not yet been run. Even if a tag-only report matches the RG-only report, each independently applied filter must be verified against the unfiltered baseline and pinned reference.

### Planning checkpoint: Enforced CI and Quality Gate 004

**Date**

```text
2026-09-23
```

A review of Gates 001-003 and the current workflow found that their PASS records are correctly scoped to their recorded commits but are not release approval. The latest reviewed merged-branch Go quality run was green, with 76.3% aggregate statement coverage against a 75% floor. Lower-coverage risk areas included throttling (33.3%), CLI (53.0%), discovery (59.9%), canonical result assembly (62.8%), and orchestration (69.5%). The `bootstrap/core-v1` branch API reported `protected: false`, so passing CI was not yet a merge requirement.

The maintenance workflows have been changed to publish generated commits on dedicated branches with PR compare links, rather than pushing directly to `bootstrap/core-v1`. The Go quality workflow has gained a separate Windows PowerShell 5.1 and Windows Go executable validation job. The [Quality Gate 004 plan](QUALITY_GATE_004_PLAN.md) defines core evidence rows and acceptance rules. These workflow changes and the branch protection setting must be verified before calling CI enforced; Gate 004 remains planned, not passed. The tag-only live pass and other roadmap boundaries remain open.

The first PR run of the Windows job passed PowerShell 5.1 parsing/helpers but exposed three tests asserting Unix `0600` file-mode bits on Windows. The Windows Go runtime reported `0666` for created CSV, XLSX and SARIF files. Go's Windows `FileMode`/`Chmod` API exposes the read-only attribute rather than the Windows access-control list, so the Unix-mode assertions were scoped to non-Windows platforms. This does **not** validate Windows report ACL privacy; Gate 004 now requires an explicit Windows ACL review before release approval. The changed PR must pass both CI jobs before merging.

### Planning checkpoint: Enforced branch checks verified

**Date**

```text
2026-09-23
```

PR #15 merged the maintenance-branch, Windows validation and Gate 004 planning changes into `bootstrap/core-v1` at `b75ca9a0397538e0f44376a225ba9c26ef66ad93`. Both required job contexts, `quality` and `windows-validation`, passed on the PR and the postmerge push run `35890146814`.

The repository's active branch ruleset `23890737` targets exactly `refs/heads/bootstrap/core-v1`, requires a pull request and both GitHub Actions job checks, prohibits deletion and non-fast-forward updates, and lists no bypass actors. The branch API reports `protected: true`. Its strict up-to-date requirement is disabled; this meets the agreed PR/checks enforcement but leaves a separate policy choice if concurrent merges become a concern. This checkpoint verifies enforced *development* CI; it is not a Gate 004 PASS, and Windows report ACL privacy is still unverified. The revised maintenance workflows have not yet been demonstrated with a postchange manual dispatch.

### Phase S: AdvisoryDev tag-only include-filter semantic comparison

**Date and provenance**

```text
2026-09-23
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 1cd44a4904dd1ea1272c2eff89782c4ab3e0f63d
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260923_170735Z
```

The user supplied `run-metadata.json`, `equivalence.json`, and a local summary computed from the unredacted target report and the Phase P unfiltered target report. The raw reference/target JSON and logs remain local and were not independently inspected. The runner used the `AdvisoryDev` management-group leaf without a CLI resource-group flag, with implicit graph, diagnostics, Advisor and Defender stages on both sides. The reference filter `examples/filters/azqr-environment-dev.yml` had SHA-256 `2e02f685c73c8810aa992045e81e5deb0e789b824ed1f6f51be94b33dce044fb`; the target filter `examples/filters/cloud-assess-environment-dev.yml` had SHA-256 `b091145d2e5c48189d3b3ff3a5fa091d39e420ea545f9edde99ba0793f12eb8f`. Both hashes match the checked-in YAML with Windows CRLF checkout line endings. The pair includes resource tag `Environment: dev` without an RG include filter.

**Execution and semantic result**

- reference, target and comparator exit codes: 0 / 0 / 0
- comparator: `equivalent = true`, no missing, extra or changed records in any enabled dataset
- recommendations: 314 / 314
- primary findings: 11 / 11
- resource types: 3 / 3
- in-scope inventory: 3 / 3
- out-of-scope inventory: 19 / 19
- Advisor: 2 / 2
- Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL and Cost: not enabled

The user's target report summary showed `complete` with one resolved scope subscription, 3 inventory, 11 graph, 3 diagnostics, 2 Advisor and 0 Defender stage records. Every requested stage completed with no warning codes. The local resource-ID set comparison against the earlier unfiltered Dev inventory returned `SelectedIdsMatchTags = True`: the three selected resources exactly matched the three with resource tag `Environment: dev`, none of the other nine baseline in-scope resources was selected, and all three selected records retained the tag. The out-of-scope count was 19. This provides non-empty, independently selective evidence for the include-tag condition, subject to the locally reported ID-set check.

**Cross-run caveat**

The RG-only Phase R pass counted three selected resources in the named group and had 3 / 3 Advisor rows, while the later tag-only run has 2 / 2. The earlier unfiltered baseline counted exactly three tagged resources in that group, and Phase S's IDs match that baseline. The two filtered runs' resource-ID sets have not yet been compared directly. The paired comparison within each run is exact; that does not establish why one Advisor row is absent across runs. The missing row's identity, whether it still targeted a selected resource, and whether Azure state changed or downstream filter behavior differed need local record-level classification. Do not silently label the difference as environmental or treat all tag-related downstream behavior as closed until this check is done.

The target SHA in this run predates the current branch merge `b75ca9a`. The intervening changed files are documentation, CI workflows, renderer tests and a SARIF comment; no Azure scan/tag-filter implementation or tag fixture content changed. The pass therefore characterizes the tag behavior exercised by the recorded target SHA; future runs should first fast-forward the local branch for exact current-branch provenance.

**Next validation boundary**

Classify the cross-run Advisor row difference from locally retained Phase R and S target reports without publishing resource IDs. Then continue remaining include/exclude tag and filter combinations and safe parent-to-child management-group traversal. Unfiltered Diagnostics caveats remain open. Quality Gate 004 is still planned, not passed.

### Phase S follow-up: Advisor row scope in the tag-only pass

**Date**

```text
2026-09-23
```

The user compared the locally retained, unredacted Phase R RG-only and Phase S tag-only target reports without sharing resource IDs. The selected inventory ID sets were identical (`SameSelectedResourceIds = True`), each with three resources. The earlier RG-only report had three Advisor rows; the tag-only report had two. Exactly one RG-only row was absent, with no added tag-only rows. The absent row was `HighAvailability`, `Medium` impact and targeted neither a selected inventory resource nor a child of one in either filtered report. The summary did not disclose whether this Advisor resource existed in either out-of-scope inventory; no raw reports or API response snapshots were independently inspected.

In the pinned AZQR source at `8e4f0577`, `internal/scanners/advisor.go` applies `filters.Azqr.IsServiceExcluded` to each Advisor resource ID. `internal/models/filters.go` excludes an Advisor row when an include-tag filter is active and no tag-scope decision is known for that ID or its recorded ancestors. Cloud Assess applies the same policy in `internal/advisor/scanner.go` and `internal/config/filters.go`. The RG-only filter need not exclude an otherwise structurally in-scope Advisor resource merely because it was not in the discovered inventory. The observed 3-to-2 difference is therefore **consistent with the source-compatible unknown tag-scope rule** when the same Advisor row is returned in both runs. The runs were separated in time; without their raw Azure API responses, a concurrent Azure Advisor change cannot be ruled out. No target-only semantic delta or defect has been demonstrated. A focused deterministic regression test now characterizes the include-tag unknown-scope case while preserving the selected-parent child case.

This closes the cross-run row *classification* with its timing limit. It does not certify every tag interaction: exclude tags, multiple tags, key/value edge cases and other downstream scope combinations remain on the roadmap. The Phase S result stays `equivalent = true` and `complete`; historical unfiltered Diagnostics warnings remain separate open evidence limits.

### Phase S evidence correction and recent-commit QA

**Date**

```text
2026-09-25
```

The review of the merged PRs #15–#17 and the pinned AZQR source identified an over-specific interpretation above. Both implementations record scope for *all* discovered inventory rows, including excluded rows (`SetResourceScope(resource.ID, !excluded)`). An Advisor ID that is neither a selected inventory ID nor its descendant can therefore match an explicitly excluded inventory ID or its descendant; it need not have *unknown* tag scope. The user's local summary did not compare the missing Advisor ID against `outOfScope`. The cross-run observation supports source/target agreement **within each run** and shows that the missing Advisor row was outside the selected inventory subtree, but cannot distinguish an excluded-scope decision, an unknown-scope decision, or a change in Azure Advisor results between runs. The previous classification of the exact cause as closed is withdrawn. No target-only mismatch has been shown.

The tag-only evidence still supports the `Environment: dev` include-filter behavior exercised in Phase S: both reports match semantically, and the user-reported ID-set comparison selects exactly the three matching resources from the earlier unfiltered Dev inventory, leaving nine nonmatching resources. The raw reports and API snapshots were not available to the reviewer, and the old Azure state cannot be reconstructed from the comparator alone. A regression assertion now covers explicitly excluded descendants as well as unknown and selected-parent descendants.

Recent-commit QA inspected merge `99d3b20` and its PR #15–#17 changes against the target specification, implementation plan, roadmap and Gate 004 plan. `git diff --check` on those changes passed; the postmerge workflow run `36085601101` passed both `quality` and `windows-validation`, including race tests, coverage floor, vet, reachable-vulnerability check, native PowerShell parser/helper checks and Windows Go tests. The active branch ruleset still requires PRs and both checks. The Gate 004 evidence matrix remains planned rather than passed. After initializing the pinned APRL submodule and obtaining Go 1.26.8 locally, focused tests and `go test -count=1 ./...` passed with this correction; the new PR must pass both required jobs again. The stale completion summary in the implementation plan was also corrected.

**Follow-up evidence boundary**

On the retained local Phase R and S target reports, check the absent Advisor row's resource ID against *both* `resources` and `outOfScope`, looking for the nearest recorded ancestor. Report only whether the nearest decision is included, excluded or absent, and whether its subscription, RG and scanner type meet structural filters; do not publish IDs. This can discriminate the recorded-scope hypotheses for the target report. Determining whether Azure returned the identical Advisor row to both scans requires historical API responses, which were not captured, or a new controlled paired run.

### Phase S follow-up: unknown Advisor tag scope confirmed in target inventory

**Date**

```text
2026-09-29
```

The user ran the read-only nearest-ancestor check from the Phase S follow-up on the locally retained Phase R RG-only and Phase S tag-only target JSON reports. The sole Advisor row present in the RG-only report but absent from the tag-only report returned `NearestRecordedScope = unknown`, `Category = HighAvailability`, and `Impact = Medium`. The check matched Advisor rows by normalized recommendation/resource IDs, then searched the tag-only target's included **and** out-of-scope inventory for an exact ID or nearest ancestor. No ID or ancestor matched. No raw IDs or reports were shared; the result is a user-supplied local summary, not an independent replay of the reports.

This resolves the previously open **excluded versus unknown recorded-scope** question for the target's tag-only inventory: the row has no recorded scope decision there. The checked-in `Environment: dev` include-tag fixture activates the fail-closed unknown-scope rule in both pinned AZQR and Cloud Assess. If Azure returned the same Advisor row to both runs and it passed structural filters, both implementations would exclude it in the tag-only run. The 3/3 RG-only and 2/2 tag-only Advisor results still come from separate scans; absent historical API response snapshots, the check cannot prove the row was returned to the tag-only scan or that the unknown-scope rule actually caused its disappearance. No source-versus-target discrepancy is evidenced within either pass. The Phase S tag-only `equivalent = true` result and its selective 3-resource coverage remain unchanged.

**Remaining evidence limit**

Treat historical Azure Advisor timing as unresolved when Gate 004 considers this cross-run observation. A new controlled paired run can test the rule under observed inputs if needed; no repeat of the already equivalent Phase S pass is required merely to classify its inventory scope. Continue the remaining tag interactions and suitable non-production nested management-group traversal.

### Phase T: AdvisoryDev exclude-tag live semantic comparison

**Date and provenance**

```text
2026-09-29
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 948264f5d9670c132adc2a6a612edc770f311a26
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260929_113319Z
```

The user supplied `equivalence.json`, `run-metadata.json`, and a local PowerShell summary of the unredacted target report and the earlier unfiltered Dev report. The raw reference/target JSON and logs remain local and were not independently inspected. The runner used the `AdvisoryDev` leaf management group, no CLI RG flag, and the implicit graph, diagnostics, Advisor and Defender stages. The reference exclude-tag fixture `examples/filters/azqr-exclude-environment-dev.yml` had SHA-256 `fd0d34a84177d2e127b921173a7934a3119401357263ed37a2d6d820e9e9679e`; the target counterpart had SHA-256 `01f58353c8f1026fb647c749969ef6c41b4007c6b1353128152bf17e15c53f00`. Both express exclusion of resource tag `Environment: dev` without an include filter.

**Execution and semantic result**

- reference, target and comparator exit codes: 0 / 0 / 0
- comparator: `equivalent = true`; no missing, extra or changed records in any enabled dataset
- recommendations: 314 / 314
- primary findings: 33 / 33
- resource types: 9 / 9
- in-scope inventory: 9 / 9
- out-of-scope inventory: 13 / 13
- Advisor: 8 / 8
- Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL and Cost: not enabled

The user's local target summary reported one scope-stage subscription, nine selected resources, zero selected with the Dev tag, and five Dev-tagged records in `outOfScope`. The selected ID set exactly matched the nine non-Dev-tagged in-scope resource IDs in the unfiltered Dev baseline from `20260923_145252Z`. This is non-empty selective evidence for the exclude-tag filter over the observed resources. The baseline is six days older, so the ID agreement is useful corroboration rather than proof that all Azure inventory remained static in the interval. `outOfScope` also contains resources outside that earlier in-scope set; its five Dev-tagged records should not be equated with the baseline's three selected Dev-tagged resources.

**Stage-health limit**

The target reported `complete_with_warnings`: Diagnostics completed with one `diagnostics_subrequest_non_success` warning and three stage records. Scope, inventory (9), graph (33), Advisor (8), and Defender (0) completed without warning codes; the optional stages were skipped. The warning's HTTP status and affected batch resource were not provided for this run, and the historical Network Watcher individual-GET evidence does not establish its cause here. Treat Diagnostics completeness as limited despite semantic equality. Advisor's 8 / 8 agreement is within this paired run; it does not resolve timing of the separate Phase R/S Advisor row.

**Next boundary**

Confirm the exact resolved subscription in a sanitized local check, classify this run's Diagnostics warning if its affected request can be identified, then continue other filter interactions and suitable non-production parent-to-child management-group traversal. Gate 004 remains planned.

### Phase T follow-up: scope and individual Diagnostics probe

**Date**

```text
2026-09-29
```

The user checked the retained Phase T target JSON locally. Across included inventory, out-of-scope inventory, findings, Advisor and Defender records, there was exactly one distinct subscription ID, and it was the expected non-production Dev subscription. This closes the Phase T target data-subscription check without publishing raw report rows.

The user then ran the read-only `tools/diagnostics-probe` against the retained Phase T target report. The tool counted one original HTTP 400 Diagnostics warning. It issued individual GETs for eight eligible inventory resources: seven succeeded, and one Network Watcher returned HTTP 400 `ResourceTypeNotSupported`. Its resource-ID hash matched the Network Watcher observed in the earlier Phase N/P individual-GET probes. Only the probe summary was shared; raw resource IDs and Azure reports remain local. This repeated observation strengthens the explanation that Network Watcher eligibility produces these Diagnostics warnings in this environment. The Phase T batch subresponse still lacks request correlation, and the later individual GET cannot prove that it was the failed batch request or rule out a transient batch failure. The target remains `complete_with_warnings`.

The pinned source and target have no dedicated missing-diagnostics recommendation for Network Watchers, as characterized in Phase O. The probe evidence therefore does not identify a missing Network Watcher recommendation delta, while the exact assessment meaning of the uncorrelated batch failure remains bounded by the warning. No additional Azure scan or comparator run was performed for this follow-up.

### Phase U: AdvisoryDev resource-group exclusion live comparison

**Date and provenance**

```text
2026-09-29
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 5ff380763f9364315dcfb28edde0958528b49426
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260929_121323Z
```

The user supplied `equivalence.json`, `run-metadata.json`, and a local summary comparing the retained target JSON with the earlier unfiltered Dev and Phase T tag-exclude target reports. Raw reports and logs remain local and were not independently inspected. Both sides ran against leaf management group `AdvisoryDev` with the implicit graph, diagnostics, Advisor and Defender stages. The reference fixture `examples/filters/azqr-exclude-dev-rg.yml` had SHA-256 `0ce146968517c0bc8d203f9753bba29d365c33564a16357b02b0a28ae7248004`; the target fixture `examples/filters/cloud-assess-exclude-dev-rg.yml` had SHA-256 `ba18dfa9d86c6578853d71817e08f8c558974c8aac94b2109b3d166418f28e83`. Each excludes the same full ARM ID of `rg-fos-FinopsHub-dev-euw`; no tag filter or CLI RG scope was used.

**Execution and semantic result**

- reference, target and comparator exit codes: 0 / 0 / 0
- comparator: `equivalent = true`, no missing, extra or changed records in enabled datasets
- recommendations: 314 / 314
- primary findings: 33 / 33
- resource types: 9 / 9
- in-scope inventory: 9 / 9
- out-of-scope inventory: 13 / 13
- Advisor: 8 / 8
- Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL and Cost: not enabled

The user's local check reported exactly three in-scope resources in the excluded group in the six-day-old unfiltered Dev baseline. All three IDs appeared in this run's `outOfScope`; none was selected. The nine selected IDs exactly matched the same-day Phase T tag-exclude target report. All populated target datasets had only the expected non-production Dev subscription ID. This is non-empty selective resource-group exclusion evidence for the observed inventory. Because the cross-run reports are separate Azure snapshots, their ID agreement corroborates selection but cannot establish that all Azure inventory or Advisor API output remained static.

**Stage-health limit**

The target reported `complete_with_warnings`. Diagnostics had one `diagnostics_subrequest_non_success` warning and three stage records; its warning message and failed batch request were not separately supplied for this run. Scope (one subscription), inventory (9), graph (33), Advisor (8), and Defender (0) completed without warnings; optional stages were skipped. The selected inventory ID set matches Phase T, where a later individual GET probe identified a Network Watcher HTTP 400, but that does not prove this run's uncorrelated batch failure had the same cause. Keep Diagnostics conclusions qualified. The within-run Advisor match does not resolve historical cross-run Azure Advisor timing.

**Next boundary**

Continue other filter interactions with discriminating inputs, especially include/exclude precedence and resource/recommendation exclusions. Parent-to-child management-group recursion, missing non-empty optional-stage data, and Gate 004 remain open.

### Phase V: AdvisoryDev recommendation exclusion live comparison

**Date and provenance**

```text
2026-09-29
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 9f87418b23e49924033f978267bea9a9750e112f
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260929_123149Z
```

The user supplied the equivalence and run-metadata JSON, plus a local check of the retained raw target/reference reports against the Phase U target report. Raw reports and logs remain local and were not independently inspected. The paired fixtures exclude recommendation `1981f704-97b9-b645-9c57-33f8ded9261a` on the `AdvisoryDev` leaf, with implicit graph, diagnostics, Advisor and Defender stages. The metadata records the reference and target fixture SHA-256 values as `f18568cdd722de0f392aad299ed759588f3b416498b3e68a4c45da488646f6e6` and `aeae852fed910ac3b2e9db16909997c61de190a80d6bfc81c6740b0aeea18b22`. These match the committed fixture bytes with Windows CRLF line endings.

**Execution and semantic result**

- reference, target and comparator exit codes: 0 / 0 / 0
- comparator: `equivalent = true`, no missing, extra or changed records in enabled datasets
- recommendations: 313 / 313; primary findings: 43 / 43
- resource types: 11 / 11; in-scope inventory: 12 / 12; out-of-scope inventory: 10 / 10
- Advisor: 10 / 10; Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL and Cost: not enabled

The user's local check found one finding for this recommendation in the earlier Phase U target report. In this pass, both target and reference have zero findings with that ID, while the earlier finding's resource remains in the new target inventory. This establishes non-empty, selective recommendation-exclusion behavior for that observed finding, in addition to within-run source/target equivalence. Separate Azure snapshots do not establish that all other service responses or findings remained stable.

**Stage-health limit and next boundary**

The user's local check reported target completeness `complete_with_warnings`; the supplied equivalence file likewise notes warnings. The specific stage warning codes and affected batch requests were not supplied for this pass, so Diagnostics results cannot be classified more narrowly here. This pass does not close the historical Diagnostics or Advisor timing limits, nested management-group traversal, other filter interactions, non-empty optional-stage gaps, or Quality Gate 004.

### Planning checkpoint: Individual-resource exclusion evidence fixture

The next paired filter pass uses `scripts/prepare-resource-exclusion.ps1` to derive one previously observed VM resource ID from the locally retained Phase U finding and Phase V inventory, without publishing the ARM ID. The helper writes paired AZQR `exclude.services` and Cloud Assess `exclude.resources` YAML into ignored `artifacts/` and makes no Azure calls. See [EQUIVALENCE.md](EQUIVALENCE.md#paired-filter-pass-exclude-one-observed-resource) for the commands and required result checks. This checkpoint prepared the test; the subsequent live result is recorded in Phase W. The existing Diagnostics warning, other filter combinations, nested management-group traversal and Gate 004 remain open.

### Phase W: AdvisoryDev individual-resource exclusion live comparison

**Date and provenance**

```text
2026-09-29
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 15ed0ca82fe8125c105a3ed6e74fad83df94565c
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
Evidence stamp: 20260929_132714Z
```

The user supplied `equivalence.json`, `run-metadata.json`, and a local check of the retained Phase U, Phase V and new target reports. The raw reports, generated filter YAML and logs remain local and were not independently inspected. The local helper generated AZQR `exclude.services` and Cloud Assess `exclude.resources` filters for one previously observed Dev VM. The metadata records their SHA-256 hashes as `ac87a6e858823bbb4cd603fc63ad2a1e2e372a271fc5e92844a586cb35443441` and `8bc729528969e53a959331056951c6d993d42a9a92364c2844eeab3287c50d4f`. Both scans used leaf management group `AdvisoryDev` and the implicit graph, diagnostics, Advisor and Defender stages.

**Execution and semantic result**

- reference, target and comparator exit codes: 0 / 0 / 0
- comparator: `equivalent = true`, no missing, extra or changed records in enabled datasets
- recommendations: 314 / 314; primary findings: 35 / 35
- resource types: 10 / 10; in-scope inventory: 11 / 11; out-of-scope inventory: 11 / 11
- Advisor: 5 / 5; Defender plan status: enabled on both sides, 0 / 0
- Policy, Defender Recommendations, Arc SQL and Cost: not enabled

The user's local target-report check confirmed that the selected VM was present in the Phase V inventory, absent from this run's selected inventory, present once in this run's `outOfScope`, and associated with zero new target findings. All other selected resource IDs matched the Phase V inventory exactly. This provides non-empty, selective exact-resource-exclusion evidence for the observed VM alongside within-run source/target equivalence. Findings and Advisor totals changed across separate Azure snapshots; their count differences alone cannot establish which rows were removed by the filter.

**Stage-health limit and next boundary**

The target reported `complete_with_warnings`: Diagnostics completed with one `diagnostics_subrequest_non_success` warning and six stage records; other enabled stages completed without warning. The warning's HTTP status, failed batch request and affected resource were not supplied for this pass. An earlier individual Network Watcher GET cannot be assumed to identify this particular batch response. The successful comparator therefore does not close Diagnostics coverage for the unsuccessful request. Nested management-group traversal, other filter interactions, non-empty optional-stage gaps and Quality Gate 004 remain open.

### Gate 004 hardening checkpoint: failure evidence and adjacent workflows

**Date**

```text
2026-09-29
```

The Gate 004 plan now states concrete evidence checks for adapter-to-report failures, Diagnostics uncertainty, missing non-empty projections, real output/privacy behavior and maintenance workflows. The later release decision includes a controlled production-scope pilot. This adds acceptance criteria; it does not pass Gate 004 or authorize a production scan.

A Diagnostics batch response with a different number of subresponses than requested previously left unmatched resources absent from the enabled-settings map, potentially creating missing-diagnostics findings without an explicit retrieval failure. The target now fails that stage on a response-count mismatch; a focused two-request/one-response test prevents a false missing-settings finding. This is an intentional safety correction for malformed API output. It is not a claim that the pinned source behaves identically under a truncated batch response, and existing successful-batch equivalence evidence is unaffected. Request-to-subresponse correlation for otherwise complete batches remains open.

A cross-package Cost access-denial test now checks the coordinator, application exit code and persisted JSON report together: the Cost stage must fail, the report must be partial, and healthy Graph findings must survive. Existing tests already cover other stage failures and individual Cost transport/authorization behavior. The maintenance workflows were also corrected to identify their generated commits as `github-actions[bot]` instead of pairing a bot name with the user's personal noreply address. Human-authored commits continue to use the user's configured identity. The workflows still need no-change and changed-output dispatch verification before their maintenance path is considered operationally exercised.

The live-equivalence runner now records a sanitized target stage-health summary in metadata schema `1.2`, including completeness, status, record counts and warning codes but no warning messages or resource IDs. This makes the recurring warning review reconstructable from the metadata without claiming that a warning code alone identifies its failing Azure request. Older `1.1` evidence bundles remain unchanged.

**Validation boundary**

PR #29 passed the hosted `quality` and `windows-validation` jobs (workflow run `36577883439`) and merged at `12d39739a70bfab03032c898f6459f0cc4e227e0`. Leave unperformed Gate 004 evidence rows open. Do not reclassify historical Diagnostics warnings as correlated from this new response-count check.

### Gate 004 output-boundary follow-up

**Date**

```text
2026-09-29
```

The existing cross-package assessment test now requests XLSX, JSON, CSV, SARIF and stdout in one run. It checks that raw subscription IDs are absent from JSON, every generated CSV table, every XLSX cell and stdout; that masked IDs are present in JSON, CSV, XLSX and stdout; and that SARIF retains the full resource ID needed to identify findings. This checks application option propagation and real renderer output together. SARIF contains sensitive resource identity by design and must be handled accordingly. The focused test and full Go suite passed locally. PR #30 passed the hosted `quality` and `windows-validation` jobs (workflow run `36578734674`) and merged at `7ce8e6c010b3fd6b1c43884c855325840f15c3c6`.

This deterministic test does not substitute for inspecting artifacts from a built CLI on supported platforms or checking Windows ACLs for unredacted reports. Those Gate 004 rows remain open, as do maintenance workflow dispatches and live Diagnostics request correlation.

### Diagnostics correlation tooling checkpoint

**Date**

```text
2026-09-29
```

An audit of the merged Gate 004 changes corrected the Phase output-boundary entry above: PR #30's two required jobs passed and its merge SHA is recorded. The ledger, gate matrix and roadmap still treat actual built-CLI artifacts, Windows ACLs, maintenance workflow dispatches and the historical batch warning correlation as open. Neither PR #29 nor #30 constitutes Gate 004 PASS.

The Diagnostics probe now offers an optional `single-batch` mode. It sends exactly one diagnostic-settings GET subrequest per ARM batch POST, so a returned subresponse can be associated with that one resource without assuming multi-request response order. Its console output contains a shortened resource-ID hash, type, status and safe Azure error code for failures, with no raw IDs or error messages. A deterministic fake-transport test exercises a successful response, an HTTP 400 with a sensitive message and an empty response array. The original individual-GET mode and request cap remain in place.

This tool has not been run against the user's Azure environment in this checkpoint. A later one-request batch failure can corroborate the Network Watcher hypothesis, but cannot identify the failing request in an old, uncaptured multi-request batch or prove Azure state has remained unchanged. Exact historical classification still requires correlated original evidence or an explicit accepted uncertainty. No Azure mutation or new scan was performed during this checkpoint.

**Validation**

The full Go suite, focused vet and diff check passed locally. PR #31 passed the hosted `quality` and `windows-validation` jobs (workflow run `36580144346`) and merged at `a2290322db919198e8eb9342eb0c2aba469974af`. The associated probe test is deterministic; no Azure batch outcome was observed during this validation.

### Gate 004 bounded-pagination checkpoint

**Date**

```text
2026-09-29
```

The Resource Graph query client previously followed any non-nil `$skipToken` without detecting an empty or repeated token. A malformed or faulty response could therefore cause unbounded requests and duplicate records. The client now stops with an error on either condition and checks cancellation before each page request. Focused tests assert that a repeated token never becomes a successful partial result and that a canceled context triggers no transport call. This is a deliberate safety correction for invalid pagination, not a claim about the pinned source's behavior under malformed responses. Normal continuation-token behavior and prior successful live equivalence are unchanged.

This deterministic boundary does not prove Azure throttling, malformed HTTP response, or end-to-end failure handling in every stage. Those Gate 004 checks remain open. No Azure requests were made for this checkpoint.

**Validation**

The focused ARG tests, full Go suite, focused vet and diff check passed locally. PR #33 passed the hosted `quality` and `windows-validation` jobs (workflow run `36581271240`) and merged at `d1fa30c10d68192b2dc963c994db3066e26712f1`. Gate 004 remains planned.

### Gate 004 Resource Graph response-shape checkpoint

**Date**

```text
2026-09-29
```

The Resource Graph HTTP transport previously decoded a successful HTTP response with missing or `null` `data` into a nil slice. The query client could then return an apparently valid empty dataset. The transport now rejects this malformed response. A test runs the HTTP transport through the query client and checks missing and `null` data fail, while a genuine empty `data: []` succeeds. The documented 2024-04-01 Resource Graph response includes a `data` field; this is a target safety correction for malformed output, not a claim of pinned-source parity under invalid responses. Existing valid live reports are unaffected.

The change adds one deterministic negative-path check. Azure throttling, timeouts, cancellation across real adapters, CLI failure persistence and other Gate 004 evidence remain open. No Azure requests were made for this checkpoint.

**Validation**

The focused ARG test, full Go suite, focused vet and diff check passed locally. PR #35 passed the hosted `quality` and `windows-validation` jobs (workflow run `36582502895`) and merged at `ddba2a9842651854606901ac25b5a58910f3c9d0`. Gate 004 remains planned.

### Phase W Diagnostics one-request batch follow-up

**Date and provenance**

```text
2026-09-29
Probe code revision: 8c4f99aa9a91c10e88cd0f27164096fe77292887
Original evidence stamp: 20260929_132714Z
Original target scan revision: 15ed0ca82fe8125c105a3ed6e74fad83df94565c
```

The user fast-forwarded their local `bootstrap/core-v1` checkout to `8c4f99a` and ran `go run ./tools/diagnostics-probe --mode single-batch` against the locally retained, unredacted Phase W `target.json`. The probe used the user's Azure authentication to issue one read-only diagnostic-settings GET subrequest in each ARM batch POST for the ten eligible inventory resources. Only the sanitized console summary was supplied; the raw report, request and response bodies and logs remain local and were not independently replayed by the reviewer.

The summary counted one HTTP 400 warning in the original target report. In the later probe, nine one-request batches succeeded and one returned HTTP 400 for `microsoft.network/networkwatchers`, with shortened resource-ID hash `4d7a2aeb63c53438`. No Azure error code was present in this batch probe's sanitized output. The hash matches the Network Watcher previously observed in individual GET probes of the broader Dev evidence, where `ResourceTypeNotSupported` was returned. The new observation directly associates a current one-request batch HTTP 400 with that Network Watcher, strengthening the explanation for the historical Phase W warning. It does not prove which subrequest failed in Phase W's original multi-request batch or that the Azure response remained unchanged between runs.

Network Watchers are included in both pinned-source and target Diagnostics request lists but have no dedicated missing-diagnostics recommendation in their shared behavior. If this Network Watcher was the sole failed historical subrequest, it could not itself produce a dedicated missing-diagnostics finding. That conditional conclusion does not rule out a different historical failed request. Keep Phase W `complete_with_warnings` and its Diagnostics assessment uncertainty; do not promote the original semantic-equivalence PASS into full Diagnostics coverage or close Gate 004 on this basis.

### Nested management-group evidence boundary and synthetic adapter test

**Date**

```text
2026-09-29
```

The user has no available non-production parent management group with nested child groups for a live paired scan. The known `Advisory` parent includes a production child. A subscription include filter is not a safe substitute for a scoped parent test: the current discovery walker asks for the descendants and direct subscriptions of every visited group before deciding which returned subscriptions pass the filter. Do not run the parent scope as an implicit workaround. No production parent was scanned and no Azure hierarchy was changed.

A new deterministic fixture drives the production Azure SDK scope adapter and the discovery walker together with mocked management-group subscription and descendant HTTP responses. It covers a synthetic root, child and leaf, an all-descendants response that repeats the leaf, a subscription include filter, disabled-subscription exclusion, and a denied leaf request that fails without returning a partial scope. It exercises actual SDK response decoding and local traversal without Azure calls. The earlier walker-only recursion test remains in place.

The fixture strengthens implementation evidence, but it is not a live source-versus-target parent-group comparison and does not establish real-world RBAC visibility or Azure behavior in an accessible nested hierarchy. Leave the live traversal row open. At Gate 004, either supply suitable live evidence or record an explicit accepted limitation and narrow the first-release management-group support claim; the target specification currently includes recursive management-group resolution. Other Gate 004 evidence can proceed independently.

**Validation**

The focused adapter test, full Go suite, focused vet and diff check passed locally. PR #38 passed the hosted `quality` and `windows-validation` jobs (workflow run `36593065309`) and merged at `06380af0f9136836d7cd251eeb4b52cd36deca48`. The live nested-scope limitation remains open.

### Deferred validation register and next filter fixture

**Date**

```text
2026-09-29
```

The operator confirmed that no suitable non-production nested management-group parent is available. Live parent-to-child traversal is now explicitly tracked as `DV-001` in [DEFERRED_VALIDATION.md](DEFERRED_VALIDATION.md), separately linked from the README, roadmap and Gate 004 plan. Its status is deferred evidence, not PASS or a preaccepted release limitation. The register states the safe resume trigger, required paired-run evidence, and the first-release scope decision if it remains unavailable. The `Advisory` production-containing parent was not scanned.

To continue with an accessible `AdvisoryDev` scope, paired same-resource-group include/exclude fixtures were prepared for pinned AZQR and Cloud Assess. The fixture puts `rg-fos-FinopsHub-dev-euw` in both lists, exercising the characterized explicit-include precedence with non-empty previously observed Dev resources. [EQUIVALENCE.md](EQUIVALENCE.md#paired-filter-pass-include-and-exclude-the-same-dev-resource-group) records the command and checks. At this checkpoint the live pass had not yet been executed; its later result is recorded below. Prior Phase R resource counts were historical expectations only.

### Same-resource-group include/exclude precedence live pass

**Date and provenance**

```text
2026-09-29
Evidence directory: C:\src\Cloud-Assess\artifacts\equivalence\20260929_160941Z
Reference: 8e4f0577f3615e6c9014c031bcad079f235369cc
Target: 92dfaeb0a8c4a124dea4c7f2fd9b85f22e6d68da
APRL: 60eaddda76541f6adbc1c5ffa686829807e55e29
```

The operator ran the paired filters from the preceding checkpoint on the `AdvisoryDev` leaf. Each filter lists the same full Dev resource-group ID in both its include and exclude resource-group fields. Run metadata records the exact commands, filter paths and SHA-256 values, implicit default stages (`graph`, `diagnostics`, `advisor`, `defender`), one scope record, and exit code 0 for the pinned reference, target and comparator. The uploaded metadata and equivalence report were inspected; the unredacted source and target reports and logs remain in the operator's ignored local evidence directory.

The comparator returned `equivalent = true`. Recommendations were 314/314, findings 11/11, resource types 3/3, selected inventory 3/3, out-of-scope inventory 19/19, Advisor 2/2 and Defender plan status 0/0, with no missing, extra or changed records in enabled datasets. Policy, Defender Recommendations, Arc SQL and Cost were not enabled. Target completeness was `complete`; scope, inventory, graph, diagnostics, Advisor and Defender stages completed without warnings.

The operator's local check of `target.json` found all three selected resources in `rg-fos-FinopsHub-dev-euw`, the same normalized selected resource IDs as the retained Phase R RG-include target report, and only the expected Dev subscription among the target's recorded data subscriptions. These checks support the observed include-over-exclude precedence with non-empty resources and paired source/target semantics. The older Phase R report is a separate Azure snapshot; ID equality is corroboration rather than a substitute for the current-run comparison. The reviewer did not independently inspect the unredacted target JSON. This pass does not exercise a different group combination, a nested management-group parent, or optional stages. `DV-001` remains deferred and Gate 004 remains open.

### Progress and roadmap review checkpoint

**Date**

```text
2026-09-29
Reviewed branch: bootstrap/core-v1 at ef3bd8b3ae5832ef6b6354f26e45725a3d4e8fea
```

The target specification, implementation plan, roadmap, Gate 004 matrix and recent live evidence were reviewed after the same-RG precedence pass. Steps 1-27 of the implementation plan are coded; step 28 has partial live coverage; step 29 has required Linux/Windows CI but no distributable release; steps 30-31 require security/license work and the plugin scope decision. Step count alone is not a completion percentage. [ROADMAP.md](ROADMAP.md#progress-estimate-and-next-checkpoint) now carries rounded, judgmental estimates for the exercised core scan, core-v1 acceptance, installable toolkit and broader AZQR parity, with the evidence limits stated alongside them. No gate status changed.

Next, strengthen a failed/throttled Azure adapter path through the query boundary with a deterministic fixture; then continue remaining filter interactions and missing optional-stage projections. This can proceed without Azure access. The nested management-group live test remains DV-001, deferred rather than passed or accepted. Reassess the estimates after Gate 004 evidence and the first-release feature boundary are decided.

### Gate 004 Resource Graph throttling characterization

**Date**

```text
2026-09-29
```

A deterministic HTTPS fixture now drives the production authenticated Azure HTTP client, Resource Graph HTTP transport and query client together. It contrasts a structured HTTP 429 response with a successful `data: []` response. With the retry budget set to zero to model an exhausted retry path, the former must return a wrapped Azure response error and nil query result; the latter must return a valid empty result. The fixture checks that each case made exactly one authenticated POST and makes no Azure calls. It closes one failed-retrieval-versus-empty-result regression path, not the default retry timing, real Azure throttling, coordinator stage health, or all Gate 004 negative paths.

The initial PR #43 CI run (`36597984445`) failed on both Linux and Windows because the fixture set `MaxRetries: 0`, which the Azure SDK interprets as its default of three retries. The mocked HTTP 429 therefore received four requests and triggered the fixture's one-request assertion. This was a test-configuration error, not an observed scanner failure. The fixture was corrected to `-1`, the SDK's documented one-try/no-retry setting. The code was reviewed and `git diff --check` passed locally. The scratch environment has no Go toolchain, so hosted checks provide executable verification. The corrected PR #43 run (`36598430521`) passed both required `quality` and `windows-validation` jobs and merged at `2afcf7146612d51ec30eab64c7c656b7acdf4ad5`. This verifies the bounded HTTP-client-to-ARG-query negative path; default retry timing, coordinator status and other failures remain open.

### Failure register and Graph critical-failure report fixture

**Date**

```text
2026-09-29
```

At the operator's request, [FAILURE_NOTES.md](FAILURE_NOTES.md) now indexes confirmed mistakes with a recurrence search, correction, prevention check, status and original evidence. It starts with the over-specific Phase S Advisor interpretation, the vague local same-RG handoff, and the PR #43 zero-retry test assumption. The notes do not replace this ledger, Git or CI. New failures should be checked against the register before a repeated workaround is attempted; record a new entry or append a recurrence without including credentials or raw Azure rows.

A cross-package test injects a critical Graph query error at the coordinator operation seam after successful scope and inventory. It drives the coordinator and application JSON renderer together, requiring exit 1, persisted `failed` completeness, a failed Graph stage, skipped Advisor, retained discovered inventory and zero fabricated findings. This complements the separate real authenticated HTTP 429-to-query fixture. It does not join that HTTP fixture to the coordinator in one test, prove CLI process behavior, or validate other adapters. No Azure requests were made. The scratch runtime has no Go toolchain; PR #45 passed required `quality` and `windows-validation` jobs (run `36600171189`) and merged at `94c954266d8961276dea1971d506fce0410b5cad`.

### Gate 004 CLI subprocess exit and report checkpoint

**Date:** 2026-09-30.

Added `runWithExecutor` around the existing command executor seam; normal `run` continues to select `executeScan`. A test-executable subprocess runs the production command dispatcher and application renderer against deterministic assessment fixtures. The parent checks actual process exits 0 (complete), 1 (execution failure), 2 (severity gate), and 3 (partial), then verifies persisted JSON completeness, retained inventory, and expected findings after the child exits.

This is not a built or installed Azure CLI test and does not exercise production credentials or join the HTTP 429 fixture to the coordinator. Other adapter failures, default retry timing, cancellation and release artifact checks remain open. No Azure requests are made. The runtime has no Go toolchain. PR #47 passed required `quality` and `windows-validation` jobs in run `36715889734` on commit `768004c7458ab858ce1769d9d0c755e6f434981f`, then merged at `d3aaf1092233ee406dcc0123e4d3d49ad67c39a3`. Local `git diff --check` also passed. Gate 004 remains planned.

Recorded two mistaken inspection paths and their recurrence prevention as FN-004 in [FAILURE_NOTES.md](FAILURE_NOTES.md).

### Gate 004 in-flight interruption checkpoint

**Date:** 2026-09-30.

Added local TLS endpoint fixtures through the production authenticated HTTP client, ARG HTTP transport and query client. The endpoint confirms receipt before cancellation; a separate caller deadline expires while the response is pending. Both must preserve context error identity, return no result and stop after one request with retries explicitly disabled. These tests do not establish default retry timing or automatic operation timeout behavior.

A separate coordinator/application fixture interrupts the noncritical Advisor stage after healthy inventory and Graph findings. It verifies exit 1, failed completeness, specific cancellation/deadline stage codes, skipped Defender, preserved JSON inventory/findings and no fabricated Advisor rows. The deadline is an expired child context created at the operation seam; this fixture does not connect the TLS request to the coordinator or test OS signals. No Azure requests are made. PR #49 passed required Linux `quality` and `windows-validation` jobs in run `36718219604`, on tested commit `a02aec3814f7f29c28ea3b70c781c58f24602551`. Local Go is unavailable; `git diff --check` passed. Merged at `df0d52bceab3a994b10cb45404037429358d52df`.

**Inspection issue:** `HTTPClientOptions.OperationTimeout` is assigned by `DefaultHTTPClientOptions` but never consumed by the HTTP implementation. Its advertised value therefore does not impose a total operation deadline. Existing caller contexts and per-attempt `Timeout` remain separate mechanisms. Track an explicit implementation/removal decision and bounded retry/body-read tests before making total-duration claims. No production timeout behavior is changed in this checkpoint.

### HTTP total-operation timeout remediation

**Date:** 2026-09-30.

Resolved the unused `OperationTimeout` configuration: a positive value now derives one caller context around each HTTP operation, covering authentication, retries and response consumption. An earlier caller deadline still wins. Defaults remain `Timeout * 10`; non-positive custom values add no operation deadline. This deliberately corrects target behavior and is not an AZQR live-equivalence claim. It does not add a whole-scan budget.

`PostStream` transfers deadline cleanup to its response body and cancels on EOF/read error/close, rather than canceling at return. Buffered Get/Post calls release context on return; request/setup/error paths also release it. The pinned Azure SDK buffers bodies by default; the ownership wrapper preserves the method contract if a body remains unread.

Focused fixtures cover a 30-second retry hint interrupted by a shorter operation deadline with retries enabled, blocked body reads for Get/Post/PostStream, earlier caller deadline and cancellation, zero/negative settings, and successful stream lifetime through EOF/close. These use the production SDK pipeline with a context-aware synthetic transport, not live Azure. Existing TLS interruption fixtures remain separate. PR #51 passed required Linux `quality` and `windows-validation` jobs in run `36720076461` on commit `bed7760bdb5feac615b350a1492744b2d85c6dc4`, merged at `d465e173155f1ae5329151a805f6e55ad476e339`. Local Go is unavailable.

Source inspection used pinned azcore `v1.23.1` retry, pipeline and body-download implementations, rather than assuming zero-retry or streaming behavior. Local `git diff --check` passed. Gate 004 remains planned; default retry sequencing, other adapters and whole-scan bounded execution remain open.

### Autonomous QA control expansion and report hardening

**Date:** 2026-09-30.

Implemented comparator missing/extra/changed/coverage/precondition guards, field-symmetry fuzzing, exact-resource exclusion properties, and an isolated four-mutation harness that requires assertion failures rather than compile failures. Existing Diagnostics-only normalization and canonical summary/ownership regressions remain in force. Found and corrected sparse/custom dataset-name ordering: custom names are sorted independently of how many standard datasets exist. Standard projection ordering is unchanged.

Introduced shared staged report writes for JSON/SARIF/CSV/XLSX. Failed rendering preserves an existing file; successful Unix replacement is mode 0600, including formerly permissive destinations. Added temporary-file cleanup, invalid-destination and symlink rejection tests; application output failure must exit 1 while retaining the assessment. CSV cells with formula-like/control/fullwidth prefixes become apostrophe-prefixed text; XLSX ordinary untrusted cells are checked for absent formulas. These are deliberate target safety corrections outside the raw assessment-semantic equivalence claim. Canonical JSON data is unchanged.

Windows CI checks report inheritance from a controlled restricted directory. This does not validate an operator's arbitrary directory ACLs. Go rename does not guarantee atomic replacement everywhere; process-kill/power-loss durability and multi-file transactions are not claimed. See QA_PROCESS for limits and OWASP spreadsheet-consumer caveats.

Added 3,001-subscription batching, failure-after-first-page and configured transient retry recovery checks; strengthened ARG query endpoint and Diagnostics GET-subrequest contracts. These are selected adapter bounds, not proof every future request is read-oriented. Retry recovery uses explicit millisecond SDK options, not a default-backoff timing claim.

CI adds bounded two-target fuzz runs, isolated comparator mutation checks, 20-minute job limits and synthetic failure-evidence retention using verified pinned upload-artifact v4.6.2. Added PR review template and QA_PROCESS feedback loop, linked roadmap/specification/Gate 004. PR #53 passed required Linux/Windows jobs in corrected run `36726089065` on head `700223a33be5de919fc4c0177224e0018970e4a8`, merged at `d36e810e56649d9e85979e924abf03e5e602dc3e`. Both bounded fuzz targets and all four comparator mutations passed; aggregate statement coverage was 79.3%, which is not a readiness measure. Local diff check and Python syntax parsing passed. No Azure calls or laptop action are required. Gate 004 and the later release decision remain open.

Initial run `36725561754` failed the Windows ACL fixture before assertions because Windows PowerShell could not load its Security module. Isolated its child PSModulePath; corrected Windows ACL assertions passed in run `36726089065`, and FN-005 records the failure and verified correction.

Logged the mistaken test-path inspection recurrence in FN-004. Development guidance was checked against official Go fuzz/rename documentation and OWASP CSV Injection; release provenance remains a later requirement, with no SLSA level claimed.

### Windows replacement privacy follow-up

**Date:** 2026-09-30.

Follow-up review found that staging a replacement could broaden an existing Windows report's DACL if its directory permitted more readers. The creation-only ACL fixture did not cover this. Copy the existing effective DACL before writing report bytes and protect it from parent inheritance; fail replacement if the DACL cannot be obtained or is absent. Add a fixture that freezes a restricted report ACL, broadens the directory, replaces the report, and checks the replacement still allows only the original test identity. New files continue to require a suitably secured directory. Owner/SACL preservation and future inherited ACL updates are not claimed.

Use already-pinned `golang.org/x/sys v0.47.0` Windows security APIs; the dependency becomes direct, with no version change. PR #54 passed required Linux/Windows CI in run `36728008504` on head `017d5487c59233033fb1d63cf185e9fe841aa53a`, merged at `ce61babc85e52fa2d8ca0aa64a9283bb37e0f2fb`. The controlled Windows creation and replacement ACL assertions passed. FN-006 records the verified gap/correction; FN-007 records the inspection-helper error.

A source search also found a Go toolchain outside PATH. It successfully selected/downloaded Go 1.26.8, so the earlier local-unavailability assumption is corrected as FN-008. Earlier checkpoints used actual hosted CI evidence; their results remain valid. After initializing the exact APRL pin, local `go test -race -count=1 ./...` passed with Go 1.26.8. Windows reportfile cross-compilation and module tidy also passed; only the already-pinned x/sys dependency moves from indirect to direct. Required hosted Windows/Linux runtime checks subsequently passed in run `36728008504`. This closes this implementation follow-up, not Gate 004 or arbitrary-directory ACL review.

### Optional-stage HTTP health through persisted reports

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `869dd3b03458ffc6ecde3a28115e8d600b314955`.

Continuation review confirmed PRs #53/#54/#55 are merged and retained the open Gate 004, release, DV-001, historical warning/Advisor and non-empty optional-stage boundaries. Added twelve synthetic TLS cases across Policy, Defender plan status and Defender Recommendations: HTTP 403, HTTP 429, HTTP 200 missing its data array, and HTTP 200 with an empty data array. Fixtures execute the production authenticated HTTP client, ARG transport/query client, selected scanner, coordinator, application and JSON writer. Healthy scope/inventory/Graph operations remain injected fixtures.

Failed retrieval must persist a failed optional stage, zero optional records, partial completeness and exit 2 while retaining healthy inventory/findings. Successful empty retrieval must persist a completed zero-record stage with complete status and exit 0. The endpoint checks authentication, query, selected subscription, POST query endpoint/version and exactly one request with retries explicitly disabled. These envelopes are fabricated contract fixtures, not recordings from Azure; they do not establish non-empty row projection, default retry sequencing, management-group discovery or complete read-only assurance.

Focused `go test -race -count=1 ./internal/app -run TestOptionalStageHTTPHealthPersistsReport` passed all cases. Full local `go test -race -count=1 ./...` also passed. Required Linux/Windows CI passed in PR #56, run `36730447316` on code head `f1e7358731f20e9b2f2430276656eeb1cfa90795`. The final documentation revision passed required run `36731029675` and PR #56 merged at `951b49c99d841aa45124d41f30e4b1cfa8095b36`. No Azure calls or laptop action are needed. Inspection path-assumption recurrence is recorded in FN-004.

### Row-decoder properties and source-shaped optional-stage projections

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `951b49c99d841aa45124d41f30e4b1cfa8095b36`.

Inspected the pinned reference Policy/Defender row builders and `graph.UnmarshalRows`. Added twelve explicit decoder shape cases and a fuzz property checking row-count conservation, partition equivalence, healthy-neighbor isolation and caller-byte ownership, with an 8 KiB input bound. Null, empty objects and null fields remain accepted zero-valued rows under the source-compatible JSON decoder; this is type tolerance, not business-record validity. No production decoder behavior changed. Existing Arc SQL numeric `vcores` rejection characterization remains; the response-shape resolution is still open.

Extended the optional-stage HTTP-to-report fixture from twelve to twenty-one cases. Literal source-reviewed expected JSON now checks one representative non-empty row per Policy/Defender stage. Mixed valid/malformed data retains the correct row and explicit malformed-row warning; all-malformed data retains a visible warning and `complete_with_warnings`, rather than a warning-free empty success. This does not execute KQL, independently run the source scanner or establish live non-empty equivalence. Fixture provenance and exact limits are in `internal/app/testdata/README.md`.

Focused race tests passed. The initial ten-second local fuzz run failed at its deadline without a property assertion or failing corpus; FN-010 retains this unresolved termination cause. An execution-count bound of 100,000 with a separate 60-second timeout passed locally. The new CI target uses that bound; existing fuzz targets are unchanged. Full local `go test -race -count=1 ./...` passed after the final fixture changes. PR #57 code head `bde30f36c834fb41766b0821c6b9966e48094926` passed required Linux/Windows CI in run `36733195533`, including the count-bounded decoder fuzz target. The final documentation revision must also pass both required checks before merge. FN-007 records the oversized inspection-output correction. No Azure calls or laptop action are required; Gate 004 remains open.

## Current boundary

The generic core scan, deterministic equivalence tooling, reproducible live runner, default/optional/resource-group/two-subscription/leaf-management-group and separate Storage/VM, RG include/exclude, tag include/exclude, recommendation-exclude and individual-resource-exclude live passes, plus enforced development CI, are complete for the behavior exercised so far. The unfiltered two-subscription, leaf-management-group, tag-exclude, RG-exclude and individual-resource-exclude passes have explicit Diagnostics warning boundaries; the recommendation-exclude pass is also `complete_with_warnings`, without supplied warning details. Network Watcher individual-GET probes and the Phase W one-request batch probe support an explanation for earlier warnings but have not mapped the original multi-request batch responses. The missing cross-run Advisor row has unknown recorded tag scope in the tag-only target inventory; historical Azure timing remains unresolved.

The next live validation boundary is:

```text
Nested management-group traversal and remaining filter combinations; retain historical Advisor timing and unfiltered Diagnostics warning caveats
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
- exclude-tag pass Diagnostics batch request correlation (individual probe identified a matching Network Watcher error)
- historical Advisor API row presence/timing across separate RG-only and tag-only runs
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
- a confirmed development mistake needs a recurrence check in [FAILURE_NOTES.md](FAILURE_NOTES.md)
- a behavior is intentionally changed from the source
- a source defect is deliberately corrected
- a new quality gate is performed
- the pinned reference/provenance changes
- live equivalence identifies a new class of delta
- a deferred feature enters or leaves scope
- a release baseline is created

Routine commits do not need individual prose entries because the Git log already records them. Significant commits and every validated milestone must be indexed here.

This makes the development process reconstructable even if the conversational history that produced it is unavailable.
