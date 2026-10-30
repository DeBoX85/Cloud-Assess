# Comparison JSON admission contract

Product successor `fix/comparison-json-admission-20261027` starts at verified PR204/fd746cbf. PR204 own quality40531800632/source40531800701 succeeded; archive41e2d576/receiptd2061659 fully read back and44members/6558197bytes restored equal originals. This contract and literal public-loader/CLI tests must be published and read back BEFORE production. No implementation credit is claimed here.

## Required behavior

Both LoadReference and LoadTarget reject repeated decoded JSON object member names anywhere in the document, including ignored sections, arrays, unknown nested objects and opaque maps. Exact and escaped spellings of the same decoded name are duplicates even when values agree. Reject before projection/comparison; CLI exit2 leaves an existing output unchanged. Errors identify the admission failure without echoing member names or values. Existing contradictory-record exit1 versus invalid-input exit2 remains distinct.

Target struct objects also reject two names which resolve to the same current typed field under pinned Go1.26.10 encoding/json v1 exact-preferred Unicode case-insensitive matching. This applies to nested structs, pointers and array elements, including scope/status/warnings/plugin metadata. Preserve one case/Unicode alias alone. Map keys and unknown names remain case-sensitive; do not fold all JSON names. Reference sections/row maps, target tags, summary category/impact maps and unknown payloads retain distinct case variants. Repeated keys in separate objects are independent. Normalized-identical array records and documented count aggregation retain their existing behavior.

No new whole-document byte, row, token or nesting limit is introduced by this narrow integrity correction. Existing io.ReadAll allocation, pinned JSON syntax/depth/type behavior and downstream schema/completeness checks remain; this is an explicit compatibility decision and leaves memory/input-budget hardening for separately justified workload/task13 evidence. Do not claim bounded total input memory, generic HTTP protection or a service security boundary. Preserve large numeric literals in ignored data without converting them to float64 merely to scan names. Keep trailing/malformed-input rejection and ordinary valid nested inputs.

## Evidence and rationale

RFC8259 sections4/8.3 recommend unique object names and compare their decoded Unicode representation; rejecting them here is a deliberate evidence-integrity correction, not a claim every duplicate violates JSON grammar. https://www.rfc-editor.org/rfc/rfc8259.html

The verified pinned toolchain's src/encoding/json/decode.go documents ordered duplicate replacement/merge and case-insensitive struct matching; fold.go binds its Unicode folding to EqualFold. Current package documentation is supplementary and may describe a later runtime: https://pkg.go.dev/encoding/json . No json/v2, compiler, SDK, dependency, source/action pin or oracle update is authorized. The pinned source is authoritative for this execution.

## Bounded sequence and exit conditions

1. Verify204 own publication/QA/retention, inventory worktrees and isolate this PRODUCT successor.
2. Freeze/read back this contract and independent literal expectations; reproduce the named rejection assertions with healthy opaque-map/single-alias controls passing BEFORE production. No generic compiler/panic failure credit.
3. Implement the smallest pre-decode admission guard; preserve schema/projection/opaque-map behavior and add strict compiling assertion reversals on both native hosts.
4. Own complete native/source QA, original controls and all additive proofs; preserve failed history and open advisories.
5. Coherent recorded204 rollback plus private latest-base additive expected-head capture, complete readback and fresh restore; retain until explicit Denis deletion approval.
6. Update PROJECT_PROGRESS/current/register/ledger/authority/restart, publish/read back, self-assess and continue task10 only after this candidate's own QA/retention.

Current state: task1 verified, task2 expectations prepared pending publication/reproduction; tasks3–6 pending. All32frozen files/44legacy/49current obligations and original controls/pins remain. Issue189 closed bounded migration, continuing retention/revisit duties. No Azure/live retail/real credentials/fixtures/RBAC/subagents/protected merge/release. Independent/live/operational/licensing/Gate004/release and wider parity/documentary obligations remain; human security assessment considered only after overall project completion.

## Implementation checkpoint

Contract/tests frozen and all12stage files fully read back at5f1ad72aec12ae062ec6f2a0309d48be67f0ed8e/treec1410457fcc1411c19f2f301f9074125d903b34e BEFORE production. Actual compiler0, healthy0/10passes and21named rejection assertions/exit1 independently verified; first missing-APRL compile was setup failure only. Frozen two test files remain byte-identical.

Iterative token admission is before both existing unmarshals; typed alias matching follows current named DTO fields and fails closed on unsupported anonymous-field shapes. Opaque maps retain exact decoded keys. UseNumber preserves ignored large numeric literals. Focused suites and nine compiling source-assertion reversals/restored baselines/hygiene pass locally. First driver invocation violated shared helper's one-package requirement; second broad subtest selector matched a sibling and strict witness admission correctly rejected it. Corrected driver anchors each path segment and verifies actual healthy baseline root tests, preserving helper/expectations. Both native proofs are additive before full suites;21workflow/sevenhandover guards pass. Own complete native/source QA, current coverage/advisories/original control conservation, coherent rollback and private retention pending; tasks1–3locally verified,4–6pending. No whole-memory bound or wider assurance claim.
