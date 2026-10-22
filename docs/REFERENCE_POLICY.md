# Reference baseline and deliberate improvements

Denis clarified this direction on 2026-10-21. Cloud Assess must match the features, functions and capabilities of his [AZQR fork](https://github.com/DeBoX85/azqr), while progressively improving quality, robustness, security, usability, performance and capability where justified. AZQR is the initial comparison baseline, not the ceiling for the product or its only source of evidence. A compatibility gap does not become complete merely because an enhancement is proposed.

## Baseline identity

The authoritative compatibility reference is **DeBoX85/azqr** at commit `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree `17d93b20c303f90f7843036be82f0dc32f3260f1`. Its parent is `28d63fdcd919428294d02c286305cf46116e9cb2`. The [fork commit/tree](https://github.com/DeBoX85/azqr/commit/8e4f0577f3615e6c9014c031bcad079f235369cc) and main ref were verified through GitHub on 2026-10-21; main currently equals this pin. This observation is not a guarantee about future branch state. APRL remains pinned at `60eaddda76541f6adbc1c5ffa686829807e55e29`.

The README and target specification name the fork. `scripts/capture-region-source.py` fetches that fork and checks the exact source commit/tree and APRL commit. `.github/workflows/import-rule-baseline.yml` imports from that fork at the same source SHA. This policy changes no pin, capture, expected result, production code or QA control.

| Evidence source | Role | Change boundary |
| --- | --- | --- |
| Denis's fork at the frozen commit | Initial feature/function/capability inventory and reproducible behavioral comparison | Preserve exact identities and healthy comparison controls; track gaps to closure |
| Later changes in Denis's fork | Preferred AZQR source for prospective feature and capability gap research | Record the observed revision and differences separately; no automatic pin or oracle refresh |
| Azure primary documentation and schemas | Platform semantics, security requirements and independent evidence | Record version/date and reconcile disagreements with the source rather than copying a source defect |
| Other projects, including upstream Azure/azqr | Secondary engineering, capability and process comparisons | Keep repository/revision attribution; they do not replace the fork or certify Cloud Assess |
| Cloud Assess independent tests and operational evidence | Validate explicit product contracts, deliberate corrections and measurable improvements | Distinguish offline, independent, live and operational evidence; retain unfulfilled obligations |

## Compatibility and improvement process

Feature, function and capability parity remain required. Behavioral equivalence applies to the pinned healthy behaviors under the applicable contracts, with already documented branding/configuration/CLI differences and deliberate corrections. Preserving unsafe or incorrect behavior is not a product goal. An enhancement does not excuse losing an existing capability, concealing unknown evidence or weakening a security boundary.

For each proposed change beyond the current contract:

1. Record the problem or opportunity, source evidence, affected capability and current behavior. Classify it as a defect correction, enhancement or process improvement. Identify whether the current user-authorized scope covers implementation.
2. Define the intended behavior, compatibility impact and success criteria before production edits. Freeze independent expectations when required. Keep historical source captures and expectations intact; add explicit correction or enhancement expectations separately.
3. Use a bounded, separately reviewable change with dependency order, risk and rollback/restoration evidence. Security or correctness corrections may precede full parity where necessary and authorized; unrelated enhancement work must not silently expand the active task.
4. Preserve the required control inventory and run the change's own exact-head QA. Compare healthy behavior and affected capabilities, and measure any performance claim under a declared workload. Record regressions, remaining gaps and acceptance obligations honestly.
5. Update the relevant contract, capability inventory, current state and final PR evidence. A future baseline refresh needs a deliberate decision and separate source/delta review, never a moving latest-branch comparison or overwritten frozen oracle.

The expanded audit may recommend improvements across all these areas. This direction authorizes recording and evaluating them; the current audit and task8-to15 dependencies remain active. It does not authorize an unrelated feature, Azure operation, source/dependency refresh, protected merge or release. Current implemented draft scope, protected acceptance and remaining live/independent/operational obligations are recorded in [CURRENT_STATE](CURRENT_STATE.md).
