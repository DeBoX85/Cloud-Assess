# Current security gate state, 2026-10-09

PR132 is UNMERGED and protected acceptance is BLOCKED by a classified vulnerability
gate failure on the unchanged Go1.26.8/x-net0.58.0 toolchain/dependency graph.
Read [HANDOVER_SECURITY_GATE_20261009.md](HANDOVER_SECURITY_GATE_20261009.md) and
[live PR132](https://github.com/DeBoX85/Cloud-Assess/pull/132) for exact failed
revision/preview/run/job/advisory evidence and latest published handover state.
Prior PR131 zero-reachable/zero-imported scans are historical, not current clearance.
No suppression, version change, protected merge or accepted-push proof occurred
in this documentation task. Finish security classification and prepare a bounded
remediation contract before further feature work. Source pins and mandatory gates
remain intact; no laptop/Azure input is needed for that offline preparation.

---

# Current handover authority, 2026-10-09

Verified implementation: [PR131](https://github.com/DeBoX85/Cloud-Assess/pull/131),
accepted `b653a3abfc35590a185531095a168a9a107e769e`, tree
`47515c474ddf400da0a68c556e9147d8d3c2576b`. Live refs can advance.
Read [HANDOVER_CHECKPOINT_20261009.md](HANDOVER_CHECKPOINT_20261009.md) for exact
parents/identity/pins, separate candidate/accepted-push proof, audit/findings
disposition, remaining tasks, access limits, workspace recovery and QA provenance.
This documentation checkpoint's final acceptance evidence is in its live PR body
on `docs/handover-20261009`; do not assume it is already accepted.

Owned reservation calculation was accepted through PR128. PR131 materially
corrected Unicode duplicate identity admission in this internal calculator
(FN097), preserving arithmetic and displayed ID normalization. All 18 runtime
controls passed on both hosts in separate candidate and accepted-push runs.
Pure/runtime acceptance is not public region execution or live resource ownership.

Earlier records below are preserved historical snapshots, superseded only in
current-state/resume instructions by the verified checkpoint above.

---

# Owned capacity reservation calculation contract

Status: IN PROGRESS, implementation candidate in
[PR128](https://github.com/DeBoX85/Cloud-Assess/pull/128). Pre-edit planning
215ada9405ce44cfe44d504cf022eb4c219ae3dd/tree712eeafa5cc6263e564f2d59f914e951745e291c
was remotely verified before production edits. Accepted PR127 baseline is
`2d641484017c92991e626683c63f5c8db35dfcdc`, tree
`f3eb3b0e04c5fb936462f5b733dcf3d1fe43c8ef`.
[PR127](https://github.com/DeBoX85/Cloud-Assess/pull/127) retains final candidate
and distinct accepted-push acceptance; those checks do not accept this proposal.
This is the separate B5 reservation calculation step under
[TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md) and
[DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md).

## Source, scope and before/after

AZQR remains8e4f0577f3615e6c9014c031bcad079f235369cc, tree
17d93b20c303f90f7843036be82f0dc32f3260f1; APRL remains
60eaddda76541f6adbc1c5ffa686829807e55e29. Source FetchReservations is
internal/scanners/plugins/region/crg/crg.go, blob
1ffb6e37cc1c21f6d8b5c9ce483f793387e0e6bc.
[REGION_RESERVATION_SOURCE.md](REGION_RESERVATION_SOURCE.md) retains actual
unchanged SDK observations; its historical pending wording is superseded by
[accepted PR126](https://github.com/DeBoX85/Cloud-Assess/pull/126).

Before, ProjectReservations formats ten already-calculated cells and cannot
distinguish an absent utilization field from known zero or a complete inventory
from a source Get/list failure. After, CalculateReservations takes explicit
requested selected subscription/region evidence, validates decoded identities and
presence, computes bounded source-valid counts/status and returns owned typed
reservations plus separate health and the existing table. No network, credentials,
SDK collector, cache, clock, coordinator or public CLI is added.

## Proposed interface and evidence health

ReservationRequest holds SubscriptionID and Region. This is an output-evidence
scope, not a promise that the source subscription-wide group API filters regions.
Future collection must reconcile whole-subscription traversal/completeness and
group locations before supplying requested region evidence.

ReservationEvidence holds Request, Status and Reservations. ReservationUsage
holds ResourceID, Region, ResponseName, ResponseRegion, SKU, Reserved, Allocated,
ReservedKnown and AllocatedKnown. Presence markers are mandatory evidence, not
zero-value assumptions. ReservationCalculation returns owned Reservations
(typed ReservationValue records), Health and Table. ReservationValue retains
canonical ResourceID and selected identity, labels, computed counts and Status.

Selected UUIDs are canonical lowercase, names nonempty and existing safe labels;
case-colliding selection rejects. Regions retain the canonical physical pattern.
Duplicate normalized requests/evidence and evidence outside requested scopes
reject. Missing expected evidence warns, rather than becoming complete empty.
Statuses are exactly complete, partial, unknown and unsupported. The latter two
cannot carry records. Complete empty is healthy; partial can preserve valid prefix
with health warnings. Status is completed or completed_with_warnings. Records is
emitted reservation count. Empty output has a nil Table but retains Health.
Nonempty table health must copy warnings independently of calculation health.

## Structural identity and deliberate corrections

Every supplied ResourceID must be a bounded structurally exact ARM reservation
path, containing these ten segments after its initial slash:
subscriptions / UUID / resourceGroups / RG / providers / Microsoft.Compute /
capacityReservationGroups / group / capacityReservations / reservation.
Fixed segment matching is case-insensitive; UUID must match the requested selected
subscription. No trailing slash, extra/missing segment, query, fragment, backslash,
percent encoding, dot segment or empty name is admitted. Each variable label uses
the existing 512-byte UTF8/control/noncharacter checks. ResourceID is at most
2048 bytes. This is structural confinement, not a claim to validate every Azure
service naming rule or arbitrary future request destination.

Fixed keywords additionally retain their exact ASCII byte lengths, so Unicode
lookalikes do not become alternate provider/type paths through Unicode folding.
At current variable-label caps a structurally admitted ID is at most1680 bytes,
below the defensive2048-byte outer bound. Do not invent an exact2048-byte valid-ID
fixture; the actual per-label boundary is tested.

Duplicate canonical case-insensitive reservation IDs reject across all evidence,
including records that would otherwise be skipped for unknown fields. Region is
required, canonical and equal to its requested evidence region. ResponseName is
optional; when present it must match the ID reservation segment case-insensitively.
Its original display case is preserved. Otherwise display falls back to the ID
reservation segment, a deliberate correction to source empty returned names.
ResponseRegion is optional; when present it must be canonical and equal to Region.
Conflicting identity/region is an error, not a silently retained wrong-resource row.
RG and group display labels come from the ID, which must have already been
reconciled with the parent/request by future collectors. This pure calculation
does not prove raw SDK/HTTP response/request identity verification occurred.

Both known counts are in [0, 1000000000000]. Unknown counts must carry zero or
reject as contradictory evidence. If either count is unknown, skip computation
with a fixed unknown-utilization/capacity warning. Never turn missing expanded
utilization into Idle. Missing SKU likewise skips with a warning because the
existing table cannot display an empty required SKU. Validate all raw identities,
texts, counts and duplicates before these skips or output allocation. Negative
counts and malformed/out-of-scope identity fail with nil result and a generic
sanitized error. The future adapter must derive allocated count from actual valid
VM evidence; this calculator does not inspect a VM list or certify null elements.

## Source-valid arithmetic and output

Available = reserved - allocated. Keep exact source precedence:

| Reserved / allocated | Available | Status |
|---|---|---|
| 4 / 0 | 4 | Idle |
| 4 / 1 | 3 | Available |
| 4 / 4 | 0 | At-Capacity |
| 4 / 5 | -1 | Over-Allocated |
| 0 / 0 | 0 | Idle |
| 0 / 1 | -1 | Over-Allocated |

No current/limit/headroom or quota threshold is reused for reservations. Preserve
declared request order and received record order. Keep existing ten columns:
Subscription, Region, Resource Group, CRG Name, Reservation Name, SKU, Reserved,
Allocated, Available, Status. Counts use canonical base10 strings. The existing
ProjectReservations remains unchanged; generated cells satisfy its contract.
Typed ResourceID is retained out of band; this adds no serialized public schema.

Cancellation before admission, during every traversal and after projection returns
nil result and the context error. Inputs remain stable during a call. Returned
typed records, table cells and warnings are independently owned; repeated and
concurrent calls do not share state.

## Work, text and intended paths

Requests and aggregate raw reservation records each have MaxAuxRows=8192 limits;
selected scope retains MaxSubscriptions=1000. Evidence count cannot exceed request
count. Every raw record consumes work/text, including unknown counts/SKU. Raw text
accounts selected names, request IDs/regions and each supplied ID/region/response
name/response region/SKU. Output accounting independently includes selected names
and emitted labels plus a128-byte numeric/status allowance per row. Each phase
retains the existing16 MiB minus64 KiB auxiliary budget. Both aggregate raw and
projected text ceilings are reachable under these reservation label caps; exercise
both independently. Reject exceeded bounds without silent truncation. Warnings
use deduplicated fixed codes/messages, without untrusted IDs or payload echoes.

Intended new production path is internal/plugins/region/reservation_runtime.go.
Add independent reservation runtime tests, an isolated compiling-control script
under scripts/tests and explicit steps to both required native jobs. Update this
document, continuity/ledger/roadmap/failure notes and index accepted PR127/display
fallback evidence in the quota contract. No shared refactor, existing reservation
projection change, new collector/public availability, schema, dependency, pin,
capture-byte or Azure-scope change.

## Acceptance, limits and rollback

Require complete literal typed records/health/ten-cell tables for all six rows
above, positive known source status/count cases and separately specified source
corrections. Synthetic IDs assembled for target decoded-evidence admission from
source entry labels are not raw Get response-ID captures or collector execution.
Exercise missing/partial/unknown/unsupported/empty, malformed/foreign/case-colliding
and duplicate identities, response name/region conflicts, name fallback, missing
SKU/counts and contradictions, extreme arithmetic, exact/over request/record/label/
raw/output bounds, skipped-input budgets, cancellation, warning/cell ownership,
sequential/concurrent isolation and quota coexistence. Critical guards require
compiling mutations rejected by named assertions, healthy/restored baselines on
both hosts. Build/panic failures are not detections.

Before protected expected-head acceptance: review exact diff/contract/source
provenance, fresh final-head Linux quality/Windows/source mandatory logs, all17
unchanged retained captures150 chunks/hash/full bytes, all three review surfaces,
live base/preview/protections/identity/tree/parents/remote bytes. Then verify signed
merge and distinct accepted-push native/source proof. Self-review and automated
status markers are not independent-person approval. Reviewed revert is rollback;
preserve prior evidence, pins and gates.

Next: publish and remotely verify the implementation checkpoint, then inspect fresh
exact-head native/source mandatory logs, review surfaces and protected acceptance
with distinct accepted-push proof.
Laptop/Azure remains unconfirmed. Live/DV001 nonproduction hierarchy/load/freshOS/
maintenance/Gate004/release/module-only advisory stay open; no resources, roles,
fixtures, production substitution, exposure or release are authorized here.

## Local implementation feedback, not acceptance

CalculateReservations, seven independent runtime tests and eighteen isolated
compiling controls are implemented. Full region race and all eighteen named
assertion controls/healthy/restored baselines pass locally. Tests include literal
typed records/ten cells/health, five known actual source status/count cases with
explicit synthetic admission IDs, unknown/partial/missing/SKU corrections,
identity/region/presence rejection, cancellation and concurrent ownership, work
limits/order, exact reachable raw/output text ceilings and one-byte overflow.
Workflow YAML,460 authored local links, whole Go suite and vet pass. Fresh exact-head
native/source acceptance remains pending until its recorded completion.

FN086 records the first control run's invalid warning-ownership mutation: removing
the sole slices.Clone call made slices unused, so compilation failed after sixteen
valid detections. No whole-run or ownership detection was claimed. Retain a
discarded clone call to preserve compilation while introducing warning aliasing;
the corrected ownership mutation fails named assertions and all eighteen controls
with healthy/restored baselines pass. Source production/captures/pins, existing
reservation projection and public interfaces remain unchanged.
