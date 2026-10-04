# Operator runbook

Status: current core-v1 candidate guidance, reviewed 2026-10-01. Gate 004 and release approval remain open. See [ACCESS_MODEL.md](ACCESS_MODEL.md) before using an identity or approving scope.

## Prepare a bounded assessment

1. Verify the candidate archive/checksum, source revision and extracted payload using [PACKAGE_INSTALL.md](PACKAGE_INSTALL.md). Retain the manifest and build information. Candidate extraction tests are not fresh-OS installation approval.
2. Check `cloud-assess --version` and `cloud-assess scan --help`. Record the version, exact executable/source SHA, intended subscription IDs, filters, effective stages, identity type and tenant. Never record secrets.
3. Choose an explicit approved subscription or management group. Omitting scope scans accessible active subscriptions, which can exceed the intended boundary. Do not combine management-group and subscription/RG CLI scope. RG scope requires exactly one subscription.
4. Independently verify intended subscription visibility, active state, hierarchy membership and effective read permissions with the scope owner. A read-only `az account list --all --output json` can help inspect an operator's session; it is not an authoritative membership/RBAC audit. Resource rows and scope counts alone do not establish full coverage, especially for empty subscriptions. Preserve this comparison privately.
5. Review filter content and enabled stages. `--stages` modifies defaults; enabling Cost adds subscription-wide billing queries even for RG/tag filters. A valid filter and an empty result do not prove the filter selected the intended scope.
6. Use a new output base in a private, access-reviewed directory. Subscription redaction is enabled by default for JSON/XLSX/CSV/stdout but is not general anonymization. SARIF retains resource identities. Names, tags, findings, warnings and logs may contain sensitive information. Review Windows directory/file ACLs and retention; controlled CI ACL fixtures do not prove an arbitrary directory private.

Example for an already authenticated Azure CLI identity, PowerShell. Replace every placeholder. This example does not log in, grant roles or run an assessment automatically:

```powershell
$env:AZURE_TOKEN_CREDENTIALS = 'AzureCLICredential'
$subscriptionId = '<approved-subscription-id>'
$outputBase = '<private-directory>\assessment-unique-run'

& '<verified-package-directory>\cloud-assess.exe' scan `
  --subscription-id $subscriptionId `
  --json --xlsx=false `
  --assessment-timeout 30m `
  --output-name $outputBase
$scanExit = $LASTEXITCODE
Write-Host "Cloud Assess exit: $scanExit"
```

Capture `$LASTEXITCODE` immediately. For a child management group replace the subscription flag with `--management-group-id '<approved-child-id>'`. Do not use `Advisory` as an implicit substitute for DV-001. Optional stages need explicit approval and access review. Use the exact checked executable's help for stage parameters and filter syntax; do not invent flags. The scan uses the identity selected by its credential chain, not necessarily the identity of an unrelated portal session. Clear stale credential/cloud overrides or deliberately select one credential source before troubleshooting unexpected identity behavior.


Filter files accept one YAML document with known neutral schema fields. Unknown or legacy keys, additional documents and a null `assessment` fail before authentication. An empty file and `assessment: {}` preserve default filter behavior. This is an intentional input correction to prevent silently discarded scope restrictions; see [AUDIT_FILTER_SCHEMA.md](AUDIT_FILTER_SCHEMA.md).

## Evaluate the result

| Exit | Meaning | Required action |
|---|---|---|
| 0 | Execution completed and any severity gate passed; warnings may remain | Check completeness, every stage, warning codes, resolved scope evidence and expected data. Do not equate it with full estate coverage or absence of risk. |
| 1 | Configuration, authentication, execution or rendering failure | Preserve stderr and any reports. A critical-stage failure may still persist partial evidence. Do not accept results as a completed assessment. |
| 2 | Configured finding severity threshold reached | Reports are rendered before this exit. Review findings; this is different from a failed retrieval. |
| 3 | A requested noncritical stage failed, leaving a partial assessment | Reports are rendered before this exit. Inspect failed stages and retain healthy data with explicit partial labeling. |

Inspect canonical JSON `scope`, `completeness` and `stages` (`name`, `status`, `records`, `warnings`) alongside the process exit. `scope.resolvedSubscriptions` includes zero-data subscriptions; compare private unredacted IDs with the approved scope. `not_completed` distinguishes failed discovery from an empty successfully resolved set. An unresolved explicit subscription produces failed scope/exit 1 before resource queries; inspect `unresolvedSubscriptionIds` without assuming the cause. Filter-only includes are selectors rather than required membership. See [SCOPE_REPORTING.md](SCOPE_REPORTING.md). Skipped opt-in stages are not assessed. `complete_with_warnings` requires classification of every warning relevant to the decision. Stage records are stage-specific, not a universal count of Azure resources. Output rendering is not one transaction across all formats; an earlier report can exist when a later renderer fails. Do not silently reuse a previous run's file as evidence for the failed run.

## Troubleshooting without widening scope

| Symptom | Check and response |
|---|---|
| Expired CLI session or token failure | Confirm tenant/selected credential source and reauthenticate deliberately if needed. `az login` changes the local authentication session; it is separate from a read-only scan. Do not amend source pins or broaden roles to fix authentication. |
| 403 or denied listing/query | Inspect operation/stage and actual Azure error code. Confirm effective access at the requested scope; ask the scope administrator for the intended read access. An empty visible set is not a successful authorization test. |
| Cost unavailable/skipped | Review subscription-level cost access, agreement/policy/provider availability and `cost_subscription_skipped`. Keep the warning and affected coverage explicit; no provider registration or billing policy change is performed by the scan. |
| Diagnostics non-success subrequest | Retain warning, HTTP status and original evidence. Historical Network Watcher probes corroborate unsupported types, but do not identify every old multi-request failed response. Follow the bounded probe procedure in [EQUIVALENCE.md](EQUIVALENCE.md); do not erase warnings based solely on another run. |
| Unsafe continuation or redirect | Advisor rejects foreign-origin, malformed and repeated continuations; the default shared HTTP transport treats redirects as retrieval errors. Preserve stage failure and confirm the configured endpoint. Do not disable the guard or forward credentials to an unexpected host. |
| Throttling/timeouts | Preserve request/stage failure and retry context. Use `--assessment-timeout` for a cooperative total assessment budget; zero default adds none. Report rendering/local initialization are outside it, and context-ignoring custom code cannot be forcibly stopped. See [PAGINATION_LIFECYCLE.md](PAGINATION_LIFECYCLE.md). Avoid automatic repeated full-estate scans; plan a smaller approved scope or scheduled rerun. |
| Fewer subscriptions/resources than expected | Reconcile approved scope, tenant/session, active state, filters and RBAC. ARG may return only visible data without a partial indicator. Stop estate-wide conclusions until visibility is independently verified. |
| ARG reports results truncated without a continuation token | Preserve the failed report and stderr. Query results are incomplete; narrow the approved scope or review query shape before rerunning. Cloud Assess does not automatically rewrite the query or use offset paging. See [ARG_COMPLETENESS.md](ARG_COMPLETENESS.md). |
| Differences across runs | Compare exact pins, effective stages, filters, intended/resolved scope and timestamps. ARG indexing and service data can change. The historical RG-versus-tag Advisor difference remains unresolved; do not infer a tag defect or filtering correctness from that difference alone. |
| Git dubious ownership in the reference checkout | Review actual ownership and the trusted repository path. Do not disable ownership protection globally or modify the pinned reference source to suppress the check. |

Ctrl+C requests cancellation through the production command context. A failed critical stage may stop before a usable assessment exists; report creation is not guaranteed on every interruption. Keep stderr and any partial reports, verify the process stopped, and use a new output base for a deliberate rerun. No Azure resource rollback is needed for the documented query paths; cancellation still has local evidence and request-completion limits. Do not delete old evidence while diagnosing.

## Preserve reviewable evidence

For a live source-versus-target pass use [EQUIVALENCE.md](EQUIVALENCE.md) and `scripts/live-equivalence.ps1`. Retain the entire private run directory: metadata, source/target JSON, comparator JSON, stdout/stderr logs and exact filter inputs. A comparator PASS establishes only its compared datasets and normalization rules, not universal coverage, permission sufficiency, production readiness or warning-free execution. Verify recorded commands and exit codes, effective default stages, scope-stage count, independent expected IDs/filter invariants and warnings before logging the milestone.

For ordinary scans retain executable/build manifest, command/configuration (without secrets), timestamps, exit code, JSON/stage health, stderr and independent approved-scope/access evidence. Hash the evidence before transferring it; keep unredacted raw Azure data outside Git and public CI artifacts. Log safe counts, classifications and private evidence references in the development ledger. Search [FAILURE_NOTES.md](FAILURE_NOTES.md) for recurrence before changes and record confirmed development mistakes separately from environmental access failures.

No Azure/laptop task is required to review this runbook. Restricted-identity live validation, successful installed Azure scans, fresh-OS operation, full security review and the approved production pilot remain release work.
