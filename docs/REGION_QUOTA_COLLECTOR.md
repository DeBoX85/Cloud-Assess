# Bounded REST quota collection contract

Status: pre-edit planning, UNACCEPTED until this proposal's exact-head gates.
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
