# Cost Management continuation audit

Audit AUD005,2026-10-04. Accepted base07011b63440e69f4a5acb128e0449943f759ee9c. Test-only reproduction, production unchanged.

## Source-grounded contract and suspected gap

Official Azure/azure-rest-api-specs stable2021-10-01 costmanagement.json defines QueryProperties.nextLink as the next page URL. The actual project queryProperties decodes only rows, so a paged successful response is indistinguishable from a complete response. Current ServiceName grouping often returns few rows; no historical live truncation or incorrect billing result is asserted. The fixture supplies literal API-shaped empty/nonempty first pages with a nonempty continuation link, plus null/empty terminal controls.

The immediate bounded safety contract is fail closed on any unconsumed nonempty continuation, accept no costs from that incomplete subscription, and propagate an explicit stage error. This does not add full pagination or follow provider-supplied authenticated URLs. Full pagination would require its own method/body, same-origin/exact-scope, API-version, cycle, page/row/byte, cancellation and late-failure design. It cannot be claimed implemented by this guard.

Coordinator Cost is noncritical: a returned error marks that stage failed and overall assessment partial while retaining other stages. Application reports partial with its existing incomplete exit behavior. A warning-free first page must not masquerade as complete previous-month costs.

## Evidence and next actions

Require compiling named failures on unchanged production with healthy terminal controls passing; formatting/compile failures do not confirm behavior. If confirmed, add NextLink to the response contract and explicitly reject nonempty continuation before returning accepted records. Add coordinator integration proof, then full fresh Linux/Windows/source gates and exact-head review. Column ordering and missing envelope are separate questions, not bundled findings yet.

Official specification: https://github.com/Azure/azure-rest-api-specs/blob/main/specification/cost-management/resource-manager/Microsoft.CostManagement/CostManagement/stable/2021-10-01/costmanagement.json . Source snapshot read during audit; API version matches project2021-10-01. No Azure request, accepted merge, full audit/release/live closure.
