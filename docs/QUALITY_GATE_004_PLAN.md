# Quality Gate 004 Plan: Core Validation Decision

Status: **planned; not passed**

This is the decision framework for the next core validation checkpoint. Quality Gates 001-003 are immutable reviews of earlier commits. A green CI run or a semantically equal report alone does not satisfy this gate.

## Intended decision

Decide whether the implemented generic scan is sufficiently characterized for a **defined first-release core scope**. The gate records an exact target commit, pinned source/APRL commits, CI runs, evidence bundle references, every accepted limitation, and the reviewer decision. It must not declare the toolkit release-ready: package installation, supply-chain material, licensing and final operational/security sign-off need a separate release gate.

Choose and document the first-release plugin and CLI boundary before passing Gate 004. An item may be implemented and validated or explicitly deferred in the target specification, implementation plan, README and ledger. A deferred item cannot silently appear as supported.

## Evidence matrix

| Area | Required evidence before PASS | Present evidence and remaining boundary |
|---|---|---|
| Pinned provenance and executable path | Clean target and reference revisions, pinned APRL/rule snapshots, passing required CI checks, successful built CLI | Gate 003 and current CI cover these checks; record the final candidate SHA and its checks again |
| Default and optional assessment stages | Representative same-input semantic comparison per required stage, including non-empty row projection or faithful sanitized fixture; inspect stage health, not just comparator exit code | Default findings, Advisor, non-empty Cost and Defender plan status have live evidence; Policy and Defender Recommendations have empty-result execution evidence only; Arc SQL numeric `vcores` is unresolved |
| Scope and filter behavior | Expected resolved/contributing subscriptions; non-empty in/out-of-scope evidence for selected subscription, RG, scanner and tag filters; parent-group recursion or explicit accepted limit | Subscription, two-subscription, resource-group CLI, leaf management group, Storage/VM, RG include/exclude, tag-only include/exclude, observed recommendation exclusion and one exact-resource exclusion have live comparisons. The RG-exclude pass moved all three baseline group resources out of scope and selected the same nine IDs as the tag-exclude pass; the target data-subscription checks confirmed only the expected Dev subscription. The recommendation-exclude pass removed the one baseline finding from both reports while the VM remained selected. The resource-exclude pass moved that VM out of selected inventory and into out-of-scope while preserving all other selected IDs. The missing cross-run Advisor row has unknown recorded scope in a user-reported local target-inventory check; Azure row presence/timing remains unverified across separate scans. Parent-to-child management-group recursion and remaining filter interactions remain open |
| Diagnostics and warning classification | For each non-success API subrequest, determine affected assessment meaning or record a bounded, explicit accepted uncertainty | Historical unfiltered and exclude-tag ARM batch HTTP 400 warnings lack request correlation. The RG-exclude and resource-exclude passes also have one non-success Diagnostics warning each, with failed request and HTTP status unverified from the supplied summaries. The recommendation-exclude pass is `complete_with_warnings`, but its warning codes were not supplied. Individual Network Watcher GETs, including a repeated matching resource hash in the exclude-tag probe, support an explanation but do not establish exact batch mapping |
| Negative paths and exit codes | Tests with denied access, throttling/timeouts, malformed rows, critical and optional stage failures, cancellation, severity exit, and report persistence; no failed retrieval represented as an empty success | Some deterministic tests exist; the risk-prone production Azure adapters, throttling and CLI paths need targeted negative-path characterization |
| Output, identity and redaction | Test real XLSX/JSON/CSV/SARIF/stdout projections and redaction on supported platforms; document SARIF identity exposure and evidence handling; verify Windows report ACL behavior explicitly | Renderers and a cross-package smoke test exist; Windows tests revealed Unix mode checks cannot establish Windows ACL privacy; real-output and ACL review remain pending |

Every row must be marked **PASS**, **ACCEPTED LIMITATION** with scope and impact, or **BLOCKED**. A blocked required row prevents Gate 004 PASS. A comparator result of `equivalent = true` with `complete_with_warnings` cannot by itself satisfy the relevant row. Keep raw unredacted Azure evidence outside Git; publish only safe metadata, counts, classification and references to locally retained bundles.

## Continuous checks supporting this decision

The Go quality workflow checks pinned provenance, formatting, module graph, branding, PowerShell validation helpers, build/help/version, race-enabled tests, aggregate statement coverage, vet and reachable known vulnerabilities. The Windows validation job adds native Windows PowerShell 5.1 parser/helper checks plus Go tests and executable help/version smoke tests. An active ruleset now requires pull requests and both job contexts on `bootstrap/core-v1`; the enforced development CI is not a Gate 004 or release PASS.

The 75% total coverage floor is a regression signal, not a release criterion by itself. Prioritize behavior-based tests in lower-covered risk areas such as throttling, discovery adapters, CLI composition, canonical result assembly and orchestration rather than raising one aggregate number without corresponding assertions.

## Later release gate

After Gate 004 and the first-release feature decision, separately verify distributable artifacts on each supported platform, clean installation, dependency and license inventory, checksums and provenance, authentication/least privilege, failure recovery, and final security/operations sign-off on the exact release SHA. Record that result as a new dated gate; do not retroactively amend Gates 001-003.
