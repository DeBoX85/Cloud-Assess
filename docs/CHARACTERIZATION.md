# Characterization Baseline

Reference repository: `DeBoX85/azqr`

Reference commit: `8e4f0577f3615e6c9014c031bcad079f235369cc`

This document tracks source behavior that Cloud Assess intentionally preserves or changes.

## Implemented deterministic contracts

- Recommendation severity ordering: High > Medium > Low.
- Filter identifiers are compared case-insensitively.
- Tag keys are case-insensitive; tag values are exact-match.
- All include tags must match.
- Any matching exclude tag removes a resource.
- Explicit subscription/resource-group include scope takes precedence over matching exclusion, matching the reference implementation.
- Child resource/service results inherit a recorded parent tag-scope decision.
- Recommendation exclusions are case-insensitive.
- Finding summary deduplicates by normalized recommendation ID plus resource ID.
- SLA-category definitions and findings are excluded from the ordinary findings summary.
- Findings without a known recommendation definition are excluded from the ordinary findings summary.
- Summary ordering is deterministic by resource type and then recommendation ID.
- Default stages are graph, diagnostics, advisor, and defender.
- Graph remains mandatory for a normal assessment.
- Invalid stage names return errors rather than terminating the process. This is an intentional reliability change.
- The quality gate counts findings at or above the configured High/Medium/Low threshold.

## Next characterization targets

1. Resource ID parsing and scope helpers.
2. Subscription-ID redaction.
3. Recommendation applicability states.
4. Cost-period date calculation.
5. Scanner registry mappings.
6. Rule loading and disabled-rule behavior.
7. Azure Resource Graph result parsing.
