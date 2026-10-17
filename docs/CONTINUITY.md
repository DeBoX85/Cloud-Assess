# Proactive continuity and interruption recovery

[CURRENT_STATE](CURRENT_STATE.md) owns current task state. [CONTINUE](CONTINUE.md) is the stable entry point; [RESTART_PROMPT](RESTART_PROMPT.md) is the copy-ready user prompt. This procedure supplements [DEVELOPMENT_WORKFLOW](DEVELOPMENT_WORKFLOW.md). Explicit current user constraints prevail.

## Ownership and checkpoint triggers

The working agent owns handoff maintenance. Do not wait for the user to request it. Persist a task plan before a bounded slice and publish coherent task-owned code and notes as work progresses. Use these triggers:

| Trigger | Required record/action |
| --- | --- |
| Before a bounded task | Objective, scope, base/proposal revision, applicable contract, expected verification and natural checkpoint |
| Material finding, decision, fix or changed next action | Facts versus hypotheses, rationale, affected files, remaining uncertainty and ledger/failure note where relevant |
| Coherent implementation slice | Publish recoverable WIP and notes to an authorized task branch; inspect diff and remotely verify content/tree/parents/identity |
| Long uninterrupted work | Aim for a durable checkpoint every 20-30 active minutes at a safe boundary; record blockers if publication is unavailable |
| Before long QA/environment work | Publish candidate first; record reproducible command/setup, revision, run ID/link and pending status as soon as available |
| After remote mutation | Read actual state back; if interrupted before confirmation, treat outcome as uncertain and inspect before retry |
| Before switching tasks, ending a response/session or planned transition | Reconcile CURRENT_STATE, task evidence and changed procedures; give the stable prompt link without being asked |

The time target is best effort during active execution, not a timer, background service, token meter or outage guarantee. Do not interrupt a non-atomic operation just to meet it. Prefer coherent batches over every-edit commits; unfinished WIP is recoverable but unverified. Existing full exact-head QA remains mandatory. Do not invent a percentage threshold for chat capacity. When reliable remaining-context information is available and low, checkpoint before another large task; otherwise use bounded slices and concise records. Split large reads/log reviews into indexed batches, preserving coverage so omitted content is not mistaken for reviewed evidence.

## Compact recovery record

Maintain the following fields in CURRENT_STATE or link to the active contract/PR when that is the authoritative home. Do not duplicate long narratives across handoffs.

- Updated timestamp with timezone; active objective, branch/PR and natural stop.
- Published recovery revision and independently verified baseline, kept distinct. Record tree/parents/identity in the PR evidence index. Avoid self-referential commit hashes in the commit being created; publication readback/evidence belongs in the PR body.
- Completed work, unfinished work and exact next action. Distinguish blocked, deferred, pending, failed and verified outcomes.
- Pending/uncertain operations: intended action, target, known result, evidence and how to inspect safely. Record none only after checking.
- Verification: exact revision, command/run IDs, status, reviewed evidence, remaining checks and expected intentional failures.
- Environment: actual inventory, pinned setup recipe, task-owned paths and reconstruction steps; local caches/paths are hints, never prerequisites for recovery.
- Evidence locations, retention limitations and any local-only work at risk. Hashes verify integrity but cannot recreate missing files. Do not commit credentials or private raw Azure evidence.
- Authorization, exclusions, protected/source pins, rollback route and applicable technical contracts, using links where possible.

## Publication and evidence

Use an isolated worktree and authorized unprotected task branch. Preserve prior draft heads when instructed by creating a successor branch. Publish code plus sufficient recovery notes before destructive local cleanup, context-heavy work or lengthy verification. Verify the remote commit, tree, parent, human identity and full changed contents. Inspect actual state after an uncertain response; do not blindly replay writes. If publication fails, say what is local-only and at risk and prioritize recovery publication when access returns. A local commit, chat summary, status message or untracked file is not a remote backup.

The PR body is the final revision-specific QA index. Retain complete required logs/captures in an approved durable evidence location with hashes and run/job identifiers. GitHub run links may expire or lose artifacts; check retention before relying on them for irreplaceable evidence. Preserve sensitive evidence outside public Git. Record expired/unavailable evidence honestly and rerun reproducible checks where necessary. Do not claim a retained hash proves unavailable contents were reviewed.

Do not create a commit solely to replace its own pending QA text with its pass and thereby start another certification cycle. Finalize outcomes in the PR body, then index them in the next material checkpoint. Every changed candidate still needs its own QA; no documentation exception weakens the existing gates.

## Recovery after lockup, outage or new chat

1. Read CONTINUE, AGENTS and CURRENT_STATE from the bookmark and any verified explicit successor. Fetch live task PR/ref state and applicable contracts before mutation. Old chat memory is supplementary.
2. Inventory actual local paths, worktrees, dirty changes, tools and running processes. Preserve unrelated work; recreate task tooling from pinned recipes when necessary. Do not reset work or refresh dependencies to make setup easier.
3. Reconcile published revision/tree/parents/identity with the workspace. Inspect pending workflow runs and uncertain operations. Classify completed, failed, cancelled, pending and unknown; a green predecessor does not certify a successor.
4. Recover only missing work from the last remotely verified checkpoint. Inspect and classify failures before rerunning. Report unrecoverable local-only work; do not invent continuity or replay completed work blindly.
5. State recovered status and exact next action, then continue within existing authorization. Update CURRENT_STATE and the task evidence index at the next material checkpoint.

## Transfer and limitations

At natural checkpoints give the user the stable restart-prompt link. Maintain explicit successor routing as described in CONTINUE before moving work to another branch. When the bookmark head must stay frozen, update its live PR continuation link after successor readback; do not move that head merely for routing. Archive dated prompts as history; do not append new current instructions to a growing historical prompt. Historical evidence remains intact.

No process guarantees zero loss if the service fails before publication. The recovery point is the last successfully published and read-back checkpoint; work since then may need reconstruction. The agent cannot save after an outage, promise unattended progress, detect every hard chat limit or initiate a replacement chat. No scheduled automation, external messaging or infrastructure is needed for this procedure. The user normally only opens a new chat and pastes the saved prompt.

## Review at each natural checkpoint

Check that the bookmark reaches the active task, next action is executable, uncertain operations are explicit, evidence statuses are revision-specific, changed local files are accounted for and current constraints agree. Automated handoff routing guards are necessary but do not prove semantic accuracy. A useful recovery check is to follow the restart prompt using only published records and confirm the same next action without consulting chat history; report it as a documentation walkthrough, not an independent-person review or real outage test.
