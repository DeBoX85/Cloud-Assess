# Development workflow research, 2026-10-17

Scope: improve development feedback and evidence handling without removing a check or advancing the AZQR, APRL, SDK, module, action or public-schema pins. Research preceded implementation. This is a bounded tooling candidate, not acceptance, a performance guarantee, or release approval.

## Primary guidance and project decisions

| Primary source inspected | Relevant guidance | Decision for Cloud Assess |
| --- | --- | --- |
| [DORA continuous integration](https://dora.dev/capabilities/continuous-integration/) | Automate builds and tests, provide prompt feedback, and keep changes integrated and working | Shorten the critical path through existing checks. Preserve protected review and post-acceptance QA; no automatic merge |
| [DORA test automation](https://dora.dev/capabilities/test-automation/) | Maintain reliable automated tests and catch defects with the fastest appropriate test category | Add a repeatable local preflight before full final-revision QA. Keep independent expected results, named mutation failures, and source comparison |
| [DORA small batches](https://dora.dev/capabilities/working-in-small-batches/) | Use small, independently testable increments; large AI-generated changes increase review difficulty | Publish coherent checkpoints, review contracts before production edits, and finish one bounded slice before expanding scope |
| [GitHub parallel jobs](https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/use-jobs) and [matrix failure behavior](https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/run-job-variations) | Independent jobs can overlap; failure and cancellation semantics require deliberate handling | Use explicit isolated jobs with stable names and exhaustive dependencies. Each mutating suite remains sequential inside its own worker |
| [GitHub job conditions](https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/control-jobs-with-conditions) and [needs semantics](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#jobsjob_idneeds) | Skipped required jobs may report success; dependency failure can skip downstream jobs | Keep the exact protected names quality/windows-validation as always-running validators that require all four corresponding workers to report success |
| [Go test caching](https://pkg.go.dev/cmd/go#hdr-Testing_flags) | Build caching and successful test-result caching are different; -count=1 explicitly reruns tests | Retain compilation/dependency caching. Explicitly rerun the native full tests and coverage tests; mutation helpers already use uncached runs |
| [Go fuzzing](https://go.dev/doc/security/fuzz/) | Fuzz targets should be deterministic and isolated; failures become reproducible corpus inputs | Preserve all five targets, iteration/time/parallel budgets, seed inputs, and failure evidence. No affected-file filtering or nightly-only relocation |
| [GitHub secure workflow use](https://docs.github.com/en/actions/reference/security/secure-use) | Protect workflow dependencies, permissions, and handling of untrusted values | Retain every existing full action SHA, read-only permissions, pull_request event and concurrency policy. Use fixed worker IDs and environment data for needs JSON |

The retrieved documentation is guidance, not a claim that every recommendation is necessary for this repository. New workflow DSL features, new dependency-lock infrastructure, dependency upgrades, merge automation, additional AI agents, automatic flaky-test retries, reduced mutation/fuzz budgets and path-based acceptance skips were not adopted in this slice. They add scope or alter the existing assurance boundary.

## Measured baseline and hypothesis

PR147 candidate b5fd314e97604a2ccd71da339c0c73c9440f2155, [run 38947446296](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38947446296): Linux quality 1,413 seconds (23m33s); Windows validation 2,140 seconds (35m40s). These are job start-to-completion times, not a service-level objective or complete developer cycle time.

The Linux job has 35 movable control suites, with 1,037 seconds of recorded step time; Windows has 34, with 1,902 seconds. Three explicit per-platform groups are balanced using those observed durations, retaining original order within each group. Expected setup, cache, queue and runner variability are not included in those weights. Actual final candidate job timings and overall elapsed time must be measured before claiming improvement. More isolated jobs can increase total runner-minutes even when wall time falls.

The original workflow is retained as canonical text in scripts/ci/testdata/quality-before-20261017.yml.txt, SHA256 1c9e08848af4aabc9b29402a0d0f264efbd679c019466c628d14048e8032e5ee. It is an independent before-change inventory. The generated workflow preserves every original platform step and its command, shell, environment and action options, except two explicit assurance improvements: add -count=1 to the full Windows test and Linux coverage test. Original race, coverage floor, fuzz, package, license and security checks remain mandatory.

## Evidence and rollback

The draft tooling PR's final body is the exact-head evidence index. At publication, native and source gates are pending. Require every worker and both original-name validators to succeed, reconstruct all 23 retained source captures, compare all 385 Linux and 380 Windows control identities with PR147, inspect complete logs and actual native package cases, and record timing and resource tradeoffs. A green aggregate alone is insufficient evidence.

Rollback is a task-owned reviewed revert of tooling/workflow changes, restoring the original workflow and its Windows block locator together. Keep historical records and unchanged product regressions. Never reset unrelated work, weaken the protected ruleset, or alter the reference to make a comparison pass.
