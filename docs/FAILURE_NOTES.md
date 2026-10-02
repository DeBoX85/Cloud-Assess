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
