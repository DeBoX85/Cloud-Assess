# PR154 final exact-head evidence, 2026-10-17

VERIFIED OFFLINE code candidate `ec191ba56c2c2c7bdf9360c8037952f1070894c8`, tree `45de4bf7e50d92107d3f1125ff6e18aca1ea5d86`, parent published runner `b64e185478fd5230d5c9e93b56b1420311cdc752`. Preview `65294bd4c89ef5c2d8a214052fb586bf989257e0` has the identical tree and ordered protected/candidate parents. Commit/tree/parent/human identity and all14complete changed files read back. PR154 is draft/open/unmerged; protected acceptance/release have not changed.

Quality [39005181223](https://github.com/DeBoX85/Cloud-Assess/actions/runs/39005181223): all eight native workers and both required gates succeeded. Every complete log, step outcome and worker runner label checked. Only the correctly conditional Linux failure-artifact step skipped. Native actual CLI/default/custom package/installed checks, race/vet/bounded fuzz, module/license/branding/documentation/provenance passed. Total coverage85.7%, unchanged floor75%. Zero reachable/imported vulnerabilities; existing one module-only advisory remains OPEN.

- linux: 409 controls; no predecessor identity missing; exactly24new; full-lines-plus-LF SHA256 `5a7eb27b96aa6f49783ab7535a6fab7930f846764dd8383c3ac255c8a27230f6`.
- windows: 404 controls; no predecessor identity missing; exactly24new; full-lines-plus-LF SHA256 `b53c00c4aca8dd35dd27101f1d87ed2e0e5ad0b719632107ecb465eb3128ce3d`.

Source [39005181209](https://github.com/DeBoX85/Cloud-Assess/actions/runs/39005181209), job117316070786, succeeded with all steps successful. All23captures reconstructed from315contiguous chunks/909865UTF-8bytes; emitted hashes verified and every full file compared with its exact candidate Git blob. AZQR8e4f0577f3615e6c9014c031bcad079f235369cc, source-tree17d93b20c303f90f7843036be82f0dc32f3260f1 and APRL60eaddda76541f6adbc1c5ffa686829807e55e29 unchanged. Complete source log: 1359869bytes, SHA256 `0ec845658737d7528ac768609e28008d9f41a56bbf385071e6c6f45965c37a28`.

All61unchanged frozen VMSS catalog scenarios and independent exact/one-over counts/work/labels/input/output text, canonical late context, sanitization, owned concurrent outputs and unsupported scoring checks passed. All24new compiling reversals failed named assertions and restored/passed on BOTH native hosts. The old385Linux/380Windows identities are retained. The original67VM oracle and source captures remain unchanged. First runner/missing-API preparation and QA-VMSS-01/02 test-isolation corrections remain in history and the separate committed preparation bundle; missing-API failures are preparation, not product passes or reproduced product defects.

Local synthetic Git rollback reverted implementation+runner as one unit:87e1ed4afbec9b877e45727ef84c6301faf1a81a exactly reproduced PR153 tree24089ad9240fa15eb41cdff37195190b838fe262 and full region tests passed. Restore83b4059f6f94e3cbeecb7f7d16b8202f4b5521a1 exactly reproduced candidate tree45de4bf7e50d92107d3f1125ff6e18aca1ea5d86 and full region tests passed. Complete command/results retained. This is local Git proof, not installed/live/dependency/platform-wide rollback or independent backup.

Worker timing: first start2026-10-17T05:08:48+00:00, last completion2026-10-17T05:20:50+00:00, observed span722seconds. Summed eight worker-job duration3850seconds is parallel job time, not elapsed, CPU or billed time and not a speedup claim. Full per-job start/end and step timestamps retain queue/setup observations in raw metadata.

| Job | ID | Complete decoded bytes | SHA256 |
| --- | --- | --- | --- |
| windows-native | 117316070760 | 42310 | `dcde9a6bd5b8369baeb01121a7b796448a7f9fbd8fa61bda1a92c64729fc148a` |
| windows-controls-3 | 117316070891 | 44530 | `751e799dbb067e3a874f7f7a402a8f1a3ccb3fdf40a5006f1503d4ba6d66f63e` |
| windows-controls-2 | 117316070892 | 45508 | `f1f9942e13b2413234e50455e4706289f3c6df2d9ac316df8127e1bdce19b66b` |
| linux-controls-2 | 117316070900 | 43010 | `82dea665250f1b624d17b0798657857c37155a9d1d6c04364d776a9028f66c73` |
| linux-native | 117316070902 | 63690 | `a584591a6da47335bec73d1be4e30c2908d4735274e9b952c1b1bff29d27864e` |
| windows-controls-1 | 117316070904 | 44451 | `b13ff070bb76b9427d0551102ab31e9b958897df4339f94c26ebb71547764d7b` |
| linux-controls-1 | 117316070938 | 41479 | `0ca19721a883d018e9b80e4ae6d542347928bd499d652f80efcdcee545abbf50` |
| linux-controls-3 | 117316071061 | 47252 | `1a84d29e85b4941819957ee011cbade12f6a4adf87f96da8c84814e5d2042e99` |
| quality | 117317292096 | 16412 | `d0ad67a5a7e4a50f3f170cb9aae857456f87bb89423cb49fc17853254d11420e` |
| windows-validation | 117318276294 | 16440 | `84dc18f825268fce66e00f3a50aeb8d64e4268aa9cde6749a861727a4f1ccf81` |

## Storage and recovery

This artifact-only branch `qa/evidence-pr154-20261017` has candidate code as parent and adds only bundle/manifest/index. It is NOT a development successor, protected proposal, release or code certification for its own different tree. The final PR154 body links its actual verified artifact commit. Existing bookmark routes PR151 -> PR152 -> PR153 -> PR154; all earlier heads remain preserved. Final outcomes belong in that PR index; committed pending text is a publication snapshot, avoiding an evidence-only certification loop.

Bounded bundle: 28members, 1877726UTF-8memberbytes, 228481bundlebytes, SHA256 `df510b19d43d4a02d76f0027d0bd0751ab2adc6a1405576875ba612ffa6a65d6`. Verify with `python scripts/verify-evidence-bundle.py --bundle docs/evidence/20261017/pr154-final.json.gz.base64 --manifest docs/evidence/20261017/pr154-final.manifest.json` before explicit restoration. No extraction or execution is performed by the verifier. Integrity is not publisher authentication, immutable storage, independent backup or a privacy certification. Logs are synthetic public CI/local evidence.

The retained portable review helper uses its own directory as evidence root and takes quality-run, source-run, candidate, preview and optional repository-directory arguments. Restore the predecessor38988652561logs from the already committed PR150/151 bundle when needed; its fixed full-line identity digests are enforced before comparison. Historical workspace paths are hints, not recovery prerequisites. The builder takes a clean artifact worktree at the exact candidate as its output directory. Do not execute restored code without review.

The diagnostic checkpoint is complete. VMSS scoring/deployment eligibility remain unsupported/not_assessed, accepted target adapters1/7 and source characterization7/7. Public coordinator/report integration remains separate. No Azure, fixture/RBAC changes, dependency/action/SDK/source/capture/public-schema refresh, subagents, protected merge/bypass or release. Independent review/backup/actual repository retention, broader operational rollback, remaining collectors/integration, live scale/permissions, module advisory, fresh OS/security/license/provenance/Gate004/release obligations remain OPEN.
