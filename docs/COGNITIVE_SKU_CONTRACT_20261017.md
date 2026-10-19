# Kind-bound Cognitive account SKU contract, 2026-10-17

Current authority: [CURRENT_STATE](CURRENT_STATE.md), the live task PR and [three-task workstream](COGNITIVE_SKU_WORKSTREAM_20261017.md). Parent is verified draft PR155/002fa14846ecf8115bd3a126080738e2f589a12c. [Reconciliation](COGNITIVE_SKU_RECONCILIATION_20261017.md) records primary facts versus deliberate target policies. This checkpoint freezes an internal API and143future scenarios; the API is absent and these scenarios have NOT executed against it. No scoring, inventory/public schema, coordinator/report/CLI, source/pin/SDK/module or Azure change.

## API and provenance boundary

```go
type CognitiveSKURequest struct {
    SubscriptionID, Region, Kind string
    LocationAliases map[string]string
}
type CognitiveSKUGroup struct {
    Known bool
    State int
    BlockedZones []string
}
type CognitiveSKUCollection struct {
    SubscriptionID, Region, Kind, ResourceType, Status string
    Groups map[string]CognitiveSKUGroup
    FailureCode string
}
func NewCognitiveSKUCollector(endpoint string, getter RESTQuotaGetter) (*CognitiveSKUCollector, error)
func (c *CognitiveSKUCollector) Collect(ctx context.Context, selected map[string]string, q CognitiveSKURequest) (*CognitiveSKUCollection, error)
func (c *CognitiveSKUCollection) Catalog(ctx context.Context, selected map[string]string, expected CognitiveSKURequest, required []string) (CognitiveCatalogEvidence, error)
type CognitiveCatalogEvidence struct {
    SubscriptionID, Region, Kind, ResourceType, Status, DeploymentEligibility string
    Values map[string]SKUAvailability
    Warnings []assessment.AssessmentWarning
}
```

The fixed resource type is microsoft.cognitiveservices/accounts. One explicit service Kind is bound to request/collection/catalog and normalized to lowercase ASCII; no inference from S0/F0, name-only inventory aggregates or model names. The caller supplies actual kind provenance. Structural binding is not authentication of caller-built evidence. Stable non-racing inputs and non-nil context are required. No implicit SKUEvidence conversion/Evidence method, registry, scoring or public integration is added.

Every successful Catalog result owns allocated Values and Warnings, canonical subscription/region/kind/type, DeploymentEligibility=not_assessed and exactly one warning:

- code: cognitive_sku_catalog_only
- message: Cognitive account SKU catalog evidence does not establish model deployment eligibility, quota or allocation capacity.

## Request, transport and pagination

Reuse NewRESTQuotaCollector's configured HTTPS origin and injected GetBoundedWithResponse, its exact-path/version continuation validator, quotaJSON/providerObject/providerArray/disk helpers and existing bounded raw traversal primitives. Do not widen existing quota paths or refactor shared collectors incidentally. No cache, credentials/discovery, POST/checkSkuAvailability, models, deployments, usages or actual network operation is performed by tests.

Before any getter, validate selected UUID scope/labels, normalized collision, nonempty labels and selected expected subscription. Canonicalize UUID/kind case; kind must match existing ASCII providerIdentifier and be nonempty/unpadded. Requested region must already be canonical lowercase and pass diskPhysical. No GLOBAL/logical/display request target. At most256aliases; each key is lowercase, printable/unpadded, not stage/global/logical; removing ASCII spaces must equal its canonical physical value. Existing physical keys cannot rebind. Copy aliases; do not guess display names or permit arbitrary mappings. Unknown wire location strings remain unknown unless this explicit alias rule resolves them. Kind is local filtering/binding, never an undocumented wire filter.

First URL is exactly GET https://management.azure.com/subscriptions/{canonical UUID}/providers/Microsoft.CognitiveServices/skus?api-version=2024-10-01 for the standard configured origin. No body or extra parameters. Continuations allow only exactly one unchanged api-version and optional nonempty $skiptoken, fixed origin/path/subscription/provider, no raw path/userinfo/fragment/opaque URL or unapproved parameter. Reuse canonical query encoding. Max8192raw URL bytes,4096decoded token bytes. Detect cycles before getter; admit at most64pages.

Caller context and wrapped getter cancellation/deadline win over status/parse/partial prefix, returning nil collection and canonical cause. Check context during traversals and immediately before return. A response with non-200 status takes precedence over a nonterminal getter error; status200 with error or nil response is request_failed. Validate1MiB per page/8MiB aggregate before decode. No unbounded/default retries are introduced. The injected bounded getter owns actual HTTP response closure; exercise real HTTP/error/oversize/cancellation closure independently using the existing bounded client seam.

Collection Status is complete/partial/unknown. Unknown means no admitted page; partial means at least one admitted page before operational failure. FailureCode is fixed, sanitized, empty on complete and nonempty on partial. A rejected page contributes no groups, seen regional scopes or counters; prior admitted diagnostics remain owned. Context terminal returns discard all prefixes. Present nextLink:null or a non-string rejects the entire page; missing/empty-string terminates. This differs deliberately from the Compute primary-null exception because this schema specifies a string and its complete example omits nextLink.

## Strict rows, regional repetition and restrictions

Require valid value array, object rows and canonical known property names; reject duplicate JSON keys, folded aliases, trailing content, invalid Unicode/control/noncharacters and existing depth/token violations. Strictly validate identity/name/kind/resourceType and supplied known metadata before filtering. Name follows existing diskName policy (printable ASCII, unpadded, at most512bytes, internal ASCII spaces allowed); kind/resourceType use ASCII identifier policy. Missing/empty/invalid identity rejects, even though primary wire fields are optional. A well-formed foreign resource type/kind is filtered after validation. Optional tier alone can be absent; wrong typed tier rejects. Unknown extensions are bounded but have no invented Compute capability semantics.

Collect groups for only the bound kind. Preserve requested-region facts separately while reducing repeated rows. For each normalized selected-kind/name identity, canonicalize/deduplicate each row's location scopes; repeated locations within ONE row are benign. Different rows with disjoint scopes combine across pages, matching the215row primary example (47identities,14disjoint regional repeats). Any overlap across rows, including identical repeats, rejects the whole page with collision. Missing/empty locations use an unknown scope sentinel; unresolved labels/GLOBAL remain distinct unknown scopes, not physical membership. Other kinds cannot collide with or contaminate the selected kind. Do not blindly reject every repeated name or overwrite a previous regional observation.

Missing or empty locations makes that group unknown; null/wrong-typed arrays reject. Unresolved display/global metadata makes group unknown, even if another row of that identity provides positive membership. Explicit recognized stage labels are nonphysical and do not establish target membership. Other explicit physical regions alone give known absence. These are conservative policies, not claims that the API defines GLOBAL as a physical deployment region or that diskPhysical is a live registry.

Missing restrictions makes group unknown; null/wrong typed restrictions rejects; [] explicitly means no restrictions. Both canonical QuotaId and NotAvailableForSubscription apply. Canonical Location unions values and restrictionInfo.locations; no usable scope means unknown. Canonical Zone requires nonempty restrictionInfo.locations and zones. Positive decimal zone strings have no leading zero; validate all supplied strings before filtering, sort lexically and deduplicate unions. Unknown restriction type/reason/scope means unknown; malformed known fields reject. Unresolved restriction locations cannot manufacture unrestricted evidence. Restriction.values alone does not invent zone scope.

Reduce all selected-kind/name row facts conservatively: unknown dominates; otherwise no explicit target membership yields state2; regional restriction for target yields state1; otherwise matching zonal union yields state3; otherwise state0. State1/2 clear zones. Absent target membership in an unrelated regional row must not overwrite an explicit target row. No account-tier fact establishes model/deployment eligibility, quota or allocation capacity.

## Catalog validation and owned projection

Pure Catalog performs zero getter/credential/network operations. Validate selected scope/expected aliases/binding and every supplied collection field/group before filtering or empty/unknown early returns. Collection subscription/kind/region/type must exactly match normalized expected binding. Validate status/fixed failure enum, complete-without-failure, partial-with-failure, unknown-without-groups. Validate all canonical group keys, Known/state consistency, state0..3 and zones: only state3 has a nonempty strictly lexically sorted unique positive set; Known=false must have state2/no zones.

Required names are validated/charged before normalization, case-fold to lowercase ASCII and reject normalized duplicates. Nil/valid partial/unknown collection yields unknown/empty. Any required Known=false group yields unknown/empty for the whole result, resolved BEFORE building potentially oversized unused output. Unknown unrequired groups do not poison known required groups; malformed unrequired groups still reject. In a complete known catalog, absent required name maps to state2. Empty requirements on a valid complete collection yield complete/empty. Clone maps/zones/warnings and preserve inputs; prove both mutation directions and concurrent call isolation.

Errors return zero catalog plus exactly Cognitive SKU catalog projection rejected [cognitive_sku_projection_invalid] or limit equivalent [cognitive_sku_projection_limit]. Request validation errors return nil collection and exactly Cognitive SKU request rejected [cognitive_sku_request_invalid]. Constructor errors are sanitized. Fixed operational enum: cognitive_sku_request_failed/status_invalid/byte_limit/aggregate_limit/invalid_response/collision/row_limit/work_limit/text_limit/page_limit/continuation_invalid/cycle. Never echo labels, names, kinds, IDs, bodies, URLs or opaque getter/failure strings. Context cause always wins before success/error return.

## Budgets and independent acceptance

At most1000selected subscriptions,256aliases,8192raw rows across the entire collection,8192groups and8192required names; labels at most512bytes (physical values additionally64). Charge all raw JSON object/array work, including filtered rows/unknown metadata, cumulatively toward65536logical entries and16711680decoded text bytes, with the existing per-page JSON depth32/token65536 guard. Account for owned regional-scope/zone state within bounded admitted input. Raw bytes/pages remain separately bounded; do not claim a decoded-text bound is independently reachable through an8MiB wire cap.

Projection separately charges every selection/alias/group/zone/required occurrence and all selected IDs/labels, expected+collection binding/status/failure, aliases, group keys/zones and required-name bytes BEFORE result allocation, at most65536entries/16711680bytes. Bound owned output binding/status/eligibility/fixed warning/map names/zones within the same text allowance. Preserve all old collectors/projections and their budgets. Mandatory supplements isolate exact/one-over scope/alias/row/group/required/label/raw-work/projection-work/input/output text, page/aggregate bytes, pagination/token/URL, late context and response closure. Use independent fixed late-return schedules; a reversal must not recalibrate its own sensitive cancellation endpoint.

[Frozen manifest](../internal/plugins/region/testdata/cognitive-sku-target-contract-v1.json):143literal cases (108collector-backed,35direct),268084bytes, SHA256c12604b2a6ac8b7fed6eb7c66397cf181bba8c0505b512bd5e844972f6aaa81b. Primary schema239116bytes SHA256e6efa93f54d62b025c4afe953b7491ea9c50f497646a6b2545f87e82d03c306a, Gitblobf0919ebcbf91755f10591b931bd80138df13f547; complete GetSkus55486bytes SHA2561c859ac36d080942f817958db9085603cc00b3ce42bbbab5d7a236cf53fc33d6, Gitblob57c0132a67d4958fb40cf20914537b034a10f39f. Both copied unchanged from REST-spec1cdcc541 under the complete MIT license retained in storage-primary-LICENSE.txt and NOTICE.md. No original AZQR capture changes. Expected results were authored from this contract/primary facts without invoking the absent target API. Integrity tests freeze bytes, not143API executions.

After this contract's own full native/source QA, task3 MUST publish an explicit JSON-tagged/unknown-field-rejecting consuming runner and inspect actual Linux and Windows missing-API compilation before production code. Retain the runner/oracle unchanged through implementation. Add compiling named reversal/restoration controls for kind/scope/resource identity, disjoint versus overlapping regions, uncertainty/restriction precedence, context/closure, ownership, limits, sanitization and absent scoring integration. Preserve all409Linux/404Windows identities plus deliberate new controls and all23captures/315chunks/909865bytes. Final outcomes/full logs/readback/artifact proof belong in the live PR index, not an evidence-only commit loop.

Rollback is a new targeted revert after dependent-work review. Revert a future implementation WITH its consuming runner as one coherent unit to the verified contract baseline; never publish the uncompilable runner-only intermediate as an operational target. Rehearse exact-tree rollback/restoration plus region tests and retain commands/results. Update development/handover records between task2 and task3, including exact next action and uncertainty. All three tasks are authorized; no routine input is required for this offline boundary. No Azure/fixture/RBAC/refresh/subagents/protected merge/bypass/release. Independent review/backup/broader operational/live/release obligations remain open.
