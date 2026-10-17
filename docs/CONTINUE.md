# Continue Cloud Assess

Bookmark this page or [the copy-ready prompt](RESTART_PROMPT.md). Read [CURRENT_STATE](CURRENT_STATE.md) for current status and exact next action. This page routes recovery; it does not grant additional authorization or certify QA.

## Start here

1. Read root [AGENTS.md](../AGENTS.md), CURRENT_STATE and the active task PR identified there. Verify live branch/head/tree/parents, protected baseline and actual worktree before editing. Do not select a branch merely because it is newest.
2. Read the applicable contract and [development workflow](DEVELOPMENT_WORKFLOW.md). Read [continuity procedure](CONTINUITY.md) for recovery, checkpointing and evidence rules. Load historical handoffs only when their details are needed.
3. Inspect pending or uncertain actions before retrying. Complete outstanding revision-specific QA and evidence review before feature work. Preserve unrelated work and existing draft heads.
4. Continue only already-authorized independent work to the recorded natural checkpoint. Update durable records proactively without waiting for a handoff request.

## Stable bookmark and successor routing

Initial bookmark branch: `docs/continuity-protocol-20261017`. Its CURRENT_STATE identifies the active task on that lineage. Keep this branch available while it is a recovery bookmark. A successor must inherit these files and publish its recoverable state before the bookmark CURRENT_STATE is updated with an explicit successor branch/PR/verified revision. Re-read the pointer after publication. Update only an authorized unprotected branch; never move protected refs for routing. A pointer-only commit does not transfer QA certification or need repeated evidence-only edits. If forwarding cannot be published, report that the bookmark is stale and provide the verified successor link to the user.

Follow only explicit published successor links; verify their identity, lineage and recorded scope. Detect loops or conflicting pointers. Inspect rather than guess when they disagree with live refs. This scheme depends on maintaining the bookmark; it is not automatic branch discovery.

The existing [2026-10-17 handoff](HANDOFF_20261017.md) remains a historical recovery snapshot. Earlier STOP, pending and workspace statements are snapshots; technical behavioral contracts remain applicable unless explicitly revised.
