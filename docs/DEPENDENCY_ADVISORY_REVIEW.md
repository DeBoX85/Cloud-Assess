# Current advisory authority, 2026-10-21

The approved product continuation is draft [PR185](https://github.com/DeBoX85/Cloud-Assess/pull/185), exact 566cc4d8b78d168e5216245a5a386f742a3a2f85, using Go **1.26.10**. Its own quality39594186699/source39594186753 passed; both native scans reported zero reachable findings, with 3 imported-not-called and 1 module-only advisory still OPEN. The Excelize prerequisite is complete offline as indexed in [CURRENT_STATE](CURRENT_STATE.md). Neither outcome accepts protected code or clears independent/live/operational/security/license obligations.

This routing supersedes the dated pending/wait instructions below. Those sections preserve historical advisory evidence, not the active toolchain or a request to replay remediation. Verify live revision-specific PR outcomes before use. No additional dependency or source-pin refresh is authorised.

---

# Historical advisory checkpoint, 2026-10-20

Read [CURRENT_STATE](CURRENT_STATE.md), [current advisory review](SECURITY_ADVISORY_REVIEW_20261020.md) and the live maintenance successor linked from PR177. The narrow Excelize update is approved; own final QA remains pending. Older advisory/acceptance states below are historical snapshots, not current clearance. Behavioral requirements remain active.

---

# Historical advisory route, 2026-10-11

For the current Go1.26.9/x-net0.60.0/x-crypto0.57.0 graph, read
[SECURITY_ADVISORY_REVIEW_20261011.md](SECURITY_ADVISORY_REVIEW_20261011.md)
and [live PR133](https://github.com/DeBoX85/Cloud-Assess/pull/133) for exact
final-head native/source/review/accepted-push outcomes. The v0.55-to0.56 SSH
assessment and its local/compiler/scanner results below are retained historical
evidence for the2026-10-01 baseline, not current remediation clearance. Existing
unsafe OpenPGP remains module-only/no known fix and unimported by either CLI.
No advisory-free graph, incident, waiver or release approval follows.

---

# Historical security remediation checkpoint, 2026-10-11

Read [SESSION_HANDOVER.md](SESSION_HANDOVER.md), the
[security remediation contract](SECURITY_REMEDIATION_20261011.md) and
[live PR133](https://github.com/DeBoX85/Cloud-Assess/pull/133) for the latest
authored head, native/source/review, protected acceptance and distinct accepted-push
outcomes. PR133 preserves the PR132 handover and supplies its bounded security
correction. PR132 is not independently accepted; do not replay or merge its failed
candidate. Accepted implementation was PR131 b653a3ab at remediation preparation;
verify the live branch before treating a proposal as accepted.

The correction proposes Go1.26.9 and x/net0.60.0 with only independently reviewed
minimum graph requirements. Isolated preparation scans do not replace final native
checks. Prior Go1.26.8/current failed/pending/zero-findings statements below are
historical observations, superseded for current resume instructions by this entry.
Contractual behavior, source pins, failures and historical evidence remain preserved.
Region feature work waits for verified security acceptance. Laptop/Azure restoration,
live scope, release/Gate004/load/fresh OS and independent-person approval remain
unestablished. No Azure writes or scope substitution are authorized.

---

# Historical security gate checkpoint, 2026-10-09

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
