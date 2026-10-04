# Cost Management continuation audit

Current status: correction VERIFIED OFFLINE and accepted through combined PR114, mergec33fb2eee93a2e5730f72612b351d0b838737d03. Original individual proposal is closed/superseded; compiling failed controls and corrected proof below are retained as historical evidence. Final combined native37188397367/source37188397374 and matching tree/preview verified. Distinct accepted-push/final disposition evidence is indexed in PR114 and the disposition PR. No historical live incidence, independent-person review or release approval.

Audit AUD005,2026-10-04. Accepted base07011b63440e69f4a5acb128e0449943f759ee9c. Reproduction confirmed; fail-closed correction proposed, unaccepted.

## Source-grounded contract and suspected gap

Official Azure/azure-rest-api-specs stable2021-10-01 costmanagement.json defines QueryProperties.nextLink as the next page URL. The actual project queryProperties decodes only rows, so a paged successful response is indistinguishable from a complete response. Current ServiceName grouping often returns few rows; no historical live truncation or incorrect billing result is asserted. The fixture supplies literal API-shaped empty/nonempty first pages with a nonempty continuation link, plus null/empty terminal controls.

The immediate bounded safety contract is fail closed on any unconsumed nonempty continuation, accept no costs from that incomplete subscription, and propagate an explicit stage error. This does not add full pagination or follow provider-supplied authenticated URLs. Full pagination would require its own method/body, same-origin/exact-scope, API-version, cycle, page/row/byte, cancellation and late-failure design. It cannot be claimed implemented by this guard.

Coordinator Cost is noncritical: a returned error marks that stage failed and overall assessment partial while retaining other stages. Application reports partial with its existing incomplete exit behavior. A warning-free first page must not masquerade as complete previous-month costs.

## Evidence and next actions

Require compiling named failures on unchanged production with healthy terminal controls passing; formatting/compile failures do not confirm behavior. If confirmed, add NextLink to the response contract and explicitly reject nonempty continuation before returning accepted records. Add coordinator integration proof, then full fresh Linux/Windows/source gates and exact-head review. Column ordering and missing envelope are separate questions, not bundled findings yet.

Official specification: https://github.com/Azure/azure-rest-api-specs/blob/main/specification/cost-management/resource-manager/Microsoft.CostManagement/CostManagement/stable/2021-10-01/costmanagement.json . Source snapshot read during audit; API version matches project2021-10-01. No Azure request, accepted merge, full audit/release/live closure.

## Confirmed reproduction and proposed correction

Test-only78dd85c9cb9f50656c605226291d027e1bce06eb native37184679263 failed the two named empty/nonempty paged assertions on Linux111384036588 and Windows111384036492. Both compiled; formatting passed; returned nil error and no warnings despite literal nonempty nextLink. Null/empty terminal controls did not fail. This demonstrates the synthetic adapter path, not historical Azure incidence.

Add NextLink to queryProperties and reject any nonempty continuation before accepting rows, without echoing its URL or performing a continuation request. Existing terminal-page mapping, valid204, queries/period/two-worker ceiling and pins remain unchanged. Add actual Cost.Scanner through Coordinator test: stage failed with visible continuation reason, zero accepted incomplete subscription costs, overall partial, healthy Advisor retained and Graph completed. Existing application partial-result/report/exit fixtures cover the subsequent surface.

Full authenticated bounded pagination remains unimplemented and must not be advertised. The intentional behavior change is that an API-provided partial page now produces a visible incomplete assessment instead of warning-free success. Column/envelope questions stay separate. Require complete corrected-head native/source proof before acceptance.
