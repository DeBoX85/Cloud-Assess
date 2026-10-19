# Meter-history supplemental acceptance plan

Frozen behavioural contract: REGION_METER_HISTORY_CONTRACT_20261019.md and67literal cases. No production API exists at preparation. Preserve the original oracle/first runner, existing551Linux/546Windows controls, original captures and pins.

Independent supplemental verification before final implementation certification:

- Constructor HTTPS root admission, nil receiver/dependency, selected scope normalisation/collision/bounds and preauthentication rejection.
- Query complete structure and calendar precision; invalid months/clocks/currency before HTTP; caller cancellation before/after query.
- Per-page1MiB exact/+1, total8MiB exact/+1,64pages exact/+1, cumulative8192rows exact/+1,32columns exact/+1, raw65536values exact/+1 including unknown metadata across pages,32depth exact/+1,128keys exact/+1 and8192generic string bytes exact/+1. Isolate each reachable bound from all other limits.
- Derived limit clarification: distinct meters cannot independently exceed8192while raw rows remain<=8192, since each row contributes at most one new meter. Test8192distinct meters and the row+1boundary; do not invent an independently reachable distinct-meter overflow or add a redundant mutation control for it.
- Raw-work proof for an otherwise valid empty response: root1+properties1+columns-array1+rows-array1+extra-array1+65531nulls=65536values. One additional null must fail, below1MiBwire and depth limits. Object keys are not values; all are independently validated and subject to the128per-object-key limit. Values in unknown metadata count.
- Depth proof: root is depth0; an extra field holding31nested arrays followed by null reaches null depth32and is allowed.32arrays reach null depth33and must fail. Keep the actual envelope/rows valid.
- Page atomicity, duplicates across pages, cost aliases, per-meter sums and exact/lower/upper numeric bounds; currency/scope/type validation even for omitted nonpositive rows; received-order sum and sorted IDs.
- Continuation origin/effective443port/path/version/query, malformed percent encodings, encoded paths, userinfo, controls, lengths and canonical loops before POST. Same query bytes across every page.
- Actual bounded HTTPClient POST request/method/content-type/ARM audience, failed and successful response closure, no redirects, oversized first/retry/challenge bodies and operation cancellation. Fake posters cannot establish these transport facts.
- Caller/wrapped cancellation and manually derived fixed late-context schedules, frozen before any reversals. No time-sensitive flaky scheduling or recalibration to weakened code.
- Own maps/slices, input snapshots, repeated/concurrent use and cost evidence copying. Complete/noncomplete collection-to-EnrichCost integration must respectively permit/deny cost computation without weakening existing arithmetic.
- Add named compiling reversal guards for independent admission/health/ownership behaviours. Each must fail its intended assertion, restore exact bytes and pass baseline; register additive generated CI steps. A compiler failure never counts as reversal proof.
- Exact final-head8workers/bothgates/source, all previous identities, full source reconstruction and every mandatory log/step/package/security result. Atomic implementation+runner rollback to verified contract tree, retest, restore candidate bytes/tree and retest. Source rollback does not prove installed/live rollback.

Preparation refinement, not a critical blocker:8192distinct-meter admission is mathematically bounded by8192raw rows. This clarification supersedes only an interpretation that every declared bound must have an independently overflowing valid fixture. It changes no frozen case, constant or production behaviour. Record this in the next between-task handover.

Budget proof 675bytes SHA256b690edc8d4951d8ad80afc12edfb8483b2f09e92001711165c2133a6e78fce9b. Preserved before production/reversals.


Preimplementation refinements: known decoder-equivalent field aliases must be admitted consistently, never ignored to hide nextLink. Pinned Go1.26.9 standard-library documentation and [official encoding/json documentation](https://pkg.go.dev/encoding/json#Unmarshal) support this matching behaviour; no decoder/toolchain migration. Pinned local ParseFloat probe confirms positive1e-999silently underflows to0; a new independent raw-positive test requires rejection, while raw0e-999and negative omission remain valid. No original literal changed. Existing PostBounded does not expose a typed byte-limit cause; history_limit classifies returned-body/raw-work counters, while a transport-aborted oversized attempt can be history_request_failed with explicit incomplete health. Never parse arbitrary provider error text to classify it.
