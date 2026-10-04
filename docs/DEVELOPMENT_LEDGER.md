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

### 2026-10-03: PR97 request library accepted, next discovery contract preserved

Final tested head `f58895fcb9cd51817d49b83f26f71b7073443739`, tree `eeb4313c18fedb890363777e095631fcb017e78f`, [run37092703311](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37092703311), quality111116220428/Windows111116220521 PASS. Required steps/full logs inspected: eleven real CLI, nineteen default/nineteen custom package cases per host without acceptance skips, docs291/9/1, coverage81.4%, four compiling AI mutations/restored baselines, full Linux race/vet/fuzz and mandatory provenance/module/inventory/branding/maintenance/comparator checks. Reachable scans reported no findings; existing module-only advisory remains open. FN055 retains ignored lone audience override and failed noncompiling mutation run37092489006; neither earlier green nor failed run certifies acceptance.

Rules23890737 active with both required checks/no bypass; current base/head/clean preview exact tree/ordered parents verified. Expected-head protected merge `aadd0d03b13a093a9ec3b9e69e3f70208d1e1714`, ordered parents9c6d0e79d05d06559fb4b543406b928013aed835 and `f58895fcb9cd51817d49b83f26f71b7073443739`; remote ref/PR/tree/parents/human author/GitHub committer verified. Source/APRL/AOR/CUSTOM/SKU pins/dependencies/source-capture bytes unchanged. Request library only; public AI is unavailable. Native push run37093027663 is separate postmerge evidence, pending at this checkpoint preparation; fresh local fetched/source/build proof is blocked by FN056.

[AI_GOVERNANCE_EXECUTION.md](AI_GOVERNANCE_EXECUTION.md) records exact next bounded discovery query/wire/filter/security/health acceptance and separate public all-format integration, based on a fresh pinned-source read and primary Microsoft paging documentation. Existing ARG Client lacks total page/row/body bounds and loses prefix results on query errors, so unchecked reuse is not accepted discovery design. Source query has no tags: unknown include-tag decisions fail closed, exclude-only unknown decisions remain included; mixed scans use recorded inventory decisions. Preserve independent source expectations and no invented tags.

Handover/roadmap/specification/implementation/request/core/QA summaries reconciled to accepted PR97. Restore an isolated development executor/clean pinned source and actual submodule/toolchain before implementation; no user laptop/Azure test needed. All code and next contract remotely preserved. Discovery/public integration/region/B6/B7, live/Gate004/release/advisory boundaries remain open. This documentation proposal requires its own exact-head native gates and protected acceptance; final proof belongs in its PR without a self-referential documentation loop.


### 2026-10-03: PR97 compiling-negative-control repair

Correction1945da7/run37092489006 passed production package tests but failed both new mutation gates, Linux111115580824/Windows111115580893, because the removed guard made audienceSet unused. FN055 records the harness mistake and unchanged requirement that a compiling fault fail the named behavioral assertion. Use a false term referencing audienceSet; fresh complete native gates are required. No production change, acceptance skip or weakened test oracle.

### 2026-10-03: Fresh-session PR97 review and cloud-boundary correction

Verified live core-v1 PR96 merge 9c6d0e79d05d06559fb4b543406b928013aed835/tree b032e44ef99aec0b2214e65ace3f15194c842a72, ordered parents/identity and only open PR97. Previous request head563c973/tree d72948606ba4f2b6921f9c379b530670585978db passed run37091739775, quality111113344322/Windows111113344512; mandatory steps/full logs read, including three AI mutations, eleven real CLI and nineteen default/nineteen custom cases per host, coverage81.4%. Existing module-only advisory is unchanged. Green CI did not cover the actual automated review's lone-audience override finding, so no merge was attempted.

FN055 correction requires all three cloud overrides or none, six partial cases with rejection/zero token/transport counts, complete-public acceptance and fourth compiling mutation. Fresh native gates remain mandatory on the corrected exact head. No source pin, fixture expectation, dependency or public plugin change. FN056 records missing current local Go/Git HTTPS/download capabilities and disabled Code Review MCP; local build/clean fetched proof cannot be claimed. Publish the focused correction through authenticated GitHub, retain dependent discovery/public implementation until a suitable executor is restored; Azure/laptop/Gate004/release deferrals remain open. Final run/merge proof belongs in PR97 and the next coherent checkpoint.


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

Focused race tests passed. The initial ten-second local fuzz run failed at its deadline without a property assertion or failing corpus; FN-010 retains this unresolved termination cause. An execution-count bound of 100,000 with a separate 60-second timeout passed locally. The new CI target uses that bound; existing fuzz targets are unchanged. Full local `go test -race -count=1 ./...` passed after the final fixture changes. PR #57 code head `bde30f36c834fb41766b0821c6b9966e48094926` passed required Linux/Windows CI in run `36733195533`, including the count-bounded decoder fuzz target. The final documentation revision passed required run `36733770483`; PR #57 merged at `05657daadc0e840d86b981f88547e2fcc4d01d72`. FN-007 records the oversized inspection-output correction. No Azure calls or laptop action are required; Gate 004 remains open.

### Maintenance proposal publication controls

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `05657daadc0e840d86b981f88547e2fcc4d01d72`.

Both maintenance workflows previously duplicated their scoped staging, bot commit and proposal push logic. Extracted the publication boundary into `scripts/publish-maintenance.sh`, called by dependency tidy and pinned rule import. It requires the expected base branch and GitHub run metadata, rejects pre-existing staged paths outside the selected maintenance scope, creates no branch/commit for unchanged output, and pushes only a run-specific proposal ref without force. Bot author and committer identity are set for the commit command, including when caller identity variables exist. A rejected push remains a failed step and cannot emit a success/PR summary. Added 20-minute job limits to both maintenance workflows.

Seven executable tests use isolated working repositories and local bare Git remotes. They cover both no-change and changed-output kinds; literal remote/base ref checks; bot identity; exact changed paths; a genuine synthetic submodule gitlink; unrelated staged changes; invalid run ID/wrong starting branch; rejected push; and divergent existing proposal refusal. Only local fixture repositories are written. The required Linux quality job now runs these tests.

Local publication tests, Bash syntax, workflow YAML parsing and diff checks passed. PR #58 code head `29bce11cc613df3578b9f6410672a65a3337c25e` passed required Linux/Windows CI in run `36735297259`, including the publication tests in the Linux quality job. The final evidence revision passed required run `36735705311`; PR #58 merged at `108ac4d16acc8f66fdac0ad91ab27fbd55b94cb6`. These tests exercise the actual publisher, not dependency resolution or full rule-import preparation, GitHub token authorization, dispatch/event delivery or an automatically created PR. Those complete maintenance workflow evidence paths remain open at Gate 004. A proposal push with `GITHUB_TOKEN` alone is not a quality-check result; open/review its PR and observe both required checks. No Azure or laptop action is required.

### Maintenance generation reproducibility and rejection checks

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `108ac4d16acc8f66fdac0ad91ab27fbd55b94cb6`.

Downloaded the exact pinned reference archive and executed the actual rule-import workflow recipe in disposable Git checkouts, with curl replaying the retrieved archive and APRL cloned from the locally initialized exact pin. The importer reproduces the committed custom/orphan Git trees and APRL gitlink without a diff. A deliberately corrupted committed rule is restored into a scoped local proposal while the baseline ref remains unchanged. Added pre-replacement source-layout guards and staged-tree/gitlink checks before publication, using the same fixed provenance values as required quality CI.

Six executable generation cases cover unchanged import, corrupted-rule repair, invalid archive, missing source directories, wrong snapshot hashes, and real `go mod tidy` removal of an unused local dependency followed by stabilization. Tests read the production workflow run blocks with step-bounded extraction; no independent copy of the generator is used. Generation fixtures run in Linux quality CI alongside publication fixtures. The first compatible local tidy fixture used Go 1.26.0; the complete six-case suite was rerun with the explicit Go 1.26.8 binary and passed. Source/module prerequisites and transport-replay boundaries are documented in MAINTENANCE.md.

Local six-case generation tests, Python syntax, workflow YAML parsing and diff checks passed. Final review reproduced a shallow-checkout fixture bootstrap failure before concluding hosted readiness. FN-011 records it; fixtures now always use depth one and permit only their disposable local bare remotes to accept the shallow bootstrap. Initial PR #59 run `36738474726` failed Linux job `109966302312` at fixture bootstrap pushes; Windows passed. The corrected six-case suite passed locally with depth-one fixtures and Go 1.26.8. Corrected code head `620de86e8756007598f1e6565c8855e923fa7d6a` passed required Linux/Windows CI in run `36739045844`, including all six generation cases. The final evidence revision passed required run `36739580927`; PR #59 merged at `532c95aa3b3fdfe4015d66567be5e57c501809fc`. These checks prove recipe/output behavior with pinned data and local Git mirrors, not hosted dispatch/token authorization or the full maintenance proposal-PR event path. Gate 004 retains those limits. No Azure or laptop action is required.

### Generated dependency and license evidence inventory

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `532c95aa3b3fdfe4015d66567be5e57c501809fc`.

Added a standard-library-only Python generator using the selected Go module graph and separate CGO-disabled Linux/Windows amd64 CLI package lists. Versioned modules are downloaded in an isolated temporary module so graph-only checksum collection cannot modify repository go.sum. The manifest includes all 49 selected modules (nine direct); 24 contribute to the Linux CLI package graph and 26 to Windows. Root license/notice text is retained with SHA-256 evidence, alongside Go 1.26.8 LICENSE/PATENTS, existing project/source-family notices, APRL source license and bundled rule/SKU provenance. No automatic legal classifications are inferred.

Nine isolated tests exercise exact notice bytes/hash and target accounting, missing/empty/symlink rejection, unreviewed replacement rejection, toolchain mismatch, failed-collection retention, stale-check retention and malformed JSON. Required Linux CI regenerates/checks freshness; existing required Windows validation remains in place. Local tests and repeat generation/check passed, with no go.mod/go.sum change. Raw third-party whitespace caused initial staged diff diagnostics; the shell did not stop its local checkpoint commit. FN-013 records this mistake and the explicit generated-notice byte-preservation attributes. The corrected staged check passes. Documented inspection-path recurrence in FN-004 and initial package-list VCS/stderr handling failure in FN-012. Code head `c91c50ced4f94ad81f74bc40a932e4deb9c2842c` passed both required jobs in PR #60 run `36742460046`, including the nine safeguards and exact inventory freshness check. Local `go mod verify` also passed. The final documentation/notice-attributes revision passed required run `36742917627`; PR #60 merged at `735a29894aeba8f24f8c56ec0518cce1b6f6d8c2`, with the merged tree matching the tested proposal.

This inventory is development evidence, not a complete package-file licensing audit, a formal SBOM, vulnerability clearance or release sign-off. All selected modules are covered conservatively, including test/tool graph modules not linked into the CLI. Additional target architectures, actual release binaries, bundled-file exceptions and distribution notice placement require release review. Gate 004 and the release gate remain open; no Azure calls or laptop action are required.

### Built CLI offline preflight and compiled dependency checkpoint

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `735a29894aeba8f24f8c56ec0518cce1b6f6d8c2`.

Added checks that execute the actual CLI built by required Linux/Windows jobs, copied into a temporary directory containing spaces. Four help/version cases check identity, documented output/redaction flags, empty stderr and absence of generated reports. Fifteen invalid-configuration cases check exit 1, precise stderr, empty stdout, no observed HTTP through a local proxy/identity tripwire, and unchanged directory/report hashes. Cases include unknown command/flag, positional arguments, malformed booleans, unknown/mandatory/deferred stages, bad stage parameters/impact, missing/malformed/invalid-scope filters and a valid scalar-tag filter followed by deferred-plugin rejection.

Build-info inspection requires the pinned Go version, expected main/package identity, host amd64 target, CGO disabled and exact compiled module versions/checksums against the selected inventory (24 Linux modules; 26 expected for Windows). The existing CI artifact builds now explicitly use CGO_ENABLED=0 to match that inventory configuration. No fixture executor, hidden scan mode or product test backdoor was added. These built-CLI cases do not claim successful scan/render behavior without Azure; the existing test-executable synthetic report/exit and coordinator-renderer fixtures remain independent evidence.

Local Linux suite passed after correcting a wrongly typed valid-tag fixture, recorded as FN-014. Toolchain-path, assumed CLI-package and shell sequencing recurrences are also logged. Local testing used the restored Go 1.26.8 toolchain with `-buildvcs=false` only for its build because local VCS stamping remains unresolved; production CI build stamping is unchanged. Required Linux/Windows run `36780706213` passed on code head `94992d058b0619f1a942c3abde48024024bb5ce9` in PR #61, including actual built-CLI execution on both hosts and exact 24/26-module matching. A local deliberately wrong expected-version run correctly failed the version assertion. The final evidence revision passed required run `36781015091`; PR #61 merged at `bd7156b6b26cc9d634c44e40eebe0af8e9240308`, with the merged tree matching the tested proposal. Printed artifact hashes are run diagnostics, not signed release provenance. No Azure/laptop action is required, and no release gate is passed.

### Development candidate packaging and isolated extraction checkpoint

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `bd7156b6b26cc9d634c44e40eebe0af8e9240308`.

Added a native Linux/Windows amd64 candidate ZIP packager consuming the actual CGO-disabled CLI. It requires trimpath, clean matching embedded VCS revision, pinned toolchain/module membership/checksums, matching CLI version and clean source checkout. Committed source bytes supply notices and inventory to avoid checkout line-ending changes. ZIP payload includes executable, complete collected notices/attribution, installation/rollback instructions, dependency inventory, build info and a candidate-status manifest with source commit/tree, source pins and file SHA-256 values. A sidecar records archive SHA-256. Exclusive same-filesystem hard-link publication rejects existing or concurrently published output paths. No GitHub Release, installer registration or automatic PATH modification is performed.

Twelve integrity/provenance fixture tests cover valid extraction/permissions, archive and payload corruption, missing notices, parent/absolute paths, case-folded duplicates, symlinks, excessive permissions, existing destination retention, unsupported target and incorrect compiled revision/dirty state/CGO/trimpath/dependency checksum. One actual-artifact integration case packages the same inputs twice, checks byte-equal ZIPs, extracts into a new directory with spaces, preserves original notice bytes, reruns the built-CLI offline suite and rejects output overwrite. The actual-artifact case is mandatory in both required CI jobs; local unit-only invocation explicitly reports it skipped.

A strict local build initially encountered the known workspace VCS-context issue. Verbose output identified the wrong root context; FN-012 records the follow-up. Candidate checks were kept strict and local integration uses an isolated Git clone. No cross-environment reproducible-build level, fresh-OS installation, signed provenance, complete license review or successful installed Azure scan is claimed. ZIP/sidecar publication is two file operations, not crash-atomic; consumers must require and verify both. The initial isolated-clone run exposed an import-time decorator error: it skipped integration despite --binary. FN-015 records that no actual package validation occurred in that run. The corrected method checks runtime arguments and must complete all thirteen cases with no integration skip. Corrected isolated-clone Linux execution completed all thirteen tests without skips on local checkpoint `069d493`, including strict source stamping and actual archive/extraction/runtime validation. The candidate SHA-256 was `09894713b6139463c5eb8c0bcb8afbf0b1ee15a194ecd4acc9fdfdc9493ee368`; this temporary fixture artifact was not published as a release. Local checkpoint `9036c3f` then passed all thirteen cases, including version-label rejection. Required Linux/Windows CI passed in PR #62 run `36783988010` on code head `1d465caeb6e76dd178cc7c0c858e89f83091f1d7`; inspected native job logs show thirteen cases and OK with no integration skip on both hosts. Final review pinned all source bytes/tree/timestamps to the captured commit SHA rather than repeatedly resolving HEAD. That final code/evidence revision must pass the native suites and required jobs before merge. Gate 004 and release approval remain open; no Azure/laptop action is required.

### Azure access model and operator recovery review

**Date:** 2026-09-30. **Starting branch:** `bootstrap/core-v1` at `ab0cc3b761956b5ea4712717b96a97d5a98f6567`.

Reconciled the roadmap overview with tested candidate packaging and the coarse 55% toolkit estimate (FN-016). PR #62's exact final head `bd852981de746cf481167f577e80a2e5061859c8` passed both required jobs in run `36784547096`; merge `ab0cc3b761956b5ea4712717b96a97d5a98f6567` has verified tree `588aeaca92e807ad50897586b095bb1f16447a8e`. This closes the prior checkpoint's final-CI prerequisite, not its release limits.

Added ACCESS_MODEL and OPERATIONS documentation after reviewing actual discovery/ARG/Diagnostics/Advisor/Defender/Policy/Arc/Cost requests, pinned credential chain, CLI flags, report and exit semantics against official Microsoft access guidance. The docs distinguish control-plane reads from local writes, stage-specific visibility and billing prerequisites, deliberate credentials/explicit scopes, sensitive reports, partial results and recovery. Cost is subscription-wide under resource filters. No Azure call, role assignment or provider registration was performed.

Source review confirmed current discovery silently omits requested IDs absent from visible listing; two focused characterization cases cover partly and wholly absent requested IDs. This is retained behavior, not a new completeness feature. ARG also documents partial visibility without a partial response signal. The canonical result's scope hash/count and returned-data IDs cannot prove intended full coverage, including empty subscriptions. Gate 004 now explicitly tracks independent intended/resolved visibility evidence and a future reporting/diagnostic contract. The target specification notes that classification against its unresolved-scope fatal expectation remains an acceptance decision; documentation does not waive that requirement. A certified minimum custom role, restricted-identity live tests, sovereign/custom-cloud validation, fresh-OS operation and production pilot remain open.

Validation: focused discovery tests passed locally with pinned Go 1.26.8; changed Markdown file links and staged whitespace checks passed. Both required jobs passed PR #63 code head `7193047bf56fbdfefeb1ff4ebfb2ff112d4dfb9f` in run `36789258708`. Final specification/evidence documentation must also pass both jobs on the exact final PR head before merge. No gate PASS is claimed.

### Eight-task offline project QA and credential-boundary hardening

**Date:** 2026-10-01 (Europe/Oslo). **Starting branch:** `bootstrap/core-v1` at `dd3cc8309b799502ba92c9d43559b776eb8b3918`. PR #63 final head `9e177fe98e2a7ea9e1221811fbf6e4d80058f50f` passed both required jobs in run `36789582586`; verified merged tree is `79196695bd8b66a3e3116916c0561bb8adc37413`.

Executed the agreed eight-task shortlist with evidence in OFFLINE_QA_REVIEW: reconciled documentation, mapped acceptance requirements to tests/open gaps, added one combined two-subscription/RG/tag/nearest-child fixture, checked operator flags against actual built CLI help and parsed the PowerShell snippet without execution, added authenticated Advisor/Cost request contracts, checked production default retry counts and sensitive SDK diagnostic output, and added incomplete ZIP/checksum publication rejection. Both required jobs now check authored documentation and native reachable vulnerabilities. Current package suite has thirteen integrity fixtures plus one mandatory native integration case.

The security review confirmed Advisor could send the ARM canary bearer header to a foreign absolute HTTPS nextLink in a synthetic transport. Initial desired regressions failed before remediation; no real token/external request was used. Advisor now validates the configured HTTPS host/port, rejects malformed/credential/fragment/protocol-relative links, detects repeated URLs and omits rejected URLs from error messages. Normal relative and same-origin public/government/custom-port fixtures pass. The default shared HTTP client also refuses redirects; two local TLS servers confirm zero forwarded requests. These are deliberate target hardening corrections, outside successful-response/source adversarial-parity claims. SDK-owned pagers, caller-injected transports and distinct unbounded continuations remain separate review boundaries. FN-017 records the defect/prevention; FN-018 records a corrected initial documentation-check error and pinned-SDK lookup recurrence.

Local focused and full race tests passed with Go 1.26.8. Final local documentation check passed 89 authored destinations, eight Cloud Assess flags and one PowerShell 7 snippet. Vet, module-tidy cleanliness and workflow YAML/job-preservation checks also passed. Unit packaging completed fourteen cases with integration explicitly skipped, not counted as package acceptance. Local govulncheck found zero reachable/imported-package vulnerabilities and three module-only x/crypto advisories (`GO-2026-6355`, `GO-2026-6354`, `GO-2026-5932`); unused affected packages and release limits are documented. An isolated clean Linux clone at local checkpoint `eadd9a29b55be73b760e85010a509ce53a547137` passed all fourteen native package cases with no skips, including interrupted checksum publication and extracted CLI checks. Its temporary candidate SHA-256 was `62b1427803b0c3769c6e9f0693ce83dd87469937c802f8335696d1ab2edbf18f`; no release artifact was published. Required Linux/Windows run `36792517368` passed PR #64 code head `b447ec46c3eee232259e7df3621227f840106451`. Inspected both native logs: fourteen package tests, no integration skips, successful 89-link/eight-flag/one-snippet checks, and zero reachable/imported-package findings with three module-only advisories. Linux statement coverage was 79.4%. That run predates a final SDK middleware correction: shared ARM options left automatic provider registration enabled. A desired synthetic read-failure regression reproduced one implicit POST after 409 MissingSubscriptionRegistration (two requests, one injected write rejection). Explicit DisableRPRegistration now blocks it; scope construction enforces the setting on copied options. One shared-pipeline and three real scope-pager fixtures check one GET/no POST, original error and unchanged caller options. Focused/full local race tests and vet were rerun after correction. FN-019 logs the missed middleware default; no historical Azure mutation is inferred. Updated final head must pass both required jobs before merge. Gate 004 and release security/operations approval remain open. No Azure/laptop task, role assignment, provider registration or release publication occurred.



Final implementation checkpoint: PR #64 head `a531237eae3bc8d4bec80bcfb6518698f1d61b9d`, including SDK registration suppression, passed required Linux/Windows run `36793880867`. Inspected both native logs: all fourteen package cases with no skips, 89 file destinations/eight flags/one snippet, and zero reachable/imported-package vulnerabilities. Linux aggregate statement coverage was 79.8%. The post-correction isolated clean Linux package suite also completed all fourteen cases at local checkpoint `dc0a04d` (temporary candidate SHA-256 `33f0be8de2dbf6c1af27004b8c2d55703d9d46a3db19af1d5aa601857965b194`). No artifact was published. Final evidence-documentation revision must pass both required jobs before merge; Gate 004/release decisions remain open.

### Independent audit of the eight unattended tasks

**Date:** 2026-10-01 (Europe/Oslo). **Audited revision:** merged `bootstrap/core-v1` at `6138e6e8dba4a3867c9e7b9c66c13ebaf8b2a79b`, tree `248e5e30468c0cf60a915ac355b3f15c4da8e4b2`. Final PR #64 head `e19470f17447fcc6573ccdbccf1c178de9588da6` passed required Linux/Windows run `36794408500`; both native logs were inspected during this audit.

Re-reviewed all eight task implementations, named tests and stated limits before starting scope reporting. An isolated clean clone with pinned APRL passed full race tests, vet, tidy cleanliness, real CLI/documentation/PowerShell checks, all fourteen native Linux package cases without skips, inventory safeguards/freshness, maintenance publisher/generator fixtures and existing comparator mutation controls. Isolated documentation controls correctly rejected broken destinations, invented flags and malformed PowerShell while accepting separately attributed Azure CLI flags. Two compiling behavioral mutations, removing Advisor origin enforcement and enabling SDK automatic registration, failed the existing security assertions; source was restored and the clean diff verified. Native Windows evidence comes from the final PR #64 job log, not a new local Windows execution. Native Linux govulncheck again found zero reachable/imported-package vulnerabilities and three module-only advisories.

The bounded code changes are confirmed by this audit, but documentation consistency needed correction: stale upcoming/completed wording and missing final run/merge indexing survived structural checks. FN-020 records the defect, recurrence and prevention. OFFLINE_QA_REVIEW contains per-task outcomes and unchanged acceptance limits. Roadmap and QA_PROCESS now distinguish tested local/native evidence from hosted/live/release gaps. No production code, source pins, comparator normalization or Azure scope was changed. Gate 004 and release approval remain open. The documentation correction PR must pass both required jobs on its exact final head; its final run/merge evidence is retained in the PR for the next ledger checkpoint, avoiding another evidence-only commit loop. Next implementation work remains intended-versus-resolved scope reporting.

### Requested/resolved scope reporting and explicit-request failure

**Date:** 2026-10-01 (Europe/Oslo). **Starting revision:** `7bbd1f661fcf49b824180599df493a13d31d78db`. PR #65 final head `b3bbf43e12fd847b013a38d9b0789a257cbcfc39` passed required Linux/Windows run `36799559903`, merged at this starting revision; reviewed/tested/merged tree `016083b8ad5832b00b7dbe7fce6d14aca375d035` was verified equal. That eight-task audit corrected documentation drift without production-code changes.

Defined and implemented the [scope-reporting contract](SCOPE_REPORTING.md). Canonical JSON adds selection mode, resolution status, explicit subscription/MG inputs, effective subscription filter lists, resolved IDs/display names (including zero-data subscriptions) and unresolved explicit IDs. Canonical construction clones/sorts arrays; JSON redaction collects every new subscription-ID list even without resource rows. Existing scopeId hashing, schemaVersion, CSV/XLSX layouts, source pins and recommendation/comparator projection remain unchanged. Strict downstream schemas must admit the additive scope object.

Selected fail-closed behavior for explicit CLI/library subscriptions in accordance with the target's fatal unresolved-request rule: after successful discovery, any missing explicit ID fails the critical scope stage before resource queries, records known resolved/unresolved scope, persists requested reports when possible and returns exit 1. Blank explicit scope IDs reject before discovery. Missing IDs are not automatically diagnosed as denied access; disabled/deleted state, tenant context or selection can cause absence. Low-level source-compatible intersection helpers and successful empty all-visible/MG/filter-only scopes remain unchanged. Existing include-over-exclude precedence is preserved. This is a deliberate target correction, logged as FN-021, not a new live-parity claim.

Offline checks cover partial/all missing requests, disabled/deleted subscriptions, duplicate/case normalization, zero-data resolved subscriptions, overlapping filters/intentional exclusions/filter-only missing includes, empty MG scope and denied listing. Application tests retain failed raw/redacted JSON with exit 1 and a resource-query tripwire. Scope array isolation/order and five distinct scope-only subscription-ID privacy cases also passed. Local pinned Go 1.26.8 full race tests and actual built CLI/help/PowerShell checks passed after submodule setup; local tests/help build use the known workspace buildvcs=false accommodation, while native CI keeps strict stamped package checks. Vet/tidy and a compiling guard-removal negative control are checked before publication. Exact final Linux/Windows CI must pass before merge; its run/merge evidence is retained in the PR for indexing at the next ledger checkpoint.

Microsoft ARG visibility guidance was rechecked against official documentation: resource results can omit inaccessible objects without a partial-result indication. Independent membership/RBAC reconciliation, restricted-identity live evidence, DV-001, historical warning/Advisor uncertainty, first-release scope and Gate 004/release approval remain open. No Azure scan, role change, provider registration or source-pin change occurred. Next autonomous review: SDK-owned pager/injected-transport credential boundaries, then pinned dependency advisory review.

### Project alignment and security sanity review, 2026-10-01

Paused normal work after PR #66 scope reporting and recorded the exact resume state in [ALIGNMENT_REVIEW.md](ALIGNMENT_REVIEW.md). PR #66 tested head `8422f78a1d247beb5baaaf1585bcfc9effea19e4` passed both required jobs in run `36810973320`, merged as `0df39bd9c3a8a0e2062f54772ac9e29750b3409b`, tree `9d3b0be415224bce256cc3e7d0070cdce3a901a2`. Previously inspected native evidence confirms fourteen package cases per host/no skips, 97 documentation destinations, eight flags/one snippet, 80.3% Linux coverage, zero reachable/imported-package vulnerabilities and three module-only advisories.

Rechecked functional alignment against pinned AZQR code, clarified normal recommendation corpora and missing feature surfaces, and retained configurable branding as unfinished. No feature deferral cancels eventual AZQR parity. A real SDK pager with an in-memory transport reproduced a foreign continuation carrying a synthetic bearer token; scope construction now validates HTTPS origin before authentication and refuses default redirects. A custom-cloud regression also proved shared HTTP adapters ignored the configured ARM token audience; ResourceManagerScope now honors audience separately from endpoint. These are documented target security/authentication corrections, with pre-fix failures and prevention checks in FN-022/FN-023. Healthy projection, rule pins, dependencies and comparator normalization are unchanged.

Focused Azure/discovery and full local race tests pass. Exact final required native CI remains blocking before merge; retain its run/head/tree/merge in the PR and index it at the next checkpoint. The new review contains a separate eight-item follow-up register and explicit inspection limits. Gate 004 and release acceptance are not passed. No Azure operation or reference-repository mutation occurred. Resume with pinned dependency advisory review; live tasks wait for a suitable environment without substituting production-containing Advisory.

### Dependency advisory review and bounded refresh, 2026-10-01

[PR #67](https://github.com/DeBoX85/Cloud-Assess/pull/67) tested head `08263690be3412c186b971d121c5e777c7c902dd` passed both required jobs in run `36813338081`; merge `8f76c46ad407f202a7aa2d9e8ae7eb8931ffaf35` matches reviewed tree `2b2f2eb80c800e1d4dfe79ccdd7dc4d15f47fbd4`. Inspected native logs confirm fourteen package cases/no skips per host, 107 documentation destinations/eight flags/one snippet, 80.5% Linux coverage, zero reachable/imported-package findings and three module-only advisories. This closes the alignment sanity-review checkpoint without declaring Gate 004 or release PASS.

[DEPENDENCY_ADVISORY_REVIEW.md](DEPENDENCY_ADVISORY_REVIEW.md) records AR-01 decisions. Updated only x/crypto v0.55.0 to first-fixed v0.56.0 for GO-2026-6355/6354, regenerated checksums and dependency/license evidence, and verified unchanged notice bytes and imported crypto package files. GO-2026-5932 remains explicitly module-only, with no fixed version and no OpenPGP imports in either supported CLI target. No vulnerability suppression, source/rule pin change or broader dependency refresh was made. Local full race and symbol/package scan pass; the latter now reports one module-only advisory. Exact-head required native CI remains blocking, with final run/head/tree/merge recorded in the PR and indexed at the next checkpoint. Next offline task: AR-02 pagination/lifecycle bounds. No Azure access is needed for this review.

### Scope pagination and opt-in assessment deadline, 2026-10-01

Indexed PR #68: final head `68150a9e405cdb09817bbbcc880a50ecca85a240` passed required run `36866099111`, merged as `a23c122a381aed0acd2f4e7daba9e6731e3540b7`, exact reviewed tree `305d8b2b34b4debf225c4ccc19083507b4578be0`. Native logs confirm fourteen package cases per host/no skips, 112 documentation destinations/eight flags/one snippet, 80.5% Linux coverage and zero reachable/imported-package findings with one residual module-only advisory. AR-01 bounded remediation is complete; release security remains open.

[Pagination/lifecycle evidence](PAGINATION_LIFECYCLE.md) records the AR-02 slice: reproduced one/two-link cycles in all three SDK scope pagers, now rejected without partial listings or continuation leakage. Advisor/SDK loops check context between pages. Added optional CLI `--assessment-timeout` and application budget with zero compatibility default, earlier-parent preservation, negative preflight rejection and cleanup before rendering. Stage invocation checks context before/after work, preserving canceled/deadline status and stopping later tasks even if a task returns nil. A real coordinator/application deadline fixture persists failed JSON with healthy prior data and exit 1; built CLI now checks seventeen preflight cases. FN-024 records the failure and regression controls. Healthy pagination/projection, pins/dependencies and normalizations are unchanged.

Local race/vet/inventory and built CLI checks pass; required native exact-head CI remains blocking before merge. Final CI/head/tree/merge is retained in the PR and indexed next checkpoint. This is cooperative assessment budgeting, not a process hard kill, render deadline, positive universal default or volume/load certification. AR-02 remains partial pending further bounds/calibration. No Azure or reference-repository operation occurred.

### Terminal pagination cancellation and service-contract follow-up, 2026-10-01

Indexed [PR #69](https://github.com/DeBoX85/Cloud-Assess/pull/69): final head `f04b6f417bd67360c4a05f461b484983f7f56290` passed both required jobs in run `36870843643`; merge `0a2dac1257cb99c390ad35eb927ac3e17171ba7d` matches reviewed tree `5760e4a6b3d0d1e814bd0abc0306ef4227f39201`. Inspected logs show fourteen native package cases per host/no skips, 120 documentation destinations/nine flags/one snippet, 80.9% Linux statement coverage and zero reachable/imported-package vulnerability findings with one residual module-only advisory. This evidence does not pass Gate 004 or release approval.

Follow-up [pagination review](PAGINATION_LIFECYCLE.md) independently reproduces successful terminal ARG/Advisor responses despite caller cancellation, now rejected before accepting response data. SDK scope's equivalent terminal fixtures already reject cancellation; no SDK production change is needed. Added distinct ARG token cancellation and healthy token reuse across batches/query invocations. FN-025 records the correction and execution recurrence. Overall stage cancellation already failed closed; this also protects direct adapter callers.

Kept source-compatible request/default behavior. A page-count flag lacks defined estate/load semantics and does not bound single-response memory; assessment timeout remains opt-in and AR-02 remains partial. Recorded matching ARG REST `$top` maximum 1000 versus source/target request 5000, plus omitted tokenless-truncation metadata, as explicit characterization tasks. Neither observation establishes a historical live mismatch. Normal Cost matches the pinned single-query scanner; region-plugin Cost paging remains separate unfinished parity work.

Final local race/vet/module/inventory and actual CLI/documentation checks passed (121 destinations/nine flags/one PowerShell snippet); exact final quality and windows-validation remain blocking. Final CI/head/tree/merge evidence stays in the PR for indexing at the next checkpoint. No Azure/laptop access, reference mutation, dependency/rule pin update or comparator normalization was used. Next offline task: ARG service-completeness characterization; branding/source-feature and deferred live/release work remain open.

### ARG page-size and explicit completeness corrections, 2026-10-01

Indexed [PR #70](https://github.com/DeBoX85/Cloud-Assess/pull/70): final head `70836a87580bf680ddb8da4f460329e46c553c25` passed both required jobs in run `36904496402`, merged as `b2494ca72b198d9bab6e36a8ff15b191e41a74dc`, exact reviewed tree `a6fe958663daffc37369e50f102521d7490126d9`. Inspected native logs confirm fourteen package cases per host/no skips, 121 documentation destinations/nine flags/one snippet, 80.9% Linux coverage and zero reachable/imported-package findings with one module-only advisory. Gate 004/release approval remains open.

[ARG_COMPLETENESS.md](ARG_COMPLETENESS.md) records matching REST 2024-04-01 contract and two deliberate target corrections: requested page size 1000 (source 5000), and failure rather than partial query success for explicit tokenless truncation. Pre-fix actual-decoder fixtures reproduce first/later-page partial success; a documented-limit fixture rejects the old request. Corrected 2501-row paging asserts every ID and stable query/scope. Valid empty/legacy-null metadata, true-with-token, invalid enum/type, empty token and nontruncated deliberate KQL limit are covered separately. Real decoder/client through inventory discovery or recommendation execution, coordinator/application and JSON rendering verifies failed completeness/exit 1, skipped later work, no incomplete findings and retained healthy prior inventory. FN-026 records failures and recurrence.

Actual pinned loader yields 560 definitions after precedence, with 344 executable non-empty query candidates; a textual operator screen found no limit/take/sample pipeline candidates in those queries or normal auxiliary query constants. This is not KQL type inference, execution, full catalog semantic validation or the count in every report. Missing/inconsistent metadata, visibility/snapshot ordering, memory/load/default policy and representative live large-result equivalence remain open. No historical Azure evidence is reclassified; pins, query text and comparator normalization are unchanged.

Final local race/vet/module/inventory, built CLI and documentation checks passed. A compiling truncation-guard removal failed both application regressions (exit 0 and Advisor ran); reviewed source was restored and focused race checks passed. Exact final quality/windows-validation remain blocking; final run/head/tree/merge stays in the PR for indexing at the next checkpoint. No Azure/laptop action was required. Next autonomous task: supported branding adjustment (AR-03), then source-feature characterization (AR-04); live/release limits remain explicit.

### Branding propagation and supported-workflow design review, 2026-10-01

Indexed [PR #71](https://github.com/DeBoX85/Cloud-Assess/pull/71): final head `acc5585c73848c92dae74c1a656b791c39bb3646` passed both required jobs in run `36907958956`, merged as `be95679ddea30fe164e67731f9c8535f3ff3368b`, exact reviewed tree `f494d7acf05cc57879d00b782719389b8483b6b9`. Inspected logs confirm fourteen native package cases per host/no skips, 134 authored destinations/nine flags/one snippet, 80.9% Linux coverage and zero reachable/imported-package vulnerability findings with one module-only advisory. This does not close large-result live evidence, Gate 004 or release approval.

[BRANDING_REVIEW.md](BRANDING_REVIEW.md) maps five active branding fields and four unused placeholders, neutral JSON/stdout/CSV data, stable SARIF integration IDs and package producer/verifier assumptions. Rebranding source defaults alone is not a supported coherent distribution workflow. Proposed default-preserving build profiles avoid repeated mutable runtime globals and source edits, with immutable decoding, filename/URL validation, real executable/report propagation, profile-aware package evidence and installation tests. Legal/provenance and assessment semantics remain separate from presentation. No profile format, command or custom package is claimed implemented.

The first PR #72 quality run failed duration-bounded comparator fuzzing at its budget boundary (FN-010 recurrence), with no saved assertion input. Preserve run `36925929037` and its retained artifact; change only that target to 100000 executions with separate 60-second timeout, as already used for decoder fuzzing, retaining assertions and requiring both jobs on the new head. The design review is complete, while AR-03 implementation/acceptance remains open. No Azure/laptop access or product input is required for its concrete next offline implementation batch. FN-027 records a read-path recurrence and research limits. Existing focused CLI/application/XLSX/SARIF race tests and documentation checks passed (139 destinations); exact-head native jobs remain blocking, with final run/head/tree/merge in the PR for next checkpoint indexing. Production code/defaults, pins, dependencies and normalizations are unchanged. Readiness estimates and deferred AR-02/live/release boundaries are unchanged.

### Immutable branding profile and native builder checkpoint

Date: 2026-10-01 (Europe/Oslo). Baseline merged PR #72: `c00b01489c7bd61fe0bbc9a9d6b31cad01b975f6`, tree `625ac49b38ff958932a1b5659c54a9ee1c60629f`; tested head `87ed0251ff634750eafa53cdf95bc3c02e701386`, run `36926967840`, quality job `110586574508` and Windows job `110586575052`, both passed with inspected logs. Both platforms ran all fourteen package cases without skips; 139 documentation destinations, nine flags and one PowerShell snippet checked. Linux coverage 80.9%; zero reachable/imported-package vulnerabilities, one previously recorded module-only advisory remains open. Merged tree matched reviewed tree; exact evidence retained in PR #72.

First two steps of AR-03 design batch now implement strict versioned JSON with five active fields, bounded/case-sensitive/duplicate/type/null/UTF-8/path/URL checks, canonical encoding, one-time immutable embedded resolution and a native Linux/Windows amd64 development builder. Builder uses explicit argv, encoded linker state, trimpath/readonly modules/strict VCS stamping/CGO disabled and exclusive output publication. CLI `branding` prints canonical public metadata without authentication. Malformed embedded state fails before CLI/auth/output work. Defaults, rule/dependency pins, canonical data schemas and comparator normalization are unchanged.

See [BRANDING_PROFILES.md](BRANDING_PROFILES.md). Both required jobs now exercise real default/custom builds and deliberately malformed embedded-state preflight under an HTTP tripwire, preservation of reports/executables, canonical metadata and help/version. FN-028 records null schema acceptance caught before publication and the local worktree stamping limitation. Final exact-head native CI and merge evidence remain blocking and are retained in the PR for next ledger checkpoint indexing.

Resume with actual custom report/data/provenance comparisons and profile-aware package/installation/tampering checks. Existing package safeguards still reject custom names; no workaround or release approval is provided. AR-03 is partial, Gate 004/release and AR-02/AR-04/live deferrals remain open. No Azure/laptop input is required for this next offline work.

### Paired immutable branding report acceptance

Date: 2026-10-02 (Europe/Oslo). Baseline PR #73 merged as `9cb3cfe0ff639d80e2cde8bbb8996208ed757adb`, tree `1d51ca05f6aa551b8fb68936c2901c697993127c`; tested head `ca687599bfd0142b8b3ce645123f8ce88213db84`, run `36930258812`, quality `110597484355` and Windows `110597483995` both passed with inspected logs. Each platform ran real default/custom/malformed builder/CLI checks and all fourteen package cases without skips; 146 documentation destinations, nine flags, one PowerShell snippet. Linux coverage 79.6% (previous 80.9%); native builder integration is separate from instrumented unit coverage. Zero reachable/imported-package vulnerability findings, unchanged one module-only advisory open. Reviewed/local/remote/merged trees matched; final evidence retained in PR #73.

Step 3 of AR-03 now has [paired report acceptance](BRANDING_REPORT_QA.md). Real production application/renderers consume deterministic synthetic non-empty data in two independently linked immutable profiles. All twelve XLSX sheets/tables, raw/redacted canonical JSON/stdout/CSV, synthetic APRL/AOR/CUSTOM labels/IDs/guidance, default timestamp prefix and explicit paths are checked. Independent XLSX ZIP/XML inspection permits title-cell differences only; a formula-like custom title must remain literal. SARIF permits only driver.name/informationUri differences, retaining all rules/results/fingerprints/automation identities. Four concurrent runs per profile preserve data and identity. Linux paired binaries use race instrumentation. HTTP tripwire execution and five changed-evidence controls complement existing actual CLI/build checks.

No production code/default, schema, pin, dependency or comparator normalization changes were needed. FN-029 records fixture/inspection preparation errors; their failed runs are not passing evidence. Local paired and focused race/vet checks passed before publication. Exact final native CI/head/tree/merge evidence remains blocking, retained in the PR for next ledger indexing.

Resume with profile/hash-aware native package producer/verifier, installation instructions and custom/tampering checks, preserving strict stamped clean source/inventory/notices/checksum/allowlist/privacy controls. Existing packaging still rejects custom names; AR-03 remains partial and Gate 004/release/live limitations unchanged. No Azure/laptop action is required for the next offline slice.

### Current-branch QA and profile-aware candidate package implementation

Date: 2026-10-02 (Europe/Oslo). Baseline PR #74 merged `5c0ec9e5b3d12340c7113c153df99db95f8df0e1`, tree `682b9355d6857a4c555bf61c518f1bbe522769aa`; tested head `8befe5e4fcbe85b625fda287d16e19ab04fbce88`, run `36933774218`, quality `110609105858`, Windows `110609106227`, both passed with inspected logs. Both ran paired reports, actual branding CLI checks and fourteen packages without skips; 155 documentation destinations, nine flags, one PowerShell snippet; Linux coverage 79.6%, zero reachable/imported vulnerability findings and unchanged one module-only advisory open. Merged tree matched reviewed tree.

Requested current-codebase/repository/documentation QA verified the live checkpoint, prior final native evidence, source pins and clean ordinary clone after runtime refresh. Full local race/vet, module/inventory consistency and authored links passed; no new production Go defect was confirmed. Current branding wording was reconciled while preserving chronological evidence. FN-030 records environment/read-path/helper-bytecode and exploratory skip limits.

[BRANDED_PACKAGES.md](BRANDED_PACKAGES.md) implements profile-aware candidate production/extraction: executable's resolved canonical profile, safe derived names, canonical profile/hash in manifest schema 2, coherent generated installation instructions and strict duplicate/type/layout/hash checks before extraction writes. Original source/target/CGO/trimpath/dependency/notice/checksum/allowlist/privacy/publication controls remain. Extraction never executes archived binaries; producer and controlled installed QA use trusted build inputs. Version-1 development manifests deliberately use their historical verifier; canonical assessment schemas are unchanged.

Both native jobs run all original fourteen cases plus five branding controls against actual default and actual custom executables, with installed profile/help/version/preflight/inventory and non-empty report acceptance already protected. Custom Unicode/quote/HTML text and URL ampersands cross the real Go/Python canonical boundary. No production Go, pin, dependency or comparator normalization changes. Final local/native exact-head results and merge remain blocking, retained in the PR for next ledger indexing. Do not mark AR-03 accepted before both final jobs pass.

After bounded five-field workflow acceptance, resume AR-04 pinned-source external YAML/KQL and missing CLI/plugin feature characterization. Release/Gate 004, fresh OS/signed provenance, SARIF consuming-service migration and deferred Azure evidence remain open. No Azure/laptop input is required for source characterization.

### Pinned-source feature characterization and accepted branding checkpoint

Date: 2026-10-02 (Europe/Oslo). Baseline PR #75 final head `c5748a4c549c5b925a4d1b51394e155be300aabf`, run `36940456214`: quality job `110630692098`, Windows job `110630691802`, both passed and logs inspected. Each platform ran nineteen default plus nineteen custom package cases without skips, paired production reports, installed CLI/profile/preflight/inventory checks and six actual-snippet failure controls. Documentation checked 167 destinations/nine flags/one PowerShell snippet; Linux coverage 79.6%; zero reachable/imported vulnerabilities and unchanged one module-only advisory open. Merge `9a022a62444b0940b3f016dc58e850b1599c6fa6` matched reviewed tree `b5fb6bbe615e5aef73c4e3349c5bca38505ce6d5`. Bounded five-field immutable development workflow AR-03 is accepted, with runtime/theme/logo/release boundaries unchanged.

[FEATURE_PARITY_CHARACTERIZATION.md](FEATURE_PARITY_CHARACTERIZATION.md) now records AR-04 pinned-source entry-point contracts and concrete implementation slices. Source YAML plugins join normal Graph recommendations automatically; internal selected plugins return separate tables. Existing target filesystem catalog loading does not implement the source plugin object schema. Source scanner commands live under scan, plugin commands are top-level plugin-only runs, and source rules output selects scanner-supported embedded plus Diagnostics definitions rather than dumping the raw catalog. Six plugins and compare/SKU/MCP remain functional parity obligations.

The inspection identified source-only caveats to handle explicitly: logged/skipped plugin failures with successful stage return, queryFile paths lacking containment/size controls, registry replacement of colliding built-ins, external queries bypassing embedded disabled/marker filters, missing guidance indexed without a guard, and region cost-history-months declared but not forwarded in CLI dispatch. These are not newly introduced target defects and this checkpoint changes no execution behavior. Each implementation slice must verify the contract and document deliberate corrections; detailed per-plugin helper/query/endpoint semantics remain work, particularly region selection.

No production code, dependencies, pins, normalization or Azure requests changed. Source HEAD was verified at `8e4f0577f3615e6c9014c031bcad079f235369cc`; reference checkout remains read-only. Documentation/link and source-fact checks precede exact-head required native CI; final head/tree/run/merge retained in the PR and indexed next checkpoint. FN-031 records inspection recurrences and the corrected early plugin-output assumption. Resume at offline rules inspection, then scanner subcommands, bounded YAML integration and plugin table/health infrastructure. No Azure/laptop input required. Gate 004, release, readiness estimates and deferred live evidence remain open.

### Offline rules command and independent reference executable capture

Date: 2026-10-02 (Europe/Oslo). Baseline PR #76 merged `fd60dd635acb7c3156a264640810a4bd8d92b59d`, reviewed/merged tree `14e08250648ec3977b4c55d517eca7d9cc440047`; final head `ef7fd526410bb10f0ea9c080b9ad7f8e20e260d9`, run `36941892097`, quality job `110635076067`, Windows job `110635075872`, both passed with inspected logs. Each platform ran nineteen default and nineteen custom package cases without skips; documentation 170 destinations/nine flags/one PowerShell snippet; Linux coverage 79.6%, bounded fuzz/mutation/race/vet passed; zero reachable/imported findings and unchanged one module-only advisory open. This indexes completed source characterization, not implemented plugin parity.

[Offline rules inspection](RULES_INSPECTION.md) implements AR-04 slice 1: source-selected scanner-supported embedded and Diagnostics recommendations, global ID replacement, sorted six-field JSON, Markdown output and local json/j flag without Azure/scan execution. The independent full pinned AZQR executable capture has 380 rows and 107 catalog types; target JSON matches its bytes exactly. Fixture hash `77d7574232a32062d87dc32f6f9fc216f6f509921cc4d75621aa3a0c5279ed08` is retained with provenance/notices. Existing reference checkout lacked APRL files, so a separate clean temporary source copy populated both exact pins for the executable capture; original source and target submodule checkouts were unchanged. FN-032 records failed preparation runs, not accepted evidence.

Literal selection/precedence/empty-guidance fixtures, complete reference-byte command tests, Markdown table escaping, invalid arguments and writer-error tests supplement actual built and installed CLI checks with HTTP/auth tripwires, all source JSON rows, exit/stderr validation and unchanged filesystem snapshots. A compiling removal of the selection predicate was rejected by independent expected rows; restored code passed focused race tests. No scan behavior, dependencies, source pins, canonical assessment schema or equivalence normalization changed, and no Azure requests occurred. Markdown formatting/safe missing guidance and reuse of the target case-insensitive marker guard are explicit limits/corrections; current pinned JSON still matches fully.

Local complete tests and exact-head required native QA precede merge; final head/tree/run/merge evidence remains in the PR for next checkpoint indexing. Existing required jobs now protect actual rules execution including default/custom package installs; another numbered gate is unnecessary. Resume at scanner-specific subcommands, then YAML Graph integration and internal plugin table/health infrastructure. No Azure/laptop input required; Gate 004/release, readiness estimates and live deferrals remain open.

### Development execution plan and scanner-specific commands

Date: 2026-10-02 (Europe/Oslo). PR #77 rules acceptance: head `d41f0d67a7c3c938ce6d5bf03e4e4c649be42651`, tree `f2c99952554e11042d90ff0dd0597bda07b70e4d`, run `36944379293`, quality `110642961367`, Windows `110642961131`; both passed with inspected logs. Merge `c79c2ffc9bfcec66da9f24dab631eb6bef6d1836` matches the reviewed tree and verified GitHub identity. Each host ran nineteen default and nineteen custom package cases without skips; four actual CLI checks; documentation 181 destinations/nine flags/one PowerShell snippet; Linux coverage 79.9%; race/vet/bounded fuzz/four comparator mutations passed. Zero reachable/imported vulnerabilities; unchanged one module-only advisory remains open.

[DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md) records B1-B7 dependencies/acceptance/fallbacks, per-slice alignment records, security and honest evidence boundaries, mandatory exact-head native QA, durable proposal checkpoints, unknown-outcome GitHub reconciliation, interruption recovery and reviewed revert rollback. This reduces risk without guaranteeing zero defects, durable scratch or unavailable live validation. Scope/pins/dependencies/gates are not silently relaxed; Azure/laptop tasks stay explicitly deferred. Independent-person review is only claimed when actually performed, never substituted by self-review wording.

[SCANNER_COMMANDS.md](SCANNER_COMMANDS.md) implements B1: independently captured 87 pinned AZQR command keys, inherited scan flags, copied single-key selection and production request mapping. Existing generic orchestration handles filtering/stages/output. Literal CLI-to-coordinator storage-specific versus generic VM-filter fixtures verify executed definitions and included/excluded inventory; disabled operations are guarded. All keys, extra arguments, exit errors, repeated-root selection isolation and installed help/preflight tripwires are checked. Removing request key mapping compiled but failed the expected definition assertion; restored focused race checks passed. FN-033 records preparation errors and recurrence prevention.

No dependencies, pins, assessment schemas, normalization, generic defaults or Azure requests changed. Exact final-head native run/merge evidence is retained in the PR and indexed at the next material checkpoint. After B1 acceptance, resume B2 bounded source-schema YAML Graph integration, then canonical plugin table/health and individual migrations. Gate 004/release, progress estimates and deferred live evidence remain open. No laptop or Azure input is currently needed.

### Bounded YAML Graph execution and uninterrupted development authority

Date: 2026-10-02 (Europe/Oslo). B1 final acceptance: PR #78 head `c3907d252b1d9cb48079efe618db1e55028a47e2`, run `36947056761`, quality job `110651469372`, Windows job `110651469128`, both SUCCESS with inspected logs. Merge `50fbaee905e7be7a7c99e174f35ea69247eb16bb`, tree `19da2f901ccd93af27add84da1b93bb13d12a4cc`, agrees with local/uploaded/preview contents and expected parents/identity; fetched working tree was clean. Each host executed five actual CLI checks and nineteen default plus nineteen custom package cases without skips, documentation 199 destinations/nine flags/one PowerShell snippet, twelve paired XLSX/CSV report tables. Linux coverage 79.9%; full race/vet/bounded fuzz/four comparator mutations and maintenance controls passed. Windows failure propagation passed six controls. Zero reachable/imported vulnerabilities; unchanged one module-only advisory stays open.

The user authorized continued execution through the agreed plan without routine continuation prompts. [DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md) records that scope and its genuine input/blocker/major-flaw boundaries. Laptop/Azure acceptance remains deferred, and no background operation beyond an active session is promised. Lost scratch was recovered from the accepted remote checkpoint; original reference remains unchanged.

B2 [YAML_GRAPH_PLUGINS.md](YAML_GRAPH_PLUGINS.md) records the pre-implementation source contract, limits/trust decisions, targeted files, acceptance and reviewed rollback. Distinct source-object parsing/discovery preserves home/current precedence, first duplicate name, sorted same-type-ID override, queryFile precedence, configured attribution and normal Graph enablement. Input byte/count/depth bounds, strict ambiguity/schema rejection, reserved names and confined regular UTF-8 query-file access are explicit target protections. Safe Unicode/space labels remain supported. Existing service-identity scans do not require a resolvable home; no substitute relative home is invented.

The unchanged pinned AZQR loader executed the literal full-field fixture in a separate temporary source checkout. Source capture SHA-256 `dbd0ba34d2aff68c8ff8dd28eb353e2906085fa43821bea24becd3188856772c` normalizes only temporary CommandPath. Complete mapped output matches target conversion, including the source recommendationTypeId omission. Literal parser/precedence/metadata/portable path/bounds tests supplement it. Catalog overlays copy caller metadata and never mutate the base; external origin remains private and uses the source plugin exclusion predicate rather than embedded disabled/development filtering.

Synthetic file-to-ARG transport serialization/decoder/executor-to-coordinator-to-JSON checks pass for healthy, excluded and denied queries, including persisted failed Graph health and absent-guidance safety. Actual built/installed CLI preflight checks were added for current/home default/custom discovery, malformed/schema/parent paths, unchanged prior reports and no observed HTTP/auth traffic; rules/help remain offline. Critical origin-removal and raised-file-limit controls compiled and failed named independent assertions; restored focused race checks passed. FN-034 records preparation/review corrections. Full final-head native acceptance and post-merge evidence remain required and are retained in the PR for the next material checkpoint.

No dependency, source pin, canonical assessment schema, normalization, internal-plugin stage enablement or Azure request changed during development. Next B3 implements real zone-mapping with canonical internal-plugin table/health/report infrastructure, then a second plugin checks shared design. Gate 004/release, estimates and live/deferred limitations remain open. No laptop/Azure input is currently required.

### B3a: Bounded zone-mapping adapter and source corrections

**Date:** 2026-10-02 (Europe/Oslo).

**Accepted B2 evidence:** PR #79 tested head `6842c50c6aa12ec10a7472483a6b528fc6938635`, tree `6d209444dbb59ace9570736a33bbd6cc8b0ff818`, native run `36963458729`, Linux quality job `110701991878` and Windows validation job `110701991724`, both successful with logs inspected. Merged `7cddf9368272dd425315bdf3635b21fb974a9fe4`; both parents, reviewed/preview/merged trees, GitHub identity and clean fetched state verified. Six actual CLI checks, nineteen default plus nineteen custom package cases on each host without acceptance skips, 208 documentation destinations/nine flags/one PowerShell snippet, race/vet/fuzz/mutation and vulnerability checks passed. Coverage 80.4% is a signal, not parity proof; unchanged one module-only advisory remains open.

**B3a scope:** [ZONE_MAPPING.md](ZONE_MAPPING.md) records the pre-implementation source/API contract and limits. The new internal zone adapter and bounded HTTP GET preserve source row values and sorting with deterministic tie breakers. The pinned source ignores documented continuation links and logs/skips subscription failures while returning success. Target corrections safely follow same-origin/scope/version links, retain healthy earlier pages/other subscriptions and return sanitized failure codes. Ambiguous/missing/null response envelopes and malformed identities fail explicitly. Selected cloud origin/audience remain a caller responsibility until production wiring.

**Security/ownership:** Five context-aware workers at most; bounded scope/pages/bytes/rows, deterministic aggregate admission, no global plugin cache and no borrowed mutable caller/result state. New bounded GET closes bodies and maintains the operation deadline. Its opt-in transport response bound applies before SDK retry/authentication policies consume attempt bodies; ordinary GET/POST calls retain existing behavior. Unsafe continuations are rejected before the authenticated getter. Redirects remain disabled on the production transport; injected code is trusted, not sandboxed.

**Evidence and QA:** Actual unchanged pinned source parser executed in a separate copy; complete six-field capture/hash and reproduction method recorded in the contract. Original reference remained clean. Literal fixtures plus synthetic authenticated pipeline cover healthy/empty/denied/malformed responses, partial later-page preservation, boundaries, cancellation, five-worker scheduling, sorting, repeated/concurrent isolation and selected-cloud token/request contracts. Compiling cross-origin and retry-body-limit negative controls failed their named assertions; restored focused race tests passed. Full local QA and exact-final-head mandatory native jobs must succeed before merge; final run/head/tree/merge evidence belongs to the PR and next material ledger checkpoint. FN-035 records repeated source/target read-path mistakes.

**Limits/resume:** This is adapter implementation, not a new exposed command or canonical schema/report change. Plugin stage availability, default scans, dependencies, source pins and equivalence normalization are unchanged. B3 remains in progress: canonical table/health/privacy/report integration, actual mixed/plugin-only CLI, honest registry list/info and a second plugin verification remain required. Live plugin validation, DV-001 and other environment-dependent acceptance remain deferred. Gate 004/release, estimates and module-only advisory stay open. No laptop/Azure input required and no Azure request made during development. Recovery baseline is PR #79 merge above; publication/rollback follow the execution plan.

### B3b: Canonical plugin-table ownership, health and report infrastructure

**Date:** 2026-10-02 (Europe/Oslo).

**Accepted B3a evidence:** PR #80 tested head `63bb1c08b5de838af291a591d45bcbef9c060489`, tree `1e6bed9f4f19b4848a6f35a9d1098c4bfccb2370`, run `36966029765`, Linux quality `110709854293` and Windows validation `110709854116`, both successful with complete logs inspected. Merged `0b76ad6b2870432357e5a38c1e8fe73d5214f11b`; exact reviewed/preview/merged trees, both parents, identity and remote/fetched state verified. Native six CLI checks, nineteen default plus nineteen custom package cases per host without skips, 214 documentation links/nine flags/one PowerShell snippet and existing provenance/module/inventory/race/vet/fuzz/mutation/maintenance/vulnerability controls passed. Total coverage 81.2%; zone package 96.4% is a regression signal, not live/feature proof. Existing module-only advisory remains open.

**Contract and change:** [PLUGIN_TABLES.md](PLUGIN_TABLES.md) was written before implementation. Ordinary canonical assessments retain schema 1.0, JSON fields and twelve report tables. Explicit plugin-bearing construction validates/deep-copies versioned table metadata/columns/rows/cells/health and uses assessment schema 1.1 with additive pluginTables. Builder folds failed/skipped health into partial completeness and warnings into complete-with-warnings while retaining critical failure. All renderers reject malformed table extensions or manually concealed completeness before mutation. Stable safe IDs/CSV suffixes, sheet collision/UTF-16 limits, row widths, status/count consistency, text/count bounds and source row order are enforced without silent truncation.

**Report/privacy:** JSON/stdout retain complete canonical tables and existing whole-result masking, now with explicit row identities. Shared CSV/Excel projection masks known/contextual subscription identities in arbitrary plugin text and status errors. Explicit failed/skipped empty tables remain visible with health on Assessment Status. CSV formula protections and Excel literal text cells remain intact. SARIF remains primary findings only, without invented plugin rule findings. Scope filtering, plugin stage/CLI wiring and real second-plugin execution remain separate acceptance work.

**QA/evidence:** Literal source zone/service-health schemas check distinct five/six-column shapes and source metadata; the service-health fixture does not certify its unmigrated adapter. Complete JSON/stdout table comparison, every plugin CSV/XLSX header/data cell, XLSX formula absence, partial application exit, empty/failed/skipped health, ownership/default compatibility and reject-before-replacement checks are synthetic/offline. Privacy, width and false-complete compiling mutations failed named assertions and were restored. Pinned Excelize source inspection exposed draft sheet Unicode-count and Unicode-fold collision mismatches, both reproduced before correction and covered by rejection, exact 31-unit actual-workbook and collision boundaries; FN-036 records them plus the fixture suffix and patch-preparation mistakes. Restored full race/vet/local checks and actual exact-final-head native acceptance are required; final evidence belongs to the PR and next checkpoint.

**Resume/limits:** No dependencies, pins, equivalence normalization, default scan availability or Azure requests changed. This is library/report infrastructure, not completed B3 or advertised internal-plugin CLI execution. Next wire zone through coordinator and actual mixed/plugin-only command selection with health/exit/discovery acceptance, then validate a second real adapter. Gate 004/release, estimates, module-only advisory and DV-001/live/laptop deferrals remain open. Recovery/rollback baseline is accepted PR #80 above; follow protected expected-head merge and reviewed revert rules. No laptop/Azure input required now.

### Fresh-session handover and outage recovery reconciliation

**Date:** 2026-10-02 (Europe/Oslo).

PR81 canonical plugin-table infrastructure accepted on head `8b1365218743d5b34fb934d24354699bf55465e2`, native run `36968213639`, Linux quality `110716475829` and Windows `110716476011`; both succeeded with logs inspected. Merge `634aa083981fc327fe5c2cc7b42815793eb197fd`, tree `4729d848509dccf13f221359add14a1ef3cf8edc`, expected parents verified. Both hosts executed six CLI and nineteen default plus nineteen custom package cases without acceptance skips; coverage 78.0% above unchanged 75% floor. Existing module-only advisory remains open.

A runtime interruption blocked unpublished B3c work. Later inventory confirmed its scratch workspace was absent; the uncertain final patch could not be inspected. No B3c implementation is accepted. PR82 checkpoint `2d413cca968784bfb7613490283bb933e3653894` preserves contract/evidence/reconstruction instructions, not exact code. Recovery rechecked PR79-81 native outcomes/logs, accepted head/tree/parents, clean original source and source-capture hashes. Fresh accepted-baseline Linux race/vet/strict build/six actual CLI tests and documentation checks (223 destinations/nine flags/one PowerShell snippet) passed after a classified missing-submodule setup correction. This was not a new Windows run, comprehensive audit, live validation or release gate.

[SESSION_HANDOVER.md](SESSION_HANDOVER.md) is the new current entry point: purpose/authority, pins/access, implemented versus pending features, exact next task, dependency-ordered backlog, security/live/release limits, toolchain/QA, verified code checkpoints and paste-ready new-session prompt. Root AGENTS.md requires proactive updates after coherent slices and material state changes, plus verified remote code/documentation checkpoints before switching tasks or ending sessions when available. README/current roadmap, execution/QA plans and FN-037 were reconciled; [ZONE_EXECUTION_CHECKPOINT.md](ZONE_EXECUTION_CHECKPOINT.md) retains PR82 recovery content. Source-grounded handover completeness review and documentation checks are required before publication; final native evidence for this documentation proposal belongs in its PR and next material index.

**Resume:** Reconstruct B3c zone coordinator first, then mixed/plugin-only CLI/list-info and actual integration/installed checks; verify against a second real adapter. No Azure/laptop input is needed. Source pins/normalization/dependencies/runtime behavior unchanged; Gate 004/release and deferred validation remain open. Rollback this documentation slice through reviewed revert after dependency review, never force/reset history or reference.

### Reconstructed zone coordinator and durable code checkpoint

**Date:** 2026-10-02. PR83 handover final head `ed081d5fff7e24f40bb44cf3927dea482ec7c112`, run `37060126509`, Linux/Windows jobs `111014451610` / `111014451903` passed with logs inspected; merged `d88f9294482eafc4e724dbe2248fc09ed4f826b3`, tree `27e0e50dc86d1b970c6ddb94cb33ac0dec857f2a`, expected parents/current remote/fetched clean state verified. PR82 recovery content retained and proposal closed without B3c code acceptance.

[ZONE_EXECUTION.md](ZONE_EXECUTION.md) defines reconstructed B3c coordinator scope: owned stage/options/name/subscription state, exact zone selection, plugin-only regular-operation skip, optional lazy operation, source metadata and validated owned tables, retained healthy rows, sanitized partial failures, context failure and requested pending headers after critical scope failure. Existing unselected schema/defaults/CLI availability are unchanged. Source pins/dependencies/normalization unchanged; no Azure requests.

Real-parser synthetic coordinator fixtures and independent literal projection tests passed focused/full local race and vet. A compiling inventory-isolation control failed its named assertion and was restored. FN-038 records read/output preparation recurrence and corrected distinct source sheet description. Actual authenticated CLI/application/installed integration remains next; adapter's prior authenticated tests are not claimed as that closure. This coherent code/documentation proposal is remotely checkpointed before more implementation accumulates; native final acceptance belongs in its PR and next material index. SESSION_HANDOVER/roadmap updated. Gate 004/release, second real adapter and live deferrals remain open.

### Zone CLI and offline registry reconstruction

**Date:** 2026-10-02. Coordinator PR84 accepted on head `88d7452df8942341513cc62a49d9df24a32e167e`, run `37061985911`, Linux quality `111020561369` and Windows `111020561056`, both actually successful with logs inspected; merge `6a920caa76576e61708717a647d67124ebb6ede6`, tree `85365fd1ab7440b4f4d19ce335f188bcbd9f6983`, identity/expected parents/current ref/fetched clean state verified. Both hosts executed six CLI and nineteen default plus nineteen custom package cases without acceptance skips, documentation 257 destinations/nine flags/one PowerShell snippet. Advisory/live/release limits remain open.

[ZONE_EXECUTION.md](ZONE_EXECUTION.md) extends the contract to shared scan/zone flags and preflight, exact internal selection, source-compatible plugin-only scope boundary, unused Graph YAML skip and honest offline list/info. YAML source path is discovery metadata; recommendation definitions/source projection unchanged. Registry text controls, stable sorted JSON metadata/capabilities and strict no-partial-list errors are explicit improvements. Source pins/dependencies/equivalence normalization unchanged; no Azure calls.

Fourteen raw/masked actual Cobra/shared-preflight/coordinator/application cases use real zone scanner and authenticated synthetic HTTP, testing selected audience/read boundary, rows, health/counts, JSON/CSV/XLSX, partial/cancel exits and persisted artifacts. Full local race/vet and named compiling unused-discovery control passed after FN-039 fixture/path corrections. Clean stamped actual CLI and native/default/custom installed acceptance are required before merge; successful live credentials/factory/Azure execution are not claimed by injection fixtures. Current handover/roadmap and failure records updated and code proposal checkpointed before more work. Next service-health second real adapter after CLI acceptance; Gate004/release and live deferrals remain open.

### Zone CLI acceptance and service-health migration checkpoint (2026-10-02)

PR85 accepted zone command, named mixed selection, bounded plugin list/info and unused-YAML isolation. Tested head `ac85fcb4754c96fe8b0046840cb05e8bb2a39d65`; tree `d6b2f5a477e454d27d21aa3bfae2a7c26e256f3e`; run `37064059199`, Linux/Windows jobs `111027374427` / `111027374648` passed, actual logs inspected. Eight built CLI, nineteen default plus nineteen custom package cases per host; 258 documentation destinations/nine flags/one PowerShell snippet, race/vet/fuzz/mutation/provenance/inventory/maintenance, 78.3% coverage above unchanged 75% floor. Existing module-only advisory remains open; no reachable/imported findings. Merge `2f88e9d5b8dca309c0f27ccb6ec4ed54787bcebb`, exact tree/ordered parents/current ref and clean fetched task branch verified. Final [PR85](https://github.com/DeBoX85/Cloud-Assess/pull/85) retains acceptance/rollback. No live zone equivalence claimed.

Next: [service-health](SERVICE_HEALTH.md), separate library adapter before public integration. The real library scanner/shared bounded POST are now locally implemented; source capture/hash, full restored race/vet and compiled page-limit control passed. Native proposal acceptance is pending, not advertised as an executable adapter. Source inspection confirmed no date cutoff despite 90-day description/2,160-hour denominator; preserve and document, do not silently invent a temporal correction. Current specification unavailability prose is being reconciled with accepted zone behavior; full internal-plugin parity remains open.

### Service-health library acceptance and integration start (2026-10-02)

PR86 tested head `e38db1400a6f06ac4fd837f24f0e4fc9074f7a28`, tree `b818c79829c50609ff55e7e698bcf462b2fb0b36`, run `37066186931`, Linux/Windows jobs `111034429914` / `111034430064` passed; actual logs inspected. Merge `a6b790ce143bc70fac8f4aa3199225ee92b0660d`, exact tree/ordered parents/current ref and clean fetched state verified. Eight built CLI and nineteen default plus nineteen custom package cases per host without acceptance skips; 264 documentation destinations/nine flags/one PowerShell snippet, 78.9% coverage, required safety/provenance/race/vet/fuzz/mutation/maintenance scans passed; existing module-only advisory remains open. Final [PR86](https://github.com/DeBoX85/Cloud-Assess/pull/86) retains evidence and rollback. Library-only acceptance does not expose the adapter or establish live parity.

Next: SERVICE_HEALTH integration contract on that accepted merge, real command/coordinator/request/report tests before public availability. FN041 records command-name prose error and path-inventory recurrence; actual command remains `zone-mapping`.

### Service-health execution candidate and accidental-stop reconciliation (2026-10-02)

Reverified core-v1 at PR86 merge `a6b790ce143bc70fac8f4aa3199225ee92b0660d`, PR85/86 exact successful native runs and both job results, preserved uncommitted integration diff and unchanged pinned AZQR checkout. No uncertain remote mutation was retried. [SERVICE_HEALTH.md](SERVICE_HEALTH.md) records the implementation/request/filter/report contract and evidence limits; SESSION_HANDOVER records the resume point.

Implemented source-ordered service-health/zone dispatch, optional captured-cloud operation, owned scope/type filters, strict/sanitized table projection, standalone and mixed/scanner commands and honest registry. Eighteen synthetic authenticated command-to-application report cases passed; review adds repeated/concurrent differing scope/type isolation and stricter operation-boundary label checks matching the scanner. Restored full race/vet and compiling metadata-boundary control passed: removing exact metadata equality compiled and failed the named source_metadata projection assertion. Clean actual executable/documentation and exact native jobs remain required before acceptance. The candidate is not accepted live evidence. Existing gate/release/temporal-query and laptop deferrals remain open. FN041 records fixture/prose corrections; no source pin/dependency/normalization changes.

### Service-health acceptance and SQL EOL contract (2026-10-02)

PR87 tested head `6c544ea989d0948bb158027ebfc8c110c626b0a2`, tree `8a15549723926129fe151f6fe94515178aba3167`, native run `37069137247`, Linux/Windows jobs `111044153941` / `111044153636` passed; actual logs inspected. Protected merge `c892227a440f431363acd626add999e32c152c0e` and exact tree/ordered parents/identity/current ref/clean fetched state verified. Nine actual CLI and nineteen default plus nineteen custom package cases per host without acceptance skips; 265 documentation destinations, 79.1% coverage. Required QA passed; zero reachable/imported vulnerabilities and existing module-only advisory still open. Final [PR87](https://github.com/DeBoX85/Cloud-Assess/pull/87) retains evidence/rollback. Service-health offline execution is accepted, live equivalence deferred.

Next [SQL_EOL.md](SQL_EOL.md) records actual source metadata/query/32 columns, subscription-only filtering and source string/null/ordering behavior. Actual unchanged source Scan produced nonempty/filtered/empty synthetic captures through the temporary TLS transport seam. Full 611-line query inspected and hash recorded; query-engine arithmetic/clock and pricing/licensing currency were not executed or certified. Preserve source estimates with explicit limits rather than infer current authoritative costs. Original source pin unchanged. Implement library separately before public integration.

### SQL EOL library candidate and request-boundary correction (2026-10-03)

[SQL_EOL.md](SQL_EOL.md) now retains complete actual source captures/query and a bounded real library scanner with strict source string/null projection, subscription scope/filtering, unchanged received ordering, truthful malformed/partial/context health and owned tables. Focused race QA passed after FN042 fixture type/tag corrections. Compiling subscription-filter omission failed its named independently captured filtered-table assertion; restored full/built/native checks are required before acceptance. Public SQL availability remains unchanged; no engine arithmetic/live pricing/credential parity claim.

Endpoint review reproduced accepted service-health and SQL draft acceptance of an empty fragment delimiter: parsed Fragment was empty while raw URL concatenation misplaced the ARG path. Both constructors now reject the delimiter before transport, with failing-before/passing-after regressions. FN043 records the narrow robustness correction and limits; no Azure call, demonstrated mutation or token leakage. Source pins/dependencies/normalization/gates unchanged. Handover, roadmap and failure records updated with exact next checks.

### SQL EOL library acceptance and execution start (2026-10-03)

PR88 exact head `58ece1ff1d6ec515151f5f664aa6c65444f70ee7`, tree `00d0478e8005c16036e62c94cddc0ae3be25649d`, run `37070978399`, Linux/Windows jobs `111050085480` / `111050085721` passed; actual logs inspected. Protected merge `f6dcbba032785ed86af63bb018cf3c2a0fcb204a` exact tree/ordered parents/human identity/current ref/clean fetched state verified. Nine built CLI and nineteen default plus nineteen custom package cases per host without acceptance skips; 269 documentation destinations/nine flags/one PowerShell snippet, coverage 79.5%. Required QA passed; prior module-only advisory/Gate004/release/live evidence remain open.

Next [SQL_EOL.md](SQL_EOL.md) execution contract governs standalone/mixed/scanner and all-adapter selection, source subscription-only filtering, owned/sanitized projection, captured cloud/options and every 32-column raw/masked report cell. SQL public availability remains unaccepted until separate exact-head native gates. No laptop/Azure input needed for this slice. Final PR88 retains evidence/rollback; previous library pending prose is chronological, not missing acceptance.

### SQL EOL execution candidate and interruption reconciliation (2026-10-03)

Reverified accepted PR87/88 refs, identity, exact trees, native jobs/logs and source pin; no ambiguous pending remote mutation or lost integration patch. The integration preserves the full source projection and subscription-only semantics through standalone/normal/scanner dispatch and source-ordered three-adapter selection. Strict owned metadata/table/health projection sanitizes provider failures and keeps ordinary inventory/other plugin rows; critical/cancel paths retain requested pending tables. Twenty-four authenticated synthetic raw/masked command/application cases compare every JSON/CSV/XLSX cell, independently assert negative-text CSV escaping and partial artifact-before-exit. Repeated/concurrent scope/filter/result isolation and pre-auth unavailable/unused options are covered.

Focused and restored full race/vet passed. Compiling metadata equality removal failed the named assertion and was restored. Ten actual built CLI checks added; clean strict build/docs and exact native acceptance remain pending. FN044 records fixture/read corrections. No Azure calls, dependency/pin/query/model/renderer change or live financial/licensing evidence. Protected reviewed revert can remove this integration while preserving the accepted PR88 library. Final run/tree/merge evidence will be kept in the proposal and indexed at the next checkpoint. Resume native acceptance, then B5 carbon; all existing live/Gate004/release deferrals remain open.

### SQL execution acceptance and carbon source/core start (2026-10-03)

PR89 exact head `fbed3340fa873fcee82ebf4131e49e94c4c3896b`, tree `179bbabe28ec5e6ab861037cb9aa19bb85d33b16`, run `37073129430`, Linux/Windows jobs `111056952522` / `111056952419` actual PASS, logs/steps inspected. Ten built CLI, nineteen default plus nineteen custom package cases on each host without acceptance skips; 272 documentation destinations/nine flags/one PowerShell snippet; 79.5% coverage. Mandatory race/vet/fuzz/mutations/provenance/branding/helpers/maintenance/packaging and zero reachable/imported vulnerability checks passed; prior module-only advisory stays open. Ruleset/current base/head/merge-preview verified. Protected merge `735cedfe68a9b91a17ec120bafd320edb8629d1e` has exact candidate tree/ordered parents/human author/GitHub merge committer; current ref/clean fetched state verified. SQL execution VERIFIED OFFLINE. Earlier pending status superseded, not live/model/release approval.

Next [CARBON_EMISSIONS.md](CARBON_EMISSIONS.md) records actual unchanged source scanner captures and pure aggregation core. Full metadata/eight cells for selected/all/101-subscription two-batch/empty/denied plus six actual arithmetic cases captured through synthetic TLS/token and temporary source transport/filter fixtures. Source ignores skip token/access metadata and hides denied report as healthy empty; target request adapter must classify these separately. Pure core preserves source sums/float formatting/conditional blanks/service latest date, with explicit finite/date/label/count/context guards and owned tables. No HTTP or public carbon availability yet. FN045 records setup/recurrence corrections; focused race passed. Compiling source-calculation negative control/restored full/build/native acceptance remain required before core acceptance, then bounded requests and separate execution/report integration. No Azure/pin/dependency/SQL change; live/Gate004/release deferrals retained.

### Carbon pure core acceptance and request-adapter handover (2026-10-03)

PR90 tested head `9cbe9569a85413b2db2c4459e5cba07c7791c7c2`, tree `dd4f7b39b27cdf1f7c3782ffbc141599871185e2`, run `37074676704`, Linux/Windows jobs `111061797849` / `111061797618` actual PASS; logs/steps reviewed. Ten CLI and nineteen default plus nineteen custom package cases per host without acceptance skips; 274 Markdown destinations/nine flags/one parsed PowerShell snippet; 79.8% coverage. Mandatory QA passed and both scans found zero reachable/imported vulnerabilities, with the prior module-only advisory still open. Protected merge `ea0aabfca055899e18baa54af24b53df90a57d98` exact candidate tree/ordered parents/current ref/human author/GitHub merge committer/clean fetched state verified.

[CARBON_EMISSIONS.md](CARBON_EMISSIONS.md) and handover now index accepted source arithmetic/core separately from the unimplemented request adapter and public execution. Reviewed official 2025-04-01 response/SDK/change-log access metadata and documented exact next bounded request/decoder/pagination/access/privacy/context/independent-negative-control checks. Official role wording differs across guides; actual scoped live access evidence remains deferred and no role changes are authorized. No new adapter code, Azure call, pin/dependency change or live/release closure in this checkpoint. Continue offline request library, then separate public execution; preserve all existing limits. Final PR90 retains exact acceptance and rollback.

### Final handover acceptance and workspace outage (2026-10-03)

PR91 final head `b591ce58ab842ef450d8862990b2ccea12fe9fb4`, tree `69bf1f743351cc4a137eb682dfc9c4df50f4f63e`, run `37075702689`, Linux/Windows `111065002687` / `111065002846` actual PASS/logs inspected. Protected merge `8f06a667a583933134f695efa414886d0af26f29` exact tree/ordered parents/current remote ref verified. Ten CLI and nineteen default/nineteen custom package cases per host; documentation 276/9/1, 79.8% coverage; mandatory QA/zero reachable-imported vulnerabilities passed, prior advisory remains open. Stale roadmap baseline and unqualified old estimates were corrected (FN046), with exact final-head QA rerun.

Final local postmerge fetch/readback stalled; follow-up terminal reported environment_offline (409), while GitHub stayed available. Stopped uncertain dependent sequence; PR91 merged state and missing final PR-note write were read back. Local final fetch/switch/source-last-check is unverified. Last known local tree was clean/equal to accepted documentation; no carbon request implementation had started, all code is remote. This outage checkpoint updates handover/failure notes through GitHub and remains a proposal until native acceptance. FN047/048 record cause/limits/content-shape correction. Restore workspace and inspect actual state before local retry, then bounded carbon request library. No Azure/laptop input, no live/gate/release closure or automatic-background recovery claim.

### Restored workspace QA and bounded carbon request candidate (2026-10-03)

PR92 exact head `79b6c3735a9ee608cec7f0c396319d0eab181465`, run `37076773339`, Linux/Windows `111068331737` / `111068331617` passed; full native steps/logs reviewed. Protected merge `62352ec9fc66abc4cc5ae3ab265ea984d5f6a27f`, tree `9c74c7e45f213dbb20161bd9f953be4a4370c729`, ordered parents/ref verified and completed clean local fetch/switch confirmed. Fresh separate clone recovered from the remote checkpoint after older scratch state returned; retained existing checkouts. AZQR clean/pinned and actual APRL checkout verified. Fresh full Linux race/vet passed. Native docs 276/9/1 and 79.8% coverage; zero reachable/imported advisory findings with prior module-only advisory open. No accepted code was lost, no Azure action or release/gate closure.

Unaccepted B5 carbon request library candidate adds exact source/REST date/report bodies, bounded response/JSON/access/pagination, immutable endpoints/audience, source arithmetic and honest empty/denied/malformed/partial/cancellation behavior. Independent source captures and authenticated synthetic source-cell/body/batch/closure/retry/budget/privacy/concurrent tests pass focused race after FN049 corrections. CARBON_EMISSIONS and handover define exact limits/rollback/next checks. Compiling denial mutation failed its named assertion; restored full race/vet, strict clean readonly/VCS-stamped CLI build, ten actual CLI tests and docs 276/9/1 passed. Code/docs WIP commit `708e94f697f0646f2cf00bb5f3a87b30ed93049c`/tree `b17beeb1c1b6776032b6a5720ba4dbc9b1811256`, parent PR92 merge and human identity verified remotely. Final head native Linux/Windows remains mandatory before acceptance. Current specification branding/plugin summaries were reconciled to existing accepted evidence (FN046 recurrence), without changing feature scope. Public plugin execution is separate; live totals/access and existing deferrals remain open.

### Carbon request library acceptance and public execution recovery (2026-10-03)

PR93 accepted bounded date/report/access/pagination library: tested head `92f235a199f96d0695230f83dcb6f4e3024bda5e`, tree `d831ebe613cb29a63e70e66b07351becd32efcd6`, native run `37081937599`, Linux/Windows jobs `111084150124` / `111084149828` actual PASS with required steps/full logs inspected. Ten actual CLI, nineteen default/nineteen custom package cases per host, docs 276/9/1, coverage 80.6%, zero reachable/imported findings; module-only advisory remains open. Protected merge `f5cae2b515481f8215b483e6c6792a85ff8c1b62`, exact tree, ordered parents (`62352ec9fc66abc4cc5ae3ab265ea984d5f6a27f`, `92f235a199f96d0695230f83dcb6f4e3024bda5e`), current ref/identity and clean fetched checkout verified. Source/APRL pins unchanged.

After interruption, live accepted ref and no open proposals matched local base; uncommitted public carbon candidate was recovered and reviewed. Added focused source-captured all-cell JSON/CSV/XLSX authenticated command tests, four-plugin dispatch, cancellation with retained earlier sums, access warnings/denial, ordinary inventory preservation, filter/scope/result ownership/concurrency, critical pending metadata and missing-operation preflight. Focused/full local race/vet, compiling isolated correlation-guard negative control, strict clean build, eleven actual CLI checks and docs 276/9/1 passed; final native/protected acceptance remains pending. Verified remote WIP e64ac74af8eb75fee91d9a407439f34a624bffb1 retains code/docs tree 374ad6eb694035a755c68ee9988e2a0592096392 and accepted PR93 parent/human identity. Candidate errors and read recurrence are recorded in FN050. Current command remains unaccepted until exact-head gates and protected merge. See CARBON_EMISSIONS and SESSION_HANDOVER for contract, next step and rollback. Live service totals/access/sovereign execution, Gate004 and release remain deferred; no Azure resources/roles were changed.

### Carbon public execution accepted; local executor recovery blocked (2026-10-03)

PR94 tested head `0ee0f59f1fea4de3690b47bdd8f514e4a0ee9b6b`, tree `80f7734522e7c7f00f7c0992cb73bb5b1a22a266`, native run `37083743138`, Linux quality `111089620311` / windows-validation `111089620449` actual PASS. Required steps/full logs inspected: eleven real CLI, nineteen default/nineteen custom package cases per host without acceptance skips, docs278/9/1, coverage80.7%, mandatory provenance/format/module/inventory/branding/PowerShell/race/fuzz/maintenance/mutation/vet checks passed. Zero reachable/imported findings; existing one module-only advisory stays open. Protected merge `ac32f3b116ef26704a06add686ea305d41c2a114`, exact tree and ordered parents (`f5cae2b515481f8215b483e6c6792a85ff8c1b62`, `0ee0f59f1fea4de3690b47bdd8f514e4a0ee9b6b`), current remote ref/merged PR/human author/GitHub committer verified. Source/APRL pins and captures unchanged.

Local postmerge fetch/switch verification is BLOCKED by an execution-service transport disconnect. The combined fetch/switch command completion is UNKNOWN. A bounded read-only pwd retry stalled and was terminated; no further code edit was attempted. Do not infer a clean checkout, current local branch or successful fetch from remote merge success. GitHub remains available and this code/docs checkpoint is remotely preserved. No Azure/laptop/credential action is needed from the operator.

Resume: restore executor, inspect actual task-owned checkout/status/branches/processes before replaying anything, preserve unrelated work, fetch live core-v1 and verify expected head/tree/ordered parents/pins, reconcile the possibly-created feat/ai-governance-core branch without reset/force, complete clean fetched verification. Re-read current handover/PR94/this recovery proposal and active rules. Then begin independent B5 AI governance source captures/pure projection, followed by bounded authenticated request library and separate public execution. Region selection follows separately; B6/B7/live/Gate004/release remain open. No speculative runtime change while local verification is unavailable.

### Post-carbon outage recovered and final documentation review (2026-10-03)

Executor became available with the recent checkout missing. A separate clean clone matched PR94 merge ac32f3b116ef26704a06add686ea305d41c2a114/tree80f7734522e7c7f00f7c0992cb73bb5b1a22a266, ordered parents and unchanged source/APRL pins; actual clean fetch/read-back completed. The earlier interrupted command's completion is still unknown. PR95 initial documentation head c7d0f0bfb93f70f95a2db4b7c6fd315eac7c6254 passed native run37084489534/jobs111091835643+111091835942; all required steps/full logs inspected, eleven actual CLI and nineteen default/nineteen custom package cases per host without acceptance skips, docs278/9/1, coverage80.7%, zero reachable/imported findings with one module-only advisory open. Fresh restored Linux full race/vet/strict build/eleven CLI/docs QA and six source-capture hashes passed.

Final semantic review corrected stale supplemental carbon unimplemented wording in the target specification and reconciled current recovery status. Updated PR95 is unaccepted pending fresh exact-head native gates, protected merge and postmerge verification. FN046/FN051 records recurrence/recovery; preserve historical outage snapshots. Next B5 AI governance source captures/pure projection remains independent of laptop/Azure access. No runtime, source/pin/dependency/equivalence or Azure mutation in this documentation fix.

### PR95 recovery accepted; AI governance pure core started (2026-10-03)

PR95 head `b066d19515bc866efd544f94e5a883ddaca133fc`, tree `44beaa7443a7ff5cef0189682032c412ac5d65a4`: final native run `37086007396`, quality `111096349779` / windows-validation `111096349892` actual PASS, required steps/full logs inspected, eleven actual CLI/nineteen default/nineteen custom package cases per host without acceptance skips, docs278/9/1, coverage80.7%, mandatory provenance/module/inventory/branding/PowerShell/race/fuzz/maintenance/mutation/vet and scans passed. Zero reachable/imported findings; existing one module-only advisory open. Active rules/current base/final head/preview tree/parents checked; protected expected-head merge `f183755ef2e7ee9c61c555dbaa1743f3f801c457` matches candidate tree and ordered parents `ac32f3b116ef26704a06add686ea305d41c2a114` / `b066d19515bc866efd544f94e5a883ddaca133fc`. Remote identity/ref and completed clean fetched/switch head/tree/parents/APRL/source verified.

B5 AI governance next: actual source Scan/parser/enrichment captures on an isolated unchanged source/APRL copy, exact SDK API/audience inspection, bounded pure aggregation/projection and complete metadata/header/cell/health/identity/limits/cancellation/ownership tests. [AI_GOVERNANCE.md](AI_GOVERNANCE.md) retains capture hashes, regeneration, reviewed source/target differences and exact next request/public obligations. Focused/full race and vet passed; compiling removed-correlation mutation failed the named assertions. Stamped CLI/docs and final native acceptance pending. Runtime registry unchanged, no AI HTTP/discovery/decoder/report execution yet. FN052 records confirmed first capture escape (one read-oriented public metrics query with fake token, HTTP401), stopped harness, failed-closed preparation correction and successful explicit SDK-owned loopback client with destination tripwire. Failed attempts do not establish live evidence. No real credential/resource writes or changed target dependency/source pin.

### AI core native Windows capture-byte finding (2026-10-03)

PR96 first candidate `6b8bf43582a01e73edf57add5a0c775d4f0ad1ae`, tree `0d29fab1129e47e192dae6f655dbf430c7901b70`, run `37087387963`: Windows job `111100395209` failed TestSourceCaptureHashes because new AI JSON fixtures lacked byte-preserving Git attributes. Correct with -text for all six input/output files, preserve their original hashes/cells, reproduce isolated autocrlf=true checkout, and obtain fresh exact-head native quality/windows-validation. FN053 retains the cause/correction/prevention; no skipped or earlier-head result certifies the corrected candidate.

### AI pure core accepted; bounded requests started (2026-10-03)

PR96 final head `cf7e14dff431a095a06271ab560fd1934c947f6a`, tree `b032e44ef99aec0b2214e65ace3f15194c842a72`, native run `37087685383`, quality `111101249440` / windows-validation `111101249557` actual PASS with required steps/full logs inspected. Eleven actual CLI/nineteen default/nineteen custom package cases per host without acceptance skips, docs285/9/1, coverage80.9%, mandatory checks/source cells/hashes/correlation mutation passed. Zero reachable findings; existing module-only advisory remains open. Active rules/current base/head/preview tree/ordered parents verified; protected merge `9c6d0e79d05d06559fb4b543406b928013aed835` matches tree, parents `f183755ef2e7ee9c61c555dbaa1743f3f801c457` / `cf7e14dff431a095a06271ab560fd1934c947f6a`, remote human identity/ref and completed clean fetched/switch verification. Original source/APRL pins/dependencies unchanged.

Next/current [AI_GOVERNANCE_REQUESTS.md](AI_GOVERNANCE_REQUESTS.md): bounded wire decoding and metrics/deployment retrieval for already discovered accounts. No Graph/CLI/report integration or public AI availability in the pure milestone; those require separate acceptance. Primary source SDK/captures establish exact API versions/audiences/window/schema, official common metrics docs corroborate batching/error meaning while their documented2023 version does not replace pinned2024-02-01. Initial production metrics contract is public cloud; different/custom ARM cloud fails before token acquisition pending separate service validation. No laptop input or live evidence claimed.

Request-slice interruption recovery, 2026-10-03: actual live/local accepted PR96 base verified; earlier focused session completed zero. Unaccepted bounded decoder/request code and explicit in-memory authenticated tests are now present. Source cells and unsafe continuation-before-authentication are asserted; no network or real token is used by target tests. Candidate compilation/fixture corrections and enrichment-failure retention repair recorded FN054; aggregate/security/full/native acceptance remains pending. Initial request WIP PR97 head da88735c087db33c29f04c09e7e81232895f184d/tree a7e9ca23e2ab2469ef7f7333dafd0e84889b2a50 remotely preserves code/tests/docs with exact accepted parent and human identity read back. Follow-up aggregate enrichment retention/global entry/byte/five-worker/later-batch/context fixtures pass focused race; three compiling mutations fail their named assertions and restored fixtures pass; bounded5000x decoder fuzz passes. Recurring Linux/Windows mutation checks and Linux fuzz are added without weakening existing gates. Final local full race/vet/strict built CLI, eleven cases, docs290/9/1, three guard mutations/restored fixtures and5000x fuzz completed zero. Candidate c4dfec1a101adfa9a488cb375edac1f38a0a2bd4/tree20a7f2e197ac67e2f87891ca07508631c6792eea identity/parent verified. Semantic review reconciles stale current plan/specification wording with accepted PR95/96 and unaccepted PR97; fresh native QA on this final documentation head/protected acceptance remains required. Public ai-gov unavailable; pins/dependencies unchanged; no live/release approval.

## Resumed AI discovery contract and executor blocker (2026-10-03)

Live bootstrap/core-v1 and PR98 reverified at `82dea5624764332153248a4d2c39b653e562c02f`, tree5f2edcfc3f3a3977fe888e6cbca0b1ee37dd9d1d, ordered parentsaadd0d03b13a093a9ec3b9e69e3f70208d1e1714 / e6950f0b0043ee3cb273715efabfe18f6c1a69a3, human author/GitHub merge committer. No open proposals at resumed start. PR98 exact-head run37093219699 passed both mandatory native jobs; its separate accepted-merge [run37093451765](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37093451765) also passed quality111118472235 and Windows111118472069 on exact82dea562. Required steps/fetched log evidence: eleven actual CLI/nineteen default/nineteen custom package tests per host, docs296/9/1, coverage81.4%, four compiling AI request mutations/restored fixtures, existing race/vet/fuzz/provenance/module/inventory/branding/maintenance/comparator and reachable scans passed. Existing module-only advisory stays open; only conditional failure-artifact upload skipped.

Independent planning refined [AI_GOVERNANCE_EXECUTION.md](AI_GOVERNANCE_EXECUTION.md): 4,096 scope/raw-row/account caps; 300-subscription/1,000-row pages; 64 total logical calls/32 pages per batch; 2MiB attempt/16MiB successful bodies and retained text; 4,096-byte opaque tokens; five-minute discovery context. Reviewed the pinned source query/field mapping/filter plus actual target ARG/filter/AI bounds. Official 2024-04-01 schema at0799fda68aab1e2c6f5f0d06f028f451b14bf0a8, blob9de6392a2cc1f1cf8a2423d3c48ed63b1fd2e4a1, requires integer count/totalRecords, string truncation enum and data. Actual official examples have false-with-token, so continuation is independent of that flag. Contract defines envelope/page admission, row correlation/unknown tags, conservative changed-total failure, sanitized owned-prefix retention and independent literal/compiling-control acceptance. This is proposed target safety behavior, not code, test PASS, live snapshot/capacity proof or a source-pin update.

FN056 recovery remains BLOCKED. Resumed local path/tool inventory returned Windows cwd and bundled Git, no callable Go; later read-only commands disconnected and recovery timed out. A subsequent recovery attempt stalled and was terminated, without inferred completion. Cloud Environment skill is available but no status/execution capability is callable; discovery found no available matching connection. No external workspace was provisioned or permissions changed. Local fetched/source/toolchain/process proof remains unavailable; GitHub/native CI are separate usable capabilities. No implementation was started or lost. Maintain handover/roadmap/failure record and publish this planning checkpoint with exact read-back/native acceptance; retain final run/tree/merge proof in its PR without self-reference commits.

Resume after usable executor restoration: inspect actual isolated state, fetch live core-v1 and unchanged AZQR/APRL, verify toolchain/pins/trees/parents/proposals, then implement discovery while ai-gov stays unavailable. Separate actual public/report integration follows. Operator laptop/Azure evidence, Gate004, release, live roles/models/load/sovereign validation and previous forensic/advisory boundaries remain open. No Azure request/resource write.

## Executor rebuilt and reproducible development workspace (2026-10-03)

Accepted PR99 baseline89087cf1d8afb7e06abfb134cd5d068cfa56186d/tree1088200d83e163f9b5b65128ac1282b37e7ac655 and ordered parents82dea562 /9878571, human identity and final native/merge-push evidence verified; no open proposal at resumed start. User explicitly authorized creating the required development workspace/tools before feature work. New Ubuntu24.04 executor responds, Git2.51.1 has HTTPS, publisher/repository/module/vulnerability-database access works. Installed task-owned exact Go1.26.8, PowerShell7.6.6 with published archive SHA256 verification and govulncheck1.8.0. Existing Python3.12.14/ripgrep/GCC13.3.0/Node24.19.0 support tooling; no system profile/proxy/credential/pin/dependency change.

Fresh task-owned target is clean at accepted89087cf/tree1088200d, unchanged AZQR8e4f057/tree17d93b20 and both actual APRL60eaddda; target AOR/CUSTOM/SKU pins also verified. Both module graphs downloaded/verified without go.mod/go.sum changes. Source AI package compiled with no tests run. Pristine target full race/vet, readonly CGO-disabled stamped CLI (revision89087cf, modified=false, hash4aec5698bc0c29af1d0a32c3c2ded2d81b565d03a191e6bd015ccdad24408401), eleven actual CLI, docs297/9/1, inventory/PowerShell/branding/report checks, nineteen default/nineteen branded packages, coverage81.4%, four AI/comparator controls/restored baselines, bounded fuzz, maintenance and reachable/imported vulnerability checks passed; prior module-only advisory remains open.

Added [bootstrap-workspace.sh](../scripts/bootstrap-workspace.sh), [boundary checks](../scripts/tests/bootstrap-workspace.py), required Linux syntax/checksum/preservation feedback and [DEVELOPMENT_WORKSPACE.md](DEVELOPMENT_WORKSPACE.md). Actual script recreated a second path containing spaces with fresh checksum-verified downloads, clean clones/submodules and fresh verified module caches; manifest records observed refs and no QA claim. Boundary fixtures reject corrupt bytes before extraction/cloning and preserve an existing destination; a valid-shell disabled-checksum control fails its named extraction assertion and restored guards pass. Documentation/handover/roadmap/FN056 now distinguish this actual local recovery from old outage snapshots. Runtime/source/capture/dependency/branding behavior is unchanged.

Publish this coherent code/docs tooling slice, verify tree/parents/identity/bytes, require independent exact-head native Linux/Windows/protected acceptance and clean fetched merge verification; final PR evidence supplies self-referential run/merge details. Then proceed to bounded AI Graph discovery while public ai-gov stays unavailable, followed by separate actual execution/report integration. No Azure scan/write, new source capture or live/Gate004/release closure. Scratch may disappear; remote recipe/checkpoints enable reconstruction, not guaranteed executor persistence.

PR100 initial native run37103069229 passed Windows111146251132 but failed Linux111146251036 in the new checksum fixture before the expected checksum assertion. The fixture checked the host prerequisite inventory before exercising its corrupt archive; its assertion initially omitted captured stderr, so the specific early prerequisite failure is not retrospectively proven. Corrected offline fixtures supply private inert rg/gcc/Git-HTTPS prerequisite shims, reject unexpected Git operations, require the existing-directory mkdir failure and retain early diagnostics. Real end-to-end provisioning above remains independent evidence. Three corrected cases, a valid-shell removed-checksum control and restored guard pass locally; final exact-head native acceptance remains required. Automated review questioned attribution, but GitHub commit API and fetched object independently show candidatec906a3 authored and committed by the mandated human. Shell push lacks credentials; connector publication is verified by exact tree and ten file readbacks, without requesting tokens.

## Bounded AI discovery candidate after accepted workspace recovery (2026-10-03)

PR100 accepted ca7fc301de421f24f62c6b976fe2d79c0bfe6780/tree019abda8, ordered parents89087cf/9cc061a, human author/GitHub merge committer and clean local core-v1/source/APRL proof verified. Final run37103412975 passed quality111147208821/Windows111147208735; separate accepted push37103646003 passed quality111147875379/Windows111147875270. Required steps/full logs inspected: eleven CLI/nineteen default/nineteen branded cases on both hosts, docs311/9/1, coverage81.4%, new bootstrap guards, old four AI/comparator controls and broad existing QA. Known module-only advisory remains open. No executor input remains necessary.

The current discovery candidate implements the exact pinned source query and official2024 schema, owned sorted UUID scope, fixed ARM audience/origin/objectArray/top/batches, strict metadata/coverage/continuations/row identity, recorded AssessmentFilter decisions, run budgets and owned partial results with truthful health. The shared production cloud guard is factored without changing metrics behavior. No normal ARG/public CLI/report/dependency/source/capture changes. Primary schema blob9de6392a was read again; Microsoft Learn corroborates paging, while the immutable schema's actual string enum and false-with-token examples govern wire behavior.

Literal mapping, actual public construction/client token/request/body closure, partial/unknown clouds before credentials,301-subscription replay, foreign batch before filter, invalid envelope/siblings/duplicates, actual tag/structural filter, two-page/two-batch retained failures/cancellation, sequential/concurrent ownership and exact/excess scope/page/request/raw-row/body/decoded-text limits pass focused race. Three new compiling discovery faults and four previous request faults are rejected by named assertions and restored baselines;5000x discovery fuzz passed. Final full local race/vet and eight compiling controls/restored baselines passed; clean stamped CLI/package/docs/native QA remains separately required on the published final candidate. During self-review, rejected projected strings were initially not charged to text; corrected before publication with an independently sized16MiB/+1-byte fixture (FN057).

Publish coherent code/tests/docs on an isolated proposal, read back tree/parents/identity/every changed byte and require exact-head native/protected acceptance and clean fetched/source proof. Public ai-gov stays unavailable; its enclosing context and standalone/mixed/scanner/all-format health/report integration is next after discovery. Existing Azure/laptop/roles/models/sovereign/load/DV001/advisory/Gate004/release boundaries stay deferred/open. No Azure scan/write, new source capture or independent-person review claimed.

## Bounded AI discovery final acceptance and public integration resume (2026-10-03)

PR101 final2258f029ed7661a77958a9b47e2ea274ae500d1e/treeb302bc664317e4b08908ab25d45eaff3d2e5aefe passed required run37105082689, quality111151894191/Windows111151894026. All required steps/full logs inspected; preview1ab4b8c2 matches tree and ordered baseca7fc301/final2258.11 actual CLI/19 default/19 branded packages per host, docs312/9/1,81.9% coverage, eight compiling controls/restored assertions, new5000x discovery fuzz and existing broad QA passed. Fresh local exact2258 full race/vet/fuzz/stamped CLI (modified=false, hashc4f00db50931bf88a838d0d247d7b7059a5d21e4a565f68c85de4de06f83e97c),11 CLI/docs passed. Earlier87d963f native PASS is historical, not correction acceptance. FN057/058 fixes and the survived first Unicode control are retained.

Expected-head protected merge accepted3c4019c88c067818f1c4d72f754fc514e5712296/treeb302bc66/ordered parentsca7fc301/2258; human author/GitHub merge committer/live ref and clean local fast-forward/source8e4f057/tree17d93b20/both actual APRL60eaddda verified. Every13 changed file and local/remote tree had exact readback. Final accepted-merge push is distinct, pending its own step/log verification at this read; final PR101 evidence supplies that subsequent result. No pins/captures/dependencies/normal ARG/public availability changed; module-only advisory remains open.

This material acceptance handover supersedes candidate-pending/current-blocker text and resumes actual public AI standalone/mixed/scanner/registry/report integration. Enforce public-cloud preflight before selected-path scope authentication, owned stable filters/enclosing lifecycle and discovery-failure propagation after enrichment. All sixteen source cells, AI Throttling empty/AI Gov nonempty, pending/failed headers, actual raw/masked all-format reports and earlier public/default behavior need independent command-to-coordinator fixtures and new final native/protected acceptance. Azure/laptop/live/roles/models/sovereign/load/DV001/Gate004/release remain deferred/open; no independent-person review or background guarantee.

PR102 review caught one active execution-plan paragraph still saying final discovery acceptance remained required despite the accepted-state header. Corrected that current-tense contradiction, updated all eight final controls and explicitly marked earlier FN058 feedback historical. FN059 records recurrence/prevention; initial successful native head92c1fd is not acceptance of this revised navigation text. Also verified accepted PR101's separate push run37105398968: quality111152783968/Windows111152783809, every required step and fetched logs passed with11 CLI/19 default/19 branded cases per host, docs312/9/1,81.9% coverage/eight controls/broad QA. Current PR101 body retains distinct evidence.

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
- live YAML/KQL extension validation on a suitable environment
- internal plugin migration/parity
- remaining internal adapters and live plugin execution (bounded list/info implemented through PR85)
- packaging/distribution
- release-specific dependency/license inventory review and distribution notices
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

Discovery pre-acceptance review additionally rejects Unicode case-fold aliases in regional DNS input and fixed ARM service path/type names, preserving safe display labels and ASCII casing. Distinct-ID Kelvin-sign/long-s fixtures cannot be masked by duplicate detection; an eighth compiling control disables the DNS ASCII guard and must fail its named assertion (FN058). Earlier candidate full local/native results do not certify this correction; final revised exact-head QA remains required.


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


## PR105 review correction and observed recovery (2026-10-03)

Initial source characterization37131392522/job111227049328 and required native37131392461/quality111227049234/windows111227049053 passed on7b1fdede. Review found path filters could skip the promised final-head source job on follow-up commits. Removed both path filters; fresh three-job final-head evidence and protected acceptance remain required (FN066). Recovered executor/source/tools in an isolated exact remote worktree, preserving old dirty corrections and unknown command completion. No auxiliary runtime, public region, source/dependency pin or Azure/live/Gate004/release change. Exact next step is publish correction, verify source chunks/hash/provenance and all final-head jobs, then migrate captured auxiliary tables separately.


## Fresh-session reconciliation after PR105 (2026-10-03)

PR105 is VERIFIED OFFLINE at accepted merge 7e4ca5c955abb7f882591a44d8bb2f74c93ae243, tree 5c79e721d8281f04ba066006a4ff5601202abdbd, ordered parents c270ea22 / c9cc112. Final candidate run37143242571 (quality111261843942, Windows111261844116) and characterization37143242555/job111261844024 passed. Separate accepted-push quality run37143573050 (quality111262849038, Windows111262849202) and characterization37143573016/job111262848443 passed. Required steps and retrieved full logs were inspected in this recovery; only conditional failure upload skipped. Native primary/AI controls, compiled CLI/package cases, race/vet/fuzz, provenance and vulnerability scans passed; coverage81.8%, zero reachable/imported findings and the existing module-only advisory remains open. Final and accepted-push auxiliary capture chunks/provenance are complete and identical: inputs f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28, outputs efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06.

Live core-v1/ref/merge/tree/parents/human identity/rules verified; no open PRs. Reconciled stale active authority summaries and preserved historical checkpoints (FN067). Next separate slice: retain exact source auxiliary fixtures and implement bounded pure Quota/Capacity Reservations under REGION_AUXILIARY.md. No service decoder/arithmetic, public region execution, live/Gate004/release closure. Local Go QA unavailable in current Windows executor; authenticated GitHub/native Actions provide final-head acceptance, separately from deferred Azure/laptop validation.


### Auxiliary implementation checkpoint

Prepared pure ProjectQuota/ProjectReservations, independent exact captured-cell/hash tests, scope/malformed/duplicate/numeric/row/text/ownership/concurrency/cancellation guards and three compiling controls required by both native jobs. No fresh local Go PASS claimed. Python control syntax and retained SHA256s checked locally; native final-head QA pending. Source captures and input contract remain separate from service decoding/arithmetic and public availability. Formatting diagnostics add gofmt diff on existing failing guard without weakening it.


### PR106 initial native failure and corrected checkpoint

Initial PR106 head4e50009de3862d68b54e28dc8c3a557b7665de96 was rejected by native Linux formatting run37146459704/job111271307391: two struct alignment spaces differed from gofmt. Applied the exact native gofmt diff. No test/gate was weakened; this head remains unaccepted and broad Linux checks were skipped after that failure. Corrected head requires fresh complete characterization/Linux/Windows evidence. All fourteen initial published files read back exactly; tree31397358d52673a5f4cd37214d8ec918745dc76b, direct parent7e4ca5c and human author/committer verified. Scope/capture/runtime boundaries unchanged.


### PR106 fixture topology correction

PR106 native Windows run37146544096/job111271573338 rejected the authored fixture loader: decoding all fourteen output branches as single tables failed on the service-availability array. Fixed the test loader to retain all14 raw JSON branches and decode only the four quota/reservation table-or-null branches; complete capture hashes and all cell/guard assertions remain unchanged. No production or captured-source byte correction. Fresh final-head native validation is required.


Automated Code Review discussion_r4174461025 independently reproduced FN069 on initial4e50009. Current e8c1cc8 RawMessage correction is published/read back with complete unchanged captures and human identity; final native gates remain pending. A manual Code Review tool is unavailable, but repository-configured automated review evidence is accessible.


## PR106 acceptance and compact active continuity (2026-10-03)

Quota/Capacity Reservations pure projections are VERIFIED OFFLINE through [PR106](https://github.com/DeBoX85/Cloud-Assess/pull/106), merge `c4643527dc948f96a565da237662a104a084f14e`, tree `38248bc19be95c0058377ed373a4bf90b20e4ca2`, ordered parents7e4ca5c/e7e5511. Final head e7e551104cc6e0e6dda25c389af9176eb5ed562c passed [run37146790757](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37146790757), quality111272343099/Windows111272342986, all mandatory steps/full logs inspected. Both tested preview06848f39 with the identical tree/parents. Source characterization [37146790810](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37146790810)/job111272298401 reproduced all14 unchanged helper branches with exact captured bytes/provenance. Both hosts passed every quota/reservation cell/empty branch, all three compiling auxiliary controls and restored assertions, existing five primary/eight AI request/four AI execution controls,12 actual CLI and19 default/19 branded packages; docs327/9/1. Linux race/vet/fuzz and82.1% coverage passed; zero reachable/imported vulnerability findings, existing module-only advisory open. Remote merge/ref/tree/parents/human author/GitHub committer verified. Fresh local Go/fetched-source evidence is unavailable in this Windows session.

Both initial failures remain retained/classified; automated earlier-head fixture finding corrected, no independent-person approval. The old handover is preserved in SESSION_HANDOVER_HISTORY_20261003.md; compact active authority now identifies accepted functionality and exact next service-availability task without stale acceptance commands. No runtime/pin/fixture/workflow changes in this continuity slice. Final documentation native proof belongs in its PR and will be separately verified. The preparation-time postmerge pending state is superseded by the separately verified proof below.


Separate accepted-merge proof completed: [run37147153459](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37147153459), quality111273357208/Windows111273357343, every mandatory step/full log inspected on exact c4643527. Source [run37147153415](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37147153415)/job111273356990 reproduced all14 unchanged branches with exact hashes/provenance/bytes. Both hosts passed all existing/new controls,12 CLI/19 default/19 branded packages/docs327/9/1; Linux race/vet/fuzz/coverage82.1% and zero reachable/imported vulnerability findings. Module-only advisory stays open. This is separate accepted-push evidence, not inferred PR success or local execution.


### PR107 remaining-work review correction

Automated review discussion_r4174503451 found a contradictory REGION_SELECTION remaining-work list still requiring accepted pure Quota/CRG migrations. Removed those pure tasks and preserved the separate future service collection/decoding obligations; labelled retained primary preparation checkpoints historical. FN059 recurrence records cause and prevention. Fresh exact-head documentation/native/source proof remains required, with no runtime/pin/capture change.


## Service-availability pure projection checkpoint (2026-10-04)

PR107 continuity is accepted at3047085ab716fa012a255607a97fe4e0070bdc3a, tree3f177778d5910da763d3e7e783f35310919360dd, ordered c4643527/5e4bec6e parents. Final native37147778438 quality111275191623/Windows111275191760 and source37147778423/job111275166806 passed all required steps with full logs/chunks/provenance verified; only conditional failure upload skipped. Complete eleven-file docs-only scope/commit/tree/parents/human identity/base/head/preview/protection and accepted live ref verified. Coverage82.1%,12 actual CLI/19 default/19 branded packages, all existing controls and zero reachable/imported vulnerability findings; existing module-only advisory remains open. Distinct PR107 accepted-push proof completed on exact3047085a: native run37166264029, quality111329667969/Windows111329667829, all mandatory steps/full logs inspected (conditional failure upload only skipped); source37166264028/job111329667989 reproduced all14 helper branches with indexed complete chunks, exact retained bytes/hashes and source/tree/APRL provenance. Both hosts passed12 CLI/19 default/19 branded packages/all existing controls; Linux race/vet/fuzz/coverage82.1% and zero reachable/imported vulnerability findings. Module-only advisory remains open. This is separate accepted-push evidence, not inferred PR green or local execution.

Service availability is IN PROGRESS under [REGION_SERVICE_AVAILABILITY.md](REGION_SERVICE_AVAILABILITY.md), unaccepted pure ProjectServiceAvailability. Preserves every captured eight-column cell/ordering/empty branch and exact registry/restriction/case behavior, with declared selected aggregate contributors but no invented row UUID or proof of underlying ownership. Validates scope/labels/counts/duplicates/entry/replicated-row/joined-cell/global-text/sheet-collision limits and owned cancellation-safe tables. Literal independent fixture/registry/boundary/ownership16-concurrency tests and ten compiling named controls are prepared for both native jobs. Python control syntax checked; fresh Go/source/native results are pending, no local Go PASS. No runtime public availability/request/arithmetic, pin/capture/dependency or Azure/live/Gate004/release/advisory change. Finish exact published code/docs readback and fresh source/Linux/Windows QA/review/protected acceptance before separate CostComparison.


### PR108 Unicode transfer correction

FN070: native Windows111331457533/Linux111331457622 rejected the degraded replacement-character fixture on original0d4b0369. Original scratch UTF8 was intact; PowerShell stdout changed U+212A/U+FFFD/U+1F600 to ASCII approximations before remote readback. Changed test source to explicit Go Unicode escapes and reread through Python ASCII-safe JSON, preserving actual intended runtime characters. No production/source/capture/pin correction or failed-candidate acceptance. All final-head native/source/control proof must be refreshed.


### PR108 automated review corrections

FN071/discussion_r4175632172: reused auxiliaryScope tolerated empty names contrary to the service input contract. Added a service-only early selected-name guard, nil-data regression and compiling removal control. FN070/discussion_r4175632173 independently confirmed the previously corrected degraded Unicode test; extended actual escaped BMP/non-BMP boundary cases to32767 units/32770 bytes and one-unit excess with compiling rune/byte-confusion controls. Ten total compiling service controls must fail named assertions/restored baselines on both hosts. No accepted auxiliary/source/pin/capture/public behavior changed; revised complete-head evidence remains pending.


## PR108 acceptance and pure CostComparison checkpoint (2026-10-04)

Service availability is VERIFIED OFFLINE through [PR108](https://github.com/DeBoX85/Cloud-Assess/pull/108), merge1957f6d71da0c257bdd4b9d118c2be4b97b1a12e/tree3ffc1b51cb9ed5e608b805a3a7e2e1c1d54ba81e, ordered3047085a/0b7c95cd parents. Final0b7c95cd6d2bb8b9c7a60149b3694a9c898fb67d native [37167296996](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167296996), quality111332833289/Windows111332833405, and source [37167296991](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167296991)/111332794154 passed. Both tested previewfefd1118 with exact candidate tree/parents. Distinct accepted-push native [37167713486](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167713486), quality111334058385/Windows111334058216, and source [37167713500](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167713500)/111334058367 passed on exact1957f6d. All mandatory steps/full logs inspected; only conditional failure upload skipped. Both hosts passed10 service/five primary/three auxiliary/eight AI request/four execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages; Linux race/vet/fuzz and82.4% coverage, zero reachable/imported findings. Source14 branches/indexed complete chunks/exact retained hashes/bytes/source-tree-APRL provenance verified for both candidate and accepted push. Complete13-file scope/remote bytes and original scratch code/test/script Git blob hashes, human commit identity/base/head/preview/rules and protected merge/ref/tree/parents verified. FN070 Unicode transport and FN071 empty-label findings corrected before acceptance, with actual escaped BMP/non-BMP32767-unit boundaries and byte/rune controls. Automated review is earlier-head evidence, not independent-person/final manual-tool approval. No local Go/source fetch, public region, live/Gate004/release or module-only advisory closure is claimed.

Current CostComparison is IN PROGRESS/unaccepted under [REGION_COST_COMPARISON.md](REGION_COST_COMPARISON.md). Implements only source exact five+region-column cells/order/first received tuple+item metadata/positive four-decimal pricing and all three nil branches. Selected contributor declarations are validated without inventing UUIDs or proving aggregate ownership; all labels/maps/float values/rows/regions/entries/decoded/projected text are bounded and owned, cancellation drops partial candidates. Independent captures and literal metadata/case/empty/zero/negative/tiny/finite/boundary/Unicode/ownership16-concurrency/cancellation tests plus nine compiling controls/restored assertions are prepared for both native hosts. Local Python syntax/original file Git blob hashes checked via ASCII-safe JSON; no fresh local Go/source proof. Publish/verify complete code/docs checkpoint, then source/native/review/expected-head merge before Inventory. No request/arithmetic/public/pin/fixture/dependency/Azure/live/Gate004/release/advisory closure.


2026-10-04 PR109 FN072: native37168867529/Linux111337510113 and Windows111337510437 passed production package tests but failed the selected-contributor control compilation (unused selected). Corrected only the mutation expression to (!selected && false); controls require named assertion failure and restored passing baseline. No acceptance/skipped-step success claimed; fresh complete-head source/native proof is pending. Source37168867514/job111337510061 passed. Source/captures/pins and production are unchanged.


## 2026-10-04: CostComparison accepted; pure Inventory checkpoint

CostComparison is VERIFIED OFFLINE through [PR109](https://github.com/DeBoX85/Cloud-Assess/pull/109), merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents. 

Final candidate1cf8e4cf6adc38ffe377894824b51f3b26431373/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, parentca16157 on accepted1957f6d7. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165402 quality111338402524/Windows111338402644 passed every mandatory step; only conditional failure upload skipped. Both tested preview77c2b59d40f89e1966402668a98efcdfad2e4d5d with identical final tree/ordered base/head parents. Full logs inspected: all nine Cost controls/restored baselines, ten service/three auxiliary/five primary/eight request/four execution controls, actual CLI/default+branded package/docs/provenance/module/inventory/branding guards; Linux race/vet/fuzz and82.5% coverage; zero reachable/imported vulnerability findings, existing module-only advisory remains open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165342 /111338403371 completed; full14 branches, indexed3-input/4-output chunks, exact retained bytes/hashes and source/tree/APRL provenance verified. Complete14-file scope/deletions/remote bytes/original scratch blobs, parent/human identity/pins/base/head/preview/current no-bypass rules23890737 verified. Semantic self-review checked actual unchanged source first-tuple/item break behavior, all source topology/cells and safe bounded identity/nonfinite/text/cancellation behavior. No callable final manual Code Review tool or independent-person approval is claimed; automatic comments currently empty. FN072 compile failure was not accepted. Expected-head protected merge and distinct accepted-push evidence follow; no fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure.

Protected acceptance verified: PR109 merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents, Denis author/GitHub committer/livecore ref and no open proposals. Distinct accepted-push native/source QA remains pending; no dependent live proof claimed.

Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513599 quality111339415060/Windows111339415184 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513582 /111339414976 passed on exact6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f. All mandatory steps/full logs inspected; only conditional Linux failure upload skipped. Both hosts passed9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages/docs/branding/provenance; Linux race/vet/fuzz82.5%, zero reachable/imported findings. Source14 complete branches, indexed chunks/both retained exact bytes/hashes and pinned source/tree/APRL provenance verified. Module-only advisory and all live/public/laptop/Azure/Gate004/release obligations stay open. No fresh local Go/source or independent-person approval claimed.

Source BuildInventorySheet/pinned ComputeCapacity/MaskSubscriptionID/InResourceID/Excel duplicate-skip and target canonical guard/SKU helpers were inspected before implementation. Microsoft Learn resourceId documentation verifies subscription vs RG/tenant formats; ID correlation deliberately does not pretend full provider validation. Current ProjectInventory preserves every ten-column raw/masked/empty source cell/order with explicit Region Inventory composition correction, safe selected IDs/duplicate/capacity/text/row scope bounds, owned cancellation failure and eight compiling controls/restored baselines. Python AST/original scratch UTF8 Git blobs verified via ASCII JSON; no new Go/native/source PASS or Inventory acceptance. Update active handover/roadmap/spec/implementation/region/service/cost/auxiliary/QA/failure with this scope and exact next action. Pins/captures/dependencies/public paths unchanged; independent calculations follow acceptance, dependent live tests stay deferred.


2026-10-04 current source-only calculation checkpoint: Inventory is VERIFIED OFFLINE through [PR110](https://github.com/DeBoX85/Cloud-Assess/pull/110), merge2dcd8dae924ce3fad5402890cdeca0ecbef6e01e/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, ordered6dfcd1be/76ac3f8d parents. 

Final76ac3f8d65acf9970305c58186e13bb19cf78ffe/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, parent accepted6dfcd1be. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985397 quality111340777403/Windows111340777269 passed every mandatory step; only conditional failure upload skipped, full logs inspected. Both tested preview4e8f21f81155a55383272fe7ade895896fb72db6 with identical tree/orderedbase+head parents. Eight Inventory controls/restored baselines and all9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution controls passed each host;12 actual CLI/19 default/19 branded packages/docs/branding/provenance/module/inventory checks. Linux race/vet/fuzz passed; 82.7% coverage, zero reachable/imported findings; existing module-only advisory stays open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985428 /111340777428 passed full14 branches/indexed3-input4-output chunks/both retained exact bytes/hashes/source-tree-APRL provenance. Original scratch Go/test/script blobs/complete16 remote files/scope/no deleted or pinned files/commit identity/tree/parent/base/preview/rules verified. Semantic self-review checked actual source helper/mask/SKU/duplicate-skip plus explicit sheet/scope/budget/product corrections. Automated comments currently empty; no callable final manual Code Review tool or independent-person approval is claimed. No fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure. 
 Protected acceptance/ref/tree/parents/Denis author/GitHub committer verified. Distinct accepted-push proof follows below.
New inventory_capture_test.go.txt/extended capture-region-source.py invoke only unchanged pure inventory/merge/normalization/predicate functions in isolated pinned source/APRL.13 structural branches/30 literal normalization cases; actual output retention/topology/hash/final source/native/remote acceptance is pending. No target calculation code, pin/public/request/live/platform change. See REGION_INVENTORY_CALCULATIONS for scope/contracts/acceptance/rollback/next action.


PR110 distinct accepted-push verified: 
Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170292930 quality111341722998/Windows111341722896 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170292817 /111341722273 passed on exact2dcd8dae924ce3fad5402890cdeca0ecbef6e01e. All mandatory steps/full logs inspected, only conditional failure upload skipped. Both hosts passed8 Inventory/9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines;12 CLI/19 default/19 branded packages/docs/branding/provenance. Linux race/vet/fuzz82.7%; zero reachable/imported findings, existing module-only advisory open. Source14 full branches/complete indexed chunks/retained hashes/bytes/source-tree-APRL provenance exact. No fresh local Go/source/public/live/laptop/Azure/Gate004/release/advisory closure.


Observed initial source capture: [37170559596](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170559596)/job111342489938 on initial9a2fffef completed both pure harnesses, unchanged original14 and new13 inventory branches/30 normalization cases, original source/APRL diff/no-extra-test-file/provenance proof. Complete indexed new2-input/4-output chunks reconstructed; original scratch UTF8 bytes/SHA256 verified, no fabricated output. Retained input cdaa6b25e3250e6d701999279f96d5ea225c927e1249ee43c32cfe0b5072b17c and output756647e2981f8cd0ac4f0d5091ba6b4a74345ad6266de1bc15f8c833afeb2006. Both prior auxiliary hashes/bytes unchanged. Literal results confirm exact subscription filtering/raw SKU keys/five maps/ASCII-space-only locations/empty-logical keys/repeated counts/source nil panics and blacklist unknown-region behavior. New Go fixture hash/topology/critical count/normalization checks are prepared but no Go PASS claimed yet. Final runner now unconditionally requires exact bytes against all four retained files. Python actual extracted byte guard passed four baseline/restored assertions and a compiling disabled-guard named failure locally; both mandatory native hosts now require it. It performs no source/network checkout in its focused test. Fresh complete final-head source/Linux/Windows all mandatory steps/full logs/remote scope/bytes/tree/identity/base/preview/rules/protected acceptance remain pending; initial-head success cannot certify these guards/tests.

Current17-file characterization checkpoint includes new actual JSONs/Go capture-hash-topology test/Python byte-guard assertions/control and both native workflow steps. No target inventory calculation runtime added. Local real Python byte-guard tests/control/restoration passed; source hash/UTF8 and original Git blobs verified. Fresh final-head required proof remains pending.

## PR111 acceptance and owned aggregation checkpoint (2026-10-04)

Inventory calculations source characterization is VERIFIED OFFLINE through [PR111](https://github.com/DeBoX85/Cloud-Assess/pull/111), merge1b17cccd0bb4ff9f08627d4ea398d3e628107afb/tree94685630597212c81c43f3f199c2068884cd3044, ordered2dcd8dae/667e680b parents. 

Final667e680bfe19c1088008b74cdaf4c9f26ce16ef6/tree94685630597212c81c43f3f199c2068884cd3044, parent9a2fffef on accepted2dcd8dae. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170957575 quality111343642564/Windows111343642362 passed every mandatory step, only conditional Linux failure upload skipped, full logs inspected. Both tested preview700c56ac82bc48a292cd11ea8d3e823cda6b281a with exact candidate tree/orderedbase+head parents. Both passed four actual byte-guard assertions/compiling disabled-guard named failure/restored baselines, new capture hash/topology assertions, all8 Inventory/9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution controls,12 actual CLI/19 default/19 branded packages/docs/branding/provenance/module/inventory. Linux race/vet/fuzz82.7%, zero reachable/imported findings; existing module-only advisory open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170957525 /111343642132 passed unchanged14 auxiliary/new13 inventory/30 normalization cases and unconditional four-file byte equality. Full indexed3+4+2+4 chunks/actual retained hashes/bytes/source-tree-APRL/original production+APRL diff/only two source test files verified. Complete17-file scope/remote readback/original scratch UTF8 Git blobs/identity/tree/parents/base/preview/no-bypass rules23890737 verified. Semantic self-review compared actual source map/case/normalization/merge/panic/predicate definitions to complete captures, safe confined runner and byte boundary/control. Automatic comments empty; no callable final manual Code Review tool or independent-person approval claimed.  No target calculation runtime or new local Go/source/public/live/laptop/Azure/Gate004/release/advisory closure.
 Accepted commit/ref/tree/parents/Denis author/GitHub committer verified. 
Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37171306079 quality111344633775/Windows111344633847 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37171306081 /111344633669 passed exact1b17cccd0bb4ff9f08627d4ea398d3e628107afb. All mandatory steps/full logs inspected, only conditional Linux failure upload skipped. Both hosts passed new captured hashes/topology, four actual byte-guard assertions/compiling guard-control/restored baselines and all8 Inventory/9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution controls,12 CLI/19 default/19 branded packages/docs/branding/provenance; Linux race/vet/fuzz82.7%, zero reachable/imported findings. Source unchanged14/new13 branches/30 normalization rows/all four files fresh exact bytes/hashes/3+4+2+4 contiguous chunks/pinned source-tree-APRL/original files and two-harness isolation verified. No fresh local Go/source/public/live/laptop/Azure/load/Gate004/release/module-only-advisory closure.

Inventory source characterization PR111 is VERIFIED OFFLINE. The bounded owned aggregation helper is IN PROGRESS/unaccepted under [REGION_INVENTORY_AGGREGATION.md](REGION_INVENTORY_AGGREGATION.md). Public region-selection remains unavailable. Finish fresh exact-head native/source QA, review and protected acceptance, then availability/latency/cost/quota/reservation arithmetic, bounded request adapters and coordinator/public/all-format execution. Azure/laptop-dependent validation remains explicitly deferred.

Prepared CalculateInventory/InventoryCounts/InventoryCalculation with seven independent focused tests, actual complete source map oracles/canonical identity correction, explicit full unsafe fixture rejection and separate literal valid subset. Ten compiling negative controls enforce isolated row/entry/decoded/projected bounds, normalization/raw-SKU/resource-count/selection/correlation guards on both hosts. Python AST and original scratch UTF8 blobs verified; no new local Go/native/source acceptance. Coherent code/docs branch publication and exact-head gates are next; source/capture/dependency/public/platform pins unchanged.


## 2026-10-04 comprehensive project audit preparation

User agreed to pause feature expansion following repeated development-process interruptions. [PROJECT_AUDIT_20261004.md](PROJECT_AUDIT_20261004.md) records the current progress, verified accepted07011b63/PR112, unmerged113/f1142c7d, source37174422877 success metadata and native37174420499 failure with no returned jobs (cause unclassified), transient unpublished availability draft limitation, full project audit phases and autonomy/input boundaries. Plan PREPARED; audit execution NOT STARTED. Historical product/runtime stability is not inferred from service interruptions. New Linux executor inventory and exact evidence reconciliation precede fixes; no acceptance/live/release claim.


Audit outage continuity preparation: added AUDIT_RECOVERY.md with explicit checkpoint fields, unknown remote-mutation handling, live-state recovery and saved restart prompt. Protocol applies throughout audit execution, which remains NOT STARTED. No outage-proof guarantee or background work claimed.


## 2026-10-04 comprehensive audit first execution checkpoint

Recovered exact accepted07011b63 signed commit/tree/all630 blobs and verified current protection/pins, preserving missing local APRL/Go limits. Planning771919b9 full native37182711852 and source37182711928 logs/preview/contiguous captures verified. Confirmed FN074 replacement-string workflow corruption in unmerged113 and published literal correctionbe867a174dbc112caba273cda620822fd9828b2e with local YAML/eight Bash syntax and exact remote-byte proof; fresh gates pending, merge paused. PROJECT_AUDIT_REPORT_20261004.md inventories all216 Go files and remaining review scope, records stale documentation and unpublished draft limits, and flags ARG drain-error behavior for reproduction rather than premature defect claim. Whole-project audit IN PROGRESS, no new live/release acceptance.


## Comprehensive audit retrieval and repair checkpoint (2026-10-04)

Repaired PR113 full native/source proof verified onbe867a17: native37183435404 quality111380420837/windows111380420966, source37183435389/111380420823. Preview8a040608 tree0878fc48 equals candidate, ordered07011b63/be867a17 parents. All mandatory checks, six capture byte/hash/chunk proofs and exact provenance pass; feature acceptance still paused.

AUD003 narrow library adapter error swallowing reproduced with eight compiling named failures on both native hosts against unchanged730a29e6. Default pinned SDK buffered pipeline mitigates actual default route, no live affected assessment claim. PR115 defensive error propagation/shared-pipeline regression head23d2ce4d/tree7438dfd7 parent730a29e6 verified identity/three bytes; full fresh QA pending. Initial recovery and selected retrieval review ongoing. Frozen07011b63 unchanged; audit-wide code/test/release/live closure not claimed. Next inspect115 QA, Cost continuation contract, remaining orchestration/Defender/test requirements and subsequent areas. Durable report/handover updated.


## Audit ARG verification, Cost completion and report review (2026-10-04)

PR11523d2ce4d full native37184407342/Linux111383253100/Windows111383253050 and source37184407425/111383253548 passed. Actual shared SDK terminal-read mitigation fixture passes; corrected injected-poster boundary no longer loses errors. Full logs, four retained capture bytes/provenance and preview87b10d5d/tree7438dfd7/ordered07011b63+23d2ce4d parents verified. Unmerged.

AUD005 Cost first-page false-success confirmed by two compiling assertions on both platforms in37184679263/test-only78dd85c9. PR116 defensive fail-closed continuation guard and actual coordinator partial/retained-stage fixture published8d54c0d9/treeed772a94/parent78dd85c9 with verified identity/3 bytes. No full pagination or live incidence claimed; fresh full QA pending. Production review and selected independent fixtures now include canonical/result/findings/gate/config/branding/redaction/report replacement/all-renderers and public zone/carbon/service-health/SQL-EOL plugin boundaries; remaining tests/AI/region/tooling/requirements keep audit IN PROGRESS. Frozenaccepted07011b63 unchanged.


## Audit Cost verification and comparator evidence checkpoint (2026-10-04)

PR1168d54c0d9 full native37184879992/Linux111384624285/Windows111384624226/source37184879983/111384624557 passed. Full mandatory logs/preview9caddc01 candidate tree/ordered parents/four contiguous complete capture bytes/hash/pins verified. Actual coordinator partial/healthy retained results regression passes; no full pagination/live incidence claim. Unmerged.

AI/region/equivalence/scanner/SKU production review extended; full test/requirements review still open. Suspected AUD006 unsafe comparator input acceptance now test-only PR1174cb8bd8451208e45e2551f72edef280de338e52f/tree3ae8c71f22fcb86fbe951810bd997fc7afd5db76/parent07011b63, exact Denis identity/two remote bytes verified. Six literal real-command invalid cases and two healthy empty controls, unchanged production, native reproduction pending. Next inspect negative control then tooling/Diagnostics/workflows. Accepted07011b63 unchanged; feature113 paused, no live/release/full-project acceptance.


## Audit comparator correction and tooling coverage (2026-10-04)

AUD006 test-only4cb8bd84/native37185359814 reproduced all6 compiling named command exit0 failures on Linux111386033586/Windows111386033687, healthy empty controls pass. Corrected733b924383561b300bcbcc68998c6c849cf2d3ab/tree92e1bb688222a8f86f50773a05cda5dc054dd363/parent4cb8bd84 Denis identity/three bytes verified on draft PR117; full native37185621965/source37185621953 pending. See comparison-input contract, no historical live incidence claim.

Production Go inventory41 areas fully read, all4 workflows and9 production scripts examined; test/oracle/requirements/docs/dependency reconciliation still open. Read-only local execution request stalled/terminated, fallback to independently verified source bytes; no product failure or uncertain mutation inferred. Frozen07011b63/proposals113-117 remain unchanged/unmerged. No live/release/full-audit acceptance.


## Audit additive-schema correction and full script review (2026-10-04)

Requirements review caught initial117733b9243 unsupported valid1.1 additive schema. Literal test-only1bf98b7a/native37185951913 reproduced one named compiling failure on Linux111387770853/Windows111387770760, previous invalid/healthy controls pass. Finale381d0ad6e45989cf64526ac5d34131dc349e8e5/tree5e91d391cead6b74507abb6ac780817b9a42afaf/parent1bf98b7a allows both declared schemas, Denis identity/three bytes verified. Fresh final full QA pending. FN076 retains prevention, no accepted regression.

All31 scripts/4 workflows read; complete Azure/discovery/config/stages/gate/throttling/findings/redaction tests read. Remaining integration/plugin/result/rules tests and specification/evidence reconciliation keep audit open. Filter unknown YAML keys/multiple documents are a new pre-auth scope question awaiting independent reproduction. Accepted07011b63 unchanged, all proposals draft/unmerged.
