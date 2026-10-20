# Continue Cloud Assess

Bookmark this page or [the copy-ready prompt](RESTART_PROMPT.md). Read [CURRENT_STATE](CURRENT_STATE.md) for current status and exact next action. This page routes recovery; it does not grant additional authorization or certify QA.

Current transfer: [2026-10-19 checkpoint](CHAT_TRANSFER_20261019.md), [audit](AUDIT_72H_20261019.md) and [discrepancies](DISCREPANCIES_20261019.md). Follow verified live-PR continuation; preserve old bookmark heads.

Current transfer: [2026-10-20 security crossroads](CHAT_TRANSFER_20261020.md). Read CURRENT_STATE and the [dependency-policy blocker](RETAIL_SECURITY_BLOCKER_20261020.md) before resuming; wait for Denis's decision. Earlier automatic task6 continuation is superseded.

## Start here

1. Read root [AGENTS.md](../AGENTS.md), CURRENT_STATE and the active task PR identified there. Verify live branch/head/tree/parents, protected baseline and actual worktree before editing. Do not select a branch merely because it is newest.
2. Read the applicable contract and [development workflow](DEVELOPMENT_WORKFLOW.md). Read [continuity procedure](CONTINUITY.md) for recovery, checkpointing and evidence rules. Load historical handoffs only when their details are needed.
3. Inspect pending or uncertain actions before retrying. Complete outstanding revision-specific QA and evidence review before feature work. Preserve unrelated work and existing draft heads.
4. Continue only already-authorized independent work to the recorded natural checkpoint. Update durable records proactively without waiting for a handoff request.

## Stable bookmark and successor routing

Initial bookmark branch: `docs/continuity-protocol-20261017`, associated with [PR151](https://github.com/DeBoX85/Cloud-Assess/pull/151). Keep the bookmark branch available. Preserve its verified head while instructed: use an explicit continuation link in the live PR body to reach a successor, rather than committing a pointer-only edit to a frozen head. Read the active branch's CURRENT_STATE and live PR before edits. A successor must publish recoverable state and be read back before forwarding is recorded. Both links must name the intended branch/revision and scope; do not infer the successor from PR number or creation date. Only update bookmark files when that head is authorized to advance. No protected ref update is needed.

The initial bookmark snapshot described CURRENT_STATE forwarding; the current user requires PR150/151 heads preserved. This live-PR successor route resolves that conflict without rewriting their commits. If forwarding cannot be published, report the stale bookmark and give the verified successor link directly. Routing is not QA certification or automatic branch discovery.

Follow only explicit published successor links; verify their identity, lineage and recorded scope. Detect loops or conflicting pointers. Inspect rather than guess when they disagree with live refs. This scheme depends on maintaining the bookmark; it is not automatic branch discovery.

The existing [2026-10-17 handoff](HANDOFF_20261017.md) remains a historical recovery snapshot. Earlier STOP, pending and workspace statements are snapshots; technical behavioral contracts remain applicable unless explicitly revised.
