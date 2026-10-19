# SQL diagnostic report contract, 2026-10-18

Frozen task2 contract for the new offline integration batch, based on reconciliation93500331 and verified SQL producta3813439. No new report API exists at freeze. This adds no public execution, scoring, eligibility inference or schema/pin/oracle refresh.

## APIs and trust boundary

```go
func ProjectSQLCapabilityTable(ctx context.Context, selected map[string]string, q SQLCapabilityRequest, collection *SQLCapabilityCollection, required []string) (assessment.PluginTable, error)
func CollectSQLCapabilityTable(ctx context.Context, collector *SQLCapabilityCollector, selected map[string]string, q SQLCapabilityRequest, required []string) (assessment.PluginTable, error)
```

Inputs remain stable with nonnil context. Project calls the unchanged pure Catalog to validate ALL input/binding/coordinates/metadata before any unknown return. Caller-built collections are structurally checked, not authenticated service provenance. Collect first validates scope/request/aliases/required via nil-collection Catalog BEFORE any transport call; then calls the existing concrete family-bound Collect and Project. Collect snapshots the validated selected map, aliases and required slice before transport so an injected getter cannot rebind the later projection through caller-owned references. Nil/mismatched collectors reject without GET; valid calls perform exactly the existing one bounded GET/API2021-11-01. No constructor/auth discovery/cache/retry/pager changes. Wrapped/caller cancellation/deadline wins, canonical cause returns, no provider diagnostics or raw error text echoes.

Existing153DB+163MI manifests, both first runners, collectors/catalogs/shared guards/modules/pins/captures remain byte-unchanged. No implicit Evidence/SKUEvidence conversion, automatic inventory coordinate reconstruction, source service-availability default, region score or CLI registry change.

## Exact table contract

Generic prevalidation failure identity: ID sql-capability-catalog, sheet SQL Capability Catalog, description SQL capability catalog diagnostics; deployment eligibility is not assessed.
Validated Database identity: ID sql-database-catalog, sheet SQL Database Catalog, description SQL Database capability catalog diagnostics; deployment eligibility is not assessed.
Validated MI identity: ID sql-mi-catalog, sheet SQL MI Catalog, description SQL Managed Instance capability catalog diagnostics; deployment eligibility is not assessed.
All have schemaVersion1.0 and metadata Name region-selection, Version1.0.0, Description Bounded SQL capability catalog diagnostics, Author Cloud Assess, License MIT, Type internal. These are new diagnostic metadata, not an optimal-region recommendation. Zero timestamps remain zero; this pure adapter invents no execution timing.

Columns, exact order:
Subscription ID; Region; Resource type; Capability coordinate; Version; Edition; Name; Family; SKU name; SKU tier; Capacity; Standalone; Instance pool; Catalog availability; Deployment eligibility.

Complete Catalog rows contain ONLY required known groups sorted by full canonical coordinate. row.SubscriptionID equals canonical catalog SubscriptionID and its first cell. Other binding cells are canonical Region/ResourceType. Preserve all coordinate/authoritative metadata/modes exactly; capacity canonical decimal; known State0 gives catalog availability available, known State2 unavailable. Final cell always not_assessed. Do not collapse repeated SKU names, infer missing fields, choose modes or turn catalog availability into deployability. Complete with no requirements has allocated empty Rows.

Every successful table has Health.Name region-selection, Records exactly row count, Status completed_with_warnings and the EXACT existing sql_capability_catalog_only limitation warning/message. Unknown Catalog (including nil, missing/unknown required coordinate or operational unknown) gives allocated empty Rows, no prefix, and a second warning: code sql_capability_catalog_unknown, message Required SQL capability catalog evidence is unknown; no capability rows are reported. It returns nil error. Unknown is not an available-empty catalog.

Invalid or limit errors produce a safe allocated empty failed table, nil warnings, Error code sql_capability_report_invalid or sql_capability_report_limit respectively, Error.Message SQL capability report could not be safely projected. Returned error SQL capability report rejected [CODE]. Existing Catalog projection_limit maps to report_limit; other validation/constructor-family rejection maps invalid. Context failures use code sql_capability_report_cancelled, same safe message and canonical context error. Caller context is rechecked before every final success/failure return, including empty/unknown results.

## Independent limits, ownership and context

Inherited public counts8192groups/required,1000selected,256aliases,512component/labels,2051full coordinates and16711680SQL text remain unchanged. New diagnostic text cap is8388608bytes (8MiB), explicitly stricter than generic PluginTable16MiB to leave room for mixed reports. Charge exactly generic validator text occurrences: SheetName, metadata Version/Description/Author/License, table Description, each Column, each row Cell and warning Message. IDs/schema/health codes use their fixed safe grammar and are not charged by generic text rules. Premeasure entire admitted output before Rows/cell allocation; cap exact succeeds, +1 fails atomically with no prefix. Independently reachable pure Project fixtures must stay below ALL inherited Catalog caps. Do not claim an output overflow through the 1MiB collector transport if the fixture cannot reach it. Validate completed output with unchanged ValidatePluginTables.

All mutable table Rows/Cells/Columns/Warnings/Error are owned per call, input unchanged and mutation directions isolated. Repeat/concurrent calls are deterministic. No parent/aggregate result or renderer is mutated. Fixed manually derived terminal schedules frozen BEFORE report mutations: empty complete Project22 = unchanged Catalog21 + final report context1; empty complete Collect57 = adapter entry1 + nil preflight Catalog9 + Collect25 + Project22, both families. Never recalibrate schedules under reversals.

## Actual report integration and QA

Use existing result.BuildWithPluginTables and existing extension schema1.1 without altering schemas/builders. JSON/CSV/XLSX/shared tabular output must preserve full coordinates, canonical values, limitation/unknown/failed health and existing subscription-ID masking; test real files/readback and formula-safe string cells where applicable. No new report may fabricate Findings/Recommendations or score changes. Test complete/unknown/failure as well as two SQL families together with existing core findings.

QA-INT-01 clarification after reading actual sarif.go: SARIF is PRIMARY-FINDINGS-ONLY. It validates plugin extension but does not serialize plugin catalog rows or plugin health/provenance. Preserve that existing behavior and test no synthetic SQL findings/rules; do not claim SARIF carries this diagnostic table or extend its public schema. JSON/tabular outputs are where plugin health is represented. This corrects the earlier reconciliation's proposed SARIF health wording, not an existing implementation defect.

Freeze independent literal report expectations before API; verify integrity without claiming runtime coverage. Then publish a consuming runner while BOTH APIs are absent and inspect BOTH actual native missing-API logs BEFORE report production. Add independent output cap exact/+1/public inherited guards, fixed contexts, HTTP/one-GET/preauth/ownership, actual renderers/redaction and result completeness controls. Preserve533Linux/528Windows predecessor identities and all23captures. New compiling reversals require named assertion failures, exact restoration and passing focused/full baselines. Rehearse atomic report implementation+runner rollback to this verified contract base, region tests, exact tree/bytes restoration and retest. Full own final8workers/bothgates/source/complete logs/control identities/capture comparison/bounded evidence/full publication readback/final handover are mandatory.

No Azure/fixture/RBAC/refresh/subagents/protected merge/bypass/release. Independent review/backup/broader operational/live/release/advisory obligations remain OPEN. Update handover between tasks. Contract publication is not product implementation or hosted certification.

## Frozen identities and local prepublication evidence

26future cases in sql_report_oracle_test.go, SHA256ab90d0ea7142f825961117e171a8d56c86cbe6a5d89249c810ef981a9a849cdb. Independent budget proof sql_report_budget_plan_test.go, SHA256b9c738407e48412dd15097108c67918ed48007909379f39555a027facd6af95e. Integrity checks reject two compiling comment-byte drifts and pass after exact byte restoration. Both-family exact/+1 fixtures are admitted by existing Catalog. Focused region preflight/tests/vet and1466authored links passed. These checks validate the pre-API contract/oracle/budget preparation, NOT execution of future report APIs. Full own hosted/source review remains pending until live PR confirms.
