# Section 3 systemic review — active record

Entry baseline: `798d7be21e393f1e85503303e3d8eb83d311913d`. No Section 3 completion or new product defect is claimed at preparation.

## Coverage and method

All15packages/94Go files (43production/51tests), with immediate data/contract boundaries. Read complete source/tests; execution alone is not code-review coverage. Independently reproduce material defects with literal/source-grounded expectations, fix in isolated commits, verify uncertain/version-sensitive contracts against pinned code or official primary documentation. Region calculation/source work is distinct from incomplete public runtime adapters.

## Findings and decisions

No confirmed Section 3 finding yet. Reportfile, redact and branding have initial complete source/test reads; integration/source-contract reconciliation remains pending, so manifest coverage is pending.

## QA and recovery

Entry native38652456544 (Linux116154267913, Windows116154267785) and source38652456593/job116154266947 report success; full logs pending reconciliation. Local Go1.26.9 was automatically fetched by the existing1.26.8launcher, without repository changes. Local checks do not replace native Windows/live evidence.

NEXT: continue full source/test review and update this record/handover after coherent slices, then reproduce/fix/QA and stop before Section 4. Keep PR142 draft/unmerged. No Section 3 production changes yet; ordinary revert of documentation preparation must account for successors. Previous fix/parent and failure evidence remains in Section 2 report/journal. A manifest preparation assertion caught an incorrect handwritten42/52split expectation before any file write; actual inventory below is authoritative (94total,43production,51tests), not a product defect.

## Baseline file coverage

| File | Baseline Git blob | Lines | Review state |
| --- | --- | ---: | --- |
| `internal/plugins/aigov_projection.go` | `408b7cda42858faeab5b18b9fcba43cf8aa80c01` | 112 | Pending |
| `internal/plugins/aigov_projection_test.go` | `c016c61fe54c6d006a3b186374245a4018fd4289` | 81 | Pending |
| `internal/plugins/carbon_projection.go` | `626c3b24d2e2adb7dec5a53c0fef7f8d216efd89` | 95 | Pending |
| `internal/plugins/carbon_projection_test.go` | `a2b0221abaa4b345d67afd7f7877829beb2562d5` | 115 | Pending |
| `internal/plugins/registry.go` | `0e7e6810d481c5c7e806ee6469daf3c03f3f33c4` | 186 | Pending |
| `internal/plugins/registry_test.go` | `202b7a54f20738d90cf1f7650280ecd235c2a96b` | 55 | Pending |
| `internal/plugins/servicehealth_projection_test.go` | `ff9b4db92f36855a592f79a7c23ecbf08fba5170` | 91 | Pending |
| `internal/plugins/sqleol_projection.go` | `a925e0545eaa0e9b4be51a58705d46c0004f56fe` | 70 | Pending |
| `internal/plugins/sqleol_projection_test.go` | `cb8616b373409d017318f1d9f5bf20131a4492e0` | 87 | Pending |
| `internal/plugins/zone/scanner.go` | `1ba211bad8d2b0558b91481d177acb35417a568d` | 314 | Pending |
| `internal/plugins/zone/scanner_test.go` | `842fa2d7be23836ec7277bc5dcd2541f9e72a836` | 331 | Pending |
| `internal/plugins/servicehealth/query.go` | `7624c0a3d23788f8c3e1565b2f35f55102b7e48f` | 72 | Pending |
| `internal/plugins/servicehealth/scanner.go` | `ffdc63d0b4669ad6ab23908deea1e69a4785ccd0` | 328 | Pending |
| `internal/plugins/servicehealth/scanner_test.go` | `1ef072fe74f1fd0024a1aaf56f38a162cead182c` | 355 | Pending |
| `internal/plugins/sqleol/query.go` | `bb408e06215f9fb2d667dbc6372033b5e83d7858` | 7 | Pending |
| `internal/plugins/sqleol/scanner.go` | `2abf52dfc1ca3fda34e97dda898b7aa1d8631cc0` | 289 | Pending |
| `internal/plugins/sqleol/scanner_test.go` | `5ad2c6c14c6ecf23820ea989a8cfb8b62dc34ae1` | 418 | Pending |
| `internal/plugins/carbon/projection.go` | `868599a28286a303671cc8a76ecd06d90f13d1e0` | 158 | Pending |
| `internal/plugins/carbon/projection_test.go` | `105315a948c25c97be421b699e13eec17c3234e1` | 190 | Pending |
| `internal/plugins/carbon/scanner.go` | `d0aff9de47afad2c1a1b8c3fcd3ff9819ac5b1fd` | 447 | Pending |
| `internal/plugins/carbon/scanner_test.go` | `f8e1e5dca08bcc1706a5a48ff804bf45044d326e` | 625 | Pending |
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
| `internal/renderers/csv/csv.go` | `8ebe602a03cf70048c495d39881c1ef25242ddf0` | 71 | Pending |
| `internal/renderers/csv/csv_test.go` | `b2f6517219ecad4c61b098150499b337d860fb80` | 134 | Pending |
| `internal/renderers/excel/excel.go` | `198b681647d0037dd60ec7562aeb3b528cdae131` | 233 | Pending |
| `internal/renderers/excel/excel_test.go` | `70a0c6cf21a57b15da8edd92bd507c6209fa6f3f` | 147 | Pending |
| `internal/renderers/json/json.go` | `78c02db82280bc9a268ae87961e9ed070ae852ee` | 200 | Pending |
| `internal/renderers/json/json_test.go` | `08b95ca48d84313fbf7dcba50bb49f4aba104146` | 235 | Pending |
| `internal/renderers/sarif/sarif.go` | `137df562533618bf40fc3c7d65f962b77410d6c4` | 224 | Pending |
| `internal/renderers/sarif/sarif_test.go` | `06e983fd422dd55265f8941ecd5afff61e1903b5` | 145 | Pending |
| `internal/renderers/tables/plugins.go` | `c9c95c5016a9ce7354cabf0ee98e6ecfec38b10d` | 161 | Pending |
| `internal/renderers/tables/tables.go` | `cdec577a766f02a56c201b76b127954ae22d2318` | 289 | Pending |
| `internal/renderers/tables/tables_test.go` | `f7dd4bded84371943f6457ea831f20e53179444b` | 142 | Pending |
| `internal/reportfile/privacy_other.go` | `6712be13d06fc37bf2756daac4a08df91c2effb6` | 5 | Pending |
| `internal/reportfile/privacy_windows.go` | `00532ec043814f2c39d7d6b3458d801f4118ef91` | 38 | Pending |
| `internal/reportfile/windows_acl_test.go` | `8c96f8f89e9f48f433067bdef47f7bab45f610da` | 43 | Pending |
| `internal/reportfile/write.go` | `161ed19e385bbdfaab5c802cd7397d741918dab5` | 44 | Pending |
| `internal/reportfile/write_test.go` | `4ef1773e1f6953a4c4c520d1b0118e509ce9e84e` | 83 | Pending |
| `internal/redact/subscription.go` | `74ca62b17ae958c03803be3b8d0550cd1a3c9852` | 25 | Pending |
| `internal/redact/subscription_test.go` | `9ce7ed9ac1ac9e8b8e25e260e4105c6b82ddd78c` | 30 | Pending |
| `internal/branding/branding.go` | `25146592055fd30be1a1f8abdabe529bc3ec48a1` | 38 | Pending |
| `internal/branding/profile.go` | `bf4ca310e6b376c421ff495e9cd58a7f5a4b8b53` | 177 | Pending |
| `internal/branding/profile_test.go` | `842b40ef97ec9f2c2a0cfbb1e3c7c147ecd6e434` | 86 | Pending |
