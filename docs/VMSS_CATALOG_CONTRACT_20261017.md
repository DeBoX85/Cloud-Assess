# VMSS catalog-only pre-edit contract, 2026-10-17

Status: proposed internal boundary and frozen future expectations, not an implemented API or scoring change. Parent: verified draft PR152/0067731b73cef1f0a9e201b680c5c85fc38f6672. [CURRENT_STATE](CURRENT_STATE.md) and the live task PR own current status and own-head evidence. [Prior reconciliation](VMSS_CATALOG_RECONCILIATION_20261017.md) remains the primary decision record. Accepted adapters stay 1/7; source characterization stays 7/7; VMSS scoring stays unsupported and public region execution unavailable.

## Evidence and scope

Primary Microsoft documentation was rechecked 2026-10-17: [VMSS FAQ](https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-faq) supports the general VM-size relationship; [SKU troubleshooting](https://learn.microsoft.com/en-us/azure/azure-resource-manager/troubleshooting/error-sku-not-available) distinguishes subscription/location/zone restrictions and deployment/capacity failures; [Resource SKUs List 2021-07-01](https://learn.microsoft.com/en-us/rest/api/compute/resource-skus/list?view=rest-compute-2021-07-01) documents the existing subscription-scoped endpoint. Earlier immutable supporting document identities remain in the reconciliation. No primary files, source captures, API versions, modules or reference pins are refreshed.

Inference/decision: node-catalog evidence can inform a VMSS diagnostic, but current inventory cannot establish a complete workload's placement, orchestration/image/network/storage constraints, quota or allocation capacity. A distinct internal catalog type makes this boundary visible. Do not introduce a resource-specific VMSS API, another parser/pager, hidden discovery or live requests.

## Proposed API, absent in this checkpoint

```go
func (c *ComputeVMSKUCollection) VMSSCatalog(
    ctx context.Context,
    selected map[string]string,
    expected ComputeVMSKURequest,
    required []string,
) (VMSSCatalogEvidence, error)

type VMSSCatalogEvidence struct {
    SubscriptionID, Region, TargetResourceType, SourceResourceType string
    Status, DeploymentEligibility string
    Values map[string]SKUAvailability
    Warnings []assessment.AssessmentWarning
}
```

This is a pure projection over the existing VMSS-bound diagnostic collection. It performs zero getter/credential/network calls; collector-backed future scenarios inject the unchanged existing collector separately. SourceResourceType is always microsoft.compute/virtualmachines, TargetResourceType always microsoft.compute/virtualmachinescalesets. Never change a caller's ResourceType to reuse VM scoring. Stable non-racing caller inputs and non-nil context are required.

Every non-error result carries canonical expected subscription/physical region, DeploymentEligibility=not_assessed, allocated Values and Warnings, and exactly this fixed warning:

- code: vmss_catalog_only
- message: VMSS catalog evidence does not establish workload deployment eligibility, quota or allocation capacity.

The distinct result type has no implicit conversion/Evidence method returning SKUEvidence and is not wired into CLI, reports or the scoring coordinator. Existing ComputeVMSKUCollection.Evidence behavior is unchanged: valid complete VMSS collections return unsupported/empty; valid nil/partial/unknown collections retain their existing unknown behavior. Future catalog-complete is not scoring-complete or deployment-ready. No score contribution is authorized by implementing this diagnostic projection.

## Validation and interpretation

1. Context cancellation/deadline wins before validation, unknown/empty return or projection. Return the canonical context cause and a zero result. Repeat checks during bounded traversal and immediately before returning.
2. Validate selected scope with the current UUID/label rules; normalize UUID case, reject normalized collisions, invalid/empty labels and an unselected expected subscription. Expected type must be canonical VMSS; physical target must already be canonical lowercase. Output uses normalized UUID. No arbitrary region aliases, logical/global/display regions or scope widening.
3. A non-nil collection must match expected subscription/region/type exactly after expected UUID normalization. Revalidate all collection names, statuses, failure enums, groups, known/state consistency and sorted unique positive zones using the frozen VM collector contract. Validate even unrequired groups and before empty-requirement early returns. Complete-with-failure, partial-without-failure, unknown-with-groups, invalid status/enum/state/zone/key or VM-bound input rejects, not rebinding. Structurally plausible caller-built data is not authenticated remote evidence; validation does not prove collection provenance.
4. Nil or valid partial/unknown collection gives catalog Status=unknown and no Values. Do not turn admitted prefixes into complete evidence. A complete collection with any required Known=false group likewise gives unknown/empty. Unknown unrequired groups do not poison required known groups, but malformed unrequired groups still reject.
5. Required names follow the existing printable ASCII/unpadded/512-byte policy, normalize case, reject normalized duplicates and allocate an owned map. Empty requirements on a valid complete collection give complete/empty. An absent required name in an explicitly complete catalog maps to state 2 (catalog unavailable), never a deployment assertion.
6. Preserve existing state meanings: 0 available catalog, 1 regional restriction, 2 unavailable catalog, 3 zone restriction. Both QuotaId and NotAvailableForSubscription apply. Unknown knowledge dominates; otherwise explicit region absence/available=false precedes regional restriction, then zonal union, then available. State 1/2 clear zones; zone lists sort lexically and deduplicate at collection admission. Missing restrictions/locations and incomplete capabilities retain uncertainty. Optional absent diagnostic capabilities alone do not manufacture unknown.
7. Reuse the existing strict collector and bounded validation primitives. The collector still admits only virtualMachines wire rows; a virtualMachineScaleSets wire row is not substituted as node catalog evidence. Existing25collector-backed vectors reuse unchanged transport/admission semantics; they are not new live observations.

Bounds: selected scope at most 1000, groups and required entries at most 8192 each, labels at most 512 bytes, logical entries at most 65536, and input text at most 16 MiB minus 64 KiB (16711680 bytes). Charge selected IDs/labels, expected and collection binding strings/failure code, group keys/zones and required-name occurrences before output allocation; charge scope/group/zone/required occurrences toward the logical-entry cap. Validate supplied values before filtering/deduplication. Bound owned output text including fixed warning/binding within the same output allowance. Exact-bound/one-over, late context and ownership tests are mandatory before production completion. These new projection input accounting rules must not silently change the existing VM Evidence budget semantics.

Invalid inputs return a zero result and fixed sanitized error VMSS catalog projection rejected [vmss_catalog_projection_invalid]; work/text/count/label/output bounds use vmss_catalog_projection_limit. Do not echo subscription IDs, labels, SKU names, URLs, error bodies or arbitrary failure codes. Operational collection failures remain diagnostics with unknown catalog, not fabricated successful data.

## Frozen independent expectations

[Target manifest](../internal/plugins/region/testdata/vmss-catalog-target-contract-v1.json): 61 literal scenarios, 25 collector-backed and 36 direct collection/scope/context cases. Bytes: 93385; SHA256: e09f3814ebc647265ea2aa7e8614e855ee24a96e4eefd4fa4897d7d739e881c0. It contains explicit inputs, expected catalog fields/maps/zones/warnings/errors and literal getter URLs/counts. Expected values were authored from this contract and primary evidence without calling the future API or deriving them from production outputs. The integrity test freezes bytes; it does not execute these 61 future scenarios.

Coverage includes both restriction reasons, union/precedence, unknown metadata, unknown required versus unrequired rows, explicit absence, empty requirements, rejected/partial collection, invalid/foreign binding, normalized scope/name collisions, zone consistency, context priority and output ownership. All successful results explicitly state not_assessed and the catalog-only warning. The separate scoring invariant retains zero VMSS denominator contribution and availability_sku_unsupported for valid complete VMSS input; no future projection expectation grants scoring support.

Before implementing the API, publish a runner consuming all 61 unchanged expectations and inspect the expected missing-API preparation failure on Linux and Windows. It must compile/run once the API exists; do not call the absent-API failure a product regression or pass. Preserve the existing 67-case VM oracle and all 385/380 controls. Add independent supplemental exact-bound/one-over, caller mutation, late cancellation and compilation-checked reversal/restoration assertions for binding, knowledge, restriction, sanitization, ownership, limits and unsupported scoring. Preserve original frozen values if a real contract error is later found: record the correction separately before changing expectations.

## QA, rollback and continuation

This checkpoint adds only contract/target expectations and an oracle-integrity test; no VMSSCatalog implementation, report/schema change, network action or score enablement. Run focused region tests/vet, static preflight, local links/whitespace and full own-head native/source QA. Compare all original control identities and reconstruct all 23 captures against pinned bytes. Final exact-revision outcomes belong in the live PR body and durable evidence snapshot, not an evidence-only commit loop.

Rollback is a new targeted revert after dependent-work review; the fixture/integrity test/contract form one unit. Preserve source captures, the old VM oracle and new frozen history. See [rollback runbook](ROLLBACK.md). After this checkpoint's own QA, the next bounded slice is runner preparation then the pure diagnostic projection under this contract. Public coordinator/report integration, scoring enablement and deployment eligibility require separate justification and tests. No Azure, fixtures/RBAC, dependency/action/SDK/source/capture/public-schema refresh, subagents, protected merge/bypass or release. Independent review/backup/live/release obligations remain tracked in [ACTION_REGISTER](ACTION_REGISTER.md).
