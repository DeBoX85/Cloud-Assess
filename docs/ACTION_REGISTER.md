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
