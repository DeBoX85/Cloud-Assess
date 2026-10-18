# Remaining-work reconciliation and integration decision, 2026-10-18

This is task1 of the newly authorized OFFLINE integration batch, after complete SQL closure160c0793/PR165. Product baseline a3813439/tree855d8a29 remains VERIFIED OFFLINE; this record is administrative analysis, not a new product candidate. Live PR141-165 remain draft/open/unmerged; protected bootstrap/core-v1 a6fee427 unchanged. All older heads/worktrees are preserved.

## Evidence and disposition

| Area | Current evidence/status | Remaining work and priority |
| --- | --- | --- |
| SQL Database and MI | PR163/165 exact-head native/source/rollback/evidence complete;316frozen cases; catalog-only | Priority1: bounded collector-to-report diagnostic integration with explicit full coordinates |
| Region inventory | inventory_aggregation.go keys raw SKUName by type/region; assessment.Resource has no SQL version/objective coordinate fields | Never infer SQL requirements from name/tier/capacity; explicit caller requirements needed |
| Region service/scoring projections | service.go assumes Available absent missing-set membership; primary.go scores explicit Comparison inputs | Do not pass unknown/catalog-only SQL evidence into source service availability or scoring |
| Internal plugin execution | registry.go advertises5implemented plugins and rejects region-selection; coordinator.go executes that registry | Full required-kind coordinator/public region execution remains separate; do not enable it in this slice |
| Report extension | result.BuildWithPluginTables validates/copies existing PluginTable extension; JSON/CSV/XLSX/table and SARIF already handle plugin diagnostics | Priority1: exercise real collector->Catalog->PluginTable->canonical result->actual renderers, including redaction and health |
| SKU families | Source characterization7/7; protected accepted target adapters1/7. Storage/VM and VMSS/Cognitive/SQL draft capabilities have progressed beyond old not-started narratives | Keep published, verified offline, protected accepted and deployment eligibility distinct; do not repeat completed collectors |
| Remaining collection/orchestration | Full Cost/history/retail, required-kind evidence coordination and ancillary commands/MCP remain open in assessment | Priority2: separate bounded contracts after this diagnostic slice; no implicit completeness or full parity claims |
| Broader QA/operations | Existing offline safeguards/native/source checks passed at recorded heads | Live permissions/scale/sovereign/credential scenarios, nonempty optional stages, Arc SQL shape, nested nonproduction MG, freshOS/operational rollback/hosted maintenance remain open |
| Acceptance/security | One module-only advisory; independent human/external review and independent backup/retention not established | No refresh/suppression, protected merge, security/license/Gate004/release inference |

Read current source paths named above plus internal/assessment/plugin_tables.go, internal/renderers/tables/plugins.go, internal/renderers/json/json.go, internal/renderers/csv/csv.go and result/plugin_tables.go. Canonical JSON and tabular renderers mask known subscription identities; row.SubscriptionID must be canonical and independently tested. Existing generic PluginTable bounds/health and result schema extension are reusable without changing public schemas. SARIF must retain health/provenance without inventing recommendation findings for catalog rows.

## Chosen next slice

Task2 freezes an independent SQL diagnostic report contract and literal expectations before API implementation. Task3 adds a bounded internal adapter taking a concrete existing SQLCapabilityCollector, selected scope, explicit SQLCapabilityRequest and full required coordinates. It must prevalidate before GET, invoke existing bounded Collect and Catalog without changing frozen semantics, preserve authoritative metadata/modes/full coordinates, and produce only a diagnostic PluginTable through the existing extension. Complete/unknown/failed/cancelled states remain explicit and atomic. No implicit SKUEvidence/Evidence conversion, source service-table Available default, automatic inventory matching, scoring or public execution.

Task4 performs own exact-head native/source QA, all533Linux/528Windows predecessor identities plus justified independent new controls, actual multi-format/redaction/health tests, coherent adapter+runner rollback/restoration, full logs/capture comparison, bounded evidence/readback and final handover. Tests use offline injected getters and temporary outputs; no Azure resource fixture/RBAC operations. Preserve all old runner/oracle/primary/capture/pin bytes and unrelated work.

## Transition and authorization

Reconciliation task COMPLETE after static documentation/link checks and full remote publication/readback. Next task2 is contract/independent expectations, not production edits. User authorized all four offline integration tasks2026-10-18; no critical input or Azure access needed at this point. Update CURRENT_STATE/ledger/handover/prompt BEFORE each transition. Scope decisions and unresolved findings must be recorded; do not guess readiness from old roadmap snapshots. No Azure/fixture/RBAC/refresh/subagents/protected merge/bypass/release. No background/outage-proof or zero-loss guarantee.
