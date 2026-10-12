# Provider/resource-type locations: pre-edit contract, 2026-10-12

Status: UNACCEPTED preparation on feat/region-provider-locations. This contract
must be published and read back before source-harness or target production edits.
Existing security PR133 is verified offline; checkpoint PR134 merged through the
protected path at50f5bef8a0bb34a0717d5adc6f6fa44039ed62c6, treeaef7b953f39838ad2038a4f8147eb2087b1331d3,
ordered6bcbca46/d4079669 parents. PR134 candidate proof completed; its DISTINCT
accepted-push source proof completed and native proof is still pending at contract
preparation. Verify [live PR134](https://github.com/DeBoX85/Cloud-Assess/pull/134).
Preparation does not transfer that proof or certify this new head.

## Requirement, scope and pinned authority

Implement the provider/resource-type location obligation selected by
[REGION_ADAPTER_INVENTORY_20261011.md](REGION_ADAPTER_INVENTORY_20261011.md), preserving
[REGION_AVAILABILITY_RUNTIME.md](REGION_AVAILABILITY_RUNTIME.md) arithmetic/global
location semantics and honest declared provider completeness. This is one adapter
boundary. SKU lists, zone counts, CostManagement/Retail, per-run coordinator,
registry/scanner/public CLI/all formats and live parity remain separate.
No target collector exists yet; this phase first captures the literal source.

All references remain AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1
and APRL60eaddda76541f6adbc1c5ffa686829807e55e29, with other catalog/SKU/source pins intact.

| Primary source | Original Git blob / selected version |
| --- | --- |
| AZQR internal/scanners/plugins/region/availability/availability.go: fetchAllProviders and CheckRegionsInParallel |55e59b9e061dd31d27cc3446ad00659bd6961174|
| AZQR internal/scanners/plugins/region/types/types.go: NormalizeRegionName and ResourceTypeLocationData.IsAvailable |f9b61fb74075e5a2297be1b898067b8cb047f9b3|
| AZQR internal/az/client_options.go: actual option factory |73693e3730f5ed51132a552b79d5e95078721a10|
| AZQR internal/az/cred.go: supported custom cloud endpoint/audience |9677ae498c8d22ebbc77a1fe3a47e17579f3c4d5|
| Azure SDK armresources/providers_client.go |armresources v1.2.0 / e1612d16090a9a66c4b2bb9fa5697e7376263992|
| Source-selected azcore runtime/pipeline.go |azcore v1.23.1 / 3d95fe30dde1c0141f3b2737d480573ec048ac33|
| Source-selected azcore runtime/transport_default_http_client.go |928e9bf92550dcd83180b7fecff40b7a50fd7ff9|

Source go.mod selects armresources1.2.0/azcore1.23.1. Target's existing graph has
armresources1.1.1 without current CLI membership. Do not substitute that SDK when
characterizing source, refresh SDKs, modify source/SDK production bytes, or advance
any pin. Future target collection should use the existing bounded ARM getter,
without new armresources CLI import or incidental dependency change.

The source calls NewProvidersClient(subscriptionID, credential, az.NewDefaultClientOptions())
and NewListPager(nil). Generated SDK and Microsoft's
[Providers List2021-04-01 contract](https://learn.microsoft.com/en-us/rest/api/resources/providers/list?view=rest-resources-2021-04-01)
agree on GET /subscriptions/{UUID}/providers?api-version=2021-04-01; nil options
do not request expand. The envelope supplies provider objects and continuation.
No tenant-scope list, register/unregister, subscription changes or write capability
belongs in the target adapter. SDK pipeline defaults must also be reviewed; GET
syntax alone is not a guarantee against SDK registration or redirect behavior.

## Literal source behavior to characterize independently

Namespace/type names are lowercased without trimming. Location normalization
removes ONLY ASCII spaces then lowercases; tabs and other whitespace are not
equivalent. Each normalized namespace/type stores a location set; normalized
duplicate type declarations overwrite earlier sets, including across pages.
Global means available for every queried region. Namespace or type absence,
known empty location sets and unknown resource types remain distinct observations.

A pager error returns nil/error and discards earlier provider pages; its caller
logs the failure and replaces provider data with an empty map before availability
calculation. It does not expose collection completeness in that comparison.
Missing/null provider fields, null pointer entries, duplicate JSON fields, malformed
envelopes, continuation, cancellation and caller fallback must be observed, not
silently normalized into target success or treated as a live incident.

Produce complete source-owned cache maps, exact canonicalized request sequences,
failure/panic/cancellation outcomes and complete caller comparison rows for
synthetic cases. Literal case declarations and exact assertions are independently
defined in the harness; target-generated results cannot become the source oracle.

## Source execution seam and safety

fetchAllProviders does not accept an injected client option. Source-selected
azcore constructs its OWN default HTTP transport; changing http.DefaultTransport
alone does not intercept it. Do not patch production code or use private SDK
linker tricks and then call it unchanged source execution.

Use the source's existing custom-cloud environment selection with a LOOPBACK TLS
fixture endpoint, explicit synthetic audience and fake TokenCredential. Trust
only the fixture certificate through Linux SSL_CERT_FILE/SSL_CERT_DIR before the
first handshake, keeping TLS verification enabled. Record this environment seam
and its limits in capture metadata. This is a proposed harness mechanism pending
actual hosted verification, not a passed TLS/source test. A separate injected
SDK NewListPager(nil) fixture checks the public-cloud request/audience contract
without calling Azure. Any default-client tripwire is supplementary, not proof
it covers the SDK's own transport.

The source capture workflow is Linux. Bind only loopback, use fixture-controlled
literal loopback continuations, short per-call contexts, fake tokens, no Azure
credential discovery and no Azure/AAD endpoints. Redirect/foreign-origin probes,
if included, must terminate on a second loopback fixture and record behavior;
never point an unsafe probe at real/external hosts. Unexpected requests fail
with bounded fixture responses/context deadlines. No retry-prone status test may
outlive the harness budget. Inspect source cloud/option/blob fingerprints and
unchanged source/APRL/go.mod/go.sum guards after the isolated run.
These are synthetic source/SDK observations, not live cloud or Windows SDK-capture
proof. Existing mandatory Windows checks test the new retained-data controls.

## Required source matrix and acceptance

| Case group | Complete observations required |
| --- | --- |
| Literal request and two pages | GET/method/body/query/audience, no expand, exact ordered next request and complete final cache |
| Normalization/set/global | namespace/type case, nested type path, spaces versus tabs/non-ASCII, repeated locations, global cross-region result |
| Presence and empty | missing/null/empty value; namespace/resourceTypes/resourceType/locations; null provider/type/location; complete empty versus failure |
| Collision | repeated normalized namespace/type within/across pages; later overwrite-to-empty and independently expected cache |
| Failure | first/later403/404/405/400, malformed JSON/wrong shapes, complete loss versus caller fallback; sanitized synthetic records |
| Cancellation and continuation | cancelled before request and after first page; literal/cyclic/relative/malformed continuation behavior with bounded test observation |
| Response alias/decode | duplicate and case-fold field aliases, invalid text/shape where safe to observe; exact raw input retained |
| Ownership topology | full cache and caller comparison rows, not counts; source/input/blob/provenance and exact request record |

Do not deliberately run an infinite source pager. Any cycle probe must be stopped
by fixture/context budget and retain its aborted/error outcome. Do not call that
successful full pagination. Unexpected source panic is preserved as a classified
record; compiler/transport setup failures are not characterization success.

Proposed files: testdata/provider_locations_capture_test.go.txt,
source-provider-locations-inputs.json and source-provider-locations-outputs.json;
extend scripts/capture-region-source.py to execute the new harness and require
byte equality, emitting complete provenance and contiguous chunks BEFORE golden
comparison. Existing17 records and eight harnesses remain byte-identical. The
strict runner must never skip missing/changed provider goldens or silently accept
a partial log. A deliberately unretained first observation is a failed preparation
run; retain its complete synthetic records, then publish exact goldens only after
independent semantic/hash/provenance review.

Add source-cell/structure tests with literal request/map/global/empty/fallback/
collision/panic/cancellation expectations, complete golden comparisons and
named compiling source-fixture mutations/restored baselines in BOTH mandatory
native hosts. Reuse existing capture byte guard, extending its actual file set;
verify newly added files are included rather than assuming filename patterns.
No implementation-mirroring count-only tests or compiler-failure credit.

## Future target corrections and proposed bounded interface

These are obligations for the subsequent production adapter; this preparation
does NOT implement or certify them. Reconcile final target details against the
accepted literal captures before changing production.

- Selected UUID/scope validation before any request. Use existing bounded
  authenticated getter/ARM endpoint and audience/closed response/no redirect
  contract; source retries/automatic registration are not copied blindly.
- Exact selected subscription provider-list path and API version; validate each
  continuation BEFORE credentials. Same HTTPS origin, selected path and version,
  no userinfo/fragment/opaque/encoded-path tricks, duplicate query aliases or cycle.
  Unsupported service continuation is honest incomplete evidence, not truncation
  labelled complete. Preserve opaque token meaning; no unrelated URL fetch.
- Owned per-run result shaped for existing AvailabilityEvidence.Locations and
  ProvidersComplete, with normalized subscription, complete/partial/unknown status
  and sanitized failure code. No mutable cache, borrowed maps or provider-derived
  metadata in errors/logs. Invalid request/cancellation returns an error/no result;
  operational response failure preserves only earlier wholly admitted pages.
- Complete empty array is successful empty collection. Missing/null/wrong-shaped
  mandatory response declarations are incomplete, never inferred complete.
  Validate optional returned provider ID structurally before equality when present,
  binding to selected subscription/namespace; absent optional ID is not fabricated.
  Namespace and nested resource-type identifiers must be bounded unambiguous ASCII.
- Preserve lowercase/ASCII-space location normalization and explicit global.
  Unknown/nonphysical malformed location metadata must not be silently made
  deployable. Characterize exact null/empty decisions before freezing target code.
- Reject normalized type collisions rather than silently last-write-wins; admission
  is page-atomic, so a conflicting page cannot partially overwrite earlier owned
  evidence. Preserve only earlier admitted pages with explicit partial health.
  Repeated locations within a valid type are source set deduplication, with raw
  work charged before deduplication.
- Proposed limits:64pages,1MiB per response,8MiB successful aggregate bytes,
  65536raw provider/type/location work entries and existing decoded-text/label
  guards, with JSON depth/token/UTF8/surrogate/Unicode field-alias checks BEFORE
  expensive projection. Charge duplicate/raw work and all response bytes, not
  only returned or selected maps. Exact/excess bounds require independent tests.
  These are safety bounds, not representative-estate/load/live certification.
- Existing availability arithmetic/health integration must prove complete/partial/
  unknown provider evidence, owned state and cancelled/error behavior. It remains
  an internal collector; registry/public region execution is not enabled.

Future target tests must use a credential canary/actual bounded transport for
destination/method/body/audience/no-write/body-close/redirect/continuation cases,
independent literal cache/health expectations, late failure, malicious scope/IDs,
normalization/collisions, exact/excess raw/byte/text/page limits, concurrent/sequential
runs and caller-data ownership. Add compiling named guards and restored assertions,
not just fixture lengths. Keep production target code out of this capture slice.

## Verification, compatibility and recovery

Require fresh complete exact-head source/native/automated-review/full authored
bytes/original Git tree/ordered parents/Denis identity/tested preview/protection,
protected expected-head merge and DISTINCT accepted-push proof for each accepted
slice. Candidate/source-only proof cannot certify target production or live work.
Record every failed setup/control/golden/hash attempt and actual disposition;
fix only justified failures without relaxing source guards or mandatory gates.

Current internal interfaces/arithmetic, public registry, SDK/dependency/toolchain,
source/APRL/catalog/SKU/license/notice pins and existing goldens stay unchanged.
Normal rollback is a reviewed new revert on then-current protected branch, never
reset/force/history deletion; later dependencies must be reviewed.

Exact next: finish PR134 accepted-push verification, publish/read back this pre-edit
contract, then build/verify the isolated provider harness, retain/review complete
original synthetic bytes and add literal source controls. Only after source
characterization acceptance begin the target adapter. Laptop/Azure access remains
unconfirmed and is not needed for this independent synthetic phase. Live nested
nonproduction DV001, optional plugins/restricted visibility/Diagnostics/Advisor/Arc,
load/freshOS/hosted maintenance/Gate004/release/licensing/human approval stay open;
no Azure resources/fixtures/roles, production-containing scope substitution,
service exposure, listener or release is authorized.
