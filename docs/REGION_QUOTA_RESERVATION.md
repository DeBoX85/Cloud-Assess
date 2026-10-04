# Region quota and reservation continuation contract

Status: IN PROGRESS, source characterization only. Accepted baseline is PR123,
`4f313ed5ce68856664537666607a2cd74b228331`, tree
`d346154a1069b33b179745c3b4b3edb0d782c07e`. Its [acceptance index](https://github.com/DeBoX85/Cloud-Assess/pull/123)
records final native/source and separate accepted-push evidence. This record does
not transfer those results to a new head or certify live region execution.

## Source identity and reviewed scope

AZQR remains `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree
`17d93b20c303f90f7843036be82f0dc32f3260f1`; APRL remains
`60eaddda76541f6adbc1c5ffa686829807e55e29`. Paths below are relative to
`internal/scanners/plugins/region/` in that exact source.

| Path | Git blob | Reviewed behavior |
|---|---|---|
| quota/quota.go | 9d98d3d2c221b925fd5b4bbde4ddf78476146491 | Name decoding, shared paging/arithmetic, VM-family selection, risk summaries |
| quota/appservice.go | 60ce83c4d2e7dbed133d1ac1389f9094ee1d39c9 | Web endpoint and exact skip list |
| quota/network.go | 8dd5eaa925f65388c966b1f4d22b9c13442828e5 | Network endpoint and exact skip list |
| quota/sql.go | 33d990e8016b9ee6a8088a203cb3a365bc69f357 | SQL endpoint and case-sensitive suffix filters |
| quota/storage.go | 62e898c9492a62bb96eef115f5ea352d0918983b | Storage endpoint and exact skip list |
| crg/crg.go | 1ffb6e37cc1c21f6d8b5c9ce483f793387e0e6bc | SDK group/list/Get traversal, counts and status precedence |

The source HTTP client at `internal/az/http_client.go`, blob
`05d85c5a50ee81f088731160020a914af83e1ca5`, supports an injected transport.
The reviewed source tests are quota/quota_test.go blob
`061ff66ed7bd275e118112a9c1c03153d64988e2` and crg/crg_test.go blob
`72e398123b1a0acf0ab51efff66f136b4549ca06`.

## Source facts and evidence limitations

Quota keeps positive limits, computes available = limit - current and headroom =
available / limit * 100. Near-limit uses strictly less than 15 percent, including
negative headroom; at/over-limit uses available <= 0. Exactly 15 percent is not
near-limit. Negative current is not rejected by source and can exceed 100 percent
headroom. The source field comment claiming a nonnegative near-limit condition
does not match execution. Zero/negative limits are skipped, not interpreted as zero
remaining capacity. Risk summaries trust stored flags, preserve order and use the
localized name when present, otherwise the resource name; percentages use %.0f.

Shared REST quota decoding accepts object or string names; JSON null yields empty
name fields. Unsupported HTTP400/404/405 returns nil with no error, including after
a successful page and discarding its prefix. Other request/decode failures return
an error with no entries. Successful empty REST pages produce a nonnil empty slice.
Continuation destinations and total work are not bounded by this source loop.
These are inspected source behaviors, not claims of live endpoint support.

Web skips CustomDomains, HostNameBindings, SslBindings, SslConnections and
Certificates. Network skips NetworkWatchers, RouteFilterRulesPerRouteFilter,
RouteFiltersPerExpressRouteBgpPeering, RoutesPerExpressRouteCircuit and
BgpCommunityFilterRulesPerRouteFilter. Storage skips TotalBlobContainers,
TotalBlobs, TotalContainers, TotalFileShares, TotalQueues and TotalTables. SQL skips
names ending in PerServer or PerDatabase. All matching is case-sensitive.
REST versions remain Web2023-01-01, Network2022-07-01, Storage2023-01-01 and
SQL2021-11-01, at the source's subscription/provider/location usages path.

VM quota separately uses the Compute SDK and case-sensitive contains("Family").
It excludes aggregate cores, skips nil/nonpositive limits and missing names, but
dereferences CurrentValue without a nil check. Its successful empty result is nil.
The next REST capture does not claim to execute or cover this SDK path.

Reservations calculate allocated from the expanded VM list and reserved from SKU
capacity. Status precedence is allocated==0: Idle; available<0: Over-Allocated;
available==0: At-Capacity; otherwise Available. Zero reserved and zero allocated is
Idle. Missing utilization is inferred zero by source, not proven idle. Location
comes from the group, lowercased; the returned reservation name comes from Get.
Group-page failure returns an error. Reservation-page failure logs and stops that
group, retaining earlier data and continuing other groups. Get failure logs and
drops that reservation. Both latter cases can return incomplete success without
structured health. Successful empty result is nil. These are source limitations;
no Cloud Assess reservation collector has yet inherited them.

The existing quota arithmetic test constructs its own UsageEntry; the reservation
status test repeats its own switch. Neither demonstrates collector execution.
AtRiskSummaries does have direct source tests. Characterization must invoke the
unchanged fetchers with synthetic credentials and a fail-closed injected transport,
record exact requests and complete outputs, and distinguish this from pure-helper
and live Azure evidence. No actual Azure request or credential is needed.

## Ordered work and acceptance boundary

1. Capture unchanged REST quota fetchers, full rows, risk summaries, exact provider
   filters, names, empty/unsupported/denied/malformed/later-page branches. Preserve
   all existing eleven source captures byte for byte. New absent goldens fail the
   existing unconditional retained-byte guard after emitting observations.
2. Retain independently reconstructed observations with hashes. Compare full
   literal outputs and requests; use compiling controls for threshold, filter,
   unsupported-after-prefix and denial semantics in both required native jobs.
3. Characterize the separate VM SDK and reservation SDK paths with strict request
   fixtures, including missing fields and incomplete Get/list traversal. Do not
   replace those fetchers with copied arithmetic and call it source execution.
4. Before target production edits, specify bounded decoded evidence, canonical
   selected identity, nonnegative counts, positive quota limits, missing utilization
   as unknown, honest completeness, ownership/cancellation and text/work budgets.
   Preserve valid source arithmetic/status while documenting deliberate corrections.
   Existing ProjectQuota/ProjectReservations remain format projections, not collectors.
5. Implement pure calculations, then guarded adapters, coordinator and public
   all-format/default/empty/partial integration. Each accepted slice needs final-head
   Linux/windows/source proof, protected merge and distinct accepted-push evidence.

This slice changes source test harnesses, retained synthetic fixtures, test controls,
runner registration and continuity records only. No target collector, public region
enablement, dependency/pin/schema change or Azure write is authorized by this record.
Rollback uses a reviewed revert against the current accepted branch; never reset
protected history. Preserve all previous audit and failed/intermediate evidence.

Laptop and Azure access are expected soon, not confirmed restored. Once confirmed,
the existing approved core/plugin live queue can proceed independently; DV001 needs
the recorded wholly nonproduction nested hierarchy. No resource creation, role
change, production substitution, service exposure or release publication. Live,
fresh-OS, maintenance, Gate004/release and the module-only advisory stay open.
