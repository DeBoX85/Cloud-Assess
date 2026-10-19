# Meter-history implementation work in progress

Programme task5, based on published ready22ada22bc01da4fa1a62d4afa1cd6195fccda9b0, following verified contract PR173/42c414ef/evidence9eebdd82 and BOTH native absent APIs in PR174/d1. No Azure operations or public integration.

Implemented private query builder, bounded collector and strict page validator in meter_history_query.go, meter_history_collector.go and meter_history_decode.go. Preserve the completed UTC-month query, selected subscription/currency admission, exact HTTP statuses, bounded safe continuations, atomic pages/received-order sums, complete/partial/unknown health, owned results and incomplete-prefix exclusion from existing EnrichCost.

Local67frozen cases and supplemental boundary/alias/numeric/scope/context/HTTP/ownership/cost-integration tests pass. Fixed manually derived late schedules Query2, post-HTTP successful Collect6 and failure-finalisation2 passed BEFORE any compiling reversals. Focused vet passed. Full region/race, compiling reversals, additive inventory, source rollback, own final native/source QA/archive/readback and task5 reevaluation remain PENDING. This is recoverable WIP, not final product certification or acceptance.

QA-MH-02 is a fixture-only amendment: direct preimplementation42c414 HTTPClient proof shows the second1s attempt cannot wait for inherited5s cost throttling. Extend that one synthetic configuration line to10s attempt/30s operation; keep all assertions, production limiter/budgets,67oracle and first strict runner unchanged. Original boundary file remains preserved atd1. This narrow amendment supersedes only its byte-preservation requirement for the timeout line; see FAILURE_NOTES and retained proof. No production workaround or expectation deletion.

Next: complete meaningful compiling controls with exact restoration and passing baselines, register additive generated CI steps without replacing551Linux/546Windows predecessors, rehearse coherent implementation+runner rollback to verified contract treea7685642/exact restoration, publish final candidate before own complete QA. Preserve all original captures/pins and prior heads. Update handover after this slice and before task6. No critical blocker; no further approval needed within the authorised programme.

## Verification checkpoint after WIP publication

Recoverable implementation:80bbcbd83fec85316c8d780036d54c33513286f5/tree d0eeb36b23163228f1cb4d3cf2051877e5ed6fba/parent22ada22. All10complete files and human author verified. PR174 now explicitly forwards recovery there; its first failing preparation head remains unchanged. Full local region race test passed113.762s. Static preflight passed after additive generated workflow registration.29compiling controls are frozen before execution and now undergoing first local verification; no passing result claimed yet. They cover request/calendar/context, evidence exclusion/copy, scope, wire/pages/JSON budgets, continuation, metadata folding, row/column/ownership/currency/numeric bounds. Existing551/546identities remain mandatory; new controls are additive.

The frozen oracle explicitly accepts both empty-declared and empty-no-columns. The contract sentence about empty arrays must be read as empty columns requiring empty rows; declared columns with zero rows are valid. Preserve both independent cases and original contract bytes. This interpretation reconciles the literal oracle with normal zero-row query responses; it does not remove any expectation.

Official Go race-detector guidance (https://go.dev/doc/articles/race_detector), checked2026-10-19, confirms -race observes executed paths only. Local success is not universal race freedom or replacement for own final hosted QA.
