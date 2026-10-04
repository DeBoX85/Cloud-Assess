# Audit: comparison input evidence

## Contract and suspected gap

The canonical target contract is internal/result/result.go schema 1.0; assessment completeness has four declared values. docs/EQUIVALENCE.md rejects partial/failed baselines and permits complete_with_warnings with explicit review. Semantic equivalence is an evidence statement, so an unknown schema or absent/unknown completeness cannot prove a healthy assessment. Pinned-reference JSON has recognized core datasets. A reference with no recognized sections cannot supply any comparison evidence.

Accepted baseline 07011b63440e69f4a5acb128e0449943f759ee9c only requires nonblank schemaVersion; its completeness switch does not reject unknown values. LoadReference accepts null, empty objects and unknown-only sections. These are source-inspection hypotheses until the native regression reproduces an exit-0 result.

## Independent regression and scope

tools/equivalence/input_evidence_test.go runs the real command function with literal independently specified input files. Six invalid cases require nonzero exit: absent/unknown target completeness, unsupported target schema, null/empty/unknown-only reference. Two healthy empty Advisor controls with complete and complete_with_warnings require exit 0. Existing nonempty semantic/coverage/partial/redaction cases remain.

This first checkpoint changes tests/documentation only, preserving accepted production. Native Linux/Windows compiling assertion failures are required before classifying AUD006 as reproduced. No Azure calls, source pins, dependency versions or normalization changes.

After reproduction, proposed scope is narrow loader validation: enforce the supported target schema and known completeness; reject reference reports with no recognized dataset. Genuine empty arrays remain valid evidence. Existing partial/failed comparisons remain non-comparable, warning review remains required. No claim of a full JSON schema validator, identity/snapshot proof, live equivalence or validation of arbitrary stage contradictions.

## Recovery and gates

Draft WIP, unaccepted and unmerged. Verify live refs and exact-head native/source evidence before acceptance. Audit PR114 contains the current project checkpoint; feature acceptance PR113 remains paused. Previous live comparison evidence is historical and is not invalidated without evidence that these malformed inputs were used.
