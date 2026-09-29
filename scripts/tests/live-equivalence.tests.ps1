#requires -Version 5.1

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$modulePath = Join-Path (Split-Path -Parent $PSScriptRoot) 'live-equivalence.helpers.psm1'
Import-Module $modulePath -Force

function Assert-SequenceEqual {
    param(
        [string] $Label,
        [AllowNull()]
        [object[]] $Actual,
        [AllowNull()]
        [object[]] $Expected
    )

    $actualText = (@($Actual) -join '|')
    $expectedText = (@($Expected) -join '|')
    if ($actualText -ne $expectedText) {
        throw "$Label was '$actualText', expected '$expectedText'."
    }
}

function Assert-Throws {
    param(
        [string] $Label,
        [scriptblock] $Action,
        [string] $MessagePattern
    )

    try {
        & $Action
    }
    catch {
        if ($_.Exception.Message -notmatch $MessagePattern) {
            throw "$Label threw '$($_.Exception.Message)', expected pattern '$MessagePattern'."
        }
        return
    }
    throw "$Label did not throw."
}

$defaults = Resolve-LiveEquivalenceStages
Assert-SequenceEqual -Label 'default stage arguments' -Actual $defaults.Arguments -Expected @()
Assert-SequenceEqual -Label 'default requested stages' -Actual $defaults.RequestedStages -Expected @()
Assert-SequenceEqual -Label 'default effective stages' -Actual $defaults.EffectiveStages -Expected @('graph', 'diagnostics', 'advisor', 'defender')
if (-not $defaults.UsesImplicitDefaults) {
    throw 'Default stage selection must be marked as implicit.'
}

$optional = Resolve-LiveEquivalenceStages -Stages @('policy,defender-recommendations', 'cost', '-advisor')
Assert-SequenceEqual -Label 'explicit stage arguments' -Actual $optional.Arguments -Expected @('--stages', 'policy,defender-recommendations,cost,-advisor')
Assert-SequenceEqual -Label 'explicit requested stages' -Actual $optional.RequestedStages -Expected @('policy', 'defender-recommendations', 'cost', '-advisor')
Assert-SequenceEqual -Label 'explicit effective stages' -Actual $optional.EffectiveStages -Expected @('graph', 'diagnostics', 'defender', 'defender-recommendations', 'policy', 'cost')
if ($optional.UsesImplicitDefaults) {
    throw 'Explicit stage selection must not be marked as implicit.'
}

Assert-Throws -Label 'unknown stage' -Action { Resolve-LiveEquivalenceStages -Stages 'unknown' } -MessagePattern 'Unknown stage name'
Assert-Throws -Label 'disabled graph' -Action { Resolve-LiveEquivalenceStages -Stages '-graph' } -MessagePattern 'graph stage is mandatory'

$singleSubscription = Resolve-LiveEquivalenceScope -SubscriptionId 'sub-1'
Assert-SequenceEqual -Label 'single-subscription arguments' -Actual $singleSubscription.Arguments -Expected @('--subscription-id', 'sub-1')
Assert-SequenceEqual -Label 'single-subscription metadata' -Actual $singleSubscription.Metadata.subscriptionIds -Expected @('sub-1')

$multiSubscription = Resolve-LiveEquivalenceScope -SubscriptionId @('sub-1', 'sub-2')
Assert-SequenceEqual -Label 'multi-subscription arguments' -Actual $multiSubscription.Arguments -Expected @('--subscription-id', 'sub-1,sub-2')
Assert-SequenceEqual -Label 'multi-subscription metadata' -Actual $multiSubscription.Metadata.subscriptionIds -Expected @('sub-1', 'sub-2')

$resourceGroups = Resolve-LiveEquivalenceScope -SubscriptionId 'sub-1' -ResourceGroup @('rg-1', 'rg-2')
Assert-SequenceEqual -Label 'resource-group arguments' -Actual $resourceGroups.Arguments -Expected @('--subscription-id', 'sub-1', '--resource-group', 'rg-1,rg-2')
Assert-SequenceEqual -Label 'resource-group metadata' -Actual $resourceGroups.Metadata.resourceGroups -Expected @('rg-1', 'rg-2')

$managementGroups = Resolve-LiveEquivalenceScope -ManagementGroupId @('mg-1', 'mg-2')
Assert-SequenceEqual -Label 'management-group arguments' -Actual $managementGroups.Arguments -Expected @('--management-group-id', 'mg-1,mg-2')
Assert-SequenceEqual -Label 'management-group metadata' -Actual $managementGroups.Metadata.managementGroupIds -Expected @('mg-1', 'mg-2')

Assert-Throws -Label 'resource groups across subscriptions' -Action {
    Resolve-LiveEquivalenceScope -SubscriptionId @('sub-1', 'sub-2') -ResourceGroup 'rg-1'
} -MessagePattern 'exactly one subscription'

Assert-Throws -Label 'mixed management-group scope' -Action {
    Resolve-LiveEquivalenceScope -SubscriptionId 'sub-1' -ManagementGroupId 'mg-1'
} -MessagePattern 'cannot be combined'

$sensitive = 'subscription-secret-1111'
$summary = Get-TargetAssessmentSummary -Report ([pscustomobject]@{
    completeness = 'complete_with_warnings'
    stages = @(
        [pscustomobject]@{ name = 'scope'; status = 'completed'; records = 1 },
        [pscustomobject]@{
            name = 'diagnostics'; status = 'completed_with_warnings'; records = 2
            warnings = @([pscustomobject]@{ code = 'diagnostics_subrequest_non_success'; message = "failed $sensitive" })
        },
        [pscustomobject]@{ name = 'policy'; status = 'skipped'; records = 0 }
    )
})
if ($summary.completeness -ne 'complete_with_warnings' -or $summary.stages.Count -ne 3 -or
    $summary.stages[1].status -ne 'completed_with_warnings' -or $summary.stages[1].records -ne 2) {
    throw 'Stage-health summary lost status or record counts.'
}
Assert-SequenceEqual -Label 'stage warning codes' -Actual $summary.stages[1].warningCodes -Expected @('diagnostics_subrequest_non_success')
Assert-SequenceEqual -Label 'omitted stage warnings' -Actual $summary.stages[0].warningCodes -Expected @()
if (($summary | ConvertTo-Json -Depth 12) -match $sensitive) {
    throw 'Stage-health summary leaked a warning message with an Azure identifier.'
}

Write-Host 'live-equivalence helper tests: PASS'
