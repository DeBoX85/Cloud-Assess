# Maintenance generation and proposal validation

The dependency workflow runs `go mod tidy`; the rule workflow imports the fixed AZQR reference archive and APRL commit. Both then use `scripts/publish-maintenance.sh` to propose changes on a separate branch. Neither pushes `bootstrap/core-v1` as its update target.

## Required checks

- Unchanged generated output creates no proposal branch or commit.
- Changed output is committed only within the selected maintenance paths, using the GitHub Actions bot author and committer identity.
- Rule import validates both source directories before replacing local rule directories. Its staged custom/orphan trees and APRL gitlink must match the recorded pins before publication.
- A failed import, wrong snapshot or rejected push remains a failure. A push rejection must not emit a success/PR summary.
- A proposal is reviewed through a PR, with `quality` and `windows-validation` observed on its exact head before merge. A bot-token branch push alone is not evidence those checks ran.

The fixed import must reproduce custom tree `674b9b3dcb443ce6dc445b48e1db47e4a0ca7082`, orphan tree `a3ff1cafbc0a74ea4e4d2cc5aa2812f7c1dab9f5`, and APRL commit `60eaddda76541f6adbc1c5ffa686829807e55e29` from reference commit `8e4f0577f3615e6c9014c031bcad079f235369cc`. An intentional provenance change requires coordinated review of the import workflow, quality-workflow provenance assertions, generation fixture pins and reference/provenance documentation. Do not relax a hash check merely to make an import pass.

## Executable fixtures

Run on Linux with Git, Bash, Python, curl, the required Go toolchain and APRL initialized at the exact pin:

```bash
python3 scripts/tests/maintenance-publication.py
python3 scripts/tests/maintenance-generation.py
```

Generation tests read the actual named workflow run blocks; their recipes are not reimplemented in the test. The fixture downloads the pinned reference archive once, then replays that archive for the workflow's curl request. Git's fixture-only URL rewrite clones APRL from the checked-out pinned submodule while preserving the production URL in `.gitmodules`. All fixture commits/pushes occur in disposable working repositories and local bare remotes. No GitHub or Azure refs/resources are changed by these tests.

Generation checks cover no-change import, repair of a deliberately corrupted rule, invalid archive/layout rejection, wrong snapshot rejection and actual Go tidy removal of an unused locally replaced dependency followed by stable re-execution. The tidy fixture disables network module retrieval; the project itself separately runs `go mod tidy` with a no-diff assertion in required quality CI. `QA_REFERENCE_ARCHIVE` may supply a locally retained archive for an offline replay; `QA_GO` may identify the exact toolchain binary. Hash checks still validate generated rule content. Network retrieval of the archive and required initialized/toolchain prerequisites must otherwise be available.

These checks do not independently validate remote APRL transport, hosted maintenance dispatch/token permissions or the subsequent proposal-PR event path. Complete hosted maintenance execution remains an open Gate 004 evidence item.

## Troubleshooting

If import fails, distinguish archive download/extraction, missing source layout, Git checkout and generated-tree mismatch before changing code. Inspect the failed run's logs and recorded pins. A layout failure occurs before replacing local rule directories. A later Git/hash failure may leave a partial staged checkout inside the disposable job; it must not proceed to publication. Keep a rejected proposal push failed, inspect the destination branch and preserve divergent remote work. Follow [QA_PROCESS.md](QA_PROCESS.md) and [FAILURE_NOTES.md](FAILURE_NOTES.md) for reproductions and confirmed corrections.

Sources: [Go module tidy semantics](https://go.dev/ref/mod#go-mod-tidy) and [GitHub token event behavior](https://docs.github.com/en/actions/concepts/security/github_token).
