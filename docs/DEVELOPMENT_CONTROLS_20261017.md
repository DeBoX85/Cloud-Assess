# Bounded development-controls contract, 2026-10-17

[CURRENT_STATE](CURRENT_STATE.md) owns active status. The user authorized continuation of the previously outlined controls work after PR150/151 exact-head review. Parent is verified draft PR151/ba9eeb9aade68144bdf1c0d51e108704659a1242; protected acceptance is unchanged. This contract is published before new evidence-verifier code.

## Scope and acceptance

- Add one compact [action register](ACTION_REGISTER.md), retaining stable IDs, state, dependencies, responsible party, next action and closure evidence. Existing finding IDs remain intact; do not renumber historic findings.
- Formalize a rollback runbook and change-record template, including stacked proposal dependencies, merge-parent selection, conflict abort, restoration, installed-artifact limitations and current protected-review rules.
- Preserve completed PR150/151 and predecessor comparison evidence beyond expiring CI: complete decoded logs, raw metadata, source reconstruction and rollback command record. Use a self-contained compressed text bundle plus a member-hash manifest. Only reviewed synthetic CI/local rehearsal evidence is eligible for this public repository; private Azure evidence remains excluded.
- Add a bounded, offline, read-only evidence verifier: reject malformed/duplicate JSON, wrong archive/decoded/member hashes or lengths, duplicate/unsafe member names and oversized expansion. It must never extract paths, execute archived commands, call a network or claim publisher authentication.
- Freeze independent test expectations before writing the verifier: one valid two-member bundle must pass; altered archive, altered declared member, missing/extra/duplicate member, wrong count/length/hash, traversal/absolute/Windows path, duplicate JSON keys and oversized gzip expansion must reject with nonzero actual CLI status. Validate the committed real bundle independently as well. Add the test through the existing additive native CI inventory; preserve every previous control.
- Assess current runtime logging against implementation/tests without changing product schemas or adding telemetry. Record existing stage timestamps, codes, completeness, output/exit behavior, privacy boundaries and unproven run/request correlation. Defer product changes to a separately justified contract.
- Verify generated workflows/static preflight, artifact integrity, negative controls, local documentation links and recovery routing, then full own-head native/source QA. The active PR body is the final exact-head evidence index. No predecessor green certifies this successor.

## Rehearsal and evidence limits

A local disposable-worktree rehearsal of PR151 already demonstrated exact candidate-to-parent tree reversal, restoration and restored static preflight. A synthetic repository separately demonstrated two-parent merge-mainline reversal/restoration and dependent-edit conflict/abort without loss. Full command results will be retained in the evidence bundle. These are local Git/static proofs, not native production-code rollback, fresh-OS, protected merge or installed-artifact rollback approval.

Run/source log hashes refer to UTF-8 decoded text returned by the connector, not the original transport archive bytes. Record exact reviewed revisions and native merge previews. A hash does not recover missing content or authenticate a publisher. Git history preserves versions while retained, but this is not an immutable compliance archive; long-term off-repository backup/ownership remains an explicit obligation.

## Rollback, constraints and next boundary

Use a new targeted revert after dependent-change review. New files and CI additions are coupled: revert inventory/generated workflow references with a removed verifier/test. Retain original evidence and do not alter source captures to make a rollback pass. No Azure, fixtures/RBAC, dependency/action/SDK/source-pin/capture/public-schema refresh, subagents, protected merge/bypass or release. No blanket runtime logger is introduced. After completed controls QA, resume only the bounded VMSS catalog-only contract/independent expectations; VMSS scoring stays unsupported.
