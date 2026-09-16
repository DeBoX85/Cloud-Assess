# Cloud Assess

Azure cloud assessment toolkit.

Cloud Assess is a new Azure assessment engine being built from a clean repository while using the behavior of a pinned Azure assessment reference implementation as a compatibility baseline.

## Current status

Early bootstrap and characterization phase. The `scan` command is intentionally a stub until the core behavior is covered by characterization tests.

## Working product identity

- Product: **Cloud Assess**
- CLI: **`cloud-assess`**
- Primary human report: **Excel**
- Canonical machine-readable report: **JSON**

Excel and JSON are both first-class outputs. Excel is intended to remain the default format for human review, while JSON is used for automation, integration, and equivalence testing.

## Reference baseline

Reference repository: `DeBoX85/azqr`

Pinned reference commit:

```text
8e4f0577f3615e6c9014c031bcad079f235369cc
```

Pinning the reference commit ensures that upstream changes do not silently change the behavior Cloud Assess is being compared against during the initial reproduction effort.

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

## Documentation

- [Target specification](docs/TARGET_SPECIFICATION.md)
- [Implementation plan](docs/IMPLEMENTATION_PLAN.md)
- [Notices and attribution](NOTICE.md)

## Development branch

Initial bootstrap work is being developed on:

```text
bootstrap/core-v1
```

## License

Cloud Assess is licensed under the Apache License 2.0 at the repository level. Incorporated or derived third-party material remains subject to its applicable license and attribution requirements. See [NOTICE.md](NOTICE.md).
