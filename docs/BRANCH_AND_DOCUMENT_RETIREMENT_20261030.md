# Branch assessment and controlled documentation retirement

30 October 2026, Europe/Berlin. Assessment requested by Denis before adding cleanup recommendations to the handover. Inventory timestamp: 2026-10-30T20:26:09.293071+00:00. Counts describe the snapshot before this assessment successor branch/PR is published; future live counts may differ. No branch, PR or evidence was deleted or closed. This is a disposition proposal, not deletion authority.

## Findings from the live inventory

Cloud-Assess had265 live branches and212 PR records:134merged,72open,6closed without merge. All four branch/PR pagination pages were checked through an empty final page. Every live development tip was fetched into a separate namespace; all265 names/SHAs matched the API snapshot. Full local commit ancestry and tip metadata were inspected; the fetched database passed git fsck. The [per-branch manifest](BRANCH_ASSESSMENT_20261030.json) records each head/tree/date/ancestry, PR relationships and proposed disposition.

| Disjoint disposition | Branches | Evidence and action |
| --- | ---: | --- |
| Keep named active/recovery roles |7| main seed/bookmark, protected baseline, current handover213, reviewed product212, preserved WIP65440, preparation089archive and frozen security185. These roles override age and ancestry |
| Hold open PR references |69| Other branch heads or bases referenced by open PRs. Reconcile explicit successors, unresolved findings and retention duties before any PR closure/branch retirement |
| Clearest retirement-review candidates |134| Each has a merged PR and its current tip is reachable from the protected accepted baseline; no open PR head/base references. Review retention and inbound links individually before proposing deletion |
| Contained but unaccepted historical refs |16| Tips reachable from the latest handover/product lineage, outside protected acceptance and not mapped to an open PR reference. Includes WIP/preparation/handover and six gate-probe refs. Hold until purpose and supersedence are reconciled |
| Hold separate history |39| Outside both accepted and current handover lineages, with no open PR reference in the prior categories. Includes evidence/outage/review snapshots and old experiments; inspect retained content and alternatives before considering retirement |

These categories sum to265. By ancestry alone,136tips are within accepted history,74additional tips are within the latest handover lineage, and55are outside both. Those are overlapping role-independent ancestry facts, not alternative deletion lists. The72open PRs are all drafts in this inventory; many are cumulative/documentary, not72independent unfinished features.

Across all265branches,37 contain54commits reachable through no other live branch head. This count excludes tag/PR/deleted/unreachable ref rescue possibilities and does not imply those commits have no other backup. It does establish why bulk age-based deletion is unjustified. The manifest records the affected branches and commit IDs.

There are259unique tip commits. Two duplicate-tip groups each contain one parent branch plus three deliberately different success/failure/skipped gate-probe branch names. Identical code therefore does not establish identical workflow evidence or a redundant recovery role.149tips have commit timestamps older than14days, but commit dates are not branch-creation, last-push or inactivity dates. No age threshold is used for disposition.

## Concrete examples

Merged/accepted candidates include docs/accepted-quota-region-resume(PR107), docs/advisor-unknown-scope-evidence(PR19), docs/arg-response-evidence(PR36) and docs/cli-process-validation(PR48). Each still needs the link/evidence/restore checks below. They are review examples, not an approved deletion batch.

Separate branches such as docs/retail-qa-outage-20261020, docs/review-resumption-evidence-20261021 and qa/comparison-input-evidence retain branch-exclusive commits in this snapshot. Do not mistake evidence and failed-experiment records for abandoned feature work. Open proposals such as docs/inventory-freeze-readiness-20261026, feat/sql-report-adapter-wip-20261018 and test/region-retail-consumer-20261020 also contain branch-exclusive commits and require explicit technical disposition.

Default main remains a historical seed. Comparing every branch only to main would misclassify much of the project as unmerged/lost. Compare accepted baseline, current product/handover lineage, open proposals and preserved evidence roles separately. Do not merge accumulated drafts to reduce the branch count.

## Private storage is a different category

SC-storage had14branches: shared main plus13intentional archive refs. Every current tip's metadata was freshly read.11archive tips are ancestors of main;2historical archive tips are outside that lineage. Main's parent graph closes using preserved original API metadata plus fresh current-tip metadata. This is ancestry/accounting, not a fresh full private payload audit. Detailed private branch records remain private.

Retain all14for now. They name capture/verification/continuity recovery points and preserve shared CI/security categories under explicit retention. Being old or contained in main does not remove those duties. No private archive deletion, retention-policy change or independent-backup closure is proposed by this assessment.

## Documentation cleanup protocol

The previous SOP already said to archive dated prompts, but it had no complete retirement inventory or consistent supersedence mechanism. PR213 rewrote three current entrypoints and added recovery controls. Adopt the following retirement protocol within PROCESS-STATE/PROCESS-EVIDENCE; this is not a new product scope or closure of those obligations.

1. Keep a small active entry set: CONTINUE, CURRENT_STATE, RESTART_PROMPT, SESSION_HANDOVER authority map, PROJECT_PROGRESS and RECOVERY_STATE. Keep behavioral contracts, acceptance records and evidence separately identified. Current status should have one authoritative home plus exact-revision PR outcomes.
2. Inventory candidate documents and classify active-state, active-contract, historical-handover, evidence, or superseded-duplicate. For each record: exact path/revision/hash, date/scope, replacement, still-binding requirements, unresolved IDs, inbound links, destination and retention/disposition. Age or filename alone is insufficient.
3. Extract any unique binding requirement or unresolved obligation into its correct current contract/register before retiring status prose. Preserve original text at an immutable revision or approved archive location. Never archive a live technical contract merely because it is old.
4. Move confirmed historical handovers out of active orientation only in small reviewed batches. Repair internal links and leave short successor redirects at important old bookmark paths. Mark archived pages with their scope/date/replacement. Do not append new instructions below old current-state paragraphs.
5. Check local links, obligation IDs, original controls, live branch roles and required authority headers. Run scripts/check-handover.py and the existing unchanged guards; their structural pass is not semantic certification. Frozen evidence consumers and bound-source hashes must not be moved or edited as a side effect.
6. Publish/read back each batch and follow its restart path from published records. Observe exact-head QA and private evidence requirements, record rollback and unresolved findings, and reconcile the current/progress/register views. Do not delete raw evidence or rewrite Git history.

## Branch retirement protocol, proposal only

After PR213's complete-log/current-AST/control/source review, final private CI capture/restore and self-assessment, prepare a dry-run batch of at most10 of the134 merged/accepted candidates. Record each branch's exact expected head, merged PR and containment proof, every open PR head/base check, current/durable bookmark role, inbound documentation/workflow references, evidence location and retention disposition, preserved immutable recovery ref or independent copy, and how its ref would be recreated. Recheck live state immediately before any action; skip a changed head instead of retrying blindly.

Present the concrete list and restoration plan for explicit Denis deletion authorization before deleting any branch. This conversation authorized assessment and handover updates, not branch deletion or mass PR closure. For an approved batch, verify remaining protected/active/evidence refs and restart routing after each bounded operation. Keep an immutable disposition receipt. Recreating a ref is not a guarantee that deleted CI/PR metadata can be restored.

Open draft PR retirement is a separate technical-disposition task. Map cumulative successors and unresolved review threads; retain failures and contract freezes. Close only an explicitly superseded proposal with an accurate successor/evidence link and actual authorization. Do not merge it just to make cleanup easier. Keep the16unaccepted-contained and39separate-history groups out of the first deletion proposal.

Longer term, use clear active/preparation/evidence roles, concise current entrypoints and explicit successor routing. Prefer durable exact revisions or recorded archive refs for historic links. Routine new WIP should advance a coherent authorized task branch where frozen-head rules permit; preserve a new successor only when materially required, rather than generating unnecessary status-only branches. Never rewrite a frozen/verified head to satisfy tidiness.

## Automation and primary sources

Automate read-only inventory, link checking and dry-run classification first. Do not enable age-based deletion. GitHub's [official branch deletion/restoration guidance](https://github.com/github/docs/blob/main/content/repositories/configuring-branches-and-merges-in-your-repository/managing-branches-in-your-repository/deleting-and-restoring-branches-in-a-pull-request.md), retrieved30October2026, describes deletion from closed/merged PRs when no other open PR references exist, and restoration of a closed PR's head branch. This is not a durable backup promise.

GitHub also supports [automatic post-merge head deletion](https://github.com/github/docs/blob/main/content/repositories/configuring-branches-and-merges-in-your-repository/configuring-pull-request-merges/managing-the-automatic-deletion-of-branches.md), subject to administration/rules. Actual settings were not inspected or changed. Reconsider that only after recovery/archive-role protections and historical-link migration are deliberate. It would not resolve72currently open draft proposals or protect unique unmerged evidence. No new scheduler, token, secret, ruleset or deletion automation is configured here.

## Exact handover and limits

PR213's own quality41106459423 and source41106459419 now report SUCCESS. Their final complete-log/current-AST/control/source review and final CI retention are still pending; metadata green is not the missing review. PR212/b0 remains the last fully reviewed product checkpoint. This new material assessment/documentation successor preserves both heads, has its own checks and QA status in its live PR, and inherits no green. First finish the recorded pending review and evidence work, then prepare the documentation inventory/branch dry-run proposal. Coordinator accounting/contracts remain the next product work under enhanced verification; no task10closure or11–15progress.

This assessment did not perform a full per-branch semantic, inbound-link or evidence-retention audit. A retirement candidate is not certified safe to delete. No root cause, original134loststreams, independent backup or global security/acceptance assurance is established. Preserve49obligations/44legacy IDs/32original frozen inputs, pins and existing guards. No Azure/live retail/real credentials/fixtures/RBAC/subagents/unrelated refresh/protected merge/bypass/release/branch deletion/mass PR closure.
