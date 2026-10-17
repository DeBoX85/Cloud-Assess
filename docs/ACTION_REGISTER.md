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
| DC-011 | Contract prepared; own QA pending | Agent | [Catalog-only contract](VMSS_CATALOG_CONTRACT_20261017.md) and61future scenarios prepared; complete own QA before runner/API work; scoring unsupported | [VMSS research](VMSS_CATALOG_RECONCILIATION_20261017.md), contract and active PR final index |
| DC-012 | Open | Agent | PR152 retained local Linux SARIF reversal/restoration with compiling regression failures. Extend platform/installed/dependency breadth before broad rollback claims | Isolated rehearsal plus behavioral tests; current proof is Git/static only |
| DC-013 | Open design obligation | Agent | Evaluate optional safe diagnostic correlation only through a separate contract if operational requirements justify it; no blanket logger/product schema change now | [Logging assessment](LOGGING_ASSESSMENT_20261017.md), then bounded design/tests if authorized |

## Per-change record

For a meaningful change record: stable task/finding ID; UTC timestamp; requirement and before/after behavior; author/actor and authorization; base/proposal/parent/tree; changed paths; dependencies; decisions and uncertainties; exact QA/run/job and evidence references; acceptance status; revert target/order and post-revert checks. Git author identity is traceability, not proof of independent review or publisher authenticity. Link existing journals/contracts instead of duplicating their content. For unknown remote outcomes, record intended target and safe inspection step until reconciled.
