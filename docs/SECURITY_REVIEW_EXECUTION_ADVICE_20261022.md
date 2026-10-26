# Security execution advice and budget constraint, 22 October 2026

Denis reports that independent human evaluation will be considered but access and budget are currently unavailable. Record this as an open, deferred assurance limitation. It does not block authorized offline development, does not mean a review was commissioned, and is not a waiver or a claim of independent-human acceptance. Do not repeatedly require purchasing an audit to continue the current work.

The separate reviewer reports incomplete source/dependency setup and a missing discoverable Go executable. This development session has a working Go1.26.10Linux toolchain and pinned dependencies, but that does not establish availability in the other chat's reported Windows environment. No automatic filesystem or chat sharing is assumed. Its reported SEC-R-001 remains an observation pending evidence; no finding has been integrated merely from this setup request.

## Recommendation, not an authorization already granted

The proposed isolated Go1.26.10 setup is reasonable with the boundaries below. Denis's request here asks for risk analysis; it does not by itself grant the separate reviewer new setup permission. Its frozen handoff remains unchanged. Denis can approve the specific setup in that chat after considering this advice. No arbitrary tooling install, dependency upgrade, system change or broader execution is implied.

The official [Go downloads](https://go.dev/dl/) list go1.26.10.windows-amd64.zip with SHA2567ecd2d7260e82c59bcbd7f1077a07cee2dea607295b858df67195a1e4e09629b, checked22October2026. Select that artifact only after confirming Windows amd64; do not assume another executor has that architecture. Verify the downloaded archive before extracting into a fresh disposable directory. Use official distribution endpoints; no unverified mirror or relaxed TLS/checksum validation.

## Material risks and mitigations

| Risk | Boundary |
| --- | --- |
| Running code can access files or credentials available to its process | A separate folder is not a security sandbox. Use the existing disposable managed executor with no real credentials exposed to the test process. Inspect the package, initialization, TestMain, subprocesses and test helpers before execution |
| Unexpected outbound requests | Permit setup downloads separately. During synthetic execution block network access using available host controls and inject confined fake transports. GOPROXY=off disables module retrieval, not application networking; changing a default HTTP client alone does not block every network path |
| CPU, memory or disk exhaustion | Start with small fixed inputs, serial execution, explicit wall-clock/process/output/data limits and available hard memory isolation. GOMEMLIMIT is a soft runtime target, not a hard process memory boundary. If a required isolation boundary is unavailable, pause that test and continue safe analysis |
| Wrong source/toolchain/dependencies or workspace contamination | Verify complete frozen185commit/tree and submodule pins, use separate caches and process-local settings, GOTOOLCHAIN=local, GOWORK=off and read-only module mode. Preserve and compare go.mod/go.sum and relevant source bytes before/after. No go get/tidy/update or expectation changes |
| Misleading result | Include a healthy control, actual execution exit, workload size, environment/limits and exact code revision. A timeout/setup failure is not a reproduced vulnerability. A bounded reproduction supports only the demonstrated claim |
| Evidence loss/disclosure | Keep sanitized minimal synthetic evidence; no public raw logs or sensitive payloads. Issue189 remains DEFERRED OPEN; hashes are not durable backup |

Go documents [toolchain selection](https://go.dev/doc/toolchain), [workspace/module behavior](https://go.dev/ref/mod) and [TestMain](https://pkg.go.dev/testing). A test-name filter does not eliminate package initialization or TestMain setup. Read-only module mode protects go.mod updates; it is not a filesystem sandbox, so separately verify go.sum/source integrity and confine writable caches/temp outputs.

No system-wide install, administrator privileges, persistent PATH or go env changes, new listener, Azure/live-retail traffic, real credential discovery, shared worktree mutation, repository/ref/PR writes, protected merge or release. Tests needing any of those stay excluded. Setup downloads do not create Azure resources or authorize paid services. Existing environment restrictions still apply; do not bypass them.

Suggested approval wording for Denis to send manually: authorize only the exact verified toolchain, unchanged complete frozen source and pinned-dependency retrieval into disposable directories; permit source-inspected, bounded, listener-free synthetic tests with fake credentials and runtime network confinement; preserve every existing restriction, record actual evidence, and pause affected tests when required isolation or integrity cannot be established. This advice is not cross-chat consent or evidence that setup/execution has occurred.
