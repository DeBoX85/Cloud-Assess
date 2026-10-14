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

# Section 1 API and lifecycle reliability correction candidate, 2026-10-14

Current state: [SESSION_HANDOVER](SESSION_HANDOVER.md), [Section 1 report](SYSTEMATIC_REVIEW_SECTION1_20261014.md), [PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). Final exact-head QA and acceptance remain pending; historical states below are not current authority.

The candidate preserves healthy semantic source data while correcting three failure contracts: direct ARG recommendation execution returns terminal caller context, preserving ordinary initiating error when only internal sibling cancellation fires; inventory adds id asc after existing subscription/resource-group sorting; ARG2024-04-01 accepts only its documented HTTP200completed-query success. Other2xx responses cannot become valid empty assessment success. The stronger inventory oracle derives projection/sorting from independent rows, with four calibration cases and compiling production controls. These deliberate reliability corrections do not guarantee changing-estate snapshot consistency.

Original source captures, pins, schemas and dependency graph remain unchanged. Both-host failing evidence precedes each production correction. The current resume action is final candidate QA/review and Section1checkpoint reconciliation, not replaying completed changes. No Azure write or release action is authorized by this audit record.

---

# Current security remediation authority, 2026-10-11

Read [SESSION_HANDOVER.md](SESSION_HANDOVER.md), the
[security remediation contract](SECURITY_REMEDIATION_20261011.md) and
[live PR133](https://github.com/DeBoX85/Cloud-Assess/pull/133) for the latest
authored head, native/source/review, protected acceptance and distinct accepted-push
outcomes. PR133 preserves the PR132 handover and supplies its bounded security
correction. PR132 is not independently accepted; do not replay or merge its failed
candidate. Accepted implementation was PR131 b653a3ab at remediation preparation;
verify the live branch before treating a proposal as accepted.

The correction proposes Go1.26.9 and x/net0.60.0 with only independently reviewed
minimum graph requirements. Isolated preparation scans do not replace final native
checks. Prior Go1.26.8/current failed/pending/zero-findings statements below are
historical observations, superseded for current resume instructions by this entry.
Contractual behavior, source pins, failures and historical evidence remain preserved.
Region feature work waits for verified security acceptance. Laptop/Azure restoration,
live scope, release/Gate004/load/fresh OS and independent-person approval remain
unestablished. No Azure writes or scope substitution are authorized.

---

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

Verified implementation: [PR131](https://github.com/DeBoX85/Cloud-Assess/pull/131),
accepted `b653a3abfc35590a185531095a168a9a107e769e`, tree
`47515c474ddf400da0a68c556e9147d8d3c2576b`. Live refs can advance.
Read [HANDOVER_CHECKPOINT_20261009.md](HANDOVER_CHECKPOINT_20261009.md) for exact
parents/identity/pins, separate candidate/accepted-push proof, audit/findings
disposition, remaining tasks, access limits, workspace recovery and QA provenance.
This documentation checkpoint's final acceptance evidence is in its live PR body
on `docs/handover-20261009`; do not assume it is already accepted.

Product requirements are unchanged by this handover. Accepted internal region
milestones extend through PR131; historical preparation/pending descriptions below
must not be read as the current implementation state. Public region execution,
full AZQR feature parity, live validation, Gate004 and release readiness are not
established. The linked checkpoint separates implemented, accepted and remaining work.

Earlier records below are preserved historical snapshots, superseded only in
current-state/resume instructions by the verified checkpoint above.

---

# PR120 availability implementation and recovery authority

The bounded pure availability constructor is implemented in PR120. Its live PR body is the exact final-head/merge/accepted-push evidence index: [PR120](https://github.com/DeBoX85/Cloud-Assess/pull/120). Determine candidate versus VERIFIED OFFLINE acceptance from that index and live refs, never from a historical pending paragraph. Entry accepted baseline is audit PR119 merge3463a72f11912039702ee60605ec4d79830c09bc/tree3bc82d8ffd4aa19167edd4e6bf02b8d031a9939a. Audit and source-only availability characterization are accepted; no entire-audit repeat is needed.

[REGION_AVAILABILITY_RUNTIME.md](REGION_AVAILABILITY_RUNTIME.md) defines the pre-edit contract, explicit scope/evidence/cancellation/response-collision corrections, bounds and independent acceptance. One comparison consumes selected owned InventoryCalculation inputs and subscription/target-bound decoded evidence; health stays separate and must be propagated by later integration. Eight complete retained-source comparisons and literal unknown/completeness fixtures are complemented by CalculateInventory-to-availability integration, independent entry/text/output boundaries, scope/evidence rejection, owned concurrent/repeated results, cancellation and ten compiling controls/restored baselines on both native hosts. No clients/global cache/registry/credentials or public region execution are added. Source captures/pins/dependencies/schemas are unchanged. Reviewed self-authored code and source contracts; no independent-person approval claimed.

Corrected8514b9f73af3aadb1a4aeca35e23d292a7abbabb/tree3af6662388963becf0d8907f42660e12475e4f20 passed Linux native37223966655/job111499731513 and source37223966664/job111499699784; full logs/six contiguous captures independently hashed and byte-matched, candidate scope/Denis identity/parent/remote bytes/local reconstructed tree/preview54744530232ff2982294be3e72e7dc882596ef41 and active no-bypass rules23890737 verified. Windows111499731689 also passed every mandatory step/full logs, including all ten compiling availability controls/restored baselines, actual CLI/default/custom packages and vulnerability scan. Existing module-only advisory remains open. Added owned-inventory integration test and this final continuity reconciliation require fresh complete final-head QA; prior head success is not transferred. Initial ee5f0b78 formatting failure37223734379/Linux111499031673 is FN068 recurrence; Windows111499031882 was superseded/cancelled under the unchanged workflow concurrency rule, not runtime product failure. No blind rerun, pin/gate relaxation or uncertain remote mutation.

Exact next: inspect final PR120 native Linux/Windows and source logs/captures/preview/protections; only then protected expected-head acceptance and distinct accepted-push proof. After verified acceptance, implement bounded latency/cost/quota/reservation calculations and adapters, then coordinator/public/all-format/default/empty/partial-health integration. Laptop access is expected soon and Azure restoration is pending confirmation. Prepare the approved-scope live handoff when access is actually restored; no Azure resources/roles/production substitution/exposure/release authorized. Live/DV-001/load/fresh-OS/hosted-maintenance/Gate004/release and module-only advisory remain open. Local Go and materialized APRL remain unavailable; Actions is execution authority. Preserve old workspaces and inspect actual remote outcomes before retrying an interrupted operation. Published code is backed up only after remote byte/object/ref verification.

Earlier checkpoint/pending/absence descriptions below are retained history, superseded by this authority and the live PR evidence.

---

# Cloud Assess Target Specification

Current audit disposition (2026-10-04): whole-project bounded offline review and hardening/source characterization are accepted through PR114, mergec33fb2eee93a2e5730f72612b351d0b838737d03, tree583496773c23902e7d7187ffc83a1c4ee6be0713. Owned inventory aggregation remains accepted through PR112; source-only availability characterization is now accepted through PR114. Runtime availability and public region-selection remain absent. Resume the separately reviewed bounded pure availability implementation after the audit disposition follow-up is accepted. Exact runtime and disposition acceptance evidence is indexed in their PR bodies; live/load/Gate004/release/advisory boundaries remain open. Earlier candidate/IN PROGRESS sequencing below is historical.

Current audit reconciliation (2026-10-04): owned inventory aggregation is VERIFIED OFFLINE through PR112, accepted07011b63. All six pure region helpers are implemented; public region-selection remains unavailable. The combined PR114 candidate contains PR113 source-only availability characterization and four audited hardening corrections, pending fresh combined exact-head gates/protected acceptance. Earlier dated IN PROGRESS/next-task paragraphs are historical. SESSION_HANDOVER is the resume authority; PROJECT_AUDIT_REPORT_20261004 records review scope and residuals.

Status: working canonical specification

Reference implementation: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

## Purpose

Cloud Assess is a read-oriented Azure assessment engine. It discovers Azure resources in a selected scope, evaluates them against recommendation sources and Azure platform signals, and emits human-readable and machine-readable reports.

The primary flow is:

```text
Authenticate
  -> Resolve scope
  -> Discover resources
  -> Apply filters
  -> Load recommendation catalog
  -> Execute assessment stages
  -> Normalize findings
  -> Build deterministic summary
  -> Render reports
  -> Evaluate optional severity gate
```

## Compatibility position

Cloud Assess targets behavioral equivalence with the pinned reference implementation, not CLI, configuration, package, or branding compatibility.

Intentional changes include:

- new product and CLI identity
- neutral configuration schema
- `GraphRecommendation` concept becomes `RecommendationDefinition`
- `GraphResult` concept becomes `Finding`
- lower-level packages return errors instead of terminating the process
- every requested assessment stage records explicit execution health
- assessment completeness is separate from finding severity
- subscription-ID masking is described explicitly as redaction rather than general anonymization

## Branding adjustment requirement

Product identity is centralized in `internal/branding`. The [immutable branding profile builder](BRANDING_PROFILES.md) now supports validated build-time presentation identity without source edits. The [branding review](BRANDING_REVIEW.md) defines the acceptance checks; [paired synthetic report evidence](BRANDING_REPORT_QA.md) now implements the custom report/data checks, and [profile-aware candidate packages](BRANDED_PACKAGES.md) implement distribution checks. The bounded five-field workflow passed coherent native CLI/report/package acceptance through PR75; runtime profiles, logos/themes and release approval remain separate requirements. A development builder alone would not establish completion.

## Core v1 scope

Core v1 includes:

- Go CLI application
- Azure SDK `DefaultAzureCredential`
- Azure Public, Government, China, and custom ARM/authority endpoints
- tenant-accessible subscription discovery
- management-group recursion
- explicit subscription and resource-group scope
- resource inventory through Azure Resource Graph
- include/exclude filtering
- scanner registry and resource-type pruning
- APRL recommendations
- Azure Orphan Resources recommendations
- custom recommendations
- YAML/KQL recommendation extensions
- diagnostics settings assessment
- Azure Advisor
- Defender for Cloud status
- Defender recommendations
- Azure Policy noncompliance
- Arc-enabled SQL reporting
- Cost Management reporting
- XLSX, JSON, CSV, SARIF, and JSON stdout output
- subscription-ID redaction
- post-report severity gate
- plugin architecture

Full internal-plugin parity remains open. Zone mapping execution, the standalone `zone-mapping` command, named mixed selection and honest `plugins list/info` are verified offline through PR84/85, service-health execution through PR87 and SQL EOL execution through PR89 and carbon request/public execution through PR93/94; see [ZONE_EXECUTION.md](ZONE_EXECUTION.md). Unnamed or unavailable plugin selections reject before Azure authentication. AI governance source captures/pure projection are verified offline through PR96; its bounded decoder/request library is verified offline through PR97 under [AI_GOVERNANCE_REQUESTS.md](AI_GOVERNANCE_REQUESTS.md). [AI Graph discovery](AI_GOVERNANCE_EXECUTION.md) is accepted offline through PR101; public/report integration is accepted offline through PR103, and [region selection primary scoring](REGION_SELECTION.md) is accepted offline through PR104 and source auxiliary characterization through PR105; [quota/reservation pure projections](REGION_AUXILIARY.md) are accepted offline through PR106; [service-availability pure sheets](REGION_SERVICE_AVAILABILITY.md) are accepted offline through PR108; [pure CostComparison](REGION_COST_COMPARISON.md) is accepted offline through PR109; [pure Inventory](REGION_INVENTORY.md) is accepted offline through PR110; [inventory calculations characterization](REGION_INVENTORY_CALCULATIONS.md) is accepted offline through PR111; [owned inventory aggregation](REGION_INVENTORY_AGGREGATION.md) is accepted offline through PR112; [source availability characterization](REGION_AVAILABILITY_CALCULATIONS.md) is IN PROGRESS/unaccepted; full region execution remains unavailable; live extension equivalence is not established by these synthetic slices.

## Output strategy

Both Excel and JSON are first-class outputs.

Excel is the default human-facing report.

JSON is the canonical machine-readable representation used for automation, integration, and equivalence testing.

Neither format replaces the other. Both are projections of the same canonical assessment result.

## Scope rules

Supported scope entry points:

- all accessible active subscriptions
- one or more management groups, recursively resolved
- one or more subscriptions
- one or more resource groups within exactly one subscription

Invalid combinations:

- management group + subscription
- management group + resource group
- resource group without subscription
- resource groups across multiple explicit subscriptions

Disabled and deleted subscriptions are excluded.

## Resource inventory

Inventory is collected before assessment execution and captures at least:

- resource ID
- subscription ID
- resource group
- location
- resource type
- resource name
- tags
- SKU name
- SKU tier
- SKU family
- SKU capacity
- kind

Discovery produces both in-scope and out-of-scope collections.

## Filtering semantics

Filters support subscription, resource group, scanner/service selection, individual-resource exclusion, recommendation exclusion, and tags.

- explicit subscription and resource-group include lists narrow scope
- an explicitly included subscription or resource group takes precedence over a matching exclusion, preserving reference behavior
- the legacy-compatible `include.resourceTypes` field contains scanner/service keys such as `aks`, `ca`, `st`, and `vnet`; it does not contain literal ARM resource-type strings
- selected scanner keys are expanded into an internal allowed ARM resource-type set
- that structural scanner/subscription/resource-group/resource-ID scope is applied both during inventory discovery and again to downstream findings
- exact individual-resource exclusions remove the matching resource
- all include tags must match
- any matching exclude tag excludes the resource
- tag keys are case-insensitive
- tag values are exact-match
- child findings inherit the nearest recorded parent tag-scope decision
- unknown tag scope is excluded for include-tag filters and retained for exclude-only tag filters

The target configuration root uses neutral terminology such as `assessment:` rather than the legacy product name. The target also renames the legacy `exclude.services` field to the clearer `exclude.resources`; this is an intentional configuration compatibility break, not a behavioral change.


Filter files accept one YAML document with known neutral schema fields. Unknown or legacy keys, additional documents and a null `assessment` fail before authentication. An empty file and `assessment: {}` preserve default filter behavior. This is an intentional input correction to prevent silently discarded scope restrictions; see [AUDIT_FILTER_SCHEMA.md](AUDIT_FILTER_SCHEMA.md).

## Recommendation catalog

Recommendation sources are normalized into one `RecommendationDefinition` model.

Initial source identifiers:

- `APRL`
- `AOR`
- `CUSTOM`
- `DIAGNOSTICS`
- `PLUGIN:<name>`

APRL remains independently pinned. Existing custom rules are retained but re-homed under neutral paths and source labels.

## Finding model

A `Finding` represents one recommendation affecting one Azure resource. It is not tied to Azure Resource Graph as an implementation mechanism.

This allows KQL rules, ARM/API diagnostics checks, and future assessment mechanisms to share one domain model.

## Core pipeline

The logical pipeline is:

```text
Initialization
Scope discovery
Resource discovery
Recommendation catalog construction
Scanner pruning
Graph recommendation execution
Diagnostics
Advisor
Defender status
Defender recommendations
Policy
Arc SQL
Cost
Plugins
Findings summary
Report rendering
Severity gate
```

Stage-level execution remains sequential unless safe concurrency is demonstrated. Work within a stage may be concurrent and bounded.

## Stage health

Each stage records one of:

- `completed`
- `completed_with_warnings`
- `skipped`
- `failed`

A failed data retrieval must never be represented as a successful empty dataset.

Assessment completeness is separately represented as:

- `complete`
- `complete_with_warnings`
- `partial`
- `failed`

Execution completeness applies to resolved/visible data; it does not independently certify intended estate visibility. Canonical JSON now records requested scope, effective subscription filters, resolved subscriptions (including zero-data subscriptions) and unresolved explicit subscription IDs. A successfully discovered scope missing any explicit CLI/library subscription fails before resource queries, persists failed status evidence when possible and returns exit 1. This is a deliberate target correction to silent subset discovery. All-visible/MG and filter-only successful empty scopes retain source behavior; ARG can still omit inaccessible resources. See [SCOPE_REPORTING.md](SCOPE_REPORTING.md), [ACCESS_MODEL.md](ACCESS_MODEL.md) and [OPERATIONS.md](OPERATIONS.md). Independent membership/RBAC and restricted-identity live evidence remain Gate 004 requirements.

## Assessment lifecycle limits

An opt-in `--assessment-timeout` / application AssessmentTimeout bounds coordinator execution through context propagation. Zero preserves the compatibility default; rendering and local CLI initialization are outside the budget. SDK scope continuation cycles now fail explicitly, and enabled stages check cancellation before/after execution. No process hard kill, universal positive default or volume cap is claimed. See [PAGINATION_LIFECYCLE.md](PAGINATION_LIFECYCLE.md).

## HTTP operation limits

The production authenticated HTTP client applies a positive `OperationTimeout` to one HTTP call, covering authentication, retry waits and response-body consumption. Defaults are ten times the per-attempt timeout. Earlier caller deadlines and cancellation remain authoritative; non-positive custom values add no operation deadline. This is not a whole-scan time limit or a guarantee that custom transports ignoring context terminate. Callers of `PostStream` must close response bodies; operation context cleanup occurs on body read completion/error or close.

ARM SDK assessment clients explicitly disable automatic resource-provider registration. A registration-required read error must remain an error, without a registration POST; caller-supplied scope options are copied before enforcing this assessment boundary.

The default shared HTTP transport does not follow redirects. Advisor metadata continuation URLs must remain on the configured HTTPS host/port, contain no user information or fragment, and not repeat. Unsafe pagination fails the Advisor stage rather than sending the ARM credential to another origin or reporting a truncated successful metadata result. These are target security corrections for unsafe response paths; scope SDK pagers now validate configured HTTPS origin before authentication and refuse redirects in their default transport. Injected transports/policies remain trusted code with separate internal redirect/destination obligations; same-origin pagination budgets remain open. Shared HTTP adapters use configured ARM audience for token scope separately from endpoint. See [ALIGNMENT_REVIEW.md](ALIGNMENT_REVIEW.md).

## Failure model

Assessment-fatal examples:

- authentication cannot initialize
- requested scope is invalid or cannot be resolved
- core resource inventory cannot be obtained
- core recommendation engine cannot execute sufficiently

Stage-local failures such as Advisor, Policy, Defender, Cost, or optional-plugin failures should normally preserve other valid results and mark the assessment partial.

## ARG behavior to preserve

- up to 300 subscriptions per request
- request at most 1,000 records per page, matching the ARG 2024-04-01 contract; this deliberately corrects the pinned source request of 5,000
- reject explicit service truncation without a continuation token rather than publish partial query success; see [ARG_COMPLETENESS.md](ARG_COMPLETENESS.md)
- skip-token pagination
- bounded rule concurrency
- malformed-row tolerance
- malformed-row visibility through warnings/completeness metadata
- unsupported logical-table skip semantics
- management-group-aware authorization where required

## Diagnostics behavior

Diagnostics uses the discovered inventory and ARM batch reads of `Microsoft.Insights/diagnosticSettings` to generate canonical findings for supported resources without diagnostic settings.

Errors are propagated rather than terminating from worker goroutines.

When an ARM batch response contains a different number of subresponses than requested, Cloud Assess treats the Diagnostics stage as failed rather than inferring that the unmatched resources lack diagnostic settings. This safety correction applies to malformed API output and is tracked separately from successful-response equivalence with the pinned source.

## Cost behavior

The source implementation defines the cost period as the previous completed UTC calendar month.

Metric: `ActualCost`

Aggregation: sum

Grouping: `ServiceName`

Cost-access failure must be distinguishable from a valid zero-cost response.

Current Cost queries aggregate each resolved subscription. Resource-group, resource-ID and tag filters do not narrow the billing query to selected resources; reports and operations must preserve that distinction.

## Summary semantics

Core findings are deduplicated by normalized recommendation ID + resource ID.

Recommendation applicability is represented explicitly as:

- `not_applicable`
- `compliant`
- `noncompliant`

SLA-category records are not counted as ordinary findings.

## Reports

Initial report formats:

- Excel, default human-readable report
- JSON, canonical machine-readable result
- CSV
- SARIF 2.1.0
- JSON stdout

The Excel report should preserve the conceptual datasets from the reference and add an assessment-status section or worksheet showing stage execution health and assessment completeness.

## Redaction

The existing masking behavior is specifically subscription-ID redaction. It does not anonymize the entire report.

Cloud Assess uses the explicit `--redact-subscription-ids` control.

When enabled, subscription-ID redaction applies to:

- Excel
- CSV
- canonical JSON files
- canonical JSON stdout

For JSON, known subscription IDs are also redacted when embedded inside serialized ARM/resource identifiers and related strings, not only in dedicated `subscriptionId` fields.

Report files are staged before replacement, preserving an existing report on rendering failure. Unix replacement files use mode 0600; New Windows files inherit destination-directory ACLs; replacements preserve and protect the previous effective DACL before writing report bytes, failing if it cannot be preserved. Multi-file exports are not transactional, and cross-platform atomic replacement or crash durability is not promised. CSV formula-like cells are prefixed as text; this includes negative numeric strings and changes the human CSV representation. Canonical JSON retains original values for machine consumers. Spreadsheet consumer/re-save behavior remains outside the neutralization guarantee. See [QA_PROCESS.md](QA_PROCESS.md).

SARIF intentionally retains stable Azure resource identities because resource IDs and fingerprints are part of its automation/baselining contract, matching the identity-bearing behavior of the pinned reference. SARIF must therefore be treated as sensitive/identity-bearing output even when subscription redaction is enabled for the other report formats.

Redaction is not full anonymization: resource names, resource groups, tags, subscription display names, and other potentially identifying metadata can remain.

## Deferred functionality

Deferred until after core v1 equivalence:

- full internal-plugin parity and remaining live extension validation
- remaining internal adapters and live validation (bounded `plugins list/info` is accepted offline through PR85)
- MCP server
- historical report compare
- alternative VM SKU utility
- public documentation site
- package-manager/public distribution
- web UI or hosted API

## Core v1 acceptance

Core v1 is complete when Cloud Assess can reproduce the pinned reference implementation's material core assessment behavior for the same Azure inputs while also providing explicit stage health and assessment completeness.

## Offline recommendation inspection

The `rules` command now lists source-selected pinned embedded and Diagnostics recommendations without Azure access, with Markdown by default and local --json/-j output. It does not load external YAML plugins or accept scan scope/report flags. JSON matches the independently captured pinned source output; Markdown escaping and safe absent-guidance handling are documented target presentation corrections. See [RULES_INSPECTION.md](RULES_INSPECTION.md). Live execution, estate coverage and full feature/release acceptance remain separate.

Scanner-specific commands now follow the pinned source command keys and reuse generic orchestration; see [SCANNER_COMMANDS.md](SCANNER_COMMANDS.md). This is offline command/selection evidence, not live validation of every service.

Bounded [YAML Graph plugin execution](YAML_GRAPH_PLUGINS.md) now joins ordinary scan orchestration after input validation before authentication. This preserves healthy source schema, precedence and external-rule filtering while applying explicit bounded/strict input protections. Explicit zone table-plugin execution is now available; external queries do not depend on that stage. Offline parser/transport/report checks do not establish live plugin equivalence.

[Zone mapping](ZONE_MAPPING.md) preserves source rows while following documented safe continuations and returning explicit partial-failure information. [Execution/CLI integration](ZONE_EXECUTION.md) is accepted offline with canonical table/schema/health/privacy/report checks. Live zone and full internal-plugin parity remain open.

[Canonical plugin tables](PLUGIN_TABLES.md) provide library/report infrastructure with additive schema 1.1 for explicitly constructed plugin-bearing assessments; ordinary assessments retain schema 1.0. Source-visible metadata/columns/rows remain reconstructable, with explicit owned identities and honest health. Zone plugin-stage/CLI wiring is accepted through PR84/85; no full internal-plugin or live parity is claimed. [Service-health](SERVICE_HEALTH.md) has accepted source-grounded library and standalone/mixed execution through PR86/87. Its temporal-query limitation is documented explicitly.


[SQL EOL](SQL_EOL.md) preserves the exact pinned 611-line query, source metadata and 32 string columns with subscription-only filtering. Library and standalone/normal/scanner execution accepted offline through PR88/89, including all-format partial/privacy tests. Static source model assumptions, live KQL/financial/licensing validation and Arc SQL numeric-vcores remain separate open boundaries.


[Carbon emissions](CARBON_EMISSIONS.md) migration starts with source-captured pure aggregation, preserving service-selected latest dates, eight columns and float/conditional display. Deterministic row order and strict finite/date/label/volume health are explicit corrections. Bounded HTTP/access/pagination and standalone/mixed/scanner command/report integration are accepted offline through PR93/94; explicit denied or failed batches retain valid prior sums with incomplete health. Live service totals/access/sovereign availability and representative load remain separate open boundaries; source silent failed batches cannot become target healthy empty output.


Region Inventory explicit composition correction: pinned source region Inventory is silently skipped in Excel if the core Inventory exists; canonical plugin sheets reserve that core name. The pure helper uses Region Inventory while preserving all ten source cells/order/description/capacity/masking. Core sheet collision rules are unchanged. Selected UUID/resource-ID correlation and arithmetic/text/work guards are deliberate safety corrections; standalone/mixed/report integration remains separate. CostComparison is VERIFIED OFFLINE through [PR109](https://github.com/DeBoX85/Cloud-Assess/pull/109), merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents. 

Final candidate1cf8e4cf6adc38ffe377894824b51f3b26431373/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, parentca16157 on accepted1957f6d7. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165402 quality111338402524/Windows111338402644 passed every mandatory step; only conditional failure upload skipped. Both tested preview77c2b59d40f89e1966402668a98efcdfad2e4d5d with identical final tree/ordered base/head parents. Full logs inspected: all nine Cost controls/restored baselines, ten service/three auxiliary/five primary/eight request/four execution controls, actual CLI/default+branded package/docs/provenance/module/inventory/branding guards; Linux race/vet/fuzz and82.5% coverage; zero reachable/imported vulnerability findings, existing module-only advisory remains open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165342 /111338403371 completed; full14 branches, indexed3-input/4-output chunks, exact retained bytes/hashes and source/tree/APRL provenance verified. Complete14-file scope/deletions/remote bytes/original scratch blobs, parent/human identity/pins/base/head/preview/current no-bypass rules23890737 verified. Semantic self-review checked actual unchanged source first-tuple/item break behavior, all source topology/cells and safe bounded identity/nonfinite/text/cancellation behavior. No callable final manual Code Review tool or independent-person approval is claimed; automatic comments currently empty. FN072 compile failure was not accepted. Expected-head protected merge and distinct accepted-push evidence follow; no fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure.

Protected acceptance verified: PR109 merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents, Denis author/GitHub committer/livecore ref and no open proposals. Distinct accepted-push native/source QA remains pending; no dependent live proof claimed.

Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513599 quality111339415060/Windows111339415184 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513582 /111339414976 passed on exact6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f. All mandatory steps/full logs inspected; only conditional Linux failure upload skipped. Both hosts passed9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages/docs/branding/provenance; Linux race/vet/fuzz82.5%, zero reachable/imported findings. Source14 complete branches, indexed chunks/both retained exact bytes/hashes and pinned source/tree/APRL provenance verified. Module-only advisory and all live/public/laptop/Azure/Gate004/release obligations stay open. No fresh local Go/source or independent-person approval claimed.



2026-10-04 historical PR111 preparation baseline, superseded below: Inventory is VERIFIED OFFLINE through [PR110](https://github.com/DeBoX85/Cloud-Assess/pull/110), merge2dcd8dae924ce3fad5402890cdeca0ecbef6e01e/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, ordered6dfcd1be/76ac3f8d parents. 

Final76ac3f8d65acf9970305c58186e13bb19cf78ffe/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, parent accepted6dfcd1be. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985397 quality111340777403/Windows111340777269 passed every mandatory step; only conditional failure upload skipped, full logs inspected. Both tested preview4e8f21f81155a55383272fe7ade895896fb72db6 with identical tree/orderedbase+head parents. Eight Inventory controls/restored baselines and all9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution controls passed each host;12 actual CLI/19 default/19 branded packages/docs/branding/provenance/module/inventory checks. Linux race/vet/fuzz passed; 82.7% coverage, zero reachable/imported findings; existing module-only advisory stays open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985428 /111340777428 passed full14 branches/indexed3-input4-output chunks/both retained exact bytes/hashes/source-tree-APRL provenance. Original scratch Go/test/script blobs/complete16 remote files/scope/no deleted or pinned files/commit identity/tree/parent/base/preview/rules verified. Semantic self-review checked actual source helper/mask/SKU/duplicate-skip plus explicit sheet/scope/budget/product corrections. Automated comments currently empty; no callable final manual Code Review tool or independent-person approval is claimed. No fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure. 
 Protected acceptance/ref/tree/parents/Denis author/GitHub committer verified. Distinct accepted-push proof follows below.
New inventory_capture_test.go.txt/extended capture-region-source.py invoke only unchanged pure inventory/merge/normalization/predicate functions in isolated pinned source/APRL.13 structural branches/30 literal normalization cases; actual output retention/topology/hash/final source/native/remote acceptance is pending. No target calculation code, pin/public/request/live/platform change. See REGION_INVENTORY_CALCULATIONS for scope/contracts/acceptance/rollback/next action.


PR110 distinct accepted-push verified: 
Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170292930 quality111341722998/Windows111341722896 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37170292817 /111341722273 passed on exact2dcd8dae924ce3fad5402890cdeca0ecbef6e01e. All mandatory steps/full logs inspected, only conditional failure upload skipped. Both hosts passed8 Inventory/9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines;12 CLI/19 default/19 branded packages/docs/branding/provenance. Linux race/vet/fuzz82.7%; zero reachable/imported findings, existing module-only advisory open. Source14 full branches/complete indexed chunks/retained hashes/bytes/source-tree-APRL provenance exact. No fresh local Go/source/public/live/laptop/Azure/Gate004/release/advisory closure.



Observed source characterization preparation is accepted through PR111; full final-head and distinct accepted-push evidence is indexed in SESSION_HANDOVER and REGION_INVENTORY_CALCULATIONS. The original observed fixture hashes remain unchanged.

## PR111 accepted and owned aggregation active (2026-10-04)

Owned inventory aggregation is VERIFIED OFFLINE through [PR112](https://github.com/DeBoX85/Cloud-Assess/pull/112), accepted07011b63. Pure source availability characterization is IN PROGRESS/unaccepted under [REGION_AVAILABILITY_CALCULATIONS.md](REGION_AVAILABILITY_CALCULATIONS.md). Public region-selection remains unavailable. Retain actual new source bytes and require unconditional six-file equality, independent oracle checks and both complete native gates before protected acceptance. Then implement bounded availability, followed by latency/cost/quota/reservation arithmetic, adapters and coordinator/public/all-format execution. Azure/laptop validation remains deferred.

Accepted PR111 final native37170957575/source37170957525 and distinct accepted-push native37171306079/source37171306081 passed all mandatory steps/full logs on the documented exact heads; unchanged four captures/13 inventory branches/30 normalization cases and all prior controls. See SESSION_HANDOVER/REGION_INVENTORY_CALCULATIONS and [PR111](https://github.com/DeBoX85/Cloud-Assess/pull/111) for commit/tree/ordered parents/job/coverage/security proof. Current aggregation contract separately defines corrections, independent budgets, actual complete five-map oracles, valid literal subset, ownership/cancellation and ten compiling controls; no new Go/native/public/live acceptance yet.


## Current checkpoint (2026-10-04)

Owned inventory aggregation is VERIFIED OFFLINE through [PR112](https://github.com/DeBoX85/Cloud-Assess/pull/112), accepted07011b63. Pure source availability characterization is IN PROGRESS/unaccepted under [REGION_AVAILABILITY_CALCULATIONS.md](REGION_AVAILABILITY_CALCULATIONS.md). Public region-selection remains unavailable. Retain actual new source bytes and require unconditional six-file equality, independent oracle checks and both complete native gates before protected acceptance. Then implement bounded availability, followed by latency/cost/quota/reservation arithmetic, adapters and coordinator/public/all-format execution. Azure/laptop validation remains deferred.

PR112 final native/source and protected merge proof: [PR112](https://github.com/DeBoX85/Cloud-Assess/pull/112), [DEVELOPMENT_LEDGER.md](DEVELOPMENT_LEDGER.md). Distinct accepted-push verification remains separately recorded there. Earlier preparation paragraphs are historical snapshots, not current unresolved gates.

## Owned latency acceptance checkpoint

# PR121 owned latency implementation checkpoint

UNACCEPTED WIP on feat/region-latency; accepted baseline remains PR120 ac88f214. Verified retained-source checkpoint7c9ecb5bb574ac85aac9641f7d76c49faa43376c/tree62336e4e13cfaae2a2a2a7f98ade0017797ecc12 passed strict fresh nine-file source37227837952/job111511077645. Windows native37227837946/job111511077776 passed full mandatory steps; Linux111511077894 stopped after compiling four oracle controls/restored baselines at two gofmt spaces (FN068 recurrence). Exact formatter diff applied. No broad Linux success transferred.

The pure latency constructor now embeds the exact observed historical dataset, decodes per call with bounded/pinned bytes, admits selected complete comparisons through Project, validates finite/canonical/bounded decoded data, deterministically computes directional cluster means, preserves all fields/order with four cloned slices and zone maps, and returns separate latency-only warnings. New tests compare24 complete valid retained-source comparisons, explicitly reject the source's negative/display inputs, verify full embedded provenance, exact nested/64-byte/1MiB boundaries, caller/concurrent/repeated ownership, cancellation before/during/after, and actual availability-to-latency-to-Project row/health integration. Thirteen compiling selected runtime controls/restored baselines are mandatory on both native hosts. No target Go execution yet for this implementation.

Next: classify fresh exact-head native feedback, correct only justified failures and inspect full Linux/Windows/source logs,84 contiguous capture chunks/hashes, all source pins, complete diff/remote bytes/tree/identity/base/preview/no-bypass protections before expected-head protected acceptance and distinct accepted-push evidence. PR121 body is exact current publication/run index. No public region execution, Azure/live/laptop/load/Gate004/release/advisory closure. Laptop/Azure access is expected soon but not confirmed; existing core live queue remains independently resumable once access and approved scope are confirmed. No Azure writes/roles/fixtures/production substitution/exposure/release. Preserve historical records/workspaces; no local Go/materialized APRL. Never call scratch-only work backed up.

## Cost-enrichment source characterization

PR122/REGION_COST_ENRICHMENT defines26 complete independent source comparisons, actual cost/types blob hashes, strict11 retained-source file equality and6 compiling fixture controls/restored semantic assertions on both mandatory native hosts. Source-only acceptance does not certify target weighted cost arithmetic, subscription/resource evidence ownership, collection completeness or public execution. Its live PR body and current SESSION_HANDOVER are exact final-head/accepted-push evidence authority. Target/runtime/collector/live obligations remain open.


## Owned weighted cost runtime

PR123/REGION_COST_RUNTIME.md implements primary weighted cost enrichment with explicit selected-subscription history and pricing evidence, separate honest cost health, bounded deterministic arithmetic and owned results. It is UNACCEPTED pending full native/source proof. Public region execution and collector completeness remain unavailable; source captures/pins remain unchanged. Mandatory gates add19 compiling runtime controls/restored baselines on both hosts, complete valid/corrected source results and independent boundary/ownership/cancellation/integration assertions. PR123 body indexes exact acceptance, not this historical pending status.
