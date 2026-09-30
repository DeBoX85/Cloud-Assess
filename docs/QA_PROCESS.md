# QA execution and feedback loop

Applies to material changes on `bootstrap/core-v1`. These controls support Gate 004 and the later release decision; they do not create a new PASS or replace missing live evidence.

## Before changing behavior

Identify the affected specification contract, concrete failure, expected behavior and evidence needed. Search existing tests, the development ledger and FAILURE_NOTES before adding overlapping checks. Independently specify expected records/counts instead of deriving the oracle from the implementation under test or assuming AZQR is always correct. Document any pinned-source correction separately from successful-response equivalence.

## Required verification and final review

Required Linux and Windows jobs remain blocking. They now include bounded comparator/filter fuzz runs, four isolated comparator mutations, synthetic large-subscription batching and later-page failure checks, report replacement failure tests, formula-like export tests, and Windows controlled-directory ACL inheritance. Seed fuzz inputs run with ordinary tests on both platforms. Generated fuzz failures are retained in synthetic CI artifacts when present; minimized reproductions must become checked-in regression seeds after review. A passing short fuzz run is not exhaustive input validation.

The mutation script compiles an isolated standard-library comparator module. Its clean baseline must pass; each selected mutation must compile and fail the named assertions. Compiler failures do not count as detected mutations. This checks four comparator faults, not mutation coverage of the entire toolkit.

Review the exact final diff after the last edit, failed attempt or interruption. Check request endpoint/method/body semantics, context and response-body ownership, redaction, output failures and compatibility for affected paths. Diagnostics fixtures reject batch subrequests other than diagnostic-settings GETs; ARG fixtures constrain its query endpoint and authenticated POST. These are focused contracts, not a complete allowlist or proof every future adapter is read-oriented. New adapters require their own request contracts.

## Failure loop

Stop dependent merge/claims, check FAILURE_NOTES for recurrence, reproduce, distinguish confirmed cause from hypothesis, fix with a regression test, then rerun affected checks. Unexpected failures cannot be erased by rerunning until green. Inspect flaky timing and retained corpus inputs; document accepted uncertainty rather than loosening assertions to obtain a pass.

## After merge and at a gate

Record tested head, CI run and merge SHA in the ledger for material checkpoints. Update roadmap and gate evidence only to the extent demonstrated. Check documentation against final behavior. Gate review records the exact candidate and all PASS, ACCEPTED LIMITATION or BLOCKED decisions. DV-001 and other live-evidence gaps remain open without the operator's environment; do not infer closure from synthetic data.

## Report guarantees and limits

JSON, SARIF, CSV and XLSX stage each file in the destination directory, then replace it only after successful rendering and close. A rendering failure preserves an existing report and removes the staged file. Unix replacement files use mode 0600; Windows files inherit directory ACLs. The controlled Windows fixture proves inheritance in its restricted test directory, not privacy in an arbitrary operator directory. Symlink destinations are rejected; parent-directory security remains an operator responsibility.

Replacement uses Go `os.Rename`; atomic replacement on every operating system and crash/power-loss durability are not promised. Temporary files can remain after process kill. Multi-file exports are not a transaction: completed files may exist when a later output fails, and application failure/path reporting must remain honest.

CSV prefixes potentially formula-like cells with an apostrophe, including leading whitespace/control characters and fullwidth formula prefixes. This changes human CSV cell representation, including negative numeric strings. Canonical JSON preserves original values and is the machine-data interface. XLSX untrusted ordinary cells must be stored as text; intentional escaped HTTP hyperlinks remain separate. Spreadsheet consumers and re-save behavior vary, so neutralization is not a universal consumer security guarantee.

## Remaining acceptance work

Continue adapter denial/malformed-response checks, decoder fuzzing, workload concurrency/cancellation bounds and independent summary invariants where existing tests lack them. Maintenance no-change/changed-output execution, full built-CLI artifact inspection, and arbitrary-directory Windows ACL review remain open. Release artifacts additionally require clean installs, dependency/license review, checksums, build provenance verification and an approved pilot. This process does not silently choose release scope or supported platforms.

## Development guidance

- [Go fuzzing](https://go.dev/doc/security/fuzz/): deterministic targets, bounded runs and reproducible regression corpora.
- [OWASP CSV Injection](https://owasp.org/www-community/attacks/CSV_Injection): CSV quoting alone does not stop spreadsheet formula interpretation; consumer behavior needs explicit limits.
- [Go os.Rename](https://pkg.go.dev/os#Rename): replacement semantics and platform-dependent atomicity.
- [SLSA build requirements](https://slsa.dev/spec/v1.2/requirements): release provenance review; no SLSA level is claimed.
