# Next region slice: pure service-availability sheets

Status: NOT STARTED. Dependency: accepted primary PR104, source capture PR105 and quota/reservation PR106. Public region-selection unavailable. This is the next bounded contract; no target code is present.

## Source behavior

Pinned AZQR8e4f057 internal/scanners/plugins/region/output/output.go BuildSvcAvailSheets and sku/provider.go plus registered compute/sql/storage/cognitiveservices providers define behavior. Use retained exact source-aux-inputs.json/source-aux-outputs.json, full service-full array and service-nil-inventory/service-no-results null branches. The fixture loader must decode this actual topology (FN069), not treat all branches as one table.

Source sorts unique target regions, resource types, implemented source regions and SKU names lexically. It emits actual Svc Avail <region> sheet names/descriptions and eight columns: ResourceType, ResourceCount, ImplementedRegions, SKUCount, SKU, SKU available, SKU zone-restricted, Service available. Preserve source unsupported-provider N/A, no-SKU N/A, missing-type/service, missing-SKU subset, zone detail stripping and unknown-suffix treatment. Plain RestrictedSKUs is unused by this helper; do not invent output semantics.

Pinned registry keys:

| Provider source file | Resource type |
|---|---|
| compute.go | microsoft.compute/virtualmachines |
| compute.go | microsoft.compute/virtualmachinescalesets |
| compute.go | microsoft.compute/disks |
| sql.go | microsoft.sql/servers/databases |
| sql.go | microsoft.sql/managedinstances |
| storage.go | microsoft.storage/storageaccounts |
| cognitiveservices.go | microsoft.cognitiveservices/accounts |

Read actual pinned files/functions before porting the lookup. Lookup does not call FetchSKUs. No SDK/client/network, availability service decision, full Scan, price or score input is implemented by this pure helper.

## Required target contract before production edits

Define an owned per-run aggregate inventory/comparison interface, selected scope/correlation limits, exact allowed input/count/region/detail types and bounded work/text before allocating sheets. Aggregate resource-type rows span subscriptions; avoid fabricated per-row UUIDs. Do not claim the source's identity-free aggregate map proves foreign-resource confinement; later collection/coordinator ownership must establish selected scope.

Preserve complete source cells and absent empty branches. Canonical region ownership metadata/health and unique IDs must be explicit corrections. Source 31-rune sheet truncation can collide for long labels; decide/reject collisions before rendering rather than silently lose tables. Restriction membership and unsupported registry support stay source-grounded.

Provisional ceilings for detailed review: at most32 targets,8192 total output rows,1000 subscriptions,65536 aggregate detail/map entries,512 UTF8 bytes per label and canonical16MiB text cap including join separators/repetition across targets. Independently calculate work limits and verify exact boundaries; do not reuse a decoded-label-only budget (FN063). These are proposed, not implemented/load-certified limits.

## Acceptance and following work

Compare every cell/order/sheet/description from source-full and nil branches. Add independent supported/unsupported/no-SKU/missing/zone/unknown/case/overlap, malformed/duplicate/foreign metadata where representable, joined/global row/text/work, ownership/concurrency/cancellation and sheet-collision cases. Critical bounds need compiling named controls/restored baselines on both native jobs. Final exact-head source characterization, full Linux/Windows QA/logs, reviewed diff/tree/identity/current base/preview/rules and protected merge remain mandatory. Preserve pins/notices and all live/Gate004/release limits.

After this slice, migrate CostComparison separately. Source sorts region/meter IDs, uses first matching meter metadata, displays positive prices at four decimals and blanks nonpositive/missing prices, with three nil-input branches. Inventory/migration calculations, real bounded request adapters and public/all-format execution remain later tasks. This record does not substitute an invented simplified score or certify retail/sovereign prices.
