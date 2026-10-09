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

REST quota and VM extensions are accepted through PR129/130, with shared JSON
Unicode duplicate admission corrected in PR131. All 25 quota collector compiling
controls passed separately on both final candidate and accepted push hosts.
Earlier pending/exclusion prose below is historical. Collection does not establish
public region execution or full Azure quota accuracy.

Earlier records below are preserved historical snapshots, superseded only in
current-state/resume instructions by the verified checkpoint above.

---

## PR129 accepted closure and VM extension proposal

Four-provider implementation accepted at PR129f0760de7/tree0c71a2ff with complete
pre/post proof in [PR129](https://github.com/DeBoX85/Cloud-Assess/pull/129), indexed in
[REGION_VM_QUOTA_COLLECTOR.md](REGION_VM_QUOTA_COLLECTOR.md). Pending/unaccepted prose
below is historical. [PR130](https://github.com/DeBoX85/Cloud-Assess/pull/130) extends
this same guarded engine to VM2024-11-01, with unknown first404/405 distinct from
other providers. Earlier VM rejection/exclusion describes PR129 scope, superseded
only after PR130 acceptance. Existing budgets/guards/ownership remain required.

# Bounded REST quota collection contract

Status: implemented candidate in [PR129](https://github.com/DeBoX85/Cloud-Assess/pull/129),
UNACCEPTED until fresh exact-head native/source/review and protected acceptance.
Pre-edit42c6bc76be9d60ff1631ef37b6ba5b146b7dd191/treea8953ea2ad0138b585024cf87f9864228fd2bd2e
original tree/all5 full remote documents/Denis/parent/ref verified before code edits.
Accepted baseline PR128:27b45950ef9e1c830e50a121781278e769a488f4,
tree64f6c15a5b8f79eddcc9f7f958618a3f535a0072.
[PR128](https://github.com/DeBoX85/Cloud-Assess/pull/128) is final recovery
authority for reservation calculations and complete pre/post native/source proof.

## Source and slice

AZQR8e4f0577f3615e6c9014c031bcad079f235369cc, APRL60eaddda76541f6adbc1c5ffa686829807e55e29 stay pinned.
Source quota.go blob9d98d3d2c221b925fd5b4bbde4ddf78476146491
and provider blobs are recorded in unchanged testdata/source-quota-inputs.json.
[REGION_QUOTA_RESERVATION.md](REGION_QUOTA_RESERVATION.md) and
[REGION_QUOTA_RUNTIME.md](REGION_QUOTA_RUNTIME.md) establish source captures
and accepted bounded calculations. Inspect exact source paths before editing.
This slice is only the four REST quota providers. VM SDK and reservations,
other region collectors, coordinator/registry/CLI/report integration remain separate.

GET /subscriptions/{UUID}/providers/{provider}/locations/{canonical-region}/usages
with exact source versions:
Network Microsoft.Network 2022-07-01;
SQL Microsoft.Sql 2021-11-01;
App Service Microsoft.Web 2023-01-01;
Storage Microsoft.Storage 2023-01-01.
Provider filtering and arithmetic stay in CalculateQuota, not duplicated in collection.
Keep received page/row order and all decoded raw rows, including filtered counters.

## Interface and request boundary

New RESTQuotaCollector holds an immutable parsed HTTPS ARM origin and trusted
RESTQuotaGetter.GetBoundedWithResponse(context,URL,limit) returning owned bytes,
closed response metadata and error. The existing azure.HTTPClient satisfies it.
Getter must honor context, bound all attempts/error reads, never follow redirects,
and use the selected cloud ARM audience. No credential discovery or client construction
inside collector; actual production audience wiring remains a future integration obligation.
Native fixture must exercise the real authenticated bounded HTTP client, assert GET,
no body, exact endpoint/query/audience/Bearer token and response close ownership.
No HTTP listener, real credential or live service needed.

Collect(ctx, selected subscriptions, one QuotaRequest) returns RESTQuotaCollection
with owned QuotaEvidence and sanitized fixed FailureCode. Validate all selected UUIDs,
case collisions/names, requested selected UUID/canonical region/four provider types
before first call. VM is rejected, not routed to REST. No retries in collector.
Origin rejects userinfo/path/query/fragment/opaque/force-query, invalid/non-ASCII
DNS host, percent/backslash/control/whitespace, ports other than absent/443.
Custom/sovereign origin allowed; do not hardcode public-cloud audience.
Every URL is validated BEFORE getter/authentication. Absolute or relative
continuation resolves to exact same HTTPS authority/path/API version, no userinfo,
fragment/opaque/escaped path or query ambiguity. Query allows exactly one api-version
and optional one nonempty bounded $skiptoken; reject other keys, duplicate keys,
cycles and escaped/alternate paths. Restricting arbitrary source NextLink is a
deliberate correction; undocumented continuation shapes stay a future live evidence gap.

## Decode, work and health

Exact HTTP200 and nonnil response metadata required. First-page404/405 is unsupported;
400/401/403/429/5xx/redirect/other2xx/read/transport error is unknown, not healthy empty.
Later failure is partial, even after a valid empty page; preserve only fully validated
earlier pages. Do not clear prefixes on unsupported later responses.
Context cancellation/deadline at any phase returns nil result/context error,
never completed empty or partial success.
Fixed sanitized codes contain no URL/body/provider error/token/subscription data.

Required nonnull value array; empty array is complete empty.
Non-null object items; names string or object value/localizedValue; missing/null counts
become explicit false presence, not zero. Wrong shape/type/fraction/overflow rejects
the whole page. Nonpositive bounded limits remain decoded for accepted calculation
policy; negative current/count beyond existing1e12 limit rejects.
All names/labels, including filtered/unknown-count rows, obey512byte UTF8/control/
replacement/noncharacter bounds. Duplicate raw names across pages reject.
Reject duplicate/casefolded JSON keys, unsafe UTF8/string escapes, trailing JSON,
depth>32, >65536 value tokens or >128 keys per object. Unknown fields are ignored
only after the same structure/text checks; nextLink has its separate8192byte bound.
No service identity inferred from a returned raw quota name beyond request scope.

Per page1MiB, aggregate8MiB, at most64 requests/pages, raw rows8192 aggregate,
text16MiB-64KiB via existing admission/accounting. Exact supported ceilings and
one-over independent fixtures required, including limits consumed by filtered rows.
At8192 raw rows with two512-byte decoded labels plus at most1000 selected names,
the existing aggregate text ceiling is unreachable inside this collector's tighter
row/label bounds. Retain that defense but do not invent an exact-ceiling fixture.
The accepted calculation has separate reachable raw/output text-boundary tests.
Each page must validate fully before prefix append; rows/bytes/pages never silently truncate.
CalculateQuota integration must retain incomplete/unsupported health and independent
ownership. Sequential/concurrent runs have no mutable shared cache.

## Acceptance and rollback

Independent literal request/full evidence/full calculation fixtures for all4 providers,
source retained positive rows, string/object names and exact empty/presence cases.
Origin/continuation/status/read/body/decode/duplicate/count/label/page/row/byte/cancellation
negative fixtures, page-atomic prefix retention, ownership and concurrency; real HTTP
client fixture for method/audience/closed response. Selected compiling mutations must
fail named assertions with healthy/restored baselines on Linux and Windows.
Whole suite/race/vet/docs/YAML, unchanged17 source captures150 chunks/full hashes/bytes,
exact trees/parents/Denis/ref/review surfaces/no-bypass protected expected-head merge
and distinct accepted-push proof. Local tests do not substitute native QA.
Rollback by reviewed revert to accepted parent; never reset protected history.
No live/release/Gate004/load/freshOS/maintenance/advisory closure is implied.

## Implementation feedback

First focused run failed an authored fixture expectation:8193-byte nextLink is
rejected during whole-page admission, so unknown/no current-page rows is correct,
not partial/current-page prefix. Fixture now distinguishes decode size rejection
from bounded unsafe continuation, preserving page atomicity. Production unchanged
for that correction. Expanded focused race passes all declared boundaries/source
integration/cancellation. First mutation run stopped on ambiguous origin anchor
(constructor and continuation), no detections/PASS credited. Use unique full
constructor expression; continue requiring compiling named assertion failures.
Final native/source/review/protected acceptance remain pending.

## Local candidate checkpoint

Eight independent collector tests cover literal requests/full evidence/full cells,
actual retained source names and absolute pagination, complete empty/presence,
request/continuation/status/JSON/count/identity/work/byte boundaries, page atomicity,
ownership/eight concurrent runs/every observed cancellation check and actual
authenticated sovereign ARM-client method/audience/token/closed body fixtures.
All21 compiling controls with healthy/restored baselines PASS locally; final full
region race/whole Go suite/vet,468 local links and workflow YAML parse PASS at
initial87722a5. Final review found case-variant authority could defer cycle
detection for one extra request. Canonicalize the admitted URL authority before
cycle tracking; new fixture/control increases final selected controls to22.
Fresh final-head proof is required; initial-head CI is superseded evidence.
Corrected full region race/all22 compiling faults/healthy-restored/whole suite/vet
PASS locally. Canonical-authority control reproduces the old candidate's cycle key
and fails the named assertion. Final authored local-file links468 PASS after
correction checkpoint; workflow unchanged and YAML remains parsed.
No local Windows/full hosted-native/live acceptance claim. Source captures/pins/
dependencies/public interfaces unchanged. Current PR129 body is exact evolving
publication/run/recovery authority; earlier planning wording remains history.

## Automated review correction

Inline4213079235/formal5449652646 on initial87722a5 found incomplete Unicode
noncharacter rejection. FN091 records confirmed inherited primary/auxiliary and
collector predicates. A shared region predicate now rejects all66 noncharacters
(FDD0–FDEF and every plane's FFFE/FFFF) while valid supplementary characters remain
supported. Nine independent collector/projection tests and23 compiling controls
are the final scope. New source/native/review acceptance required; earlier heads
and finite focused tests do not certify all project Unicode handling.
Final corrected local full region race/all23 compiling controls/healthy-restored,
whole suite/vet/existing five primary controls PASS;468 links/YAML/staged diff PASS.
This includes empty-port origin admission and complete shared noncharacter checks.
Final-head native/source/review/protected merge/distinct push still required.
