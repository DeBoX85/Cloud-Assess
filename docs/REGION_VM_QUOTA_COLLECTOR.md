## Implementation checkpoint

[PR130](https://github.com/DeBoX85/Cloud-Assess/pull/130) implements this bounded
contract, UNACCEPTED pending final checks. Pre-edit97609c6a/treecbfa01c1/parentf076
original tree/all5 full remote files/Denis/ref verified before production edits.
Three independent VM test functions cover actual19 retained scenarios/full literals,
VM-specific boundary semantics and actual sovereign authenticated-client ownership.
All25 collector compiling controls must pass with restored baselines on both hosts.
FN093 fixture assumed empty table, corrected without product change; focused VM
PASS locally, broad local/native/source/review/acceptance still pending.
Current PR body is exact evolving run/recovery authority; prior prose is history.

# Bounded VM quota collection contract

Status: IN PROGRESS, pre-edit contract. Accepted baseline PR129 is
f0760de7eaf125e388b335320e49132b52e25da8, tree
0c71a2ffaeac4424d3b810be2fcd924f047530bb. The proposal PR body will hold exact
planning/candidate/run/merge/push evidence. Earlier source-only status in
[REGION_VM_QUOTA_SOURCE.md](REGION_VM_QUOTA_SOURCE.md) is historical; PR125 accepted
that characterization. [REGION_QUOTA_RUNTIME.md](REGION_QUOTA_RUNTIME.md) governs
existing VM Family filtering/arithmetic. [REGION_QUOTA_COLLECTOR.md](REGION_QUOTA_COLLECTOR.md)
governs the accepted shared request/decode/work/ownership engine.

## Source and implementation boundary

AZQR8e4f0577f3615e6c9014c031bcad079f235369cc/tree17d93b20c303f90f7843036be82f0dc32f3260f1
and APRL60eaddda76541f6adbc1c5ffa686829807e55e29 remain pinned.
Actual source quota/quota.go blob9d98d3d2c221b925fd5b4bbde4ddf78476146491
FetchVMQuota uses armcompute/v6 v6.4.0 UsageClient.ListPager. Retained19-case SDK
inputs/outputs prove GET /subscriptions/{UUID}/providers/Microsoft.Compute/locations/{region}/usages
with api-version2024-11-01, absolute continuation and canonical encoded skiptoken.
No latest-version substitution or source/dependency refresh.

Extend the existing RESTQuotaCollector to admit QuotaType VM using that exact
request. This deliberately replaces the target's absent VM transport with the
existing bounded authenticated getter, avoiding an additional SDK pipeline and
registration/retry defaults. Existing four providers keep their behavior.
VM first404/405 remains unknown, unlike the four REST providers' unsupported policy:
the retained SDK reports failures, not successful unsupported responses. Later
failure retains only validated prior pages as partial evidence, a correction to
source prefix loss. Fixed failure codes remain sanitized. No retries in collector.

The unchanged shared engine validates selected UUID/region/provider/origin and
all continuation destinations/path/version/query/cycles BEFORE credentials; custom
cloud endpoints use the caller's selected ARM audience. Bound64 pages,1MiB/page,
8MiB total,8192 raw rows,512-byte labels, JSON depth32/tokens65536/128 object keys.
Filtered counters consume budgets. Unknown counts remain explicit presence flags;
null items/malformed counts/negative current/overflow reject whole pages. Missing
current becomes unknown evidence at calculation rather than source panic. Positive
limits and case-sensitive Family containment remain CalculateQuota responsibility.
Raw rows and missing localized labels remain preserved; existing calculation's raw
name display fallback is already an explicit correction. Complete empty is owned
empty evidence, not an inferred unavailable service. Cancellation/deadline returns
nil result and context error. All shared guards/isolation policies remain required.

## Independent acceptance

Use actual unchanged retained VM inputs and exact recorded request sequences for
selection/localized/paged/empty/missing-current/null-item/negative/malformed and
first/later status cases. Assert full literal raw evidence and calculated table cells,
health/warnings, zero-versus-missing presence and preserved prefixes. Valid source
rows retain source meaning; unsafe negative arithmetic is rejected, not relabelled
parity. Actual source captures remain byte-identical and hash-checked.

VM-specific synthetic fixtures constrain same-origin/path/version/selected scope,
no authenticated foreign request, unknown first404/405 and partial later failures,
malformed whole-page atomicity and cancellation. Real authenticated bounded HTTP
client fixture verifies sovereign ARM audience/exact VM GET/no body/closed response
with synthetic credential and injected transport; no listener or live Azure call.
Reuse shared work-boundary/concurrency/control coverage without duplicating every
engine fixture. New compiling controls change VM version and VM status policy;
require named assertion failures plus healthy/restored baselines on Linux/Windows.
Existing23 collector controls and mandatory gates remain unchanged in meaning.

Affected production: quota_collector.go provider routing and first-page status policy.
Tests: quota_collector_test.go invalid-provider fixture, new vm_quota_collector_test.go,
existing collector mutation runner. Documentation: this contract, continuity/roadmap/
ledger/failure index and shared collector scope clarification. No modules, source
captures/pins, registry, public CLI/schema/report changes or reservation collection.
Whole region race/whole suite/vet/docs/staged original tree/full remote bytes/Denis
identity/parents/ref review, exact candidate native/source/review/protected merge and
distinct accepted-push proof required. Local checks are not full native acceptance.
Rollback through reviewed revert on current accepted branch, never history reset.

## Recovery and limits

PR129 completed pre/post native and source acceptance: candidate37703322161 jobs
113071785275/113071785549, source37703322500/job113071696381; PR jobs actually checked
preview9df5a96fa9a09a5b145ab129f27a6bf577e5e55e, exact tree0c71a2ff and parents27b45950/54b4957.
Accepted-push37705050463 jobs113077299205/113077298952 and source37705050524/job113077298666
actually checked f0760de7. All23 controls both hosts;17 captures150 chunks/full hashes/
bytes; automated review completed54b4957, not human approval. Full final authority is
[PR129](https://github.com/DeBoX85/Cloud-Assess/pull/129). No need to replay unchanged gates.

Laptop/Azure restoration is expected but NOT confirmed. Live/DV001/load/freshOS/
maintenance/Gate004/release/module-only advisory remain open. No resources/roles/
new fixtures/production substitution/exposure/release authorization. Internal helpers
are not public region execution. Next after bounded VM acceptance: reservation and
remaining adapters, coordinator/public/all-format integration. No blocking input now.
Publish and remotely verify coherent checkpoints; raw scratch logs/cache are not
backups. Inspect uncertain remote outcomes before retries, preserve historical records.
