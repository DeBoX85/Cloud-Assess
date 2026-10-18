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

# Section 1 systemic review, 2026-10-14

## Current state and next action

**Section 1 code review and correction QA are complete on verified code candidate `00ba2716f76366f6663de15531fdef4de6b812f8`; this is the planned user re-evaluation checkpoint.** No protected acceptance, completed whole-codebase audit or release readiness is claimed.

The accepted baseline remains `a6fee427e19feeef28691737ea9bc5026086efa0`. Read [live PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142) and current branch refs for the actual candidate head and latest QA/review results. Earlier preparation records are history and must not be replayed.

The next action is the user's planned progress re-evaluation. The documentation-only successor records verified code evidence on00ba2716; it does not inherit exact-head workflow acceptance. Inspect fresh branch refs and any successor checks before merging. Do not repeat the already completed original cancellation or inventory corrections and their original reproduction runs. Do not start Section 2 before the planned user re-evaluation.

## Coverage and intent

All 95 baseline Go files in nine package directories received a source/test pass and fresh full-content/Git-blob verification: 38 production files and 57 test files. The candidate adds two Go test files, giving 97 current Section 1 Go files. The new regression-control Python helper and changed CI execution path are also reviewed for their immediate correctness; broader tooling review remains Section 4.

| Package | Production | Baseline tests | Review |
| --- | ---: | ---: | --- |
| `cmd/cloud-assess` | 4 | 12 | Reviewed; correction QA verified at00ba2716 |
| `internal/app` | 1 | 9 | Reviewed; correction QA verified at00ba2716 |
| `internal/config` | 2 | 3 | Reviewed; correction QA verified at00ba2716 |
| `internal/assessment` | 8 | 0 | Reviewed; correction QA verified at00ba2716 |
| `internal/azure` | 8 | 7 | Reviewed; correction QA verified at00ba2716 |
| `internal/arg` | 6 | 11 | Reviewed; correction QA verified at00ba2716 |
| `internal/discovery` | 5 | 6 | Reviewed; correction QA verified at00ba2716 |
| `internal/throttling` | 1 | 1 | Reviewed; correction QA verified at00ba2716 |
| `internal/orchestration` | 3 | 8 | Reviewed; correction QA verified at00ba2716 |

The review traces flags and known-field filter parsing, prepared owned scope, SDK discovery, ARG batching/pagination and rule workers, ordered stage execution and application exits/rendering. It checks the target specification's semantic AZQR compatibility, read-oriented Azure operations, adjustable branding, explicit stage health and separate completeness/severity. Source-compatible behavior is retained unless an explicit, evidenced reliability correction is documented.

Full fixture/query ownership, collector-specific request budgets, business-record validity, standalone plugin/rendering contracts and delivery/security/license readiness remain assigned to Sections 2–4. Reading specific dependency functions here does not complete their package review.

## Confirmed findings and corrections

### SR-F01, P2: terminal cancellation can be reported incorrectly by the rule executor

The baseline can return nil error after canceled/expired work is skipped. Its first corrective candidate also retained the baseline's plain-worker-error priority, allowing an unwrapped transport error to mask the caller's terminal context.

Correction retains the parent context separately. Check context before scheduling and after receiving work. After joining workers, prioritize **parent** cancellation/deadline, then preserve the initiating ordinary query error when only the executor's sibling context was canceled. Successful ordinary empty-definition behavior stays unchanged; terminal work returns nil findings with a wrapped context error.

Public impact qualification: the stage runner independently checks context around tasks, and the coordinator returns terminal context for failed assessments. The direct/composed executor was inaccurate; these reproductions do not establish false-complete public scans.

| Evidence | Head | Native run and jobs | Outcome |
| --- | --- | --- | --- |
| Original terminal-work cases | `432f2034aa3eed973670ab5b3a31cd1775d94611` | 38561749374; Linux115844262674, Windows115844263005 | Compiled unchanged production; five expected cases failed |
| Remaining caller-priority cases | `a77e558f06bee40edc23f6e450f356f4cbc5bd87` | 38565950638; Linux115859094976, Windows115859094769 | Canceled/deadline cases returned plain query failure; ordinary-error control passed |
| First correction | `5dd97b0f530d578e17767a4eea14510b72ee6f6a` | Parent432f2034 | Original cases corrected; caller-priority coverage was insufficient |
| Parent-priority correction | `e65e66cd74d8e08ebfe1ea4452f6391a65e12878` | Parenta77e558f | Parent retained separately; full code QA verified at00ba2716 |

### SR-F02, P2: inventory pagination has non-unique sort keys

The target and pinned source order inventory by subscription ID and resource group. Multiple resources share both keys. Microsoft's primary guidance documents duplicate/missing paginated records even in static environments when ties reorder.

Correction `9d1872b1a1737da05db512addae8479e530aa238`, parent `e012a94d8077abf180917408156b12cfa69f47ec`, appends `id asc` only. Projection, scope, filtering and batching remain unchanged. This is a deliberate target reliability improvement; it does not guarantee a snapshot in a changing estate.

The initial 1001-resource fixture inferred complete pages from an ID token. Review correctly found that insufficiently independent. That old fixture and its native results are retained as historical evidence, **not** as final proof. Its replacement computes projected ID values and sorts independently constructed rows using the query's actual ordering. The second page uses another permitted input ordering for equal keys. Calibration checks tied keys, unique ascending and descending keys, and a non-unique projected ID. It is a narrow deterministic model, not a KQL engine or live Azure reproduction.

The stronger discovery tests passed on both hosts in 38565950638. Permanent compiling controls then remove the production tie-breaker and replace its ID projection with a non-unique subscription field. Each must fail the actual discovery regression after compilation, followed by exact-byte restoration and a healthy rerun. The old token-presence-only fixture would not detect the latter error.

Inventory semantic projection compares normalized resource-ID keys, not the source's ambiguous tied row order. Original source captures and independent source oracles remain immutable.

### SR-F03, P2: ARG transport admits undocumented 2xx statuses as complete queries

The shared HTTP client accepts generic 2xx success. The ARG2024-04-01 endpoint documents completed-query success as HTTP200 and other statuses as unexpected errors. The endpoint-specific transport must enforce that contract before accepting data.

Regression head `57ef9f93b9d0618e8fde87eb9323f91349b0ff54` adds actual shared HTTPClient/mocked-transport cases for200 versus201–206 before changing the production transport. Linux native38567759652/job115864719992 compiled and demonstrated all six unexpected statuses returning a successful empty result. Windows job115864720321 also compiled and demonstrated the same six failures; both full logs were read.

The correction adds an endpoint-specific HTTP200 guard before decoding. It closes the response through the existing defer and retains valid empty200 behavior. A fourth compiling/restoration control removes this guard and must recover the independently declared failure. Correction commit1f93442eca7e949c9b31f7af20b594c91ff10b62, parent57ef9f93b9d0618e8fde87eb9323f91349b0ff54, contains the guard and fourth control. Code QA is verified at00ba2716; the checkpoint documentation successor is a separate head.

## Review feedback disposition

Review of1243461f8ee884d9b0071e13bb219ab8f403ac45 found:
- 4270844237: caller cancellation lost to plain worker error. Reproduced independently on both hosts and corrected in e65e66cd.
- 4270844250: paging fixture independence gap. Replaced with projected-row sorting, calibration and production mutation controls.
- 4270844258: obsolete current resume instruction. The current report is reconciled to completed code QA and the user's checkpoint.

These are technical findings assessed against the user's task. References to old AGENTS procedures in automated comments do not reinstate the instructions the user revoked. The first review is not clearance for later commits. The new review on00ba2716 completed2026-10-14T13:20:05.783287Z with no major issues reported (marker6168441214 and result6169320382); it is automated review, not human or release approval.

Import cleanup `a98f2170d39c9dd8a53335c2bd147f4f9fc6496b` groups standard and external imports in the CLI. It changes no behavior.

## Architecture and residual boundaries

- Prepared requests own filter slices/maps, stage configuration and selected plugins. Per-scan state is local; tested repeated/concurrent coordinator use retains caller ownership.
- Invalid explicit scope and unresolved requested subscriptions fail before inventory. Retained scope evidence distinguishes missing selection from valid zero-data scope, without claiming RBAC or estate coverage.
- SDK scope requests constrain HTTPS origin before authentication, reject redirects/cycles and disable automatic provider registration. Injected transports/policies remain trusted code.
- ARG limits batches to300subscriptions and pages to1000rows, rejects explicit tokenless truncation and repeated tokens, and retains source-specific row/unsupported-table behavior.
- Scope/inventory/Graph failures stop execution; ordinary optional failures preserve other data and mark partial health. Context interruption stops further stages.
- Application reporting persists available health/results before execution/partial/severity exits. Rendering and local CLI initialization are outside the opt-in assessment deadline by contract.
- Exact registered plugin names are validated before AI-specific preflight; case normalization does not bypass that ordering.
- Positional resource-ID helpers match pinned source semantics exactly. They are not general ARM-ID validators; endpoint-specific admission and business-row implications remain separate.
- Request volume/memory budgets across all collectors and plugins remain Section2/3 work. No universal boundedness, hard kill, live parity or load guarantee is inferred.
- Source-compatible missing-ID stopping and tolerant null/empty rows need their business-validity/completeness implications assessed in Section2. Preserve their intentional compatibility status while assessing residual risk.

## Primary-source checks

Checked2026-10-14:
- Microsoft ARG paging: https://learn.microsoft.com/en-us/azure/governance/resource-graph/concepts/paging-results
- ARG REST2024-04-01 response/page-size/wire contract: https://learn.microsoft.com/en-us/rest/api/azureresourcegraph/resourcegraph/resources/resources?view=rest-azureresourcegraph-resourcegraph-2024-04-01
- Kusto sort defaults and direction: https://learn.microsoft.com/en-us/azure/data-explorer/kusto/query/orderoperator
- Go context at the actual project taggo1.26.9: official golang/go src/context/context.go, Git blobd00ac67e38296573d20520dbd7aceb2fd1a52906. A generic package-doc page resolved togo1.27.2, so the pinned source was used instead.
- Original AZQR8e4f0577f3615e6c9014c031bcad079f235369cc: inventory query and models.go positional helpers; models blob3724182c006dcb8adc42fd2d6fe877d32ebf1aea.

REST examples use string-enum truncation metadata despite generic Boolean terminology in the prose paging article. Decode according to the actual versioned wire contract. No dependency/toolchain/source-pin change is made.

## QA, counting and backtracking

Expected red tests are not infrastructure failures. Successful baseline/restored controls do not turn an otherwise failed or canceled workflow into an accepted candidate. Read all job steps and complete logs on the exact final head.

The old Linux314control count consists of313named Go mutation rejections plus one syntactically valid shell bootstrap control. Windows has309Go-labeled controls. Do not call all314Go compiling controls. The candidate adds four Go controls; expected totals are318Linux (317Go plus one shell) and313Windows, with every old identity preserved.

Each reproduction and correction is retained in Git with its parent. Use ordinary revert commits after dependency assessment and appropriate QA; do not reset protected history, rewrite independent source fixtures or conceal failed review/QA. Candidate/new-source and distinct accepted-push evidence remain different states.

Read [active journal](SYSTEMATIC_REVIEW_CHANGELOG_20261014.md) and [handover](SESSION_HANDOVER.md). PR142 provides the actual live candidate and QA outcome. After Section1is verified, stop for user re-evaluation; Sections2–4 and overall release work remain open.

## Reproduction log integrity

The full returned UTF-8 logs were hashed for recovery, retaining platform line endings. Character counts and UTF-8 byte counts are different where BOM/non-ASCII appears.

| Reproduction | Job | Characters | SHA-256 of returned UTF-8 |
| --- | --- | ---: | --- |
| Original cancellation Linux | 115844262674 | 64338 | `d9f6408fc815bc109df14de24cd9eb60f950ac55bc07e7e248a8e53c939cf0a5` |
| Original cancellation Windows | 115844263005 | 51809 | `d1d4ebbe88f4d64fa8625e1d63cbd5626ec20ecc742e3991f0a913530cf165a5` |
| Initial weaker paging Linux | 115848958752 | 63105 | `d6185362816586d8a6c977011ac4f61142564792bf82714058a81a3237cdf8ff` |
| Initial weaker paging Windows | 115848959368 | 50413 | `7d7e742e441672b60dffccc739a7b7ca0898d3134146f031040f89fdc4a39881` |
| Parent priority Linux | 115859094976 | 62888 | `8d95dfbf39cfa126f59d7e96e29b2c30265998f7b905c787cd52346e41c8f4b3` |
| Parent priority Windows | 115859094769 | 50901 | `a29e345092b5f5c4d8b59a2d902d51b0bf3fd0447eb6fbdb6de4100316971f5b` |
| ARG status Linux | 115864719992 | 65138 | `bb5f8be4e1f1b64f6e47daafb3861555243f3a7ba4c94ae6d03e86bb4180b1c4` |
| ARG status Windows | 115864720321 | 51286 | `f7bbe97fd233d2eedcfe60d0d404ecc6fcad8d355f4c3e87cdddf29d5deea7b7` |
