# Concrete offline sequence and first reevaluation

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Requirements: [acceptance matrix](MVP_ACCEPTANCE_MATRIX_20261019.md). Programme safeguards: [programme](MVP_PROGRAMME_20261019.md). User authorisation is the 2026-10-19 conversation. All rows are offline; no row grants Azure, protected acceptance, release, dependency/source refresh or new public schema authority.

## First reevaluation, after tasks1-3

The provisional plan's repeated "next missing collector" rows are superseded: all seven characterised families already have target draft collectors or diagnostic representations. Their remaining limitations are integration and eligibility semantics. Rebuilding them would add cost without closing a real gap. The remaining public plugin is region selection, which depends on currently absent region meter-history/retail requests and owned coordination. Move those before public exposure. Keep historical source captures and existing pure arithmetic intact.

Core-v1 and eventual full plugin parity stay separately visible. MCP/compare/alternative-VM are deferred by the target specification and do not displace integration work. No agreed requirement is deleted. If later public execution requires a schema refresh or an unresolved correctness/security choice, stop at that boundary instead of assuming authorisation.

| Task | Concrete deliverable and purpose | Required dependency / closure |
| --- | --- | --- |
| 1 | Verify live baseline, proposals, evidence and local/source state | Reconciliation in programme; full administrative QA remains separately tracked |
| 2 | Map MVP/core versus full parity requirements and evidence | Acceptance matrix published/read back; no weighted percentage |
| 3 | Rank the actual gaps and record reevaluation/supersedence | This sequence, research record and between-task handover |
| 4 | Freeze region meter-history query/decoder contract and independent literals, including UTC calendar boundaries, declared currency, column identity and bounded paging | Pinned source plus official API docs; separate from core Cost; no request to Azure |
| 5 | Implement the bounded meter-history request/collection slice, integrate existing cost evidence without false completeness, and complete own QA/rollback/evidence | Task4 independently frozen/verified; compiling reversals and both native/source gates; reevaluate here |
| 6 | Freeze retail metadata/pricing contract and independent literals for relevant meters/regions, currencies, UoM, tiers and conflicts | Current Microsoft documentation reconciled against pinned source; no preview/API refresh by inference |
| 7 | Implement bounded unauthenticated retail transport/decoder and cost integration | Task6; prevent ARM credentials reaching public pricing; partial/error states explicit; full QA/evidence |
| 8 | Freeze required-kind inventory/evidence coordination contract, explicitly retaining unsupported VMSS scoring and unresolved SQL coordinates | Inspect complete current inventory model and source; no SKU-name inference |
| 9 | Implement justified private coordinate/identity preservation and required-kind evidence admission | Task8; no public schema refresh; fail/unknown rather than invented eligibility; full QA/evidence |
| 10 | Implement bounded owned per-run coordinator across existing calculations/collectors with injected dependencies | Tasks5/7/9; preserve component health and caller cancellation; reevaluate here |
| 11 | Add supported region CLI/registry execution using the existing public contracts | Task10; early validation before authentication; standalone/mixed/scanner selections and source flag/default reconciliation |
| 12 | Complete required region output/health/exit agreement and ancillary parity reconciliation | Task11; all source-required tables including FullCost disposition; formats follow existing explicit contracts; types utility only if confirmed required |
| 13 | Exercise representative synthetic multi-subscription workloads, paging, cancellation, limits and injected permission failures end to end | Working supported path; retain reproducible fixtures and existing controls; no live scale claim |
| 14 | Rehearse packaged default/custom-brand installed workflows and supported rollback with available isolated platforms | Existing package workflow; record fresh-OS/operational limitations rather than pretending CI is a pristine customer machine |
| 15 | Reconcile exact evidence and residual advisory/security/licensing/live/independent-review obligations; publish readiness and restart checkpoint | No release or acceptance inference; reevaluate completed versus deferred scope and remaining operator decisions |

Rows4-15 are named deliverables, not a promise that each fits one commit or one session. Required preparation/native observations and QA are part of each task, not optional extra scope. Reassessment can split, replace or cancel a superseded row when evidence justifies it, with recorded rationale. Do not fill the list with speculative features or redundant tests. Normal checkpoints do not require user reapproval.

## Research record and remaining uncertainty

Reviewed pinned AZQR `internal/scanners/plugins/region/cost/cost.go` at8e4f0577: history anchors to first UTC day, queries ActualCost grouped by ResourceId/ResourceGuid/MeterCategory/MeterSubcategory/Meter, parses named columns and sums meter weights across resources; retail metadata and regional pricing are separate phases. Existing target `internal/cost/scanner.go` uses a different resource-cost result contract; region `cost_enrichment.go` already performs owned weighted arithmetic. Its comment explicitly leaves history ownership/currency/unit/clock/paging admission to a future collector. Do not substitute pure evidence constructors for actual collection.

Microsoft [Retail Prices documentation](https://learn.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices), retrieved2026-10-19Oslo, confirms an unauthenticated public API, distinct currency/type/tier/UoM metadata and NextPageLink paging. Current documentation is supporting evidence, not permission to change a pinned API or use retail prices as negotiated customer rates. [Cost query documentation](https://learn.microsoft.com/en-us/rest/api/cost-management/query/usage) supports the history-query investigation; resolve exact pinned SDK request/response and documented continuation before implementing. One version-specific web open failed; do not treat it as read evidence. Pinned source is the direct equivalence authority.

Task4 must resolve actual SDK version/API, history upper-bound behaviour (source comment says1-12 but function only clamps below1), calendar precision versus existing core end-of-month nanoseconds, mixed/absent currency, negative/zero costs and malformed-column behaviour. Safe candidate corrections require explicit independent expectations and rationale; no silent rewrite of source captures. These are bounded contract questions with established fail-closed/explicit-unknown approaches, not currently an unresolved critical blocker. No production cost changes have begun.

## Drift/conflict check

Purpose remains read-oriented assessment equivalence and useful reports with adjustable branding. No scope expanded. Provisional collector/ancillary order is superseded as above. No critical blocker found in initial reconciliation. Mandatory full exact-head QA remains pending for the administrative programme, and product work cannot inherit PR170 or predecessor certification. Before production, finish required administrative/contract validation and publish handover. Independent human review/backup/live/release/advisory obligations remain open.


## 2026-10-20 security crossroads

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Autonomous batch paused at task6 own security QA; [blocker/options/decision](RETAIL_SECURITY_BLOCKER_20261020.md) controls the next action.177contract is frozen but unverified;178first consumer intentionally lacks production APIs; tasks7implementation/8contract not started. No dependency/pin refresh or gate bypass. Earlier next-task/no-blocker statuses are historical, not permission to proceed through this pause.


## Approved narrow security maintenance, 2026-10-20

Denis authorised the bounded Excelize update; the earlier decision pause is historical within this exception. See [CURRENT_STATE](CURRENT_STATE.md), [remediation contract](EXCELIZE_SECURITY_REMEDIATION_20261020.md) and [advisory review](SECURITY_ADVISORY_REVIEW_20261020.md). Exact fixed commit plus necessary transitive updates only; own full native/source QA pending. No retail production/protected merge/release. Preserve original failed177and first consumer178.


## Task7 retail checkpoint, 2026-10-20

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Narrow security maintenance/task6 prerequisites are COMPLETE VERIFIED OFFLINE in180/9fe12612; final quality39440501477/source39440501617 and cf93d214evidence supersede dated pending/wait snapshots. Task7 private retail implementation is IN PROGRESS, preliminary52frozen+15prepared+additional boundary cases pass locally. Final controls/rollback/full own QA/evidence remain PENDING; task8 research only. Follow180explicit live successor and [implementation record](REGION_RETAIL_IMPLEMENTATION_20261020.md). Continue existing authorised independent batch, preserving all old heads/pins/gates and live/review/acceptance boundaries.
