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

# Region owned inventory aggregation

Status: VERIFIED OFFLINE through PR112 merge07011b63440e69f4a5acb128e0449943f759ee9c/tree72bebbf4885a15a74a5f43dfe57c68fe979c1e7c. Public region-selection unavailable.

## Source, requirement and intended before/after

Pinned AZQR8e4f0577f3615e6c9014c031bcad079f235369cc selection.go buildInventoryForSubscription/mergeInventory and types.NormalizeRegionName are the authority. PR111 retains actual13 inventory/merge/nil branches and30 normalization/predicate cases, input cdaa6b25e3250e6d701999279f96d5ea225c927e1249ee43c32cfe0b5072b17c/output756647e2981f8cd0ac4f0d5091ba6b4a74345ad6266de1bc15f8c833afeb2006; every fresh source run requires four captured files byte-for-byte. Source counts resources once (ignores capacity/tags/group/name/tier/family/SLA), lowercases type without trimming, preserves raw SKU case/whitespace and removes only ASCII spaces then Unicode-lowercases locations. Empty SKU has no SKU map; empty/logical locations remain private count keys. Merge adds counts across selected source inventories.

Before: accepted pure sheets and source calculation oracles, no target aggregate constructor. After: CalculateInventory(ctx, selected display-name map, []assessment.Resource) returns an owned InventoryCalculation containing five-map InventoryCounts per canonical selected UUID, a separately owned Aggregate, actual sorted contributor UUIDs and Records count. InventoryCounts matches all five source map roles using int64 values. Empty calls produce all selected five-map empties plus an empty aggregate, no nil-source panic. No report/public table/schema/CLI/registry, discovery/SKU/ARM/retail/credential/client/network/predicate/availability/latency/cost/quota/reservation arithmetic or dependency/pin/platform change. Source physical blacklist is characterized but not yet migrated; do not invent Azure catalog completeness.

## Explicit safety and identity corrections

Require valid selected UUID/nonempty safe display names, at most1000 with no case aliases, even on empty input. Every resource UUID must be selected and match its resource-ID segment using the accepted Inventory identity/masking-prefix guard. Foreign resources reject, rather than source silently filtering a global list; caller supplies selected decoded resources. Lowercase UUID identity matches accepted current selection behavior and intentionally corrects source's case-sensitive filtering: the source lowercase-alpha branch is empty while its uppercase branch has one resource; target correlates both spellings and must match the uppercase branch's real count. Reject duplicate IDs case-insensitively rather than double-counting; no mutable external merge API or unproven contributor declaration.

Read only ID/SubscriptionID/Type/Location/SKUName. Type and ID nonempty; Location/SKU may be empty; other resource fields are ignored exactly as the source aggregation does. Raw consumed strings are safe UTF8<=512 bytes with no controls/replacement/U+FFFE/U+FFFF; thus source control-bearing location input fails explicitly, not normalized into a different region or partial inventory. Standalone private normalization behavior is tested against all30 source vectors including controls, while the constructor rejects unsafe records before aggregation. Normalized type/location keys also must remain safe<=512 bytes after Unicode casing, an explicit transformation-bound correction. Preserve valid raw SKU keys, formula-like strings and empty/global/logical/nonASCII-space location keys privately; later request/table composition must establish representable region/service health without silently dropping them. Input values stay stable during the call.

## Bounded owned result and acceptance

At most8192 resource inputs/1000 selected UUIDs. Validate raw/normalized label text and all identity/duplicate inputs before counter-map assembly; decoded text<=16MiB-minus64KiB. Output structural entries globally<=65536 including subscription/contributor keys and every new nested map/count key in both per-subscription and aggregate maps. Charge actual output key bytes with the same text ceiling before each insertion; no truncation or returned partial maps on cancellation/limit/malformed error. int64 counts cannot overflow under the resource cap. Each selected inventory and aggregate owns all five maps; sorted contributors contain only UUIDs with actual validated resources. Records is input projection count, not collection completeness. Fixed sanitized errors return nil result; cancellation preserves context identity. No shared per-run state/caches/alias between aggregate/subscription/caller/another invocation.

Errors: nil result with fixed scope_invalid/input_limit/input_invalid/text_limit/output_limit/output_invalid code text; cancellation retains ctx.Err identity. An output_invalid defensive lookup protects assembly if preflight is inadvertently weakened. Counts are int64; structural limits include selected UUID keys and actual contributor UUID keys as well as both inventories' nested keys. Decoded scope display names count toward input text; output names are absent, so output charges canonical IDs instead.

Acceptance: actual captured subscription-two and uppercase-alpha complete five-map oracles; source lower-case mismatch explicitly corrected, normalized30 vectors, empty five-map shape/selected/no-contributor behavior; literal mixed type/SKU/location/resource-not-capacity counts, aggregate-vs-subscription ownership, ignored fields, raw whitespace/case/empty/logical inputs; foreign/mismatched/duplicate/alias/names/raw+normalized safety/row/scope/global-entry/decoded+projected budgets; owned16-concurrency/input/output/cancellation. The complete malformed source resource fixture must fail with no partial output, acknowledging its captured tab key rather than hiding it. The mixed six-resource literal valid subset explicitly excludes only the captured tab-bearing record; this is a separately specified safety-valid fixture, not altered source oracle bytes. Literal8192 resources/1000 scopes/exact65536 output entries and1813 versus1814 projected512-byte-key records exercise independent boundaries; repeated labels with5400 decoded512-byte IDs isolate decoded text from output limits. Ten controls cover empty selected names, selected-resource preflight, embedded UUID, row cap, global entries, decoded text, projected key text, ASCII-space normalization, raw SKU keys and resource counts. Every control must compile, fail named assertions with no panic/build failure, restore source in finally and pass restored baselines on both native hosts. Fresh final-head source14+13/30/four files, full native mandatory jobs/logs/race/vet/fuzz/built CLI/package/security/all prior controls/complete scope/bytes/tree/parents/human identity/base/preview/rules and expected-head protected merge/accepted-push required. No local Go/source, live, load, public, Gate004/release/advisory closure from this pure library.

Rollback protected reviewed revert after dependency review. Next availability/latency/cost/quota/reservation arithmetic and bounded request adapters/coordinator defaults/empty/partial-health/public/all-format execution. User Azure/laptop-dependent tests remain deferred; no input needed now.


## Current checkpoint (2026-10-04)

Owned inventory aggregation is VERIFIED OFFLINE through [PR112](https://github.com/DeBoX85/Cloud-Assess/pull/112), accepted07011b63. Pure source availability characterization is IN PROGRESS/unaccepted under [REGION_AVAILABILITY_CALCULATIONS.md](REGION_AVAILABILITY_CALCULATIONS.md). Public region-selection remains unavailable. Retain actual new source bytes and require unconditional six-file equality, independent oracle checks and both complete native gates before protected acceptance. Then implement bounded availability, followed by latency/cost/quota/reservation arithmetic, adapters and coordinator/public/all-format execution. Azure/laptop validation remains deferred.

## PR112 inventory aggregation accepted (2026-10-04)

[PR112](https://github.com/DeBoX85/Cloud-Assess/pull/112) is VERIFIED OFFLINE. Accepted commit `07011b63440e69f4a5acb128e0449943f759ee9c`, tree `72bebbf4885a15a74a5f43dfe57c68fe979c1e7c`, ordered parents `1b17cccd0bb4ff9f08627d4ea398d3e628107afb` and `8eb72c4dd49d3f8ceb400ab4ef8111d1ee04b3cf`. Denis author, GitHub merge committer and live protected ref verified.

Final candidate `8eb72c4d` passed [native run37172490247](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37172490247), Linux111348132665 and Windows111348132783, and [source run37172490240](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37172490240), job111348132671. All mandatory steps/full logs inspected; only conditional Linux failure-evidence upload skipped. All ten compiling aggregation controls/restored baselines and all prior controls passed on both hosts, along with actual CLI/default/custom packages, docs/branding/provenance/module/inventory. Linux race/vet/bounded fuzz and82.7% total coverage passed. Zero reachable/imported vulnerability findings; existing module-only advisory remains open. Source14 auxiliary/13 inventory branches/30 normalization vectors and all four retained files reproduced exactly from complete3+4+2+4 indexed chunks, with pinned source/tree/APRL and unchanged production/module/APRL proof. All15 changed files/remote hashes, three original scratch blobs, complete scope/pins/identity/base and tested preview `3eee48d6` with identical tree/ordered parents verified. Self-review inspected actual source/target contracts and bounds; no independent-person review claimed. Inline automatic threads empty.

Distinct accepted-push proof: exact07011b63440e69f4a5acb128e0449943f759ee9c passed native [37172980065](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37172980065), Linux111349577702 and Windows111349577424, and source [37172980083](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37172980083), job111349577747. All mandatory steps/full logs inspected, only conditional Linux failure upload skipped. Both passed ten aggregation and all prior compiling controls/restored baselines, actual CLI/default/branded package checks, docs/branding/provenance/module/inventory and the real byte-guard tests/control/restoration. Linux race/vet/bounded fuzz82.7%, zero reachable/imported vulnerability findings. Source14 auxiliary/13 inventory/30 normalization and all four retained files exactly reconstructed from contiguous3+4+2+4 chunks with pinned source-tree-APRL/isolation. Module-only advisory and Azure/laptop/live/public/load/Gate004/release remain open. No local Go/source/Azure/laptop/live/public/Gate004/release/advisory closure.

Next: source-only availability arithmetic characterization under [REGION_AVAILABILITY_CALCULATIONS.md](REGION_AVAILABILITY_CALCULATIONS.md), then bounded target availability. No public execution is added. Protected reviewed revert after dependency review is the rollback route.
