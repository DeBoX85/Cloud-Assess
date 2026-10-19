# Review and QA complete; PR140 verified offline, 2026-10-14

PR140 merged at a6fee427e19feeef28691737ea9bc5026086efa0 and completed
candidate AND distinct actual accepted-push verification on both native hosts.
The protected branch still points to that verified merge. Signed commit/raw SHA1,
human identities, tree1b13eec931b510592e42608fe8605ebdc287b0a9 and ordered parents
6f923c1d/01e8de62 matched; protected expected-head merge used active ruleset23890737
and required integration15368 checks, with no bypass.

Three confirmed defects are corrected:
| Severity | Defect | Verified correction |
| --- | --- | --- |
| P2 | ARG accepted a valid JSON prefix followed by extra data | Validate buffered and remaining bytes as JSON whitespace through EOF |
| P1 | Disk collector rejected documented nextLink:null final pages | Accept omitted, empty or null terminators; keep nonstring and unsafe-link rejection |
| P2 | Disk collector lost HTTP status classification from actual bounded-client errors | Preserve terminal context priority, then classify non-200 response status before remaining transport failures |

Each defect reproduced after unchanged production compilation on Linux and Windows.
ARG FN124 includes18bad suffix cases and4compiling controls. Disk FN126/FN127 use
three full immutable Microsoft2021-07-01 examples plus actual bounded-client and
partial-prefix fixtures;23compiling disk controls, restored baselines, vet and
Linux race/native Windows passed. FN125 preserves the UTF-8 mutation diagnostic
failure and its exact harness correction. Full formatter bytes matched both hosts.
No original source capture or independent fixture was rewritten to resemble target
code; only the mistaken null-terminal target expectation was explicitly corrected
from primary evidence.

Complete final candidate QA38500431232:
Linux115645943577/full100868/60steps (conditional failure upload only skipped)/
314controls; Windows115645943804/full81672/45steps/all-success/309controls.
Every earlier291Linux/286Windows identity retained plus23disk controls.
Candidate source38500431285/115645861513/full1360507:7successful steps,
23records315chunks/SHA256/fresh full candidate UTF8 and Gitblob equality.
Final automated marker6159203561 completed on01e8de6 at
2026-10-14T02:23:10.953464Z; no-major-issues comment6159475528 reviewed
01e8de6244. Both original threads resolved with specific correction QA evidence.
This automated review is not independent-human or release approval.

DISTINCT actual merge-head PUSH native38502782779:
Linux115653354934/full100555/60steps (same conditional skip)/314controls;
Windows115653354646/full81129/45steps/all-success/309controls.
All314/309candidate identities independently preserved, not only matching counts.
Actual PUSH source38502782765/115653354694/full1359850:7successful steps,
23records315chunks, each SHA256 and fresh full a6fee427 merge bytes/Gitblobs matched.
Both runs explicitly event=push/head=a6fee427/branch=bootstrap/core-v1/completed
success; final protected ref matched. No candidate/push substitution.
Original AZQR8e4f0577/tree17d93b20/APRL60eaddda, source/SDK/modules and
Go1.26.9 remain pinned. Dependency scan still reports0reachable/0imported and
1module-only advisory; the graph is not declared advisory-free.

Progress measures distinct scopes: source SKU captures7/7; target SKU adapters
1/7 (14.3%), with Compute disks now verified internally. These are not a weighted
overall project percentage. Public region execution is still unavailable.
Remaining6types, coordinator/public integration, FullCost/other adapters,
AR-02/load/fresh-OS/live/Gate004/licensing/human/release work remain open.
Overall release ETA requires those contracts, integration and live/laptop QA;
the completed review/first-disk phase has no remaining acceptance step.

Continuation: full immutable primary Storage2023-01-01 schema/common/example
read and Gitblob checked. The248-row example repeats5names across kinds/regions;
source discovery retains kind, while SKU counts lose it. Storage grouping,
kind-aware projection, display/stage locations and no-schema-paging rules need
a concrete pre-edit contract and independent expectations before production code.
This checkpoint/Storage-analysis branch remains task-owned DRAFT/UNACCEPTED;
its prose records verified PR140 evidence and does not accept new Storage code.
Laptops/Azure scope have not been reverified; no Azure action or live claim.

Historical checkpoints below remain as historical evidence; this outcome governs.

---

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
