# Development execution, QA and recovery plan

Date: 2026-10-02 (Europe/Oslo). Starting accepted branch: `bootstrap/core-v1` at `c79c2ffc9bfcec66da9f24dab631eb6bef6d1836`, tree `f2c99952554e11042d90ff0dd0597bda07b70e4d`, merged PR #77. This plan governs future implementation batches; it does not declare Gate 004, release readiness or full feature parity. [TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md) remains the product contract, [ROADMAP.md](ROADMAP.md) orders work, and [QA_PROCESS.md](QA_PROCESS.md) defines existing mandatory checks.

## Goal, scope and authority

Deliver functional equivalence to pinned AZQR with independently adjustable branding, preserving APRL/AOR/CUSTOM recommendation libraries and validated generic assessment meaning. Equivalent features do not require identical architecture or every source defect. Every deliberate difference must be source-grounded, documented and tested. The five-active-field immutable branding workflow and offline rules command are accepted bounded milestones; they do not certify the entire toolkit.

Source pins remain AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`, APRL `60eaddda76541f6adbc1c5ffa686829807e55e29`. Never advance pins, refresh dependencies, change normalization, add platforms, expose an HTTP service or expand live scope as an incidental feature fix. Those changes need their own justified contract/evidence record. The user's current authorization covers repository development, QA and publication through protected PRs; it does not authorize creating Azure fixtures, changing roles or scanning production to replace a missing test scope.

The user currently cannot run laptop/Azure tests. Continue deterministic engineering without those inputs. Preserve explicit deferrals and report a blocker only when no independent authorized work remains. No user silence or elapsed time is an approval, answer or live-test result.

## Alignment controls and per-slice acceptance record

Before production edits, record these items in the slice document or PR:

1. Specification requirement and pinned-source file/function/trigger; selected scope, effective defaults, data/provenance and report/error semantics.
2. Concrete before/after behavior and independently defined expected outputs. Include source omissions or defects without assuming the source is always correct.
3. Intended changed files/interfaces and exclusions. Separate unrelated refactors, dependency refreshes and feature additions.
4. Acceptance criteria, tests, negative controls and evidence boundaries; distinguish injected execution, actual CLI/package execution and live Azure comparisons.
5. Existing known-good baseline, compatibility/rollback consequences, dependencies and resume point.

At final review compare the implemented diff with this record and TARGET_SPECIFICATION. Unexplained new behavior, missing requirements, changed source labels/pins, added write capabilities or wider scope stop dependent merge/acceptance claims. Narrow or split the slice, document a deliberate correction, or defer an uncertain feature; do not normalize away the mismatch or silently redefine success.

Keep batch authorization separate from PR size: a batch may include several steps, while each PR remains self-contained, working, reviewable and independently revertible. Avoid speculative stubs or unused common APIs. Build shared plugin behavior with one real plugin, then confirm the abstraction against a second before generalizing.

## Ordered offline implementation batches

| Batch | Deliverable and dependency | Required acceptance | Fallback when blocked |
|---|---|---|---|
| B1 | Scanner-specific `scan <key>` commands, reusing existing orchestration/selection | Independently captured source command keys; inherited flags; generic-vs-specific precedence; invalid input before observed auth; no state leakage; synthetic CLI-to-coordinator scope/result checks; actual installed CLI help/preflight; unchanged generic rules output | Keep generic scan available. Preserve selection helpers and defer only disputed command semantics; do not invent aliases or rewrite assessment meaning. |
| B2 | Bounded source-schema YAML plugin parser/discovery and normal Graph integration | Inline/queryFile equivalence; defaults/duplicate precedence; source attribution; supported type/ID overrides; exclusions; unknown/duplicate/malformed YAML; byte/query limits; traversal/symlink policy; actual query-to-report fixture and unchanged default scan without plugins | If discovery trust policy is unresolved, characterize explicit-path loading first. Do not ship auto-discovery, weaken confinement or claim source enablement compatibility without a recorded decision. |
| B3 | Canonical plugin table/health/report contract plus zone-mapping, mixed and plugin-only dispatch, honest list/info | Metadata/selection/unknown names; literal nonempty and empty rows; partial failures distinct from success; all requested tables; collision-safe filenames/sheets; privacy/formula safety; JSON ownership; allowed request destinations/methods/audiences; pagination/cancellation; actual executable fixture where feasible | Keep unavailable plugins unavailable and fail explicit requests clearly. Do not advertise a registry entry as executable before its adapter/report path works. |
| B4 | Service-health and SQL EOL, each a separate migration | Source query/header/order/formatting, source empty branches, filters, malformed rows, denial/throttling/partial failures, cancellation and independent tables | Retain source literals and synthetic fixtures; record live gaps. Do not label absent Azure data as equivalent nonempty behavior. |
| B5 | Carbon, AI governance and region selection, separately | Source calculations, clocks/windows, batching, metrics versus ARM audiences, caches, pagination/size bounds, quotas/prices and all output sheets; per-run isolation; every required input's error/partial health | Start with pure calculations/decoders and request fixtures if service access is unavailable. Unknown contracts stay open; region selection does not become a simpler invented score. |
| B6 | Offline report comparison and alternative VM SKU commands; separate MCP design/implementation | Source input/format/scoring/duplicate semantics and file failure/privacy checks; for MCP, transport/exposure/auth/lifecycle design and no unintended listener | Defer network exposure until explicit design review; bundled SKU provenance is not execution parity. Generic semantic comparison is not a substitute for the source's operator command. |
| B7 | Remaining offline acceptance/release preparation | AR-02 volume/load/metadata limits, artifact/inventory/license review, operations/runbooks, hosted maintenance controls and exact candidate evidence reconciliation | Produce provisional evidence and list remaining environment/approval boundaries. Never label candidate ZIP QA as fresh-OS or signed release validation. |

Use [FEATURE_PARITY_CHARACTERIZATION.md](FEATURE_PARITY_CHARACTERIZATION.md) for source entry-point contracts. Detailed per-plugin endpoint/body/column/query semantics must be reviewed before its migration is accepted. A fallback changes the next implementation step, not the feature obligation or accepted scope without documentation.

## Verification and security protections

For each meaningful behavior, use literal source/schema-grounded expectations and complete row/output comparison where feasible, rather than count-only or target-generated oracles. Keep independent source captures with pins, hashes, notices and explicit regeneration rules. Synthetic evidence is labelled synthetic. Shared assumptions among reviewers and fixtures remain a risk; independent review must inspect the source contract as well as the target implementation.

Cover healthy/nonempty/empty, malformed/denied/throttled, cancellation, later-page failure, output failure and partial-data paths according to the affected adapter. Add tests to the existing required jobs. Avoid mirrored or redundant tests; coverage percentage is a regression signal, not proof of correctness. New critical guards need a compiling negative control that fails the named assertion, followed by restored-code checks. Existing comparator mutation coverage does not certify every new guard.

Verify all new authenticated destinations and continuations before credentials are sent; preserve cloud endpoint/audience selection, redirect controls and disabled automatic registration. GET alone is not a no-write proof; read-oriented query POSTs require verified endpoint/body contracts. Use synthetic token canaries and request tripwires. Do not log tokens, authorization headers, raw tenant evidence or resource payloads into Git. Arbitrary injected application code is trusted, not sandboxed.

Plugin caches, selected scope/regions, clients and results must not leak across sequential or concurrent runs. Use different input sets, repeated runs, caller-data ownership assertions and race tests. Bound parsing/output/API work from an explicit contract; investigate exhaustion without silently truncating data. Preserve honest warning/completeness/exit behavior and healthy results when optional work fails.

Reports retain existing privacy/replacement controls. Review additive canonical table/schema changes, all formats, filenames/sheet collisions, formula-like cells, redaction and identity-bearing SARIF. Preserve the existing limitations: multi-file exports are not transactions; crash durability, arbitrary-directory ACL privacy and consumer re-save safety are not universally certified.

## Per-PR QA and publication sequence

1. Read current branch HEAD/tree, applicable instructions, relevant failure notes and dependency state. Inventory paths before reads; verify actual submodule contents and toolchain, not just gitlink metadata.
2. Develop on an isolated branch/check-out. Run focused checks as fast feedback. After correction, rerun affected tests; broaden when interfaces, defaults, dependencies or failure scope justify it.
3. Review the final diff, contract coverage, security boundaries, fixture provenance, documentation and rollback consequences. A separate author/reviewer workflow can be trialled when explicitly available/authorized; record who reviewed and what they checked. Self-review and independent expected data are not a claim of independent-person review.
4. Commit and publish the coherent candidate on a proposal branch. Verify human commit identity `Denis Bogunic <134433647+DeBoX85@users.noreply.github.com>`; automated maintenance commits retain bot identity. Compare exact local/published trees and verify the intended base.
5. Require actual success, with logs inspected, from Linux `quality` and `windows-validation` on the final head. Required controls include pinned provenance, formatting/module/inventory, actual CLI/default/custom package acceptance, race/vet/coverage, bounded fuzz/mutation, documentation and vulnerability scans. No skipped installed test or success on an earlier head substitutes for acceptance.
6. Confirm current base/merge preview matches the reviewed candidate. If base changes, rebase/reconcile explicitly and obtain checks for the new candidate. Merge with an expected-head condition, then verify merged commit/tree and clean fetched state.
7. Retain tested head, base, tree, run/jobs, reviewer/limits and merge SHA in the PR; index final evidence in the next material ledger checkpoint. Avoid endless evidence-only commits. Update current roadmap/status wording without rewriting historical QA snapshots.

Full native QA is mandatory at each merge. Focused development tests are not a waiver. A clean installed candidate must retain strict VCS stamping, CGO-disabled native target, readonly modules and the documented build profile; scratch build-context failures are not a reason to weaken artifact provenance.

## Logs, checkpoints and feedback

The ledger records material behavior, milestones, findings and accepted limits. FAILURE_NOTES records confirmed development mistakes, symptoms, cause versus hypothesis, correction, recurrence and prevention. Use PR/workflow logs for exact execution evidence, not unlimited copied transcripts. Retain failure artifacts; minimized relevant fuzz inputs become regression seeds after review. Never rerun until green without classifying the failure or relabel a timeout as a proven assertion defect.

At a checkpoint record: baseline and proposal refs/SHAs/trees; requirement/batch; changed files and source pins; last completed check and failures; pending native run IDs; known uncertainties and live deferrals; exact next action; rollback target and dependent work. Mark states as NOT STARTED, IN PROGRESS, VERIFIED OFFLINE, VERIFIED LIVE, BLOCKED or DEFERRED with explicit evidence. These labels describe the bounded task, not global release status.

Measure reviewable work completed, time spent implementing/verifying, rework and escaped defects across comparable batches before revising pace/ETA. Number of PRs, token consumption or green coverage cannot establish feature completion. Calendar forecasts remain conditional on active sessions, scope, test environments and approval dependencies.

## GitHub/session interruption and backup recovery

| Event | Required response | What must not be assumed |
|---|---|---|
| API timeout during commit/ref/PR/merge | Treat outcome as unknown. Query the remote ref, commit/tree and PR state with bounded retries/backoff. Continue only after identity is established; check expected parent/head before retrying a mutation. | A lost response is not proof of failure. Blind retries can create duplicates or overwrite work. |
| GitHub read/write capability unavailable | Keep local Git commits and uncommitted diff status, stop dependent publication/merge and record the blocked operation. Try an already authorized sufficient connector; do not request credentials, probe a browser session or change branch protections. | A local commit or proposed tree is not a durable remote backup or accepted merge. |
| Remote branch moved | Fetch, preserve the local candidate on a named branch, inspect divergence and reconcile in a new reviewed candidate. Re-run relevant final QA; never force-update core-v1. | Matching file counts or commit titles do not establish matching content. |
| Scratch/tool/runtime lost | Rebuild from the last verified remote proposal/base, read checkpoint/PR/ledger, verify source/toolchain and reproduce outstanding checks. If unpublished work is gone, recreate from the contract and rerun QA. | Scratch is not durable storage. Unpublished work may be unrecoverable; do not claim a backup exists. |
| User connection/intervention unavailable | Continue within existing repository authorization; publish coherent WIP checkpoints to proposal branches where possible. Defer dependent questions/live tests and return a precise status/resume point. | No background continuation after the session ends is promised. User silence is not approval. |
| CI fails or timing/network glitches recur | Preserve failed run/artifacts; classify build/test/product/environment cause, reproduce where feasible, fix with evidence and rerun the corrected candidate. | An unexplained rerun passing does not erase the failure or justify weaker assertions. |

GitHub proposal branches/commits are the primary durable checkpoint. Publish coherent WIP before substantial work accumulates when connection allows, but do not merge it before acceptance. Do not publish secrets or raw Azure captures as backup. If no durable publication succeeds, state exactly which local work is at risk and provide a reviewable diff/status; do not promise recovery. Independent source/reference checkouts remain unchanged; use temporary copies for compilation/capture when necessary.

## Rollback and incident correction

Before publication identify the accepted base and whether later work depends on the change. Normal development rollback uses a new reviewed revert commit/PR on current core-v1, not reset/force-push or history deletion. For a merged PR, inspect parents and diff before selecting the correct mainline; the expected operation is `git revert -m 1 <merge-sha>` only after verifying parent 1 is the intended core-v1 baseline. Plain commits use ordinary revert. Test the resulting candidate on current code, review dependent behavior and retain original/revert evidence. Do not automatically revert a historical merge without checking subsequent dependencies.

Pre-merge rejection leaves core-v1 unchanged and preserves the proposal/evidence; remove only task-owned temporary files. A user working tree must not be reset or cleaned. For installed candidate rollback, preserve and reselect a previously verified artifact; do not silently replace executables or delete report evidence. Existing partial-report/replacement guarantees still apply. No Azure rollback is expected from this offline development because Azure mutations are not performed.

If a security defect is found, stop affected acceptance/release claims, contain the faulty path, record a sanitized reproduction and prioritize a fix/revert under protected review. Notify the user in this conversation; do not message third parties or invent a security incident from a hypothetical risk.

## Explicit environment-dependent queue

- DV-001: suitable non-production nested management-group traversal; production-containing Advisory is not a substitute.
- Nonempty live Policy/Defender Recommendations, Arc numeric response shape, restricted-identity visibility, representative load and plugin-specific live comparisons.
- Historical Diagnostics request correlation and Advisor cross-run uncertainty, including the unclassified recommendation-exclusion warning.
- Actual installed Azure operation, fresh OS, hosted maintenance execution, approved release publication/provenance/licensing/operations and controlled production pilot.

Every item stays separately tracked in [DEFERRED_VALIDATION.md](DEFERRED_VALIDATION.md), alignment register, Gate 004/release evidence and the roadmap as applicable. If one becomes critical for a particular feature's acceptance, pause that acceptance and proceed with independent tasks. Only ask for laptop/Azure input when available and materially necessary; no credentials in chat.

## Primary development references and limits

- [Google small-change/review guidance](https://google.github.io/eng-practices/review/developer/small-cls.html): self-contained changes, related tests and separating unrelated refactors.
- [GitHub required status checks](https://docs.github.com/en/pull-requests/how-tos/merge-and-close-pull-requests/troubleshooting-required-status-checks): latest candidate status and head/merge-check distinction; our acceptance additionally requires real native success, not skipped/neutral tests.
- [Git revert](https://git-scm.com/docs/git-revert): corrective history and merge-mainline consequences.
- [SLSA requirements](https://slsa.dev/spec/v1.2/requirements): release provenance is a separate evidence boundary; no achieved SLSA level is claimed.

These controls reduce risk and improve auditability/recovery. They cannot guarantee defect-free code, independent reviewer availability, continuous connectivity, durable unpublished scratch work or unavailable live evidence. Broader authorization should reduce handoffs, not remove acceptance boundaries.

## Ongoing authorization and execution checkpoint

The user confirmed continuation through this plan without routine “continue” handoffs. Proceed with independently reviewable steps under existing QA/protected-merge authority; stop for genuine missing specifications, input-dependent decisions, critical blockers or major flaws. Defer laptop/Azure-dependent acceptance and continue independent offline work. This does not authorize new live scopes, credentials, Azure mutations or weakening gates, and does not promise background operation after the active session ends.

B1 is VERIFIED OFFLINE through PR #78, native run `36947056761`, merged `50fbaee905e7be7a7c99e174f35ea69247eb16bb`; B2 follows [YAML_GRAPH_PLUGINS.md](YAML_GRAPH_PLUGINS.md). A fresh source recovery from that remote merge replaced the lost scratch checkout. Source, pins, mandatory gates and live deferrals were rechecked before implementation.

B2 is VERIFIED OFFLINE through PR #79, merged `7cddf9368272dd425315bdf3635b21fb974a9fe4`, tree `6d209444dbb59ace9570736a33bbd6cc8b0ff818`, native run `36963458729`; required Linux/Windows jobs and exact trees were verified. B3 starts with [ZONE_MAPPING.md](ZONE_MAPPING.md), the bounded adapter contract. It deliberately leaves public plugin availability unchanged until canonical table/health/privacy/report/CLI acceptance. Resume after B3a acceptance at that integration, then verify a second plugin against the shared framework.

B3a is VERIFIED OFFLINE through PR #80, merged `0b76ad6b2870432357e5a38c1e8fe73d5214f11b`, tree `1e6bed9f4f19b4848a6f35a9d1098c4bfccb2370`, run `36966029765`. [PLUGIN_TABLES.md](PLUGIN_TABLES.md) governs B3b canonical/report infrastructure, including additive schema and conservative consumer limits. After its acceptance, resume zone coordinator/CLI/discovery integration, then the second real adapter. Synthetic second-shaped tables do not close that adapter verification.
