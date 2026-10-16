# Current continuation authority: bounded VM collector implementation, 2026-10-16

The user authorized continuation after completed draft PR146/37954dd. Its own full Linux/Windows/source gates passed and the final evidence index remains in [PR146](https://github.com/DeBoX85/Cloud-Assess/pull/146). Continue the bounded VM catalog collector under the frozen [contract](REGION_VM_VMSS_SKU_CONTRACT_20261016.md); VMSS scoring remains unsupported. Read [the active collector record](REGION_VM_SKU_COLLECTOR_20261016.md) and the live draft implementation PR for actual current head and QA.

This publication is test preparation: a runner consumes all67unchanged frozen scenarios and the four immutable primary files; the collector API is absent. Inspect the expected missing-API compilation on both mandatory native hosts before production edits. That failure is preparation, not product regression reproduction or acceptance. Then complete the already-authorized collector, independent boundaries/HTTP/context/ownership/projection tests and compiling reversals with restored baselines; require fresh complete corrected-head native/source evidence before STOP at the collector checkpoint.

Keep prior PR141-146 draft/open/unmerged and protected bootstrap/core-v1 a6fee427 unchanged. Accepted target adapters remain1/7; source characterization7/7. Preserve the frozen67-case manifest/primary bytes, source23captures/pins/modules/public schemas, all353Linux/348Windows prior controls and audit corrections. No Azure operations, fixtures/RBAC, refresh, sub-agents, protected merge or release. Coordinator/public integration and VMSS eligibility remain outside this slice. Earlier pre-edit stop/branch/pending wording below is historical.

---

# Current continuation authority: VM/VMSS pre-edit contract, 2026-10-16

The user requested all relevant handover/status updates and continuation to the next natural stopping point. [Storage collector completion](STORAGE_COLLECTOR_COMPLETION_20261016.md) and the [full PR145 evidence snapshot](STORAGE_COLLECTOR_PR145_EVIDENCE_20261016.md) now record its completed exact-head Linux/Windows/source proof. Do not rerun the completed Storage slice or treat the historical pending wording below as current.

Active branch: docs/vm-vmss-contract-20261016, based on verified PR145 a14128bc94bf531dd1c28553db14d2be011b2214. Read [the VM/VMSS pre-edit contract](REGION_VM_VMSS_SKU_CONTRACT_20261016.md). This slice freezes primary Compute schema/examples and literal future expectations; it adds no VM collector, production change, new executable control or public integration. VM size catalog and VMSS deployment eligibility are explicitly separate; VMSS evidence remains unsupported in the first proposed collector. Future manifest scenarios have NOT executed against the absent API.

This successor requires its OWN complete Linux/Windows/source gates. At publication these are pending; its live draft PR is the final exact-head commit/tree/parent/preview/log evidence index. Finish those gates, update that index, then STOP at the contract/expectation checkpoint before writing VM production code. A later implementation must first publish compiling tests consuming the frozen oracle and inspect the expected missing-API preparation; compilation failure is never a product pass.

PR141-145 remain draft/open/unmerged at their recorded heads; protected bootstrap/core-v1 remains a6fee427e19feeef28691737ea9bc5026086efa0. Accepted target adapters remain 1/7; Storage is an unaccepted candidate, and this contract adds no accepted adapter. Source characterization remains 7/7. Preserve all audit corrections, 353 Linux / 348 Windows controls, original 23 source captures, pins/dependencies/public schemas and open roadmap/live/security/license/release obligations. No Azure operations, fixtures/RBAC, refresh, sub-agents, protected merge or release authorization. Earlier dated next actions and stop instructions below are historical unless their behavioral contracts remain applicable.

---

## Storage collector Windows QA correction, 2026-10-16

Candidate `c2d5c5df48fa0e4939b5b3a3b7b67f9bd391ca8c` failed Windows native job `116989921866` in run `38906095075`: the new mutation helper decoded raw CRLF checkout bytes while its multiline anchors used LF. Six controls ran before `stage-physical` rejected the mismatched anchor. The helper restored the original source in its `finally` block. This is QA-ST-03, a control portability defect, not a product pass or frozen-oracle change. Full failed Windows log: 36,704 characters, 36,706 UTF-8 bytes, SHA256 `8bf7ffc77d1bcfdc47a788cf544e8189481f4141ad3feaf388df80aaa0380337`.

Correction normalizes only the helper's in-memory anchor text to LF; exact original checkout bytes remain the restoration authority. Production collector, frozen71-case oracle, primary example and all25control assertions remain unchanged. All25controls passed locally with both LF and simulated CRLF, including exact CRLF restoration followed by original LF restoration; fresh complete Linux, Windows and source gates must pass on the corrected published head. Local logs: storage-controls-portable-lf.log: 2215 bytes, SHA256 `ceec94ce21e61bf5b791b65ed735599def5ba85e2df3cce8e2b5135f45fd4b4f`; storage-controls-portable-crlf.log: 2215 bytes, SHA256 `ceec94ce21e61bf5b791b65ed735599def5ba85e2df3cce8e2b5135f45fd4b4f`. The first CRLF launcher omitted the pinned Go PATH and stopped before testing; the corrected launcher completed all25controls. The live draft PR145 body remains the final exact-head evidence index. No protected merge or next-family implementation is authorized by this correction.

# Storage implementation candidate prepared, 2026-10-16

Status at publication: implementation/local verification complete; full own exact-head Linux/Windows/source QA pending. The live PR145 final evidence index provides actual final SHA/tree/parent/preview, jobs/full logs/control identities and final disposition. Keep draft/unmerged and do not transfer prior evidence to this candidate. STOP after complete corrected-head proof before coordinator/public integration.

## Implementation and retained contract

New storage_sku_collector.go implements the exact published API. Collect validates scope/aliases/physical target before one bounded Storage2023-01-01 global GET, preserves status/context priority, rejects undocumented paging, validates/bounds the entire JSON envelope including filtered/ignored metadata, and admits all rows atomically. It unions SKU+kind regional/restriction evidence without overwrite; preserves unknown optional metadata, explicit display mappings, stage exclusion, scoped source-compatible restriction states and capability conflicts. Missing kind/type remains unqualified rather than falsely unavailable. Operational failure exposes empty unknown evidence, never a valid prefix or unsupported success.

Evidence validates/copies stable collection and required-kind inputs under entry/text limits. Only inventory-required kinds participate, with no caller-map/slice aliases. Any unresolved required group yields exactly unknown/empty for the Storage entry in the existing name-only consumer. It does not infer inventory kind, build the required-kind coordinator, change percentage arithmetic/schema, register a plugin or add CLI/network execution. Full decoder/unknown/output returns preserve terminal context.

The original71scenario JSON and exact primary fixture are byte-identical to cd5d234. The full manifest runner from f0707bb and boundary/client tests from07eb37d remain unchanged through implementation. Supplemental review tests retain late unknown cancellation, nonrequired-kind isolation and hard-restriction precedence. Existing audit/source/Compute controls and assertions are unchanged.

## Local verification

- All71frozen manifest scenarios and the entire primary248-row fixture passed, including exact complete grouped maps/projection maps and request/context/error expectations.
- All body/row/group/JSON/alias/scope/projection boundary, actual bounded-client body-closure/audience/status-error, repeated/concurrent ownership and deterministic in-work cancellation tests passed.
- All25new production reversals compiled, failed their named independent assertion, restored exact original bytes and passed restored baselines. Control output2215bytes SHA256ceec94ce21e61bf5b791b65ed735599def5ba85e2df3cce8e2b5135f45fd4b4f.
- Complete local go test -race ./... and go vet ./... passed serially AFTER the source-mutating controls completed/restored. Race log2618bytes SHA256e3e7348322690479672758db7e209f9879e9df804ff8176fa2d6207f987a648c; vet log0bytes SHA256e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855. Focused corrected log11705bytes SHA256b018de0bd7de392cfda5edc3f26cabf9f9413c678d571eb20af4baaf221e2e24.

The workflow adds exactly one mandatory Storage-control step to each existing native job. All prior steps/controls, fuzz/property budgets, coverage floor and source/capture checks remain unchanged. Whole-job wall budgets increase by5minutes each, Linux25->30 and Windows30->35, solely to accommodate the25additional compile/fail/restore controls; jobs remain bounded. This is new mandatory work, not suppression or a relaxed assertion. Expected unique control identities are old328+25=353Linux and old323+25=348Windows, with exact prior identity retention required.

## QA-ST-01: implementation review cancellation correction

The first unpublished implementation passed all71frozen cases and published boundary tests. Supplemental TestStorageSKULateUnknownCancellation independently reproduced a missing final terminal-context check on the nil/unknown Evidence return. Full failing local log349bytes SHA256eee81479507ec3846a9c8af8dd24d65e4ee5459fd3f0402cebdd66b2f0a3fed5. This was a new unaccepted feature gap, not a previously accepted audit defect. A private incomplete-return helper now checks ctx.Err on every unknown path. The unchanged supplemental regression and a compiling reversal/restoration control pass; no oracle expectation changed.

QA-ST-02: first local control invocation rejected its request-version anchor because gofmt aligned the constant with two spaces. It restored original source and was not credited as a product pass. The exact helper anchor was corrected, then the complete25control run succeeded. First helper log51bytes SHA25607e44a3e8945f7053a8b1fe716823f21507a3330468e20759a241ed11d75323b. Prevent recurrence by inspecting formatted production before freezing literal anchors; mutations must compile and fail named tests, not merely return nonzero.

## Publication, recovery and rollback

The predecessor test-only heads f0707bb and07eb37d are durable and intentionally missing the API; their red is preparation only. Both first-host native missing-symbol logs are retained: Linux116981113328/33822characters/33824bytes SHA256e29436b13a727d3518be60f5570ac1905450db06aad6c70ae4adf7882732da87; Windows116981113477/33189characters/33191bytes SHA256c98e350c03b5f9e73a32e3920ea0ac546476738a8b7175114cf92f0dc62c9f17. Additional boundary-preparation38904120528 is also expected red; source38904120413 succeeded separately. No preparation result accepts production.

Local shell fetch works and push lacks credentials; connected expected-head publication is used without changing credentials/permissions. Verify exact blobs/tree/parent/Denis identity and live ref after publication. Historical Linux/Windows workspaces and local metadata-only commits remain preserved. If interrupted, inspect live PR145 and current recovery authority before replaying work.

Rollback uses a targeted successor commit after dependency review. Remove/revert the production collector, its mutation helper, two native execution steps and corresponding job-budget allowance together, retaining independent tests and history; the API regression should intentionally return red. Do not reset/force protected history, delete expectations or refresh captures. No schema/data/Azure rollback is required for this internal library.

Next after completed checkpoint evaluation: VM/VMSS pre-edit contract under REGION_COMPUTE_SKU_SOURCE_20261013.md and REGION_SKU_TARGET_ANALYSIS_20261013.md. Required-kind coordinator and public/all-format execution remain deferred. Accepted target adapters remain1/7 until protected disposition; this Storage candidate is unmerged. Source captures7/7, remaining adapters/FullCost/ancillary parity/AR-02/load/live/nonempty/sovereign/visibility/Arc/freshOS/hosted maintenance/module-only advisory/human/license/security/Gate004/release remain open. No sub-agents or Azure operations.

---

## Historical test preparation

## Published preparation and exact budget interpretation

Preparation f0707bb1983dd0d9578cee84ce06999155d20641/treeb1678a8663e00e3d7a246270ada95b7846ff4aa1 is published in draft PR145. Native38903459610/Linux116981113328 and Windows116981113477 both stop on the three declared missing API symbols; this is expected preparation, not an accepted product failure/pass. Source38903459779 succeeded separately. Shell push lacked credentials; connected publication/readback succeeded without identity/permission changes. The metadata-equivalent local ac671a6 marker is preserved; the task worktree now uses genuine remote ancestry.

Before production, boundary tests define the published contract's accounting precisely: body work charges each object and each array element once, including ignored metadata; lexical token counting remains the existing quotaJSON guard. Projection input text charges collection SubscriptionID/Region/Status/FailureCode plus every group/unqualified/required SKU/kind label once. Projection logical-entry work charges each outer SKU/unqualified/required entry and each group/required-kind entry once, max65536. Validation and projection may revisit stable entries, but output keys are separately bounded by8192/512bytes; no unbounded allocation or work is inferred from an empty result. Output text has its own counter. These are testable operational definitions, not changed semantic expectations.

The raw1MiB bound dominates the16MiB-64KiB decoded-body text bound, and8192groups/rows/output names dominate their larger aggregate text limits; tests must state these dominance relations rather than claim impossible independent +1 witnesses. Projection inputs can reach16711680bytes and tests exercise its exact boundary:50collection identity/status bytes plus2040*(1SKU+15kind)labels of512bytes, reducing one label by50. They also exercise65536/65537logical entries, bytes/rows/depth/properties/tokens/body work/aliases/scope/kinds/names, malformed filtered metadata, real bounded-client closure/status/audience, per-call ownership/concurrency and deterministic in-work cancellation. The original71manifest oracles are unchanged.

# Bounded Storage SKU collector, 2026-10-16

Status: test preparation only; collector API not implemented. The user authorized continuation after completed contract PR144/cd5d234. All71frozen manifest scenarios are now consumed by independent Go tests with exact grouped maps/evidence/request/status/context assertions and the complete immutable primary example. No oracle/source capture changed. Primary integrity test checks the original SHA256 and literal248/5/8/31premises.

Local Go1.26.9 on the recovered isolated Linux worktree compiles dependencies but the new tests stop at the expected missing symbols StorageSKUGroup, NewStorageSKUCollector and StorageSKURequest. This is missing-API test preparation, not product regression reproduction or acceptance. The complete log hash is 3f5cae8fda8db54e27aec7844e40c221cc018ef1d8ff693a106a7d050c9f537a. Both-host published preparation inspection is pending. No prior green certifies this intentionally uncompilable test checkpoint.

Entry cd5d2340d86e56e334410afca0aca1ba3f5fe0de/treee67cc2e9384ab82c73ba291f827f4db24e220806. Local checkout is /workspace/scratch/c04791396c67/Cloud-Assess-storage, separate from previous worktrees. Go1.26.9/cache/APRL60eaddda were verified; shell fetch succeeded. Current branch feat/storage-sku-collector-20261016. Human identity Denis Bogunic with recorded noreply address. No sub-agents.

Next: publish/read back this test preparation and inspect missing-API native evidence. Then implement under [the frozen contract](REGION_STORAGE_SKU_CONTRACT_20261016.md), retaining every manifest expectation; add independent boundary/ownership/cancellation/actual-client cases and meaningful compiling reversal/restoration controls, obtain complete exact-head native/source QA and stop before required-kind coordinator/public integration. Do not merge protected history or alter any original capture/dependency/pin/public schema.

Recovery: inspect actual live branch/PR state before retrying a mutation. The original contract fixtures and new tests must remain byte-identical through production unless a separately evidenced expectation defect is explicitly recorded before changing production. Targeted ordinary successor/revert commits preserve regression history; no reset/force/protected rewrite. Do not infer accepted target adapters2/7 from an unmerged proposal; accepted remains1/7 and source captures7/7. Azure/live, six adapters/integration/FullCost, load/freshOS/advisory/human/license/security/Gate004/release remain open.
