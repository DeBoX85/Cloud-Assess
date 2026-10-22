# Section 4 corrective checkpoint, 2026-10-16

Review and contract reconciliation are complete; final corrective acceptance is IN PROGRESS. All 11 baseline Go files in the four remaining packages, the independent SR-F12 regression and all 62 initial script/workflow/support files received full reads. [Contracts and exact manifest](SYSTEMATIC_REVIEW_SECTION4_CONTRACTS_20261016.md) and [final assessment](SYSTEMATIC_REVIEW_FINAL_ASSESSMENT_20261016.md) reconcile architecture, data/provenance, current documentation, findings and remaining work. The complete four-section index accounts for all 251 baseline Go files and all 41 packages, plus 12 added Go regressions. This does not credit Section 4 completion before final QA or measure project completion.

SR-F12 production correction is implemented after the independent test/contract publication at 10efcd48f8fc3dd832350f03c7b4abd37bf0d5ad / tree 5e41c5dafb6a8f29c7f6a40e1a284f21200713fc, the direct parent of this candidate. The private `put` rejects conflicting projected fields by permanently setting noncomparability, keeps the first row for diagnostic deltas only and emits one sanitized sorted dataset note. Identical normalized records remain deduplicated and resource-type count aggregation is unchanged. No identity, numeric/provenance/SLA normalization, public schema or source input changed. Original independent regression blob 126f45bc1980a750245474d70ee7fe490f087917 is unchanged.

Both native hosts reproduced exactly the nine declared child failures in native 38856474803: seven false exit 0 and two missing-comparability-precondition assertions. Linux formatting passed; no build-failure or data-race marker appears in either complete log. All earlier completed steps passed; later skipped steps are not accepted. Source 38856474640 passed separately with eleven harnesses, all 23 captures / 315 chunks / 909,865 bytes and source/tree/APRL identities independently reconstructed/hashed/compared. Its full log is indexed in the contracts record. Earlier c8a15c7 native 38855113862 was canceled when superseded; it has no accepted native proof.

| Complete reproduction log | Job | Characters | UTF-8 bytes | SHA-256 |
| --- | ---: | ---: | ---: | --- |
| linux | 116829220849 | 69360 | 69365 | `ae6c321d702dbb72b47c1a7708fb75a28730e82af1bc10ef3bf5038d500c3e27` |
| windows | 116829220610 | 56115 | 56117 | `ccd9e52cb638c3857397e326756a1b3812040976dd9a68281c2e5e5634ffd6c1` |

The new `scripts/tests/section4-regressions.py` control compiles the deliberately disabled conflict branch, restores actual old last-row overwrite behavior, requires the named independent assertion failure, restores exact original bytes and requires the clean baseline. Both existing native jobs now run it; every previous step/control/timeout/floor/pin remains. Local control/restoration, all four Section 4 package race tests and vet pass. The full local race suite and full vet subsequently passed; their log hashes are recorded below. Source-mutating controls and source-consuming tests were run serially.

Full exact-head Linux, Windows and source QA is PENDING AT PUBLICATION. The live [PR142 evidence index](https://github.com/DeBoX85/Cloud-Assess/pull/142) records actual final candidate/tree/parent/tested preview/runs/log hashes/control identities and completion, under QA_PROCESS's rule avoiding documentation-only CI loops. Do not use earlier green as this candidate's acceptance. Keep PR142 draft/unmerged and protected a6fee427 unchanged. No protected merge/accepted-push or independent-person approval is claimed.

NEXT: finish this exact head's full QA and final diff/ref readback; if complete in PR142, STOP at the final audit checkpoint for user re-evaluation/disposition. Do not replay SR-F12, start feature code or merge. After the audit disposition resume Storage contract/tests in draft PR141 at 6cd002cb77fcead75a602b06a75625a178a84e49, then remaining region work. No Azure operations, dependency/pin/capture/schema refresh or sub-agents.

Rollback: targeted ordinary revert of this candidate's production guard and new runner/workflow steps together, after checking later dependencies. Preserve the independent regression, all original controls and historical records; the known failure intentionally becomes red. No protected reset, assertion removal, source capture regeneration or data migration.

## Completed local corrective QA before publication

Full local `go test -race -count=1 ./...` passed all 41 package directories (the development builder has no Go test file and is exercised by native binary integration). Full `go vet ./...` passed. The Section 4 compiling reversal/restoration control ran and restored exact production bytes BEFORE the focused and full suites. Pinned formatting is clean, all 50 current Python files parse, shell syntax and authored links pass, the independent SR-F12 regression remains blob126f45bc1980a750245474d70ee7fe490f087917, and dependency/source data are unchanged. Full native/Linux/Windows/source gates still require their own exact-head results.

| Local log | Bytes | SHA-256 |
| --- | ---: | --- |
| section4-control-local.log | 96 | `fa09ae7addba9f013b3bac988404bd21e6db72d30ab2d5e273807fdbb9ecb44d` |
| section4-corrected-local.log | 266 | `b0314d355c6e3d4cef9db29ee8b63eb2e97705649908117edce4ad218aa55d86` |
| section4-full-race-local.log | 2618 | `180f2676d81e3cb9f4e8d80171897a13bd5d75b4df2dcc7b583253599b949bb6` |
| section4-vet-local.log | 0 | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| section4-inventory-local.log | 31 | `e258f36d40a08392b39a931c079fda94ef0d79aedb906fd0dd8599087cdd2988` |

## Historical entry and published pre-edit contract below

# Section 4 systemic review, 2026-10-16

Status: IN PROGRESS. The user re-evaluated the completed Section 3 checkpoint and authorized continuing in proper order. This supersedes earlier instructions to wait before Section 4, without changing the draft/unmerged PR142, no-Azure or no-refresh boundaries. No merge or release approval is inferred.

Entry head: `98b8915fdcb934435b568fd5c61b5500d26b592c`; tree `255fbfce40e2b05318cfe91c55fd08747f79dbbe`. Protected accepted baseline remains `a6fee427e19feeef28691737ea9bc5026086efa0`. Both live refs and PR142 draft/unmerged state were freshly verified. Section 3 full exact-head native38841447642/source38841447701 evidence remains indexed in PR142. Sections 1–3 are complete; Section 4 is not.

## Scope and sequence

Review all source/tests in `internal/equivalence`, `tools/brand-build`, `tools/diagnostics-probe` and `tools/equivalence`; then all Python/PowerShell/shell scripts, CI, dependency/license/provenance checks, configuration/queries and current documentation. Reconcile the whole inventory against the preserved baseline manifest and Sections 1–3 contract maps. Initial reads are not completion credit.

Follow QA_PROCESS and DEVELOPMENT_EXECUTION_PLAN: independently specify a substantiated failure and correction contract before production edits, retain failing expectations, reproduce after compilation, isolate corrections, add meaningful compiling reversal/restoration controls, inspect full exact-head Linux/Windows/source QA and publish verified recovery evidence. No source capture, pin, dependency or schema refresh; no Azure operations. No sub-agents.

## Workspace recovery

The old Windows workspace and caches are absent. A fresh Linux clone at `/workspace/scratch/c04791396c67/Cloud-Assess` has genuine remote ancestry and an exact entry tree, clean before edits. The existing checksum-verifying bootstrap is provisioning separate tools/reference/accepted target under `/workspace/scratch/c04791396c67/dev`. Provisioning is not a test pass. No Azure command is invoked.

## Active review notes

- Read all four remaining Go production packages and their existing tests. Comparator map insertion can overwrite an earlier record with the same normalized identity. Conflicting duplicate input is a suspected false-equivalence boundary; independent literal regression and prior-contract reconciliation follow before a production change. Existing documented resource-type count aggregation and identical duplicates must be distinguished.
- Read maintenance/source workflows, native QA workflow and scoped maintenance publisher. Maintenance publication is proposal-only; hosted maintenance execution and independent-person/release approval remain deferred. No workflow dispatch or dependency generation is authorized by this review.
- Remaining scripts, data/provenance and full documentation reconciliation are pending. No Section 4 package receives final completion credit.

## SR-F12: conflicting duplicate equivalence rows, P2, reproduced and UNFIXED

`internal/equivalence/projection.go:put` overwrites records at the normalized key without checking projected fields. Literal independent real-command Advisor rows with the same identity and different impact reproduce false exit0 when the last row agrees; identical contradictory rows on both sides also yield false PASS. Reversed order can yield a changed-field exit1 but still incorrectly remains comparable. All nine conflict cases fail their assertions after compilation locally on Go1.26.9/Linux. Normalized identical duplicate and resource-type2+3=5 controls pass. No historical Azure report incidence is established.

BEFORE production edits, correction contract: retain existing identity and field normalization; when `put` encounters the same key with unequal projected fields, set the owning projection non-comparable permanently and add one dataset-level diagnostic without raw identities or values. Keep the first row only for diagnostic deltas, never use deltas as a valid comparison after the precondition fails. Repeated identical normalized records remain deduplicated. Preserve `addCountRecord` resource-type count aggregation; no identity expansion, source correction, public JSON schema change, SLA-rule rewrite, numeric normalization change or dependency change. The real command must write an explicit non-equivalent report and exit1, not silently pass or panic. Later identical rows must never repair a conflicted projection.

Keep `tools/equivalence/duplicate_evidence_test.go` unchanged after this pre-edit checkpoint. Add a compiling reversal/restoration control for the private conflict branch to both native jobs; old behavior must compile and fail the named assertion, followed by exact-byte restoration and passing baseline. Full corrected exact-head native Linux/Windows/source QA remains mandatory. This regression-only head intentionally has red native gates; both-host failure logs must be inspected before correction. Source QA success is separate evidence, never corrective acceptance.

Rollback: ordinary targeted revert of production correction and its new control/workflow steps together after checking dependent changes. Preserve independent expectations, which should turn red again. No reset/force/protected merge, Azure operation, dependency/pin/schema/capture refresh. Section4 remains incomplete and PR142 draft/unmerged.

Local reproduction log: 3677 bytes, SHA256 `71ea80820b59d924a18ce3a68b76389f0c031df99ea37f30e8c2cf946aec0f29`. First local attempt failed setup because this new worktree had no APRL contents; it is not reproduction credit. Initialized the unchanged60eaddda submodule and reran successfully through compilation, yielding only the nine intended conflict assertions. The setup log was inspected before rerunning.

## Recovery and rollback

Resume from the actual live PR142 head and this report, not old scratch paths. Preserve historical logs, captures and independent expectations. Use ordinary targeted revert commits after inspecting dependent controls; never reset protected history. Keep PR142 draft/unmerged. After the full audit, reconcile its disposition before resuming Storage pre-edit contract/tests in draft PR141; Storage is not implemented by this review.
