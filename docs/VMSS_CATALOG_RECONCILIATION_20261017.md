# VMSS catalog and deployment boundary reconciliation, 2026-10-17

Status: pre-edit decision record, not a new collector or eligibility implementation. Parent is PR149 fcc153659d8e1333cf0a3c4741f2ba984890c9df; PR149 full native/source QA has completed. This decision record is published with the documentation handoff and requires that successor's own checks. The current user authorized continued development after the safeguards checkpoint. [CURRENT_STATE](CURRENT_STATE.md) remains the current resume authority. The live task PR supplies final exact-revision outcomes.

## Primary evidence and observed code

The [Microsoft VMSS FAQ](https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-faq) describes broad VM-size support for scale sets. Its [orchestration comparison](https://learn.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-orchestration-modes) likewise lists broad SKU support for both Flexible and Uniform modes. The [SKU troubleshooting guidance](https://learn.microsoft.com/en-us/azure/azure-resource-manager/troubleshooting/error-sku-not-available) separately explains that catalog visibility can coexist with a deployment failure. These statements support a catalog comparison, not a successful workload placement claim.

Immutable primary document revisions inspected: [FAQ at 40aaa64a2771acda033ca09289e1e777241f201e](https://github.com/MicrosoftDocs/azure-compute-docs/blob/40aaa64a2771acda033ca09289e1e777241f201e/articles/virtual-machine-scale-sets/virtual-machine-scale-sets-faq.yml) and [orchestration modes at 47479904bc1e7aa0c02dbfc484482e183e26c363](https://github.com/MicrosoftDocs/azure-compute-docs/blob/47479904bc1e7aa0c02dbfc484482e183e26c363/articles/virtual-machine-scale-sets/virtual-machine-scale-sets-orchestration-modes.md). These are supporting evidence links, not changes to the pinned AZQR, APRL or retained source captures.

[VMSS List SKUs](https://learn.microsoft.com/en-us/rest/api/compute/virtual-machine-scale-sets/list-skus?view=rest-compute-2026-03-02) requires a named existing scale set and resource group. It cannot silently become a subscription/target-region catalog or justify a resource creation. The viewed page actually displays api-version 2026-03-01; the view selector is not an API-version authority. Preserve the existing Resource SKUs 2021-07-01 request and its frozen primary schema; no SDK/API fallback or live calls are introduced here.

Pinned source compute.go (Git blob ee9300e5f8652591fe43121d0973321742ee70e3) uses vmSizes for VM and VMSS. Its synthetic observations do not establish deployment validity. The current target ComputeVMSKUCollector can collect VMSS-bound node-catalog diagnostics but Evidence deliberately returns unsupported for VMSS. The current assessment.Resource contains SKU name/tier/family/capacity and kind; it does not contain a complete orchestration, placement, image/generation, disk/network or capacity-reservation deployment specification. CalculateAvailability consumes catalog/provider evidence and inventory counts; its percentage is not a deployment certificate, and unsupported evidence is excluded with a warning.

Primary file identities were checked against GitHub directory listings and exact Git object bytes, rather than treating rendered connector text as byte-preserving:

| Document | Bytes | Git blob | SHA256 |
| --- | ---: | --- | --- |
| FAQ | 50812 | 574c740a163a68473a9b76213f6a5c5963815694 | 8b899f688d34c32dcdde8385caca485d82a6f0c577d555bb71507ae6c8ee76a7 |
| Orchestration modes | 20859 | 6320cb6b0b86c27fdf30a09d30604228da34dbcc | 275708d14b346255cda0e7e247cd4f98379f9c024a5df00eff50efde26778284 |

Initial connector text produced different byte counts for these Unicode documents. It was not credited as immutable byte evidence; exact Git bytes resolved the mismatch and matched both original blob identities. No third-party document is copied into the product or notice inventory here.

## Reconciled decision

| Question | Decision |
| --- | --- |
| Can VM sizes inform a VMSS catalog comparison? | Yes, primary guidance supports that limited interpretation; source parity can be pursued under an explicitly named catalog contract |
| Does the catalog establish eligibility for every intended deployment? | No; specific workload constraints and allocation/capacity are distinct and are not provided by current inventory |
| Should the current unsupported projection be removed immediately? | No; changing it requires a new frozen catalog contract and independent literal expectations, with explicit limitation reporting |
| Is a resource-specific VMSS API an acceptable silent substitute? | No; its resource identity and scope differ and it would not answer the requested target-region question |
| Does preserving unsupported scoring mean Azure lacks VMSS support? | No; it is an implementation/evidence boundary, not a service-support assertion |

This clarifies the earlier VM/VMSS contract rather than rewriting its frozen expectations. Broad support evidence resolves the general catalog relationship. Deployment eligibility remains unproven. A future catalog adapter must label its evidence and limitations so a catalog pass cannot be reported as workload deployability. Until that contract and implementation pass their own exact-head gates, retain the existing unsupported/empty VMSS projection and all its compiling controls.

## Next bounded development contract

Prepare a separate pre-edit VMSS catalog contract using the existing bounded Resource SKUs collector as the diagnostic source. Define an explicit catalog-only projection and subscription/region/type binding; do not mutate a caller-owned collection or silently relabel VM evidence as VMSS. Keep true deployment eligibility out of that contract. Before enabling score contributions, define a visible limitation/warning at the existing internal evidence/coordinator boundary; public schema changes remain outside current scope. If safe limitation reporting requires a new public schema or unsupported inference, retain unsupported scoring and proceed with another bounded SKU family instead.

Freeze independent literal cases before production changes:

| Cases | Required assertion |
| --- | --- |
| Complete explicit unrestricted catalog; regional and zone restrictions | Exact catalog state/map/zones; both restriction reasons and precedence retained |
| Missing restriction/capability/location metadata; unknown required name; absent name | Missing knowledge stays unknown; explicit absence remains catalog unavailable; no deployment claim |
| Partial/unknown collection, operational status and malformed response | No complete scoring map; diagnostic prefix never becomes complete evidence |
| Wrong selected subscription/target/type; forged complete collection | Sanitized rejection before projection; no implicit rebinding |
| Empty requirements; mixed known/unknown unrequired rows | Explicit literal outcomes without poisoning unrelated required groups or manufacturing deployment eligibility |
| Caller cancellation/deadline and owned output | Context priority even before early returns; no shared maps or slices |
| Consumer integration | Unsupported remains excluded before enablement; future catalog-only contributions retain an explicit limitation and existing denominator semantics |
| Production reversals | Named compiling mutations must prove binding, knowledge, restriction, context, ownership and unsupported-before-enable assertions on Linux and Windows |

Reuse validated transport/admission bounds rather than copying a second parser or pager. No live call, resource creation, fixture/RBAC edit, module/source refresh, protected merge or release is authorized. Adapter acceptance stays 1/7 and source characterization 7/7; this decision record adds no accepted feature.

## Checkpoint and recovery

This successor is documentation-only. Require its own full native/source QA, all original 385/380 control identities, file/tree/parent/identity readback and live PR evidence index. Read CURRENT_STATE and the final PR before repeating work. The next action after this reconciliation is the catalog-only pre-edit contract/expectation slice described above; if it cannot report limitations within current internal boundaries, proceed to the next SKU-family contract while preserving VMSS unsupported. Independent review/live/release obligations remain open. This is a natural evidence boundary, not a request for routine user approval.
