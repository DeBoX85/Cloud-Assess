# Development Failure Notes

This register records confirmed mistakes in our implementation, test design, evidence interpretation, or instructions. It is a troubleshooting index, not a substitute for the development ledger, Git diff, CI logs, or a formal quality gate. It contains no credentials or unredacted Azure report rows.

## On each failure

1. Stop the dependent claim or merge. Search this file and the development ledger for the same symptom or assumption (`rg -n -i '<term>' docs/FAILURE_NOTES.md docs/DEVELOPMENT_LEDGER.md`).
2. Record the date, observed symptom, actual cause if known, impact, correction, and prevention check. Mark an unconfirmed cause as a hypothesis. Link the PR, run or ledger section that holds the evidence. A short entry is enough for a small error; related attempts can share one ID.
3. If it recurs, append the new occurrence to the existing entry and strengthen the prevention check. If the cause differs, give it a new ID and cross-reference the earlier one.
4. Re-run the affected check and record whether the correction passed. Leave an item open if verification is still pending. Do not turn a failed or partial Azure assessment into a passing baseline.

Record user-visible instruction errors and errors that affected code, tests, evidence, security, or a quality gate. A command typo caught and corrected immediately can use one compact entry; do not copy large tool logs into this file. Git and CI retain the exact attempts. Revisit this register during QA and before a gate decision.

## FN-001: Advisor cross-run explanation was too specific

- **Date:** 2026-09-25; later local scope check 2026-09-29.
- **Mistake and impact:** An earlier interpretation treated an Advisor row absent from a later tag-filter run as a closed unknown-scope classification before comparing against excluded inventory. The separate runs also could not establish identical Azure API responses.
- **Correction:** The Phase S evidence correction withdrew the exact-cause claim. The operator later found no selected or excluded inventory ancestor in the tag-only report, establishing unknown *recorded* scope, while historical Advisor API timing remains unresolved.
- **Prevention:** Check both selected and out-of-scope ancestors, and distinguish within-run source/target equivalence from across-run causal explanation. Do not close a live-data cause without correlated inputs.
- **Evidence:** [Phase S correction and follow-up](DEVELOPMENT_LEDGER.md#phase-s-evidence-correction-and-recent-commit-qa). **Status:** Interpretation corrected; historical timing open.

## FN-002: Local Azure instruction lacked the exact command

- **Date:** 2026-09-29.
- **Mistake and impact:** After preparing the same-RG precedence fixture, the response pointed to a runbook section without giving the operator the exact PowerShell block or requested return files. The operator had to ask what to run.
- **Correction:** Supplied the exact `AdvisoryDev` command and named `run-metadata.json`, `equivalence.json`, and the console output. The paired run passed and was recorded in the ledger.
- **Prevention:** For a local-only validation handoff, include copyable commands, safe scope, expected evidence directory and exact non-sensitive outputs requested. Check the command against the merged branch before sending it.
- **Evidence:** [Same-RG precedence pass](DEVELOPMENT_LEDGER.md#same-resource-group-includeexclude-precedence-live-pass). **Status:** Corrected.

## FN-003: Azure SDK zero retry assumption in a CI fixture

- **Date:** 2026-09-29.
- **Mistake and impact:** The new Resource Graph HTTP 429 fixture used `MaxRetries: 0` assuming it disabled retries. The SDK applies its default of three retries to zero, so the test received four requests and failed both CI jobs. This was a test configuration error; no scanner defect or Azure request resulted.
- **Correction:** Changed the fixture to `MaxRetries: -1`. Corrected Linux and Windows CI jobs passed, and PR #43 merged.
- **Prevention:** Inspect the exact dependency version's option semantics before writing an edge-case fixture, especially when zero values can mean defaults. When a test fails, compare the observed request count with SDK retry policy before altering product behavior.
- **Evidence:** [Gate 004 throttling checkpoint](DEVELOPMENT_LEDGER.md#gate-004-resource-graph-throttling-characterization), initial run `36597984445`, corrected run `36598430521`. **Status:** Corrected.

## FN-004: Assumed source paths during inspection

- **Date:** 2026-09-30.
- **Mistake and recurrence:** Tried reading nonexistent `cmd/cloud-assess/scan.go`, then repeated the path assumption with `internal/assessment/stage.go`. Both were inspection errors; no files were changed by the failed reads.
- **Correction:** Used `rg --files` and read `command.go` and `types.go`.
- **Prevention:** Inventory paths before reading an unfamiliar package; use symbol search to locate implementations.
- **Recurrence, 2026-09-30 QA expansion:** Tried reading nonexistent `internal/equivalence/compare_test.go` even though the inventory showed `projection_test.go`. Corrected through symbol search. Strengthened prevention: use the enumerated path exactly, and search test function names before guessing a test filename. No product files were changed by the failed read.
- **Recurrence, optional-stage checkpoint:** Guessed `internal/assessment/assessment.go` and then repeated the already recorded `stage.go` path during inspection. Neither exists. Located `types.go` using the file inventory. No product mutation resulted; the prevention rule was not followed and remains required.
- **Recurrence, dependency inventory:** Guessed nonexistent `internal/rules/README.md` despite the path-inventory rule. Corrected by reading the enumerated provenance and embed files. No product change resulted.
- **Recurrence, built-CLI checkpoint:** Assumed an `internal/cli` package during inventory/search; the CLI actually lives under `cmd/cloud-assess`. Corrected from enumerated paths before editing.
- **Evidence:** CLI process-test and QA expansion inspection sessions. **Status:** Corrected; recurrences recorded in this entry.

## FN-005: Windows ACL fixture inherited another PowerShell module path

- **Date:** 2026-09-30.
- **Mistake and impact:** The Go ACL fixture launched Windows PowerShell 5.1 with all environment variables inherited from a PowerShell 7 workflow shell. In run `36725561754`, Get-Acl could not autoload Microsoft.PowerShell.Security, so ACL assertions never ran and Windows CI failed. No product ACL defect was established.
- **Cause:** The inherited PowerShell 7 module path interfered with Windows PowerShell autoloading through the intermediate Go process. Removing it restored the native module path, matching Microsoft documented startup behavior.
- **Correction:** Remove PSModulePath from only the child test process environment, allowing Windows PowerShell to construct its native module path.
- **Prevention:** Isolate version-specific shell environments when invoking a different PowerShell runtime. Keep actual ACL assertions blocking; do not skip them to obtain green CI.
- **Evidence:** PR #53, initial run `36725561754`, Windows job `109921410714`. **Status:** Corrected; run `36726089065` passed both required jobs, including the Windows ACL assertions.

## FN-006: Staged replacement initially omitted existing Windows DACL preservation

- **Date:** 2026-09-30.
- **Mistake and impact:** PR #53 tested new report inheritance from a secure directory but omitted replacing an already restricted report inside a broader directory. Staged replacement could inherit broader access and weaken the existing report DACL. Detected during follow-up review, before any live Azure use or release claim.
- **Correction:** Copy the existing effective Windows DACL to the staged file before writing data, protect it against new parent inheritance, and fail if it cannot be preserved. Add a Windows fixture that deliberately broadens the parent while keeping the old file restricted.
- **Prevention:** Privacy tests must cover creation and replacement independently, including different parent/file permissions. Test successful output and failure preservation, not only nominal private creation.
- **Evidence:** PR #54, tested head `017d5487c59233033fb1d63cf185e9fe841aa53a`; required Linux/Windows run `36728008504` passed, including the replacement ACL regression. Merge `ce61babc85e52fa2d8ca0aa64a9283bb37e0f2fb`. **Status:** Corrected and verified for the controlled fixture.

## FN-007: Inspection helper regex was overescaped

- **Date:** 2026-09-30. A JavaScript tool helper used an invalid regex when searching Windows API declarations; it failed before running any tool. Replaced it with plain string matching. No repository mutation resulted. **Status:** Corrected.

- **Tool-helper recurrence, optional-stage checkpoint:** Workflow discovery returned an empty run list immediately after an update; the helper attempted to store an undefined run ID and was rejected. No repository mutation resulted from the failed lookup. Guard empty lists and repeat read-only discovery before storing or querying a run ID.

- **Inspection-output recurrence, decoder checkpoint:** Printed the complete recursive reference tree when only three source paths were needed, causing output truncation. Retrieved and filtered the metadata before emitting it on the follow-up. No content was published from the truncated output. Use path-filtered metadata and bounded source reads.

- **Tool-helper recurrence, generation evidence update:** Incorrect escaping of Markdown backticks in a JavaScript template produced a syntax error before any tool ran. Built the command as plain lines, used Python character construction for backticks and validated the SHA before substitution. No files were changed by the failed call.

- **Tool-argument recurrence, dependency inventory:** Initially passed `repo` instead of the connector-required `repository_full_name` when creating a tree. Binding validation rejected the request before mutation. Corrected the argument and compared the resulting Git tree exactly with the local staged tree before publishing.

## FN-008: Toolchain availability was inferred from PATH alone

- **Date:** 2026-09-30. Earlier checkpoints said local Go was unavailable after checking PATH. A later source search found an existing Go 1.26.0 toolchain outside PATH, which successfully selected/downloaded Go 1.26.8. Earlier executable verification was performed in CI, not locally.
- **Correction/prevention:** Verify known workspace toolchain locations and the actual selected version before concluding local testing is unavailable. Local tests are now possible; an initial broader run then failed setup because the APRL gitlink was uninitialized. Focused report/config tests passed; initializing the exact pinned submodule restored broader validation, and local `go test -race -count=1 ./...` then passed. Check toolchain, module dependencies and embedded gitlink contents together; continue required hosted Windows/Linux checks. **Status:** Environment availability corrected; previous CI evidence remains valid.

- **Toolchain-path recurrence, built-CLI checkpoint:** The previously verified `/root/go/pkg/mod` toolchain was absent after the session gap. Initial build commands failed. Re-enumerated the workspace toolchain, selected Go 1.26.8 and used an explicit task-scoped GOPATH cache; did not claim executable validation from the failed attempt.

## FN-009: Documentation upload assumed command output was complete

- **Date:** 2026-09-30. An evidence-update tree was constructed from shell output without first verifying complete content. The ledger exceeded the output limit, so the remote tree differed from the local staged tree. Detected before creating a commit or branch; the incomplete tree was never published.
- **Correction/prevention:** Read large files in bounded chunks and require exact local/remote tree equality before publishing. **Status:** Corrected after the tree equality check.

## FN-010: Time-bounded decoder fuzz run failed at its budget boundary

- **Date:** 2026-09-30. The first local Go 1.26.8 decoder fuzz command (`-fuzztime 10s -parallel 2`) reached 148,697 executions and exited 1 with `context deadline exceeded` at 10.31 seconds. No property assertion or saved failing input was reported. Focused race tests and seed cases had passed. No previous fuzz deadline entry was found in this register.
- **Assessment:** Inspected the selected Go toolchain's fuzz coordinator cancellation/error handling. A coordinator time-budget termination issue is a hypothesis; its exact cause is unconfirmed. The failed run is not counted as a pass or evidence of a decoder defect.
- **Mitigation and evidence:** Use an explicit execution-count budget (`-fuzztime 100000x`) with a separate 60-second test timeout and two workers for the new CI target. It passed locally with exactly 100,000 executions, including the retained interesting-input cache, in 7.485 seconds. Property assertions and decoder production code are unchanged. Required Linux/Windows run `36733195533` on PR #57 code head `bde30f36c834fb41766b0821c6b9966e48094926` also passed, including count-bounded fuzzing. Do not suppress command failures or retry until green; keep time-bound failure visible if it recurs.
- **Status:** Execution-count control verified locally and in required hosted CI; original time-bound termination cause remains unresolved.

- **2026-10-01 recurrence:** PR #72 head `5c97810124e12502c994a97849c39b715df52c82`, run `36925929037`, quality job `110583114668`, failed comparator FuzzCompareFieldSymmetry at 11.00s with `context deadline exceeded` after 83,850 executions. No property assertion or saved failing input appeared; retained synthetic failure artifact `11194181343` contains coverage evidence. The cause remains unconfirmed. Apply the existing execution-count mitigation to this target only: 100000x, two workers, separate 60-second timeout, unchanged properties/corpus and no failure suppression or blind retry. Both required jobs must pass the changed final head.

- **Comparator mitigation local evidence:** Local comparator count-bounded runs passed exactly 100000 executions twice (6.688s, then 7.389s with the retained 49-input baseline cache); focused comparator race tests also passed. This verifies the mitigation locally, not the root cause or final native result.

## FN-011: Generation fixtures initially assumed full Git history

- **Date:** 2026-09-30. The first local generation suite passed from a full-history checkout, but CI checks out depth one. Final review reproduced the fixture bootstrap push failure with a depth-one file-URL clone: the disposable bare remote rejected it with `shallow update not allowed`. No GitHub or Azure refs were changed by the reproducer.
- **Correction:** Generation fixtures now always clone depth one and explicitly allow a shallow bootstrap in their disposable bare remote. This setting is never applied to origin/GitHub or the project repository. The proposal/base-ref assertions remain unchanged. Unexpected fixture-command failures now include bounded captured stderr in the test failure to expose the underlying Git error.
- **Prevention:** Match CI checkout depth and fixture-server assumptions during local recipe tests; do not infer hosted readiness from full-history fixture results. Check this register when a bootstrap push fails. Original PR #59 run `36738474726` failed Linux job `109966302312` during fixture bootstrap pushes; Windows passed. The corrected six-case suite passed locally with depth-one fixtures and Go 1.26.8. Corrected code head `620de86e8756007598f1e6565c8855e923fa7d6a` passed required Linux/Windows run `36739045844`, including all six generation cases. **Status:** Corrected and verified locally and in required hosted CI.

## FN-012: Inventory package listing initially hid VCS-status errors

- **Date:** 2026-09-30.
- **Mistake:** Initial package-list command used implicit VCS stamping, which failed with `error obtaining VCS status: exit status 128`. The first Python wrapper exposed only the command exit, withholding useful stderr. No evidence files were produced.
- **Correction:** Package dependency listing uses `-buildvcs=false` because its purpose is target membership; source pins are checked explicitly through Git. Command failures now report bounded stderr and have a five-minute timeout. This does not disable stamping in product builds or change repository ownership settings.
- **Prevention:** Do not rely on incidental executable stamping to determine dependency membership. Preserve command diagnostics, validate source pins separately, and check that failed collection retains previous inventory files.
- **Evidence:** Local initial inventory collection and corrected generation, plus retention fixture. **Status:** Corrected for inventory collection; underlying local VCS-stamping environment issue is not claimed resolved.

- **Packaging follow-up:** A strict local `-buildvcs=true` build again failed. Bounded verbose output showed Go asking Git for status from a different workspace root instead of the project checkout. Kept candidate provenance strict and used an isolated clone for package validation rather than relaxing the package revision check or changing repository trust/ownership settings. This is a local workspace-context limitation, not a product source defect.

## FN-013: Inventory notice whitespace check did not stop a local commit

- **Date:** 2026-09-30.
- **Mistake:** Ran `git diff --cached --check` followed by a commit without enforcing the first exit code. Raw retained third-party notices contain CRLF and trailing spaces; Git reported whitespace diagnostics but the local checkpoint commit still ran. The earlier unstaged-only check had also omitted new files.
- **Correction:** Preserve source notice bytes intentionally with path-specific `-text -whitespace` attributes for generated NOTICES.md only. Use blocking staged checks before publication; freshness/hash checks continue verifying third-party evidence. Authored files retain whitespace checks.
- **Recurrence, built-CLI setup:** Initial shell commands continued to metadata decoding after a missing cached Go binary failed the build, producing a secondary JSON error. Follow-up builds enforce the build exit before tests. No report/binary validation claim was made from that attempt.
- **Prevention:** Stage new files before checking and enforce command failure before the next mutation. Document intentional source-byte exceptions rather than silently claiming every initial check passed.
- **Evidence:** Local checkpoint commit `2037596`, notice diagnostics and corrected staged diff check. **Status:** Corrected before merge; no production behavior affected.

## FN-014: Built-CLI valid-filter fixture initially used the wrong tag schema

- **Date:** 2026-09-30.
- **Mistake:** A fixture labeled valid used `Environment: [dev]`; Cloud Assess tag values are scalar strings. Its case expected the deferred-plugin error but correctly received a YAML type error. The first local suite failed one of fifteen preflight cases; help/version and compiled dependency inspection passed. No Azure traffic was observed.
- **Correction:** Reviewed the actual filter structure and used `Environment: dev`. Re-ran the entire built-CLI suite; all cases passed locally.
- **Prevention:** Ground supposedly valid inputs in the target schema, and keep the expected error precise so earlier unrelated failures cannot masquerade as successful validation.
- **Evidence:** Initial local built-CLI output and corrected suite. **Status:** Corrected locally; required Windows/Linux validation recorded in the ledger.

## FN-015: Package integration skip condition was evaluated before CLI parsing

- **Date:** 2026-09-30.
- **Mistake:** Used `unittest.skipUnless(BINARY)` as a decorator while BINARY was initialized to None and populated later by argument parsing. The first isolated-clone run supplied --binary but still skipped the actual-artifact case. Its twelve unit fixtures passed; it did not validate packaging/install behavior.
- **Correction:** Evaluate missing-binary handling inside the integration method, after argument parsing. Require a complete thirteen-test native run with no skipped integration before publication/merge.
- **Prevention:** Avoid import-time skip decisions for runtime arguments. Inspect executed/skipped counts and require artifact validation instead of treating any exit-zero test invocation as equivalent.
- **Evidence:** First isolated-clone output `Ran 13 tests ... OK (skipped=1)` with --binary, then corrected native runs. **Status:** Corrected before PR publication; final run evidence indexed in the ledger.

## FN-016: Roadmap overview lagged the completed packaging checkpoint

- **Date:** 2026-09-30.
- **Mistake:** The detailed roadmap paragraph described tested candidate ZIP packaging while its overview still said no packaging and retained a 45% toolkit estimate. The conversation had already moved the coarse estimate to 55%.
- **Correction:** Reconcile the overview with PR #62's final passing evidence and partial packaging boundary. Preserve open fresh-OS, publication, license/security/operations and release decisions; do not present packaging tests as release approval.
- **Prevention:** When logging a milestone, review the top-level status, progress basis, current boundary and next tasks together. Search this register for documentation drift before reporting progress.
- **Evidence/status:** Confirmed by source review at `ab0cc3b761956b5ea4712717b96a97d5a98f6567`; corrected in the access/operations review change. Existing FN-007 retrieval-output limits were also reviewed: bounded follow-up source reads replace truncated batches; no unseen lines support claims. A later spec-edit command omitted its heredoc terminator and failed Python parsing before writes. Corrected the command delimiter and reran the staged checks; check heredoc closure before dispatching multiline commands.

## FN-017: Advisor continuation URLs crossed the credential boundary

- **Date:** 2026-10-01 (Europe/Oslo).
- **Confirmed defect:** Advisor metadata trusted arbitrary absolute HTTPS nextLink values. A controlled synthetic transport received the ARM canary header on a foreign-host second request. Historical successful scans do not establish that this boundary was safe. No actual credential/network disclosure was performed.
- **Correction:** Parse URLs and require the configured HTTPS host/port; reject credentials, fragments, malformed/protocol-relative links and repeated continuation URLs without echoing their content. Add authenticated foreign-origin/cycle/normal-pagination regressions. The default shared HTTP transport additionally refuses redirects; a two-local-TLS-server fixture checks no forwarded request. Caller-injected transports and SDK-owned pagers remain separate review boundaries.
- **Prevention:** Treat service-provided URLs as data requiring destination validation before authentication/transport. Check continuation and HTTP redirect policies together, preserving custom configured endpoints rather than hardcoding public Azure. Do not equate read-oriented verbs with safe credential destinations.
- **Evidence/status:** Foreign-host/port/userinfo/protocol-relative/fragment regressions failed before remediation and pass after it; all local race tests pass. Required final CI evidence is indexed in the ledger. Adversarial-response hardening is a deliberate target correction; no claim of pinned-reference adversarial parity.

## FN-018: Runbook flag checker initially mixed different CLIs

- **Date:** 2026-10-01 (Europe/Oslo).
- **Mistake:** The initial checker compared every runbook flag against Cloud Assess help, including `--all` from the separately documented `az account list` command. It failed locally before publication.
- **Correction:** Exclude inline commands explicitly belonging to Azure CLI from the Cloud Assess flag set. Eight actual Cloud Assess flags and the PowerShell snippet then passed local help/parser checks. No scan ran.
- **Prevention:** Attribute examples to their command before validating syntax/options; do not suppress specific unknown flags to force a pass. The checker has a bounded authored-document scope, not a claim of universal Markdown understanding.
- **Recurrence check:** FN-004 also applied during SDK review: an unversioned online page and guessed module layout did not match the selected cache. Read go.mod and enumerate pinned SDK files before inspection; only v1.23.1 source and pinned documentation support implementation claims.

## FN-019: Read-oriented review omitted SDK automatic-registration middleware

- **Date:** 2026-10-01 (Europe/Oslo).
- **Mistake:** Reviewing explicitly called GET/query methods did not account for the ARM SDK's default provider-registration policy. Shared ARM options did not set DisableRPRegistration, allowing an implicit write attempt after certain service errors. No evidence of historical Azure mutation is asserted.
- **Reproduction:** Synthetic valid ARM-resource GET through the actual shared SDK options returned 409 MissingSubscriptionRegistration; the pipeline attempted POST /subscriptions/sub-fixture/providers/Microsoft.Test/register. The desired regression failed with two requests and one injected/rejected write. No real credential or network was used.
- **Correction:** Explicitly disable registration in shared ARM options and in scope-client construction, using a copy so caller options remain unchanged. The pipeline now returns its original 409 with one GET/no POST. Three scope pager failure fixtures confirm typed original error, one read request and no caller-option mutation.
- **Prevention:** Audit middleware and SDK defaults as well as explicit endpoint verbs; pin implicit remediation/registration off for assessment clients and test registration-required negative paths. This differs from FN-017 destination validation and complements its credential-boundary checks.
- **Evidence/status:** Initial failing and corrected focused fixtures, rerun full local race/vet checks; exact updated final CI must pass before merge. Registration configuration is a deliberate target hardening correction, not a claim about pinned-reference adversarial behavior.

## FN-020: Eight-task documentation review left stale current-status wording

- **Date:** 2026-10-01 (Europe/Oslo).
- **Mistake:** Despite the documentation-consistency task, the roadmap still called the completed PR #57 work the next checkpoint. QA_PROCESS ambiguously listed tested local maintenance execution and built-artifact inspection as open, and final PR #64 documents retained pending-CI wording without a final run/merge record. Green structural link/help/parser checks did not validate semantic status.
- **Correction:** Reconcile completed local/native evidence with remaining hosted, live and release work; index final PR #64 head/run/merge/tree; document the independent eight-task audit and keep its limits explicit.
- **Prevention:** Review every current-status/next-work paragraph against recorded evidence, separately from historical snapshots. Finish code and documentation together; retain exact final run/merge evidence in the PR and index it in the next ledger checkpoint instead of repeatedly committing self-referential evidence. Both jobs remain required on every changed final head.
- **Recurrence check:** FN-016 also concerns stale roadmap status. FN-007 applied to overlarge audit reads and an unsupported generic job-log URL; bounded follow-up reads and the purpose-built job-log tool supplied the evidence. One local mutation-control attempt lacked Go on PATH, recurring the toolchain setup issue in FN-004/FN-012; its finally block restored source, then explicit pinned PATH/GOPATH allowed both controls to run. Generated Python bytecode was removed before the clean-checkout checks.
- **Evidence/status:** Audit of merged `6138e6e8dba4a3867c9e7b9c66c13ebaf8b2a79b`; full clean-clone race/native package checks passed. Final PR #64 logs independently confirm required jobs, all fourteen package cases and bounded documentation/vulnerability results. No new production-code defect was confirmed within the eight-task audit.

## FN-021: Explicit subscription requests could silently scan a smaller scope

- **Date:** 2026-10-01 (Europe/Oslo).
- **Confirmed boundary defect:** Subscription discovery returns the intersection of requested IDs, visible active subscriptions and filters. The coordinator previously accepted that subset, so an absent/disabled/deleted explicit subscription could disappear while execution completed. A scope hash/count and resource rows did not identify every resolved zero-data subscription. This conflicted with the target's fatal unresolved-request expectation; no historical live report is reclassified from this code finding alone.
- **Correction/decision:** Add canonical JSON requested/filter/resolved/unresolved scope evidence. Fail the coordinator before resource queries when any explicit CLI/library subscription is unresolved, preserve failed status JSON when rendering succeeds and return exit 1. Keep low-level discovery, all-visible/MG empty scopes, filter-only selectors and include-over-exclude precedence unchanged. This is a deliberate target scope-safety correction rather than a comparator normalization. Do not infer an RBAC cause from absence.
- **Prevention:** Verify scope independently of data-bearing rows, include zero-data subscriptions, preserve failed-listing versus valid-empty states, and apply privacy tests to every new ID-bearing field. Independently reconcile MG/filter-only intended membership; resolved discovery does not certify complete access.
- **Evidence/status:** Offline coordinator fixtures drive the real subscription discovery helper for absent, disabled/deleted and valid filtered cases; application fixtures assert failed JSON/exit 1 and no inventory query. A compiling removal of the guard must fail the resource-query tripwire. Canonical ownership/order and JSON scope-only redaction tests cover the new metadata. Required exact-head Linux/Windows checks remain blocking.
- **Execution recurrence:** After runtime refresh the prior /tmp clone was absent; restored a worktree from the live branch instead of trusting old local state. FN-004 path guessing recurred in two read commands; actual file inventory supplied the correct files. An initial focused test invocation preceded APRL submodule initialization and failed setup; initialized the unchanged pinned gitlink and reran the checks. No incomplete setup run counts as passing QA. Full race checks were rerun after the final code edits; toolchain flags only accommodate the known local VCS context, not production packaging.

## FN-022: SDK scope next links bypassed the application origin boundary

- **Date:** 2026-10-01 (Europe/Oslo).
- **Reproduction:** In-memory real SDK subscription pager followed a foreign HTTPS nextLink with a synthetic bearer header (`calls=2 foreignBearer=1`). No real token/network used. Advisor-specific safeguards did not cover SDK-owned pagers; search FN-017 when continuation or credential-boundary work recurs.
- **Correction:** Copied scope options prepend an origin guard before SDK authentication; reject unsafe HTTPS destinations and use a redirect-refusing default client. Fifteen pager cases plus named/custom origin and caller-ownership fixtures cover the boundary. Injected policies/transports remain trusted and must enforce their internal behavior; this is not a universal network sandbox.
- **Prevention/status:** Inspect actual SDK request construction and policy order, reproduce with synthetic bearer tripwires and test all pager entry points. Same-origin paging budget/cycles remain AR-02. Full local race suite passes; final native jobs remain required.
- **Review execution recurrence:** Some batched reads exceeded output limits and guessed module/plugin paths were absent (FN-004/FN-007). Follow-up reads used bounded source regions and actual file inventory. Unseen/truncated text is not treated as evidence; review claims are limited to inspected paths and explicit tests.

## FN-023: Shared ARM HTTP token scope ignored explicit audience

- **Date:** 2026-10-01 (Europe/Oslo).
- **Reproduction:** Explicit custom audience differed from endpoint; ResourceManagerScope still requested the endpoint scope, contrary to the SDK audience contract.
- **Correction:** Use configured ARM audience for OAuth scope; retain endpoint fallback when absent. Named-cloud, custom audience and fallback regressions pass. This is an intentional authentication correction with no live sovereign/private deployment claim.
- **Prevention:** Test audience and endpoint as distinct values, compare shared HTTP adapters with SDK options, and do not equate environment parsing tests with actual authentication validation. See ALIGNMENT_REVIEW for remaining custom-cloud limits.

## FN-024: SDK continuation cycles and cancellation boundaries

- **Date:** 2026-10-01 (Europe/Oslo).
- **Reproduction:** All three SDK scope pagers followed one/two-link cycles until the synthetic seven-request tripwire. Existing ARG/Advisor checks did not cover these pagers. Temporarily restoring the old stage invocation also failed a new cancellation/later-task tripwire: a nil task error could report success and start another task despite canceled context.
- **Correction:** Per-listing normalized continuation sets reject cycles without partial listing or URL leakage. Context checks between SDK/Advisor pages and before/after stage tasks preserve cancellation. Opt-in assessment timeout derives an earlier context deadline and preserves failed reports; no arbitrary default or page ceiling is assumed.
- **Prevention:** Exercise the actual SDK pager and coordinator/application report path, use independent finite call counts and bounded runaway fixtures, test healthy paging/opaque values/earlier caller deadline/zero compatibility, and inspect actual CLI preflight. See PAGINATION_LIFECYCLE for precise limits and remaining AR-02 work. Do not mistake deadline propagation for preemption of injected context-ignoring code.
- **Review note:** A combined read again exceeded the requested output limit; follow-up bounded reads supplied required stage/application context. No unseen text is treated as reviewed evidence (FN-007).

## FN-025: Terminal response could hide adapter cancellation

- **Date:** 2026-10-01 (Europe/Oslo).
- **Reproduction:** ARG and Advisor synthetic transports cancel the caller context and return non-empty terminal pages with nil error. Both direct adapters returned success before correction. This was a deterministic injected-transport finding, not evidence of a historical Azure incident. The stage runner's existing post-task check already protects overall assessment status.
- **Correction:** Check context after a successful request, before accepting/decoding the terminal response. Return cancellation with no partial adapter result. Keep request values, healthy projections, pins and normalizations unchanged.
- **Prevention/evidence:** New independent terminal-page fixtures fail on the original implementation and pass after correction. Actual SDK scope terminal cancellation already passes for all three entry points, retained as regression evidence. ARG distinct-token cancellation and per-batch/per-query opaque-token isolation complement existing cycle and healthy paging tests. Page/memory/load certification remains open; see PAGINATION_LIFECYCLE.
- **Execution recurrence:** FN-004/FN-007 recurred: guessed arc/cost/documentation filenames were absent and combined reads truncated. Actual inventory and smaller source regions supplied the needed evidence; unseen output is not counted as reviewed. A runtime refresh removed the new scratch worktree and changed the available shell name. Restored from verified live PR #69 baseline and reran the original failing fixtures before reapplying changes. No disappearing scratch run is the final QA evidence. Use supported shell and verify paths/toolchain/submodule after environment changes.

## FN-026: Explicit ARG truncation was ignored and requested page size exceeded contract

- **Date:** 2026-10-01 (Europe/Oslo).
- **Reproduction:** Actual response decoder plus Query accepted `resultTruncated="true"` without a token, returning partial rows as success on first and later pages. A synthetic matching-contract endpoint rejected source-compatible `$top=5000`; matching REST 2024-04-01 specifies maximum 1000. This demonstrates code/contract boundaries, not a historical Azure incident or observed service rejection.
- **Correction:** Request 1000 rows, model optional string truncation metadata and fail explicit tokenless truncation with no partial query result. Preserve 300-subscription batching, continuation/query text, valid empty results and earlier cancellation/cycle safeguards. Treat malformed supplied metadata as failure. Missing/null metadata retains legacy behavior and is not certified complete. Deliberate source differences are recorded in ARG_COMPLETENESS and target specification.
- **Prevention/evidence:** Pre-fix fixtures fail independently. A 2501-row paged fixture validates all IDs and stable query/scope; true-with-token, false terminal, invalid metadata, empty token and deliberate KQL limit have separate cases. Production coordinator/application fixtures route truncated inventory and Graph through the actual decoder/query client, retain failed JSON/exit 1 and healthy prior inventory, and stop later stages. A compiling removal of the guard made both fixtures return exit 0 and run Advisor; the regressions rejected it and source was restored before focused race checks. Do not infer estate visibility or live equivalence from synthetic results.
- **Review mistakes/recurrence:** Initial catalog export used canonical recommendation JSON, which omits Query, so that export could not support query screening. Corrected it to explicitly export actual Query strings before reporting operator counts. An application fixture initially guessed `stages.Inventory`, causing a compile failure; actual coordinator constant is `StageResourceInventory`. Corrected from source inventory and reran focused/full tests; no failed compilation counts as QA. A combined source/status read truncated again (FN-007); bounded earlier source reads and explicit own-file edits supplied relevant evidence. Check JSON tags and declared constants before using generated audit data or writing fixtures.

## FN-027: Branding review read-path recurrence

- **Date:** 2026-10-01 (Europe/Oslo).
- **Review execution mistake:** Guessed internal/renderers/renderers.go while tracing renderer propagation; that file does not exist. This repeats FN-004. Actual file inventory then supplied application and concrete XLSX/SARIF/JSON/CSV/table paths. No missing-file output counts as reviewed source.
- **Correction/prevention:** Inventory before path-specific reads, report exactly which fields have consumers, and distinguish reserved struct fields from implemented customization. Use bounded reads and actual producer/verifier paths rather than assuming package identity follows CLI defaults.
- **Publication guard mistake:** Compared Python Unicode code-point length with JavaScript UTF-16 string length for workflow text containing an emoji. The mismatch guard aborted before any remote mutation. Corrected to code-point counting; exact local/remote Git tree equality remains the authoritative byte-preservation check.
- **Design wording correction:** Initial review implied that stable SARIF fingerprints prevent identity changes from rebranding. Driver.name is also a consumer tool identity; official GitHub SARIF guidance distinguishes tool names and fingerprints. Corrected before merge: preserve emitted identities but require separate consumer migration validation, without claiming cross-tool alert continuity.
- **Research limit:** Version-specific pkg.go.dev linker documentation was unavailable; used the installed pinned Go 1.26.8 linker source contract instead. A general Go release-notes fetch did not establish a design requirement and was excluded from evidence. No production defect or custom-brand security incident is claimed by this review.

## FN-028: Null schema version retained the initialized default

- **Date:** 2026-10-01 (Europe/Oslo).
- **Reproduction:** New rejection fixture for `{"schemaVersion":null}` failed because Go JSON unmarshalling null into an initialized integer leaves that integer unchanged. Initial schemaVersion was 1, so validation incorrectly accepted null.
- **Correction/prevention:** Reject null explicitly before accepting schemaVersion; retain independent parser and actual-builder null rejection tests. Defaults apply to omitted presentation fields, never null. The pre-correction failing run is not passing evidence. No deployed/custom-brand incident is claimed.
- **Execution recurrence:** FN-004/FN-007 recurred during implementation inspection: guessed internal/app/app.go was absent and a combined read truncated. Actual inventory and bounded regions supplied the relevant source; missing/truncated text was not treated as inspected. Use inventory first.
- **Local environment limit:** Strict stamped builder testing in the scratch Git worktree failed VCS discovery. Pinned Go traced Git status from /workspace rather than recognizing the worktree .git file. Preserve strict builder/native CI stamping, test locally in an ordinary clone, and do not weaken packaging safeguards or count the failed stamped build as passing evidence.
- **Publication preparation recurrence:** A combined JSON export of eighteen changed files exceeded tool output limits; JSON parsing failed before any remote mutation. Re-export file contents in bounded chunks and verify exact local/remote tree equality before publishing. The failed export is not publication evidence.

## FN-029: Branding report fixture preparation mistakes

- **Date:** 2026-10-02 (Europe/Oslo).
- **Compile failure:** The synthetic Resource literal initially included SubscriptionName, a field absent from assessment.Resource. Remove it using the actual type declaration; the failed compile is not QA evidence. This repeats the declaration-check lesson in FN-026.
- **Checker failure:** The raw output artifact glob also matched the separately saved raw.stdout.json, incorrectly counting sixteen instead of fifteen production report artifacts. Exclude the explicit test stdout-evidence suffix, then retain an exact fifteen-artifact check and twelve CSV tables. The original failed run does not establish a product defect.
- **Inspection recurrence:** A combined result/fixture/renderer read truncated (FN-007). Later bounded actual type and renderer regions supplied the needed evidence; unseen output was not counted as reviewed.
- **Review correction:** A synthetic AZQR GitHub guidance link would trigger the existing broad legacy-branding source guard. Use a fixed Microsoft guidance URI while retaining independently checked APRL/AOR/CUSTOM labels/IDs. No production guard exception or source-pin change was needed. Synthetic provenance is not live/catalog parity evidence.
- **Prevention:** Keep fixed non-empty fixture data, independent ZIP/XML and byte/structural comparisons, explicit two-field SARIF normalization, five evidence mutation controls and exact-head Linux/Windows execution. Do not normalize away changed assessment rows or claim SARIF consumer migration from stable emitted fingerprints.

## FN-030: Scratch refresh and package integration preparation

- **Date:** 2026-10-02 (Europe/Oslo).
- **Environment recurrence:** Initial /bin/bash/current scratch worktree was unavailable after a runtime refresh. Recovered using the supported bash entry point and a new ordinary clone, verifying live branch `5c0ec9e5b3d12340c7113c153df99db95f8df0e1`, pinned submodule/toolchain and prior PR #74 native results. Do not trust disappearing scratch results or relax VCS stamping.
- **Read-path recurrence:** A cache listing incorrectly used a path relative to the new repository; corrected to the verified absolute cache. No missing listing counts as evidence.
- **QA limits:** Unit-only package execution without --binary intentionally skipped the installed case; those exploratory runs are not full package acceptance. Final local/native default and custom runs must provide real clean stamped executables, run all nineteen cases without skips and inspect installed identity.
- **Review correction:** Importing the new profile helper could create untracked bytecode before the strict clean-checkout check. Disable helper bytecode creation rather than weakening cleanliness or adding ignored debris. Reconcile current branding status separately from historical append-only checkpoint wording.
- **Cross-platform review correction:** Installed CLI text decoding used the host locale, which can misdecode Go UTF-8 Unicode profile text on Windows. Make the actual CLI subprocess decoder explicitly UTF-8; retain Unicode custom native acceptance rather than assuming a UTF-8 locale.
- **Confirmed workflow failure-masking risk:** Windows multi-command Python steps previously relied on the final external exit code; a failed earlier branding check could be hidden by a later successful report check. Add immediate LASTEXITCODE checks after each Windows branding/package command and verify the actual snippets with a synthetic first-command failure before native CI. This is a QA control correction, not an observed historical Windows failure.
- **Resolved-profile URL correction:** The initial Python helper accepted a hostname containing a space and a malformed percent escape, both demonstrated before correction. Reject invalid authority whitespace/backslashes/brackets syntax and malformed percent escapes; retain negative fixtures. This affects extraction metadata validation, not an observed Azure request or a valid Go-exported production profile. A standalone inspection import also created helper bytecode; removed only those newly generated files and use disabled bytecode for controlled helper imports.
- **Source coherence correction:** Generated INSTALL.md initially read implicit HEAD rather than the already captured source revision. Pass the captured revision explicitly, matching all other committed payload reads, so a concurrent checkout change cannot mix installation text into evidence for another source commit.

## FN-031: Feature inspection path/output and early contract assumption

- **Date:** 2026-10-02 (Europe/Oslo).
- **Recurrence:** FN-004/FN-007 recurred: assumed source commands and target loader/plugin paths were absent; recursive rule inventories and combined source reads exceeded output bounds. Corrected using actual tracked entry-point paths and smaller reads. Missing/truncated portions do not count as inspected; region helpers and full embedded query semantics remain explicitly incomplete.
- **Assessment correction:** Early commentary associated YAML plugins with table outputs. Source Graph and plugin execution entry points establish the distinction: YAML recommendations join ordinary Graph findings; internal scanners return separate tables. Corrected the characterization before publication or implementation. Do not use the existing recommendation-array filesystem loader as proof of plugin-object compatibility.
- **Prevention:** Inventory only relevant top-level Go files before reads, inspect dispatch plus production execution rather than infer behavior from type names or flags, and distinguish static source characterization from executable/live acceptance. Retain inspected source pin, selection/precedence/error/report contracts and remaining per-plugin verification before implementing migration.

## FN-032: Rules reference and artifact preparation assumptions

- **Date:** 2026-10-02 (Europe/Oslo).
- **Read recurrence:** FN-004/FN-007 recurred in combined registry/documentation reads and a guessed singular diagnostic-settings filename. Correct paths/bounded regions supplied the required contracts; truncated text was not treated as reviewed. Inventory before selecting files and avoid long-paragraph combined reads.
- **Connector correction:** GitHub file fetch returns decoded source text, not necessarily API JSON. An initial JSON parse failed before source retrieval was assessed or any mutation occurred. Read the result shape; the subsequent online pinned renderer read corroborated local source behavior.
- **Reference preparation:** Correct source HEAD did not mean its APRL gitlink was populated. Initial source build failed the embed directive. Verify actual pinned submodule contents before compilation; use a separate clean temporary copy for capture without changing the retained reference. The initial failed build is not reference output evidence.
- **Native artifact preparation:** The first target binary used a generic CGO-enabled/non-trimpath build. The existing compiled-inventory checker encountered an empty Go build setting without a Value key. Rebuilt using the required CGO-disabled/trimpath/stamped native QA profile; all four actual CLI checks then passed. Do not weaken inventory checks or treat an unsupported exploratory build as the accepted candidate.
- **Prevention/evidence:** Capture expected rules from the actual unchanged pinned source executable, retain its output hash and compare the complete target bytes. Validate synthetic selection/error boundaries independently and run the real rules command through existing installed default/custom package tripwires before exact-head native acceptance.

## FN-033: Development-plan and scanner fixture preparation

- **Date:** 2026-10-02 (Europe/Oslo).
- **Tool preparation:** An extra JavaScript parenthesis caused a syntax error before any tool or mutation ran. Corrected the call; keep orchestration expressions simple and validate returned shape before interpreting results.
- **Read recurrence:** FN-004/FN-007 recurred in combined reads with insufficient output budget and an assumed `ALIGNMENT_REGISTER.md` filename. Located `ALIGNMENT_REVIEW.md` using file inventory and used bounded relevant sections; truncated output was not treated as complete review. Match the output budget to selected text rather than merely shrinking the budget.
- **Fixture correction:** The first synthetic scanner coordinator fixture configured only enabled operations. Production coordinator validation requires all operation functions regardless of stage selection, so it rejected the fixture before execution. Reviewed that contract and configured disabled-operation guards that fail if invoked. No production validation was weakened; corrected fixture passed both selection cases.
- **Build hygiene recurrence:** Direct Python syntax compilation created an untracked task-owned `__pycache__` before the exploratory CLI build. Removed it, checked clean Git state and rebuilt the stamped binary before package acceptance. Use an external Python cache prefix for syntax checks, as already noted in FN-030.
- **Prevention/evidence:** Use the actual pinned executable for scanner keys, exercise production request mapping and coordinator independently, and verify fixture dependency validation before running. The compiling mapping-removal negative control failed the expected definition assertion; restored code passed focused race checks. Logs distinguish these preparation failures from accepted feature evidence.

## FN-034: YAML recovery and contract preparation review

- **Date:** 2026-10-02 (Europe/Oslo).
- **Recovery:** The previously accepted scratch checkout was gone, so the initial workdir call could not start. Inventory confirmed the absence; recovered a fresh checkout from the verified remote PR #78 merge/tree and restored the exact APRL input from an existing pinned local copy. Unpublished scratch is not treated as durable backup. No source reference or user working tree was reset.
- **Read recurrence:** FN-004/FN-007 recurred with assumed source catalog/target recommendation/steps filenames. Corrected using inventory and bounded actual files; missing paths supplied no review evidence. Check filenames before each new package read, even when a similarly named package was reviewed previously.
- **Name assumptions:** Initial reserved-name placeholders for three internal plugins were wrong. Inspected actual source metadata and corrected them to carbon-emissions, ai-gov and region-selection before acceptance, with literal tests for all six. An initial ASCII-only plugin-name policy was unnecessarily restrictive for safe source labels; replaced it with bounded UTF-8 labels supporting spaces, rejecting controls/path separators/ambiguous outer whitespace/reserved names.
- **Portable-path review:** An intermediate Windows device guard matched multi-digit names such as COM12 too broadly. Consulted Microsoft naming guidance and pinned Go Windows implementation; corrected single-digit/superscript rules and added positive COM12 plus negative device/character/path tests. These are prepublication review corrections, not deployed defects.
- **Field capture:** Pinned source accepts recommendationTypeId but omits it from YAML conversion; source execution confirmed the omission and the target matches it. Temporary CommandPath is the sole normalized source-capture field. Do not infer projected fields from the schema declaration alone.
- **Evidence/prevention:** Literal source mapping plus actual pinned-loader capture match; compiling origin-removal and file-bound-relaxation mutations failed named assertions, then restored focused race checks passed. Directories, source pins, caller/catalog ownership and input policies are reviewed separately from successful transport/report fixtures. Invalid candidates fail before observed authentication rather than being silently skipped.

## FN-035: Source/target read-path recurrence during zone characterization

- **Date:** 2026-10-02 (Europe/Oslo).
- **Confirmed recurrence:** Several read commands mixed target package paths into the pinned source workdir, assumed nonexistent singular HTTP/stage/decoder filenames, or requested output beyond the bounded display. These failed reads were not implementation failures and did not modify either reference or target. They are not evidence of inspection. The original reference remained clean.
- **Correction/prevention:** Separate source and target commands with explicit absolute workdirs. Inventory actual tracked paths before selecting reads, check path membership before reading and limit one contract/read region per output budget. Missing/truncated text is excluded from reviewed evidence and must be re-read correctly. Repeat occurrence is a process defect; merely recognizing FN-004/FN-007 did not prevent it.
- **QA separation:** Intentional unsafe-continuation and retry-body-limit mutations were compiling negative controls, not production fixes. Both failed their named assertions, original files were restored in a finally block, and restored focused/full checks are required. Source capture ran only in the separate test copy; retained source was not altered.

## FN-036: Plugin-table fixture and format-contract corrections

- **Date:** 2026-10-02 (Europe/Oslo).
- **Patch preparation:** One multi-operation patch repeated a new file target and was rejected before mutation; another patch used a context line differing from gofmt output. Inspect actual state and issue one operation per path with verified context; failed patch output is not applied code.
- **Fixture filename assumption:** Actual application tests initially requested `.sarif.json`; production writes `.sarif`. Both raw/masked cases failed on the missing filename after earlier outputs passed. Checked the real application suffix and corrected only the fixture, then reran the affected suite. Existing filename behavior was not changed.
- **Confirmed draft validation defect:** Initial plugin sheet guard counted Unicode runes and allowed 16 supplementary characters (32 UTF-16 units). Pinned Excelize v2.11.0 `checkSheetName` uses a 31 UTF-16-unit limit. A focused case demonstrated that the draft builder accepted this invalid sheet. Corrected the guard before publication; retain 32-unit rejection and exact 31-unit positive boundaries. General online wording about characters was insufficient; read the pinned implementation and check the actual consumer contract.
- **Read-path recurrence:** A source follow-up again guessed shortened command/types paths; actual source inventory located `cmd/azqr/commands` and `internal/plugins/interface.go`. No failed/truncated read counts as inspected evidence. This continues FN-035/FN-004/FN-007, not a new source change.
- **QA controls:** Privacy, row-width and false-complete compiling negative controls failed their named assertions and were restored. These intentional failures are separated from the real sheet-validation defect and fixture error. Required final native gates still apply after correction.

- **Second confirmed draft consumer mismatch (FN-036):** Lowercasing sheet names is weaker than pinned Excelize's Unicode EqualFold lookup. `Coſts` passed the draft reserved-name guard but aliases core `Costs`; the consumer can reuse an existing sheet index. Focused rejection case demonstrated the defect before correction. Use the same Unicode fold relation for reserved/dynamic sheets and duplicate headers, with reserved-name and cross-plugin regression cases. No draft was published or exposed in CLI before this correction.

- **Draft ordering correction (FN-036):** Sorting composed CSV keys can reverse prefix-related plugin names because the separator sorts after a hyphen. Compare plugin ID and then table ID explicitly; retain the `zone` versus `zone-mapping` case. This affects generic framework order, not the source row order or any published adapter.

- **Privacy review correction (FN-036):** Draft human masking used MustCompile for a dynamically sized identity pattern. Although IDs were quoted, excessive identity volume could exhaust construction or panic; this was a review risk, not an observed production panic. Restrict collected values to valid UUIDs, bound distinct plugin-mask identities to 4,096, compile with controlled error handling and fail before replacement rather than fall back to raw text. Build one masker per projection and retain a literal 4,097-ID rejection/raw/default-compatibility case. Existing JSON/global-volume limitations are not claimed resolved.

## FN-037: Unpublished code loss and recovery-process corrections

- **Date:** 2026-10-02 (Europe/Oslo). **Confirmed:** Runtime calls stalled; later inventory found the B3c scratch workspace absent. Cause/timing of disappearance is not established. The design checkpoint survived remotely but exact unpublished edits did not. This is not a demonstrated product defect.
- **Process gap:** Substantial coherent work was not remotely code-checkpointed before interruption. Recovery notes support reconstruction, not exact restoration. Preserve verified remote code plus documentation after each coherent slice and before task changes/session end; explicitly report unavailable publication and local risk. Root AGENTS.md and SESSION_HANDOVER now enforce that cadence during active work.
- **Recovery setup mistake:** Started fresh-clone tests before APRL submodule initialization; embedded input was missing. Restored exact pin, then full race/vet/build/actual CLI/documentation passed. Verify provenance and required embedded paths before compilation. No failed setup was claimed as product assertion evidence.
- **Publication helper correction:** create_branch on an existing checkpoint returned definite 422 with no ref mutation; used the separate non-forced update_ref after read-back. Check exact tool operation rather than treating description wording as proof of update semantics. Remote resulting identity/tree/parent/file and PR head were verified.
- **Read recurrence:** Missing workdir and truncated output supplied no inspection evidence. Re-enumerated actual directories and used bounded contract reads. Continue FN-004/FN-007/FN-035/FN-036 prevention rather than assuming recognition alone prevents recurrence.
- **Handover preparation:** Copying the retained recovery document introduced an extra blank line at EOF; staged whitespace review found it, corrected before commit. Authored content is retained with normalized final newline.
- **Evidence:** PR82 checkpoint `2d413cca968784bfb7613490283bb933e3653894`, retained ZONE_EXECUTION_CHECKPOINT; fresh SESSION_HANDOVER and handover QA. **Status:** Recovery reconciled and process corrected; lost B3c implementation still requires reconstruction and full acceptance.

## FN-038: Reconstruction read and metadata preparation corrections

- **Date:** 2026-10-02. A combined follow-up again requested a target zone path from the source workdir; failed read supplied no evidence. Re-read the actual tracked target path in the separate workdir. Larger combined documentation/code output was truncated and was not treated as a complete review; bounded follow-ups supplied the relevant contracts. Continue path-membership/output-budget checks from FN-035.
- **Metadata draft correction:** Initial pending table reused plugin metadata description as sheet description. Actual source mapping.go defines different literal text. Corrected before publication and added independent literal checks for both, source headers/author/version/license and cell values. This is an unpublished preparation defect, not a deployed semantic regression.
- **Verification process:** Local multi-command verification initially relied on final command status; inspected each focused/full log, mutation result and empty vet output independently. Use blocking shell/error checks for dependent validation/publication steps. No failed test was relabelled successful.
- **Prior postmerge correction indexed:** PR83 local proposal and connector-authored commit had equal trees but different SHAs, so local fast-forward correctly rejected divergence. Preserved proposal, switched to new task-owned branch at fetched accepted merge, verified tree/clean status; no force/reset. Prefer new accepted checkout over assuming connector commits descend from local commits.
- **Control:** Deliberate inventory-enabling mutation compiled and failed the named plugin-only DiscoverResources tripwire; restored code passed full race/vet. This intentional failure is distinct from the metadata/read corrections.

## FN-039: CLI reconstruction fixture/path preparation corrections

- **Date:** 2026-10-02. Guessed app/app.go during inspection; inventory located app/scan.go. Failed read supplied no evidence. Continue tracked-path checks from FN-038; recognizing the prior rule did not prevent recurrence.
- **Source-file metadata fixture:** Adding discovery SourceFile correctly changed the in-memory plugin object, but initial literal expected object omitted the new path. Focused test rejected the mismatch; updated explicit expected fixture path without altering original source projection capture.
- **Report fixture assumptions:** Initially guessed underscore CSV separator; production uses dot plus table key. After correcting from existing CSV code/tests, the initial XLSX comparison still omitted three branded preface rows. Corrected expected branded title/sheet/header positions and full cell comparison against actual renderer contract. These were fixture defects, not production changes; restored tests passed. Preserve format/layout inspection before new artifact assertions.
- **Verification:** Intentional unused-discovery mutation compiled and failed named pre-credential assertion; restored full race/vet passed. Local temporary paths/fixture errors are not live evidence. No unverified CLI proposal is accepted or merged before exact-head native checks.

## FN-040: Repeated source/target path mixing during service-health characterization

Date: 2026-10-02. During continuation, several reads guessed absent paths (assessment filters, target/source graph files, Azure credential filename) or mixed source-only models into the target workdir. Those failed reads are not evidence; actual inventories and explicit separate workdirs supplied the correct files. This repeats FN-038/039 rather than a new product defect. Prevention: inventory first, constrain each read to files returned from that specific repository, and never append speculative unrelated paths to a useful command. Large combined output was also truncated; excluded portions must be reread in bounded sections before reliance.

The first actual source-capture test failed compilation because git archive does not include APRL submodule contents. Restored the verified pinned APRL archive in the temporary source copy and reran successfully; no assertion/source-parity success was claimed from the setup failure. Verify embedded submodule assets before source builds. Retained reference unchanged; temporary transport-only helper is synthetic characterization infrastructure, not a source edit or live scan.

Service-health bounded POST test initially assumed an oversized 503 body would proceed to retry-delay deadline as the GET case does. Focused QA rejected that expectation: the SDK POST body-download policy encountered the byte guard first and failed immediately, with exactly nine bytes read/one closure/one attempt. Corrected the independently asserted expectation to that bounded immediate failure, and added a separate small-body 503 case for actual retry-delay cancellation. No production guard was weakened. Check method-specific middleware ordering before borrowing a failure expectation across methods; retained failed log precedes restored full QA.

## FN-041: Command-name documentation inferred from compressed context

Date: 2026-10-02. The PR86 specification reconciliation called the accepted command `zone`, relying on compressed continuation wording. Actual Cobra Use, source registration and built tests use `zone-mapping`; no alias `zone` exists. Corrected current documentation and final PR prose after direct code/test inspection. Command/feature descriptions must distinguish shorthand from executable syntax, and documentation structural checks alone do not certify every command name. During the same integration read, an absent orchestration `zone_execution_test.go` was guessed rather than taking `zone_test.go` from inventory; failed read excluded and actual file inspected. This repeats FN040: inventory in the correct workdir before choosing file paths.

Service-health integration's initial include-precedence expectation was wrong for a single scanner key: the characterized source/target rule lets a scanner-specific command win; include.resourceTypes overrides normal multi-key/default selection. Inspected existing SelectedKeys and pinned LoadFilters conditions, corrected the include fixture to exercise default/all-key selection and retained scanner-specific coverage; production precedence is unchanged. The failed targeted log is retained. A subsequent combined target read appended source-only models/filters and guessed result/json.go; excluded those failed reads. Enforce explicit absolute source paths and existence-checked per-repository reads instead of relying on a workdir label alone.

During accidental-stop QA, the added concurrent coordinator test initially assigned the scanner method directly to an operation with a concrete filter pointer signature. Go correctly rejected interface-versus-concrete function type assignment; a wrapper adapts the existing compatible filter value, with no production API change. Restored focused race checks passed. Also replaced guessed renderer paths with the actual inventory; absent-file reads remain excluded from evidence. The compiling metadata-boundary mutation separately failed its named assertion and was restored; it is an intentional negative control, not a product failure.

## FN-042: SQL EOL fixture type/tag and bounded upload corrections

Date: 2026-10-03 (Europe/Oslo). The new truncation fixture initially assigned bool instead of the existing ARG Response pointer-to-string field; Go rejected compilation. Corrected to the actual shared schema. The source capture decoder initially omitted sheet_name JSON tag, so complete metadata/table comparisons failed despite matching rows; added the actual source tag after inspecting retained capture/service test. No product guard/oracle was weakened; restored focused race tests passed. Inspect shared types and serialized field names before borrowing fixtures, and separate setup/compile failures from product assertions.

Service integration publication attempted one large JSON output; the tool truncated it and JSON parsing failed before any GitHub mutation. Reconstructed full file content in bounded chunks; published tree matched exact local Git tree, then identity/parent/ref were verified. Use bounded reads for large ledger/handover uploads. A later source read again guessed the registry under scanners/plugins and a target read guessed credentials.go; inventories established actual internal/plugins/registry.go and azure/credential.go. Failed reads supplied no evidence. This repeats FN040/041; constrain reads to inventory rather than appending guessed paths.

## FN-043: Empty fragment endpoint validation edge case

Date: 2026-10-03 (Europe/Oslo). Verified current SQL draft and accepted service-health constructor accepted https://host# because url.Parse reports an empty Fragment. Appending the ARG path to the original string would place the intended path after the fragment delimiter. Independently added constructor regression cases failed before correction. Reject any literal # in these root origins before HTTP construction; focused restored tests pass, and native gates remain mandatory. No endpoint call to Azure occurred, no credential exfiltration/resource write is demonstrated. Parsed-URL zone construction is not affected by this raw concatenation mechanism. Review delimiter presence as well as parsed values before raw URL concatenation; prefer structured construction when a future shared helper is justified.


## FN-044: SQL integration fixture validation and repeated read-path correction

Date: 2026-10-03 (Europe/Oslo). The initial SQL coordinator filter fixture supplied the unqualified resource group name unrelated; existing preflight correctly rejected it. Corrected the test to a full subscription/resourceGroups ID, preserving the production validation and SQL subscription-only behavior. Restored focused/full race and vet passed. Inspect config validation before constructing orchestration fixtures; a setup failure is not a product semantic defect.

A target read again guessed assessment/resource.go and result/assessment.go; inventories provided actual types.go and the needed existing files. These failed reads supplied no evidence. This repeats FN040-042: choosing one actual inventory path does not justify appending another guessed path. A larger combined documentation read was truncated; bounded reads established the relevant current contracts. Use inventory-selected bounded reads consistently. Intentional compiling metadata-boundary mutation was detected and restored; it is a negative control, not an accidental product change.


## FN-045: Carbon source-capture transport setup and read/tool recurrence

Date: 2026-10-03 (Europe/Oslo). The initial synthetic custom cloud lacked ARM audience, and the unchanged SDK client factory rejected configuration before a request. After setting audience, dynamically replacing http.DefaultTransport did not establish the SDK's synthetic TLS trust: the test retried and failed x509 validation. Corrected the temporary client-options Transport seam explicitly and added 20-second capture contexts/45-second test timeout; actual unchanged carbon scanner/calculation captures then passed in 0.009s. Only temporary characterization HTTP setup changed; retained source stays clean. Do not claim source execution success from setup/TLS failures; wire the actual SDK transport and bound fixtures before starting retrying calls. A later attempted signal found no matching PID in that tool's namespace; the original failed test completed and its final log/session was read back, not reported cancelled.

Read-path recurrence continued during source/target/SDK investigation (target execution plan in source workdir, absent endpoint/client/response_types/scanners paths, SDK models in source workdir). Each failed read was excluded and corrected using bounded actual files in its own workdir. FN040 prevention remains necessary: separate each repository/module call and select only returned inventory paths; an inventory printed after a failed guess does not satisfy inventory-first. The SQL branch tool call initially used repo_full_name where repository_full_name was required, rejected during argument binding before mutation; corrected operation succeeded and exact ref readback followed (PR89). Inspect tool argument contracts rather than copying a different operation's key.


FN045 recurrence follow-up: a later SDK inspection appended target-only internal/azure/cloud.go in the SDK workdir. That read failed and was excluded; the actual target file was read separately. This confirms the inventory/workdir safeguard was not yet followed consistently, rather than a new product defect. Future reads must keep target/source/SDK calls separate and include only known files for that workdir.


## FN-046: Stale current-state wording found during final documentation QA

Date: 2026-10-03 (Europe/Oslo). Final bounded first-paragraph review found ROADMAP still labelled PR81 as the current baseline, while later appended checkpoints correctly recorded PR90. Its original percentage table was also under an unqualified progress heading and could be mistaken for a fresh assessment despite later implemented scanner/plugin features. Corrected the top baseline/date, explicitly labelled retained estimates historical (no invented replacement percentages/ETA), and reconciled IMPLEMENTATION_PLAN's immediate next task/DV-001 deferral. This was documentation drift, not a runtime code change. Structural link/flag checks alone did not detect it; inspect current-tense first summaries separately from chronological evidence before each handover. The correction creates a new final proposal head, requiring both native jobs again; old-head success is not new-head acceptance.


## FN-047: Execution workspace offline during postmerge reconciliation

Date: 2026-10-03 (Europe/Oslo). Environment failure, not a demonstrated product defect. PR89/90/91 exact native acceptance and remote merge/tree/parents/ref were verified; final local fetch/readback stalled. Stopped the pending orchestration sequence and read actual state before retry. Terminal then returned environment_offline / Environment is not connected (409); GitHub remained available. Local final fetch/switch/pin acknowledgement is missing, so it is not accepted evidence. PR91 readback confirmed its unacknowledged final postmerge note had not been written. Preserve remote accepted code plus this outage checkpoint, restore workspace, inspect current state and only then resume local operations. Code checkpoint/handover protects reconstruction, not automatic workspace reconnection. No unpublished carbon request code was at risk.

Verification precaution: future-dependent mutations must require a terminal's completed zero exit, not merely a yielded session or an awaited tool object. The aborted sequence did not publish an unverified clean-fetch claim. Earlier yielded fetches were subsequently completed/read back successfully; do not generalize them to this outage.

## FN-048: GitHub content fetch returned decoded Markdown

Date: 2026-10-03 (Europe/Oslo). The outage read attempted to JSON.parse a GitHub contents response assuming raw REST base64 metadata, but the connector returned decoded Markdown. Parsing failed before any mutation. Read the exact accepted blob URL, verified Markdown shape and used returned content directly. Do not infer output shape solely from URL; inspect the connector's actual structured content before decoding. No repository data or test oracle changed due to the failed parse.


## FN-049: Carbon candidate warning health validated before normalization

Date: 2026-10-03. Initial library tests found that adding access/malformed warnings to a still-completed table before canonical validation caused valid rows to be discarded as invalid output. Normalize completed-with-warnings/failed health before validation, preserve denied/prior/later valid rows, and retain independent canonical assertions for all access/failure states. Focused race rerun passed; no affected candidate was accepted. The new exact request fixture also counted nine top-level fields where the pinned source/REST query has eight without skipToken; corrected the fixture count after reviewing the actual schema, with all field/value assertions retained. A digits-only UUID uppercase test could not exercise aliases; use distinct uppercase/lowercase alphabetic UUIDs. One orchestration call had a JavaScript quote syntax error before any tool execution, corrected without file mutation. Recurrence: inspect semantic health and literal schema, not just rows/counts or a green compilation.

FN046 recurrence follow-up: current TARGET_SPECIFICATION still described branding package acceptance as pending despite verified PR75, and its plugin paragraph only named the older zone milestone. Reconciled current summaries to accepted bounded branding/service-health/SQL evidence without closing runtime themes, live validation or release. Review first summaries across specification, plan, roadmap and handover together; link validation alone does not identify stale meaning.


## FN-050: Carbon execution boundary and fixture review

Date: 2026-10-03. Before acceptance, review found a display-value positivity check would reject valid source previous emissions that round to 0.00. Accept normalized nonnegative display text and retain a tiny-positive source regression; source negative previous values remain blank. Focused projection race passed. New orchestration test preparation assumed exclude.resourceTypes existed; compilation rejected that nonexistent field. Checked the actual config schema and removed the unsupported fixture rather than adding a new feature. Supported include/scanner tests retained, focused race rerun passed. Older availability tests intentionally changed as carbon becomes a candidate executable; preserve ai-gov rejection and identify YAML rows by name rather than assuming registry index zero. These are candidate review/test errors, no affected code accepted.

Read/output recurrence (FN040/044/045/009): one pre-interruption guessed sqleol.go read failed (actual sqleol_projection.go), and one source-command inventory used the target workdir before correcting it. This recovery also guessed nonexistent PLUGIN_EXECUTION.md; inventory confirmed the actual execution docs. Several combined read outputs truncated; split relevant reads and inspect actual small contracts rather than assuming omitted output. No edit or acceptance conclusion depended on those failed reads. Prevention remains inventory-first, explicit pinned-source versus target workdirs, bounded per-read output and exact returned-content checks.


FN046/FN050 follow-up: final public-carbon review found ACCESS_MODEL still claimed all plugin execution unavailable and omitted the plugin-only Graph exception. Reconciled current operator guidance to accepted zone/service/SQL and candidate public-carbon state, with read-oriented endpoints and separate aggregate/access/live limits. A target-doc read was also mistakenly batched in the pinned-source workdir; the failed paths were then read with the explicit target workdir before any conclusion/edit. Keep source inventory/reads and target documentation in separate calls; check the workdir at each call, not only initial inventory. The changed final documentation head requires fresh mandatory native checks; earlier CI is not acceptance.


## FN-051: Execution service disconnected after verified protected merge

Date: 2026-10-03. PR94 final native checks and remote postmerge identity/tree/parents/ref passed. The subsequent local fetch/switch command lost exec-server transport; its completion is unknown. A read-only retry stalled and was terminated. This is an environment/recovery failure, not evidence of a product assertion failure; retain both proven QA and unverified local state. Recurrence of FN037/047: remote code/docs WIP and final proposal preserve reconstructable accepted content; never claim an unpublished patch or completed fetch survives. Resume by inspecting actual status/branches before retry, not resetting or assuming the target branch exists. No source/Azure mutation occurred.

FN050 read follow-up: source SDK paths were absent from the selected target cache and a guessed internal/az/client.go read failed; actual source inventory located client_options.go/http_client.go. No contract conclusion depended on the failed reads. Do not infer SDK API versions from current REST search results; inspect exact pinned SDK. Preliminary AI-source observations remain separate from captures/acceptance.


FN051 recovery resolution (2026-10-03): executor is available but the recent checkout disappeared. A separate clean clone from verified PR94 plus exact PR95 fetch now completes local head/tree/parents/pins verification and fresh full race/vet/strict built CLI/docs QA. The original fetch/switch completion remains unknown; no user working tree was reset. Keep remotely verified code/docs as the reconstruction baseline rather than trying to recover a nonexistent directory.

FN046 recurrence: supplemental TARGET_SPECIFICATION carbon paragraph still said HTTP/public execution unimplemented despite its earlier current summary correctly indexing PR93/94. Corrected the whole canonical feature paragraph, preserving live/access/load boundaries. Review current-tense semantic claims throughout canonical documents, including appended feature summaries, not only the first paragraph or structural links. Refresh exact native gates on the updated final docs head. One combined read exceeded the aggregate output budget; split the omitted QA_PROCESS read and inspect it fully before relying on it (FN009/FN050).


## FN-052: AI source-capture transport did not confine SDK requests

Date: 2026-10-03. The first temporary source capture set http.DefaultTransport, assuming the Azure SDK used it. azcore runtime owns its own defaultHTTPClient/transport, so one read-oriented metrics POST reached the public regional service using the explicitly synthetic token and received HTTP401 AuthenticationFailed. No real credential, tenant token or resource-write operation was used; this is a confirmed development harness failure, not live feature evidence. Stop the harness immediately and inspect the exact pinned SDK transport implementation before retrying. The corrected temporary SDK copy adds only a runtime.CaptureDefaultHTTPClient hook; every original SDK file and the scanner bytes are verified identical. Install a client routing only approved fixture origins to loopback and assert foreign-destination rejection before any capture. The corrected capture completed zero and observed Graph, metrics and ARM requests only on the local TLS server, with exact methods/versions/scopes/window/body assertions. Never infer transport confinement from changing net/http globals.

The first hook preparation also failed because a copied module directory was read-only. Its shell would have continued into the old harness; the command was interrupted (exit130) before requests, with an empty capture log and absent hooks verified. Make preparation and verification separate completed steps and use fail-closed shell execution. Preserve these failures rather than rerunning without classification. Initial wire fixture used values where the pinned SDK uses per-resource value, and a row-index assertion assumed lowercase deployment preceded Unknown; corrected from SDK schema and source sorting without changing production or source expectations.

FN009/FN040/FN050 recurrence: combined long documentation/inventory output was truncated and several source/target/SDK helper paths were guessed incorrectly. Exact needed source/SDK functions were subsequently read in bounded outputs from verified inventory before conclusions. Keep source and target calls separate, avoid recursive rule-file inventories when locating Go entry points, and set aggregate output budgets deliberately. Source capture JSON uses sheet_name/table, not guessed PascalCase keys; inspect returned keys before decoding.


## FN-053: New source captures omitted byte-preserving Git attributes

Date: 2026-10-03. PR96 first native run `37087387963`, Windows job `111100395209`, failed TestSourceCaptureHashes/source-all.json. The new AI testdata JSON lacked the -text attributes already used for the other independently captured fixtures, allowing Windows checkout newline conversion. Linux/source bytes/hashes were unchanged. Add `internal/plugins/aigov/testdata/*.json -text`, including both captured outputs and literal wire inputs; retain the exact hashes and assertion. Reproduce a core.autocrlf=true isolated checkout and verify all six restored files byte-for-byte before fresh native gates. Do not normalize hashes, skip Windows or rerun the same candidate without a correction. This is a packaging/provenance oversight, not a source-cell mismatch or flaky test.


## FN-054: AI request candidate retention and test preparation

Date: 2026-10-03. Post-interruption review verified live core-v1 at PR96 merge 9c6d0e79d05d06559fb4b543406b928013aed835 and completed the earlier focused test session (zero). The request slice was still unpublished and unaccepted. Review found enrichment aggregate/text failures could occur before metrics were projected, losing otherwise valid received points. Retry metric-only projection on enrichment projection failure, retaining the original failed health and source-default metadata rather than claiming full enrichment. The dedicated regression and full/native acceptance remain required. Initial request compilation removed an unused pending table; new test compilation exposed an incorrect helper name (excludeAll versus actual excludeFilter). A fixture assertion expected seven decoded points, but literal capture inspection shows eight, including two incomplete points; corrected the count and final status index without changing capture bytes or expected source cells. A combined documentation read again exceeded output budget (FN009): omitted content was not validation evidence. These are confirmed candidate development mistakes, not accepted behavior. Prevention: inventory actual helpers, inspect literal fixture entries, test partial-data invariants across projection layers, require completed process status and remotely preserve coherent WIP before extended work.

FN054 verification follow-up: independent aggregate fixture now proves two accounts retain all18000 count-one points each when enrichment plus repeated metric identity text exceeds the pure budget, with default metadata and original failed health. Global entries/bytes/workers and later-batch/cancellation paths also pass focused race. Three compiling guard mutations fail named assertions, restored tests pass, and5000x decoder fuzz passes. A private explicit transport seam now exercises production constructor audience selection; manual injected-client setup alone would not have tested that wiring. Unknown cloud names/partial custom config and ambiguous account paths are rejected before credentials. No source fixture bytes or acceptance guards were weakened.

FN054 constructor follow-up: the new production-audience test initially guessed management.azure.com as the ARM OAuth audience, causing its fake credential to reject deployment authentication. The target correctly uses shared ResourceManagerScope, grounded in pinned azcore ARM runtime configuration (management.core.windows.net audience; management.azure.com endpoint). Corrected the literal scope expectation to management.core.windows.net/.default and required completed health, both actual requests, both scopes and closed bodies; production scope logic/source captures unchanged. Endpoint and token audience are separate contracts: inspect SDK cloud service configuration instead of deriving one from the other. The failure remained a candidate test failure and is not acceptance.

FN046 semantic-documentation recurrence: IMPLEMENTATION_PLAN current introduction still directed source characterization/recovery and its final current-labelled paragraph said AI unimplemented after accepted PR96. Reconciled the current summary, labelled the old outage paragraph historical and added the actual request-to-public integration order; TARGET_SPECIFICATION now distinguishes accepted pure core from unaccepted request/public obligations. Local link checks cannot detect stale semantic status. The final documentation head must obtain fresh native checks rather than borrowing earlier CI.

## FN-055: PR97 ignored a lone ARM audience override

Date: 2026-10-03. Fresh-session review read live refs and the actual automated review comment [4171403211](https://github.com/DeBoX85/Cloud-Assess/pull/97#discussion_r4171403211), then inspected the constructor and shared cloud helper. Final candidate 563c973 had passed both native jobs, but only authority/endpoint presence was compared. The shared helper discards a lone AZURE_RESOURCE_MANAGER_AUDIENCE override, so the constructor could return a public-cloud scanner for a partial/custom configuration. This is a confirmed unaccepted-candidate validation gap; no target scan or real credential was used in this session.

Require all three overrides together or none, before reading the cloud helper or constructing clients. Existing exact public authority/endpoint/audience checks still reject unsupported complete configurations. The regression enumerates all six partial combinations using public-looking values, requires constructor rejection and zero credential/transport calls, and accepts the complete public triple. Extend the isolated mutation gate with a compiling removed-audience-presence guard that must fail the named regression. Existing source fixtures, pins, dependencies and normal public behavior remain unchanged. Fresh exact-head native jobs and log inspection are mandatory; earlier green CI cannot certify this correction.

## FN-056: Fresh Windows executor lacks the previous development toolchain

Date: 2026-10-03. The historical Linux workspace paths do not exist here. Bundled Git 2.53.0 reports no remote-https helper; task-owned clone attempts did not produce a checkout. Go is absent from PATH and inspected runtime/common locations. Shell HTTPS to go.dev is blocked. The Code Review plugin is installed with its MCP disabled and exposes no callable review tool, so GitHub code/diff/review reads provide the review evidence. No credentials, browser session or protection changes were attempted. Use authenticated GitHub objects to publish this focused correction and required native Actions for executable QA. Do not claim local race/build or clean fetched verification in this executor. Substantial next feature work needs a restored development toolchain or an authorized equivalent execution workspace; raw Azure/laptop validation remains separately deferred.

Recovery-helper failure from the previous PR97 body is retained: an isolate helper was absent after a new message, producing TypeError after successful checkpoint publication. Re-read actual remote state and reconstruct helpers from persistent records; never retry an uncertain mutation blindly.

FN055 native follow-up: run37092489006 on correction1945da7 failed Linux111115580824 and Windows111115580893 at the new mutation gate. Production ordinary/race package tests passed, but removing the audience guard left audienceSet unused, so the mutation did not compile. The gate correctly refused to count this as detected behavior. Replace the guard with authoritySet != endpointSet || audienceSet && false to preserve the variable reference while ignoring its value. Require the named partial-cloud assertion, restored baseline and fresh complete native jobs. This is a confirmed mutation-harness authoring error, not product behavior or flaky infrastructure; no production test or acceptance condition weakened.

FN055 final proof: corrected headf58895f/run37092703311 passed both mandatory jobs111116220428/111116220521; all four mutations compiled, failed named assertions and restored baselines passed on each host. Protected mergeaadd0d03b13a093a9ec3b9e69e3f70208d1e1714 matches exact testedeeb4313c18fedb890363777e095631fcb017e78f tree/ordered parents and human identity. Earlier failed run remains recorded.

FN056 next-state clarification: all requested focused correction and protected-merge work that can be validated through existing native CI completed without operator laptop/Azure input. The missing clean local development/source checkout and test toolchain remains an execution blocker for the next substantial implementation. The separate AI_GOVERNANCE_EXECUTION contract is durable research/planning, not untested implementation or proof that public AI exists. Postmerge native push checks are separate remote evidence and do not establish unavailable local proof.

FN056 resumed observation (2026-10-03): initial read-only Get-Location/Get-Command completed with Windows cwd/bundled Git and no callable Go. Later commands failed with exec-server transport disconnected and a recovery timeout; another read-only recovery attempt stalled and was terminated. Its completion is unknown. This is an observed execution dependency, not a diagnosed product defect; prior missing HTTPS helper/download failures were not retested after disconnect. Cloud Environment runtime SKILL.md was read, but no environment_status/execution tool is exposed in this session. Plugin search/dependency metadata did not identify an available matching connection; unconnected alternative hosting/workspace services were not provisioned. Repeated speculative shell recovery does not establish readiness. Continue independent GitHub/source/schema planning and durable publication, then request a usable executor separately from deferred operator laptop/Azure tests. No credentials or permissions/gates/pins changed.

Resumed documentation reconciliation: accepted PR98 and its own successful merge push run37093451765 are now indexed in current handover/roadmap/ledger. The earlier plan's instruction to finish PR98 gates is superseded, not replayed. Required wire fields/string enum/false-with-token examples were independently checked against the immutable official 2024-04-01 API schema before specifying AI discovery acceptance; a malformed metadata page cannot be treated as healthy empty discovery. New limits and tests are a planned contract, not implementation or executed evidence.

FN056 actual local recovery (2026-10-03): user requested creating the required workspace/tools. The new Linux executor responds and permitted network downloads/Git HTTPS now work. Built a fresh isolated workspace rather than assuming older checkouts/tool paths are valid; installed exact Go1.26.8 and PowerShell7.6.6 after publisher SHA256 checks, govulncheck1.8.0, verified clean live core-v1/pinned AZQR/actual APRL/provenance and both module graphs. Full local race/vet/stamped CLI/docs/package/branding/security/mutation/fuzz/maintenance checks passed within their offline boundaries. Added the checked-in reproducible bootstrap, actual second fresh provisioning at a space-containing path and corrupt-download/preservation/valid-shell negative-control feedback. This resolves the current missing-local-development proof; previous disconnects/unknown command completion remain historical. The recipe checks readiness from observed files/tool output and refuses existing destinations. Linux recovery does not close native Windows/operator Azure/live/Gate004/release evidence or guarantee scratch uptime. See DEVELOPMENT_WORKSPACE and final tooling PR evidence before resuming discovery.

PR100 initial native run37103069229 passed Windows111146251132 but failed Linux111146251036 in the new checksum fixture before the expected checksum assertion. The fixture checked the host prerequisite inventory before exercising its corrupt archive; its assertion initially omitted captured stderr, so the specific early prerequisite failure is not retrospectively proven. Corrected offline fixtures supply private inert rg/gcc/Git-HTTPS prerequisite shims, reject unexpected Git operations, require the existing-directory mkdir failure and retain early diagnostics. Real end-to-end provisioning above remains independent evidence. Three corrected cases, a valid-shell removed-checksum control and restored guard pass locally; final exact-head native acceptance remains required. Automated review questioned attribution, but GitHub commit API and fetched object independently show candidatec906a3 authored and committed by the mandated human. Shell push lacks credentials; connector publication is verified by exact tree and ten file readbacks, without requesting tokens.

## FN057: discovery rejected-text budget omitted before publication (2026-10-03)

Self-review of the unaccepted AI discovery decoder found that wireString rejected an oversized/unsafe projected string before its length was charged, allowing invalid rows to bypass the separate decoded-text budget. Raw rows and successful bodies were already bounded, but the specified text accounting included rejected strings and unused sku_tier. Corrected the row decoder to charge decoded projected string bytes before safe-label rejection. Independently sized rejected-tier fixtures exercise exactly16MiB including owned scope labels/continuations and one-byte excess; restored focused race passes. No faulty discovery code was previously accepted or publicly executed, and pinned source/capture/dependency bytes are unchanged. Prevention: review every budget against excluded/invalid/secondary data, and pair each critical boundary with a literal accepted/excess fixture. Final broad/native/protected acceptance remains required.

## FN058: discovery Unicode fold alias before acceptance (2026-10-03)

Review of published but unaccepted discovery head87d963f found regional lowercasing could turn the Kelvin sign into ASCII k before DNS validation, and Unicode-insensitive fixed-path/type comparison could accept long-s aliases. Fixed validation now permits ASCII case changes in regions and fixed ARM names while retaining safe display-label Unicode. Named fixtures use distinct IDs so ordinary duplicate rejection cannot mask the unsafe region. An additional compiling mutation disables the regional ASCII guard and must fail that assertion, then restored cases pass. Initial CI is not final-correction acceptance; rerun all exact-head native gates. No public AI path, real Azure credential, source/capture/dependency update or accepted faulty discovery is claimed.

FN058 feedback correction: the first Unicode mutation survived because an authoring replacement did not match the escaped Go Unicode literal, leaving the region fixture with the existing good account ID. Duplicate detection masked the missing DNS guard as anticipated. The mutation script rejected this outcome; an exact patch now assigns region-alias before rerunning. No compile failure or survived control is credited as a PASS. Eight final compiling controls and restored assertions must complete.

## FN059: active acceptance sentence left pending after header update (2026-10-03)

PR102 Code Review identified an active AI discovery acceptance paragraph still instructing final native/protected acceptance after PR101 had completed it. Updating the status header/final section did not remove that contradictory active instruction. Corrected the paragraph to exact2258/run37105082689 acceptance and separate future public integration gates, and explicitly labelled earlier FN058 notes historical. Initial92c1fd native success does not certify revised documentation. Prevention: scan active authority sections for every pending/current imperative after an accepted merge, and reconcile each with exact final PR evidence while retaining clearly labelled historical snapshots.


## FN060: public AI integration test oracles used wrong existing contracts (2026-10-03)

The first58-case report run failed masked rows because the authored expected identity was *** instead of the existing xxxxxxxx-xxxx-xxxx-xxxx-xxxxx plus seven-character suffix. It also failed the healthy all-plugin case because a service-health PendingTable was incorrectly fabricated as completed without its successful Service Issues branch. Actual redactor and service projection contracts were inspected; the fixture now uses the literal established mask and a real bounded service scanner. Corrected58 cases and focused checks passed. No production masking/service behavior, expected source cells, validation guard or acceptance condition was weakened. First failures remain in task logs; candidate is unaccepted until broad/native final-head QA. Prevention: inspect actual existing contracts before creating cross-feature fixture defaults, and retain real adapters for healthy branch oracles. Earlier ai-gov-unavailable tests now target genuinely unavailable region-selection.


FN060 ownership-fixture follow-up: the first recorded-tag concurrency test expected complete despite omitting the account's DeploymentSet. Existing projection correctly marked missing enrichment as unavailable. Supplied an explicit successful empty deployment set, preserving the strict complete oracle; restored ownership/race checks passed. This is fixture contract correction, not suppressing missing-input health.

## FN061: single-branch recovery fetch did not create a tracking branch (2026-10-03)

The bootstrap clone's configured fetch mapping covers only core-v1. A Git fetch of feat/ai-public-execution populated FETCH_HEAD, so the subsequent origin-based tracking switch failed. Explicit remote-ref fetch then succeeded, but tracking setup still rejected the unmapped remote branch. Inventory proved local HEAD remained02b1ffd with the correct tree and new owned untracked test retained. Switched without tracking to fetched exact82fbda5, then verified d0f425ca/tree/direct032740d parent/human identity, preserving both local branches and the new test. No unknown switch completion, reset/force/protection/credentials or product failure was inferred. Prevention: inspect single-branch fetch mappings, fetch the required exact ref explicitly and verify the object before a nontracking task branch switch.


## FN062: region primary capture input bands and test formatting (2026-10-03)

Initial literal synthetic cases labelled band-80/band-60/below-60 actually scored60/40/39.65 in unchanged source. Corrected the explicit inputs to actual80/60/59.65 and regenerated through original source methods; never hand-edited captured output or target thresholds. The first target test invocation failed Go vet because a table struct used a %q formatting verb; changed only diagnostic formatting to %#v and reran focused race tests. No failed run is acceptance. Earlier guessed read paths and oversized batched output recurred from FN004/FN007; bounded reads of inventoried files supplied the actual reviewed functions, not omitted output. Prevention: inspect actual original scores and threshold inputs, compare every cell with independent captures, retain byte attributes/hashes and classify compilation/fixture errors before using QA claims. Pure joined-cell validation uses actual UTF-16 units rather than a conservative two-units-per-rune assumption; boundary fixtures require bounded valid rows to pass. Source production, pins and credentials unchanged.


## FN063: region aggregate budget omitted joined delimiters (2026-10-03)

Automated Code Review on PR104 head596f102 found detail separators were checked per cell but absent from the pre-projection total text count. Repeated4096-entry empty lists could evade label-byte accounting and cause large CPU/allocation work before final canonical output validation. No accepted or public region implementation exists. TestPrimaryAggregateDetailWorkRejectedBeforeProjection demonstrably failed on original code. Added65536 total detail/mapping-entry cap and counted UTF-8 detail separators/mapping arrows and separators before loops/joins. A separate256-comparison fixture with32x511-byte labels per detail column is below decoded-label budget but exceeds joined output budget; it must fail region_text_limit before projection. Compiling removal of each guard fails its named assertion, followed by restored checks. Final corrected-head native/review QA is mandatory; earlier green run37129101097 remains historical proof only. Prevention: bound global structural entry work and output delimiters, test empty/short-entry amplification and distinguish early rejection from late validation. Pure source cells/hashes, dependencies, pins and privacy guards unchanged.

## FN064: remote PR creation returned an unknown timeout outcome (2026-10-03)

The first PR104 creation call and follow-up connector reads timed out despite verified remote code596f102. Treated the response as unknown: read public GitHub REST without credentials and observed no open proposal before retrying. Retry created exactly PR104; actual head/base/body verified. No blind repeated mutation, duplicate PR, force update or credential request. Prevention: inspect uncertain remote state through permitted read-only capabilities before retry; retain verified branch/tree/code as recovery evidence even when proposal publication is temporarily unknown.

## FN065: executor stopped returning after completed correction QA (2026-10-03)

Focused/full race, vet and all five compiling region controls completed successfully. The subsequent local commit/stamped-build/package command never returned an execution result/session ID; a fresh pwd health check also did not return. Wrapper waits were ended; underlying command outcome is unknown, not assumed cancelled, successful or rolled back. Preserve existing checkout/branches/processes. Reconstruct the correction from verified remote596f102 and bounded reviewed edits through GitHub, verify complete remote code/docs/identity/tree/parent and require fresh native gates on that exact head. Do not reuse an unknown local hash/build or request credentials/laptop tests. Inspect actual files/logs/HEAD before any retry after executor recovery; source pins and remote checkpoints remain recovery authority. No alternative executor capability is exposed. This is an observed tool execution interruption, not a claim about its infrastructure cause or a product fault.


## FN066: characterization path filters could omit final-head proof (2026-10-03)

PR105 automated review identified that the dedicated workflow ran only for runner/harness/workflow paths. A later documentation-only commit could therefore lack the explicit final-head characterization required by the contract, despite earlier successful output. Removed both PR and accepted-push path filters while retaining branch restrictions, read-only permissions, pinned actions and existing mandatory gates. Require a fresh successful characterization plus Linux/Windows jobs on the corrected complete head; prior7b1fdede success is not correction acceptance. Prevention: align workflow scheduling with the evidence contract, including documentation-only follow-ups, and inspect actual tested head/preview rather than relying on previous run status. No capture/source/runtime/pin change.

Interruption recovery preserved the dirty old596f102 task checkout and fetched published7b1fdede into an isolated worktree with explicit remote mapping (FN061); actual source/pins/tools and acceptedc270ea22 objects were verified. Previous unknown commit/build completion remains unknown. The missing tracking-ref error was resolved with an explicit non-destructive fetch; no reset, force or protection change.


## FN067: accepted PR105 absent from active continuity (2026-10-03)

Live readback found PR105 merged at7e4ca5c while active handover/roadmap/source records still required its acceptance. Reverified final/native and separate accepted-push jobs/full logs, complete captured chunks/provenance, exact refs/tree/parents/identity/protection; reconciled active summaries while retaining historical proof. This repeats FN059. Prevention: refresh every active authority summary from final PR/run evidence after merge before further work.

Current executor is Windows, with Git/Python/PowerShell but no observed Go/ripgrep/prior Linux workspace or callable managed environment/Code Review tool. Git HTTPS failed to connect. Use authenticated GitHub/native Actions without bypassing network policy, requesting credentials, weakening checks or claiming fresh local Go/Azure evidence. Historical unknown local command outcomes remain unknown.


## FN068: native formatter rejected two authored alignments (2026-10-03)

Initial PR106 head4e50009de3862d68b54e28dc8c3a557b7665de96 was rejected by native Linux formatting run37146459704/job111271307391: two struct alignment spaces differed from gofmt. Applied the exact native gofmt diff. No test/gate was weakened; this head remains unaccepted and broad Linux checks were skipped after that failure. Corrected head requires fresh complete characterization/Linux/Windows evidence. Local Go/gofmt is unavailable; the retained failure diff supplies exact formatting correction. Prevention: keep formatting blocking and print its diff for bounded recovery, then validate every corrected head through native QA rather than hand-claiming format/test success.
