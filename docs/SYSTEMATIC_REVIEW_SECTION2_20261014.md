# Section 2 active systemic review, 2026-10-14

User authorized continuation at23:23Europe/Berlin on2026-10-14. Section1 remains complete; Section2 is now in progress, not a completed checkpoint.

Baseline/task head: `7911764ebd9fff9968589df8ee4c953a16a03d1c`, tree `1f09d2b00264f6efe30531625a8e097eb00bc1f3`. Draft PR142 remains unmerged; protected baseline is `a6fee427e19feeef28691737ea9bc5026086efa0`. Section1 documentation successor native38574597736 and source38574597703 now completed/success; fresh status and both native job steps checked before continuation.

Scope is all13assigned package directories and51Go files (25production,26tests): rules, scanners, findings, stages, result, gate, advisor, arcsql, cost, defender, diagnostics, policy and skus. Full source/test review and coverage reconciliation are pending. Read dependency contracts as needed without crediting outside modules as fully reviewed.

Review priorities: canonical finding/applicability/severity/health semantics; recommendation catalog precedence and extension inputs; stage cancellation/optional failures; collector scope, pagination, row/empty/malformed/completeness contracts; provider-specific endpoints and request/memory budgets; owned deterministic results and semantic pinned-source compatibility.

Findings: no new confirmed defect yet. SR-S04 is suspected: caller cancellation/deadline during an optional final task can be masked by a plain task error, leaving stage_failed/partial rather than terminal-context/failed. Independent regression declared before production correction; healthy ordinary-error control preserves partial semantics. Suspicions require independent failing reproduction or primary-source contract verification before classification/correction. Keep source captures/pins/dependencies immutable unless an evidenced, separately reviewed change is necessary.

NEXT: fetch exact baseline file contents, review each production module with all assigned tests, double-check uncertain API/source contracts, record findings and exact repro/fix/rollback evidence. Publish durable intermediate checkpoints when findings or completed package coverage change. Stop after Section2corrective QA/reconciliation for the user's next re-evaluation; do not startSection3automatically.

Audit remains1/4sections complete and9/41finally reviewed packages. Section2planning budget3–4active review hours; Sections2–4overall9–15active hours plus material fixes and repeated QA. This is effort, not a background schedule or release ETA. No Azure action, laptop access or live parity is implied.

## Baseline file coverage

| File | Git blob | State |
| --- | --- | --- |
| `internal/advisor/metadata.go` | `97dd140e5f0bf49b1377ed4c7d713cb6cfdf7be2` | Source/test reviewed; final correction QA pending |
| `internal/advisor/metadata_http_test.go` | `b8a81318ef8f86e7a5c2748b761a21f5fe62ba5d` | Source/test reviewed; final correction QA pending |
| `internal/advisor/pagination_context_test.go` | `e9afac187046b6addde7d5e6100a58f22694fd07` | Source/test reviewed; final correction QA pending |
| `internal/advisor/scanner.go` | `a636cd01e963469ac1a959da415bfc1c893f944c` | Source/test reviewed; final correction QA pending |
| `internal/advisor/scanner_test.go` | `b7a0d87604cc26613dfe71514bbd751e3dc5e2c2` | Source/test reviewed; final correction QA pending |
| `internal/arcsql/scanner.go` | `d4f9f0230c26842e02d62fbda9730132a050c59f` | Source/test reviewed; final correction QA pending |
| `internal/arcsql/scanner_test.go` | `4cd0d0130a7652adc66ba6c0881b66611524f5ad` | Source/test reviewed; final correction QA pending |
| `internal/cost/continuation_test.go` | `2a2f9e4ce072da9d2ef7232886dd7e7515a4aaf5` | Source/test reviewed; final correction QA pending |
| `internal/cost/http_contract_test.go` | `2d93b9d1801eef11b9c9545c172e2c94c3e8149b` | Source/test reviewed; final correction QA pending |
| `internal/cost/period.go` | `ee7f36b63acef11043e2f4e42722499ccb658426` | Source/test reviewed; final correction QA pending |
| `internal/cost/period_test.go` | `192d0e94f67b3631d58e3372240b3b6fa0796b8f` | Source/test reviewed; final correction QA pending |
| `internal/cost/scanner.go` | `aaa7d4f85d576e7e76f2965513c23d53ea5f11ef` | Source/test reviewed; final correction QA pending |
| `internal/cost/scanner_test.go` | `f4e7a4b3d908c37c41cc67d2be78643f1dc0d0fc` | Source/test reviewed; final correction QA pending |
| `internal/defender/scanner.go` | `7e843f0e826f5c7a1b928eb22c4de525c6bd48f5` | Source/test reviewed; final correction QA pending |
| `internal/defender/scanner_test.go` | `a72483ee1c6c6f0f4aad402328e8c2d58cf89d9b` | Source/test reviewed; final correction QA pending |
| `internal/diagnostics/recommendations.go` | `9fb8f241d68502dd425ebb6e2f40f8414dcc5ba6` | Source/test reviewed; final correction QA pending |
| `internal/diagnostics/scanner.go` | `b9c3cb25c28b1282c06866344effe0590539b13e` | Source/test reviewed; final correction QA pending |
| `internal/diagnostics/scanner_test.go` | `7bfcb1ac7a56ad33a14ec79a91de08ff78951bd3` | Source/test reviewed; final correction QA pending |
| `internal/findings/applicability.go` | `366dde9e246f1a8e92fd58f975dec70652ea98b2` | Source/test reviewed; final correction QA pending |
| `internal/findings/applicability_test.go` | `068c908fa55bbcff01202294876dcf3d8198a0c3` | Source/test reviewed; final correction QA pending |
| `internal/findings/summary.go` | `e4fc5b01bd25c668a85c3b77584b1a1d3223823a` | Source/test reviewed; final correction QA pending |
| `internal/findings/summary_test.go` | `62cad5632c5402f649c1a4d88e7d7f203f584924` | Source/test reviewed; final correction QA pending |
| `internal/gate/gate.go` | `35e3d887aff375b3ca63b673dfd090641dc5c4b4` | Source/test reviewed; final correction QA pending |
| `internal/gate/gate_test.go` | `45b96ba6241d71a11f7437960f424e0e68b52d4d` | Source/test reviewed; final correction QA pending |
| `internal/policy/scanner.go` | `1a549afefd7be50e0262336f733d8d5f76901a3d` | Source/test reviewed; final correction QA pending |
| `internal/policy/scanner_test.go` | `74d77a818ba3d0164d141b3ee1b0656503d7e53e` | Source/test reviewed; final correction QA pending |
| `internal/result/plugin_tables.go` | `9703654451984758e91b8c962ae5626f31e79b96` | Source/test reviewed; final correction QA pending |
| `internal/result/plugin_tables_test.go` | `4a52a5a19bcc8658a2b8aa5b1832a3d89ca26a07` | Source/test reviewed; final correction QA pending |
| `internal/result/result.go` | `f099c293a44707535ffe624c201eddcf240847f0` | Source/test reviewed; final correction QA pending |
| `internal/result/result_test.go` | `bac7ad307b4c7f79d19831e1229d36293d9ef97e` | Source/test reviewed; final correction QA pending |
| `internal/result/scope.go` | `a65bc4f42d7c973cb399fd3a06c7f30019bfb7fa` | Source/test reviewed; final correction QA pending |
| `internal/result/scope_test.go` | `ed3ae54acdd9b33e90ab318836810fd3ac1523c4` | Source/test reviewed; final correction QA pending |
| `internal/rules/catalog.go` | `e2d0a9ba82323ccb01dd837d14f7c25cf5d22c6d` | Source/test reviewed; final correction QA pending |
| `internal/rules/catalog_test.go` | `e67a3d090bead57e968e2464c2de205b16494f12` | Source/test reviewed; final correction QA pending |
| `internal/rules/embedded.go` | `5c3e0ce214e8d1571b9d1977b8f4a7f1d0dd326f` | Source/test reviewed; final correction QA pending |
| `internal/rules/embedded_test.go` | `82611bcda9beb4d44c9a817f79379ad94d4591d4` | Source/test reviewed; final correction QA pending |
| `internal/rules/inspection.go` | `97a998b23916134245f23fd15f0de1d4134897e7` | Source/test reviewed; final correction QA pending |
| `internal/rules/inspection_test.go` | `8ff8d0177243a95bab178cabb8b4d4bebebe125f` | Source/test reviewed; final correction QA pending |
| `internal/rules/provenance.go` | `0315d420dfc05b39cf03a6e782769d630d66ef8a` | Source/test reviewed; final correction QA pending |
| `internal/rules/yaml_plugins.go` | `5a43f569f3ceceab829b9a26dd1d46c6732b5076` | Source/test reviewed; final correction QA pending |
| `internal/rules/yaml_plugins_test.go` | `d02f2967b009aa549c422017ce8292b3a851b9e8` | Source/test reviewed; final correction QA pending |
| `internal/scanners/registry.go` | `24512b10be30824c62f11bf83bc66e0b04e11eb6` | Source/test reviewed; final correction QA pending |
| `internal/scanners/registry_test.go` | `16f8f679331a50d1346c10292bdb72e73674fb71` | Source/test reviewed; final correction QA pending |
| `internal/skus/skus.go` | `396cc4b7172f536b08d91355f8f64e9e4f0b97a1` | Source/test reviewed; final correction QA pending |
| `internal/skus/skus_test.go` | `6d1b6e6d89e4ac3a47c49fa54fbf34ec31db79dc` | Source/test reviewed; final correction QA pending |
| `internal/stages/config.go` | `307eabbb30b2b988cdc6e2fc750b3bce55d993d9` | Source/test reviewed; final correction QA pending |
| `internal/stages/config_test.go` | `804ccc644facc788b357fb36e54d06fe21c50e05` | Source/test reviewed; final correction QA pending |
| `internal/stages/options.go` | `d0138077d1784f5e82cde165af8eb62b410e55d2` | Source/test reviewed; final correction QA pending |
| `internal/stages/options_test.go` | `8245b9a14a468022ded93bdf62ab23538ecd6e16` | Source/test reviewed; final correction QA pending |
| `internal/stages/runner.go` | `7f3ec14d2f7b150619446984ddd9659b9761d8e4` | Source/test reviewed; final correction QA pending |
| `internal/stages/runner_test.go` | `6e6b6f6bbeb0db79a0b4d2c1a21cbd5ef46c3be7` | Source/test reviewed; final correction QA pending |

## SR-S04 reproduction preparation

New internal/stages/runner_terminal_context_test.go exercises canceled/deadline caller context during an optional final task that returns a plain transport error, followed only by a disabled task. Expected terminal health is failed with the appropriate context code. An ordinary-error control must remain partial/stage_failed. Returned record/warning evidence must remain retained. Production runner remains unchanged at baseline blob 7f3ec14d2f7b150619446984ddd9659b9761d8e4. Pending native reproduction; do not credit this as a confirmed failure until full job evidence is inspected.

## Additional independent reproduction preparation

SR-S05: scanner ByKey returns a shallow Service copy, sharing ResourceTypes with the global registry. Independent expected Storage mapping and both ByKey/All mutation paths are declared before correction; private registry access is used only to restore the old implementation's state after the negative test, not to derive the expectation. Pending native confirmation.

SR-S06: Cost and Diagnostics worker collectors can lose caller cancellation during the final request. Independent canceling clients return either a valid response or a plain transport error. Cost must return the caller terminal error; Diagnostics must return it before inferring missing-setting findings. These are direct-library postconditions; public stage runner/coordinator have separate interruption guards. Production collectors remain unchanged; native repro pending.

Reproduction-only predecessor6b34a7a3 has a running native gate and successful source gate. It may be superseded by this expanded regression-only head under normal workflow concurrency; canceled predecessors are not counted as confirmed or accepted evidence. No production fix has been made.

## Endpoint response-contract reproduction preparation

SR-S07: Diagnostics stream decoding currently consumes only the first JSON object. Tests declare valid whitespace controls, four trailing-content variants and exact body-close ownership before correction. The pinned source has the same loose decoder, independently read; this is proposed target reliability hardening on malformed responses, not healthy-data divergence.

SR-S08: a valid-looking diagnostic setting from a different requested batch is admitted into the shared enabled map. Declared one-resource/two-batch fixture retains owned/case-insensitive healthy control and requires explicit error on foreign batch evidence before generating findings. Microsoft2021-05-01-preview List is scoped to resourceUri; actual returned ID sample belongs to that parent. Reference also admits parsed IDs without scope binding. No live API misbehavior is claimed.

SR-S09: Cost adapter ignores endpoint-specific HTTP completion status. Official Azure SDK armcostmanagement/v1.1.1 (the pinned source dependency) query_client.go, blobfcce0881a833b237333f35de9dba97976b36d153, explicitly allows only200/204 for API2021-10-01. The generic latest Microsoft Learn page resolves2026-06-01 and versioned2021-10-01 page was inaccessible, so the actual tagged primary SDK was used. Declared tests exercise the shared production HTTPClient/mocked transport:200row and204empty healthy controls,201/202/203/205/206 negative cases. No Azure call.

All are still suspected pending native compiling reproduction. Native predecessors may be canceled by normal concurrency; full evidence on this final regression-only head is authoritative. No production correction yet.

## Completed manual review and contract checks

All51baseline Go files (25production/26tests) received full source/test review. Full returned UTF-8 content was independently hashed as Git blobs and matched every baseline tree entry. Six new regression files were read back byte-for-byte and independently blob-hashed. No section checkpoint completion is credited until reproductions/corrections/final QA are reconciled.

Reviewed invariants: source catalog precedence and exact recommendation identity; external YAML first-name/sorted override rules, strict one-document schema and confined/portable file inputs; scanner pruning/structural scope; SLA exclusion and normalized finding dedup; severity separate from health; ordered stages and critical/optional outcomes; deep-owned canonical/plugin tables; provider-specific mappings/filter/dedup and empty/malformed/error paths.

Pinned-source checks at8e4f0577f3615e6c9014c031bcad079f235369cc:
- All282literal registry key/name/resource-type tokens match target after removing only target struct-field labels.
- Advisor, Arc SQL, Policy and both Defender KQL queries match after trimming each line's indentation/trailing formatting.
- Diagnostics support table74entries/41dedicated recommendation IDs, wording and LearnMore links match exactly; nil-support semantics retained.
- known_skus.yaml target full blob independently hashed and matches original a2d97a60ec445ce4023a1a5a9b6c4dca76e31f5a.
- Summary definitions/normalized recommendation+resource dedup/SLA exclusion/count semantics reviewed against actual original source92c6214347fcb030eec527fd52c64300a68722c0.
- Cost previous completed UTC month, ActualCost/Sum/ServiceName and positional row mapping match original b0d27e8c0f74b92f595520d82fd06793c7e27523. Subscription names and explicit continuation failure are documented target corrections.
- Original Arc SQL decoder also uses string vcores despite query toint; existing numeric-row rejection/warning is a known semantic-compatibility boundary, not a newly introduced divergence.

Primary contract2026-10-14: Microsoft DiagnosticSettings List2021-05-01-preview is resource-scoped with200success and resource-parent setting IDs:
https://learn.microsoft.com/en-us/rest/api/monitor/diagnostic-settings/list?view=rest-monitor-2021-05-01-preview
Official SDK Cost API2021-10-01 completion is200/204, read at actual pinned dependency tagv1.1.1:
https://github.com/Azure/azure-sdk-for-go/blob/sdk/resourcemanager/costmanagement/armcostmanagement/v1.1.1/sdk/resourcemanager/costmanagement/armcostmanagement/query_client.go
Versioned Cost Learn page was inaccessible; latest resolved2026-06-01 was not treated as2021-10-01 evidence.

Residuals requiring explicit later acceptance: Graph-backed collectors tolerate structurally valid empty/null rows as source does; strict business-row identity admission is not universally guaranteed. Policy and Defender composite dedup intentionally retains first source-compatible row; this is not newest-assignment/history assurance. Arc SQL numeric vcores can omit otherwise valid rows with warnings. Diagnostics presence only does not certify enabled categories/destinations and non200subrequests can retain missing-setting findings with explicit uncertainty. Core Cost is subscription-wide and does not honor RG/tag/resource billing narrowing. Same-origin Advisor pagination, overall ARG/Cost/Diagnostics byte/row/work volume and total scan duration have no universal positive cap; per-operation deadlines, bounded workers, cycle rejection and opt-in assessment deadline are different limits. Load/AR-02 and live validation remain open; no arbitrary new limits or healthy-source/schema change is introduced by the audit.

Cross-module trace for SR-S04: Cost is last enabled regular stage when Plugin is disabled. Coordinator returns execution/context error only when Runner reports completeness failed. A plain error masking caller cancellation at that point can therefore yield partial result/nil coordinator error. This is a source trace/inference until integration regression is run. An independently prepared coordinator test will verify the correction; the existing direct Runner regression reproduces the underlying condition.
