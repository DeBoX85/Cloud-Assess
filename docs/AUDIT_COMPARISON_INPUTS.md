# Audit: comparison input evidence

## Contract and suspected gap

The canonical target contract is internal/result/result.go schema 1.0; assessment completeness has four declared values. docs/EQUIVALENCE.md rejects partial/failed baselines and permits complete_with_warnings with explicit review. Semantic equivalence is an evidence statement, so an unknown schema or absent/unknown completeness cannot prove a healthy assessment. Pinned-reference JSON has recognized core datasets. A reference with no recognized sections cannot supply any comparison evidence.

Accepted baseline 07011b63440e69f4a5acb128e0449943f759ee9c only requires nonblank schemaVersion; its completeness switch does not reject unknown values. LoadReference accepts null, empty objects and unknown-only sections. AUD006 CONFIRMED with six compiling real-command assertion failures on both native platforms, with unchanged production.

## Independent regression and scope

tools/equivalence/input_evidence_test.go runs the real command function with literal independently specified input files. Six invalid cases require nonzero exit: absent/unknown target completeness, unsupported target schema, null/empty/unknown-only reference. Two healthy empty Advisor controls with complete and complete_with_warnings require exit 0. Existing nonempty semantic/coverage/partial/redaction cases remain.

Test-only4cb8bd8451208e45e2551f72edef280de338e52f/tree3ae8c71f22fcb86fbe951810bd997fc7afd5db76/parent07011b63 reproduced all six exit-0 faults in native37185359814 on Linux111386033586 and Windows111386033687. Formatting/compilation passed, only the six named assertions failed; healthy empty controls did not fail. Complete logs inspected. These malformed files are synthetic, not historical live capture incidence. No Azure calls, source pins, dependency versions or normalization changes.

The correction is narrow loader validation in internal/equivalence/projection.go: enforce the supported target schema and known completeness; reject reference reports with no recognized dataset. Genuine empty arrays remain valid evidence. Existing partial/failed comparisons remain non-comparable, warning review remains required. No claim of a full JSON schema validator, identity/snapshot proof, live equivalence or validation of arbitrary stage contradictions.

## Recovery and gates

Draft WIP, unaccepted and unmerged. Verify live refs and exact-head native/source evidence before acceptance. Audit PR114 contains the current project checkpoint; feature acceptance PR113 remains paused. Previous live comparison evidence is historical and is not invalidated without evidence that these malformed inputs were used.

## Corrected checkpoint

The loader rejects missing/unsupported target schema versions, missing/unknown completeness values, and reference inputs with zero recognized core sections. Both allowed healthy completeness values and valid empty/null known sections retain their established behavior; partial/failed remain non-comparable. This deliberately closes an offline QA false-success boundary without changing assessment scanning, source pins or normalization. Full fresh Linux/Windows/source gates and tested-preview identity must be verified for the corrected head before acceptance. New supported schemas require an explicit reviewed loader update.

## Review correction: additive schema compatibility

Before acceptance, requirements reconciliation identified a valid second canonical schema: result.PluginSchemaVersion1.1, additive plugin tables under docs/PLUGIN_TABLES.md. The first unmerged733b9243 guard incorrectly admitted only1.0. A literal nonempty valid plugin-table command fixture now requires successful core-only comparison for1.1 against healthy empty Advisor evidence. This checkpoint changes tests/documentation only on733b9243 to obtain an independent compiling compatibility failure before correcting the guard. No accepted branch affected, no historical incidence. Future correction must allow both declared schemas without claiming generic core comparison compares plugin cells.
