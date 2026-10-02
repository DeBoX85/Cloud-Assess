# Zone execution recovery checkpoint

Date: 2026-10-02 UTC. Status: BLOCKED by execution-runtime file-tool availability, not Azure or a missing user decision.

## Accepted remote baseline

GitHub live bootstrap/core-v1 was independently read after the runtime stall at `634aa083981fc327fe5c2cc7b42815793eb197fd`, tree `4729d848509dccf13f221359add14a1ef3cf8edc`. Preserve that accepted baseline.

| Slice | PR | Tested head | Native run | Linux quality / Windows validation | Merge |
| --- | --- | --- | --- | --- | --- |
| B2 bounded YAML Graph | #79 | 6842c50c6aa12ec10a7472483a6b528fc6938635 | 36963458729 | 110701991878 / 110701991724 | 7cddf9368272dd425315bdf3635b21fb974a9fe4 |
| B3a zone adapter | #80 | 63bb1c08b5de838af291a591d45bcbef9c060489 | 36966029765 | 110709854293 / 110709854116 | 0b76ad6b2870432357e5a38c1e8fe73d5214f11b |
| B3b canonical tables/reports | #81 | 8b1365218743d5b34fb934d24354699bf55465e2 | 36968213639 | 110716475829 / 110716476011 | 634aa083981fc327fe5c2cc7b42815793eb197fd |

All listed native jobs succeeded on exact final heads; logs, identity, base/preview/tree and both merged parents were inspected and verified. Six actual CLI checks, nineteen default plus nineteen custom installed package cases per host passed without acceptance skips for these slices. Latest documentation acceptance checked 223 local destinations, nine flags and one PowerShell snippet without scanning. Coverage 78.0% above unchanged 75% floor; default per-package instrumentation excludes cross-package validator calls in result/application tests. This is not feature/live certification. Existing module-only OpenPGP advisory remains open; Gate 004/release and live deferrals remain open.

PRs retain full evidence and postmerge/resume/rollback. B2 and B3a evidence was indexed by subsequent material ledger entries; PR81 final evidence must be indexed by the next material checkpoint, avoiding an endless evidence-only commit cycle.

## Unpublished B3c work and uncertain edit

Task: wire actual zone execution/health into `scan --plugin zone-mapping` and top-level `zone-mapping`, with offline honest `plugins list/info`. B3 overall still requires a second real adapter after zone acceptance. No Azure request was made.

Task-owned workspace: `/workspace/scratch/d726e08ce991/cloud-assess-yaml`, branch `feat/zone-execution-cli`, based on the accepted merge above. Original read-only source: `/workspace/scratch/d726e08ce991/repo-read-SfxMIl/azqr`, pin `8e4f0577f3615e6c9014c031bcad079f235369cc`, previously verified clean. Source parser captures used a separate source test copy; do not modify the retained reference.

Confirmed local edits before the stall:
- docs/ZONE_EXECUTION.md pre-implementation contract;
- internal/stages/config.go: copied configuration, plugin-only configuration/validation;
- internal/plugins/registry.go: fresh implemented-zone metadata/pending table and sanitized typed zone-result projection;
- internal/orchestration/operations.go: optional zone operation, lazily constructed adapter with captured ARM endpoint/options;
- internal/orchestration/coordinator.go: copied selected names, named-selection validation, plugin-only inventory skip, pending table retention, zone invocation and canonical table construction;
- cmd/cloud-assess/command.go: shared scan flag binding, zone-only command, regular --plugin selector, shared stage preflight and mode-dependent YAML discovery.

The initial coordinator/stage/zone changes passed focused tests before the subsequent CLI edits. That is not final B3c acceptance. CLI integration after refactoring had not been compiled or fully tested.

The last attempted multi-file patch would add SourceFile to rules.YAMLPlugin/discovery, add cmd/cloud-assess/plugins.go for offline list/info and update the old unnamed-plugin preflight test. Its outcome is UNKNOWN. The tool call stalled without a result, was terminated after bounded waits, and the subsequent read-only git status/file check also stalled and was terminated. Do not infer that the patch failed, succeeded or only partly applied. Do not blindly retry it or claim its files are backed up. This checkpoint records boundaries and design; it does not contain an exact backup of the unpublished diff. Scratch work may be lost.

## Resume procedure

1. Recover file-tool capability and inspect actual workspace/ref/status and candidate files before any edit. Treat the last patch outcome as unknown. If scratch is gone, clone the accepted remote baseline and reconstruct only the bounded B3c change from this contract, then rerun QA.
2. Verify actual tracked paths before source/target reads with explicit separate workdirs. A guessed source plugins.go follow-up failed during characterization and supplied no evidence; this recurs in FN-035/FN-036 and should be indexed in the next material failure note. The stall itself is an environment issue, not a proven product defect.
3. Review completed coordinator/CLI changes, repair any incomplete helper references/patch, format and run focused compilation/tests. Verify default behavior and caller ownership; do not assume an initial green test covered later CLI edits.
4. Add real Cobra/preflight/coordinator/application tests with synthetic authenticated zone transport. Tripwire inventory, Graph and every auxiliary operation in plugin-only mode; check mixed/scanner operation, selected scope/cloud/audience, healthy/empty/denied/malformed/later-page/cancelled rows, table/overall health and exit behavior. Verify invalid output retains a sanitized failed placeholder and other core data. Verify registry/source-file metadata, terminal controls, malformed YAML, unknown/conflicting/unused options, branded home precedence and help without Azure.
5. Update built native CLI/installed acceptance for the new offline commands and changed unnamed-plugin error. Add a compiling isolation/health negative control, confirm its named assertion rejects, restore and rerun checks.
6. Update relevant plan/roadmap/spec/alignment/ledger/failure notes and index PR81 exact acceptance. Preserve historic gates and explicit laptop/live deferrals.
7. Publish a coherent isolated candidate with exact local/uploaded tree and correct GitHub identity. Require actual exact-final-head Linux quality and Windows validation success with inspected logs, current base/preview/tree, expected-head protected merge and postmerge verification. No relaxed gates or unverified WIP merge.

## B3c contract

Pinned source normal scan selects internal plugins by --plugin. Top-level zone-mapping runs subscription discovery and zone only. Zone rows cover selected subscriptions and do not apply resource-group/type/tag filtering to location mappings; subscription filters and management-group discovery still apply.

Only zone is currently implemented. Unknown/unmigrated names and unnamed plugin-stage requests must fail before credentials. Collapse repeats deterministically using exact source names. Reject selected-plugin/-plugin contradiction, re-enabled regular stages in plugin-only mode and nonempty zone plugin.target-regions. Regular scans retain mandatory Graph and B2 YAML integration; zone-only skips unused Graph YAML discovery.

Create fresh pending metadata/headers for requested plugins before critical discovery. Preserve successful rows with explicit incomplete table/stage/overall health and artifact-before-partial exit. Preserve cancellation as context failure. Do not copy raw provider error bodies or continuation URLs. Valid typed failure UUIDs may be retained for correlation under existing masking. Capture production endpoint/audience and construct optional zone adapter only when invoked so ordinary scans gain no zone-only endpoint prerequisite. Injected code remains trusted, not sandboxed.

Offline list/info advertises implemented internal zone plus currently discovered bounded YAML Graph plugins, never all six source internal scanners as available. Preserve supported metadata/counts/resource types/source path, stable sorting, safe terminal text and no query/Azure execution. Strict YAML discovery errors must not yield a misleading partial registry list. JSON metadata is a small offline extension; keep source text functionality.

## Recovery and rollback limits

This proposal is a checkpoint only, not an implementation merge or full B3 QA claim. Accepted core-v1 remains the verified baseline. Recover/rebuild WIP with source/pins/toolchain checked; no background continuation is promised after a turn ends. No user laptop/Azure input or credentials are needed for the next offline task.

For an accepted merged slice, rollback is a reviewed revert on current history with dependencies and intended parent 1 checked, never reset/force or reference changes. DV-001, nonempty live optional stages, actual zone/plugin comparison, restricted visibility/load, release/maintenance/operations and existing advisory remain separately open.
