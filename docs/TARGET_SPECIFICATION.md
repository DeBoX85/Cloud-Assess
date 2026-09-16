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

Exact behavioral parity for internal plugins is deferred until after the core engine is equivalent.

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

Filters support subscription, resource group, resource type, individual-resource exclusion, recommendation exclusion, and tags.

- include lists narrow scope
- exclusion takes precedence
- all include tags must match
- any matching exclude tag excludes the resource
- tag keys are case-insensitive
- tag values are exact-match

The target configuration root will use neutral terminology such as `assessment:` rather than the legacy product name.

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

## Failure model

Assessment-fatal examples:

- authentication cannot initialize
- requested scope is invalid or cannot be resolved
- core resource inventory cannot be obtained
- core recommendation engine cannot execute sufficiently

Stage-local failures such as Advisor, Policy, Defender, Cost, or optional-plugin failures should normally preserve other valid results and mark the assessment partial.

## ARG behavior to preserve

- up to 300 subscriptions per request
- 5,000 records per page
- skip-token pagination
- bounded rule concurrency
- malformed-row tolerance
- unsupported logical-table skip semantics
- management-group-aware authorization where required

## Diagnostics behavior

Diagnostics uses the discovered inventory and ARM batch reads of `Microsoft.Insights/diagnosticSettings` to generate canonical findings for supported resources without diagnostic settings.

Errors are propagated rather than terminating from worker goroutines.

## Cost behavior

The source implementation defines the cost period as the previous completed UTC calendar month.

Metric: `ActualCost`

Aggregation: sum

Grouping: `ServiceName`

Cost-access failure must be distinguishable from a valid zero-cost response.

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

Cloud Assess should use explicit naming such as `--redact-subscription-ids`.

## Deferred functionality

Deferred until after core v1 equivalence:

- full internal-plugin parity
- MCP server
- historical report compare
- alternative VM SKU utility
- public documentation site
- package-manager/public distribution
- web UI or hosted API

## Core v1 acceptance

Core v1 is complete when Cloud Assess can reproduce the pinned reference implementation's material core assessment behavior for the same Azure inputs while also providing explicit stage health and assessment completeness.
