# Immutable branding profiles: development builder

Date: 2026-10-01 (Europe/Oslo). This implements the first two steps of [BRANDING_REVIEW.md](BRANDING_REVIEW.md): strict profile parsing, a supported native development builder and actual executable identity checks. AR-03 remains open. Paired custom report/data checks are now implemented in [BRANDING_REPORT_QA.md](BRANDING_REPORT_QA.md); profile-aware distribution is the next acceptance step. Existing candidate packaging still expects cloud-assess and rejects custom identities; do not bypass that validation or describe custom executables as release-approved packages.

## Build and inspect

From the repository root, with pinned Go 1.26.8 and initialized pinned submodules:

```text
go run ./tools/brand-build --profile examples/branding/example.json --output .build/example-brand
```

Linux creates `.build/example-brand/example-cloud`; Windows creates `.build/example-brand/example-cloud.exe`. The builder only supports native Linux/Windows amd64, requires Git VCS stamping and uses trimpath, readonly modules and CGO disabled. It never overwrites an existing executable. Choose a new output directory for another build. `--version` defaults to dev and accepts a bounded alphanumeric candidate label. Omit `--profile` to use the unchanged built-in identity. A dirty checkout is allowed for development builds but remains ineligible for strict candidate packaging.

Run the built executable with `branding` to print its fully resolved canonical JSON profile, or `--help` and `--version` to inspect presentation identity. These commands do not authenticate or scan. The profile is embedded at build time, validated once at process startup and immutable for that process. Malformed linked state fails before command construction, authentication or report replacement, including help/version requests. No runtime environment or scan flag changes branding.

The builder executes Go directly with explicit arguments. Public profile metadata is canonicalized and base64url encoded into one dedicated linker string; display text cannot supply shell or linker options. See Go's [-X contract](https://pkg.go.dev/cmd/link). Pinned local Go 1.26.8 linker source confirms it. The builder invokes the Go executable on PATH from the repository root; toolchain/PATH and build environment remain trusted operator inputs, not a sandbox for untrusted compilers.

## Profile contract

```json
{
  "schemaVersion": 1,
  "productName": "Example Cloud Toolkit",
  "cliName": "example-cloud",
  "reportTitle": "Example Cloud Assessment",
  "reportFilePrefix": "example_report",
  "websiteURL": "https://example.test/tool"
}
```

Only schemaVersion is required; omitted presentation fields use the original defaults. Null is not omission. Unknown or duplicate keys, wrong types, non-object/trailing JSON, unsupported schema version, invalid UTF-8 and input over 16 KiB are rejected. Keys are case-sensitive. Empty or padded display text, control/format characters and Unicode replacement characters are rejected; productName/reportTitle are limited to 256 UTF-8 bytes. The parser explicitly rejects duplicate keys rather than relying on Go JSON's [default last-value behavior](https://pkg.go.dev/encoding/json).

cliName and reportFilePrefix must match `[a-z][a-z0-9_-]{0,63}` and cannot be Windows device names CON, PRN, AUX, NUL, COM0-9 or LPT0-9 (case-insensitive). This excludes drive/UNC paths, separators, traversal, leading option syntax and trailing dot/space. websiteURL must be an absolute HTTPS URL with a host and no user information, controls or padding, limited to 2048 bytes. URLs are public metadata and are never fetched during validation. Do not store secrets in the profile. ShortName, CompanyName, SupportURL and LogoAltText remain unused reserved fields and are deliberately not accepted.

## Evidence and remaining acceptance

`internal/branding/profile_test.go` covers defaults, canonical round trip, duplicate/type/version/UTF-8/control/path/URL rejections, malformed embedded state and concurrent copy isolation. Both required native CI jobs run `scripts/tests/branding.py`: real default/custom builds (including spaces/quotes and shell metacharacters in display text), canonical identity, help/version, scan preflight, invalid-profile output prevention, existing executable preservation and a deliberately malformed linked executable. A local HTTP tripwire and isolated Azure configuration detect authentication attempts in those CLI checks. The malformed executable must preserve an existing report.

Existing application/XLSX/SARIF consumers use the immutable resolved identity. The subsequent [paired report acceptance](BRANDING_REPORT_QA.md) checks actual custom/default JSON/CSV/stdout, synthetic source provenance, all XLSX titles, SARIF metadata/results/fingerprints/automation namespace, default filenames and explicit output precedence. Next implement profile/hash-aware package producer/verifier, installation instructions and native tampering checks while preserving all existing package cases, notices, stamped source/inventory checks and privacy controls. Changing SARIF driver.name still requires separate consumer integration migration validation.
