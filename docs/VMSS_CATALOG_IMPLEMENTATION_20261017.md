# VMSS catalog diagnostic implementation, 2026-10-17

Current authority: [CURRENT_STATE](CURRENT_STATE.md) and this branch's live draft PR. Frozen behavior and61literal cases are in the [contract](VMSS_CATALOG_CONTRACT_20261017.md). Parent is verified draft PR153/ccfee397c4c12da80b12d0b8c5d0043836d3f543, tree24089ad9240fa15eb41cdff37195190b838fe262. This proposal remains unaccepted and scoring stays unsupported.

## Preparation

The first revision adds only a consuming test runner and recovery records. Explicit JSON-tagged wire structs cover every underscored manifest field, with recursive unknown-field rejection and trailing-JSON rejection. All61cases retain their original bytes/hash. Full catalog equality includes binding, status, not_assessed, maps/zones and fixed warnings; errors require a zero result and exact sanitized or canonical context errors. Collector-backed cases verify literal URLs/counts and1MiB request bounds with an injected getter. Direct cases cannot call a getter. Inputs are snapshotted; the ownership case probes input mutations and independent output storage. A separate consumer test checks unsupported scoring after catalog projection using the manifest's literal arithmetic/health/warning.

The API is absent. Local uncached region compilation fails at undefined VMSSCatalogEvidence and VMSSCatalog only. These are missing-API preparation failures, not a product pass or reproduced product defect. No61-case execution is claimed. Publish this runner and inspect both Linux and Windows actual preparation logs before adding production code. Record exact preparation head/tree/parent, runs, jobs and full log hashes in the live PR body; preserve the first runner in history.

## Authorized completion

Implement only the contract's pure diagnostic projection using existing validation primitives. No parser/pager/credential/network/CLI/report/coordinator or SKUEvidence conversion is added. Preserve old VM Evidence budget/semantics. Supplement the frozen cases with independent exact-bound/one-over count/work/text/label/output, late cancellation/deadline, sanitization, owned map/zone/warning and unchanged unsupported scoring checks. Publish compiling reversal/restoration controls for binding, knowledge, restriction, sanitization, ownership, limits and unsupported scoring through the conserved native CI inventory.

Require fresh complete final-head Linux/Windows/source QA, all predecessor385/380control identities plus deliberate new controls, full job/step/log review and reconstruction/readback of all23captures. Archive complete synthetic evidence with the existing bounded verifier; final exact-head outcome belongs in the PR body without an evidence-only commit loop. Update compact recovery records at material checkpoints. No independent-person or live proof is implied.

## Rollback and limits

Use a new targeted revert under [ROLLBACK](ROLLBACK.md), reviewing dependencies and reverting this implementation with its runner as a coherent unit. A runner-only checkpoint is deliberately uncompilable until the proposed API exists; it is unsuitable as an operational release target. Preserve the separately verified PR153 contract baseline, prior frozen oracles/source captures and every existing draft head. Re-run relevant behavior and full required QA after reversal; do not reset shared history.

No Azure operations, fixture/RBAC changes, dependency/action/SDK/source/capture/public-schema refresh, subagents, protected merge/bypass or release. Broader independent review/backup/rollback/live/release obligations remain in [ACTION_REGISTER](ACTION_REGISTER.md).
