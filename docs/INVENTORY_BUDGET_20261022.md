# Required inventory budget design, 22 October 2026

Task 8 remains provisional. This design supplies a concrete input ledger and four independent exact/excess admission recipes before any new API implementation. It does not modify the original50cases, collectors, source pins, public schema or existing accounting. Output accounting and the final consumer contract remain separate freeze prerequisites.

## Meaning of a budget

Use separate input and owned-output counters for each BuildRequiredInventory or AdmitRequiredEvidence call: at most65536logical entries and16711680UTF8 bytes per counter. Never reset a counter for another family, Cognitive kind, collection, requirement, nested map or helper invocation. Existing stricter projector limits still apply. These are logical data bounds, not a hard RSS limit, timeout, instruction count or proof of Azure completeness.

An entry is one dynamic map member, list member or explicitly listed row below. Fixed scalar fields belong to their row and do not each create another entry. Text counts every listed string occurrence using UTF8 byte length, including repeated values and empty strings, before deduplication/filtering. Fixed Go field names are not strings supplied by a caller and do not charge this typed-input ledger. Unknown JSON properties already belong to collector admission, not these typed collections. Tags/SLA are ignored, never traversed/copied, and cannot supply coordinates; do not claim this boundary limits the caller's original object graph.

Charge counts before allocating sorted-key scratch, owned records or result maps. For an oversized bounded container, reject from its length before copying its keys. Keep context checks in loops and after sorting/projection. Bounded validation may traverse data more than once; entry accounting is not a count of loop iterations. Use a one-subscription selected map when invoking existing scoped projectors, preventing repeated1000-scope validation for every supplied collection. Do not change the projectors' own validation/accounting.

## Build input ledger

| Data | Entries | Text occurrences |
| --- | --- | --- |
| Selected scopes | One per map member | Raw ID and display name |
| Inventory declarations | One per row | Raw ID and status |
| Raw resources | One per row | ID, SubscriptionID, ResourceGroup, Name, Type, Location, Kind, SKUName, SKUTier and SKUFamily |

Capacity is a bounded numeric field in its resource row. All retained strings, including currently unused tier/family labels, receive scalar validation and text charges. This entry counter cannot reach65536 under1000selected scopes,1000declarations and8192resources: its maximum is10192. Record that domination honestly; do not weaken outer caps or invent a public exact65536Build-input test. Its text bound still requires separate reachable fixtures. Normalized keys and owned copies charge the output ledger when retained, not an invented second raw-input charge.

## Admission input ledger

Only the requested subscription/source slice of the immutable validated plan is consumed. Exported mutable diagnostics are not consulted. Count all requirements in that slice, every supplied binding and every nested supplied payload row, even when its SKU is unrequired or its status will project to unknown. Other selected scopes do not become evidence for this request.

| Data | Entries | Text occurrences |
| --- | --- | --- |
| Scoped context | One | Canonical selected ID, selected display name and inventory status |
| Request | One | SubscriptionID, SourceRegion and TargetRegion |
| Each scoped requirement | One plus each SKU, each Storage kind-map member, each required-kind list member and each ResourceID | SubscriptionID, SourceRegion, ResourceType, Cognitive-only Kind, all SKU strings, Storage kind-map keys/kind values and ResourceIDs |
| Each collection binding | One | SubscriptionID, SourceRegion, TargetRegion, ResourceType and Cognitive-only Kind |
| Each typed payload | One | Every actual string metadata field, including Status/FailureCode and inner Kind/ResourceType where present |
| Disk values | One per value plus each zone | Every map key and zone |
| VM/Cognitive groups | One per group plus each zone | Every group key and zone |
| Storage groups | One per SKU map member plus each kind member | SKU key and every kind key |
| Storage UnqualifiedNames | One per member, even a false value which will be rejected | Every key |
| SQL groups | One per group | Coordinate key and Version, Edition, Name, Family, SKUName, SKUTier, Standalone and InstancePool |

All fields represented by a row are validated under the correct existing family contract. SQL coordinate keys retain their existing2051-byte bound; ordinary retained labels retain512bytes and stricter metadata rules. Incomplete/unknown results do not exempt supplied data from accounting. Do not flatten Cognitive kinds, collapse Storage requirements, infer SQL coordinates or admit VMSS scoring. A global-budget failure returns nil with required_evidence_limit, not a partial prefix or truncated success. The raw original proposal is retained; this defines the accounting candidate, not an implementation result.

## Independent exact/excess recipes

All recipes use one selected synthetic subscription aaaaaaaa-1111-4111-8111-111111111111 named One, complete inventory, eastus source and westeurope target. There is one Cognitive resource/requirement for each kind c0 through cN, all with required SKU s0 and canonical IDs ending `/microsoft.cognitiveservices/accounts/cN`. Each binding has exactly one complete kind-bound Cognitive collection, Known=true/state0 groups and empty zones. Unrequired group names are synthetic lowercase ASCII. They must still be fully accounted and validated.

The required s0 is absent from those complete catalogs. Therefore the successful result has complete admission, unsupported/empty scoring, one catalog per exact kind with s0 state2 and deployment eligibility not_assessed. This is a negative catalog result, not a deployment claim. Excess variants fail atomically with required_evidence_limit and nil result.

| Recipe | Fixed input construction | Independent expectation |
| --- | --- | --- |
| Entry exact | Eight kinds. Group counts8192 seven times, then8150; keys g00000 onward per collection |42context/request/requirement/binding/payload entries +65494groups =65536 |
| Entry first excess | Change only last group count to8151 |65537entries; each collection still meets8192local cap |
| Text exact | Four kinds,8192groups each. Keys normally510bytes: g plus five decimal digits, then x padding. First three global keys have length6 and fourth391 | Fixed non-key text1631bytes + key text16710049bytes =16711680;32790entries |
| Text first excess | Change only fourth key length391 to392 |16711681bytes; each scalar remains within512bytes and entries remain32790 |

The entry recipe's non-key text is3163bytes and group keys total392964bytes, well below the text bound. The text recipe's context47 + request52 + four times(requirement201 + binding90 + inner metadata92) gives the literal1631-byte overhead. Four times8192times510 is16711680; shortening three keys by504 and one by119 removes exactly1631. This makes the global text boundary independently reachable without violating a local group, scalar or global entry limit.

[Four literal recipes](evidence/20261022/inventory-budget/budget-recipes.json) and [the independent checker](evidence/20261022/inventory-budget/budget-recipe-check.py.txt) retain those numbers and the exact string construction. The checker constructs/hash-checks the synthetic keys and verifies arithmetic, uniqueness, scalar bounds and expected dispositions. It imports no product code and ran before any new consumer/API. It is specification arithmetic evidence, not proof that an absent implementation enforces the budget.

## Validation and owned-output work still open

Before freeze, map the logical owned-output ledger to every public diagnostic and private immutable snapshot/accessor copy. Count repeated retained strings/entries separately rather than treating two owned views as free. Include statuses, fixed reasons/warnings, catalog coordinate metadata and zones; counters must cover all families together. Define alias/duplicate/invalid/limit precedence under the canonical traversal and terminal context rules. Prepare independent output, Build-text, cancellation, ownership and forged-plan/accessor cases. A single scalar or existing tighter projector cap must not mask a claimed whole-call boundary.

The input recipes above resolve a concrete part of the budget prerequisite. Task8is not frozen and task9production remains unstarted. Complete the remaining applicable-document reconciliation and freeze/verify the whole contract before publishing the strict consumer and inspecting BOTH native missing APIs. Preserve all existing controls and original expectations.
