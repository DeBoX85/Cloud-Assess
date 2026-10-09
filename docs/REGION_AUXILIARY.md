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

# Region quota and capacity reservation pure projections

Current audit disposition (2026-10-04): whole-project bounded offline review and hardening/source characterization are accepted through PR114, mergec33fb2eee93a2e5730f72612b351d0b838737d03, tree583496773c23902e7d7187ffc83a1c4ee6be0713. Owned inventory aggregation remains accepted through PR112; source-only availability characterization is now accepted through PR114. Runtime availability and public region-selection remain absent. Resume the separately reviewed bounded pure availability implementation after the audit disposition follow-up is accepted. Exact runtime and disposition acceptance evidence is indexed in their PR bodies; live/load/Gate004/release/advisory boundaries remain open. Earlier candidate/IN PROGRESS sequencing below is historical.

Current status: pure auxiliary sheets and owned inventory aggregation are accepted offline through PR106/108/109/110/112. Earlier sequencing is historical. Availability source characterization is proposed in the combined audit candidate; runtime availability/request/public execution is still absent.

Status: VERIFIED OFFLINE through PR106 mergec4643527dc948f96a565da237662a104a084f14e. Public region-selection unavailable. Original pre-acceptance requirements below are the continuing contract; completed proof follows.

## Contract and source authority

Pinned AZQR8e4f0577f3615e6c9014c031bcad079f235369cc internal/scanners/plugins/region/output/output.go BuildQuotaSheet/BuildCRGSheetFromRows define labels, descriptions, received ordering, formatting, flag precedence and nil empty branches. [Source characterization](REGION_SOURCE_CHARACTERIZATION.md) calls unchanged functions with literal synthetic inputs. Retain source-aux-inputs.json/source-aux-outputs.json exact bytes covering all fourteen branches; implement only quota/reservations here.

Before: primary scoring only. After: ProjectQuota and ProjectReservations accept stable already-decoded selected inputs, return optional owned canonical tables, no HTTP/discovery/cache/CLI/registry change. Empty success returns nil, matching source; failure returns zero rows/source headers and failed health with fixed error text. Canonical metadata is the owning region Metadata(), IDs quota/reservations, schema1.0, replacing source helper default/absent metadata deliberately.

Quota preserves Current/Limit/Available/HeadroomPct and flags without recalculating arithmetic or deriving flags from rounded text. Over-limit wins, then near-limit, then OK. Reservations preserve ten supplied source cells, signed availability and status. Service decoders/arithmetic consistency remain future contracts. A required out-of-band SubscriptionID UUID must be selected; display Subscription must match exactly. Attach canonical lowercase UUID for later privacy processing, preserving cells. Same names across distinct UUIDs remain valid.

## Input corrections and bounds

Reject foreign/malformed/case-alias duplicate scope IDs; duplicate quota tuple(UUID,region,type,resource) or reservation tuple(UUID,region,resource group,group,name), using structured keys; control/replacement/invalid UTF8 labels; invalid ASCII regional IDs; nonfinite/excess numbers; invalid reservation width/count/status; excessive rows/text. These deliberate pure-input corrections are not live/parser parity.

Maximum8192 rows per call,1000 subscriptions,512 UTF8 bytes per label, counts within1e12 (Current/Limit/Reserved/Allocated nonnegative; Available signed), finite HeadroomPct within plus/minus1e12. Reservation numeric strings must be canonical base10. Text budget before row allocation counts scope names and all supplied labels/cells, capped at MaxPluginTextBytes minus64KiB for canonical metadata/formatting. Final canonical validation remains mandatory. No truncation or partial successful rows. Empty input still checks scope/context. Caller input must stay stable during a call; no input map/slice escapes.

## Acceptance and boundaries

Compare every quota/reservation source cell, columns, description/order and nil empties; verify hashes, canonical metadata/health/correlation, precedence independent of rounding, negatives, same-name scopes, foreign/malformed/duplicates, finite/count/row/aggregate budgets, owned columns/cells, cancellation and repeated/concurrent isolation. Both required hosts run compiling selected-scope/work-limit controls with named failures/restored tests, alongside all existing guards and characterization. Full native final-head Linux/Windows jobs/logs, semantic review, exact tree/parent/identity/bytes/base/preview/rules and expected-head protected merge remain required.

No fresh local Go QA, Azure/laptop, quota/SKU/cost/latency fidelity, public/all-sheet integration, load/fresh-OS, Gate004/release proof. Pins/notices unchanged. Rollback: protected reviewed revert after dependency review. Service availability PR108/CostComparison PR109 are accepted offline. Current separate pure Inventory under REGION_INVENTORY is unaccepted, followed by calculations/bounded service adapters/public execution.

SHA256 inputs f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28; outputs efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06. [Microsoft capacity documentation](https://learn.microsoft.com/en-us/dotnet/api/azure.resourcemanager.compute.models.capacityreservationutilization.currentcapacity?view=azure-dotnet) distinguishes billed reserved capacity from allocated resources; it does not justify modifying the captured helper cells or adding modern API fields here.


## Accepted proof

Quota/Capacity Reservations pure projections are VERIFIED OFFLINE through [PR106](https://github.com/DeBoX85/Cloud-Assess/pull/106), merge `c4643527dc948f96a565da237662a104a084f14e`, tree `38248bc19be95c0058377ed373a4bf90b20e4ca2`, ordered parents7e4ca5c/e7e5511. Final head e7e551104cc6e0e6dda25c389af9176eb5ed562c passed [run37146790757](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37146790757), quality111272343099/Windows111272342986, all mandatory steps/full logs inspected. Both tested preview06848f39 with the identical tree/parents. Source characterization [37146790810](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37146790810)/job111272298401 reproduced all14 unchanged helper branches with exact captured bytes/provenance. Both hosts passed every quota/reservation cell/empty branch, all three compiling auxiliary controls and restored assertions, existing five primary/eight AI request/four AI execution controls,12 actual CLI and19 default/19 branded packages; docs327/9/1. Linux race/vet/fuzz and82.1% coverage passed; zero reachable/imported vulnerability findings, existing module-only advisory open. Remote merge/ref/tree/parents/human author/GitHub committer verified. Fresh local Go/fetched-source evidence is unavailable in this Windows session.

Next [service-availability slice](REGION_SERVICE_AVAILABILITY.md) remains NOT STARTED. Final acceptance imperatives above are fulfilled for this implementation; apply them again to future changes.


Separate accepted-merge proof completed: [run37147153459](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37147153459), quality111273357208/Windows111273357343, every mandatory step/full log inspected on exact c4643527. Source [run37147153415](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37147153415)/job111273356990 reproduced all14 unchanged branches with exact hashes/provenance/bytes. Both hosts passed all existing/new controls,12 CLI/19 default/19 branded packages/docs327/9/1; Linux race/vet/fuzz/coverage82.1% and zero reachable/imported vulnerability findings. Module-only advisory stays open. This is separate accepted-push evidence, not inferred PR success or local execution.
