# Quality Gate 001: Core Foundation Review

Date: 2026-09-16

Status: **PASS for the implemented foundation**

Reference repository: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Validated target code commit: `dc82bc81755acae7ad84759483b59509f0889e2d`

Quality-gate workflow run: `35095316848`

## Purpose

This gate was performed before beginning the next implementation phase. Its purpose was to loop back over the Cloud Assess work completed so far, compare implemented behavior against the pinned reference, identify accidental drift, strengthen automated controls, and distinguish verified behavior from future work.

This is not a declaration of full product equivalence. Diagnostics, Advisor, Defender, Policy, Arc SQL, complete Cost API behavior, report renderers, final CLI orchestration, and live-Azure source-versus-target equivalence remain future work.

## Areas reviewed

- repository and branch health
- scanner registry
- filtering and scope semantics
- Azure resource-ID helpers
- resource discovery
- subscription and management-group discovery abstractions
- recommendation source provenance and loading
- Azure Resource Graph request batching and pagination
- ARG row normalization and recommendation execution
- findings summary behavior
- severity-gate behavior
- Azure cloud selection
- authenticated HTTP retry/throttling behavior
- cost-period calculation
- branding boundaries
- licensing and attribution
- documentation consistency
- CI effectiveness

## Material defect found and remediated

### Scanner/resource filtering scope

The initial target implementation interpreted `include.resourceTypes` as literal ARM resource types. The pinned reference actually interprets these values as scanner/service keys such as `aks`, `ca`, `st`, and `vnet`.

The pinned reference then expands the selected scanner keys into a runtime set of allowed ARM resource types and applies that structural scope both when partitioning inventory and when filtering downstream findings.

The target also initially reapplied only exact-resource and tag filtering to downstream ARG findings, rather than subscription, resource-group, scanner resource-type, and exact-resource structural scope.

This was a material parity defect because it could admit inventory or findings outside the intended assessment scope.

Remediation:

- scanner selection semantics now match the pinned reference
- `include.resourceTypes` is explicitly documented as scanner/service keys
- selected scanners are expanded into a runtime allowed ARM resource-type set
- structural subscription, resource-group, resource-type, and exact-resource scope is applied during inventory discovery
- the same structural scope is reapplied to downstream findings
- explicit subscription/resource-group include precedence is preserved
- parent tag-scope inheritance is preserved
- regression tests cover all of these behaviors

## Additional quality issues found and remediated

### Documentation drift

The target specification previously described exclusion precedence too broadly and the characterization document listed already-implemented behavior as future work.

Both documents were corrected to describe the verified source behavior and current implementation state.

### Malformed ARG row observability

The target tolerated malformed individual ARG rows, but initially discarded them silently. The reference logs this condition.

Cloud Assess now preserves tolerant decoding while surfacing malformed-row counts as structured rule warnings. This is an intentional observability improvement suitable for the target's stage-health/completeness model.

### CI coverage

The earlier workflow ran ordinary unit tests and `go vet`, but did not enforce several important repository invariants.

The quality gate now checks:

- exact pinned APRL revision
- exact imported AOR tree identity
- exact imported custom-rule tree identity
- `gofmt` cleanliness
- clean `go mod tidy` result
- executable-source branding boundaries
- race-enabled Go tests
- `go vet`

The strengthened workflow initially caught four pre-existing formatting issues. They were corrected before this gate was marked passed.

### Third-party licensing

The preliminary NOTICE identified MIT-derived material but did not preserve the full incorporated MIT license texts in a dedicated target artifact.

`THIRD_PARTY_LICENSES.md` now preserves the relevant license texts for:

- Microsoft Azure Quick Review derived/reference material
- Azure Proactive Resiliency Library v2
- Azure Orphan Resources

`NOTICE.md` now links the license artifact and distinguishes repository-level Apache-2.0 licensing from third-party obligations.

A generated dependency-license inventory is still required before an actual release or redistribution.

## Provenance verification

The quality gate verified exact Git identities:

- APRL: `60eaddda76541f6adbc1c5ffa686829807e55e29`
- AOR snapshot tree: `a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5`
- custom-rule snapshot tree: `674b9b3dcb443ce6dc445b48e1db47e4a0ca7082`

The scanner registry also has an automated invariant requiring exactly 87 scanner keys and tests for the specialized workload mappings and dual Redis registrations.

## Behaviors reviewed with no material defect found

The implemented portions of the following areas were consistent with the pinned reference or represented documented intentional improvements:

- default stage enable/disable configuration
- mandatory Graph stage for normal assessment
- resource-ID component parsing
- findings normalization and deduplication
- SLA finding exclusion from the ordinary summary
- severity threshold/gate ordering
- Azure Public, Government, China, and custom cloud selection
- DefaultAzureCredential construction strategy
- ARG request batching at 300 subscriptions
- ARG result page size of 5,000 and skip-token pagination
- management-group-aware ARG authorization option
- unsupported logical-table warning/skip semantics
- bounded rule execution concurrency
- deterministic post-concurrency ordering
- resource inventory ARG projection
- Disabled/Deleted subscription exclusion
- deterministic scope hashing
- previous-completed-UTC-month cost window
- retry and proactive throttling configuration for the implemented HTTP layer

## Intentional differences confirmed

The following differences remain intentional rather than parity defects:

- new Cloud Assess product identity
- neutral `assessment:` configuration root
- `GraphRecommendation` replaced conceptually by `RecommendationDefinition`
- `GraphResult` replaced conceptually by `Finding`
- lower-level operations return errors instead of terminating via `log.Fatal`
- deterministic ordering is added where source map/concurrent ordering was unstable
- malformed ARG rows are surfaced as structured warnings rather than log-only diagnostics
- management-group traversal uses a visited-group guard to avoid duplicate traversal
- the legacy `exclude.services` configuration concept is renamed to `exclude.resources`

## Gate limitations and residual work

This gate certifies only the code implemented at this point in the project.

It does not yet certify:

- live authentication against an Azure tenant
- real Azure subscription or management-group enumeration through final SDK adapters
- end-to-end live Resource Graph behavior
- diagnostics ARM batch equivalence
- Advisor equivalence
- Defender status or recommendation equivalence
- Azure Policy equivalence
- Arc SQL equivalence
- complete Cost Management API equivalence
- internal-plugin equivalence
- final JSON/XLSX/CSV/SARIF report parity
- CLI exit-code parity
- performance parity at large Azure scale
- complete dependency/license inventory for a distributable release

These items are future implementation and characterization work, not failures of this quality gate.

## Gate result

**PASS for the implemented foundation.**

The project may proceed to the next implementation phase from this baseline, provided subsequent work remains subject to the strengthened CI gate and the source-versus-target characterization discipline defined in this repository.
