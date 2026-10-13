# Compute disk SKU target reconciliation, 2026-10-13

Status: proposed pre-edit contract, UNACCEPTED task-owned WIP. PR139 merged at
6f923c1d and its distinct accepted-push native/source proof is VERIFIED OFFLINE.
Verify live refs and this published concrete contract before production edits.
This document specifies proposed behavior; no target collector exists yet.

First target SKU slice: Microsoft.Compute/disks only. VM/VMSS and all other SKU
families remain deferred. No public region execution, report coordinator, cache,
SDK/module changes, Azure calls, quota/deployability or performance claim.

## Source and primary schema authority

Preserve accepted original AZQR8e4f0577/tree17d93b20/APRL60eaddda captures verbatim.
Original disk request uses2023-03-01; target deliberately selects the documented
Resource SKUs version2021-07-01. The synthetic source request does not establish
service validity or justify changing the original oracle.

Microsoft REST documentation:
https://learn.microsoft.com/en-us/rest/api/compute/resource-skus/list?view=rest-compute-2026-03-02
Immutable Microsoft schema inspected in full:
https://github.com/Azure/azure-rest-api-specs/blob/1cdcc5414764d20e8815554c8b89708c10747c18/specification/compute/resource-manager/Microsoft.Compute/Compute/stable/2021-07-01/skus.json
Gitblob b85286923a6382356ab773d90faac382e73740ec.
Schema has required value array, nextLink paging, resourceType/name/locations,
Location or Zone restrictions, QuotaId or NotAvailableForSubscription reasons,
restriction.values for Location scope and restrictionInfo.locations/zones.
Missing optional schema fields do not establish unrestricted availability.
The stricter target admission below is a deliberate conservative policy.

## Internal request and owned state

New ComputeDiskSKURequest has SubscriptionID and Region. Constructor validates
configured HTTPS ARM origin through existing NewRESTQuotaCollector and accepts
RESTQuotaGetter.GetBoundedWithResponse. No widening of existing quota endpoints.
Collect accepts context, stable selected subscription map and request. Before
authentication, validate unique canonical selected UUID, lowercase canonical
physical-format region through existing regionID/costPhysical and no global/logical
region. First GET path:
 /subscriptions/{lowercase UUID}/providers/Microsoft.Compute/skus
Query is one api-version=2021-07-01 and one location-only filter
 location eq '{region}'
Use url.Values canonical encoding; independent tests assert its exact escaped URL.
No body, includeExtendedLocations, registration, tenant discovery or Azure writes.

Owned result binds canonical SubscriptionID and requested Region. Status is
complete/partial/unknown, Values map[string]SKUAvailability and sanitized FailureCode.
All raw names admitted as bounded printable ASCII nonblank strings with no leading/
trailing whitespace; keys use lowercase name, without source trim behavior.
Reject normalized disk-name collisions within/across pages even when identical.
No borrowed map/slice or shared mutable cache.

Invalid scope/request/origin and cancellation/deadline return error/nil result.
Operational failures are unknown before any admitted page, partial after any whole
admitted page, including admitted empty pages. Preserve previous-page diagnostics;
never partially commit the failing page. Successful explicit empty value completes.
This first supported family never labels HTTP404/405 as proven unsupported.

An Evidence method returns owned SKUEvidence: only complete collections expose
deep-copied values. Partial/unknown project unknown and EMPTY values, as required
by CalculateAvailability. Caller must also bind collection subscription/region
before assigning it to AvailabilityEvidence.SKUs["microsoft.compute/disks"].
No fourth arithmetic status, partial-as-complete or absence-as-success conversion.

## Schema and restriction reconciliation

Use whole JSON admission first: valid UTF8/no invalid-surrogate replacement/
controls/noncharacters, canonical known property spelling, no case-folded duplicate
object keys, bounded depth/tokens/keys/string lengths. Require envelope object,
explicit non-null value array and non-null row objects. Present nextLink must be a
string; omitted or empty ends paging, present null fails.

Every row requires valid resourceType and name and charges raw work/text. Rows
of other resource types are excluded after their identity and JSON admission,
without mapping VM rows to disks. Matching disks require explicit non-null locations
and restrictions arrays. Empty locations lists no physical-region availability.
Accept ASCII canonical physical-format locations only; normalize ASCII case,
reject display strings/global/logical/aliases, charge duplicates before dedup.
If requested region is absent, map SKUUnavailable. This is complete negative
listing evidence, not a claim the SKU does not exist elsewhere.

Restrictions must be non-null objects, known exact type Location/Zone and known
exact reasonCode QuotaId/NotAvailableForSubscription. Missing/null/wrong/unknown
type/reason fails the page; no silently unrestricted unknown future enum.
Location restriction requires explicit scope in values and/or restrictionInfo.locations.
Each supplied array must contain valid physical-format names. At least one supplied
scope must be nonempty; union the supplied location scopes conservatively. If
requested region matches either, set SKURestricted. Empty/unscoped restriction
fails, rather than falsely treating it as another-region evidence.
Zone restriction requires restrictionInfo.locations nonempty and zones nonempty.
Validate and charge every location/zone; zone labels are bounded ASCII positive
decimal identifiers without leading zeros. Other-region zone restrictions are
ignored after full validation. Matching zones are deduplicated and sorted in an
owned blocked-zone slice. SKURestricted wins over SKUZoneRestricted, which wins
over SKUAvailable; SKUUnavailable wins when region is absent. Do not let later
restrictions overwrite an earlier regional prohibition.
Optional metadata follows the explicit matrix below. Its zones do not establish
deployment capacity or alter physical-region membership.

## Bounds, pagination and acceptance

Reuse transport bounds1MiB per page,8MiB successful aggregate,64pages. Add65,536raw
row/location/restriction/restriction-scope/zone/locationInfo/capability entries,
charged before filtering/dedup, and existing auxTextBudget/MaxLabelBytes text guard.
Charge all relevant nested values, including discarded rows and metadata. Existing
whole JSON guard limits32depth,65,536tokens/page,128keys/object independently.
Admission is atomic and checks cancellation throughout and after the final page.

Continuation accepts safe relative/absolute HTTPS references only at exact configured
origin, selected first path and version. Permit at most one unchanged exact decoded
location filter and at most one nonempty skiptoken up to4096bytes. A continuation
may omit the filter as refined below. No foreign subscription, unknown/case-aliased/duplicate
query, userinfo/fragment/opaque/raw or encoded path, redirect or cycle. Compare
canonicalized URLs before the next authenticated request. The 64th final page
completes; page65 continuation is partial without another request. Unsupported
continuation is failure, not silently successful truncation.

Independent tests use literal maps/states/blocked zones/full URLs and full arithmetic
records: healthy empty, other types, missing requested region, restrictions.values
without restrictionInfo, conservative union, both reasons, multiple restrictions,
region-over-zone precedence, unrelated-region zones, duplicate-zone raw charging,
canonical name collision, malformed/null/missing/alias/Unicode/extended metadata,
first/later error/page atomicity, all HTTP/cancellation/transport-close outcomes.
Test exact and first-excess each byte/page/work/text/name/token/depth/key/skiptoken
budget with no tighter guard masking the assertion. Real bounded ARM transport
fixtures assert token audience, GET/body/path/query, disabled registration/redirects,
response closure and no request before rejected scope/continuation.
Verify returned map/slice isolation across sequential/concurrent runs and Evidence
projection mutation. Stable-input callers are required; no racing caller-map claim.

Named compiling production mutations must cover restriction-values scope, both
reason codes, region/zone precedence, type filter, collision, partial projection,
atomic page, scope/region/origin/version/filter, redirects/cycles/bounds/cancellation/
response closure and output ownership, restoring baseline after each control.
Preserve all existing source oracles/pins/native/security checks. Require fresh
full Linux/Windows/source/review/file/tree/identity/protected expected-head merge
and distinct actual accepted-push proof. Synthetic listing evidence never proves
live deployability, quota, zonal capacity or current service behavior.


## Optional wire metadata admission matrix

This matrix closes the optional-shape detail before production edits. Mandatory
disk name/resourceType/locations/restrictions and restriction semantics above stay.
Omitted optional metadata is allowed. Present null or a wrong declared shape fails
the whole page. For every known object, reject case-fold aliases of canonical wire
property names before decoding. Unknown properties remain ignored semantically
only after complete JSON validation and bounded recursive work/text admission.

| Property/context | Supplied admission | Projection |
| --- | --- | --- |
| SKU tier/size/family/kind | Bounded printable string, empty allowed | None |
| SKU apiVersions | Non-null array of bounded printable strings, empty allowed | None; it cannot override the pinned request version |
| SKU capacity | Non-null object; supplied minimum/maximum/default are int64 JSON integers; supplied scaleType is exact Automatic/Manual/None | None; no capacity/quota inference |
| SKU capabilities | Non-null array of non-null objects; supplied name/value are bounded printable strings | None |
| SKU costs | Non-null array of non-null objects; supplied meterID/extendedUnit are bounded printable strings; supplied quantity is int64 JSON integer | None; no price inference |
| SKU locationInfo | Non-null array of non-null objects; supplied location is a canonical physical-format region; zones is a non-null zone-label array; extendedLocations is a non-null bounded string array; supplied type is a bounded string | None; only SKU locations establishes physical-region membership |
| locationInfo zoneDetails | Non-null array of non-null objects; exact wire Name is a non-null zone-label array, supplied capabilities follows the capability shape | None; lower-case name is an alias of Name and is rejected |
| Restriction values | Supplied non-null string array; Location entries are canonical physical-format regions; Zone entries are bounded printable strings with no additional inferred meaning | Location scope only |
| restrictionInfo | Supplied non-null object; supplied locations is a non-null physical-region array; supplied zones is a non-null zone-label array | Both required nonempty for Zone; Location needs nonempty scope from values and/or locations |
| Unknown fields, including extended metadata | Whole JSON admission and recursively charged raw arrays/objects/strings | Never interpreted as physical availability |

All supplied metadata strings are charged before normalization/dedup/filtering.
For metadata arrays charge every entry, including ignored/duplicate values.
Recursively charge each JSON object and array element throughout every row, including
unknown fields and filtered non-disk rows, in the same65,536aggregate raw-work budget;
object property keys and their strings also count toward decoded text. The mandatory
whole JSON token/depth/key/UTF8/string bounds run before typed page projection.
The per-call work and text counters are committed only with an admitted page;
invalid pages are not retried and no next page is fetched after failure.

None of the irrelevant int64 metadata values is used for availability or clamped.
Non-integer/overflow values fail shape admission; negative metadata integers are
retained only as ignored wire data, never interpreted as a quota or capacity.
Strings in non-projected metadata use existing bounded UTF8 printable admission;
only identity/location/zone strings require the stricter ASCII grammar.
For filtered non-disk rows, require canonical resourceType/name and complete JSON/
recursive raw-work/text bounds; irrelevant disk-specific field shapes do not decide
disk availability. This distinction needs an explicit independent fixture.

## Continuation refinement

The first request always sends the exact region filter. A same-origin/path/version
continuation may omit the filter, because each row is independently checked against
the requested region. If present, it must appear once with the exact first decoded
filter value. Permit at most one nonempty bounded skiptoken; reject every other
query property, case alias or duplicate. Retaining a server continuation without
filter does not make a nonmatching returned location available. This is an explicit
target decision, separate from unsafe original SDK continuation observations.
Canonical URL comparison still detects cycles before another authenticated GET.

Independent first literal URL for public ARM/eastus/subscription
11111111-1111-4111-8111-111111111111:
https://management.azure.com/subscriptions/11111111-1111-4111-8111-111111111111/providers/Microsoft.Compute/skus?%24filter=location+eq+%27eastus%27&api-version=2021-07-01
This is the proposed url.Values encoding, not a retained original-source request.


## Result and diagnostic details

Collection values and every blocked-zone slice are owned; complete/partial/unknown
results use an allocated empty map when no rows were admitted. Evidence() returns
a new owned snapshot on every call. A nil collection or an unrecognized status
projects unknown with an allocated empty map. Complete empty remains complete.
Stable collection inputs are required during snapshotting; no racing caller mutation
safety is claimed. The later coordinator must bind the request subscription/region
before using the snapshot and must not treat retained diagnostic rows as complete.

Use sanitized fixed codes for request/status/bytes/schema/identity/collision/text/
work/pages/unsafe-continuation/cycle failures, without response/error-derived text.
Cancellation or deadline, including a getter error wrapping either context error,
returns nil result and the context error, even after earlier admitted pages.
The maximum nextLink wire string is8192bytes. It charges the total decoded text
budget while using that separate limit rather than MaxLabelBytes; the decoded
skiptoken remains capped at4096bytes. First/continuation method/body/version/audience
checks and every response close are independently asserted with the actual bounded
transport before complete evidence is projected.


## Independent first literal cases declared before target code

For all cases use the selected UUID above and requested eastus. Output expectations
are literal proposed target results, not observations of an implementation.

| Wire rows or page outcome | Complete collection values | Snapshot status/values |
| --- | --- | --- |
| value=[] and no nextLink | {} | complete / {} |
| disks, name=Premium_LRS, locations=[eastus], restrictions=[] | premium_lrs: state0, blockedZones=[] | complete / full owned map |
| Same disk, locations=[westus], restrictions=[] | premium_lrs: state2, blockedZones=[] | complete / full owned map |
| Same eastus disk, Location/QuotaId restriction values=[eastus], no restrictionInfo | premium_lrs: state1, blockedZones=[] | complete / full owned map |
| Same eastus disk, Location values=[westus], restrictionInfo.locations=[eastus], reason=NotAvailableForSubscription | premium_lrs: state1, blockedZones=[] | complete / full owned map |
| Same disk with two matching Zone restrictions, zones=[2,1] then [1,3], both scoped eastus | premium_lrs: state3, blockedZones=[1,2,3] | complete / full owned map |
| Same disk with matching Zone restriction then matching Location restriction | premium_lrs: state1, blockedZones=[] | complete / full owned map |
| virtualMachines row named Standard_A0 with valid identity and otherwise irrelevant fields | {} | complete / {} |
| First admitted available disk page, second malformed page | diagnostic prefix premium_lrs: state0, blockedZones=[]; collection partial | unknown / {} |
| First admitted empty page, next request fails | diagnostic prefix {}; collection partial | unknown / {} |
| Cancellation after an admitted page | nil result and context error | no snapshot accepted |

Blocked zones belong only to state3. Regional restriction/unavailability must clear
the zone slice so irrelevant earlier restrictions do not leak into final evidence.
Additional literal shape/null/alias/collision/bounds/transport cases remain required
before acceptance. Tests must assert complete maps, request URLs/counts, collection
binding/health/failure codes and owned snapshot fields, not just result length.
