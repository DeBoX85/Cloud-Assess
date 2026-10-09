# Current handover authority, 2026-10-09

Read [SESSION_HANDOVER.md](SESSION_HANDOVER.md) and
[HANDOVER_CHECKPOINT_20261009.md](HANDOVER_CHECKPOINT_20261009.md) before following
an implementation status, baseline, executor description or next action below.
Accepted implementation extends through [PR131](https://github.com/DeBoX85/Cloud-Assess/pull/131),
`b653a3abfc35590a185531095a168a9a107e769e`. The handover PR on
`docs/handover-20261009` holds its own exact publication/CI/merge/push state.
Verify later live refs and PR evidence; this header does not assume that proposal
is accepted. New-chat instructions: [NEW_CHAT_PROMPT_20261009.md](NEW_CHAT_PROMPT_20261009.md).

PR114/119 audit and PR120-131 bounded region milestones are accepted offline;
historical IN PROGRESS/UNACCEPTED/pending baselines or next tasks below are snapshots,
not current instructions to replay completed work. The checkpoint identifies
each accepted milestone and its PR evidence. Preserve the underlying source
contracts, explicit corrections, bounds, mandatory tests and historical results.
Public region execution, full feature parity, live access, Gate004 and release
readiness remain unestablished. No new whole-project or independent-person review
is claimed by this documentation task.

---

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

The exact default SDK synthetic scope is
`https://management.core.windows.net//.default`: azcore v1.23.1 ARM runtime defines
the public audience with a trailing slash and appends `/.default` in its pipeline.
Do not normalize that observation to the REST client's explicitly supplied scope
or infer endpoint host from the token audience. Initial local fixture asserted the
REST scope and failed before transport (FN081); correct only the independent exact
SDK expectation, not the production pipeline or accepted audience set.

The injected credential and transport must honor context cancellation themselves.
Initial synthetic stubs ignored context and allowed a pre-cancelled request (FN081);
that is a fixture defect, not evidence of an SDK/source cancellation defect. Corrected
stubs return ctx.Err before producing a token or accepting a request. The SDK also
encodes/orders continuation queries: supplied `api-version=2024-11-01&$skiptoken=...`
is observed as `%24skiptoken=...&api-version=2024-11-01`. Require that exact request,
preserving the original supplied nextLink in input. Do not normalize away a mismatch.

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

## Observed source and independent fixture record

Planning50736112/treeb2a75e7d/parent28c7d996 and observatione72e9ebf/treea70511f2/
parent50736112 were verified remotely against complete original local trees/bytes,
ordered parents and Denis identity. [PR125](https://github.com/DeBoX85/Cloud-Assess/pull/125)
is the exact evolving final-head/merge/push recovery index. Source37243637952/
job111557127792 passed seven unchanged harnesses and source/module/APRL isolation,
emitted15 files/123 contiguous chunks, then failed the declared absent-golden check.
Classify that run as evidence gap, never source PASS. All chunks independently
reconstructed and SHA256/full UTF8 compared; original13 unchanged. New2 match
actual local corrected source observations. FN081/FN082 retain failed fixture and
local-command attempts without transferring success to those attempts.

| Capture | SHA256 | Bytes/chunks |
|---|---|---|
| source-vm-quota-inputs.json | 5b5319746110e4e86cc6753ec5b20a3306d6290d03b0e829682d2511324d3a1c | 7730/3 |
| source-vm-quota-outputs.json | b6311db11c7a62c74818ef694c2991cd1ab46d3a8c0bf1337ed5fcd7cb480c32 | 9708/4 |

Nineteen scenarios retain five arithmetic rows; case-sensitive Family containment;
skipped missing names/value/limits and nonpositive limits; missing/localized display
and risk fallback; nil empty; two pages; first/later400/403/404/405 and malformed
responses; pre-cancellation; actual missing-current/null-item panics. Full literal
records include every row/flag/percentage/risk/request/error/panic/nil distinction,
plus exact SDK/source/options premises. Separate hashes preserve original input and
output bytes. Five compiling controls alter family selection, source panic marker,
unsupported-as-error, nil-empty and actual SDK request version; each must fail the
named complete assertion with healthy/restored baselines. Local source19-case,
focused full literals/hashes and five compiling controls/restored passed. Final-head
native/source and accepted-push remain required; bounded oracle sensitivity does not
certify target runtime, all source branches or source memory safety.
