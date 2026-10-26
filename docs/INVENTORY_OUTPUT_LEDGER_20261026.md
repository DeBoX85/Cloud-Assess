# Inventory ownership and budget preparation, 26 October 2026

Task8 preparation only. No future API exists or was executed. This supplements the [input ledger](INVENTORY_BUDGET_20261022.md) and preserves the original three provisional files and50cases. It specifies logical retained data, not Go heap allocation, cumulative caller memory, execution time or authenticated Azure provenance. Final task8freeze still requires the remaining validation/expectation/documentary work and own QA.

## Ownership accounting decision

Keep independent whole-call input/output ceilings of65536entries and16711680UTF8bytes. A dynamic map/list member is one entry; fixed scalar fields charge text but do not each create entries. The explicitly listed request row is one entry. Root structs, nil/empty containers and fixed boolean/numeric fields create no extra entries. Count every retained string occurrence, including repeated keys, fixed statuses, warnings and not_assessed labels. Do not count JSON property names: this is a typed-data ledger, not serialized-byte accounting.

Build retains two separately owned logical views: public diagnostics and private immutable snapshots of scopes, resources and requirements. Plan.Status is one public root string. A scope view is a canonical-ID map member containing Name and Status; do not add an uncharged duplicate ID field. Each resource view retains the ten raw strings in the input ledger and separate canonical SubscriptionID, SourceRegion, ResourceType and Kind. Nonphysical locations retain their raw Location but canonical SourceRegion is empty and Physical=false; they create no regional requirement. Canonical Kind is retained diagnostically even for unsupported types; it does not invent an adapter or eligibility.

| Build output, per view | Entries | Text |
| --- | --- | --- |
| Scope member |1| Canonical ID, display name, inventory status |
| Resource row |1| Ten raw fields plus four canonical fields |
| Requirement row |1| SubscriptionID, SourceRegion, ResourceType, Cognitive-only Kind |
| Required SKU list | One per member | Every canonical SKU |
| Storage kind map/list | One per map member plus each kind list member | SKU map key and every kind, including missing-kind marker |
| ResourceID list | One per member | Every canonical resource ID |

Thus Build output is two times the sum of those rows plus the single Plan.Status string. Fixed requirement flags are booleans, not a retained string array. External diagnostic mutation cannot change either private view. Neither raw Tags nor SLA is traversed or copied. Raw-input strings and normalized-output strings are separate accounting occurrences even when Go might share immutable string backing storage.

Admission retains public Request/SKUs and a separate immutable private Request/SKUs. Reserve one further full Request/SKUs copy for a single Evidence() result before constructor success. This third logical view means an accessor can return a deep copy without adding an unbounded shared counter or changing its no-error signature. Each accessor invocation is independently bounded by that reservation. Callers retaining arbitrarily many results are outside this per-call guarantee. There is no lifetime memory quota or claim that concurrent caller writes are safe.

Catalogs remain owned public diagnostics. They are not scoring input to Evidence(); mutable catalog diagnostics cannot rebind or manufacture scoring. A future accessor that exports another private catalog snapshot would require a new explicit budget, rather than treating another copy as free. The owned coordinator must not treat mutable diagnostics as authenticated eligibility.

| Admission output | Copies | Entries and text |
| --- | --- | --- |
| Request |3| One row each; SubscriptionID/SourceRegion/TargetRegion each time |
| SKUs map |3| One member per resource type; type and status text |
| SKU values/zones |3| One per value and zone; keys and every zone occurrence |
| ScopeStatus and Status |1| Both root strings |
| Catalog binding |1| One list member; all five binding strings plus Status/DeploymentEligibility |
| Typed catalog payload |1| One row; every actual metadata string |
| Catalog values |1| One per map member and each zone; SQL key plus all eight SQL group strings |
| Catalog warnings |1| One per warning; Code and Message |

All families share the same counters; no reset for another kind, type or catalog. Missing catalogs still charge their typed empty diagnostic metadata/warnings. Unsupported/unknown scoring maps still charge type/status. SQL values remain actual collected coordinates, never guessed matches.

Charge planned global retention and copies before retaining/cloning them. Existing projectors may allocate their own bounded temporary results; those existing methods are unchanged. Do not claim this ledger caps every transient allocation. Check cancellation before/after projection and before adopting its result. A later rejection returns nil and no partial global result; it does not retroactively erase temporary allocations.

## Independent boundary recipes

[Recipes](evidence/20261026/inventory-output/budget-recipes.json) and the [independent checker](evidence/20261026/inventory-output/check-budget-recipes.py.txt) define literal arithmetic before implementation. The checker imports no product code and does not regenerate the original50cases. Recipe scalars remain within512bytes; cardinalities remain within existing outer limits. All IDs are synthetic and no request is sent.

Build text fixtures use subscription aaaaaaaa-1111-4111-8111-111111111111, display One and complete declaration. Resources have IDs ending microsoft.example/widgets/r0000 through rNNNN, group rg, global Location and unsupported type microsoft.example/widgets. Name, Kind, SKUName, SKUTier and SKUFamily normally contain512ASCIIletters each. Canonical source is empty, canonical Kind repeats the512-byte raw kind. No regional requirement exists. Generic display Name is not inferred from ID. One raw resource totals2740bytes, and one output resource view3313bytes; the output scope view is47bytes and Plan.Status8bytes.

| Recipe | Literal construction | Expected accounting |
| --- | --- | --- |
| Build input text exact |6500rows; first2149Names length1, next248, remaining512|16711680input bytes; output limit also exceeded |
| Build input first excess | Change only that248to249 |16711681input bytes; required_inventory_limit |
| Build output text exact |2600rows; first504Names length1, next45, remaining512|16711680output bytes;6866072input bytes |
| Build output first reachable excess | Change only45to46 |16711682output bytes; required_inventory_limit |
| Build output entries exact |8191VM rows, one unique physical source r00000... per row, one SKU g00000... and one ID per requirement; four complete selected scopes |65536output entries;8199input entries |
| Build output first reachable excess | Add a fifth complete empty selected scope |65538output entries; required_inventory_limit |

Build input count remains dominated by10192outer entries. Under the prescribed duplicate ownership, Build input text is also dominated for successful results: raw retained fields alone must be represented twice in output. The input-exact fixture therefore does not mean successful Build admission. Public identical limit errors cannot independently establish which internal counter rejected it. Do not claim isolated causal Build-input enforcement from these recipes. For these fixed complete-status Build fixtures output counts/text move in steps of2; no plus-one API fixture exists. The65536/16711680ceilings stay unchanged; first reachable excess is not a relaxed ceiling.

Admission output recipes use8191VM requirements for entry tests, or8190for text tests, plus one SQL Database resource. All VM resources share source eastus and use distinct six-byte g00000... SKUs. The SQL resource has safe generic SKU x and unresolved coordinates. Complete VM groups are Known=true/state3 with13650total sorted numeric zones: the first5459groups have zones1and2 in the8191recipe; for8190groups the first5460have two. All other groups have one zone. SQL catalog input has complete known group v/e/a, fields Version=v, Edition=e, Name=a, SKUName=a, SKUTier=e, Capacity=1, Standalone=InstancePool=not_applicable. SQL scoring remains unsupported/empty and deployment not_assessed.

| Recipe | Literal construction | Expected accounting |
| --- | --- | --- |
| Admission output entries exact |8191VM keys,13650one-character zones, one SQL group |65536output entries;38234input entries |
| Admission output entries first excess | Add complete SQL group v/e/b with corresponding Name/SKUName=b |65537output entries;38235input entries |
| Admission output text exact |8190six-byte VM keys; zones normally400digits; first546zone occurrences length512, next404, remaining400; SQL Family=ff |16711680output bytes;65533output entries |
| Admission output text first excess | Change only SQL Family ff to fff |16711681output bytes;65533output entries |

The text fixture has790fixed output bytes before Family, and three copies of VM key/zone strings. Zones total5521156bytes; output is790+2+3*(49140+5521156)=16711680. Increasing Family by one changes only the singly retained SQL diagnostic occurrence. Input and each existing projector remain below their respective limits; all eight existing projector calls accepted these four fixtures on isolated frozen product185 with Go1.26.10Linux, readonly cached modules and a clean environment. The [probe](evidence/20261026/inventory-output/projector-probe.go.txt) and [execution record](evidence/20261026/inventory-output/projector-proof.json) preserve exact construction and results. These are expectation recipes, not future API execution or deployment evidence.

## Remaining freeze obligations

Existing projector acceptance is verified. Finalize exact validation/error order and additional literal/ownership/context expectations, reconcile applicable documents, and run own native/source QA before task8freeze. Do not change the original50cases, source pins, legacy behavior or retained controls. Task9still first publishes its strict consumer and inspects BOTH actual native missing APIs before production.
