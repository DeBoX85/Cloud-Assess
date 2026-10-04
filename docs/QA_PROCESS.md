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


Region primary projection candidate adds scripts/tests/region-primary-mutation.py to both required native jobs. Selected-identity, unknown-SKU-denominator, comparison-work-limit, aggregate-detail-work and joined-separator-bytes faults must compile and fail their named assertions; restored baselines must pass. See [REGION_SELECTION.md](REGION_SELECTION.md) for independent capture hashes, complete42-score/25-cell checks and the explicitly unavailable public feature.


The [region source characterization runner](REGION_SOURCE_CHARACTERIZATION.md) is a separate development candidate. Its exact-head characterization job must produce complete provenance/chunks and unchanged pinned-source/APRL/module guards before promoting source fixtures. Both ordinary required quality/windows-validation jobs and reviewed protected acceptance remain mandatory; source-only execution does not certify auxiliary target behavior or live Azure.


## Historical PR105 correction and recovery (superseded below) (2026-10-03)

Initial source characterization37131392522/job111227049328 and required native37131392461/quality111227049234/windows111227049053 passed on7b1fdede. Review found path filters could skip the promised final-head source job on follow-up commits. Removed both path filters; fresh three-job final-head evidence and protected acceptance remain required (FN066). Recovered executor/source/tools in an isolated exact remote worktree, preserving old dirty corrections and unknown command completion. No auxiliary runtime, public region, source/dependency pin or Azure/live/Gate004/release change. Exact next step is publish correction, verify source chunks/hash/provenance and all final-head jobs, then migrate captured auxiliary tables separately.


## Accepted source runner and auxiliary projection gate

PR105 is VERIFIED OFFLINE at accepted merge 7e4ca5c955abb7f882591a44d8bb2f74c93ae243, tree 5c79e721d8281f04ba066006a4ff5601202abdbd, ordered parents c270ea22 / c9cc112. Final candidate run37143242571 (quality111261843942, Windows111261844116) and characterization37143242555/job111261844024 passed. Separate accepted-push quality run37143573050 (quality111262849038, Windows111262849202) and characterization37143573016/job111262848443 passed. Required steps and retrieved full logs were inspected in this recovery; only conditional failure upload skipped. Native primary/AI controls, compiled CLI/package cases, race/vet/fuzz, provenance and vulnerability scans passed; coverage81.8%, zero reachable/imported findings and the existing module-only advisory remains open. Final and accepted-push auxiliary capture chunks/provenance are complete and identical: inputs f9695cfa0bd662b2dbc52addb68a9208b57be1e60dabe31c99cb5651d2962c28, outputs efe06093eddece6b54617c806bdb74727aed5228a5a96f7300068b22bd77bc06.

Prior PR105 pending prose describes preparation. The next [auxiliary contract](REGION_AUXILIARY.md) requires independent complete cells/hash checks and compiling guard controls on both required hosts. No fresh local Go or automated Code Review tool evidence is claimed in this Windows recovery.

Both required hosts run [region-auxiliary-mutation.py](../scripts/tests/region-auxiliary-mutation.py): selected identity, row work and aggregate text controls must compile, fail named assertions and restore. Formatting failure now prints gofmt's diff before the unchanged failing exit, providing correction evidence when local Go is unavailable.


## Accepted quota/reservation proof

Quota/Capacity Reservations pure projections are VERIFIED OFFLINE through [PR106](https://github.com/DeBoX85/Cloud-Assess/pull/106), merge `c4643527dc948f96a565da237662a104a084f14e`, tree `38248bc19be95c0058377ed373a4bf90b20e4ca2`, ordered parents7e4ca5c/e7e5511. Final head e7e551104cc6e0e6dda25c389af9176eb5ed562c passed [run37146790757](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37146790757), quality111272343099/Windows111272342986, all mandatory steps/full logs inspected. Both tested preview06848f39 with the identical tree/parents. Source characterization [37146790810](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37146790810)/job111272298401 reproduced all14 unchanged helper branches with exact captured bytes/provenance. Both hosts passed every quota/reservation cell/empty branch, all three compiling auxiliary controls and restored assertions, existing five primary/eight AI request/four AI execution controls,12 actual CLI and19 default/19 branded packages; docs327/9/1. Linux race/vet/fuzz and82.1% coverage passed; zero reachable/imported vulnerability findings, existing module-only advisory open. Remote merge/ref/tree/parents/human author/GitHub committer verified. Fresh local Go/fetched-source evidence is unavailable in this Windows session.

Earlier PR106 candidate/pending statements are historical after this acceptance. The next [service-availability contract](REGION_SERVICE_AVAILABILITY.md) needs its own source-cell/structural/aggregate/sheet and final native proof; successful quota tables do not certify that work.


Accepted service-availability slice PR108: REGION_SERVICE_AVAILABILITY defines full independent source cells/empty topology, literal registry and boundary/collision/selected aggregate declaration/ownership/cancellation tests. Ten compiling controls and restored named assertions are added to both mandatory native jobs. Every final code/docs head needs source characterization plus successful Linux/Windows/full logs; prior quota/continuity QA is not acceptance. Input aggregate declaration checking does not prove selected underlying resource ownership, and projection-completed health does not certify service-input completeness.


Historical preacceptance PR108 review strengthened selected-empty/null and exact mixed BMP/non-BMP UTF16 boundaries; ten compiling service controls/restored named tests now include empty names and byte/rune confusion. FN070/071 preserve the initial failures/review. Earlier native/source runs do not certify this corrected complete head.


PR108 final/accepted-push source/native/full-control proof completed and is indexed in REGION_SERVICE_AVAILABILITY. Current CostComparison candidate adds nine compiling named controls/restored baselines to both native jobs; REGION_COST_COMPARISON defines independent complete-source/received-order/float/scope/entry/text/ownership/cancellation acceptance. New complete code/docs head needs separate source/Linux/Windows QA; prior service success is not cost acceptance. Declared aggregate contributors do not prove actual selected resource ownership or complete pricing health.


PR109 initial native37168867529 failed only at the selected-contributor mutation compile error after passing package tests (FN072); no acceptance or skipped-step success is claimed. The corrected control retains selected through (!selected && false). Production/source/capture/pins are unchanged. Fresh complete-head source/native gates and all nine compiling controls/restored baselines remain required.


Current pure Inventory candidate: require complete captured raw/masked/empty ten columns/cells/order/description, documented Region Inventory composition correction, pinned known/unknown/case/nonpositive/VMSS capacity and mask semantics, selected/embedded UUID/duplicates/source masking prefix/text/work/product overflow, optional/real Unicode, owned16-concurrency/cancellation no partial rows. Eight compiling named controls/restored baselines run on both required native hosts. Decoded-budget regression includes1000x512-byte scope labels and4000 resources whose projected output stays below limit, so final output validation cannot mask a removed early guard. Full source/native exact-head logs/remote scope/bytes/tree/identity/base/preview/rules/protected acceptance remain mandatory. No projected-only-control/local-Go/load/live claim. CostComparison is VERIFIED OFFLINE through [PR109](https://github.com/DeBoX85/Cloud-Assess/pull/109), merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents. 

Final candidate1cf8e4cf6adc38ffe377894824b51f3b26431373/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, parentca16157 on accepted1957f6d7. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165402 quality111338402524/Windows111338402644 passed every mandatory step; only conditional failure upload skipped. Both tested preview77c2b59d40f89e1966402668a98efcdfad2e4d5d with identical final tree/ordered base/head parents. Full logs inspected: all nine Cost controls/restored baselines, ten service/three auxiliary/five primary/eight request/four execution controls, actual CLI/default+branded package/docs/provenance/module/inventory/branding guards; Linux race/vet/fuzz and82.5% coverage; zero reachable/imported vulnerability findings, existing module-only advisory remains open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169165342 /111338403371 completed; full14 branches, indexed3-input/4-output chunks, exact retained bytes/hashes and source/tree/APRL provenance verified. Complete14-file scope/deletions/remote bytes/original scratch blobs, parent/human identity/pins/base/head/preview/current no-bypass rules23890737 verified. Semantic self-review checked actual unchanged source first-tuple/item break behavior, all source topology/cells and safe bounded identity/nonfinite/text/cancellation behavior. No callable final manual Code Review tool or independent-person approval is claimed; automatic comments currently empty. FN072 compile failure was not accepted. Expected-head protected merge and distinct accepted-push evidence follow; no fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure.

Protected acceptance verified: PR109 merge6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f/treeb4087835fcead110dd41ef3f5b1fad0a660d31be, ordered1957f6d7/1cf8e4cf parents, Denis author/GitHub committer/livecore ref and no open proposals. Distinct accepted-push native/source QA remains pending; no dependent live proof claimed.

Distinct accepted-push native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513599 quality111339415060/Windows111339415184 and source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169513582 /111339414976 passed on exact6dfcd1be54f3c8199aa3478eb62b4bb38f6a956f. All mandatory steps/full logs inspected; only conditional Linux failure upload skipped. Both hosts passed9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution compiling controls/restored baselines,12 actual CLI/19 default/19 branded packages/docs/branding/provenance; Linux race/vet/fuzz82.5%, zero reachable/imported findings. Source14 complete branches, indexed chunks/both retained exact bytes/hashes and pinned source/tree/APRL provenance verified. Module-only advisory and all live/public/laptop/Azure/Gate004/release obligations stay open. No fresh local Go/source or independent-person approval claimed.



2026-10-04 current source-only calculation checkpoint: Inventory is VERIFIED OFFLINE through [PR110](https://github.com/DeBoX85/Cloud-Assess/pull/110), merge2dcd8dae924ce3fad5402890cdeca0ecbef6e01e/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, ordered6dfcd1be/76ac3f8d parents. 

Final76ac3f8d65acf9970305c58186e13bb19cf78ffe/tree5f9e76032cb816121b97febacd19ae3e40aefe6f, parent accepted6dfcd1be. Native https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985397 quality111340777403/Windows111340777269 passed every mandatory step; only conditional failure upload skipped, full logs inspected. Both tested preview4e8f21f81155a55383272fe7ade895896fb72db6 with identical tree/orderedbase+head parents. Eight Inventory controls/restored baselines and all9 cost/10 service/3 auxiliary/5 primary/8 request/4 execution controls passed each host;12 actual CLI/19 default/19 branded packages/docs/branding/provenance/module/inventory checks. Linux race/vet/fuzz passed; 82.7% coverage, zero reachable/imported findings; existing module-only advisory stays open. Source https://github.com/DeBoX85/Cloud-Assess/actions/runs/37169985428 /111340777428 passed full14 branches/indexed3-input4-output chunks/both retained exact bytes/hashes/source-tree-APRL provenance. Original scratch Go/test/script blobs/complete16 remote files/scope/no deleted or pinned files/commit identity/tree/parent/base/preview/rules verified. Semantic self-review checked actual source helper/mask/SKU/duplicate-skip plus explicit sheet/scope/budget/product corrections. Automated comments currently empty; no callable final manual Code Review tool or independent-person approval is claimed. No fresh local Go/source/laptop/Azure/public/Gate004/release/advisory closure. 
 Protected acceptance/ref/tree/parents/Denis author/GitHub committer verified. Distinct accepted-push native/source remains pending.
New inventory_capture_test.go.txt/extended capture-region-source.py invoke only unchanged pure inventory/merge/normalization/predicate functions in isolated pinned source/APRL.13 structural branches/30 literal normalization cases; actual output retention/topology/hash/final source/native/remote acceptance is pending. No target calculation code, pin/public/request/live/platform change. See REGION_INVENTORY_CALCULATIONS for scope/contracts/acceptance/rollback/next action.
