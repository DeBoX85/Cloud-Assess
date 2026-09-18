# Quality Gate 002: Integration and Runnable-Path Audit

Date: 2026-09-18

Reference implementation: `DeBoX85/azqr`

Pinned reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Validated Cloud Assess code/CI baseline: `ec24ddb0eaab7d4f8902a4eb962e9f392499ee9f`

Validation workflow: Go quality gate run `35295328518`

Result: **PASS for the implemented generic core scan and integration layer**

## Purpose

This gate was performed after an interrupted development session and specifically re-audited the most recent integration work rather than assuming earlier package-level green checks were sufficient.

The reviewed path was:

```text
CLI
  -> configuration and stage validation
  -> Azure credential/runtime composition
  -> assessment coordinator
  -> stage health/completeness
  -> canonical assessment result
  -> XLSX / JSON / CSV / SARIF / stdout rendering
  -> severity and process-exit semantics
```

The audit compared material behavior back to the pinned reference where compatibility matters and separately classified deliberate target changes.

## Verified behaviors

### Scope and Graph orchestration

The audit re-verified that:

- explicit CLI subscription IDs are added to the include filter before discovery, matching the pinned reference initialization behavior
- resource-group CLI scope is represented as full ARM resource-group IDs and requires one explicit subscription
- normal scans begin with all registered scanner keys before filter-driven scanner selection
- the two-phase Graph process is preserved:
  1. build the applicable recommendation catalog for selected scanners
  2. prune service scanners to deployed resource types
  3. always include the generic resource scanner
  4. execute only the reduced rule set
- Diagnostics recommendation definitions are added to the same primary catalog as their findings so summary generation does not discard them

No behavioral defect was found in these recent coordinator decisions.

### Failure, cancellation, and exit semantics

The audit re-verified that:

- critical scope/inventory/Graph failures stop the assessment and produce `failed` completeness
- noncritical stage failures preserve later valid stages and produce `partial`
- warnings alone produce `complete_with_warnings`
- reports are written before severity-gate exit code 2 or partial-assessment exit code 3
- a critical stage may still return a partial canonical result that can be persisted before exit code 1
- `SIGINT` / `SIGTERM` cancellation flows from the CLI context through the coordinator to Azure operations
- process exit semantics are:
  - `0`: success
  - `1`: execution/configuration/authentication/render failure
  - `2`: severity gate failed
  - `3`: partial assessment

No defect was found in the recent cancellation or exit-code wiring.

## Findings and remediation

### 1. JSON/stdout subscription redaction gap

**Severity:** blocking privacy/compatibility issue

The CLI defaulted `--redact-subscription-ids=true`, but the application originally passed redaction only to XLSX/CSV table projections. Canonical JSON and JSON stdout still emitted raw subscription IDs.

This also diverged from the pinned reference, where JSON/stdout are constructed from the same masked report tables as the human-oriented reports.

**Remediation:**

- added explicit JSON renderer options
- wired CLI/application redaction through JSON file and stdout output
- redact subscription IDs in dedicated fields and when embedded in ARM/resource strings
- added detection for subscription IDs appearing only in stage warning/error text, such as a skipped Cost subscription with no cost records
- added direct renderer and application regression tests
- strengthened the cross-package smoke test to verify redaction in both JSON and the XLSX ImpactedResources sheet

SARIF intentionally remains identity-bearing because its resource identities/fingerprints are part of the automation/baselining contract and the pinned reference also retains raw identities there.

### 2. Deferred plugin stage could be enabled misleadingly

**Severity:** product-integrity / UX issue

Stage configuration contained the plugin stage even though production plugin execution is intentionally deferred. A user could therefore request the stage and only discover at execution time that there was no runnable implementation.

**Remediation:**

- the current core-v1 CLI rejects explicit plugin-stage execution before Azure authentication
- the error explicitly states that the plugin stage is unavailable in the current core-v1 build
- the target specification and implementation plan now distinguish plugin architecture from production plugin execution/parity

### 3. Default report timestamp selected too late

**Severity:** minor compatibility drift

The application originally selected the default report filename after assessment execution. The pinned reference selects its default name during initialization before the long-running scan.

**Remediation:**

- default output name is now selected before assessment execution
- a regression test advances the clock during a fake assessment and verifies the scan-start timestamp is retained

### 4. Cross-package integration coverage was insufficient

**Severity:** quality-coverage gap

Coordinator, application runner, renderers, and CLI had good package-level characterization but were not sufficiently proven as one path.

An end-to-end test had been started immediately before the interrupted session. During this audit a second overlapping smoke test was briefly introduced; CI correctly exposed the duplicate test declaration. The two tests were reconciled into one authoritative test.

**Remediation:**

The retained cross-package test now executes:

```text
fake Azure Operations
  -> real Coordinator
  -> real Application Runner
  -> canonical Result
  -> real JSON Renderer
  -> real Excel Renderer
```

It verifies:

- phase-two recommendation execution
- complete assessment status
- JSON result structure
- JSON subscription-ID redaction
- Assessment Status as the first Excel worksheet
- ImpactedResources worksheet presence
- subscription-ID redaction in the Excel finding row

### 5. CI did not explicitly prove the final executable

**Severity:** quality-coverage gap

Package tests compiled the command package but did not explicitly demonstrate that the actual executable could be built and invoked.

**Remediation:**

The permanent quality gate now:

- builds `./cmd/cloud-assess`
- runs `cloud-assess --help`
- runs `cloud-assess scan --help`
- runs `cloud-assess --version`

These checks execute before the race suite.

### 6. GitHub Actions checkout runtime was deprecated

**Severity:** non-blocking CI hygiene

`actions/checkout@v4` emitted a Node 20 deprecation warning and was being forced onto Node 24 by GitHub-hosted runners.

**Remediation:**

- moved the workflow to `actions/checkout@v7`, which uses the supported Node 24 runtime

### 7. Repository status documentation was stale

**Severity:** documentation/integrity issue

The README and implementation documents still described `scan` as a stub and placed the project around the pre-orchestration milestone.

**Remediation:**

- README now describes the runnable generic core path and its limits
- implementation plan advances the completion boundary through CLI/exit semantics
- characterization documents orchestration, canonical-result/report, redaction, and executable contracts
- target specification now states the exact redaction/SARIF behavior and deferred plugin execution

## Permanent quality gate

The validated baseline passed all of the following:

1. recursive checkout with pinned APRL submodule
2. exact APRL/AOR/custom/SKU source-data provenance checks
3. `gofmt` enforcement
4. clean `go mod tidy`
5. executable branding-boundary checks
6. actual `cloud-assess` binary build
7. root help smoke test
8. scan help smoke test
9. version smoke test
10. `go test -race -count=1 ./...`
11. `go vet ./...`

Workflow run `35295328518` completed successfully.

## Known limits not certified by this gate

This gate does **not** claim full source-product equivalence or production readiness.

Remaining work includes:

- semantic source-versus-target whole-report comparison
- live Azure regression against the same stable test environment
- live resolution of the Arc SQL `vcores` response shape
- production external/YAML plugin execution
- internal plugin parity/migration
- scanner-specific CLI commands
- `rules` command
- `plugins list/info` command surface
- release packaging/distribution
- final distributable-product security and dependency/license review

## Gate decision

The recent orchestration, application, rendering, and generic CLI work is qualitatively sound after the remediation above.

The project may proceed to the equivalence/live-Azure validation phase using this baseline. Any material source-versus-target delta found during live validation must still be classified as:

- source-compatible target defect
- intentional target improvement
- source defect intentionally corrected
- environmental/data timing difference
- unsupported/deferred feature
