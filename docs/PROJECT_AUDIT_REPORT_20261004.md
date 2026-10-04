# Comprehensive project audit: progress and findings

Started2026-10-04. Status IN PROGRESS. Feature expansion paused. This is an incremental audit record, not project acceptance.

## Frozen baseline and recovered workspace

Accepted bootstrap/core-v1:07011b63440e69f4a5acb128e0449943f759ee9c/tree72bebbf4885a15a74a5f43dfe57c68fe979c1e7c, PR112. Live protections23890737 require quality/windows-validation, no bypass. Open proposals initially113/f1142c7d and114/771919b9; no accepted ref advanced.

Current Linux review workspace was recovered from authenticated GitHub file APIs. All630 baseline blobs independently match their expected Git SHA1, all reconstructed directory trees match the full untruncated885-entry recursive inventory, and the original signed commit object reconstructs to exact07011b63. Raw third-party notice bytes were transferred through base64 after text patch transfer normalized embedded CRLF; final notice blob40904e84 and all630 hashes match. Local snapshot is for review, not new Go execution proof.

APRL submodule contents are NOT materialized locally, so Git status explicitly reports missing internal/rules/upstream/aprl; do not claim a fully populated clean clone. Remote APRL60eaddda/tree3ea4d285 and source8e4f0577/tree17d93b20 were verified. Target gitlink60eaddda, AOR treea3ff1caf, CUSTOM tree674b9b3d and SKU bloba2d97a60 match governing pins. Native Actions check out and verify APRL separately.

Git/Python3/ripgrep/Node available; no Go command found. Native Actions remain the Go/build/race/vet/package/security executor. Prior Windows scratch drafts/store absent; unpublished availability implementation is not a recoverable accepted artifact.

## Evidence inspected

Planning head771919b90d82dc565fb229d17608423eb09172ad has production bytes identical to accepted07011b63 and four documentation changes. Native [37182711852](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37182711852), quality111378339400/windows111378339442: all mandatory steps success, only conditional Linux failure upload skipped. Full decoded logs read. Both use preview7ba8438521b75b2f144c3b043a7a32a91ebef6d9/treefa3f3fefd39e21bd2cd498f1f7596b7434056a22, ordered07011b63/771919b9 parents, identical candidate tree.

Existing complete test suite, Linux race/vet/format/module/provenance/coverage82.7%, bounded existing fuzz targets,12 actual CLI cases,19 default and19 branded package cases, prior compiling mutation controls/restored baselines, PowerShell and documentation checks passed. This is fresh native evidence on the documentation candidate/preview, not a new accepted-push result and not complete code-review closure. The accepted config fuzz target still uses10s; its historical timeout remains classified separately.

Source [37182711928](https://github.com/DeBoX85/Cloud-Assess/actions/runs/37182711928)/111378328791 full log inspected.14 auxiliary branches/13 inventory branches/30 literal normalization cases. Four files reconstructed from contiguous3+4+2+4 chunks and compared to retained complete UTF8 content and local SHA256 bytes. Provenance source8e4f0577/tree17d93b20/APRL60eaddda exact. This baseline runner does not include unmerged113's availability captures.

Both native vulnerability scans report zero reachable/imported-package findings and one module-only advisory; do not call the dependency graph advisory-free.

## Findings register

| ID | Category / impact | Evidence | Status |
|---|---|---|---|
| AUD001 / FN074 | Process, blocking PR113 native acceptance | Raw JavaScript replacement-string dollar-apostrophe expansion reproduces entire malformed27592-character test.yml from preceding e4f2f766. YAML rejects line341; truncated command/duplicated tail prevents job creation. Accepted baseline unaffected | Repaired on unmerged113 atbe867a174dbc112caba273cda620822fd9828b2e; complete remote bytes/identity/parent verified. Intended13430-character YAML parses and eight explicit Bash blocks pass bash -n. Native37183435404 quality111380420837/windows111380420966 and source37183435389/job111380420823 all mandatory steps passed; full logs inspected, preview8a040608/tree0878fc48 with ordered07011b63/be867a17 parents and identical candidate tree verified. All six captures reconstructed/hashed exactly. Merge remains paused |
| AUD002 | Documentation/recovery, stale accepted resume point | Accepted handover/spec/plan still say aggregation preparation/PR111 despite accepted112; newer continuity was only unmerged113 | Active audit checkpoint supersedes these resume instructions. Final authoritative reconciliation pending, history retained |
| AUD003 | Defensive library transport correction; default SDK route already mitigates | Test-only730a29e6 native37183817468 failed all eight named terminal-read assertions on Linux111381519563/Windows111381519469 against unchanged production; compilation/formatting passed | CONFIRMED injected-poster boundary. Pinned SDK bodyDownloadPolicy already protects default HTTPClient; no live false-complete claim. PR115 corrected head23d2ce4daee7474904975156b4c97fb75c2a6843/tree7438dfd7, parent730a29e6, identity and all3 bytes verified. Fresh full QA pending |
| AUD004 | Recovery limitation | Prior availability draft existed only in transient Windows workspace/store; verified remote holds source contract but no production implementation | Recorded; reconstruct later from verified source contract, feature implementation paused |

No security incident, product-wide instability, independent-person approval, live equivalence or release readiness is inferred.

## Requirement and review matrix

The review inventory below covers all216 first-party Go files. Status NOT REVIEWED means code review remains open even when existing tests passed. Tests and code inventory alone do not establish adequate coverage. Scripts31, workflows4, docs58, bundled rule/source/fixture data and root legal/configuration/build files are separately in scope.

Initial requirements: branding/profile/package; credential/cloud/read orientation; discovery/intended scope/filtering; inventory/catalog/scanners; Graph/Diagnostics/Advisor/Defender/Policy/Arc/Cost; stage completeness/lifecycle/exits; canonical model/summaries; every report/privacy/replacement contract; YAML/internal plugins; operator tooling; maintenance/provenance/dependencies; live/release scope/deferrals. Reconcile each to specification, code, independent tests and exact evidence during phaseB.

| Area | Production Go files | Test files | Review status |
|---|---:|---:|---|
| cmd/cloud-assess | 4 | 11 | NOT REVIEWED |
| internal/advisor | 2 | 3 | NOT REVIEWED |
| internal/app | 1 | 9 | NOT REVIEWED |
| internal/arcsql | 1 | 1 | NOT REVIEWED |
| internal/arg | 6 | 9 | IN PROGRESS |
| internal/assessment | 8 | 0 | NOT REVIEWED |
| internal/azure | 8 | 7 | IN PROGRESS |
| internal/branding | 2 | 1 | NOT REVIEWED |
| internal/config | 2 | 3 | NOT REVIEWED |
| internal/cost | 2 | 3 | NOT REVIEWED |
| internal/defender | 1 | 1 | NOT REVIEWED |
| internal/diagnostics | 2 | 1 | NOT REVIEWED |
| internal/discovery | 5 | 6 | IN PROGRESS |
| internal/equivalence | 3 | 2 | NOT REVIEWED |
| internal/findings | 2 | 2 | NOT REVIEWED |
| internal/gate | 1 | 1 | NOT REVIEWED |
| internal/orchestration | 3 | 7 | IN PROGRESS |
| internal/plugins/aigov | 5 | 5 | NOT REVIEWED |
| internal/plugins | 4 | 5 | NOT REVIEWED |
| internal/plugins/carbon | 2 | 2 | NOT REVIEWED |
| internal/plugins/region | 6 | 7 | NOT REVIEWED |
| internal/plugins/servicehealth | 2 | 1 | NOT REVIEWED |
| internal/plugins/sqleol | 2 | 1 | NOT REVIEWED |
| internal/plugins/zone | 1 | 1 | NOT REVIEWED |
| internal/policy | 1 | 1 | NOT REVIEWED |
| internal/redact | 1 | 1 | NOT REVIEWED |
| internal/renderers/csv | 1 | 1 | NOT REVIEWED |
| internal/renderers/excel | 1 | 1 | NOT REVIEWED |
| internal/renderers/json | 1 | 1 | NOT REVIEWED |
| internal/renderers/sarif | 1 | 1 | NOT REVIEWED |
| internal/renderers/tables | 2 | 1 | NOT REVIEWED |
| internal/reportfile | 3 | 2 | NOT REVIEWED |
| internal/result | 3 | 3 | NOT REVIEWED |
| internal/rules | 5 | 4 | NOT REVIEWED |
| internal/scanners | 1 | 1 | NOT REVIEWED |
| internal/skus | 1 | 1 | NOT REVIEWED |
| internal/stages | 3 | 3 | IN PROGRESS |
| internal/throttling | 1 | 1 | NOT REVIEWED |
| tools/brand-build | 1 | 0 | NOT REVIEWED |
| tools/diagnostics-probe | 1 | 1 | NOT REVIEWED |
| tools/equivalence | 1 | 1 | NOT REVIEWED |

Review coverage is an inventory, not a completed review. Initial inspection traced shared HTTP lifetime/destination boundaries, scope SDK adapters/recursion/cycle guards, intended-scope recording, inventory decoding, ARG paging/metadata and stage interruption/completeness. Critical tests were sampled against actual call paths; full package/test review remains open.

## Exact next action and publication boundary

Finish native repair evidence for PR113 without merging it. Reproduce AUD003 with a targeted terminal-read-error regression on unchanged production code, inspect neighboring Diagnostics streaming behavior, then decide the smallest correction. Complete scope/inventory semantic validation review and build the requirement-to-evidence matrix. Continue remaining production packages/scripts/docs review. No local Go run is claimed. Publish material findings and test/fix candidates separately through existing gates; update this record and SESSION_HANDOVER at coherent boundaries.

## Retrieval review checkpoint

PR113 full repaired native/source logs inspected. Six capture files exactly reproduced from contiguous3+4+2+4+5+4 chunks; each complete UTF8 content and local SHA256 matches retained bytes. New availability input32345bea44a162a1518f9022c93dba3cf5640980c94fb1d3734e172b780cc245/outputd3e8186fb80ffc9a47191228e5cad222636ffdfea1de3c837747bf896acbaef9. All previous controls plus four compiling availability oracle mutations/restored baselines pass both hosts. Count-budget scope fuzz completed500000 executions; property/corpus/workers unchanged. Source14 auxiliary/13 inventory/30 normalization/11 availability/eight probes, exact source/tree/APRL provenance. Zero reachable/imported vulnerabilities, module-only advisory open. This verifies repaired proposal evidence, not feature acceptance or audit-wide closure.

Manual production review has covered CLI/app entry points, all Azure/discovery/stages/rules production files, ARG transport/client/executor, Advisor/Policy/Arc/Cost adapters and throttling, with selected independent fixtures. Orchestration and Defender remainder, comprehensive tests, reports/privacy, plugins, tooling, workflow security, requirements/docs/dependencies remain open. Review matrix status remains IN PROGRESS until associated tests/requirements are reconciled. Cost Management continuation/column handling is a newly identified contract question, not yet a confirmed finding; inspect official pinned API/source and independent fixtures before classification.

AUD003 deliberately retains SDK buffering and all source pins. Default shared pipeline mitigation is source-grounded in azcore1.23.1 bodyDownloadPolicy and internal1.12.0 Payload. New actual pipeline fixture must verify it natively. Earlier Diagnostics inspection does not imply the same production reachability. The original eight failures are retained at immutable730a29e6. No full-project review completed, accepted branch advanced, independent-person review or live approval.
