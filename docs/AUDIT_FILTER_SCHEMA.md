# Audit: filter schema admission

## Contract and suspected scope gap

TARGET_SPECIFICATION.md defines the neutral assessment root, include.subscriptions and exclude.resources; legacy azqr/exclude.services compatibility is intentionally not required. Filtering must represent the operator's intended scope rather than silently ignore a misspelled or unsupported scope field.

Accepted07011b63440e69f4a5acb128e0449943f759ee9c internal/config/load.go uses yaml.Unmarshal without unknown-field or document-count admission and permits a null assessment pointer. Source inspection suggests invalid/legacy filters can be treated as defaults while assessment startup proceeds. This is an unconfirmed execution hypothesis until the independent native regression fails.

## Independent test-only checkpoint

cmd/cloud-assess/filter_schema_audit_test.go passes five literal files to the real executeScanWithFactories path. Legacy root, misspelled include key, legacy exclusion key, null assessment and second document must fail as filter input before credential/operations factories. Healthy empty file, empty neutral assessment and known include/exclude syntax must reach the deliberately failing credential boundary once. Plugin-only zone selection isolates filter admission from unrelated YAML extension discovery. No credential created, Azure transport called or report produced.

This checkpoint preserves accepted production. Test-only native compiling assertion failures on Linux/Windows are required before classifying AUD007. It does not prove historical live overscan. Expected narrow correction after reproduction: known fields, one document, non-null assessment; preserve valid empty files/defaults and all supported filter semantics. No source-pin/dependency changes or scan-scope expansion.

## Recovery

Draft WIP, unaccepted/unmerged, whole-project audit in PR114. Verify actual refs/head/tree/identity and logs after interruption. Dependent acceptance waits for reproduction, corrected exact-head full native/source QA and final semantic review. No release/live/Gate004 claim.
