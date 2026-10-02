# Bounded YAML Graph plugins

Date: 2026-10-02 (Europe/Oslo). B2 starting baseline is PR #78 merge `50fbaee905e7be7a7c99e174f35ea69247eb16bb`, reviewed tree `19da2f901ccd93af27add84da1b93bb13d12a4cc`. This pre-implementation contract follows [DEVELOPMENT_EXECUTION_PLAN.md](DEVELOPMENT_EXECUTION_PLAN.md). It concerns YAML recommendation queries, not compiled internal table plugins or shell/subprocess execution.

## Source and target decision

Pinned AZQR `8e4f0577f3615e6c9014c031bcad079f235369cc`: `internal/plugins/interface.go`, `yaml.go`, `loader.go`, `internal/pipeline/stage_graph_scan.go`, and `internal/graph/graph_scanner.go` define the contract. Preserve object-schema metadata/query mapping, version default 1.0.0, queryFile precedence over inline query, sorted recursive case-sensitive yaml/yml discovery, first plugin name found across home then current directories, plugin-name-sorted recommendation overrides by resource type and ID, attribution to configured name, normal Graph enablement and recommendation exclusions. External queries do not use embedded disabled/manual/development predicates. Unsupported scanner resource types remain unexecuted as in the source. Existing safe projection/health/redaction apply.

Default target directories are `$HOME/.<immutable-cli-name>/plugins` then `./plugins`; Windows home is the standard Go user-home location. If the standard home cannot be resolved, only the current-directory root is used; no relative substitute home is invented and existing service-identity scans do not gain a home-directory prerequisite. Only an actual scan discovers files, before credential creation. Rules/help/branding remain offline and do not load plugins. Branding directory selection uses the already accepted CLI-name profile. Ordinary scans with no plugin directories retain existing behavior. Operators must control and review both directories and their KQL: discovery intentionally causes queries to execute against the selected Azure scope. This is not an untrusted-code sandbox or a least-privilege certification.

Deliberate target protections: candidate errors fail preflight rather than source debug-log-and-skip; missing directories alone are skipped. Non-directory roots, unreadable paths, malformed/unknown/duplicate fields, aliases/merge keys, multiple documents and invalid required values fail explicitly without echoing YAML/query payloads. Plugin names are bounded UTF-8 labels (128 bytes maximum, no surrounding whitespace, controls or path separators) and cannot impersonate APRL/AOR/CUSTOM/DIAGNOSTICS or the six reserved internal plugin names. Resource type must be nonempty without surrounding whitespace, avoiding a target-only expansion of otherwise inactive source type matching. Required IDs/descriptions/query text cannot be whitespace-only. Boolean fields and optional string recommendationTypeId match the source schema. The pinned loader accepts recommendationTypeId but does not assign it to the converted definition; the target likewise leaves that projected field empty. Duplicate query IDs within one resource type retain the source last-query-wins behavior; duplicates across sorted plugin names retain the later-name-wins behavior. Duplicate plugin names retain first-discovered-wins.

Local inputs are bounded: 1 MiB per YAML/query file, 16 MiB cumulative bytes per discovery, 256 candidate files, 4,096 query definitions and 8,192 visited entries, depth at most 16. These are explicit initial operator limits, not measured Azure or arbitrary-hostile-filesystem performance guarantees. Crossing a limit fails rather than silently truncating. Candidate/query files must be regular files, with no final symlink. Query files are portable relative paths confined to the declaring YAML's directory, including nested subdirectories, never parent/absolute/drive/UNC/device paths, control/reserved characters and device names including superscript variants. Query paths are at most 4,096 bytes and 16 components. `os.Root` provides containment at open time, including intermediate symlink escape defense. Confined intermediate links may be followed; directory walking does not descend symlink entries. Roots themselves are operator-selected directories. Privileged mount manipulation, malicious concurrent filesystem changes and special-device availability are outside this trusted-directory contract.

## Implementation and acceptance boundaries

Implement distinct bounded parser/discovery in rules, leaving recommendation-array loading unchanged. Overlay parsed definitions on a fresh catalog per run and retain external origin outside the public assessment schema. CLI preflight passes only verified definitions into existing orchestration; ordinary scan stages and report schemas remain unchanged. No source pin/dependency/normalization updates or Azure operations are required for development. Do not advertise missing internal plugins as executable.

Acceptance: literal full-field source mapping, inline/queryFile preservation and precedence; first-name versus sorted override and same-type-ID collisions; unsupported/scanner/metadata/exclusion behavior; bounded/ambiguous/malformed input and traversal/symlink controls; no caller/catalog/state leakage; actual CLI rejection before observed authentication with unchanged report files; healthy synthetic YAML-to-ARG-to-canonical-report plus failure health; unchanged offline rules capture and native default/custom packages. Critical mapping/containment changes need compiling negative controls. Required quality and Windows jobs must pass the exact final head and logs must be inspected before merge. Evidence distinguishes synthetic transport from live Azure.

## Recovery and next task

Rollback is a reviewed revert PR on current core-v1 after checking dependent work; preserve evidence and rerun required QA. GitHub response loss is reconciled by remote head/tree/PR state before retrying. Last durable baseline remains PR #78 until this proposal is published. If blocked, continue independently testable parsing/contracts without claiming production acceptance. After B2 acceptance, resume B3 canonical internal table/health and zone-mapping. Laptop/Azure validation remains deferred, including DV-001 and live nonempty datasets.

Primary containment rationale: [Go traversal-resistant file APIs](https://go.dev/blog/osroot), [YAML 1.2.2 mapping-key uniqueness](https://yaml.org/spec/1.2.2/), [Windows file/device naming](https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file). No comprehensive quality guarantee or live plugin-equivalence claim is made.

## Reference capture

The actual unchanged pinned loader was executed with the literal full-field test fixture in a separate temporary source checkout. `internal/rules/testdata/yaml-source-conversion.json` retains its complete metadata and recommendation output, with only the temporary CommandPath replaced by `<fixture-path>` for reproducibility. Capture SHA-256 `dbd0ba34d2aff68c8ff8dd28eb353e2906085fa43821bea24becd3188856772c`. Full mapped output is compared semantically in tests, alongside independent literal expectations; public target Graph validation metadata is additive. Source reference checkout remains unchanged.

## Operator example

Create a reviewed YAML file in the selected home plugin directory or `./plugins`:

```yaml
name: operator-controls
queries:
  - aprlGuid: operator-storage-control
    description: Review the selected storage resources
    recommendationResourceType: Microsoft.Storage/storageAccounts
    recommendationControl: Security
    recommendationImpact: Medium
    query: |
      resources
      | where type =~ 'Microsoft.Storage/storageAccounts'
      | project id, subscriptionId, resourceGroup, name, type, location
```

An ordinary storage scan includes this query automatically. It intentionally reports every storage row returned by this illustrative query; replace it with the actual reviewed control predicate. A `queryFile: kql/control.kql` field overrides the inline query and must refer to a confined regular UTF-8 file. Remove files from the discovery directories to disable them; the separate internal-plugin stage flag does not control YAML Graph queries. Preserve the original AZQR library if migrating configurations, review inputs and place copies in the neutral branded directory. The `rules` command remains the pinned catalog/Diagnostics inspection surface and does not list these external definitions.
