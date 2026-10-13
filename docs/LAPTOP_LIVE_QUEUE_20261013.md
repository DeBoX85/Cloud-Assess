# Laptop restoration and read-only live validation queue, 2026-10-13

Status: access expected within one hour in the user's 2026-10-13 message, NOT yet
confirmed. This queue is preparation, not an executed scan or restored credential.
Accepted development is PR137 af18468fe1dfc697838115c0a7585ba675fdec97 at preparation;
reverify the remote branch before use. PR138 is draft source characterization,
not a live candidate. Public region-selection is still unavailable.

## First return: inspect without changing the workspace

Confirm the laptop is available and identify the actual existing Cloud-Assess and
azqr-reference checkout paths. Preserve dirty files and interrupted work. Do not
reset, clean, force-push or blindly replay earlier commands. Do not paste access
tokens or environment credentials.

In PowerShell, after setting $TargetRepo and $ReferenceRepo to the actual paths:

```powershell
Get-Command git, go, az -ErrorAction SilentlyContinue | Select-Object Name, Source
go version
git -C $TargetRepo status --short --branch
git -C $TargetRepo rev-parse HEAD
git -C $TargetRepo submodule status --recursive
git -C $ReferenceRepo status --short --branch
git -C $ReferenceRepo rev-parse HEAD
git -C $ReferenceRepo submodule status --recursive
az account show --query '{tenantId:tenantId,subscriptionId:id,name:name,state:state}' --output json
az account list --all --query '[].{id:id,name:name,tenantId:tenantId,state:state}' --output json
```

These are local inspection/account-list commands. Actual missing tooling/login,
dirty checkout or scope selection is classified before recovery instructions.
No login secrets are requested. Installed Go may differ from the required1.26.9;
inspection does not certify build compatibility.

Required original source pin8e4f0577f3615e6c9014c031bcad079f235369cc and
APRL60eaddda76541f6adbc1c5ffa686829807e55e29. Both initialized APRL checkouts
must match their recorded gitlinks. The runner rejects dirty source/target trees
and mismatched pins. Verify a clean accepted target checkout before a live run;
do not use the draft source pilot merely because its CI starts.

## Scope and run order

| Order | Scope and evidence | Boundary |
| --- | --- | --- |
|1| Existing approved nonproduction subscription plus haz-rg-aitestbed-wus, if still accessible and confirmed | Same identity/scope for source and target; default Graph/diagnostics/Advisor/Defender first |
|2| Same confirmed scope, adding policy, defender-recommendations and cost | Classify health, empty data, denial and findings separately; no fixture creation |
|3| Two explicitly approved distinct nonproduction subscriptions, if available | Confirm IDs and both resolved scopes; the previously described Advisory Prd child is not assumed authorized |
|4| Confirmed existing AdvisoryDev leaf | Verify actual group ID, descendants and permissions before running; leaf evidence does not close nested traversal |
|5| Existing public plugin live cases under their accepted contracts | Separate commands/filters and source outputs; region-selection remains unavailable |
|Deferred| DV001 wholly nonproduction nested parent | Advisory contains a production child; a Dev include filter does not prevent parent discovery through Prd |

Do not change resources, roles, hierarchy, budgets or deploy fixtures. No public
listener, production scope substitution or release activation is authorized.

## Existing core runner invocation

After checkout/access/scope verification, use the actual accepted
scripts/live-equivalence.ps1. This is a template, not an instruction to run before
the above conditions. $ApprovedSubscriptionId must be the confirmed nonproduction
subscription, not a placeholder or the CLI's incidental default.

```powershell
& (Join-Path $TargetRepo 'scripts/live-equivalence.ps1') `
    -TargetRepo $TargetRepo `
    -ReferenceRepo $ReferenceRepo `
    -SubscriptionId @($ApprovedSubscriptionId) `
    -ResourceGroup @('haz-rg-aitestbed-wus')
```

The follow-up uses the same confirmed scope and adds
-Stages @('policy', 'defender-recommendations', 'cost'). The helper adds these to
default stages rather than replacing defaults. Inspect the runner's effective
stage metadata before interpreting results. Its plugin stage is a separate
selection and does not provide a region-selection runtime.

Default output: artifacts/equivalence under the target checkout. Retain the full
private bundle, including run-metadata.json, equivalence.json, source/target JSON,
stdout/stderr and exit codes. Read the actual files to assess comparator result,
requested/resolved/contributing scope, completeness, warnings, counts and timing.
Empty/no-permission data and a PASS comparator are not blanket feature validation.
These bundles contain unredacted Azure identifiers and must not be committed or
published to public GitHub. Update ledger and deferred gates only from reviewed
actual evidence; retain failed attempts as history.

## Resume handoff

When access returns, report availability and the two checkout paths, then inspect
the above state before choosing the next live command. Missing laptop access does
not block ongoing synthetic SKU work. No live ETA or release readiness increase
is inferred from expected availability.
