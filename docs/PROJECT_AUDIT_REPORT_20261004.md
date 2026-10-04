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
| AUD003 | Defensive library transport correction; default SDK route already mitigates | Test-only730a29e6 native37183817468 failed all eight named terminal-read assertions on Linux111381519563/Windows111381519469 against unchanged production; compilation/formatting passed | CONFIRMED injected-poster boundary. Pinned SDK bodyDownloadPolicy already protects default HTTPClient; no live false-complete claim. PR115 corrected head23d2ce4daee7474904975156b4c97fb75c2a6843/tree7438dfd7, parent730a29e6, identity and all3 bytes verified. Native37184407342 Linux111383253100/Windows111383253050 and source37184407425/111383253548 passed. Full logs and four capture bytes/provenance inspected; preview87b10d5d/tree7438dfd7 with ordered07011b63/23d2ce4d parents verified. Remains draft/unmerged |
| AUD005 | Product completeness, Cost adapter | Official stable2021-10-01 spec defines nextLink; native37184679263 failed two compiling named paged empty/nonempty assertions on Linux111384036588 and Windows111384036492 | CONFIRMED synthetic first-page false-success. PR116 corrected8d54c0d921825996bbe9c663a8f48f1a6d95ca00/treeed772a94 parent78dd85c9, identity/3 changed bytes verified. Fail closed before accepted costs, no provider URL followed/echoed. Actual coordinator partial/retained healthy stages fixture passed. Native37184879992/Linux111384624285/Windows111384624226 and source37184879983/111384624557 fully passed; full logs/four contiguous capture bytes and hashes/provenance inspected. Preview9caddc01fbe6c95df1cc2766cb16c90835a9d000/treeed772a9455cf11d490b310b6778d11977ac262a4 matches candidate, ordered07011b63/8d54c0d9 parents verified. Remains draft/unmerged. Full pagination still unimplemented, no historical live incidence claim |
| AUD006 | QA evidence correctness, comparator false success | Test-only4cb8bd84/native37185359814 failed all six compiling real-command assertions on Linux111386033586 and Windows111386033687; format passed, healthy empty controls did not fail | CONFIRMED synthetic exit0 for absent/unknown completeness, unsupported schema and no recognized reference datasets. PR117 final correctione381d0ad6e45989cf64526ac5d34131dc349e8e5/tree5e91d391cead6b74507abb6ac780817b9a42afaf/parent1bf98b7a Denis identity/three remote bytes verified. Both supported schemas1.0/1.1 retained, fresh final gates pending; no historical live incidence inferred |
| AUD004 | Recovery limitation | Prior availability draft existed only in transient Windows workspace/store; verified remote holds source contract but no production implementation | Recorded; reconstruct later from verified source contract, feature implementation paused |

No security incident, product-wide instability, independent-person approval, live equivalence or release readiness is inferred.

## Requirement and review matrix

The review inventory below covers all216 first-party Go files. Status NOT REVIEWED means code review remains open even when existing tests passed. Tests and code inventory alone do not establish adequate coverage. Scripts31, workflows4, docs58, bundled rule/source/fixture data and root legal/configuration/build files are separately in scope.

Initial requirements: branding/profile/package; credential/cloud/read orientation; discovery/intended scope/filtering; inventory/catalog/scanners; Graph/Diagnostics/Advisor/Defender/Policy/Arc/Cost; stage completeness/lifecycle/exits; canonical model/summaries; every report/privacy/replacement contract; YAML/internal plugins; operator tooling; maintenance/provenance/dependencies; live/release scope/deferrals. Reconcile each to specification, code, independent tests and exact evidence during phaseB.

| Area | Production Go files | Test files | Review status |
|---|---:|---:|---|
| cmd/cloud-assess | 4 | 11 | IN PROGRESS |
| internal/advisor | 2 | 3 | IN PROGRESS |
| internal/app | 1 | 9 | IN PROGRESS |
| internal/arcsql | 1 | 1 | IN PROGRESS |
| internal/arg | 6 | 9 | IN PROGRESS |
| internal/assessment | 8 | 0 | IN PROGRESS |
| internal/azure | 8 | 7 | IN PROGRESS |
| internal/branding | 2 | 1 | IN PROGRESS |
| internal/config | 2 | 3 | IN PROGRESS |
| internal/cost | 2 | 3 | IN PROGRESS |
| internal/defender | 1 | 1 | IN PROGRESS |
| internal/diagnostics | 2 | 1 | IN PROGRESS |
| internal/discovery | 5 | 6 | IN PROGRESS |
| internal/equivalence | 3 | 2 | IN PROGRESS |
| internal/findings | 2 | 2 | IN PROGRESS |
| internal/gate | 1 | 1 | IN PROGRESS |
| internal/orchestration | 3 | 7 | IN PROGRESS |
| internal/plugins/aigov | 5 | 5 | IN PROGRESS |
| internal/plugins | 4 | 5 | IN PROGRESS |
| internal/plugins/carbon | 2 | 2 | IN PROGRESS |
| internal/plugins/region | 6 | 7 | IN PROGRESS |
| internal/plugins/servicehealth | 2 | 1 | IN PROGRESS |
| internal/plugins/sqleol | 2 | 1 | IN PROGRESS |
| internal/plugins/zone | 1 | 1 | IN PROGRESS |
| internal/policy | 1 | 1 | IN PROGRESS |
| internal/redact | 1 | 1 | IN PROGRESS |
| internal/renderers/csv | 1 | 1 | IN PROGRESS |
| internal/renderers/excel | 1 | 1 | IN PROGRESS |
| internal/renderers/json | 1 | 1 | IN PROGRESS |
| internal/renderers/sarif | 1 | 1 | IN PROGRESS |
| internal/renderers/tables | 2 | 1 | IN PROGRESS |
| internal/reportfile | 3 | 2 | IN PROGRESS |
| internal/result | 3 | 3 | IN PROGRESS |
| internal/rules | 5 | 4 | IN PROGRESS |
| internal/scanners | 1 | 1 | IN PROGRESS |
| internal/skus | 1 | 1 | IN PROGRESS |
| internal/stages | 3 | 3 | IN PROGRESS |
| internal/throttling | 1 | 1 | IN PROGRESS |
| tools/brand-build | 1 | 0 | IN PROGRESS |
| tools/diagnostics-probe | 1 | 1 | IN PROGRESS |
| tools/equivalence | 1 | 1 | IN PROGRESS |

Review coverage is an inventory, not a completed review. Initial inspection traced shared HTTP lifetime/destination boundaries, scope SDK adapters/recursion/cycle guards, intended-scope recording, inventory decoding, ARG paging/metadata and stage interruption/completeness. Critical tests were sampled against actual call paths; full package/test review remains open.

## Exact next action and publication boundary

Inspect PR117 finale381d0ad exact-head full native/source/preview/capture proof. Continue whole-project production/test/requirements review, especially tools, Diagnostics and scripts/workflows. PR113/115/116 have verified full proposal evidence but remain unmerged pending audit disposition; frozen accepted07011b63 unchanged. No local Go, live or release acceptance claim. Publish coherent code/docs checkpoints with exact-head evidence before switching or stopping.

## Retrieval review checkpoint

PR113 full repaired native/source logs inspected. Six capture files exactly reproduced from contiguous3+4+2+4+5+4 chunks; each complete UTF8 content and local SHA256 matches retained bytes. New availability input32345bea44a162a1518f9022c93dba3cf5640980c94fb1d3734e172b780cc245/outputd3e8186fb80ffc9a47191228e5cad222636ffdfea1de3c837747bf896acbaef9. All previous controls plus four compiling availability oracle mutations/restored baselines pass both hosts. Count-budget scope fuzz completed500000 executions; property/corpus/workers unchanged. Source14 auxiliary/13 inventory/30 normalization/11 availability/eight probes, exact source/tree/APRL provenance. Zero reachable/imported vulnerabilities, module-only advisory open. This verifies repaired proposal evidence, not feature acceptance or audit-wide closure.

Manual production review has covered CLI/app entry points, all Azure/discovery/stages/rules production files, ARG transport/client/executor, Advisor/Policy/Arc/Cost adapters and throttling, with selected independent fixtures. Orchestration and Defender remainder, comprehensive tests, reports/privacy, plugins, tooling, workflow security, requirements/docs/dependencies remain open. Review matrix status remains IN PROGRESS until associated tests/requirements are reconciled. Cost Management continuation/column handling is a newly identified contract question, not yet a confirmed finding; inspect official pinned API/source and independent fixtures before classification.

AUD003 deliberately retains SDK buffering and all source pins. Default shared pipeline mitigation is source-grounded in azcore1.23.1 bodyDownloadPolicy and internal1.12.0 Payload. New actual pipeline fixture must verify it natively. Earlier Diagnostics inspection does not imply the same production reachability. The original eight failures are retained at immutable730a29e6. No full-project review completed, accepted branch advanced, independent-person review or live approval.

## Cost and report review checkpoint

AUD005 test-only78dd85c9 proved paged empty/nonempty first responses are accepted without warnings on both platforms. Guard in PR116 rejects nonempty nextLink, preserving null/empty terminal mappings and valid204. This intentionally reports partial data rather than implementing full pagination; coordinator fixture must prove failed Cost/overall partial with healthy Advisor/Graph retained. Official API snapshot inspected through Azure/azure-rest-api-specs stable2021-10-01 definitions/SubscriptionQueryGrouping example. Missing envelope/column order are still separate unconfirmed questions.

Production source review extended to all canonical assessment/result/findings/gate/config/branding/redaction/reportfile/renderers and public zone/carbon/service-health/SQL-EOL adapters plus plugin projection boundary. Read core result ownership/ordering fixtures, findings fixtures, JSON embedded/scope/error-only identity tests, CSV formula/private output tests, actual XLSX text/hyperlink/workbook tests, SARIF record tests and Windows ACL inheritance/replacement fixtures. Existing native branded all-format checks also passed on115. Remaining comprehensive plugin/core integration tests, comparator/tooling/pinned-data/requirements review keep these areas IN PROGRESS, not accepted full audit completion.

Documented report limitations retained: SARIF identity-bearing by design; report staging protects existing files on rendering failure but multi-file exports are not transactional and universal crash durability/atomic replacement is not claimed; Windows fixture proves controlled ACL inheritance and prior DACL preservation, not arbitrary operator-directory privacy. JSON masks known IDs in serialized related strings, core human tables preserve narrower inherited dedicated-ID masking; audit should clarify any broader operator expectation rather than claim complete anonymization. Table constructor assumes canonical Build summary for core trusted library inputs; default production always supplies it. No exploit or live failure inferred from trusted injection/malformed manual object hypotheses.

## Current comparator and plugin review checkpoint (2026-10-04)

All production AI-governance, six accepted region helper files, orchestration/Defender remainder, equivalence projection/compare/types, scanner registry and SKU helpers were read. Existing tests were reviewed or sampled according to each area; comprehensive requirements/test adequacy reconciliation remains open. Public region service execution and unpublished availability implementation are absent from accepted code and are not treated as audited implemented features.

PR117 test-only negative control preserves accepted production and uses literal real-command inputs. Native failure evidence is pending; no confirmed AUD006 execution result yet. Known comparator boundary remains generic semantic comparison, not source operator-command parity or proof of identical Azure snapshots. Possible nested registry slice ownership and contradictory manually constructed stage metadata remain separate source-inspection questions, not claimed production incidents.

PR116 fresh full QA and source evidence passed as above. Linux82.7% coverage, race/vet/existing fuzz/mutations, both host actual CLI/packages/docs/branding succeeded; zero imported/reachable vulnerabilities and one module-only advisory remain. This is proposal-specific QA, not project-wide review completion or new accepted-push evidence.

## Comparator correction and tooling review (2026-10-04)

AUD006 now reproduced/corrected as above. All first-party production Go paths in the41-area inventory have been read, including ARG finding mapping, Diagnostics support metadata/scanner, branding builder and Diagnostics/equivalence tools. This is source coverage, not test/requirements or full-audit closure. Diagnostics cancellation is guarded by the stage runner before/after execution on the default route; injected direct-scanner cancellation behavior remains a lower-level ownership question. Pinned non-success diagnostic-subrequest warnings deliberately preserve inherited missing-setting findings and historical request-correlation uncertainty, not healthy completeness.

All four workflows reviewed for events/permissions/pins/native failures and proposal-only maintenance paths. Source capture isolates test-only harnesses and guards bytes/provenance. Production scripts bootstrap, maintenance publisher, branding canonicalization, package/inventory generation, PowerShell evidence/scope/stage helpers and exclusion preparation were read. Remaining31-script adequacy work concerns22 synthetic harnesses/tests plus source capture integration already examined. Native gates execute focused fixtures, not hosted-maintenance dispatch/token behavior. Packaging checksums prove integrity under a trusted artifact input, not signed publisher identity/fresh OS qualification. Live runner evidence is private/read-oriented, retains scope/pins/exits and requires the operator; it was not executed in this audit.

One local read-only execution request stalled and was terminated; review continued using already verified saved source bytes and native GitHub tools. No mutation was retried blindly, no source changed, and this process observation is not a product failure.

## Additive schema and harness checkpoint (2026-10-04)

The initial unmerged733b9243 comparator guard incorrectly admitted only core1.0, overlooking declared plugin schema1.1. Requirements reconciliation caught this before acceptance. Test-only1bf98b7a native37185951913/Linux111387770853/Windows111387770760 independently failed the one named valid additive-schema assertion; formatting/compilation passed and earlier invalid/healthy checks did not fail. Final correctione381d0ad accepts result.SchemaVersion and PluginSchemaVersion. Exact final QA pending; prior successful733b9243 Linux run is superseded, not final acceptance. No accepted code affected. FN076 records cause/prevention.

All31 scripts (9 production,22 harnesses) and4 workflows now read. Harness review confirms actual isolated builder/CLI/package execution, named selected mutation baselines/failures/restoration, real recipe extraction/local Git proposal rejection, controlled byte/oracle and XML workbook comparisons. The source runner/source harness fixtures remain synthetic pure captures, not live service evidence. Existing native logs confirm expected cases, but representative load/fresh OS/hosted token dispatch remain deferred. All Azure/discovery/filter/stage/gate/throttling/findings/redaction test files have now been read as well; remaining Go integration/plugin/result/rules tests still open.

Filter YAML unknown-field and multiple-document handling is a newly identified scope-configuration question: yaml.Unmarshal ignores unknown keys while the target explicitly uses a neutral schema, so accidental legacy/typo fields could be ignored. Independent pre-auth/load regression is the next investigation, not yet a confirmed live scope incident or accepted correction.
