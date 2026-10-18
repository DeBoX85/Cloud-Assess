# Active action register

Current status authority: [CURRENT_STATE](CURRENT_STATE.md). Detailed evidence belongs in the linked PR/contract; this register is a compact index, not a replacement ledger. Update after material progress and at checkpoints. Responsible party is the agent unless an independent reviewer, authorized operator or release owner is explicitly required. Status labels do not imply protected acceptance. Keep closed rows for traceability; archive older closed rows by link when needed.

| ID | State | Responsible party | Dependency / next action | Closure evidence |
| --- | --- | --- | --- | --- |
| DC-001 | Verified offline | Agent | PR150 own-head QA completed; remain draft/unmerged | [PR150 final index](https://github.com/DeBoX85/Cloud-Assess/pull/150), head 4c1ddc0 |
| DC-002 | Verified offline | Agent | PR151 own-head QA completed; remain draft/unmerged | [PR151 final index](https://github.com/DeBoX85/Cloud-Assess/pull/151), head ba9eeb9 |
| DC-003 | Verified offline (PR152) | Agent | Retained bundles/verifier verified under [controls contract](DEVELOPMENT_CONTROLS_20261017.md); preserve final proof | [PR152 final index](https://github.com/DeBoX85/Cloud-Assess/pull/152) |
| DC-004 | Verified offline; operational limits open | Agent | Publish rollback runbook and retained local rehearsal; distinguish remaining operational exercises | Controls record and retained command evidence |
| DC-005 | Verified offline assessment | Agent | Bounded assessment completed; no product schema change | Logging assessment with explicit limitations |
| DC-006 | Open | Agent / repository administrator if needed | Workflow failure-artifact retention is 14 days; actual repository-wide setting is not verified. Verify accessible setting, record owner/retention for off-repository backup | Explicit setting evidence and retention/restore record; CI links alone insufficient |
| DC-007 | Open | Independent reviewer | Before protected acceptance review gate logic and retained check inventory; no subagents authorized | Named independent review and disposition |
| DC-008 | Open | Agent | Region collector-to-report integration must cover success/partial/unsupported/failure, scoring/report/exit agreement and completeness | Future bounded integration contract and exact-head tests |
| DC-009 | Deferred | Authorized operator / agent | Representative live permissions/clouds/optional stages, nested non-production traversal and workload validation need actual scope/access | [Deferred validation](DEFERRED_VALIDATION.md) and separate live evidence |
| DC-010 | Open before release | Release owner / reviewer | Version-stamped releases, retained trustworthy artifacts, fresh OS, installed rollback, licensing/security/provenance/Gate004 | [Gate004](QUALITY_GATE_004_PLAN.md), [package limits](PACKAGE_BUILDS.md) and release decision |
| DC-011 | Verified offline PR154 | Agent | Pure VMSS diagnostic checkpoint and24new controls per platform complete; scoring unsupported; preserve exact-head evidence | [Implementation record](VMSS_CATALOG_IMPLEMENTATION_20261017.md), frozen contract and active PR final index |
| DC-012 | Open | Agent | PR152 retained local Linux SARIF reversal/restoration with compiling regression failures. Extend platform/installed/dependency breadth before broad rollback claims | Isolated rehearsal plus behavioral tests; current proof is Git/static only |
| DC-013 | Open design obligation | Agent | Evaluate optional safe diagnostic correlation only through a separate contract if operational requirements justify it; no blanket logger/product schema change now | [Logging assessment](LOGGING_ASSESSMENT_20261017.md), then bounded design/tests if authorized |
| DC-014 | Verified offline PR157 | Agent | Cognitive batch complete:143cases,46new controls/platform, full native/source/atomic source rollback verified; preserve draft head | [Three-task workstream](COGNITIVE_SKU_WORKSTREAM_20261017.md), active PR final index |
| DC-015 | Task1 verified offline; task2 frozen, own QA pending | Agent | PR158 complete; task2 316 independent future cases and primary integrity before runtime; full own QA then Database/MI with handover between each | [SQL workstream](SQL_CAPABILITY_WORKSTREAM_20261017.md), active live PR index |

## Per-change record

For a meaningful change record: stable task/finding ID; UTC timestamp; requirement and before/after behavior; author/actor and authorization; base/proposal/parent/tree; changed paths; dependencies; decisions and uncertainties; exact QA/run/job and evidence references; acceptance status; revert target/order and post-revert checks. Git author identity is traceability, not proof of independent review or publisher authenticity. Link existing journals/contracts instead of duplicating their content. For unknown remote outcomes, record intended target and safe inspection step until reconciled.

QA-SQL-01: task2 Windows license checkout correction; reproduced failure/fixed-byte baseline passed; full successor QA/archive pending. Task3/4 gated on task2 closure. See CURRENT_STATE and FAILURE_NOTES.

SQL task2 CLOSED offline at PR160/46223812, final own QA/archive/readback complete. QA-SQL-01/02/03 resolved/clarified with retained evidence. Same batch task3Database NEXT, task4MI pending; publish/inspect both-host missing APIs first. Independent/live/release rows remain OPEN.

SQL Database task3: local42compiling reversals/exact restoration and full SQL baseline complete; final candidate preflight/own hosted QA/atomic rollback/archive pending. Task4MI gated on complete task3 and handover. Final identity expectations497Linux/492Windows.

SQL task3 Database CLOSED offline at PR163/e2acce427/evidence026bbb61. Task4MI NEXT, gated by first163case runner and BOTH actual native missing constructors before production. Current state updated before transition; independent/live/release rows remain OPEN.

## Final MI candidate preparation after environment recovery, 2026-10-18

Fresh sequential42Database+36MI compiling reversals, named assertions, exact restoration, per-control baselines and final SQL baselines passed; focused Go1.26.9 region preflight/tests/vet passed. All316frozen cases and independent MI bounds/context/HTTP/modes/metadata/ownership/projection proof pass; first runners/oracles/Database42control bytes/pins/captures preserved. Both actual missing constructors39078083638/Linux117529286563/Windows117529286547 were inspected before production and complete logs reinspected on recovery. Older Windows local logs are unavailable in this Linux session; fresh retained proof replaces them without inheriting certification. Task3bundle110members verified before explicit restore. CURRENT_STATE/live final MI implementation PR own exact next action/outcomes. Final own hosted533Linux/528Windows controls/all8workers/bothgates/source/23captures, atomic MI+runner rollback/restoration, archive/readback/final handover remain pending at publication. Same final task4; tasks1/2/3complete; no new batch, critical input or authorization/boundary change.

## SQL completion and authorized next offline batch, 2026-10-18

## SQL batch complete, 2026-10-18 UTC

All four SQL tasks are COMPLETE VERIFIED OFFLINE: reconciliation PR158, frozen contract PR160, Database PR163 and Managed Instance PR165. Final product code is a3813439223ef6240e1aff8ae677e40a46011ce0, tree855d8a2994285f4b007616e06b3381dd93755503, sole parenta3738eec35bd0074b960a3ca072271e4e0818f80. PR165 remains draft/open/unmerged; protected bootstrap/core-v1 remains a6fee427e19feeef28691737ea9bc5026086efa0. Read [completion checkpoint](SQL_BATCH_COMPLETED_20261018.md) and [live final PR165](https://github.com/DeBoX85/Cloud-Assess/pull/165).

Own quality39095502229 passed all8native workers/both required gates, all10complete logs/steps reviewed. Exact533Linux/528Windows control identities conserved with no missing/extra/duplicates;78cumulative SQL additions include42Database+36MI named compiling reversals/exact restorations. Coverage86.4% with unchanged75%floor. Own source39095502254 reconstructed23captures/315chunks/909865bytes against exact Git blobs. All316frozen literal cases execute unchanged. Both native absent-API preparation failures were inspected before each implementation. Full cumulative20file readback, identical-tree/ordered-parent previewa8e249ca and atomic MI+runner rollback to Database tree6a9f557e/restoration/region tests passed.

Complete86member evidence archive is published and fully read back at0ec0d402392c056c0ab3970769d72d261d097324, tree7142fce4621e024ed6bff6a99221beedab01b610, parent exact final product code. Bundle2854937UTF8memberbytes/440377bytes, SHA256038c2617f4520dc12ffe1141591a35a01f3b5b4aa0aaed0410232bafb66b46dc. Bounded integrity verification precedes any explicit allowlisted restoration. Artifact and this administrative completion record are NOT product candidates or certification of new code trees.

The user authorized a DISTINCT next offline batch on2026-10-18: reconcile remaining work, define a bounded collector-to-report integration contract/independent expectations, implement the justified slice and complete its own QA/evidence/handover. First next action is repository/roadmap reconciliation, not assumed SQL scoring or public region selection. Publish/update handover BETWEEN each task. No critical input or Azure access is needed for reconciliation; resolve later scope using evidence and preserve unsupported eligibility/scoring. No new SQL four-task batch.

Current executor is Windows; recover from published refs, not old Linux scratch paths. Retained workspaces under C:/workspace/sql-recovery-20261017 and the original dirty checkout are preserved. This completion worktree is sql-final-closure. Previously unavailable Windows logs may now be present; final certification uses the retained exact-head archive, not memory. No uncertain remote mutation is claimed completed without readback. No Azure/fixture/RBAC/refresh/subagents/protected merge/bypass/release. Independent review/backup/broader operational/live/security/license/Gate004/release and one module-only advisory remain OPEN; accepted adapters1/7/source characterization7/7 unchanged.

## Offline integration reconciliation complete, 2026-10-18

SQL closure160c0793 and PR165 remain complete. New batch task1 reconciled live PRs/protected ref and current inventory/catalog/registry/result/rendering boundaries. [Reconciliation and ranked work](OFFLINE_INTEGRATION_RECONCILIATION_20261018.md) selects bounded internal explicit-coordinate SQL collector-to-report diagnostics. Task2 NEXT: freeze independent report contract/expectations before adapter production. Task3 implementation and task4 exact-head QA/evidence/handover follow. No automatic inventory match, deployment eligibility, scoring or public execution; no critical input/Azure access. Product baselinea3813439 remains unchanged. Preserve all prior heads/frozen inputs/pins and exclusions.

## SQL report contract freeze, 2026-10-18

New offline batch task1 reconciliation93500331 COMPLETE; task2 freezes [SQL report contract](SQL_REPORT_CONTRACT_20261018.md),26independent literal future report cases and an independent pure-output8MiB exact/+1 proof for BOTH families under all unchanged Catalog bounds. No report API exists and future cases are not runtime coverage. Byte-addressed expectation/budget Go files have exact -text checkout protection and integrity hashes. Existing153DB+163MI oracle/runner/production/capture/pin bytes unchanged.

QA-INT-01: direct sarif.go review confirms primary-findings-only output omits plugin rows/health. Preserve it and test no synthetic findings; JSON/tabular outputs carry diagnostic health. This supersedes the earlier proposed SARIF health wording. New APIs will prevalidate before transport, use existing Collect/Catalog, preserve full coordinates/modes, output only safe diagnostic tables and enforce atomic bounds/ownership/context; no public/scoring/eligibility integration.

Own contract-head QA is pending at publication; read its live draft PR for exact outcomes. After verified contract evidence/handover, task3 first publishes consuming runner while BOTH report APIs are absent and inspects actual Linux+Windows missing APIs before production edits. Task4 completes final own native/source/renderer/control/rollback/evidence/handover. Fixed late schedules Project22/Collect57 manually frozen before reversals. No critical input/Azure access; all prior exclusions remain.


## 2026-10-19 authorised programme

OP-001: task1 repository/evidence reconciliation complete; publication/readback pending. OP-002: MVP acceptance matrix mapped; administrative QA pending. OP-003: concrete ranked sequence/first reevaluation recorded; OP-004 independent cost-history contract next. User-authorised implementation follows without reapproval if justified and independent. See [programme](MVP_PROGRAMME_20261019.md). Earlier pending/task-next snapshots are historical; PR170 original SQL/report work remains complete.
