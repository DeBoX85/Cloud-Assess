# Offline rules inspection

Date: 2026-10-02 (Europe/Oslo). Implementation baseline: merged PR #76, `fd60dd635acb7c3156a264640810a4bd8d92b59d`.

`cloud-assess rules` prints a Markdown table. `cloud-assess rules --json` or `cloud-assess rules -j` prints a JSON array to stdout. The command loads the bundled pinned catalog and Diagnostics definitions without credentials, scope discovery, Azure queries or report creation. Use the custom immutable build's CLI name when branding changes it. The JSON flag is local to this command; scan/report flag scope is unchanged.

```bash
cloud-assess rules
cloud-assess rules --json
```

## Selection and output

The command follows the pinned source's offline inspection selection, not a raw filesystem catalog dump: visit all registered scanner resource types in scanner order; add eligible embedded Graph recommendations and then Diagnostics definitions; replace duplicate IDs globally; sort by recommendation ID. Existing target embedded-rule disabled/marker guards apply. Rules without executable query text can still appear as inspection metadata, matching the source; appearance does not prove a query executes or a resource is present. Specialized APRL workloads and external YAML plugins are not promoted into this command.

JSON has the source's six fields: category, impact, learnMoreUrl, recommendation, recommendationId and resourceType. Source provenance remains in the underlying catalog; this projection follows the source and does not add a provenance field. At the current pins the command emits 380 rows, with byte-identical JSON to the pinned reference executable. The Markdown heading counts all 107 embedded catalog resource types, including types outside the selected executable rows. It does not describe Azure estate coverage or a scan's recommendation count.

Markdown presentation deliberately uses proper table delimiters, escapes pipes/HTML/newlines in cells and displays the guidance URI as text. It is not byte-identical to the source's Markdown formatting. Missing guidance yields an empty string rather than indexing an absent first link. JSON retains original text and URLs. The existing target case-insensitive query-marker predicate is reused; source marker matching is case-sensitive. The pinned output matches exactly, but hypothetical uppercase-marker extensions are not certified as source-identical. External definitions are not loaded here.

Unexpected arguments, scope/report flags and malformed booleans fail through CLI validation. Writer errors propagate to the command and process exit 1. Successful inspection returns 0. This command creates no output file and makes no claim about authentication, subscription visibility or live scan success. Redirected stdout is a shell operation owned by the operator, separate from secure report staging.

## Independent reference evidence and regression checks

The existing read-only reference checkout at AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc` lacked populated APRL files and its first build failed the embed directive. A separate temporary source clone used that exact AZQR revision and APRL `60eaddda76541f6adbc1c5ffa686829807e55e29`, copied from the already pinned target submodule. Both original checkouts were unchanged. The temporary source was clean and compiled with Go 1.26.8, readonly modules and strict VCS stamping. The source command `rules --json` produced the independent fixture, not target code or a rewritten source renderer.

Fixture: `cmd/cloud-assess/testdata/rules-reference.json`, 129537 bytes, SHA-256 `77d7574232a32062d87dc32f6f9fc216f6f509921cc4d75621aa3a0c5279ed08`. Git attributes preserve the captured bytes without platform newline conversion. It contains public pinned recommendation metadata only. Existing [notices](../NOTICE.md) and [third-party license texts](../THIRD_PARTY_LICENSES.md) apply. Keep this independent capture when changing target selection; only regenerate under reviewed source-pin/provenance changes.

- Command tests compare complete JSON bytes for long/short/explicit-boolean flags, verify the pinned Markdown heading and 380 rows, reject invalid arguments and propagate output errors.
- Literal synthetic expectations cover unsupported scanner types, all disabled/development/manual predicates, empty query metadata, first guidance link, missing guidance, Diagnostics replacement and cross-type ID precedence, sorted order and independent returned state. Removing the selection guard caused the expected-row regression to fail; restored code passed focused race tests.
- The actual built CLI runs from an isolated directory with synthetic HTTP/authentication tripwires and no configured Azure identity. Tests compare all source rows, check Markdown, invalid arguments, exit codes, empty stderr on success and unchanged filesystem snapshots. These checks also execute inside existing default/custom isolated package installations on both native platforms.

Final required Linux/Windows results on the exact reviewed head, merged tree and workflow IDs are retained in the PR and indexed at the next ledger checkpoint. Deterministic offline command parity is not full AZQR feature completeness or Gate 004/release acceptance. Next implement scanner-specific subcommands, then bounded external YAML Graph integration and internal plugin table/health execution as ordered in [FEATURE_PARITY_CHARACTERIZATION.md](FEATURE_PARITY_CHARACTERIZATION.md).
