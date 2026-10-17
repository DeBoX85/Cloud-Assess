# Rollback and restoration runbook

[CURRENT_STATE](CURRENT_STATE.md) owns authorization/status. This runbook expands [the existing development procedure](DEVELOPMENT_EXECUTION_PLAN.md#rollback-and-incident-correction); it does not authorize protected merges, Azure operations or release. [Action/change records](ACTION_REGISTER.md) and [retained rehearsal evidence](EVIDENCE_RETENTION.md) provide traceability.

## Select the correction boundary

| Situation | Route and required checks |
| --- | --- |
| Unmerged proposal | Protected baseline is unchanged. Preserve its branch and evidence; stop adopting it or prepare a successor corrective commit. An ancestor proposal can still be included by a later stacked PR: inspect the complete ancestry/diff before choosing the replacement. |
| Individual commit | Inspect commit/parent, affected files and subsequent dependents; make a new targeted revert or forward correction on a task branch based on the actual current target. A clean Git application alone is not behavioral proof. |
| Merge commit | Inspect every parent and the merge diff. Choose mainline only after proving which parent represents the retained baseline. Squashed ordinary commits use ordinary revert, not a merge-mainline flag. |
| Dependent changes/conflicts | Identify helpers, tests, CI inventory/generated files, contracts and schema/data consumers that depend on the change. Revert a reviewed group in dependency order or use a forward fix. Do not choose ours/theirs wholesale or erase unrelated work. |
| Installed candidate | Retain and deliberately select a previously verified artifact with source/version/manifest/hash and reviewed compatibility/security. Preserve current executable and report evidence. No silent executable replacement. Installed/fresh-OS rollback remains unverified release work. |

## Execute only in a task-owned isolated checkout

1. Record incident/change ID, reason, affected revision, intended before/after behavior, authorization and safe containment. Inventory local work and pending remote operations. Confirm current baseline/proposal refs and exact rollback target. Do not reset/clean a user worktree.
2. Publish current recoverable WIP where authorized; keep original commit/PR/evidence reachable. Create an isolated task branch/checkout from the intended current target, not a remembered historical path. Record branch, starting HEAD/tree/parents and known-good comparison.
3. Inspect the proposed commit(s), merge parents and downstream consumers. For an ordinary commit the operation is `git revert <commit-sha>`. For a genuine merge, `git revert -m 1 <merge-sha>` is appropriate only when parent1 is independently verified as the intended mainline. These are templates, not instructions to execute against protected refs.
4. If Git conflicts, preserve diagnostics and inspect changes. In this task-owned clean-start rehearsal/candidate, `git revert --abort` should restore the recorded starting state; verify HEAD, tracked/untracked state and file content. Never blindly abort another person's ongoing operation. Resolve only after dependency review and recheck the entire resulting diff.
5. Verify the intended behavior, retained unrelated features, regression controls, reporting/privacy/security and executable/package compatibility. Run focused tests plus mandatory full exact-head Linux/Windows/source QA. Reverting a correction may intentionally restore a known defect and red tests; keep those expectations and do not call that an acceptable candidate. Prefer a forward fix or explicitly keep the candidate unaccepted.
6. Publish/read back commit/tree/parents/identity and all changed contents; record original/revert/restore links and evidence. Protected acceptance, when separately authorized, follows normal review/gates and distinct accepted-push verification. No force push/reset/bypass.
7. To reintroduce reverted work, inspect current dependencies and merge topology first. A reviewed revert-of-revert or new corrected patch may be needed; merely merging the same branch again may omit previously reverted changes. Repeat behavioral/full QA rather than assuming reversibility implies correctness.

## Local rehearsal completed on 2026-10-17

- Actual revision PR151/ba9eeb9aade68144bdf1c0d51e108704659a1242 was reversed in a disposable worktree. Tree f722248deee6065fa81e28f81658a614837f2c5d returned exactly to parent tree fd649a6d461c2ec84ea505222379133bae0ec507.
- Revert d9a2630238fd145209ecbe86ae8573fb9d073fcf was itself reverted as 0dbbc5aa3090469366fc3689e30a06e2932e3693, restoring the exact original tree. Static preflight passed on the restored checkout. Local marker branch preserves that chain; the durable command record is in the evidence bundle.
- A separate synthetic repository verified two-parent merge selection, mainline reversal/restoration and deliberate dependent-edit conflict. Abort restored the exact dependent HEAD/content with a clean tree. No product branch was mutated by that exercise.

This proves local Git/tree/static behavior for the documentation revision and the synthetic topology. It does not prove production-code behavior after rollback, Windows Git rollback, installed-package recovery, protected acceptance or live operations. DC-012 and DC-010 retain those obligations.

## Sources and limits

[Git revert documentation](https://git-scm.com/docs/git-revert) defines corrective commits and merge-mainline consequences; [faulty merge guidance](https://www.kernel.org/pub/software/scm/git/docs/howto/revert-a-faulty-merge.html) explains later reintroduction. Contract-specific rollback limitations remain active. Hashes and successful static checks do not establish a safe release. No Azure resource rollback is expected from this offline work because it performs no Azure operations.
