# Source-versus-Target Equivalence Runbook

Status: development validation procedure

Reference implementation: `DeBoX85/azqr`

Pinned reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

## Purpose

The reference implementation and Cloud Assess intentionally emit different JSON schemas.

The pinned reference JSON is a collection of report-table rows. Cloud Assess JSON is the canonical assessment domain model with stage health and completeness metadata.

For that reason, raw JSON equality is not a valid equivalence test.

The development-only equivalence harness projects both reports into the same semantic datasets and compares stable assessment meaning.

## Development tool

Run:

```bash
go run ./tools/equivalence \
  --reference reference.json \
  --target target.json \
  --output equivalence.json
```

Exit codes:

```text
0 = semantically equivalent for the compared core datasets
1 = one or more semantic differences were found
2 = invalid invocation, unreadable input, or invalid report structure
```

This tool is deliberately outside the public `cloud-assess` command tree. It is a reproduction/regression harness, not the deferred historical-report comparison feature.


## Automated Windows live runner

For the normal Windows development workflow, use:

```powershell
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -SubscriptionId <subscription-id>
```

For a resource-group-scoped first pass:

```powershell
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -SubscriptionId <subscription-id> `
  -ResourceGroup <resource-group>
```

Multiple subscriptions are accepted as a PowerShell array or a comma-separated value:

```powershell
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -SubscriptionId <subscription-1>,<subscription-2>
```

One or more resource groups may be supplied when exactly one subscription is selected:

```powershell
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -SubscriptionId <subscription-id> `
  -ResourceGroup <resource-group-1>,<resource-group-2>
```

For management-group validation:

```powershell
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId <management-group-id>
```

Multiple management groups may likewise be supplied as an array or comma-separated value.

The runner:

1. verifies the reference checkout is exactly at the pinned AZQR commit;
2. requires both reference and Cloud Assess working trees to be clean;
3. verifies both APRL submodules are initialized at the pinned APRL revision;
4. records the exact source and target commit SHAs and branches;
5. records Go/PowerShell versions and Azure cloud endpoint variables that affect runtime behavior;
6. records the selected scope, requested stage controls, and resolved effective stages;
7. hashes any reference/target filter files used;
8. runs the pinned reference with JSON enabled and masking disabled;
9. runs Cloud Assess with canonical JSON enabled and redaction disabled;
10. captures stdout and stderr separately for both scans;
11. runs the semantic comparator;
12. writes a machine-readable `run-metadata.json` alongside the two reports and equivalence report.

When `-Stages` is omitted, the runner passes no stage override to either executable and records the resolved default set (`graph`, `diagnostics`, `advisor`, and `defender`) with `usesImplicitDefaults: true`. Explicit additions and `-stage` removals are resolved and recorded separately from the original requested controls.

By default, the evidence is written beneath:

```text
artifacts/equivalence/<UTC timestamp>/
```

The `artifacts/` directory is ignored by Git because live-equivalence evidence contains unredacted Azure identifiers and may contain other sensitive assessment metadata.

A typical evidence bundle contains:

```text
reference.json
target.json
equivalence.json
run-metadata.json
reference.stdout.log
reference.stderr.log
target.stdout.log
target.stderr.log
equivalence.stdout.log
equivalence.stderr.log
```

The runner performs read-oriented assessment operations only; it does not provision or modify Azure resources.

Starting with run-metadata schema `1.2`, the runner also records `targetAssessment`: the target completeness and each stage's name, status, record count and warning codes. It deliberately omits warning messages and resource identifiers from this summary. Earlier metadata bundles remain valid without this property. A warning code still requires inspection of the private target report and, where necessary, request-correlated evidence before its impact is classified.

### Investigating Diagnostics batch warnings

If the target Diagnostics stage records `diagnostics_subrequest_non_success` warnings, the batch status alone does not identify the affected resources. From the clean Cloud Assess checkout used for the pass, probe the retained target inventory with individual, read-only diagnostic-settings GET requests:

```powershell
go run ./tools/diagnostics-probe --target C:\src\Cloud-Assess\artifacts\equivalence\<run-stamp>\target.json
```

The probe uses Cloud Assess's normal Azure credential and endpoint selection and the same supported-resource table and diagnostic-settings API version as the scanner. It caps the number of requests at 100 unless `--max-requests` is raised deliberately. Its console JSON contains only the number of eligible resources, the original HTTP 400 warning count, successes, and failures with resource type, HTTP status, Azure error code, and a shortened SHA-256 hash of the resource ID. Raw resource IDs are not printed. Keep the original target report private.

Individual GET responses may differ from the earlier ARM batch subresponses or from later Azure state. If the two HTTP 400 responses are not reproduced, retain the original warnings as unresolved and investigate batch-specific behavior; do not reclassify an uncertain diagnostic finding as confirmed simply because the probe succeeded.

To investigate the batch endpoint itself without assuming that a multi-request response preserves request order, use one read-only diagnostic-settings GET inside each ARM batch:

```powershell
go run ./tools/diagnostics-probe --mode single-batch --target C:\src\Cloud-Assess\artifacts\equivalence\<run-stamp>\target.json
```

The command makes one batch POST per eligible resource, with one GET subrequest in each batch. A non-success subresponse is therefore associated with the sole submitted resource; the summary prints only its shortened resource-ID hash, type, HTTP status and safe Azure error code. It caps requests at 100 by default and uses the scan's ARM batch and diagnostic-settings API versions. The POST executes a read-only GET subrequest; it is still an Azure API call and should be run only in an approved scope. A one-request batch may behave differently from the original multi-request batches, and Azure state can change. A matching resource hash is corroboration for a repeated failure, not proof of which resource failed in a prior unrecorded multi-request response. Retain original warning uncertainty until the original response is correlated or an explicit limitation is accepted.

### Preparing the pinned reference checkout

Create a dedicated clean checkout for the reference:

```powershell
git clone https://github.com/DeBoX85/azqr.git C:\src\azqr-reference
Set-Location C:\src\azqr-reference
git checkout 8e4f0577f3615e6c9014c031bcad079f235369cc
git submodule update --init --recursive
```

Also initialize the target submodule after cloning Cloud Assess:

```powershell
git submodule update --init --recursive
```

Do not make local edits in either checkout before an equivalence run. The runner intentionally refuses dirty repositories because an uncommitted state cannot be reconstructed later from the recorded SHA.

### Filter-file note

Cloud Assess intentionally does not require legacy configuration compatibility, so the pinned reference and Cloud Assess do not necessarily consume the same filter-file schema.

For filtered equivalence runs use:

```powershell
-ReferenceFilters C:\path\reference-filters.yaml
-TargetFilters C:\path\target-filters.yaml
```

The two files must express the same assessment intent. Their SHA-256 hashes are captured in `run-metadata.json`.

The first live validation should use no filter file so configuration translation cannot obscure core assessment parity.

For the first service-selection pass on the `AdvisoryDev` non-production leaf management group, the checked-in pair at `examples/filters/azqr-storage-vm.yml` and `examples/filters/cloud-assess-storage-vm.yml` includes only the `st` and `vm` scanner keys. The pinned source uses `azqr:` and Cloud Assess uses `assessment:`; both files select Storage Accounts and Virtual Machines. The known Dev inventory has both types, so this pass can exercise non-empty inventory and findings without including the Network Watcher in the diagnostic-settings batch. The absence of its warning in a filtered pass does not close the unfiltered Diagnostics warning boundary.

```powershell
Set-Location C:\src\Cloud-Assess
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters .\examples\filters\azqr-storage-vm.yml `
  -TargetFilters .\examples\filters\cloud-assess-storage-vm.yml
```

The runner records hashes of both filter files. Review the resolved subscription, enabled stages, warning codes and semantic datasets before treating this pass as evidence for scanner selection. This fixture does not cover resource-group, tag, individual-resource or recommendation filters; those require separate passes.

The next Dev filter pass tests resource-group inclusion independently. The checked-in pair `examples/filters/azqr-dev-rg.yml` and `examples/filters/cloud-assess-dev-rg.yml` includes the full ARM ID of `rg-fos-FinopsHub-dev-euw` in the Dev subscription. A local review of the unfiltered Dev target inventory found three in-scope resources in that group. Keep the management-group scope and default stages the same as the earlier Dev passes; do not also supply the CLI resource-group flag, which would make the filter's effect harder to distinguish.

```powershell
Set-Location C:\src\Cloud-Assess
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters .\examples\filters\azqr-dev-rg.yml `
  -TargetFilters .\examples\filters\cloud-assess-dev-rg.yml
```

A second pair, `examples/filters/azqr-environment-dev.yml` and `examples/filters/cloud-assess-environment-dev.yml`, includes resources tagged `Environment: dev` regardless of resource group. Run that pair separately only after confirming the unfiltered Dev inventory contains both matching and nonmatching resources; otherwise the live pass cannot demonstrate that the tag filter excludes anything. All three resources observed in the named group carry this tag, so combining the two include conditions there would not independently validate tag filtering. For both passes, inspect target stage health and the resolved subscription, as well as exact comparator counts. A warning-free filtered pass cannot resolve the earlier unfiltered Diagnostics warnings.

The unfiltered Dev inventory check confirmed 3 matching and 9 nonmatching resources. To reproduce the separate tag-only pass, use the same `AdvisoryDev` scope and default stages, with no CLI resource-group flag:

```powershell
Set-Location C:\src\Cloud-Assess
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters .\examples\filters\azqr-environment-dev.yml `
  -TargetFilters .\examples\filters\cloud-assess-environment-dev.yml
```

The Phase S tag-only pass selected exactly the same three tagged resource IDs as the earlier unfiltered baseline; the other nine resources were excluded and both reports matched semantically. If the Dev inventory changes in a later run, investigate a different count before treating a semantically equal result as tag-filter coverage. A later local comparison confirmed the tag-only and RG-only selected ID sets also matched. Advisor returned two rows on each side in the tag-only pass versus three on each side in the earlier RG-only pass. A subsequent local check of the missing row against both selected and out-of-scope tag-only target inventory returned `NearestRecordedScope = unknown`: no resource or ancestor decision was recorded. Both pinned source and target exclude unknown downstream scope with an include-tag filter if the row is returned and structurally eligible. Because the reports came from separate Azure scans, the check cannot establish whether Azure returned that row in the later run. See the dated Phase S follow-ups in the ledger.

### Paired filter pass: exclude the Dev tag

The pair `examples/filters/azqr-exclude-environment-dev.yml` and `examples/filters/cloud-assess-exclude-environment-dev.yml` excludes resources tagged `Environment: dev` without an include filter. Run it against the same `AdvisoryDev` leaf scope with the implicit default stages. The earlier unfiltered Dev target report had three matching and nine nonmatching in-scope resources. Compare selected IDs to the retained unfiltered baseline and inspect excluded tag records. The old counts are an expectation, not a fixed acceptance condition; obtain a contemporaneous unfiltered baseline if the inventory has changed or a stronger current-state claim is needed.

```powershell
Set-Location C:\src\Cloud-Assess
git fetch origin bootstrap/core-v1
git merge --ff-only origin/bootstrap/core-v1
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters .\examples\filters\azqr-exclude-environment-dev.yml `
  -TargetFilters .\examples\filters\cloud-assess-exclude-environment-dev.yml
```

The runner performs read-only assessment requests and writes reports under the local ignored `artifacts/` directory. It does not provision or change Azure resources. Preserve its `run-metadata.json`, source/target JSON, logs, and `equivalence.json` locally. Review `equivalent`, completeness, scope-stage subscription count, stage warning codes, and the selected/out-of-scope resource ID sets. The known Network Watcher may re-enter this wider scope; classify any Diagnostics warning without assuming the individual GET probe identifies a historical batch subrequest. Advisor can vary across separate Azure scans, so do not attribute a cross-run row difference to the filter without observed matching API input. Do not add the RG or include-tag filter to this pass.

The `20260929_113319Z` pass completed with `equivalent = true` for nine selected resources and non-empty findings and Advisor data. The user's local ID check matched the nine non-Dev-tagged in-scope IDs from the earlier unfiltered baseline; none of the nine selected resources had the Dev tag. The target was `complete_with_warnings` because Diagnostics recorded one `diagnostics_subrequest_non_success` warning. Do not treat that warning as resolved by the exact semantic match. See Phase T in the ledger for the evidence and remaining limits.

A local follow-up confirmed that all populated target data belonged to the expected Dev subscription. The read-only individual-GET probe found one eligible Network Watcher returning HTTP 400 `ResourceTypeNotSupported`, with the same hashed ID as in earlier probes. This supports the existing explanation but does not map the batch warning to that resource; retain `complete_with_warnings`.

### Next paired filter pass: exclude the Dev resource group

The pair `examples/filters/azqr-exclude-dev-rg.yml` and `examples/filters/cloud-assess-exclude-dev-rg.yml` excludes the full ARM ID of `rg-fos-FinopsHub-dev-euw` without an include filter. This tests resource-group exclusion independently of the completed RG include and tag exclude passes. The earlier unfiltered Dev inventory contained three in-scope resources in that group and nine outside it. All three group resources carried `Environment: dev`, so the selected IDs should match the nine from the Phase T tag-exclude report if the inventory has remained stable. These older counts and ID sets are comparison points, not unconditional acceptance criteria.

```powershell
Set-Location C:\src\Cloud-Assess
git fetch origin bootstrap/core-v1
git merge --ff-only origin/bootstrap/core-v1
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters .\examples\filters\azqr-exclude-dev-rg.yml `
  -TargetFilters .\examples\filters\cloud-assess-exclude-dev-rg.yml
```

Keep both unredacted reports and logs in the ignored local evidence directory. Review comparator equality, the exact resolved subscription, stage completeness and warnings, the selected IDs against the earlier baseline and Phase T report, and whether the three previously included group resources are now recorded out of scope. Resource-group exclusion may still include the Network Watcher elsewhere in Dev; do not treat a repeated Diagnostics warning as resolved by semantic equality or an individual GET. Advisor rows can vary between scans, so assess source/target agreement within this run before interpreting cross-run differences. Do not add the CLI RG flag or a tag filter to this pass.

The `20260929_121323Z` pass returned `equivalent = true` with nine selected resources, 33 findings and eight Advisor rows on each side. The user's local check found all three previously included resources in the excluded group in the target `outOfScope` set, none selected, and an exact match between the selected IDs and the same-day Phase T tag-exclude pass. Only the expected Dev subscription contributed data. The target remained `complete_with_warnings` with one Diagnostics subrequest warning; see Phase U in the ledger. The prior Network Watcher probe does not correlate this run's batch failure to a resource.

### Paired filter pass: exclude one observed recommendation

The pair `examples/filters/azqr-exclude-vm-backup-recommendation.yml` and `examples/filters/cloud-assess-exclude-vm-backup-recommendation.yml` excludes pinned APRL recommendation `1981f704-97b9-b645-9c57-33f8ded9261a`. Its pinned KQL targets Azure virtual machines without backup, and the user's local Phase U report showed one APRL/AOR finding with this ID. The fixture tests recommendation exclusion without narrowing inventory by resource group, tag, or scanner. It was run on the non-production `AdvisoryDev` leaf with implicit default stages using:

```powershell
Set-Location C:\src\Cloud-Assess
git fetch --quiet origin bootstrap/core-v1
if ($LASTEXITCODE -ne 0) { throw 'Fetch failed' }
git merge --ff-only origin/bootstrap/core-v1
if ($LASTEXITCODE -ne 0) { throw 'Fast-forward failed' }
.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters .\examples\filters\azqr-exclude-vm-backup-recommendation.yml `
  -TargetFilters .\examples\filters\cloud-assess-exclude-vm-backup-recommendation.yml
```

Retain the raw reports and logs locally. Verify that the selected recommendation ID had a finding in the earlier Phase U report, that the same VM resource remains in the new target inventory, and that no target finding with this ID remains. Compare other findings and inventory against the baseline only after checking for intervening Azure changes; an old one-row count alone cannot prove why a current row is absent. Review comparator equality, stage coverage, exact scope subscription and completeness separately. A Diagnostics warning or cross-run Advisor count difference must be classified on its own evidence. The source and target recommendation catalogs may still list the excluded definition; do not assume the catalog row count must drop when execution is filtered.

The `20260929_123149Z` pass returned `equivalent = true`: 313/313 recommendations, 43/43 findings, 12/12 inventory, 10/10 out-of-scope and 10/10 Advisor, with no enabled-dataset deltas. Reference, target and comparator exited 0. The user's local raw-report check found one baseline finding with the excluded ID, zero in each new report, and the affected resource still selected in the target inventory. The new target had `complete_with_warnings`; its warning codes and exact affected requests were not supplied. This is evidence for the observed recommendation exclusion, subject to that stage-health limit and the separate-snapshot boundary. See Phase V in the ledger. The observed recommendation catalog count is 313 on both sides; the preceding Phase U pass had 314, but the filter effect is established by the finding and inventory checks rather than a catalog-count inference.

### Paired filter pass: exclude one observed resource

The pinned source calls exact resource-ID exclusions `azqr.exclude.services`; Cloud Assess calls them `assessment.exclude.resources`. To avoid committing an unredacted resource ID, `scripts/prepare-resource-exclusion.ps1` derives a VM from the observed Phase U backup finding, checks that the same ID remains in the newer Phase V target inventory under the expected Dev subscription, and writes both YAML files under ignored `artifacts/`. It reads local JSON only and makes no Azure request. On the user's Windows checkout, after fetching the branch containing this helper:

```powershell
Set-Location C:\src\Cloud-Assess
git fetch --quiet origin bootstrap/core-v1
if ($LASTEXITCODE -ne 0) { throw 'Fetch failed' }
git merge --ff-only origin/bootstrap/core-v1
if ($LASTEXITCODE -ne 0) { throw 'Fast-forward failed' }

$prepared = .\scripts\prepare-resource-exclusion.ps1 `
  -BaselineTargetJson 'C:\src\Cloud-Assess\artifacts\equivalence\20260929_121323Z\target.json' `
  -CurrentTargetJson 'C:\src\Cloud-Assess\artifacts\equivalence\20260929_123149Z\target.json' `
  -ExpectedSubscriptionId 'c09f96df-19de-4c49-80ff-0c1a94d93ab2'
$prepared | Format-List

.\scripts\live-equivalence.ps1 `
  -ReferenceRepo C:\src\azqr-reference `
  -ManagementGroupId AdvisoryDev `
  -ReferenceFilters $prepared.ReferenceFilters `
  -TargetFilters $prepared.TargetFilters
```

Keep both generated filter files with the local evidence bundle. The preparer prints only a resource ID hash and counts, plus local paths. The runner records the filter paths and SHA-256 hashes. Do not publish the generated YAML or unredacted reports. For a discriminating result, verify the previously selected resource moves out of target inventory and into `outOfScope`, related findings disappear or are consistently filtered on both sides, and unrelated inventory remains. Review resolved and contributing subscriptions, stage warnings and the semantic comparator separately. Azure may change between runs, so do not infer causation solely from a changed total count.

The `20260929_132714Z` pass returned `equivalent = true` with 314/314 recommendations, 35/35 primary findings, 11/11 selected and 11/11 out-of-scope resources, and 5/5 Advisor rows; all enabled datasets had no deltas. The user's local check confirmed the previously selected VM was excluded, appeared once out of scope, had zero remaining target findings, and that every other selected resource ID matched the earlier Phase V target inventory. Target completeness was `complete_with_warnings`: Diagnostics had one `diagnostics_subrequest_non_success` warning; the failed request and HTTP status were not supplied. This establishes the observed exact-resource-exclusion behavior, with the Diagnostics warning and separate Azure snapshot limits preserved. See Phase W in the ledger.

A later `single-batch` probe against this retained target inventory returned nine successful one-request ARM batches and one Network Watcher HTTP 400, with the same shortened resource-ID hash as earlier individual GET failures. It corroborates the likely cause but cannot map the original multi-request warning to that resource. The original target remains `complete_with_warnings`; see the Phase W Diagnostics follow-up in the ledger.

## Required run conditions

For a meaningful live comparison:

1. Build or run the reference at the exact pinned commit.
2. Run Cloud Assess from the branch/revision being validated.
3. Use the same Azure identity and cloud environment.
4. Use the same subscription, management-group, resource-group, filters, and stage selection.
5. Run the two assessments close together and avoid making Azure changes between them.
6. Generate JSON from both tools.
7. Disable subscription-ID masking/redaction so full Azure resource identity is available.
8. Do not enable reference plugins while Cloud Assess plugin execution remains deferred.

Both inputs must contain full, unredacted subscription identities. The harness detects the reference/Cloud Assess subscription mask marker and rejects redacted reports as non-comparable rather than producing misleading record deltas.

A target assessment with overall `partial` or `failed` completeness is rejected as a valid equivalence baseline.

A target assessment with `complete_with_warnings` remains comparable, but the stage warnings must be reviewed alongside the semantic diff.

## Baseline commands

For a single subscription using the default stage set:

### Pinned reference

```bash
azqr scan \
  --subscription-id <subscription-id> \
  --json \
  --xlsx=false \
  --mask=false \
  --output-name reference
```

### Cloud Assess

```bash
cloud-assess scan \
  --subscription-id <subscription-id> \
  --json \
  --xlsx=false \
  --redact-subscription-ids=false \
  --output-name target
```

Then:

```bash
go run ./tools/equivalence \
  --reference reference.json \
  --target target.json \
  --output equivalence.json
```

The same approach applies to management-group or resource-group scope.

## Optional stages

Stage coverage itself is compared. A stage enabled on only one side produces a coverage mismatch before record-level analysis.

For optional stages, pass the same stage selection to both tools. Examples include:

- `policy`
- `cost`
- `arc`
- `defender-recommendations`

The first live pass should use the default core stages. Optional stages should then be enabled deliberately so failures and data-access requirements remain easy to classify.

## Compared semantic datasets

The harness currently compares:

- recommendation summaries
- primary impacted findings
- resource-type counts
- in-scope inventory
- out-of-scope inventory
- Azure Advisor
- Defender for Cloud plan status
- Defender recommendations
- Azure Policy noncompliance
- Arc-enabled SQL
- Cost Management

External/internal plugin results are intentionally excluded until production plugin execution is implemented in Cloud Assess.

## Primary finding identity

Primary findings are matched by normalized:

```text
Recommendation ID + Resource ID
```

Once matched, stable semantic fields such as category, impact, source, resource type, recommendation text, subscription/resource identity, parameters, and learn-more URL are compared.

SLA findings are not treated as ordinary impacted findings. They are instead used to reproduce the SLA value attached to inventory rows, matching report behavior.

## Known intentional normalization

The harness suppresses only documented source/target differences that should not count as equivalence failures.

### Legacy custom-rule provenance

The pinned reference labels its embedded custom rule corpus as:

```text
AZQR
```

Cloud Assess deliberately uses:

```text
CUSTOM
```

The harness maps legacy `AZQR` provenance to `CUSTOM` for non-diagnostics recommendation IDs.

### Diagnostics provenance

The pinned reference also labels diagnostics recommendations as `AZQR`.

Cloud Assess identifies the real mechanism/source explicitly:

```text
Source: DIAGNOSTICS
Validation mechanism: Azure Resource Manager
```

Known diagnostic recommendation IDs therefore normalize from reference `AZQR` to target `DIAGNOSTICS`.

The reference ImpactedResources table hard-codes `Azure Resource Graph` as the validation mechanism even for diagnostics. Validation-mechanism text is therefore not used as an equivalence field.

### Diagnostics subscription display name

The pinned reference Diagnostics scanner creates primary findings from inventory resources but does not populate `GraphResult.SubscriptionName`. Cloud Assess intentionally fills the already-discovered subscription display name.

For known Diagnostics recommendation IDs, the equivalence harness therefore ignores `subscriptionName` when comparing primary findings. This normalization is deliberately narrow: subscription display names remain strict comparison fields for ordinary Graph findings and auxiliary datasets.

This behavior was confirmed by the first live Azure comparison on 2026-09-22, where all 11 reported finding deltas were this exact reference omission.

### Cost subscription display name

The pinned reference Cost stage does not pass subscription display name into its scanner configuration, leaving the Cost report subscription-name column empty.

Cloud Assess deliberately corrects this by using the already-discovered subscription display name.

Cost equivalence therefore compares subscription ID, period, service, value, and currency, but not subscription display name.

### Azure Policy timestamp

Policy-state timestamps are excluded from semantic equality because they are not the identity or assessment conclusion and may differ between near-consecutive queries.

### Ordering and target-only metadata

The following are not semantic differences:

- JSON field ordering
- row ordering
- Cloud Assess schema version
- generated timestamp
- Cloud Assess scope hash
- Cloud Assess stage start/finish/duration metadata
- Cloud Assess warning/completeness metadata when the assessment remains comparable

## Normalization rules

For identity and Azure-enumeration fields, comparison is case-insensitive and whitespace-trimmed.

Descriptive text has surrounding/repeated whitespace normalized but otherwise retains its content.

Numeric strings such as Cost values and Arc SQL vCore values are normalized before comparison so formatting-only differences such as `12.300000` versus `12.3` do not create false deltas.

Duplicate resource-type counts with identical subscription display name and resource type are aggregated because the pinned reference JSON does not include subscription ID in that table.

## Diff model

Each dataset reports:

- `referenceCount` / `targetCount`: total projected records on each side
- `missingCount` / `extraCount` / `changedCount`: quick triage totals
- `coverageMismatch`: enabled on one side but not the other
- `missing`: present in the reference but absent from Cloud Assess
- `extra`: present in Cloud Assess but absent from the reference
- `changed`: same semantic identity exists on both sides but one or more compared fields differ

A report is equivalent only when:

- both projections are valid for comparison
- stage/dataset coverage matches
- there are no missing records
- there are no extra records
- there are no changed semantic fields

## Delta classification

Every live delta should be classified before code is changed:

1. **Target defect**: Cloud Assess fails to reproduce material reference behavior.
2. **Intentional target improvement**: behavior differs deliberately and is documented.
3. **Source defect corrected by target**: pinned behavior is demonstrably incorrect and the correction is intentional.
4. **Environmental/timing difference**: Azure state changed or APIs returned time-sensitive data between runs.
5. **Access/authorization difference**: the two executions did not have materially identical access.
6. **Deferred/unsupported feature**: the delta belongs to functionality not yet in the current equivalence boundary.
7. **Source ambiguity requiring live evidence**: current example is the Arc SQL `vcores` response shape.

Do not add a new normalization rule merely to make a failing comparison green. Any new ignored/normalized difference must first be justified and documented as an intentional compatibility decision.

## Recommended validation sequence

Use progressively broader scopes:

1. a small stable subscription or resource group
2. the existing non-production Azure test environment
3. targeted Terraform fixtures for behavior absent from that environment
4. multi-subscription scope
5. management-group scope

For each pass, retain:

- reference JSON
- Cloud Assess JSON
- equivalence JSON
- exact source and target commit SHAs
- scope/stage/filter parameters
- classification notes for every non-empty delta

These artifacts form the reproducible evidence set for declaring core-v1 equivalence.
