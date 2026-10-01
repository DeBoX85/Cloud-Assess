# Pagination and assessment lifecycle checkpoint

Date: 2026-10-01 (Europe/Oslo). Baseline `a23c122a381aed0acd2f4e7daba9e6731e3540b7`, merged PR #68. This implements a bounded part of AR-02, not a complete load/volume or release decision. No Azure request, reference mutation, rule/dependency change or comparator normalization occurs.

## Pagination behavior

ARG already rejects empty/repeated skip tokens per subscription batch. Advisor metadata already rejects repeated URLs and unsafe origins. The three SDK scope pagers previously kept following repeated continuations; synthetic one-link and two-link cycles reached a seven-request fixture tripwire with no cycle-specific error. They now maintain independent continuation sets per listing, reject repeated continuations and return an error with no partial listing. Scheme/host case is normalized; path/query values remain significant and opaque. Error messages omit supplied links/tokens. Existing origin/redirect boundaries remain enforced separately.

Tests cover all three SDK scope listing operations, non-empty accumulated rows discarded on cycle, valid finite multi-page listing, independent repeated listings, host-case normalization and distinct query tokens. SDK and Advisor loops check caller context between pages. A synthetic Advisor getter that ignores context but produces distinct URLs stops after caller cancellation rather than continuing to its runaway tripwire. These checks do not preempt a transport blocked inside an individual call, bound response size or place a fixed page ceiling on a legitimate large estate.

## Opt-in assessment deadline

The CLI accepts `--assessment-timeout 30m`; Go duration syntax also accepts values such as `90s` or `2h`. The application API uses `ScanOptions.AssessmentTimeout`.

- Positive values derive a context deadline around coordinator assessment execution, including scope, inventory, Graph/auxiliary queries, their paging, authentication-token requests and retries.
- Zero is the compatibility default and adds no deadline. Earlier caller deadlines/cancellation remain authoritative.
- Negative/malformed CLI values fail before authentication or report replacement. Negative application values fail before assessment execution.
- Local CLI credential/client construction precedes the application budget. Report rendering is outside it so status artifacts can still be persisted. The derived timer is canceled when assessment returns, with deferred cleanup for exceptional exits.
- The stage runner checks context before and after each enabled task. Canceled work cannot report complete success merely because a task ignored cancellation and returned nil error. Later tasks are skipped using existing failure semantics.
- Deadline/cancellation failures preserve already collected data and explicit stage error codes, render requested reports when possible and return exit 1. They do not become successful empty datasets or severity-only failures.

This is cooperative cancellation, not an OS hard kill or a strict whole-process/rendering deadline. An injected task/transport that blocks forever while ignoring context can still block. Unbounded distinct continuations with zero timeout remain a lifecycle/load limitation; operators should choose an approved assessment budget. A universal positive default, volume/page caps, resource/concurrency/backoff calibration and successful installed Azure/sovereign runs remain open. No arbitrary estate-size ceiling was introduced without evidence.

## Validation and evidence

A real coordinator/application fixture waits for the configured deadline in Advisor after healthy inventory/Graph results. It asserts failed completeness/exit 1, persisted JSON, retained earlier findings/resources and `assessment_deadline_exceeded`. Tests also cover zero budget, an earlier parent deadline, negative API/CLI values, flag mapping and context checks before/after tasks. A compiling mutation restoring the old stage invocation fails the later-task tripwire; reviewed code was restored. Native built CLI checks now contain seventeen preflight cases, adding malformed/negative timeouts and preserving report sentinels with zero observed authentication requests.

Local full race suite, vet, unchanged tidy/inventory freshness and built CLI checks passed under pinned Go 1.26.8. Final native quality and windows-validation must pass the exact published head before merge. Final run/head/tree/merge stays in the PR for indexing at the next ledger checkpoint. Gate 004 and release approval remain open.

PR #68 head `68150a9e405cdb09817bbbcc880a50ecca85a240` passed both required jobs in run `36866099111` and merged at the baseline above, tree `305d8b2b34b4debf225c4ccc19083507b4578be0`. Inspected native logs show fourteen package cases per host/no skips, 112 documentation destinations/eight flags/one snippet, 80.5% Linux coverage and zero reachable/imported-package findings with one residual module-only advisory. The advisory upgrade does not close OpenPGP or release-security review.

Primary contracts: [Go context deadline/cancellation](https://pkg.go.dev/context#WithTimeout), [signal cancellation](https://pkg.go.dev/os/signal#NotifyContext), [QA process](QA_PROCESS.md), [operations](OPERATIONS.md). The existing CLI signal context remains unchanged.

## Follow-up review: terminal cancellation and page policy

The follow-up starts from merged PR #69, commit `0a2dac1257cb99c390ad35eb927ac3e17171ba7d`, on 2026-10-01 (Europe/Oslo). Independent synthetic transports canceled the caller context while returning a successful, non-empty terminal response. Before correction ARG returned a result and Advisor returned metadata, both with nil error. Their adapters now check context after a successful request and before accepting its response. Both return cancellation and no partial result. The stage runner already protects assessment-level status; this closes the direct adapter API gap as well. Cancellation remains cooperative and cannot be made atomic with return or preempt an injected blocked call.

The same terminal-response fixture passes through all three actual SDK scope pagers without production changes; the pinned SDK pipeline already rejects that cancellation. Added ARG fixtures verify that fresh distinct tokens stop after three requests on cancellation and that legitimate identical tokens are independent across subscription batches and repeated query invocations. Existing healthy SDK paging, Advisor continuations, empty scope and request-option tests remain in force.

Reviewed production wiring shows inventory/Graph and Policy/Defender/Arc SQL use the shared ARG implementation. Advisor metadata is its separate pager. Normal Cost requests previous-month ServiceName aggregates once per subscription, matching pinned `internal/scanners/cost.go`; the region-selection plugin has a separate Cost pagination path and is an outstanding plugin-parity obligation, not the normal Cost stage. This review does not certify every plugin/response-volume path.

**Page-policy decision:** keep the existing zero assessment-budget compatibility default and do not add a page-count flag in this slice. There is no measured healthy estate/query distribution to justify a fixed ceiling, and page count does not bound a single large response or memory use. Cycle guards plus a positive operator-approved assessment timeout address distinct lifecycle risks without silently truncating healthy data. They do not establish a memory/row ceiling. Any future configurable page policy must define per-query/per-subscription-batch and SDK-listing semantics, fail explicitly rather than publish truncated success, propagate through all production stage clients and include healthy-boundary/load evidence. A default policy and volume/load calibration remain open; this decision does not close AR-02.

Two service-completeness questions are now explicit follow-ups:

1. Both pinned AZQR and Cloud Assess request `$top=5000` with ARG API `2024-04-01`. Microsoft's [matching REST contract](https://learn.microsoft.com/en-us/rest/api/azureresourcegraph/resourcegraph/resources/resources?view=rest-azureresourcegraph-resourcegraph-2024-04-01) specifies a maximum of 1000. The existing constant is the requested source-compatible value, not a verified service-enforced maximum. No large live result has established whether Azure rejects, clamps or honors 5000 in this environment. Do not infer that historical small-result equivalence validates this boundary. Request value remains unchanged pending characterization and a documented parity decision.
2. Both response models follow `$skipToken` and omit `resultTruncated`, `count` and `totalRecords`. Microsoft's [large-data guidance](https://learn.microsoft.com/en-us/azure/governance/resource-graph/concepts/work-with-data) describes truncation without a token for some query shapes. A terminal response therefore is not by itself proof of all matching records. Characterize the matching REST contract and supported catalog query shapes before adding fail/warn handling; preserve deliberate KQL limits and valid empty results. No historical report is reclassified solely from this static finding.

Local final race/vet/module/inventory and actual CLI/documentation checks passed before publication (121 documentation destinations, nine flags and one PowerShell snippet). Exact final native Linux/Windows jobs remain blocking, with final evidence in the PR and indexed at the next ledger checkpoint. FN-025 records reproduction, correction and execution recurrence. No Azure request or source/dependency/normalization change occurs.

## Resume point

AR-02 is partial: cycle safeguards and configurable cooperative assessment budgeting are implemented; default/load/volume calibration remains. Next autonomous check should characterize ARG request-size and tokenless-truncation semantics against the matching service contract and catalog query shapes. Page/volume configuration awaits defined semantics and load evidence. Branding (AR-03), source-feature characterization (AR-04), historical evidence gaps and deferred live/release validation remain on the roadmap.
