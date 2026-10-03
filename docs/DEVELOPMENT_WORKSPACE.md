# Isolated offline development workspace

Status: executor recovery VERIFIED OFFLINE through PR100 mergeca7fc301de421f24f62c6b976fe2d79c0bfe6780, with required exact-head and separate accepted-merge native Linux/Windows QA and clean local/source/pin readback. This supplies the missing development tools and fresh source/target checkout proof from FN056. It does not close Azure/laptop, Gate004 or release validation.

## Recreate the workspace

Use [bootstrap-workspace.sh](../scripts/bootstrap-workspace.sh) on Linux amd64. The current executor is Ubuntu24.04 with working Git HTTPS, curl/CA trust, tar, sha256sum, Python3.12.14, ripgrep, GCC13.3.0 and Node24.19.0/npm. These prerequisites remain executor-provided; the script checks its required commands before downloading. It does not install system packages or change global profiles/proxy/credential settings.

```bash
bash scripts/bootstrap-workspace.sh /workspace/scratch/new-cloud-assess-workspace
source /workspace/scratch/new-cloud-assess-workspace/env.sh
cd /workspace/scratch/new-cloud-assess-workspace/target
```

The destination must not exist and its parent must exist. Failure retains a partial task-owned workspace for inspection; choose a fresh destination after resolving the failure. No existing checkout is reset or cleaned. Paths containing spaces have been exercised end to end. Re-read live handover/refs/open proposals after creation; a manifest records an observed checkout, not a promise that core-v1 will remain unchanged.

| Tool | Version/source and verification |
| --- | --- |
| Go/gofmt/compiler/race tooling | Exact1.26.8 required by target go.mod/workflow; official go.dev Linux amd64 archive, SHA256 d0f743b33e8d8945e6b1f432edd15785c70507121d6e2a723b21285eddf8b57b. GOTOOLCHAIN=local prevents silent toolchain replacement. |
| PowerShell | Official [7.6.6 release](https://github.com/PowerShell/PowerShell/releases/tag/v7.6.6) Linux x64 archive, SHA256 ddbc4a2d113bbd46d283cfedcbcd117a70caefd7673f41f2b4e0000badf103bc from the publisher release-asset digest. This is local PowerShell7 tooling, not Windows PowerShell5.1 proof. |
| govulncheck | v1.8.0, matching native CI; installed into workspace tools/bin with normal Go module checksum verification. Vulnerability-database access was exercised by the real scan. |
| Git/Python/ripgrep/GCC/Node | Executor-provided; Git2.51.1 has its HTTPS helper, Python3.12.14, GCC13.3.0 and Node24.19.0 verified. GCC was exercised by the full Go race build/tests. Python scripts use the standard library. |

Tools and Go module/build caches are workspace-local. env.sh sets readonly module flags, Go compilation parallelism4 and GOMAXPROCS4; it preserves the executor's managed network settings. Both target and unchanged reference dependencies are downloaded and verified before work. No Go/module/source pin update is performed. Bootstrap itself does not run an Azure command or assessment.

The recipe clones current bootstrap/core-v1 into target, configures the recorded human Git identity there, and clones AZQR into a separate azqr-reference checkout at8e4f0577f3615e6c9014c031bcad079f235369cc. Both actual APRL checkouts are verified at60eaddda76541f6adbc1c5ffa686829807e55e29; target AOR/CUSTOM/SKU provenance is checked too. A workspace-manifest.json records tool hashes, observed HEAD/tree/parents and qaExecuted:false. Downloads and tool verification alone are not a build or QA pass. Reference captures must use isolated copies; never edit the pinned checkout.

## Verified recovery evidence

Historical PR100 workspace: /workspace/scratch/d06563362caa/cloud-assess-dev, with target and azqr-reference under that directory. Source env.sh for every fresh shell. The bootstrap additionally recreated /workspace/scratch/d06563362caa/rebuilt workspace end to end, using fresh checksum-verified downloads, fresh clones/submodules and fresh verified module caches. Both target checkouts observed accepted89087cf1d8afb7e06abfb134cd5d068cfa56186d/tree1088200d83e163f9b5b65128ac1282b37e7ac655; unchanged source8e4f057/tree17d93b20c303f90f7843036be82f0dc32f3260f1 and actual source/target APRL were clean. Older checkouts/tools were left intact.

On the pristine accepted runtime/tree: full Linux race/vet, strict CGO-disabled readonly stamped CLI, eleven actual CLI cases, docs297/9/1, dependency inventory, PowerShell parsing/helpers/native-failure checks, default/custom branding/report checks, nineteen default plus nineteen branded package tests, coverage81.4%, four AI/comparator compiling controls/restored baselines, bounded comparator/scope/ARG/AI fuzz and maintenance fixtures passed. Reachable/imported vulnerability scans reported zero; the known module-only advisory remains open. The binary's vcs.revision was89087cf1d8afb7e06abfb134cd5d068cfa56186d and vcs.modified=false, SHA2564aec5698bc0c29af1d0a32c3c2ded2d81b565d03a191e6bd015ccdad24408401. Local logs are under workspace/logs; final exact-head native evidence belongs in the proposal PR.

The unchanged source AI package also compiled with go test -run '^$'; no source tests ran and no new source capture or live source equivalence is claimed. Compilation proves the cached pinned SDK/toolchain path is available for subsequent engineering.

[Bootstrap boundary checks](../scripts/tests/bootstrap-workspace.py) exercise argument handling, preserved existing directories and corrupt-download rejection before extraction/cloning, with explicit local shims and no external test request. A valid-shell disabled-checksum control must fail the named extraction boundary; restored checks pass. Linux CI runs syntax and these guards. Actual fresh provisioning was exercised separately in this session, not by those synthetic boundary fixtures.

## Build and validation

After verifying a clean committed candidate, use the existing [QA process](QA_PROCESS.md) and [.github workflow](../.github/workflows/test.yml). For focused feedback:

```bash
go test -race -count=1 ./...
go vet ./...
mkdir -p .build
CGO_ENABLED=0 go build -trimpath -buildvcs=true -o .build/cloud-assess ./cmd/cloud-assess
python3 scripts/tests/built-cli.py --binary .build/cloud-assess
python3 scripts/tests/documentation.py --binary .build/cloud-assess --powershell pwsh
```

These commands are development feedback; all existing mandatory exact-head native Linux/Windows gates and protected PR acceptance remain required. Linux does not provide local Windows-native/ACL/fresh-OS evidence. Actual Windows validation remains with native CI. GitHub connector publication is available independently of any shell push credential setup; never request tokens or reconfigure protections to compensate for transport failure.

Scratch paths/caches can disappear. The checked-in recipe and verified remote code/docs checkpoint provide reconstruction, not persistent executor uptime. If the session disconnects again, inspect actual paths/status/processes and rebuild from live accepted refs; do not infer that an interrupted command completed. Resume bounded AI Graph discovery under [AI_GOVERNANCE_EXECUTION.md](AI_GOVERNANCE_EXECUTION.md) after this tooling slice's acceptance. All operator Azure/laptop/live roles/models/load/sovereign, Gate004/release and forensic/advisory boundaries remain open.

Primary setup references: [official Go installation](https://go.dev/doc/install), [Go downloads](https://go.dev/dl/) and the official PowerShell release above. Publisher checksum agreement is verified; no signed release/distribution approval or perpetual workspace availability is claimed.


## Current recovery after scratch loss (2026-10-03)

Actual inventory found the previous tools/checkouts absent. Accepted bootstrap recreated /workspace/scratch/d06563362caa/cloud-assess-recovery-20261003 from live accepted032740d/treebf5880, unchanged pinned AZQR8e4f057/tree17d93b20 and both actual APRL60eaddda. Go1.26.8/PowerShell7.6.6 publisher hashes, govulncheck1.8.0, module graphs/provenance and clean source/target were verified. Source env.sh per shell. Provisioning manifest remains qaExecuted:false. Separate fresh full race/vet/stamped CLI,11 actual CLI/docs307/9/1/eight compiling controls/restored assertions and5000x discovery fuzz passed; binary hash2aa4396a9a76bdb7f32e0dc7f67f20958d5c7852ad2f05de4ed8ff538259f6b2. Current target branch feat/ai-public-execution contains an unaccepted public AI candidate; inspect actual Git status/current refs after interruptions. No Azure scan or new source capture was run.
