# Section 2 QA correction: Cost local response name, 2026-10-14 UTC

Current candidate7ce5cf2273edc6e95183ad3a8eabfa9fe6a11c38 FAILED native38645472872 on both hosts at the Section2control healthy Cost baseline. The new HTTP response local collided with the existing decoded queryResult local. This is QA-S2-01, an introduced compile defect, not an expected compiling mutation rejection or one of the baseline product findings. No passing acceptance is claimed.

Full returned Linux116131114362 log32401characters SHA-256e1b2405ac4e98a4345c38927190217527a3e718a1f5c3ab8a69e37ac9eb7f363; Windows116131113918 log32221characters SHA-256de5bc2f083ffe07634230daf35c3a7a80b6f25213fdd3e72833c57ca5bd53781. First two Section2controls passed/restored on both; later controls and full QA were not reached. Source38645472929 passed separately, with23captures/315chunks verified against retained baseline blobs; it does not clear native failure.

Correction: rename only the HTTP local to httpResponse, retain decoded response, and update the exact Cost-status mutation anchor. Fresh full candidate QA and exact-head review are required. Read live PR142/current refs for actual successor; do not transfer previous results. The Section2checkpoint remains pending and Section3has not started.

---

# Section 2 systemic review, started 2026-10-14

## Current state and next action

All manual source/test review and both-host independent reproductions are complete. Six P2 finding groups are corrected in separate commits; full corrected-candidate QA and automated review are pending. This is not a completed Section 2 checkpoint, protected acceptance, whole-codebase clearance or release approval.

Baseline: `7911764ebd9fff9968589df8ee4c953a16a03d1c`, tree `1f09d2b00264f6efe30531625a8e097eb00bc1f3`. Production correction ancestry ends at `cece5532c1ac79fd2f4964b1594bbafc06a0c61d`. Read [live PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142) and current refs for the actual candidate with helpers/documentation and its exact-head results. The accepted branch remains `a6fee427e19feeef28691737ea9bc5026086efa0`.

NEXT: inspect complete Linux/Windows/source QA and review on the latest corrected candidate, address genuine failures, then reconcile coverage and publish the Section 2 checkpoint. Do not replay completed reproductions or fixes. Stop for user re-evaluation before Section 3.

## Coverage

All 51 baseline Go files across 13 packages received full source/test review: 25 production and 26 tests. Returned full UTF-8 contents were independently Git-blob hashed against all baseline entries. The candidate adds six test files inside this section (57 current files), one coordinator integration test in Section 1, one Python control helper and two CI execution steps. The immediate helper/workflow paths are reviewed; broad tooling remains Section 4.

| Package | Production | Baseline tests | State |
| --- | ---: | ---: | --- |
| `internal/rules` | 5 | 4 | Reviewed; corrected-candidate QA pending |
| `internal/scanners` | 1 | 1 | Reviewed; corrected-candidate QA pending |
| `internal/findings` | 2 | 2 | Reviewed; corrected-candidate QA pending |
| `internal/stages` | 3 | 3 | Reviewed; corrected-candidate QA pending |
| `internal/result` | 3 | 3 | Reviewed; corrected-candidate QA pending |
| `internal/gate` | 1 | 1 | Reviewed; corrected-candidate QA pending |
| `internal/advisor` | 2 | 3 | Reviewed; corrected-candidate QA pending |
| `internal/arcsql` | 1 | 1 | Reviewed; corrected-candidate QA pending |
| `internal/cost` | 2 | 4 | Reviewed; corrected-candidate QA pending |
| `internal/defender` | 1 | 1 | Reviewed; corrected-candidate QA pending |
| `internal/diagnostics` | 2 | 1 | Reviewed; corrected-candidate QA pending |
| `internal/policy` | 1 | 1 | Reviewed; corrected-candidate QA pending |
| `internal/skus` | 1 | 1 | Reviewed; corrected-candidate QA pending |

No package outside this section is credited as newly fully reviewed. Audit remains 1/4 sections complete and 9/41 finally verified packages until corrected QA/reconciliation completes. These are counts, not a weighted project percentage.

## Confirmed findings and corrections

### SR-F04: Caller cancellation can be hidden by a final optional stage error, P2

Observed: Runner only inspected caller context after a nil task error. Canceled/deadline callers were recorded as partial/stage_failed. The coordinator regression returned partial/nil error when final Cost canceled the caller and returned a plain error.

Correction: Always inspect caller context after task execution and give its terminal error priority. Preserve records/warnings and ordinary optional-error behavior.

Fix commit `3c5860adee63e710e1d7015dbdc631522f276f42`, parent `ac08820a66cdcbf55dfe39314f03b2138b11a48c`. State: corrected candidate, final verification pending.

### SR-F05: Scanner accessors expose mutable global resource-type slices, P2

Observed: ByKey copied Service records but shared ResourceTypes. Mutating the returned Storage slice through ByKey or All changed a fresh lookup and could change later scan scope.

Correction: Deep-copy each service's ResourceTypes in ByKey; All inherits owned slices. Registry content and ordering stay unchanged.

Fix commit `b871736a17b7d3375fc1be424b23655e9bcef044`, parent `3c5860adee63e710e1d7015dbdc631522f276f42`. State: corrected candidate, final verification pending.

### SR-F06: Core worker collectors lose terminal caller context, P2

Observed: Cost and Diagnostics accepted a valid-looking final response after caller cancellation, or returned a plain transport error that masked cancellation. Diagnostics inferred a missing-setting finding after interruption.

Correction: Check caller context before work and after workers join, check worker context after receiving jobs, and prioritize the caller over plain errors. Cost checks again after sorting. Diagnostics checks inference boundaries and discards inferred findings on interruption. Ordinary internal sibling cancellation keeps its initiating error.

Fix commit `0c21ac4c08300ee23f7cc9da3651f0819fd9c2b5`, parent `b871736a17b7d3375fc1be424b23655e9bcef044`. State: corrected candidate, final verification pending.

### SR-F07: Diagnostics accepts trailing batch-response content, P2

Observed: A valid first batch object followed by {}, null, junk or an incomplete object produced findings with nil error. Empty/whitespace suffix and one-close controls passed.

Correction: Require EOF after one JSON response using the same decoder; preserve underlying read-error/deadline causes with %w. Keep the existing body-close defer and valid whitespace behavior.

Fix commit `fe39724654dd944297b3ff73de8b0ab5f733f34f`, parent `6c897ef200c381ab8292e1129940da07dd084bd9`. State: corrected candidate, final verification pending.

### SR-F08: Diagnostics admits setting evidence outside the requested batch, P2

Observed: A one-resource batch returning another batch's setting was admitted and produced inferred findings with nil error. The owned/case-insensitive control passed.

Correction: Bind each parseable setting parent to the normalized requested IDs inside scanBatch. Reject foreign evidence before inferring findings. This proves batch membership, not universal per-subresponse correlation.

Fix commit `d12baf35aaaa70e204405594a06d3aed42c29488`, parent `fe39724654dd944297b3ff73de8b0ab5f733f34f`. State: corrected candidate, final verification pending.

### SR-F09: Cost accepts undocumented completion statuses, P2

Observed: Using the actual shared HTTPClient and mocked responses, 201/202/203/205/206 with valid-looking cost rows were admitted as completed queries. The corrected fixture's 200/204 controls passed and every case reached the transport.

Correction: Retain and inspect the HTTP response. Accept only the pinned API's documented 200/204 statuses; reject missing or other responses. Keep valid no-content and existing error/skip handling.

Fix commit `cece5532c1ac79fd2f4964b1594bbafc06a0c61d`, parent `d12baf35aaaa70e204405594a06d3aed42c29488`. State: corrected candidate, final verification pending.


SR-F06 covers two worker collector implementations under one terminal-context finding group. Seven distinct compiling counterfactual controls cover the six groups; findings are not a count of all possible failure paths or the unreviewed codebase.

## Reproduction evidence and fixture correction

Final regression-only head `ac08820a66cdcbf55dfe39314f03b2138b11a48c` retains unchanged production blobs and all independent expectations. Native run [38642216627](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38642216627) compiled and failed only Linux race tests and Windows Go packages at the declared regressions. Jobs 116120286515 and 116120287313 both completed/failure. Formatting/build and all earlier steps passed; no build failure or data-race warning was observed. Source run 38642216664 completed/success, separate from failing native reproduction.

Both hosts confirmed all six groups, including coordinator partial/nil cancellation impact, every Cost status reaching the transport, and healthy 200/204/ordinary-error/whitespace/owned-setting controls. Source captures, pins, dependencies, schemas and healthy KQL remain unchanged.

Initial head `2ea0b11fbbaf142d4c4ff7db1d20f5693ffc13bc`, native 38640793764, confirmed the first five groups on both hosts. Its Cost status test was insufficient: a 1-second attempt limit was below the global Cost limiter's 5-second spacing, so 200/204 failed before transport. No status repro is credited from that trial. Corrected the fixture to a 10-second attempt limit and required calls==1 in every case; production throttling/timeouts and expected status outcomes were untouched. The corrected-fixture run reproduced actual bad-status admission on both hosts.

| Returned log | Job | Characters | SHA-256 of exact returned UTF-8 |
| --- | --- | ---: | --- |
| Initial Linux | 116115891040 | 171089 | `5f437158028648c5121440dec5ac4b9969d4046f043f707ea75803f3a2a553be` |
| Initial Windows | 116115890658 | 158522 | `81ee811fccf7db58ffd3ddbec3509fed39902b3f2d38b95416c53c32fbef9cad` |
| Corrected-fixture Linux | 116120286515 | 175182 | `3bf50be5f6b9769a233f457bab828d0c7b5eabbfc61e02e2a1897bb0f81c3c8c` |
| Corrected-fixture Windows | 116120287313 | 161030 | `bab7539a44bf6a2b4124e67c5ba4e5fffda2544377a445fa46d9542c8a3963cc` |

Hashes retain original line endings/BOM. Earlier native predecessors superseded by normal workflow concurrency are not accepted evidence. No blind retries or gate relaxation occurred.

## Primary and pinned-source checks

Checked during this review:

- Exact 282 registry key/name/resource-type literal tokens match source after removing target field labels. 87 scanner keys and the separate Redis entries are preserved.
- Advisor, Arc SQL, Policy and both Defender queries match the actual source after per-line indentation/trailing-format normalization; no semantic query change.
- Diagnostics support table: all 74 entries, 41 dedicated recommendation IDs, text and links match exactly, including supported types without dedicated recommendations.
- known_skus.yaml full target content independently Git-blob hashed and matches source `a2d97a60ec445ce4023a1a5a9b6c4dca76e31f5a`.
- Summary source `92c6214347fcb030eec527fd52c64300a68722c0`: normalized recommendation/resource dedup, SLA exclusion and recommendation-impact counts preserved. Impact totals count recommendation/resource pairs, not unique estate resources.
- Cost source `b0d27e8c0f74b92f595520d82fd06793c7e27523`: previous completed UTC month, ActualCost/Sum/ServiceName and positional row mapping preserved. Subscription names and unconsumed-continuation failure are established target corrections.
- Source Arc SQL also uses string vcores despite toint projection. The existing numeric-row skip/warning is an explicit semantic-compatibility boundary, not a new divergence.
- Microsoft [DiagnosticSettings List 2021-05-01-preview](https://learn.microsoft.com/en-us/rest/api/monitor/diagnostic-settings/list?view=rest-monitor-2021-05-01-preview): resourceUri scope, HTTP 200 and parent-qualified setting ID verified.
- Official [Cost SDK v1.1.1 query client](https://github.com/Azure/azure-sdk-for-go/blob/sdk/resourcemanager/costmanagement/armcostmanagement/v1.1.1/sdk/resourcemanager/costmanagement/armcostmanagement/query_client.go), blob `fcce0881a833b237333f35de9dba97976b36d153`, generated from API 2021-10-01: successful 200/204 contract verified. Versioned Learn page was inaccessible; the latest page resolved 2026-06-01 and was not treated as 2021-10-01 evidence.
- Defender REST assessment documentation returned 2020-01-01 and names azurePortalUri; that does not establish the ARG azurePortal wire representation used by the pinned query. No portal-link defect or field rename is inferred from a different API projection; live schema checks remain separate.

## Architecture and residual boundaries

The implementation keeps typed collector interfaces in consuming packages, separates collection from canonical/report projections, returns explicit errors from libraries, bounds worker concurrency and owns mutable request/result data at composition boundaries. Stage sequencing, source precedence, selected scope, scanner pruning, SLA/dedup/severity and plugin health remain intact. The registry copy repair closes a demonstrated global-state ownership leak.

Successful canonical results retain copied nested slices/maps, deterministic key ordering and UTC timestamps. Optional failures preserve other datasets; critical failures and terminal caller context stop execution. Cost's materialized Post response body is consumed/closed by the shared client. Diagnostics owns PostStream body closure; framing success now includes EOF.

Known limits are explicit and are not claimed fixed by this review:

- Source-compatible Graph decoders tolerate structurally valid empty/null rows; strict business-row identity/completeness is not universal.
- Arc SQL numeric vcores may omit valid rows with warnings. Changing that behavior needs a deliberate compatibility correction and evidence.
- Policy/Defender first-row composite dedup is not a latest-assignment/history guarantee. Defender namespace-level type/name/link projections retain source behavior.
- Diagnostics presence does not certify enabled categories or destinations. Non-200 subrequests may retain source-compatible missing-setting findings with uncertainty warnings. Batch membership does not prove a complete named response/request correlation contract.
- Core Cost is subscription-wide; RG/resource/tag filters do not narrow billing totals. Positional row and value formatting are source-compatible.
- Per-operation deadlines, bounded workers, cycle detection and opt-in assessment timeout do not establish a universal whole-scan byte/row/work/time limit or process hard kill. Same-origin Advisor pagination and core collector/ARG volume limits require representative load/AR-02 and further contract decisions.
- Canonical core library builders assume validated domain input; plugin extension validators do not claim validation of every manually constructed core result.
- Live scope/RBAC/estate membership, fresh OS/operator and source equivalence, module-only advisory disposition, human/license/security/Gate004 and release work remain open. No Azure operation occurred.

## Controls, change isolation and rollback

`scripts/tests/section2-regressions.py` runs seven independent healthy-baseline, compiling mutation rejection, exact raw-byte restoration and restored-baseline checks on both hosts. Mutations cover Runner context, registry ownership, Cost/Diagnostics caller priority, Diagnostics EOF and batch scope, and Cost completion status. Mutation failure must include its declared diagnostic and cannot be a build failure. New helper AST/imports and all seven exact anchors were checked locally; local Go is unavailable, so native execution remains authoritative.

Full candidate control totals expected: Linux 325 (324 Go-labeled controls plus one shell bootstrap control), Windows 320 Go-labeled controls. Every previous identity must remain present; counts alone are insufficient.

- SR-F04: revert `3c5860adee63e710e1d7015dbdc631522f276f42` against parent `ac08820a66cdcbf55dfe39314f03b2138b11a48c`; paths `internal/stages/runner.go`.
- SR-F05: revert `b871736a17b7d3375fc1be424b23655e9bcef044` against parent `3c5860adee63e710e1d7015dbdc631522f276f42`; paths `internal/scanners/registry.go`.
- SR-F06: revert `0c21ac4c08300ee23f7cc9da3651f0819fd9c2b5` against parent `b871736a17b7d3375fc1be424b23655e9bcef044`; paths `internal/cost/scanner.go`, `internal/diagnostics/scanner.go`.
- SR-F07: revert `fe39724654dd944297b3ff73de8b0ab5f733f34f` against parent `6c897ef200c381ab8292e1129940da07dd084bd9`; paths `internal/diagnostics/scanner.go`.
- SR-F08: revert `d12baf35aaaa70e204405594a06d3aed42c29488` against parent `fe39724654dd944297b3ff73de8b0ab5f733f34f`; paths `internal/diagnostics/scanner.go`.
- SR-F09: revert `cece5532c1ac79fd2f4964b1594bbafc06a0c61d` against parent `d12baf35aaaa70e204405594a06d3aed42c29488`; paths `internal/cost/scanner.go`.

Use ordinary targeted revert commits after checking subsequent dependencies. SR-F06/F07/F08 share Diagnostics; SR-F06/F09 share Cost. Helper anchors and regression expectations depend on the fixes, so restoring old defects will intentionally fail them and requires an explicit rollback disposition. Do not reset protected history, rewrite independent fixtures, or claim healthy rollback without QA. No data migration or Azure rollback exists.

Self-review caught a batch guard placed at the first enable-map assignment in an unpublished draft, outside its owning function. It was corrected inside scanBatch before any ref advance. Prepared `7838ec5c`/`4ffc43c9` objects were abandoned and are not task-branch ancestry or accepted evidence; final SR-F08/SR-F09 above are authoritative. This draft issue is not an additional product defect.

Supplemental ordinary trailing-read error and context.DeadlineExceeded cases were declared in `6c897ef200c381ab8292e1129940da07dd084bd9` before SR-F07's framing correction. They verify cause retention/body closure on the corrected candidate; they are not separately credited as an old-native reproduction.

Read [change journal](SYSTEMATIC_REVIEW_CHANGELOG_20261014.md) and [handover](SESSION_HANDOVER.md). After final corrected-candidate QA/review and coverage reconciliation, publish the Section 2 checkpoint and stop for user re-evaluation. Sections 3–4 remain approximately 6–11 active review hours plus material corrections/repeat QA; overall release ETA and weighted project percentage remain unverified.

## Baseline file integrity and coverage

| File | Baseline Git blob | Review state |
| --- | --- | --- |
| `internal/advisor/metadata.go` | `97dd140e5f0bf49b1377ed4c7d713cb6cfdf7be2` | Full source/test reviewed; final QA pending |
| `internal/advisor/metadata_http_test.go` | `b8a81318ef8f86e7a5c2748b761a21f5fe62ba5d` | Full source/test reviewed; final QA pending |
| `internal/advisor/pagination_context_test.go` | `e9afac187046b6addde7d5e6100a58f22694fd07` | Full source/test reviewed; final QA pending |
| `internal/advisor/scanner.go` | `a636cd01e963469ac1a959da415bfc1c893f944c` | Full source/test reviewed; final QA pending |
| `internal/advisor/scanner_test.go` | `b7a0d87604cc26613dfe71514bbd751e3dc5e2c2` | Full source/test reviewed; final QA pending |
| `internal/arcsql/scanner.go` | `d4f9f0230c26842e02d62fbda9730132a050c59f` | Full source/test reviewed; final QA pending |
| `internal/arcsql/scanner_test.go` | `4cd0d0130a7652adc66ba6c0881b66611524f5ad` | Full source/test reviewed; final QA pending |
| `internal/cost/continuation_test.go` | `2a2f9e4ce072da9d2ef7232886dd7e7515a4aaf5` | Full source/test reviewed; final QA pending |
| `internal/cost/http_contract_test.go` | `2d93b9d1801eef11b9c9545c172e2c94c3e8149b` | Full source/test reviewed; final QA pending |
| `internal/cost/period.go` | `ee7f36b63acef11043e2f4e42722499ccb658426` | Full source/test reviewed; final QA pending |
| `internal/cost/period_test.go` | `192d0e94f67b3631d58e3372240b3b6fa0796b8f` | Full source/test reviewed; final QA pending |
| `internal/cost/scanner.go` | `aaa7d4f85d576e7e76f2965513c23d53ea5f11ef` | Full source/test reviewed; final QA pending |
| `internal/cost/scanner_test.go` | `f4e7a4b3d908c37c41cc67d2be78643f1dc0d0fc` | Full source/test reviewed; final QA pending |
| `internal/defender/scanner.go` | `7e843f0e826f5c7a1b928eb22c4de525c6bd48f5` | Full source/test reviewed; final QA pending |
| `internal/defender/scanner_test.go` | `a72483ee1c6c6f0f4aad402328e8c2d58cf89d9b` | Full source/test reviewed; final QA pending |
| `internal/diagnostics/recommendations.go` | `9fb8f241d68502dd425ebb6e2f40f8414dcc5ba6` | Full source/test reviewed; final QA pending |
| `internal/diagnostics/scanner.go` | `b9c3cb25c28b1282c06866344effe0590539b13e` | Full source/test reviewed; final QA pending |
| `internal/diagnostics/scanner_test.go` | `7bfcb1ac7a56ad33a14ec79a91de08ff78951bd3` | Full source/test reviewed; final QA pending |
| `internal/findings/applicability.go` | `366dde9e246f1a8e92fd58f975dec70652ea98b2` | Full source/test reviewed; final QA pending |
| `internal/findings/applicability_test.go` | `068c908fa55bbcff01202294876dcf3d8198a0c3` | Full source/test reviewed; final QA pending |
| `internal/findings/summary.go` | `e4fc5b01bd25c668a85c3b77584b1a1d3223823a` | Full source/test reviewed; final QA pending |
| `internal/findings/summary_test.go` | `62cad5632c5402f649c1a4d88e7d7f203f584924` | Full source/test reviewed; final QA pending |
| `internal/gate/gate.go` | `35e3d887aff375b3ca63b673dfd090641dc5c4b4` | Full source/test reviewed; final QA pending |
| `internal/gate/gate_test.go` | `45b96ba6241d71a11f7437960f424e0e68b52d4d` | Full source/test reviewed; final QA pending |
| `internal/policy/scanner.go` | `1a549afefd7be50e0262336f733d8d5f76901a3d` | Full source/test reviewed; final QA pending |
| `internal/policy/scanner_test.go` | `74d77a818ba3d0164d141b3ee1b0656503d7e53e` | Full source/test reviewed; final QA pending |
| `internal/result/plugin_tables.go` | `9703654451984758e91b8c962ae5626f31e79b96` | Full source/test reviewed; final QA pending |
| `internal/result/plugin_tables_test.go` | `4a52a5a19bcc8658a2b8aa5b1832a3d89ca26a07` | Full source/test reviewed; final QA pending |
| `internal/result/result.go` | `f099c293a44707535ffe624c201eddcf240847f0` | Full source/test reviewed; final QA pending |
| `internal/result/result_test.go` | `bac7ad307b4c7f79d19831e1229d36293d9ef97e` | Full source/test reviewed; final QA pending |
| `internal/result/scope.go` | `a65bc4f42d7c973cb399fd3a06c7f30019bfb7fa` | Full source/test reviewed; final QA pending |
| `internal/result/scope_test.go` | `ed3ae54acdd9b33e90ab318836810fd3ac1523c4` | Full source/test reviewed; final QA pending |
| `internal/rules/catalog.go` | `e2d0a9ba82323ccb01dd837d14f7c25cf5d22c6d` | Full source/test reviewed; final QA pending |
| `internal/rules/catalog_test.go` | `e67a3d090bead57e968e2464c2de205b16494f12` | Full source/test reviewed; final QA pending |
| `internal/rules/embedded.go` | `5c3e0ce214e8d1571b9d1977b8f4a7f1d0dd326f` | Full source/test reviewed; final QA pending |
| `internal/rules/embedded_test.go` | `82611bcda9beb4d44c9a817f79379ad94d4591d4` | Full source/test reviewed; final QA pending |
| `internal/rules/inspection.go` | `97a998b23916134245f23fd15f0de1d4134897e7` | Full source/test reviewed; final QA pending |
| `internal/rules/inspection_test.go` | `8ff8d0177243a95bab178cabb8b4d4bebebe125f` | Full source/test reviewed; final QA pending |
| `internal/rules/provenance.go` | `0315d420dfc05b39cf03a6e782769d630d66ef8a` | Full source/test reviewed; final QA pending |
| `internal/rules/yaml_plugins.go` | `5a43f569f3ceceab829b9a26dd1d46c6732b5076` | Full source/test reviewed; final QA pending |
| `internal/rules/yaml_plugins_test.go` | `d02f2967b009aa549c422017ce8292b3a851b9e8` | Full source/test reviewed; final QA pending |
| `internal/scanners/registry.go` | `24512b10be30824c62f11bf83bc66e0b04e11eb6` | Full source/test reviewed; final QA pending |
| `internal/scanners/registry_test.go` | `16f8f679331a50d1346c10292bdb72e73674fb71` | Full source/test reviewed; final QA pending |
| `internal/skus/skus.go` | `396cc4b7172f536b08d91355f8f64e9e4f0b97a1` | Full source/test reviewed; final QA pending |
| `internal/skus/skus_test.go` | `6d1b6e6d89e4ac3a47c49fa54fbf34ec31db79dc` | Full source/test reviewed; final QA pending |
| `internal/stages/config.go` | `307eabbb30b2b988cdc6e2fc750b3bce55d993d9` | Full source/test reviewed; final QA pending |
| `internal/stages/config_test.go` | `804ccc644facc788b357fb36e54d06fe21c50e05` | Full source/test reviewed; final QA pending |
| `internal/stages/options.go` | `d0138077d1784f5e82cde165af8eb62b410e55d2` | Full source/test reviewed; final QA pending |
| `internal/stages/options_test.go` | `8245b9a14a468022ded93bdf62ab23538ecd6e16` | Full source/test reviewed; final QA pending |
| `internal/stages/runner.go` | `7f3ec14d2f7b150619446984ddd9659b9761d8e4` | Full source/test reviewed; final QA pending |
| `internal/stages/runner_test.go` | `6e6b6f6bbeb0db79a0b4d2c1a21cbd5ef46c3be7` | Full source/test reviewed; final QA pending |
