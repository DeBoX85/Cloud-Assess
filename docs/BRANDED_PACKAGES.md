# Profile-aware development candidate packages

Date: 2026-10-02 (Europe/Oslo). This completes the implementation of step 4 in [BRANDING_REVIEW.md](BRANDING_REVIEW.md), pending final exact-head native CI acceptance. It extends the [validated profile builder](BRANDING_PROFILES.md) and [paired report acceptance](BRANDING_REPORT_QA.md). Development candidates are not release-approved distributions, and no Azure scan or fresh-OS installation is claimed.

## Produce and inspect

From a clean checkout with pinned Go 1.26.8 and initialized pinned submodules:

```text
go run ./tools/brand-build --profile examples/branding/example.json --output .build/example-brand
python scripts/package-candidate.py --binary .build/example-brand/example-cloud
```

On Windows, the binary path is `.build/example-brand/example-cloud.exe`. The builder's version defaults to dev; pass the same explicit --version to builder and packager when changing it. Choose new output directories rather than overwriting existing candidates. The producer reads the executable's fully resolved canonical profile through its offline `branding` command, validates it, and verifies its version, clean stamped source revision, native target/CGO/trimpath and reviewed compiled dependencies before producing files. The executable is trusted build input: producing a package executes its branding/version commands. This is not a sandbox for arbitrary downloaded binaries.

The safe cliName determines the executable and archive root, for example example-cloud-dev-linux-amd64.zip or example-cloud-dev-windows-amd64.zip. The generated INSTALL.md uses that name throughout its commands and checksum examples. Original LICENSE/NOTICE/third-party/dependency notice bytes are retained. The package includes BRANDING_PROFILE.json with canonical Go-compatible JSON bytes and PACKAGE_MANIFEST.json with brandingSHA256. The manifest is now schemaVersion 2; version-1 development candidate archives require their historical verifier and are deliberately not interpreted as this new contract. This does not modify canonical assessment/report schema versions.

The Python resolved-profile validator requires exactly six fields, including schemaVersion, rejects duplicate keys, malformed/noncanonical encoding and unsafe filename components, and mirrors the build profile's public text/URL boundaries. It reproduces Go JSON HTML/Unicode separator escaping for canonical hash equality. Native custom checks include quoted Unicode display text and URL ampersands to test this across the actual Go/Python boundary. It does not accept partial operator input; the Go builder resolves omitted defaults first.

## Extraction and trust boundary

verify_extract validates archive/size/path/member/mode/checksum constraints, strict manifest JSON, canonical profile/hash and profile-derived root/executable layout before creating destination paths. It rejects unsafe/duplicate profile keys and mismatched names even when test fixtures recompute outer/payload hashes. It never executes the archived executable. Installed executable identity is checked separately against the packaged expected profile in controlled native QA. A supplied checksum/manifest is integrity evidence, not publisher authentication: an attacker who replaces the entire archive and all untrusted hashes can forge them. Keep trusted independent checksums, signed provenance/release review as separate outstanding requirements. See Python's [ZIP extraction guidance](https://docs.python.org/3/library/zipfile.html) and [JSON duplicate-key behavior](https://docs.python.org/3/library/json.html).

The existing fourteen package cases remain, supplemented by five resolved-profile/hash/layout/duplicate/encoding controls. Both native jobs run all nineteen cases with a real default binary and again with a real custom binary, without skipping the installed case. Isolated installs run the existing offline CLI preflight/report-preservation/HTTP-tripwire/compiled-inventory checks, now also comparing the installed `branding` command to BRANDING_PROFILE.json. Repeated packaging must be byte-reproducible, notices unchanged, wrong version rejected, existing files preserved and publication races/interrupted checksum bundles retained and rejected as before. No automatic PATH, credential, service or Azure configuration changes occur.

## QA checkpoint and next work

Current-branch QA first verified merged PR #74 and its final native evidence, rebuilt the pinned ordinary clone after a scratch runtime refresh, and ran full race/vet, documentation links, module consistency and inventory freshness. No new production Go defect was confirmed. Current branding status wording was reconciled where append-only checkpoints had left older future-work descriptions. See FN-030 for environment/preparation limits.

Final local/native results and exact head/tree/merge are retained in the PR and indexed in the next ledger checkpoint. AR-03 can be marked accepted for this bounded immutable five-field development workflow only after both required native jobs pass. Runtime per-scan profiles, logos/themes and existing SARIF consuming-service migration remain outside this acceptance. Release/Gate 004, Azure-dependent evidence and full AZQR feature parity remain open.

Next autonomous task after acceptance: characterize the pinned source's external YAML/KQL execution and missing CLI/plugin features into concrete parity slices (AR-04), preserving deferred live-validation items. No laptop/Azure input is needed for source characterization.

A related QA finding corrected Windows native-command failure propagation: each command in the branding/report and default/custom-package blocks now checks LASTEXITCODE immediately. Both jobs execute the actual Windows snippets with injected first/second failures and healthy success (six cases); no prior masked historical failure is asserted.
