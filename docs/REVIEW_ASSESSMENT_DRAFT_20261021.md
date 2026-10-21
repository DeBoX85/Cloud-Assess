# Assessment working draft, 2026-10-21

INCOMPLETE REVIEW, not final disposition, product acceptance or release approval. The [review progress](REVIEW_PROGRESS_20261021.md) records read/evidence coverage and remaining work. The final report must reconcile remaining data/reference contracts, finish its own review-tree QA and collect unresolved questions at the end. This working draft makes the current assessment concrete without calling the audit finished.

## A. Intended outcome

Cloud-Assess is a read-oriented Azure assessment engine: authenticate, resolve authorised scope, discover/filter inventory, assess recommendations and platform signals, preserve explicit stage health/completeness, and produce useful human and machine reports. Denis's [AZQR fork](REFERENCE_POLICY.md) is the pinned initial feature/function/capability comparison. Behavioral parity is required; identical CLI/config/package identity is not. Deliberate safety corrections and further improvements need their own explicit contracts and independent expectations. Additional primary sources and comparable projects inform that work.

The existing core-v1 acceptance set and the broader AZQR parity goal must remain distinct. A smaller core acceptance set does not remove deferred parity requirements. Source bugs are observations, not requirements to reproduce unsafe behavior. Catalogue presence is not deployment eligibility, successful empty data is not failed retrieval, and visible scope is not proof of intended estate visibility.

## B. Actual progress

The foundation is substantial: the core pipeline, reports, branding/packages, five public plugins and many region calculations/collectors exist. Region integration remains unfinished. The protected baseline remains a6fee427; subsequent verified draft implementations are unmerged candidates.

Tasks1-5 and security/task6 prerequisites are complete offline. Task7 retail has its own corrected183QA and formal evidence closure; the recovered local executions remain historical. The approved narrow Go1.26.10continuation185has its own successful native/source QA,622Linux/617Windows conserved control identities,87.6%coverage against the unchanged75%floor, and zero reachable findings in both scans. Three imported-not-called and one module-only advisory remain OPEN. These measures are validation evidence, not a completion percentage or assurance that no defect remains.

Task8's three recovered drafts/50scenarios remain provisional. The audit has read all326first-party Go files, scripts/workflows and first-party rules. Data/reference/applicable documentation reconciliation is still in progress; complete raw-source reads and reversible semantic-view reads are tracked separately. All45initial PR bodies and the seven-day allocated logs are retained/reviewed, including failed steps inside cancelled runs. No service-outage cause is inferred from cancellation or scratch loss.

## C and F. Remaining work and roadmap

| Order | Deliverable | Exit evidence |
| --- | --- | --- |
| Current audit | Finish remaining input/reference contracts, decision synthesis, risk/blocker register, research recommendations and final questions | Exact read scope/limits, retained findings and own final review-tree QA; no inherited185green |
| Task8 | Freeze required-kind identity/evidence contract and50independent scenarios | All-required Storage kinds, exact Cognitive kind, no invented SQL coordinates, unsupported VMSS scoring; contracts reviewed before production |
| Task9 | Preserve private inventory identity/coordinates and admit evidence correctly | Published strict consumer and BOTH native missing-API observations before production, independent expectations, own complete QA/rollback |
| Task10 | Connect existing components through a bounded owned per-run coordinator | Tasks5/7/9; health, cancellation, ownership, endpoint/credential separation and completeness propagated end to end |
| Task11 | Expose supported region CLI/registry execution | Early validation, defaults/flags, standalone/mixed/scanner selections and actual command evidence |
| Task12 | Complete outputs/health/exit agreement and required ancillary parity | FullCost disposition and all required tables; format-specific privacy/identity contracts preserved |
| Task13 | Exercise representative synthetic workloads | Reproducible scale, paging, limits, cancellation and injected permissions; measured budgets without a live-scale claim |
| Task14 | Exercise installed default/custom packages and rollback | Native package execution; fresh-OS/customer operational limits explicitly retained |
| Task15 | Reconcile readiness and residual obligations | Exact supported path/acceptance set, independent/live/security/license/Gate004 decisions; no automatic merge or release |

Existing collectors must be reused through their contracts rather than rebuilt from stale missing-adapter counts. Frozen source captures/oracles/control identities/pins and all old draft heads remain intact. Material coding starts from verified185product lineage with current handover explicitly carried, never this review/evidence tree. Tasks may require several coherent slices; their count does not imply equal effort or a calendar ETA.

## D. Risks and mitigations

| Risk | Verified basis / practical consequence | Current mitigation and required follow-up |
| --- | --- | --- |
| Ambiguous tag filters | R7-F01: conflicting case-equivalent keys produce nondeterministic include/exclude decisions; also present in pinned source | Avoid conflicting normalized keys; specify deterministic rejection/duplicate policy and verify both hosts before using ambiguous filters |
| Comparator evidence concealment | R7-F03: duplicate JSON fields can overwrite populated/partial evidence and yield a false equivalent result | Use provenance-bound generated comparison inputs with independent ambiguity checks; contract strict bounded admission before treating arbitrary inputs as trusted |
| Mutation proof admission | R7-F04: four SQL/scope drivers admit actual named Go panic output as assertion proof | No historical panic occurrence established; require explicit rejection of panic/build/setup/timeout/infra outcomes in future proof admission, preserving restoration/baseline evidence |
| Scope and evidence integration | Region's supported public path and owned coordinator are absent; SQL/VMSS catalogue data cannot establish eligibility | Complete tasks8-12 with explicit unknown/partial/unsupported states and independent end-to-end expectations |
| Handover drift / recovery | Append-heavy current/pending/STOP narratives, raw-byte archive normalization incident and corrected inventory count | One compact current authority, explicit successor routing, raw-byte manifests/readback, inspect uncertain operations before retry; proposed semantic consistency checks |
| Resource use and lifecycle | Global volume/workload acceptance and several caller-impact observations remain open | Task13 must define/measure representative budgets; do not infer universal memory or hard-termination guarantees from component caps |
| Dated rule/source content | Frozen captures/rules preserve old SLA, retirement, pricing and source behavior | Keep provenance/compatibility separate from current platform correctness; changes require primary evidence and separate review, without unrelated pin refresh |
| Release assurance | Residual advisories, independent review/backup, live permissions/clouds/workloads, fresh OS/operational rollback and Gate004 remain open | Preserve explicit owners/evidence gaps; offline automation cannot close these obligations |

R7-F02 is reclassified: the supplemental history plan explicitly defines root depth0, and the observed empty-container response satisfies that contract. The audit probe assumed a different root-inclusive convention; its red result is retained as an unsupported review expectation, not a product defect or blocker. Generic HTTP error-body volume, positional ID helper callers, canonical sort ties and redaction fallback remain observations pending demonstrated impact. They are not promoted to proven vulnerabilities or unmitigated critical blockers without evidence.

## E. Blockers and decision boundaries

No newly established critical issue currently lacks a mitigation that requires stopping all authorised independent work. The audit itself remains unfinished. Task8contract closure and task9's required preparation/native observations are dependencies that block their production work. Region public support depends on coordination and output/exit integration.

Independent-person approval, representative live/sovereign/permission evidence, intended estate visibility, fresh supported OS and customer operational rollback cannot be manufactured by further offline testing. They remain readiness/acceptance blockers where the governing requirement demands them. No Azure scope/fixtures/RBAC or release authorisation is requested or inferred now. The current programme can continue offline while preserving those boundaries.

## Quality and development-speed priorities under review

1. Automate evidence reconciliation and compact current-state consistency: exact head/tree/parents/both identities, all changed bytes, native/source outcomes, full control sets and source hashes. Emit a reviewable evidence index tied to the candidate; preserve complete raw logs and failed trials. Avoid status-only certification commits that trigger another unchanged QA loop.
2. Strengthen mutation-result admission through a separately reviewed shared contract. Distinguish the intended compiling assertion failure from panic, compilation, setup, cancellation, timeout and infrastructure outcomes. Keep baseline, exact restoration and bounded complete rejected output. Static inventory is broader than the four reproduced drivers and needs its own validation before migration.
3. Address the demonstrated tag/comparator ambiguities with independent cases and explicit compatibility decisions. Bind every proposed regression to its complete governing contract/supplemental plan before assigning defect status. Add permutation/alias/unknown-state properties where they expose real behavior, rather than tests that merely restate implementation.
4. Keep focused local feedback early and full exact-head native/source gates at closure. Use isolated worktrees for mutators; never parallelize source-mutating controls with source-consuming suites in one checkout. Cache compilation/dependencies while preserving uncached test execution and all existing budgets/control identities.
5. Measure the integration critical path and rework, not only CI wall time. On the single observed185run, Windows controls dominate; ideal redistribution across its three existing shards can save only54seconds. More workers require a measured runner-cost/elapsed-time tradeoff and complete identity conservation. Evidence automation and less handover reconstruction are stronger initial speed candidates.
6. Separate current rule-content assurance, compatibility comparisons and future improvements. Retain provenance/freshness/precedence evidence and independent examples, use Denis's fork first, and evaluate other tools without importing weaker workflow defaults or expanding product scope by analogy.

[Research](REVIEW_RESEARCH_20261021.md) binds primary guidance and selected comparison-project revisions. These are proposals for final evaluation, not a claim that the improvements have been implemented or independently approved. Final questions will be consolidated after the remaining review.
