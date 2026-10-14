# Section 2 active systemic review, 2026-10-14

User authorized continuation at23:23Europe/Berlin on2026-10-14. Section1 remains complete; Section2 is now in progress, not a completed checkpoint.

Baseline/task head: `7911764ebd9fff9968589df8ee4c953a16a03d1c`, tree `1f09d2b00264f6efe30531625a8e097eb00bc1f3`. Draft PR142 remains unmerged; protected baseline is `a6fee427e19feeef28691737ea9bc5026086efa0`. Section1 documentation successor native38574597736 and source38574597703 now completed/success; fresh status and both native job steps checked before continuation.

Scope is all13assigned package directories and51Go files (25production,26tests): rules, scanners, findings, stages, result, gate, advisor, arcsql, cost, defender, diagnostics, policy and skus. Full source/test review and coverage reconciliation are pending. Read dependency contracts as needed without crediting outside modules as fully reviewed.

Review priorities: canonical finding/applicability/severity/health semantics; recommendation catalog precedence and extension inputs; stage cancellation/optional failures; collector scope, pagination, row/empty/malformed/completeness contracts; provider-specific endpoints and request/memory budgets; owned deterministic results and semantic pinned-source compatibility.

Findings: none newly confirmed at this starting record. Suspicions require independent failing reproduction or primary-source contract verification before classification/correction. Keep source captures/pins/dependencies immutable unless an evidenced, separately reviewed change is necessary.

NEXT: fetch exact baseline file contents, review each production module with all assigned tests, double-check uncertain API/source contracts, record findings and exact repro/fix/rollback evidence. Publish durable intermediate checkpoints when findings or completed package coverage change. Stop after Section2corrective QA/reconciliation for the user's next re-evaluation; do not startSection3automatically.

Audit remains1/4sections complete and9/41finally reviewed packages. Section2planning budget3–4active review hours; Sections2–4overall9–15active hours plus material fixes and repeated QA. This is effort, not a background schedule or release ETA. No Azure action, laptop access or live parity is implied.

## Baseline file coverage

| File | Git blob | State |
| --- | --- | --- |
| `internal/advisor/metadata.go` | `97dd140e5f0bf49b1377ed4c7d713cb6cfdf7be2` | Pending full review |
| `internal/advisor/metadata_http_test.go` | `b8a81318ef8f86e7a5c2748b761a21f5fe62ba5d` | Pending full review |
| `internal/advisor/pagination_context_test.go` | `e9afac187046b6addde7d5e6100a58f22694fd07` | Pending full review |
| `internal/advisor/scanner.go` | `a636cd01e963469ac1a959da415bfc1c893f944c` | Pending full review |
| `internal/advisor/scanner_test.go` | `b7a0d87604cc26613dfe71514bbd751e3dc5e2c2` | Pending full review |
| `internal/arcsql/scanner.go` | `d4f9f0230c26842e02d62fbda9730132a050c59f` | Pending full review |
| `internal/arcsql/scanner_test.go` | `4cd0d0130a7652adc66ba6c0881b66611524f5ad` | Pending full review |
| `internal/cost/continuation_test.go` | `2a2f9e4ce072da9d2ef7232886dd7e7515a4aaf5` | Pending full review |
| `internal/cost/http_contract_test.go` | `2d93b9d1801eef11b9c9545c172e2c94c3e8149b` | Pending full review |
| `internal/cost/period.go` | `ee7f36b63acef11043e2f4e42722499ccb658426` | Pending full review |
| `internal/cost/period_test.go` | `192d0e94f67b3631d58e3372240b3b6fa0796b8f` | Pending full review |
| `internal/cost/scanner.go` | `aaa7d4f85d576e7e76f2965513c23d53ea5f11ef` | Pending full review |
| `internal/cost/scanner_test.go` | `f4e7a4b3d908c37c41cc67d2be78643f1dc0d0fc` | Pending full review |
| `internal/defender/scanner.go` | `7e843f0e826f5c7a1b928eb22c4de525c6bd48f5` | Pending full review |
| `internal/defender/scanner_test.go` | `a72483ee1c6c6f0f4aad402328e8c2d58cf89d9b` | Pending full review |
| `internal/diagnostics/recommendations.go` | `9fb8f241d68502dd425ebb6e2f40f8414dcc5ba6` | Pending full review |
| `internal/diagnostics/scanner.go` | `b9c3cb25c28b1282c06866344effe0590539b13e` | Pending full review |
| `internal/diagnostics/scanner_test.go` | `7bfcb1ac7a56ad33a14ec79a91de08ff78951bd3` | Pending full review |
| `internal/findings/applicability.go` | `366dde9e246f1a8e92fd58f975dec70652ea98b2` | Pending full review |
| `internal/findings/applicability_test.go` | `068c908fa55bbcff01202294876dcf3d8198a0c3` | Pending full review |
| `internal/findings/summary.go` | `e4fc5b01bd25c668a85c3b77584b1a1d3223823a` | Pending full review |
| `internal/findings/summary_test.go` | `62cad5632c5402f649c1a4d88e7d7f203f584924` | Pending full review |
| `internal/gate/gate.go` | `35e3d887aff375b3ca63b673dfd090641dc5c4b4` | Pending full review |
| `internal/gate/gate_test.go` | `45b96ba6241d71a11f7437960f424e0e68b52d4d` | Pending full review |
| `internal/policy/scanner.go` | `1a549afefd7be50e0262336f733d8d5f76901a3d` | Pending full review |
| `internal/policy/scanner_test.go` | `74d77a818ba3d0164d141b3ee1b0656503d7e53e` | Pending full review |
| `internal/result/plugin_tables.go` | `9703654451984758e91b8c962ae5626f31e79b96` | Pending full review |
| `internal/result/plugin_tables_test.go` | `4a52a5a19bcc8658a2b8aa5b1832a3d89ca26a07` | Pending full review |
| `internal/result/result.go` | `f099c293a44707535ffe624c201eddcf240847f0` | Pending full review |
| `internal/result/result_test.go` | `bac7ad307b4c7f79d19831e1229d36293d9ef97e` | Pending full review |
| `internal/result/scope.go` | `a65bc4f42d7c973cb399fd3a06c7f30019bfb7fa` | Pending full review |
| `internal/result/scope_test.go` | `ed3ae54acdd9b33e90ab318836810fd3ac1523c4` | Pending full review |
| `internal/rules/catalog.go` | `e2d0a9ba82323ccb01dd837d14f7c25cf5d22c6d` | Pending full review |
| `internal/rules/catalog_test.go` | `e67a3d090bead57e968e2464c2de205b16494f12` | Pending full review |
| `internal/rules/embedded.go` | `5c3e0ce214e8d1571b9d1977b8f4a7f1d0dd326f` | Pending full review |
| `internal/rules/embedded_test.go` | `82611bcda9beb4d44c9a817f79379ad94d4591d4` | Pending full review |
| `internal/rules/inspection.go` | `97a998b23916134245f23fd15f0de1d4134897e7` | Pending full review |
| `internal/rules/inspection_test.go` | `8ff8d0177243a95bab178cabb8b4d4bebebe125f` | Pending full review |
| `internal/rules/provenance.go` | `0315d420dfc05b39cf03a6e782769d630d66ef8a` | Pending full review |
| `internal/rules/yaml_plugins.go` | `5a43f569f3ceceab829b9a26dd1d46c6732b5076` | Pending full review |
| `internal/rules/yaml_plugins_test.go` | `d02f2967b009aa549c422017ce8292b3a851b9e8` | Pending full review |
| `internal/scanners/registry.go` | `24512b10be30824c62f11bf83bc66e0b04e11eb6` | Pending full review |
| `internal/scanners/registry_test.go` | `16f8f679331a50d1346c10292bdb72e73674fb71` | Pending full review |
| `internal/skus/skus.go` | `396cc4b7172f536b08d91355f8f64e9e4f0b97a1` | Pending full review |
| `internal/skus/skus_test.go` | `6d1b6e6d89e4ac3a47c49fa54fbf34ec31db79dc` | Pending full review |
| `internal/stages/config.go` | `307eabbb30b2b988cdc6e2fc750b3bce55d993d9` | Pending full review |
| `internal/stages/config_test.go` | `804ccc644facc788b357fb36e54d06fe21c50e05` | Pending full review |
| `internal/stages/options.go` | `d0138077d1784f5e82cde165af8eb62b410e55d2` | Pending full review |
| `internal/stages/options_test.go` | `8245b9a14a468022ded93bdf62ab23538ecd6e16` | Pending full review |
| `internal/stages/runner.go` | `7f3ec14d2f7b150619446984ddd9659b9761d8e4` | Pending full review |
| `internal/stages/runner_test.go` | `6e6b6f6bbeb0db79a0b4d2c1a21cbd5ef46c3be7` | Pending full review |
