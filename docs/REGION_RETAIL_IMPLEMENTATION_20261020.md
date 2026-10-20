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


## Final candidate preparation and reevaluation

WIP181/44114c6357caed81d5153b69be5a231e7eabcbec/tree1f81f07f218520c29f8b7aae9b9b2961f605c592 is preserved with all17complete files/lineage/Denis identity read back. It is a recovery snapshot, not the corrected final candidate. Follow its explicit live successor for final own-head outcomes.

Additional review tightened the remaining-run getter bound before I/O, while retaining max+1 overflow detection/accounting. Exhaustion of the32MiB/4,096request/1,048,576token/65,536binding budgets is exercised globally; raw-row and matching counters have explicitly labelled direct boundary checks. Matching work is normally masked: each admitted raw row costs one lookup and each binding one operation, so65,536rows plus65,536bindings bounds effective matching work to131,072 before the independent1,048,576guard. Do not claim the larger limit was reached end to end. A single unrepresentable triple in URL batching fails explicitly instead of truncating. New cost integration confirms a2to4retail price produces100%weighted difference and a conflicting target prefix contributes no evidence, preserving all unrelated Comparison fields/input bytes.

QA-RETAIL-03: initial complete-only reversal survived an insufficiently specific integration assertion; corrected selection uses unchanged frozen zero-price status/map assertion. A failure-entry context reversal was observably redundant because the final check still returned caller cancellation; remove only that non-proving reversal, retain the failure test and both runtime guards. Required-bool null/number/string tests now explicitly exercise pricing pages, where metadata's separate primary=true rule cannot mask the type check. Original surviving logs/plans retained. No expectation was weakened.

QA-RETAIL-04: pinned Go time/format_rfc3339.go and format.go confirm the fallback accepts single-digit hours, comma fractions and out-of-range timezone fields, and truncates fractions after nanoseconds. New independent literal tests fail WIP code for malformed admission/missing-date ordering before correction. Runtime now bounds RFC3339 grammar, compares parsed UTC instants plus exact residual fractional digits, treats absence as older than ANY dated row (including year0000) and excludes subnanosecond future instants. All52frozen cases remain unchanged. Four compiling date reversals establish these behaviors independently. No standard/pin refresh or source API rewrite.

All42final compiling reversals fail named behavioral assertions without compilation/panic errors, restore exact original source bytes and pass restored baselines locally. Full repository uncached preflight (generated workflow/strict gates/handover/tooling guards, formatting, every Go package test and vet) PASSED serially after mutations stopped. No overlapping local QA claim. The additional CI step uses only explicit per-platform additions: Linux controls-1, Windows controls-2; old580/575identities must be conserved plus42, for622Linux/617Windows. Protected gate logic/names/triggers and all previous commands remain unchanged.

Coherent rollback/restoration is next, then final candidate publication/readback/full own native/source/log/evidence. Task7 remains IN PROGRESS until those complete; task8 contract must follow reevaluation of the final actual API/state. Latest WIP/previous success never certifies this changed candidate. No critical blocker or scope change found; conservative unknown/not-assessed still mitigates missing kind/SQL deployment coordinates for the upcoming inventory contract.


Coherent rollback/restoration PASSED for all25task-owned code/test/workflow/documentation files to exact verified180. The old region package compiles/tests successfully without retail APIs/runners; exact final bytes restored and every retail test passes again. Generated Python bytecode caches were excluded from the final owned-file set. Full source hashes/file dispositions and complete logs are retained. Local source/runner rollback is not installed/live rollback. Final own-host gates/evidence remain PENDING at authoring; live final PR owns actual outcomes.
