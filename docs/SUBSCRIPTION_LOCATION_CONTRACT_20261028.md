# Subscription locations and zone evidence contract

This is a new task10 prerequisite on a PRODUCT successor of verified PR208/69637552e1875ff533a81f8e762ff5dc52a71d01. No location collector or coordinator is implemented at this preparation checkpoint. Publish and read back the independent consumer and bound files before production. Missing APIs are compilation observations, never runtime mutation proof.

## Bounded batch

1. Audit live206–208 heads, own QA indices and private retention; inventory actual worktrees/tools/uncertain operations and preserve unrelated work.
2. Bind the pinned locations/version/physical-region/fallback behavior and explicit target corrections below.
3. Freeze independent wire literals, boundary recipes and strict Go consumer; publish/read back before production.
4. Observe actual missing APIs on BOTH native hosts, then implement on a separate PRODUCT successor without changing frozen expectations.
5. Verify focused lifecycle/ownership tests, strict compiling source-assertion reversals, OWN full BOTH-native/source QA and coherent rollback.
6. Retain new raw evidence privately on latest-main with additive expected-head updates, full readback/fresh restore, then reconcile progress/handover/register and self-assess. No task10 closure follows from this prerequisite alone.

## First-party and pinned source

AZQR8e4f0577f3615e6c9014c031bcad079f235369cc: selection.go2997f4c3d3617dc5a2e4d2330c5151e6e207f414 and internal/az/locations.go201deebe626bdf4d2fd6346ea79af3372970df41. Source filters metadata.regionType exactly Physical, lowercases/sorts names, keeps subscription-scoped logical-to-physical maps, and uses the unchanged33-entry static count table only when no physical region has ANY API mapping. Fallback supplies counts, never mappings. Mixed mapped/unmapped regions do not get per-region fallback. Parser/orchestration do not traverse nextLink; current same-version official schema declares it. This is a static completeness observation, not demonstrated Azure truncation or severity.

- https://learn.microsoft.com/en-us/rest/api/resources/subscriptions/list-locations?view=rest-resources-2022-12-01
- https://learn.microsoft.com/en-us/azure/reliability/availability-zones-overview

A subscription's mappings can differ from another subscription. A count, mapping or location membership does not prove service/SKU availability, quota or deployment eligibility. No Azure/API-version/SDK/source-pin/static-list/oracle refresh is authorized.

## Exact private API

The package is internal/plugins/region. NewSubscriptionLocationCollector(endpoint string, getter RESTQuotaGetter) returns (*SubscriptionLocationCollector,error). Origins follow the existing bounded ARM origin contract; nil and typed-nil getters are rejected. No client, credential, Context or run cache is created/retained.

SubscriptionLocationRequest contains SubscriptionID string. Collect(ctx context.Context, selected map[string]string, q SubscriptionLocationRequest) returns (*SubscriptionLocationCollection,error). Validate and own the selected scope before I/O; normalize request UUID to lowercase; require selected nonempty scope names. Nil/invalid caller/collector and cancellation return nil result and sanitized error. Wrapped dependency cancellation/deadline errors become canonical context errors; caller cancellation wins.

SubscriptionLocationEvidence contains SubscriptionID, Status, FailureCode, CountSource strings; Regions []string; ZoneCounts map[string]int; ZoneMappings map[string]map[string]string. Collection exposes the same fields as diagnostics but Evidence() returns a separate owned private snapshot. Every mutable map/slice, including empty views, belongs to the call/accessor. Nil/zero/fabricated collections return unknown/incomplete with empty unbound data, never complete. Diagnostic/accessor mutation cannot promote status, redirect subscription or change future runs. No authenticated Azure provenance is claimed.

Status is complete after the entire admitted traversal, partial after at least one admitted page followed by an operational failure, unknown before any admitted page. A successful explicit empty value page is complete. Error text/status/body/URLs/IDs are not copied into new diagnostics. Operational failure returns collection plus nil error, explicit bounded FailureCode and only previous admitted pages. Cancellation returns nil/error and no later request or successful publication.

CountSource is api for a complete traversal with any physical mapping, fallback for a complete traversal with none (including empty), incomplete for partial/unknown. Complete api counts equal each physical row's mapping length, including explicit zero for unmapped peers. Complete fallback counts use exactly zone-counts-pinned.json, zero for unlisted physical names, and mapping output stays empty. Never infer fallback from a partial prefix. In incomplete results counts/maps describe only admitted API prefixes; they are not complete/fallback eligibility. The future coordinator must retain status/provenance and cannot treat incomplete zero as observed absence of zone support.

## Request and pagination

First GET is exactly /subscriptions/{normalized UUID}/locations?api-version=2022-12-01. Use existing RESTQuotaGetter with1MiB response bound. Getter must bound every attempt, use selected ARM audience, close bodies, propagate Context and never redirect. Offline fake getters do not certify production transport or caller visibility.

Continuation may be relative or absolute HTTPS with exact configured origin (host compared case-insensitively), exact first subscription/path, no encoded path/userinfo/fragment/opaque URI, exactly one api-version=2022-12-01 and at most one nonempty bounded $skiptoken. Other/aliased/duplicate query keys, control/noncharacter/invalid UTF8 decoded tokens, unsafe URLs and canonical cycles fail closed before another call. Maximum link8192bytes/token4096bytes. Preserve decoded token meaning by canonical encoding. Unknown service continuation formats yield partial/unknown, not successful truncation; no universal backend-format claim. No registration, writes or credential requests.

## Wire admission and identity

Require one JSON object and explicit non-null value array. Require non-null location objects with canonical name and metadata object containing nonempty string regionType; Physical is projected, every other nonempty bounded string is ignored only after row/mapping validation. Names are nonempty bounded ASCII region identifiers normalized to lowercase, without trimming/display-space repair. Nonphysical names such as global are valid raw names but not physical output. Preserve source Physical-only selection; no live registry or future-enum support claim.

availabilityZoneMappings may be omitted or null, matching source optional-array behavior; both mean empty/unobserved, not a mapping. A present array must contain non-null objects with canonical nonempty logicalZone and physicalZone bounded text. These strings are opaque: preserve case/Unicode distinctions, do not infer numeric zone names or physical-ID equality. Reject repeated exact logical keys, even identical values; distinct logical keys mapping to one physical value are not rejected by an invented bijection rule. Unknown fields are ignored only after whole-JSON validation.

Reject duplicate normalized location identities across all raw rows/pages, including nonphysical rows. Each page is validated against prior state and its own identities before committing anything. Invalid later page retains previous pages only. Optional id, if present, must be non-null/nonempty string with fixed ASCII /subscriptions/{selected UUID}/locations/{normalized row name} segments before case-insensitive equality; no foreign/encoded/malformed identity. Optional subscriptionId must be valid selected UUID. Missing optional identities are not fabricated.

Whole-JSON guard rejects invalid UTF8, controls/noncharacters/replacement runes, exact decoded duplicate member names (including escapes), trailing JSON, depth beyond32, more than65536 value nodes or128keys per object. Root depth0; value-node accounting includes containers/scalars, excludes keys/closing tokens. Strings/keys at most512bytes except root nextLink8192 and location id2048. Typed property sets are envelope(value,nextLink), location(name,metadata,availabilityZoneMappings,id,subscriptionId), metadata(regionType), mapping(logicalZone,physicalZone). These fields require canonical property spelling and reject Unicode/case aliases even alone. Opaque unknown maps retain case-sensitive distinct keys; do not reuse quotaJSON's global folding policy. Existing contracts/callers are unchanged. Ignored numeric literals use UseNumber and need no float conversion.

## Deliberate bounds and failures

Reuse64admitted pages/1MiB per getter result/8MiB accumulated in-bound returned bodies. Stop before call65. At most8192raw work entries per collection: charge one per location plus one per mapping entry before physical filtering/deduplication, including ignored rows. Text uses the existing projected auxiliary text budget, charged before retaining fields. Whole JSON/page limits remain independent; earliest exceeded guard wins. Page/work/projection state is admitted atomically. No universal wire-retry/RSS/heap/total-lifetime/timer claim; an uncooperative getter cannot be forcibly stopped safely by this sequential caller.

FailureCode is locations_request_failed for non-context getter error; locations_status_invalid for missing/non200 response; locations_byte_limit for per-body/aggregate overflow; locations_page_limit before excess page; locations_work_limit for raw-work overflow; locations_invalid_response for JSON/shape/field/identity/text failure; locations_identity_duplicate for normalized row/exact logical-key collision; locations_unsafe_continuation for invalid link; locations_continuation_cycle for repeated canonical link. Initial CountSource remains incomplete on failure. No failure is relabeled complete emptiness. Check Context before/after getter, while decoding/deriving/copying and immediately before publication.

## Required assurance and remaining scope

Independent consumers must cover literal fallback/mixed/page-boundary cases, optional/invalid wire identity, exact/escaped/typed-alias versus opaque-map behavior, atomic duplicates, status/error/cancellation, origins/continuations, at-limit/plus-one body/page/raw-work recipes, selected-map/diagnostic/accessor ownership and concurrent runs. Freeze/readback before production; strict compiling assertion-specific reversals must restore source/pins/frozen tests/checkout exactly. No compiler/panic/timeout credit as runtime mutation evidence.

Retain all32original frozen files/44legacy/49current obligations, nine207bound files, existing controls and all pins. Never run prepare_literals.py. Exact coordinator request/result, currency/ownership and cumulative retained-work contracts remain before coordinator production;11–15 and ancillary/documentary/independent/live/operational/security/licensing/Gate004/release stay separate. Human security assessment considered AFTER overall completion. Issue189/DC-CI-EVIDENCE-20261022 remains closed bounded migration with ongoing private latest-main additive lease/preservation/readback/fresh-restore/revisit/retention-until-explicit-Denis-deletion duties. No automatic archival/independent-backup/DC006closure or blanket security/ETA/percentage/background claim.
