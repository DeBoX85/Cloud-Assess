# Storage SKU target analysis, 2026-10-14

Status: analysis only; no target Storage collector or complete pre-edit contract.
Follow the verified disk acceptance checkpoint before further production edits.
Original7/7 source capture coverage remains separate from target acceptance.

Primary Microsoft specs independently retrieved in full at immutable commit
1cdcc5414764d20e8815554c8b89708c10747c18:
- storage2023-01-01/storage.json blob890f54786ed199f0805e2fa11df085198291a9c2
- common.json blob3cdf2ad7ddec91e15cdd8070ee3989e2ab6c5c57
- examples/SKUList.json blobea659d8e4c2408b470a31fe1152e3d40c675af43
All full UTF8 Gitblob identities reproduced. Exact approved source request version
2023-01-01 has matching primary schema; synthetic source success still does not
certify service validity, sovereign/live behavior or deployability.

GET /subscriptions/{subscriptionId}/providers/Microsoft.Storage/skus is global
per selected subscription, with no region-filter parameter. Operation has
x-ms-pageable nextLinkName:null; response definition includes only value.
Do not blindly reuse Compute continuation/version/filter rules.

Primary SkuInformation requires only name. resourceType usuallystorageAccounts,
kind is Storage/StorageV2/BlobStorage/FileStorage/BlockBlobStorage, tierStandard/
Premium, locationsarray, name/value capabilities, restrictions type documented
as location with values and QuotaId/NotAvailableForSubscription reasons.
Missing schema-optional locations/restrictions do not prove global unrestricted
availability. Original storage-missing-locations capture marksStandard_LRS
available, a target correction must be separately declared, preserving oracle.

Primary example has248rows but only5names. Standard_LRS/GRS/RAGRS each62rows
across Storage and BlobStorage; ZRS/Premium31rows. Same name/kind repeats by
location. A name-only duplicate-reject/last-wins rule would reject normal example
data or lose valid regional information. Define grouped identity/aggregation
and inventory kind provenance before claiming kind-qualified availability.
Existing InventoryCounts.SKUsByTypeAndRegion and AvailabilityEvidence.SKUs
carry name keys without kind; any unsupported/ambiguous scope must remain honest.

Schema location description permits display labels (e.g. East US), while the
example mostly uses canonical IDs and includes eastus2(stage). Define explicit
validated region reconciliation and nonphysical-stage handling; do not use an
unbounded/ad-hoc string replacement or assume a static format is a live catalogue.

Remaining concrete decisions: whole-response/row schema and null policy; grouping
and duplicate semantics across regions/kinds/restrictions; absent/malformed versus
complete empty; exact response/status/error/context classification; no undocumented
paging; fixed raw/decoded/item/byte limits; ownership and scope; literal positive
complete-map/primary-example/scoring cases plus compiling independent controls.
No production code until concrete contract and independent expectations published.
No original SDK/module/pin/capture changes, public activation or Azure writes.


## Verified inventory-kind provenance

Accepted discovery already projects ARG kind into a typed inventory row and copies
it to assessment.Resource.Kind (internal/discovery/resources.go). Region inventory
renders resource.Kind and charges it among bounded labels (inventory.go).
InventoryCounts.SKUsByTypeAndRegion and current SKUEvidence.Values still aggregate
by SKU name without kind. Raw kind is available; a future coordinator must preserve
it for a scoped, bounded per-SKU required-kind projection rather than reconstruct
kind from a SKU name or a rendered/redacted cell.

Proposed reconciliation to develop before production edits: collect and own typed
SKU+kind groups, combine legitimate same-group regional rows, and project only
against the source resource kinds actually required for each SKU. Availability
across different required kinds needs an explicit conservative aggregation rule;
missing/ambiguous kind needs honest unknown evidence. This is a design option,
not an implemented interface, accepted contract or verified availability promise.
