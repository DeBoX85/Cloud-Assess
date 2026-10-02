# Repository agent continuity and QA

Before changing this repository, read `docs/SESSION_HANDOVER.md`, its current-state/authority map and the linked development execution/QA plans. Verify live refs and open proposals; a historical hash or old workspace path is not necessarily current. Explicit user instructions take precedence over repository guidance.

Maintain the handover proactively during active work: update it after each coherent implementation slice, material finding/decision, accepted merge, changed resume point, and before switching tasks or ending a session. Keep task state, verified baseline, uncertainty, exact next action and evidence links accurate. Update the development ledger, roadmap and failure notes where material. Preserve historical gate snapshots; do not relabel missing/live evidence as passed.

Publish task-owned WIP code and documentation checkpoints to an authorized remote branch when available, and verify commit/tree/parent/identity by read-back. WIP is unaccepted and stays unmerged until mandatory exact-head QA. A local commit/design note is not a durable code backup. If publication is unavailable, report what is at risk; do not promise background work or outage-proof recovery.

After interruption, inspect actual state before retrying an uncertain edit or remote mutation. Reconstruct missing work from the last verified remote contract/base and rerun checks. Use actual path inventory and separate explicit source/target workdirs. Never reset unrelated user work, force protected history, alter the pinned reference or weaken required checks.

Use the human Git identity recorded in the handover for development; maintenance bot rules remain separate. Azure-dependent validation requires its recorded scope/access and remains deferred while unavailable. Do not substitute production-containing management groups or create Azure fixtures without explicit authorization. Continue authorized independent tasks until genuine input is needed.

Resolve uncertain technical contracts by inspecting current/pinned code and primary documentation. Distinguish tested facts, deliberate corrections and hypotheses. Record confirmed mistakes with recurrence/prevention; do not duplicate implementation-mirroring tests or invent success/percentage/ETA claims. Keep updates concise and communicate a genuine blocker precisely.
