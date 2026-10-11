# Current security remediation authority, 2026-10-11

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

# Dependency and license evidence

The [manifest](dependencies/inventory.json) records the selected Go module graph, module/download checksums, directness, Linux/Windows amd64 CLI membership, input hashes, and bundled-source revisions. The [generated notices](dependencies/NOTICES.md) preserve full discovered module-root license/notice texts, Go standard-library/toolchain LICENSE/PATENTS, project/source attribution and the pinned APRL license. Current collection includes 49 modules, nine direct, with 24 contributing packages to Linux CLI builds and 26 to Windows CLI builds. Modules without CLI membership remain listed conservatively because the selected graph also includes test/tool dependencies.

## Regeneration and review

Use Python 3 and Go **1.26.9**, with the APRL submodule initialized at the recorded pin. From the repository root:

```sh
python3 scripts/dependency-inventory.py
python3 scripts/tests/dependency-inventory.py
python3 scripts/dependency-inventory.py --check
```

Use `--go /absolute/path/to/go` if Go is outside PATH. Generation downloads exact selected versions into the Go cache and uses an isolated temporary module for those downloads. It does not upgrade dependencies, run Azure scans, or modify go.mod/go.sum. Package membership uses `go list -deps` with CGO disabled and explicit target OS/architecture. `-buildvcs=false` avoids incidental VCS stamping during package listing; explicit source pin/dirty checks remain mandatory. Generation must finish collecting all evidence before output files are written. The two output writes are not a crash-atomic transaction; rerun after an interrupted write.

Review both generated files whenever dependencies, notices, Go patch version, embeds or bundled source pins change. Commit the updated evidence with the change. CI runs safeguard tests and `--check`, rejecting missing or stale files. Do not suppress a stale result simply to make a dependency-refresh PR green: regenerate and review that proposal. Existing automatic tidy proposals still require this human review step when their dependency graph changes.

The collector rejects replaced/unversioned modules until their provenance policy is explicitly implemented. Missing, empty or symlink root license evidence is a failure. It does not infer SPDX identifiers from a few words or mark absent data as approved. A newly embedded source family or target requires explicit generator coverage and notice review. Current rule roots match the normal reference catalog; specialized APRL workloads outside that catalog are not claimed automatically executed. The SKU data derives from the pinned AZQR baseline.

## Release boundary

This inventory is evidence for review, not a formal SPDX/CycloneDX SBOM or a legal clearance. Root license discovery does not prove every package file lacks additional terms. Review package-specific exceptions, secondary embedded assets, modifications, required notices and applicable obligations before redistribution. The standard library is recorded separately from the module graph. GitHub Actions and independently downloaded QA tools are build/test infrastructure, not enumerated as CLI modules here.

For each release candidate, regenerate on the exact release tree/toolchain, confirm the actual supported platform/build flags against recorded membership, inspect binaries/packages and retain required full notices beside artifacts. Other architectures, CGO configurations and release packaging are not covered by this amd64 snapshot. Keep exact release SHA, artifact checksums, build provenance, vulnerability results and licensing/operations sign-off in the separate release gate. No release or Quality Gate 004 approval follows from this check alone.

Generated NOTICES.md intentionally preserves upstream line endings and trailing whitespace. Its path-specific Git attributes disable text conversion and whitespace diagnostics for that file only; authored files retain ordinary checks. Notice SHA-256 values and freshness checks verify retained evidence.

The module graph/download semantics follow the [official Go Modules Reference](https://go.dev/ref/mod#go-mod-download); package-list semantics follow the [Go command documentation](https://pkg.go.dev/cmd/go#hdr-List_packages_or_modules).

The [dependency advisory review](DEPENDENCY_ADVISORY_REVIEW.md) records the selected x/crypto v0.56.0 upgrade, unchanged license evidence and residual module-only OpenPGP advisory. Native reachability checks remain mandatory; absence of reachable findings does not certify an advisory-free module graph or release security.
