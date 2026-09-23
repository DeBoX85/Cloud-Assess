# Cloud Assess

Azure cloud assessment toolkit.

Cloud Assess is a new Azure assessment engine built from a clean repository while using the behavior of a pinned Azure assessment reference implementation as a compatibility baseline.

## Current status

The core generic `cloud-assess scan` path is implemented and covered by deterministic, cross-package, race-enabled, and executable smoke tests.

The current core path includes Azure authentication, subscription and management-group discovery, resource inventory, filtering, scanner pruning, pinned recommendation execution, Diagnostics, Advisor, Defender, Azure Policy, Arc SQL, Cost, stage health/completeness, severity gating, and XLSX/JSON/CSV/SARIF/stdout rendering.

The project is **not yet release-complete**. The deterministic semantic source-versus-target harness has produced equivalent live baselines for the default stages, optional Policy/Defender Recommendations/Cost execution, resource-group scope, and two-subscription scope. Cost and Defender plan status now have non-empty live evidence; Policy and Defender Recommendations have empty-result evidence only. The two-subscription pass completed with two unresolved Diagnostics HTTP 400 subrequest warnings. Management-group traversal, Arc SQL, and other scenarios absent from the test environment remain to be validated.

External/plugin execution, scanner-specific CLI commands, `rules` / `plugins` CLI surfaces, packaging, generated dependency/license inventory, and final release/security review also remain outstanding.

## Working product identity

- Product: **Cloud Assess**
- CLI: **`cloud-assess`**
- Primary human report: **Excel**
- Canonical machine-readable report: **JSON**

Excel and JSON are both first-class outputs. Excel is the default format for human review, while JSON is used for automation, integration, and equivalence testing.

Subscription-ID redaction is enabled by default for XLSX, CSV, JSON, and JSON stdout. SARIF intentionally retains stable Azure resource identities for automation/baselining and should be treated as identity-bearing output.

Store reports in an access-controlled directory. On Windows, Go's Unix-style `0600` file mode does not establish a private Windows ACL; report ACL behavior remains a first-release security review item.

## Reference baseline

Reference repository: `DeBoX85/azqr`

Pinned reference commit:

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

Pinning the reference commit ensures that upstream changes do not silently change the behavior Cloud Assess is compared against during the initial reproduction effort.

## Design goals

- read-oriented Azure assessment
- behavioral equivalence for core assessment functionality
- clean product branding and namespace
- explicit assessment-stage health and completeness
- Azure SDK-native authentication
- scalable Azure Resource Graph assessment
- Excel, JSON, CSV, and SARIF output
- extensible rule and plugin architecture
- clear error propagation instead of lower-level process termination
- required open-source attribution retained

## Development and validation

Active development branch:

```text
bootstrap/core-v1
```

The repository quality gate verifies pinned source-data provenance, Go formatting, module consistency, branding boundaries, PowerShell validation helpers, executable build/help smoke tests, race-enabled tests, a minimum statement-coverage floor, `go vet`, and reachable-vulnerability scanning. External GitHub Actions are pinned to immutable commit SHAs.

The generic scan path is suitable for controlled test-environment validation. It should not yet be treated as production/customer-equivalent until the remaining live-coverage, packaging, dependency/license, security, and operational boundaries are complete.

## Documentation

- [Target specification](docs/TARGET_SPECIFICATION.md)
- [Implementation plan](docs/IMPLEMENTATION_PLAN.md)
- [Development ledger](docs/DEVELOPMENT_LEDGER.md)
- [Execution roadmap](docs/ROADMAP.md)
- [Characterization baseline](docs/CHARACTERIZATION.md)
- [Source-versus-target equivalence runbook](docs/EQUIVALENCE.md)
- [Quality Gate 001](docs/QUALITY_GATE_001.md)
- [Quality Gate 002](docs/QUALITY_GATE_002.md)
- [Quality Gate 003](docs/QUALITY_GATE_003.md)
- [Quality Gate 004 plan (not yet passed)](docs/QUALITY_GATE_004_PLAN.md)
- [Notices and attribution](NOTICE.md)

## License

Cloud Assess is licensed under the Apache License 2.0 at the repository level. Incorporated or derived third-party material remains subject to its applicable license and attribution requirements. See [NOTICE.md](NOTICE.md) and [THIRD_PARTY_LICENSES.md](THIRD_PARTY_LICENSES.md).
