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
