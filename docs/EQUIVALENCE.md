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
