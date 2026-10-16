# Section3 final corrective checkpoint and stopping rule, 2026-10-16

All94 baseline Go files (43production/51tests) and both independent added regressions have full reads; contract/data/integration reconciliation covers all15Section3packages. Read the [active report](SYSTEMATIC_REVIEW_SECTION3_20261015.md) and [contract map](SYSTEMATIC_REVIEW_SECTION3_CONTRACTS_20261016.md). Sections1-2 stay complete; no Section4 review is claimed.

This checkpoint implements SR-F11 after both-host compiling reproduction and the published pre-edit contract. Its only production change is the zone private duplicate-key guard using Unicode SimpleFold equivalence, matching pinned Go1.26.9 encoding/json. Long-s field overwrites and sigma duplicates reject; unique Unicode metadata remains accepted; already-rejected Kelvin continuation remains a control. The independent zone regression is byte-identical to b82de398 (blob651e4782e1af6a8cebc443b4c8755e59e7f96c02). The new compiling reversal restores strings.ToLower, requires the named behavioral failure, restores exact bytes and passes baseline. Both existing native jobs already run this Section3 runner. SR-F10 remains fixed at7384439; its original regression, canonical IDs, resource identities, dedup and fingerprint code are unchanged.

Own final-head Linux/Windows/source checks are initially PENDING AT PUBLICATION. The top current disposition in [PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142) is the authoritative exact authored commit/tree/parent/tested-preview/run/job/full-log-hash/control/capture index. Never transfer7384439 or a reproduction-head pass. QA_PROCESS permits final evidence in the PR to avoid a self-referential documentation/CI loop. Once that exact checkpoint has all required proof, Section3 is complete:3/4audit sections and37/41package directories by unweighted count, not overall project progress or defect-free certification.

RESUME ONLY by verifying live PR142/task/protected refs and that checkpoint's actual gates. Finish outstanding corrective QA if still pending; after verified success STOP for the user's re-evaluation. Do not republish an already-existing correction, replay Sections1-2, start Section4, merge, mark ready or resume Storage/features. PR142 stays draft/unmerged; protected bootstrap/core-v1 remains a6fee427e19feeef28691737ea9bc5026086efa0. No protected acceptance/push evidence exists for this unmerged audit.

Preserve AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1, APRL60eaddda76541f6adbc1c5ffa686829807e55e29, Go1.26.9/x-net0.60.0, all source captures/schema/dependencies. No Azure operations, RBAC/fixtures, production-scope substitute, pin/capture/dependency refresh or release approval. Source SKU coverage7/7 is distinct from accepted target adapters1/7; public region remains unavailable. FullCost/remaining adapters/coordinator, AR-02/load, live/freshOS, module-only advisory, human/license/security/Gate004/release remain open. Section4 and later Storage PR141 require the documented later authorization/sequence.

## Correction, reproduction and rollback index

- SR-F10: fix7384439c4058feb8c3f7a2241a6b1a8e45950163/tree434f498ee2668cfb80e61687fe01a1a50ae0d372, parent e9dd8d182054a768ff6ab65a91072f7fc11f5615. Full native38837602499/source38837602502 proof is retained below. New final head still needs its own complete gates.
- SR-F11: regression/contract b82de3982c1bdec21a0a50bcbfa21a2a2aa08069/tree1f93508d21ccc8841476967231c17a7e52ca82b0/parent7384439. Native38840067841 FAILED after compilation on Linux116777875249 and Windows116777874877 with exactly physical_zone_long_s, mappings_long_s and unknown_sigma_duplicate assertions. No build failure/race warning; later steps skipped are not acceptance. Source38840067774/job116777874590 success is separate metadata evidence, not final corrective capture credit.
- Full reproduction log Linux65940characters/65944UTF8bytes SHA2569697dd36d8dd8d0672011f09264dbe6eb18d04e2666080d00f234ff3043d6e7d; Windows52546/52548 SHA256e3156da5ccc941bbace07ff960862efc299e9eb104fef778041a2d865e2040e2.
- SR-F11 correction is this checkpoint, directly parented by b82de398; obtain its exact SHA/tree from the verified PR index. Local corrected full zone suite and vet, both Section3 compiling controls/restored baselines, pinned gofmt and immutable-regression check pass. Full native/source proof remains mandatory.
- Rollback: ordinary targeted revert commits after dependency review. Revert each affected production expression/helper case together, keep independent expectations and historical captures, and deliberately expect the corresponding regression to turn red. Later zone/SARIF controls share one runner, so do not blindly revert the entire earlier helper/workflow when successors depend on it. Never force/reset protected history or delete failing expectations.

## QA-S3-02 local orchestration correction

The first corrected zone suite was accidentally launched concurrently with its source-mutating control. Although both exited0, that suite result is discarded because compile/read timing was uncertain. The control restored exact source bytes; the complete suite was rerun serially afterward and passed, followed by vet and empty pinned formatter output. No concurrent mutation/test suite result is used for acceptance. Prevention: await control completion/restoration before any other source-consuming checks. Earlier PowerShell-PATH/SeSecurityPrivilege trials and QA-S3-01 formatting failure remain below; none was hidden or used to relax a test.

| Local log | Bytes | SHA256 |
| --- | ---: | --- |
| zone-corrected-suite.log | 66 | f62f2bb226228bf852582703947be2bbceb993d1250629626dbf900fe2b574bf |
| zone-corrected-suite-serial.log | 66 | 3b62de273177f0cd9093bc37684c125395bd587b9665e73920d580428dc730cc |
| zone-corrected-vet.log | 0 | e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 |
| section3-corrected-control.log | 177 | 7d43439b3143b352d4e97f6e773dfd71bad759340da93eb33fc2bcf335c2d103 |

Manual final diff/contract reconciliation is a bounded self-review. Earlier automated review comments are earlier-head evidence; no callable independent final Code Review tool or independent-person approval is claimed. All mandatory native controls, coverage/fuzz/security/source guards stay intact. Only the PR's verified current-head evidence can close corrective QA.

---

## Historical records below

# SR-F10 full corrective QA and SR-F11 published reproduction, 2026-10-16

SR-F10 correction7384439c4058feb8c3f7a2241a6b1a8e45950163/tree434f498ee2668cfb80e61687fe01a1a50ae0d372/parent e9dd8d182054a768ff6ab65a91072f7fc11f5615 completed native38837602499 SUCCESS on Linux116770101098 and Windows116770100656; source38837602502/job116770100933 also SUCCESS. Every mandatory step passed; only conditional Linux failure-evidence upload skipped. All326Linux and321Windows named controls contain every previous325/320 identity plus exactly the SARIF compiling reversal/restoration control. Both native full Go suites, Linux race85.2%coverage/five bounded fuzz/vet, actual CLI/default+custom packages/isolated install/branding/PowerShell/docs/provenance/inventory and scans passed. Both scans report0reachable/0imported-package/1module-only advisory; the last remains open.

Full exact returned log hashes: Linux102675characters/102679UTF8bytes SHA256309ed6c59f0c4b23fc5c1d8d6fee1d16de3936b70f2cabe61b25e928e9d88d76; Windows84106/84108 SHA2561d154ea37adc6d414f18d43da1eb8010672515939122ca17df3f1298f27f5350; source1359863/1359939 SHA256073baad731d40ce0d768ba04f98fd484853b7433c339f593911a993599efdc35.23captures/315unique contiguous chunks independently reconstructed/hashed/byte-matched plus all23fresh published files909865bytes independently fetched/matched. Pinned source/tree/APRL and tested preview5e6970718289fe0e25b2b7a28467510535776a10 exact fix tree and ordered protected-base/fix parents verified. No build/race/error marker in completed corrective logs.

Current published reproduction head is b82de3982c1bdec21a0a50bcbfa21a2a2aa08069/tree1f93508d21ccc8841476967231c17a7e52ca82b0/parent7384439. Full six-file readback/independent local tree/Denis author+committer and expected-ref publication verified. PR142 remains draft/unmerged; protected base remains a6fee427. It adds independent zone regression and pre-edit contract with production zone unchanged. SR-F11 remains UNFIXED; both-host compiling reproduction is pending. Prior SARIF green is not this changed head acceptance.

All94 baseline Go initial reads and15-package contract/integration reconciliation are now recorded in [contract map](SYSTEMATIC_REVIEW_SECTION3_CONTRACTS_20261016.md). No final Section3package/section credit until SR-F11 corrected full QA. NEXT: inspect this exact reproduction head native failures, apply isolated decoder-equivalent zone guard/control under the already published contract, verify full corrected Linux/Windows/source gates and publish the completed Section3 stopping checkpoint. STOP before Section4 for user reevaluation. No Azure, dependency/pin/capture/schema change.

---

# Section 3 continuation: SR-F11 reproduction and pre-edit correction contract, 2026-10-16

Current published correction head is 7384439c4058feb8c3f7a2241a6b1a8e45950163, tree434f498ee2668cfb80e61687fe01a1a50ae0d372, parent e9dd8d182054a768ff6ab65a91072f7fc11f5615. All eight published changed files were independently read back and blob/tree/parent/Denis identity matched. PR142 remains draft/unmerged; protected bootstrap/core-v1 remains a6fee427e19feeef28691737ea9bc5026086efa0. SR-F10 fixed locally and both-host compiling reversal/restoration controls passed; full native38837602499 remains pending, source38837602502 succeeded separately. Full log/capture reconciliation is still required.

All remaining52 baseline Go files have now received complete initial reads: AI10 and region42, including all1961 lines of other_sku_capture_test.go. Thus initial reads cover94/94 baseline files, plus added regressions; this is not final package or Section3 acceptance. Contract/data/integration reconciliation remains active. The baseline manifest below preserves exact entry blobs and reading history. Region source-only seven-family captures do not certify more than the accepted one target disk adapter; public region execution remains unavailable.

## SR-F11: zone Unicode-fold duplicate keys, P2, REPRODUCED and UNFIXED

The documented ZONE_MAPPING.md contract rejects duplicate JSON keys including case variants, while allowing unknown additional properties. Zone validateJSON uses strings.ToLower; pinned Go1.26.9 encoding/json fold.go/decode.go use Unicode SimpleFold equivalence for struct-field matching. Distinct spellings physicalZone/phyſicalZone and availabilityZoneMappings/availabilityZoneMappingſ pass the zone guard but the decoder accepts the second alias and overwrites the first value. Literal independent cases demonstrate a forged physical zone or a healthy empty mapping list. A nested unknown Σ/ς duplicate also escapes the declared duplicate-key guard. This is a target safety-guard defect, not a reason to modify the pinned source or captures.

Independent regression internal/plugins/zone/unicode_key_test.go requires zero admitted rows and the exact zone_invalid_response failure for each duplicate. Unique unknown Unicode metadata remains accepted with all original rows. Kelvin nextLinK is already rejected by the existing guard and remains a passing control; no continuation-erasure defect is claimed. Local Windows Go1.26.9 compiles and reproduces three named assertion failures. Both-host reproduction is pending.

BEFORE production edits: replace only the private zone duplicate-key equivalence with the same minimum Unicode SimpleFold orbit used by pinned encoding/json, preserving original API fields, Unicode metadata acceptance, source rows, pagination, prefix retention, cancellation, row/byte/depth/worker limits and sanitized failures. No shared JSON refactor. Keep all independent expectations. Add a compiling control that restores strings.ToLower at this one guard, requires the named Unicode assertion failure, restores exact bytes and requires the clean baseline, in both existing native Section3 jobs. Review carbon/AI/region siblings: carbon/AI reject noncanonical known fields via EqualFold, and regional JSON guards already use SimpleFold; no equivalent overwrite reproduced there.

Correction acceptance requires independent duplicate/unique controls, complete zone tests, both native compiling controls, full corrected exact-head Linux/Windows/source QA with fresh retained-capture byte/hash proof, and unchanged pins/dependencies/schema. Isolate production correction after this regression/contract checkpoint. Rollback is an ordinary targeted revert after checking helper-anchor successors; preserving the independent regression intentionally restores red. Do not delete failing expectations or force/reset history.

Local Section3 package trial passed14 packages but reportfile failed because PowerShell was initially absent from PATH; a corrected-PATH focused trial then failed Set-Acl due to executor SeSecurityPrivilege. This restricted local environment does not establish a Windows ACL product defect or a full local suite pass. Native Windows reportfile QA remains mandatory; no ACL expectation is relaxed.

NEXT: complete SR-F10 full-log QA, publish/read back this regression/contract, obtain both-host SR-F11 reproduction, apply isolated correction and control, finish contract reconciliation and corrective full QA. Stop at completed Section3 for user reevaluation before Section4. No Azure operations or dependency/pin/capture refresh.

---

# Section 3 resumed: SR-F10 correction candidate, 2026-10-16

Live PR142/task ref verified at e9dd8d182054a768ff6ab65a91072f7fc11f5615; protected bootstrap/core-v1 remains a6fee427e19feeef28691737ea9bc5026086efa0. PR142 remains draft/unmerged. The four requested handover documents and root AGENTS.md were read. This entry supersedes the earlier unfixed/paused resume paragraphs; historical evidence is retained.

SR-F10 production correction now passes local Windows Go1.26.9 regression and full SARIF package tests. Results use the published descriptor ID through the private builder argument. Canonical IDs, dedup keys, resource identities and fingerprint code are unchanged; the original regression file/expectations are byte-identical to e9dd8d1. scripts/tests/section3-regressions.py deliberately restores the old result-ID expression, requires a compiling named assertion failure, restores exact production bytes and requires the clean test to pass. Both native jobs now run this control. These are LOCAL corrective results only; full exact-head Linux/Windows/source acceptance is pending. Section3 remains incomplete; no package receives final credit.

Fresh handoff-head native38835324513 failed only the two declared mixed-case assertions after compilation on Linux116763053490 and Windows116763053350; no build-failure or race-warning marker appeared in complete returned logs. Source38835324499 succeeded separately. Later native steps skipped are not credited.

Recovery: current short-path workspace C:/workspace/ca3-20261016/repo was reconstructed from the exact e9dd8d1 archive. All773 tracked file blobs matched before edits; reconstructed Git tree fbfdf3dcb71d997518d53c2da195a497e7265ee2 matched. Its local root commit is a reconstruction marker, not remote ancestry. Go1.26.9 official Windows archive SHA256 d722201a9c0c086d1610e111c48203009af690892ed708072bd5ae20160e7a59 was verified. Shell Git networking failed, but direct fetch and authenticated GitHub publication work. APRL archive extraction first hit Windows path length; same pinned archive re-extracted to the short path. Pinned module downloads do not refresh go.mod/go.sum. None of these setup events is a product finding.

Parent for this isolated correction: e9dd8d182054a768ff6ab65a91072f7fc11f5615. Identify its resulting commit/tree from live PR142 and read-back before resuming. Rollback is an ordinary targeted revert of the production/helper/workflow changes together, after dependency review; that intentionally restores the known red original regression, whose expectations must not be removed. No protected reset, pin/capture/schema/dependency change or Azure action.

NEXT: inspect fresh exact-head native/source controls and full outcomes; finish AI/region complete reads plus all94-file contract/data/integration reconciliation, independently fix any substantiated finding, obtain complete corrective QA and publish the completed Section3checkpoint. STOP before Section4 for the user.

Local log bytes and SHA256:

- sarif-baseline-proxy: 785 bytes, 42617cc4fcfac2dfbcaf6cb9ba57b31571af5be569b62d6f901d64f9b2b48935
- section3-control: 93 bytes, 732e82dba98a06831cc2db36a587f640066ddc8d64cdacf99c55382425bf7c10
- sarif-suite: 69 bytes, 965b4edb6305ea551f55e400142cec634cce87998840ef0b7ce5122b8e58af8b

---

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
| `internal/plugins/aigov_projection.go` | `408b7cda42858faeab5b18b9fcba43cf8aa80c01` | 112 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov_projection_test.go` | `c016c61fe54c6d006a3b186374245a4018fd4289` | 81 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/carbon_projection.go` | `626c3b24d2e2adb7dec5a53c0fef7f8d216efd89` | 95 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/carbon_projection_test.go` | `a2b0221abaa4b345d67afd7f7877829beb2562d5` | 115 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/registry.go` | `0e7e6810d481c5c7e806ee6469daf3c03f3f33c4` | 186 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/registry_test.go` | `202b7a54f20738d90cf1f7650280ecd235c2a96b` | 55 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/servicehealth_projection_test.go` | `ff9b4db92f36855a592f79a7c23ecbf08fba5170` | 91 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/sqleol_projection.go` | `a925e0545eaa0e9b4be51a58705d46c0004f56fe` | 70 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/sqleol_projection_test.go` | `cb8616b373409d017318f1d9f5bf20131a4492e0` | 87 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/zone/scanner.go` | `1ba211bad8d2b0558b91481d177acb35417a568d` | 314 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/zone/scanner_test.go` | `842fa2d7be23836ec7277bc5dcd2541f9e72a836` | 331 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/servicehealth/query.go` | `7624c0a3d23788f8c3e1565b2f35f55102b7e48f` | 72 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/servicehealth/scanner.go` | `ffdc63d0b4669ad6ab23908deea1e69a4785ccd0` | 328 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/servicehealth/scanner_test.go` | `1ef072fe74f1fd0024a1aaf56f38a162cead182c` | 355 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/sqleol/query.go` | `bb408e06215f9fb2d667dbc6372033b5e83d7858` | 7 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/sqleol/scanner.go` | `2abf52dfc1ca3fda34e97dda898b7aa1d8631cc0` | 289 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/sqleol/scanner_test.go` | `5ad2c6c14c6ecf23820ea989a8cfb8b62dc34ae1` | 418 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/carbon/projection.go` | `868599a28286a303671cc8a76ecd06d90f13d1e0` | 158 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/carbon/projection_test.go` | `105315a948c25c97be421b699e13eec17c3234e1` | 190 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/carbon/scanner.go` | `d0aff9de47afad2c1a1b8c3fcd3ff9819ac5b1fd` | 447 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/carbon/scanner_test.go` | `f8e1e5dca08bcc1706a5a48ff804bf45044d326e` | 625 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/decoder.go` | `c31875f5576165af668ce95eaa2408767321a203` | 454 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/decoder_test.go` | `76e9df63a13a056b7799877dbeaf4b0282c71b3c` | 143 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/discovery.go` | `63e8d826cd8e4cb4fd51a6ea86a5c4cc40a86660` | 314 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/discovery_test.go` | `9f32084ed4059b885f158ef28044403b1c206503` | 633 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/execution.go` | `558d79cdb2526dcff5256741061737c3a2bc8f7b` | 70 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/execution_test.go` | `e63932442c8fce7f51b4f9a88657d8ffd49c0bd6` | 113 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/projection.go` | `90ceacb770826aff716c4d1ebc2e106ba0b6211c` | 312 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/projection_test.go` | `20454b562bb0645b53157958e97f249e5de66acb` | 345 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/scanner.go` | `c00e0ee00be57cf8bc7256322fbc62722f64a768` | 463 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/aigov/scanner_test.go` | `3da3767d71fefb91d76576f712044c742bc929be` | 749 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/auxiliary.go` | `15ab8323b5b4dd7974105eb15994d777f160ad35` | 260 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/auxiliary_test.go` | `1526c4508c2a25a009305b9accf7eb5daa4c05b9` | 342 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/availability.go` | `94360cfbe6a59a3debdf1a8f002cf7b062f8e264` | 308 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/availability_capture_test.go` | `a212b8f605c28d8eac7137c893de3dc4f89ddaa5` | 146 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/availability_test.go` | `695105d71c28c64980de19334f1a4370c3057869` | 329 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/compute_disk_sku_collector.go` | `d7a5b8b2e1facd7a36b157acc9eadd4d91a57465` | 633 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/compute_disk_sku_collector_test.go` | `76ce22e83bc9fe9fe4bdb51520782a0fed4cc97e` | 652 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/compute_disk_sku_review_test.go` | `8789739cb653347be8b596e9fccb7e4836410212` | 65 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/compute_sku_capture_test.go` | `2ef251506b2c6752a02685016aacee513bdf1978` | 193 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/cost.go` | `dc7c43cdcfc46c2fa67e00aaafbc5cf0faeec149` | 235 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/cost_enrichment.go` | `94a9cc17e8d70a02852274a69bfd6a6f5eb193d5` | 278 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/cost_enrichment_capture_test.go` | `92cdcf06990aa8a5139c1966ee9b2ff9a62eec84` | 96 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/cost_enrichment_test.go` | `c64816dd1daa8be6695b2d10076c267854955a7b` | 468 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/cost_test.go` | `48712680a4eb91785d889a3c2796bfb3b6e623a3` | 251 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/inventory.go` | `6157171769aaca1a4e78a026a4dc7f02a610bb87` | 134 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/inventory_aggregation.go` | `9828cf3a30056e2456028b6e746e1e23afd5782c` | 197 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/inventory_aggregation_test.go` | `b23fbf28124fdb72a78d5b86c85fff897e4cd30c` | 316 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/inventory_capture_test.go` | `c7f85db283725e33b2d89a4df85dfdf6e8f8638e` | 107 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/inventory_test.go` | `ac97771893275ad42fc5272a435e8a9dae57d718` | 294 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/latency.go` | `1a4a037e8d6c1c80aed1d7d51c66a5d5d1ae1090` | 217 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/latency_capture_test.go` | `877e215f1212f85096033acf9636618620479d13` | 99 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/latency_test.go` | `60ea39a005cf1f1aff3f13c3889700706fc46dc2` | 258 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/other_sku_capture_test.go` | `c73799952e61070813757283c13dc66079408ee7` | 1961 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/primary.go` | `3a7fb0b9601a525f91c385d0016ce0b8ff532303` | 322 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/primary_test.go` | `6219219eca501b00af02a8af6fb23bb4a51c5d78` | 335 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/provider_locations_capture_test.go` | `87112d5d3876fba242dbd2e0bb6ec9e24b35085f` | 183 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/provider_locations_collector.go` | `8c85347d96488f2d4b1440da3f98a64ab3657413` | 324 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/provider_locations_collector_test.go` | `9d1180be0e50e0f7fca5d30bbda00cb25eeb3ca6` | 547 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/quota_capture_test.go` | `78a92620225e5e317e7af84dba73d6971828c326` | 80 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/quota_collector.go` | `e2a44936371ec66600f33e2378b80ec689ee8b3f` | 352 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/quota_collector_test.go` | `70a85e863018ef976e4f54a92318f3503d26b3f1` | 543 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/quota_runtime.go` | `a2be4a680e2c3d705ed1198d1cbb509a04e44713` | 234 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/quota_runtime_test.go` | `5ac53f03b1614a386b37fd92ed8071885b3efab2` | 353 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/reservation_capture_test.go` | `31377c86e8b54423d799279c6d29ca3870cc1e72` | 114 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/reservation_collector.go` | `c5b9c83e5d41be27798d7ab208f2568cca437c49` | 357 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/reservation_collector_test.go` | `d617b44f6701b59234fbfb3866c08035a26fdf02` | 566 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/reservation_runtime.go` | `bc6e23eca781e5b5ad81be2f416e2f9c36ec58c8` | 269 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/reservation_runtime_test.go` | `1ec9bffe0cddf1e624c2b92be434ed8ee7952370` | 421 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/service.go` | `e2e53165c2a4cb597c2bc36c239557189ad68ccd` | 356 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/service_test.go` | `f683d274a5b208d4d54a7d1081de3f43aaba1581` | 377 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/vm_quota_capture_test.go` | `18476fb14c4385d67e7d7fd16ccabbcc7ba8b52c` | 77 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/plugins/region/vm_quota_collector_test.go` | `47546d5e6b3d79e9d0baa230f6739a3b75fbfa72` | 218 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/csv/csv.go` | `8ebe602a03cf70048c495d39881c1ef25242ddf0` | 71 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/csv/csv_test.go` | `b2f6517219ecad4c61b098150499b337d860fb80` | 134 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/excel/excel.go` | `198b681647d0037dd60ec7562aeb3b528cdae131` | 233 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/excel/excel_test.go` | `70a0c6cf21a57b15da8edd92bd507c6209fa6f3f` | 147 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/json/json.go` | `78c02db82280bc9a268ae87961e9ed070ae852ee` | 200 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/json/json_test.go` | `08b95ca48d84313fbf7dcba50bb49f4aba104146` | 235 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/sarif/sarif.go` | `137df562533618bf40fc3c7d65f962b77410d6c4` | 224 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/sarif/sarif_test.go` | `06e983fd422dd55265f8941ecd5afff61e1903b5` | 145 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/tables/plugins.go` | `c9c95c5016a9ce7354cabf0ee98e6ecfec38b10d` | 161 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/tables/tables.go` | `cdec577a766f02a56c201b76b127954ae22d2318` | 289 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/renderers/tables/tables_test.go` | `f7dd4bded84371943f6457ea831f20e53179444b` | 142 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/reportfile/privacy_other.go` | `6712be13d06fc37bf2756daac4a08df91c2effb6` | 5 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/reportfile/privacy_windows.go` | `00532ec043814f2c39d7d6b3458d801f4118ef91` | 38 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/reportfile/windows_acl_test.go` | `8c96f8f89e9f48f433067bdef47f7bab45f610da` | 43 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/reportfile/write.go` | `161ed19e385bbdfaab5c802cd7397d741918dab5` | 44 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/reportfile/write_test.go` | `4ef1773e1f6953a4c4c520d1b0118e509ce9e84e` | 83 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/redact/subscription.go` | `74ca62b17ae958c03803be3b8d0550cd1a3c9852` | 25 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/redact/subscription_test.go` | `9ce7ed9ac1ac9e8b8e25e260e4105c6b82ddd78c` | 30 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/branding/branding.go` | `25146592055fd30be1a1f8abdabe529bc3ec48a1` | 38 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/branding/profile.go` | `bf4ca310e6b376c421ff495e9cd58a7f5a4b8b53` | 177 | Full read and contract review complete; exact-head QA disposition in PR142 |
| `internal/branding/profile_test.go` | `842b40ef97ec9f2c2a0cfbb1e3c7c147ecd6e434` | 86 | Full read and contract review complete; exact-head QA disposition in PR142 |

Reproduction log `sarif-reference-reproduction.log`: 785bytes, SHA-256`70df2911c589990e7411b3ca9fa7f1e213b0bc1558304ee2092a42de7746496e`; declared assertion failure, not compile failure.

Baseline local Section3 race suite completed exit0 across all15packages before adding the new failing regression; full log 2145bytes SHA-256`6fe5512aeefcad71e613dbcaded8c8de6394af1b77d31d8ba048da4ae7c37e20`. This does not replace source review or native Windows acceptance.

## QA-S3-01 and exact regression checkpoint

Regression-only head `ed30ef46910df31593656ba0910a6242173201e8`, tree `ae80fcfaa8b6b84a56483387317fa9f2119acfef`, parent `dc1d03214e5c15abaa41338181eee7a575c00552`, contains unchanged production and one independent SARIF test. Native38722060006 failed. Windows116379246417 compiled and failed only the two declared mixed-case assertions (same-case control passed). Complete returned log51425characters/51427UTF-8bytes SHA-256`daed817edea815a68706e4ecff33dbe8a733bc9e0fc4fac18c1834855432a5b8`. Source38722060018 passed separately.

Linux116379246868 stopped at formatting, not the product assertions: new test alignment/anonymous struct layout had not been formatted before publication. This is introduced QA-S3-01, not a baseline product defect, compiler failure, Linux product reproduction or acceptance. Complete log50824characters/50826bytes SHA-256`959d5169b7545b45ce7cfd74bf8c726953816384d78c2ad7001a943e3dd8afc3`. Corrected with the pinned Go1.26.9 gofmt; expectations and production remain unchanged. Prevention: run pinned formatter and inspect its empty diff before every Go checkpoint, including regression-only WIP. A later green check cannot erase the failed trial.

The formatted local test still compiles and fails both declared mixed-case assertions while its control passes. Full Linux native reproduction remains pending; do not transfer Windows proof. NEXT: publish the formatted regression checkpoint, continue AI/region source/test and remaining contract reconciliation, then apply the documented isolated SARIF correction and compiling reverse/restoration control. Final full exact-head Linux/Windows/source QA and the Section3 stopping checkpoint remain outstanding. Initial complete reads now cover42/94baseline Go files; AI10and region42 remain. These are reading counts, not final verified package credit.
