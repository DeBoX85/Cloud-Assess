# AI governance discovery and public execution plan

Status: implementation not started. This is the next B5 contract after request-library PR97 acceptance, not evidence that ai-gov is publicly executable. Start with [SESSION_HANDOVER.md](SESSION_HANDOVER.md); source/cell contracts are in [AI_GOVERNANCE.md](AI_GOVERNANCE.md) and [AI_GOVERNANCE_REQUESTS.md](AI_GOVERNANCE_REQUESTS.md). No laptop/Azure test is needed for the planned deterministic engineering.

## Exact next implementation slice

First implement and independently accept bounded source-specific Graph account discovery, preserving truthful partial data and recorded filter decisions. Keep ai-gov unavailable in CLI/registry during that library slice. Then integrate actual standalone, mixed and scanner execution plus all-format reports as a separate reviewable slice. Do not combine region selection, dependency updates, release work or new network services.

Pinned source: DeBoX85/azqr at 8e4f0577f3615e6c9014c031bcad079f235369cc, internal/scanners/plugins/aigov/aigov.go, discoverOpenAIResources/Scan/processBatch. Fresh-session read verified the exact query and source mapping. Retained capture_test.go.txt already exercises actual unchanged source Scan with synthetic Graph rows, metrics/deployments, empty and unknown-include-tag branches; its captures remain independent source evidence. No new source capture or local execution is claimed by this plan.

```kusto
resources
| where type =~ "Microsoft.CognitiveServices/accounts"
| where isempty(kind) or kind contains "openai" or kind contains "aiservices"
| project id, subscriptionId, resourceGroup, location, type, name, sku.name, sku.tier, kind
| order by subscriptionId, resourceGroup
```

Preserve projected wire names id/subscriptionId/resourceGroup/location/type/name/sku_name/sku_tier/kind. Target LocatedAccount needs source Account fields and region; SKU display uses sku_name, not sku_tier. Missing/empty kind or SKU uses the accepted source display defaults. Source silently skips JSON decode failures; target must report malformed or incomplete discovery visibly.

## Discovery request and data protections

Use exact read-oriented ARM POST /providers/Microsoft.ResourceGraph/resources?api-version=2024-04-01, fixed validated origin/audience, shared no-redirect bounded transport and no provider registration. Inspect actual shared ARG behavior before choosing reuse: internal/arg.Client currently batches 300 subscriptions and follows skip tokens, but does not impose a total page/row/body limit and returns nil results on query error. Reusing it without additional bounds would lose discovery-prefix data and could create unbounded work. Do not incidentally refactor every normal Graph query.

Owned normalized UUID scope must reject invalid or case-duplicate IDs, foreign account IDs, descendants and contradictory subscription/RG/name/type/region identity before metrics authentication. Preserve exact source account labels, validate bounded safe text/path segments and reject aliases/duplicates explicitly. Enforce page/total row/text/body/request/time budgets tied to accepted MaxAccounts/MaxTextBytes and transport limits; specify concrete values before implementation. Empty scope needs no Graph/metrics/deployment tokens or fabricated rows.

Keep Graph continuations as bounded opaque body tokens sent only to the fixed endpoint. Reject duplicate/cyclic/empty/ambiguous continuation metadata and truncation without a token. Query JSON must reject malformed UTF-8, duplicate/case aliases, error envelopes, invalid arrays and mismatched response coverage. Keep valid prior-page accounts on later denial/cancellation/limits with failed health. Shared retry limits and successful-byte accounting must remain explicit; trusted injected clients must honor context and byte/redirect contracts.

Microsoft [large-result guidance](https://learn.microsoft.com/en-us/azure/governance/resource-graph/concepts/work-with-data) corroborates 1,000-record paging and token/truncation obligations; the pinned source/API and existing target completeness protections govern exact implementation. Do not replace pinned API versions with a current documentation example.

## Filter and ownership contract

Source discovery calls IsServiceExcluded(resource ID). Target config.AssessmentFilter already applies structural scope and recorded inventory decisions. The query provides no tags. An unknown account remains excluded with include-tag filters, while exclude-only tag filters preserve unknown scope, matching the retained actual-source include-tag capture. Mixed scans use the cloned recorded inventory decisions. Standalone scans must preserve the source unknown-scope behavior; never invent tags or call SetResourceScope(id, true) merely because an account appeared in Graph.

Exercise subscription/RG/exact-resource exclusions, scanner/type selection, known included/excluded tag decisions, unknown include/exclude-only decisions, case-insensitive ARM correlation and filter cancellation. Filter callbacks remain trusted application code; snapshots/results must be independently owned across sequential and concurrent runs. Do not overwrite the normal inventory or permit a plugin to widen selected scope.

## Separate public integration

Affected entry points are actual inventoried files: internal/orchestration/operations.go and coordinator.go; internal/plugins/registry.go with an owned AI projection validator; cmd/cloud-assess/command.go and plugins.go plus existing standalone/mixed/report tests. Add only the AI operation, pending-table dispatch, alphabetical selection and standalone command after real discovery/request acceptance.

Production public-cloud validation must occur before the selected AI path can authenticate to an unintended service. Inspect command preflight, credential/factory timing and coordinator scope discovery, not just the metrics constructor. Preserve the all-or-none authority/ARM endpoint/ARM audience check and exact distinct Monitor/ARM scopes; sovereign/custom support remains explicitly unavailable until separately validated.

Canonical AI tables retain all sixteen exact source columns/cells, metadata, AI Throttling empty-discovery versus AI Gov nonempty branch, real subscription correlation, limits and source formatting. Validate injected output identity/headers/sheet/labels/numeric cells, sanitize errors and own mutable rows/health. Pending tables survive critical failures; optional failure cannot discard independent assessment data. Unselected ordinary output stays schema1.0; selected table results use the accepted additive schema.

Require authenticated synthetic command-to-coordinator-to-JSON/CSV/XLSX checks for raw and default-masked output, source complete cells, empty/default branches, partial enrichment/discovery, denied/malformed/later-page/cancellation cases, stdout/report failure and concurrent ownership. SARIF must preserve its existing identity contract without inventing plugin recommendations. Retain formula/text/Unicode/sheet-collision/report replacement guarantees.

## Acceptance and recovery

Publish coherent code/tests/docs WIP on a new task branch, verifying tree/ordered parent/human identity and every changed file. Discovery guards need compiling named negative controls and restored baselines; broad native Linux/Windows jobs, full race/vet/build/package/docs/provenance/source hashes and reviewed exact-head logs remain mandatory. Final public integration additionally needs actual compiled command/help/preflight and all-format synthetic execution. No synthetic case closes live roles/models/metric totals/volume/sovereign evidence, Gate004 or release.

The current Windows executor has no callable Go toolchain, Git HTTPS helper or shell download access. FN056 distinguishes native CI evidence from unavailable fresh local development/source checkout proof. Restore an appropriate isolated development executor, inspect actual paths/processes/status, fetch live accepted core-v1 and unchanged source/APRL, verify clean tree/parents/pins/toolchain and re-read open PRs before beginning this implementation. Never reset unrelated work, request credentials in chat or weaken a gate. This is an execution dependency, separately from the operator's deferred laptop/Azure tests.

Rollback: reviewed protected revert of each isolated slice after checking current dependencies and intended first parent. Preserve existing pure/request libraries, source captures and four earlier public plugins. The remaining order is discovery, public AI execution, then separately characterized region selection, B6 and B7.
