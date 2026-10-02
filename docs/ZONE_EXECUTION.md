# Zone coordinator execution contract

Status: B3c coordinator slice, 2026-10-02. The accepted zone adapter and canonical tables are wired through a typed coordinator request; CLI/list-info wiring remains a separate slice. This document records contract and local evidence, not final native acceptance. Final head/run/merge evidence belongs in the proposal and next material ledger index.

Pinned source `8e4f0577f3615e6c9014c031bcad079f235369cc`, `internal/pipeline/builder.go:BuildPluginOnly` performs subscription discovery and plugin execution without inventory/Graph. `internal/scanners/plugins/zone/mapping.go` defines metadata, five headers, Zone Mapping sheet and its distinct description. Source skips failures/continuations; [ZONE_MAPPING.md](ZONE_MAPPING.md) records target corrections. No source/library/dependency pins or equivalence normalization change.

## Behavior and ownership

`orchestration.Request.InternalPlugins` supports exact `zone-mapping` only; repeats collapse deterministically. `PluginOnly` selects a fresh all-disabled regular configuration with plugin enabled. Named selection enables plugin in regular mode while Graph remains mandatory. A supplied regular config cannot re-enable any regular stage in plugin-only mode. Unnamed plugin requests, unavailable names, missing selected operation and unused target-regions reject before discovery. CLI must separately reject selected-plugin/explicit-disabled-stage contradictions before credentials when it is implemented.

Prepare owns stage enablement/options maps, name slices, filters and existing request definitions. Registered stage option values are immutable scalars. Zone receives a copied discovered-subscription map. Requested pending metadata/headers exist before scope discovery, so critical failure/cancellation cannot silently omit requested output. Zone-only skips inventory and all regular stages. Mixed mode retains ordinary assessment execution and adds the zone table. Subscription discovery/filtering still applies; resource-group/type/tag filters do not filter subscription location mappings.

The optional production operation captures selected ARM origin and token options and creates its zone HTTP adapter only when invoked. Unselected scans gain no zone-only prerequisite. Existing required operations contract is unchanged. Injected operations remain trusted and must honor context; this is not a sandbox.

## Projection and health

Preserve source metadata and both distinct descriptions, five header/cell values, canonical row identity and adapter order. Return fresh metadata/columns/row cells/health. Sanitized typed failure codes and validated UUID correlation are projected; raw error messages, provider bodies and continuation URLs are never copied. Row/failure counts are bounded before allocation, full table validation applies, and invalid trusted output becomes a failed empty placeholder rather than discarding other assessment data.

Successful empty/nonempty results are complete. Subscription/later-page failure retains healthy rows and yields failed table/stage with partial assessment. Cancellation remains a context failure with retained rows and failed assessment. Critical scope failure retains requested skipped table headers. Default unselected results retain schema 1.0; plugin-bearing coordinator results use schema 1.1. Application persistence/exit and raw/masked JSON/CSV/Excel support are inherited from canonical infrastructure; actual CLI-to-application acceptance is still required.

## Verification and limits

Focused tests run the real zone parser/scanner via a synthetic bounded getter through the coordinator, asserting exact selected endpoint/API/byte bound and literal rows. Every regular operation is a plugin-only tripwire. Cases cover owned config/names/options/subscription map, repeat/concurrent runs, mixed Graph execution, valid empty data, malformed response, failed subscription with retained rows, invalid output placeholder, critical discovery, cancellation, unknown/unnamed/unused/conflicting-mode requests and missing selected operation before discovery. Direct metadata/projection tests independently specify source descriptions/columns/cells and sanitized failures.

A compiling mutation enabling inventory in plugin-only mode failed the named DiscoverResources tripwire; code was restored before full local race/vet checks. Existing adapter pagination/authenticated cloud/audience/bounded-body tests continue to run, but this new coordinator fixture does not itself exercise successful authenticated transport or actual command dispatch. Those integration/installed cases belong to the next CLI slice. No Azure requests or live plugin equivalence are claimed.

Base: handover merge `d88f9294482eafc4e724dbe2248fc09ed4f826b3`, tree `27e0e50dc86d1b970c6ddb94cb33ac0dec857f2a`. Require both final-head native jobs and exact tree/identity/base/merge verification before acceptance. Rollback by reviewed revert with dependencies/mainline checked. Next: shared normal/plugin-only flags/preflight, actual zone CLI, honest offline list/info and full authenticated command/application/report/exit acceptance, followed by a second real adapter. Gate 004/release and live deferrals remain open.
