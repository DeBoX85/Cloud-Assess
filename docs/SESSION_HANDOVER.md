# Fresh-session handover: start here

Updated: 2026-10-03 (Europe/Oslo). Owner: repository maintainer/development assistant. This is the current navigation and restart record, not a replacement for specifications, evidence or source inspection.

## Current state and immediate action

- Product: Cloud Assess, functionally equivalent to the pinned AZQR reference with independently adjustable branding. Architecture/code need not be identical. Preserve required attribution.
- Primary repository: [DeBoX85/Cloud-Assess](https://github.com/DeBoX85/Cloud-Assess). Accepted development branch: `bootstrap/core-v1`; do not assume `main` is the development baseline.
- Current accepted runtime baseline: PR90 merge `ea0aabfca055899e18baa54af24b53df90a57d98`, tree `dd4f7b39b27cdf1f7c3782ffbc141599871185e2`. SQL execution and carbon pure aggregation are accepted offline. Carbon HTTP/access/pagination is the next unimplemented task; registry/CLI remains unavailable. This documentation checkpoint can be newer than the runtime hash. Verify live refs; older candidate/pending paragraphs below are historical.
- Accepted: generic assessment, reports/equivalence tooling, bounded branding, rules inspection, scanner commands, B2 YAML Graph execution, B3a zone adapter, B3b canonical plugin tables.
- B3c coordinator accepted through PR84: head `88d7452df8942341513cc62a49d9df24a32e167e`, run `37061985911`, Linux/Windows jobs `111020561369` / `111020561056`, merge `6a920caa76576e61708717a647d67124ebb6ede6`, tree `85365fd1ab7440b4f4d19ce335f188bcbd9f6983`. Actual native results/logs, identity/tree/parents/current ref verified. CLI/list-info accepted through PR85: tested head `ac85fcb4754c96fe8b0046840cb05e8bb2a39d65`, run `37064059199`, Linux/Windows jobs `111027374427` / `111027374648`, merge `2f88e9d5b8dca309c0f27ccb6ec4ed54787bcebb`, exact tree `d6b2f5a477e454d27d21aa3bfae2a7c26e256f3e`; logs, identity, ordered parents, current remote and clean fetched state verified. Eight built CLI and nineteen default plus nineteen custom package cases per host passed without acceptance skips; coverage 78.3%. Service-health execution accepted through PR87: tested head `6c544ea989d0948bb158027ebfc8c110c626b0a2`, run `37069137247`, Linux/Windows jobs `111044153941` / `111044153636`, merge `c892227a440f431363acd626add999e32c152c0e`, tree `8a15549723926129fe151f6fe94515178aba3167`. Exact native logs/identity/tree/parents/current ref/clean fetched state verified. SQL EOL library accepted through PR88: tested head `58ece1ff1d6ec515151f5f664aa6c65444f70ee7`, run `37070978399`, jobs `111050085480` / `111050085721`, merge `f6dcbba032785ed86af63bb018cf3c2a0fcb204a`, tree `00d0478e8005c16036e62c94cddc0ae3be25649d`, native/postmerge verified. SQL EOL execution accepted through PR89: head `fbed3340fa873fcee82ebf4131e49e94c4c3896b`, run `37073129430`, Linux/Windows jobs `111056952522` / `111056952419`, merge `735cedfe68a9b91a17ec120bafd320edb8629d1e`, exact tree/ordered parents/current ref/clean fetched checkout verified. Carbon core accepted through PR90: head `9cbe9569a85413b2db2c4459e5cba07c7791c7c2`, run `37074676704`, jobs `111061797849` / `111061797618`, merge `ea0aabfca055899e18baa54af24b53df90a57d98`, exact native/postmerge verified. Next bounded carbon request adapter; see [CARBON_EMISSIONS.md](CARBON_EMISSIONS.md). The earlier service-health library passed exact-head run `37066186931`, jobs `111034429914` / `111034430064`, tested head `e38db1400a6f06ac4fd837f24f0e4fc9074f7a28`; postmerge tree/parents/ref/fetched state verified. See [SERVICE_HEALTH.md](SERVICE_HEALTH.md).
- Unfinished B3c edits were lost with a missing scratch workspace. The last attempted patch was never verified. Reconstruct from the accepted implementation, not from an assumed surviving patch. That historical loss does not invalidate the later reconstructed PR84/85 acceptance.
- The user cannot currently run laptop/Azure tests. The next task is offline and needs neither. Continue authorized independent work; preserve live deferrals.
- Gate 004, release readiness and full AZQR parity remain open. Do not infer completion from coverage, PR count, estimates or green CI.

## How to establish truth in a new session

1. Read this file, root `AGENTS.md`, the execution plan, target specification and current roadmap checkpoint. Fetch live `bootstrap/core-v1`, open PRs and the latest handover/checkpoint. Inspect actual files/status before edits; inventory paths rather than guessing names.
2. Confirm accepted head/tree, proposal base/head/tree, merge parents and Git identity. Check whether another session or maintainer advanced the branch. Preserve unrelated user work; use an isolated checkout/branch.
3. Follow the document authority map below. Earlier chat exports and historical ledger/gate paragraphs are reference evidence, not automatically current decisions. Where current status contradicts implementation or final PR evidence, classify and reconcile it before relying on it.
4. Verify toolchain, source pins and APRL submodule before build/test. Missing scratch directories or previous executable paths are expected recovery possibilities, not proof of a product defect.
5. Restate the bounded next task, expected behavior, acceptance tests and deferrals in the proposal. Resume independent work without requiring repeated routine "continue" messages.

| Question | Authoritative entry point |
| --- | --- |
| Intended product, scope and deliberate differences | [TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md); current explicit user instructions prevail over an older document |
| Implementation definition of done | [IMPLEMENTATION_PLAN.md](IMPLEMENTATION_PLAN.md) |
| Ordered batches, authorization, recovery and rollback | [DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md) |
| Current work/dependencies | This handover and latest [ROADMAP.md](ROADMAP.md) checkpoint; verify against live refs and final evidence |
| What happened and why | [DEVELOPMENT_LEDGER.md](DEVELOPMENT_LEDGER.md), final PR bodies and linked run logs |
| Mandatory development QA | [QA_PROCESS.md](QA_PROCESS.md), active ruleset and `.github/workflows/test.yml` |
| Source features versus target gaps | [FEATURE_PARITY_CHARACTERIZATION.md](FEATURE_PARITY_CHARACTERIZATION.md), [ALIGNMENT_REVIEW.md](ALIGNMENT_REVIEW.md), actual pinned source |
| Evidence gaps and gate decisions | [QUALITY_GATE_004_PLAN.md](QUALITY_GATE_004_PLAN.md), [DEFERRED_VALIDATION.md](DEFERRED_VALIDATION.md) |
| Previous mistakes and prevention | [FAILURE_NOTES.md](FAILURE_NOTES.md) |
| Environment/operations/permissions | [ACCESS_MODEL.md](ACCESS_MODEL.md), [OPERATIONS.md](OPERATIONS.md), [EQUIVALENCE.md](EQUIVALENCE.md) |

Historical Gate 001-003 records are snapshots. Chronological roadmap/alignment paragraphs can describe features before later implementation; use dated final evidence and current code. An unmerged proposal is not accepted behavior. A successful synthetic case does not establish live equivalence. Source bugs are not automatically target requirements: characterize, document and test deliberate corrections.

## Repository identity, provenance and access

- Reference repository: [DeBoX85/azqr](https://github.com/DeBoX85/azqr), pinned `8e4f0577f3615e6c9014c031bcad079f235369cc`. The user's local checkout is named `azqr-reference`; that directory name is not the GitHub repository name.
- APRL submodule `internal/rules/upstream/aprl`: `60eaddda76541f6adbc1c5ffa686829807e55e29`.
- Other provenance checked by CI: AOR tree `a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5`, CUSTOM tree `674b9b3dcb443ce6dc445b48e1db47e4a0ca7082`, SKU blob `a2d97a60ec445ce4023a1a5a9b6c4dca76e31f5a`. Recheck exact paths/pins in the workflow and [DEPENDENCIES.md](DEPENDENCIES.md); do not update as an incidental fix.
- Human development identity: `Denis Bogunic <134433647+DeBoX85@users.noreply.github.com>`. Configure locally in task-owned clones, verify remote author/committer. Maintenance bot identity is separately defined by maintenance workflows.
- User authorized repository implementation, documentation, QA, proposals and protected merges. That is not permission to write Azure resources, create test hierarchies/change roles, publish a release, expose new services or substitute production scopes for non-production validation.
- Use an available authenticated GitHub connector or authorized Git transport. Read-only public fetch may work while push does not. Never request tokens/credentials in chat, weaken protection or silently probe a browser session. If capabilities fail, preserve a verifiable checkpoint and identify the actual blocker.
- Active ruleset verified: ID `23890737`, protecting `bootstrap/core-v1`; PRs and required `quality` plus `windows-validation` (integration `15368`), no deletion/non-fast-forward. Verify live rules rather than assuming a legacy branch-protection endpoint describes all enforcement. No independent-person review is implied by self-review or zero required approvals.

## Product invariants and feature boundaries

Normal assessment includes authentication/discovery, scope/filtering, inventory, pinned recommendations, Diagnostics, Advisor, Defender plans/recommendations, Policy, Arc SQL and Cost with explicit stage health. Excel and canonical JSON are first-class; CSV, SARIF and JSON stdout also exist. Check [TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md) and actual CLI for defaults/flags rather than inventing them.

- Recommendation libraries APRL/AOR/CUSTOM are retained. CUSTOM is the renamed AZQR recommendation source; renaming is not removal. Specialized APRL files retained outside the normal reference catalog are not automatically executed by normal scans in either implementation.
- Ordinary results retain schema `1.0` and twelve core tables. Explicit canonical plugin-table builds use schema `1.1`; source metadata/columns/rows plus owned health/identity fields. Failed/skipped plugin tables cannot masquerade as complete success. See [PLUGIN_TABLES.md](PLUGIN_TABLES.md).
- Bounded YAML/KQL recommendation plugins run through the normal Graph path; these differ from explicitly selected internal table plugins. See [YAML_GRAPH_PLUGINS.md](YAML_GRAPH_PLUGINS.md). Safe strict input behavior intentionally differs from source silent skips.
- Six eventual internal plugins: `zone-mapping`, `service-health`, `sql-eol`, `carbon-emissions`, `ai-gov`, `region-selection`. Zone execution/CLI/list-info is accepted through PR85; service-health library/execution through PR86/87. SQL EOL execution/registry is accepted offline through PR89; carbon, AI governance and region selection remain unavailable. Do not list unmigrated plugins as executable.
- Five active branding fields have accepted immutable build/profile/package/report behavior. Runtime branding, logos/themes and release approval remain separate. See [BRANDING_PROFILES.md](BRANDING_PROFILES.md) and [BRANDED_PACKAGES.md](BRANDED_PACKAGES.md).
- Default subscription-ID masking applies to JSON/stdout, CSV and Excel. SARIF intentionally retains resource identity. Redaction is not anonymization. Raw Azure evidence stays private outside Git; report-directory security and Windows ACL limitations are documented.
- Output replacement protects existing files on render failure, but multi-file export is not transactional; crash durability and universal rename atomicity are not claimed. Formula/text/Unicode/sheet/privacy limits must survive new adapters.
- Read-oriented assessment can use POST for ARG queries or ARM batch envelopes; method alone does not prove a resource mutation. Verify endpoint/body/subrequest contracts, no automatic provider registration, allowed cloud/audience, redirects, cancellation, body/page/row budgets and sanitized errors.
- Injected transports/plugins are trusted application code, not sandboxed. Same-origin guards are not DNS/IP allowlists. Never conceal incomplete data by normalization, truncation or fabricated empty success.

## Accepted recent implementation and evidence

| Slice | Tested head | Native run; Linux / Windows jobs | Accepted merge |
| --- | --- | --- | --- |
| B1 scanner commands, PR78 | c3907d252b1d9cb48079efe618db1e55028a47e2 | 36947056761; see [PR78](https://github.com/DeBoX85/Cloud-Assess/pull/78) | 50fbaee905e7be7a7c99e174f35ea69247eb16bb |
| B2 YAML Graph, PR79 | 6842c50c6aa12ec10a7472483a6b528fc6938635 | 36963458729; 110701991878 / 110701991724 | 7cddf9368272dd425315bdf3635b21fb974a9fe4 |
| B3a zone adapter, PR80 | 63bb1c08b5de838af291a591d45bcbef9c060489 | 36966029765; 110709854293 / 110709854116 | 0b76ad6b2870432357e5a38c1e8fe73d5214f11b |
| B3b tables/reports, PR81 | 8b1365218743d5b34fb934d24354699bf55465e2 | 36968213639; 110716475829 / 110716476011 | 634aa083981fc327fe5c2cc7b42815793eb197fd |

Final [PR79](https://github.com/DeBoX85/Cloud-Assess/pull/79), [PR80](https://github.com/DeBoX85/Cloud-Assess/pull/80) and [PR81](https://github.com/DeBoX85/Cloud-Assess/pull/81) retain exact trees, reviewed limits, native logs, identity, merge parents and rollback. All six required jobs were reinspected during 2026-10-02 recovery. Each host executed six built CLI tests and nineteen default plus nineteen custom package cases without acceptance skips. Latest coverage was 78.0%, above unchanged 75% floor; coverage is not a parity/security proof.

Source captures retained in Git:
- `internal/rules/testdata/yaml-source-conversion.json`: SHA-256 `dbd0ba34d2aff68c8ff8dd28eb353e2906085fa43821bea24becd3188856772c`; actual unchanged source loader, only temporary CommandPath normalized.
- `internal/plugins/zone/testdata/source-rows.json`: SHA-256 `f75b14ea9012511ed976563984748a63d0bf8baa360c6135e936f3cbcf8aab42`; unchanged source parser values under corresponding public field names.

Fresh recovery QA on the accepted baseline passed full Linux race tests, vet, strict native build, six real CLI tests and documentation checks (223 local destinations/nine flags/one PowerShell snippet). This was not a fresh Windows run or comprehensive security audit. See [ZONE_EXECUTION_CHECKPOINT.md](ZONE_EXECUTION_CHECKPOINT.md) and [PR82](https://github.com/DeBoX85/Cloud-Assess/pull/82). PR82 is a documentation recovery record, closed without implementation acceptance after its content was retained in PR83. PR83 exact head `ed081d5fff7e24f40bb44cf3927dea482ec7c112`, run `37060126509`, Linux/Windows jobs `111014451610` / `111014451903` passed; merge `d88f9294482eafc4e724dbe2248fc09ed4f826b3` and exact tree/parents verified. Documentation checked 253 destinations; no runtime behavior changed in that slice.

## Ordered remaining work

Follow B1-B7 in the execution plan; B1/B2 are accepted. Each remaining feature should be a coherent reviewable slice, with its own contract/source evidence/tests. Do not expand scope to accelerate apparent progress.

1. **B3c zone execution:** coordinator accepted through PR84; preserve its request ownership, pending/partial/context health and default schema invariants. See ZONE_EXECUTION and final PR evidence.
2. **B3c CLI and registry:** VERIFIED OFFLINE through PR85; actual native checks and protected acceptance completed. Successful live credential/scan remains deferred.
3. **Second real adapter:** real service-health library accepted through PR86, with actual source capture and bounded authenticated request tests. Execution integration accepted through PR87; raw/masked reports and actual native/default/custom package gates passed. The earlier six-column table fixture alone did not close this task. This is B4 and the B3 abstraction check, counted once.
4. **B4:** service-health and SQL EOL accepted offline through PR87/89, including negative/partial/request/report cases. Live credential/model evidence remains deferred.
5. **B5:** carbon emissions, AI governance and region selection separately. Source calculations, windows/clocks, batching/caches, audience differences, input failures and all output sheets need independent expectations. Do not replace source scoring with an invented simplification.
6. **B6:** source operator compare, alternative VM SKU, then separately reviewed MCP exposure/auth/lifecycle design and implementation. Existing semantic comparator is not the source operator compare command. No incidental public website/service work.
7. **B7 / acceptance:** unresolved volume/load/metadata bounds, hosted maintenance, artifact/license/provenance/security/operations evidence; Gate 004 decision and release gate for explicitly agreed scope. Live tasks can resume when an appropriate environment is available; never claim closure from fixtures.

Each item must state its dependencies and whether it can proceed offline. Percentages in older planning paragraphs are coarse historical estimates with about ten-point uncertainty; do not count commits/tests to update them or treat them as approval. No reliable calendar ETA is established.

## Retained B3c reconstruction contract and accepted offline tests

Details are in the retained checkpoint. Required source behavior: normal scan selects internal names via `--plugin`; top-level zone command performs subscription discovery and zone only. Subscription filters/management-group discovery apply. RG/type/tag filters do not filter location mappings. Source ignores locations continuation and logs/skips failures; target adapter deliberately follows safe pages and reports incomplete health.

Reconstruct these interfaces from actual current files, not guessed patch contexts:
- `internal/stages/config.go`: fresh plugin-only config, cloned caller config and validation.
- New internal registry/projection: fresh zone metadata/pending table, typed rows and sanitized error codes; validated UUID correlation remains maskable.
- `internal/orchestration/operations.go`: optional lazy zone operation capturing selected ARM origin/audience; no ordinary-scan zone-only prerequisite.
- `internal/orchestration/coordinator.go`: copied names/mode, preflight, inventory skip in plugin-only mode, pending requested table before critical discovery, row retention and truthful table/stage/result health.
- `cmd/cloud-assess/command.go` and registry commands: shared flag binding/preflight, separate normal/plugin-only paths, no unused YAML discovery for zone-only, offline list/info with strict discovery errors and safe terminal text.

Unknown/unmigrated names, unnamed plugin-stage requests, selected-plugin/disabled-plugin contradiction, re-enabled regular stages in plugin-only mode and unused nonempty zone target-regions must reject before credentials. Collapse repeated exact names deterministically. Preserve normal mandatory Graph/B2 YAML behavior. Registry advertises implemented plugins only; adding source-file metadata must not execute queries or allow terminal control injection.

Acceptance must cover: actual Cobra/preflight through real coordinator/application and synthetic authenticated zone transport; every disabled-operation tripwire; selected scope/cloud/audience; mixed/scanner behavior; config/request ownership and repeat/concurrent isolation; healthy/empty/denied/malformed/later-page/cancelled data; table/overall health and artifact-before-partial exit; requested pending/failed empty headers after discovery failure; invalid output sanitized placeholder with other data retained; registry metadata/text/JSON/args/help/branded discovery/no-auth; raw/masked reports. Update existing built CLI unnamed-plugin expectations only when behavior changes. A compiling isolation/health negative control must fail its named assertion, then restored code must pass. Initial pre-CLI green tests from lost work are not reusable final evidence.

## Deferred validation, uncertainty and risks

| Item | Limit and next closure requirement |
| --- | --- |
| DV-001 nested management groups | No suitable non-production hierarchy. `AdvisoryDev` is a leaf; `Advisory` contains production. Do not scan that parent or create/change hierarchy. Synthetic traversal tests do not close live parity. See deferred register. |
| Nonempty optional-stage evidence | Cost and Defender plans have nonempty live evidence. Policy and Defender Recommendations have only empty live execution evidence plus synthetic rows. Arc SQL numeric `vcores` response shape/live equivalence remains open. |
| Diagnostics historical HTTP400 | Later GET/single-batch probes reproduced Network Watcher failures, but original batch failures lacked resource correlation. Retain `complete_with_warnings`; do not claim historical identity from matching counts/hashes. Phase V warning details also remain unclassified. |
| Tag/RG Advisor difference | Include-tag and RG selected the same three resources; Advisor was 2 versus 3 across separate scans. Missing row's nearest recorded scope was unknown. Both source/target passed each paired run; cross-run Azure timing/returned membership is unverified. This is not proof of tag-filter defect or live closure. |
| Live scopes/filters already exercised | Default/optional stages, RG, two subscriptions, leaf MG, scanner selection, separate tag/RG include/exclude, recommendation and individual-resource exclusion, same-RG precedence. Exact inputs/results/caveats are in ledger/equivalence. Do not erase these passes or generalize beyond their datasets. |
| AR-01 dependency advisory | Existing module-only OpenPGP advisory remains open; accepted scans had zero reachable/imported findings. Recheck official advisory/current graph before a new claim. Module-only does not mean universally safe or release-approved. |
| AR-02 lifecycle/volume | Several deadlines/page/cycle/cancellation guards and ARG completeness fixes exist. Universal default budget, global metadata/volume/load and context-ignoring trusted code limits remain open. |
| AR-03 branding | Bounded five-field workflow accepted. Runtime/logo/theme/SARIF consumer/release boundaries not certified. |
| AR-04 full functionality | Rules/scanner/YAML accepted; internal execution/registry and later ancillary commands remain. Characterization is not implementation acceptance. |
| AR-05 through AR-08 | Historical warning/timing, DV-001, restricted visibility/nonempty stages/Arc, hosted maintenance/fresh OS/release/pilot. Keep separately tracked; do not flatten into a completed generic QA task. |
| Deployment/release | Successful candidate ZIP extraction on runners is not fresh-OS install, installed Azure execution, signed publisher provenance, legal sign-off or approved production pilot. |

Private raw live evidence is on the user's laptop/outside Git; metadata/equivalence files were supplied in earlier sessions. A fresh model must not assume those attachments remain accessible. Ask only for the particular missing evidence when closure depends on it. Historical local locations include `C:\src\Cloud-Assess`, `C:\src\azqr-reference`, and timestamped `artifacts\equivalence` bundles documented in the ledger/runbook. No credentials are needed for offline engineering.

## Rebuild and QA procedure

Clone/fetch into a task-owned directory. Public reads can use Git; authenticated publication uses available authorized capability. Locate source files with `git ls-files`/`rg --files`, confirm tracked membership, and use explicit separate source/target workdirs. Use a temporary source copy for characterization tests; retained pinned reference stays unchanged.

```bash
# Run from a clean task-owned Cloud-Assess checkout after verifying its ref.
git submodule update --init --recursive
# Verify expected commit before compiling embedded APRL input.
git -C internal/rules/upstream/aprl rev-parse HEAD
go version                         # required CI toolchain: go1.26.8
go test -race ./...
go vet ./...
# Keep build output outside the strict clean checkout; adjust destination locally.
CGO_ENABLED=0 go build -trimpath -mod=readonly -buildvcs=true -o /tmp/cloud-assess-handover ./cmd/cloud-assess
python3 scripts/tests/built-cli.py --binary /tmp/cloud-assess-handover
python3 scripts/tests/documentation.py --binary /tmp/cloud-assess-handover --powershell pwsh
```

These Linux examples are a recovery starting point, not the entire required workflow. Read the current workflow for formatting, module consistency/provenance/inventory, branding/report/package tests, PowerShell helper/failure checks, coverage, bounded fuzz/mutations and vulnerability checks. Go 1.26.8, Python and PowerShell must be discovered/installed using permitted capabilities; previous scratch paths are not durable. Ensure clean stamped builds and required checks. Do not weaken tests for an environment issue or claim a compiled mutation was detected if it only failed compilation.

Before publication: review exact diff/contract, failure recurrence, privacy/security and source evidence; stage all new files; enforce diff-check exit codes. Publish exact candidate tree with correct identity and parent. Observe actual `quality` and `windows-validation` success on final head, inspect logs, verify current base/preview/tree and expected-head protected merge. Verify merged parents/tree/current remote and fetched state. Do not merge stale or unverified WIP.

## Outage prevention, update cadence and rollback

Root `AGENTS.md` directs future repository agents to read and maintain this file. This is an active-session procedure, not an unattended scheduler or guaranteed recovery service.

- Before production edits, update the bounded task contract, accepted base and pending checks. After every coherent slice, material decision/failure, accepted merge or changed next action, update this handover and relevant ledger/roadmap/failure notes. Include changed interfaces, uncertainty, tests tied to commit, and exact next action.
- Publish a task-owned remote WIP **code plus documentation** checkpoint before changing tasks, substantial further work, or ending a session when GitHub is available. Verify remote identity/head/tree/parent by read-back. Mark WIP unaccepted and unmerged. Do not publish secrets/raw Azure captures.
- A local commit, planned upload or design-only note is not a durable code backup. If publication fails, explicitly state what is only local and at risk. A dropped response is an unknown outcome: inspect refs/commit/PR state before retrying. Preserve failed logs and classify cause versus hypothesis.
- On recovery, reread live state and inspect actual files. Reconstruct absent work from verified baseline/contract and rerun QA. Do not blindly replay an uncertain patch, reuse stale tests, promise background progress or rely on model memory.
- Rollback accepted changes through a reviewed revert PR on current history. Inspect dependencies and correct mainline parent first; never force/reset protected history, user work or pinned source. Preserve proposal and acceptance evidence.
- Keep this file compact enough to read; detailed history belongs in linked records. Record its verified-against implementation commit rather than attempting a self-referential containing-commit SHA. Final native/merge evidence first goes in PR, then the next material ledger/handover update.

PR82's durable checkpoint at `2d413cca968784bfb7613490283bb933e3653894` documented the interruption and recovery. Its content is retained as the linked local checkpoint in this proposal. Once integrated, PR82 is superseded as a separate proposal, not accepted as B3c code. FN-037 records lost unpublished work, missing-submodule test setup, bounded-read/branch-helper corrections and prevention. Review FN-004/FN-007/FN-035/FN-036 before new path/format/consumer investigations.

## Paste-ready prompt for a new session

> Continue development of https://github.com/DeBoX85/Cloud-Assess on bootstrap/core-v1. Start by reading AGENTS.md and docs/SESSION_HANDOVER.md from the live development branch, then the linked specification, execution plan, roadmap, ledger, QA process and failure/deferred records. Use the handover's actual URL/ref if its proposal is not merged yet. Verify live branch/open PR state, exact accepted/proposal commits and evidence before editing. Treat historical chat and stale status as context, not proof. Our goal is functional AZQR parity against the pinned reference with adjustable branding and high-quality, secure, read-oriented behavior. Preserve explicit differences, attribution, pins and unresolved evidence limits. Repository work and protected PR publication/merges are authorized; do not weaken gates or write Azure resources. I currently cannot run laptop/Azure tests, so continue independent offline work and keep dependent validation deferred. Resume the exact next task in the handover and work until genuine input/specification/access is required. Maintain the handover, ledger/roadmap/failure records and verified remote code checkpoints after coherent slices. Report verified state, next action and blockers; do not claim background execution or recoverability of unpublished work.

Handover completeness review on 2026-10-02 checked this record against live core-v1/ruleset and final PR79-82 evidence, execution/QA plans, current feature boundaries and deferred records. The fresh-session reader can recover goal, pins, accepted state, lost-work boundary, next ordered task, required tests, authorization and risks without chat. This was self-review against repository/source evidence, not independent-person approval. Local documentation validation passed 253 destinations/nine flags/one PowerShell snippet; final native acceptance for this documentation slice is retained in its PR and indexed at the next material update.

New-session startup is complete only when the model can identify the product contract, verified baseline, implemented features versus unaccepted WIP, next dependency-ordered slice, applicable checks, open risks and authorization without guessing. If the repository is inaccessible, state that blocker and request this handover plus the specific linked files needed; do not invent their contents or request credentials.

Service-health current slice: interruption reconciliation reconfirmed PR85/86 native success, core-v1 at `a6b790ce143bc70fac8f4aa3199225ee92b0660d`, clean unchanged pinned reference and preserved integration work. Source capture and exact query hashes recorded in SERVICE_HEALTH. A real bounded library scanner and bounded POST helper are implemented with source table projection, typed scope/schema validation, explicit partial rows/health, per-run budgets and synthetic authenticated request checks. Full restored race/vet and compiling page-limit control passed. Library native publication/acceptance completed through PR86; no public service-health availability or live equivalence is accepted. Resume at coordinator/registry/CLI/report integration under SERVICE_HEALTH contract, then final native proposal gates.

Integration resume checkpoint: service-health standalone/mixed/scanner/dual dispatch, source type-filter scope, owned/sanitized projection and all raw/masked formats are implemented. Eighteen authenticated synthetic Cobra-to-application cases passed before interruption; restored focused/full race and vet plus the compiling metadata-boundary negative control passed. Clean nine-case actual CLI/documentation and exact-head native acceptance are the remaining checks. Projection review tightened region/type labels to the scanner’s existing nonempty 512-byte/control-free contract. Repeated/concurrent selected scope/type isolation is explicitly tested. No live credential/factory execution or Azure parity is claimed. Publish code plus docs before starting SQL EOL; keep final run/merge in PR and index at the next material checkpoint.

Current next action after PR87: implement bounded SQL EOL library under SQL_EOL contract, preserving exact 611-line source query/32 columns, subscription-only filtering, string/null projection and received row ordering. Actual unchanged source scanner captures nonempty/filtered/empty tables with synthetic TLS/token, not KQL evaluation or live pricing/licensing validation. Final library native acceptance must precede separate public integration. Service integration pending-check paragraphs above are historical; PR87 acceptance supersedes them.

SQL EOL current candidate: source capture/query and real bounded library implemented under SQL_EOL; public registry/CLI still service-health/zone only. Focused race and compiling subscription-filter negative control passed after FN042 fixture corrections. Empty-fragment endpoint acceptance reproduced in the new SQL and existing service constructor, narrowly corrected with regressions (FN043). Final restored full race/vet, clean actual built CLI/docs and native Linux/Windows before library acceptance; code/docs remote checkpoint precedes integration. Source pin unchanged, no Azure access.

SQL EOL final local checkpoint: restored full race/vet passed on the complete library and endpoint correction. Resume clean stamped actual CLI/docs, exact remote candidate readback, both native gates/logs and protected merge; successful library acceptance does not advertise SQL EOL. The next separate integration must preserve subscription-only semantics and all 32 raw/masked cells with mixed/partial/context behavior.

Current resume: PR88 library/native acceptance completed. Implement SQL execution under SQL_EOL’s 32-column/subscription-only contract; previous library pending-check paragraphs are historical. Candidate preflight/ownership/request/report/full/built/native checks remain required before public acceptance. No Azure/laptop dependency for this integration.


Current SQL execution checkpoint: optional selected operation, honest registry and standalone/normal/scanner dispatch implemented. Twenty-four authenticated synthetic raw/masked command-to-application cases compare every source-captured JSON/CSV/XLSX cell; CSV independently expects its existing apostrophe protection for negative savings text. Healthy regular inventory survives SQL failure; three-adapter order/dedup, prior rows and pending cancellation tables survive. Owned repeated/concurrent scope/filter copies and subscription-only SQL semantics are tested. Focused/full restored race and vet passed, and a compiling metadata-boundary negative control failed its named assertion. Ten actual built CLI cases plus documentation and exact native gates remain required before merge. No production credential-factory, live KQL, pricing/licensing or Azure completeness proof is claimed. FN044 records fixture/read corrections. Publish code and this checkpoint together, verify exact tree/identity/parent, then resume native acceptance; next B5 carbon characterization.


Current B5 checkpoint (2026-10-03): PR89 native/merge acceptance completed; earlier SQL pending-check prose is superseded. Carbon actual unchanged source scanner/all-column and six-case arithmetic captures passed after correcting synthetic SDK transport (FN045). Pure bounded aggregation/projection core is implemented without HTTP/registry/CLI; focused race passes. Next finish compiling source-calculation negative control/restored full race/vet, clean built/docs, publish exact code+docs and native acceptance. Then implement separate bounded date/report adapter with pagination/access metadata and truthful failed batches, followed by execution/reports. Carbon source silently hides denied batch and ignores token; both captured, not target requirements. SQL final evidence is indexed in SQL_EOL/ledger/PR89. All laptop/Azure/Gate004/release/advisory deferrals remain open; no input needed for this offline core.


Carbon core local final checkpoint: compiling previous-display mutation failed the actual source calculation 3 assertion; restored full race/vet passed. Clean stamped built CLI/docs and exact-head native acceptance remain pending. All pure-core/request/CLI and live limits above remain distinct.


Current resume after core acceptance: PR90 tested head `9cbe9569a85413b2db2c4459e5cba07c7791c7c2`, tree `dd4f7b39b27cdf1f7c3782ffbc141599871185e2`, run `37074676704`, Linux/Windows jobs `111061797849` / `111061797618` actual PASS; logs/steps reviewed. Ten CLI and nineteen default plus nineteen custom package cases per host without acceptance skips; 274 Markdown destinations/nine flags/one parsed PowerShell snippet; 79.8% coverage. Mandatory QA passed and both scans found zero reachable/imported vulnerabilities, with the prior module-only advisory still open. Protected merge `ea0aabfca055899e18baa54af24b53df90a57d98` exact candidate tree/ordered parents/current ref/human author/GitHub merge committer/clean fetched state verified. Earlier carbon pending QA paragraphs are superseded. The accepted source-grounded pure core is not a complete plugin. Next implement the bounded authenticated date/report library under CARBON_EMISSIONS's latest contract and independent fixtures; no adapter code exists yet. Then separately integrate public execution/all-format reports. User laptop/Azure access is not needed for these tasks. Live totals/service roles/visibility and existing Gate004/release/AR/DV remain open. Final PR90 evidence and this handover checkpoint make the interruption restart concrete; verify live refs/open proposals rather than relying on old local paths.
