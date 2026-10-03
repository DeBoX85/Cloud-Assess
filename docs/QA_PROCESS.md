# QA execution and feedback loop

Applies to material changes on `bootstrap/core-v1`. These controls support Gate 004 and the later release decision; they do not create a new PASS or replace missing live evidence.

See [DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md) for ordered batches, scope boundaries, durable checkpoints, GitHub interruption recovery and reviewed rollback.

## Session continuity check

Start with [SESSION_HANDOVER.md](SESSION_HANDOVER.md), verify current refs/proposals and reconcile accepted versus pending work. Maintain that record with code and evidence after coherent slices and material changes, not only at the final response. Before ending or switching tasks, verify remote code-checkpoint read-back when available; explicitly identify unpublished work at risk otherwise. A fresh-session review must identify goal, source pins, accepted baseline, next task, mandatory checks and unresolved evidence without relying on chat.

## Before changing behavior

Identify the affected specification contract, concrete failure, expected behavior and evidence needed. Search existing tests, the development ledger and FAILURE_NOTES before adding overlapping checks. Independently specify expected records/counts instead of deriving the oracle from the implementation under test or assuming AZQR is always correct. Document any pinned-source correction separately from successful-response equivalence.

## Required verification and final review

Required Linux and Windows jobs remain blocking. They now include bounded comparator/filter fuzz runs, four isolated comparator mutations, synthetic large-subscription batching and later-page failure checks, report replacement failure tests, formula-like export tests, and Windows controlled-directory ACL inheritance. Seed fuzz inputs run with ordinary tests on both platforms. Generated fuzz failures are retained in synthetic CI artifacts when present; minimized reproductions must become checked-in regression seeds after review. A passing short fuzz run is not exhaustive input validation.

The mutation script compiles an isolated standard-library comparator module. Its clean baseline must pass; each selected mutation must compile and fail the named assertions. Compiler failures do not count as detected mutations. This checks four comparator faults, not mutation coverage of the entire toolkit.

Review the exact final diff after the last edit, failed attempt or interruption. Check request endpoint/method/body semantics, context and response-body ownership, redaction, output failures and compatibility for affected paths. Diagnostics fixtures reject batch subrequests other than diagnostic-settings GETs; ARG fixtures constrain its query endpoint and authenticated POST. These are focused contracts, not a complete allowlist or proof every future adapter is read-oriented. New adapters require their own request contracts.

## Failure loop

Stop dependent merge/claims, check FAILURE_NOTES for recurrence, reproduce, distinguish confirmed cause from hypothesis, fix with a regression test, then rerun affected checks. Unexpected failures cannot be erased by rerunning until green. Inspect flaky timing and retained corpus inputs; document accepted uncertainty rather than loosening assertions to obtain a pass.

## After merge and at a gate

Record tested head, CI run and merge SHA in the ledger for material checkpoints. Final run/merge evidence can first be retained in the PR and indexed in the next ledger checkpoint, avoiding a self-referential documentation/CI loop. Every changed final head still requires both blocking jobs. Update roadmap and gate evidence only to the extent demonstrated. Check documentation against final behavior and review current-tense status separately from historical checkpoints. Gate review records the exact candidate and all PASS, ACCEPTED LIMITATION or BLOCKED decisions. DV-001 and other live-evidence gaps remain open without the operator's environment; do not infer closure from synthetic data.

## Report guarantees and limits

JSON, SARIF, CSV and XLSX stage each file in the destination directory, then replace it only after successful rendering and close. A rendering failure preserves an existing report and removes the staged file. Unix replacement files use mode 0600; New Windows files inherit directory ACLs; replacement files preserve the existing effective DACL before report bytes are written and protect that DACL from broader parent inheritance. A missing/unreadable DACL fails replacement. Owner/SACL preservation is not claimed; preserved DACLs are snapshots, so later directory inheritance changes do not automatically apply. The controlled Windows fixture proves inheritance in its restricted test directory, not privacy in an arbitrary operator directory. A separate replacement assertion broadens the directory while the existing report remains restricted, checking that replacing it does not grant new readers. Symlink destinations are rejected; parent-directory security remains an operator responsibility.

Replacement uses Go `os.Rename`; atomic replacement on every operating system and crash/power-loss durability are not promised. Temporary files can remain after process kill. Multi-file exports are not a transaction: completed files may exist when a later output fails, and application failure/path reporting must remain honest.

CSV prefixes potentially formula-like cells with an apostrophe, including leading whitespace/control characters and fullwidth formula prefixes. This changes human CSV cell representation, including negative numeric strings. Canonical JSON preserves original values and is the machine-data interface. XLSX untrusted ordinary cells must be stored as text; intentional escaped HTTP hyperlinks remain separate. Spreadsheet consumers and re-save behavior vary, so neutralization is not a universal consumer security guarantee.

## Decoder feedback checks

The row-decoder fuzz target uses a 100,000-execution budget, a separate 60-second test timeout, two workers and an 8 KiB input bound. Seed/shape cases also run in ordinary Linux/Windows tests. Keep failures as failures; a time-budget termination without an assertion/corpus is distinct from a property counterexample and requires investigation (FN-010). Literal stage projection expectations must have reviewed source/schema provenance and explicitly distinguish fabricated rows from live captures.

## Maintenance publication checks

The tidy and pinned-rule workflows call the same tested publisher. It starts on `bootstrap/core-v1`, stages only the selected maintenance paths, rejects unrelated staged files, commits with bot identity and pushes a separate run-specific proposal ref without force. No-change output produces no branch/commit. Rejected publication must stay failed and omit a success summary. Local bare-remote tests cover the publication boundary, including gitlinks. [MAINTENANCE.md](MAINTENANCE.md) documents additional recipe/output fixtures and staged rule-hash rejection checks, together with the remaining hosted transport/token/event limits.

A branch push using `GITHUB_TOKEN` does not automatically trigger push workflows. Open a reviewable PR with an authorized identity and observe `quality` and `windows-validation` on its exact head; never equate a published proposal with checked or merged changes. No complete maintenance workflow PASS is claimed from publisher tests alone.

## Dependency evidence checks

The [inventory generator](DEPENDENCIES.md) collects versioned module checksums and full root license/notice texts plus bundled-source pins. Nine isolated safeguards and a required Linux freshness check reject missing evidence or stale output. Preserve original notice bytes rather than stripping third-party whitespace; generated NOTICES.md has a path-specific Git exception only. Stage new authored files and run a blocking staged diff check before publication. Regenerate and review dependency-refresh proposals; do not treat this snapshot as release-specific licensing approval.

## Built artifact checks

[BUILT_CLI_VALIDATION.md](BUILT_CLI_VALIDATION.md) distinguishes real compiled-CLI offline preflight/metadata checks from the existing injected test-executable synthetic report checks. Required jobs execute the actual CGO-disabled host amd64 artifact from a copied path with spaces. Fifteen invalid inputs must reject with the expected error before observed HTTP/authentication, preserve report/directory hashes and emit no stdout report. Compiled module membership, versions and checksums must match the inventory. These cases do not replace successful Azure execution, real package installs or release provenance.

## Candidate package checks

[PACKAGE_BUILDS.md](PACKAGE_BUILDS.md) adds strict compiled/source revision matching, committed-byte notices, archive/payload checksums and new-directory extraction validation. Both required jobs must pass actual package integration alongside integrity/provenance fixtures; a local unit-only run is insufficient. Existing directories/packages are retained on rejection. Do not infer signed publisher identity from checksums or fresh-machine behavior from a directory isolated on an existing runner.

## Remaining acceptance work

Continue adapter denial/malformed-response checks, decoder fuzzing, workload concurrency/cancellation bounds and independent summary invariants where existing tests lack them. Hosted maintenance dispatch/token/PR execution remains open; local no-change/changed-output publication and generation fixtures have passed. Actual built-CLI preflight/module inspection and candidate ZIP integrity/extraction have passed on Linux and Windows; successful installed Azure scans, fresh-OS operation and arbitrary-directory Windows ACL review remain open. Release artifacts additionally require release-specific dependency/license review, publisher provenance and an approved pilot. This process does not silently choose release scope or supported platforms.

## Development guidance

- [Go fuzzing](https://go.dev/doc/security/fuzz/): deterministic targets, bounded runs and reproducible regression corpora.
- [OWASP CSV Injection](https://owasp.org/www-community/attacks/CSV_Injection): CSV quoting alone does not stop spreadsheet formula interpretation; consumer behavior needs explicit limits.
- [Go os.Rename](https://pkg.go.dev/os#Rename): replacement semantics and platform-dependent atomicity.
- [SLSA build requirements](https://slsa.dev/spec/v1.2/requirements): release provenance review; no SLSA level is claimed.

- [Microsoft PSModulePath guidance](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_psmodulepath?view=powershell-7.5): isolate module paths when a Go/Python intermediate process starts Windows PowerShell from PowerShell 7.

- [GitHub GITHUB_TOKEN event guidance](https://docs.github.com/en/actions/concepts/security/github_token): proposal pushes and subsequent quality workflow execution are separate evidence.

## Documentation and credential-boundary feedback

Both required jobs run `scripts/tests/documentation.py`: authored local file destinations, documented Cloud Assess flags against real built help and operator PowerShell snippets through the parser (Linux PowerShell 7 and Windows PowerShell 5.1). Snippets are never executed. This detects structural drift; it does not replace semantic review, check all external/anchor links or anonymize retained notices.

Use synthetic bearer canaries when reviewing service-provided continuations, redirects and diagnostic output. Destination changes must be validated before sending credentials; do not trust a query-style request's verb alone. The current Advisor/default shared-client guards and their explicit SDK/injected-transport limits are recorded in [OFFLINE_QA_REVIEW.md](OFFLINE_QA_REVIEW.md). Native Linux and Windows jobs run reachable-vulnerability scans; record module-only advisories separately rather than calling the graph advisory-free.

Include SDK middleware defaults in read-oriented review. Shared ARM and scope options now disable automatic provider registration; registration-required error fixtures must confirm no implicit POST. An explicit GET in adapter source is insufficient evidence if its SDK pipeline can remediate by writing. Record such corrections separately from successful-response equivalence (FN-019).

## AI request boundary regressions

Both required jobs run the isolated [AI request mutations](../scripts/tests/ai-request-mutation.py): partial-cloud audience, foreign continuation, foreign metric identity and enrichment-limit data loss must compile and fail their named assertions, with healthy/restored fixtures passing. Linux adds bounded AI decoder fuzzing; all native package tests include explicit in-memory authenticated audience/request/cancellation/paging/budget checks and independent captured source cells. See [AI_GOVERNANCE_REQUESTS.md](AI_GOVERNANCE_REQUESTS.md) for precise limits and public/cloud/live boundaries. These extend development feedback without changing Gate004 or release acceptance.

## Development workspace recovery checks

[DEVELOPMENT_WORKSPACE.md](DEVELOPMENT_WORKSPACE.md) and [bootstrap-workspace.sh](../scripts/bootstrap-workspace.sh) recreate isolated Linux amd64 tools/clones/caches with exact publisher hashes/pins. Required Linux CI runs shell syntax, three [boundary cases](../scripts/tests/bootstrap-workspace.py) and the valid-shell removed-checksum control/restored guard. Those fixtures use explicit local shims, not actual remote downloads; the real fresh workspace provisioning has separate recorded local evidence. Preserve existing destinations and reject corrupt downloads before extraction/cloning. The manifest's qaExecuted:false must not be relabelled as a QA pass. Actual local build/race/profile checks, mandatory native Windows acceptance and deferred Azure/fresh-OS/release boundaries stay distinct.

## AI discovery regression checks

The same isolated AI mutation script now checks three additional compiling faults: foreign batch accounts reaching the filter, stopping on false-with-token, and discarding valid discovery prefix on later failure. Eight total controls require their named assertion failure and restored baseline, on both required hosts; compile/import errors cannot qualify. Linux additionally runs5000x FuzzDiscoveryEnvelope. Literal source query/row expectations and official wire fields, actual AssessmentFilter tag/structural decisions, production in-memory ARM audience/confinement/body closure, continuation/batch/health/ownership/cancellation and exact/excess scope/page/request/raw-row/body/text limits are synthetic library evidence. Public/report/live integration is separately required.

Discovery pre-acceptance review additionally rejects Unicode case-fold aliases in regional DNS input and fixed ARM service path/type names, preserving safe display labels and ASCII casing. Distinct-ID Kelvin-sign/long-s fixtures cannot be masked by duplicate detection; an eighth compiling control disables the DNS ASCII guard and must fail its named assertion (FN058). Earlier candidate full local/native results do not certify this correction; final revised exact-head QA remains required.


## Public AI execution feedback

Both required native jobs run [ai-execution-mutation.py](../scripts/tests/ai-execution-mutation.py). Four isolated faults must compile and fail named pre-auth factory, retained discovery failure, recorded-tag snapshot and projected subscription correlation assertions; restored baselines must pass. Actual Cobra/factories/coordinator/report tests use confined authenticated in-memory Graph/Monitor/ARM responses and independent source cells, not Azure or default transport. Twelve built-CLI tests additionally check AI registry/help and unsupported cloud preflight with no observed tripwire traffic/report changes, while the complete public triple reaches the deliberately invalid credential selection. Independent normal findings/SARIF identity, formula text, default masking, failed/pending source headers and output errors remain covered. Synthetic/native evidence does not close Azure/laptop/Gate004/release validation.
