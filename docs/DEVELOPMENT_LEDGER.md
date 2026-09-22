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
| `EQUIVALENCE.md` | Procedure and normalization rules for source-vs-target comparison | Validation runbook |
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

## Current boundary

The deterministic/local development phases through semantic equivalence tooling and the reproducible live-equivalence runner are complete and quality-gated.

The next major phase is:

```text
Live Azure source-versus-target regression
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

- live Azure equivalence
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
