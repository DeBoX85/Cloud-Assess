# Current handover authority — 2026-10-16

Read [handover checkpoint](HANDOVER_CHECKPOINT_20261016.md), [restart prompt](NEW_CHAT_PROMPT_20261016.md) and [Section 3 active report](SYSTEMATIC_REVIEW_SECTION3_20261015.md) first. This entry supersedes earlier current-state, pause and next-action paragraphs below; historical contracts, failures and evidence remain preserved.

The user requested a new-conversation handoff during Section 3. Audit Sections 1–2 are complete; Section 3 is incomplete, with 42/94 baseline Go files initially read and no Section 3 package finally verified. Audit counts remain **2/4 sections (50%), 22/41 package directories (53.7%)**; neither measures weighted overall project completion. Section 3 has 15 packages/94 baseline Go files (43 production,51 tests), plus one added SARIF regression. AI governance10Go files and region42Go files remain, alongside contract/data/architecture reconciliation across the initial reads.

Verified code/test head `b293e9533b06667d15b5e7efb6a20b39a9ea8427`, tree `71a1adb10c20ac23b0a55bd361c88b2204a2b048`, is draft/unmerged in [PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). Protected `bootstrap/core-v1` stays `a6fee427e19feeef28691737ea9bc5026086efa0`. **SR-F10 is reproduced and UNFIXED**: SARIF mixed-case result IDs do not match their published rule descriptors. Formatted regression-only native38829684903 fails the declared two assertions on BOTH Linux116745335663 and Windows116745335627 after compilation; Linux formatting passes and no race warning was observed. Source38829684879 succeeds separately. Red is reproduction, never acceptance. Earlier ed30ef4 Linux formatting failure is QA-S3-01 and is corrected in b293e95; do not repeat it.

NEXT: verify live refs/PR142 and actual files; implement the documented isolated SR-F10 correction without mutating canonical IDs or fingerprints, add compiling reversal/restoration control to both required jobs, complete ALL remaining Section 3 source/test/contract review, fix substantiated findings, obtain full corrected exact-head Linux/Windows/source QA, then publish the Section 3 checkpoint and STOP for re-evaluation before Section 4. Continuation of Section 3 is already authorized. Do not start Section 4, merge protected history, repeat Sections 1–2 or resume feature development as part of this handoff.

Remaining estimate: Section 3 **5–8 active hours**, Section 4 **2–4 active hours**, plus material fixes/repeat QA. This replaces the earlier optimistic6–11hour combined estimate; it is not a wall-clock completion promise. Overall project percentage/release ETA are unverified. After the FULL audit resume Storage contract/tests on draft [PR141](https://github.com/DeBoX85/Cloud-Assess/pull/141), then remaining adapters/integration. Source captures7/7, accepted target SKU adapters1/7, public region unavailable. FullCost/other parity, AR-02/load, live/fresh OS, module-only advisory disposition, human/license/security/Gate004/release remain open. No Azure access/write, fixture/RBAC change, dependency/pin/capture/schema update or production-scope substitution is authorized by this handoff.

This documentation successor changes no Go/Python/workflow/data. Its new checks are separate and initially pending; the known unfixed regression can keep native CI red. Inspect live PR142 for the actual documentation head/results. Earlier passing code evidence must not be transferred to a failing/changed head.

---

## Historical record below

# Section 3 authorized and in progress, 2026-10-15

The user authorized continuation at 05:04 Europe/Oslo. Complete the full Section 3 module review, independent reproduction/primary-source checks, isolated fixes and full corrective QA; stop for re-evaluation before Section 4. Earlier stop instructions below are historical.

Entry head `798d7be21e393f1e85503303e3d8eb83d311913d` is draft/unmerged in [PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). Its native38652456544 and source38652456593 now completed/success on that exact head (fresh run/job metadata; full log reconciliation follows). Protected baseline stays `a6fee427e19feeef28691737ea9bc5026086efa0`. Audit remains **2/4 complete (50% by section count), 22/41 package directories (53.7% by count)**; Section 3 is not complete. These are unweighted audit counts, not project completion.

Section 3 inventory: **15 package directories, 94 Go files (43 production/51 tests), 24439 lines including tests**. All local files independently Git-blob checked against the exact entry commit. Read the [active Section 3 report](SYSTEMATIC_REVIEW_SECTION3_20261015.md) for file-level coverage, findings, fixes, QA and rollback. Source pins/captures/dependencies/schema are unchanged; no Azure action or laptop access assumed.

NEXT: complete report/privacy/branding and plugin contracts, then small adapters, AI governance and region calculations/collectors, including all tests. Independently substantiate suspected defects before production changes, keep compatibility limits explicit, publish coherent code/notes checkpoints and read them back. Require full corrective Linux/Windows/source QA; remain unmerged and stop before Section 4.

After the FULL audit resume Storage contract/tests and remaining region adapters/integration. Source captures7/7, accepted target adapters1/7, public region unavailable. FullCost/other parity, AR-02/load, live/fresh OS, module-only advisory disposition, human/license/security/Gate004/release remain. Prior remaining-audit budget6–11active hours plus substantial fixes/repeat QA; weighted overall progress and release ETA unverified. Do not replay completed Sections 1–2 or begin feature work during this audit.

---

## Historical records below — superseded resume instructions

# Section 2 stopping checkpoint, review started 2026-10-14

Section 2 source/test review and corrective code QA are complete on `1c5d300ce2f726a36b0bc4e3eac09fc82588f337`, tree `a8147abfea11cbcd21151be447e51347c331a0bc`. This documentation-only successor records that verified code evidence and the user's planned re-evaluation stop. It does not inherit exact-head workflow acceptance. PR142 remains draft/unmerged; protected accepted baseline stays `a6fee427e19feeef28691737ea9bc5026086efa0`.

Coverage: all 13 Section 2 package directories, 51 baseline Go files (25 production/26 tests), six added tests inside the section and one coordinator integration test in Section 1. Whole audit: **2/4 sections complete (50% by section count)**, **22/41 package directories reviewed (53.7% by directory count)**. These are unweighted review counts, not overall project completion or defect-free certification. Current repository: 260 Go files and 48 Python files; the original inventory remains traceable.

Six baseline P2 finding groups (SR-F04–SR-F09) were reproduced after compilation on both hosts and corrected in isolated commits. Seven compiling/restoration controls and the final-Cost coordinator regression passed. Introduced QA-S2-01 (Cost local-name collision) failed candidate 7ce5on both hosts and was corrected in 4284; it was never accepted. Earlier automated no-major-issues feedback did not override that failure. Exact heads, parents, logs/hashes, fixture correction, source contracts and rollback dependencies are in the [Section 2 report](SYSTEMATIC_REVIEW_SECTION2_20261014.md) and [journal](SYSTEMATIC_REVIEW_CHANGELOG_20261014.md).

NEXT: stop for the user's re-evaluation. On authorized continuation, start Section 3's 15 plugin/reporting/platform packages; Section 4's four Go directories and all remaining tooling/delivery/data/provenance/docs follow. Do not replay completed fixes. Inspect live PR142/task/protected refs and any documentation-successor checks before a merge; no accepted-push proof exists for the unmerged audit. A new chat should read this current record and the current reports, then validate actual run heads/results.

Remaining audit estimate: **6–11 active review hours**, plus substantial findings and repeated QA. Overall weighted project progress and release ETA remain unverified. Source-compatible Arc SQL numeric decoding, core business-row/volume constraints and other residuals are explicitly recorded; review completion does not close them silently.

Project state remains: accepted region target adapters 1/7 (14.3% of that specific scope), source captures7/7, public region unavailable. Six adapters, coordinator/public region, FullCost/other parity, AR-02/load, live/fresh OS, module-only advisory disposition, human/license/security/Gate004 and release remain. After the entire audit, resume Storage contract/tests and remaining adapter/integration work; PR141 retains draft Storage analysis. Laptop/Azure restoration is unconfirmed. No Azure action or source pin/capture/dependency/schema/KQL change occurred.

Normal targeted revert commits must account for dependent helpers and tests, especially shared Diagnostics/Cost edits and the later HTTP variable correction. Do not reset protected history or rewrite retained source fixtures. Scratch and tool stores can disappear; the repository history is the recovery record.

## Verified code-head QA evidence

All evidence below targets `1c5d300ce2f726a36b0bc4e3eac09fc82588f337`, tree `a8147abfea11cbcd21151be447e51347c331a0bc`. The eventual documentation successor has separate checks and must not be treated as exact-head accepted by these results.

Native run [38649413582](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38649413582). Linux job 116144226683 completed successfully with 62 steps: 61 successful and the conditional failure-evidence step appropriately skipped. Race tests, 85.2% total statement coverage, bounded fuzzing, vet, builds/smoke tests, installation/package checks and all compiling controls passed. Its 325 control identities contain every previous 318 identity plus exactly seven Section 2 controls. Identity comparison removes only log timestamps and terminal CR; hashes preserve the complete returned UTF-8 content.

Source run [38649413551](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38649413551), job 116144226354, completed successfully. All 23 captures/315 chunks were reconstructed, index/count/duplicate/chunk hashes and Git blobs independently checked, and all 23 fresh candidate files matched. AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree `17d93b20c303f90f7843036be82f0dc32f3260f1`, and APRL `60eaddda76541f6adbc1c5ffa686829807e55e29` remain unchanged.

Windows job 116144227611 completed successfully, all 47 steps. Native Windows Go tests, CLI/preflight, PowerShell 5.1 checks, builds/smoke, branding, package/isolated-install and all compiling controls passed. Its 320 identities contain all previous 313 identities plus exactly the same seven Section 2 controls. Overall native run completed/success at the stated code head. Both fresh vulnerability scans report 0 reachable, 0 imported-package and 1 module-only advisory; that advisory remains open.

| Exact returned log | Job | Characters | UTF-8 bytes | SHA-256 |
| --- | ---: | ---: | ---: | --- |
| Linux | 116144226683 | 101775 | 101779 | `3c1af7d4e581c77bead3bcfc8acd0822b8067d88017056ca3647490c59f6f9e6` |
| Windows | 116144227611 | 82321 | 82323 | `ab0dc13c5f79a81d5fd7d7fdb465fbcc79a66f0429c5a7a7dadc3590d6b57aa8` |
| Source captures | 116144226354 | 1359878 | 1359954 | `b6c80e8f7169a0188f9c0b147929b14282cd0391a053135f9a9be32bdff0ca1f` |

This evidence closes corrective code QA for Section 2. It does not merge PR142, provide protected accepted-push proof, close release/live gates or accept an unverified documentation-successor head. CI triggered by this documentation checkpoint is separately pending at publication; inspect live PR142/checks for its later outcome. No code, source captures, dependency versions, query/schema or workflow settings change in this checkpoint.

---

## Historical records below — superseded progress and resume actions

The checkpoint above is current authority. Earlier in-progress/failed/pending snapshots below are retained for audit and rollback history, not instructions to replay completed work. The module/file inventory tables are reconciled to the completed checkpoints; baseline blob identities remain unchanged.

# Section 2 CI envelope correction in progress

This is current recovery authority. Code head4284ffeccc4bfd7c2bf7e257311b310937a29226 is still unaccepted: native38646061400 Linux116133085260 was canceled after1213seconds (22:13:01–22:33:14UTC), with timeout-minutes20 configured. Task head stayed4284and no newer workflow superseded it. This timing and GitHub's documented job timeout behavior support a CI-envelope timeout classification, not a product assertion failure. Check-run annotations are unavailable through the connector's allowed Fetch endpoints, so the cause is inferred from configuration, elapsed time and complete logs.

All325Linux control identities were present, including all318previous ones; race/coverage/fuzz/vet passed and govulncheck printed0reachable/0imported-package/1module-only advisory before cancellation. The canceled gate is not accepted. Returned full log102984characters/102988UTF-8 bytes, SHA-256b436bbeef0e848451675682bddca340063e37f3039bed6e713dc701cd9f86079. Windows116133085006 completed/success with all47steps. Complete returned log83777characters/83779UTF-8 bytes SHA-25612dd4081a0f592ef20d511ed01c9947e041fde9d64ea6b13da85ffbb2f62e5e4; all320control identities were present, with all313previous ones retained. Its complete scan reported0reachable/0imported-package/1module-only advisory. This successful job does not override the canceled Linux job or accept the overall run.

QA-S2-02 correction: increase only the bounded Linux job envelope from20to25minutes for the added controls/full-suite runtime. Retain every required test, control identity, coverage floor, individual Go/Python/fuzz timeout, pins, source checks, permissions and failure propagation. Windows remains30minutes. This changes CI execution allowance, not application timeouts or product/load acceptance criteria. Primary guidance:
https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#jobsjob_idtimeout-minutes

The earlier7ce5compile defect was corrected in4284; its HTTP and decoded Cost response locals are now distinct. Both hosts' seven new controls and ordinary full Go suites passed, and source38646061358/23captures/315chunks matched fresh retained bytes/pins; automated review4284reported no major issues. These partial results do not turn a canceled native run into acceptance.

NEXT: inspect complete native38649413582/source38649413551 and review on the already published candidate1c5d300ce2f726a36b0bc4e3eac09fc82588f337. The workflow-envelope change and preceding Windows outcome are already recorded. Continue with this candidate's actual-head verification, then save Section2checkpoint and stop for user re-evaluation. Section2is not complete yet; Section3has not started. PR142 stays draft/unmerged, protectedbaselinea6fee427unchanged, no Azure action or dependency/source-capture/pin change.

---

# Section 2 QA correction: Cost local response name, 2026-10-14 UTC

Superseded candidate7ce5cf2273edc6e95183ad3a8eabfa9fe6a11c38 FAILED native38645472872 on both hosts at the Section2control healthy Cost baseline. The new HTTP response local collided with the existing decoded queryResult local. This is QA-S2-01, an introduced compile defect, not an expected compiling mutation rejection or one of the baseline product findings. No passing acceptance is claimed.

Full returned Linux116131114362 log32401characters SHA-256e1b2405ac4e98a4345c38927190217527a3e718a1f5c3ab8a69e37ac9eb7f363; Windows116131113918 log32221characters SHA-256de5bc2f083ffe07634230daf35c3a7a80b6f25213fdd3e72833c57ca5bd53781. First two Section2controls passed/restored on both; later controls and full QA were not reached. Source38645472929 passed separately, with23captures/315chunks verified against retained baseline blobs; it does not clear native failure.

Correction: rename only the HTTP local to httpResponse, retain decoded response, and update the exact Cost-status mutation anchor. Fresh full candidate QA and exact-head review are required. Read live PR142/current refs for actual successor; do not transfer previous results. The Section2checkpoint remains pending and Section3has not started.

---

# Section 2 correction candidate, review started 2026-10-14

This is the current authority; older Section 1 stop/preparation paragraphs below are historical. The user authorized Section 2 continuation at 23:23 Europe/Berlin. Full manual review covers all 13 assigned packages and 51 baseline Go files (25 production, 26 tests), with independent full-content/Git-blob checks. Six P2 finding groups reproduced after compilation on both hosts and are corrected in separate commits. Full corrected-candidate QA and review are pending; do not credit Section 2 complete yet.

Read [Section 2 report](SYSTEMATIC_REVIEW_SECTION2_20261014.md), [change journal](SYSTEMATIC_REVIEW_CHANGELOG_20261014.md) and [live PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). Production correction ancestry ends at `cece5532c1ac79fd2f4964b1594bbafc06a0c61d`; current candidate refs may include subsequent helper/documentation or review fixes. PR142 stays draft/unmerged; protected baseline is `a6fee427e19feeef28691737ea9bc5026086efa0`.

Finding groups: SR-F04 terminal stage health; SR-F05 nested registry ownership; SR-F06 caller context in Cost/Diagnostics workers; SR-F07 one Diagnostics response/EOF; SR-F08 requested-batch setting scope; SR-F09 documented Cost 200/204 success. Seven compiling controls and one coordinator integration regression are retained. The first Cost status fixture failed before transport and is not status proof; corrected 10-second fixture/calls==1 reproduced all five invalid statuses on both hosts. A misplaced guard was caught in an unpublished draft and corrected before ref advance.

Reproduction-only `ac08820a66cdcbf55dfe39314f03b2138b11a48c`: native38642216627 failed only the declared Go tests on Linux116120286515 and Windows116120287313, after compilation/formatting/build; source38642216664 passed separately. Full logs/hashes, predecessor fixture classification, fix/parent and rollback dependencies are in the report. These failing runs are evidence, not acceptance.

NEXT: inspect complete latest-head native/source QA and automated review; address real findings, reconcile coverage and publish the Section 2 checkpoint. Do not replay completed work. Stop for the user's re-evaluation before Section 3. Audit remains 1/4 sections complete and 9/41 finally verified packages while candidate QA is pending.

Section 1 code00ba2716 and documentation7911764 passed both hosts/source; exact-current success is never inherited. Source captures/pins/dependencies/schemas remain unchanged; no Azure action, live parity or release approval. Source-compatible numeric Arc SQL and other business-row/volume limits are explicitly residual. One module-only advisory remained in prior scans; inspect fresh results rather than assuming current clearance.

Project state remains: accepted target region adapters1/7, source captures7/7, public region unavailable. Six adapters, coordinator/public region, FullCost/other parity, AR-02/load, live/fresh OS, advisory and human/license/security/Gate004/release work remain. After the entire audit, resume Storage contract/tests and remaining adapter/integration work; draft PR141 holds Storage analysis. Laptop/Azure restoration is unconfirmed. Sections 3–4 planning budget is 6–11 active review hours plus material fixes/repeat QA; no reliable overall weighted percentage or release ETA.

Recovery: freshly inspect task/protected refs, PR142 and actual run heads/results; read this current report and journal. Scratch/tool stores may disappear. Use normal targeted revert commits after dependency/QA assessment; never reset protected history or modify retained source captures to hide regressions.

---

# Section 2 authorized and in progress, 2026-10-14

The user authorized continuation at23:23Europe/Berlin. Section1 remains complete; Section2 review now covers13packages/51Go files. Read [Section2active report](SYSTEMATIC_REVIEW_SECTION2_20261014.md) for the latest coverage, findings and next step. Earlier Section1stop instructions below are historical and do not cancel this authorization.

Entry baseline is7911764ebd9fff9968589df8ee4c953a16a03d1c, whose native38574597736 and source38574597703 now completed/success. Both native jobs/steps freshly checked. PR142 remains draft/unmerged and protected code staysa6fee427e19feeef28691737ea9bc5026086efa0. Audit remains1/4sections complete; Section2not yet verified. Proceed with source/test review, independent reproduction and primary-source checks, then corrective QA and the Section2re-evaluation stop. No Azure action is authorized by this entry.

---

# Section 1 stopping checkpoint, 2026-10-14

This is the authoritative current state; earlier preparation paragraphs below are retained history. Section 1 code review and corrective QA are complete on code candidate `00ba2716f76366f6663de15531fdef4de6b812f8`, tree `e5098c01f631e0186d244631c61b4fad95cdda44`. This documentation-only successor does not claim exact-head QA inherited from its parent. PR142 remains draft and unmerged; protected baseline is still `a6fee427e19feeef28691737ea9bc5026086efa0`.

Coverage: all nine Section 1 packages, 95 baseline Go files plus two new regression files, reviewed and reconciled. Audit checkpoints completed: **1/4 (25% by sections)**. Package coverage: **9/41 (22.0% by directory count)**. Neither is a weighted overall-project percentage, a claim that all modules are reviewed, or defect-free certification.

Three confirmed Section 1 P2 defects corrected: terminal ARG executor context/error priority, nonunique inventory paging order, and undocumented ARG201–206 completed-query admission. Four compiling counterfactual/restoration controls are retained on Linux and Windows. First-review gaps in the correction and test oracle were reproduced/corrected; the current review on00ba2716 completed with no major issues reported. Historical repro runs/fix parents/hashes remain in the report and journal.

NEXT: stop for the user's planned re-evaluation. Do not start Section 2 automatically or replay completed fixes. Once continuation is authorized, review Section 2's 13 assessment/collector packages against the coverage map, including business-row completeness and request/memory budgets. Sections 2–4 need approximately **9–15 active review hours**, plus material findings and repeated QA. Overall project release ETA remains unverified.

Project scope remains: accepted region target adapters1/7, source captures7/7, public region unavailable; six adapters, coordinator/public region, FullCost/other parity, AR-02/load, live/fresh-OS, advisory disposition, Gate004/license/security/human and release checks remain. After the entire audit resume Storage contract/tests, then remaining adapters and integration. Storage analysis remains on draft PR141. Laptop/Azure restoration is unconfirmed; no Azure call or write occurred.

Recovery: freshly inspect PR142/task/protected refs; use the latest Section 1 report and change journal, retaining the distinction between verified code head and documentation successor. Follow documented normal revert commits with dependency review and appropriate QA; do not reset protected history or rewrite pinned source fixtures.

## Verified code-head evidence

All following evidence is for `00ba2716f76366f6663de15531fdef4de6b812f8`, not this documentation-only successor:
- Native [run38570549786](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38570549786), pull_request event, completed/success. Linux job115874677557:61steps, all success except conditional failure-evidence step appropriately skipped. Windows115874677774:46steps, all success.
- Complete returned Linux log:100982characters; SHA-256 `997e3543589bae24fec402b113caffa6fe6538e9d2a064e39ff95918d3b00dee`. Windows:82461characters/82463UTF-8 bytes; SHA-256 `fa29a05b2d4d79cd49708bbd00c6385de2041d875a8c8ac1e3015063dce24912`. Hashes retain original line endings; control identities alone normalize terminal CRLF.
- All prior control identities retained; four new compiling Go regressions/restorations passed. Linux318total (317Go-labeled controls plus one shell bootstrap control), Windows313Go-labeled controls.
- Race suite, formatting/module/dependency/provenance/docs checks, bounded fuzzing, vet, native CLI/branding/package checks passed; total Linux statement coverage85.2%. Both vulnerability scans:0reachable,0imported-package,1module-only advisory. Advisory disposition remains open.
- Source [run38570546281](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38570546281), same actual head/pull_request event, completed/success; job115874269869 all7steps success. Complete log1360512characters; SHA-256 `21914898e23c78a09179f39dbc1651537dbf175f676a0d623c6dcb743a0690cf`. Reconstructed315chunks/23captures with per-capture SHA-256; every capture matched fresh candidate file bytes, API Git-blob SHA and baseline retained blob.
- Source/AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree `17d93b20c303f90f7843036be82f0dc32f3260f1`, APRL `60eaddda76541f6adbc1c5ffa686829807e55e29` unchanged.
- Automated review result6169320382 and marker6168441214 completed2026-10-14T13:20:05.783287Z on00ba2716 with no major issues reported. This is not human or release approval.
- All17changed candidate files were freshly read back and their full authored contents/Git-blob hashes matched; final diff inspected. Protected branch and PRdraft/unmerged state independently re-read.
- This successor changes only eight Markdown records. It records completed parent-code QA and the user stopping point. Fresh successor checks must be inspected before any merge; no accepted-push QA exists for these corrections.

---

# Systematic module review, 2026-10-14

## Section 1 current work

All 95 baseline Go files in its nine packages received source/test review and full baseline content/blob checks. The candidate adds two Go test files and one Python regression helper; current totals are253Go and47Python. Baseline file inventory below remains754blobs/251Go/46Python for traceability.

Three confirmed findings and the automated review's corrective feedback are documented in [Section 1 report](SYSTEMATIC_REVIEW_SECTION1_20261014.md). Final native/source/review evidence and checkpoint reconciliation are pending. Audit completion remains0/4sections and0/41finally verified packages at candidate preparation. Previous initial-read counts are historical.

| Candidate delta | Purpose | State |
| --- | --- | --- |
| `internal/arg/executor_cancellation_test.go` | Terminal work, caller priority and ordinary-error control | Reviewed; code QA verified at00ba2716 |
| `internal/discovery/inventory_paging_order_test.go` | Value-based projection/sorting and model calibration | Reviewed; code QA verified at00ba2716 |
| `scripts/tests/section1-regressions.py` | Four compiling counterfactuals and exact-byte restoration | Reviewed/AST checked; both-host controls verified at00ba2716 |
| `.github/workflows/test.yml` | Execute these controls on both hosts | Relevant paths reviewed; both-host execution verified at00ba2716 |

No section outside Section1is claimed complete. Sections2/3still own collector/plugin budgets and business-row semantics; Section4reconciles all tooling/configuration/query/dependency/provenance/delivery coverage and this candidate delta.

## Status and baseline

**In progress, not a completed whole-codebase review.** Authorized by the user's request for a full systemic module-by-module assessment. The earlier three corrected defects are findings from a focused review; they are not a codebase-wide total.

Review baseline: `a6fee427e19feeef28691737ea9bc5026086efa0`, tree `1b13eec931b510592e42608fe8605ebdc287b0a9`, protected branch `bootstrap/core-v1`. The recursive inventory was not truncated. Findings introduced by later changes need separate delta review.

Inventory: 754 regular Git blobs; 112 production Go files, 139 Go test files, 46 Python files, four PowerShell files, two shell scripts and four CI workflows. There are 41 Go package directories, including developer tools. Queries, schemas, embedded data, configuration, dependency manifests, provenance, licenses and current operational/design documentation are also in scope. Historical status documents must be distinguished from current product requirements.

## Estimated effort

| Work | Active hours |
| --- | ---: |
| Module source and test review | 7–11 |
| Cross-module architecture, scripts, CI, dependencies and configuration | 2–4 |
| Integrated offline QA and targeted reproductions | 1–2 |
| Findings, coverage reconciliation and final assessment | 2–3 |
| Total review | **12–20** |

Planning allowance including ordinary remediation and repeat native QA: **2–4 working days**. Major redesign or substantial new findings require re-estimation. These are effort estimates, not a background execution schedule. Live Azure parity and fresh-OS/operator validation require appropriate access and are separate from offline review completion.

## Intent and acceptance criteria

Preserve a clean Azure assessment engine with semantic compatibility to pinned AZQR, Excel as primary human output and JSON as canonical machine output. Preserve read-oriented operations, Azure SDK authentication, explicit health/completeness, scalable assessment, controlled extensibility and attribution. Compatibility does not require reproducing AZQR's original JSON schema.

Assess correctness, cohesion and dependency direction, ownership and cancellation, bounded work, safe request scope and credentials, failure propagation, deterministic outputs, default identity redaction, platform-specific report privacy, API contracts, maintainability and meaningful tests. Review Go, Python and PowerShell practices in their respective contexts. Do not apply stylistic changes that alter deliberate semantic compatibility without evidence.

Review static architecture first, then trace executable paths from CLI/configuration through scope discovery, request execution, normalization, health aggregation and rendering. Validate failure paths as well as success, especially partial data, malformed responses, paging, throttling, cancellation and local file failures.

The existing accepted CI evidence remains historical baseline evidence; no new QA run is claimed by this checkpoint. Offline acceptance does not imply live estate coverage, fresh-OS readiness or release approval.

## Module coverage

A file read is not a completed module review. Initial source review means the listed source has been examined, while intent reconciliation, call-site tracing, test adequacy and fresh execution may remain open.

| Package directory | Production files | Test files | Review state |
| --- | ---: | ---: | --- |
| `cmd/cloud-assess` | 4 | 12 | Section 1 verified at 7911764 |
| `internal/advisor` | 2 | 3 | Section 2 verified at 1c5d300 |
| `internal/app` | 1 | 9 | Section 1 verified at 7911764 |
| `internal/arcsql` | 1 | 1 | Section 2 verified at 1c5d300 |
| `internal/arg` | 6 | 11 | Section 1 verified at 7911764 |
| `internal/assessment` | 8 | 0 | Section 1 verified at 7911764 |
| `internal/azure` | 8 | 7 | Section 1 verified at 7911764 |
| `internal/branding` | 2 | 1 | Pending |
| `internal/config` | 2 | 3 | Section 1 verified at 7911764 |
| `internal/cost` | 2 | 4 | Section 2 verified at 1c5d300 |
| `internal/defender` | 1 | 1 | Section 2 verified at 1c5d300 |
| `internal/diagnostics` | 2 | 1 | Section 2 verified at 1c5d300 |
| `internal/discovery` | 5 | 6 | Section 1 verified at 7911764 |
| `internal/equivalence` | 3 | 2 | Pending |
| `internal/findings` | 2 | 2 | Section 2 verified at 1c5d300 |
| `internal/gate` | 1 | 1 | Section 2 verified at 1c5d300 |
| `internal/orchestration` | 3 | 8 | Section 1 verified at 7911764 |
| `internal/plugins` | 4 | 5 | Pending |
| `internal/plugins/aigov` | 5 | 5 | Pending |
| `internal/plugins/carbon` | 2 | 2 | Pending |
| `internal/plugins/region` | 15 | 27 | Pending |
| `internal/plugins/servicehealth` | 2 | 1 | Pending |
| `internal/plugins/sqleol` | 2 | 1 | Pending |
| `internal/plugins/zone` | 1 | 1 | Pending |
| `internal/policy` | 1 | 1 | Section 2 verified at 1c5d300 |
| `internal/redact` | 1 | 1 | Pending |
| `internal/renderers/csv` | 1 | 1 | Pending |
| `internal/renderers/excel` | 1 | 1 | Pending |
| `internal/renderers/json` | 1 | 1 | Pending |
| `internal/renderers/sarif` | 1 | 1 | Pending |
| `internal/renderers/tables` | 2 | 1 | Pending |
| `internal/reportfile` | 3 | 2 | Initial source review; integration verification pending |
| `internal/result` | 3 | 3 | Section 2 verified at 1c5d300 |
| `internal/rules` | 5 | 4 | Section 2 verified at 1c5d300 |
| `internal/scanners` | 1 | 1 | Section 2 verified at 1c5d300 |
| `internal/skus` | 1 | 1 | Section 2 verified at 1c5d300 |
| `internal/stages` | 3 | 3 | Section 2 verified at 1c5d300 |
| `internal/throttling` | 1 | 1 | Source/test pass; correction QA pending |
| `tools/brand-build` | 1 | 0 | Pending |
| `tools/diagnostics-probe` | 1 | 1 | Pending |
| `tools/equivalence` | 1 | 2 | Pending |

## Initial observations — not confirmed new defects

- Assessment model declarations, severity handling and plugin-table validation have received an initial source pass. Model consumers and validator tests outside this package still need tracing.
- Azure source and all seven package test files have received an initial pass. Default HTTP redirects are disabled; ARM options disable automatic provider registration; operation timeouts and bounded methods have explicit tests. Caller-controlled transports are documented trusted code.
- The shared HTTP client exposes bounded and unbounded response methods. Trace all callers before assessing whether each endpoint's memory budget is appropriate. No defect is established merely by the existence of an unbounded method.
- Resource-ID helpers are positional and their current tests cover a top-level ID and empty input. Trace nested IDs, subscription-level resources and deliberately source-compatible use before changing behavior.
- Report staging source and both test files have received an initial pass. Existing-output preservation on render failure, Unix private permissions, symlink rejection and Windows DACL preservation are represented in tests. Fresh-OS validation, rename behavior and integration paths remain open.
- No additional defect has yet been confirmed or reproduced in this new systematic review. This is not a conclusion that the reviewed code or repository is defect-free.

## Four-section review plan and checkpoint policy

User instruction recorded 2026-10-14 at 13:38 Europe/Berlin. Audit takes priority over further Storage implementation. Review stops at each section checkpoint for a progress re-evaluation before advancing; this is a review boundary, not new permission for Azure actions or release.

| Section | Scope | Review effort | Natural stopping point |
| --- | --- | --- | --- |
| 1. Foundations, scope and execution | CLI/app/config, assessment contracts, Azure/ARG requests, discovery, throttling and orchestration; intent and architecture baseline | 3–5 active hours | Trace a scan from user inputs to validated scope and request execution; reconcile ownership, cancellation, bounded work and partial failure. Publish reviewed files, findings, fixes, verification and estimate. |
| 2. Assessment logic and collectors | Rules/scanners, findings, stages, results, severity gates and all core auxiliary datasets | 3–4 active hours | Reconcile data and health from collectors to canonical assessment, including malformed/empty/partial/failure paths and source semantic expectations. Publish checkpoint. |
| 3. Plugins, reporting and platform protections | Every plugin including internal-only region, all renderers, redaction, branding and report staging | 4–7 active hours | Trace plugin results through every output and platform boundary. Reconcile internal/public feature limits, deterministic behavior, privacy and source expectations. Publish checkpoint. |
| 4. Tooling, delivery and final integrated QA | Equivalence and developer tools, all Python/PowerShell/shell scripts, CI, dependency/license/provenance checks, configuration, queries and current docs; whole-system reconciliation | 2–4 active hours | Required offline QA and final inventory reconciliation; report closed findings, deferred risks, verified architecture/intent alignment and separate live/release gaps. Publish final assessment. |

Total estimated review effort stays **12–20 active hours**. With ordinary fixes and repeat native QA, allow **2–4 working days**. No measured active-hours completion is claimed. Stage budgets include checkpoint reporting and review QA; remediation depends on findings. Each checkpoint revises the estimate from actual work rather than treating elapsed calendar time as progress.

A section is complete only when its mapped source/tests have been examined, relevant call paths and intent checked, findings resolved or explicitly deferred with reasons, appropriate verification performed, and the checkpoint saved. Cross-section dependencies stay open with an owner section; the final audit cannot be called complete while they remain unassessed.

All 41 Go package directories map exactly once:

- Section 1: `cmd/cloud-assess`, `internal/app`, `internal/config`, `internal/assessment`, `internal/azure`, `internal/arg`, `internal/discovery`, `internal/throttling`, `internal/orchestration`.
- Section 2: `internal/rules`, `internal/scanners`, `internal/findings`, `internal/stages`, `internal/result`, `internal/gate`, `internal/advisor`, `internal/arcsql`, `internal/cost`, `internal/defender`, `internal/diagnostics`, `internal/policy`, `internal/skus`.
- Section 3: `internal/plugins`, `internal/plugins/zone`, `internal/plugins/servicehealth`, `internal/plugins/sqleol`, `internal/plugins/carbon`, `internal/plugins/aigov`, `internal/plugins/region`, `internal/renderers/csv`, `internal/renderers/excel`, `internal/renderers/json`, `internal/renderers/sarif`, `internal/renderers/tables`, `internal/reportfile`, `internal/redact`, `internal/branding`.
- Section 4: `internal/equivalence`, `tools/brand-build`, `tools/diagnostics-probe`, `tools/equivalence`.

Non-Go scripts/workflows belong to section 4. Queries/configuration/fixtures receive semantic review alongside their owning module, then completeness/provenance reconciliation in section 4. Architectural relationships are traced across all sections.

## Current progress and evidence discipline

- Accepted project baseline: PR140 merge `a6fee427e19feeef28691737ea9bc5026086efa0`; protected ref freshly rechecked.
- Core scan and multiple plugin migrations are implemented and accepted offline; release readiness remains incomplete.
- Region source SKU captures: 7/7. Accepted target SKU adapters: 1/7, Compute disks, or 14.3% of that specific adapter scope. Public region execution remains unavailable.
- New systematic audit: **0/4 sections complete; 0/41 packages complete**. Section 1 has started. Initial pass covers 19 production Go files and nine test files across assessment/Azure/reportfile; reportfile belongs to section 3. These counts describe initial reads, not fully verified coverage or an overall project percentage.
- No reliable weighted overall project completion percentage is available: the remaining features, live evidence, security/license review and release gates have different effort and no agreed weighting. Do not substitute adapter or audit percentages.
- Freshly rechecked PR140 merged status and accepted-push native `38502782779` and source `38502782765`: both event=push, head=a6fee427, completed/success. Earlier complete-log verification is documented separately in PR141; this checkpoint rechecked metadata, not a repeat of full log reconstruction.
- PR141 remains open/draft at `6cd002cb77fcead75a602b06a75625a178a84e49`. It contains detailed prior acceptance evidence and Storage analysis. Do not merge it implicitly, call it accepted, or credit it as production Storage code.
- Outstanding project work: six further SKU adapter types, coordinator/public region integration, FullCost/other migrations, AR-02/load, fresh-OS/live parity, dependency module-only advisory disposition, Gate004, license/security/human review and release. No release ETA is inferred from the audit estimate.

For uncertain or potentially outdated facts, record the question before relying on it; verify against current primary documentation, a version-pinned API specification or directly inspected dependency/source. Record source URL, version/commit, access date, exact claim and result. Reproduce behavior where applicable. Unresolved uncertainty stays explicitly open and cannot silently become a contract or fixture expectation.

Primary references rechecked on 2026-10-14:
- Go review guidance: https://go.dev/wiki/CodeReviewComments
- Python PEP 8: https://peps.python.org/pep-0008/
- Microsoft PSScriptAnalyzer: https://learn.microsoft.com/en-us/powershell/utility-modules/psscriptanalyzer/overview?view=ps-modules
- GitHub Actions security: https://docs.github.com/en/actions/reference/security/secure-use

These inform review; no lint/analyzer execution or architecture compliance is claimed merely from reading guidance. Preserve intentional compatibility and platform constraints instead of mechanically imposing style rules.

See [active change journal](SYSTEMATIC_REVIEW_CHANGELOG_20261014.md) for findings, correction evidence, commit ancestry and recovery/revert procedure.

## File-level initial source/test coverage

Only files actually examined in this initial pass are marked initial-pass. All other files below remain pending; fixture/data coverage needs semantic/provenance review rather than a claim that every captured JSON line is executable source.

| Path | Baseline blob | State |
| --- | --- | --- |
| `.gitattributes` | `3a487e37de77618987ec8f7752e2aa76a17edba6` | Pending |
| `.github/go-mod-tidy-trigger` | `21323c5a37b1bf0ef39a3dfa152600b3cc7b84fa` | Pending |
| `.github/pull_request_template.md` | `33f07fc5bcb2b0d08b3ce4c109b947b42518941a` | Pending |
| `.github/rule-sync-trigger` | `851d604358a9cff997656f9d4c07e8b73dbf0a20` | Pending |
| `.github/workflows/import-rule-baseline.yml` | `d8677b2da84ed72abae7dd6a66346f084e010204` | Pending |
| `.github/workflows/region-source-characterization.yml` | `cfddb4bfcf1fc2276288bafb925d48d4816c8867` | Pending |
| `.github/workflows/test.yml` | `4d4f02f5116e496a1d520a1afa067c596b817912` | Pending |
| `.github/workflows/tidy.yml` | `5d90829578fea27b72f273aace54b23e3aad7ae9` | Pending |
| `.gitignore` | `2244a330457ec2970946ee1e5884f9dac1209eea` | Pending |
| `.gitmodules` | `eba5947b59a32a2e2af085c1705ef5847b75e82d` | Pending |
| `AGENTS.md` | `cc4191adaaf1f4f4a18ba7021e696112ee0845fa` | Pending |
| `LICENSE` | `261eeb9e9f8b2b4b0d119366dda99c6fd7d35c64` | Pending |
| `NOTICE.md` | `9982f8f54cc3c15181c3c120801c2599afb4c7f4` | Pending |
| `README.md` | `e26277addbd403cf47c7a7f133270ddf2f66953a` | Intent read; reconciliation pending |
| `THIRD_PARTY_LICENSES.md` | `d8a2e12097e44a66d27683a118bfc1f30cdb5afe` | Pending |
| `cmd/cloud-assess/aigov_execution_test.go` | `26737a531244a72859ba49f38d88bd6a97ae3191` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/aigov_preflight_test.go` | `bd7768201aa4d0e4b6f50ad1c9919a6893f45a00` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/carbon_execution_test.go` | `334bfc568d3009b92f37bb76c1eea1db811d3226` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/command.go` | `529406d69d9add6dbaa94c665737d9f7a4816ddc` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/command_test.go` | `e2bd19e9dbaa38f50337101862e856263151da12` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/filter_schema_audit_test.go` | `981fa1e465b67d1a11c38c4efdde24035e5c709e` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/main.go` | `f63e26324925d093dc8178274e06fa745ffc6d14` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/plugins.go` | `d6eae94b62a49f9f863557235cbfa4535a77dbec` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/plugins_test.go` | `3c4ed15d27290135894bf27947d3d489985cf543` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/process_test.go` | `ece5916832edf0563507b9e21e3c4f95be50b1e0` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/rules.go` | `32896773fa8f52799347cc9f2067da0bc6f2d8fa` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/rules_test.go` | `f6abda3b4a2fb35a0ff1752e109ea8cd489df0ce` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/scanners_test.go` | `7ac24f32caef0b84a824e80025ddf6a1037d7758` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/servicehealth_execution_test.go` | `057bda19fe3511d456d53124a85d43d15b31bb06` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/sqleol_execution_test.go` | `18d6997d2a772d23f55499e81aadfc7f579c761e` | Section 1 reviewed; code QA verified at00ba2716 |
| `cmd/cloud-assess/testdata/rules-reference.json` | `992b347802179b23bb611440e77d5be59d052f01` | Pending |
| `cmd/cloud-assess/testdata/scanner-keys-reference.json` | `2719845c92c1098a4a22d4dc4fa6f390ec1198e5` | Pending |
| `cmd/cloud-assess/zone_execution_test.go` | `bfd21f46475553f00bd9c47d45a1f60c2affa53f` | Section 1 reviewed; code QA verified at00ba2716 |
| `docs/ACCESS_MODEL.md` | `7fae8708891ed906bf1dbfbcd7f122d2d5d5f042` | Pending |
| `docs/AI_GOVERNANCE.md` | `524934212cd531b8b7a4e57e47df4902f3897e74` | Pending |
| `docs/AI_GOVERNANCE_EXECUTION.md` | `686096f2a2bdcb633e1ff3b9ccfb43838ef98df5` | Pending |
| `docs/AI_GOVERNANCE_REQUESTS.md` | `72be43e1d062b663b1e51eb627ea6497d5c06125` | Pending |
| `docs/ALIGNMENT_REVIEW.md` | `a26595bf491c64d862e6111f69cdb14a8bc377fb` | Pending |
| `docs/ARG_COMPLETENESS.md` | `ddf493afe9fa76ebf223ce7847b9d952e9ee9e33` | Pending |
| `docs/AUDIT_ARG_RESPONSE_COMPLETION.md` | `ad0ed6133b9833093b5cbbe85f204f1c660930a4` | Pending |
| `docs/AUDIT_COMPARISON_INPUTS.md` | `a90eda1e2d4863605770f0615404c5fd0d4b46df` | Pending |
| `docs/AUDIT_COST_CONTINUATION.md` | `e6fbe4b1176937a826389e6e8b17f9feb0b1524e` | Pending |
| `docs/AUDIT_FILTER_SCHEMA.md` | `a2001608c9a9d9221845467ecb774485c35365b7` | Pending |
| `docs/AUDIT_RECOVERY.md` | `6172ed1317a13eec5146f3a6efa192e0350e34a6` | Pending |
| `docs/AUDIT_REVIEW_COVERAGE_20261004.md` | `34dedc5619215ea111dc47e9154f5bf7c21e6e9d` | Pending |
| `docs/BRANDED_PACKAGES.md` | `f4f0d513acb159a192f1b0abc25deb0ea196558f` | Pending |
| `docs/BRANDING_PROFILES.md` | `cd3d9b669501c4d51c97eebee84b98ba8948aa63` | Pending |
| `docs/BRANDING_REPORT_QA.md` | `2dfff2faec842fd50a517ba04b7b391c56da634c` | Pending |
| `docs/BRANDING_REVIEW.md` | `fca1b3baf7f7ae9304d4bcee973551a27d81a610` | Pending |
| `docs/BUILT_CLI_VALIDATION.md` | `cecdea86ca4223f9f292bcfdcdd648724257253a` | Pending |
| `docs/CARBON_EMISSIONS.md` | `5753528f76576e3b8b756231289533ce4f3f24e2` | Pending |
| `docs/CHARACTERIZATION.md` | `1c76c869e1b179ed7e0490c9996451953b9a278c` | Pending |
| `docs/CODE_REVIEW_QA_20261014.md` | `a0937fece243458f6d89e3be3c486531bb6df1ce` | Pending |
| `docs/DEFERRED_VALIDATION.md` | `f3c6d36af489a13cb0c6fdf15da1b55407ee438f` | Pending |
| `docs/DEPENDENCIES.md` | `26aade0b3c684dc257c08c0214b32ed82dec80a1` | Pending |
| `docs/DEPENDENCY_ADVISORY_REVIEW.md` | `9908b5e834608e05eb5957be13ca870564df057a` | Pending |
| `docs/DEVELOPMENT_EXECUTION_PLAN.md` | `7be46c09f11aa4cf898c3ec3862258b4422b39b9` | Pending |
| `docs/DEVELOPMENT_LEDGER.md` | `d5f073c59b5cfc74ee25b844abd76c593e3e1b7e` | Pending |
| `docs/DEVELOPMENT_WORKSPACE.md` | `dd304120c3b0623d55e7cf38f4d4c72801a9a1e8` | Pending |
| `docs/EQUIVALENCE.md` | `284488121547798b249c068cb6111d7bb61203f5` | Pending |
| `docs/FAILURE_NOTES.md` | `9d1b8dd52d700b678958c254723dd47da455cd4c` | Pending |
| `docs/FEATURE_PARITY_CHARACTERIZATION.md` | `a0e7ed76f73f2b260bcae1dae6abc864841eb186` | Pending |
| `docs/HANDOVER_CHECKPOINT_20261009.md` | `6712a9ffaad3b7b3c7b533bf71046195ab9898ae` | Pending |
| `docs/HANDOVER_CHECKPOINT_20261011.md` | `75d4587a93525b280d235cb9ca97b3c15a4b6351` | Pending |
| `docs/HANDOVER_CHECKPOINT_20261012.md` | `8b1ebd266a5ddeb3a8d29e5890fd0e33ff3667c4` | Pending |
| `docs/HANDOVER_SECURITY_GATE_20261009.md` | `433b4ba28497f694a3c5227850505666048fc4b3` | Pending |
| `docs/IMPLEMENTATION_PLAN.md` | `006bf2546d7f876ef0a8deac816d9fdd196f4d4c` | Pending |
| `docs/LAPTOP_LIVE_QUEUE_20261013.md` | `b4c95fc63f99ba4e71169e339262196efa445e4e` | Pending |
| `docs/MAINTENANCE.md` | `9a0abb3bb513c7a0deb3a0022f3595bcfd191372` | Pending |
| `docs/NEW_CHAT_PROMPT_20261009.md` | `83405e432cf3ad5fd02b4836f782cf222288649a` | Pending |
| `docs/NEW_CHAT_PROMPT_20261011.md` | `fec2422547b6349eebc38d3a4349cb763041e274` | Pending |
| `docs/NEW_CHAT_PROMPT_20261012.md` | `4c3886fac903b77e9c72dce93694f62819613835` | Pending |
| `docs/OFFLINE_QA_REVIEW.md` | `b64082f2e3dd746dacc040afc6ee1e1454e4bb7b` | Pending |
| `docs/OPERATIONS.md` | `75c17f1f0a049a58926d5e5f0b4d59942036d40b` | Pending |
| `docs/PACKAGE_BUILDS.md` | `e8dc7d1dec700ad3822b75be72aa39a5ccef7082` | Pending |
| `docs/PACKAGE_INSTALL.md` | `0cd9c20bee25023cb632613352292e2ddebd8e0d` | Pending |
| `docs/PAGINATION_LIFECYCLE.md` | `60efd2722356d139196bb0c504b59b0cc6982271` | Pending |
| `docs/PLUGIN_TABLES.md` | `c5d1d64dcfd020b8ca727584d483b240910b2620` | Pending |
| `docs/PROJECT_AUDIT_20261004.md` | `fd16192dccd447fee65bd98db3bb539b47c73553` | Pending |
| `docs/PROJECT_AUDIT_REPORT_20261004.md` | `f4f9479d434ff554940699e1068c62af6c9b6f49` | Pending |
| `docs/QA_PROCESS.md` | `4bac7fc8743b0b3403dc98c7ceb442b267d24d31` | Pending |
| `docs/QUALITY_GATE_001.md` | `abfcdd6eaed301987bb09b4eba09a8662657182e` | Pending |
| `docs/QUALITY_GATE_002.md` | `3dc62ee21c4161b56bdb35c4d972f40d0190e9e2` | Pending |
| `docs/QUALITY_GATE_003.md` | `bf2429ac2d40ef4dde89cc3c6912dc19c440dc04` | Pending |
| `docs/QUALITY_GATE_004_PLAN.md` | `aee13cfcc85673695e424b6c05914e88c1121d4c` | Pending |
| `docs/REGION_ADAPTER_INVENTORY_20261011.md` | `715766ac2b29e61ac1f450293f24bc9c330086ff` | Pending |
| `docs/REGION_AUXILIARY.md` | `18c2b9c24df4b00e4b5352411456590ad54d17b4` | Pending |
| `docs/REGION_AVAILABILITY_CALCULATIONS.md` | `73b1e863bafc136fafbae9c9ffdbd3df0089e71f` | Pending |
| `docs/REGION_AVAILABILITY_RUNTIME.md` | `84a345b440cf34d352325caa725e7435afbe282d` | Pending |
| `docs/REGION_COMPUTE_DISK_SKU_COLLECTOR_20261013.md` | `169dca2f8ddfb74060422a28c4d788ea98a967c8` | Pending |
| `docs/REGION_COMPUTE_SKU_SOURCE_20261013.md` | `a06e0004221ca036e8cd70adc93312b3871ff824` | Pending |
| `docs/REGION_COST_COMPARISON.md` | `2587226af9ee360444a0c516aedad7f7279696cd` | Pending |
| `docs/REGION_COST_ENRICHMENT.md` | `49a2f58f16f022d3500618b3be1ce2307e58b6ca` | Pending |
| `docs/REGION_COST_RUNTIME.md` | `f93fd5ee408680b93ac41028790bb6cec22b3153` | Pending |
| `docs/REGION_INVENTORY.md` | `f987bab1aaa4854ba32f27b0afaa4748fc5d2af0` | Pending |
| `docs/REGION_INVENTORY_AGGREGATION.md` | `b0b2ee5b3fccff1a36763f4ede06d074dd3a48ce` | Pending |
| `docs/REGION_INVENTORY_CALCULATIONS.md` | `c35560425125135655c59b393a9c24ac1ea79760` | Pending |
| `docs/REGION_LATENCY.md` | `31c47b35f29cff7dd5abcca5255442fa008d737b` | Pending |
| `docs/REGION_OTHER_SKU_SOURCE_20261013.md` | `70323b04cfa23cad1519c5ef2419bff77934c0e3` | Pending |
| `docs/REGION_PROVIDER_LOCATIONS.md` | `bde5a1a30761dce0051db14cc8d80c2ddfe7c406` | Pending |
| `docs/REGION_PROVIDER_TARGET_RECONCILIATION_20261012.md` | `960d48b125372be41770464f8378a8cecf8fded3` | Pending |
| `docs/REGION_QUOTA_COLLECTOR.md` | `feaffb7fdd1d53e3fbd91c047b751b3c583f1e2b` | Pending |
| `docs/REGION_QUOTA_RESERVATION.md` | `737d3f9257e238b08071d09600944466a8afa1e0` | Pending |
| `docs/REGION_QUOTA_RUNTIME.md` | `a739f47f19cc009aa9698da3573b3a6c685ac716` | Pending |
| `docs/REGION_RESERVATION_COLLECTOR.md` | `c04b455d7eb2e5eee8a84f8cd88305b639b7eaa6` | Pending |
| `docs/REGION_RESERVATION_RUNTIME.md` | `7e4e9c750c0f5794140a9940a531ec3cf093d31b` | Pending |
| `docs/REGION_RESERVATION_SOURCE.md` | `7754a4f97b60bf50b795b76b49417248c65b3c9f` | Pending |
| `docs/REGION_SELECTION.md` | `9ac6582492f7a1b92f898a1f20e27c9d8029df28` | Pending |
| `docs/REGION_SERVICE_AVAILABILITY.md` | `3828af8f46568568caa6e122519584717e4853ec` | Pending |
| `docs/REGION_SKU_TARGET_ANALYSIS_20261013.md` | `06a17ead247a0ac5f47bc8a6ee816a2a0eaa9312` | Pending |
| `docs/REGION_SOURCE_CHARACTERIZATION.md` | `347d89b8069200dc263b81e6f8ae9821c689106e` | Pending |
| `docs/REGION_VM_QUOTA_COLLECTOR.md` | `ca010b2b3e28c11f5607fd6c7a0bb379ea73790a` | Pending |
| `docs/REGION_VM_QUOTA_SOURCE.md` | `6634641ab60db8f36710fd0bad8f609e80184493` | Pending |
| `docs/ROADMAP.md` | `838f78156d78206317f6727ef596e32afdf35608` | Pending |
| `docs/RULES_INSPECTION.md` | `bd9b2b64e66d43758e58a9810beab5566d1862fb` | Pending |
| `docs/SCANNER_COMMANDS.md` | `4d0762f02c80a34e51b1b965ef6462f07ca41eb0` | Pending |
| `docs/SCOPE_REPORTING.md` | `5042bdf6d3ef2170beb64ed80e38988925882084` | Pending |
| `docs/SECURITY_ADVISORY_REVIEW_20261011.md` | `e91526ebb59f47a39924cb533d7046c7bbf592e1` | Pending |
| `docs/SECURITY_REMEDIATION_20261011.md` | `64877d95e3bde94a270dd15b5e4d507a51737677` | Pending |
| `docs/SERVICE_HEALTH.md` | `fd6483668f724e348148f494add049455d69b781` | Pending |
| `docs/SESSION_HANDOVER.md` | `e20a713cccba47914e7f370721118a759dba9f95` | Pending |
| `docs/SESSION_HANDOVER_HISTORY_20261003.md` | `421c8cb628c90b6991863185f004eb8005690f6f` | Pending |
| `docs/SQL_EOL.md` | `5aa5c8bad73c97f5be6a0f6369787e97707ed49e` | Pending |
| `docs/TARGET_SPECIFICATION.md` | `c37a63034c127b4f9345601d54cd883514b5776f` | Pending |
| `docs/YAML_GRAPH_PLUGINS.md` | `8f906ed55d0a7f71a0de0ba120cbbc89dc9ced84` | Pending |
| `docs/ZONE_EXECUTION.md` | `6ab5662bc863ceaa79652a35bdf3f9e541f66357` | Pending |
| `docs/ZONE_EXECUTION_CHECKPOINT.md` | `5cf16551c200ca19c607145698e53c8146b2da38` | Pending |
| `docs/ZONE_MAPPING.md` | `f60c08d242676bd6b3beb644b6996a81fde7e0c1` | Pending |
| `docs/dependencies/NOTICES.md` | `879f502c29e16a1c33fea5093ee8c490efe97abc` | Pending |
| `docs/dependencies/inventory.json` | `7739907a5c651e334ccb8a45219e34a1ed8adb72` | Pending |
| `examples/branding/example.json` | `9043c50ea35e7089fd1697eae62b322322b49d16` | Pending |
| `examples/filters/azqr-dev-rg.yml` | `eb229caec0e5a6102790428df2d3f69d54bd571f` | Pending |
| `examples/filters/azqr-environment-dev.yml` | `9d65a480529bdae7cdedb0822b98b66a0e72b41e` | Pending |
| `examples/filters/azqr-exclude-dev-rg.yml` | `365fa1f67b92c3bdb5c7431b5b9e98b2d4cc496b` | Pending |
| `examples/filters/azqr-exclude-environment-dev.yml` | `7a0a743879e95cc90633e1302e7773be6281ca27` | Pending |
| `examples/filters/azqr-exclude-vm-backup-recommendation.yml` | `5dc977dd46ace0494f38cf97ed113ec289207ae1` | Pending |
| `examples/filters/azqr-include-exclude-dev-rg.yml` | `d3ed4c6fb88015c57ba18c81674f9b44baae482e` | Pending |
| `examples/filters/azqr-storage-vm.yml` | `c57c5b11fd7e2b071e197a2c96340abb9711adb8` | Pending |
| `examples/filters/cloud-assess-dev-rg.yml` | `ccdc8198ef7fa29682d3ebed73f440ca35b9228b` | Pending |
| `examples/filters/cloud-assess-environment-dev.yml` | `e6c33f61e563de5249319cc99e40390fc6141b59` | Pending |
| `examples/filters/cloud-assess-exclude-dev-rg.yml` | `4f94fa02f623af24c46f42fb8e3b92948d59e44c` | Pending |
| `examples/filters/cloud-assess-exclude-environment-dev.yml` | `1b98b22324df3f833f95afda497da3926350e5bd` | Pending |
| `examples/filters/cloud-assess-exclude-vm-backup-recommendation.yml` | `eac3d0c1fcb8e1ad43be3319408709e25bed65ad` | Pending |
| `examples/filters/cloud-assess-include-exclude-dev-rg.yml` | `8bdaec233576a6a40904180464aaa48a3d9e6644` | Pending |
| `examples/filters/cloud-assess-storage-vm.yml` | `907a7c2f7f8746af07b49afc1739a581121fc8f2` | Pending |
| `go.mod` | `cb0d71e91bf9ce889c50a1ae5b5da9a962a4e502` | Pending |
| `go.sum` | `26a56ee20d2be4fa9ee61ea93af8da6788e400b2` | Pending |
| `internal/advisor/metadata.go` | `97dd140e5f0bf49b1377ed4c7d713cb6cfdf7be2` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/advisor/metadata_http_test.go` | `b8a81318ef8f86e7a5c2748b761a21f5fe62ba5d` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/advisor/pagination_context_test.go` | `e9afac187046b6addde7d5e6100a58f22694fd07` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/advisor/scanner.go` | `a636cd01e963469ac1a959da415bfc1c893f944c` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/advisor/scanner_test.go` | `b7a0d87604cc26613dfe71514bbd751e3dc5e2c2` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/app/arg_completeness_test.go` | `d18d1c9971b2e947d1360431de5de66accafd8d3` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/branding_report_test.go` | `a85b5d4806e27de8cde0f6750093b6e914eb94b8` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/e2e_test.go` | `c6f2dbb636fed2697cab6b716cf2eec762e32f93` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/optional_stage_http_test.go` | `1758ffd07d33ee7f17ab8c9ac21c7905bbc44bd2` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/plugin_tables_test.go` | `e6f021a83e2d6b62df5af65122baac1708bd7dff` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/scan.go` | `7a8400c35c19ca7c3e1ce75fcf843b73de4a155b` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/scan_test.go` | `5bfc08b37a18d7b1a4fe320e86735aca51f0866a` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/scope_test.go` | `e82fc2304cd70abecb91ba29ce28e092035a958d` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/testdata/README.md` | `f1e0fe410f7a6b4008c12356f6c31b0c113eae0b` | Pending |
| `internal/app/testdata/optional-stage-projections.json` | `074ee2a6918414bbb578ef945836fe6152e765fa` | Pending |
| `internal/app/timeout_test.go` | `20a50b7df1fa03506a66c6a2e1340b1363c46213` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/app/yaml_plugins_test.go` | `3b46b7aedae9a9b944d893416012c603d6fcec36` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arcsql/scanner.go` | `d4f9f0230c26842e02d62fbda9730132a050c59f` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/arcsql/scanner_test.go` | `4cd0d0130a7652adc66ba6c0881b66611524f5ad` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/arg/body_completion_test.go` | `a72f3191415ed7291785b1808b3fdef07ab16289` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/client.go` | `19c6d785b64b8dfd992e47975abdeefc5e03efa7` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/client_test.go` | `7de891afdf0262813347d5a4a289a9cb8a70371e` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/completeness_test.go` | `369578ea6871c1de732b2ab0b274f501e5378355` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/errors.go` | `7911871d473ec7773c9af0d8519ab27a69039c67` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/errors_test.go` | `4662fd2c280f266413a4cd0fa1f30c95c559fb41` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/executor.go` | `23422ce441240a46b60583b88ed30fb1e3b092e5` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/executor_test.go` | `0f0203f8b1993f127b0ee11bedddd4463b49ef9b` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/findings.go` | `68a9a7bbc7941bacd64966d9bf9cb0f2be778561` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/findings_test.go` | `16d076666f8565bd34c3d3151831e6c6cc1a6b50` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/http_transport.go` | `5928b3abda9c90b46d41f19aa4d49b432258af13` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/http_transport_suffix_test.go` | `420bf27921d2b12008926bb2fba09da679464ca1` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/http_transport_test.go` | `51920ccb610fab38c91069e85f62b17969409a7c` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/pagination_context_test.go` | `736346f36b0d998c02650c0dd2542010f69ebd36` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/rows.go` | `e432b57f7351879a92f1563cbfe09b3a8ebf6d07` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/rows_property_test.go` | `279d2ec7e4bde004761258012ea412882b418f72` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/arg/rows_test.go` | `c30ac9df118bb87c9479ef12585e0d8be3ffe296` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/advisor.go` | `600c865cf1e575c7eb1ea3b276425b4d4c40118c` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/arcsql.go` | `7805b19d055e20fea2efa08245e19d7b96af746d` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/cost.go` | `1ea0f29b29b434462725a9b5830500dbc8d802d2` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/defender.go` | `81342db7f4ccd410667b2eb67cdec527faa5db7f` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/plugin_tables.go` | `8582cd8d953ccb39726e67e8bd0b55ef7e94c179` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/policy.go` | `cf2c01227b6809b8286439c4557d48aa4ba12c5b` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/severity.go` | `eca6a403f225234868146b83a49535a54efd6037` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/assessment/types.go` | `5219513683b189cba4498eaff379f31fb67965ea` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/arm_boundary.go` | `e3606208cf4a1f1a80d99494b95eb9f1bd6a75cc` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/arm_boundary_test.go` | `9a1958596435f8e38eb5f2d46f5387b2879839d4` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/bounded_get.go` | `71413fe6db9a4bedd4715fecc59458b175923a9d` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/bounded_get_test.go` | `7679ce7fbbd0efb10be7b5561a393be8b2db5452` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/bounded_post_test.go` | `f111eab3f9982cf7939a58fec094cf1311d46123` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/client_options.go` | `af6ea4b0f6d3eba0bf0f5d47f4d0e8c1973d28d8` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/client_options_test.go` | `a40acbc5b66ef314348321132af8b5eaacc7f7ba` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/cloud.go` | `c5d9a2a080bc689180bfa15041c5636df13f2a44` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/cloud_test.go` | `d56a6172defc3af21166f0468a7fe6d44345a51d` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/credential.go` | `5afa22163e2bae2f03aa427ff5744eddc6b9023a` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/errors.go` | `d28cde817367348b86877961c0b76ab74afcd21d` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/http_client.go` | `a096e389ea6124faa4253f98186562d069b2b5b9` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/http_client_test.go` | `6367f9ae4eefdec7125212c09d51942c02571c4a` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/resourceid.go` | `1033f240db30820066a4e09af50097c8c29149e8` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/azure/resourceid_test.go` | `66a16170e9be8277dd83d94f9eb3c1a896246946` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/branding/branding.go` | `25146592055fd30be1a1f8abdabe529bc3ec48a1` | Pending |
| `internal/branding/profile.go` | `bf4ca310e6b376c421ff495e9cd58a7f5a4b8b53` | Pending |
| `internal/branding/profile_test.go` | `842b40ef97ec9f2c2a0cfbb1e3c7c147ecd6e434` | Pending |
| `internal/config/filter_property_test.go` | `a6d146f1216fa04825ddb9c9f1b40edb77e3e533` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/config/filters.go` | `896a08214bf42d63fa543a7777f59ce4e63dea89` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/config/filters_test.go` | `03fee8985cf11b60be583f3f06a89c3ba6d4c84c` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/config/load.go` | `46a03a41bf947dad526a3a828aaecb3687174dd8` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/config/load_test.go` | `f4ca80921d24affa2b368ba9ec739bbb7e11e0df` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/cost/continuation_test.go` | `2a2f9e4ce072da9d2ef7232886dd7e7515a4aaf5` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/cost/http_contract_test.go` | `2d93b9d1801eef11b9c9545c172e2c94c3e8149b` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/cost/period.go` | `ee7f36b63acef11043e2f4e42722499ccb658426` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/cost/period_test.go` | `192d0e94f67b3631d58e3372240b3b6fa0796b8f` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/cost/scanner.go` | `aaa7d4f85d576e7e76f2965513c23d53ea5f11ef` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/cost/scanner_test.go` | `f4e7a4b3d908c37c41cc67d2be78643f1dc0d0fc` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/defender/scanner.go` | `7e843f0e826f5c7a1b928eb22c4de525c6bd48f5` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/defender/scanner_test.go` | `a72483ee1c6c6f0f4aad402328e8c2d58cf89d9b` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/diagnostics/recommendations.go` | `9fb8f241d68502dd425ebb6e2f40f8414dcc5ba6` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/diagnostics/scanner.go` | `b9c3cb25c28b1282c06866344effe0590539b13e` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/diagnostics/scanner_test.go` | `7bfcb1ac7a56ad33a14ec79a91de08ff78951bd3` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/discovery/azure_scope.go` | `57efd23ef843fe916788d1b757842a1263dda2a0` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/azure_scope_test.go` | `c97ce7cd715b431b7be30ac39b00c98e57767b7c` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/management_groups.go` | `543c870651d9a4ef2a77b5879162a977c9fbb86a` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/management_groups_test.go` | `739abd7868270c71576597071ca544cb654bee9e` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/pager_cycle.go` | `a0cb960c6198507f7169abed61bcae6e0a9d81cc` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/pager_cycle_test.go` | `306fc3c24494d962b5b20428eb76580b4fdb1da9` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/pager_security_test.go` | `02c4be2faf09c4c3a575db33a57cba2e715c729c` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/resources.go` | `82f177ccf0d8232f64ca40690a131c5da9c31aca` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/resources_test.go` | `8a9a351608e3461868f389f7d3fffaf642f69945` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/subscriptions.go` | `c91a85d1fc43231acbb71cbbeae2099396b335c1` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/discovery/subscriptions_test.go` | `9a728c5df3099cc22571a6920c2d8afbc456c6e2` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/equivalence/compare.go` | `1d0b77de42d21c555c96bad6f473964d765fb544` | Pending |
| `internal/equivalence/compare_guard_test.go` | `e59fe464f75e415870fef0bcb6cc096c38ae55c0` | Pending |
| `internal/equivalence/projection.go` | `cb7f31799442cba5227e13096d7f4c9dcd7d95c2` | Pending |
| `internal/equivalence/projection_test.go` | `3915be9d5907fced529204bb5ecf18b78df7a012` | Pending |
| `internal/equivalence/types.go` | `f7be9b21a3aeabc654730ca4086e0b4e7a3f5ecb` | Pending |
| `internal/findings/applicability.go` | `366dde9e246f1a8e92fd58f975dec70652ea98b2` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/findings/applicability_test.go` | `068c908fa55bbcff01202294876dcf3d8198a0c3` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/findings/summary.go` | `e4fc5b01bd25c668a85c3b77584b1a1d3223823a` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/findings/summary_test.go` | `62cad5632c5402f649c1a4d88e7d7f203f584924` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/gate/gate.go` | `35e3d887aff375b3ca63b673dfd090641dc5c4b4` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/gate/gate_test.go` | `45b96ba6241d71a11f7437960f424e0e68b52d4d` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/orchestration/aigov_test.go` | `bc8b62daac1199850b224cfb8a1962745222a020` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/carbon_test.go` | `92f21c2937f4d9a5fd4e970f9c0370e644b4f2f3` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/coordinator.go` | `ffd707e55b58f486525aa5a884cebe29ee807c00` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/coordinator_test.go` | `8fd79a0a2bb2f2d6589f141e4a264d5d4a94b2cd` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/cost_completion_test.go` | `ba66c3d68abbce8faa0341a53602d58461ebe68b` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/operations.go` | `db99c5d23dbc08715280421b5c1ae237e237077b` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/scope.go` | `76c5849a372ff6ed1818428f734f3252c4e3827e` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/scope_test.go` | `c24febd19ede835a5a5bb66d8e23d2a53956a5fb` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/servicehealth_test.go` | `3ca35440ad643206fda80027e53ed4ec84e230f2` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/sqleol_test.go` | `11d8ba70a4d5415a82d929ebf2ca18ecc3b3da12` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/orchestration/zone_test.go` | `f79682ee7cf4a9d5c8a7a864b34410710990fd95` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/plugins/aigov/decoder.go` | `c31875f5576165af668ce95eaa2408767321a203` | Pending |
| `internal/plugins/aigov/decoder_test.go` | `76e9df63a13a056b7799877dbeaf4b0282c71b3c` | Pending |
| `internal/plugins/aigov/discovery.go` | `63e8d826cd8e4cb4fd51a6ea86a5c4cc40a86660` | Pending |
| `internal/plugins/aigov/discovery_test.go` | `9f32084ed4059b885f158ef28044403b1c206503` | Pending |
| `internal/plugins/aigov/execution.go` | `558d79cdb2526dcff5256741061737c3a2bc8f7b` | Pending |
| `internal/plugins/aigov/execution_test.go` | `e63932442c8fce7f51b4f9a88657d8ffd49c0bd6` | Pending |
| `internal/plugins/aigov/projection.go` | `90ceacb770826aff716c4d1ebc2e106ba0b6211c` | Pending |
| `internal/plugins/aigov/projection_test.go` | `20454b562bb0645b53157958e97f249e5de66acb` | Pending |
| `internal/plugins/aigov/scanner.go` | `c00e0ee00be57cf8bc7256322fbc62722f64a768` | Pending |
| `internal/plugins/aigov/scanner_test.go` | `3da3767d71fefb91d76576f712044c742bc929be` | Pending |
| `internal/plugins/aigov/testdata/capture_test.go.txt` | `4ceeb21fb1ab73b17719440753afb36fd0bac1dd` | Pending |
| `internal/plugins/aigov/testdata/deployments-input.json` | `2e071f26e617982d775d8a1bd3b508de7c53d24b` | Pending |
| `internal/plugins/aigov/testdata/metrics-input.json` | `b6c7407eb3364d385f42d39d291d6c3da4201152` | Pending |
| `internal/plugins/aigov/testdata/source-all.json` | `b018113ea8d9452659751c66fb1e6a5ba17f8a36` | Pending |
| `internal/plugins/aigov/testdata/source-deployment-denied.json` | `a93a7f65d04110b9df491125a088a3a43b04a28f` | Pending |
| `internal/plugins/aigov/testdata/source-empty.json` | `49196550c01e0dd866d9d73919bb165185bfec1f` | Pending |
| `internal/plugins/aigov/testdata/source-unknown-include-tag.json` | `49196550c01e0dd866d9d73919bb165185bfec1f` | Pending |
| `internal/plugins/aigov_projection.go` | `408b7cda42858faeab5b18b9fcba43cf8aa80c01` | Pending |
| `internal/plugins/aigov_projection_test.go` | `c016c61fe54c6d006a3b186374245a4018fd4289` | Pending |
| `internal/plugins/carbon/projection.go` | `868599a28286a303671cc8a76ecd06d90f13d1e0` | Pending |
| `internal/plugins/carbon/projection_test.go` | `105315a948c25c97be421b699e13eec17c3234e1` | Pending |
| `internal/plugins/carbon/scanner.go` | `d0aff9de47afad2c1a1b8c3fcd3ff9819ac5b1fd` | Pending |
| `internal/plugins/carbon/scanner_test.go` | `f8e1e5dca08bcc1706a5a48ff804bf45044d326e` | Pending |
| `internal/plugins/carbon/testdata/capture_test.go.txt` | `6c7d1bcdb70761c8f955e335d29f1a7f50d86225` | Pending |
| `internal/plugins/carbon/testdata/source-all.json` | `4a15cc42ae1b2a3feea8bc076ee7c3b9ec506df7` | Pending |
| `internal/plugins/carbon/testdata/source-batched.json` | `00c4f21c40937d67226fa1871ebf4f1fa1f023a1` | Pending |
| `internal/plugins/carbon/testdata/source-calculations.json` | `5464439fd1deeaabbf9bdc5434c0dab311847e89` | Pending |
| `internal/plugins/carbon/testdata/source-denied.json` | `5f2649fa1110a4206e1c0b7aa8db6464b0c99e7d` | Pending |
| `internal/plugins/carbon/testdata/source-empty.json` | `5f2649fa1110a4206e1c0b7aa8db6464b0c99e7d` | Pending |
| `internal/plugins/carbon/testdata/source-filtered.json` | `ed56b94a9d0a199b95def82693830f33ef998f64` | Pending |
| `internal/plugins/carbon_projection.go` | `626c3b24d2e2adb7dec5a53c0fef7f8d216efd89` | Pending |
| `internal/plugins/carbon_projection_test.go` | `a2b0221abaa4b345d67afd7f7877829beb2562d5` | Pending |
| `internal/plugins/region/auxiliary.go` | `15ab8323b5b4dd7974105eb15994d777f160ad35` | Pending |
| `internal/plugins/region/auxiliary_test.go` | `1526c4508c2a25a009305b9accf7eb5daa4c05b9` | Pending |
| `internal/plugins/region/availability.go` | `94360cfbe6a59a3debdf1a8f002cf7b062f8e264` | Pending |
| `internal/plugins/region/availability_capture_test.go` | `a212b8f605c28d8eac7137c893de3dc4f89ddaa5` | Pending |
| `internal/plugins/region/availability_test.go` | `695105d71c28c64980de19334f1a4370c3057869` | Pending |
| `internal/plugins/region/compute_disk_sku_collector.go` | `d7a5b8b2e1facd7a36b157acc9eadd4d91a57465` | Pending |
| `internal/plugins/region/compute_disk_sku_collector_test.go` | `76ce22e83bc9fe9fe4bdb51520782a0fed4cc97e` | Pending |
| `internal/plugins/region/compute_disk_sku_review_test.go` | `8789739cb653347be8b596e9fccb7e4836410212` | Pending |
| `internal/plugins/region/compute_sku_capture_test.go` | `2ef251506b2c6752a02685016aacee513bdf1978` | Pending |
| `internal/plugins/region/cost.go` | `dc7c43cdcfc46c2fa67e00aaafbc5cf0faeec149` | Pending |
| `internal/plugins/region/cost_enrichment.go` | `94a9cc17e8d70a02852274a69bfd6a6f5eb193d5` | Pending |
| `internal/plugins/region/cost_enrichment_capture_test.go` | `92cdcf06990aa8a5139c1966ee9b2ff9a62eec84` | Pending |
| `internal/plugins/region/cost_enrichment_test.go` | `c64816dd1daa8be6695b2d10076c267854955a7b` | Pending |
| `internal/plugins/region/cost_test.go` | `48712680a4eb91785d889a3c2796bfb3b6e623a3` | Pending |
| `internal/plugins/region/data/latency.json` | `363e872a014e0a4f92b5e43ecc9ee6cad3962bc3` | Pending |
| `internal/plugins/region/inventory.go` | `6157171769aaca1a4e78a026a4dc7f02a610bb87` | Pending |
| `internal/plugins/region/inventory_aggregation.go` | `9828cf3a30056e2456028b6e746e1e23afd5782c` | Pending |
| `internal/plugins/region/inventory_aggregation_test.go` | `b23fbf28124fdb72a78d5b86c85fff897e4cd30c` | Pending |
| `internal/plugins/region/inventory_capture_test.go` | `c7f85db283725e33b2d89a4df85dfdf6e8f8638e` | Pending |
| `internal/plugins/region/inventory_test.go` | `ac97771893275ad42fc5272a435e8a9dae57d718` | Pending |
| `internal/plugins/region/latency.go` | `1a4a037e8d6c1c80aed1d7d51c66a5d5d1ae1090` | Pending |
| `internal/plugins/region/latency_capture_test.go` | `877e215f1212f85096033acf9636618620479d13` | Pending |
| `internal/plugins/region/latency_test.go` | `60ea39a005cf1f1aff3f13c3889700706fc46dc2` | Pending |
| `internal/plugins/region/other_sku_capture_test.go` | `c73799952e61070813757283c13dc66079408ee7` | Pending |
| `internal/plugins/region/primary.go` | `3a7fb0b9601a525f91c385d0016ce0b8ff532303` | Pending |
| `internal/plugins/region/primary_test.go` | `6219219eca501b00af02a8af6fb23bb4a51c5d78` | Pending |
| `internal/plugins/region/provider_locations_capture_test.go` | `87112d5d3876fba242dbd2e0bb6ec9e24b35085f` | Pending |
| `internal/plugins/region/provider_locations_collector.go` | `8c85347d96488f2d4b1440da3f98a64ab3657413` | Pending |
| `internal/plugins/region/provider_locations_collector_test.go` | `9d1180be0e50e0f7fca5d30bbda00cb25eeb3ca6` | Pending |
| `internal/plugins/region/quota_capture_test.go` | `78a92620225e5e317e7af84dba73d6971828c326` | Pending |
| `internal/plugins/region/quota_collector.go` | `e2a44936371ec66600f33e2378b80ec689ee8b3f` | Pending |
| `internal/plugins/region/quota_collector_test.go` | `70a85e863018ef976e4f54a92318f3503d26b3f1` | Pending |
| `internal/plugins/region/quota_runtime.go` | `a2be4a680e2c3d705ed1198d1cbb509a04e44713` | Pending |
| `internal/plugins/region/quota_runtime_test.go` | `5ac53f03b1614a386b37fd92ed8071885b3efab2` | Pending |
| `internal/plugins/region/reservation_capture_test.go` | `31377c86e8b54423d799279c6d29ca3870cc1e72` | Pending |
| `internal/plugins/region/reservation_collector.go` | `c5b9c83e5d41be27798d7ab208f2568cca437c49` | Pending |
| `internal/plugins/region/reservation_collector_test.go` | `d617b44f6701b59234fbfb3866c08035a26fdf02` | Pending |
| `internal/plugins/region/reservation_runtime.go` | `bc6e23eca781e5b5ad81be2f416e2f9c36ec58c8` | Pending |
| `internal/plugins/region/reservation_runtime_test.go` | `1ec9bffe0cddf1e624c2b92be434ed8ee7952370` | Pending |
| `internal/plugins/region/service.go` | `e2e53165c2a4cb597c2bc36c239557189ad68ccd` | Pending |
| `internal/plugins/region/service_test.go` | `f683d274a5b208d4d54a7d1081de3f43aaba1581` | Pending |
| `internal/plugins/region/testdata/aux_capture_test.go.txt` | `e79411f6b378b76b5f676bf8c48e2bac549cb45a` | Pending |
| `internal/plugins/region/testdata/availability_capture_test.go.txt` | `25f3064df53ecfe557f4d825896b00ca93dd6e52` | Pending |
| `internal/plugins/region/testdata/capture_test.go.txt` | `ef8950b2abbaab17fcd4be24e0a68048ff56e611` | Pending |
| `internal/plugins/region/testdata/compute_sku_capture_test.go.txt` | `024dd96e96c48e12a7c76831ab5a58c9e56185f5` | Pending |
| `internal/plugins/region/testdata/cost_enrichment_capture_test.go.txt` | `1a865fcda1a149e884ae3ef543f34de6e9b90c49` | Pending |
| `internal/plugins/region/testdata/inventory_capture_test.go.txt` | `246736fdd46a2d071edcbdc3cd161df58059e5ef` | Pending |
| `internal/plugins/region/testdata/latency_capture_test.go.txt` | `bb1f843fa2371d6e1983b0626381a4e4faabd314` | Pending |
| `internal/plugins/region/testdata/other_sku_capture_test.go.txt` | `a9aedb09c9553eadc2e30bef7bd81a6821fe61ae` | Pending |
| `internal/plugins/region/testdata/provider_locations_capture_test.go.txt` | `755bc85e0a2c47860d3433cefeba716cae898113` | Pending |
| `internal/plugins/region/testdata/quota_capture_test.go.txt` | `a25aa04792e0cd873b156914e70cee9f9a76cda6` | Pending |
| `internal/plugins/region/testdata/reservation_capture_test.go.txt` | `5dfa4022ef513b6aea31271e08eaf5146e3e01fb` | Pending |
| `internal/plugins/region/testdata/source-aux-inputs.json` | `0eb15c3eee36d5bd402a988e46aa4ec9962a1970` | Pending |
| `internal/plugins/region/testdata/source-aux-outputs.json` | `7fbdf1bfa4e8656441a1c300bcfb936e0954067c` | Pending |
| `internal/plugins/region/testdata/source-availability-inputs.json` | `661bf5de672c6f3ddeda1783d4b4d75496f5eb47` | Pending |
| `internal/plugins/region/testdata/source-availability-outputs.json` | `d23e678ca8ecfac76841ef9370846f7d13eacd2e` | Pending |
| `internal/plugins/region/testdata/source-compute-sku-inputs.json` | `bf4e214a4493bfa6b612dd5181ae47db4ef4a598` | Pending |
| `internal/plugins/region/testdata/source-compute-sku-outputs.json` | `13f3492599bc3716cdf9abeca112f03d77ba89a7` | Pending |
| `internal/plugins/region/testdata/source-cost-enrichment-inputs.json` | `f629df479d264ef91fb4281755c124e0056e3932` | Pending |
| `internal/plugins/region/testdata/source-cost-enrichment-outputs.json` | `5e6368d2e4f0544184a48fe0713a42e97588174d` | Pending |
| `internal/plugins/region/testdata/source-inputs.json` | `4f1d725b7ec1eb6776af26ea0b7b11ddfb39ce27` | Pending |
| `internal/plugins/region/testdata/source-inventory-inputs.json` | `660b1dbf7f182179706247a233066e25791596ba` | Pending |
| `internal/plugins/region/testdata/source-inventory-outputs.json` | `67a78477a13f979278be4723892a300a82242dea` | Pending |
| `internal/plugins/region/testdata/source-latency-data.json` | `363e872a014e0a4f92b5e43ecc9ee6cad3962bc3` | Pending |
| `internal/plugins/region/testdata/source-latency-inputs.json` | `77d8bafebc9dbf210d481b287ed0e6ebdd170da1` | Pending |
| `internal/plugins/region/testdata/source-latency-outputs.json` | `946868f074f4717c2422fa0144fc7fc75a78c91f` | Pending |
| `internal/plugins/region/testdata/source-other-sku-inputs.json` | `07e5783b172a5871071d9cc55ff0224c6f807e0b` | Pending |
| `internal/plugins/region/testdata/source-other-sku-outputs.json` | `c49f80c973603b28965f70f3b93bd9ea20956452` | Pending |
| `internal/plugins/region/testdata/source-primary.json` | `ea5a30eae95518e06569b3bf466f4328b0d1193a` | Pending |
| `internal/plugins/region/testdata/source-provider-locations-inputs.json` | `3c9eb938921ea14f95d772b84a13c446b0da6bd1` | Pending |
| `internal/plugins/region/testdata/source-provider-locations-outputs.json` | `007a88c6b9ef2cc28c63aebf722a4693808fd97b` | Pending |
| `internal/plugins/region/testdata/source-quota-inputs.json` | `f4e8cd4ba70dcef555cc6441325e5eda23bde58c` | Pending |
| `internal/plugins/region/testdata/source-quota-outputs.json` | `06ff99b184458bff46f02d9a1df5e01585baf9e3` | Pending |
| `internal/plugins/region/testdata/source-reservation-inputs.json` | `8240a99b52515a955cd174b49188d6e55723b716` | Pending |
| `internal/plugins/region/testdata/source-reservation-outputs.json` | `6d3aa17eb688be5fff8a0c910080795c0b27024e` | Pending |
| `internal/plugins/region/testdata/source-vm-quota-inputs.json` | `9d593a0ac1fcca636a2b2c597fe52fbf87d275d6` | Pending |
| `internal/plugins/region/testdata/source-vm-quota-outputs.json` | `7c810647dfee30cd9028219efc67fc348934bd0f` | Pending |
| `internal/plugins/region/testdata/vm_quota_capture_test.go.txt` | `f90d337ce7fa5e7940e5167bd8860ac3993dcb53` | Pending |
| `internal/plugins/region/vm_quota_capture_test.go` | `18476fb14c4385d67e7d7fd16ccabbcc7ba8b52c` | Pending |
| `internal/plugins/region/vm_quota_collector_test.go` | `47546d5e6b3d79e9d0baa230f6739a3b75fbfa72` | Pending |
| `internal/plugins/registry.go` | `0e7e6810d481c5c7e806ee6469daf3c03f3f33c4` | Pending |
| `internal/plugins/registry_test.go` | `202b7a54f20738d90cf1f7650280ecd235c2a96b` | Pending |
| `internal/plugins/servicehealth/query.go` | `7624c0a3d23788f8c3e1565b2f35f55102b7e48f` | Pending |
| `internal/plugins/servicehealth/scanner.go` | `ffdc63d0b4669ad6ab23908deea1e69a4785ccd0` | Pending |
| `internal/plugins/servicehealth/scanner_test.go` | `1ef072fe74f1fd0024a1aaf56f38a162cead182c` | Pending |
| `internal/plugins/servicehealth/testdata/source-output.json` | `2b3f846ecada1657f18a267347f4146bfd362ebe` | Pending |
| `internal/plugins/servicehealth_projection_test.go` | `ff9b4db92f36855a592f79a7c23ecbf08fba5170` | Pending |
| `internal/plugins/sqleol/query.go` | `bb408e06215f9fb2d667dbc6372033b5e83d7858` | Pending |
| `internal/plugins/sqleol/query.kql` | `1857e4320d16c20b16dd368ee43d56a84ad5206f` | Pending |
| `internal/plugins/sqleol/scanner.go` | `2abf52dfc1ca3fda34e97dda898b7aa1d8631cc0` | Pending |
| `internal/plugins/sqleol/scanner_test.go` | `5ad2c6c14c6ecf23820ea989a8cfb8b62dc34ae1` | Pending |
| `internal/plugins/sqleol/testdata/source-empty.json` | `20f7fdf1a1c2c631d6162fe8de656a6846bc220e` | Pending |
| `internal/plugins/sqleol/testdata/source-filtered.json` | `73698041e621001011a08c62ace3f3e694e99498` | Pending |
| `internal/plugins/sqleol/testdata/source-input.json` | `cc6a0937194fa89f73d1681e8eae27d05cfc3c28` | Pending |
| `internal/plugins/sqleol/testdata/source-output.json` | `93adcc1de871cd99e9a2acef7f4c196009318164` | Pending |
| `internal/plugins/sqleol_projection.go` | `a925e0545eaa0e9b4be51a58705d46c0004f56fe` | Pending |
| `internal/plugins/sqleol_projection_test.go` | `cb8616b373409d017318f1d9f5bf20131a4492e0` | Pending |
| `internal/plugins/zone/scanner.go` | `1ba211bad8d2b0558b91481d177acb35417a568d` | Pending |
| `internal/plugins/zone/scanner_test.go` | `842fa2d7be23836ec7277bc5dcd2541f9e72a836` | Pending |
| `internal/plugins/zone/testdata/source-rows.json` | `f9812b012cbe5c704a85808a70ed802e5eac513a` | Pending |
| `internal/policy/scanner.go` | `1a549afefd7be50e0262336f733d8d5f76901a3d` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/policy/scanner_test.go` | `74d77a818ba3d0164d141b3ee1b0656503d7e53e` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/redact/subscription.go` | `74ca62b17ae958c03803be3b8d0550cd1a3c9852` | Pending |
| `internal/redact/subscription_test.go` | `9ce7ed9ac1ac9e8b8e25e260e4105c6b82ddd78c` | Pending |
| `internal/renderers/csv/csv.go` | `8ebe602a03cf70048c495d39881c1ef25242ddf0` | Pending |
| `internal/renderers/csv/csv_test.go` | `b2f6517219ecad4c61b098150499b337d860fb80` | Pending |
| `internal/renderers/excel/excel.go` | `198b681647d0037dd60ec7562aeb3b528cdae131` | Pending |
| `internal/renderers/excel/excel_test.go` | `70a0c6cf21a57b15da8edd92bd507c6209fa6f3f` | Pending |
| `internal/renderers/json/json.go` | `78c02db82280bc9a268ae87961e9ed070ae852ee` | Pending |
| `internal/renderers/json/json_test.go` | `08b95ca48d84313fbf7dcba50bb49f4aba104146` | Pending |
| `internal/renderers/sarif/sarif.go` | `137df562533618bf40fc3c7d65f962b77410d6c4` | Pending |
| `internal/renderers/sarif/sarif_test.go` | `06e983fd422dd55265f8941ecd5afff61e1903b5` | Pending |
| `internal/renderers/tables/plugins.go` | `c9c95c5016a9ce7354cabf0ee98e6ecfec38b10d` | Pending |
| `internal/renderers/tables/tables.go` | `cdec577a766f02a56c201b76b127954ae22d2318` | Pending |
| `internal/renderers/tables/tables_test.go` | `f7dd4bded84371943f6457ea831f20e53179444b` | Pending |
| `internal/reportfile/privacy_other.go` | `6712be13d06fc37bf2756daac4a08df91c2effb6` | Initial pass; module verification pending |
| `internal/reportfile/privacy_windows.go` | `00532ec043814f2c39d7d6b3458d801f4118ef91` | Initial pass; module verification pending |
| `internal/reportfile/windows_acl_test.go` | `8c96f8f89e9f48f433067bdef47f7bab45f610da` | Initial pass; module verification pending |
| `internal/reportfile/write.go` | `161ed19e385bbdfaab5c802cd7397d741918dab5` | Initial pass; module verification pending |
| `internal/reportfile/write_test.go` | `4ef1773e1f6953a4c4c520d1b0118e509ce9e84e` | Initial pass; module verification pending |
| `internal/result/plugin_tables.go` | `9703654451984758e91b8c962ae5626f31e79b96` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/result/plugin_tables_test.go` | `4a52a5a19bcc8658a2b8aa5b1832a3d89ca26a07` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/result/result.go` | `f099c293a44707535ffe624c201eddcf240847f0` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/result/result_test.go` | `bac7ad307b4c7f79d19831e1229d36293d9ef97e` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/result/scope.go` | `a65bc4f42d7c973cb399fd3a06c7f30019bfb7fa` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/result/scope_test.go` | `ed3ae54acdd9b33e90ab318836810fd3ac1523c4` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/catalog.go` | `e2d0a9ba82323ccb01dd837d14f7c25cf5d22c6d` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/catalog_test.go` | `e67a3d090bead57e968e2464c2de205b16494f12` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/custom/azure-resources/AAD/domainServices/kql/domain-003.kql` | `e12039b0a994213b4a6f8853f29f6f65d1a763c0` | Pending |
| `internal/rules/custom/azure-resources/AAD/domainServices/recommendations.yaml` | `0003baeba38685767edb681b531619db4e811331` | Pending |
| `internal/rules/custom/azure-resources/AVS/privateClouds/kql/avs-003.kql` | `543d4f8db8286ba960bac40b0ce64b4687adad30` | Pending |
| `internal/rules/custom/azure-resources/AVS/privateClouds/recommendations.yaml` | `f3fddddb2617845821bbe411f7b9eccdfd3cb190` | Pending |
| `internal/rules/custom/azure-resources/AnalysisServices/servers/kql/as-002.kql` | `da9a3af67b00d94480c1b425671157c503713190` | Pending |
| `internal/rules/custom/azure-resources/AnalysisServices/servers/recommendations.yaml` | `7a2638323db5a0f6e98c5c3f27134a0c92110a23` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/kql/apim-003.kql` | `bdc5b005821fc63c0abcea8f8046ea2d950a2a8d` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/kql/apim-004.kql` | `32efdde8773b5dc9b52ff5f03c038bdfa50eff70` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/kql/apim-008.kql` | `bb6e327b6cf25f3a3defcd8aecb316fab66d7b8f` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/kql/apim-009.kql` | `9dba63ba7c645d5bea44ef275e457d19bd66c3ff` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/kql/apim-010.kql` | `1b4e22917e575672c41aea9dfc845d9f0ae1ff59` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/kql/apim-011.kql` | `f22073a606e762058800a30223b70058e8b0beb0` | Pending |
| `internal/rules/custom/azure-resources/ApiManagement/service/recommendations.yaml` | `572dcb965f80ee9b2d86661cbc60b9bdaa1b66f6` | Pending |
| `internal/rules/custom/azure-resources/App/containerApps/kql/ca-003.kql` | `ae274fa5c7602deb3b25c0ff73387a2ff891d497` | Pending |
| `internal/rules/custom/azure-resources/App/containerApps/kql/ca-008.kql` | `42484f0417b618d657f691f3fc9077ed604ad57d` | Pending |
| `internal/rules/custom/azure-resources/App/containerApps/kql/ca-009.kql` | `36fed5044fa549464399f250557dc22e3f83340c` | Pending |
| `internal/rules/custom/azure-resources/App/containerApps/kql/ca-010.kql` | `dc10a9d6dedb0da30a484b723901dc362c4381c6` | Pending |
| `internal/rules/custom/azure-resources/App/containerApps/kql/ca-011.kql` | `330615d90906dae4a9d069f2b03a0d449aad3c28` | Pending |
| `internal/rules/custom/azure-resources/App/containerApps/recommendations.yaml` | `4f12bdb415dfb43808ceffaa3dbd3406461065d0` | Pending |
| `internal/rules/custom/azure-resources/App/managedenvironments/kql/cae-003.kql` | `6f72c72a13b940a80ec13ab2439103f9b0ae083c` | Pending |
| `internal/rules/custom/azure-resources/App/managedenvironments/kql/cae-004.kql` | `3983fe757891dc812297ec9964c676835a1605c6` | Pending |
| `internal/rules/custom/azure-resources/App/managedenvironments/recommendations.yaml` | `071aaa6f58e7f1a590202a4f666c91c38cbd3e7f` | Pending |
| `internal/rules/custom/azure-resources/AppConfiguration/configurationStores/kql/appcs-003.kql` | `f29b9e30b6d24578e806bcf1f00be9df8e8130e8` | Pending |
| `internal/rules/custom/azure-resources/AppConfiguration/configurationStores/kql/appcs-004.kql` | `2734c65c8c2dd8f9325f5fda7a36e3fba65b0615` | Pending |
| `internal/rules/custom/azure-resources/AppConfiguration/configurationStores/kql/appcs-008.kql` | `ea8fb672dfb62f0c6773c21ca7f732b2ab0f60e1` | Pending |
| `internal/rules/custom/azure-resources/AppConfiguration/configurationStores/recommendations.yaml` | `83db233163b5645eb217a8b926dab760de0b14cd` | Pending |
| `internal/rules/custom/azure-resources/Automation/automationAccounts/kql/aa-003.kql` | `10f8af780570d0b9c04eddc83869fca8ae7aa42f` | Pending |
| `internal/rules/custom/azure-resources/Automation/automationAccounts/recommendations.yaml` | `151bb3452a0c174af24a8e3d8ddacc99e861b469` | Pending |
| `internal/rules/custom/azure-resources/Batch/batchAccounts/kql/ba-003.kql` | `18100342351f9788fcedf7e46a459dfb6366967e` | Pending |
| `internal/rules/custom/azure-resources/Batch/batchAccounts/recommendations.yaml` | `b420cf970077b6efcc4809916affb41a38c25436` | Pending |
| `internal/rules/custom/azure-resources/Cache/Redis/kql/redis-003.kql` | `58653a572e705abcda965c5aee9b40c5b502843e` | Pending |
| `internal/rules/custom/azure-resources/Cache/Redis/kql/redis-008.kql` | `911782eec58df24ae53745de10e2c0955d85b35c` | Pending |
| `internal/rules/custom/azure-resources/Cache/Redis/kql/redis-009.kql` | `a22d02999f0c33c3fb3327a688c9b32ebc1a140a` | Pending |
| `internal/rules/custom/azure-resources/Cache/Redis/kql/redis-011.kql` | `17817238845a8a1308aed68570721fee9545c0c9` | Pending |
| `internal/rules/custom/azure-resources/Cache/Redis/recommendations.yaml` | `3eab97b1e78df571eaf24361b11cd914eebe0d14` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/kql/redis-010.kql` | `7f87c2870bfcb2c5414d96c4622171b0a8e47d6b` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/kql/redis-012.kql` | `201fae10ffef1fd0f5c12c8a7309b7eba9fe34d1` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/kql/redis-013.kql` | `81b86fd38c3f7777f3d1d570f68eb5b8406b73ad` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/kql/redis-014.kql` | `86b02acf14b5cbabf8e8ddb52c20225d6423f6c9` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/kql/redis-015.kql` | `5b5f659da4889b71732d6c287d7b0ebe1dd55cdf` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/kql/redis-016.kql` | `7cd7567aea80ca638a49b9552848ff0170038ff5` | Pending |
| `internal/rules/custom/azure-resources/Cache/redisEnterprise/recommendations.yaml` | `7081806c723631e179cb817340de0997fee5fa27` | Pending |
| `internal/rules/custom/azure-resources/Cdn/profiles/kql/afd-003.kql` | `f743b35e080e230ea2223cb9266e839a5c8064eb` | Pending |
| `internal/rules/custom/azure-resources/Cdn/profiles/recommendations.yaml` | `a225cca34f06f1206ebc505b2b7fddf54f42545b` | Pending |
| `internal/rules/custom/azure-resources/CognitiveServices/accounts/kql/aif-003.kql` | `5d0cd7935933b0b04ba7a0c4fda71e5211a6e2d0` | Pending |
| `internal/rules/custom/azure-resources/CognitiveServices/accounts/kql/aif-004.kql` | `9ccb69a01687cc1c530b2648aa39742cdc5c8075` | Pending |
| `internal/rules/custom/azure-resources/CognitiveServices/accounts/kql/aif-008.kql` | `c3a789eaa578cc5c1fa9bd5c4ed0264618f0d76f` | Pending |
| `internal/rules/custom/azure-resources/CognitiveServices/accounts/recommendations.yaml` | `1213617433ee5d7096ad497f7cdd84266795da95` | Pending |
| `internal/rules/custom/azure-resources/Compute/virtualMachineScaleSets/kql/vmss-003.kql` | `e4769f25ac1ab1f07a967862530f8b20ce027c9b` | Pending |
| `internal/rules/custom/azure-resources/Compute/virtualMachineScaleSets/recommendations.yaml` | `62417429d784ddaff8e0c0f375ee44f941eef86b` | Pending |
| `internal/rules/custom/azure-resources/Compute/virtualMachines/kql/vm-003.kql` | `b516fa1e7f5222b6a7fa2862170549c791d5bf99` | Pending |
| `internal/rules/custom/azure-resources/Compute/virtualMachines/recommendations.yaml` | `db1f08e82df8b973e980956db56a3655b9708fc3` | Pending |
| `internal/rules/custom/azure-resources/ContainerInstance/containerGroups/kql/ci-002.kql` | `99ac6c3801e101d469422e9781e0428007523590` | Pending |
| `internal/rules/custom/azure-resources/ContainerInstance/containerGroups/kql/ci-003.kql` | `b72d7e5afab942c3fab8851a21af1eadd6069cf9` | Pending |
| `internal/rules/custom/azure-resources/ContainerInstance/containerGroups/kql/ci-004.kql` | `841f1a768878e2d037a47e60b7182c21d0d6a856` | Pending |
| `internal/rules/custom/azure-resources/ContainerInstance/containerGroups/recommendations.yaml` | `2b7122afd3302fcdf12e0e92d43cf66ff65bd7e7` | Pending |
| `internal/rules/custom/azure-resources/ContainerRegistry/registries/kql/cr-003.kql` | `1d313e11ea838df4c58d448d088a7bab8cf3595e` | Pending |
| `internal/rules/custom/azure-resources/ContainerRegistry/registries/kql/cr-004.kql` | `044e47b3206038403bfed5eda169f46316b27a73` | Pending |
| `internal/rules/custom/azure-resources/ContainerRegistry/registries/kql/cr-008.kql` | `f89d83e914c8b668473eff92b317eb8cf5cdb142` | Pending |
| `internal/rules/custom/azure-resources/ContainerRegistry/registries/kql/cr-010.kql` | `ae3d871955f557f5e687f01ec0816bc2de1541cd` | Pending |
| `internal/rules/custom/azure-resources/ContainerRegistry/registries/recommendations.yaml` | `e3fc14422bfeb0a4012e159958c1f6d266224791` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/kql/aks-003.kql` | `412098003acddd65889d0b3bd659830022d546e0` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/kql/aks-004.kql` | `ba20f04c30d637329fda9eebde8fdceaa8232848` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/kql/aks-007.kql` | `4a275e4eb4e83c017e47172dea8a37afc7a1be28` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/kql/aks-010.kql` | `73273b3172b196eaba3654e04be6c179ead8e8f6` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/kql/aks-012.kql` | `89fe88302908bb5f95d17e9f3e271f5a78f1800c` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/kql/aks-016.kql` | `9993f780bf0863f57afa182d31b2238ab9a8eeea` | Pending |
| `internal/rules/custom/azure-resources/ContainerService/managedClusters/recommendations.yaml` | `5f5072c0e19c9a44a971652bb0fdd9ed1e1862d4` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/kql/mysql-003.kql` | `507117c2cfb7086401ea814908a6ea5786dbed74` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/kql/mysql-004.kql` | `bcd28f35f2c6403b134b5e98d58cd44fc1796767` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/kql/mysql-007.kql` | `2c830283eef358d824de14ba86215ad94b9b8d4d` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/kql/mysqlf-003.kql` | `bde06e5b50b4b9a61aece4f898b18cf7fa6bfab8` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/kql/mysqlf-004.kql` | `c39df628825924954ef7385d32702e859f68d520` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/kql/mysqlf-005.kql` | `16693b5cb27a0e74f86dddec6c909024d390e117` | Pending |
| `internal/rules/custom/azure-resources/DBforMySQL/servers/recommendations.yaml` | `6a389ff8fa6c57f6b24977022f6bfb8bc4dfe629` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/flexibleServers/kql/psqlf-003.kql` | `544dffcc363a65752f7c57974f0ebef9ba2d8cf3` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/flexibleServers/kql/psqlf-004.kql` | `a0176e585f0d469230e7deb92f400b7c26143bcb` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/flexibleServers/kql/psqlf-005.kql` | `8a7e905cd0988a827dd962b357e233e081c170c6` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/flexibleServers/recommendations.yaml` | `f93ade38ae5b8a685d06b2066d58da8354c39d9f` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/servers/kql/psql-003.kql` | `3cee839736753fe152f1c001e935662f31a23f97` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/servers/kql/psql-004.kql` | `9035912b37b02f384c6c1870e454a6423f33225f` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/servers/kql/psql-008.kql` | `cf99687a04e28e07752e9e1b7590f8f36f468f86` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/servers/kql/psql-009.kql` | `a748c1421445141b0b868fe6999680191713918b` | Pending |
| `internal/rules/custom/azure-resources/DBforPostgreSQL/servers/recommendations.yaml` | `4efbef4fb9ebda587ce4dc45c92a511c09e570ce` | Pending |
| `internal/rules/custom/azure-resources/Dashboard/grafana/kql/amg-002.kql` | `9f977fffec7d094e40c0e9dbe7132ba5253e5088` | Pending |
| `internal/rules/custom/azure-resources/Dashboard/grafana/kql/amg-004.kql` | `9b5d50b17626505bb1f1c4265ef4ed83aee397d0` | Pending |
| `internal/rules/custom/azure-resources/Dashboard/grafana/kql/amg-005.kql` | `ab3c04a37a45b7509af1039b832750f32da120eb` | Pending |
| `internal/rules/custom/azure-resources/Dashboard/grafana/recommendations.yaml` | `cffbf2d1308ab3db564b37b82fd5e5cadfeccd19` | Pending |
| `internal/rules/custom/azure-resources/DataFactory/factories/kql/adf-003.kql` | `85ff2ef0e2ee4bd8300572729f9c5afa7b4deefd` | Pending |
| `internal/rules/custom/azure-resources/DataFactory/factories/recommendations.yaml` | `fab61e892d5c013a94ce199bbac23e37f2b1f945` | Pending |
| `internal/rules/custom/azure-resources/Databricks/workspaces/kql/dbw-003.kql` | `90ae086a02328c67ace5b4159eabb12a3c382de1` | Pending |
| `internal/rules/custom/azure-resources/Databricks/workspaces/kql/dbw-004.kql` | `d581b3bf59eddaaa7cd7865781992df87320e13d` | Pending |
| `internal/rules/custom/azure-resources/Databricks/workspaces/kql/dbw-007.kql` | `fb22753ac0c22c7cf148b55e6a1dd5220c084d99` | Pending |
| `internal/rules/custom/azure-resources/Databricks/workspaces/recommendations.yaml` | `ff831f56582d5fc6235e617c3f9a989177ab903a` | Pending |
| `internal/rules/custom/azure-resources/Devices/IotHubs/kql/iot-003.kql` | `61d20dc3be920d6840920e00472cbe9da1d2357b` | Pending |
| `internal/rules/custom/azure-resources/Devices/IotHubs/recommendations.yaml` | `6c3803dc5a88687ce4444acb706e4791d2e420d9` | Pending |
| `internal/rules/custom/azure-resources/DocumentDB/databaseAccounts/kql/cosmos-003.kql` | `81d6bd21ff85bf23c3f01052aa3147fbbacdacf9` | Pending |
| `internal/rules/custom/azure-resources/DocumentDB/databaseAccounts/kql/cosmos-004.kql` | `c137bf01f5c89c669ee3543edb33fd9e5f7a6599` | Pending |
| `internal/rules/custom/azure-resources/DocumentDB/databaseAccounts/kql/cosmos-008.kql` | `e86c41426ea303275c77e6e31eacf305fc44b669` | Pending |
| `internal/rules/custom/azure-resources/DocumentDB/databaseAccounts/kql/cosmos-009.kql` | `301785d9b4ae533a99748d24df450f4afb47b411` | Pending |
| `internal/rules/custom/azure-resources/DocumentDB/databaseAccounts/recommendations.yaml` | `ee5a2a0e3958972fbbd5b596c44e729a1c815d80` | Pending |
| `internal/rules/custom/azure-resources/EventGrid/domains/kql/evgd-003.kql` | `d7d7c6e2b5a37ae31bd8be947d7163429ff60deb` | Pending |
| `internal/rules/custom/azure-resources/EventGrid/domains/kql/evgd-004.kql` | `dcda594d8d42173ffaab18ae14500da56d1bef58` | Pending |
| `internal/rules/custom/azure-resources/EventGrid/domains/kql/evgd-008.kql` | `15b091d5d1855da899046c9b867497134022ebaf` | Pending |
| `internal/rules/custom/azure-resources/EventGrid/domains/recommendations.yaml` | `f2bd00c7917cc9691aad4362fb4e4324b097752c` | Pending |
| `internal/rules/custom/azure-resources/EventGrid/topics/kql/evgt-003.kql` | `3c76259a933bf1d8bfa59ca2023750ab9df13580` | Pending |
| `internal/rules/custom/azure-resources/EventGrid/topics/recommendations.yaml` | `6d4c0554e0dea333cbd87dcf7061e41c8f1402db` | Pending |
| `internal/rules/custom/azure-resources/EventHub/namespaces/kql/evh-003.kql` | `398728cc08fdf469db37b44fcb271043526d159f` | Pending |
| `internal/rules/custom/azure-resources/EventHub/namespaces/kql/evh-004.kql` | `f7f1f156ec80695ec57db4c6ad5fd656f6f8a075` | Pending |
| `internal/rules/custom/azure-resources/EventHub/namespaces/kql/evh-008.kql` | `58de8a4462b0f000015f6580789cdf3b9e371d71` | Pending |
| `internal/rules/custom/azure-resources/EventHub/namespaces/recommendations.yaml` | `1e5ff63643f5eae4bab777db8c76dd999c3a3f8f` | Pending |
| `internal/rules/custom/azure-resources/Fabric/capacities/kql/fabric-003.kql` | `3979c05a132d6ab01d21fd65157b41d81df2bab3` | Pending |
| `internal/rules/custom/azure-resources/Fabric/capacities/kql/fabric-004.kql` | `b6a32353879504af9d38805908e9c904d8a8d95b` | Pending |
| `internal/rules/custom/azure-resources/Fabric/capacities/kql/fabric-005.kql` | `6bc8669c5abb96b21658c2cccdbc8e09304216e2` | Pending |
| `internal/rules/custom/azure-resources/Fabric/capacities/kql/fabric-006.kql` | `8686b59adcff6531b39a55384efa8eecd96fc192` | Pending |
| `internal/rules/custom/azure-resources/Fabric/capacities/recommendations.yaml` | `9fe516ed1364fe6cc6d4c63b4b42a1d42f54c4ae` | Pending |
| `internal/rules/custom/azure-resources/Insights/components/kql/appi-003.kql` | `d640034875a6b6f56dd2962f7733d30027c1f140` | Pending |
| `internal/rules/custom/azure-resources/Insights/components/recommendations.yaml` | `82e58083195c1918e894c5e234d8ed5cc53e920e` | Pending |
| `internal/rules/custom/azure-resources/KeyVault/vaults/kql/kv-003.kql` | `71296db30d34a8919336efe901de0d65861f175e` | Pending |
| `internal/rules/custom/azure-resources/KeyVault/vaults/recommendations.yaml` | `74ce3e6064e51fc348ceae007a1953fac3b1ef88` | Pending |
| `internal/rules/custom/azure-resources/Kusto/clusters/kql/dec-002.kql` | `4b72325facf77016670c101b2114a1c25f26b702` | Pending |
| `internal/rules/custom/azure-resources/Kusto/clusters/kql/dec-003.kql` | `060f35cef4dfadc153381f3f123e81f67b0641b5` | Pending |
| `internal/rules/custom/azure-resources/Kusto/clusters/kql/dec-004.kql` | `37099a360fcc98404f925e563d6e78aa9d2f2b7e` | Pending |
| `internal/rules/custom/azure-resources/Kusto/clusters/kql/dec-008.kql` | `8fd49f611a3b1598310beca4cf34aaf5bb03969e` | Pending |
| `internal/rules/custom/azure-resources/Kusto/clusters/kql/dec-009.kql` | `e9ec19b6c18a7cada11eb4c8acca58b5c77bd967` | Pending |
| `internal/rules/custom/azure-resources/Kusto/clusters/recommendations.yaml` | `ef23aba80fd68a4170b4fd1057df7175a0915530` | Pending |
| `internal/rules/custom/azure-resources/Logic/workflows/kql/logic-003.kql` | `a5a40692d072f31f3e5ac67c80d6be04ebbe9795` | Pending |
| `internal/rules/custom/azure-resources/Logic/workflows/kql/logic-004.kql` | `6ab9c85176ee5b57bd574af223ead41425b8f758` | Pending |
| `internal/rules/custom/azure-resources/Logic/workflows/recommendations.yaml` | `59157577b7725de08fcd4199b95814bdd02f727b` | Pending |
| `internal/rules/custom/azure-resources/MachineLearningServices/workspaces/kql/hub-003.kql` | `011258b4de86508072358e4601cfdcffae700d05` | Pending |
| `internal/rules/custom/azure-resources/MachineLearningServices/workspaces/kql/hub-004.kql` | `147bd6f853c262c8720cb2903580702cadc90f76` | Pending |
| `internal/rules/custom/azure-resources/MachineLearningServices/workspaces/kql/hub-005.kql` | `baec02d96dad7c1d9a60967c4e528d19ae1cc14d` | Pending |
| `internal/rules/custom/azure-resources/MachineLearningServices/workspaces/recommendations.yaml` | `ad4c7298ab79ec5eafc418ef1e2628992d076769` | Pending |
| `internal/rules/custom/azure-resources/NetApp/netAppAccounts/kql/netapp-003.kql` | `39e33321d202a7277eef3232ff333d265e830f11` | Pending |
| `internal/rules/custom/azure-resources/NetApp/netAppAccounts/recommendations.yaml` | `8a68263fe1def7a7db9ab9b9a38cd09cdc102415` | Pending |
| `internal/rules/custom/azure-resources/Network/applicationGateways/kql/agw-103.kql` | `bf28875e63b7ec9015a3d97a509f28cda4a8444b` | Pending |
| `internal/rules/custom/azure-resources/Network/applicationGateways/recommendations.yaml` | `2c8eb9587f9333797fd21f9df740336b5878ebe4` | Pending |
| `internal/rules/custom/azure-resources/Network/azureFirewalls/kql/afw-003.kql` | `95e462f76d4ccbf70165f667afd585881293a4e5` | Pending |
| `internal/rules/custom/azure-resources/Network/azureFirewalls/recommendations.yaml` | `638bfb4112b9ff12ddd9e3166f239bd76c897227` | Pending |
| `internal/rules/custom/azure-resources/Network/bastionHosts/kql/bastion-003.kql` | `62a9d5a96a6862dffb22c1a8f3b61417b6636462` | Pending |
| `internal/rules/custom/azure-resources/Network/bastionHosts/recommendations.yaml` | `e8706815843d88de0793f5b437da8760a08e08df` | Pending |
| `internal/rules/custom/azure-resources/Network/ddosProtectionPlans/kql/ddos-003.kql` | `7811edf0a9fd6d7bbc2dec5c223bc4acc9a4eae5` | Pending |
| `internal/rules/custom/azure-resources/Network/ddosProtectionPlans/recommendations.yaml` | `b91b8c7c251f69bef225cf473ad34b07e5a48628` | Pending |
| `internal/rules/custom/azure-resources/Network/dnsResolvers/kql/dnsres-003.kql` | `f2ed766eab5e1da8c03edcac5ec16e7432abb543` | Pending |
| `internal/rules/custom/azure-resources/Network/dnsResolvers/recommendations.yaml` | `6ccad6f184d0fe1fbb4d4935b57dab29c1dab5c3` | Pending |
| `internal/rules/custom/azure-resources/Network/dnsZones/kql/dnsz-003.kql` | `4844e56231917ed972041d691092a64ec4e0a978` | Pending |
| `internal/rules/custom/azure-resources/Network/dnsZones/recommendations.yaml` | `b079a0a3af1a1fa302bedc3f09369daae3c40942` | Pending |
| `internal/rules/custom/azure-resources/Network/expressRouteCircuits/kql/erc-003.kql` | `808dc203a278d5400dbd177b57307fc7d2e7ac1b` | Pending |
| `internal/rules/custom/azure-resources/Network/expressRouteCircuits/recommendations.yaml` | `5e49235da7b9b0e12c07b42e661191603d160e39` | Pending |
| `internal/rules/custom/azure-resources/Network/expressRouteGateways/kql/erg-003.kql` | `9e09a1f8f7baad208345625edbfd6562ea75afe5` | Pending |
| `internal/rules/custom/azure-resources/Network/expressRouteGateways/recommendations.yaml` | `1a90eb48ed7e781e7764c68c81f1d0521b6dbc49` | Pending |
| `internal/rules/custom/azure-resources/Network/loadBalancers/kql/lb-003.kql` | `e4976aed94d1d441ed0ac575f9889f96747489de` | Pending |
| `internal/rules/custom/azure-resources/Network/loadBalancers/recommendations.yaml` | `037bc28d381eb629c2d9030049ae857328a0fa55` | Pending |
| `internal/rules/custom/azure-resources/Network/natGateways/kql/ng-003.kql` | `ca54b9bb6e2fe5779375cec62e9d7e4da261d18c` | Pending |
| `internal/rules/custom/azure-resources/Network/natGateways/recommendations.yaml` | `7f7deebad049130a42d522bb4dc5e7c45a061403` | Pending |
| `internal/rules/custom/azure-resources/Network/networkWatchers/kql/nw-003.kql` | `3c717bb904cb8b05e71bbe0294cac954236ae2ac` | Pending |
| `internal/rules/custom/azure-resources/Network/networkWatchers/recommendations.yaml` | `155aeb7ffd345bc541fe28e72bce95eddd05d119` | Pending |
| `internal/rules/custom/azure-resources/Network/privateEndpoints/kql/pep-003.kql` | `40dad8c7c101b18d21be38ffa3b96dbdf288a02b` | Pending |
| `internal/rules/custom/azure-resources/Network/privateEndpoints/recommendations.yaml` | `253c75f5555cb74ec5ae73fa6886dfe9fbcf048b` | Pending |
| `internal/rules/custom/azure-resources/Network/trafficManagerProfiles/kql/traf-003.kql` | `d02fc96bad8e181a4c08e8cf47a74905f2a265b3` | Pending |
| `internal/rules/custom/azure-resources/Network/trafficManagerProfiles/kql/traf-009.kql` | `7d9b2895e7843a5f5e269ffc1e74bbbd6f8eed13` | Pending |
| `internal/rules/custom/azure-resources/Network/trafficManagerProfiles/recommendations.yaml` | `22b2fbd9ae437483f25a2c9a91e5957983e0706f` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualNetworkGateways/kql/vgw-004.kql` | `86f93d815fd9e1d3075a5a7794047db178b6472e` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualNetworkGateways/kql/vgw-005.kql` | `dff66ad7231131452cff7ecbdef674ffcc5b048c` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualNetworkGateways/recommendations.yaml` | `60a9c392f7d0bae91413188b5631f8698fc249c3` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualNetworks/kql/vnet-009.kql` | `6133975cf111bd6b37c4c32c11ca0f275a5884d7` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualNetworks/recommendations.yaml` | `4c4862ad15c7e34f4871d8e3562ad0254f808f54` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualWans/kql/vwa-003.kql` | `8669b5f5e2138f2623cc838375746e61656714ce` | Pending |
| `internal/rules/custom/azure-resources/Network/virtualWans/recommendations.yaml` | `8abd505f85df8d74e95223e0b9d63871d498e237` | Pending |
| `internal/rules/custom/azure-resources/NetworkFunction/azureTrafficCollectors/kql/ntc-003.kql` | `5ab434bfb863c51f95f13055c7249c280786ed64` | Pending |
| `internal/rules/custom/azure-resources/NetworkFunction/azureTrafficCollectors/recommendations.yaml` | `c45c3074211b23a061e7634411a8b99e23c41266` | Pending |
| `internal/rules/custom/azure-resources/OperationalInsights/workspaces/kql/log-003.kql` | `8bd0215952983ce70df509cf11d13c60522b16d7` | Pending |
| `internal/rules/custom/azure-resources/OperationalInsights/workspaces/recommendations.yaml` | `4c86348720d90bf770326724a1c8dad4bb65b890` | Pending |
| `internal/rules/custom/azure-resources/RecoveryServices/vaults/kql/rsv-003.kql` | `f0bee317c90e6e1cd2657d7dad4d1f37016d23d6` | Pending |
| `internal/rules/custom/azure-resources/RecoveryServices/vaults/recommendations.yaml` | `726639b793ff3288fa4e12be5be2987e7217f4a8` | Pending |
| `internal/rules/custom/azure-resources/Resources/kql/resources-001.kql` | `ce479a40a26683b4ce58fac6bcb09c0fe4507ff2` | Pending |
| `internal/rules/custom/azure-resources/Resources/kql/resources-002.kql` | `8c0082a3f78bf282b4b9002b1cfaeca7cae229e1` | Pending |
| `internal/rules/custom/azure-resources/Resources/recommendations.yaml` | `deea07f2575c98311cba07cda6c6ccf6b6bc2f2f` | Pending |
| `internal/rules/custom/azure-resources/Search/searchServices/kql/srch-002.kql` | `1794271e6901efc5c276e509bfa75fd114724ac0` | Pending |
| `internal/rules/custom/azure-resources/Search/searchServices/kql/srch-004.kql` | `5ba4b45703b5014e6fda5cbc32b7cd5b02e1c174` | Pending |
| `internal/rules/custom/azure-resources/Search/searchServices/kql/srch-005.kql` | `c1bd0693fe58fc97d472c7c541d1b4d11fcc19fa` | Pending |
| `internal/rules/custom/azure-resources/Search/searchServices/recommendations.yaml` | `b1a49215168783d19aa60072c39b4b18ea2b1106` | Pending |
| `internal/rules/custom/azure-resources/ServiceBus/namespaces/kql/sb-003.kql` | `4ccf3508c37f2777cde3c4914dcc9caf1560560d` | Pending |
| `internal/rules/custom/azure-resources/ServiceBus/namespaces/kql/sb-004.kql` | `42a974734ed1c79999534e4d64f0be82ce288bcd` | Pending |
| `internal/rules/custom/azure-resources/ServiceBus/namespaces/kql/sb-008.kql` | `f1c69cc1b4a1b04e533b122c37c70e71f5bfc601` | Pending |
| `internal/rules/custom/azure-resources/ServiceBus/namespaces/recommendations.yaml` | `09c997ce404517c396fd1c490d2d06d02ef9f2dc` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/SignalR/kql/sigr-003.kql` | `e7a1fb3fcf5b2c70679e2889431df6de8811f3fd` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/SignalR/kql/sigr-004.kql` | `3b8fe6839de6733c7e521c1eb90acc2450f6d957` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/SignalR/recommendations.yaml` | `4b9415a586d952ba86b9cf0fc0bd9a781cce0517` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/webPubSub/kql/wps-002.kql` | `b1caafc45eb35144c7d175efd049c47f4789f80f` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/webPubSub/kql/wps-003.kql` | `ce6574ff75f29cb2d573d362136ac3b9b9652058` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/webPubSub/kql/wps-004.kql` | `1f28ed4b3734d06c6d08530fbaa8d62ac441cbe6` | Pending |
| `internal/rules/custom/azure-resources/SignalRService/webPubSub/recommendations.yaml` | `60cf0ff2ae4ab0badbde7e1f9af4a3cc7d61ae55` | Pending |
| `internal/rules/custom/azure-resources/Sql/managedInstances/kql/sqlmi-003.kql` | `7cbc3dedff378fbff63e7d469125babdb269f9c5` | Pending |
| `internal/rules/custom/azure-resources/Sql/managedInstances/recommendations.yaml` | `98e99fb16fecb94e5064fd52625cda6c3abab33a` | Pending |
| `internal/rules/custom/azure-resources/Sql/servers/kql/sql-004.kql` | `ffc7bf96b500a9d9ac6cb8144fa45b7dfc87f76e` | Pending |
| `internal/rules/custom/azure-resources/Sql/servers/kql/sql-008.kql` | `87eb5e65e8ebe562f6503cb562e92ec2c03875e0` | Pending |
| `internal/rules/custom/azure-resources/Sql/servers/kql/sqldb-003.kql` | `c8ac5f0cb4041fa1ab095e9008e28d23d021bc0b` | Pending |
| `internal/rules/custom/azure-resources/Sql/servers/recommendations.yaml` | `4b7bd8eaea34d563c7d4bfadd08a0c4fbd437a5e` | Pending |
| `internal/rules/custom/azure-resources/Storage/storageAccounts/kql/st-003.kql` | `ebd37489302db389c483c2ca0a3293d8e087674a` | Pending |
| `internal/rules/custom/azure-resources/Storage/storageAccounts/kql/st-007.kql` | `140ff82edc5213f50b1bf16d5e415cd878c0e1ff` | Pending |
| `internal/rules/custom/azure-resources/Storage/storageAccounts/kql/st-009.kql` | `c496f7ffa9edf08d1f2359595c057edd8dae3871` | Pending |
| `internal/rules/custom/azure-resources/Storage/storageAccounts/kql/st-010.kql` | `02ecf1a569cc0734d05785ed77a95f29026a853c` | Pending |
| `internal/rules/custom/azure-resources/Storage/storageAccounts/recommendations.yaml` | `98421154aba9f7c277dd74383f7f1114ff6bea9a` | Pending |
| `internal/rules/custom/azure-resources/StreamAnalytics/streamingJobs/kql/asa-003.kql` | `7a42fcf40cecd89bf0d469619675bf0ae33c2d62` | Pending |
| `internal/rules/custom/azure-resources/StreamAnalytics/streamingJobs/recommendations.yaml` | `d2d3323201aac8d29a93f4da37070adee20ac07a` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/kql/syndp-002.kql` | `05cbe4df46693c144de34d193deafc8de25cb97b` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/kql/synsp-002.kql` | `b2aba70e53ab639382181afdec6550bcd95f71ed` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/kql/synw-002.kql` | `01c24b8006a3a4fe45bb8f927a3d2b2d45f6a763` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/kql/synw-003.kql` | `41f58d97b611a02ec105ce3e1aa5f478a4f62e42` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/kql/synw-006.kql` | `5354d3615f7c49453dfefda648c002e34d46682b` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/kql/synw-007.kql` | `df46a20f9df45141282ac99b9c15a518f8dd514c` | Pending |
| `internal/rules/custom/azure-resources/Synapse/workspaces/recommendations.yaml` | `8bed8d330b061e28fc2098ee078a05dbf0b41d36` | Pending |
| `internal/rules/custom/azure-resources/Web/serverFarms/kql/asp-003.kql` | `23dfd5753a7441fc2f1144c70f3c8f15eed2557f` | Pending |
| `internal/rules/custom/azure-resources/Web/serverFarms/recommendations.yaml` | `0c92d46752ce080dd3c42fe7babab7614e668fd8` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/app-003.kql` | `33e709b0858a5e4e5f4821fbd52e2742c3fca3e8` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/app-007.kql` | `db8c864860337eda32037ea0df1f47f79d2fbcfd` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/app-009.kql` | `5eac3e63ebb11e1b2b16071ee6648581ccb069fc` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/app-010.kql` | `d3910eeee3e544ea099cd86fc0ced3763ae6ef6a` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/app-015.kql` | `002f35bbf1a2a312cd0ddaca2705047b2d85304d` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/func-007.kql` | `fe80601893bf35fdd8fe484fb2f08ee8dbebb3a1` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/func-009.kql` | `3bb0dc9574001daf0c9b32b80a47706b59394be9` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/func-010.kql` | `fb54f190efc89fd78d22432e06050322828b2706` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/func-013.kql` | `44abad2e9e0fd4c778e2a496cb67fec20816ab2d` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/logics-007.kql` | `c7e94edb038fb62f61712dc7d0a70c6dcea4110d` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/logics-009.kql` | `6d292ba36e1c079b0c3377819b51299efdbc7ac9` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/logics-010.kql` | `f5b6ed4dd275e23b272b5a27e6ea16c2a2034697` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/kql/logics-013.kql` | `5a7cbef0b6e93bd18a42ac2632da334ac440533b` | Pending |
| `internal/rules/custom/azure-resources/Web/sites/recommendations.yaml` | `109ad04ed1e8f6b8127dcab01b27f6b7acda5730` | Pending |
| `internal/rules/embedded.go` | `5c3e0ce214e8d1571b9d1977b8f4a7f1d0dd326f` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/embedded_test.go` | `82611bcda9beb4d44c9a817f79379ad94d4591d4` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/inspection.go` | `97a998b23916134245f23fd15f0de1d4134897e7` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/inspection_test.go` | `8ff8d0177243a95bab178cabb8b4d4bebebe125f` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/provenance.go` | `0315d420dfc05b39cf03a6e782769d630d66ef8a` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/testdata/yaml-source-conversion.json` | `9300ce7cdfe2042eff327ad115155b552b5d9993` | Pending |
| `internal/rules/upstream/orphan-resources/Compute/kql/2b3c4d5e-6f7a-8b9c-0d1e-2f3a4b5c6d7e.kql` | `6da9ddf8a618384f0e4eca13e9aaff0a10020979` | Pending |
| `internal/rules/upstream/orphan-resources/Compute/kql/3c4d5e6f-7a8b-9c0d-1e2f-3a4b5c6d7e8f.kql` | `08d24b6a265942641352d2070f7d14837bcdf07e` | Pending |
| `internal/rules/upstream/orphan-resources/Compute/queries.yaml` | `23fdaf13541808518a96f53891fbaaf68d13ef7f` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/0b1c2d3e-4f5a-6b7c-8d9e-0f1a2b3c4d5e.kql` | `4215fffb60a32d57d459bb96bae055a325e5ae3e` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/0d1e2f3a-4b5c-6d7e-8f9a-0b1c2d3e4f5a.kql` | `aea3a1d04c698b438f6a586d0e361cbfd538d0e8` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/1e2f3a4b-5c6d-7e8f-9a0b-1c2d3e4f5a6b.kql` | `2ac0db732349dba54e2406ffd0cee9690cc6f4d0` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/2f3a4b5c-6d7e-8f9a-0b1c-2d3e4f5a6b7c.kql` | `74137eccaea750c45987af81c9346d6181e6f250` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/3a4b5c6d-7e8f-9a0b-1c2d-3e4f5a6b7c8d.kql` | `e50560fafb4fea6d16f4a11deee3cba243667f3c` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/4b5c6d7e-8f9a-0b1c-2d3e-4f5a6b7c8d9e.kql` | `8cb4802e57f8b33bc3814a9469f1fb99c7730d35` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/5c6d7e8f-9a0b-1c2d-3e4f-5a6b7c8d9e0f.kql` | `9c373c53e5a6fed9745acfaf2b174ee6f292b2d0` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/5e6f7a8b-9c0d-1e2f-3a4b-5c6d7e8f9a0b.kql` | `36990fac301c44bd9666a64ec56281d90c5a88f5` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/6d7e8f9a-0b1c-2d3e-4f5a-6b7c8d9e0f1a.kql` | `d2f6cd71c3d27bf321ef02c2a33e0f93851bb118` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/6f7a8b9c-0d1e-2f3a-4b5c-6d7e8f9a0b1c.kql` | `9f21e6081fc0294305998317a09407c7fc766576` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/7a8b9c0d-1e2f-3a4b-5c6d-7e8f9a0b1c2d.kql` | `bf8b66d78d62035c138f382d5d24b9b741d28960` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/7e8f9a0b-1c2d-3e4f-5a6b-7c8d9e0f1a2b.kql` | `009dd271b6c3ab10227a98ef34ffd2460cb21d7f` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/8b9c0d1e-2f3a-4b5c-6d7e-8f9a0b1c2d3e.kql` | `7bdbd0020a7779da5a131271b0e3e2be19dc03e9` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/8f9a0b1c-2d3e-4f5a-6b7c-8d9e0f1a2b3c.kql` | `5b22faa9f724954a6ef975f4d824c5e3488e95e1` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/9a0b1c2d-3e4f-5a6b-7c8d-9e0f1a2b3c4d.kql` | `ace9b0b38d3feb6befd92ef7099c34b0759b10bb` | Pending |
| `internal/rules/upstream/orphan-resources/Network/kql/9c0d1e2f-3a4b-5c6d-7e8f-9a0b1c2d3e4f.kql` | `c45af1f400a22d1693a1cd65d423b3ef4e4dc31b` | Pending |
| `internal/rules/upstream/orphan-resources/Network/queries.yaml` | `e24681041174b32e8925bd32ad66bb1eedc0dbe3` | Pending |
| `internal/rules/upstream/orphan-resources/Resources/kql/1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f.kql` | `992b8329d3c468a664c8a0399ae8ac73e077342f` | Pending |
| `internal/rules/upstream/orphan-resources/Resources/queries.yaml` | `32b7ec48ecfd73dbdb8f3fa440dd2ef5e8a9ca80` | Pending |
| `internal/rules/upstream/orphan-resources/Sql/kql/4d5e6f7a-8b9c-0d1e-2f3a-4b5c6d7e8f9a.kql` | `c7dbb4b59892ec76290453760666a65b2773cdd3` | Pending |
| `internal/rules/upstream/orphan-resources/Sql/queries.yaml` | `09e41a15a167cc3c5124d7c5de2916b2b064360c` | Pending |
| `internal/rules/upstream/orphan-resources/Web/kql/1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d.kql` | `dd311d4dfe4176c677289cbefdd1b2952066730a` | Pending |
| `internal/rules/upstream/orphan-resources/Web/kql/2d3e4f5a-6b7c-8d9e-0f1a-2b3c4d5e6f7a.kql` | `5b24c04785f2a7643e456422649e2a5fcb797233` | Pending |
| `internal/rules/upstream/orphan-resources/Web/kql/3e4f5a6b-7c8d-9e0f-1a2b-3c4d5e6f7a8b.kql` | `6943b999b51378e6a299205d8b5c3e6258f12426` | Pending |
| `internal/rules/upstream/orphan-resources/Web/queries.yaml` | `2b64eb02bef2594ac77cce67d804fbd14218167b` | Pending |
| `internal/rules/yaml_plugins.go` | `5a43f569f3ceceab829b9a26dd1d46c6732b5076` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/rules/yaml_plugins_test.go` | `d02f2967b009aa549c422017ce8292b3a851b9e8` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/scanners/registry.go` | `24512b10be30824c62f11bf83bc66e0b04e11eb6` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/scanners/registry_test.go` | `16f8f679331a50d1346c10292bdb72e73674fb71` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/skus/known_skus.yaml` | `a2d97a60ec445ce4023a1a5a9b6c4dca76e31f5a` | Pending |
| `internal/skus/skus.go` | `396cc4b7172f536b08d91355f8f64e9e4f0b97a1` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/skus/skus_test.go` | `6d1b6e6d89e4ac3a47c49fa54fbf34ec31db79dc` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/stages/config.go` | `307eabbb30b2b988cdc6e2fc750b3bce55d993d9` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/stages/config_test.go` | `804ccc644facc788b357fb36e54d06fe21c50e05` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/stages/options.go` | `d0138077d1784f5e82cde165af8eb62b410e55d2` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/stages/options_test.go` | `8245b9a14a468022ded93bdf62ab23538ecd6e16` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/stages/runner.go` | `7f3ec14d2f7b150619446984ddd9659b9761d8e4` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/stages/runner_test.go` | `6e6b6f6bbeb0db79a0b4d2c1a21cbd5ef46c3be7` | Section 2 reviewed; code QA verified at 1c5d300 |
| `internal/throttling/policy.go` | `0558c864c40152c4df848fe4c9c2a10fead10e0c` | Section 1 reviewed; code QA verified at00ba2716 |
| `internal/throttling/policy_test.go` | `753939b0001fd7709b3ec2fde99cf2e197e707e0` | Section 1 reviewed; code QA verified at00ba2716 |
| `scripts/bootstrap-workspace.sh` | `53be7d565640c8d873760f8e063163e231233f86` | Pending |
| `scripts/branding-profile.py` | `f97eb370157f18d3b36b355f479413ebbff445b9` | Pending |
| `scripts/capture-region-source.py` | `59c5a4915600dda8c786aa8b0643bd8af7365bfb` | Pending |
| `scripts/dependency-inventory.py` | `f76d92a0d36d5ee506a5f5a7972a0a5211ae7812` | Pending |
| `scripts/live-equivalence.helpers.psm1` | `bb2b1590e403375c963797daff462c8753f26a21` | Pending |
| `scripts/live-equivalence.ps1` | `5d4f64248659afc0afedca15098ae4ee012a6e72` | Pending |
| `scripts/package-candidate.py` | `be29991465ba8d847ef8477fd2ce3ed34858c905` | Pending |
| `scripts/prepare-resource-exclusion.ps1` | `74160968e8b6a913342d8b210ee03f0cec1d698a` | Pending |
| `scripts/publish-maintenance.sh` | `2ca62cc146154a88b67a37550a14882cd82ef462` | Pending |
| `scripts/tests/ai-execution-mutation.py` | `feb4a7e50eb5541c7fd11e20cba3a76ef51a39b4` | Pending |
| `scripts/tests/ai-request-mutation.py` | `ddd1895658aab8f1c4d5541824bfe74a812e4244` | Pending |
| `scripts/tests/arg-response-framing-mutation.py` | `06315565951f7fa07f760111e1335e45ffba5c08` | Pending |
| `scripts/tests/bootstrap-workspace.py` | `3a3650f9d5e45a61c4a104f495a9d2f4a2c88729` | Pending |
| `scripts/tests/branded-package.py` | `5b729447d473a83c77d40bd57c073bcdaeaa62c7` | Pending |
| `scripts/tests/branding-reports.py` | `6c204c8d348ebe19f8d9d431e82133709ec7db91` | Pending |
| `scripts/tests/branding.py` | `d17d968b181ea3d1c43180fc99880de2b9a22cf2` | Pending |
| `scripts/tests/built-cli.py` | `455eb6152187ce2d59e37d068b29d2a7216a9cbd` | Pending |
| `scripts/tests/compute-disk-sku-collector-mutation.py` | `4ad481d7d24dfd288533a7794c6de418a911f5ed` | Pending |
| `scripts/tests/dependency-inventory.py` | `b5e4e6b8bd178772c7a98e0facd6e3707ae00ffd` | Pending |
| `scripts/tests/documentation.py` | `73f737436e754ce32e788d742613fcc16b4668ba` | Pending |
| `scripts/tests/live-equivalence.tests.ps1` | `ac0e45d3f0285245c25de8a3ee2007098a42c2a8` | Pending |
| `scripts/tests/maintenance-generation.py` | `f32611354703065e13bdb13bd49944a634c26749` | Pending |
| `scripts/tests/maintenance-publication.py` | `2a82eddda9c35cda8c038ace21b53ae7f9b67691` | Pending |
| `scripts/tests/package-candidate.py` | `5f7aa9b00d644c4dc3f3b0259109ac36d0efddf2` | Pending |
| `scripts/tests/qa-mutation.py` | `d6e8ee02f9e47773be9e059b430b232eab8a586b` | Pending |
| `scripts/tests/qa-temp-cleanup.py` | `c44624027c09b377dd6701f529367729f4f9fa3d` | Pending |
| `scripts/tests/qa_temp_cleanup.py` | `7953e2ba3c2c60d6bdb06d62f6b5ddb771e7d3db` | Pending |
| `scripts/tests/region-auxiliary-mutation.py` | `3b46159ba695da177f5cdccf5670570f37e808ca` | Pending |
| `scripts/tests/region-availability-capture-mutation.py` | `1f08bbd4653c541f8ffa644e02839c25ff9ff145` | Pending |
| `scripts/tests/region-availability-mutation.py` | `a6301ee4ed107c84948e673c958a705deae597fe` | Pending |
| `scripts/tests/region-compute-sku-capture-mutation.py` | `f748aa76ac01c1f4d9f31b4acfb83344063b1531` | Pending |
| `scripts/tests/region-cost-enrichment-capture-mutation.py` | `32c97a20f0fb9a7da32c7c12b5e20ea40e2c7d54` | Pending |
| `scripts/tests/region-cost-enrichment-mutation.py` | `e7c98a0134cf6a3c87bac3136f035f1d6dc6f2cc` | Pending |
| `scripts/tests/region-cost-mutation.py` | `7cd73153c307aee9b044af60633e03f5207b9df7` | Pending |
| `scripts/tests/region-inventory-aggregation-mutation.py` | `e68828ae080cac91003973138ecc8533c8c7e2af` | Pending |
| `scripts/tests/region-inventory-mutation.py` | `8bb7956a8ad5ad4c17eb3be67e399b608b93d69d` | Pending |
| `scripts/tests/region-latency-capture-mutation.py` | `ea325dc97d26fdaa47c2d8eee28b232409bfcf03` | Pending |
| `scripts/tests/region-latency-mutation.py` | `7eccf70f359803e6360d0cf6e21299773176342a` | Pending |
| `scripts/tests/region-other-sku-capture-mutation.py` | `f01d6dda19fd2e6ea9b4d5e2c21f3d25a29d4d49` | Pending |
| `scripts/tests/region-primary-mutation.py` | `e4ed5303df304bcac0e91f25bcd7eb67d75e13b1` | Pending |
| `scripts/tests/region-provider-locations-capture-mutation.py` | `145bea96f8fe587dc38e296bc67b33eb81cd234b` | Pending |
| `scripts/tests/region-provider-locations-collector-mutation.py` | `8d2111616f8bcbb330c48cd707d1d9591d188452` | Pending |
| `scripts/tests/region-quota-capture-mutation.py` | `a0a014e37edd8945a1dcf27761587d363b1bdd44` | Pending |
| `scripts/tests/region-quota-collector-mutation.py` | `7b9fb20c8e922806facdeffc87b4ff6014f637f3` | Pending |
| `scripts/tests/region-quota-runtime-mutation.py` | `a48202c1d82eafc5bc82d7b4332a62333bee019f` | Pending |
| `scripts/tests/region-reservation-capture-mutation.py` | `4f942f1af28038e9d670c7eb074a33d4b6ff1a31` | Pending |
| `scripts/tests/region-reservation-collector-mutation.py` | `7323d0edbc61ffd8137d63ffc545c7907107e768` | Pending |
| `scripts/tests/region-reservation-runtime-mutation.py` | `fbdcb37335b4627f9db4c62e664082c9d98dfb21` | Pending |
| `scripts/tests/region-service-mutation.py` | `488cd60ebb9198f942605e86b604fbe620dcee25` | Pending |
| `scripts/tests/region-source-capture.py` | `e6a4873624db9f5fbfb63bc360c93be0e4d14af4` | Pending |
| `scripts/tests/region-vm-quota-capture-mutation.py` | `13fa511c57da0d64c82ab26a12291e1346e1772d` | Pending |
| `scripts/tests/workflow-command-failures.py` | `4aefe8afa6702b44f6183ba59bb9f84af250a273` | Pending |
| `tools/brand-build/main.go` | `875d824c50c48319faab2106c0d7ddacaaadb98c` | Pending |
| `tools/diagnostics-probe/main.go` | `ea2b806be70686e6524242f021f8df5751c04f5b` | Pending |
| `tools/diagnostics-probe/main_test.go` | `abeef1ac10f7083f85ee27f34b4354ebce3dcbf9` | Pending |
| `tools/equivalence/input_evidence_test.go` | `3d0a89b227cb523f2ffeb9b169f4d196bb894e95` | Pending |
| `tools/equivalence/main.go` | `66b8f529c8893ec87ac4f19ffa766885be9f80b8` | Pending |
| `tools/equivalence/main_test.go` | `89be2708f8ca84a8113df978e9cc0d9246479cca` | Pending |

## Review limitations and execution note

The first broad remote-read batch failed before usable completion. It is not credited as reviewed coverage; subsequent retrieval used smaller batches. One guessed test filename returned 404; the actual inventory path was then read. Neither failed retrieval changes the repository or constitutes a product defect. Future coverage must use the inventoried paths and actual examined content.
