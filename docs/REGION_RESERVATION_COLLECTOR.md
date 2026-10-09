# Current handover authority, 2026-10-09

Verified implementation: [PR131](https://github.com/DeBoX85/Cloud-Assess/pull/131),
accepted `b653a3abfc35590a185531095a168a9a107e769e`, tree
`47515c474ddf400da0a68c556e9147d8d3c2576b`. Live refs can advance.
Read [HANDOVER_CHECKPOINT_20261009.md](HANDOVER_CHECKPOINT_20261009.md) for exact
parents/identity/pins, separate candidate/accepted-push proof, audit/findings
disposition, remaining tasks, access limits, workspace recovery and QA provenance.
This documentation checkpoint's final acceptance evidence is in its live PR body
on `docs/handover-20261009`; do not assume it is already accepted.

This collector is VERIFIED OFFLINE through merged PR131, not still pre-edit.
Exact candidate and accepted-push runs, checkout revisions and job IDs are indexed
in the linked checkpoint and PR131 body. Final accepted corrections include
structural returned IDs, Unicode duplicate identities/JSON aliases and physical
ASCII region admission. Public region execution and production audience wiring
remain separate integration work. Historical contract/status paragraphs follow.

Earlier records below are preserved historical snapshots, superseded only in
current-state/resume instructions by the verified checkpoint above.

---

# Bounded capacity-reservation collection contract

Status: IN PROGRESS, pre-edit. Accepted PR130 baselinea79464db59fa0b3e72115d2f07a18f8621bd437f,
tree55802384dfd3df4f7f852c030f3aa75d96ec3acd. Current proposal PR body will be exact
planning/candidate/native/source/review/merge/push authority. See
[REGION_RESERVATION_SOURCE.md](REGION_RESERVATION_SOURCE.md) and
[REGION_RESERVATION_RUNTIME.md](REGION_RESERVATION_RUNTIME.md) for accepted source
characterization and pure calculations. Those historical pending descriptions are
superseded by PR126/128 accepted closure, not new acceptance evidence for this slice.

## Source and actual execution boundary

AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1,
APRL60eaddda76541f6adbc1c5ffa686829807e55e29 unchanged. Actual source
internal/scanners/plugins/region/crg/crg.go blob1ffb6e37cc1c21f6d8b5c9ce483f793387e0e6bc,
FetchReservations, armcompute/v6 v6.4.0 and35 retained actual SDK scenarios establish:
1. subscription provider capacityReservationGroups list GET;
2. validated group resourceID/capacityReservations list GET;
3. validated reservation resourceID GET with $expand=instanceView.
All use source api-version2024-11-01; query encoding matches retained requests.
No version/dependency/source-capture refresh. Current Microsoft Learn Get documentation
corroborates resource hierarchy/expanded utilization, but its2025-04-01 examples are
not substituted for the pinned2024 request or treated as pinned observations.

New ReservationCollector reuses the existing bounded response getter and constructor
origin policy. Collect(ctx,selected subscriptions,one ReservationRequest) emits owned
ReservationCollection(Evidence,FailureCode). No clients/credentials discovered inside,
no SDK/provider registration/new retry pipeline/cache/dependency/listener. Getter must
honor context, bound all attempt/error reads, close responses, reject redirects and
use selected cloud ARM audience. Actual production audience wiring remains integration.
Single requested canonical physical region: list selected subscription groups, fully
validate every raw group, then collect only matching group regions. Nonmatching groups
consume raw budgets; no data discarded silently as malformed filtered input.

## Identity and request confinement

Selected UUID/name/case-collision/region admission precedes first call. Required group
ID has exact eight-segment ARM group path, selected subscription and safe nonempty RG/
group labels. Optional returned group name must agree; required location canonicalizes
ASCII case and must be a physical region. Duplicate casefold group IDs reject entire
current page, including outside requested region. Missing/null group metadata/null
items are incomplete evidence, never source panic or quietly healthy skipped data.

Reservation summaries require bounded nonempty name. Derive exact reservation URI
only from validated parent group and name; optional summary ID must match. Duplicate
casefold reservation IDs across pages reject whole current page before downstream Gets.
Get optional ID/name must match requested reservation; optional location must match
parent/request region. Missing response ID falls back to validated request identity:
this is a declared request binding, not proof that a returned identity was verified.
Missing returned name uses existing calculation fallback. Group and Get location
conflicts are explicitly rejected. Existing2048-byte ID and512-byte label structural
policies apply; escaping is performed by URL construction, never string-concatenated
untrusted query/destination. Source permissive slash extraction is not reused.

All continuation links are validated BEFORE authentication: same HTTPS configured
origin/exact chain path, api-version2024-11-01 once, optional bounded single skiptoken,
no foreign path/scope/provider/version/query/userinfo/fragment/escaped alternate path.
Canonical query/authority keys detect cycles. Get has only fixed version/instanceView
query and no pagination. Preserve exact GET/no body. Custom/sovereign origin allowed.

## Decode, presence, work and health

Strict bounded JSON reuses quota token walker, with a separate explicit id-string cap
2048 for reservation identities; quota strings remain512. Unknown fields remain subject
to structure/text/noncharacter/UTF8/duplicate-casefold-key/trailing/depth32/tokens65536/
128-object-key checks. Required list value array, no null item; current page fully
validated before downstream requests. Get must be a nonnull object. Counts are int64,
nonnegative and at most1e12. Missing/null SKU capacity remains ReservedKnown=false.
Missing/null expanded utilization or allocated array remains AllocatedKnown=false,
never inferred Idle. Actual empty allocated array is known zero. Every allocated VM
item needs nonnull exact ARM VM or VMSS-instance ID in selected subscription; reject
null/missing/non-VM/malformed/duplicate references, do not count ambiguous objects.
This validates reference structure, not actual VM existence/ownership/location.

Budgets:1MiB per response,8MiB aggregate successful bodies,256 getter calls total,
64 pages per list chain;8192 aggregate raw groups/summaries/allocated references.
All filtered/unknown rows consume work. No silent truncation. Existing text defense
remains, but no impossible exact-ceiling fixture inside tighter byte/work caps.
No collector retries. Global request/byte/work exhaustion stops further calls with
honest prefix health. Per-group list/Get operational failures preserve valid earlier
results, continue independent peers where safe and mark partial even with zero rows.
First group-list failure unknown;404/405 are unknown failures, not unsupported success.
Group-list later failure preserves prior complete work. Invalid whole current page
contributes no downstream requests/current-page records. Fixed first failure code
contains no raw body/URL/token/subscription/provider exception.
Cancellation/deadline anywhere returns nil result/context error, never partial success.
Run state/results owned; sequential/concurrent calls cannot leak mutable state.

## Independent acceptance and correction provenance

Literal three-level exact request/evidence/calculated cells/health fixtures, all four
source-valid statuses and known zero versus absent counts. Actual retained source
empty/denied/malformed/null/conflicting identity/missing fields/pagination/cancellation
requests are provenance. Most source arithmetic fixtures intentionally have EASTUS
group and WESTUS Get plus synthetic-vm IDs. They are not healthy target identity proofs.
Use separately labelled corrected synthetic counterparts for healthy identity/count
oracles, retaining original source files/hashes unchanged. Do not call adapted data an
unchanged capture or a live fixture. Source all35-case capture still independently gated.

Tripwires constrain scope/parent/returned identity/region and all continuations before
credentials, plus denial/read/redirect/status/error/prefix/cancellation/malformed atomicity,
raw budgets/exact permitted ceilings/one-over, ownership/concurrency. Actual sovereign
HTTP-client audience/token/GET/no body/closed responses needed. Selected compiling
controls for request version/expand, identity/region, presence/allocated references,
pagination/prefix/size/work/cancellation must fail named independent assertions with
healthy/restored baselines on BOTH native hosts. Shared quota25 controls preserved.

Affected code: new reservation_collector.go/tests/control runner and native registration;
quota strict-walker wrapper for ID mode without changing quota policy. Continuity docs
and this contract updated. No pure arithmetic/registry/public CLI/schema/report/source/
modules/pins/other adapters. Full region race/suite/vet/docs/staged original tree/full
remote file/Denis/parents/ref/provenance, exact final native/source/review/protected
expected-head merge and distinct accepted-push proof required. Rollback reviewed revert
on current accepted branch, never force reset. No local-only backup or human-review claim.

## PR130 accepted closure and recovery

PR130 final authority: [PR130](https://github.com/DeBoX85/Cloud-Assess/pull/130).
Candidatec0a11376/tree55802384/parent97609c6a actually tested preview1cfa564dcf6b1e963f6bdb70d4d6339a0be34962
with parentsf0760de7/c0a11376. Native37723750988 jobs113137349623/113137349427 and
source37723750987/job113137092286 SUCCESS. Distinct actual accepted-pusha79464db:
native37725227530 jobs113141759835/113141760070, source37725227544/job113141760080 SUCCESS.
Full52/37 steps/all25 collector controls both hosts/eight source harnesses17 records150
chunks/full hashes/bytes verified. Automated6051685481 completedc0a at2026-10-08T03:42:37.240613Z,
formal/inline empty/no findings, not human/independent-person approval. Signed merge
and706 original blobs/raw commit/fsck verified; source/candidate actually materialized
APRL separate from accepted snapshot gitlink metadata. PR130 need not be replayed.

Laptop/Azure NOT confirmed; live/DV001/load/freshOS/maintenance/Gate004/release/module-only
advisory OPEN. No resources/roles/new fixtures/production substitution/exposure/release.
No blocking offline input. Next remaining adapters/coordinator/public/all-format integration.
Publish coherent remotely verified checkpoints and inspect uncertain outcomes before retries.

## Implementation checkpoint

PR131 implements the internal bounded contract above. Eight focused test functions
and17 compiling healthy/fault/restored controls are present; corrected local named
checks pass. FN095 records the initial work fixture and uncompiled fault failure.
Independent peer recovery asserts that list/Get/later-list failure preserves prior
rows and continues a separate valid group. Mandatory Windows/Linux/source/review/
protected acceptance and distinct push remain pending; exact evolving revisions,
run IDs and completed checks are in the PR131 body. Historical source captures,
module pins and public execution boundaries are unchanged.

Self-review correction:18th compiling control independently disables response
status checking. Non200 responses now carry healthy bodies and nil getter errors;
200/getter-error checked separately. Earlier17-control PASS is historical. Current
corrected head/native/source/review evidence must be obtained, see PR131/FN095.

FN096 correction validates optional returned summary/Get IDs structurally before
case-insensitive equality. Unicode fixed-segment lookalike regressions and19th
compiling structural-bypass control required. This confirmed proposal-code defect
was found before acceptance; prior-head evidence remains historical only.

FN097 automated review and separate collector/runtime regressions confirm lowercase
duplicate keys disagree with EqualFold for Unicode labels. Material correction also
changes accepted internal runtime duplicate admission to canonical SimpleFold keys;
arithmetic and displayed lowercase ResourceID unchanged. This explicitly supersedes
earlier unchanged-runtime scope.21 controls and10 focused test functions are required
on the corrected head; earlier native/source/review evidence does not transfer.

FN098 extends correction to the shared JSON guard: regionFoldKey canonicalizes
Unicode aliases consistently with Go struct-field decoding before duplicate
rejection. Literal collector and both shared modes reject duplicate sku/properties
aliases.22 collector faults and retained25 quota controls required; only quota
control anchor follows corrected key expression.11 focused functions now cover
these findings. Mandatory final-head native/source/review proof remains required.

FN099 enforces the declared physical/ASCII region contract via unchanged pinned
physical-location exclusions, canonical format and ASCII-case byte length. Caller,
raw group and optional Get location share predicate; no live catalogue or changed
pure runtime format claim.24 collector faults now required on corrected head.
