# Required inventory contract v1

Specification freeze candidate, not implementation or acceptance. The live candidate PR owns publication and OWN native/source QA. Task9cannot start from predecessor green. This consolidates the original provisional contract, input/output ledgers, validation order, scoped documentary reconciliation and deterministic cancellation protocol without rewriting any original50or supplemental26expected results.

## Internal shapes for the strict consumer

```go
type InventoryScopeDeclaration struct { SubscriptionID, Status string }
type RequiredInventoryScope struct { Name, Status string }
type RequiredInventoryResource struct {
 Record assessment.Resource
 SubscriptionID, SourceRegion, ResourceType, Kind string
 Physical bool
}
type RequiredInventoryRequirement struct {
 SubscriptionID, SourceRegion, ResourceType, Kind string
 SKUs []string
 Kinds map[string][]string
 ResourceIDs []string
 MissingSKU, MissingKind, UnrepresentableSKU, UnrepresentableKind bool
 CatalogOnly, Unsupported bool
}
type RequiredInventoryPlan struct {
 Resources []RequiredInventoryResource
 Requirements []RequiredInventoryRequirement
 Scopes map[string]RequiredInventoryScope
 Status string
 // Separate immutable private snapshots and constructor validation state.
}
type RequiredEvidenceRequest struct { SubscriptionID, SourceRegion, TargetRegion string }
type RequiredCollectionInput struct {
 SubscriptionID, SourceRegion, TargetRegion, ResourceType, Kind string
 Storage *StorageSKUCollection
 Disk *ComputeDiskSKUCollection
 VM *ComputeVMSKUCollection
 Cognitive *CognitiveSKUCollection
 SQL *SQLCapabilityCollection
}
type RequiredCatalogBinding struct {
 SubscriptionID, SourceRegion, TargetRegion, ResourceType, Kind string
 Status, DeploymentEligibility string
 Cognitive *CognitiveCatalogEvidence
 SQL *SQLCatalogEvidence
 VMSS *VMSSCatalogEvidence
}
type RequiredEvidenceAdmission struct {
 Request RequiredEvidenceRequest
 ScopeStatus, Status string
 SKUs map[string]SKUEvidence
 Catalogs []RequiredCatalogBinding
 // Separate immutable private request/scoring snapshot.
}
func BuildRequiredInventory(context.Context, map[string]string, []assessment.Resource, []InventoryScopeDeclaration) (*RequiredInventoryPlan, error)
func AdmitRequiredEvidence(context.Context, *RequiredInventoryPlan, RequiredEvidenceRequest, []RequiredCollectionInput) (*RequiredEvidenceAdmission, error)
func (*RequiredEvidenceAdmission) Evidence() (RequiredEvidenceRequest, map[string]SKUEvidence)
```

The package remains internal/plugins/region. These are private coordination interfaces, not report/CLI/public JSON schema changes. Canonical-ID comparisons/sorting preserve original ResourceIDs in diagnostic requirements. Every returned mutable view is owned. Zero/nil admission Evidence returns zero request and an owned nonnil empty map; forged/zero/nil plan admission rejects. No authentication/provenance claim follows from constructor validation.

Private test-core names are buildRequiredInventoryObserved and admitRequiredEvidenceObserved, taking the corresponding public arguments plus final observe func(string). Public wrappers pass nil. These names are test seams for the frozen phase protocol, never exported context options or production caller hooks. Healthy public-wrapper and private-core nil-observer results must independently match the same literal oracle. No global observer. Actual phase-specific/mutation/ownership proof belongs to task9after actual missing APIs are observed on both hosts.

## Normative evidence and precedence

The [descriptor](evidence/20261026/inventory-freeze/contract-descriptor.json) binds exact files/hashes, product185and unchanged AZQR/APRL pins. Current user restrictions and live revision-specific QA own authorization/status. For behavior, this v1shape plus [validation order](INVENTORY_VALIDATION_ORDER_20261026.md), [input ledger](INVENTORY_BUDGET_20261022.md), [output ledger](INVENTORY_OUTPUT_LEDGER_20261026.md) and [cancellation protocol](INVENTORY_CANCELLATION_PROTOCOL_20261026.md) refine the unchanged provisional text. Original expected results remain mandatory. An actual conflict is a gate failure requiring reconciliation, never permission to edit an expectation to fit code.

Descriptor phase labels explicitly map inventory/build to inventory and admission/admit to admission. Old prepared_against/frozen_before_production and supplemental OPEN schedule text are historical preparation metadata, not current authorization or runtime proof. The descriptor binds the resolved phase protocol separately; it does not mark any future test executed. Fixtures require exact keys/types/presence, duplicate/alias/trailing/nonfinite/Unicode rejection before typed decoding. Semantically invalid typed values remain allowed so negative scenarios reach the API. Schema preflight is not API enforcement proof.

All50original scenarios,26supplemental literal scenarios,13procedural obligations, four original input recipes and ten newer budget recipes are required by task9. The sixteen reachable phase schedules cover success at three phases plus error at before-return for both APIs and both context causes. Build strided/dominated budgets retain their disclosed limits; no invented plus-one or RSS claim. Preserve all seven resource types, ALL Storage kinds, exact Cognitive kinds, unresolved SQL coordinates and VMSS unsupported scoring/not_assessed.

## Gate to task9

Before production: freeze/read back this specification and integrity inputs; complete OWN full native/source QA, retain full logs/captures/control conservation and private restore evidence, reconcile handover; then create a product185lineage worktree carrying current guidance and frozen inputs. Publish a strict consumer with no production API, and inspect BOTH actual native missing-API failures before implementation. Preserve first consumer/oracles and existing622Linux/617Windows identities. Compiling mutation evidence must distinguish intended assertions from panic/build/infrastructure/timeout, addressing the open R7-F04assurance obligation before relying on new mutation proofs.

Task8freeze is a specification checkpoint, not whole documentary audit, product acceptance, public region support, live Azure proof, security clearance or release. Human assessment consideration remains after overall project completion. Every ancillary/live/operational/licensing/Gate004/release and confirmed-finding obligation remains. No Azure/live retail/real credentials/fixtures/RBAC/subagents/unrelated refresh/protected merge/bypass/release.
