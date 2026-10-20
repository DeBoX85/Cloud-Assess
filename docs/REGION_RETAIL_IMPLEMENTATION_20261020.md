# Retail implementation checkpoint, 2026-10-20

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Task7 is in progress under the frozen [retail contract](REGION_RETAIL_CONTRACT_20261020.md). Task8 remains research only and follows task7 reevaluation. No live retail or Azure request was made.

## Verified starting point

Draft PR180 at9fe1261260f55d9796917686a4f0bc10d1cd9490/tree145d410ff9f225ce1f5a275bf5a13bd49cf411fb is COMPLETE VERIFIED OFFLINE for narrow security maintenance and task6 prerequisites. Quality39440501477 and source39440501617 passed.580Linux/575Windows prior controls preserved. The final [security evidence/handover](https://github.com/DeBoX85/Cloud-Assess/blob/cf93d2149e194e4c29622336d3344e41ccd72301/docs/evidence/20261020/PR180_EXCELIZE_SECURITY_EVIDENCE.md) supersedes dated pending/decision snapshots. Residual3imported-not-called and1module-only advisories remain open; no advisory-free claim.

## Implementation and current checks

Isolated retail-implementation-20261020 starts at180. New private two-phase collector retains source retail fields, original/target bindings, raw admitted rows, UoM mismatches and counts. Atomic strict page admission, bounded global budgets, safe continuation validation, parsed dates, deterministic minimum-tier/latest-price conflict exclusion and complete-only owned evidence preserve existing EnrichCost/weighted arithmetic and public schemas. No production caller is wired yet; task7 tests explicitly feed Evidence into existing cost enrichment.

The preserved178first consuming runner is byte-identical. Its actual missing-API failures on BOTH native hosts were inspected before any production edit. All52frozen literal outcomes and all15independently prepared supplemental tests pass locally. The prepared supplemental SHA256f88d36c860233026b461c07a053ebdb337e2f4412ec46648b6e8b2872201a21b is preserved in prior preparation evidence; gofmt only changes its formatting. Additional adversarial/continuation/exact-budget/ownership/conflict tests also pass in one uncached TestRetail run. These preliminary tests do not certify the unfinished candidate.

Research of pinned Go1.26.9 net/http/transport.go, shouldRetryRequest, confirms transparent GET retries on reused connections and separate HTTP/2 retry paths. The independent default transport explicitly enables HTTP/1 only and disables keep-alives, compression and proxies; no cookie jar/credentials/redirect or application retry.30second client/header/TLS timeouts and max+1 bounded reads/body closure. This conservative one-attempt choice costs extra connection overhead; live performance is not established. No network probe or dependency change.

## Remaining before task7 closure

Review edge cases and masked/derived bounds, add compiling reversal/restoration CI controls without losing any580/575previous identities, run full serial preflight and source+consumer rollback/restoration. Publish final candidate and require its own full native/source gates, complete logs/control conservation/source reconstruction and durable evidence. Then reevaluate scope/conflicts/superseded tasks and update handover BEFORE task8 required-kind inventory/evidence contract. Do not label WIP as verified, accepted or MVP complete.

## Development observations

QA-RETAIL-02: first development compile found production retailRun type colliding with the independently prepared test helper. Renamed only production type to retailCollectionRun. Frozen test expectations unchanged; compile failure is not behavioral defect proof. An exec-server disconnect happened before an uncertain shell call; actual status was checked before importing missing files. No uncertain remote mutation.

The raw65,536row cap is normally masked by the lower JSON-token/wire limits for full typed rows. Direct decoder counter-bound tests are separately labelled; do not claim an unreachable full-page integration proof. Request text4MiB is likewise constrained by8,192*512byte identities plus required regions and65,536cells. Keep every explicit bound, report actual reachability honestly.

No Azure/fixtures/RBAC/subagents/protected merge/bypass/release/unrelated refresh. Protected a6fee427 and AZQR/APRL pins remain unchanged. Independent review, live filter acceptance, fresh-OS/operational recovery and broader acceptance obligations stay open. Local-only work must be published/read back before long QA. Source rollback is not installed/live rollback.
