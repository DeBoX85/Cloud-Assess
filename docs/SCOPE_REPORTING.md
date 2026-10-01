# Requested and resolved scope evidence

Core-v1 canonical JSON now has an additive `scope` object. This records selection/discovery evidence; it is not an independent subscription-membership, RBAC or resource-completeness audit. Older results/library fixtures may omit the object. The existing schemaVersion and scopeId hashing remain unchanged; strict downstream schemas must allow the new field. CSV/XLSX layouts are unchanged; inspect canonical JSON for this evidence.

| Field | Meaning |
|---|---|
| `selection` | `subscriptions`, `management_groups` or `accessible_subscriptions`, based on explicit assessment inputs |
| `status` | `not_completed` before/after unsuccessful discovery; `resolved` after successful discovery with no missing explicit subscription; `unresolved` when an explicit subscription did not enter the resolved set |
| `requestedSubscriptionIds` | Explicit CLI/library subscription inputs, trimmed, case-normalized, deduplicated and sorted |
| `requestedManagementGroups` | Explicit management-group inputs, normalized in the report; actual discovery still receives the original inputs |
| `includedSubscriptionIds`, `excludedSubscriptionIds` | Effective subscription filter configuration, including explicit CLI subscriptions appended to include filters; these lists are selection inputs, not discovered estate membership |
| `resolvedSubscriptions` | Discovered active subscriptions after selection/filtering, each with `subscriptionId` and `subscriptionName`; includes subscriptions with zero resource/finding rows |
| `unresolvedSubscriptionIds` | Explicit requested subscriptions absent from the successfully resolved set; empty when listing itself failed because absence cannot then be classified |

## Explicit requests fail closed

The coordinator treats any unresolved explicit subscription as a critical scope failure before inventory, Graph, Diagnostics or optional-stage resource queries. The scope stage is `failed`, completeness is `failed`, the application returns exit 1, and requested reports are persisted when rendering succeeds. The error contains `scope_requested_subscription_unresolved` and a count, not raw subscription IDs. The ordinary stage error code remains `stage_failed`; `scope.status` and the unresolved list provide the structured classification. Blank explicit subscription/management-group identifiers reject before discovery.

An unresolved ID can result from visibility, tenant/session, disabled/deleted state or selection. No specific authorization diagnosis is inferred. This deliberately corrects the previous pinned-source-compatible behavior that silently scanned only the visible subset of an explicit request. Underlying discovery helpers retain their characterized intersection behavior; the assessment coordinator enforces the product boundary. Recommendation projection, library pins and comparator normalizations are unchanged.

## Filters, empty scopes and hierarchy limits

Explicit CLI subscriptions are also include-filter inputs. Existing include-over-exclude precedence is retained. Filter-only includes/excludes are selectors, not an independent expected membership list: a successful filtered empty set remains valid, even if an include ID was absent. Use explicit subscription inputs when each intended ID must resolve. Resource-group requests already require one explicit subscription and use the same check.

All-accessible and management-group discovery retain successful empty-result behavior. `resolved` means discovery completed under the selected identity and filters, not that hidden subscriptions/resources do not exist. A management-group failure retains `not_completed`; a successful empty hierarchy response cannot independently prove intended membership. DV-001 nested live traversal and restricted-identity evidence remain open. ARG can still omit inaccessible resources within resolved subscriptions. Compare the JSON scope with independently approved membership, including zero-data subscriptions.

Microsoft's [ARG permissions guidance](https://learn.microsoft.com/en-us/azure/governance/resource-graph/overview#permissions-in-azure-resource-graph), rechecked 2026-10-01, documents that accessible-subscription results can be returned without a partial-result indication. This motivates independent access reconciliation; the toolkit's explicit-request failure policy is a product decision, not an Azure service guarantee.

## Privacy and validation

JSON redaction collects IDs from all subscription scope lists and resolved entries, including IDs appearing nowhere in resource rows. The CLI redaction setting applies to persisted JSON and JSON stdout; private unredacted evidence remains necessary for exact ID reconciliation. Display names, management-group names, tags, raw service errors and SARIF remain sensitive. This is not general anonymization.

Offline regressions cover partial/all missing explicit IDs, disabled/deleted IDs, valid zero-data scopes, case/deduplication, overlapping include/exclude precedence, intentional exclusion and filter-only missing includes, empty management-group discovery, denied listing, blank inputs, owned/sorted canonical arrays, scope-only ID redaction and application exit/report persistence without resource queries. Existing empty-scope and live semantic evidence remain bounded by their original inputs. Exact final Linux/Windows CI must pass before merge. No Azure call or role change is required for these tests; no new live-equivalence or release PASS is claimed.
