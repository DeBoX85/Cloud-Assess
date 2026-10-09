# Cloud-Assess restart prompt, 2026-10-09

Copy the text below into the new chat. It is intentionally detailed. All revisions
are discovery pointers to verify live, not assumptions that nothing has changed.

---

Resume development of Cloud-Assess from our verified handover. Do the recovery
checks first, then continue authorized independent work to the next natural
re-evaluation point. Do not ask for confirmation already given in this prompt.

Repository: https://github.com/DeBoX85/Cloud-Assess
Accepted development branch: `bootstrap/core-v1`.
Goal: functional equivalence to the pinned AZQR reference with independently
adjustable branding, accurate behavior, read-oriented Azure operations and strong
QA. Deliberate source corrections must be documented and independently tested.
No simplified replacement for the source's region scoring or operator features.

## First actions and authority

1. Read the live root `AGENTS.md` before changing anything.
2. Locate the latest accepted and proposed checkpoint through live branch refs,
   open/merged PRs and their bodies. The handover branch is
   `docs/handover-20261009`, PR132:
   https://github.com/DeBoX85/Cloud-Assess/pull/132
   Verify its current head and acceptance evidence rather than assuming it is
   merged. If later work exists, reconcile it.
3. Read the latest `docs/SESSION_HANDOVER.md`,
   `docs/HANDOVER_CHECKPOINT_20261009.md`, `docs/AUDIT_RECOVERY.md`,
   `docs/PROJECT_AUDIT_20261004.md`, `docs/PROJECT_AUDIT_REPORT_20261004.md`,
   `docs/AUDIT_REVIEW_COVERAGE_20261004.md` and all relevant linked records.
4. Reconcile `docs/TARGET_SPECIFICATION.md`, `docs/IMPLEMENTATION_PLAN.md`,
   `docs/ROADMAP.md`, `docs/DEVELOPMENT_EXECUTION_PLAN.md`, `docs/QA_PROCESS.md`,
   `docs/DEVELOPMENT_LEDGER.md`, `docs/FAILURE_NOTES.md`,
   `docs/DEFERRED_VALIDATION.md` and current `REGION_*` contracts. New current-state
   headers supersede old pending/resume wording; preserve historical evidence.
5. Inventory the actual executor, tools, worktrees, tracked/untracked files,
   materialized submodules and source copies. Old scratch paths are hints only.
6. Verify live refs, exact commits/trees/ordered parents, Denis author identity,
   source pins, active protections and actual executed-head native/source CI.
   Before retrying an interrupted commit/ref update/push/PR creation/merge, inspect
   remote state to determine whether it already succeeded. Response loss means
   unknown outcome, not failed mutation.

Start with a concise recovery report: current accepted baseline, active proposals,
last completed phase, unresolved findings/evidence, any lost or unpublished work,
and exact next action. Then resume without unnecessarily repeating verified work.

## Last verified implementation and evidence

Before the documentation handover, the accepted implementation was merged PR131:
https://github.com/DeBoX85/Cloud-Assess/pull/131

- Accepted commit: `b653a3abfc35590a185531095a168a9a107e769e`.
- Tree: `47515c474ddf400da0a68c556e9147d8d3c2576b`.
- Ordered parents: `a79464db59fa0b3e72115d2f07a18f8621bd437f`, then
  `4b5a088835f2b97715347929d865574b94cefd3d`.
- Denis Bogunic author: `134433647+DeBoX85@users.noreply.github.com`.
  Merge committer GitHub; verified/valid signature. Use Denis identity for
  development commits; existing maintenance-bot rules remain separate.
- No open proposals were observed before creating the documentation handover.
- Ruleset `23890737`: active, no bypass, mandatory `quality` and
  `windows-validation`, integration `15368`. Do not weaken any mandatory gate.

PR131's final body is its complete historical accepted evidence index:

| Evidence | Run | Jobs | Executed object |
|---|---|---|---|
| Final candidate native | `37730322104` | Linux `113158076872`, Windows `113158077129` | Preview `3c30675a2727a63eca2b52b29e904a9a94c83688` |
| Final candidate source | `37730322119` | `113157725843` | Same preview |
| Distinct accepted push native | `37732157523` | Linux `113163501776`, Windows `113163501861` | Accepted `b653a3abfc35590a185531095a168a9a107e769e` |
| Distinct accepted push source | `37732157557` | `113163501600` | Same accepted commit |

All succeeded. Metadata candidate head `4b5a0888` is authored, not the actual PR
preview checkout. Preview tree equals final implementation tree and ordered
parents are accepted `a79464db` / authored `4b5a0888`. Push logs separately prove
actual accepted branch checkout. Full logs were inspected at acceptance: 53 Linux
steps (only conditional failure-evidence retention skipped), 38 Windows steps all
successful; all 24 reservation collector, 25 quota collector and 18 reservation
runtime compiling fault/named assertion/restored controls on BOTH hosts. Full Go
suite, Linux race/vet, 84.5% coverage, bounded fuzz, real CLI/help/preflight/module,
branding/docs/PowerShell/package/isolated-install checks passed. Do not describe
those as Windows race, live Azure, fresh OS, representative estate load or release
approval. Zero reachable/zero imported vulnerability findings still leave one
module-only advisory open.

Pins unchanged: AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`, source tree
`17d93b20c303f90f7843036be82f0dc32f3260f1`, APRL
`60eaddda76541f6adbc1c5ffa686829807e55e29`, actual development APRL tree
`3ea4d285f0fb00e0a90edf4e5d186c90c4c6296d`. Compute SDK remains
`armcompute/v6 v6.4.0`, API `2024-11-01`. Candidate and accepted source proof
separately checked eight unchanged harnesses, 17 capture records and 150 unique
contiguous chunks against independent SHA256 and full retained remote bytes.
The 35 reservation source scenarios are not 35 target collector fixtures.

## Audit and completed development

The comprehensive whole-project offline audit was completed and accepted through
PR114, hardening merge `c33fb2eee93a2e5730f72612b351d0b838737d03`, and PR119,
disposition merge `3463a72f11912039702ee60605ec4d79830c09bc`. Do not restart it or
independently merge superseded PR113/115-118. The frozen first-party inventory
covered all 216 Go files and audit corrections brought it to 222 Go files,
32 scripts and four workflows. Requirements, docs, data/provenance, operational
tooling and independent failure/integration tests were included. This is offline
self-review, not independent-person approval, a line-by-line third-party audit,
all-input certification or live/release acceptance. Later additions have separate
bounded review records. Do not relabel historical audit counts as a fresh review
of every file in the current repository.

AUD001/002/003/005/006/007 were resolved. AUD004 records a lost unpublished
availability draft, later reconstructed and accepted through PR120; no accepted
code loss is known. Full Cost pagination remains unimplemented: fail-closed
handling of unconsumed continuation is accepted, not full pagination parity.

Accepted post-audit region slices through PR131: owned availability PR120;
latency PR121; weighted-cost characterization/runtime PR122/123; REST/VM quota
and reservation characterization PR124-126; owned calculations PR127/128;
REST quota collection PR129; VM quota collection PR130; reservation collection
PR131. Their exact evidence lives in their PRs and governing contract documents.
Internal components do not make public `region-selection` executable.

PR131 implements bounded group/list/expanded-Get collection using existing
transport, explicit presence/unknown/partial health, safe continuations, validated
scope/identities, cancellation and owned state. Bounds are 256 getter calls,
64 pages per chain, 1 MiB per response, 8 MiB aggregate successful bytes and
8192 aggregate raw work. Preserve strict decode/Unicode/duplicate/budget guards.

Resolved FN096-099: returned-ID structure before equality; Unicode duplicate
identity keys including internal runtime; shared JSON alias overwrite admission;
physical/ASCII region contract. Keep all regressions and compiling controls.
FN094 executor disconnect/path/cache issues and FN095 fixture/compilation/status
coverage mistakes are process/environment/test-setup records. Failed reproductions
remain failed history. No live incident, product-wide instability or missing
independent-person approval is inferred. Final automated marker `6052469777`
completed on `4b5a0888` at `2026-10-08T05:05:12.526071Z`; it is not human approval
or an exhaustive code-review report. Check original review commit identity when
GitHub has rebased inline comments.

## Exact next work and constraints

PR132 initial automated review identified omitted status headers in linked region
and QA records (inline `4228431509`, formal `5467981898`, original head `304e839`).
The correction covers all REGION_* entry points and QA/parity/alignment/workspace
records, preserving 37 existing histories plus checkpoint/prompt. Earlier initial
candidate runs are superseded history, not corrected-head acceptance. Read PR132's
latest full candidate/accepted-push evidence and final review state. This finding
is documentation correctness, not a new product defect or another whole-project
audit. Corrected documentation QA ran before this prompt update; a second prompt
QA and fresh full native/source checks are required afterward.

No new implementation slice had started after PR131. First finish any pending
acceptance/recovery checks for the docs handover or later live proposal. Then
inventory remaining region resource/SKU/availability and pricing/history adapters
against pinned source, including completeness and per-run coordinator integration.
Select a bounded independent offline task from the live specification/roadmap and
publish/verify a pre-edit contract before production edits. Do not repeat completed
calculations, captures or quota/reservation collectors. Final public region paths
need standalone/mixed/scanner/registry/all-format/empty/default/partial-health
integration. B6 operator comparison/alternative-SKU tools and separately scoped
MCP design, plus B7 load/maintenance/release preparation, remain to reconcile.
No listener exposure or release publication is authorized merely by naming them.

Laptop and Azure access were expected soon but have NOT been confirmed. Continue
offline work without input where possible. Defer live checks until identity/access
and the approved whole resolved scope are confirmed. DV-001 requires a wholly
non-production nested management-group hierarchy. AdvisoryDev is a leaf; Advisory
also contains Prd. A Dev subscription include filter does not authorize discovery
through that production-containing parent. Never substitute it for missing scope.

Live optional-stage/plugin parity, restricted visibility, Diagnostics/Advisor
uncertainty/Arc numeric shape, nested MG/DV-001, representative load, fresh OS,
hosted maintenance, Gate004, dependency advisory and release/pilot/provenance/license
evidence remain open. Historical live RG/leaf checks have their own revisions and
limits. SARIF identity exposure, narrower core human masking, multi-file report
nontransactionality and platform/crash/ACL limits remain documented.

Proceed autonomously within existing repository authorization: code/doc fixes,
tests, task-owned coherent publication and protected expected-head acceptance only
after full mandatory exact-head evidence. Keep Azure operations read-only. Do not
create Azure resources/fixtures, change roles, substitute production scopes, expose
services or publish a release without relevant authorization. Ask only for genuinely
blocking input. Do not assume user silence or elapsed time is approval/access.

## Interruption safety and stopping point

Scratch root was `/workspace/scratch/03055aa2c763`. Handover edits were in
`handover-20261009`; clean implementation snapshot `audit/accepted-b653a3ab`.
Some older worktrees have historical physical HEADs and staged current files.
Preserve unrelated edits/untracked source fixtures, and do not reset/clean them.
APRL gitlink metadata is not a materialized checkout. Go 1.26.8 was at
`tools/go/bin/go` under scratch root, absent default PATH; local PowerShell was
absent. Re-inventory tools rather than assume paths survive.

Publish coherent checkpoints and verify remote full file bytes, original tree,
ordered parents, identity and ref. Record completed/pending work, findings,
revisions, run/job IDs, uncertain operation outcomes and exact next action in
handover/logs and PR bodies. Current docs already underwent a document QA pass;
the prompt and complete candidate require an additional QA pass and their own
final native/source acceptance evidence in the handover PR. Read that evidence
before claiming completed publication/acceptance.

Raw logs/cache/scratch/local commits are transient, not durable backups. No accepted
code or implementation WIP was known lost at the last checkpoint; historical
availability draft loss is separately recorded. Never describe scratch-only work
as backed up, claim tests/review/live approval that did not occur, or promise
background execution after a turn ends. Classify failed checks before rerunning,
preserve evidence and source pins, never suppress faults to get green.

Continue to the next natural re-evaluation point with a verified, remotely published
checkpoint. If access remains unavailable, leave live gates explicitly open and
continue useful independent work until a real dependency blocks progress.
