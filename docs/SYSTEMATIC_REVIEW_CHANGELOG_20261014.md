# Active Section 1 changes and evidence, 2026-10-14

Current state and next action are governed by [Section 1 report](SYSTEMATIC_REVIEW_SECTION1_20261014.md) and [live PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). This is a correction candidate; final full QA/review and section reconciliation are pending.

| Entry | Change | Commit and parent | Verification state |
| --- | --- | --- | --- |
| SR-03 | Original cancellation expectations before production edits | 432f2034aa3eed973670ab5b3a31cd1775d94611; parente4064ba12e2b2a7299bad093f7876cc73244b5bd | Both-host expected failures in38561749374 |
| SR-04 | First cancellation correction | 5dd97b0f530d578e17767a4eea14510b72ee6f6a; parent432f2034 | Original cases pass; later caller-priority review gap retained |
| SR-05 | Initial tied-page fixture, now superseded | e012a94d8077abf180917408156b12cfa69f47ec; parent5dd97b0f | Both-host failures recorded; old fixture too dependent on ID token, not final oracle |
| SR-06 | Unique inventory tie-breaker | 9d1872b1a1737da05db512addae8479e530aa238; parente012a94d | Stronger model/control verification follows |
| SR-07 | Standard/external import organization | a98f2170d39c9dd8a53335c2bd147f4f9fc6496b; parent9d1872b1 | Style only |
| SR-08 | First candidate and review | 1243461f8ee884d9b0071e13bb219ab8f403ac45; parenta98f2170 | Review5517988620 found three issues; not clearance; native38564627993 superseded |
| SR-09 | Parent-priority expectations and stronger modeled rows | a77e558f06bee40edc23f6e450f356f4cbc5bd87; parent1243461 | Caller priority fails and discovery model/calibration passes on both in38565950638 |
| SR-10 | Correct parent priority and add three restoration controls | e65e66cd74d8e08ebfe1ea4452f6391a65e12878; parenta77e558f | Three controls passed both hosts in38567759652; complete workflow still intentionally red |
| SR-11 | Completed-query HTTP status expectations before guard | 57ef9f93b9d0618e8fde87eb9323f91349b0ff54; parente65e66cd | Both-host six unexpected-status failures in38567759652 |
| SR-12 | HTTP200guard and fourth control | 1f93442eca7e949c9b31f7af20b594c91ff10b62; parent57ef9f93 | Final verification pending |
| SR-13 | Rewrite active report/resume steps, update coverage and clarify control counts | Successor of1f93442e; exact documentation commit from PR142/history | Readback and final candidate gates required |

Review issue disposition:4270844237(parent priority) reproduced and corrected;4270844250(model independence) corrected with value-based sorting/calibration and compiled counterfactual controls;4270844258(obsolete resume instruction) corrected by rewriting the current report. Do not resolve review discussions merely because changes are authored; use actual final corrective evidence.

DOC-02: old314Linux aggregate includes one shell control. Preserve the count, but describe it accurately as313Go+one shell before the four additions. Windows309was Go-labeled. This is counting terminology, not a dropped or failed control.

All source captures and previous expected source outputs remain immutable. Superseded self-authored target fixtures remain in historical commits for backtracking; they are explicitly replaced rather than presented as independent proof. Every control requires a healthy baseline, compiled failing variant, exact production-byte restoration and a healthy restored test. Go build failure, timeouts or missing diagnostic are not accepted detections.

Rollback: identify the exact isolated fix and parent, inspect later dependencies, create an ordinary revert on a task branch, rerun the meaningful regressions and required gates, and use protected PR acceptance. A reverted defect should fail retained expectations. Do not reset protected history, blindly revert a multi-change merge, rewrite fixtures or perform Azure rollback.

At this checkpoint, a scratch-only Python patch attempt was rejected because two operations targeted the same file; no repository mutation occurred. It was corrected with a normal update and local AST parsing passed. Hosted Go execution remains the authority; no local Go run is claimed. Complete reproduction-log hashes and distinctions between character/UTF-8 byte counts are in the Section 1 report.

Historical records below preserve earlier states and must not govern current restart steps.

---

# Systematic review findings and change journal, 2026-10-14

## Recovery pointers

Accepted code baseline: `a6fee427e19feeef28691737ea9bc5026086efa0`.
Audit branch: `audit/systematic-module-review-20261014`.
Initial inventory checkpoint: `cb9d18669e7f9fefb997b53062c2907a14911a49`, parent accepted baseline.
Read the latest branch head and handover before resuming. All current audit changes are documentation only. Runtime code, tests, dependencies and original source captures have not been modified.

## Progress journal

| Entry | Scope and outcome | Evidence and recovery |
| --- | --- | --- |
| SR-20261014-01 | Created full inventory and initial review coverage.0/41packages complete;19production and9test files received initial pass. | Commit cb9d18669e7f9fefb997b53062c2907a14911a49, parent a6fee427. Audit file fetched back and exact content matched. |
| SR-20261014-02 | User requested four sections, explicit stopping points, verification of uncertain/current facts, active correction/revert notes and interruption handover. Mapped41packages exactly once to9/13/15/4package sections; added progress and recovery record. | Documentation successor of cb9d18669e7f9fefb997b53062c2907a14911a49. Exact resulting commit is discoverable via branch history and commit message; read back files before crediting persistence. No code fixes or new test run claimed. |

## Findings and questions

| ID | Classification | Evidence/status | Next action |
| --- | --- | --- | --- |
| DOC-01 | Recovery documentation drift | Accepted-branch handover still said PR140 pending; freshly inspected PR140 is merged at a6fee427 and both accepted-push runs show success. Corrected historical full-verification record exists on OPEN/DRAFT PR141. | Current audit handover now distinguishes accepted code from unmerged notes and supersedes older status on this task branch. Still unmerged into protected branch. |
| Q-01 | Unresolved architecture question | Shared Azure client provides bounded and unbounded methods. Initial code pass alone does not establish unsafe endpoint use. | Section1 traces callers and response budgets; section2/3 resolve specific collector/plugin paths. |
| Q-02 | Unresolved compatibility question | Resource-ID helpers are positional; tests cover top-level ID and empty input. | Section1/2 trace nested/subscription-level IDs and pinned source intent before declaring defect or changing behavior. |
| Q-03 | Deferred platform validation | Report staging/DACL code and existing tests reviewed initially; fresh-OS/operator evidence still outstanding. | Section3 static/native assessment; separate fresh-OS validation when access exists. |

No newly confirmed production defect has been found or fixed in this audit checkpoint. Initial observations are not closed findings and do not establish defect-free code.

## Previously corrected findings, retained for backtracking

These were found in the earlier focused review, not the new systematic review:

| Finding | Before | Accepted correction and evidence |
| --- | --- | --- |
| FN124, P2 | ARG admitted one valid JSON value followed by extra buffered/streamed content. | Validate remaining bytes as whitespace through EOF.18bad suffix cases and4compiling controls; accepted through PR140. |
| FN126, P1 | Compute disk collector rejected documented final nextLink:null. | Accept omitted/empty/null terminators while retaining nonstring/unsafe destination rejection. Three pinned primary2021-07-01 examples and first/later-page cases; accepted through PR140. |
| FN127, P2 | Actual bounded-client error+non200response lost correct SKU HTTP classification. | Context terminal priority, then non200status, then remaining transport/read failure. Actual403/429/500/307 and partial-prefix pairs plus200error/context cases; accepted through PR140. |

Corrected candidate01e8de6244214204e4b34530c62e42c2c83221e3 merged as a6fee427e19feeef28691737ea9bc5026086efa0. Final native314Linux/309Windows named controls and source23records/315chunks were fully checked previously. PR141 preserves full evidence. A new session must reverify current refs; these historic counts are not new QA or audit coverage.

## Mandatory record for each future correction

For every substantive change, record:
- Stable finding ID, section/module/files, severity, intended behavior, observed failure and user impact.
- Baseline commit and affected blobs; minimal reproducible input and failing evidence before correction.
- Primary-source contract, pinned API/dependency version and access date where relevant; uncertainty and alternatives.
- Exact fix commit/parent, before/after behavior, rationale and affected interfaces; unrelated changes in separate commits.
- Regression expectations declared independently; exact QA commands, platform, run/job/head, result and retained failure logs.
- State: suspected, confirmed, fixed-unverified, verified, accepted or explicitly deferred. No implied acceptance from a passing candidate.
- Revert commit or precise commit(s) to revert, dependencies/ordering, schema/data effects and verification needed after rollback.

## Rollback and backtracking

Use the exact fix commit and parent, inspect its diff and subsequent dependencies, then create a normal revert commit on a task branch. Run relevant regression and required gates before protected acceptance. A merge revert needs explicit mainline/parent assessment and may remove multiple accepted changes; prefer a targeted corrective/revert commit where appropriate. Do not execute rollback merely because a plan is recorded.

The documentation-only initial checkpoint can be reverted by reverting cb9d18669e7f9fefb997b53062c2907a14911a49, subject to later documentation dependencies. SR-20261014-02 is the next commit with message "docs: split systemic audit into four sections and update recovery logs"; its SHA is in Git history. If undoing both, consider reversing later dependent documentation first. No data migration or Azure rollback exists for these notes.

## Verification ledger

At this checkpoint, independently re-read protected/audit refs, PR140 merged status, PR141 open/draft/head, both accepted-push run metadata and all documentation edited. Previously verified complete logs are linked as historic evidence, not rechecked full logs. Newly authored package mapping was checked for uniqueness and exhaustive coverage of the recursive, nontruncated41package inventory. Output file readback and modified-path checks are required before reporting save success.

Primary guidance opened2026-10-14: Go CodeReviewComments, Python PEP8, Microsoft PSScriptAnalyzer overview and GitHub Actions secure-use reference. Endpoint schemas, evolving dependencies/advisories and uncertain claims must be checked when relied upon; these generic references do not prove compliance.

## Interruption handling

Record last completed step and next dependency before leaving work. After transport/service failure, do not assume a mutation failed or succeeded: inspect branch ref, commit object and exact changed content. Retry only after reconciling outcome. Failed reads are not coverage. This repository journal is durable; tool-side stores/local scratch may disappear between sessions.
