# Compute SKU source characterization: pre-edit contract, 2026-10-13

Status: UNACCEPTED source-only preparation. Publish/read back before harness edits.
Provider target PR137 merged at af18468fe1dfc697838115c0a7585ba675fdec97;
candidate full native/source/review proof is complete. Distinct accepted-push source
38334004104/115089899014 is verified; native38334004115 is pending at preparation.
Read [live PR137](https://github.com/DeBoX85/Cloud-Assess/pull/137) for final outcomes.
Candidate proof cannot certify accepted-push or this next proposal.

## First slice and source authority

Observe original common pager, shared restriction helper, registry/cache behavior
and Compute VM/VMSS/disk providers. No target SKU collector, coordinator/public
region activation, SDK refresh, pin change or Azure operation belongs in this slice.
Storage, Cognitive Services and SQL provider request captures follow separately.
Metadata may verify all seven registered types without claiming those APIs executed.

Pinned source AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1
and APRL60eaddda76541f6adbc1c5ffa686829807e55e29 remain unchanged.

| Source under internal/scanners/plugins/region/sku | Original Git blob |
| --- | --- |
| provider.go |7f71b296d332001ae12e49a021b6fee44f88a970|
| cache.go |5be9c158fd38c609b06ffc67537cd07b838fc9d5|
| restrictions.go |8d053daf58c7c8a18f39b579e95d4a53f8ebd55a|
| compute.go |ee9300e5f8652591fe43121d0973321742ee70e3|

Actual pinned directory inventory has seven files and seven registered resource
types across FOUR ARM namespaces (Compute, Storage, CognitiveServices, Sql), with
FIVE endpoint/response families because Compute uses both vmSizes and skus.
The earlier inventory's “five provider families” wording must not be read as five
ARM namespaces. Preserve historical text and record the clarified current count.

## Literal requests versus service-validity uncertainty

VM and VMSS both call:
GET /subscriptions/{subscription}/providers/Microsoft.Compute/locations/{region}/vmSizes?api-version=2023-03-01.
Disk calls:
GET /subscriptions/{subscription}/providers/Microsoft.Compute/skus?api-version=2023-03-01&$filter=location%20eq%20%27{region}%27.
The source hardcodes the public management.azure.com origin. Capture literal URL,
raw path/query, method/body/audience and ordered continuations unchanged.

Fresh [Resource SKUs documentation](https://learn.microsoft.com/en-us/rest/api/compute/resource-skus/list?view=rest-compute-2026-04-01)
specifies2021-07-01 and a location filter; current
[VM Sizes documentation](https://learn.microsoft.com/en-us/rest/api/compute/virtual-machine-sizes/list?view=rest-compute-2026-04-01)
marks vmSizes deprecated. Preserve source2023 literals during characterization.
This discrepancy is not proof that Azure rejects2023, and synthetic success does
not establish current service support or deployability. Before target collection,
reconcile official schema/version/semantics and publish deliberate corrections;
never silently substitute a version in the pinned source or its oracle.

## Supported source execution seam

Reuse the accepted REST quota harness's original az.NewHttpClient with fake
TokenCredential and explicit supported HttpClientOptions.Transport. Original SKU
providers take that client directly. Inspect pinned constructor/options before
coding; do not patch source/SDK bytes or assume default HTTP globals alone suffice.
The fixture transport handles every allowed synthetic request, records exact
requests and returns bounded declared bodies. Fake public ARM audience/token,
short contexts and no credential discovery. Tripwire default HTTP globals as a
supplementary guard. No actual Azure/AAD/external HTTP request or registration.

Explicit no-retry fixture options differ from source production defaults; label
that seam and do not claim default retry policy, sovereign behavior or live parity.
Unsafe/foreign continuation probes may reach ONLY the injected fake transport.
Bound cycle probes by fixture error/context; never run an infinite source pager.
Unexpected method/body/destination outside the declared fixture matrix is a harness
failure, not source acceptance. Keep source/module/APRL/untracked guards intact.

## Independently declared complete observation matrix

| Group | Required original observations |
| --- | --- |
| Registry | Exact seven type keys, case lookup, unknown lookup, duplicate registration panic without mutating original registry |
| VM/VMSS | Exact shared vmSizes requests, single/two pages, lowercase names, empty-name skips, missing/null/wrong envelopes/items, full returned availability maps |
| Disk | Literal encoded filter, name/tier fallback, ignored resource-type metadata, collision/empty-name behavior, full states/zones and original page loss on error |
| Restriction helper | Location reason precedence; Zone order/duplicates and empty zones; available capability false/empty/mixed case; combined restrictions/capabilities; full state plus zones |
| Common pager | First/later status/decode failure, complete empty versus source empty-success shapes, duplicate/lowercase keys within/across pages, relative/foreign/malformed/cyclic continuation, cancellation |
| Cache | Subscription/type/normalized-region identity, cached reuse, unsupported types, failed-call retry, VM versus VMSS request separation, full caller map ownership/mutation observations |
| Concurrency | Bounded controlled same-key coalescing and separate-key requests; source returned-map sharing recorded without deliberately introducing a data race |
| Decode/text | Single and duplicate ASCII/Unicode field aliases, null entries and invalid text where safely observable; complete raw input retained |

Independently declare expected full maps, states, zones and exact request sequences.
Do not derive expected values from the implementation being tested. Do not replace
unknown/incomplete/collision behavior with empty success. Source comments are not
evidence that a returned size is currently deployable or quota-sufficient.
Cache map iteration can make different-state trim collisions nondeterministic:
record/classify that limitation instead of selecting a convenient winner, modifying
source or using flaky probabilistic golden assertions. Deterministic raw captures
must retain exact bytes; any separately observed nondeterminism must be explicitly
bounded and cannot be disguised as deterministic full-input certification.

## Retention and native source controls

Proposed compute_sku_capture_test.go.txt injects ONLY a test into original sku
package and invokes original functions/providers/cache with supported transport.
Add tenth harness and source-compute-sku-inputs.json/source-compute-sku-outputs.json.
Strict runner then requires21 records (original19 untouched plus two), emits full
raw provenance/contiguous chunks before unconditional equality. Extend the actual
manifest/full-loop byte guard and its independent compiling omission/subset checks.

First unretained observation is a FAILED strict-golden preparation run after
successful harness/provenance/guards, not acceptance. Classify any earlier
compiler/fixture failure from actual logs. Reconstruct original raw bytes and
independently SHA256/full remote Git identities before retention. Add literal
complete source-cell comparisons and compiling named fixture controls/restored
assertions on BOTH mandatory native hosts. Compiler failure never counts as
mutation detection; inspect both match and replacement expressions after edits.

Require fresh complete final-head source/native/review/fullpublication/identity/
tree/preview/protections, protected expected-head acceptance and distinct
accepted-push proof. No source-only result certifies future target SKU runtime.
Keep existing25provider production and all earlier controls/security budgets/pins.

## Exact next and remaining boundaries

Verify PR137 distinct push completion, publish/read back this source contract,
then build the isolated Compute/common/cache harness. No new target production code.
Other SKU families, target bounded collectors, zone/history/Retail/coordinator/
public execution, Full Cost pagination, B6/B7/module-only advisory/live/load/freshOS/
maintenance/Gate004/release/licensing/human approval remain open. Laptop/Azure
restoration unconfirmed; no fixtures/resources/roles/production scope substitution/
listener exposure/release. Update handover/ledger/roadmap/failure history after each
material slice and before ending; remote WIP stays unaccepted until its own QA.
