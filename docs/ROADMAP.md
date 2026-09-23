# Cloud Assess Roadmap

Status: working execution roadmap, 2026-09-23

Active development branch: `bootstrap/core-v1`

Reference implementation: `DeBoX85/azqr` at `8e4f0577f3615e6c9014c031bcad079f235369cc`

This document orders the remaining work and defines the evidence needed to close it. [TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md) defines the product contract; [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) records the implementation sequence and definition of done; [DEVELOPMENT_LEDGER.md](DEVELOPMENT_LEDGER.md) records what actually happened. Historical quality gates remain snapshots, not statements that every later boundary has passed.

## Starting point

- The generic scan, its canonical reports, the deterministic semantic comparator, and the Windows live-equivalence runner are implemented and CI-gated.
- Default-stage subscription, optional-stage, resource-group, and two-subscription live passes returned `equivalent = true` for the compared data.
- Cost and Defender plan status have non-empty live row-level equivalence evidence. Policy and Defender Recommendations have only empty-result live execution evidence. Arc SQL has not been live-validated.
- The two-subscription pass resolved both requested subscriptions and exercised non-empty inventory, findings, and Advisor records in each. Its target Diagnostics stage reported two HTTP 400 ARM batch subrequests, so those diagnostic-setting conclusions remain uncertain. See Phase N in the ledger.
- Management-group traversal, representative filtered scopes, final feature boundaries, packaging, and release review remain open.

The next work item is the Diagnostics investigation. Management-group planning and deterministic fixture preparation can proceed alongside it. Core equivalence and a release decision require the open validation limits to be resolved or explicitly accepted with a narrowed scope.

## Execution order and completion gates

| Order | Workstream | Completion evidence | Dependency or input |
|---|---|---|---|
| 1 | Investigate Diagnostics HTTP 400 responses | Affected resource types and IDs identified privately, response cause classified, finding impact assessed, and a regression fixture or documented API limitation recorded | Locally retained Phase N reports and access to the same non-production scope |
| 2 | Validate management-group traversal | Pinned source and target runs on the same management group; resolved descendants, per-subscription contribution, stage health, and semantic datasets reviewed | Suitable non-production management group and read access to its child subscriptions |
| 3 | Close filter and scope combinations | Paired source/target filter fixtures and semantic evidence for include/exclude, tags, scanner selection, and resource/recommendation exclusions | Representative existing resources; provision fixtures only if separately approved |
| 4 | Close missing stage data | Non-empty Policy and Defender Recommendations comparisons; Arc SQL response shape resolved and characterized | Appropriate live resources or sanitized targeted fixtures |
| 5 | Strengthen failure and output checks | Tests/evidence for partial stages, permission failures, warnings, redaction, severity exit codes, and report consistency | Results from steps 1-4 and controlled negative-path fixtures |
| 6 | Set and implement the first-release feature boundary | Required plugin and CLI work implemented; any omissions expressly re-scoped in the specification and ledger | Decision on first-release scope after core equivalence |
| 7 | Build distributable artifacts and complete release review | Repeatable builds, license/dependency inventory, checksums/provenance, clean-machine tests, and security/operations sign-off | Feature boundary and selected target platforms |
| 8 | Record release candidate baseline | Pinned source/target provenance, completed evidence matrix, accepted limitations, release SHA, CI runs, and rollback instructions indexed in the ledger | All required gates pass |

The order indicates dependencies, not promised dates. A suitable test scope may let steps 2-4 proceed concurrently; a missing environment must remain an explicit evidence gap rather than a fabricated success.

## 1. Diagnose the two HTTP 400 responses

1. Preserve the unredacted Phase N evidence bundle outside Git. Use its recorded commits, Azure environment, and stage health to reproduce the same requests.
2. Map the failing ARM batch subresponses to their diagnostic-settings resource requests. Perform read-only individual requests as needed; do not assume an HTTP 400 means that a diagnostic setting is absent.
3. Distinguish invalid/unsupported resource paths, transient Azure behavior, authorization issues, and target request-construction defects. Compare the pinned reference's request and finding behavior.
4. If Cloud Assess has a defect, add a focused regression fixture, fix it, and rerun the relevant comparison. If Azure or the source behavior is limiting, retain the warning and document which conclusions cannot be trusted.

**Exit:** the two responses and their possible effect on findings are classified. A semantically equal report by itself does not prove those two diagnostic-setting assessments are correct.

## 2. Finish scope and filter equivalence

- Run a management-group pass with the pinned AZQR checkout and Cloud Assess on the same identity, cloud, stages, and stable non-production hierarchy. Confirm which child subscriptions were resolved, which produced data, and whether either scan returned warnings or partial results. Cover nested groups when the test hierarchy provides them.
- Exercise the no-explicit-scope path against a controlled test tenant if its accessible subscription set is known; verify active/disabled subscription handling without unexpectedly scanning unrelated subscriptions.
- Validate the remaining filter interactions in paired source/target schemas: subscription and resource-group include/exclude precedence, service/scanner selection, individual-resource and recommendation exclusions, include/exclude tags, and parent/child tag scope. Cover multiple resource groups in one subscription where data permits.
- Reuse existing reference integration fixtures and sanitized recorded API fixtures for cases the live environment cannot exercise. Provisioning Terraform/Azure fixtures is a separate, opt-in step because it changes Azure resources and may incur cost.

**Exit:** each exercised scope has exact commits, parameters, unredacted local JSON, a semantic diff, per-subscription data checks, stage-health review, and classified non-empty deltas. Record untested scope variants as open.

## 3. Finish stage equivalence and behavior coverage

- Find existing non-empty Azure Policy noncompliance and Defender Recommendations data in an approved test scope, or develop targeted sanitized fixtures that validate their row projections. The previous empty-result passes establish execution and coverage only.
- Exercise Arc-enabled SQL with a real or faithful recorded response containing numeric `vcores`; resolve the pinned source's response-decoder ambiguity, add a characterization test, and classify any intentional target correction.
- Retain the non-empty Cost and Defender plan-status baselines and rerun only when changed code or a new scope creates a concrete regression question. Exercise Cost access failure separately from a valid zero-cost result.
- Characterize invalid input, authentication/authorization failures, API throttling, partial noncritical stages, malformed rows, and warning-to-completeness behavior with focused fixtures. Verify that reports are retained before partial and severity-gate exits.
- Cover recommendation applicability (not applicable, compliant, noncompliant), finding deduplication, and diagnostic settings present/absent with existing source fixtures or representative live data. Verify each supported cloud/authority/ARM endpoint variant with an appropriate test environment, or record it as unvalidated if none is available.
- Check real-output projections for XLSX, JSON, CSV, SARIF, and stdout, including redaction boundaries and the identity-bearing nature of SARIF. Avoid putting unredacted Azure reports into public CI artifacts.

**Exit:** required stages have representative non-empty or targeted evidence, known API limitations are classified, and source differences are explicitly documented rather than silently normalized.

## 4. Decide and finish first-release features

The target specification places external/YAML plugin execution and internal-plugin parity after core equivalence. The implementation plan permits remaining CLI/plugin gaps only when the release target explicitly defers them. Before declaring a release candidate, record one choice for each item: implement for the first release, or defer with a documented user-visible limitation.

| Item | Implementation if included | Acceptance evidence |
|---|---|---|
| External/YAML plugins | Load and validate definitions; execute them through production orchestration with bounded scope, stage health, and canonical results | Executable plugin fixture, malformed-input/failure tests, report projection, and read-only live test where feasible |
| Internal plugins | Migrate selected pinned-source plugins only after the common plugin path works | Named plugin parity matrix and source/target evidence for each included plugin |
| Scanner-specific CLI commands | Reuse generic orchestration without changing established assessment meaning | Help/flag behavior and end-to-end command tests |
| `rules` CLI | Expose catalog inspection with consistent provenance and neutral naming | CLI tests against the pinned catalog and documented output |
| `plugins list/info` CLI | Describe available plugins and enablement state honestly | CLI tests for list, lookup, errors, and plugin metadata |

Keep the current early error for explicit plugin-stage requests until a working production plugin path exists. Do not present internal-plugin parity as already achieved.

## 5. Package, secure, and operate the release

- Decide supported operating systems, architectures, installation method, version scheme, and whether package-manager distribution belongs in the first release. Build versioned binaries and test a clean installation on each supported platform.
- Produce reproducible build metadata, checksums, and artifact provenance. GitHub Actions [artifact attestations](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/use-artifact-attestations) are one option to evaluate once the release workflow exists.
- Generate and review a dependency and license inventory, including bundled rule data and transitive dependencies; retain the required notices and source pins. Recheck reachable vulnerabilities, action SHAs, Go patch level, and package-specific negative paths rather than relying only on the current aggregate coverage floor.
- Review authentication and least-privilege setup, timeout/retry/throttle behavior, error handling, local report permissions, redaction and SARIF exposure, evidence retention, and recovery/rollback instructions. Verify default operations remain read-oriented.
- Update the user guide, CLI examples, deployment/upgrade notes, supported scope matrix, known limitations, and troubleshooting guidance. Keep quality-gate snapshots immutable and record new results in the ledger.

**Exit:** a release candidate is installable and supportable from its published artifacts, and a final security, licensing, and operational gate has recorded results on the exact release SHA.

## Evidence and decision rules

For every new live pass, retain the pinned reference SHA, target SHA, APRL SHA, scope/stages/filter hashes, identity/cloud context, unredacted source and target reports, equivalence report, stage-health details, and classification of every delta or warning. Keep raw Azure evidence outside Git. A comparator exit code of zero with `complete_with_warnings` requires a warning review before declaring coverage complete.

Azure Resource Graph can return results only for resources the querying identity can read, so explicitly check resolved and contributing subscriptions instead of assuming the requested scope was fully observed. See [Microsoft's Resource Graph permission guidance](https://learn.microsoft.com/en-us/azure/governance/resource-graph/overview).

Before provisioning a test fixture or exposing a sensitive evidence bundle outside the local test environment, select the scope, expected cost, handling method, and access explicitly. Repository documentation and deterministic tests can proceed without that decision.

## Inputs needed when each boundary is reached

1. For Diagnostics: the retained Phase N target/reference reports or a local, sanitized account of which two read-only requests failed, including resource type and ARM error code.
2. For management-group testing: a suitable non-production management-group ID, descendant access, and the same scanner identity used for both executions.
3. For missing stage data: existing safe Policy, Defender Recommendations, and Arc SQL examples; any Azure fixture creation requires separate authorization.
4. For release scope: a decision about internal plugins, scanner-specific commands, `rules`, and `plugins list/info` if these cannot all be included in the first release.
5. For packaging: target operating systems/architectures and distribution expectations.

Do not request credentials in chat. The Azure scans run under an already configured identity in the user's environment.

## Explicitly deferred product work

Historical-report comparison, alternative VM SKU selection, MCP integration, a public documentation site, a web UI/hosted API, and package-manager/public distribution remain outside the currently validated generic core path. Revisit them only after the first-release boundary is recorded; new scope requires a corresponding specification, plan, and ledger update.
