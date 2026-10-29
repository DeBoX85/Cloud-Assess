# Security QA crossroads, 2026-10-20

Current authority: [CURRENT_STATE](CURRENT_STATE.md). Implementation is paused for Denis's decision. This is a mandatory-gate/dependency-policy blocker, not a demonstrated exploitable incident or a severity/CVSS assessment. Do not weaken a scanner, pin its vulnerability database backward, remove Excel output, suppress findings or refresh dependencies without authorisation.

## Decision received, 2026-10-20

Denis approved the narrow approach. The earlier wait-for-decision statements are historical; see [remediation contract](EXCELIZE_SECURITY_REMEDIATION_20261020.md), [current advisory review](SECURITY_ADVISORY_REVIEW_20261020.md) and CURRENT_STATE. Implementation of this approved maintenance is underway, with own exact-head QA pending. The failed177evidence remains failed; no suppression/bypass or retail production is authorised by this update.

## Trigger and exact state

User authorised an autonomous execute/validate/verify/reevaluate batch, tasks6-8. Task6 contract candidate is draft [PR177](https://github.com/DeBoX85/Cloud-Assess/pull/177), b85655ac43d0f902ab7e3a0d32cc38c928c93e8c/tree7e874e15c073f3388a5833da1de9e5316abb37c1, sole parent verified176/1ededa501e663d888a905d6585dafff0d11e9087. All10complete changed files/tree/parent/recorded human identity were read back. Its52future literals and three exact pinned source copies are frozen. Production, go.mod/go.sum, workflows, historical captures and pins equal176. Task6 is NOT complete or verified offline: full quality39350527288 failed both native security scans. Source39350527219 passed; full23capture reconstruction/evidence review is recorded separately as it completes.

Task7 strict consumer is draft [PR178](https://github.com/DeBoX85/Cloud-Assess/pull/178), 9e15db3df8cd34c35842394e26acaee0ac5d01ef/treecb5c637f8e47b419b41e70dd56c6063e593db730, sole parent177. BOTH complete native logs in39350624037 show the intended missing NewRetailCollector/RetailPricingRequest before ANY retail production edit: Windows118393688896 completed00:59:17Z and Linux118393689193 completed00:59:25Z. This expected failed preparation is not runtime coverage/product regression/certification. The full preparation run is not a certified replacement for177. Retail production has not started; task8 has only read-only dependency research.

Last historically fully verified development head remains176/1ededa50 (product bytes175/eb554ee8). Its old scan success does not establish current advisory clearance. All drafts remain unmerged; protected bootstrap/core-v1 staysa6fee427e19feeef28691737ea9bc5026086efa0.

## Failure evidence and reproducibility

Unchanged workflow invokes pinned govulncheck v1.8.0 against the current Go vulnerability database. Both current native jobs failed at that step, not at the new integrity test: Linux118393375951 and Windows118393375552. They report:

| Advisory | Reported condition | Current/fixed version reported by scanner |
| --- | --- | --- |
| [GO-2026-6688](https://github.com/golang/vulndb/blob/master/data/osv/GO-2026-6688.json), GHSA-wp2g-vpjj-g53r | Mutually referencing array formulas can exhaust the stack | v2.11.0 / v2.11.1-0.20260911060113-ea12859e43c6 |
| [GO-2026-6684](https://github.com/golang/vulndb/blob/master/data/osv/GO-2026-6684.json), GHSA-2j4c-ffch-9f23 | Malformed encrypted OLE/CFB workbook decryption can panic | v2.11.0 / v2.11.1-0.20260915055537-22f76f9acb94 |

Both native scans report two affected vulnerabilities, four additional imported-package findings and one module-only finding; the latter are not six additional demonstrated exploits. Previous176native logs on2026-10-19 at20:15/20:16Z reported0affected/0imported/1module-only. No dependency or product change explains this new result; it is consistent with changed advisory data. Exact publication timing/cause is not independently proven. Fetched OSV records contain zero-valued published/modified fields, so those fields are not reliable disclosure timestamps.

Reproduce without any Azure access: in an isolated exact177 checkout with pinned Go1.26.9/modules, run the unchanged workflow command `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...`. Retain output/database metadata and do not blind-rerun until a relevant change. Both hosted executions already reproduce the failure; no additional network scan is needed simply to repeat it.

Local focused AND serial full-repository preflight/tests/vet/formatting/static/gate/handoff passed. Four frozen-artifact drift checks compiled, failed their named integrity assertion, restored exact bytes and passed baseline. QA-RETAIL-01: first broad local test overlapped fixture mutation; its result was excluded even though green. Full preflight was rerun serially after complete restoration and passed. No hosted checkout was affected. Retain both logs/scheduling note; this is not the cause of the hosted vulnerability failure.

## What is and is not established

The production renderer creates a new workbook and writes cells; code inspection found no production OpenFile/OpenReader/Decrypt/CalcCellValue call. Tests reopen generated workbook outputs. This narrows plausible attack paths but does not close the gate or prove non-exploitability. GO-2026-6688 lacks symbol-specific import metadata in the fetched OSV entry, so broad traces through many exported Excelize APIs must not be treated as independent exploit demonstrations. GO-2026-6684 explicitly includes File.Write among affected symbols; the renderer calls File.Write. No malicious-file exploit was executed. A self-reviewed reachability argument is insufficient authority to suppress the findings.

Primary lookup: direct pkg.go.dev/vuln and vuln.go.dev JSON opens failed; exact upstream golang/vulndb JSON and both qax-os/excelize fixing commits were successfully fetched through GitHub and retained. Empty search results do not negate the advisories. Upstream fix ea12859e43c64d498ecf839263c5f917391f3316 changes ANCHORARRAY recursion to share resolver context;22f76f9acb94b85b3eb9c4365ab4f750cebbc565 adds decryption validation/error checks. Upstream compare reports22f76is four commits ahead ofea128with0behind, so the later candidate includes both fixes.

## Concrete decision proposal

Recommended: authorise a NARROW security-maintenance exception to the no-dependency-refresh rule, on a separate draft branch based on the frozen177 contract (not the intentionally uncompilable178 tree). Evaluate exact upstream commit22f76f9acb94b85b3eb9c4365ab4f750cebbc565 / module versionv2.11.1-0.20260915055537-22f76f9acb94 first; do not select an arbitrary latest revision. The repository releases API currently lists affectedv2.11.0 as latest stable, so this proposal deliberately requires a reviewed pseudo-version rather than pretending a newer stable exists.

Proposed direct requirement change (NOT APPLIED): github.com/xuri/excelize/v2 v2.11.0 -> v2.11.1-0.20260915055537-22f76f9acb94. Read the complete upstream delta, verify module sums/provenance, and constrain any transitive changes. The fixed upstream go.mod requests mscfb1.0.8 (current root1.0.7) and x/image0.44.0 (not currently a root requirement), and updates crypto/net/text lower minima already exceeded by current root pins. Actual minimal-version-selection impact and sums have NOT been resolved or installed. Stop if it requires toolchain/unrelated refresh; return with exact diff. This is a proposal, not a verified safe upgraded candidate.

On approval: first create a bounded dependency/remediation contract and exact recovery branch; preserve all source oracles/control identities. Add appropriate independent upstream-regression/renderer/roundtrip/package/branding controls without accepting malicious remote workbooks in production. Re-run full native/source/security/licensing/module/install QA and review full logs. Require current scans to clear these findings and separately inventory imported/module-only advisories. Rehearse dependency+sum rollback and restore on the isolated branch. Complete task6 QA/evidence anew at the changed head; only then rebase by additive successor/merge-free lineage the preserved strict consumer and begin task7. No original head rewrite, protected merge, bypass, release, Azure or source/API refresh follows from this exception.

Alternative: keep pins and pause this batch pending an independently reviewed remediation/reachability decision or stable upstream release. This leaves the security gate failed and tasks6-8 incomplete. Replacing/removing Excel output, vendoring a silent patch or suppressing scanner findings are not routine authorised mitigations and are not recommended shortcuts.

Exact question for Denis: May the agent evaluate and implement this narrow Excelize security update, including only necessary reviewed transitive/sum changes, while keeping all existing tests/gates/pins and draft-only boundaries? Otherwise retain the pause and specify the preferred remediation path.

## Preservation and recovery

All prior heads,52literals, strict first runner, original captures/source pins and unrelated dirty worktrees remain preserved. No retail production API or dependency modification exists. New supplemental tests and task8 research are retained as inert evidence, not certified implementation. This checkpoint is documentation/evidence only and cannot inherit176certification. Live177final index owns exact run/archive/checkpoint references; inspect current states before retrying any operation. No uncertain remote mutation remains after verified readback.

Undoing the task6 candidate restores exact176tree; that removes the new contract but does not remediate a vulnerability in the unchanged dependency. Do not roll back protected branches or mark historic176green as current clearance. Recovery after a new chat starts at this document and CURRENT_STATE; wait for Denis's explicit decision before dependency edits or dependent retail work. No notification has been scheduled/sent.
