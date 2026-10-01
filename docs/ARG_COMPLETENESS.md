# ARG page-size and explicit truncation review

Date: 2026-10-01 (Europe/Oslo). Baseline: merged PR #70, `b2494ca72b198d9bab6e36a8ff15b191e41a74dc`. Scope: production shared ARG request/response adapter, normal pinned catalog, and critical inventory/Graph failure persistence. No Azure access, reference/rule/dependency pin change or comparator normalization occurs. This closes a bounded service-contract slice of AR-02; live/load and release acceptance remain open.

## Contract and deliberate differences from pinned AZQR

Microsoft's [ARG REST 2024-04-01 contract](https://learn.microsoft.com/en-us/rest/api/azureresourcegraph/resourcegraph/resources/resources?view=rest-azureresourcegraph-resourcegraph-2024-04-01) specifies `$top` from 1 through 1000 and a string `resultTruncated` enum (`true`/`false`). It defines `$skipToken` as the continuation and also returns `count`/`totalRecords`. The matching API version is used by both implementations. [Large-data guidance](https://learn.microsoft.com/en-us/azure/governance/resource-graph/concepts/work-with-data) describes result shapes without a continuation token. Its general wording is less precise than the versioned REST response contract; the implementation follows the string enum and preserves normal true-with-token pagination.

| Boundary | Pinned AZQR | Cloud Assess after this correction |
|---|---|---|
| Requested `$top` | 5000 | 1000, matching the documented maximum |
| Subscription batching | Up to 300 | Unchanged, up to 300 |
| Valid continuation | Append rows and follow token | Unchanged, with existing token/cancellation safeguards |
| Explicit `resultTruncated=true`, no token | Field ignored; terminal rows can be returned as success | Explicit query error, no partial query result |
| Explicit `true`, empty token | Token handling is source-owned | Existing empty-token failure retained |
| Explicit `false`, no token | Terminal success | Terminal success, including valid empty results |
| Missing/null truncation field | Not modeled | Optional field preserves legacy/injected transport behavior |
| Invalid enum or wrong JSON type | Field ignored | Explicit metadata/decode error; raw enum value omitted from client error |
| Deliberate KQL limit with nontruncated response | Query as supplied | Query as supplied; no rewriting or automatic expansion |

These are deliberate target request/completeness corrections. Page size is not an estate/total-row ceiling: the fixture retrieves 2501 distinct rows over three pages. This does not prove Azure snapshot stability or live equivalence after the correction. A later paired live pass must record exact commits and classify any non-empty delta; no historical result is reclassified based on these fixtures alone. Source `$top=5000` may have been rejected, clamped or honored in past environments; static documentation cannot establish which occurred.

## Query screening and limits

The actual `rules.LoadPinnedCatalog()` loader returns 372 APRL, 23 AOR and 167 CUSTOM imported definitions, with 560 definitions after precedence. Applying `rules.IsExecutableEmbedded(definition, nil)` and excluding empty Query strings yields 344 candidate queries: 159 APRL, 23 AOR and 162 CUSTOM. A normal scan selects a resource-type/filter subset, so 344 is not the recommendation count in every report and is not a new parity/coverage percentage.

A temporary export explicitly included each Query string because canonical RecommendationDefinition JSON deliberately omits Query. After stripping line comments, a case-insensitive textual screen for `\|\s*(limit|take|sample)\b` found no candidate pipeline operators in those 344 queries, nor in inventory, Advisor, Policy, Defender status/recommendations and Arc SQL query constants. This is a candidate screen, not a KQL parser, output-type inference, complete semantic review, or Azure execution. It does not establish that every output includes a scalar column or that every response can be paged. Specialized APRL files, unimplemented built-in plugins and future external queries are outside that normal catalog.

The new guard uses explicit service truncation metadata only. Count/totalRecords consistency, missing metadata, inaccessible-resource omission, changing data across pages, query ordering, single-response memory limits and load/concurrency/backoff calibration remain separate limits. A missing/null field is not proof of completeness. No inferred count mismatch triggers a new policy in this slice. Existing execution completeness describes observed work; independently reconcile intended estate membership/access.

## Independent validation

Before the fix, actual response decoding followed by Query accepted a non-empty tokenless `true` response, both as a first page and after earlier healthy rows. A documented-limit fixture rejected the original request of 5000. Corrected tests assert explicit failure with no partial result and collect all 2501 unique, ordered fixture rows at requested page size 1000. True-with-token pages continue until a false terminal page. Empty, absent/null legacy metadata, invalid enum/type, empty token and deliberate limited-query cases are tested separately.

Application fixtures use the real HTTP response decoder and ARG query client. Inventory truncation runs through `discovery.DiscoverResources`; Graph truncation runs through `arg.ExecuteRecommendations`. Both pass through the production coordinator/application/report renderer, assert exit 1, failed stage/completeness, skipped later Advisor, no incomplete findings and persisted JSON. Healthy inventory remains in the Graph-failure report. No authenticated Azure call is made and this is not installed live-CLI proof.

Final local race/vet/module/inventory, actual CLI and documentation checks passed before publication. A compiling negative control removing only the truncation guard made both application fixtures report exit 0 and run Advisor; the tests rejected it, and reviewed source was restored before focused race checks. Final exact-head quality and windows-validation remain blocking before merge; native logs and final run/head/tree/merge are retained in the PR, then indexed at the next ledger checkpoint. Gate 004 and release approval remain open. FN-026 records reproduction, correction and review mistakes. See [PAGINATION_LIFECYCLE.md](PAGINATION_LIFECYCLE.md), [TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md), [QA_PROCESS.md](QA_PROCESS.md) and [ROADMAP.md](ROADMAP.md).

## Resume

Continue offline with supported branding adjustment (AR-03) and source-feature characterization (AR-04). AR-02 remains partial for default/page/volume policy, omitted/inconsistent response metadata, load limits and representative large-result live evidence. Defer live work until a suitable approved environment is available; do not use the production-containing Advisory parent to bypass DV-001.
