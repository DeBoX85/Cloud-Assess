# Narrow dependency/advisory review, 2026-10-20

Current authority: [CURRENT_STATE](CURRENT_STATE.md); [bounded approved contract](EXCELIZE_SECURITY_REMEDIATION_20261020.md). User approved this exception at16:14Oslo. Published exact candidate/native/source/evidence outcomes belong to the live maintenance successor linked from PR177. These local results are not protected acceptance or release clearance.

## Actual minimum delta

Root go.mod changes only Excelize2.11.0 ->2.11.1-0.20260915055537-22f76f9acb94 and required indirectmscfb1.0.7 ->1.0.8. Final tidy go.sum adds six required checksum lines and removes six obsolete-version lines relative to177. Original dependency pairs remain recoverable from immutable177and retained rollback evidence; obsolete sums must not remain in the active file. Complete selected graph additionally raises x/image0.38.0 ->0.44.0 as an upstream requirement. All other versions, Go1.26.9, workflow/action/tool pins, AzureSDK/source/APRL/capture/public-schema bytes are unchanged. There is no replace/exclude/vendor/suppression. Root and selected graph must not be conflated.

Unchanged generator recreated49module inventory and full notices. Only three module entries/version headings/input hashes change; license/notice hashes, two-platform CLI membership and bundled source evidence are unchanged. x/image is selected but not imported by either CLI. mscfb0.8 changes only mscfb.go: reject an oversized untrusted mini-FAT sector count before allocation. No legal approval or independent audit is inferred.

Full Excelize baseline-to-fixed patch retained:45files/1461additions/825deletions, including unrelated-in-functionality but inseparable intervening upstream commits. All production hunks/file inventory reviewed: formula evaluation/cycle/Unicode/count bounds, column/shared-string/row/style validation, encryption input/spin/allocation guards, streamed duration/float precision, conditional-format validation, worksheet panes, table/form relationships, chart/fontXML ordering and comment/refactor changes. Source renderer uses strings and explicit hyperlink formulas; no generic incoming workbook/decrypt/formula-evaluation feature added. Regression expectations remain independent. Upstream workflow changes do not alter our CI/toolchain. Full raw patch plus affected test-hunk inventory retained; self-review is not independent-human review.

Go resolves exact22f76f9acb94b85b3eb9c4365ab4f750cebbc565 and verifies module origin/sums. Both fixed advisories and upstream original fixes are retained in the earlier blocker archive. [Go module documentation](https://go.dev/ref/mod#pseudo-versions) supports canonical commit/timestamp/base-version validation; the actual tool-resolved version is used, not a handwritten latest tag. Current-module verification passed.

## Local current scanner evidence

Pinned govulncheck1.8.0, Go1.26.9, symbol/source mode, live database last-modified2026-10-20T00:20:56Z:0symbol-reachable,3imported-not-called and1module-only findings. GO-2026-6684 andGO-2026-6688are no longer findings. Scanner returned successfully without suppression or backdated data. Native Linux/Windows exact-head scans must independently confirm actual outcomes.

| Residual advisory | Local classification | Disposition |
| --- | --- | --- |
| GO-2026-6685 | Imported package, no called vulnerable symbol detected | OPEN; no waiver or exploitability conclusion |
| GO-2026-6687 | Imported package, no called vulnerable symbol detected | OPEN; no waiver or exploitability conclusion |
| GO-2026-6695 | Imported package, no called vulnerable symbol detected | OPEN; no waiver or exploitability conclusion |
| GO-2026-5932 | Module-only OpenPGP | OPEN; previously documented residual, no new integration |

Raw JSON contains consideration of many historical OSV entries; that is not the number of applicable findings. Retained review classifies each unique current advisory by highest-detail trace and distinguishes module/package/symbol findings. No advisory-free graph, practical exploit proof, broader remediation or release approval is claimed. A changed reachable result or unrelated needed refresh reopens reevaluation rather than silently expanding this update.

## Behavioural proof and correction record

Two bounded child-process public-API tests reproduce the ORIGINAL dependency's circular ANCHORARRAY fatal recursion and malformed standard-encryption header panic; both compile and fail by named assertion. The same unaltered tests pass on the updated pair. Child-only stack bound/timeout protect QA, not production settings. A synthetic8192byte aligned container passes exact valid decrypt roundtrip before corrupting only its unique header-size field. Initial tiny41/32byte fixture roundtrip mismatches occurred in both versions and were excluded from defect proof; the precondition was corrected with a realistic CFB stream size, not a weaker expected error. Cause of tiny-payload mismatch was not diagnosed as a product issue; no encrypted-input feature was added.

Additional actual-renderer cases preserve exact long/leading-zero numeric-looking strings, Unicode/non-BMP text, formula-looking text, header and nonzero styles across both dependency pairs. Original renderer/report/branding fixtures are unchanged. Atomic old dependency+sum rollback reproduced both safety failures while old ordinary/new compatibility report tests passed; updated pair and all test bytes restored exactly and regressions passed. This is isolated source/dependency recovery, not installed/live rollback.

Local full uncached preflight/tests/vet/format/static/gate/handoff passed before the final additive compatibility case; focused renderer suite passed afterwards. Final preflight and hosted full gates still required. Read-only import graph listing needed buildvcs=false after incidental local VCS metadata lookup failure, as already used by the unchanged inventory generator; no build/CI flags or validation rules changed. Tests and actual builds retain their original checks.

No Azure/tenant/permission request, original oracle/capture rewrite, protected merge, bypass or release. Complete candidate readback, own all-worker/both-gate/source results, control conservation, package/install proof and durable archive are pending at this authored checkpoint.


## QA-EXCELIZE-02 checksum freshness correction

Initial179/9e163066was published before running the exact workflow tidy check. Its complete Linux-native118686117581log in39439734077shows go.sum drift and exit1; preserve this failure. Go get alone left obsolete original-version sums and omitted one selected x/image go.mod sum. Local exact tidy check reproduced the same diff. Corrected successor retains every runtime version/test and changes only go.sum plus inventory input hash and accurate documentation. Original notices bytes unchanged. Require own full fresh QA at successor, never relabel179green or suppress tidy. Before future dependency publication run go mod tidy, verify go mod tidy -diff is empty and regenerate/check inventory. No broader refresh is required.
