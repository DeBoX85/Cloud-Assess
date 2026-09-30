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
