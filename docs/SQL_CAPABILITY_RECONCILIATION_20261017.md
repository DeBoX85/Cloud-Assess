# SQL capability reconciliation, 2026-10-17

Current authority: [CURRENT_STATE](CURRENT_STATE.md), active live PR and [four-task workstream](SQL_CAPABILITY_WORKSTREAM_20261017.md). Task1 is analysis only, based on verified draft PR157/13097058b5fab00fe083bb06a62f27fb3d2107b0. No production/scoring/public-schema, source capture, SDK/module/action/pin, protected merge, release or Azure operation. This reconciles evidence before freezing a new contract; proposed policies below are not source equivalence claims.

## Primary evidence and immutable source

Original AZQR commit8e4f0577f3615e6c9014c031bcad079f235369cc, sql.go blob1cb1e6430f11476e66e8bbb13bc0cc2ad5e0a43f. All67retained SQL source cases (35Database/32MI), their source harness and original23capture files remain unchanged.

Microsoft REST specifications commit1cdcc5414764d20e8815554c8b89708c10747c18, path specification/sql/resource-manager/Microsoft.Sql/SQL/stable/2021-11-01:

| File | Git blob | UTF-8 bytes | SHA256 |
| --- | --- | --- | --- |
| LocationCapabilities.json | d6c7b7fb88fa594932ba62c5523c343ab85108dd | 38062 | 8246593e56c6d2708186a3603f9d92ee64ef2cc292931c7db09753471e506040 |
| examples/LocationCapabilityListByLocation.json | 9815b0a083dd402d2b7b17e29bb9b09df392dfce | 119773 | f4b992d87b9e6e8bae6a7693092b4028dfa60bc7a79f132a8b658acebc9b57dd |
| common/v1/types.json | f7decac418742c0cd8fb0c6b44d1330a5b6cb46d | 5904 | 709781a78cdfb4ac27edcd4bdc7709e4667b3e33eaf4f91ae1d2c21a42b9a1c3 |

Complete files and independent primary analysis are retained with task evidence. The MIT license already retained in internal/plugins/region/testdata/storage-primary-LICENSE.txt belongs to this same REST-spec commit. Task2 will copy immutable fixtures with explicit attribution before runtime tests. These historical examples are not today's live inventory or a service-support guarantee.

Official [2021-11-01 API documentation](https://learn.microsoft.com/en-us/rest/api/sql/capabilities/list-by-location?view=rest-sql-2021-11-01) was checked online. GET /subscriptions/{id}/providers/Microsoft.Sql/locations/{location}/capabilities?api-version=2021-11-01 returns one nested LocationCapabilities object, not a value list. Optional include is documented with a bounded enum; the source omits it. The operation has no x-ms-pageable or nextLink in the pinned schema. Keep that same documented version; no version/SDK refresh is needed. Future target must not follow an undocumented continuation.

## Confirmed source behavior versus required correction

| Concern | Original observed/code behavior | Primary fact / target implication |
| --- | --- | --- |
| Database identity | Maps lowercase sku.name; later observations overwrite earlier ones | Complete example has26objectives but only HS_Gen4(12) and HS_Gen5(14); retain version/edition/objective plus SKU metadata and capacity |
| Database status | Only case-insensitive leaf Disabled becomes unavailable; ancestors ignored and other strings default available | Root, version, edition and objective each have status; all governing levels must participate |
| MI identity | Builds tier-code + underscore + family, including a future-tier fallback | Family already has an explicit sku string; example Hyperscale/Gen4 supplies HS_Gen4, while source fallback constructs hyperscale_gen4 |
| MI sizing | Family presence becomes available; vcore capabilities ignored | Example has19vcore leaves under4explicit SKU families, with standaloneSupported/instancePoolSupported and per-leaf status |
| Status vocabulary | Unrecognized/missing values can become available | Schema enum is Visible, Available, Default, Disabled; visibility alone is not documented deployment eligibility |
| Region | URL embeds caller region; response name not bound | Primary uses display name East US; request must remain canonical physical and any supplied display alias must be explicit/validated |
| Inventory | Rows retain tier/capacity; current aggregate SKU index is name-only | Aggregate loses the capability coordinates needed to bind an exact SQL shape; no incidental inventory change or implicit scoring conversion |

The complete example contains Database26distinct objective coordinates and MI19distinct version/edition/family/vcore coordinates. Repeated SKU names across sizes are valid, not duplicate resource collisions. The full independent analysis retains all45leaf records with ancestor statuses and supplied dimensions. Do not infer live Hyperscale MI availability from this historical example.

## Proposed conservative target policies for task2

Use explicit full coordinates: Database version/edition/objective, MI version/edition/family/vcore, preserving authoritative SKU name and capacity separately. Bound, canonicalize and validate each coordinate component; reject overlapping normalized full coordinates, but permit repeated SKU names at different coordinates. Missing MI sku must not trigger an undocumented name reconstruction. Retain optional standalone/pool information as diagnostics, not deployment approval.

Recognize only canonical Available/Default as positive catalog status and Disabled as negative. Visible, absent or future status remains unknown. Combine all relevant ancestors conservatively, with unknown dominating; malformed supplied metadata rejects before filtering. These are deliberate fail-closed policies, not a claim that the schema defines a placement algorithm. A missing requested coordinate stays unknown; do not manufacture unavailable evidence from incomplete catalog knowledge.

Use a distinct internal catalog-only result with not_assessed deployment eligibility and a fixed limitation warning. No Evidence conversion into existing name-only SKUEvidence, registry/report/coordinator/CLI integration or SQL scoring. Explicit coordinate knowledge is necessary; quota, allocation, license, pool mode, storage and zone requirements remain separate. Raw/source names remain available in immutable characterization evidence.

Proposed transport: one bounded GET through existing configured origin and RESTQuotaGetter, selected UUID/canonical region/resource family validated before credentials, strict JSON and cumulative work/text limits, owned results, terminal context priority and exact URL/no-body/no-filter checks. No cache/retry policy changes; reject a present undocumented continuation. Final concrete errors/budgets/schema/admission/ownership and supplemental exact/one-over tests belong in task2's pre-edit contract before API implementation.

Local preparation passed the focused region preflight (uncached package tests/vet plus workflow/gate/handover/whitespace guards),1450authored relative links and independent complete primary byte/Gitblob analysis. Existing tests are not execution of a new SQL API; no such API exists at this checkpoint. Own hosted/source verification is still required.

## Transition and remaining obligations

Task1 requires its own full final-head8workers/bothgates/source review, all455Linux/450Windows control identities and23captures/315chunks/909865bytes, full readback and retained proof. Then update development/handover BETWEEN tasks before task2. Tasks3/4 must publish consuming runners and inspect actual Linux/Windows missing-API preparation before production edits. Each implementation needs independent boundaries/reversals/rollback/full own QA.

No critical user input is currently needed. Preserve all prior draft heads and protecteda6fee427. Independent review, independent off-repository backup/retention, installed/broader operational rollback, live permissions/scale/fresh-OS/security/license/Gate004/release and existing module-only advisory remain open. Source characterization7/7 and accepted target adapters1/7 remain unchanged; full parity is not established.
