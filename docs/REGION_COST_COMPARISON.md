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

# Region CostComparison: bounded pure retail-price sheet

Status: VERIFIED OFFLINE through PR109 merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f. Original pre-acceptance contract below remains applicable; its pending preparation notes are historical. Baseline accepted PR108 merge1957f6d71da0c257bdd4b9d118c2be4b97b1a12e/tree3ffc1b51cb9ed5e608b805a3a7e2e1c1d54ba81e. Public region-selection remains unavailable. This is a pure decoded-data helper, not pricing collection, source arithmetic or financial validation.

## Source and alignment before production edits

Pinned AZQR8e4f0577f3615e6c9014c031bcad079f235369cc internal/scanners/plugins/region/output/output.go BuildCostComparisonSheet and types/types.go define behavior. Retained source-aux-inputs.json/source-aux-outputs.json hashes f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28/efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06 are unchanged. Decode only cost-full (one table) and cost-nil/cost-no-meters/cost-no-prices (null) from the actual fourteen-branch capture; verify branch names against retained keys. Every source header/description/row/cell is an independent oracle.

Before: region has primary, quota/reservation and service-availability pure projections, with no CostComparison helper. After this slice: ProjectCostComparison consumes selected display names and already decoded CostSheetData, returning an owned canonical table or source nil branch. Intended files: internal/plugins/region/cost.go/cost_test.go, scripts/tests/region-cost-mutation.py, both native workflow jobs and continuity docs. No network/credential/client/retail API/Cost Management/history/calculation/public registry/CLI/report integration, new dependency/source pin, platform or live scope is added.

Preserve CostComparison sheet name and description "Retail price comparison for resource meters across Azure regions". Fixed columns MeterId, ServiceName, MeterName, ProductName, SKUName; append lexically sorted <region>-RetailPrice columns from every pricing-map entry, including unknown meter keys. Sort output meters by MeterID. Prices strictly greater than zero use Go %.4f; zero/negative/absent prices are blank. Tiny positive values may display0.0000 and must stay distinct from an absent blank. Supplied meter ProductName/ServiceName source fields are not fallback metadata.

Metadata matching is exact/case-sensitive on MeterName/ProductID/SKUName. For each price item in received order, source scans meters in received order and stops at the first matching meter, even if its metadata already exists. Only the first item's service/product names for that meter are used. Different meter IDs sharing a tuple do not all receive metadata; preserve the first-input association with an indexed tuple lookup without source quadratic work. Output sorting must not change that received-order decision.

Nil data, zero meters or no regions has no table. Validate all provided input first so malformed/oversized data cannot masquerade as absent success, an explicit safety correction. Pure-completed health certifies projection only; later collection must establish request/currency/unit/coverage health. No price/currency normalization or retail/sovereign equivalence is asserted.

## Input, corrections and ownership

CostSheetData declares contributor SubscriptionIDs, MeterInputs, RegionPricing map[meterID]map[region]float64 and ordered PriceItems. CostMeter has MeterID/MeterName/ProductID/SKUName; CostPriceItem has MeterName/ProductID/SKUName/ServiceName/ProductName. Only source-consumed fields are represented. Labels are safe UTF8, at most512 bytes with no controls/replacement/U+FFFE/U+FFFF; MeterID and pricing-map IDs are nonempty, other source fields may be empty. Regions are lowercase ASCII at most64 bytes. Prices are finite with absolute value at most1e12, including negative values preserved as blank. Duplicate meter IDs reject rather than source unstable equal-ID ordering; distinct IDs sharing metadata tuples are permitted.

Selected UUIDs/names and declared contributor UUIDs must be valid and nonempty names; aliases/duplicate contributors/foreign declarations reject. Nonempty meters require contributor identities. Declarations are not proof of actual meter ownership; later selected collection/coordinator must construct the aggregate. Aggregate output has no invented row UUID. Do not confuse display names with identities. Inputs remain stable during projection; output owns cells/columns/metadata/health and no mutable per-run cache is shared.

Canonical schema1.0, owning region Metadata, ID cost-comparison and source name/description are deliberate metadata corrections. Errors return one fixed five-column safe header-only failed table with sanitized code/message, no provider text or partial rows. Cancellation preserves context error identity and drops candidate rows.

## Bounds and acceptance

Provisional ceilings:1000 selected/contributor identities,8192 meters,32 output regions,65536 structural entries globally including selected/contributor/meter/price-item/outer-pricing/inner-region entries. Cap all lengths before indexing, bound decoded label text independently to16MiB minus64KiB and charge exact projected cell bytes across rows before report-row allocation, reserving64KiB for fixed bounded table/header/health metadata. Individual safe512-byte labels and bounded numeric cells stay below32767 UTF16 units. Whole-set canonical validation still follows. These are target bounds, not load certification; no silent truncation.

Require full captured cell/header/sheet/description/empty/hash comparisons, received-order first-tuple/first-item exact matching, case/blank/zero/negative/tiny-positive pricing, unknown-meter region union, duplicate/scope/malformed/nonfinite/count/entry/region/text bounds, ownership/repeated concurrency/cancellation and nine compiling controls with restored named baselines: empty selected names, foreign contributors, meter/entry/decoded-text limits, first received meter/item metadata, zero-price absence and nonfinite values on both native hosts. Preserve unchanged source characterization and all existing controls. Every final code/docs head requires successful source/Linux/Windows mandatory steps/full logs, complete diff/bytes/tree/parents/human identity/base/preview/rules and expected-head protected merge. Python syntax/byte proof is distinct from unavailable local Go/source evidence. All laptop/Azure/Gate004/release/module-only-advisory deferrals remain open.

Rollback is a protected reviewed revert to PR108 after dependency checks. Following work is separate Inventory and remaining availability/latency/cost/quota/reservation calculations, then bounded ARM/retail adapters with truthful health and public/all-format execution. No simplified region score or pricing completeness is substituted.


PR108 exact candidate/accepted-push evidence: Service availability is VERIFIED OFFLINE through [PR108](https://github.com/DeBoX85/Cloud-Assess/pull/108), merge1957f6d71da0c257bdd4b9d118c2be4b97b1a12e/tree3ffc1b51cb9ed5e608b805a3a7e2e1c1d54ba81e, ordered3047085a/0b7c95cd parents. Final0b7c95cd6d2bb8b9c7a60149b3694a9c898fb67d native [37167296996](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167296996), quality111332833289/Windows111332833405, and source [37167296991](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167296991)/111332794154 passed. Both tested previewfefd1118 with exact candidate tree/parents. Distinct accepted-push native [37167713486](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167713486), quality111334058385/Windows111334058216, and source [37167713500](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37167713500)/111334058367 passed on exact1957f6d. All mandatory steps/full logs inspected; only conditional failure upload skipped. Both hosts passed10 service/five primary/three auxiliary/eight AI request/four execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages; Linux race/vet/fuzz and82.4% coverage, zero reachable/imported findings. Source14 branches/indexed complete chunks/exact retained hashes/bytes/source-tree-APRL provenance verified for both candidate and accepted push. Complete13-file scope/remote bytes and original scratch code/test/script Git blob hashes, human commit identity/base/head/preview/rules and protected merge/ref/tree/parents verified. FN070 Unicode transport and FN071 empty-label findings corrected before acceptance, with actual escaped BMP/non-BMP32767-unit boundaries and byte/rune controls. Automated review is earlier-head evidence, not independent-person/final manual-tool approval. No local Go/source fetch, public region, live/Gate004/release or module-only advisory closure is claimed.

Historical pre-acceptance checkpoint: Cost code/tests/nine controls are prepared; Python syntax and original UTF8 Git blob hashes verified. New Go/native/source evidence is pending. Source strings and capture bytes were not hand-edited; Go Unicode escapes and Python ensure_ascii preserve test characters (FN070 prevention).


Historical FN072 checkpoint, superseded by accepted proof below: PR109 initial native37168867529 failed only at the selected-contributor mutation compile error after passing package tests (FN072); no acceptance or skipped-step success is claimed. The corrected control retains selected through (!selected && false). Production/source/capture/pins are unchanged. Fresh complete-head source/native gates and all nine compiling controls/restored baselines remain required.


## Accepted proof and next task

CostComparison is VERIFIED OFFLINE through [PR109](https://github.com/DeBoX85/Cloud-Assess/pull/109), merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents. 

Final candidate1cf8e4cf6adc38ffe377894824b51f3b26431373/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, parentca16157 on accepted1957f6d7. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165402 quality111338402524/Windows111338402644 passed every mandatory step; only conditional failure upload skipped. Both tested preview77c2b59d40f89e1966402668a98efcdfad2e4d5d with identical final tree/ordered base/head parents. Full logs inspected: all nine Cost controls/restored baselines, ten service/three auxiliary/five primary/eight request/four execution controls, actual CLI/default+branded package/docs/provenance/module/inventory/branding guards; Linux race/vet/fuzz and82.5% coverage; zero reachable/imported vulnerability findings, existing module-only advisory remains open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165342 /111338403371 completed; full14 branches, indexed3-input/4-output chunks, exact retained bytes/hashes and source/tree/APRL provenance verified. Complete14-file scope/deletions/remote bytes/original scratch blobs, parent/human identity/pins/base/head/preview/current no-bypass rules23890737 verified. Semantic self-review checked actual unchanged source first-tuple/item break behavior, all source topology/cells and safe bounded identity/nonfinite/text/cancellation behavior. No callable final manual Code Review tool or independent-person approval is claimed; automatic comments currently empty. FN072 compile failure was not accepted. Expected-head protected merge and distinct accepted-push evidence follow; no fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure.

Protected acceptance verified: PR109 merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents, Denis author/GitHub committer/livecore ref and no open proposals. Distinct accepted-push native/source QA remains pending; no dependent live proof claimed.

Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513599 quality111339415060/Windows111339415184 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513582 /111339414976 passed on exact6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f. All mandatory steps/full logs inspected; only conditional Linux failure upload skipped. Both hosts passed9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages/docs/branding/provenance; Linux race/vet/fuzz82.5%, zero reachable/imported findings. Source14 complete branches, indexed chunks/both retained exact bytes/hashes and pinned source/tree/APRL provenance verified. Module-only advisory and all live/public/laptop/Azure/Gate004/release obligations stay open. No fresh local Go/source or independent-person approval claimed.

Current independent task: pure Inventory under [REGION_INVENTORY.md](REGION_INVENTORY.md), unaccepted; then remaining calculations, bounded ARM/retail adapters and public/all-format integration.
