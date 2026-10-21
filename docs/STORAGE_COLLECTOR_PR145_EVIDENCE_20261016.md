# PR145 final evidence snapshot, retrieved 2026-10-16

The bounded, kind-aware Storage SKU collector checkpoint is complete as an offline draft candidate. The collector makes one scoped, bounded read, validates the response atomically, preserves name/kind groups and regional restrictions, and projects owned evidence. It does not enable the public region coordinator.

Candidate: `a14128bc94bf531dd1c28553db14d2be011b2214`
Tree: `488bd4099d39a9f17b2034ba499f8d73710315b8`
Parent: `c2d5c5df48fa0e4939b5b3a3b7b67f9bd391ca8c`
Merge preview: `eca7f72af6b7450ed84b0620c9aa14eeec0e133b`, identical tree; ordered parents are protected `a6fee427e19feeef28691737ea9bc5026086efa0` and this candidate. All three full job logs show that preview checkout. API-created commit identity, parent/tree, all 15 implementation files and all 11 correction files were read back and verified. Local metadata-equivalent commit markers are f9f0631 and be4b5c8; their trees match the published candidates.

| Exact-head gate | Run / job | Result |
| --- | --- | --- |
| Linux Go1.26.9 | [38907421820 / 116994897339](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38907421820/job/116994897339) | Passed, 64 successful steps; only failure-evidence upload skipped |
| Windows Go1.26.9 | [38907421820 / 116994897883](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38907421820/job/116994897883) | Passed, all 50 steps successful |
| Pinned source | [38907421620 / 116994406387](https://github.com/DeBoX85/Cloud-Assess/actions/runs/38907421620/job/116994406387) | Passed, all 7 steps successful |

Validation includes the unchanged 71-scenario manifest; immutable 248-row primary Storage example; exact limit and limit-plus-one cases; malformed and filtered metadata; actual bounded HTTP client with offline transport/credential; cancellation priority; ownership and concurrency; required-kind isolation; complete package/race suite, vet, coverage and bounded fuzzing. Linux total coverage is 85.4%, region 93.6%. Authored documentation checked 1,150 local-file links. Default and custom native package variants each passed 19 tests and installed offline CLI checks on both platforms; those installation cases were not skipped. Linux's existing native-Windows sharing-lock unit test is platform-skipped; Windows runs that case.

Control identity comparison retained every previous control: Linux 328 + 25 = 353; Windows 323 + 25 = 348. New controls all compiled, failed their named assertion, restored exact original source and passed the named baseline; final Storage baselines passed. The same 25 controls passed locally with LF and simulated CRLF, including byte-exact restoration.

Added control identities: request-version, getter-bound, pre-auth-scope, no-paging, json-admission, display-alias, stage-physical, regional-union, kind-grouping, missing-locations, restriction-scope, restriction-union, unknown-group-gating, unqualified-gating, required-kind, atomic-admission, status-error-pair, context-priority, body-bytes, row-limit, raw-work, projection-work, projection-text, capability-conflict, late-unknown-context.

Sorted unique control identity set SHA256 (UTF-8, newline after each identity):
- Linux baseline: `b3821c2a01de445b6e29b21b18f0ee77c1b2a548b5790f8071eaad91a19626ef`
- Linux candidate: `9e8ba7becd8a3853166301f03d2702aceada2429968637235ad0cf37e916b8f2`
- Windows baseline: `4981209416fbb0b86e1b510ba47431528c4e711405309ef9874fbee91e8d5a9f`
- Windows candidate: `42d908451d77cf582da574a22f21a3ab93b4b47325ba38387b7d6bdc1f273d22`

Full decoded job-log fingerprints preserve retrieved BOM and native newlines:

| Job | Characters | UTF-8 bytes | SHA256 |
| --- | ---: | ---: | --- |
| Linux | 106908 | 106912 | `9492a6865f4669534920abe5cf694dba46e4fc06e7e22c2ea1c74afe8ba64da6` |
| Windows | 86742 | 86744 | `045ff378a319718d3b976ec687828063440b9bda94674c0fd501e535dcaca932` |
| Source | 1359797 | 1359873 | `d1e9cd4818bcc70fa447d7a1d760e7306ab3999ec852a643a805e95cd07798f2` |

Source proof is complete: 11 capture harnesses, 23 files, 315 unique contiguous chunks, 909,865 reconstructed UTF-8 bytes. Every reconstructed file matched its declared SHA256 and a fresh complete file fetched from this candidate. Provenance remains AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`, tree `17d93b20c303f90f7843036be82f0dc32f3260f1`, APRL `60eaddda76541f6adbc1c5ffa686829807e55e29`. The frozen contract blob remains `11e7cfefa657a22d9f1bb933106ce5d0ed1c31ef`; primary example SHA256 remains `aceb9c06f233a6db9b4adb726e8f7951037531887dd28c949a24dee5b354830e`.

Both fresh vulnerability scans report zero reachable vulnerabilities and zero in imported packages, plus one module-only advisory. That advisory remains open. Offline automated QA and self-review are not independent-person review, live Azure validation or release approval.

Failure history is retained:
- Test preparation f0707bb / native38903459610 and 07eb37d / native38904120528 intentionally failed compilation because the new collector API did not yet exist. Their separate source runs passed; neither native preparation is credited as product acceptance.
- QA-ST-01: an additional local review test reproduced a late unknown-result cancellation gap; the collector was corrected before publication, and the independent regression plus compiling reversal passed. This is a new candidate implementation gap, not a new baseline audit finding or oracle rewrite.
- QA-ST-02: initial local helper anchor used the wrong gofmt spacing; it stopped and restored source. Corrected complete control rerun passed.
- QA-ST-03: c2d5c5d failed Windows job116989921866 in run38906095075 after six controls because CRLF checkout bytes did not match LF multiline anchors. Source was restored. Full failed log: 36,704 characters / 36,706 UTF-8 bytes, SHA256 `8bf7ffc77d1bcfdc47a788cf544e8189481f4141ad3feaf388df80aaa0380337`. Superseded Linux job116989922179 was cancelled by workflow concurrency when this correction was pushed. Its source run38906095088 passed separately but is not this candidate's evidence.
- The correction normalizes in-memory anchor text only and restores exact original checkout bytes. Production and frozen expectations are unchanged. Both local LF/CRLF logs: 2,215 bytes, SHA256 `ceec94ce21e61bf5b791b65ed735599def5ba85e2df3cce8e2b5135f45fd4b4f`. The initial CRLF launcher omitted pinned Go PATH and ran no tests; the corrected launcher passed all 25 controls.

[Collector checkpoint](https://github.com/DeBoX85/Cloud-Assess/blob/a14128bc94bf531dd1c28553db14d2be011b2214/docs/REGION_STORAGE_SKU_COLLECTOR_20261016.md), [frozen contract](https://github.com/DeBoX85/Cloud-Assess/blob/a14128bc94bf531dd1c28553db14d2be011b2214/docs/REGION_STORAGE_SKU_CONTRACT_20261016.md), and the session handover/recovery records contain the implementation contract, local proof, failure history and stop boundary. Their pending-at-publication wording is historical; this body is the final exact-head evidence index. No evidence-only commit was added after these gates.

STOPPED at the requested natural checkpoint. PR141-145 remain draft, open and unmerged. Protected bootstrap/core-v1 is unchanged; accepted target adapters remain 1/7, with this proposed Storage collector still unaccepted. Source characterization covers 7/7 families. No Azure operation, dependency/SDK/pin/source-capture/public-schema refresh, coordinator integration, protected merge or release was performed. The workflow adds one mandatory 25-control step per platform and five minutes of job budget for that added work; existing checks, floors and fuzz budgets remain intact.

The next distinct slice is the VM/VMSS pre-edit contract using REGION_COMPUTE_SKU_SOURCE_20261013.md and REGION_SKU_TARGET_ANALYSIS_20261013.md, with immutable primary examples and independent expectations before implementation. Required-kind coordinator and public/all-format execution remain deferred. Remaining migration families, FullCost and prior roadmap/live/security/license/release obligations remain open.
