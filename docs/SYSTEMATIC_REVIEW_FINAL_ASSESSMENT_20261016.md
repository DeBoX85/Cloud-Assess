# Systematic review final assessment, 2026-10-16

Status at preparation: Section 4 review/reconciliation is complete; SR-F12 is corrected locally after both-host reproduction; its own full corrective QA is pending. This report is not completion or acceptance evidence until the exact candidate's Linux, Windows and pinned-source outcomes are verified in [PR142](https://github.com/DeBoX85/Cloud-Assess/pull/142). PR142 remains draft/unmerged and protected `bootstrap/core-v1` remains `a6fee427e19feeef28691737ea9bc5026086efa0`.

## Goal and bounded conclusion

The code remains aligned with the intended read-oriented Azure assessment engine and material behavioral equivalence to pinned AZQR, with adjustable immutable presentation branding. Excel and canonical JSON remain first-class outputs. Collection, explicit health/completeness, owned canonical results, format projection and severity decisions remain separate concerns. The review found concrete reliability and evidence-integrity defects that require corrections and retained independent regressions; it does not establish defect-free software or full feature/live parity.

The accepted baseline contains 251 Go files across 41 package directories: 112 production and 139 tests. The four review scopes account for every baseline path, with no missing package or baseline Go path. The audit has added 12 Go regression files, bringing the current total to 263 without deleting original expectations. Section 4 also reconciles 62 tracked script/workflow/support files, plus the new corrective helper, authored configuration/examples, provenance and current documentation. Per-file baseline blobs and initial-read history remain in the section manifests.

| Scope | Packages | Baseline Go files | Review record |
| --- | ---: | ---: | --- |
| Foundations, scope and execution | 9 | 95 | [Section 1](SYSTEMATIC_REVIEW_SECTION1_20261014.md), original manifest in [review plan](SYSTEMATIC_MODULE_REVIEW_20261014.md) |
| Assessment logic and collectors | 13 | 51 | [Section 2](SYSTEMATIC_REVIEW_SECTION2_20261014.md) |
| Plugins, reporting and platform protections | 15 | 94 | [Section 3](SYSTEMATIC_REVIEW_SECTION3_20261015.md), [contract map](SYSTEMATIC_REVIEW_SECTION3_CONTRACTS_20261016.md) |
| Tooling, delivery and integrated reconciliation | 4 | 11 | [Section 4](SYSTEMATIC_REVIEW_SECTION4_20261016.md), [contracts and manifest](SYSTEMATIC_REVIEW_SECTION4_CONTRACTS_20261016.md) |

These are unweighted review counts, not an overall project completion percentage or a quantified probability of correctness. Earlier section QA is preserved as exact historical proof; it cannot certify a changed final head.

## Finding disposition and rollback index

| Finding | Material problem and correction | Correction ancestry |
| --- | --- | --- |
| SR-F01 | ARG terminal caller context lost to scheduling/sibling or plain query errors; preserve parent context priority while retaining ordinary initiating errors | `5dd97b0`, then `e65e66c` |
| SR-F02 | Non-unique inventory paging sort can duplicate/omit tied rows; append projected resource ID as a deterministic tie-breaker | `9d1872b` |
| SR-F03 | ARG accepted undocumented 2xx responses as completed queries; enforce documented HTTP200 | `1f93442` |
| SR-F04 | Stage Runner could lose caller terminal state after task execution; check and prioritize caller context | `3c5860a` |
| SR-F05 | Scanner registry returned shared nested resource-type slices; copy them at public accessors | `b871736` |
| SR-F06 | Cost and Diagnostics workers could turn caller interruption into data/plain failure; enforce terminal priority through inference/result boundaries | `0c21ac4` |
| SR-F07 | Diagnostics accepted trailing batch JSON/content; require EOF after one document, preserving read causes/body closure | `fe39724` |
| SR-F08 | Diagnostics admitted settings from outside its requested batch; enforce requested parent membership | `d12baf3` |
| SR-F09 | Cost accepted undocumented completion statuses; preserve HTTP object and enforce200/204 | `cece553`, naming/control correction `4284ffe` |
| SR-F10 | SARIF mixed-case result rule IDs did not reference the published descriptor spelling; pass actual descriptor identity privately | `7384439` |
| SR-F11 | Zone lowercase key guard missed decoder-equivalent Unicode duplicates; use pinned Go SimpleFold equivalence | `98b8915` |
| SR-F12 | Comparator overwrote conflicting normalized records and could return false PASS; published independent regression/contract requires permanent noncomparability and sanitized diagnostic | Regression `10efcd4`; private guard and compiling control in this checkpoint, direct parent `10efcd4` |

The first eleven groups have their earlier complete section-specific corrective proof. SR-F12 has compiling reproduction on Linux and Windows, with the corrected local package suite and compiling reversal/restoration control passing. The [active report](SYSTEMATIC_REVIEW_SECTION4_20261016.md) and [change journal](SYSTEMATIC_REVIEW_CHANGELOG_20261014.md) retain actual final corrective SHA, exact parent/tree, full logs, failures and later disposition. Findings are groups, not a count of all possible faults; SR-F06 covers two collectors.

Rollback is a new targeted corrective/revert commit after inspecting current dependent changes. Preserve independent regressions, historical capture bytes and records. Diagnostics and Cost fixes share files; later helper anchors also depend on earlier corrections. Reverting a production fix intentionally makes its independent regression red. Do not blindly revert entire shared helper/workflow history, delete assertions, reset protected history or imply an untested rollback is healthy. No schema/data migration or Azure rollback belongs to this offline audit.

## Final corrective acceptance requirements

The final code/documentation candidate needs complete exact-head Linux and Windows native jobs and the independent pinned-source job, with logs inspected. Require every previous compiling control identity plus the new SR-F12 control, valid named assertion failure under deliberate reversal, exact restoration and clean rerun. Prior section evidence is not transferred to a changed head.

Required verification includes formatting/module/inventory/provenance, full Go suites, Linux race/vet/coverage and five bounded fuzz runs, actual native CLI/preflight/default and custom branding/report checks, candidate integrity and isolated installed execution, PowerShell and authored documentation, security scans and all existing region/source controls. Windows controlled ACL and file-lock fixtures are native evidence, not arbitrary-directory or fresh-OS certification.

Source proof must reconstruct all 23 captures from 315 unique contiguous chunks, verify declared hashes and exact retained/published bytes, and confirm unchanged source/tree/APRL identities. Verify the tested preview has the exact candidate tree and expected base/head parents. Record full returned-log hashes without discarding line endings/BOM. Re-read current head/base and PR draft/unmerged state after QA. The final live PR evidence index supplies the actual outcome to avoid a self-referential evidence-only commit loop.

## Residual engineering and validation limits

- Core and plugin contracts intentionally differ where the pinned source can silently lose data, infer findings after interruption or accept unsafe input. The original source captures remain unchanged and corrections are separately described.
- Core Graph business-row admission, manually constructed canonical core domain objects and comparator inputs are not universal schema/identity validators. Positional source helpers are not general ARM parsers.
- A successful visible scope is not proof of intended estate/RBAC coverage. Missing requested explicit subscriptions fail, but all-visible/MG visibility still needs independent expected membership.
- Diagnostics presence does not certify category/destination enablement; warned non-success findings and historical request correlation retain explicit uncertainty. Same-origin paging/whole-scan volume and context-ignoring injected code are not universally bounded by a hard kill.
- Core Cost remains subscription-wide under resource/RG/tag filters. Arc SQL numeric-vCore shape remains unresolved. Policy/Defender source dedup is not a latest-history guarantee.
- SARIF retains stable resource identities. Other output masking is subscription-ID redaction, not anonymization. Multi-file exports and package/sidecar publication are not crash-atomic transactions.
- Development packages carry integrity and unsigned build evidence. Fresh OS, signed publisher provenance, package-specific licensing, supported-platform decisions and operational/security approval remain release work.
- The module-only advisory remains open. Only the final current scan can establish current reachable/imported-package results; no dependency refresh or suppression is authorized.
- Source SKU captures cover 7/7 families; accepted target SKU adapters cover 1/7. Internal region libraries are not public region execution or live deployability/capacity proof.

FullCost/history/retail collection, the six remaining SKU adapters, coordinator/public region integration, ancillary source commands/B6, AR-02/load, representative live/sovereign/credential behavior, non-empty Policy/Defender Recommendations, Arc SQL shape, nested non-production management-group traversal, fresh OS, hosted maintenance and human/license/security/Gate004/release decisions remain open. Existing historical live passes are preserved and are not invalidated without evidence of affected inputs.

## Next development checkpoint

Finish this audit's exact-head corrective QA and publish its final evidence/disposition. Keep PR142 draft/unmerged and stop at the completed audit checkpoint for re-evaluation; do not infer merge or release approval.

After the audit disposition, resume the Storage pre-implementation contract and independent tests on draft [PR141](https://github.com/DeBoX85/Cloud-Assess/pull/141), currently `6cd002cb77fcead75a602b06a75625a178a84e49`. No Storage collector is implemented there. Resolve grouping repeated names across kinds/regions, kind-aware projection, display/stage locations, paging and absent optional metadata against the inspected schema/example before production edits. Then continue remaining SKU collectors and region integration in bounded, independently verified slices.

No Azure operations, fixture/RBAC changes, dependency refresh, source pin/capture/schema refresh, protected merge or release publication was performed by this audit continuation. The review is a bounded self-review with independent expected inputs and compiling controls, not independent-person approval.
