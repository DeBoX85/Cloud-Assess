# Cloud-Assess: consolidated A–F assessment, 22 October 2026

This delivers the current A–F assessment and recommendations. It is an evidence-based project assessment, not independent security assurance, completed whole-repository audit, protected acceptance or release approval. The expanded audit remains open for applicable-document reconciliation and final disposition. Task8 remains provisional. This changed assessment/handoff candidate needs its own full native/source QA; PR188 success does not certify it.

## A. Goal

Deliver a read-oriented Azure assessment tool with adjustable branding and feature/function/capability parity against Denis's fork, DeBoX85/azqr, pinned at8e4f0577f3615e6c9014c031bcad079f235369cc. Keep the healthy baseline behavior reproducible and preserve source captures. Deliberate safety corrections and improvements are encouraged through separately specified contracts, independent expected results and their own QA. Other tools and Microsoft guidance inform decisions without replacing the compatibility baseline.

The user outcome is a usable assessment workflow with truthful scope, findings, optional-stage health, completeness and useful reports. An available SKU catalog must not imply deployability; failed retrieval must not become successful empty evidence. Azure writes, scope substitution, fixture/RBAC changes, protected merges and release remain outside current authorization.

## B. Actual progress

The implementation foundations are substantial. Authentication/scope, inventory/filtering, recommendation execution, core stages, human/machine reports, branding/package support and five public plugins exist. The region feature has many implemented collectors/calculations, but lacks the supported end-to-end public execution path. This integration boundary is the largest remaining implementation uncertainty.

Protected baseline a6fee427 is unchanged. Product continuation remains draft PR185 at566cc4d8, which carries the approved Go1.26.10 patch over preserved PR183 retail implementation. Later review branches are not product coding bases. Tasks1–5, security/task6 prerequisites and task7 are complete offline; task7's recovered historical executions and formal closure are distinct from new execution. No new live validation is claimed.

PR185's own native/source QA passed. PR188 at e210ca1f5275929a32ba5bb352a99db1d9d3fcd0 now also has its own fully inspected quality39695096155/source39695096162 success:10quality jobs/full logs/steps,622Linux/617Windows complete control lines equal183/185,87.6%coverage against unchanged75%floor,23captures/315chunks/909865bytes equal exact candidateGit. Both native scans reported0reachable/3imported-not-called/1module-only vulnerabilities; the residuals remain OPEN. These are revision/run observations, not a guarantee that future advisories or undiscovered bugs do not exist.

The seven-day review retains its fixed14–21October observation window:45initial captured PR bodies,211runs/719jobs,691allocated historical full logs and28unallocated jobs. All43successful quality runs conserve predecessor controls and93successful source runs reconstruct their captures. Cancelled/unallocated work remains uncertified. Failed steps inside cancellations were preserved; no outage root cause was inferred from cancellation alone. Later PR187/188 verification is separate from those historical counts.

Read coverage:795complete raw files including all326first-party Go files, plus29separate reversible semantic views. Scripts/workflows/capture templates and first-party rule assets were read. All135remaining product183paths are explicitly classified; some applicable-document reconciliation remains incomplete. Generated known-SKU/latency data was mechanically checked against the pin, not manually read scalar by scalar. External APRL/vendor code and schemas are not a whole-upstream audit. See REVIEW_BOUNDARIES_20261022.md and its inventory. This is the exact limitation on the request to read the codebase; no blanket whole-repository completion claim is justified yet.

## C. Remaining work

1. Complete applicable-document reconciliation and final audit disposition. Keep source reads, generated/reference checks and historical narratives separate; define the justified corrective sequence for the three confirmed findings.
2. Task8: finalize required-kind identity/evidence contract and50independent scenarios. Preserve all required Storage kinds, exact Cognitive kind, SQL coordinate provenance and unsupported VMSS scoring.
3. Task9: publish the strict consumer and inspect BOTH actual native missing APIs before production. Preserve typed original identity/coordinates rather than reconstructing them from name-only or display/redacted rows.
4. Tasks10–12: owned per-run coordinator, public region CLI and all required output/health/exit integration, including FullCost disposition. Reuse existing collectors through their contracts.
5. Tasks13–14: representative synthetic workloads/budgets and actual installed default/custom package plus rollback exercises. Keep live/fresh-OS/customer-operation limitations explicit.
6. Task15: readiness reconciliation, supported path/platform decision and independent/live/operational/security/license/Gate004 obligations. No automatic protected acceptance or release.
7. Wider parity: report compare, alternative-vm-sku and separately designed MCP remain open. Reconcile the source types utility behavior as the acceptance matrix requires. These are not closed by finishing the region sequence, bundled SKU data or tools/equivalence.

The core-v1 acceptance milestone and full AZQR-parity milestone remain different. Any first-release deferral must stay visible in the full-parity backlog. Future enhancements need a separate backlog so completion has a finite definition.

## D. Risks

| Risk | Established evidence | Mitigation and exit requirement |
| --- | --- | --- |
| R7-F01 scope ambiguity | Conflicting case-equivalent tag keys change inclusion nondeterministically; reproduced locally, inherited from source | Avoid ambiguous filters; define deterministic admission before a narrow correction, preserve healthy/source cases and verify both native platforms |
| R7-F03 false equivalence | Duplicate JSON names can overwrite contradictory or partial evidence; independent compiling reproduction | Admit only independently ambiguity-checked provenance-bound comparison inputs; specify strict bounded admission before correction |
| R7-F04 mutation-proof weakness | Four SQL/scope drivers accept real named panic output as assertion proof; broader18of39driver static inventory lacks explicit panic rejection | Preserve historical outcomes without claiming historical panics; define shared result classification, retain complete rejected output and all identities/restoration checks |
| Integration truthfulness | Catalogs and calculations exist without complete public region coordination | Tasks8–12 must propagate identity, required inputs, unknown/partial/failure, cancellation and ownership through actual outputs/exits |
| Evidence loss/disclosure | Recent complete logs remain local; private destination unavailable, public upload not approved | DC-CI-EVIDENCE-20261022/issue189 is deferred open, with mandatory revisit and closure proof; never claim durable retention from hashes or green QA |
| Semantic handover drift | Active documents contained obsolete security waits, task7pending and SQL pending states despite passing heading/link tests | Current routing corrected; maintain one authoritative state and propose semantic generation/checks separately |
| Dated platform content | Pinned SLA/retirement/pricing/SKU/latency content proves baseline behavior, not current cloud accuracy | Separate content freshness/provenance review from compatibility; no unrelated refresh |
| Operational/security acceptance | Residual advisories, independent review, live permissions/visibility, workload/fresh-OS/rollback/legal decisions open | Assign owners and evidence at readiness; offline tests cannot certify missing evidence classes |

R7-F02 is not a confirmed defect: the governing supplemental history plan counts root depth0, and the observed behavior complies. The original failed probe used an unsupported expectation; its evidence remains preserved. Generic HTTP-body volume, positional-ID callers, canonical sort ties and redaction fallback remain unconfirmed observations requiring demonstrated caller impact before severity or corrective claims.

## E. Blockers, deferrals and stop rule

No newly established unmitigated critical issue currently requires stopping all independent offline work. Audit disposition gates task8; task9's required native observations gate production; coordination gates public-region support. These are dependencies, not missing user permissions for ordinary authorized work.

Private evidence retention is expressly DEFERRED OPEN by Denis, not resolved. It must be reconsidered at every handover/end of run, when a destination becomes available, and before task15/readiness/protected acceptance/release. No GitHub repository creation or Azure activity is assumed. Loss of a local file must be recorded and investigated rather than silently reconstructed from its hash.

Live estate visibility/permissions, representative workloads where live evidence is required, independent-person review, fresh supported OS, installed operational rollback and security/license acceptance remain later acceptance blockers. Do not infer them from a source rollback or a green development run. Stop and report an unexpected issue that requires Denis and has no safe mitigation or explicit deferral. Preserve uncertain operations and inspect before retrying.

## F. Roadmap and quality-preserving pace

The immediate priority is finish the audit's open documentary/disposition work, then the required-kind contract and regional integration sequence. Before each slice record its exact base, governing contract, independent expectation, exit evidence and rollback route. Update the handover between tasks. Product work resumes from verified185 with current guidance explicitly carried, never a review/archive/recipe tree.

| Priority | Improvement | Expected gain | Required proof before adoption |
| --- | --- | --- | --- |
| First | Read-only resumable evidence reconciliation | Less manual log recovery, incomplete transfer and repeated lineage work | Full pagination/attempt/head/tree/control/source binding; missing/duplicate/reordered/tampered inputs and raw-byte round trips |
| First | Compact machine-readable current state and generated resume views | Prevent obsolete pending/STOP/toolchain facts from surviving heading checks | Deliberately inconsistent state must fail; historical contracts/evidence preserved |
| First | Shared mutation-outcome classification, then staged driver adoption | Make negative proof distinguish assertion from panic/compile/setup/timeout/cancellation failures | Independent real-output controls, exact restoration and passing baselines; conserve all mutation identities |
| Next | Separate narrow tag/comparator corrections | Deterministic scope and trustworthy parity inputs | Both orders, equal/conflicting aliases, malformed/unknown/partial cases and unchanged healthy behavior on both platforms |
| Next | Measure recovery/rework and several comparable CI runs | Optimize actual elapsed development rather than test-count optics | Preserve full gates, isolate mutators and report runner cost separately |
| Readiness | Independent security review, release evidence, private retention and restore | Establish assurance beyond self-review/CI | Exact candidate scope, owners, retained results and dispositions; no inferred release approval |

Do not cut required controls to accelerate development. Use focused local tests early, one coherent published final candidate, full exact-head native/source QA at closure and final PR-body outcomes rather than a status-only certificate commit loop. The one measured PR185 three-shard model offers only54seconds ideal rebalance improvement; it does not justify a large speedup promise. Removing repeated evidence/handover reconstruction is the stronger first target.

## Research and applicability

Primary Go security guidance supports source/binary vulnerability scanning, reviewed updates, fuzzing/race/vet. GitHub secure-use guidance supports least privilege, immutable action pins and manual review of logging exposure; masking is not guaranteed. Microsoft testing guidance supports explicit strategy, entry/exit criteria and early layered feedback. Those principles support the ranked recommendations, not a compliance claim or authorization to refresh pins or run Azure.

Sources rechecked22October2026:
- https://go.dev/doc/security/best-practices
- https://docs.github.com/en/actions/reference/security/secure-use
- https://learn.microsoft.com/en-us/azure/well-architected/operational-excellence/testing

REVIEW_RESEARCH_20261021.md retains the prior NIST/SLSA/SRE research and selected revision-bound GitHub comparisons. Denis's fork is authoritative for compatibility; Azure/azqr upstream research stays separately attributed. The fork's release build/checksum/SBOM patterns inform packaging research, Prowler's grouping informs isolated feedback, and ScoutSuite's reviewer guidance reinforces independent review. None is a complete audit of those projects or permission to copy weaker path skips/unpinned actions. Native support needs execution evidence, not merely a cross-compiled artifact matrix.

## Questions and decisions retained

These are later acceptance/ownership questions, not reasons to halt the current offline work:
1. Who will independently review product correctness, security boundaries, gate logic and retained controls?
2. Which OS/architecture combinations must have native support evidence for the intended milestone?
3. When storage is available, who owns private evidence access, retention, backup and restore verification? Issue189 must remain open meanwhile.
4. What representative estate sizes, workload mixes, duration/memory budgets define task13 acceptance?
5. Who will independently attest intended subscription/resource membership and permissions before live estate-coverage claims?
6. Who owns residual-security risk, licensing, fresh-OS/installed rollback and Gate004/release decisions?
7. Who reviews future time-sensitive rule content and intentional source differences?

End-of-run reminder: consider commissioning an independent security audit of the codebase before protected acceptance/release. Suggested scope includes authenticated destination/redirect/pagination boundaries, parsers/budgets, scope/identity filtering, report privacy, comparison evidence, dependencies and CI trust. This assessment and zero reachable scan findings do not replace that audit; no audit has been commissioned by this reminder.
