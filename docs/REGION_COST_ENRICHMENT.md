# Current handover authority, 2026-10-09

Read [SESSION_HANDOVER.md](SESSION_HANDOVER.md) and
[HANDOVER_CHECKPOINT_20261009.md](HANDOVER_CHECKPOINT_20261009.md) before following
an implementation status, baseline, executor description or next action below.
Accepted implementation extends through [PR131](https://github.com/DeBoX85/Cloud-Assess/pull/131),
`b653a3abfc35590a185531095a168a9a107e769e`. The handover PR on
`docs/handover-20261009` holds its own exact publication/CI/merge/push state.
Verify later live refs and PR evidence; this header does not assume that proposal
is accepted. New-chat instructions: [NEW_CHAT_PROMPT_20261009.md](NEW_CHAT_PROMPT_20261009.md).

PR114/119 audit and PR120-131 bounded region milestones are accepted offline;
historical IN PROGRESS/UNACCEPTED/pending baselines or next tasks below are snapshots,
not current instructions to replay completed work. The checkpoint identifies
each accepted milestone and its PR evidence. Preserve the underlying source
contracts, explicit corrections, bounds, mandatory tests and historical results.
Public region execution, full feature parity, live access, Gate004 and release
readiness remain unestablished. No new whole-project or independent-person review
is claimed by this documentation task.

---

# Cost-enrichment source characterization and runtime contract

Status: IN PROGRESS, source-only characterization. Accepted baseline fbbb819657e3fb555cef83b39a52284b59cb482f/treeaa1f68d3b72e5e0eb2d1a7f5532e602d5cefdf91 (PR121). Live accepted ref, merged latency acceptance and no open proposals verified before work. Existing audit/availability/latency and CostComparison sheet acceptance are not repeated. Laptop/Azure restoration is expected soon, not confirmed.

## Source recorded before target production edits

TARGET_SPECIFICATION requires the pinned region assessment, including historical-spend-weighted cross-region cost difference. The existing CostComparison projection is a sheet, not primary-comparison enrichment. AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1: internal/scanners/plugins/region/cost/cost.go ApplyCostDiffs blob82b357d2c00dd65c08bc5908000311664431c0e6, tests9df6f7622f42a2f26a2911f36bddfb353961aea9, types/types.go blobf9b61fb74075e5a2297be1b898067b8cb047f9b3. APRL60eaddda76541f6adbc1c5ffa686829807e55e29 and all other pins/dependencies/schemas stay unchanged.

Source uses supplied subscription meter history as an ownership index, last duplicate MeterID wins. It normalizes region lookup by lowercasing/removing ASCII spaces, skips nonphysical logical identifiers, and iterates shared region pricing only for indexed meters. HistoricalCost zero becomes weight1. Missing either price skips that meter. If both prices are strictly below0.0001, difference0 contributes its weight to the denominator; a source price strictly below0.0001 with a real target is skipped. Otherwise difference=(target-source)/source*100, weighted by historical spend. A positive total weight sets both primary cost fields. Equal regions/prices and free target remain valid. There is no subscription identity check, context, numeric validation or budget. Negative prices can be treated as free, negative history can make the total nonpositive, and nil/shared-absent/no-pair/logical cases retain earlier AvgCostDifference/HasCostData. It mutates input and map iteration determines float summation order. These are source facts, not target requirements to inherit defects.

Capture unchanged ApplyCostDiffs only, with HTTP tripwires, actual cost/types file Git-blob hashes, complete untouched inputs and complete outputs.26 cases/26 comparisons cover ownership, weighted arithmetic, zero history, both/source near-zero, exact threshold, free/equal/same regions, missing source/target/empty/nil, duplicate received-order weights, negative source values, display/logical names, stale-field retention, empty results and unbound subscription. Source-invalid branches stay retained. Existing nine files remain byte-identical. The runner emits all11 captures then unconditionally verifies every retained byte. Initial two missing new goldens intentionally fail after observation: declared evidence gap, not source PASS. No optional missing-file path, pin/source-production/module/APRL changes or target-generated goldens.

## Independent expected source results

Complete comparison fields/arrays/zone map/order must be compared against independent literals, not only the percentage. Owned-only100; weighted62.5; zero-history unit weight-12.5; both-near-zero weighted25; source-near-zero skipped100; exact0.0001 boundary100; free target-100; both-free/equal/same-region0. Duplicate9-then1 gives50;1-then9 gives90. Negative history and missing/nil/logical/unowned cases retain the initial-5/true; negative prices return0/true. Nil without earlier cost remains0/false. Display labels stay unchanged with100. Unbound second subscription is also enriched100, proving that supplied meter membership alone does not establish selected subscription ownership. Nonfinite source outputs cannot be retained as JSON; target rejection must be tested directly later, not claimed observed by this capture.

## Target-runtime contract for the following slice

Before: no primary cost constructor. After the separate runtime slice: owned complete comparisons plus cost-only stage health. Admit selected/canonical complete comparisons through existing Project before copying. Evidence must explicitly bind each historical meter list to its selected subscription, with collection-complete declarations; shared prices require explicit complete/partial/unavailable status. The future collector must establish actual selected resource/subscription ownership and pricing/UoM completeness. Declared evidence is not live collection proof.

Preserve source owned-meter weighting, zero-history unit weight, exact near-zero threshold/denominator, valid free/same-region behavior and all unrelated fields/order. Deliberately reject duplicate per-subscription meter IDs and negative/nonfinite history/prices; reject display/case comparison identities rather than rewriting them. Reset both cost fields before each pair so unavailable/no-eligible evidence cannot reuse previous results. Canonical deterministic meter ordering avoids map-order float variation. Deep-clone all comparison slices/maps and health; no mutable global cache or exposed evidence. Cancellation before/during/after work returns no partial result. Cost health remains separate from availability/latency and must be propagated by the future coordinator. Missing/partial/unsupported/ineligible pricing and incomplete history produce explicit bounded warnings, never false completeness.

Use existing MaxComparisons/MaxSubscriptions/MaxCostMeters/MaxCostRegions/MaxCostEntries/text/cell limits and an explicit comparison-times-owned-meter work budget before nested arithmetic. Validate finite intermediate sums/ratios and projected output, with independent boundary fixtures and compiling negative controls that cannot be masked by unrelated guards. Record exact numeric limits in the runtime contract before editing production code. This source-only slice does not settle collector endpoint/continuation/currency/UoM/clock behavior. Those need pinned-source and official API contracts, separate request/canary/failure-page tests and honest evidence health before collection acceptance.

## Acceptance, recovery and next action

This slice changes source harness/runner, two actual observed JSONs, independent whole-comparison source-oracle tests and compiling fixture controls on both mandatory hosts, plus current documentation. It adds no target arithmetic, credentials, clients, Azure I/O, registry/CLI/public execution or dependency/pin/schema refresh. Existing sheet/latency controls remain mandatory. Require fresh strict11-file source proof, full exact-head Linux/Windows jobs/steps/logs, complete indexed captures/hashes/UTF8 bytes, source isolation, full diff/remote bytes/commit/tree/ordered parents/Denis identity/base/preview/no-bypass protections, expected-head protected merge and distinct accepted-push proof. Self-review is not independent-person approval; synthetic characterization is not live Azure parity/load/fresh-OS/release approval.

Next: recover actual emitted observations, independently verify and retain them, implement complete independent source oracles/controls, and complete protected source-only acceptance. Then implement the bounded runtime under the contract above; quota/reservation/adapters/coordinator/public/all-format integration remain after it. Existing installed core/plugin live queue can resume independently when laptop/Azure access and approved scope are confirmed. DV001 still requires a wholly nonproduction nested hierarchy; Advisory includes production and is not a substitute. No Azure fixtures/resources/roles/production substitution/exposure/release authorized. Existing module-only advisory and Gate004/live/load/maintenance/release obligations remain open.

Rollback is a protected reviewed revert after dependency review. Read AGENTS/current handover/live proposals before any interrupted mutation retry. Publish coherent WIP and verify bytes/objects/ref; scratch-only work is not backed up. Preserve historical evidence and old workspaces. No local Go/pwsh/materialized APRL; native Actions executes proof.

## Retained-source checkpoint authority

# PR122 cost-source characterization and recovery authority

Source-only cost characterization is implemented in PR122 on feat/region-cost-enrichment-source. Its live PR body is the exact final-head/merge/accepted-push index: https://github.com/DeBoX85/Cloud-Assess/pull/122. Determine unaccepted candidate versus VERIFIED OFFLINE source-only acceptance there and from live refs. Entry accepted baseline remains PR121 fbbb819657e3fb555cef83b39a52284b59cb482f/treeaa1f68d3b72e5e0eb2d1a7f5532e602d5cefdf91; accepted audit/availability/latency/sheet review is not repeated.

Observationaebcde3ee1bbf00418bcde2801ca9c28f2522d37/tree8cd6ff83691f15d0db4ba80afc9002334cf0f57e/parentfbbb8196 passed all5 pure pinned source harnesses and source/module/APRL isolation in source37234582329/job111531040978, then failed the declared absent-golden guard after emitting all11 observations. This is a planned evidence gap, not source PASS/product defect. All109 indexed chunks reconstructed and independently SHA256/UTF8-byte checked; existing9 unchanged. Actual cost inputs112a462e1bb2090afb39322684264189a8d1977b7c7cab938c8b217535ed670d (46496bytes/16chunks), outputs0ae38cf25c54462e3b7314c4f74c31a5bbe0cfac6d22af72554198ffeb447526 (24534/9) retained. Runtime-sourcecost82b357d2/typesf9b61fb7 Git blobs are independently hash-checked inside source execution.26 cases/26 entire comparisons and actual stale/duplicate/negative/unbound behavior have independent literals. Six selected compiling source-oracle controls/restored baselines are mandatory on both native hosts. Strict unconditional11-file equality remains unchanged; no optional bypass/production/pin/dependency/schema/CLI/credential/Azure changes.

REGION_COST_ENRICHMENT.md defines source facts, independent expected results and target correction obligations. Missing pricing retains earlier cost state; duplicate history is last-wins; negative prices can be treated as free; per-call supplied meters do not independently bind comparison subscription identity. These are pure-helper contract limitations, not evidence of historical live misuse. Target runtime remains NOT IMPLEMENTED; source invalid branches stay retained. Next verify full exact-head Linux/Windows/source logs, all109 chunks/hash/bytes, scope/identity/tree/pins/base/preview/protections, protected source-only acceptance and distinct push proof. Then record exact bounded numeric/evidence interfaces and implement owned selected/context-aware honest-health cost enrichment, followed by quota/reservation/adapters/coordinator/public/all-format integration. PR122 body is current publication/run index; native observation status may be superseded by unchanged concurrency policy and must not be transferred to this head.

Laptop/Azure restoration expected soon but not confirmed. Existing core/plugin live queue independently resumes once access/approved scope are confirmed; DV001 wholly nonproduction hierarchy/live/load/freshOS/hostedmaintenance/Gate004/release/module-only advisory remain open. No Azure resources/roles/fixtures/production substitution/exposure/release authorized. No local Go/pwsh/materialized APRL; Actions executes proof. Self-review, no independent-person approval. Old workspaces/history preserved. Publish/read back coherent checkpoints; scratch-only work is not backed up. Inspect live outcomes before retrying uncertain operations.

Earlier pending entries below are retained history superseded by this authority and live PR evidence.
