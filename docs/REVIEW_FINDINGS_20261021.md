# Expanded review findings

IN PROGRESS. Findings are separate from the approved Go security patch in PR185. This report is not a complete codebase audit or permission to merge/release. Exact source base for these initial observations is product183/458d5c5296c321c8e5e961f18eb5422d43c2eb18. Preserve frozen expectations and publish an explicit correction contract before any later product fix.

## R7-F01: conflicting case-equivalent filter tags change scope nondeterministically

Confirmed local reproduction; proposed P2 correctness/integrity finding, not a demonstrated exploit or live-tenant incident. internal/config/filters.go normalizeTags iterates a map and writes normalised keys without collision checks. LoadFilters accepts include.tags with env:prod and ENV:dev. Those are contradictory constraints for one case-insensitive tag name, yet one value silently overwrites the other.

The unchanged code compiles under Go1.26.10; an independent isolated test loads the same YAML256times and evaluates the same resource tagged env:prod. Result:21included/235excluded/0rejected, then the named assertion fails because a resource satisfying only one conflicting value was admitted. The count is an observed run, not a deterministic expected ratio. [Original test](evidence/20261021/seven-day-review/tag-conflict-test.go.txt) and [full reproduction](evidence/20261021/seven-day-review/tag-conflict-reproduction.log) are retained. The test resides in a disposable audit-probes worktree, not PR185; no existing expectation/production byte changed. Windows reproduction and a corrective candidate are not yet claimed.

The target specification says all include tags must match and tag names are case-insensitive. Microsoft [tag documentation](https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources) independently confirms case-insensitive names and case-sensitive values. The primary page was fetched after search. Existing tests check ordinary case-insensitivity and exact values but do not cover colliding case-equivalent input keys. No existing finding covering this case was located in current failure/systematic-review records.

Recommended later correction: reject conflicting normalised filter keys before credential construction or scope discovery; define equal-value duplicates and provider tag collisions explicitly. Preserve ordinary single-key/source-equivalent semantics and all current tests. Independent literals should cover include/exclude, insertion order, same-value duplicates, exact value case and whitespace/Unicode normalisation policies before implementation. Pinned AZQR8e4f0577internal/models/filters.go was read completely and uses the same map-overwrite normalization. This is inherited source ambiguity; introduction within the last seven days is not established. Document a later rejection as a deliberate safety correction rather than source parity. Until corrected, avoid ambiguous tag-filter files. This is an acceptance risk, not an authorisation to silently widen scope or weaken checks.

## Observations requiring further validation

R7-O01: generic authenticated HTTP Get/Post and PostStream error paths read complete bodies without byte caps, whereas newer bounded methods enforce attempt/body limits. Trace callers and existing explicit lifecycle/volume contracts before assigning priority or proposing a compatible bound.

R7-O02: structural resource-ID helpers parse positional/slash components without full prefix validation. Inventory, filter and findings caller validation must be traced before claiming an actual scope escape or malformed-result defect.

R7-O03: the canonical result sort keys omit some secondary fields and compare selected strings after case/space normalisation. Verify whether legitimate ties can produce different report order for the same logical input multiset before proposing a deterministic tie-break contract. No defect/changed expectation is claimed yet.

R7-P01: repeated current/status/stop narratives remain in normative documents and PR bodies despite authority headers. A machine-validated single active state with generated compact routing pages is a candidate recommendation, subject to full continuity review. Historical facts must remain immutable and separately accessible.

R7-P02: implemented draft work, protected acceptance and broader/core-v1 scope are frequently adjacent in status prose. Report them as separate acceptance states; do not turn control counts, coverage or source capture counts into an overall percentage. Roadmap priority and release acceptance need explicit decision ownership.


R7-P03 confirmed documentation drift: reviewed README at185still stated no retail implementation had started, while exact183contains it and own QA passes. Existing authority/header tests passed this contradiction. The current review corrects that active statement and removes superseded waits from compact CONTINUE/CURRENT_STATE/RESTART_PROMPT; frozen earlier prose remains accessible through Git. Proposed semantic state generation/validation is not yet implemented, so this does not certify every historical narrative as reconciled.
