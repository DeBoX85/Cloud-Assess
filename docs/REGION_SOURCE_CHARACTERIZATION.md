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

# Reproducible pinned region source characterization

Status: VERIFIED OFFLINE through PR105 merge7e4ca5c955abb7f882591a44d8bb2f74c93ae243. Public region-selection unavailable. See [REGION_SELECTION.md](REGION_SELECTION.md).

## Purpose and bounded contract

The local executor stopped returning after completed correction race/vet/negative-control QA (FN065). GitHub publication and native Actions remain usable. This slice creates reproducible development execution for the next source-grounded auxiliary migration; it changes no production assessment behavior or dependencies. It is not an unattended uptime guarantee or a replacement for the operator's laptop/Azure validation.

[The runner](../scripts/capture-region-source.py) creates a task-owned temporary checkout of DeBoX85/azqr8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1 and initializes actual APRL60eaddda76541f6adbc1c5ffa686829807e55e29. It verifies identities, adds only [the retained harness](../internal/plugins/region/testdata/aux_capture_test.go.txt) as an untracked test in that copy, calls unchanged pure source output helpers and asserts original source/APRL files and module graph remain unchanged. Temporary state is cleaned up by the runner; target/reference/user checkouts are untouched. Command/time limits and explicit GOTOOLCHAIN=local apply on the workflow.

The dedicated workflow uses contents:read, checkout without persisted credentials, already supported Ubuntu24.04/Go1.26.8 and the existing pinned checkout/setup actions. It runs for every core-v1 PR candidate and accepted push, including documentation-only follow-ups. Path filters are deliberately absent so earlier-head captures cannot replace final-head evidence. Existing protected required quality/windows-validation jobs remain mandatory and unchanged. Git public clone/module download are development inputs; no Azure API, credential, scan, role change or resource write is invoked. Standard HTTP tripwires complement inspection that these pure helpers construct no SDK/service client; this is not a universal SDK transport sandbox (FN052).

## Exact source helper scope

Fourteen branches call only internal/scanners/plugins/region/output/output.go:
- BuildSvcAvailSheets: nonempty two targets plus nil inventory/no results. Actual sheet prefix is Svc Avail followed by a space and target region, despite some historical comments using SvcAvail_. Test missing service, missing SKU, zone restriction detail stripping, unsupported SKU provider, no SKUs and source treatment of an unknown-suffix detail. Original provider registry is used without FetchSKUs.
- BuildCostComparisonSheet: nonempty, nil, no meters and no prices. Dynamic region columns and meter IDs sort lexically; metadata uses first matching meter name/product ID/SKU; zero price is blank and small positive price formats to four decimals. This is sheet formatting, not pricing retrieval/selection or weighted historical-cost calculation.
- BuildQuotaSheet and BuildCRGSheetFromRows: full and empty. Quota flag precedence, received ordering/rounding and negative headroom/availability cells are literal inputs; the harness does not establish service decoder/status arithmetic.
- BuildInventorySheet: raw, masked and empty. Includes source SKU-derived capacity and VM scale-set multiplication, formula-like synthetic name, subscription/resource-ID masking and ten literal column names. Full mixed-report Inventory collision behavior remains a separate integration obligation.

The literal synthetic typed inputs are serialized by the unchanged source module alongside complete outputs. No source tables/cells/metadata are normalized. Helpers return default/absent metadata on auxiliary tables; future canonical target ownership must be documented separately. This is not full source Scan, score input completeness, service availability, price/quota/latency accuracy or live equivalence.

## Evidence protocol and future runner acceptance

Successful source execution emits explicit source/tree/actual APRL provenance and two UTF-8 JSON files, source-aux-inputs.json/source-aux-outputs.json. REGION_CAPTURE_JSON log records contain file name, original SHA256, ordered chunk index/count and JSON-escaped content. Reassemble only after the entire job succeeds, all chunks/provenance match and source-copy guards pass; retain exact bytes and byte attributes, then compare every independently captured cell in a separate target migration. Log chunks avoid requiring a functioning local artifact downloader; they contain only synthetic fixtures.

Required review: inspect actual runner/harness/workflow diff; verify fixture branch coverage, no unintended HTTP/credential/SDK construction, exact pins, copied source cleanliness, bounded subprocesses and untracked-file guard. Native characterization must pass on final published head; ordinary Linux/Windows quality jobs must also pass. No source output/hash is claimed before actual successful execution. Final PR evidence records exact tree/parents/identity/current base/preview and protected acceptance.

After accepted runner/captures, migrate auxiliary pure tables/calculations in bounded, separately reviewable slices with complete source expectations, malformed/foreign/duplicate/limit/ownership/report and compiling guard controls. Retain both source Scan empty sentinel branches, original unstable ordering defects, singleton cache/state, unforwarded history flag, debug JSON privacy and mixed Inventory name collision as explicit later decisions. Bounded authenticated ARM versus unauthenticated retail adapters and actual public execution/report integration remain separate.

Microsoft Learn review on2026-10-03 consulted the [Retail Prices API overview](https://learn.microsoft.com/rest/api/cost-management/retail-prices/azure-retail-prices): the public commercial API is unauthenticated, prices are retail USD without negotiated discounts, responses page at1000 with NextPageLink and filter/version/currency semantics matter. This current documentation informs later destination/auth/access guards; it does not replace pinned SDK/source behavior, refresh embedded latency/zone data or certify live prices. Sovereign support requires a separate supported-contract decision.

## Historical interruption review and correction (superseded by accepted proof) (2026-10-03)

Initial PR105 head7b1fdede passed source characterization run37131392522/job111227049328 and required run37131392461/quality111227049234/windows111227049053. The automated review correctly identified that path filtering could skip characterization after a documentation-only follow-up. Both PR and push path filters are removed; branch filters, permissions, pinned actions, pins and gate strength are unchanged. This new complete code/documentation head requires its own characterization and both required jobs before acceptance. Earlier green results are historical evidence only. FN066 records the defect.

The executor responds again. A new isolated worktree was fetched at exact published7b1fdede without modifying the older dirty task checkout. Actual source8e4f057/APRL60eaddda, accepted mergec270ea22/tree/parents and Go1.26.8/PowerShell7.6.6 were verified locally. Earlier unknown command completion remains unknown; the recovered old edits are retained. Next reassemble final successful source output, verify hashes/provenance, accept this runner through protected merge and begin the auxiliary target migration.


## Accepted recovery proof

PR105 is VERIFIED OFFLINE at accepted merge 7e4ca5c955abb7f882591a44d8bb2f74c93ae243, tree 5c79e721d8281f04ba066006a4ff5601202abdbd, ordered parents c270ea22 / c9cc112. Final candidate run37143242571 (quality111261843942, Windows111261844116) and characterization37143242555/job111261844024 passed. Separate accepted-push quality run37143573050 (quality111262849038, Windows111262849202) and characterization37143573016/job111262848443 passed. Required steps and retrieved full logs were inspected in this recovery; only conditional failure upload skipped. Native primary/AI controls, compiled CLI/package cases, race/vet/fuzz, provenance and vulnerability scans passed; coverage81.8%, zero reachable/imported findings and the existing module-only advisory remains open. Final and accepted-push auxiliary capture chunks/provenance are complete and identical: inputs f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28, outputs efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06.

Prior candidate/pending paragraphs are historical; final evidence supersedes their acceptance imperatives. Next separate slice: retain exact source auxiliary fixtures and implement bounded pure Quota/Capacity Reservations under REGION_AUXILIARY.md. No service decoder/arithmetic, public region execution, live/Gate004/release closure.
