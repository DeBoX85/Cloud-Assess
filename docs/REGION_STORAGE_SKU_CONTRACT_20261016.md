# Storage SKU collector pre-edit contract, 2026-10-16

Status: pre-implementation proposal. This slice defines behavior and independent expectations; it adds no collector, coordinator or public execution. The user authorized continuation after [PR143](https://github.com/DeBoX85/Cloud-Assess/pull/143)'s verified reconciliation. Start from integration `8e482c8d7c3a3b438ad4688be1a9718c5880a85a`, tree `0f3e658e7cb5d4c239eae0c05b11d4fba298aaba`, not PR141's older code. All audit corrections and independent regressions must remain.

The natural stop is a published/read-back contract and expectation checkpoint with its own full native/source QA. Stop before implementing the new API. The live successor PR supplies its actual head/tree/parents, final outcomes and complete log evidence; statements pending at publication are snapshots, not transferable acceptance.

## Evidence and authority

Primary sources are Azure/azure-rest-api-specs at immutable commit `1cdcc5414764d20e8815554c8b89708c10747c18`:

| File | Verified Git blob | Relevant evidence |
| --- | --- | --- |
| [storage.json](https://github.com/Azure/azure-rest-api-specs/blob/1cdcc5414764d20e8815554c8b89708c10747c18/specification/storage/resource-manager/Microsoft.Storage/stable/2023-01-01/storage.json) | `890f54786ed199f0805e2fa11df085198291a9c2` | Skus_List, StorageSkuListResult, SkuInformation, Restriction, SKUCapability |
| [common.json](https://github.com/Azure/azure-rest-api-specs/blob/1cdcc5414764d20e8815554c8b89708c10747c18/specification/storage/resource-manager/Microsoft.Storage/stable/2023-01-01/common.json) | `3cdf2ad7ddec91e15cdd8070ee3989e2ab6c5c57` | SKU names and tiers; extensible enum metadata |
| [SKUList.json](https://github.com/Azure/azure-rest-api-specs/blob/1cdcc5414764d20e8815554c8b89708c10747c18/specification/storage/resource-manager/Microsoft.Storage/stable/2023-01-01/examples/SKUList.json) | `ea659d8e4c2408b470a31fe1152e3d40c675af43` | Entire248-row example, not a truncated first-page sample |

All three complete retrieved files were independently Git-blob hashed. The unchanged example is retained as [the primary fixture](../internal/plugins/region/testdata/storage-primary-2023-01-01-SKUList.json), SHA256`aceb9c06f233a6db9b4adb726e8f7951037531887dd28c949a24dee5b354830e`. Its exact [MIT license](../internal/plugins/region/testdata/storage-primary-LICENSE.txt) is included. The example has five SKU names, eight name/kind groups,31rows per group and31distinct location labels, including eastus2(stage). It establishes valid example structure, not a current live region catalogue, account-kind support matrix or deployability.

Pinned AZQR remains `8e4f0577f3615e6c9014c031bcad079f235369cc` / tree `17d93b20c303f90f7843036be82f0dc32f3260f1`. StorageProvider, FetchPagedSKUs and CheckStandardRestrictions were read directly. Source uses name-only last-row assignment, ignores account kind and restriction values, falls back to tier for a missing name, treats missing locations as global, and follows generic nextLink. Preserve all original source captures; the corrections below are deliberate target behavior, not regenerated source observations.

Existing target interfaces were inspected directly: inventory.go preserves assessment.Resource.Kind; inventory_aggregation.go discards kind in SKU counters; availability.go consumes name-keyed SKUEvidence with only complete/unknown/unsupported status. No per-SKU unknown state exists there. The [original analysis](REGION_STORAGE_SKU_TARGET_ANALYSIS_20261014.md) remains historical evidence.

## Proposed internal API

These names are a contract for the next implementation, not existing symbols:

```go
type StorageSKURequest struct {
    SubscriptionID string
    Region string
    LocationAliases map[string]string
}
type StorageSKUGroup struct {
    Known bool
    State int
}
type StorageSKUCollection struct {
    SubscriptionID string
    Region string
    Status string
    Groups map[string]map[string]StorageSKUGroup
    UnqualifiedNames map[string]bool
    FailureCode string
}
func NewStorageSKUCollector(endpoint string, getter RESTQuotaGetter) (*StorageSKUCollector, error)
func (c *StorageSKUCollector) Collect(ctx context.Context, subscriptions map[string]string, request StorageSKURequest) (*StorageSKUCollection, error)
func (c *StorageSKUCollection) Evidence(ctx context.Context, requiredKinds map[string][]string) (SKUEvidence, error)
```

Groups are lowercase SKU then lowercase account-kind keys, projected for the collection's one target Region. Known=false always carries State=SKUUnavailable(2), preventing a zero-value SKUAvailable from masquerading as usable evidence. Missing kind/resourceType produces an UnqualifiedNames entry for the admitted name, not a fabricated group. Operationally unknown collections have empty owned maps. No new public result schema or region plug-in registration.

The future coordinator must supply requiredKinds from selected raw assessment.Resource records for the same subscription and source-region slice, preserving kind before inventory aggregation loses it. It must correlate collection subscription/target region with its requested scope, and reconcile required SKU identities with existing inventory counters. It must not infer kind from a SKU name, rendered cell or another subscription. Case-normalized required-kind sets may merge multiple source rows, but existing raw SKU counters and scoring denominators are not changed. Coordinator construction/validation and public integration are separate future work.

## Request and endpoint contract

Constructor delegates HTTPS ARM-origin and nonnil-getter admission to the existing bounded transport contract. No credentials are acquired for invalid inputs. Selected subscription IDs are normalized GUIDs; selected names must be nonempty; scope collisions or limits reject. Request.Region must already be a lowercase physical-format region admitted by diskPhysical, not a display label, global or a geographic aggregate. This validates format, not live region existence.

Exactly one collector invocation of GetBoundedWithResponse is allowed, with bound1048576:

```text
GET {validated-origin}/subscriptions/{normalized-selected-guid}/providers/Microsoft.Storage/skus?api-version=2023-01-01
```

No location filter, region path, provider registration, SDK refresh or follow-up request. Getter contracts still require selected ARM audience, bounded attempt bodies, body closure, cancellation and no redirects. Transport-managed bounded retries, if any, are distinct from the one collector invocation.

The primary operation has only documented200 completion and x-ms-pageable.nextLinkName=null. A root nextLink property of any type, including empty string or null, causes storage_sku_unexpected_paging after lexical/object admission. Never follow it, even same-origin. This is a strict target completeness policy for an undocumented field, not a claim that the schema explicitly forbids additional properties. Compute's documented null terminator does not apply to Storage.

Precedence after getter returns: caller ctx.Err; getter context.Canceled/DeadlineExceeded; present non200 status; remaining getter error/nil response; raw byte limit; lexical/envelope validity; unexpected paging; row admission and grouping. Parent cancellation/deadline wins even over status403. Context errors return nil collection and an errors.Is-compatible context sentinel, without opaque getter text. Other operational faults return an unknown collection and sanitized FailureCode with nil Go error. Invalid caller inputs return nil collection and a sanitized error. No partial status/prefix exists because this endpoint admits one whole envelope atomically.

## Wire and optional-metadata admission

Require one complete UTF-8 JSON object through EOF, canonical known property spellings, no duplicate keys under the existing Unicode SimpleFold equivalence, bounded nesting/fields/tokens and no forbidden control/replacement/noncharacter strings. Unknown additional metadata is allowed only under the same bounds. Do not parse through trailing garbage or silently case-match an aliased known property.

value must be a present nonnull array; every item must be an object. name is required by primary schema and must be a nonempty printable ASCII label, with no leading/trailing whitespace, max512bytes. No fallback to tier. SKU/kind normalization is ASCII lowercase, without trimming or lossy Unicode folding. Names/kinds remain extensible; do not restrict values to the current example's five names/two kinds.

| Field | Missing | Present but null/wrong type | Semantic treatment |
| --- | --- | --- | --- |
| resourceType | Mark name unqualified | Reject whole envelope | Well-formed other resource types are excluded after bounded metadata admission; storageAccounts matches without ASCII case |
| kind | Mark name unqualified | Reject whole envelope | Nonempty ASCII identifier; different kinds never substitute |
| locations | Mark group unknown | Reject whole envelope | Empty array means no physical location advertised; never global |
| restrictions | Mark group unknown | Reject whole envelope | Explicit empty array establishes no declared restrictions |
| capabilities | Allowed | Reject whole envelope | Validate objects/known scalar fields; informational capabilities do not imply global availability |
| tier | Allowed | Reject whole envelope | String metadata only; never SKU identity |

All present list entries and known fields must have correct JSON types; this applies even to filtered rows. Capability name/value and restriction type/values/reasonCode are schema-optional. Missing fields needed to interpret an admitted restriction, or an unnamed/value-missing capability that might encode availability, make its group unknown rather than fabricating success. Unknown bounded metadata otherwise has no availability effect.

## Location reconciliation

Canonical response IDs are ASCII-lowercased and admitted by diskPhysical; they are never certified as actual live locations. Display labels resolve only through explicitly supplied LocationAliases. Alias keys must be unique already-lowercase printable ASCII labels with no edge whitespace, max512bytes; values must be lowercase physical-format canonical IDs. Reject normalized collisions, noncanonical keys, stage keys and any alias that remaps an already-valid canonical ID to a different ID. At most256aliases. This is caller-supplied mapping, not automatic ARM discovery or evidence of live membership.

A response label that fails physical-format admission and has no explicit alias makes the containing group unknown. Do not delete arbitrary spaces or punctuation. Exactly a physical-format identifier followed by literal (stage), without extra text, is recognized as a nonphysical stage label and ignored for physical-target matching; it is never stripped into a real region. This applies to locations and restriction values. Unknown display restriction values remain unknown because their scope cannot be established.

## Grouping and restriction semantics

Group all storageAccounts rows by normalized name+kind before projection. Each group combines the union of physical locations and restrictions. Exact duplicates are idempotent. Regional row order cannot change any result. Repeated names/kinds/regions are not identity errors. No last-row or first-row overwrite. Any semantic incompleteness in a group's admitted rows makes that entire group Known=false, even if another row appears unrestricted. An unqualified row poisons only projection for its SKU name, not unrelated names.

For complete groups, derive State for the target region:

1. If target is absent from the group's union of advertised physical locations, SKUUnavailable.
2. If an applicable location restriction has QuotaId, SKUUnavailable. If applicable reason is NotAvailableForSubscription, SKURestricted unless an unavailable condition also exists. Match documented strings without ASCII case.
3. An available capability with value=false also makes the group SKUUnavailable. value=true adds no location. Other named capabilities, such as supportsnetworkacls=false, do not block availability.
4. Otherwise SKUAvailable. Storage does not infer SKUZoneRestricted from an undocumented Compute restrictionInfo shape.

Restriction values scope applicability to the named locations. Do not apply a westus-only restriction to eastus. Missing/empty restriction values, unknown restriction type/reason, or unresolved location labels make the group unknown. Present malformed restriction fields reject the whole response. Even a restriction on a separate regional row participates in the group's union, so an unrestricted duplicate cannot erase it.

For available capability, only explicit true/false strings are interpretable. Identical normalized duplicate capability entries in one row are idempotent; conflicting values in the same row make the group unknown. Across separate complete rows, a false capability is a conservative group-wide unavailable condition. These are explicit policy decisions; capability availability is inherited source behavior, not documented live capacity certification.

The source's QuotaId->Unavailable and NotAvailableForSubscription->Restricted numeric mapping is retained while correcting restriction location scope. These labels are source-compatible assessment states, not a promise that either restriction is temporary or removable. Unknown metadata always takes precedence over positive or negative state projection in this bounded contract.

## Projection into the current name-only consumer

Evidence receives requiredKinds, not an arbitrary union of all service kinds. Validate labels/bounds, reject colliding normalized SKU keys, and deduplicate identical normalized kind entries. Missing provenance is represented by an empty kind list or an empty kind marker and yields unknown evidence, not a guessed kind. Malformed nonempty identifiers reject with a sanitized input error. Empty requirements produce complete empty evidence only for a complete valid collection.

For each requested SKU, require every inventory-required kind. A missing group in an otherwise complete qualified catalogue is unavailable; a different kind is not a substitute. Known unavailable dominates restricted; restricted dominates available across the required kinds. Do not count kinds separately in the existing SKU denominator.

If ANY required SKU has unresolved kind provenance, UnqualifiedNames membership or a Known=false required group, return exactly SKUEvidence{Status:"unknown", Values:empty}. An operationally unknown/nil collection also projects unknown/empty. This deliberately loses otherwise known SKU results because the existing consumer has no per-SKU unknown state. Never place incomplete values under Status:"complete" and never let a missing map entry turn unknown into confirmed unavailability. Fully known projections contain exactly the requested normalized SKU keys, owned maps, state0/1/2 and empty blocked-zone slices. Never return mixed unknown/complete values or unsupported status for an operational error.

The existing CalculateAvailability contract is retained: complete-empty marks requested names unavailable; unknown marks all names unknown with availability_sku_unknown and completed-with-warnings. With all checks unknown, the retained source percentage is100. That number is not a positive Storage result. The new consumer test checks literal counts, details, health/warning and this exact limitation. Any future percentage/schema redesign is a separate reviewed change.

## Finite bounds and ownership

| Boundary | Limit |
| --- | ---: |
| Selected subscriptions | Existing MaxSubscriptions1000 |
| Raw single response | 1048576bytes |
| Storage rows / distinct groups | 8192 each |
| JSON depth / properties per object | 32 /128 |
| Lexical tokens / aggregate object-array work | 65536 each |
| Decoded scalar/property label | 512bytes, except known unexpected nextLink may be lexically admitted under existing8192-byte guard before unconditional rejection |
| Aggregate decoded/input/output text | Existing auxTextBudget=16MiB-64KiB, with separate bounded counters per phase |
| Location aliases | 256 |
| Required SKU names / kinds per name | 8192 /16 |
| Aggregate projection entries/work | 65536 |

Charge ignored/filtered metadata and duplicate rows before projection. Enforce raw/body/item/work bounds before growing results. Malformed or over-budget response admits no group, including a valid prefix. Invalid output state or oversized manually modified collection must reject safely during Evidence; owned results are not a sandbox for malicious trusted callers.

Inputs must remain stable during a call. Own scope/aliases/call state and all returned maps/slices; no caller-owned map escapes. No shared run cache or mutation of inventory, response fixtures or source captures. Cancellation is checked during validation/grouping/copying and immediately before return; no worker/JSON parsing loop is a universal hard kill for a context-ignoring injected getter.

Failure codes are storage_sku_request_invalid (caller input), storage_sku_status_invalid, storage_sku_request_failed, storage_sku_byte_limit, storage_sku_invalid_response, storage_sku_unexpected_paging, storage_sku_row_limit, storage_sku_work_limit and storage_sku_text_limit. Projection input/output corruption uses storage_sku_projection_invalid or storage_sku_projection_limit. No subscription IDs, URLs, bodies, tokens or opaque transport text in diagnostics.

## Independent expectations and implementation gate

[storage-target-contract-v1.json](../internal/plugins/region/testdata/storage-target-contract-v1.json) contains71independent literal collection/projection scenarios, including complete primary-example maps, order reversal, cross-kind substitution, restrictions, unknown optional metadata, explicit display aliases, stage exclusion, atomic malformed-prefix rejection, paging/status/context precedence and canary failures. Raw body strings preserve duplicate and noncanonical wire keys. Context/getter-failure labels describe future mock behavior; they are not current execution evidence.

The manifest is a pre-edit oracle, not an implemented parser or test runner. Do not claim its71cases passed against a collector that does not exist. This checkpoint adds only an executable three-case consumer-boundary test around unchanged CalculateAvailability: mixed complete, complete-empty and unknown. Existing full suites and source harnesses still run unchanged.

Before collector implementation, add compiling tests consuming every manifest scenario and the exact primary fixture. Publish/inspect the expected compile failure on the missing API as test preparation, not a product pass; then implement production. Require explicit boundary vectors at limit and limit+1 for raw bytes, rows, depth, properties, lexical tokens, total work/text, aliases, required kinds and output ownership. Include invalid UTF8/control/Unicode alias, ignored metadata budget, cancellation during decode/copy and actual bounded-client status/error pairs. Add targeted compiling reversal/restoration controls for grouping, required-kind projection, unknown gating, no-paging, location aliases/stage, restriction scope/precedence, atomic admission, status/context priority and budgets. Mutated production must compile and fail named independent assertions, then restore exact bytes and pass baseline. Never regenerate the manifest from the new implementation or loosen an expectation to obtain green.

Any discovered contract mistake must be reported with primary evidence and corrected in its own documented checkpoint before production, preserving the old expectation and reason in history. Full Linux/Windows/native/source QA, exact-head/readback/preview checks and final bounded self-review are mandatory. No live/Azure, source capture refresh or protected acceptance is inferred.

## Recovery, residuals and rollback

The live draft PR for branch `docs/storage-sku-contract-20261016` records the frozen candidate and own QA. Maintain SESSION_HANDOVER, ledger, roadmap, failure notes, execution and QA plans plus current restart records. After full checkpoint verification STOP before collector implementation. PR141/142/143 stay draft/unmerged at their existing heads. Protected baseline remains a6fee427. No dependency/SDK/source/public-schema refresh or Azure operation.

Rollback the contract/test/fixture additions through a normal targeted successor commit, preserving original audit and integration history. The new primary example and license are separate from all23unchanged AZQR captures. No reset/force/protected rewrite.

Source SKU captures7/7, accepted target adapters1/7, public region unavailable. FullCost/history/retail, six adapters/coordinator/public integration, AR-02/load/live/nonempty/sovereign/visibility/Arc, freshOS, hosted maintenance, module-only advisory and human/license/security/Gate004/release remain open. This checkpoint does not close them.
