# Current security gate state, 2026-10-09

PR132 is UNMERGED and protected acceptance is BLOCKED by a classified vulnerability
gate failure on the unchanged Go1.26.8/x-net0.58.0 toolchain/dependency graph.
Read [HANDOVER_SECURITY_GATE_20261009.md](HANDOVER_SECURITY_GATE_20261009.md) and
[live PR132](https://github.com/DeBoX85/Cloud-Assess/pull/132) for exact failed
revision/preview/run/job/advisory evidence and latest published handover state.
Prior PR131 zero-reachable/zero-imported scans are historical, not current clearance.
No suppression, version change, protected merge or accepted-push proof occurred
in this documentation task. Finish security classification and prepare a bounded
remediation contract before further feature work. Source pins and mandatory gates
remain intact; no laptop/Azure input is needed for that offline preparation.

---

# Dependency advisory review

Date: 2026-10-01 (Europe/Oslo). Baseline: `bootstrap/core-v1` at `8f76c46ad407f202a7aa2d9e8ae7eb8931ffaf35`, the merged alignment/security review in PR #67. This completes the bounded AR-01 advisory investigation and prepares the reviewed dependency refresh. Exact final native CI remains mandatory before merge; release security and license approval remain separate.

## Advisory decisions

| Advisory | Affected code | Decision / evidence |
|---|---|---|
| GO-2026-6355 / CVE-2026-56855 | x/crypto/ssh, malicious-peer connection deadlock; fixed v0.56.0 | Raise the selected indirect module from v0.55.0 to the first fixed v0.56.0. Neither current CLI imports ssh; remove the known affected module version without claiming an exploited application path |
| GO-2026-6354 / CVE-2026-78662 | x/crypto/ssh, pre-establishment channel request deadlock; fixed v0.56.0 | Same bounded upgrade; do not add SSH features or vendor a separate fix |
| GO-2026-5932 | x/crypto/openpgp and related packages, unsafe/unmaintained; all versions, no known fix | Remains a module-only advisory. Neither platform CLI imports these packages. Retain the module because other required packages are used; no suppression, no claim of an advisory-free graph and no unnecessary replacement dependency |

Sources: [GO-2026-6355](https://pkg.go.dev/vuln/GO-2026-6355), [GO-2026-6354](https://pkg.go.dev/vuln/GO-2026-6354), [GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932). These are Go vulnerability database records reviewed on the date above. Re-evaluate when source/platform/plugin imports or database records change. Any future OpenPGP integration needs a maintained implementation and its own reviewed contract; upgrading x/crypto alone cannot repair this package.

## Change and impact assessment

`go mod why -m golang.org/x/crypto` identifies Excel rendering through excelize/md4. CGO-disabled Linux/Windows amd64 `go list -deps` also finds pkcs12, its internal rc2 package and ripemd160. None of those four imported package directories changed between v0.55.0 and v0.56.0. Full downloaded-module file comparison found changes confined to acme, SSH, their tests and go.mod. This limits the expected application impact but does not replace full native checks or certify all cryptographic uses.

`go get golang.org/x/crypto@v0.56.0` followed by tidy changed only that selected dependency and its module/download checksums. Go remains 1.26.8; Azure SDK versions, recommendation/SKU pins, application code and comparator normalizations are unchanged. The selected graph remains 49 modules with 24 Linux and 26 Windows CLI modules. Generated inventory updates input hashes/version/checksums; generated notices change only the x/crypto heading. Its retained LICENSE SHA-256 remains `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`.

Existing required jobs already run reachable vulnerability checks natively on Linux/Windows plus the real executable/compiled-module and package-integrity checks. No redundant implementation-mirroring version test or weakened advisory gate was added. Module-only findings remain visible in scanner output. The inventory remains evidence, not complete SPDX/license clearance.

## Validation and follow-up

Local final full race suite and govulncheck under Go 1.26.8 passed. Native Linux scan reports zero reachable/imported-package findings and one module-only advisory after upgrade, compared with three before. Vet, nine inventory safeguards/freshness cases, the actual built CLI/compiled dependency checks and documentation validation (112 destinations, eight flags, one PowerShell snippet) also passed locally; final native quality and windows-validation must pass the exact published PR head before merge. The PR retains final run/head/tree/merge evidence for indexing at the next ledger checkpoint.

AR-01's three-advisory investigation is closed by this bounded remediation and explicit residual decision, contingent on required final CI. The residual OpenPGP record and full release dependency/license/security decision remain open. Next autonomous work is AR-02: same-origin pagination cycles and total scan lifecycle, with independent adversarial fixtures and preserved healthy pagination. Azure-dependent items remain deferred without substituting the production-containing Advisory parent.

No Azure scan, role change, reference-repository mutation, hosted maintenance dispatch or release publication occurred. The runtime had discarded the previous alignment worktree and rejected its old shell path; a fresh clone of the verified live merge and the available bash entry point restored execution. Failed environment setup commands are not counted as QA evidence.
