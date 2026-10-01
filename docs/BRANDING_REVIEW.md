# Branding customization review and implementation boundary

Date: 2026-10-01 (Europe/Oslo). Baseline: merged PR #71, `be95679ddea30fe164e67731f9c8535f3ff3368b`. This is the completed AR-03 design review, not implemented customization or release approval. Production defaults and assessment behavior are unchanged. No Azure/laptop access is needed to implement the next offline batch.

## Current propagation, verified from code

| Field or identity | Current consumer | Boundary |
|---|---|---|
| ProductName | Cobra root long help | Source default only |
| CLIName | Cobra root usage/version; SARIF tool.driver.name | Source default only; package naming is separate |
| ReportTitle | Cobra short help; every rendered XLSX sheet title | Source default only |
| ReportFilePrefix | Application timestamped default output basename | Explicit --output-name already overrides the path; custom default prefix is not validated today |
| WebsiteURL | SARIF tool.driver.informationUri | Source default only; not fetched by renderer |
| ShortName, CompanyName, SupportURL, LogoAltText | No executable consumers found | Reserved placeholders; accepting configuration for these would silently imply effects they do not have |
| JSON/stdout, CSV tables | Canonical result/table projection | No product branding field; retain neutral assessment data/schema |
| SARIF automation ID and finding fingerprint namespace | Fixed cloud-assess/ and cloudAssessFinding/v1 identifiers | Treat as stable assessment integration identity, separately from display/tool name |
| Candidate executable, ZIP root and expected version output | Package producer/verifier and installed CLI tests | Hard-coded cloud-assess/cloud-assess.exe; source-brand edits alone can conflict with validation |
| Go module/import path, upstream source labels, links and legal notices | Build identity, recommendation provenance and distribution evidence | Retain original technical/legal attribution; customer presentation branding must not rewrite them |

Inspection covered internal/branding, CLI command construction, application output naming/render calls, XLSX sheet rendering, SARIF driver/identity generation, neutral JSON/CSV/table/result/redaction layers, package producer/verifier and actual CLI/package tests. This is a bounded propagation audit, not proof of working customization. There is currently no runtime profile flag, supported custom builder, custom package manifest or end-to-end custom-brand test. A centralized struct is not completion of the user's adjustment requirement.

## Recommended implementation

Use a validated build profile first. The operator edits one small profile and invokes a supported builder; no Go source edits. Rebuild is acceptable for this proposed workflow and gives each distributed executable a fixed identity. Runtime per-scan branding can be considered separately if required; it is not assumed to be a current user requirement. Avoid a mutable global/environment value repeatedly read while rendering, which could mix identities across outputs or concurrent scans.

The proposed first profile supports the five active fields above, with defaults for omitted fields. Use a versioned JSON object with explicit keys for productName, cliName, reportTitle, reportFilePrefix and websiteURL. Reject unknown/duplicate keys, wrong types, invalid UTF-8, control characters and excessive input before building. Do not claim to implement the four unused placeholders, image embedding, arbitrary themes or website/marketing generation. Proposed parsing limits and examples are implementation design, not a presently available configuration format.

The builder must pass arguments directly to Go, without invoking a shell or accepting arbitrary compiler/linker arguments from profile data. An encoded canonical profile can populate one dedicated string variable through Go's supported [-X linker option](https://pkg.go.dev/cmd/link); pinned local Go 1.26.8 src/cmd/link/doc.go confirms that contract. Decode/validate once and fail clearly on malformed embedded state rather than silently falling back. Keep caller-independent immutable values and unchanged default profile behavior. The exact encoder/loader must be covered by built-binary tests, not only unit tests that substitute a Branding struct.

Filename fields must be safe single components on both supported platforms: reject separators, dot/traversal forms, drive/UNC syntax, controls, Windows reserved device names, trailing dot/space and overlong labels. Derive executable/package names from a restricted ASCII slug, independently of a Unicode display name. Validate HTTPS information URLs without user information; do not fetch logos/URLs or introduce network work during branding validation. No profile value may alter scope, credentials, endpoints, rule loading, stages, filtering, redaction or report permissions. Profiles are public presentation metadata, not a place for secrets.

Canonical JSON/stdout/CSV content remains brand-neutral. XLSX title, SARIF driver metadata, CLI help/version, default output basename and package naming must share the same resolved profile. Preserve SARIF rule IDs, findings, locations, fingerprints and automation namespace to preserve emitted record identities. Changing SARIF driver.name can still change the tool identity in a consuming service; unchanged fingerprints do not guarantee cross-tool alert continuity. GitHub [documents the tool-name field](https://docs.github.com/en/code-security/reference/code-scanning/sarif-files/sarif-support) separately from fingerprints. Migration of an existing consumer integration needs separate validation. Explicit output paths remain operator-controlled and retain existing privacy/replacement checks. Keep upstream recommendation labels, technical module identity and license/notice bytes intact.

Packaging must understand the resolved profile instead of bypassing its current safeguards. Record the canonical profile and hash in candidate evidence; validate executable version/help/identity against it, derive archive/executable names safely, generate installation instructions matching the artifact, and reject mismatched/tampered profile or package contents. Preserve strict clean checkout, stamped source revision, reviewed dependency inventory, payload allowlist, private report handling, checksums and isolated installation tests. Do not mark a branded executable distributable until this path works on Linux and Windows.

## Next offline batch and acceptance checks

1. Implement profile model/parser, safe component/URL validation, deterministic canonical encoding and unchanged defaults; reject unsupported/duplicate keys before build/output replacement.
2. Implement supported builder plus one-time embedded-profile loading. Build both default and synthetic custom executables, including names with spaces in display text; verify help/version and invalid-state preflight with no authentication requests.
3. Exercise custom profile through actual application/renderers: all XLSX titles, SARIF tool metadata and default report paths. Independently compare canonical JSON/CSV/stdout data, recommendation provenance and SARIF result/fingerprint/automation identities with the default run. Verify explicit output paths still win. Test profile isolation/concurrent rendering and negative path traversal/control/URL cases.
4. Integrate profile-aware candidate producer/verifier, manifest/hash and installation instructions. Native default and custom packages must pass isolated extraction, checksums, notices and installed executable checks, with tampering/unsafe-name negative controls. Keep existing fourteen default package cases and both required CI jobs.
5. Document a real supported command only after it exists; record exact tested head/tree/native jobs and deliberate boundaries. AR-03 acceptance requires coherent CLI/report/package propagation; a partial builder or passing default CI is not enough.

No input from Denis is required to implement this default-preserving proposal. Any later demand for runtime per-scan branding, actual logo embedding, or replacement SARIF integration identities is a distinct product choice to record before extending the scope. Do not create another numbered quality gate solely for branding; add observable default/custom/negative checks to existing quality and windows-validation, retaining release approval separately.

## Validation and resume

Focused existing CLI/application/XLSX/SARIF race tests and authored-document checks passed before publication (139 destinations); final exact-head native quality/windows-validation remain blocking. Their broader QA does not certify custom branding that is not yet implemented. Final CI/head/tree/merge evidence stays in the PR and is indexed at the next ledger checkpoint. Readiness estimates, AR-02 metadata/load/live limits, AR-04 feature parity, DV-001 and Gate 004/release status are unchanged. Resume with the concrete five-step offline branding implementation batch above.

## QA execution follow-up

The first PR #72 native quality run failed comparator duration-bounded fuzzing at its shutdown boundary, without a property assertion or saved failing input. FN-010 records exact failed head/run/job and retained artifact; the failure is not counted as passing. The existing decoder execution-count mitigation is applied to the comparator target only, preserving its assertions and two-worker parallelism with 100000 executions and a separate 60-second timeout. This is a QA workflow correction discovered during the review, not a branding or comparator production change. Both native jobs must pass the updated final head before merge.

Local comparator count-bounded runs passed exactly 100000 executions twice (6.688s, then 7.389s with the retained 49-input baseline cache); focused comparator race tests also passed. This verifies the mitigation locally, not the root cause or final native result.

## Implementation checkpoint

The first two steps are now implemented in [BRANDING_PROFILES.md](BRANDING_PROFILES.md): validated immutable profiles, native development builder and actual default/custom executable checks. The historical audit above describes the PR #71 baseline. Runtime profile inspection now exists via the `branding` command. Custom report/data acceptance and profile-aware packaging remain the next work; AR-03 and release approval are still open.
