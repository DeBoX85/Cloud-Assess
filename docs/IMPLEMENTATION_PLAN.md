# Cloud Assess Implementation Plan

Status: working implementation plan

Reference implementation: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Current implementation milestone: deterministic foundation, Diagnostics, Advisor, and Defender (steps 1-15) implemented and quality-gated. Azure Policy is the next subsystem.

## Principle

Cloud Assess will be rebuilt behavior-first, not by blindly copying the source tree.

For each major subsystem:

```text
Observe source behavior
  -> Create characterization tests
  -> Define target interface
  -> Implement target subsystem
  -> Run equivalence tests
  -> Record intentional differences
```

## Core v1 development order

1. Repository bootstrap
2. Branding abstraction
3. Neutral domain model
4. Characterization-test harness
5. Configuration and filters
6. Azure authentication/environment
7. Scope discovery
8. Resource discovery
9. Scanner registry
10. Rule catalog
11. Azure Resource Graph client
12. Core recommendation execution
13. Diagnostics
14. Advisor
15. Defender
16. Policy
17. Arc SQL
18. Cost
19. Stage health and assessment completeness
20. Canonical findings summary
21. JSON renderer
22. CSV renderer
23. Excel renderer
24. SARIF renderer
25. CLI integration
26. Severity and exit-code gates
27. Equivalence harness
28. Live Azure regression suite
29. CI and packaging
30. Security/license review
31. Plugin migration

## Current completion boundary

Steps 1-15 are implemented at the subsystem/characterization level and pass the repository quality gate. Their deterministic contracts are available for later orchestration, but the end-to-end `scan` command is not yet complete.

Diagnostics returns canonical recommendation definitions, findings, and warnings. Advisor remains a separate auxiliary dataset combining ARG recommendation instances with Advisor metadata. Defender remains two separate auxiliary datasets: plan/tier status and unhealthy security recommendations. These subsystems are intentionally not wired into the placeholder CLI until stage orchestration and canonical result assembly are implemented.

The next implementation target is Azure Policy (step 16).

## Characterization levels

### 1. Deterministic unit behavior

No Azure connection required.

Priority cases:

- filter precedence
- tag matching
- resource-group validation
- stage defaults and validation
- severity thresholds
- previous-calendar-month cost period
- subscription-ID redaction
- finding deduplication
- summary calculation
- recommendation applicability
- scanner registry mappings
- recommendation normalization
- report field mappings

### 2. Fixture-based components

Sanitized Azure/API fixtures should test mappings such as:

- ARG row -> Finding
- Advisor response -> Advisor record
- Defender pricing row -> Defender plan status
- Defender assessment -> Defender recommendation
- Policy state -> Policy record
- resource row -> Resource
- diagnostic-settings batch response -> diagnostic finding

### 3. Golden reports

Generate deterministic canonical assessment data and compare semantic output for:

- JSON
- CSV
- XLSX
- SARIF

Excel comparison should focus on worksheet names, headers, rows, ordering, counts, redaction and relevant formatting contracts rather than raw XLSX bytes.

### 4. Live Azure equivalence

For representative Terraform scenarios:

```text
Deploy known fixture
Run pinned reference
Run Cloud Assess
Normalize both outputs
Compare findings and datasets
Destroy fixture
```

Primary finding comparison key:

```text
Recommendation ID
Resource ID
Category
Impact
Source
```

## Existing reference fixtures to reuse

The pinned source already includes useful integration scenarios for:

- authentication
- Redis Enterprise
- VNet DNS
- generic resource findings
- Storage diagnostic settings
- Storage HTTPS
- Storage TLS
- Storage immutable versioning

These form the initial live-equivalence suite.

## Additional characterization scenarios

Add tests for:

- subscription include/exclude
- resource-group include/exclude
- resource-type include
- resource exclusion
- recommendation exclusion
- include tags
- exclude tags
- include/exclude tag precedence
- multiple subscriptions
- management-group recursion
- recommendation not applicable
- recommendation compliant
- recommendation noncompliant
- duplicate findings
- diagnostic setting present/absent
- Advisor data present
- Defender pricing status
- Defender unhealthy recommendation
- Policy noncompliance
- Arc SQL presence
- Cost available
- Cost unauthorized/unavailable

## Intentional differences from source

The following differences are deliberate and should not fail equivalence tests:

- new product and CLI identity
- no legacy CLI compatibility requirement
- no legacy configuration compatibility requirement
- `GraphResult` replaced by `Finding`
- `GraphRecommendation` replaced by `RecommendationDefinition`
- lower-level packages return errors instead of calling `log.Fatal` or silently returning nil datasets
- stage status and assessment completeness are explicit
- subscription masking is named subscription-ID redaction
- custom rule source/path is neutral rather than legacy-branded
- plugin permissions are documented independently
- Diagnostics findings identify `Azure Resource Manager` as their validation mechanism instead of inheriting a generic Azure Resource Graph label
- non-success Diagnostics subrequests retain source-compatible finding semantics while producing explicit uncertainty warnings
- malformed diagnostic-setting IDs become warnings rather than panic-prone parsing
- Advisor metadata is retrieved through the shared authenticated ARM HTTP layer rather than adding the `armadvisor` SDK dependency
- malformed Advisor ARG rows become explicit warnings
- Advisor records are sorted deterministically after normalization
- malformed Defender ARG rows become explicit warnings
- Defender status and recommendation records are sorted deterministically after normalization

## Output strategy

JSON is implemented before Excel because it is easier to normalize and compare in automated equivalence tests.

This does not make JSON the preferred human report.

Excel remains the default end-user report. Both Excel and JSON are first-class outputs generated from the same canonical assessment result.

## Exit semantics

Tool execution failures, partial assessments, and severity-gate failures must be distinguishable.

Conceptually:

```text
0 = complete successful assessment and gate passed
1 = execution/configuration failure
2 = quality/severity gate failed
3 = partial assessment because a requested stage failed
```

Exact numeric values can be finalized with the CLI implementation.

## Licensing and attribution

Cloud Assess is currently licensed under Apache 2.0 at the repository level.

Reused or derived MIT-licensed source and recommendation material must retain applicable copyright and license notices, including Microsoft-originated source, APRL material, Azure Orphan Resources material, and other third-party dependencies as required.

A generated/maintained notice inventory should be part of the release process.

## Definition of done for core v1

Core v1 is complete when:

- repository builds independently
- no legacy product-facing branding remains
- required legal attribution remains
- scope/filter behavior is equivalent
- rule loading and scanner pruning are equivalent
- normalized core findings are materially equivalent for reference fixtures
- Diagnostics, Advisor, Defender, Policy, Arc SQL and Cost are reproduced
- stage failures are explicitly visible
- assessment completeness is explicit
- Excel and JSON are both first-class outputs
- CSV and SARIF work
- redaction and severity gates work
- unit and characterization tests pass
- selected live Azure equivalence tests pass
- security and license checks pass
