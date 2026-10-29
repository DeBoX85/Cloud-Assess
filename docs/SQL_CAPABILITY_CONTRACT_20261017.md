# SQL capability catalog contract, 2026-10-17

Frozen pre-edit task2 contract; not a runtime implementation. Parent is task1 after its own final full QA. All facts link to SQL_CAPABILITY_RECONCILIATION_20261017.md and immutable REST-spec1cdcc541 SQL2021-11-01. No scoring/public integration or deployment assertion.

## Frozen API

```go
type SQLCapabilityRequest struct {
 SubscriptionID, Region, ResourceType string
 LocationAliases map[string]string
}
type SQLCapabilityGroup struct {
 Known bool
 State int
 Version, Edition, Name, Family, SKUName, SKUTier string
 Capacity int64
 Standalone, InstancePool string
}
type SQLCapabilityCollection struct {
 SubscriptionID, Region, ResourceType, Status string
 Groups map[string]SQLCapabilityGroup
 FailureCode string
}
type SQLCapabilityCollector struct { /* private transport/family */ }
func NewSQLDatabaseSKUCollector(endpoint string, getter RESTQuotaGetter) (*SQLCapabilityCollector,error)
func NewSQLManagedInstanceSKUCollector(endpoint string, getter RESTQuotaGetter) (*SQLCapabilityCollector,error)
func(c *SQLCapabilityCollector) Collect(ctx context.Context,selected map[string]string,q SQLCapabilityRequest)(*SQLCapabilityCollection,error)
func(c *SQLCapabilityCollection) Catalog(ctx context.Context,selected map[string]string,expected SQLCapabilityRequest,required []string)(SQLCatalogEvidence,error)
type SQLCatalogEvidence struct {
 SubscriptionID,Region,ResourceType,Status,DeploymentEligibility string
 Values map[string]SQLCapabilityGroup
 Warnings []assessment.AssessmentWarning
}
```

Task3 adds Database constructor/runtime and shared structural reader/catalog only. Task4 publishes the MI consuming runner before its constructor/MI traversal exists, inspects both native missing symbols and then adds bounded MI support. Shared existing collectors, modules/pins/captures and the first runners/oracles remain unchanged. Neither collection nor catalog provides an Evidence method or implicit SKUEvidence conversion.

Fixed Database type microsoft.sql/servers/databases; MI microsoft.sql/managedinstances. Request type must exactly match the constructor family. Canonical UUID/region/type binding is revalidated by pure Catalog. Stable non-racing inputs and nonnil context are preconditions; structural binding does not authenticate caller-built evidence. Selected map1000/labels512, nonempty valid labels/UUIDs, normalized collision rejection and selected request required. Aliases<=256 use unchanged cognitiveAliases provenance rules: lowercase printable/unpadded keys, canonical physical values, removal of ASCII spaces must equal physical value, no stage/global/logical/rebinding. Copy aliases. Region must already be canonical lowercase physical. Bound request ID/type length before normalization.

## Transport, root and completeness

One GET, no body/filter/include parameter:
https://management.azure.com/subscriptions/{canonical UUID}/providers/Microsoft.Sql/locations/{canonical region}/capabilities?api-version=2021-11-01

Use existing NewRESTQuotaCollector origin/injected bounded getter. No discovery/authentication/network/cache in tests and no pager/retry changes. Root is one LocationCapabilities object, never a value array/envelope; a supplied value envelope field is invalid_response. The primary operation is not pageable. Any canonical present nextLink, even null/empty, is continuation_invalid and never followed; folded known keys are invalid_response. Strict quotaJSON guards precede semantic interpretation, preserving depth32/token65536/label512/known continuation8192 behavior. One MiB maximum response,65536raw logical entries,16711680decoded text and8192selected leaf entries/groups. Raw walk charges the entire payload including other root families and unknown extensions before owned result allocation. Reuse vmSKUBudget privately, remapping only fixed work/text failures. No decoded-text overflow claim through a1MiBwire cap.

Caller context and wrapped getter cancellation/deadline win over status/decode/diagnostics and return nil collection plus canonical cause. Status non200 wins over ordinary getter error. Status200 with error/nilresponse is request_failed. Check context during traversals and immediately before every return. Actual injected bounded HTTP seam owns response closure; test status/success/oversize/read-error/cancellation/closure independently with offline transports and selected audience.

Root name: supplied string must be valid bounded JSON; resolve canonical physical name or explicit validated alias to requested region. Known different physical region rejects invalid_response. Missing/unresolved/nonphysical root name makes the body incomplete. Do not infer East US without the explicit alias. Root status participates in every leaf. Validate supplied known metadata before returning incomplete. A missing selected root array or selected child array is incomplete; explicit [] is a healthy empty list; null/wrong types reject. Unknown returned fields stay opaque but bounded. Arrays and names in other branches remain structurally validated against the pinned wire shapes; they do not create selected-family entries.

## Known wire shapes and status semantics

Validate canonical recognized field names, supplied primitive/container wire types recursively through all28pinned capability definitions plus common Sku. Unknown fields are allowed under the raw bounds. This is bounded selected-schema shape validation, not a general JSON Schema engine or certification of every semantic enum/range. Sku.name is required by common Sku when a sku object is supplied. int32 fields must parse as integer within signed32bit range; known floating fields must be finite. Known optional strings/bools/arrays/objects reject supplied null/wrong types. Selected coordinate names additionally follow nonempty ASCII providerIdentifier, <=512bytes; required selected identities cannot be anonymous. Validate side metadata (e.g. zoneRedundant, computeModel and nested license/storage properties) before statuses or filtering. Optional status absent is unknown, supplied status must be a string; only exact canonical Available/Default/Disabled are recognized, Visible/lowercase/future/Unicode strings stay unknown. Unknown dominates a Disabled ancestor, otherwise any Disabled yields known State2, all Available/Default yield known State0. No State1/3 or blocked-zone inference.

Database path: root -> server version -> edition -> service objective. Canonical key version/edition/objective, all components lowercase ASCII; retain Version/Edition/Name and authoritative sku.name/tier/family/capacity. A supplied sku requires nonempty valid name; missing whole sku, absent tier or capacity makes that leaf unknown (zero unknown fields, no fallback). Tier must match edition when supplied. Family is optional for Database and is not inferred. Capacity must be positive signed32bit when supplied. Other known sku size metadata is validated but not used for placement. Standalone/InstancePool are both not_applicable.

MI path: root -> MI version -> edition -> family -> vcore. Key version/edition/family/vcore. Use family.sku string directly, never hand-coded GP/BC/NGP or raw-edition fallback. Missing sku/value makes leaf unknown; supplied sku is valid ASCII identifier; supplied value positive signed32bit. Family is the exact canonical family node name, SKUTier is edition. Name is vcore node name, which is not guessed from numeric value. Optional standaloneSupported/instancePoolSupported retain supported/unsupported/unknown strings, do not silently select a mode or establish deployment eligibility. All governing statuses include family/vcore. Repeated SKU names across different coordinates are valid. Duplicate normalized sibling version/edition/family/leaf identities are collision, including identical duplicates and disjoint descendants under a duplicate parent.

Collection status complete/unknown only, no partial root. A fully validated admitted body is complete even if some leaves are unknown. Incomplete metadata/list knowledge or any operational/shape/limit/collision failure returns unknown, empty Groups and a fixed failure code; rejected body contributes no diagnostic prefix. No JSON null collection success. Fixed failure suffixes request_failed,status_invalid,byte_limit,invalid_response,row_limit,work_limit,text_limit,continuation_invalid,collision,incomplete under sql_capability_. Complete has no failure; unknown requires a fixed nonempty failure. Never echo URLs, identities, provider errors or returned reason text.

## Pure catalog, ownership and errors

Always validate selected scope/aliases/expected binding and ALL supplied collection fields/groups, including unrequired or empty requirements, before unknown early returns. Group map keys must equal canonical coordinate fields (3Database/4MIcomponents, max2051bytes); metadata/flags/state and canonical fields must be consistent. Known=false implies State2; metadata may be absent/zero only where the collector's declared missing-field policy allows it. Known=true requires valid authoritative SKUName/tier/positive capacity. Database tier equals edition, modes not_applicable; MI family required, tier equals edition, modes supported/unsupported/unknown. All direct strings/fields must be bounded and valid; no unrecognized failure enum, noncanonical binding or inconsistent status/failure/groups.

Required full coordinates are valid/canonicalized ASCII, reject normalized duplicates; no lookup by name-only SKU or inferred inventory tier/capacity. Nil/valid unknown collection, ANY missing required coordinate or required Known=false makes whole catalog unknown/empty. Unknown unrequired entries do not poison known required entries. Valid complete with no required coordinates is complete/empty. Copy Values and Warnings; scalar entry values own no mutable slices/pointers. Never mutate selected/aliases/required/collection; prove repeated and concurrent calls, both mutation directions and isolated maps/warnings.

Every nonerror catalog result owns allocated Values/Warnings, canonical binding, DeploymentEligibility=not_assessed and exactly one warning:
code sql_capability_catalog_only
message SQL capability catalog evidence does not establish deployment eligibility, quota or allocation capacity.

Constructor: invalid SQL capability origin or getter.
Request rejection: SQL capability request rejected [sql_capability_request_invalid], nil collection.
Catalog rejection: zero catalog and SQL capability catalog projection rejected [sql_capability_projection_invalid] or [sql_capability_projection_limit]. Context cause always wins before any success/error return.

## Projection accounting and independent reachability

Before output allocation, charge every selected/alias/group/required occurrence and all supplied selected IDs/labels, expected/collection binding/status/failure, alias keys/values, group coordinate keys plus eight string metadata fields, and required coordinates. Bound groups/required8192, aliases256, selected1000, component/metadata labels512, full coordinate2051 and total16711680text bytes. Raw collector work has independent65536exact/one-over tests. Projection occurrence count is provably <=17640 from admitted public counts, so its generic65536work guard is not independently reachable; test all public counts exact/+1 and maximal joint admission rather than claiming a false independent overflow.

Output separately bounds binding/status/eligibility/warning/key/metadata text within16711680. Under valid counts/component limits and required-key input charging, this output cap is implied by the input cap: small group counts cannot approach16MiB, large counts charge enough extra required-key bytes to exceed fixed warning overhead. The independent bound calculation gives Database maximum16693718 and MI16689011under valid input caps, below16711680. Retain the inequality/maximum calculation and tests; do not invent a reachable output-overrun fixture. Input text overflow remains independently reachable with valid long direct coordinates/metadata and empty requirements. Freeze literal byte-overhead accounting, exact/+1 cases and control-sensitive tests before code. No shared production primitive may be weakened to pass.

## Acceptance and rollback

Separate unchanged Database and MI manifests freeze full expected collection/catalog/error, exact one-GET URL/count/limit, input preservation and ownership; complete primary example expectations retain all26/19coordinate entries, not target-produced output. Integrity tests freeze manifest/schema/example/common/license bytes. Cases have NOT run before API creation. Each implementation first publishes a strict explicit JSON-tagged unknown-field/trailing-rejecting runner and inspects real Linux+Windows missing symbols before production. Fixed late-return context schedules must not recalibrate themselves under reversal. Supplemental limits/raw JSON/HTTP/ownership/full coordinates/ancestor statuses/malformed unrequired metadata and no-scoring tests must be independent. Compiling reversals require named assertion failures and exact byte restoration/passing baselines; preserve every455/450predecessor identity plus literal new controls through additive CI inventory and all23captures.

Rehearse atomic Database implementation+runner rollback to verified contract tree and exact restoration; MI implementation+runner rollback to verified Database tree and restoration, running region tests both sides. Do not publish an uncompilable runner-only intermediate as an operational rollback target. Full final-head8workers/2gates/source/full logs, remote tree/parent/identity/full-byte readback and durable evidence are mandatory per task. Final live PR outcomes supersede publication-time pending text; update handover BETWEEN each task. No Azure/refresh/subagents/protected merge/release. Independent review/backup/broader operational/live/security/license/Gate004/release/advisory obligations remain open.

## Frozen oracle and immutable primary identities

- `sql-database-capability-contract-v1.json`: 153 cases (124 collector/29 direct), 362752 bytes; SHA256 `67b664515b0407d57a037d401bd830662db49b4b177c2a67140701a34587a5fd`.
- `sql-managed-instance-capability-contract-v1.json`: 163 cases (134 collector/29 direct), 399973 bytes; SHA256 `4035965123c5c72c93c099c7accd6c7be8d2d5a84a7041a8359f0ec1ecd720a7`.

316 independent future cases; no target API exists or was invoked. Pre-freeze review corrected literal-null default handling and explicitly declared both required coordinates in the same-SKU case. No frozen expectation was rewritten to fit runtime behavior.

- `sql-primary-2021-11-01-schema.json`: 38062 bytes, SHA256 `8246593e56c6d2708186a3603f9d92ee64ef2cc292931c7db09753471e506040`, Git blob `d6c7b7fb88fa594932ba62c5523c343ab85108dd`.
- `sql-primary-2021-11-01-example.json`: 119773 bytes, SHA256 `f4b992d87b9e6e8bae6a7693092b4028dfa60bc7a79f132a8b658acebc9b57dd`, Git blob `9815b0a083dd402d2b7b17e29bb9b09df392dfce`.
- `sql-primary-2021-11-01-common.json`: 5904 bytes, SHA256 `709781a78cdfb4ac27edcd4bdc7709e4667b3e33eaf4f91ae1d2c21a42b9a1c3`, Git blob `f7decac418742c0cd8fb0c6b44d1330a5b6cb46d`.

All three primary files are unchanged from Azure/azure-rest-api-specs commit `1cdcc5414764d20e8815554c8b89708c10747c18`, specification/sql/resource-manager/Microsoft.Sql/SQL/stable/2021-11-01. Complete MIT notice is retained in `storage-primary-LICENSE.txt`, 1076 bytes, SHA256 `b37363bf2295aa0acde60ce78a60034a71a0955075fc01aa09cc23eadff44436`.
