# Go security gate blocker, 2026-10-21

Decision received2026-10-21: Denis approved the narrow Go1.26.9to1.26.10update with full QA and rollback. Earlier await-decision statements below are historical within this exception. Follow [the active plan](REVIEW_PLAN_20261021.md); implementation and own QA remain pending.

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Await Denis's decision on a narrow Go toolchain security exception. This is a mandatory QA and pinned-toolchain boundary, not an established exploit or severity assessment. Task7 formal closure remains open; task8 remains provisional.

## Verified recovery and failure

Live draft PR184 is d8fd8ce207939400a5e5e7c3b114bdcda0fbb8f1, tree d815098c8a3ebc1776a893f89a28785c72d5e38b, sole parent product PR183/458d5c5296c321c8e5e961f18eb5422d43c2eb18. Fetched Git confirms author AND committer Denis Bogunic <134433647+DeBoX85@users.noreply.github.com>. All55 changed files were compared byte-for-byte with fetched Git blobs. Preview0b034c74d0926fd85f34f99ca552a5924721c2b8 has the same tree and ordered protected/candidate parents; its automated GitHub committer is expected and distinct from the human candidate identity. Recovered manifest35payloads/270383bytes passes length/SHA256 checks. Static preflight passes unchanged workflow/gate/handoff checks. This does not close full evidence review or establish full native QA.

Own source39589963002 reports SUCCESS. Own quality39589962815 is not successful: Linux native119188168541 and Windows native119188168488 both completed FAILURE at unchanged govulncheck v1.8.0. Other workers were still in progress when inspected; final whole-run inventory remains to be checked. Native steps following the scan were skipped, so no full622/617 conservation pass is claimed for184. No rerun was requested.

Both full decoded native logs and the upstream advisory/blob identity are retained in [hosted evidence](evidence/20261021/go-security-blocker/hosted-evidence.json). Decoding the JSON content strings to UTF-8 reproduces the original returned bytes, including Windows CRLF. Windows42325bytes/SHA256482abe83afafdee02f961dfa796841d5498e1b633ba098c9ddf3cffe165ad78b; Linux64796bytes/SHA25612aea243a022f0496df74c765732937871bb5db7fbeeb5407efd5a82f62e075a. Each includes checkout identity, the failure and cleanup footer. No truncated prefix is accepted.

## Primary evidence and limits

Both scans report GO-2026-6809, compromised intermediate-certificate distrust in crypto/x509, found in Go1.26.9 and fixed in Go1.26.10. The reported trace is internal/azure/http_client.go:187 to x509.Certificate.Verify. They separately report3imported-not-called and1module-only findings. Existing183scan success remains a historical result and cannot certify184against current advisory data.

The [upstream Go advisory](https://github.com/golang/vulndb/blob/master/data/osv/GO-2026-6809.json), fetched through GitHub at blob73ca6db7d76ba5134c46b564416c2cfc5face7d0, identifies CVE-2026-111943, Certificate.Verify and fixed versions1.26.10/1.27.3. It explains that most platforms use system roots, but narrower cases may retain distrust gaps; the fix adds a denylist. No actual exploit or affected deployed environment was tested. Its zero-valued publication/modification fields cannot establish advisory timing. [Go release history](https://go.dev/doc/devel/release#go1.26.10) independently lists1.26.10released2026-10-19with crypto/x509 security fixes. Direct pkg.go.dev/vuln and vuln.go.dev retrieval failed; the upstream record was successfully fetched instead. Microsoft Learn was consulted for Azure/Go certificate context, but it is not authority for this specific Go advisory.

## Concrete proposed exception

Recommend authorising only Go1.26.9 to1.26.10 across the active toolchain declaration, reproducible bootstrap/checksum metadata, generated CI inputs/workflows and relevant operator documentation. Audit actual pin consumers before edits; preserve frozen historical source captures, original baseline workflow, evidence and expected controls. Rebuild affected metadata through the existing generators, with any required generator-input change explicitly reviewed. No SDK/module/action/AZQR/APRL refresh, scanner/database suppression, TLS verification workaround, protected merge or release.

The current user's instruction preserves pins. The prior [Excelize exception](EXCELIZE_SECURITY_REMEDIATION_20261020.md) explicitly keeps Go1.26.9 and excludes broad toolchain refresh; it does not establish permission for this separate pin change. Therefore the fix is identified but implementation awaits a new narrow decision. Repeating the same scan without a relevant change is not mitigation.

If approved, start a material security successor from exact product183, carry current handover explicitly, retain184as failed historical evidence, obtain official1.26.10checksums/provenance, review its security delta, and run full exact-head native/source QA with unchanged controls/frozen inputs and rollback/restoration. Do not use this documentation/evidence tree as the product coding base. After successful remediation, reconcile both task7archives and update PR183formal evidence closure before task8freeze/task9consumer. Any additional reachable finding or unrelated required change reopens the decision boundary.

This checkpoint contains documentation and retained evidence only. Own full native/source QA for a changed checkpoint is not claimed. All old draft heads, dirty worktrees, protecteda6fee427, AZQR8e4f0577 and APRL60eaddda remain preserved. Independent/live/operational/security/license/acceptance obligations remain open. No Azure/live retail/fixtures/RBAC/subagents/merge/bypass/release or notification occurred.
