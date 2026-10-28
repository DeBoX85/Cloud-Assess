# Owned region run coordinator: contract preparation

Task10 is IN PROGRESS, not implemented or accepted. Product base is verified PR205/08a220a72edd902eb3f7eafc14b8e95f64a46f12/tree089172703ce24ad49c134f5d36893a8a981fa5ad. This document defines the behavioural boundary before production or a strict API consumer. Exact Go request/result signatures must follow the binding review below, not be guessed from a generic callback runner.

## Bounded batch

1. Verify205live final index/successors; inventory surviving worktrees, recover a clean product checkout and exact compiler/modules/APRL without touching unavailable links.
2. Trace pinned source orchestration against current target constructors/calculators and record concrete ownership, health and adapter boundaries.
3. Freeze independent literal scheduling/cardinality/terminal expectations and a separate reference model with adversarial controls. This is specification evidence, not execution of a future Go API.
4. Reconcile lifecycle, cancellation, shared cost, input/output budgets and unresolved bindings before a strict consumer/API freeze. No source-generated expected results.
5. Publish/read back the preparation candidate; own full native/source QA, original control conservation, private additive capture/readback/fresh restore. Do not inherit205green.
6. Reconcile dashboard/current/register/ledger/authority/restart, self-assess and continue the strict consumer/implementation sequence when ready. Task11 waits for the complete task10 implementation and its own QA/retention, not merely this contract preparation.

## Source and target phase map

Pinned source: DeBoX85/azqr8e4f0577f3615e6c9014c031bcad079f235369cc, internal/scanners/plugins/region/selection.go, blob2997f4c3d3617dc5a2e4d2330c5151e6e207f414. Reviewed Scan lines32–510 and region discovery511–569; this is targeted orchestration review, not new whole-repository audit credit. Existing captured pure arithmetic remains authoritative; no oracle refresh.

| Source responsibility | Existing target boundary | Coordinator requirement |
| --- | --- | --- |
| Inventory collected once, per-subscription counters and aggregate | CalculateInventory, BuildRequiredInventory, ProjectInventory | Validate/snapshot selected scope/resources before dependency calls; preserve declared completeness and all resource kinds |
| Subscription locations/physical regions and AZ mapping | AvailabilityEvidence accepts scoped zone counts; public execution/discovery adapter not wired | Do not invent a registry from region strings or copy a refreshed zone list. Explicit validated region/pair planning and subscription-scoped mappings must be bound before API freeze |
| Provider and required SKU availability per comparison | ProviderLocationCollector; typed Storage/VM/disk/Cognitive/SQL collectors; AdmitRequiredEvidence; CalculateAvailability | Invoke typed injected dependencies; admit coordinates before scoring; retain unknown/partial, catalog-only/not_assessed and unsupported VMSS scoring |
| Per-subscription meter history, one shared retail pass | MeterHistoryCollector, RetailCollector, EnrichCost | History once per active subscription; deduplicate meter IDs/physical regions for one shared retail collection; use Evidence accessors, not mutable diagnostic status to promote eligibility |
| Five quota families on unique source and target regions | RESTQuotaCollector/VM quota adapter, CalculateQuota | Deduplicate explicit (subscription,region,family) requests; validate every request before any dependency invocation |
| Whole-subscription capacity reservation inventory | ReservationCollector currently takes a regional request; CalculateReservations | Resolve complete request coverage explicitly. Comparing only source/target regions does not establish whole-subscription CRG coverage; no silent parity closure |
| Latency, cost and recommendation scoring | EnrichLatency, EnrichCost, Project | Reuse pinned calculations; do not reset unrelated comparison fields or let a successful later phase hide upstream health |
| Service, cost, quota, reservation and inventory sheets | Existing projectors; complete output ordering/FullCost is task12 | Carry owned inputs/diagnostics for later projection; task10 alone does not close all output/format obligations |

The source stores credentials/client/cache on the scanner and launches a goroutine per subscription. The target will use fresh per-run state and an explicitly bounded schedule. Matching source concurrency is not a capability-parity requirement. No production transport, credential discovery, network call, listener or Azure operation is part of this preparation.

## Normative ownership and execution requirements

- Input selection/resources/declarations and work requests are validated and owned before injected dependencies can retain or mutate arguments. No caller map/slice or result from a previous run may become run state. Input stability during snapshot construction remains a caller obligation.
- A run has no mutable global cache, stored Context, retained credential or background worker. Context is passed explicitly to every dependency. Initial implementation may be sequential with maximum concurrency1; any later parallelism needs an explicit bounded worker/lifecycle contract and tests.
- Work planning is deterministic and checked before execution. Selected scope maximum1000, comparisons maximum8192 and all existing component limits remain mandatory. Planned collector callbacks maximum65536; this counts calls to injected collectors, not HTTP attempts/pages performed inside a collector. Existing per-collector paging/body/attempt bounds remain separate. No universal request/RSS/total-lifetime-memory claim.
- Shared provider/history collection is per run/subscription; required-kind evidence remains per coordinate/kind. Retail receives only its existing MeterIDs/Regions/Now request, with no ARM transport or credential parameter. Any actual public-pricing adapter must preserve the existing unauthenticated transport contract.
- Explicit request planning must cover selected physical source/target coordinates. No fabricated target defaults, dates, complete empty collections, quota coverage, SQL coordinates or deployment eligibility. No SKU/name inference. ALL Storage kinds and exact Cognitive kinds remain.
- Partial/unknown evidence and disabled components remain explicit. Preserve the complete component health set; do not replace it with the last successful calculator's health. Nil/forged/mismatched collection results are invalid, not healthy emptiness.
- Unexpected dependency errors stop later scheduling, produce a safe failed diagnostic and must not publish a successful comparison result. Input/provider error text, tokens, resource names and IDs are not echoed by new coordinator diagnostics. Intended partial/unknown collection states remain data with their explicit health.
- Caller cancellation/deadline is checked before external work, immediately after every dependency, through derived work and before publication. A caller cancellation wins over a simultaneous dependency error. No later calls or successful output after a terminal cancellation. A dependency that ignores Context cannot be forcibly interrupted safely by a sequential caller; this is an explicit cooperative dependency contract, not a wall-clock guarantee.
- Existing typed collections and their accessors establish structural admission only, not authenticated service provenance. Diagnostic exports cannot overwrite private eligibility snapshots.

## Binding review still required before production

The following are concrete work, not a request for Denis to supply Azure access:

1. Freeze the exact private work-plan/request/result shapes and diagnostic result on failure/cancellation; bind raw selected input to existing inventory/admission snapshots.
2. Bind region discovery/AZ counts/mapping completeness without silently refreshing the source fallback list, and distinguish public registry/CLI wiring in task11.
3. Resolve complete reservation request coverage using existing regional collection without mistaking source/target-only coverage for the source's whole-subscription inventory.
4. Freeze cumulative retained-output/work charging, including catalog diagnostics, after inspecting actual closed target shapes. Per-component limits are not a cumulative run limit. The65536callback limit above is a deliberate initial design bound, not proof of overall bounded memory.
5. Bind declared history currency to the existing USD-only retail contract and shared cross-subscription diagnostics. No conversion or negotiated-price inference. A non-USD history cannot silently become USD evidence.
6. Then publish a strict consumer with independent literals BEFORE implementation, observe actual missing APIs at intended preparation assertions, implement bounded slices, compiling reversals, own full native/source/rollback/private restore and reevaluate before11.

## Independent preparation evidence

The scheduling/cardinality model and literal cases under docs/evidence/20261028/coordinator are a finite reference specification. They prove their stated algebra/terminal schedules only. They do not invoke a coordinator, validate real collector adapters, prove real timers, heap bounds, all phase interleavings or live Azure behaviour. Future Go consumer tests must enforce the same literal outcomes against the actual private API without changing expected values to fit code.

## First-party context guidance

Go context documentation and its contexts-and-structs article require explicit per-operation context propagation and explain derived-context cancellation. Azure SDK Go guidelines require Context as the first parameter for I/O and substantial CPU work. These support the execution design, not new SDK/version authority:

- https://pkg.go.dev/context
- https://go.dev/blog/context-and-structs
- https://azure.github.io/azure-sdk/golang_introduction.html
- https://github.com/DeBoX85/azqr/blob/8e4f0577f3615e6c9014c031bcad079f235369cc/internal/scanners/plugins/region/selection.go

All32frozen files/44legacy/49current obligations, original controls and module/action/source pins remain. Human security assessment is considered AFTER overall project completion, not an offline blocker. DC-CI-EVIDENCE-20261022 issue189 remains a closed bounded migration with ongoing private retention/revisit/latest-base lease/readback/restore duties. No Azure/live retail/real credentials/fixtures/RBAC/subagents/unrelated refresh/protected merge/release. No percentage/ETA, blanket security or completed-task10 claim.
