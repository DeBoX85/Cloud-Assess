# MVP acceptance and parity matrix, 2026-10-19

Current status: [CURRENT_STATE](CURRENT_STATE.md). Authorisation and reevaluation: [programme](MVP_PROGRAMME_20261019.md). This is a requirement-to-evidence map, not a percentage or release approval. Code inspected at product PR170/8c47d810 and pinned AZQR8e4f0577. Protected baselinea6fee427 remains unchanged. A path establishes implementation presence, not correctness or live certification.

## Acceptance sets

The existing [target specification](TARGET_SPECIFICATION.md#core-v1-scope) defines core-v1 material assessment equivalence with explicit health/completeness and adjustable branding. It explicitly defers full internal-plugin parity and MCP/compare/alternative-VM utilities until after core-v1 equivalence. The user's wider goal remains pinned-AZQR functional equivalence. Track both without deleting deferred requirements or assuming a smaller MVP is approved. No automatic CLI spelling/config/package identity compatibility is required by the specification.

| Requirement | Present implementation / evidence pointer | Remaining acceptance work | Can advance offline? |
| --- | --- | --- | --- |
| Authentication, cloud endpoints and scope discovery | internal/azure, internal/scope; access/scope contracts; historical equivalence | Representative credential/cloud failures and independently expected scope visibility, nested MG DV-001 | Regression/procedure yes; real access no |
| Inventory, filtering, catalog, YAML/KQL and scanner pruning | internal/graph, internal/rules, internal/scanners; core command/tests | Required actual-input equivalence and outstanding lifecycle/load scenarios | Synthetic integration yes; live result no |
| Advisor, diagnostics, Defender plans/recommendations, Policy, Arc SQL, core Cost | Dedicated core packages, stage orchestration, retained audit corrections | Non-empty optional-stage datasets, Arc SQL shape and exact live completeness; unresolved historical warning correlation | Reconciliation and controlled response tests yes |
| Canonical JSON, XLSX, CSV, SARIF, stdout and severity gate | internal/assessment, internal/renderers, cmd/cloud-assess; PR170 renderer regressions | Whole supported workflow exit/health/output agreement and live paired comparison | Yes, with existing format-specific contracts |
| Adjustable build branding and installed packages | internal/branding, tools/brand-build; native package checks | Fresh supported OS, operational rollback, publisher/release acceptance | Package rehearsals yes; independent release approval no |
| Reliability/security fixes | Four-section review and retained regressions; PR170 scope masking | Independent review, open module-only advisory disposition and required release security/license gates | Investigation yes, no unapproved refresh/suppression |
| Zone mapping, service health, SQL EOL, carbon, AI governance | Five registry entries and standalone commands in cmd/cloud-assess | Representative live plugin behaviour and access where required | Integration tests yes; tenant proof no |
| Region selection sixth public plugin | Pure calculations/outputs and bounded libraries in internal/plugins/region; registry intentionally excludes it | Missing cost collection, run coordinator, full evidence propagation and public/report integration | Yes, via bounded contracts and exact-head QA |
| Region SKU families | Disk/Storage/VM/Cognitive/SQL DB/SQL MI collectors plus VMSS catalog diagnostics | Correct per-family integration; no inferred SQL coordinates or VMSS eligibility; protected adapter acceptance remains1/7 | Coordination and justified metadata handling yes |
| Region cost parity | ProjectCost and EnrichCost exist; pinned region/cost source has historical query + retail phases | Bound history/retail requests, consistent currency/unit/tier/identity/completeness, then FullCost output reconciliation | Yes; do not confuse with existing core resource-cost scanner |
| Ancillary source utilities | rules/plugins/scanner commands present; source types utility needs target equivalence check | types behaviour; MCP/compare/alternative VM remain explicitly deferred | types reconciliation yes; no speculative deferred expansion |
| Operational release | DEFERRED_VALIDATION, QUALITY_GATE_004_PLAN and safeguards matrix | Independent review/backup, live scale/permissions, fresh OS, installed rollback, legal/security/Gate004/release owner | Prepare/rehearse only within access boundaries |

## Implemented region inventory, not an accepted-adapter percentage

| Family | Actual product path | Integration limitation |
| --- | --- | --- |
| Compute disks | compute_disk_sku_collector.go | Existing complete/unknown evidence projection; coordinator still absent |
| Storage | storage_sku_collector.go | Keep kind-aware keys and complete/unknown semantics |
| Compute VM | compute_vm_sku_collector.go | SKU catalog evidence is not allocation/capacity proof |
| Compute VMSS | vmss_catalog.go | Diagnostic catalog only; deployment scoring unsupported |
| Cognitive | cognitive_sku_collector.go, cognitive_sku_catalog.go | Kind-bound request and catalog; preserve kind identity |
| SQL Database | sql_capability_collector.go, sql_capability_catalog.go | Full coordinates; inventory raw SKU name is insufficient |
| SQL Managed Instance | sql_managed_instance_collector.go, sql_capability_catalog.go | Explicit family SKU and coordinates; no inferred mode/eligibility |

Provider locations, quota and reservations collectors and pure availability/latency/cost/quota/reservation calculations also exist. SQL collector-to-report diagnostics are complete offline in PR170. Do not count them as public region support or create replacement collectors to satisfy a stale task list.

## Evidence limits and closure test

PR170 final recorded551Linux/546Windows controls, coverage87.2%,23captures and exact-head native/source certification remain its own evidence. This matrix has not rerun or re-reviewed those full logs. New candidates require their own mandated validation. Live run conclusions were checked during task1, with both success; that is not equivalent to full current log review.

An MVP readiness decision must state which acceptance set is met, exact supported user path, required same-input parity evidence, explicit completeness/health, supported installation, unresolved findings and required human decisions. Do not declare readiness from test counts, coverage, source-characterisation counts or draft publication. No reliable overall percentage/calendar ETA follows from these unweighted rows.
