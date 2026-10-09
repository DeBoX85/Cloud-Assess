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

# Owned weighted cost runtime contract

Status: IN PROGRESS, implemented UNACCEPTED runtime; native tests pending. Planning4597c3afdca634524c0a4441191202d544f2f026/tree8233f7bf4db240ef4592bcf627ae9578929bcb2d/parent201f2692 was remotely published and full four-file bytes/Denis identity/local full tree verified before production edits. Accepted baseline PR122 201f26921c67741a1433697d9489e6458f32bd64, tree2751881f56fe5f326826e92aeb5ca615b6f898c8, ordered fbbb8196/dbed672e parents. Its candidate native37234900124/source37234900149 and distinct accepted-push native37235617666 (Linux111533995467/Windows111533995215), source37235617640/job111533995150 passed complete mandatory logs. PR122 is the retained source acceptance index; all11 files/109 chunks remain unchanged. No open proposals at entry. This contract is recorded before production edits.

## Requirement and independent source contract

TARGET_SPECIFICATION region assessment requires historical-spend-weighted primary cost enrichment. REGION_COST_ENRICHMENT records unchanged AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1 cost.go ApplyCostDiffs blob82b357d2c00dd65c08bc5908000311664431c0e6 and types blobf9b61fb74075e5a2297be1b898067b8cb047f9b3. APRL60eaddda76541f6adbc1c5ffa686829807e55e29 remains unchanged. CostComparison sheet is already implemented and is not this primary arithmetic. Before: no target enrichment. After: pure EnrichCost returns owned complete comparisons and separate cost-only health. No clients, credentials, Azure requests, collector, public registry/CLI, dependencies, pins or schema changes.

API: HistoricalCostMeter has MeterID and HistoricalCost. CostHistoryEvidence has Complete and Meters. CostEnrichmentEvidence has History map[canonical selected UUID]CostHistoryEvidence, PricingStatus (complete/partial/unavailable/unsupported), and RegionPricing map[meterID]map[canonical region]float64. Nil evidence is unavailable. Missing history is unknown, explicit incomplete history is not used for arithmetic. Prices marked partial can contribute eligible pairs but always warn. Unavailable/unsupported must contain no prices. Complete describes declared evidence only, not observed Azure ownership, collector currency/UoM/tier/time/page correctness or live completeness. The later collector must establish these separately before integration acceptance.

Admit complete comparisons through Project first. Additionally reject noncanonical comparison UUIDs rather than silently rewriting. Evidence UUIDs must match selected normalized scope; no cross-subscription aggregate weights. Duplicate meter IDs per subscription reject; the same meter in two subscriptions is valid and retains distinct weights. All numeric history/prices must be finite and in [0,1e12], deliberately rejecting source negative/nonfinite behavior. Meter IDs are nonempty safe UTF8 labels <=512 bytes. Pricing regions are canonical <=64 bytes. A canonical logical region is admitted but yields cost_region_unsupported rather than fabricated data. The source logical exclusion list is independently recorded from pinned types.IsPhysicalRegion.

For each comparison reset AvgCostDifference=0 and HasCostData=false. Preserve all other fields and order; deep-clone all four detail slices and zone map, with per-call owned health. For owned meters in sorted MeterID order: history zero becomes weight1; missing either price warns cost_price_missing and skips; both prices strictly<0.0001 contribute zero numerator and weight denominator; source strictly<0.0001 with real target warns cost_meter_ineligible and skips; otherwise ((target-source)/source)*100 times history. Positive total weight yields HasCostData. No eligible pair warns cost_no_eligible_meters. Same/free/equal prices remain valid. No source stale-state, duplicate-last-wins or unbound subscription behavior is inherited.

Cost warnings are deduplicated bounded static codes, sorted canonically, with StageCompletedWithWarnings; records=len(comparisons). Missing history warns cost_history_unavailable; incomplete history warns cost_history_partial; partial/unavailable/unsupported pricing warns cost_pricing_partial/cost_pricing_unavailable/cost_pricing_unsupported. Health describes cost only and must not replace availability/latency/coordinator completeness. Empty comparison output retains empty owned slice and completed health after evidence validation. All failures/cancellation return nil result, never partial output. Inputs must remain stable during calls; arbitrary concurrent mutation by callers is not supported.

## Exact bounds before production edits

Existing MaxComparisons8192/MaxSubscriptions1000/MaxCostMeters8192/MaxCostRegions32/MaxCostEntries65536 apply. MaxCostMeters bounds the sum of all historical entries, including repeated IDs across subscriptions, and independently pricing meter map size. MaxCostRegions bounds unique pricing regions. MaxCostEntries counts history declarations+historical entries+pricing outer entries+price cells. All decoded scope-name/meter/region text is charged to existing auxTextBudget (MaxPluginTextBytes minus64KiB), labels to512 bytes and regions64. Project enforces full detail/text/UTF16/output bounds before and after enrichment. MaxCostWork=1048576 counts comparison-times-owned-history-meters across all comparisons, charged before cloning/calculation even for incomplete or missing pricing. Checked subtraction prevents integer overflow. Sorting and validation have at most bounded admitted entries; no silent truncation.

Individual monetary values<=1e12 and source denominator>=0.0001 bound per-meter diff magnitude<=1e18, weighted term<=1e30, sum<=8192e30 and total weight<=8192e12. Check finite intermediate sum/ratio and final Project despite implied finite upper bounds. No claim of independently triggerable intermediate infinity under admitted bounds. No arbitrary upper-percent clamp or invented conversion. The later recorded proven lower-domain roundoff correction below is a numerical identity, not accepted negative pricing. Limits are engineering admission limits, not Azure service quotas or currency conversion.

## Independent verification and recovery

Reuse actual retained source inputs/outputs and independently asserted PR122 whole-comparison literals for every valid source arithmetic case; explicitly test corrected invalid/stale/unbound cases. Full literal warning/health and preservation assertions, selected-subscription distinct weighting, deterministic received-order reversal, zero/free/exact threshold/missing/ineligible cases. Independently construct exact/one-over historical, nested-entry, region, text and work boundaries. Verify cancellation before/during/final return, caller/result/evidence ownership, repeated and concurrent calls under race tests. Availability -> pinned latency -> cost -> Project integration compares a complete independently calculated row and preserves each earlier health. Numeric rejection includes NaN/infinities, negative and overlimit prices/history; no false JSON source observation claim.

Add compiling selected fault controls with named assertions/restored baselines to both mandatory hosts, covering scope, evidence binding, duplicates, negative inputs, work limit, weighted division, zero-history, free denominator, threshold, stale reset, incomplete health, ownership and terminal cancellation. Full Linux/Windows existing gates and strict fresh11-file source proof remain mandatory. Self-review is not independent-person review. Protected expected-head merge and separate accepted-push proof follow full diff/remote bytes/tree/parents/Denis identity/source pins/preview/no-bypass verification.

Current proposal branch feat/region-cost-enrichment-runtime is UNACCEPTED. Publish/read back coherent contract and code checkpoints; PR body is exact run/head/merge recovery index. Next implement runtime/tests/controls, classify failures before retries, inspect full native/source proof, then protected acceptance. Rollback is reviewed protected revert after dependent-code inspection. Preserve all historical records/workspaces. Local Go/pwsh/materialized APRL absent; Actions is native execution authority. Scratch-only files are not backed up. Laptop/Azure access expected soon, not confirmed; live scope/DV001/load/freshOS/maintenance/Gate004/release/module-only advisory remain open. No Azure resources/roles/fixtures/production substitution/exposure/release authorized.

## Runtime checkpoint

EnrichCost and tests are implemented on PR123 https://github.com/DeBoX85/Cloud-Assess/pull/123. Twenty complete captured valid/corrected cases, distinct selected-subscription weighting, exact/one-over meter/entry/region/label/work boundaries, numeric/status/identity/canonical/duplicate rejection, adversarial floating sum ordering, all four slices/map/health ownership, concurrent/repeated/cancelled calls and availability -> latency -> cost -> full literal Project row are covered by authored tests, not yet executed. Fourteen compiling runtime fault controls/restored baselines are required on both hosts. Original six cost source-oracle controls and all prior gates stay mandatory. Aggregate decoded text is implied by outer/cell/512-byte/64-byte bounds within auxTextBudget; no independently triggerable aggregate overflow is claimed. Intermediate infinities are likewise implied impossible by admitted numeric/work bounds, while defensive checks remain. Existing Project enforces complete output text/UTF16. Local Python compile/static anchor/YAML checks are separate from native Go execution. Next classify actual exact-head feedback, apply justified corrections and verify full native/source/protected/accepted-push proof.

## Numerical domain correction recorded before its production edit

Review reproduced IEEE754 accumulation of eight valid nonnegative historical weights and free target prices yielding -100.00000000000001 instead of mathematical -100. Project correctly refuses values below-100, so direct division can reject valid all-free evidence. Enforce the proven nonnegative-price/nonnegative-weight lower domain with math.Max(-100,weighted/total); no upper clamp, no accepted negative evidence or omitted warning. Independent fixed literal weights test expects complete -100/free result; a compiling removed-domain-floor control must fail it on both native hosts. Python reproduction is arithmetic evidence; native Go failure without floor is not claimed until that control executes. This adds a fifteenth runtime control.

## Corrected candidate recovery

Initial fb21aab6816dca9adeaeaf7c4f037438c830b470/tree1ece35b6 passed strict source37237310705/job111538915495, all11/109 exact fresh bytes/hashes, but native37237310701/Linux111538916983/Windows111538917181 failed the ownership fixture panic (FN079). Source-literal restricted/zone detail arrays are empty; only the separate ownership probe now explicitly populates both and adjusts total/confirmed SKU arithmetic. No production/source guard was relaxed. FN078 lower-domain correction and fifteenth compiling control are separately reviewed. Executable Go1.26.8/gofmt recovered at /workspace/scratch/d06563362caa/cloud-assess-recovery-20261003/tools/go/bin, verified actual version; current focused cost runtime tests pass locally. Local package race/control results remain pending; Windows/pwsh/materialized APRL absent. Earlier PATH-only unavailable Go descriptions are historical (FN008). Fresh full candidate native/source evidence remains required, neither initial source success nor focused local pass transfers acceptance.

Local corrected evidence: focused TestCostRuntime suite and full region package race tests passed using actual recovered Go1.26.8. Free-roundoff removed-floor control compiled and failed its independent complete-comparison regression, confirming FN078's native arithmetic boundary. First full controls stopped at zone-ownership because removing maps.Clone left maps import unused (FN072); rejected, not counted. Correct isolated mutation preserves the import reference while returning the caller map; full15-control result follows PR body. Fresh complete native/source acceptance still pending.

All fifteen corrected runtime mutations compiled, failed their selected named assertions without build failures/panic, and restored/passed local baselines. Local region race suite passed. This is focused Linux development evidence; full fresh Linux/Windows/source candidate and distinct accepted-push evidence are still mandatory. Final candidate/merge/run proof will be indexed in live PR123 body to avoid evidence-only commit loops; determine acceptance there and from live refs. Current pending historical entries remain preserved.

## Final control scope and evidence authority

Final suite has18 compiling controls, adding exact nested-entry/unique-region guards and incomplete-history computation rejection to the earlier15. These are independently triggerable guards; their named boundary/health tests distinguish removed admission/completeness. Full local control result and exact final Linux/Windows/source jobs/accepted merge/push proof are indexed in PR123 live body. Handover/ledger/roadmap preserve failed and intermediate heads; no earlier pass transfers to final head. Remaining adapters/public/live/Gate004/release obligations remain unchanged. No further target production edit occurred in this verification update.

Final local evidence: all18 compiling runtime mutations failed their selected named tests without panic/build failure and restored/passed. Earlier focused cost suite and full region race tests passed; runtime Go production bytes are unchanged since that race pass. Full candidate native/source/package and separate accepted-push acceptance remain pending in PR123 evidence index.

## Tiny-positive historical weight correction before production edit

Automated Codex review on cdfc53bc identified P2 https://github.com/DeBoX85/Cloud-Assess/pull/123#discussion_r4179481979: admitted math.SmallestNonzeroFloat64 historical cost with source1/target1.001 underflows diff*weight to0 while total remains positive, falsely reporting0%. This is an unaccepted target numeric defect, not a service interruption or historical live evidence. Adopt the review's documented-minimum alternative: MinCostWeight=1e-12 in declared historical monetary units. Exactly zero retains source unit-weight1; positive history must be>=1e-12 and<=1e12, finite. Retail prices retain [0,1e12] and existing near-zero rules. This engineering admission bound avoids subnormal weighted products for admitted nonzero price differences (baseline>=0.0001, representable adjacent values). No inferred currency/rounding/conversion, no silent truncation or false0% success. Future collector must explicitly propagate rejection/health and declare currency/UoM separately. Full nil-result error [value_invalid] below minimum.

Add independent exact-minimum complete comparison and one-below/subnormal/nonfinite rejection, plus nineteenth compiling removed-floor control that fails the named boundary test and restores passing baseline. Ordinary source arithmetic/goldens stay unchanged. Automated review is useful independent tool feedback on an earlier implementation commit, not independent-person approval or final-head approval. Final candidate source/native/19controls and protected/accepted-push proof are still mandatory.

FN080 corrected local proof: focused cost runtime suite, full region race suite and all19 compiling controls/restored named baselines PASS using actual Go1.26.8. The tiny-history control specifically starts with math.SmallestNonzeroFloat64, removes only the minimum guard and fails named TestCostRuntimeHistoricalWeightFloor, so compiling zero-percent acceptance cannot hide under another admission check. Exact1e-12 full comparison and immediately preceding representable value are independently tested. Fresh final-head native/source/push proof remains mandatory. Previous18-control Linux111540661062 passed complete native run37237906970 but does not certify this production correction; Windows111540661216 was still running before new publication. No test-gate/pin/source-capture relaxation.
