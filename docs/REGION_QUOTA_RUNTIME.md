# Owned region quota calculation contract

Status: IN PROGRESS, pre-edit contract. Accepted baseline is PR126,
`4768999aa10d85ddc69a357b81c9123e8c853b88`, tree
`be7d5256c6657fda436a5438c548c1e2fc4cc269`. The
[PR126 acceptance index](https://github.com/DeBoX85/Cloud-Assess/pull/126)
retains exact candidate and distinct accepted-push evidence. This contract does
not transfer that evidence to its proposal or declare region execution available.

## Requirement and source authority

This is the pure quota portion of B5 region functional equivalence under
[TARGET_SPECIFICATION.md](TARGET_SPECIFICATION.md) and
[DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md).
[REGION_QUOTA_RESERVATION.md](REGION_QUOTA_RESERVATION.md),
[REGION_QUOTA_SOURCE.md](REGION_QUOTA_SOURCE.md) and
[REGION_VM_QUOTA_SOURCE.md](REGION_VM_QUOTA_SOURCE.md) retain source contracts,
actual unchanged-fetcher observations and independently specified oracles.
Reservations remain a separate subsequent calculation slice.

AZQR stays `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree
`17d93b20c303f90f7843036be82f0dc32f3260f1`; APRL stays
`60eaddda76541f6adbc1c5ffa686829807e55e29`. Quota arithmetic and VM selection
come from `internal/scanners/plugins/region/quota/quota.go`, blob
`9d98d3d2c221b925fd5b4bbde4ddf78476146491`. REST provider filters retain the
other exact source blobs indexed in REGION_QUOTA_RESERVATION. Source caller
`internal/scanners/plugins/region/selection.go` uses VM, Network, SQL,
**App Service** and Storage labels; the space is significant.

Before: ProjectQuota validates and formats already-calculated rows, trusts flags
and cannot tell missing provider evidence from a successful empty response.
After: a bounded decoded-evidence calculation owns rows, calculates source-valid
values and carries explicit collection health into the existing table projection.
This does not add an Azure collector, coordinator dispatch or public CLI option.

## Proposed interface and selected evidence

CalculateQuota takes context, a selected subscription ID/display-name map,
explicit requested queries and decoded evidence. QuotaRequest holds SubscriptionID,
Region and QuotaType. QuotaEvidence holds Request, Status and Usages. QuotaUsage
holds ResourceName, LocalizedName, Current, Limit, CurrentKnown and LimitKnown.
QuotaCalculation returns owned Rows, Health and Table. No credential/client/cache,
clock or network argument is involved.

Subscription UUIDs are canonical lowercase and must belong to the selected map.
Selected names must be nonempty and satisfy existing label rules. Case-colliding
selection or duplicate normalized requests/evidence is rejected. Regions use the
existing canonical physical-region pattern. QuotaType is exactly one of the five
source labels. Evidence outside the explicit request set is rejected. Output
preserves declared request order and then received usage order.

Statuses are complete, partial, unknown and unsupported. Complete empty evidence
is valid, distinct from an absent response. Missing expected evidence, partial,
unknown and unsupported responses produce fixed deduplicated health warnings.
Unknown/unsupported evidence cannot carry rows. Partial evidence can carry its
valid prefix, with warnings retained. Status is completed or
completed_with_warnings, never a fabricated StagePartial schema value. Records is
the number of emitted rows. With no rows Table remains nil, but Health survives.
Nonempty Table health receives a separate copy of calculation warnings.

## Counts, missing fields and deliberate corrections

Known Current is in [0, 1000000000000]. Known Limit is in
[-1000000000000, 1000000000000]. Unknown numeric fields must carry zero, otherwise
the evidence contradicts its presence marker and is rejected. Missing counts or
missing raw resource identity skip the row with an explicit warning. Known
nonpositive limits skip arithmetic as in source, with unusable-quota health.
Negative current, malformed identity/text, duplicates and exceeded budgets fail
the calculation without a result. These corrections prevent source negative
usage, nil dereference and absent-field zero from becoming false capacity claims.

For admitted positive limits: available = limit - current; headroom =
float64(available) / float64(limit) * 100; at/over = available <= 0; near =
available * 100 < limit * 15. Integer comparison makes the exact threshold stable
within existing count bounds. Exact15 is not near, below15 is near, and near also
includes at/over. Existing projection rejects absolute headroom above 1e12; retain
this boundary explicitly, without silent clamping or widening the accepted API.
Source permits extreme ratios outside this bounded target domain.

| Current / limit | Available | Headroom | Near | At/over | Table status |
|---|---|---|---|---|---|
| 85 / 100 | 15 | 15 | false | false | OK |
| 86 / 100 | 14 | 14 | true | false | Near Limit |
| 100 / 100 | 0 | 0 | true | true | At/Over Limit |
| 110 / 100 | -10 | -10 | true | true | At/Over Limit |

VM retains case-sensitive contains Family, excluding aggregate cores. Network
skips NetworkWatchers, RouteFilterRulesPerRouteFilter,
RouteFiltersPerExpressRouteBgpPeering, RoutesPerExpressRouteCircuit and
BgpCommunityFilterRulesPerRouteFilter. SQL skips case-sensitive PerServer and
PerDatabase suffixes. App Service skips CustomDomains, HostNameBindings,
SslBindings, SslConnections and Certificates. Storage skips TotalBlobContainers,
TotalBlobs, TotalContainers, TotalFileShares, TotalQueues and TotalTables.
Case variants remain distinct exactly as in source. Filtered records still count
toward work/text validation; filtering is not a budget bypass.

## Identity, formatting and ownership

QuotaRow retains raw ResourceName as duplicate identity and gains an optional
DisplayName. CalculateQuota supplies LocalizedName, falling back to ResourceName.
ProjectQuota validates/budgets DisplayName and renders it when present; existing
callers with an empty DisplayName retain their previous rendering. Different raw
keys can share a localized label without an identity collision. Exact duplicate
raw names within the same query are rejected. Source raw-name matching remains
case-sensitive. Owned calculation rows retain raw identity; current canonical
PluginTable has only subscription identity and display cells. This change does
not claim a new serialized machine quota-identity field or change public schema.

Keep existing nine columns, integer formatting, %.1f%% headroom and status
precedence. No additional risk-summary output is introduced; AtRiskSummaries is
a separate source helper and is not the current auxiliary table contract.
Returned rows/table/warnings must not alias caller slices/maps or each other.
No per-run state is shared. Cancellation before admission, during loops and after
projection returns no result and the context error.

Requests and aggregate raw usages each have MaxAuxRows=8192 admission limits;
selected scope retains MaxSubscriptions=1000. Validate raw texts/counts before
filtering and output allocation. Raw evidence and projected output each retain
the existing text budget, MaxPluginTextBytes (16 MiB) minus 64 KiB, including
selected names/request labels/raw names/localized labels and the existing
128-byte per-row numeric/status allowance as applicable. Each label retains the
512-byte UTF8/control/noncharacter rules. Every supplied evidence record must be
admitted; no silent truncation. Errors and warnings use fixed generic messages
without echoing untrusted payloads, tokens or tenant identifiers.

## Intended changes, acceptance and rollback

Intended production paths are internal/plugins/region/quota_runtime.go and the
small optional DisplayName projection change in auxiliary.go. Add independent
quota runtime tests, a compiling negative-control runner under scripts/tests,
and explicit runner steps to both native jobs. Update this contract, handover,
roadmap, ledger and failure notes. No source captures, source production, module
versions, pins, reservation projection, registry, CLI, schema or Azure scopes
change. All 17 source captures and the unconditional source-byte guard remain.

Acceptance requires full literal rows/flags/health/table checks for the table
above; positive bounded cases against retained actual REST/VM outputs; separately
labelled negative-current/missing-field/identity corrections; all five source
filters and case variants; duplicate scope/query/evidence/raw keys and equal
display labels; missing/partial/unknown/unsupported/complete-empty evidence;
count/ratio/text/work exact and over limits, including filtered rows;
cancellation, owned results, repeat/concurrent calls and race checks. Critical
scope, identity, arithmetic/filter, budget, missing-health, cancellation and
ownership safeguards require isolated compiling mutations failing named
assertions, plus healthy/restored baselines on both native hosts. Compiler errors
or panics do not count as detected controls.

Review the final diff against this contract before final-head Linux quality,
windows-validation and pinned source execution with complete mandatory logs.
Inspect formal reviews, inline comments AND issue comments for automated review
state on the exact head. Protected expected-head merge and a distinct accepted
push need fresh evidence; planning checks cannot accept later code. No independent
person review or live Azure proof is implied. Rollback is a reviewed revert of
this slice, retaining historical evidence and accepted PR126 compatibility.

Next action: publish and verify this pre-edit contract, then implement the pure
quota calculation and its acceptance controls. Laptop/Azure access is not yet
confirmed. Live/DV001/load/fresh-OS/hosted-maintenance/Gate004/release and the
existing module-only advisory remain open; defer them while continuing offline.
