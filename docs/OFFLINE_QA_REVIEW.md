# Offline project QA checkpoint

Review date: 2026-10-01 (Europe/Oslo). Starting revision: `dd3cc8309b799502ba92c9d43559b776eb8b3918` on `bootstrap/core-v1`. This bounded review executes the eight agreed autonomous tasks and the development QA workflow. It is not Gate 004 PASS, a penetration test or release/security approval.

## Eight-task outcome

| Task | Implemented evidence | Remaining boundary |
|---|---|---|
| Documentation consistency | Reconciled roadmap/review links and PR #63 final CI evidence. New CI checker validates authored README/top-level docs local file links. | It checks file destinations, not all anchors, external URLs or retained third-party notices. Historical gate snapshots remain unchanged. |
| QA evidence mapping | The requirement-to-test table below indexes covered behavior and open acceptance decisions. | A covered slice is not a complete gate row. |
| Offline scope/filter checks | New discovery fixture combines two included subscriptions/RGs, overlapping exclusions, exact tag values and a separately excluded child under a selected parent; checks exact selected IDs and nearest recorded downstream scope. | No new live equivalence evidence or change to existing filter semantics. |
| Runbook validation | Actual built CLI help validates eight Cloud Assess flags; PowerShell parser checks the operator snippet without executing it. Pinned identity source confirms credential settings. | Local PowerShell 7 and code-checkpoint Windows 5.1 checks passed; final required checks must pass before merge. This is syntax/help validation, not authentication or a scan. |
| Read-oriented request checks | Authenticated HTTP fixtures now constrain Advisor metadata GET and Cost subscription-query POST, complementing existing ARG/Diagnostics and SDK discovery fixtures. Registration-required failure fixtures now enforce no automatic provider-registration POST. | Not a universal allowlist or complete adapter audit. |
| Retry policy | Production default five retries yields six attempts for exhausted 429/503 and recovery on attempt six; 403 is not retried. Server millisecond retry hints keep fixtures bounded. Existing cancellation/retry-wait/body tests remain. | Wall-clock exponential timing, whole-scan budgeting and all service-specific retry behavior remain open. |
| Sensitive-output review | Canary bearer header absent from tested SDK request/response/error logs and 403 error text; upstream error message still appears, documenting the private-log boundary. Fixed unsafe Advisor pagination and disabled redirects in the default shared HTTP transport. | Not general log anonymization. Injected transports and SDK-owned pagers have separate destination/redirect policies. Names, tags, IDs, raw service errors and SARIF remain sensitive. |
| Interrupted package publication | Missing-checksum fixture rejects before extraction. Native integration injects failure on checksum publication after ZIP creation, rejects consumption/retry overwrite and preserves the orphan ZIP. | Two-file publication remains non-atomic; this tests detection and retention, not automatic repair or fresh-OS installation. |

## Requirement-to-test map

Paths below refer to the repository root. They identify concrete evidence, not an aggregate coverage claim.

| Gate/release concern | Existing or added tests | Still open |
|---|---|---|
| Comparable semantic inputs | `internal/equivalence/compare_guard_test.go`, `projection_test.go`; bounded fuzz and selected mutation checks | Live missing scenarios, accepted normalization/uncertainty review |
| Scope and filter behavior | `internal/config/filters_test.go`; `internal/discovery/resources_test.go`, `subscriptions_test.go`, `azure_scope_test.go` | Intended-versus-resolved reporting contract, restricted identities, DV-001 live nested hierarchy |
| Failed retrieval versus valid empty data | `internal/app/optional_stage_http_test.go`, `e2e_test.go`; ARG continuation/envelope tests | Non-empty live optional stages, all adapter combinations |
| Exit codes/report persistence | `cmd/cloud-assess/process_test.go`; `internal/app/scan_test.go`, `e2e_test.go` | Successful installed CLI Azure scan on each selected platform |
| Query destinations and credential boundary | ARG/Diagnostics request contracts; `internal/advisor/metadata_http_test.go`; `internal/cost/http_contract_test.go`; `internal/azure/http_client_test.go` | SDK-owned pager/injected transport policies, whole adapter/security review |
| Cancellation/retries/lifecycle | `internal/azure/http_client_test.go`; ARG and application interruption fixtures | Whole-scan budget and service/load behavior |
| Privacy and output handling | All-renderer application fixture; `internal/reportfile/write_test.go`, `windows_acl_test.go`; bearer canary fixture | Arbitrary-directory ACLs, raw logs/evidence handling, SARIF exposure acceptance |
| Built/package integrity | `scripts/tests/built-cli.py`, `package-candidate.py`; dependency inventory safeguards | Fresh-OS validation, publication, signed provenance, licensing sign-off |
| Maintenance controls | `scripts/tests/maintenance-publication.py`, `maintenance-generation.py` | Controlled hosted workflow dispatch/token/proposal review paths |
| Documentation drift | `scripts/tests/documentation.py` in Linux/Windows CI | Semantic review remains manual; parser does not execute examples |
| Known reachable vulnerabilities | Pinned `govulncheck` in both native jobs | Module-only advisories and exact release artifact/security decision |

## Security findings and corrections

**Confirmed:** Advisor metadata accepted an arbitrary HTTPS absolute continuation. A synthetic transport showed a second request with an ARM canary bearer header to a foreign host before the fix. No real token or external network was used. Advisor now permits only HTTPS continuations on the configured host/port, rejects user information, fragments, malformed/protocol-relative URLs and detects repeated URLs. Rejection errors omit supplied URLs. Same-origin public, government-host and custom-port fixtures plus existing relative pagination pass; these fixtures do not certify sovereign deployment. Distinct endlessly generated continuations remain a whole-scan budgeting concern.

The default shared HTTP transport also refuses redirects, returning the 3xx as a retrieval error. A two-server local TLS fixture proves zero redirected requests. A caller supplying its own transport must enforce its own redirect/destination policy. These are deliberate target hardening changes for unsafe responses, not claims about pinned-source handling of adversarial responses. Healthy data projection and comparator normalizations are unchanged. See FN-017.

**Dependency advisories:** Local Go 1.26.8 symbol/package scan found zero reachable/imported-package findings. Module-only findings are `GO-2026-6355` and `GO-2026-6354` (`x/crypto/ssh`, fixed in v0.56.0) and `GO-2026-5932` (unmaintained `x/crypto/openpgp`, no fixed version). Current selected module is v0.55.0; the Linux scanner does not import those packages. Do not call this an advisory-free dependency graph. Windows receives its own native scan; upgrading reviewed dependencies/inventory remains a separate maintenance checkpoint. Re-evaluate reachability when imports, plugins, platforms or dependencies change.

**SDK implicit registration:** Final source review found `DisableRPRegistration` unset in the shared ARM options. A synthetic valid ARM-resource GET through those options received 409 `MissingSubscriptionRegistration` and triggered a registration POST before returning the injected write rejection (two requests, one write attempt). This is a configuration-path reproducer, not evidence that historical subscription/MG scans changed Azure resources. Shared options now explicitly disable registration; scope discovery also forces it off on a copied options structure, preserving caller configuration. The regression now returns the original 409 with one GET and no POST; three real SDK scope-pager fixtures preserve the original error and caller options. This closes the identified implicit-action boundary, not all SDK-owned pager destination policies. See FN-019.

## Validation record

Local focused tests and full `go test -race -count=1 ./...` passed with pinned Go 1.26.8 (`-buildvcs=false` only for this workspace's known stamping context). Local authored-file links and real-CLI help/PowerShell parsing passed: 89 destinations, eight flags, one snippet on the completed documentation revision. Unit packaging invocation completed fourteen cases with the native integration explicitly skipped; that invocation is not package acceptance. The isolated clean Linux clone at local checkpoint `eadd9a29b55be73b760e85010a509ce53a547137` subsequently executed all fourteen native cases with no skips, including interrupted publication and extracted CLI checks. Both native CI jobs must execute all fourteen cases, plus the project provenance/module/coverage/fuzz/mutation/vet/vulnerability workflow, on the exact final PR head before merge. Strict stamped packaging uses an isolated clean clone locally and native CI; product provenance checks are not relaxed.

Required Linux/Windows code-checkpoint run `36792517368` passed PR #64 head `b447ec46c3eee232259e7df3621227f840106451`. Native logs on both hosts show all fourteen package cases with no integration skip, 89 file destinations, eight flags and one parsed snippet. Both native vulnerability scans reported zero reachable/imported-package findings and three module-only findings. Linux aggregate statement coverage was 79.4%, above the 75% regression floor; this is not behavioral completeness. That run predates the final SDK registration correction. Focused and full local race tests plus vet were rerun after the correction; the exact updated final head must pass both jobs before merge.



Final implementation checkpoint: PR #64 head `a531237eae3bc8d4bec80bcfb6518698f1d61b9d`, including SDK registration suppression, passed required Linux/Windows run `36793880867`. Inspected both native logs: all fourteen package cases with no skips, 89 file destinations/eight flags/one snippet, and zero reachable/imported-package vulnerabilities. Linux aggregate statement coverage was 79.8%. The post-correction isolated clean Linux package suite also completed all fourteen cases at local checkpoint `dc0a04d` (temporary candidate SHA-256 `33f0be8de2dbf6c1af27004b8c2d55703d9d46a3db19af1d5aa601857965b194`). No artifact was published. Final evidence-documentation revision must pass both required jobs before merge; Gate 004/release decisions remain open.

The development ledger records code-head/run evidence after verification. No Azure scan, role assignment, provider registration, production pilot or release publication occurred.

## Primary guidance

- [Pinned Azure SDK retry options](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore@v1.23.1/policy#RetryOptions)
- [Pinned ARM automatic-registration option](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azcore@v1.23.1/arm/policy#ClientOptions)
- [Go HTTP redirect policy](https://pkg.go.dev/net/http#Client)
- [OWASP SSRF prevention](https://cheatsheetseries.owasp.org/cheatsheets/Server_Side_Request_Forgery_Prevention_Cheat_Sheet.html)
- [Go vulnerability checking](https://go.dev/doc/security/vuln/)
- [GO-2026-6355](https://pkg.go.dev/vuln/GO-2026-6355), [GO-2026-6354](https://pkg.go.dev/vuln/GO-2026-6354), [GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932)
