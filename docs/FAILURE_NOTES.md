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
- **Evidence:** CLI process-test inspection session. **Status:** Corrected; recurrence recorded in this entry.
