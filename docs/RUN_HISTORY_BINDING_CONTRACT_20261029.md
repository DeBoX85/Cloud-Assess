# Owned per-run history and currency binding

Task10 prerequisite from frozen verified PR210/7ad8233ed47d6a658405d5879380371d4551afed. This contract and independent consumers must be published/read back and actual missing APIs observed before production. No coordinator, public CLI, new transport, pricing query, source refresh, FX or absolute FullCost is implemented here.

## Why this boundary is required

Existing `MeterHistoryRequest` declares subscription, months, explicit clock and three-uppercase-letter currency. `MeterHistoryCollection` exposes subscription/currency/status/meters; its `Evidence()` projects only completeness and weights from mutable diagnostic fields. `CostHistoryEvidence` carries no currency. Existing `EnrichCost` computes relative regional percentages with each subscription's own weights. Retail remains USD only, with private eligibility. Consequently the coordinator needs an owned, explicitly scoped/currency-labelled admission snapshot before constructing its cost input. Do not replace the completed history or retail APIs.

The binder accepts declared dependency data. It cannot authenticate service provenance or prove the response used the declared months/clock: the old collection does not carry those fields. Future orchestration must associate one observed dependency return with its exact planned request. No cryptographic attestation, Azure visibility or invoice reconciliation is inferred from structural admission. A forged NEW binding with no private state is never admitted.

## Frozen API

```go
func AdmitRunHistory(context.Context, map[string]string, MeterHistoryRequest,
    *MeterHistoryCollection) (*RunHistoryBinding, error)

type RunHistoryEvidence struct {
    Admitted bool
    Request MeterHistoryRequest
    Status, FailureCode string
    ObservedMeters int
    History CostHistoryEvidence
}
type RunHistoryBinding struct {
    Request MeterHistoryRequest
    Status, FailureCode string
    Meters []HistoricalCostMeter // diagnostic, separately owned
    // unexported owned admission state
}
func (*RunHistoryBinding) Evidence() RunHistoryEvidence
```

This is pure admission with explicit Context, no getter, credential, cache, worker, clock default or I/O. Request and collection inputs remain stable during the call. No Context is retained. Validate the complete selected scope, request, envelope and ALL meters, including incomplete prefixes, before publication. Return nil plus a constant sanitized error on invalid input, never echoing input/collection fields. Cancellation/deadline takes precedence, returns the exact caller terminal error and no binding; nil Context is invalid without panic. Check through validation/copy and before publication; cooperative Context observations are not timer/heap guarantees.

## Admission and provenance

- Retain existing scope maximum1000, UUID/label/text admission, nonempty names and case-normalized duplicate rejection. Request subscription must be selected; validate existing `historyRequestValid` semantics (months1–12, nonzero UTC year2–9999, syntactic three uppercase ASCII currency characters). No ISO currency registry is inferred: `ZZZ` remains syntactically valid, not an Azure service-support claim.
- Output request uses lowercase subscription ID and UTC clock preserving the instant/nanoseconds; months/currency are unchanged. Collection subscription is the same selected UUID case-insensitively; collection currency equals request currency exactly. No defaulting, trimming, silent USD relabelling or conversion.
- Status is exactly complete/partial/unknown. Complete requires empty failure code. Partial/unknown require exactly one existing collector code: history_limit, history_continuation_invalid, history_response_invalid, history_request_failed or history_decode_invalid. Unknown must have no meters; partial may have an empty admitted prefix. Nil, zero, foreign, malformed or contradictory envelopes are invalid.
- At most8192 meters, each nonempty, at most512UTF8bytes, no controls/replacement/noncharacters; exact opaque IDs preserve case/Unicode and received order. Repeated exact IDs reject, including equal values; no folding/aggregation/last-write repair. Reuse existing projected auxiliary text admission, not an independent whole-run memory certificate.
- Every weight is finite, nonnegative, at most1000000000000; positive weights are at least0.000000000001. Zero remains zero, preserving existing `EnrichCost`'s explicit zero-as-unit-weight behavior. Do not apply positive scale-invariance claims to that special zero rule. Reject negative/NaN/infinite/tiny-positive/oversized weights even in a partial prefix.

Valid admission sets Admitted=true and privately binds the entire normalized request, status/code, observed meter count and cost eligibility. Complete produces History.Complete=true with its owned ordered meters, including an explicitly complete empty history. Partial/unknown produce History.Complete=false and a non-nil empty History.Meters; diagnostic prefixes remain available only through separately owned Meters. Neither exported diagnostics, original input, previous Evidence output nor another run may change private eligibility/coordinates/currency/weights. Each Evidence call deep-copies; nil/zero/diagnostic-only forged bindings return Admitted=false, unknown status, no coordinates/currency, zero observed count and incomplete non-nil empty history. No fabricated complete-empty result.

## Currency and arithmetic expectations

Independent integration literals use two subscriptions with separate USD/NOK histories. With meter regional differences +100%/-50%, weights9/1 yield85%, and1/9 yield-35%; multiplying the latter positive weights by1000 still yields-35%. The NOK request label remains NOK; no FX or cross-subscription monetary sum occurs. Only existing per-subscription relative `EnrichCost` arithmetic is invoked; the target's USD-only retail policy and frozen source expectations remain unchanged. Coordinator health must preserve incomplete history, and task12 still owns FullCost/absolute output semantics.

Microsoft first-party documentation read29October2026 distinguishes billing/pricing currencies and documents USD retail prices versus reference estimates for other currencies. These facts motivate explicit labels; they do NOT validate a conversion or invoice-equivalence algorithm. The relative-weight interpretation above is an engineering inference constrained by existing target arithmetic, not financial or negotiated-price assurance:

- https://learn.microsoft.com/en-us/rest/api/cost-management/retail-prices/azure-retail-prices
- https://learn.microsoft.com/en-us/azure/cost-management-billing/microsoft-customer-agreement/microsoft-customer-agreement-faq
- https://learn.microsoft.com/en-us/rest/api/cost-management/query/usage?view=rest-cost-management-2021-10-01

## Verification and remaining work

Freeze literal admission/alias/ownership/partial/empty/currency cases and at-limit/plus-one scope/meters/string/weight recipes. Prove missing APIs and separate fixture syntax without runtime credit, then implement on an isolated PRODUCT successor with unchanged consumers. Strict faults must compile, reach intended original assertions, restore source/pins/consumers/status and pass; compiler/panic/timeout failures never count. Own full Linux/Windows/source QA, original622/617 and301witnesses/host, coherent rollback/private readback/fresh restore remain mandatory for each changed head.

Per-call bounds are NOT cumulative retained-run accounting. Complete closed-shape retained-work/copy charging, exact coordinator request/result and structural phase ownership still precede coordinator production/task11. All32frozen files/44legacy/49current IDs/nine207/nine209/two210supplement files/original controls/expectations/pins and opaque distinctions remain; never prepare_literals.py. Separate compare/alternative-vm-sku/MCP/types/documentary/independent/live/operational/security/licensing/Gate004/release remain. Human security assessment only AFTER overall completion. DC-CI-EVIDENCE-20261022 https://github.com/DeBoX85/Cloud-Assess/issues/189 remains CLOSED bounded migration with ongoing latest-private-main additive expected-head/no-force preservation/full readback/fresh restore/revisit/retention until explicit Denis deletion approval. No automatic archival/independent-backup/DC006closure. No Azure/live retail/real credentials/fixtures/RBAC/subagents/unrelated refresh/protected merge/bypass/release/background/ETA/percentage.
