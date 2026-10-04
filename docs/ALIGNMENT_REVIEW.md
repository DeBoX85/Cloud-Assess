# Project alignment and security sanity review

This register preserves historical alignment evidence and explicit deferrals. Current offline implementation/recovery status is in SESSION_HANDOVER and PROJECT_AUDIT_REPORT_20261004; historical open feature sequencing below is not a current availability declaration. Live alignment gaps remain open.

Date: 2026-10-01 (Europe/Oslo). Starting live revision: `0df39bd9c3a8a0e2062f54772ac9e29750b3409b` (`bootstrap/core-v1`, merged PR #66). Reference checkout verified at `8e4f0577f3615e6c9014c031bcad079f235369cc`; APRL remains `60eaddda76541f6adbc1c5ffa686829807e55e29`. No reference repository changes, dependency refresh, Azure request, role change or release publication occurred.

## Resume checkpoint

The eight unattended tasks and their subsequent audit were complete. Scope reporting was merged in [PR #66](https://github.com/DeBoX85/Cloud-Assess/pull/66): tested head `8422f78a1d247beb5baaaf1585bcfc9effea19e4`, required [run 36810973320](https://github.com/DeBoX85/Cloud-Assess/actions/runs/36810973320), merge above, tree `9d3b0be415224bce256cc3e7d0070cdce3a901a2`. Native logs previously inspected show all fourteen package cases on both hosts, no skips, 97 documentation destinations/eight flags/one parsed snippet, 80.3% Linux coverage and zero reachable/imported-package vulnerabilities, with three module-only advisories.

The next normal task was SDK-owned pagination and injected-transport credential boundaries, then pinned dependency advisory review. This sanity review addresses the SDK scope-pager destination defect below. Resume with dependency advisory maintenance after this review's exact-head QA; do not re-open completed tasks merely to generate another evidence commit. Injected transports remain trusted code with an explicit obligation to enforce their internal redirect/destination behavior.

## Alignment conclusion and feature matrix

The implemented common assessment path remains aligned with the intended functionally equivalent AZQR toolkit. It is not full feature parity or release-ready. Different architecture, neutral schemas and branding are permitted; missing functions remain obligations for eventual parity. Deferral from a first core release does not cancel them. Branding customization is an additional user requirement and is not complete simply because identity constants are centralized.

| Function / requirement | Pinned source and target evidence | Status / next boundary |
|---|---|---|
| Scope, inventory, scanner selection and filters | Source commands/scanners; target discovery/config/registry/coordinator; paired live passes and deterministic fixtures | Implemented within recorded coverage; DV-001, restricted-identity coverage and remaining combinations open |
| APRL, AOR and AZQR-authored recommendations | Reference `internal/graph/graph_scanner.go` selects APRL Azure resources, orphan resources and AZQR Azure resources in that order; target `internal/rules/embedded.go` selects the same normal corpora, with AZQR renamed CUSTOM | Retained; pinned snapshot and catalog fixtures protect provenance. AZQR was not removed as a recommendation library |
| Specialized APRL workloads | Source embeds specialized workloads but normal NewScanner scan types do not select them; target pin retains them outside its normal catalog | Normal-scan behavior agrees; presence on disk is not automatic execution. Do not silently expand the default catalog |
| Diagnostics and Advisor | Production adapters plus row/report fixtures and live non-empty comparisons | Implemented; historical batch correlation and cross-run Advisor timing remain unresolved |
| Defender plans, recommendations, Policy, Cost and Arc SQL | Production stages and explicit health; Cost and plan data live non-empty; fabricated optional-stage HTTP-to-report cases | Implemented with limits; non-empty live Policy/recommendations, Arc numeric vcores and deployment coverage open |
| XLSX, JSON, CSV, SARIF, stdout, redaction and exits | Canonical result, renderer/application/process tests and native package checks | Implemented; SARIF intentionally identity-bearing, redaction is not anonymization; installation/live release boundaries open |
| YAML/KQL extensions | Filesystem catalog loader exists; target CLI does not expose a production external-definition execution path | Partial infrastructure only; do not equate loader tests with executable feature support |
| Built-in plugins | Source registration names: ai-gov, carbon-emissions, region-selection, service-health, sql-eol, zone-mapping | Not migrated; each needs its own functional characterization and permissions/report contract |
| Plugin discovery and metadata | Source commands/plugins.go list/info and registry; target rejects explicit plugin stage before authentication | Not implemented; explicit failure is safer than false support but does not satisfy parity |
| Scanner-specific CLI and rules listing | Source commands/scanners.go and rules.go; target now includes root/scan/branding/rules | Generic scanner filtering exists; offline rules inspection implemented, specialized commands remain missing |
| Compare, alternative VM SKU, MCP | Source commands/compare.go, alternative_vm_sku.go, mcpserver.go | Source features deferred beyond core; retain eventual parity work separately from non-source web UI/site enhancements |
| Easily adjustable branding | internal/branding consumed by CLI/report layers | Validated immutable profile, native builder and paired report checks now available; profile-aware packages implemented; final native acceptance required |
| Packaging, security and operations | Candidate ZIP, inventory/notices, strict stamped builds and blocking native jobs | Development evidence only; signed publication, fresh OS, license/security decisions and controlled pilot remain open |

## Confirmed defects and narrowly scoped remediation

### SDK scope continuation could carry a bearer token to a foreign origin

A pre-fix real subscription SDK pager driven by an in-memory transport returned a foreign HTTPS nextLink. The regression failed with `calls=2 foreignBearer=1` and no error. Only a synthetic bearer value was used; this is not evidence of historical credential disclosure. Pinned azcore runtime accepts absolute next links and its ARM bearer policy is not an application origin allowlist.

Scope construction now copies options, preserves no-registration behavior, prepends a per-call HTTPS origin guard before SDK authentication, and uses a redirect-refusing default client over Go's standard transport (including environment proxy and standard dial/TLS timeouts). SDK-specific idle-pool/HTTP2 tuning is not preserved by that replacement. Host comparison is case-insensitive and port-specific; user information, fragments and backslashes are rejected. Invalid configured scope endpoints fail before client creation. Errors omit response-supplied URLs. Existing caller policy arrays are copied. Fifteen actual SDK pager cases cover all three discovery operations, unsafe destinations and accepted same-origin continuation. Public/Government/China/custom-port option fixtures do not certify live sovereign deployments.

Caller-supplied transports and policies are trusted application code. A transport can redirect internally after the SDK policy runs; the guard is not a sandbox around injected code, a DNS/IP allowlist or a whole-scan budget. SDK repeated/distinct same-origin continuations and total scan bounds remain follow-up items. Healthy scope projection and equivalence normalization are unchanged.

### Shared HTTP adapters ignored the selected ARM token audience

A pre-fix custom-cloud regression requested `https://management.example.test/.default` despite an explicitly configured different audience. The pinned SDK defines Audience as the access-token target and Endpoint as the service URL, and SDK scope clients already honor that distinction. The shared ResourceManagerScope now uses the configured audience, trimming its trailing slash before appending `/.default`, with its existing endpoint fallback only when no audience is supplied. Named-cloud and custom/fallback tests pass. This intentionally corrects authentication configuration; no claim of live Government/China/private-cloud validation is made. A custom SDK cloud still requires an audience, and incomplete overrides retain the documented source-compatible fallback behavior.

## Quality/security assessment and remaining work register

Review inspected current specifications/planning/gate records, recent commits, pinned command/catalog/plugin entry points and security-sensitive target paths: credential/cloud/options, SDK discovery, scope safety, coordinator/application exit/report persistence, result ownership, redaction, report staging, request destinations, comparator guards and workflow controls. Existing full tests supplement this targeted source inspection. This is not a line-by-line proof of every dependency, a penetration test, live estate audit or Gate 004 PASS.

The existing gates are suitable for their stated bounded historical/development purposes. Blocking native jobs, race/vet, provenance, actual executable/package checks, vulnerability scanning, independent negative cases, fuzzing and mutation controls provide useful regression protection. Statement coverage and green CI do not establish parity, security or release acceptance. Two reproduced findings show why adversarial middleware/path review must complement those controls. New regressions belong in existing required jobs; another numerically named gate is unnecessary for these fixes. Gate 004 and the separate release decision remain necessary and open.

| ID | Follow-up | Can proceed without user Azure/laptop? | Closure evidence |
|---|---|---|---|
| AR-01 | Investigation/remediation recorded in [DEPENDENCY_ADVISORY_REVIEW.md](DEPENDENCY_ADVISORY_REVIEW.md); x/crypto upgrade and explicit residual advisory decision | Yes | Inventory and local scans complete; exact-head native CI required before merge; release review remains open |
| AR-02 | SDK cycle/page cancellation, ARG/Advisor terminal cancellation and opt-in assessment timeout implemented; see [PAGINATION_LIFECYCLE.md](PAGINATION_LIFECYCLE.md) | Yes | Partial: documented page size/explicit tokenless truncation corrected; default/volume/load bounds and broader metadata/live completeness remain open |
| AR-03 | Immutable five-active-field build workflow accepted in PR #75 | Complete for this bounded development workflow | Final native run 36940456214; CLI/report/default and custom installation accepted; runtime/logo/theme/release boundaries remain separate |
| AR-04 | [Pinned-source characterization](FEATURE_PARITY_CHARACTERIZATION.md) complete for entry points and sequencing; production parity still open | Yes for implementation | Offline rules command implemented; next scanner subcommands; per-plugin endpoint/column fixtures, health, reports and eventual live evidence remain required |
| AR-05 | Reconcile historical Diagnostics warnings, recommendation-exclusion warning and Advisor cross-run uncertainty | Partly; historical raw evidence or suitable later live pass required | No inference from counts or later uncorrelated probes; preserve warning status |
| AR-06 | DV-001 nested non-production management-group equivalence | No suitable live scope currently available | Explicit deferred validation record; synthetic hierarchy is not live closure |
| AR-07 | Non-empty live optional stages, Arc response shape and restricted-identity visibility | Requires suitable Azure evidence | Same-input raw reports, expected membership and classified deltas |
| AR-08 | Hosted maintenance dispatch, fresh OS, signed release and controlled production pilot | Offline controls partly possible; environment/approval later | Actual hosted/published/installed artifact evidence and separate release decision |

## Validation and primary references

Focused Azure/discovery tests, final full race suite, vet, tidy with unchanged go.mod/go.sum, nine inventory safeguard cases/inventory freshness, actual built CLI checks and documentation validation (107 destinations, eight flags, one PowerShell snippet) pass locally under pinned Go 1.26.8. A compiling mutation removing the SDK origin guard failed the new unsafe-continuation regressions; reviewed source was restored before final race validation. Local `-buildvcs=false` accommodates the already recorded workspace context; native CI and packaging must continue strict clean stamped builds. Exact final PR head must pass quality and windows-validation before merge. Final workflow/head/tree/merge evidence is retained in that PR and indexed at the next ledger checkpoint, avoiding self-referential documentation-only CI loops.

- [Pinned SDK cloud contract](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore@v1.23.1/cloud)
- [OWASP destination validation and redirect prevention](https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html)
- [Go vulnerability checking](https://go.dev/doc/security/vuln/)
- [QA process](QA_PROCESS.md), [Gate 004](QUALITY_GATE_004_PLAN.md), [failure notes](FAILURE_NOTES.md), [deferred live validation](DEFERRED_VALIDATION.md)

Subsequent checkpoint: [dependency advisory review](DEPENDENCY_ADVISORY_REVIEW.md) records AR-01 and advances the autonomous resume point to AR-02. Earlier next-task statements above describe the starting checkpoint.

Subsequent AR-04 implementation: [offline rules inspection](RULES_INSPECTION.md) now matches all 380 pinned reference JSON rows and adds actual CLI/package checks. Its exact-head native acceptance is retained in the PR; remaining plugin migrations and live/ancillary parity are unchanged.

The [development execution plan](DEVELOPMENT_EXECUTION_PLAN.md) governs future batches. [Scanner commands](SCANNER_COMMANDS.md) now implement source command registration and scanner-key selection using existing orchestration; final native acceptance is required, and all plugin/live/ancillary boundaries remain open.

B2 [YAML Graph implementation](YAML_GRAPH_PLUGINS.md) preserves discovered query integration, metadata conversion, sorted override and external filtering, with explicit strict/bounded input and confinement differences from the pinned source. The actual pinned loader capture and synthetic ARG-to-report fixtures are offline evidence. Internal table plugins, ancillary commands and live/release boundaries remain open; source characterization tables above describe their starting state.

B3a [zone adapter](ZONE_MAPPING.md) preserves pinned source row projection and adds API-documented continuation handling and explicit partial-failure information. This is an unexposed adapter: canonical plugin tables, privacy/rendering, CLI execution and plugin discovery remain open. Source silently omitted continuations/subscription failures; those corrections are documented rather than normalized away. Synthetic/source-parser evidence does not certify live plugin behavior.

B3b [canonical table infrastructure](PLUGIN_TABLES.md) preserves ordinary schema 1.0 and adds schema 1.1 only for explicit plugin-bearing assessments. Source metadata/columns/rows are owned and retained; target health/version/identity fields and bounded format validation are documented additions. Synthetic zone/service-health-shaped fixtures verify report/privacy infrastructure, not real plugin command or second-adapter parity. Those remaining B3 boundaries and live/release limits stay open.


Current internal-plugin advance (2026-10-03): the earlier unmigrated-plugin table is a historical review snapshot. Zone/service-health/SQL EOL public execution is accepted offline through PR85/87/89; carbon pure/request libraries through PR90/93. Carbon public execution is the current unaccepted proposal with local full/race/build/report/source/security checks passed, final native/protected acceptance pending. AI governance, region selection, remaining ancillary commands and live parity remain open. See SESSION_HANDOVER and CARBON_EMISSIONS for exact current evidence, contract and deferrals.
