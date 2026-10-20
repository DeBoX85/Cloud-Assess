# VM/VMSS SKU pre-edit contract, 2026-10-16

Status: unaccepted pre-edit contract and future expectations. No proposed collector API or VM/VMSS production implementation exists in this slice. User authorized handover/status reconciliation and continuation after [Storage completion](STORAGE_COLLECTOR_COMPLETION_20261016.md). Branch docs/vm-vmss-contract-20261016 starts at PR145 a14128bc94bf531dd1c28553db14d2be011b2214. Protected bootstrap/core-v1 remains a6fee427e19feeef28691737ea9bc5026086efa0. Keep this successor draft/unmerged. Its live PR body is its final exact-head evidence index; earlier green results do not certify it.

## Primary authority and deliberate corrections

Preserve all 23 original captures and AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1/APRL60eaddda76541f6adbc1c5ffa686829807e55e29. Original compute.go blob ee9300e5f8652591fe43121d0973321742ee70e3 calls vmSizes2023-03-01 for both VM and VMSS; retain its observed requests/maps verbatim. Synthetic source success does not prove service validity, deployment eligibility, quota or capacity.

Fresh primary references inspected2026-10-16:

- [VM Sizes List](https://learn.microsoft.com/en-us/rest/api/compute/virtual-machine-sizes/list?view=rest-compute-2026-04-01) marks vmSizes deprecated and directs clients to Resource SKUs.
- [Resource SKUs List](https://learn.microsoft.com/en-us/rest/api/compute/resource-skus/list?view=rest-compute-2026-03-02) documents the subscription catalog,2021-07-01 and optional location filter.
- [VM size overview](https://learn.microsoft.com/en-us/azure/virtual-machines/sizes/overview) applies sizes to VMs and both scale-set orchestration modes. It does not prove every SKU eligible for a particular VMSS.
- [VMSS List SKUs](https://learn.microsoft.com/en-us/rest/api/compute/virtual-machine-scale-sets/list-skus?view=rest-compute-2026-03-02) requires an existing named scale set/resource group; it is not a target-region subscription catalog.

Immutable schema authority: Azure/azure-rest-api-specs commit1cdcc5414764d20e8815554c8b89708c10747c18 under specification/compute/resource-manager/Microsoft.Compute/Compute/stable/2021-07-01. Complete schema and three complete examples are preserved unchanged in testdata. [Primary identities](VM_VMSS_PRIMARY_IDENTITIES_20261016.json) records byte counts, original Git blobs and SHA256. Original blobs independently checked against primary directory listings:

| Retained suffix after compute-primary-2021-07-01- | Original Git blob |
| --- | --- |
| skus.json | b85286923a6382356ab773d90faac382e73740ec |
| ListAvailableResourceSkus.json | 9d403f67699adfedf2714ada315ba5215f793c81 |
| ListAvailableResourceSkusForARegion.json | def2413e26e98f88fdbc46448b1b1511fbdf4fb6 |
| ListAvailableResourceSkusWithExtendedLocations.json | f1c254f72d650831db106561306a1e102a718a02 |

Copyright(c)2017Microsoft/MIT, with [complete original license](../internal/plugins/region/testdata/storage-primary-LICENSE.txt). Root NOTICE and exactly its generated section are synchronized; no module inventory refresh. The schema's external common-types references are not vendored or executed here; they cover shared subscription/version/error types, while the full retained skus.json supplies the SKU definitions used below.

All three examples contain virtualMachines Standard_A0/Standard_A1 in westus, omit restrictions and end with null nextLink. The first row has zoneDetails with exact uppercase wire property Name. The extended example adds EdgeZone locations Las Vegas/Seattle/Portland. Do not rewrite omissions into empty arrays. Proposed complete diagnostic groups for both names are Known=false/state2/empty zones; projection for required A0/A1 is unknown/empty. These primary examples are not independent proof of unrestricted deployment.

## Proposed internal boundary

First implementation supports the VM size catalog only. VMSS deployment eligibility/scoring remains unsupported pending a separate primary-supported contract. Do not silently copy a VM catalog map into a VMSS scoring key or claim this contract completes either production adapter.

Proposed ComputeVMSKURequest has SubscriptionID, Region, ResourceType. Allowed request types are canonical microsoft.compute/virtualmachines and microsoft.compute/virtualmachinescalesets. Both may collect virtualMachines catalog diagnostics, but VMSS-bound collections MUST project unsupported/empty. This labels underlying node catalog diagnostics without asserting VMSS eligibility.

NewComputeVMSKUCollector(endpoint,getter) validates configured HTTPS ARM origin using the existing safe constructor and RESTQuotaGetter.GetBoundedWithResponse. Collect(ctx,stableSelectedMap,request) validates unique normalized selected UUIDs/nonempty display names, exact selected subscription, lowercase physical-format target and resource type before credentials. Canonicalize UUID case; reject absent scope, normalized collisions, global/logical/stage/display targets and arbitrary aliases. No widening of accepted quota/provider/disk paths, no shared mutable cache, and no SDK/source edits.

First GET has empty body, configured origin, path /subscriptions/{lowercase UUID}/providers/Microsoft.Compute/skus and canonical url.Values query: %24filter=location+eq+%27{region}%27&api-version=2021-07-01. [The future manifest](../internal/plugins/region/testdata/compute-vm-target-contract-v1.json) includes the complete literal public-ARM eastus URL. No vmSizes fallback, resource-specific VMSS call, includeExtendedLocations, registration, tenant discovery, Azure call or public command belongs here.

ComputeVMSKUCollection owns canonical SubscriptionID/Region/ResourceType, Status complete/partial/unknown, Groups map[name]ComputeVMSKUGroup{Known,State,BlockedZones}, and sanitized FailureCode. Unknown groups are Known=false,state2,empty zones. Every empty map/slice is allocated. Evidence(ctx,requiredNames) is the sole supported scoring projection. Stable caller inputs are required; no racing caller mutation guarantee.

## Admission and optional metadata

Require valid UTF8, single complete JSON document/EOF, no invalid-surrogate replacement, controls/noncharacters, case-folded duplicate keys or noncanonical known-property spelling. Enforce depth32,properties128/object,tokens65536/page independently. Unknown properties may be retained only after full raw-work/text/JSON validation. Canonical zoneDetails.Name is uppercase; SKU/capability name is lowercase.

Root is object with explicit non-null value array; rows are non-null objects. Every row requires nonempty unpadded printable ASCII name/resourceType, bounded512bytes. No tier/size fallback or trimming. Matching virtualMachines names lowercase; any normalized duplicate within/across pages rejects the entire candidate page, including identical duplicates. Validate supplied known metadata before filtering other types. Wire virtualMachineScaleSets rows do not become VM rows. Page admission is atomic.

| Known property | Shape when supplied | Policy |
| --- | --- | --- |
| tier,size,family,kind | Non-null bounded string | Diagnostic only, no identity fallback |
| locations | Non-null array of ASCII physical-format regions; ASCII case normalized | Missing makes matching group unknown; empty explicitly not listed; display/logical/stage entries reject page |
| restrictions | Non-null array of non-null canonical objects | Missing unknown; empty explicitly unrestricted catalog classification |
| capabilities | Non-null object array; supplied name/value non-null strings | Missing either makes matching group unknown; complete unknown capabilities diagnostic |
| capacity | Non-null object; supplied minimum/maximum/default nonnegative int64 JSON integer; scaleType string | Diagnostic, never allocation/scale-set capacity evidence |
| costs | Non-null object array; supplied meterID/extendedUnit strings, quantity int64 JSON integer | Diagnostic, not Retail price |
| apiVersions | Non-null string array | Diagnostic, not version negotiation |
| locationInfo | Non-null object array; supplied location physical region, zones positive-decimal string array, zoneDetails object array | Does not replace locations or grant availability |
| extendedLocations/type in locationInfo | Non-null bounded string array / bounded string | Opaque EdgeZone metadata, not physical aliases |
| zoneDetails.Name/capabilities | Non-null positive-decimal string array / capabilities as above | Exact uppercase Name, metadata only |
| restrictionInfo.locations/zones | Non-null physical-region / positive-decimal zone arrays | Restriction scope, always validated |

Missing optional diagnostic metadata does not poison evidence; present null/wrong shape rejects the page, including filtered rows. Future capacity scaleType/extended-location type strings are opaque bounded metadata. Restriction enum uncertainty affects interpretation and makes its group unknown instead. Zone labels are positive decimal strings without sign/leading zeros, maximum512bytes. Duplicates charge work/text before deduplication.

## Restrictions and existing consumer

Recognized restriction literals are exact schema enums Location/Zone and QuotaId/NotAvailableForSubscription. Location scope is union of supplied values and restrictionInfo.locations; either matching target sets state1 for either reason. At least one supplied scope must be nonempty; otherwise group unknown. Zone requires nonempty restrictionInfo.locations and zones; target zones union/deduplicate/lexically sort. Zone values, if supplied, are validated physical-region strings but cannot replace restrictionInfo scope. Other-region restrictions are ignored only after validation. Missing/unknown enum or valid but uninterpretable scope produces unknown group; malformed supplied shape rejects page.

Precedence: unknown dominates; otherwise region absent from explicit locations is state2; explicit available capability false is state2; applicable regional prohibition state1; applicable zone prohibition state3; otherwise state0. States1/2 clear zones. Available capability name and true/false compare ASCII case-insensitively; conflicting values or future value produce unknown. Other complete capabilities do not manufacture restrictions. Do not permit a later row/restriction to overwrite an earlier prohibition.

Evidence validates collection binding, canonical names, status/fixed codes, Known/state consistency and owned sorted unique positive zones on state3 only, plus required-name bounds. Invalid forged complete collections return sanitized error, including complete with FailureCode, noncanonical/duplicate keys, bad states/zones or unknown status. Nil or valid partial/unknown projects unknown/empty. Valid complete VMSS-bound collection projects unsupported/empty. Complete VM projection normalizes required names, rejects normalized duplicates, returns unknown/empty if any REQUIRED group is unknown, and maps required absent names to state2. Unknown unrequired rows do not poison projection. Empty VM requirements produce complete/empty; VMSS stays unsupported/empty. Check context before every unknown/unsupported early return and final copy. Return new maps/slices each time.

Future coordinator must bind subscription/region/type before assignment. CalculateAvailability stays unchanged: unknown counts/warnings remain; all-unknown percentage can be100 under existing denominator policy and is not availability proof. No coordinator, report/public schema or scoring change here.

## Paging, failure and exact bounds

Omitted/empty/null nextLink ends paging; other types reject page. Continuation resolves to HTTPS at exact configured origin/path/subscription/version and SAME mandatory location filter. Exactly one version/filter, at most one nonempty decoded $skiptoken4096bytes, no unknown/duplicate/aliased query, userinfo/fragment/opaque/encoded path/redirect. Full nextLink limit8192bytes. Canonicalized URL cycle detection occurs before next authenticated GET. Never follow unsafe continuation or accept truncated paging as complete.

Context/terminal getter errors dominate. Present non200 status dominates nonterminal error/decode, including response+error pairs. All statuses including404/405 are operational unknown, not proven unsupported API. No response/error-derived diagnostic text. Invalid request/origin/projection returns sanitized error. After whole admitted pages, operational failure yields partial with previous diagnostic groups; before any admitted page, unknown/empty. Empty admitted page counts. A page whose nextLink is unsafe is admitted before partial continuation failure without another request. A malformed page contributes no groups. Cancellation/deadline discards prefix and returns nil/error with errors.Is compatibility.

| Independent bound | Inclusive limit | First-excess requirement |
| --- | ---: | --- |
| Selected subscriptions / scalar label | 1000 / 512bytes | +1 before getter or admission |
| Body/page / successful aggregate bytes | 1MiB / 8MiB | +1, no tighter guard masking assertion |
| Pages | 64 |64th final completes; continuation65 partial without request65 |
| Matching groups / required names |8192 each |8193 within independent raw/page budgets |
| Raw collection work |65536 | Every object AND array element, filtered/unknown data included; +1 |
| Decoded collection text |16711680bytes | All property names/decoded strings including duplicates before filtering |
| Projection logical work |65536 | Each group, zone and required entry once, no repeated-pass double charging |
| Projection input / output text |16711680bytes each | Binding strings, keys/zones/raw requirements each once; output independently bounded |
| JSON depth / keys / tokens |32 /128 /65536 | Exact and +1 independently |
| Decoded skiptoken / nextLink |4096 /8192bytes | Exact and +1, no request after rejection |

Raw work/text span admitted pages and candidate validation before atomic commit; JSON budgets per page. Root nextLink alone may exceed scalar512 but still charges text. Projection limits are independent of wire limits and need fixtures that actually reach them. The collection decoded-text cap is dominated by the tighter8MiB aggregate wire cap for valid JSON; do not claim an end-to-end exact text-cap case is reachable. Test that accumulator through a private accounting seam and independently test the public aggregate bound; never weaken production wire bounds to force reachability. Literal boundary sizes must not merely refer to implementation constants. No RSS/latency claim.

Fixed operational codes: vm_sku_request_failed, vm_sku_status_invalid, vm_sku_byte_limit, vm_sku_aggregate_limit, vm_sku_invalid_response, vm_sku_collision, vm_sku_work_limit, vm_sku_text_limit, vm_sku_page_limit, vm_sku_continuation_invalid, vm_sku_cycle. Error codes: vm_sku_request_invalid, vm_sku_projection_invalid, vm_sku_projection_limit. Unknown/partial diagnostics never become complete scoring evidence.

## Frozen expectations and stopping point

The manifest freezes67literal full collection/evidence/request scenarios before production. The three complete primary examples use body_file plus /responses/200/body. These FUTURE cases have NOT executed against an implementation. Current preparation checks JSON/unique names/state consistency, complete primary Git identities and authored links only, not collector correctness.

Later implementation first publishes a compiling runner consuming every frozen case and inspects expected missing-API compilation on BOTH native hosts as preparation. Then implement bounded collector, independently test every exact/excess bound, actual bounded HTTP client/audience/GET/body/status+error/context/closure/redirect path, atomic paging, forged projection, complete map/scoring and sequential/concurrent ownership. Fixture credentials/transports are offline; no real credential discovery/network.

Compiling production reversals must cover version/filter/scope, JSON admission/uppercase Name, primary unknown handling, missing restrictions, region absence, both reasons/scope union, zone precedence/ownership, required unknown/unrequired isolation, VMSS unsupported projection, collision, atomic page/prefix gating, cancellation including late unknown/unsupported, body/aggregate/pages/work/text/continuation/cycle and output ownership. Every mutation compiles, fails its named independent assertion, restores original bytes including CRLF and passes baseline. Preserve all353Linux/348Windows existing identities and all mandatory gates.

For THIS checkpoint: publish/read back changed files and exact tree/parent/human identity; require own complete Linux/Windows/source gates, full logs/control identities/scans/packages/source reconstruction; index final evidence then STOP. No VM production edit before that boundary. Correct any discovered contract inconsistency before implementation, retaining preparation history. Rollback uses targeted revert after successor review, never reset unrelated work, force protected history or delete expectations to hide failure.

Next distinct slice after evaluation: bounded VM catalog implementation; VMSS scoring remains unsupported pending separate eligibility evidence/contract. Storage remains unaccepted; accepted adapters1/7, source characterization7/7. Other families/coordinator/public execution/FullCost/module-only advisory/live/nonempty/load/sovereign/freshOS/Arc/maintenance/Gate004/license/security/independent-person/release obligations remain open. No Azure operations/fixtures/RBAC, dependency refresh, sub-agents, protected merge or release.

Local preparation checked67unique future cases, request/projection/error consistency, all4complete primary Git blobs/byte counts/SHA256, exact root-NOTICE-only generated replacement and1,203authored links. No local full Go suite is claimed in this recovered Windows environment; own fresh hosted gates remain pending at publication.
