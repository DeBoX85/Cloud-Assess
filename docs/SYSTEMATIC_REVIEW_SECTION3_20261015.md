# Current handover authority — 2026-10-16

Read [handover checkpoint](HANDOVER_CHECKPOINT_20261016.md), [restart prompt](NEW_CHAT_PROMPT_20261016.md) and [Section 3 active report](SYSTEMATIC_REVIEW_SECTION3_20261015.md) first. This entry supersedes earlier current-state, pause and next-action paragraphs below; historical contracts, failures and evidence remain preserved.

The user requested a new-conversation handoff during Section 3. Audit Sections 1–2 are complete; Section 3 is incomplete, with 42/94 baseline Go files initially read and no Section 3 package finally verified. Audit counts remain **2/4 sections (50%), 22/41 package directories (53.7%)**; neither measures weighted overall project completion. Section 3 has 15 packages/94 baseline Go files (43 production,51 tests), plus one added SARIF regression. AI governance10Go files and region42Go files remain, alongside contract/data/architecture reconciliation across the initial reads.

Verified code/test head `b293e9533b06667d15b5e7efb6a20b39a9ea8427`, tree `71a1adb10c20ac23b0a55bd361c88b2204a2b048`, is draft/unmerged in [PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). Protected `bootstrap/core-v1` stays `a6fee427e19feeef28691737ea9bc5026086efa0`. **SR-F10 is reproduced and UNFIXED**: SARIF mixed-case result IDs do not match their published rule descriptors. Formatted regression-only native38829684903 fails the declared two assertions on BOTH Linux116745335663 and Windows116745335627 after compilation; Linux formatting passes and no race warning was observed. Source38829684879 succeeds separately. Red is reproduction, never acceptance. Earlier ed30ef4 Linux formatting failure is QA-S3-01 and is corrected in b293e95; do not repeat it.

NEXT: verify live refs/PR142 and actual files; implement the documented isolated SR-F10 correction without mutating canonical IDs or fingerprints, add compiling reversal/restoration control to both required jobs, complete ALL remaining Section 3 source/test/contract review, fix substantiated findings, obtain full corrected exact-head Linux/Windows/source QA, then publish the Section 3 checkpoint and STOP for re-evaluation before Section 4. Continuation of Section 3 is already authorized. Do not start Section 4, merge protected history, repeat Sections 1–2 or resume feature development as part of this handoff.

Remaining estimate: Section 3 **5–8 active hours**, Section 4 **2–4 active hours**, plus material fixes/repeat QA. This replaces the earlier optimistic6–11hour combined estimate; it is not a wall-clock completion promise. Overall project percentage/release ETA are unverified. After the FULL audit resume Storage contract/tests on draft [PR141](https://github.com/DeBoX85/Cloud-Assess/pull/141), then remaining adapters/integration. Source captures7/7, accepted target SKU adapters1/7, public region unavailable. FullCost/other parity, AR-02/load, live/fresh OS, module-only advisory disposition, human/license/security/Gate004/release remain open. No Azure access/write, fixture/RBAC change, dependency/pin/capture/schema update or production-scope substitution is authorized by this handoff.

This documentation successor changes no Go/Python/workflow/data. Its new checks are separate and initially pending; the known unfixed regression can keep native CI red. Inspect live PR142 for the actual documentation head/results. Earlier passing code evidence must not be transferred to a failing/changed head.

---

## Historical record below

## Current verified reproduction evidence

Exact formatted code/test head b293e9533b06667d15b5e7efb6a20b39a9ea8427, tree71a1adb10c20ac23b0a55bd361c88b2204a2b048, parent ed30ef46910df31593656ba0910a6242173201e8. No Section3production correction exists. Native [38829684903](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38829684903) completed/failure. Both hosts compile and fail ONLY TestRuleReferencesUsePublishedDescriptorIdentity/upper_finding and /upper_descriptor; same-identity control passes. Linux Verify formatting succeeded. No build failure or data-race warning in complete returned logs. Later full gates were not reached and are not credited. Source [38829684879](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38829684879) completed/success separately.

| Host | Job | Full characters | Exact UTF-8 bytes | SHA-256 |
| --- | ---: | ---: | ---: | --- |
| Linux | 116745335663 | 64933 | 64937 | `be17229af65f4d016371bc7324c8b6bb53d6b94b0bd197bc50dad302a8c1b47a` |
| Windows | 116745335627 | 51425 | 51427 | `5742e9e63e5f15ccbb1696596bcf708646592cd9a57caa08cf427872fdce3b1f` |

Earlier Linux formatting failure QA-S3-01 is retained below and is not product reproduction. Current pinned Go1.26.9 formatter was rechecked with empty output on the published regression. The launcher and toolchain are ephemeral workspace dependencies, not repository modifications. An old absolute toolchain path disappeared and was restored by the existing Go launcher; re-check tool availability in a new environment. Baseline local race-suite evidence (all15packages before adding this regression) is preserved below, not final acceptance.

---

## Earlier Section 3 records and baseline manifest

# Section 3 systemic review — active record

Entry baseline: `798d7be21e393f1e85503303e3d8eb83d311913d`. No Section 3 completion or new product defect is claimed at preparation.

## Coverage and method

All15packages/94Go files (43production/51tests), with immediate data/contract boundaries. Read complete source/tests; execution alone is not code-review coverage. Independently reproduce material defects with literal/source-grounded expectations, fix in isolated commits, verify uncertain/version-sensitive contracts against pinned code or official primary documentation. Region calculation/source work is distinct from incomplete public runtime adapters.

## Findings and decisions

### SR-F10 — SARIF result ID differs from its published descriptor, P2

Confirmed locally after compilation on unchanged production at entry798d7be: definitions and findings are joined case-insensitively by the canonical summary, but SARIF writes the finding's original ID while publishing the definition's ID. Two independent casing cases fail exact descriptor-reference assertions; the same-case control passes. Test: `internal/renderers/sarif/rule_reference_test.go`, `TestRuleReferencesUsePublishedDescriptorIdentity`. No production fix yet.

Pinned AZQR `internal/renderers/sarif/sarif.go` at8e4f0577 also uses the raw finding ID after a lower-case join. This is a proposed deliberate correction to inherited behavior, not a target-only source divergence. Official OASIS SARIF2.1.0+Errata01 §§3.27.5/3.27.7/3.52.4 establish descriptor-reference identity; https://docs.oasis-open.org/sarif/sarif/v2.1.0/sarif-v2.1.0.html. Use the published descriptor's spelling in SARIF results while retaining canonical input identity, rule metadata, dedup, resource identities and normalized fingerprints. No format/schema/pin/dependency change.

Before production edits, acceptance is: same-case and both casing directions resolve to the actual published descriptor; rendering does not mutate canonical input; existing results/fingerprints/empty/SLA/report privacy behavior remains. Change only the SARIF private result builder/call site, add the independent regression and a compiling reversal/restoration control in both native jobs. Both-host reproduction is now verified at b293e95; compiling correction controls and full corrected-head QA remain pending. Red is evidence, not acceptance. Rollback must account for the new expected identity and helper anchor; ordinary targeted reverts only.

Reportfile/redact/branding, all renderers, root plugin projection boundaries, zone, service-health and SQL-EOL source/tests have initial complete reads. Carbon projection and request-scanner source/tests are also initially read. Contract/integration reconciliation remains pending; no final package completion is credited yet. Human core projection redaction boundaries and zone JSON alias handling remain review questions, not confirmed findings.

## QA and recovery

Entry native38652456544 (Linux116154267913, Windows116154267785) and source38652456593/job116154266947 report success; full logs pending reconciliation. Local Go1.26.9 was automatically fetched by the existing1.26.8launcher, without repository changes. Local checks do not replace native Windows/live evidence.

NEXT: continue full source/test review and update this record/handover after coherent slices, then reproduce/fix/QA and stop before Section 4. Keep PR142 draft/unmerged. No Section 3 production changes yet; ordinary revert of documentation preparation must account for successors. Previous fix/parent and failure evidence remains in Section 2 report/journal. A manifest preparation assertion caught an incorrect handwritten42/52split expectation before any file write; actual inventory below is authoritative (94total,43production,51tests), not a product defect.

## Baseline file coverage

| File | Baseline Git blob | Lines | Review state |
| --- | --- | ---: | --- |
| `internal/plugins/aigov_projection.go` | `408b7cda42858faeab5b18b9fcba43cf8aa80c01` | 112 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/aigov_projection_test.go` | `c016c61fe54c6d006a3b186374245a4018fd4289` | 81 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/carbon_projection.go` | `626c3b24d2e2adb7dec5a53c0fef7f8d216efd89` | 95 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/carbon_projection_test.go` | `a2b0221abaa4b345d67afd7f7877829beb2562d5` | 115 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/registry.go` | `0e7e6810d481c5c7e806ee6469daf3c03f3f33c4` | 186 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/registry_test.go` | `202b7a54f20738d90cf1f7650280ecd235c2a96b` | 55 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/servicehealth_projection_test.go` | `ff9b4db92f36855a592f79a7c23ecbf08fba5170` | 91 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/sqleol_projection.go` | `a925e0545eaa0e9b4be51a58705d46c0004f56fe` | 70 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/sqleol_projection_test.go` | `cb8616b373409d017318f1d9f5bf20131a4492e0` | 87 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/zone/scanner.go` | `1ba211bad8d2b0558b91481d177acb35417a568d` | 314 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/zone/scanner_test.go` | `842fa2d7be23836ec7277bc5dcd2541f9e72a836` | 331 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/servicehealth/query.go` | `7624c0a3d23788f8c3e1565b2f35f55102b7e48f` | 72 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/servicehealth/scanner.go` | `ffdc63d0b4669ad6ab23908deea1e69a4785ccd0` | 328 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/servicehealth/scanner_test.go` | `1ef072fe74f1fd0024a1aaf56f38a162cead182c` | 355 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/sqleol/query.go` | `bb408e06215f9fb2d667dbc6372033b5e83d7858` | 7 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/sqleol/scanner.go` | `2abf52dfc1ca3fda34e97dda898b7aa1d8631cc0` | 289 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/sqleol/scanner_test.go` | `5ad2c6c14c6ecf23820ea989a8cfb8b62dc34ae1` | 418 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/carbon/projection.go` | `868599a28286a303671cc8a76ecd06d90f13d1e0` | 158 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/carbon/projection_test.go` | `105315a948c25c97be421b699e13eec17c3234e1` | 190 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/carbon/scanner.go` | `d0aff9de47afad2c1a1b8c3fcd3ff9819ac5b1fd` | 447 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/carbon/scanner_test.go` | `f8e1e5dca08bcc1706a5a48ff804bf45044d326e` | 625 | Initial complete read; contract/QA reconciliation pending |
| `internal/plugins/aigov/decoder.go` | `c31875f5576165af668ce95eaa2408767321a203` | 454 | Pending |
| `internal/plugins/aigov/decoder_test.go` | `76e9df63a13a056b7799877dbeaf4b0282c71b3c` | 143 | Pending |
| `internal/plugins/aigov/discovery.go` | `63e8d826cd8e4cb4fd51a6ea86a5c4cc40a86660` | 314 | Pending |
| `internal/plugins/aigov/discovery_test.go` | `9f32084ed4059b885f158ef28044403b1c206503` | 633 | Pending |
| `internal/plugins/aigov/execution.go` | `558d79cdb2526dcff5256741061737c3a2bc8f7b` | 70 | Pending |
| `internal/plugins/aigov/execution_test.go` | `e63932442c8fce7f51b4f9a88657d8ffd49c0bd6` | 113 | Pending |
| `internal/plugins/aigov/projection.go` | `90ceacb770826aff716c4d1ebc2e106ba0b6211c` | 312 | Pending |
| `internal/plugins/aigov/projection_test.go` | `20454b562bb0645b53157958e97f249e5de66acb` | 345 | Pending |
| `internal/plugins/aigov/scanner.go` | `c00e0ee00be57cf8bc7256322fbc62722f64a768` | 463 | Pending |
| `internal/plugins/aigov/scanner_test.go` | `3da3767d71fefb91d76576f712044c742bc929be` | 749 | Pending |
| `internal/plugins/region/auxiliary.go` | `15ab8323b5b4dd7974105eb15994d777f160ad35` | 260 | Pending |
| `internal/plugins/region/auxiliary_test.go` | `1526c4508c2a25a009305b9accf7eb5daa4c05b9` | 342 | Pending |
| `internal/plugins/region/availability.go` | `94360cfbe6a59a3debdf1a8f002cf7b062f8e264` | 308 | Pending |
| `internal/plugins/region/availability_capture_test.go` | `a212b8f605c28d8eac7137c893de3dc4f89ddaa5` | 146 | Pending |
| `internal/plugins/region/availability_test.go` | `695105d71c28c64980de19334f1a4370c3057869` | 329 | Pending |
| `internal/plugins/region/compute_disk_sku_collector.go` | `d7a5b8b2e1facd7a36b157acc9eadd4d91a57465` | 633 | Pending |
| `internal/plugins/region/compute_disk_sku_collector_test.go` | `76ce22e83bc9fe9fe4bdb51520782a0fed4cc97e` | 652 | Pending |
| `internal/plugins/region/compute_disk_sku_review_test.go` | `8789739cb653347be8b596e9fccb7e4836410212` | 65 | Pending |
| `internal/plugins/region/compute_sku_capture_test.go` | `2ef251506b2c6752a02685016aacee513bdf1978` | 193 | Pending |
| `internal/plugins/region/cost.go` | `dc7c43cdcfc46c2fa67e00aaafbc5cf0faeec149` | 235 | Pending |
| `internal/plugins/region/cost_enrichment.go` | `94a9cc17e8d70a02852274a69bfd6a6f5eb193d5` | 278 | Pending |
| `internal/plugins/region/cost_enrichment_capture_test.go` | `92cdcf06990aa8a5139c1966ee9b2ff9a62eec84` | 96 | Pending |
| `internal/plugins/region/cost_enrichment_test.go` | `c64816dd1daa8be6695b2d10076c267854955a7b` | 468 | Pending |
| `internal/plugins/region/cost_test.go` | `48712680a4eb91785d889a3c2796bfb3b6e623a3` | 251 | Pending |
| `internal/plugins/region/inventory.go` | `6157171769aaca1a4e78a026a4dc7f02a610bb87` | 134 | Pending |
| `internal/plugins/region/inventory_aggregation.go` | `9828cf3a30056e2456028b6e746e1e23afd5782c` | 197 | Pending |
| `internal/plugins/region/inventory_aggregation_test.go` | `b23fbf28124fdb72a78d5b86c85fff897e4cd30c` | 316 | Pending |
| `internal/plugins/region/inventory_capture_test.go` | `c7f85db283725e33b2d89a4df85dfdf6e8f8638e` | 107 | Pending |
| `internal/plugins/region/inventory_test.go` | `ac97771893275ad42fc5272a435e8a9dae57d718` | 294 | Pending |
| `internal/plugins/region/latency.go` | `1a4a037e8d6c1c80aed1d7d51c66a5d5d1ae1090` | 217 | Pending |
| `internal/plugins/region/latency_capture_test.go` | `877e215f1212f85096033acf9636618620479d13` | 99 | Pending |
| `internal/plugins/region/latency_test.go` | `60ea39a005cf1f1aff3f13c3889700706fc46dc2` | 258 | Pending |
| `internal/plugins/region/other_sku_capture_test.go` | `c73799952e61070813757283c13dc66079408ee7` | 1961 | Pending |
| `internal/plugins/region/primary.go` | `3a7fb0b9601a525f91c385d0016ce0b8ff532303` | 322 | Pending |
| `internal/plugins/region/primary_test.go` | `6219219eca501b00af02a8af6fb23bb4a51c5d78` | 335 | Pending |
| `internal/plugins/region/provider_locations_capture_test.go` | `87112d5d3876fba242dbd2e0bb6ec9e24b35085f` | 183 | Pending |
| `internal/plugins/region/provider_locations_collector.go` | `8c85347d96488f2d4b1440da3f98a64ab3657413` | 324 | Pending |
| `internal/plugins/region/provider_locations_collector_test.go` | `9d1180be0e50e0f7fca5d30bbda00cb25eeb3ca6` | 547 | Pending |
| `internal/plugins/region/quota_capture_test.go` | `78a92620225e5e317e7af84dba73d6971828c326` | 80 | Pending |
| `internal/plugins/region/quota_collector.go` | `e2a44936371ec66600f33e2378b80ec689ee8b3f` | 352 | Pending |
| `internal/plugins/region/quota_collector_test.go` | `70a85e863018ef976e4f54a92318f3503d26b3f1` | 543 | Pending |
| `internal/plugins/region/quota_runtime.go` | `a2be4a680e2c3d705ed1198d1cbb509a04e44713` | 234 | Pending |
| `internal/plugins/region/quota_runtime_test.go` | `5ac53f03b1614a386b37fd92ed8071885b3efab2` | 353 | Pending |
| `internal/plugins/region/reservation_capture_test.go` | `31377c86e8b54423d799279c6d29ca3870cc1e72` | 114 | Pending |
| `internal/plugins/region/reservation_collector.go` | `c5b9c83e5d41be27798d7ab208f2568cca437c49` | 357 | Pending |
| `internal/plugins/region/reservation_collector_test.go` | `d617b44f6701b59234fbfb3866c08035a26fdf02` | 566 | Pending |
| `internal/plugins/region/reservation_runtime.go` | `bc6e23eca781e5b5ad81be2f416e2f9c36ec58c8` | 269 | Pending |
| `internal/plugins/region/reservation_runtime_test.go` | `1ec9bffe0cddf1e624c2b92be434ed8ee7952370` | 421 | Pending |
| `internal/plugins/region/service.go` | `e2e53165c2a4cb597c2bc36c239557189ad68ccd` | 356 | Pending |
| `internal/plugins/region/service_test.go` | `f683d274a5b208d4d54a7d1081de3f43aaba1581` | 377 | Pending |
| `internal/plugins/region/vm_quota_capture_test.go` | `18476fb14c4385d67e7d7fd16ccabbcc7ba8b52c` | 77 | Pending |
| `internal/plugins/region/vm_quota_collector_test.go` | `47546d5e6b3d79e9d0baa230f6739a3b75fbfa72` | 218 | Pending |
| `internal/renderers/csv/csv.go` | `8ebe602a03cf70048c495d39881c1ef25242ddf0` | 71 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/csv/csv_test.go` | `b2f6517219ecad4c61b098150499b337d860fb80` | 134 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/excel/excel.go` | `198b681647d0037dd60ec7562aeb3b528cdae131` | 233 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/excel/excel_test.go` | `70a0c6cf21a57b15da8edd92bd507c6209fa6f3f` | 147 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/json/json.go` | `78c02db82280bc9a268ae87961e9ed070ae852ee` | 200 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/json/json_test.go` | `08b95ca48d84313fbf7dcba50bb49f4aba104146` | 235 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/sarif/sarif.go` | `137df562533618bf40fc3c7d65f962b77410d6c4` | 224 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/sarif/sarif_test.go` | `06e983fd422dd55265f8941ecd5afff61e1903b5` | 145 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/tables/plugins.go` | `c9c95c5016a9ce7354cabf0ee98e6ecfec38b10d` | 161 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/tables/tables.go` | `cdec577a766f02a56c201b76b127954ae22d2318` | 289 | Initial complete read; contract/QA reconciliation pending |
| `internal/renderers/tables/tables_test.go` | `f7dd4bded84371943f6457ea831f20e53179444b` | 142 | Initial complete read; contract/QA reconciliation pending |
| `internal/reportfile/privacy_other.go` | `6712be13d06fc37bf2756daac4a08df91c2effb6` | 5 | Initial complete read; contract/QA reconciliation pending |
| `internal/reportfile/privacy_windows.go` | `00532ec043814f2c39d7d6b3458d801f4118ef91` | 38 | Initial complete read; contract/QA reconciliation pending |
| `internal/reportfile/windows_acl_test.go` | `8c96f8f89e9f48f433067bdef47f7bab45f610da` | 43 | Initial complete read; contract/QA reconciliation pending |
| `internal/reportfile/write.go` | `161ed19e385bbdfaab5c802cd7397d741918dab5` | 44 | Initial complete read; contract/QA reconciliation pending |
| `internal/reportfile/write_test.go` | `4ef1773e1f6953a4c4c520d1b0118e509ce9e84e` | 83 | Initial complete read; contract/QA reconciliation pending |
| `internal/redact/subscription.go` | `74ca62b17ae958c03803be3b8d0550cd1a3c9852` | 25 | Initial complete read; contract/QA reconciliation pending |
| `internal/redact/subscription_test.go` | `9ce7ed9ac1ac9e8b8e25e260e4105c6b82ddd78c` | 30 | Initial complete read; contract/QA reconciliation pending |
| `internal/branding/branding.go` | `25146592055fd30be1a1f8abdabe529bc3ec48a1` | 38 | Initial complete read; contract/QA reconciliation pending |
| `internal/branding/profile.go` | `bf4ca310e6b376c421ff495e9cd58a7f5a4b8b53` | 177 | Initial complete read; contract/QA reconciliation pending |
| `internal/branding/profile_test.go` | `842b40ef97ec9f2c2a0cfbb1e3c7c147ecd6e434` | 86 | Initial complete read; contract/QA reconciliation pending |

Reproduction log `sarif-reference-reproduction.log`: 785bytes, SHA-256`70df2911c589990e7411b3ca9fa7f1e213b0bc1558304ee2092a42de7746496e`; declared assertion failure, not compile failure.

Baseline local Section3 race suite completed exit0 across all15packages before adding the new failing regression; full log 2145bytes SHA-256`6fe5512aeefcad71e613dbcaded8c8de6394af1b77d31d8ba048da4ae7c37e20`. This does not replace source review or native Windows acceptance.

## QA-S3-01 and exact regression checkpoint

Regression-only head `ed30ef46910df31593656ba0910a6242173201e8`, tree `ae80fcfaa8b6b84a56483387317fa9f2119acfef`, parent `dc1d03214e5c15abaa41338181eee7a575c00552`, contains unchanged production and one independent SARIF test. Native38722060006 failed. Windows116379246417 compiled and failed only the two declared mixed-case assertions (same-case control passed). Complete returned log51425characters/51427UTF-8bytes SHA-256`daed817edea815a68706e4ecff33dbe8a733bc9e0fc4fac18c1834855432a5b8`. Source38722060018 passed separately.

Linux116379246868 stopped at formatting, not the product assertions: new test alignment/anonymous struct layout had not been formatted before publication. This is introduced QA-S3-01, not a baseline product defect, compiler failure, Linux product reproduction or acceptance. Complete log50824characters/50826bytes SHA-256`959d5169b7545b45ce7cfd74bf8c726953816384d78c2ad7001a943e3dd8afc3`. Corrected with the pinned Go1.26.9 gofmt; expectations and production remain unchanged. Prevention: run pinned formatter and inspect its empty diff before every Go checkpoint, including regression-only WIP. A later green check cannot erase the failed trial.

The formatted local test still compiles and fails both declared mixed-case assertions while its control passes. Full Linux native reproduction remains pending; do not transfer Windows proof. NEXT: publish the formatted regression checkpoint, continue AI/region source/test and remaining contract reconciliation, then apply the documented isolated SARIF correction and compiling reverse/restoration control. Final full exact-head Linux/Windows/source QA and the Section3 stopping checkpoint remain outstanding. Initial complete reads now cover42/94baseline Go files; AI10and region42 remain. These are reading counts, not final verified package credit.
