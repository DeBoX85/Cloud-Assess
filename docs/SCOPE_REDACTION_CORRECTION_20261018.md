# Assessment Status scope redaction correction, 2026-10-18

Bounded correction discovered during authorized internal SQL report integration. Original SQL four-task batch is complete; report integration remains task3 in progress. This branch starts from verified contract PR166/054aba88 and intentionally contains no report APIs or consuming runner. This is the proposed stronger report fallback base, pending own full QA, not protected acceptance.

## Evidence before production

Report WIP PR168/57af1f64 retained the original renderer. BOTH actual native failures were inspected and published/read back in PR168 BEFORE renderer production edits: quality39162596987 Windows117778845796 (36929UTF8bytes SHA256401278962902f35969760b9deb8a312c4c7189f8f3a55cceca846cfda2990b96) and Linux117778845844 (51434bytes SHAcdb6fad13f67a2964c33e4806b979d6064de65d22dcb633a6189e3c9a8f3266e). Named independent Assessment Status scope regressions and shared/actual CSV/actual XLSX SQL tests failed, without compilation or race failures. Full logs and pre-edit proof are retained for the evidence bundle. Source-only success cannot certify this WIP.

## Correction and independent expectations

The existing core Assessment Status ScopeID cell bypassed selected subscription masking. Apply the existing plugin masker plus a scope-specific fallback to every fresh core/plugin status row. The fallback retains only UUID candidates already present in ScopeID, recognizes semantic subscription context case-insensitively, and intersects candidates with canonical scope/core record subscription IDs. Preserve unrelated UUIDs, disabled-redaction values, schema/columns and canonical result bytes. Core-only reports do not acquire the plugin4096identity cap or an estate-sized identity map; existing plugin bounds remain unchanged. This is a precise renderer correction, not a universal anonymization or scale claim.

Independent regressions cover core-only/mixed reports; path/case/bare/embedded known IDs; unrelated UUID retention; contextual scope without discovery;5001subscription core selection; and all10canonical core record identity sources. The first regression file is byte-identical to the pre-fix published PR168 file. Three frozen compiling reversals cover status call, context inference and scope-known identity. Every reversal must fail a named assertion, restore exact original bytes and pass its baseline; full restored scope tests run last. Existing533Linux/528Windows identities are conserved, three additive controls yield536/531.

## QA and transition

Local tests/control outcomes and full own8workers/bothgates/source/complete logs/capture/control/archive/readback remain revision-specific. At publication QA is pending. Do not certify reportWIP using this branch. After this correction is verified and its evidence/handover published, build the final report candidate on this corrected base, preserving first report runner and frozen oracles. Amend the report implementation rollback target explicitly: atomic report+runner reversal retains this independent fix, rather than reintroducing the known privacy regression by returning to raw contract63545e0c. Original contract/head/frozen expectations remain preserved.

No Azure/input required. No fixtures/RBAC/dependency/action/SDK/source/capture/public-schema refresh, subagents, protected merge/bypass/release. All prior heads/unrelated work retained. Independent review/live/security/license/release/module-only advisory obligations remain open.
