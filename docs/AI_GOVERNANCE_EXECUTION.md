# AI governance discovery and public execution plan

Status: bounded Graph discovery IN PROGRESS as an unaccepted library proposal after accepted workspace PR100. The pure/request libraries are accepted through PR96/97; ai-gov remains unavailable publicly. Start with [SESSION_HANDOVER.md](SESSION_HANDOVER.md); source/cell contracts are in [AI_GOVERNANCE.md](AI_GOVERNANCE.md) and [AI_GOVERNANCE_REQUESTS.md](AI_GOVERNANCE_REQUESTS.md). No laptop/Azure test is needed for the planned deterministic engineering.

## Exact next implementation slice

The candidate implements bounded source-specific Graph account discovery in internal/plugins/aigov/discovery.go; independently accept it after final native QA, preserving truthful partial data and recorded filter decisions. Keep ai-gov unavailable in CLI/registry during that library slice. Then integrate actual standalone, mixed and scanner execution plus all-format reports as a separate reviewable slice. Do not combine region selection, dependency updates, release work or new network services.

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

## Concrete discovery limits and response acceptance

This contract was independently specified before code in PR99. It is implemented by the current unaccepted discovery candidate. These limits belong only to AI discovery; existing normal ARG queries and accepted metrics/deployment budgets remain unchanged. They are target safety limits, not Azure capacity guarantees or measured load acceptance.

| Work | Limit and accounting |
| --- | --- |
| Owned subscription scope | 4,096 normalized UUIDs, matching the accepted projection scope cap; reject invalid/case-duplicate keys before requests. Subscription display names use the existing 512-byte safe-label rule. |
| Request batch/page | At most 300 subscriptions; sorted normalized UUIDs; fixed query above; `resultFormat: objectArray`; `$top: 1000` on every page. No tenant/MG wildcard, `$skip`, facets or widened scope. |
| Logical work | 64 discovery POSTs across the run, at most 32 pages for one subscription batch, sequential requests. Check before sending; retries remain inside the shared client's maximum six attempts per operation. |
| Received rows | 4,096 raw rows across all pages/batches, counted before filtering and including rejected/excluded rows. Selected unique accounts also cannot exceed accepted MaxAccounts. This prevents excluded/malformed rows from bypassing the work budget. |
| Response bytes | 2 MiB per attempt through PostBounded; 16 MiB across successful bodies, charged before JSON decoding, including a malformed successful body. Every body is closed. |
| Retained text | 16 MiB across owned scope labels, decoded projected strings and retained continuation tokens. IDs at most 2,048 bytes; projected labels at most 512 bytes; region at most 64 ASCII bytes. Reuse accepted UTF-8/control/unsafe-path validation; preserve display labels and normalize the regional DNS key to lowercase. |
| Continuation | At most 4,096 UTF-8 bytes, nonempty when present, no controls; preserve exact opaque value in JSON body. Track tokens per subscription batch; reject repetition/cycles. An HTTPS-looking token remains body data, never a request destination. |
| JSON work | Existing strict validator: depth16, 64 keys per object, 1,048,576 tokens, bounded arrays; then enforce discovery's 1,000-row page and 4,096-row run limits. Reject duplicate/case-aliased known fields, invalid UTF-8, trailing JSON and error envelopes. |
| Time | Five-minute discovery context, honoring earlier caller deadline/cancellation; shared default 30-second attempt timeout and bounded retry/operation policy. Injected clients/filter callbacks remain trusted cooperative code. |

Discovery and later Scan have separate budgets. These limits do not establish a whole-public-AI five-minute or shared-request cap; public integration must define its enclosing context and preserve earlier assessment deadlines. Exceeding any limit returns explicit failed discovery health, never healthy truncation, and preserves the already accepted prefix. Tests must cover the allowed boundary and one excess, including excluded rows and oversized late bodies.

The official [2024-04-01 schema](https://github.com/Azure/azure-rest-api-specs/blob/0799fda68aab1e2c6f5f0d06f028f451b14bf0a8/specification/resourcegraph/resource-manager/Microsoft.ResourceGraph/ResourceGraph/stable/2024-04-01/resourcegraph.json), blob9de6392a2cc1f1cf8a2423d3c48ed63b1fd2e4a1, was read independently. It requires `totalRecords`, `count`, `resultTruncated` and `data`; `resultTruncated` is the string enum `"true"`/`"false"`, despite prose guidance describing a Boolean. Require nonnegative int64 counts, count equal to raw data length, count no greater than totalRecords, and an object-array data value. Missing/null/string/fractional/negative count metadata fails visibly. Do not cap or allocate using totalRecords.

The schema's actual first/next-page examples contain `resultTruncated: "false"` **with a continuation token**. Always follow a valid token regardless of that flag, preserving exactly the same query and subscription batch. A present null/empty/aliased token is invalid. A tokenless `"true"` response is incomplete. Track cumulative raw rows for that subscription batch; tokenless `"false"` termination also requires cumulative count equal to reported totalRecords. If totalRecords changes between pages, record failed coverage rather than assuming a stable estate snapshot. Keep validated data but stop dependent completeness claims. This is a conservative target correction: Azure Graph pagination is not certified as a snapshot, and synthetic tests cannot certify intended estate visibility.

Validate the whole envelope/continuation and raw page budget before admitting any current-page accounts. An invalid envelope admits no current-page data; earlier accepted pages remain. Within an accepted envelope, validate each row before filtering: exact CognitiveServices account path; requested-batch subscription ownership; consistent subscription/RG/name/type; safe region; and empty/OpenAI/AIServices kind as the query defines case-insensitively. Missing/null kind and SKU strings retain the source's empty/default behavior. Wrong typed/unsafe/nonmatching rows are skipped with counted incomplete warnings; duplicate normalized account identities are skipped visibly, never overwrite the first account. Contradictory/foreign/descendant identities never reach the filter or metrics. Validate and budget sku_tier even though display uses only sku_name. Unknown safe fields do not supply tags or scope decisions.

Return an independently owned account slice plus explicit discovery health and error. A successfully queried, correctly filtered empty selection is completed. Invalid rows make completed-with-warnings discovery; later transport/metadata/cancellation/limit failure makes failed discovery with the valid prefix and sanitized cause. Health records count retained selected accounts; warning messages contain counts, not raw provider payloads/tokens. Cancellation must remain detectable with errors.Is. Public integration must carry this health into the final table even if enrichment of retained accounts succeeds.

## Independent discovery acceptance cases

The current discovery tests exercise these synthetic obligations; final exact-head native/protected acceptance remains pending. This is not public/report or live equivalence evidence.

| Case | Observable oracle |
| --- | --- |
| Literal source mapping | Synthetic rows distinguish sku_name from sku_tier and mixed casing; all Account/Region fields equal reviewed literals. Empty/null kind and SKU preserve accepted defaults; valid source filter cases are retained. |
| Request confinement | Explicit in-memory production transport observes only the fixed read-oriented ARM POST, exact API/query/options/batch and ARM audience. All six partial cloud override combinations, custom/sovereign configurations and invalid scope fail before credential/transport calls. Empty scope performs zero calls. |
| Batch/continuation | A 301-subscription scope produces reviewed 300/1 membership; token replay keeps the same first batch/query/top. A URL-looking token never changes destination. Literal false-with-token example continues; cycle/null/empty/tokenless-true and contradictory counts fail before subsequent calls. |
| Envelope and row trust | Invalid encoding, duplicate/case aliases, trailing JSON, error envelopes, wrong array/count types and excess budgets admit no invalid page. Foreign batch/account, descendant, contradictory labels and duplicate identity never reach metrics; valid sibling rows have explicit incomplete health. |
| Filter behavior | Actual AssessmentFilter exercises subscription/RG/resource/type selection, recorded include/exclude tag decisions and unknown include/exclude-only behavior. No invented tags or SetResourceScope widening; cancellation after a callback remains visible. |
| Prefix retention | Literal two-page and two-batch fixtures preserve first-page account fields after denial, malformed metadata, changed totals, canceled request or body/row/page/request exhaustion. No later account is fabricated; failed status/error persists. |
| Ownership/concurrency | Mutating returned fields cannot change a subsequent result or caller scope/filter snapshot; two concurrent runs with disjoint literal accounts/tokens share no state, including seen IDs, budgets and warnings. |
| Compiling controls | Remove batch-ownership validation, stop on false despite a token, and discard prefix on a later failure in isolated copies. Each must compile and fail its named behavioral assertion; clean/restored cases pass. Keep four earlier request controls. |

Use existing captured source outputs for source defaults/tag semantics and the immutable official examples for wire paging; construct new explicitly synthetic discovery fixtures with reviewed literal expectations. Do not generate expected results from the target decoder. Focused discovery race tests, independently literal mapping/request/filter/prefix cases, scope/page/request/raw-row/body/text boundary/excess checks, three compiling discovery controls/restored baselines and5000x decoder fuzz passed in the restored Linux workspace. Preserve the four earlier request controls. Final broad/stamped/native jobs, package/docs/provenance and protected acceptance remain required. No new capture or public execution/live PASS is claimed.

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

Historical FN056 Windows/tool disconnect is retained in the failure notes. Its local development dependency is now resolved by a fresh isolated Ubuntu24.04 workspace with checksum-verified Go1.26.8/PowerShell7.6.6, Git HTTPS, clean accepted target/pinned AZQR and actual APRL, verified modules and full local Linux QA. [DEVELOPMENT_WORKSPACE.md](DEVELOPMENT_WORKSPACE.md) supplies the actual recipe/paths/proof boundaries. PR100 independent native/protected acceptance and clean fetched source/target proof completed at ca7fc301de421f24f62c6b976fe2d79c0bfe6780; its separate merge-push native run37103646003 also passed. Current discovery acceptance must use its own final head and evidence. Never reset unrelated work, request credentials or weaken a gate. Executor readiness is separate from deferred operator laptop/Azure tests and does not certify Windows-native/live/Gate004/release evidence.

Rollback: reviewed protected revert of each isolated slice after checking current dependencies and intended first parent. Preserve existing pure/request libraries, source captures and four earlier public plugins. The remaining order is discovery, public AI execution, then separately characterized region selection, B6 and B7.

## Current library implementation and limits

NewDiscovery uses the same public-cloud guard as accepted metrics construction, capturing only ARM scope; no constructor requests a token. Discover owns normalized sorted subscription membership before callbacks, uses per-call seen IDs/tokens/counters and returns owned LocatedAccount values plus StageExecution. Invalid scope/empty scope performs no request. Production transport enforces per-attempt bodies/closure/redirects/retries; injected client/filter code is trusted and must cooperate with deadlines and confinement. Callers must supply a stable, owned filter snapshot; this interface cannot sandbox or clone arbitrary callbacks. The library never changes recorded inventory decisions or returns provider payload/token text in health/errors.

Strict page metadata and continuation checks precede admitting any current-page account. Invalid sibling rows and normalized duplicates are skipped visibly, while later request/page/coverage/cancellation/budget failure retains accepted accounts with failed health. All received rows count before filtering; projected text includes rejected fields and unused sku_tier. Exact16MiB decoded-text/one-byte excess fixtures use independently sized literal strings, including rejected oversized tiers, scope labels and tokens. Successful-body exact2MiB/16MiB and one excess are separate from retry/error-body limits. Metrics/deployment budgets remain separate, and public integration must propagate discovery failure even after successful enrichment. No whole-public-run budget or Azure snapshot/estate visibility is certified.

Self-review found an initial unaccepted row decoder did not charge rejected projected string bytes. It was corrected before publication and now has the independent exact/excess text fixture; FN057 records prevention. Normal Graph behavior, source/capture bytes/pins/dependencies, CLI availability and reports are unchanged. Reviewed protected revert remains the rollback route after dependency inspection.

Discovery pre-acceptance review additionally rejects Unicode case-fold aliases in regional DNS input and fixed ARM service path/type names, preserving safe display labels and ASCII casing. Distinct-ID Kelvin-sign/long-s fixtures cannot be masked by duplicate detection; an eighth compiling control disables the DNS ASCII guard and must fail its named assertion (FN058). Earlier candidate full local/native results do not certify this correction; final revised exact-head QA remains required.
