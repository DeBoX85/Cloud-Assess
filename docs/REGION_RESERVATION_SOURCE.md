# Reservation SDK source characterization

Status: IN PROGRESS, source-only and unaccepted. Baseline PR125 is
cd81283a5250f2ffb9d943e4697a7a88aa395560, tree2077e879ed2073686a4334b5b6eb707fe85cd1c0.
Its final exact-head and distinct accepted-push evidence is indexed in
[PR125](https://github.com/DeBoX85/Cloud-Assess/pull/125).

## Source and request premises

AZQR stays8e4f0577f3615e6c9014c031bcad079f235369cc,
tree17d93b20c303f90f7843036be82f0dc32f3260f1; APRL stays
60eaddda76541f6adbc1c5ffa686829807e55e29.
Unchanged internal/scanners/plugins/region/crg/crg.go blob
1ffb6e37cc1c21f6d8b5c9ce483f793387e0e6bc contains FetchReservations.
crg_test.go blob72e398123b1a0acf0ab51efff66f136b4549ca06 repeats status arithmetic
and does not execute that fetcher. models.GetResourceGroupFromResourceID extracts
between fourth/fifth slash without validating semantic subscription/provider identity.

Compute SDK armcompute/v6 v6.4.0 group subscription list, reservation group list
and expanded Get use api-version2024-11-01. Inspect actual generated request code
before encoding expectations. Get uses expand instanceView. Inject synthetic
credential and transport, honoring context and accepting only explicit sequential
GET URLs/no body/canary bearer. Exact default public SDK token scope is
https://management.core.windows.net//.default, distinct from endpoint
https://management.azure.com. Fixtures disable provider registration and retries;
this does not prove default source caller middleware safety or every cloud.
Fail closed for unlisted requests and real HTTP fallback.

## Inspected behavior and actual capture obligations

Group items with missing ID/name are skipped; null group item dereferences nil.
Reservation summaries with missing name are skipped; null summary dereferences nil.
These nil panics are inspected, not executed evidence until capture.
Get response name is used instead of summary name. Lowercase location comes from
group instead of reservation; SKU missing becomes empty. No response identity check
is performed. Reserved defaults0 unless SKU capacity exists. Allocated defaults0
unless expanded utilization exists, otherwise length of allocated VM list, including
null elements. Missing utilization is not evidence of zero allocation.

Status precedence: allocated==0 Idle; available<0 Over-Allocated;
available==0 At-Capacity; otherwise Available, where available=reserved-allocated.
Zero/zero is Idle; negative capacity is not rejected. Group page error returns nil
and error, discarding all prior entries. Reservation page error stops that group
but retains prior entries/continues others. Get error drops that reservation and
continues. The latter two can yield unmarked incomplete success. Empty success is
nil. Cancellation within list/Get must be recorded honestly, including swallowed
inner errors if observed. No explicit page/work/identity/text bound exists here.

Capture actual unchanged fetcher with healthy Idle/Available/At-Capacity/Over-Allocated/
zero/negative counts, missing utilization/capacity/fields, group-versus-Get identity,
missing/skipped IDs/names, malformed group IDs, null group/summary, empty pages,
group and reservation pagination, malformed/denied first/later group/list/Get,
continuation failure and cancellation. Recover only at harness boundary to mark
Panicked=true, never as successful empty collection. Full input/output snapshots
retain every entry field, Failed/Panicked, request sequence and nil-empty state.
Do not substitute copied arithmetic for collector execution.

## Acceptance and target correction boundary

Publish this pre-edit contract and acceptedPR125 continuity first. Add an eighth
isolated source harness/two captures, keeping prior15 bytes unchanged. Emit actual
observations before unconditional retained-byte equality; missing goldens cause
declared evidence-gap failure, never source PASS. Independently reconstruct complete
indexed chunks with SHA256/UTF8 checks, retain actual bytes, then full independently
specified literal outputs/requests and hash/policy premises. Compiling controls
must detect changed status, panic, partial success, identity/location and SDK request
expectations through named assertions, with healthy/restored baselines both hosts.
Formatting/build/panic failures do not count as detected oracle controls.

Final head needs fresh full Linux/windows/source mandatory logs, exact source/APRL/
modules unchanged, full diff/tree/bytes/commit parent/Denis/protections/preview, protected
expected-head merge and distinct accepted-push proof. No old-head acceptance transfer.
Record failures and uncertain outcomes before retry; preserve all historical evidence.

Later target contract must distinguish unknown utilization from idle, reject malformed/
negative/unselected/conflicting identities/counts, bound requests/text/work and preserve
healthy partial results with structured incomplete health. Valid source arithmetic/status
should remain, deliberate corrections documented/tested separately. This slice adds no
target collector, new CLI/schema/dependency/pin or public region enablement. Reviewed
revert is rollback; never reset protected history.

Laptop/Azure expected within24 hours, NOT confirmed restored. No offline input needed.
Approved live queue is separate after confirmed access/scope, including DV001 wholly
nonproduction hierarchy. No Azure resources/roles/fixtures/production substitution/
exposure/release. Live/load/freshOS/maintenance/Gate004/release/module-only advisory open.
Self-review/source-independent literals are not independent-person approval.

## Local source observations

Unchanged FetchReservations35 scenarios PASS locally with strict controlled SDK fixtures.
Null group/summary panics are now executed observations. Missing utilization becomes
Idle, null allocated VM counts as one, malformed group ID returns unmarked nil success,
list/Get failures retain incomplete success and inner cancellation is swallowed.
Later group failure drops prefix. Exact observations are not target corrections.
Inputs3e8e9870bb4a4d4a91490ac400dccfc28e7070ca4b8a2d22434176d1c61b3c09/47327bytes;
outputs4393a4ae45d0b5bacd116aa48bc44ed5698a4cde5b991b2e6e6ab6429974c2d5/31727bytes.
Hosted capture/full retained-byte/literal/native acceptance remains pending.

## Hosted observations and independent source-oracle acceptance

PR126 is the exact evolving final-head/merge evidence index. Source37569270680/
job112623993394 executed8 harnesses/source-module-APRL isolation, emitted17 files/
150 contiguous chunks then failed declared two absent goldens. Retain this evidence-gap
failure, never a source PASS or product/service failure. All SHA256/full UTF8 bytes
verified, prior15 unchanged; new inputs47327bytes/16chunks and outputs31727bytes/11chunks
match local actual source observations and are now retained unchanged.

Complete35-case literal outputs include all entry fields, requests, failed/panicked/
nil-empty state. Independent hash/policy checks preserve source/SDK/fixture premises.
Seven compiling controls modify status/panic/partial success/inner cancellation/identity/
location/SDK request expectations and must fail the named full-record assertion; each
healthy/restored baseline passes. Both native hosts require them. Final local focused
full assertions/seven controls/restored/full region race pass; final source rerun with
only runtime nil-dereference recovery produces unchanged bytes. A network/other panic
is a harness failure, never recovered as source input behavior.

Final exact-head native/source and distinct accepted-push verification remain pending
until recorded in the live PR126 body. Controlled SDK middleware is not a default
source safety/target collector/every-cloud/live claim. Complete source observations
are correction obligations for the later bounded owned target contract.
