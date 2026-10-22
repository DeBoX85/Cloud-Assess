# Cloud-Assess security reviewer handoff, 22 October 2026

Role: security assessment in a separate chat. The original development chat remains the sole production/integration owner. This document is a review contract and copy-ready instruction, not evidence the assessment has been performed. Do not treat the prior general read ledger or green QA as your own security coverage.

## Frozen code under review

Repository: https://github.com/DeBoX85/Cloud-Assess

Product candidate: draft PR185, branch fix/go-1.26.10-security-20261021, exact commit566cc4d8b78d168e5216245a5a386f742a3a2f85, tree39c44bf8e1d477623b901105fa752da1e7be9669, sole parent1348503ba415e630082004b785d802982b4c8b50. That carry descends from preserved PR183/458d5c5296c321c8e5e961f18eb5422d43c2eb18. Review this fixed product tree, not the latest documentation tree or a moving branch. Verify live objects/ref state before work; if they differ, report the difference and retain the frozen scope rather than silently switching. The development owner alone authorizes a scope revision.

Compatibility baseline: Denis's fork https://github.com/DeBoX85/azqr, pin8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1. APRL60eaddda76541f6adbc1c5ffa686829807e55e29. Protected bootstrap/core-v1a6fee427e19feeef28691737ea9bc5026086efa0 remains untouched. Source bugs are observations, not requirements to reproduce unsafe behavior; deliberate corrections require separate contracts/QA. No pin/expectation refresh.

Supplemental process scope: exact PR190 commitb2c524665e06e3a9adabee5ae1cd75ca96b0344c/tree5b375530feb3da91ab360f634444f869ea9c0f0c, including its additive handoff guard. Label that review separately from product185. Any later process/test/source changes need a separately identified delta review; this handoff does not certify itself.

Read current guidance from the verified documentation handoff revision separately: AGENTS.md, CONTINUE.md, CURRENT_STATE.md, ASSESSMENT_AF_20261022.md, AUDIT_DISPOSITION_20261022.md, TASK_CARRY_FORWARD_20261022.json, REFERENCE_POLICY.md and REVIEW_FINDINGS_20261021.md. These guide scope but must not replace the frozen product code. Current review packaging/own QA must be checked in its live draft PR; documentation success does not certify the product and vice versa.

## Permissions and separation

- Read-only repository/ref/PR inspection. Do not commit/push/create/close issues or PRs, alter branches/shared documents, merge, release, update dependencies/actions/SDKs or change source pins/oracles/controls. Return findings to the development chat for integration.
- Isolated local checkout/copy and disposable synthetic reproduction code are allowed. Never modify an existing development worktree. Record every test's exact base, isolation, command, expected/actual behavior, actual exit and limits; keep healthy controls and original bytes. A tool/setup/compile failure is not a vulnerability reproduction.
- No Azure requests, live retail queries, real credentials, fixture deployments, RBAC, production scans or external attack traffic. Fake credentials and explicit confined transports only. Research primary documentation is allowed; uncertain claims at90%confidence or below require primary evidence.
- No subagents, external messaging, new servers/listeners or scheduled work. This is another AI review, not independent-human approval or certification.
- Private evidence destination is unavailable. [DC-CI-EVIDENCE-20261022/issue189](https://github.com/DeBoX85/Cloud-Assess/issues/189) is DEFERRED OPEN. Do not publish raw logs, sensitive payloads or exploit detail to the public repository. Retain minimal safe synthetic evidence locally, disclose its temporary retention and return a sanitized report. Never claim hashes provide durable storage. Stop/report if an important finding cannot be safely handled within these boundaries.

## Method and review areas

Use four checkpoints with a complete own file/entrypoint ledger:
1. Scope/threat model: assets, attacker-controlled inputs, entrypoints, trust boundaries, security requirements, exclusions and exact revision inventory.
2. Contextual source analysis: all first-party production/tests/tooling/workflows/package code; trace complete paths, not isolated suspicious lines. Classify generated/reference files and third-party scope honestly; do not claim all transitive dependency source is manually audited.
3. Adversarial verification: independent minimal healthy/attack controls, compiling synthetic reproductions, actual exits, bounded raw evidence and precise preconditions; reject or downgrade unsupported hypotheses. Do not edit production or frozen tests to demonstrate a finding.
4. Report: confirmed findings, rejected hypotheses, remaining concerns, evidence limits, recommended remediation order and delta-review obligations. No blanket secure/release-ready verdict.

Cover credentials/audience/endpoint/cloud selection; redirects/continuations/SSRF and credential forwarding; read-oriented operation/no implicit writes; selected-scope/resource identity and ambiguous filtering; JSON/YAML/file/schema admission and resource exhaustion; HTTP bodies/retries/timeouts/cancellation; concurrency/ownership/caches; report masking/formulas/paths/permissions/staging; errors/logging/secrets; dependency/source/artifact provenance and CI trust; mutation/comparator/evidence integrity. Inspect test seams and injection assumptions, not only happy paths. Existing source and binary scans, fuzz/race/vet are supporting evidence; confirm tool availability/version before claiming new executions. Do not install/upgrade arbitrary tooling as an incidental review step.

Known starting points: R7-F01conflicting normalized tag keys, R7-F03duplicate comparison fields, R7-F04panic-versus-assertion mutation admission. They are confirmed general findings; reproduce/security-classify against frozen scope before severity. R7-F02is explicitly NOT a defect after root-depth0contract reconciliation. R7-O01generic body bounds, O02positional ID callers, O03sort ties and O04redaction fallback remain unconfirmed. Three imported-not-called and one module-only advisory remain OPEN from prior native scans; current database/source/binary results need new dated execution if claimed. Do not infer historical exploit or leak from a theoretical path.

Primary starting guidance:
- https://cheatsheetseries.owasp.org/cheatsheets/Secure_Code_Review_Cheat_Sheet.html
- https://go.dev/doc/security/best-practices
- https://docs.github.com/en/actions/reference/security/secure-use
- https://csrc.nist.gov/pubs/sp/800/218/final

## Finding exchange contract

For each finding return: temporary SEC-R-nnnID; frozen commit/tree and exact paths/lines; threat/preconditions; observed consequence; independent reproduction and healthy control; actual command/exit/environment; severity rationale and uncertainty; dependencies/affected development tasks; safe mitigation; suggested correction and retest criteria; evidence location/hash and retention limit. Distinguish confirmed/rejected/observation/blocked. Do not renumber or silently close existing R7 IDs.

The development owner deduplicates against existing findings, records acceptance/block decisions in the carry-forward register and alone implements separate corrective candidates. A fixed-snapshot review does not certify subsequent source changes: return reviewed path/blob coverage so the owner can request relevant delta review. No assumed cross-chat synchronization; Denis must bring the report/link back to the development chat.

Stop and report an unexpected issue requiring Denis when there is no safe mitigation or explicit deferral. Preserve the last verified scope and uncertain operation, without retrying mutations or widening access. At the end of each run include issue189retention status and remind Denis that independent human security review remains recommended before acceptance/release.
