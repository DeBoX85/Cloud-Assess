# Required inventory preparation, 22 October 2026

Task 8 remains provisional. This checkpoint reconciles the recovered proposal against existing code and six additional complete document reads. It freezes no API or expectation and starts no production implementation. The separate security chat keeps its fixed product185 scope; any future product delta needs separate review.

## Current code and contract decisions

The code reference is product185 `566cc4d8b78d168e5216245a5a386f742a3a2f85`. The relevant blobs are recorded in [the preparation inventory](evidence/20261022/inventory-preparation.json). Review191 carries the same relevant production files. The recovered three-file proposal remains byte-for-byte unchanged under `docs/evidence/20261021/recovered-local/task8-provisional/`.

| Boundary | Existing behavior | Task 8 decision to retain |
| --- | --- | --- |
| Resource identity | `assessment.Resource` already retains Kind and SKU name/tier/family/capacity; inventory aggregation counts resources and drops some coordinates | Add a distinct private plan; preserve raw fields and separately normalized join keys. Keep legacy aggregate/projection behavior unchanged |
| Inventory completeness | Pure aggregation counts supplied records; it does not prove collection coverage | Require explicit per-subscription declarations. Missing declaration means unknown. Complete empty input differs from missing evidence |
| Storage | `StorageSKUCollection.Evidence` intersects every required kind, bounds kinds per SKU at16 and returns unknown when a required kind is missing | Retain ALL required kinds, including the missing-kind marker; never select the first kind or flatten kinds |
| Cognitive | `Catalog` binds canonical subscription, target region, resource type and exact canonical account kind | Keep separate kind-bound catalogs even when both contain S0; catalog completeness never establishes model/quota/deployment eligibility |
| VM versus VMSS | VM Evidence excludes VMSS from scoring; a separate VMSSCatalog preserves diagnostics | Keep VMSS scoring unsupported/empty and deployment eligibility not_assessed |
| SQL Database versus MI | Catalog validates actual full coordinate keys and family-specific policies; inventory has no complete SQL coordinate specification | Preserve both types, keep collected coordinate catalogs diagnostic, and infer no resource match from generic SKU, tags or capacity |
| Disk evidence | Existing `Evidence()` deep-copies complete values but does not validate arbitrary caller-built metadata or enum values | Validate the new private admission boundary before calling it. This is a coordination prerequisite, not proof the current collector emits malformed values |
| Scope and ownership | Existing Inventory guards selected subscription and masking prefix; Catalog methods expect stable input | Validate full new binding, reject aliases/duplicates and malformed unions atomically, and keep private snapshots separate from mutable diagnostics |
| Legacy corrections | Existing projection permits negative capacity fallback and up to1e12; private proposal rejects negative capacity and caps at2^31-1 | Label that as a deliberate stricter private admission rule, not source parity or a change to legacy inventory. Retain raw values in existing report paths |

The code inspection confirms these boundaries without requiring new cloud documentation or changing pins. Caller-built structural evidence does not authenticate Azure provenance. A complete catalog is not proof of workload placement.

## Preparation gaps before freeze

The [case inventory](evidence/20261022/inventory-preparation.json) records all50 unique provisional names and their exact member hashes. This is structural accounting, not execution or a new complete manual read of every scalar. The recovered Python preparation helper is retained evidence; do not rerun it to regenerate expectations.

The final freeze must resolve these explicit requirements and supply independent literal expectations or separately frozen boundary proofs:

1. Define a deterministic error taxonomy. The proposal says oversized cells are invalid but also defines limit errors. Case23 currently expects required_inventory_invalid for a513-byte cell. State whether individual malformed/oversized labels use invalid while whole-call count/text/work ceilings use limit; do not change the literal later to fit implementation.
2. Specify the whole-call entry/text accounting formula for raw records, normalized copies, requirement bindings, private snapshots, diagnostics and ALL supplied collection data. Per-family projector limits do not replace global accounting. Prove independent exact/plus-one count, work and text boundaries; disclose dominated limits.
3. Freeze additional nested and extension resource-ID/type/group examples, including malformed suffixes and alias/percent-encoding policy. Existing inventoryIdentity only proves selected-subscription correlation and masking syntax. Generic display Name must not be reverse-engineered from an ID.
4. Freeze mixed-union rejection, unused malformed catalog rows, explicit invalid per-family statuses, duplicate canonical bindings and wrong inner metadata. Storage and SQL accept complete/unknown collections; VM and Cognitive also admit bounded partial diagnostics. A generic partial rule must not override those typed validators.
5. Specify missing/unrepresentable SKU/kind disposition and canonicalization exactly. Safe raw labels remain diagnostic; no fuzzy service/type/region/SQL coordinate inference. Show that nonphysical/global rows remain retained without manufacturing regional requirements.
6. Require nonnil-context handling before calling existing methods, and independently prove terminal cancellation at early/late boundaries, atomic invalid-tail rejection, deep ownership, diagnostic-mutation resistance and concurrent repeatability.
7. Reconcile the remaining applicable documents and preserve every residual before declaring task 8 frozen. Six more documents are reconciled in [the documentary record](DOCUMENT_RECONCILIATION_20261022.md);117 of the original123 remaining documentation paths still lack complete raw-read credit.

These are bounded preparation tasks with conservative unknown/rejection policies. No critical choice currently requires Denis. Do not implement the proposed APIs until the contract and independent expectations pass their own checks and the task transition is recorded.

## QA and transition

PR191's two runs report success and all11 complete logs were retrieved. Content checks conserve622Linux/617Windows full controls,87.6%coverage/floor75%,0reachable/3imported-not-called/1module-only residuals OPEN, and23captures/315chunks/909865bytes equal exact Git. However the quality job119912379705 still reports in_progress/null in two job inventories, while its steps and footer are complete. Keep overall certification OPEN; neither alter the observed metadata nor infer completion from a run badge. This discrepancy can be paused safely while documentation preparation continues. [The safe summary](evidence/20261022/pr191-qa-observation.json) contains hashes/metadata only. Recent raw logs remain local under [issue189](https://github.com/DeBoX85/Cloud-Assess/issues/189).

After task 8 freeze and own verification, task 9 first publishes the strict consumer and inspects BOTH actual native missing-API observations before production edits. Preserve all622/617 existing control identities, source captures, pins, previous heads and unrelated worktrees. Tasks10–15 and all ancillary, live, operational, security, licensing, Gate004 and acceptance obligations remain in the unchanged44legacy/49current [task register](TASK_CARRY_FORWARD_20261022.json).

No Azure, live retail, fixtures/RBAC, dependency/API/source refresh, subagents, protected merge/bypass or release. Revisit deferred CI evidence issue189 at every handover/end of run and before readiness/acceptance/release. Consider an independent human security audit before acceptance/release; the separate AI review does not satisfy that obligation.
