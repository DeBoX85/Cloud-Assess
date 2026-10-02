# Service-health migration contract

Status: IN PROGRESS, library adapter first; CLI/registry availability remains zone-only until separate integration acceptance. Base: PR85 merge `2f88e9d5b8dca309c0f27ccb6ec4ed54787bcebb`. This is B4 and the second real adapter check for B3; do not count the same migration twice.

## Pinned source and meaning

Source `internal/scanners/plugins/servicehealth/servicehealth.go` at AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`: metadata service-health / 0.1.0-beta / Azure Quick Review Team / MIT. Copy the exact Resource Graph KQL, six headers, sorting (percentage, subscription, region, type), two-decimal percentage and decimal integer formatting. Only selected subscription scope and scanner resource-type exclusions apply; RG, resource, recommendation and tag filters cannot be inferred from aggregated rows without resource identities.

The source query describes 90 days and uses a 2,160-hour denominator/duration cap, but contains no date cutoff. It sums durations and caps rather than calculating overlap unions. ServiceIssue durations determine the percentage; the event-count join additionally includes PlannedMaintenance. Preserve this source behavior explicitly; a corrected temporal definition is a separate characterized design change, not an incidental parity fix or a guaranteed rolling availability measurement.

Healthy nonnil data, including an empty array, uses **Service Issues** and the source last-90-days description. Source nil result/data uses **Service Health Availability** and its alternate description; target retains that table identity with failed health because missing data is not proven empty success. Decode failures retain valid rows with explicit warnings. Malformed metadata/request/later-page failures must not become empty success or leak provider bodies.

## Intended slice and acceptance

Implement a real bounded library adapter, exact query, source metadata/projection and canonical owned table health. Keep registry/CLI unchanged until coordinator/report integration is verified. Preserve ordinary schema 1.0; explicit plugin tables remain schema 1.1. No pin/dependency/normalization changes.

Synthetic acceptance: actual unchanged source scanner capture with a temporary transport-only seam; complete metadata/header/row comparison; independently literal query hash, sorting/formatting, scanner-type filtering, nil versus empty, malformed/foreign-subscription rows, denial/throttling, pagination/later-page failure, cancellation, bounds, per-run/concurrent isolation, authenticated read-oriented POST/cloud/audience/body and closure controls. A compiling negative control must fail its named assertion; restored focused/full checks and final required Linux/Windows jobs before merge. Report/CLI integration and nonempty live parity remain a subsequent boundary.

A temporary source archive initially lacked embedded APRL submodule data; copying an archive of the verified pinned APRL into that temporary copy restored source capture compilation. The retained reference checkout is unchanged. Capture execution succeeded with synthetic local TLS/token only; never Azure. Source algorithm/metadata/query are unchanged; a helper in temporary internal/az replaces its transport only.

## Resume, recovery and rollback

Finish the adapter/tests and independent capture evidence, publish code plus this contract/handover, verify exact native jobs and protected merge. Then integrate service-health alongside zone through coordinator/registry/CLI and raw/masked reports. Keep user laptop/Azure validation deferred. Roll back an accepted slice through a protected revert PR, preserving source pins and dependent table contracts.

## Current implementation and evidence bounds

The library scanner owns a per-run ARG transport budget: 3,000 selected UUID subscriptions, 64 request pages, 1,000 rows/page, 65,536 total rows, 2 MiB HTTP attempt/page, 16 MiB projected JSON data, 8 KiB row and 512-byte region/type fields. The underlying shared client bounds every retry/authentication attempt before middleware reads it; the selected HTTPS origin is validated before construction. The whole query has a five-minute context budget; each operation timeout is 30 seconds/attempt and 300 seconds total with the shared defaults, intentionally bounded independently of the source's 120-second default. Caller context can further shorten this. Injected transports remain trusted; row budgets cannot constrain arbitrary code inside them. Global ordinary-scan workload/AR-02 remains open.

ARG preserves selected scope, sorted 300-subscription batches, object-array format, documented top 1,000, safe opaque skip tokens and explicit tokenless truncation/cycle rejection. No management-group authorizationScopeFilter is added to this query. Completed page rows are preserved on later failures; incomplete query health is failed, provider text omitted, cancellation identity retained. Source type filtering applies after strict validation. Missing/null/duplicate/case-mismatched fields, invalid UUID/outside selected subscription, nonfinite/out-of-range percentage, negative/fractional/overflow counts and unsafe text become counted malformed-row warnings; source silently defaults/skips many of those shapes. Stable source-key sorting preserves page order for otherwise tied rows; no aggregate arithmetic is reimplemented.

Source capture `internal/plugins/servicehealth/testdata/source-output.json` SHA-256 `b310fef3dd1eff03cf251e3e948fcebff275a25581966e48c609f0ee8b432549`; exact copied query SHA-256 `b09ab97088e0c670d0199e83cac467aca7d6cae3c3fbae856d9d885660c36aff`. Capture retains entire source metadata and output (source internal type enum 1 maps explicitly to canonical internal string), no target-generated oracle. Empty nonnil source data was actually executed and uses Service Issues. Nil-data branch is source-inspected, not a live response capture.

Restored full race/vet passed including bounded POST/GET controls and scanner cases. A compiled mutation permitting a 65th page failed the named retained-page count assertion; restored code passed. Clean documentation/native gates and publication remain pending. Public service-health integration remains unavailable; zone registry unchanged.

The bounded POST path also rejects invalid/trailing/non-object JSON envelopes and duplicate/case-aliased top-level data/continuation accounting before the shared ARG decoder. Unknown unique metadata is permitted with a 128-field limit. Global data-byte and subscription-count controls retain valid completed pages and reject scope before transport. The page budget bounds continuation-token metadata retained by the shared client; the 16 MiB figure specifically bounds projected JSON rows, not all process memory.
