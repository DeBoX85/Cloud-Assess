# Quality Gate 003: Post-Live-Validation Repository Audit

Date: 2026-09-23

Reference implementation: `DeBoX85/azqr`

Pinned reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

Audited Cloud Assess baseline: `5cc911ea853cbcdf402a43cb0c96d359f8ece81d`

Validated remediation commit: `954dd75584f96dff1f5017b952c0fd11def40a53`

Validation workflow: Go quality gate run `35847062551`

Result: **PASS for the remediated core-v1 repository and permanent quality gate**

## Purpose

This gate was performed after the first default-stage, optional-stage, and resource-group-scoped live Azure equivalence passes. Its purpose was to check for implementation drift, review every repository Markdown document and the relevant GitHub Actions history, rerun the deterministic quality controls independently, and strengthen controls that were no longer sufficient.

This gate does not replace the recorded live-equivalence passes and does not claim release readiness. It validates the repository and evidence tooling used to continue that work.

## Audit scope

- target specification, implementation plan, characterization baseline, equivalence runbook, quality-gate records, and development ledger
- all other tracked Markdown and repository licensing/notice files
- current and recent GitHub Actions execution logs
- CLI, orchestration, stage health, scope, rendering, and semantic-comparison paths
- live-equivalence runner scope and evidence metadata
- pinned rule/source provenance
- Go formatting, module graph, build, tests, race detector, vet, coverage, and vulnerability status
- GitHub Actions triggers, permissions, action references, and maintenance workflows

## Material finding and remediation

### Outdated Go patch release

The repository and CI selected Go `1.26.0`. An independent `govulncheck` source scan reported 19 reachable standard-library vulnerability findings through Cloud Assess call paths.

Remediation:

- raised the repository and CI baseline to Go `1.26.8`
- reran the full race-enabled suite and executable checks under Go `1.26.8`
- added pinned `govulncheck v1.8.0` execution to the permanent quality gate

Result after remediation:

```text
Your code is affected by 0 vulnerabilities.
```

The scanner still identified vulnerability records in imported/required modules whose affected symbols are not called by Cloud Assess. Those non-reachable records remain subject to normal dependency maintenance and release review.

## Evidence-runner findings and remediation

### Implicit stages serialized as null

The live runner correctly executed the default stage set when `-Stages` was omitted, but `run-metadata.json` serialized `stages` as null. Phase L remained reconstructable from its commands and dataset coverage, but the metadata did not meet the runbook's intended evidence contract.

Remediation:

- metadata schema advanced from `1.0` to `1.1`
- `stages` now contains the resolved effective stage set
- `stageSelection.requested` preserves explicit stage controls
- `stageSelection.usesImplicitDefaults` distinguishes default selection from an explicit override

### Broader scope could not be expressed by the runner

The product and pinned reference both support string-slice scope flags, while the runner accepted only one subscription, management group, or resource group.

Remediation:

- runner accepts multiple subscriptions
- runner accepts multiple management groups
- runner accepts multiple resource groups within exactly one subscription
- invalid mixed management-group/subscription scope and resource groups across subscriptions are rejected before execution

### PowerShell behavior was only parser-checked

The permanent gate parsed the live runner but did not exercise its evidence-selection logic.

Remediation:

- extracted deterministic scope/stage resolution into `live-equivalence.helpers.psm1`
- added executable PowerShell tests for defaults, additions, removals, unknown stages, mandatory Graph, all supported scope forms, and invalid combinations
- expanded CI parsing to recurse across both `.ps1` and `.psm1` files

## Documentation reconciliation

The README, implementation plan, characterization targets, and ledger current boundary still described the first live source-versus-target comparison as future work despite the recorded successful passes.

They now distinguish:

- completed default-stage live equivalence
- completed optional-stage execution/equivalence, including non-empty Cost data
- completed resource-group-scoped equivalence
- remaining multi-subscription/management-group traversal work
- remaining non-empty Policy, Defender Recommendations, and Defender plan-status evidence
- unresolved Arc SQL and release-level work

Historical Quality Gates 001 and 002 remain unchanged audit snapshots.

## Permanent quality-gate strengthening

The gate now additionally:

- runs for pull requests targeting `bootstrap/core-v1`
- uses Go `1.26.8`
- pins external GitHub Actions to full immutable commit SHAs
- executes the PowerShell helper tests
- enforces at least 75% total statement coverage
- runs `govulncheck v1.8.0`

The maintenance workflows use the same pinned checkout and Go setup revisions.

## Validation results

The remediation passed:

1. exact APRL/AOR/custom/SKU provenance checks
2. `gofmt` enforcement
3. clean `go mod tidy`
4. branding-boundary enforcement
5. recursive PowerShell parse checks
6. executable live-runner helper tests
7. CLI build and root/scan help/version smoke tests
8. repeated shuffled Go tests (`-shuffle=on -count=3`)
9. full race-enabled tests
10. total statement coverage of 77.2%, above the 75% floor
11. `go vet ./...`
12. `actionlint`
13. `govulncheck ./...` with zero reachable vulnerabilities

An advisory `staticcheck` run reported only ST1005 capitalization warnings for error strings beginning with Azure service/product names such as Azure, Advisor, Defender, Policy, and Excel. No correctness finding was reported; those user-facing names were not mechanically lowercased in this remediation.

## Gate limitations and residual work

This gate does not certify:

- raw replay of the intentionally untracked Phase J/K/L live evidence bundles
- multi-subscription or management-group live equivalence
- non-empty Policy, Defender Recommendations, or Defender plan-status row projection
- Arc SQL `vcores` response-shape resolution
- production plugin execution or internal-plugin parity
- scanner-specific, `rules`, or `plugins list/info` CLI surfaces
- packaging/distribution
- generated dependency/license inventory
- final release-level security and operational review

Lower-coverage packages, especially throttling, CLI composition, discovery adapters, canonical result assembly, and production Azure adapters, remain priority areas for additional negative-path characterization even though the aggregate floor passes.

## Gate decision

**PASS for the remediated core-v1 repository and permanent quality gate.**

After this gate is merged, live validation may continue with multi-subscription or management-group scope. Release readiness remains explicitly deferred until the residual boundaries above are closed or deliberately re-scoped.
