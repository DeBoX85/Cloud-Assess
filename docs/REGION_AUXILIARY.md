# Region quota and capacity reservation pure projections

Status: IN PROGRESS; unaccepted candidate based on PR105 merge7e4ca5c955abb7f882591a44d8bb2f74c93ae243. Public region-selection unavailable.

## Contract and source authority

Pinned AZQR8e4f0577f3615e6c9014c031bcad079f235369cc internal/scanners/plugins/region/output/output.go BuildQuotaSheet/BuildCRGSheetFromRows define labels, descriptions, received ordering, formatting, flag precedence and nil empty branches. [Source characterization](REGION_SOURCE_CHARACTERIZATION.md) calls unchanged functions with literal synthetic inputs. Retain source-aux-inputs.json/source-aux-outputs.json exact bytes covering all fourteen branches; implement only quota/reservations here.

Before: primary scoring only. After: ProjectQuota and ProjectReservations accept stable already-decoded selected inputs, return optional owned canonical tables, no HTTP/discovery/cache/CLI/registry change. Empty success returns nil, matching source; failure returns zero rows/source headers and failed health with fixed error text. Canonical metadata is the owning region Metadata(), IDs quota/reservations, schema1.0, replacing source helper default/absent metadata deliberately.

Quota preserves Current/Limit/Available/HeadroomPct and flags without recalculating arithmetic or deriving flags from rounded text. Over-limit wins, then near-limit, then OK. Reservations preserve ten supplied source cells, signed availability and status. Service decoders/arithmetic consistency remain future contracts. A required out-of-band SubscriptionID UUID must be selected; display Subscription must match exactly. Attach canonical lowercase UUID for later privacy processing, preserving cells. Same names across distinct UUIDs remain valid.

## Input corrections and bounds

Reject foreign/malformed/case-alias duplicate scope IDs; duplicate quota tuple(UUID,region,type,resource) or reservation tuple(UUID,region,resource group,group,name), using structured keys; control/replacement/invalid UTF8 labels; invalid ASCII regional IDs; nonfinite/excess numbers; invalid reservation width/count/status; excessive rows/text. These deliberate pure-input corrections are not live/parser parity.

Maximum8192 rows per call,1000 subscriptions,512 UTF8 bytes per label, counts within1e12 (Current/Limit/Reserved/Allocated nonnegative; Available signed), finite HeadroomPct within plus/minus1e12. Reservation numeric strings must be canonical base10. Text budget before row allocation counts scope names and all supplied labels/cells, capped at MaxPluginTextBytes minus64KiB for canonical metadata/formatting. Final canonical validation remains mandatory. No truncation or partial successful rows. Empty input still checks scope/context. Caller input must stay stable during a call; no input map/slice escapes.

## Acceptance and boundaries

Compare every quota/reservation source cell, columns, description/order and nil empties; verify hashes, canonical metadata/health/correlation, precedence independent of rounding, negatives, same-name scopes, foreign/malformed/duplicates, finite/count/row/aggregate budgets, owned columns/cells, cancellation and repeated/concurrent isolation. Both required hosts run compiling selected-scope/work-limit controls with named failures/restored tests, alongside all existing guards and characterization. Full native final-head Linux/Windows jobs/logs, semantic review, exact tree/parent/identity/bytes/base/preview/rules and expected-head protected merge remain required.

No fresh local Go QA, Azure/laptop, quota/SKU/cost/latency fidelity, public/all-sheet integration, load/fresh-OS, Gate004/release proof. Pins/notices unchanged. Rollback: protected reviewed revert after dependency review. Next separate pure sheets: service availability/cost, inventory/calculations, then bounded service adapters and public execution.

SHA256 inputs f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28; outputs efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06. [Microsoft capacity documentation](https://learn.microsoft.com/en-us/dotnet/api/azure.resourcemanager.compute.models.capacityreservationutilization.currentcapacity?view=azure-dotnet) distinguishes billed reserved capacity from allocated resources; it does not justify modifying the captured helper cells or adding modern API fields here.
