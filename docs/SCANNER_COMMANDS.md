# Scanner-specific scan commands

Date: 2026-10-02 (Europe/Oslo). AR-04 slice B1 follows [DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md). Accepted starting baseline is PR #77, `c79c2ffc9bfcec66da9f24dab631eb6bef6d1836`. Rollback is a reviewed revert of this slice on the then-current branch, with dependency checks and required native QA; generic scan remains the fallback before acceptance.

## Source contract and acceptance

Pinned AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`, `cmd/azqr/commands/scanners.go`, registers sorted service keys under `scan`, uses no positional arguments and passes exactly one scanner key to the shared scan. Its persistent scan flags apply to child commands. `models` selection treats a single scanner key as authoritative over `include.resourceTypes`; generic scan continues to use that filter. This does not bypass other resource, subscription, tag, recommendation or stage filters.

Before implementation, independently captured `/tmp/azqr-pinned-rules scan --help` from a clean pinned executable with APRL `60eaddda76541f6adbc1c5ffa686829807e55e29`. Extracted exactly 87 `Scan` command keys into `cmd/cloud-assess/testdata/scanner-keys-reference.json`, SHA-256 `398fed601cae3785f660c5572edfdf58d631380cb9a02be14ed262bec567d861`. The JSON fixture is compared semantically, so platform line endings are irrelevant. Regeneration must compile that unchanged source with pinned embedded inputs; do not derive the oracle from the target registry. Attribution remains in repository notices.

Implementation adds only command registration, inherited flags and explicit scanner-key mapping into existing production orchestration. Generic scan defaults, pinned libraries, stages, filters, schemas, dependency graph, normalization and branding remain unchanged. `scan arc` selects the Arc service registry; enabling the separate optional Arc SQL stage still requires its existing stage configuration. No new top-level aliases or automatic optional-stage enablement are introduced.

## Usage

```bash
cloud-assess scan --help
cloud-assess scan st --help
cloud-assess scan vm --subscription-id SUBSCRIPTION_ID --json --xlsx=false
```

Successful assessment commands require the existing Azure identity/permissions. Help and invalid-input checks used for this offline slice perform no Azure assessment. Inherited scope, filter, output, timeout and quality-threshold flags work before or after the scanner key. Do not confuse a scanner key with an Azure resource-type string.

## QA and limits

- All 87 commands match the independent source capture, map their selected key, inherit representative flags and retain application exit codes.
- A fabricated CLI-to-production-coordinator fixture contrasts storage-specific scan against generic scan under a VM include filter, verifying executed definitions, selected inventory and excluded inventory. Disabled operation guards fail if called; no credentials or Azure requests occur.
- Extra arguments fail before execution; failure exit codes propagate; scanner selection does not leak into a subsequent generic command on the same root.
- Actual built/installed default and custom executables verify all command keys plus representative inherited help, rejected arguments, unavailable plugin stage and negative timeout through HTTP/auth tripwires and unchanged filesystem snapshots.
- A compiling negative control removed request scanner-key mapping. The independent coordinator assertion failed on VM definitions in the storage case; restored code passed focused race tests. This is one meaningful mutation, not whole-project mutation coverage.

Exact-head Linux quality and native Windows validation remain mandatory before merge. Existing generic orchestration, filter, report and package checks remain required. Current source command parity and synthetic selection evidence do not certify live execution of all 87 services or close nonempty optional-stage, nested-group, plugin, security/release or other deferred boundaries. No Azure/laptop input is required for this bounded implementation.
