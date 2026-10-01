# Azure access model

Status: source and Microsoft documentation review, 2026-10-01. This is practical operator guidance for the current core-v1 build, not a live-certified minimum custom role or production approval.

## Identity and scope

Use a dedicated read-oriented identity with access to the explicitly approved subscriptions and resources. Subscription-level **Reader**, or equivalent inherited effective permissions, is a practical control-plane baseline. Reader grants control-plane reads, not resource data-plane access. Owner and Contributor are not required to run an assessment. Granting roles is an administrator task outside this toolkit; do not expand access automatically to cure a scan failure.

Management-group scans also need read access to the selected group's hierarchy and subscription membership. Management Group Reader addresses hierarchy visibility; it does not substitute for resource read access in descendant subscriptions. Access to the tenant root is not a prerequisite for targeting a readable child group. A partially accessible hierarchy must be validated against independently approved membership. Do not use the production-containing `Advisory` parent to work around deferred non-production validation DV-001.

Cloud Assess uses Azure Identity `DefaultAzureCredential` (pinned SDK v1.14.1). Its chain tries environment, workload identity, managed identity, Azure CLI, Azure Developer CLI and Azure PowerShell credentials. Select a specific source with `AZURE_TOKEN_CREDENTIALS`, for example `AzureCLICredential` for an operator's existing CLI session or `ManagedIdentityCredential` for an Azure-hosted identity. Authentication success does not prove resource, hierarchy or billing authorization. Do not put tokens or client secrets in reports or Git.

## Requests and stage permissions

The following describes implemented request semantics, not a claim that every adapter has a complete security audit. Query-style POSTs below retrieve data. The normal scan does not create Azure resources, enable Defender plans, assign policy, deploy agents, register providers, configure diagnostics or execute remediation. Authentication caches and local report files can be written.

| Stage | Implemented Azure operation | Read-access guidance and boundary |
|---|---|---|
| Scope | ARM subscription listing and management-group subscription/descendant GET pagers | Visibility of intended subscription IDs and, when selected, the target hierarchy. Listing errors fail discovery; successful listing does not prove hidden scopes were included. |
| Inventory and Graph | POST Resource Graph `resources`, API `2024-04-01`; pinned recommendation KQL | Read access to the queried resources. ARG can omit inaccessible objects without a partial-results signal. Includes the renamed AZQR recommendation catalog and supported pinned APRL/AOR queries. |
| Diagnostics | POST ARM `batch`, API `2020-06-01`, containing GET `Microsoft.Insights/diagnosticSettings`, API `2021-05-01-preview` | Diagnostic-settings read permission on eligible resources. No settings are created. Unsupported types and denied subrequests are different causes; retain warnings. |
| Advisor | GET Advisor metadata (`2020-01-01`) and query `AdvisorResources` | Read access to Advisor metadata/recommendations and affected resources. Resource scope is filtered downstream; indexed data and timing can differ between runs. |
| Defender plan status | Query `SecurityResources` pricing rows | Access to subscription-level security/pricing objects. RG-only resource access does not establish complete subscription plan visibility. |
| Defender recommendations | Query unhealthy assessments in `SecurityResources` | Microsoft documents Reader/Security Reader for viewing security health and recommendations. Security Reader alone does not establish general inventory access. No dismissal or remediation calls. |
| Policy | Query noncompliance in `PolicyResources` and definition/container joins | Read visibility of policy state, definitions and relevant resources. No policy evaluation trigger, assignment or remediation. Empty results do not certify compliance without coverage checks. |
| Arc SQL | Query Arc SQL server instances, machines and extensions through ARG | Control-plane read visibility. No SQL database connection or extension deployment. Numeric `vcores` response-shape validation remains open. |
| Cost | POST subscription `Microsoft.CostManagement/query`, API `2021-10-01` | Cost Management Reader, or equivalent effective cost-query access, at each resolved subscription plus billing-agreement prerequisites. EA charge-view policy and CSP enablement may matter. Queries previous completed UTC month ActualCost grouped by ServiceName. |

Graph, Diagnostics, Advisor and Defender plan status are enabled by default. Policy, Defender recommendations, Arc and Cost are opt-in. Plugin execution remains unavailable in this core-v1 build. `--stages` changes the defaults; it is not an exclusive allowlist. Graph remains mandatory.

**Cost scope:** resource-group, resource-ID and tag selection do not narrow the Cost query to those resources. Cost is aggregated across each resolved subscription. Do not label it RG-only or tag-only spend. Disable Cost when subscription-wide billing data is outside the approved evidence scope.

## Coverage is separate from execution health

Microsoft documents that ARG can return only accessible subscriptions/resources without indicating partial results when some requested subscriptions are accessible. Underlying subscription discovery intersects requested IDs with visible active subscriptions and filters. The coordinator now fails before resource queries if any explicit CLI/library subscription is absent from the resolved set. Filter-only selectors and all-visible/MG successful empty discovery retain existing behavior; see [SCOPE_REPORTING.md](SCOPE_REPORTING.md). A missing ID may reflect authorization, subscription state, tenant/session context or filtering; the scanner cannot infer the cause from absence alone.

`complete` and exit 0 describe execution against resolved/visible data. They do not prove every intended subscription, resource or security/policy row was visible. Independently record the approved subscription set and applicable permissions before accepting estate-wide conclusions. Canonical JSON exposes `scope` with requested inputs, effective subscription filter lists, resolved subscriptions and unresolved explicit IDs, alongside the existing hash/count. Resolved entries include subscriptions with zero resource rows. A failed listing is `not_completed`, distinct from successful empty resolution. Redacted output limits exact identity comparison; use private unredacted JSON for independent scope reconciliation. The reporting contract has offline coverage; representative restricted-identity live visibility evidence remains Gate 004/release work.

Successful empty data, exact source equivalence and read-only roles do not independently close that coverage boundary. ARG indexing is eventually consistent, so comparisons between different times also require caution.

## Cloud and access failures

Public Azure is the default. Government/China aliases and custom ARM/authority environment configuration exist in source, but this checkpoint does not certify those deployments. Review `internal/azure/cloud.go` before configuring custom endpoints/audiences; do not treat an arbitrary override as validated sovereign-cloud support.

Generic access denial fails the affected stage. Cost specifically retains a warning (`cost_subscription_skipped`) for Azure error codes `MissingRegistrationForResourceProvider`, `MissingSubscriptionRegistration`, `DisallowedOperation` and `NotFound`. This is not a rule that all HTTP 403 responses are successful empty results. Shared ARM options explicitly set `DisableRPRegistration`, and the scope adapter enforces that setting on copied options. Registration-required errors are returned rather than triggering an SDK registration POST. See [OPERATIONS.md](OPERATIONS.md) for response handling.

Advisor metadata continuations are restricted to the configured HTTPS host/port and repeated URLs fail. The default shared HTTP transport refuses redirects rather than moving requests outside their checked destination. Injected transports and SDK-owned pagers have separate policies. Tested SDK logs and response errors omit the bearer header, but raw service error content remains sensitive; this is not general log anonymization. See [OFFLINE_QA_REVIEW.md](OFFLINE_QA_REVIEW.md).

## Sources and validation boundary

Reviewed implementation: `internal/azure`, `internal/discovery`, `internal/arg`, `internal/diagnostics`, `internal/advisor`, `internal/defender`, `internal/policy`, `internal/arcsql`, `internal/cost`, stage configuration and CLI dispatch. SDK credential selection is grounded in pinned module source `azidentity@v1.14.1/default_azure_credential.go`.

Microsoft references, checked 2026-09-30:

- [Reader role](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles/general#reader)
- [Management and governance roles](https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles/management-and-governance#management-group-reader)
- [ARG permissions and consistency](https://learn.microsoft.com/en-us/azure/governance/resource-graph/overview)
- [Defender permissions](https://learn.microsoft.com/en-us/azure/defender-for-cloud/permissions)
- [Cost access and billing prerequisites](https://learn.microsoft.com/en-us/azure/cost-management-billing/costs/assign-access-acm-data)

No Azure calls, role assignments or production pilot were performed for this review. Minimum custom-role certification, restricted-identity live testing and full adapter/security review remain open.
