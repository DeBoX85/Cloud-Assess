# Audit: filter schema admission

## Confirmed AUD007 and contract

TARGET_SPECIFICATION.md defines the neutral assessment root, include.subscriptions and exclude.resources; legacy azqr/exclude.services compatibility is intentionally not required. Filtering must represent the operator's intended scope rather than silently ignore a misspelled or unsupported scope field.

Accepted 07011b63440e69f4a5acb128e0449943f759ee9c internal/config/load.go used yaml.Unmarshal without unknown-field or document-count admission. RebuildIndexes restored a null assessment to default configuration. Five literal invalid files reached credential creation with their restrictions discarded in the real executeScanWithFactories path. This establishes configuration admission failure, not historical live overscan.

## Independent negative control

Test-only commit 926ac5cedf4432b58f30a5ea76c744eb5f2ab247, tree 0a89f834a3ff834a82e20296da2e6b78b64fc124, preserved accepted production. Native run 37186278979 compiled/formatted successfully, then failed all five named assertions on Linux job 111388751676 and Windows job 111388751827.

cmd/cloud-assess/filter_schema_audit_test.go passes legacy root, misspelled include key, legacy exclusion key, null assessment and a second document to actual preflight. Each must fail as filter input before credential/operations factories. Healthy empty file, empty neutral assessment and known include/exclude syntax reached the deliberately failing credential boundary once and did not fail their controls. Plugin-only zone selection isolates filter admission from unrelated YAML extension discovery. No credential created, Azure transport called or report produced.

## Narrow correction

LoadFilters now uses yaml.Decoder.KnownFields, requires a non-null assessment and requires EOF after the first document. Empty files retain initialized defaults. Supported nested filter keys, indexes, include/exclude precedence, validation and no-file defaults remain unchanged. Unknown legacy/misspelled keys and additional documents fail before authentication. This deliberate target configuration correction does not alter pins, dependencies or authorized scan scope.

The independent regression remains in the proposal. Full exact-head Linux/Windows and source characterization gates are pending; no acceptance, release or live equivalence claim.

## Recovery

Draft WIP, unaccepted/unmerged, whole-project audit in PR114. Verify refs/head/tree/identity and full logs after interruption. Acceptance waits for corrected exact-head full native/source QA and final semantic review. No release/live/Gate004 claim.
