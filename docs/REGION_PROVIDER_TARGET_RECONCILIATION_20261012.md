# Provider-location target collector reconciliation, 2026-10-12

Status: pre-edit contract, unaccepted until this proposal's own QA. No target production
collector has been added. Read [REGION_PROVIDER_LOCATIONS.md](REGION_PROVIDER_LOCATIONS.md)
for the source contract and [PR135](https://github.com/DeBoX85/Cloud-Assess/pull/135)
for accepted characterization. [PR136](https://github.com/DeBoX85/Cloud-Assess/pull/136)
merged at bf4991af7162c8076a80da827a2c80c77af3eb25; its distinct accepted-push native
38294494765/source38294494675 proof is pending at preparation. Verify its live body
before dependent acceptance. Do not transfer candidate results.

## Request and owned interface

Use the existing RESTQuotaGetter/GetBoundedWithResponse transport contract and a
validated HTTPS ARM origin, without adding armresources imports or changing the
dependency graph. The first request is exactly GET
/subscriptions/{lowercase-selected-UUID}/providers?api-version=2021-04-01, without
body or expand. The [Microsoft request contract](https://learn.microsoft.com/en-us/rest/api/resources/providers/list?view=rest-resources-2021-04-01)
was freshly inspected: 200 supplies ProviderListResult, including resource types and
locations. A metadata location is not proof of deployability or a live location catalogue.

Proposed internal ProviderLocationRequest contains SubscriptionID. Collect accepts
context, stable selected subscription map and request; result contains normalized
SubscriptionID, Status (complete/partial/unknown), owned Locations keyed by
namespace/type, and sanitized FailureCode. Collection is subscription-wide and is
reused across target regions; it does not invent TargetRegion or SKU/zone evidence.
The later coordinator must bind the selected subscription and set
AvailabilityEvidence.ProvidersComplete only when Status is complete. In this slice,
independent calculation tests may explicitly construct AvailabilityEvidence to
prove the boundary. No public registry/scanner/CLI activation.

Invalid origin, request/scope and cancellation return error/no result. Operational
failure returns unknown before any page was admitted, partial after any wholly
admitted page, retaining only previous pages. Successful explicit empty value array
is complete. No borrowed maps, shared per-run cache, raw response error text or
provider-derived diagnostic text escapes.

## Accepted source observations and deliberate target decisions

| Source observation | Target decision |
| --- | --- |
| Lowercase namespace/type without trimming | Require bounded ASCII namespace and slash-separated type identifiers; lowercase after validation, reject surrounding whitespace/non-ASCII/empty segments |
| Remove ASCII spaces from location, then lowercase | Preserve that transform for bounded ASCII display location strings; accept explicit global and physical-format normalized names, reject logical exclusions; no claim of catalogue membership |
| Repeated locations form a set | Deduplicate valid locations only after charging every raw entry |
| Duplicate normalized types overwrite, including empty across pages | Reject all normalized type collisions within/across pages; never union or overwrite |
| Missing/null envelope/fields often become empty success | Require explicit non-null value/resourceTypes/locations arrays, nonempty namespace/type strings and non-null objects/string entries; empty arrays remain valid |
| Null provider/type panic | Reject the whole page without panic |
| Later failure drops all pages; caller substitutes empty | Retain only prior admitted pages with partial health; no fabricated complete empty collection |
| Caller produces output after cancellation | Cancellation/deadline returns nil result/error, including after last page |
| Relative continuation fails in source SDK | Safe relative references may be resolved against the validated first URL; require exact origin/path/version before credentials |
| Cycle stops only because synthetic fixture returns 400 | Reject canonical repeated continuation before another authenticated request |
| SDK behavior depends on options and default transport | Use existing bounded ARM client, disabled registration/redirect contract, fake-token canary for independent transport tests |

Optional returned provider ID, if present, must be a nonempty structurally valid
/subscriptions/{selected UUID}/providers/{namespace} identity with fixed ASCII
segments and namespace binding before case-insensitive equality. Missing ID is
allowed and is not fabricated. Null/wrong-type/present-empty ID is rejected.
Unknown metadata fields may be ignored only after whole JSON admission, including
duplicate/Unicode case-fold aliases and string/token/depth limits.

Each page is decoded and validated in owned temporary state. Validate shapes, raw
work, text, optional IDs and collisions against both current page and prior state
before committing any of its locations. A bad continuation discovered after a
valid page preserves that page as partial; cancellation still returns no result.
Namespace-only repeats with disjoint normalized types are allowed; duplicate
complete type keys are rejected even when sets are identical.

## Fixed bounds and continuation rules

64 successful pages, 1 MiB per response, 8 MiB successful aggregate bytes and 65,536
raw provider/type/location entries per call. Charge a provider object, a type object,
and each location entry before deduplication. Count empty providers/types as work.
Use existing MaxLabelBytes and decoded text guard for every projected identifier,
returned ID and raw location; charge text before normalization/deduplication.
Use the existing whole JSON validator for UTF-8, invalid surrogate replacement,
controls/noncharacters, Unicode folded duplicate object keys, maximum depth 32,
65,536 value tokens per page and 128 keys per object. These limits are independent:
the earliest exceeded guard wins. Add provider-specific helpers only where the
existing validator cannot express this contract, preserving existing callers.

Continuation permits only HTTPS exact configured origin/selected first path,
one exact api-version=2021-04-01 and at most one nonempty bounded $skiptoken.
No expand, foreign subscription, userinfo, fragment, opaque URI, encoded path,
unknown/aliased/duplicate query, redirects or cycles. Preserve decoded skiptoken
meaning with canonical re-encoding. Validate before calling the authenticated getter.
An unsupported continuation is incomplete evidence, never successful truncation.
Do not make registration, tenant discovery, Azure writes or credentials requests
during offline fixtures.

## Independent production acceptance matrix

Literal cases must cover one/two pages, complete empty, global and nested types,
display normalization, duplicate locations, namespace repeats and type collision
within/across pages, empty overwrite rejection, missing/null/wrong mandatory shapes,
null entries, optional ID structure/binding, full page rejection and preserved
earlier pages, all HTTP status/error/body-close outcomes, safe/unsafe relative and
absolute continuation, aliases/Unicode/text boundaries and canonical cycles.

For each page/aggregate bytes, page count, raw work, decoded text, label, JSON
token/depth/key and skiptoken bounds test exact acceptance and first excess with
independent literal inputs. Construct multi-page fixtures so a tighter page guard
does not mask an aggregate/work assertion. A successful 64th final page completes;
a continuation needing page 65 is partial without the extra request.

Use the actual bounded ARM client with credential canary and fake transport to
assert GET/body/audience/path/query/registration absence/redirect rejection and
response closure. Test cancellation before, between and during decode/admission,
deadline errors, sequential/concurrent run isolation and caller mutation of returned
maps. Inputs must remain stable during the call; no claim of racing caller-map safety.
Independent availability tests must check complete/partial/unknown health and full
comparison fields using literal expected records.

Add compiling named production mutations with restored baselines to both required
hosts. A failed compilation is not mutation detection. Keep all existing source
goldens/harnesses/pins/controls/security checks and Go1.26.9/x-net0.60 intact.
Require fresh complete exact-head Linux/Windows/source/review/publication/tree/
identity/preview/protection, protected expected-head merge and distinct actual
accepted-push proof. No local Go, live Azure, load, fresh-OS or release proof inferred.

## Exact next and deferrals

Publish and read back this reconciliation before target production edits. Finish
PR136 distinct accepted-push verification, then implement only this internal
collector and independent bounded tests/controls. SKU/zone/Retail/history adapters,
coordinator/public execution, Full Cost pagination, B6/B7 and residual module-only
GO-2026-5932 remain separate. Laptop/Azure authorization is unconfirmed. No Azure
resources/fixtures/roles, production-containing scope substitution, listener exposure
or release. Record actual failures and maintain handover/ledger/roadmap after the slice.

## Final target admission details, 2026-10-13

The implemented internal candidate uses ProviderLocationRequest and an owned
ProviderLocationCollection with SubscriptionID, Status, Locations and FailureCode.
No target region is invented. The coordinator/public path remains absent.
Optional nextLink absence or empty string terminates; present null/wrong shape is
incomplete. Namespace/type segments begin with ASCII alphanumeric and otherwise
use ASCII alphanumeric/dot/underscore/hyphen; nested types use nonempty segments.
Locations admit ASCII display strings, remove ASCII spaces, lowercase, then require
at most64normalized characters/physical-format or explicit global. Returned IDs
are structurally validated before ASCII-case-insensitive binding. All projected
raw text and duplicate work are charged; normalized joined keys remain bounded.

Existing whole JSON admission is reused unchanged: token/depth/key/Unicode limits
apply also to ignored metadata. Collect's8MiB response cap dominates the existing
roughly16MiB decoded-text guard; exact/excess text tests exercise the decoder seam
independently and are not claims that a real collection reaches that bound.
The eight-page duplicate-location fixture has65,536raw provider/type/location
entries but only eight returned type maps, proving work is charged before dedup.
Its first excess rejects the entire eighth page, retaining seven exact earlier maps.

Self-review corrected wrapped cancellation/deadline transport text to standard
sentinels. Tests require nil result and exact sentinel, not just errors.Is.
The temporary formatting workflow was preparation only and is removed from the
final acceptance candidate. Fresh complete final-head native/source/review/protected
and distinct accepted-push proof remains required; no prior pilot proof transfers.
