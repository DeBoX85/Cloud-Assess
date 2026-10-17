# Section 4 contracts and coverage reconciliation, 2026-10-16

Review state: initial full reads and contract reconciliation completed for the four remaining Go packages and the 62 tracked script/workflow/support files at regression head 10efcd48. SR-F12 was reproduced on both native hosts and is now corrected locally under the published contract; final integrated QA and final documentation readback are still required. This is not Section 4 completion, protected acceptance or independent-person review.

## Tool contracts

| Package | Reviewed behavior and boundary |
| --- | --- |
| internal/equivalence | Eleven core datasets with fixed ordering, stage coverage and field deltas; input schema/completeness/redaction preconditions; documented provenance/display-name/timestamp/number normalization. The SR-F12 correction rejects conflicting projected records as non-comparable while retaining sanitized diagnostics. Identical normalized records may deduplicate; resource-type counts deliberately aggregate. This core harness does not implement source historical Excel comparison or compare plugin tables. Loaders are not full schema/identity/size validators. |
| tools/equivalence | Real-command exits 0 for equivalent, 1 for differences or noncomparability, and 2 for invocation or I/O error. Literal SR-F12 tests cover both sides, both orders and repeated trailing rows, plus identical normalization and 2+3 count aggregation. Output includes raw diagnostic identities by design and uses ordinary local file writes; it does not inherit renderer replacement/crash guarantees. Inputs/outputs require private operator storage. |
| tools/brand-build | Version/profile/native-target checks precede build output. Explicit argv plus encoded canonical linker metadata avoids shell evaluation. Strict clean stamped candidate verification is later enforced by packager, not promised by this development builder alone. Exclusive same-filesystem hard-link publication preserves an existing executable; required default/custom actual-binary tests cover quoted display text and invalid profiles. |
| tools/diagnostics-probe | Explicit operator-only retained report, eligible inventory and deliberate request bound; individual diagnostic-settings GET or one GET per ARM batch. Context/client timeout and sanitized hash/type/status summaries. Summary failures are evidence, not a scan-completeness assertion; successful process return means a summary was written. A current probe cannot retroactively identify an old uncorrelated multi-request failure. No live invocation performed. |

## Script, workflow and provenance contracts

| Family | Reconciled behavior and acceptance limits |
| --- | --- |
| bootstrap-workspace | New task directory only; exact official Go/PowerShell checksums precede extraction, actual source/APRL identities and module verification, readonly modules/local caches. Partial failed workspace is retained, never reset. Recovered archives match the documented checksums. Provisioning alone is not QA. |
| live-equivalence and exclusion preparer | Exact clean target/reference and initialized APRL checks; scope/stages/filter hashes/exit/status evidence captured. PowerShell 5.1 native stderr handled separately from exit. Matched source/target flags, subscription-wide Cost, local unredacted reports. Exclusion preparer checks one expected subscription/resource and YAML quoting without Azure calls. Examples remain scope-specific historical fixtures, not newly authorized scopes. |
| immutable branding and packages | Go/Python canonical resolved profile boundary, safe component/version/URL validation and no shell interpolation. Packager captures a clean committed revision/tree, validates compiled modules/checksums/native target/CGO/trimpath, includes unchanged notices and canonical branding. Extractor validates all members before writes, rejects traversal/aliases/symlinks/modes/hash/layout errors. Repeated same-runner ZIP bytes do not certify cross-host rebuilds; unsigned hashes are not publisher authentication; orphan ZIP on sidecar failure remains detectable and is not overwritten. |
| dependency evidence | Selected fixed module graph, platform-specific CLI membership, original root license bytes and bundled rules/SKU/toolchain notices; no guessed SPDX/legal approval. Disposable module downloads cannot mutate project sums. `--check` is read-only freshness validation. Module-only advisory and package-specific release review remain open; no dependency update executed. |
| source characterization | Isolated exact AZQR 8e4f0577 / tree 17d93b20 / APRL 60eaddda; eleven added synthetic harnesses, unchanged production/module inputs. 23 records emitted as indexed JSON chunks before unconditional retained-byte checks. This is Linux synthetic source evidence, not live Azure or Windows SDK execution. All 315 chunks were reconstructed at 10efcd48 and match 909,865 retained bytes and Git blobs; final correction head needs its own proof. |
| compiling controls | Named baseline, altered production/source expectation, compiling named assertion failure and restoration. In-place Sections 1–3/ARG/disk controls must run sequentially with other source consumers; copied region/AI controls are isolated. Compiler failures are not detections. No claim of exhaustive mutation coverage. All prior controls remain mandatory. |
| maintenance | Actual workflow recipes with source-tree checks and scoped bot proposal publication; no protected push target. Local bare-remote tests cover no change, scoped repair, invalid layout/snapshot, existing/divergent proposal and rejected push. Hosted dispatch/token/PR event validation stays deferred. Neither dispatch nor dependency generation was invoked. |
| CI | Immutable action pins, native Linux/Windows jobs, source isolation, explicit shell failure propagation, deterministic formatter/module/inventory/CLI/package/docs checks, full tests/race/vet/coverage/five bounded fuzz runs and vulnerability scans. Read-only PR permissions. Concurrency may cancel earlier heads; canceled work is never transferred as a pass. Both final native jobs and source gates are still required. |
| QA harnesses | Offline CLI HTTP/auth tripwires, real 12-sheet/12-CSV/JSON/stdout/SARIF paired report checks, scope/privacy/partial inputs, extracted native executable tests and controlled Windows ACL/lock checks. Temporary-directory cleanup retries only bounded Windows sharing violations; other failures remain failures. Synthetic QA is separate from fresh-OS/installed Azure evidence. |

## Whole-system intent and outstanding development

Target specification, implementation/execution plan, current session authority, three completed audit sections, operations/access, packaging/dependencies, parity characterization and Gate 004/deferred registers were reconciled by current contract. Historical status/next-action paragraphs remain historical; they are not silently promoted to current requirements or fresh passes. The latest session/PR evidence index remains the resume authority.

The architecture continues to separate validated CLI/configuration, scope and inventory, bounded queries/collectors, owned canonical results/health, then format projection and post-report severity. Excel and JSON remain first-class; branding remains presentation-only. Failure health is separate from finding severity and estate visibility. Schema 1.0 ordinary and additive 1.1 plugin-table contracts are unchanged. Read-oriented query POSTs and ARM GET envelopes are distinguished from mutations, and automatic provider registration remains disabled. Full pinned-source parity is the project goal, with deliberate documented safety corrections rather than silent source-bug copying.

Unchanged source-data diff against accepted a6fee427 covers go.mod/go.sum, rules/APRL/AOR/custom, SKU data, all region source captures and latency. Source/target captures/pins are not refreshed. Full source/transport contracts from Sections 1–3 remain authoritative; final exact-head QA must retain their controls. No unresolved Sections 1–3 question is reclassified as an additional confirmed defect without independent evidence.

After the audit's final disposition, resume draft PR141 Storage pre-implementation contract and independent expectations. Grouping repeated names across kind/region, kind-aware projection, display/stage locations, paging and omitted metadata need concrete expectations before implementation. Only Compute disk target SKU adapter is accepted (1/7); source SKU captures cover 7/7. The six remaining target families, region coordinator/public execution and FullCost/retail/history collection remain open. The generic comparator is not the source historical-compare command. B6/other parity, AR-02/load, restricted-visibility/nonempty stage/Arc numeric/sovereign/live parity, fresh OS, hosted maintenance, module-only advisory disposition, legal/security/human/Gate 004/release all remain separate. No Azure operations, new scope, fixture provisioning or dependency refresh is authorized here.

## Evidence and local QA notes

The baseline union check accounts for all 251 original Go paths across the original Section 1 manifest in SYSTEMATIC_MODULE_REVIEW, Section 2 and 3 reports and this Section 4 manifest, across 41 package directories. It is a coverage-index reconciliation, not a new claim to have rerun previous sections or a weighted overall project percentage. The current audit adds 12 Go test files; the original 251 paths/pins and prior regressions remain retained.

Fresh local read-only dependency-inventory verification passed: 49 selected modules, 9 direct, 24 Linux CLI members and 26 Windows CLI members. Python AST parsing of all 49 initial tracked Python files and the added Section 4 runner and shell syntax checks passed. Dependency safeguards (9 cases), bootstrap (3 cases) plus checksum reversal, Windows workflow failure propagation (6 cases), maintenance proposal (7 cases), PowerShell helper tests and 1031 authored local links passed. These are focused local results; full final native acceptance is still mandatory.

One local command mistakenly invoked a PowerShell test file with Python and returned SyntaxError. It did not run the test or modify production. The same file was subsequently invoked correctly with pinned PowerShell 7.6.6 and passed. The failed invocation is retained as QA-S4-01; it is not a product defect or accepted test result. Prevention: use each script's extension/shebang and explicit interpreter, and inspect every command's exit before treating a batch as passed. The earlier uninitialized-submodule setup failure was also excluded from reproduction credit. Native PowerShell 5.1 proof remains separate.

Source reproduction run 38856474640 completed successfully with all eleven harnesses. Full returned log 1,359,793 characters / 1,359,869 UTF-8 bytes, SHA-256 `134b48c7df3b0e7c835fb87194569962f6eaa2934201105900bafa1673462918`. All 23 captures / 315 unique contiguous chunks / 909,865 bytes match declared SHA-256 values, retained files and exact-head Git blobs; pinned source/tree/APRL match. Full corrected-head source proof is still required.

## Exact initial-read manifests

Manifest hashes below identify reviewed regression-head bytes, not final acceptance. Corrections/new tests/helpers need delta review and final QA. All eleven baseline Go files and the new independent regression were fully read; all 62 script/workflow/support files were fully read. Source-owned configurations and fixtures were reconciled alongside Sections 1–3 and unchanged data/provenance checks, rather than regenerating expectations.

| Go path | Git blob | Lines |
| --- | --- | ---: |
| `internal/equivalence/compare.go` | `1d0b77de42d21c555c96bad6f473964d765fb544` | 116 |
| `internal/equivalence/compare_guard_test.go` | `e59fe464f75e415870fef0bcb6cc096c38ae55c0` | 84 |
| `internal/equivalence/projection.go` | `cb7f31799442cba5227e13096d7f4c9dcd7d95c2` | 657 |
| `internal/equivalence/projection_test.go` | `3915be9d5907fced529204bb5ecf18b78df7a012` | 362 |
| `internal/equivalence/types.go` | `f7be9b21a3aeabc654730ca4086e0b4e7a3f5ecb` | 122 |
| `tools/brand-build/main.go` | `875d824c50c48319faab2106c0d7ddacaaadb98c` | 97 |
| `tools/diagnostics-probe/main.go` | `ea2b806be70686e6524242f021f8df5751c04f5b` | 251 |
| `tools/diagnostics-probe/main_test.go` | `abeef1ac10f7083f85ee27f34b4354ebce3dcbf9` | 122 |
| `tools/equivalence/duplicate_evidence_test.go` | `126f45bc1980a750245474d70ee7fe490f087917` | 83 |
| `tools/equivalence/input_evidence_test.go` | `3d0a89b227cb523f2ffeb9b169f4d196bb894e95` | 100 |
| `tools/equivalence/main.go` | `66b8f529c8893ec87ac4f19ffa766885be9f80b8` | 84 |
| `tools/equivalence/main_test.go` | `89be2708f8ca84a8113df978e9cc0d9246479cca` | 102 |

| Script/workflow/support path | Git blob | Lines |
| --- | --- | ---: |
| `.github/go-mod-tidy-trigger` | `21323c5a37b1bf0ef39a3dfa152600b3cc7b84fa` | 1 |
| `.github/pull_request_template.md` | `33f07fc5bcb2b0d08b3ce4c109b947b42518941a` | 17 |
| `.github/rule-sync-trigger` | `851d604358a9cff997656f9d4c07e8b73dbf0a20` | 2 |
| `.github/workflows/import-rule-baseline.yml` | `d8677b2da84ed72abae7dd6a66346f084e010204` | 72 |
| `.github/workflows/region-source-characterization.yml` | `cfddb4bfcf1fc2276288bafb925d48d4816c8867` | 36 |
| `.github/workflows/test.yml` | `c9b40be73cb5e43e355a244ed32a9e548e565fae` | 467 |
| `.github/workflows/tidy.yml` | `5d90829578fea27b72f273aace54b23e3aad7ae9` | 37 |
| `scripts/bootstrap-workspace.sh` | `53be7d565640c8d873760f8e063163e231233f86` | 127 |
| `scripts/branding-profile.py` | `f97eb370157f18d3b36b355f479413ebbff445b9` | 64 |
| `scripts/capture-region-source.py` | `59c5a4915600dda8c786aa8b0643bd8af7365bfb` | 123 |
| `scripts/dependency-inventory.py` | `f76d92a0d36d5ee506a5f5a7972a0a5211ae7812` | 153 |
| `scripts/live-equivalence.helpers.psm1` | `bb2b1590e403375c963797daff462c8753f26a21` | 176 |
| `scripts/live-equivalence.ps1` | `5d4f64248659afc0afedca15098ae4ee012a6e72` | 424 |
| `scripts/package-candidate.py` | `be29991465ba8d847ef8477fd2ce3ed34858c905` | 219 |
| `scripts/prepare-resource-exclusion.ps1` | `74160968e8b6a913342d8b210ee03f0cec1d698a` | 123 |
| `scripts/publish-maintenance.sh` | `2ca62cc146154a88b67a37550a14882cd82ef462` | 59 |
| `scripts/tests/ai-execution-mutation.py` | `feb4a7e50eb5541c7fd11e20cba3a76ef51a39b4` | 50 |
| `scripts/tests/ai-request-mutation.py` | `ddd1895658aab8f1c4d5541824bfe74a812e4244` | 52 |
| `scripts/tests/arg-response-framing-mutation.py` | `06315565951f7fa07f760111e1335e45ffba5c08` | 50 |
| `scripts/tests/bootstrap-workspace.py` | `3a3650f9d5e45a61c4a104f495a9d2f4a2c88729` | 119 |
| `scripts/tests/branded-package.py` | `5b729447d473a83c77d40bd57c073bcdaeaa62c7` | 31 |
| `scripts/tests/branding-reports.py` | `6c204c8d348ebe19f8d9d431e82133709ec7db91` | 192 |
| `scripts/tests/branding.py` | `d17d968b181ea3d1c43180fc99880de2b9a22cf2` | 108 |
| `scripts/tests/built-cli.py` | `455eb6152187ce2d59e37d068b29d2a7216a9cbd` | 426 |
| `scripts/tests/compute-disk-sku-collector-mutation.py` | `4ad481d7d24dfd288533a7794c6de418a911f5ed` | 69 |
| `scripts/tests/dependency-inventory.py` | `b5e4e6b8bd178772c7a98e0facd6e3707ae00ffd` | 128 |
| `scripts/tests/documentation.py` | `73f737436e754ce32e788d742613fcc16b4668ba` | 76 |
| `scripts/tests/live-equivalence.tests.ps1` | `ac0e45d3f0285245c25de8a3ee2007098a42c2a8` | 109 |
| `scripts/tests/maintenance-generation.py` | `f32611354703065e13bdb13bd49944a634c26749` | 220 |
| `scripts/tests/maintenance-publication.py` | `2a82eddda9c35cda8c038ace21b53ae7f9b67691` | 153 |
| `scripts/tests/package-candidate.py` | `5f7aa9b00d644c4dc3f3b0259109ac36d0efddf2` | 239 |
| `scripts/tests/qa-mutation.py` | `d6e8ee02f9e47773be9e059b430b232eab8a586b` | 42 |
| `scripts/tests/qa-temp-cleanup.py` | `c44624027c09b377dd6701f529367729f4f9fa3d` | 108 |
| `scripts/tests/qa_temp_cleanup.py` | `7953e2ba3c2c60d6bdb06d62f6b5ddb771e7d3db` | 15 |
| `scripts/tests/region-auxiliary-mutation.py` | `3b46159ba695da177f5cdccf5670570f37e808ca` | 49 |
| `scripts/tests/region-availability-capture-mutation.py` | `1f08bbd4653c541f8ffa644e02839c25ff9ff145` | 51 |
| `scripts/tests/region-availability-mutation.py` | `a6301ee4ed107c84948e673c958a705deae597fe` | 54 |
| `scripts/tests/region-compute-sku-capture-mutation.py` | `f748aa76ac01c1f4d9f31b4acfb83344063b1531` | 91 |
| `scripts/tests/region-cost-enrichment-capture-mutation.py` | `32c97a20f0fb9a7da32c7c12b5e20ea40e2c7d54` | 54 |
| `scripts/tests/region-cost-enrichment-mutation.py` | `e7c98a0134cf6a3c87bac3136f035f1d6dc6f2cc` | 63 |
| `scripts/tests/region-cost-mutation.py` | `7cd73153c307aee9b044af60633e03f5207b9df7` | 62 |
| `scripts/tests/region-inventory-aggregation-mutation.py` | `e68828ae080cac91003973138ecc8533c8c7e2af` | 66 |
| `scripts/tests/region-inventory-mutation.py` | `8bb7956a8ad5ad4c17eb3be67e399b608b93d69d` | 61 |
| `scripts/tests/region-latency-capture-mutation.py` | `ea325dc97d26fdaa47c2d8eee28b232409bfcf03` | 51 |
| `scripts/tests/region-latency-mutation.py` | `7eccf70f359803e6360d0cf6e21299773176342a` | 57 |
| `scripts/tests/region-other-sku-capture-mutation.py` | `f01d6dda19fd2e6ea9b4d5e2c21f3d25a29d4d49` | 100 |
| `scripts/tests/region-primary-mutation.py` | `e4ed5303df304bcac0e91f25bcd7eb67d75e13b1` | 53 |
| `scripts/tests/region-provider-locations-capture-mutation.py` | `145bea96f8fe587dc38e296bc67b33eb81cd234b` | 75 |
| `scripts/tests/region-provider-locations-collector-mutation.py` | `8d2111616f8bcbb330c48cd707d1d9591d188452` | 67 |
| `scripts/tests/region-quota-capture-mutation.py` | `a0a014e37edd8945a1dcf27761587d363b1bdd44` | 52 |
| `scripts/tests/region-quota-collector-mutation.py` | `7b9fb20c8e922806facdeffc87b4ff6014f637f3` | 70 |
| `scripts/tests/region-quota-runtime-mutation.py` | `a48202c1d82eafc5bc82d7b4332a62333bee019f` | 57 |
| `scripts/tests/region-reservation-capture-mutation.py` | `4f942f1af28038e9d670c7eb074a33d4b6ff1a31` | 58 |
| `scripts/tests/region-reservation-collector-mutation.py` | `7323d0edbc61ffd8137d63ffc545c7907107e768` | 127 |
| `scripts/tests/region-reservation-runtime-mutation.py` | `fbdcb37335b4627f9db4c62e664082c9d98dfb21` | 63 |
| `scripts/tests/region-service-mutation.py` | `488cd60ebb9198f942605e86b604fbe620dcee25` | 64 |
| `scripts/tests/region-source-capture.py` | `e6a4873624db9f5fbfb63bc360c93be0e4d14af4` | 186 |
| `scripts/tests/region-vm-quota-capture-mutation.py` | `13fa511c57da0d64c82ab26a12291e1346e1772d` | 52 |
| `scripts/tests/section1-regressions.py` | `203f4552e0453774163f91187055edae86b6fe06` | 110 |
| `scripts/tests/section2-regressions.py` | `1892ba62673ce998c39ec0bd2bea4668c141b97b` | 150 |
| `scripts/tests/section3-regressions.py` | `b23f5f3dd335310eb02c6f0c622521b33bba369d` | 83 |
| `scripts/tests/workflow-command-failures.py` | `4aefe8afa6702b44f6183ba59bb9f84af250a273` | 49 |
