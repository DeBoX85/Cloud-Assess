# Cloud Assess Target Specification

Status: working canonical specification

Reference implementation: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

## Purpose

Cloud Assess is a read-oriented Azure assessment engine. It discovers Azure resources in a selected scope, evaluates them against recommendation sources and Azure platform signals, and emits human-readable and machine-readable reports.

The primary flow is:

```text
Authenticate
  -> Resolve scope
  -> Discover resources
  -> Apply filters
  -> Load recommendation catalog
  -> Execute assessment stages
  -> Normalize findings
  -> Build deterministic summary
  -> Render reports
  -> Evaluate optional severity gate
```

## Compatibility position

Cloud Assess targets behavioral equivalence with the pinned reference implementation, not CLI, configuration, package, or branding compatibility.

Intentional changes include:

- new product and CLI identity
- neutral configuration schema
- `GraphRecommendation` concept becomes `RecommendationDefinition`
- `GraphResult` concept becomes `Finding`
- lower-level packages return errors instead of terminating the process
- every requested assessment stage records explicit execution health
- assessment completeness is separate from finding severity
- subscription-ID masking is described explicitly as redaction rather than general anonymization

## Branding adjustment requirement

Product identity is centralized in `internal/branding`. The [immutable branding profile builder](BRANDING_PROFILES.md) now supports validated build-time presentation identity without source edits. The [branding review](BRANDING_REVIEW.md) defines the acceptance checks; [paired synthetic report evidence](BRANDING_REPORT_QA.md) now implements the custom report/data checks, and [profile-aware candidate packages](BRANDED_PACKAGES.md) implement distribution checks pending final native acceptance. This requirement remains open until coherent native CLI/report/package propagation passes; a development builder alone is not completion.

## Core v1 scope

Core v1 includes:

- Go CLI application
- Azure SDK `DefaultAzureCredential`
- Azure Public, Government, China, and custom ARM/authority endpoints
- tenant-accessible subscription discovery
- management-group recursion
- explicit subscription and resource-group scope
- resource inventory through Azure Resource Graph
- include/exclude filtering
- scanner registry and resource-type pruning
- APRL recommendations
- Azure Orphan Resources recommendations
- custom recommendations
- YAML/KQL recommendation extensions
- diagnostics settings assessment
- Azure Advisor
- Defender for Cloud status
- Defender recommendations
- Azure Policy noncompliance
- Arc-enabled SQL reporting
- Cost Management reporting
- XLSX, JSON, CSV, SARIF, and JSON stdout output
- subscription-ID redaction
- post-report severity gate
- plugin architecture

Exact behavioral parity for internal plugins is deferred until after the core engine is equivalent. The current core-v1 executable therefore rejects explicit plugin-stage execution before Azure authentication rather than advertising a stage that cannot yet run.

## Output strategy

Both Excel and JSON are first-class outputs.

Excel is the default human-facing report.

JSON is the canonical machine-readable representation used for automation, integration, and equivalence testing.

Neither format replaces the other. Both are projections of the same canonical assessment result.

## Scope rules

Supported scope entry points:

- all accessible active subscriptions
- one or more management groups, recursively resolved
- one or more subscriptions
- one or more resource groups within exactly one subscription

Invalid combinations:

- management group + subscription
- management group + resource group
- resource group without subscription
- resource groups across multiple explicit subscriptions

Disabled and deleted subscriptions are excluded.

## Resource inventory

Inventory is collected before assessment execution and captures at least:

- resource ID
- subscription ID
- resource group
- location
- resource type
- resource name
- tags
- SKU name
- SKU tier
- SKU family
- SKU capacity
- kind

Discovery produces both in-scope and out-of-scope collections.

## Filtering semantics

Filters support subscription, resource group, scanner/service selection, individual-resource exclusion, recommendation exclusion, and tags.

- explicit subscription and resource-group include lists narrow scope
- an explicitly included subscription or resource group takes precedence over a matching exclusion, preserving reference behavior
- the legacy-compatible `include.resourceTypes` field contains scanner/service keys such as `aks`, `ca`, `st`, and `vnet`; it does not contain literal ARM resource-type strings
- selected scanner keys are expanded into an internal allowed ARM resource-type set
- that structural scanner/subscription/resource-group/resource-ID scope is applied both during inventory discovery and again to downstream findings
- exact individual-resource exclusions remove the matching resource
- all include tags must match
- any matching exclude tag excludes the resource
- tag keys are case-insensitive
- tag values are exact-match
- child findings inherit the nearest recorded parent tag-scope decision
- unknown tag scope is excluded for include-tag filters and retained for exclude-only tag filters

The target configuration root uses neutral terminology such as `assessment:` rather than the legacy product name. The target also renames the legacy `exclude.services` field to the clearer `exclude.resources`; this is an intentional configuration compatibility break, not a behavioral change.

## Recommendation catalog

Recommendation sources are normalized into one `RecommendationDefinition` model.

Initial source identifiers:

- `APRL`
- `AOR`
- `CUSTOM`
- `DIAGNOSTICS`
- `PLUGIN:<name>`

APRL remains independently pinned. Existing custom rules are retained but re-homed under neutral paths and source labels.

## Finding model

A `Finding` represents one recommendation affecting one Azure resource. It is not tied to Azure Resource Graph as an implementation mechanism.

This allows KQL rules, ARM/API diagnostics checks, and future assessment mechanisms to share one domain model.

## Core pipeline

The logical pipeline is:

```text
Initialization
Scope discovery
Resource discovery
Recommendation catalog construction
Scanner pruning
Graph recommendation execution
Diagnostics
Advisor
Defender status
Defender recommendations
Policy
Arc SQL
Cost
Plugins
Findings summary
Report rendering
Severity gate
```

Stage-level execution remains sequential unless safe concurrency is demonstrated. Work within a stage may be concurrent and bounded.

## Stage health

Each stage records one of:

- `completed`
- `completed_with_warnings`
- `skipped`
- `failed`

A failed data retrieval must never be represented as a successful empty dataset.

Assessment completeness is separately represented as:

- `complete`
- `complete_with_warnings`
- `partial`
- `failed`

Execution completeness applies to resolved/visible data; it does not independently certify intended estate visibility. Canonical JSON now records requested scope, effective subscription filters, resolved subscriptions (including zero-data subscriptions) and unresolved explicit subscription IDs. A successfully discovered scope missing any explicit CLI/library subscription fails before resource queries, persists failed status evidence when possible and returns exit 1. This is a deliberate target correction to silent subset discovery. All-visible/MG and filter-only successful empty scopes retain source behavior; ARG can still omit inaccessible resources. See [SCOPE_REPORTING.md](SCOPE_REPORTING.md), [ACCESS_MODEL.md](ACCESS_MODEL.md) and [OPERATIONS.md](OPERATIONS.md). Independent membership/RBAC and restricted-identity live evidence remain Gate 004 requirements.

## Assessment lifecycle limits

An opt-in `--assessment-timeout` / application AssessmentTimeout bounds coordinator execution through context propagation. Zero preserves the compatibility default; rendering and local CLI initialization are outside the budget. SDK scope continuation cycles now fail explicitly, and enabled stages check cancellation before/after execution. No process hard kill, universal positive default or volume cap is claimed. See [PAGINATION_LIFECYCLE.md](PAGINATION_LIFECYCLE.md).

## HTTP operation limits

The production authenticated HTTP client applies a positive `OperationTimeout` to one HTTP call, covering authentication, retry waits and response-body consumption. Defaults are ten times the per-attempt timeout. Earlier caller deadlines and cancellation remain authoritative; non-positive custom values add no operation deadline. This is not a whole-scan time limit or a guarantee that custom transports ignoring context terminate. Callers of `PostStream` must close response bodies; operation context cleanup occurs on body read completion/error or close.

ARM SDK assessment clients explicitly disable automatic resource-provider registration. A registration-required read error must remain an error, without a registration POST; caller-supplied scope options are copied before enforcing this assessment boundary.

The default shared HTTP transport does not follow redirects. Advisor metadata continuation URLs must remain on the configured HTTPS host/port, contain no user information or fragment, and not repeat. Unsafe pagination fails the Advisor stage rather than sending the ARM credential to another origin or reporting a truncated successful metadata result. These are target security corrections for unsafe response paths; scope SDK pagers now validate configured HTTPS origin before authentication and refuse redirects in their default transport. Injected transports/policies remain trusted code with separate internal redirect/destination obligations; same-origin pagination budgets remain open. Shared HTTP adapters use configured ARM audience for token scope separately from endpoint. See [ALIGNMENT_REVIEW.md](ALIGNMENT_REVIEW.md).

## Failure model

Assessment-fatal examples:

- authentication cannot initialize
- requested scope is invalid or cannot be resolved
- core resource inventory cannot be obtained
- core recommendation engine cannot execute sufficiently

Stage-local failures such as Advisor, Policy, Defender, Cost, or optional-plugin failures should normally preserve other valid results and mark the assessment partial.

## ARG behavior to preserve

- up to 300 subscriptions per request
- request at most 1,000 records per page, matching the ARG 2024-04-01 contract; this deliberately corrects the pinned source request of 5,000
- reject explicit service truncation without a continuation token rather than publish partial query success; see [ARG_COMPLETENESS.md](ARG_COMPLETENESS.md)
- skip-token pagination
- bounded rule concurrency
- malformed-row tolerance
- malformed-row visibility through warnings/completeness metadata
- unsupported logical-table skip semantics
- management-group-aware authorization where required

## Diagnostics behavior

Diagnostics uses the discovered inventory and ARM batch reads of `Microsoft.Insights/diagnosticSettings` to generate canonical findings for supported resources without diagnostic settings.

Errors are propagated rather than terminating from worker goroutines.

When an ARM batch response contains a different number of subresponses than requested, Cloud Assess treats the Diagnostics stage as failed rather than inferring that the unmatched resources lack diagnostic settings. This safety correction applies to malformed API output and is tracked separately from successful-response equivalence with the pinned source.

## Cost behavior

The source implementation defines the cost period as the previous completed UTC calendar month.

Metric: `ActualCost`

Aggregation: sum

Grouping: `ServiceName`

Cost-access failure must be distinguishable from a valid zero-cost response.

Current Cost queries aggregate each resolved subscription. Resource-group, resource-ID and tag filters do not narrow the billing query to selected resources; reports and operations must preserve that distinction.

## Summary semantics

Core findings are deduplicated by normalized recommendation ID + resource ID.

Recommendation applicability is represented explicitly as:

- `not_applicable`
- `compliant`
- `noncompliant`

SLA-category records are not counted as ordinary findings.

## Reports

Initial report formats:

- Excel, default human-readable report
- JSON, canonical machine-readable result
- CSV
- SARIF 2.1.0
- JSON stdout

The Excel report should preserve the conceptual datasets from the reference and add an assessment-status section or worksheet showing stage execution health and assessment completeness.

## Redaction

The existing masking behavior is specifically subscription-ID redaction. It does not anonymize the entire report.

Cloud Assess uses the explicit `--redact-subscription-ids` control.

When enabled, subscription-ID redaction applies to:

- Excel
- CSV
- canonical JSON files
- canonical JSON stdout

For JSON, known subscription IDs are also redacted when embedded inside serialized ARM/resource identifiers and related strings, not only in dedicated `subscriptionId` fields.

Report files are staged before replacement, preserving an existing report on rendering failure. Unix replacement files use mode 0600; New Windows files inherit destination-directory ACLs; replacements preserve and protect the previous effective DACL before writing report bytes, failing if it cannot be preserved. Multi-file exports are not transactional, and cross-platform atomic replacement or crash durability is not promised. CSV formula-like cells are prefixed as text; this includes negative numeric strings and changes the human CSV representation. Canonical JSON retains original values for machine consumers. Spreadsheet consumer/re-save behavior remains outside the neutralization guarantee. See [QA_PROCESS.md](QA_PROCESS.md).

SARIF intentionally retains stable Azure resource identities because resource IDs and fingerprints are part of its automation/baselining contract, matching the identity-bearing behavior of the pinned reference. SARIF must therefore be treated as sensitive/identity-bearing output even when subscription redaction is enabled for the other report formats.

Redaction is not full anonymization: resource names, resource groups, tags, subscription display names, and other potentially identifying metadata can remain.

## Deferred functionality

Deferred until after core v1 equivalence:

- full internal-plugin parity and remaining live extension validation
- the `plugins list/info` command surface
- MCP server
- historical report compare
- alternative VM SKU utility
- public documentation site
- package-manager/public distribution
- web UI or hosted API

## Core v1 acceptance

Core v1 is complete when Cloud Assess can reproduce the pinned reference implementation's material core assessment behavior for the same Azure inputs while also providing explicit stage health and assessment completeness.

## Offline recommendation inspection

The `rules` command now lists source-selected pinned embedded and Diagnostics recommendations without Azure access, with Markdown by default and local --json/-j output. It does not load external YAML plugins or accept scan scope/report flags. JSON matches the independently captured pinned source output; Markdown escaping and safe absent-guidance handling are documented target presentation corrections. See [RULES_INSPECTION.md](RULES_INSPECTION.md). Live execution, estate coverage and full feature/release acceptance remain separate.

Scanner-specific commands now follow the pinned source command keys and reuse generic orchestration; see [SCANNER_COMMANDS.md](SCANNER_COMMANDS.md). This is offline command/selection evidence, not live validation of every service.

Bounded [YAML Graph plugin execution](YAML_GRAPH_PLUGINS.md) now joins ordinary scan orchestration after input validation before authentication. This preserves healthy source schema, precedence and external-rule filtering while applying explicit bounded/strict input protections. Internal table-plugin execution is still unavailable; external queries do not depend on that stage. Offline parser/transport/report checks do not establish live plugin equivalence.
