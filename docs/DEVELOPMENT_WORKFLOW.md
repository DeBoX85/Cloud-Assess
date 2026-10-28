# Development and QA workflow

Current implementation/acceptance state is in [CURRENT_STATE.md](CURRENT_STATE.md). Research and baseline measurements are in [the research record](DEVELOPMENT_WORKFLOW_RESEARCH_20261017.md). The original parallel workflow was verified in PR148; later changes require their own exact-revision QA. Verification does not imply protected acceptance.

## Continuity

Follow [the continuity procedure](CONTINUITY.md) proactively during active work. [CONTINUE](CONTINUE.md) routes recovery and [RESTART_PROMPT](RESTART_PROMPT.md) is the stable copy-ready prompt. Update compact state and publish coherent WIP after meaningful slices, before long QA/environment work and before transitions; aim for 20-30 active-minute safe-boundary checkpoints. Preserve uncertain outcomes explicitly and verify remote publication. Do not wait for a user handoff request. No zero-loss, precise token-threshold or background-work guarantee is implied.

## Development-control records

Use [ACTION_REGISTER](ACTION_REGISTER.md) for compact ownership/dependency/closure tracking, [ROLLBACK](ROLLBACK.md) for targeted reversal/restoration and [EVIDENCE_RETENTION](EVIDENCE_RETENTION.md) for durable synthetic evidence. [Logging assessment](LOGGING_ASSESSMENT_20261017.md) records runtime boundaries without changing product schemas. The [controls contract](DEVELOPMENT_CONTROLS_20261017.md) defines this bounded tooling slice. Follow frozen-head successor routing rather than rewriting preserved PR heads.

## Developer loop

Use [REFERENCE_POLICY](REFERENCE_POLICY.md) to distinguish compatibility gaps, defect corrections and enhancements. Denis's pinned AZQR fork anchors feature/function/capability parity; additional primary evidence and project comparisons can justify improvements. Define their impact and independent expected behavior before implementation, retain frozen source evidence and use separate bounded review/QA. New work must fit the authorized scope and dependency sequence.

1. Read AGENTS.md, CURRENT_STATE, the applicable contract and relevant failure notes. Verify live refs and the actual worktree before edits. Preserve unrelated work.
2. Pick a small behavioral slice and its natural stop. Resolve uncertain semantics against pinned source and primary documentation. Freeze independent literal expectations before production changes where required by the contract.
3. Work in an isolated checkout. Publish coherent WIP checkpoints for recovery. WIP and local checks are not acceptance.
4. Run `python scripts/dev-preflight.py --package ./internal/plugins/region` for focused changes, or omit --package for all Go packages. Repeated --package arguments are supported. Use --static-only for workflow/whitespace checks without Go execution. Prerequisites are the existing pinned Go toolchain, Python, Git and initialized pinned submodules.
5. Inspect the final diff, expected results and test coverage before publishing the full validation candidate. Run meaningful controls appropriate to the behavior; avoid duplicating implementation details in tests. Do not count absent-API preparation or compilation failures as product passes.
6. Require complete final-revision Linux, Windows and pinned-source QA. Review failures before rerunning; no blind retries, skip flags or scope-based acceptance shortcuts. Changed code invalidates earlier-head certification.
7. Update [PROJECT_PROGRESS](PROJECT_PROGRESS.md), CURRENT_STATE and one coherent ledger entry after each completed bounded task, material finding/dependency change, published milestone and handover/end-of-run. Keep the dashboard brief and reconcile candidate QA, verified milestones, ordered remaining tasks and blockers against the obligation register. Update technical specifications when contracts actually change; use links instead of duplicating long status narratives across every document. Preserve immutable historical gate/evidence snapshots. The final PR body records the actual commit/tree/preview, complete logs and outcomes without an evidence-only commit loop.

The local preflight checks Git whitespace, generated workflow conservation, strict gate negative cases and, unless static-only, Go formatting plus uncached tests and vet for selected packages. Focused feedback does not replace the full repository, race, fuzz, mutation, platform or source gates.

## Parallel CI and protected checks

| Protected check | Required isolated jobs |
| --- | --- |
| quality | linux-native; linux-controls-1; linux-controls-2; linux-controls-3 |
| windows-validation | windows-native; windows-controls-1; windows-controls-2; windows-controls-3 |

Each worker starts with its own recursive checkout and the same pinned Go setup/cache. The three control workers run different suites sequentially; they never mutate a shared checkout. Native workers retain full tests, installed CLI/default and custom package checks, platform-specific protections and scans. Linux additionally retains full race, coverage, fuzz, module/license and tooling checks. The existing independent pinned-source workflow is unchanged and remains mandatory completion evidence.

The two protected job IDs and visible check names stay unchanged. Their validators use `always()` and reject any non-success dependency, missing/extra/duplicate identity or malformed dependency data. No ruleset edit, bypass or write permission is needed. The workflow's original PR/push triggers and cancellation of superseded runs remain unchanged; a cancelled run is never acceptance evidence.

## Maintaining the check inventory

The canonical original workflow and scripts/ci/quality-plan.json define the conserved platform inventory and explicit partition. Run `python scripts/ci/generate_quality_workflow.py --write` after a deliberately reviewed inventory change, then run preflight. The generator does not parse arbitrary YAML: it uses the fixed original step-block grammar and retains commands/options verbatim. A byte-for-byte generated-file check rejects manual workflow drift.

Add new checks to quality-plan.json's explicit per-platform additions arrays using only name, script and worker. Scripts must exist under scripts/tests and workers must be native or controls-1/2/3. No conditional execution, ignored errors or replacement of baseline names is allowed. Regenerate test.yml, regenerate gate-contract-probe.yml if gate bodies change, then run preflight and both platforms' full QA. New required checks must be added deliberately to this inventory/generation mechanism and its independent tests; never replace the frozen baseline with a weaker workflow to make drift disappear. Existing gate tests exercise actual subprocess exit codes and every worker's skipped/cancelled/failure cases. Native-command tests execute the actual Windows package/branding blocks with injected first/second command failures.

For final QA, inspect all four worker logs per platform plus both validator logs. Compare the union of control identities, with no duplicates or missing identities, to the predecessor. Record each complete log hash, all job conclusions and step outcomes, source capture reconstruction, actual native install checks, vulnerability scan distinctions, first/last job timestamps, queue/setup time and summed worker time. Do not compare only one shard or report summed CPU/runner time as elapsed speedup.

## Remaining boundaries

This tooling slice adds no feature parity, accepted target adapter, live Azure evidence, dependency/advisory closure, protected merge or release. Supported-platform, fresh-OS, load, hosted maintenance and human/security/license/Gate004 decisions remain open. The next bounded task is recorded in CURRENT_STATE; preserve unsupported VMSS scoring until a separate justified contract and tests resolve it.

## Hosted gate regression probes

The generated gate-contract-probe.yml runs only on isolated qa/gate-contract-*-success, -failure and -skipped push branches. It copies both production gate bodies verbatim under distinct probe names and uses genuine Ubuntu/Windows workers. All four dependencies succeed in the success case. In the other cases controls-1 deliberately fails its command or is skipped before scheduling; both actual validators must execute and fail. These expected red runs do not replace or weaken the product workflow. Cancellation is covered by local negative tests, not claimed as a hosted experiment.

After publishing the candidate, create three fresh isolated probe refs at its predecessor, then advance each to the same candidate to generate push events. Do not reuse a ref without inspecting its state. Preserve raw run/job metadata and complete gate logs. Run scripts/ci/verify_gate_probe.py with --run-json, --jobs-json, --logs-dir, --case and --sha. It rejects wrong revisions, missing jobs, wrong native runners, unrelated failures and missing validator diagnostics, and reports complete gate-log hashes. Record all three cases in the task PR. Ordinary task branches/PRs do not trigger this workflow.

## Review and handoff safeguards

scripts/tests/handoff.py checks stable authority headers on the principal resume/specification entrypoints; adversarial cases reject missing authority, obsolete per-task stops and loss of the historical/behavioral distinction. This is a routing guard, not a semantic proof of every paragraph. Review changed documents and their links manually as well. Current status belongs in CURRENT_STATE; immutable historical evidence remains historical.

Self-review, automated checks and independent review are separate evidence. No independent analysis service is callable in the current connection and no human review is recorded; record that obligation as open before protected acceptance/release. Continue authorized independent development without inventing a review or a new approval requirement. Additional live/release work requires its actual scope and access, and cannot be replaced by offline fixtures.

## Retrospective reconciliation and complete-log transfer

Freeze request-time PR/time/ref scope for fork concerns. Inspect heads/trees/parents/successor records; distinguish product/preparation/artifact/recipe branches; compare carried-forward code/frozen expectations and record inaccessible/unreachable evidence limits. Preserve failures and their replacements separately. Never reset/merge branches merely to force one lineage.

A successful fetch/write alone does not establish a complete log. Check chunk offsets/sizes, read saved bytes back, require complete runner footer/job/checkout identity, full control/capture inventories and hashes. Retain/reject partial copies and refetch before acceptance. Two48002byte prefixes were detected; cause unknown. Never weaken inventory assertions to accept a green job's incomplete local log.

PR-filtered run helpers are insufficient for a time-window audit: enumerate all events/pagination, compare to per-head runs, include superseded/cancelled heads and push recipes/probes, and inspect preceding boundary records. Preserve unavailable logs only as explicit no-runner/no-step skipped/cancelled metadata; never call them passed. The2026-10-19 audit closed this gap in its all-event addendum.


## 2026-10-20 security crossroads

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Autonomous batch paused at task6 own security QA; [blocker/options/decision](RETAIL_SECURITY_BLOCKER_20261020.md) controls the next action.177contract is frozen but unverified;178first consumer intentionally lacks production APIs; tasks7implementation/8contract not started. No dependency/pin refresh or gate bypass. Earlier next-task/no-blocker statuses are historical, not permission to proceed through this pause.


## Approved narrow security maintenance, 2026-10-20

Denis authorised the bounded Excelize update; the earlier decision pause is historical within this exception. See [CURRENT_STATE](CURRENT_STATE.md), [remediation contract](EXCELIZE_SECURITY_REMEDIATION_20261020.md) and [advisory review](SECURITY_ADVISORY_REVIEW_20261020.md). Exact fixed commit plus necessary transitive updates only; own full native/source QA pending. No retail production/protected merge/release. Preserve original failed177and first consumer178.


## QA-EXCELIZE-02, 2026-10-20

Initial narrow maintenance179failed unchanged go mod tidy/git-diff gate: go get left stale go.sum lines. Original179head/full failed Linux log preserved. Corrected successor removes only obsolete sums/adds required x/image module checksum and regenerates inventory input hash; versions, tests, source pins and notice bytes unchanged. Verify exact atomic pair rollback/restoration and full own native/source QA. Prevention: dependency publication additionally requires go mod tidy followed by an empty go mod tidy -diff and regenerated/check inventory; ordinary dev-preflight does not cover this dependency-specific gate. No suppression/bypass/retry without cause.


## Task7 retail checkpoint, 2026-10-20

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Narrow security maintenance/task6 prerequisites are COMPLETE VERIFIED OFFLINE in180/9fe12612; final quality39440501477/source39440501617 and cf93d214evidence supersede dated pending/wait snapshots. Task7 private retail implementation is IN PROGRESS, preliminary52frozen+15prepared+additional boundary cases pass locally. Final controls/rollback/full own QA/evidence remain PENDING; task8 research only. Follow180explicit live successor and [implementation record](REGION_RETAIL_IMPLEMENTATION_20261020.md). Continue existing authorised independent batch, preserving all old heads/pins/gates and live/review/acceptance boundaries.
