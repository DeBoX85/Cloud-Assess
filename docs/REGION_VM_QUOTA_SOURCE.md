# VM quota SDK source characterization

Status: IN PROGRESS, source-only. Accepted baseline PR124 is
`28c7d996bf040e3410dce5aef1cf3051dd187696`, tree
`43eb87ddd5c89a0592f571b8f8bcaf89aba97821`; its final native/source and separate
accepted-push evidence are indexed in [PR124](https://github.com/DeBoX85/Cloud-Assess/pull/124).
The next source checkpoint must not inherit earlier test results.

## Exact source and trigger

AZQR stays `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree
`17d93b20c303f90f7843036be82f0dc32f3260f1`, APRL
`60eaddda76541f6adbc1c5ffa686829807e55e29`. The unchanged
`internal/scanners/plugins/region/quota/quota.go` blob
`9d98d3d2c221b925fd5b4bbde4ddf78476146491` contains FetchVMQuota and
AtRiskSummaries. Source go.mod pins armcompute/v6 v6.4.0, whose UsageClient GET
uses `/subscriptions/{id}/providers/Microsoft.Compute/locations/{region}/usages`
with api-version2024-11-01. This SDK path was inspected before the capture; actual
requests must be independently recorded and checked.

FetchVMQuota includes case-sensitive names containing Family, excludes cores and
lowercase family, skips missing names/value/limit or nonpositive limit. CurrentValue
is dereferenced after those checks without testing nil. A null page item also
dereferences nil. Missing localized name becomes empty, unlike REST string names.
Available, percentage and flags follow the source arithmetic described in
[REGION_QUOTA_RESERVATION.md](REGION_QUOTA_RESERVATION.md). Successful empty returns
nil entries. SDK/decode/page failure returns nil and an error, discarding a prior
prefix; unlike REST,400/404/405 are not explicitly converted to unsupported nil
success. No explicit page/work bound is applied by this source loop.

## Characterization and correction boundary

Invoke the actual unchanged FetchVMQuota with synthetic token credentials and an
injected SDK transport accepting only declared GET requests, no body and the
synthetic bearer. Fail closed for unmatched requests and real network fallback.
Disable retries and automatic provider registration in the supplied fixture options
so this observation cannot issue registrations or hide request counts. The source
accepts caller options, so this is a controlled fixture policy, not evidence of its
default caller middleware safety, every cloud or target adapter confinement.

Use literal valid/filtered/missing/negative/threshold/empty/paged/malformed/denied/
unsupported/cancelled response bodies. Recover only at the harness boundary to
record Panicked=true for malformed source input; never call a recovered panic a
successful collector result. The harness's test passes when expected source behavior
is retained, not when the source itself is safe. Do not reproduce arithmetic in a
replacement helper and label it unchanged collector execution.

Retain complete inputs, source blob metadata, entries, risk strings, Failed/Panicked
and exact request sequence. Independent full output literals must distinguish nil
from empty; separate hashes preserve original fixture bytes. Selected compiling
controls must fail named assertions, with healthy/restored baselines, on both native
hosts. Preserve all13 existing source captures and unconditional retained-byte
equality. New absent goldens deliberately fail after observations are emitted;
classify that evidence gap before retaining independently reconstructed bytes.

Before target runtime implementation, define canonical selected identity, bounded
nonnegative counts, positive limits, missing counts as invalid/unknown, ownership,
cancellation and honest stage completeness. Valid source Family selection and
arithmetic are parity evidence; nil panic, invalid numeric input and silent data
loss are source limitations requiring explicit corrections. This source-only slice
adds no target collector or public region enablement. CRG SDK characterization is
separate pending work, followed by pure calculations, guarded adapters and complete
coordinator/report integration.

Changed files are a source harness, two retained fixtures, independent source tests,
compiling control script, runner/workflow registration and continuity records. Pins,
modules, production, schema and CLI stay unchanged. Final head requires native
Linux/windows and strict source success, protected expected-head merge and distinct
accepted-push proof. Revert through reviewed history against current accepted branch,
never force/reset it. Exact revisions/runs are retained in the current PR index.

Laptop/Azure access is expected soon, not confirmed restored. Approved core/plugin
live checks can resume independently after confirmed access/scope; DV001 requires
the recorded wholly nonproduction hierarchy. No Azure resources, roles, fixtures,
production substitution, service exposure or release. Live/load/fresh-OS/maintenance/
Gate004/release/module-only advisory stay open. Self-review and independent literal
oracles do not establish independent-person review. Preserve earlier workspaces and
historical evidence; scratch-only work is not a verified backup.
