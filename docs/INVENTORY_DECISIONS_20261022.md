# Inventory admission decisions, 22 October 2026

This continues task 8 preparation. No production API, frozen expectation, collector, public schema or source pin changes. The original three provisional files and all50 cases remain unchanged. [The preparation record](INVENTORY_PREPARATION_20261022.md) remains the open checklist; this record resolves specific design questions without claiming the whole task is frozen.

## Error classification

Keep case23's required_inventory_invalid for a513-byte retained cell. An individual label must satisfy its declared format and maximum length; malformed or oversized labels, invalid IDs, unsupported status values, inconsistent payload bindings and invalid retained capacity are invalid input. Whole-call scope/resource/collection cardinality, accumulated entries/work and accumulated input/output text ceilings use the corresponding limit error. Caller cancellation/deadline takes priority and returns nil with the context cause. Do not use provider/parser error text in diagnostics.

Preflight outer cardinality before walking supplied data. For bounded map traversal use sorted raw keys, and keep array traversal in supplied order; never let map iteration select the diagnostic. A multi-defect input can hit the first applicable check in this prescribed traversal, but always rejects atomically. The final contract must specify the exact validation/accounting sequence and independent single-fault boundary cases before implementation. This resolves the cell-versus-whole-call classification; it does not claim global accounting or its boundary proofs are complete.

## Resource-ID parsing is not identity validation

The already pinned azcorev1.23.1 parser omits empty path segments, trims each segment and accepts a final unnamed child type. Seven synthetic cases ran successfully with Go1.26.10linux/amd64 against the unchanged frozen product go.mod/go.sum and existing dependencies, using GOPROXY=off/GOWORK=off/GOTOOLCHAIN=local and read-only module mode. No listener, HTTP request, credential constructor, production package test suite or source change was involved. The [probe](evidence/20261022/inventory-decisions/resource-id-probe.go.txt) contains independent literal expected parser outputs; its [summary](evidence/20261022/inventory-decisions/summary.json) records exact identities, command and actual exit. This checks SDK behavior, not the future admission API.

| Synthetic parser case | Actual parser result | Required new private boundary |
| --- | --- | --- |
| Normal VM ID | Named VM leaf, expected subscription/group/type | Admit when selected scope and declared fields agree |
| Doubled slash inside provider path | Empty segment disappears | Reject before SDK parsing using the existing empty-segment guard |
| Space-padded final name | Name is trimmed to vm | Reject any segment that differs from its TrimSpace form before parsing, preventing lossy identity joins |
| VM extensions type with no final name | Parse succeeds with empty leaf name | Reject an unnamed final resource instance even when a display Name was supplied separately |
| Nested SQL database | Type Microsoft.Sql/servers/databases, group rg, leaf d | Compare the full parsed type; never infer SQL version/edition/objective from that identity |
| Extension provider after a VM | Type Microsoft.Insights/diagnosticSettings, inherited group rg | Compare the terminal provider's parsed type, not the base VM type; unsupported service still remains unsupported |
| Literal vm%2Fchild name | Percent sequence remains literal | Do not decode/clean an opaque resource name or claim it equals a different Azure ID |

These are stricter requirements for the new private join boundary. Existing projection and aggregate behavior remains unchanged. An SDK parse success does not prove the ID exists or the caller owns it. Microsoft documents distinct nested/extension shapes in [resource functions](https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/template-functions-resource). No resource-specific Azure request is required to implement structural checks.

Before freezing: add independent future-admission literals for those accepted/rejected shapes, including ID/type/group disagreement after parsing; prove canonical duplicate handling and raw-record ownership; preserve the existing masking-prefix guard and opaque names. Do not credit these SDK executions as execution of those absent future APIs.

## Additional applicable-document reconciliation

Two further complete raw reads are reconciled, with exact183blobs matching192bytes in the summary:

| Document | Binding requirements | Historical states and remaining work |
| --- | --- | --- |
| REGION_STORAGE_SKU_CONTRACT_20261016.md | Required kinds originate in selected raw records, ALL kinds participate; missing/unknown required provenance makes whole family unknown; complete absent group is unavailable; global projector bounds and owned atomic results | Pre-implementation/STOP/old adapter counts are dated snapshots. Storage is completed offline work. Task8/9 must supply correct required-kind sets and matching bindings; tasks10–15/public/live/acceptance stay open |
| REGION_VM_VMSS_SKU_CONTRACT_20261016.md | Bind subscription/region/type; partial diagnostic prefixes cannot score; whole-result required unknown handling, zones and complete validation; VMSS unsupported scoring; distinguish dominated text from reachable wire limits | Old absent-API/STOP/control counts precede completed VM/VMSS work. Preserve all frozen captures and existing controls; do not rebuild collectors or claim deployment eligibility |

Cumulative raw coverage is803files, including the same326Go files;29 semantic views and127 other classified paths remain. Of the original123 documentation paths,115 still lack complete raw-read credit. No upstream SDK whole-source review or complete documentary audit is claimed. All44legacy/49current obligations remain visible.

## Next actions

Finish the remaining applicable-document reconciliation, global budget formula/proofs, typed union/status/unused-row cases, missing-coordinate canonicalization and context/ownership/forgery controls; publish and verify the final task8contract before task9's strict consumer and BOTH native missing-API observations. Keep Storage all-kinds, exact Cognitive kind, SQL no inferred coordinates and VMSS unsupported scoring.

PR191's earlier job/run contradiction has now resolved to completed/success for all10jobs. Rechecked the final gate log, which equals the earlier complete log byte-for-byte, and reran the strict unmodified verifier: all logs/steps,622Linux/617Windows controls,87.6%coverage and23exact source captures pass. Residual3imported-not-called/1module-only advisories remain open. PR192 and this successor retain their own separate QA state in their live PR bodies. Recent raw logs remain local under issue189 DEFERRED OPEN.
