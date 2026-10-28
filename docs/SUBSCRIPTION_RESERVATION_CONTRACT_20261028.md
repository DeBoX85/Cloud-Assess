# Subscription-wide reservation collection

Task10 prerequisite, specification preparation only. Base: verified offline PR206/4627a4d9c2853d9337e8c0f959f241a6d319b416/treea8d7b71cfaa960536644162b62c090cf67ac72c7. No new production API is implemented at this checkpoint. Freeze and publish/read back this contract, independent expectations and strict consumer before implementation. Own final BOTH-native/source QA and private retention are required for the resulting candidate; PR206green cannot certify changed code.

## Confirmed scope gap

Pinned AZQR8e4f0577f3615e6c9014c031bcad079f235369cc `internal/scanners/plugins/region/crg/crg.go`, blob1ffb6e37cc1c21f6d8b5c9ce483f793387e0e6bc, was read in full for this binding review. FetchReservations lists every subscription group, then reservation summaries and instance-view details, without a source/target-region filter. Its skipped child failures and missing-count inference remain known source limitations, not requirements to reproduce unsafe completeness.

The target's existing `ReservationCollector.Collect` enumerates group pages but only follows groups in request.Region. Two calls for the comparison's source and target cannot establish whole-subscription coverage: a third group's region need not occur in inventory or comparison targets. Repeating the regional method over a presumed catalogue does not prove that every returned group was covered and repeats listing work. The existing regional API remains supported with unchanged filtering, request shape, limits, errors and original controls.

The bounded extension uses the existing collector/transport/decoders with a separate subscription request and an explicit overall enumeration status. It traverses every admitted physical group location once, using pinned2024-11-01 URLs and the existing Get instanceView expand. No SDK/version/source-pin/oracle refresh, live request, credential discovery or Azure fixture is authorized.

## Frozen API boundary

The strict consumer is retained as immutable `.go.txt` before production and compiled with a temporary Go overlay. The implementation stage copies its exact frozen bytes to the corresponding internal test paths; expectations must not be regenerated.

All types live in the existing private `internal/plugins/region` package. They do not change public CLI/report schemas.

```go
type SubscriptionReservationRequest struct { SubscriptionID string }

type SubscriptionReservationEvidence struct {
    SubscriptionID, Status, FailureCode string
    Requests []ReservationRequest
    Regional []ReservationEvidence
}

type SubscriptionReservationCollection struct {
    SubscriptionID, Status, FailureCode string
    Regions []string
    Reservations []ReservationUsage
    // Additional unexported owned admission state is required.
}

func (c *ReservationCollector) CollectSubscription(
    ctx context.Context, selected map[string]string,
    request SubscriptionReservationRequest,
) (*SubscriptionReservationCollection, error)

func (c *SubscriptionReservationCollection) Evidence() SubscriptionReservationEvidence
```

The request selects one subscription, not a region list. UUID normalization/selected-scope admission follow the existing regional policy. Nil context or an uninitialized collector is rejected safely before I/O. A cancelled/deadline context wins over invalid input or a simultaneous dependency error. Invalid caller state returns nil result and a fixed safe error (`request_invalid` for nil context/uninitialized collector/invalid request, `scope_invalid` for invalid selection); neither request fragments nor provider diagnostics are echoed.

Status is `complete`, `partial` or `unknown` and describes bounded whole-subscription enumeration, not permission completeness, deployment eligibility, reserved/allocated field presence or live service certification. Complete empty subscription is distinct from unknown/partial empty output. Missing counts retain the existing Known flags and must not become inferred zero capacity.

## Owned result and regional bridge

Regions contains each admitted observed group region once in first-service-encounter order, including a group whose admitted reservation list is empty. Reservations retains admitted details in group/list/Get encounter order, with each detail bound to its group region. No region is invented from a target list or physical-region catalogue.

Evidence returns independently owned requests and regional evidence, ordered by first observed region. Within each region, reservation detail order follows encounter order. Regional regrouping is deliberate deterministic target behavior; it does not claim the source's concurrent subscription ordering is stable. Requests and Regional have equal lengths and matching subscription/region coordinates. All Regional statuses equal the overall enumeration status. A partial observed empty region stays partial. If no group was admitted, Requests/Regional are empty and the overall status remains explicit.

Exported collection diagnostics are separate copies from its private admission snapshot. Mutating SubscriptionID/Status/FailureCode/Regions/Reservations cannot promote failed collection, redirect coordinates or alter accessor evidence. Accessor calls and separate runs do not alias mutable slices. Nil or caller-fabricated/zero collections return unknown evidence with empty subscription/requests/regional values; this cannot masquerade as an admitted complete selected subscription. This is structural provenance within this package, not authenticated Azure provenance.

The coordinator must retain overall status even when no regional rows exist. It can pass the accessor's Requests/Regional to the existing CalculateReservations, but an empty table alone never establishes a healthy complete subscription. The complete coordinator and required outputs remain separate subsequent work.

## Shared bounds and failure semantics

Preserve existing limits globally across the one subscription traversal:256getter calls,64pages per list chain,1MiB getter body bound,8MiB accounted response bodies and8192aggregate work units across admitted group rows, summary rows and VM references. Do not restart these counters per region/group. Overflow may receive one extra body within the existing per-call bound before total-byte rejection; no8MiB universal wire/lifetime-memory claim. Existing identity, continuation, JSON and fixed failure-code policies remain.

At most one initial subscription-group listing traversal is used. Existing continuation validation precedes transport, including every region's child list. No guessed nextLink, redirect, endpoint/audience change or peer subscription access. The injected getter keeps its existing every-attempt byte/deadline/body-close/no-redirect obligations.

Initial operational failure returns an owned unknown collection, not an empty successful result. Once a valid page has established an admitted prefix, subsequent operational/shape/budget failures return partial diagnostic evidence. Pages remain atomic; no rows from an invalid page are adopted. Existing peer recovery remains: an individual child/list/Get failure need not erase separately admitted later peers, but overall status never returns to complete. The first safe failure code is retained. Cancellation/deadline returns nil result and the canonical context error with no later calls or successful evidence publication. A wrapped context error from a getter must not leak its surrounding provider text.

The collector stores no per-call context, credentials, result cache or worker. Initial traversal is sequential. Cancellation is cooperative; a getter that ignores context cannot be forcibly interrupted safely. No wall-clock/RSS guarantee or background execution is implied.

## Independent expectations before production

Freeze literal outcomes for a complete three-region subscription including a region outside a comparison pair; complete empty subscription versus observed empty group; repeated region across distinct resource groups; initial and later operational/shape failures; valid peer recovery after one failed Get; forged/diagnostically mutated collection access; source-compatible regional filtering; every dependency cancellation boundary; two independent/concurrent runs; foreign/duplicate identities and unsafe continuations; exact and plus-one call/page/byte/work limits. Expected values are authored independently, never generated from the candidate or refreshed source captures.

First reproduce the insufficiency of pair-only regional coverage using the existing API and a healthy regional control; that is a proposed whole-subscription requirement, not a claim that the regional API is defective. Then compile the strict new consumer to observe only the intended missing API before implementation; missing-API compilation is preparation evidence, not a runtime assertion pass. After implementation, require healthy full consumer execution, compiling assertion-specific reversals/restoration, original controls and own full QA/retention before closure.

## Current evidence and safeguards

PR206 final index verifies its own quality40616376872/source40616376953, private captureeffaf284/receiptaf59ad63 and122fresh-restored members/9874097bytes equal originals. It is a verified preparation checkpoint, not a production coordinator. Source capture files,32frozen descriptor files,44legacy/49current obligations and all original expectations/pins stay unchanged. Never run prepare_literals.py.

The new dated native scans record0reachable/5imported-not-called/1module-only advisories. GO-2026-6866 andGO-2026-6860 concern selected x/net/html and were first-party classified from full official records after PR205; all five package and one module residual remain OPEN. No dependency update/suppression/gate change. Reevaluate after integration and before readiness/acceptance/release; source-mode observations are not a security guarantee. Human security assessment is considered AFTER overall project completion.

Issue189/DC-CI-EVIDENCE-20261022 is a closed bounded migration with ongoing private latest-main additive expected-head/no-force preservation, full readback/fresh restore and retention-until-Denis-deletion-approval duties. No automatic archival/independent-backup/DC006closure. No Azure/live retail/real credentials/fixtures/RBAC/subagents/unrelated dependency/action/SDK/source-pin/oracle refresh/protected merge/bypass/release. Separate coordinator/CLI/FullCost/health/workload/installed/readiness, ancillary, documentary, independent, operational, security, licensing and Gate004 obligations remain.
