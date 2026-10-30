## Branch and document retirement assessment addendum, 30 October 2026

[The live inventory and controlled retirement protocol](BRANCH_AND_DOCUMENT_RETIREMENT_20261030.md) and [per-branch manifest](BRANCH_ASSESSMENT_20261030.json) now qualify the cleanup recommendation:265development branches at collection,134merged/accepted review candidates,69otheropenPRreferences,39separate-historyholds,16containedunacceptedrefs and7named roles.37branches have54commits not reachable through other live branch heads. All14private storage branches are held under evidence retention. No deletion, PR closure or automation enabled. First finish PR213 final review/CI retention and this successor's own checks; then prepare small dry-run disposition batches. Preserve historical contracts/evidence, repair links, require explicit branch-deletion authorization and verify restoration. Source of truth remains CURRENT_STATE/RECOVERY_STATE/live exact-head PR.


# Handover process review after disaster recovery

30 October 2026, Europe/Berlin. Scope: reduce reconstruction work and detect drift earlier without weakening QA, inventing automated persistence or changing protected infrastructure. Root cause of the workspace rollback is unknown; the controls below address observed loss and confusion, not a proven platform cause.

## Assessment

The existing SOP already required remote WIP publication, full readback, private evidence capture, expected-head updates and20–30active-minute checkpoints. Those rules helped preserve the product and much evidence. They were insufficient in practice: long verification batches produced local-only raw evidence; several documents repeated changing state; branch names, preparation green and implementation completion were easy to conflate; handoff guards checked routing/IDs, not semantic consistency. More repeated narrative would increase that last problem.

| Observed weakness | Implemented mitigation | Residual risk / validation |
| --- | --- | --- |
| Contradictory growing entrypoints | Rewrite CONTINUE, CURRENT_STATE and RESTART_PROMPT as concise orientation; archive prior text via immutable b0 links, preserve applicable contracts | Source-backed review is still needed; links/hashes do not prove semantics |
| Ambiguous revision roles | RECOVERY_STATE.json distinguishes accepted, verified product, WIP, preparation, security review and this handover proposal | Local consistency checker cannot observe live changes; fresh remote readback remains mandatory |
| Local-only evidence survives until end of long QA | Phase gates require completed evidence batches to be durably captured/read back/restored before moving to the next task; explicit at-risk inventory before long work/status handoff | Process enforcement during active execution, not unattended uploader; an interrupted atomic batch can still be lost |
| Preparation/source/byte proof mistaken for runtime/acceptance | Explicit state vocabulary and revision-bound QA/retention fields; no inherited green, no status-only self-certification commits | Final complete native/source checks remain mandatory |
| A new chat trusts stale context | Published-only orientation and enhanced verification for at least next3completed material slices | Same-agent walkthrough is not independent review or full outage simulation |
| Uncertain writes retried or shared storage overwritten | Before/after operation journal with expected parent; inspect live state after uncertain response; additive immutable categories and lease/non-force updates | Journal itself must be durably published before sensitive transition; no platform transaction guarantee |
| Lost originals recreated as originals | Explicit loss/provenance fields and immutable partial/final captures | Replacement execution can validate behavior but cannot restore historical authenticity |
| Structural handoff tests miss meaning | New local check-handover.py checks revision relationships, registry counts and current entrypoint markers; adversarial negative fixtures required for this tool | Additional check, not a replacement for existing seven guards/33workflow guards/full gates; no automatic CI integration claimed |

## Small operating loop

1. Orient from published records, verify live refs and inspect actual worktree/uncertain actions.
2. Record one bounded objective, source/contract claims, independent expectations, expected command exits and natural stop. Publish code/contract WIP before long execution.
3. Execute one coherent slice. Preserve full command streams, including failures. For long suites, retain completed, closed evidence batches at safe boundaries; do not archive a live file prefix as a complete log or disrupt an atomic mutation/restoration to meet a timer.
4. Publish evidence to its authorized private category, verify all remote bytes and restore into a new directory. If final QA is incomplete, publish a partial capture with explicit missing members and healthy_quality_accepted:false. Partial capture is recoverability, never acceptance.
5. Review full own-head native/source QA and frozen controls, preserve rollback/forward evidence, publish final capture/receipt and update the live PR index. Only then move to the next dependent development task.
6. Reconcile compact state/progress/register and changed procedures at a material checkpoint. Read the restart path as a fresh reader using only published data. Status replies must distinguish local-only, published-unverified and remotely verified results.

The20–30minute safe-boundary target remains a secondary reminder. The primary boundary is before expanding work, long QA/setup, context-heavy review, task switch or session end. If durable storage fails, keep local work, explicitly record what is at risk, and avoid expanding an irreplaceable-evidence backlog. Read-only research or reproducible work can continue where independent; no invented user approval is required.

## Source of truth and supersedence

CURRENT_STATE owns concise human state; RECOVERY_STATE records exact role/evidence identities; the live task PR finalizes its own exact-head outcome. These are complementary views that must agree. Historical documents retain their date/scope and behavioral requirements; newer timestamps alone do not override a frozen contract. Resolve contradictions through exact source, independent expectations and revision-specific evidence. Record the superseded claim and reason once, with links from other entrypoints. Never interpret an unreviewed historical paragraph as automatic authority.

No self-referential commit hash belongs in its own tree. Publish intended state with pending status, then record actual commit/tree/parents/identities/readback/QA in the live PR and private receipt. Carry final outcomes into the next material checkpoint. Preserve frozen PR212 as a verified product base while this handover proposal has separate pending checks.

## Recommended independent-backup work, still open

SC-storage protects against local workspace loss but shares an account/provider failure domain with Cloud-Assess. An owner-controlled independent copy should eventually cover both repositories' history/refs, required submodule/LFS objects, PR/review/issue metadata, workflow/run/job metadata, raw evidence, private receipts and recovery tools. Test restores in a separate directory/environment and document actual recovery-point/time measurements. A mirror alone does not back up all GitHub metadata; a hash-only manifest cannot reconstruct missing bytes. Privacy, destination access, retention and restore credentials need a concrete owner decision before provisioning anything.

No independent destination, scheduled backup, workflow permission, new secret, branch protection or external messaging is configured by this change. CI evidence automation remains a separately tracked design: least-privilege cross-repository writer, private payload review, completed-log detection, resumable immutable batches, idempotent expected-head writes, failure reporting and offline restore verification must precede deployment. Do not silently add a broad token or a workflow that exposes private logs publicly.

## Primary-source check

GitHub's official [repository backup guidance](https://github.com/github/docs/blob/main/content/repositories/archiving-a-github-repository/backing-up-a-repository.md), retrieved30October2026, distinguishes mirror history backups, separate LFS retrieval and selected-metadata archives, with restoration limitations. Its official [retention policy source](https://github.com/github/docs/blob/main/data/reusables/actions/about-artifact-log-retention.md), retrieved the same day, describes default90day retention and configurable limits. Actual repository settings were not inspected or changed. These sources support retaining full data outside ephemeral Actions links and separating Git objects from metadata backups; they do not certify this project's backup completeness.

## Completion and limits of this improvement

Implemented here: concise entrypoints, detailed incident/recovery/supersedence handover, structured role/evidence state, enhanced next-slice protocol, phase-based durable-evidence conditions, local consistency checker and explicit open backup/automation work. Existing pins/frozen expectations/gates remain unchanged. Publication/readback/local checks and this proposal's hosted QA status are recorded in its live PR; no acceptance is inferred here.

These changes reduce ambiguity and the amount of work exposed between checkpoints. They cannot guarantee zero loss, save after a hard outage, repair the underlying platform, automatically start another chat or make independent-human/security/operational acceptance unnecessary.
