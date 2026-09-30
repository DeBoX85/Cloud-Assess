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

## FN-008: Toolchain availability was inferred from PATH alone

- **Date:** 2026-09-30. Earlier checkpoints said local Go was unavailable after checking PATH. A later source search found an existing Go 1.26.0 toolchain outside PATH, which successfully selected/downloaded Go 1.26.8. Earlier executable verification was performed in CI, not locally.
- **Correction/prevention:** Verify known workspace toolchain locations and the actual selected version before concluding local testing is unavailable. Local tests are now possible; an initial broader run then failed setup because the APRL gitlink was uninitialized. Focused report/config tests passed; initializing the exact pinned submodule restored broader validation, and local `go test -race -count=1 ./...` then passed. Check toolchain, module dependencies and embedded gitlink contents together; continue required hosted Windows/Linux checks. **Status:** Environment availability corrected; previous CI evidence remains valid.

## FN-009: Documentation upload assumed command output was complete

- **Date:** 2026-09-30. An evidence-update tree was constructed from shell output without first verifying complete content. The ledger exceeded the output limit, so the remote tree differed from the local staged tree. Detected before creating a commit or branch; the incomplete tree was never published.
- **Correction/prevention:** Read large files in bounded chunks and require exact local/remote tree equality before publishing. **Status:** Corrected after the tree equality check.
