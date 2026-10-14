# Section 1 review working checkpoint, 2026-10-14

**In progress; not accepted and not a completed section checkpoint.**

Baseline: a6fee427e19feeef28691737ea9bc5026086efa0. Scope covers nine package directories and95baseline Go files (38production,57tests). All received a source/test pass across the initial and current sessions. Fresh correction QA, unresolved architecture questions and final coverage reconciliation remain pending. Files fetched solely for specific cross-module contracts outside this scope are not credited as fully reviewed packages.

## Findings

### SR-F01: ARG recommendation executor can silently return success on terminal context, P2

Confirmed by native run38561749374 on unchanged production executor blob23422ce441240a46b60583b88ed30fb1e3b092e5 at regression-only head432f2034aa3eed973670ab5b3a31cd1775d94611.
Linux job115844262674/full64338characters failed only the race test step; Windows115844263005/full51809failed only Go packages step. Both compiled and failed the four canceled/deadline definitions/empty cases plus cancellation during filtering; error was nil. This is expected regression failure, not a formatter/build/infrastructure failure.

Public-impact qualification: stages.Runner checks ctx.Err before and after each enabled task and marks it failed. This finding does not establish false-complete public scans. It is an inaccurate library postcondition and cancellation boundary for direct/composed executor use.

Correction: inspect context before scheduling, check it after receiving a job, and inspect it after joining workers. Preserve firstErr before final context checking so internally canceling sibling workers does not replace the initiating query error with context.Canceled. Return nil findings on terminal context. Ordinary successful empty definitions still retain existing empty-slice semantics.

State: confirmed; correction prepared/published in the successor commit; fresh verification pending. Regression expectations were authored before production correction and remain unchanged.

### SR-F02: inventory pagination uses non-unique ordering, P2 reliability risk

ResourceInventoryQuery orders only by subscriptionId and resourceGroup. Multiple resources share both keys. The unchanged pinned source internal/scanners/resources.go uses the same query, independently read at8e4f0577f3615e6c9014c031bcad079f235369cc.

Microsoft's primary guidance, accessed2026-10-14, explicitly documents that ties can reorder across pages even in static environments and recommends adding a unique ID to ordering:
https://learn.microsoft.com/en-us/azure/governance/resource-graph/concepts/paging-results

Contract proposed before correction: append id asc as the final sort key; preserve projection, filters, batch scope and page size. This is an intentional reliability correction to source behavior, not a claim of byte-identical source query. It does not guarantee a snapshot in a changing estate.

A1001resource/two-page deterministic integration fixture is declared before the query change. It exercises actual target HTTP decoder, ARG pager and discovery with all resources sharing subscription/RG. Without a unique tie-breaker, its service model returns a duplicate and omits the final ID while retaining page counts. The fixture models documented allowed reordering; it does not execute KQL or claim live Azure reproduction.

Semantic-compatibility check: projectTargetResources uses normalized resource-ID keys and compareDataset sorts those keys. Complete inventory identity/data, rather than tied source row order, is the semantic contract. Original source captures remain immutable. State: confirmed static/design risk; synthetic reproduction pending; no production query correction yet.

## Architecture and intent observations

- CLI flags/configuration feed a prepared owned request; scanner expansion and inventory tag decisions are kept in runtime indexes.
- Filters use a single known-field YAML document; legacy/unknown keys, null assessment and multiple documents fail before credential factories. Include precedence and nearest recorded parent tag scope match source contracts.
- Explicit unresolved subscription scope fails before inventory; reports retain requested/resolved/unresolved evidence. Filter-only and all-visible empty scopes intentionally preserve source behavior and do not certify RBAC/estate coverage.
- SDK scope discovery constrains HTTPS origin before authentication, disables provider registration and rejects continuation cycles. Default HTTP redirects are disabled. Trusted injected transports/policies have explicitly documented obligations.
- ARG batches at300subscriptions and1000rows, checks tokenless truncation and repeated tokens, preserves malformed-row/unsupported-table source semantics and propagates transport failures.
- Coordinator owns per-scan state, stage configuration and mutable plugin inputs. Scope and inventory/Graph are critical; optional failures retain other valid datasets and expose partial health.
- Stage runner handles cancellation before/after enabled tasks. Application renders available status/results before returning failure/partial/severity exits; local rendering is outside assessment-timeout budget by contract.
- Plugin names require exact registered spelling, so AI-specific preflight checks are not bypassed through name normalization.
- Positional resource-ID behavior supports source-compatible parent structural scope; broad arbitrary-ID parsing remains a question to reconcile across collectors, not a demonstrated security admission bug.
- Bounded/unbounded HTTP methods and scope pager total budgets must be traced through remaining collector/plugin sections. Existing docs explicitly avoid a universal volume cap; no new safety guarantee is inferred.
- Cross-module references checked here include stages.Runner, plugin registry and inventory equivalence projection/comparison. Their remaining standalone module review stays assigned to their planned sections.

## Evidence and remaining work

Regression-only source run38561749344 completed/success. It is separate from the failing native reproduction. No protected code change or acceptance occurred.

Next: publish the small SR-F01 correction with retained repro evidence, run SR-F02's independently declared fixture against the old query, then correct only the query ordering. Inspect new native QA and source evidence on the final candidate, review the exact diff and checkpoint Section1 honestly before asking for progress re-evaluation.

If service interruption occurs, read the latest audit branch and PR142 before retrying a mutation. Regression-only commit432f2034 is retained for backtracking; revert actual correction commits using normal Git history and repeat QA rather than resetting protected history.

## Source-contract checks after initial investigation

Pinned AZQR models.go blob3724182c006dcb8adc42fd2d6fe877d32ebf1aea contains the same positional resource-ID helpers. Their parent-type/name extraction is deliberately source-compatible, not general ARM-ID validation. Preserve it here; remaining collectors must validate their own contracts.

Official context source for the actual project toolchain, golang/go taggo1.26.9 src/context/context.go blobd00ac67e38296573d20520dbd7aceb2fd1a52906, confirms terminal Err semantics. A generic pkg.go.dev page resolved togo1.27.2, so it was not treated as version-specific evidence. The tool could not open the versioned package page; the pinned official source was read instead. No toolchain upgrade is made.

Microsoft REST2024-04-01 documents top maximum1000, opaque skipToken and string-enum resultTruncated examples. Its prose paging page uses Boolean terminology, so wire decoding follows versioned REST contracts/examples rather than assuming a Boolean JSON field.

Import organization in cmd/cloud-assess/command.go is corrected separately: standard-library imports together, external imports together. No semantic behavior changes are attributed to that cleanup.

Existing source-compatible missing-ID stopping and tolerant null/empty row decoding remain explicit compatibility behavior. Section2must assess business-record validity and completeness implications; these are not silently relabeled as newly fixed defects here.

## Two-host reproduction and correction preparation

SR-F02 regression-only query evidence at e012a94d8077abf180917408156b12cfa69f47ec:
native38563129439 Linux115848958752/full63105characters and Windows115848959368/full50413both compiled and failed ONLY inventory_paging_order_test's duplicate0999assertion in Go test steps. Both reported internal/arg passing, including the cancellation correction and added remaining-definitions case. Later mandatory steps skipped after intentional test failure are not credited as passes.

SR-F01 fix commit5dd97b0f530d578e17767a4eea14510b72ee6f6a, parent432f2034aa3eed973670ab5b3a31cd1775d94611.
SR-F02 fix commit9d1872b1a1737da05db512addae8479e530aa238, parente012a94d8077abf180917408156b12cfa69f47ec.
Import-only cleanupa98f2170d39c9dd8a53335c2bd147f4f9fc6496b, parent9d1872b1a1737da05db512addae8479e530aa238.
Each correction is isolated for backtracking/revert. Original source captures and pins are unchanged; neither source oracles nor earlier independent expectations were adjusted to make tests pass.

Final candidate includes all these changes plus recovery notes. Full exact-head native/source QA and review remain pending at preparation. Use live PR142 for current immutable head and subsequent results, not earlier regression runs.
