Set-StrictMode -Version Latest

function Get-NonEmptyValues {
    param(
        [AllowNull()]
        [string[]] $Values
    )

    $resolved = @()
    foreach ($value in @($Values)) {
        if (-not [string]::IsNullOrWhiteSpace($value)) {
            $resolved += $value.Trim()
        }
    }
    return @($resolved)
}

function Resolve-LiveEquivalenceScope {
    [CmdletBinding()]
    param(
        [AllowNull()]
        [string[]] $SubscriptionId,

        [AllowNull()]
        [string[]] $ResourceGroup,

        [AllowNull()]
        [string[]] $ManagementGroupId
    )

    $subscriptions = @(Get-NonEmptyValues -Values $SubscriptionId)
    $resourceGroups = @(Get-NonEmptyValues -Values $ResourceGroup)
    $managementGroups = @(Get-NonEmptyValues -Values $ManagementGroupId)

    if ($managementGroups.Count -gt 0) {
        if ($subscriptions.Count -gt 0 -or $resourceGroups.Count -gt 0) {
            throw 'Management-group scope cannot be combined with subscription or resource-group scope.'
        }

        return [pscustomobject]@{
            Arguments = @('--management-group-id', ($managementGroups -join ','))
            Metadata = [ordered]@{
                managementGroupIds = @($managementGroups)
            }
        }
    }

    if ($subscriptions.Count -eq 0) {
        throw 'At least one subscription ID or management-group ID is required.'
    }
    if ($resourceGroups.Count -gt 0 -and $subscriptions.Count -ne 1) {
        throw 'Resource-group scope requires exactly one subscription ID.'
    }

    $arguments = @('--subscription-id', ($subscriptions -join ','))
    $metadata = [ordered]@{
        subscriptionIds = @($subscriptions)
    }
    if ($resourceGroups.Count -gt 0) {
        $arguments += @('--resource-group', ($resourceGroups -join ','))
        $metadata.resourceGroups = @($resourceGroups)
    }

    return [pscustomobject]@{
        Arguments = @($arguments)
        Metadata = $metadata
    }
}

function Resolve-LiveEquivalenceStages {
    [CmdletBinding()]
    param(
        [AllowNull()]
        [string[]] $Stages
    )

    $stageOrder = @(
        'graph',
        'diagnostics',
        'advisor',
        'defender',
        'defender-recommendations',
        'arc',
        'policy',
        'cost',
        'plugin'
    )
    $enabled = @{
        graph = $true
        diagnostics = $true
        advisor = $true
        defender = $true
        'defender-recommendations' = $false
        arc = $false
        policy = $false
        cost = $false
        plugin = $false
    }

    $requested = @()
    foreach ($value in @($Stages)) {
        foreach ($rawToken in ($value -split ',')) {
            $token = $rawToken.Trim().ToLowerInvariant()
            if ([string]::IsNullOrWhiteSpace($token)) {
                continue
            }

            $requested += $token
            $stageName = $token
            $stageEnabled = $true
            if ($stageName.StartsWith('-')) {
                $stageEnabled = $false
                $stageName = $stageName.Substring(1)
            }

            if (-not $enabled.ContainsKey($stageName)) {
                throw "Unknown stage name: $stageName"
            }
            $enabled[$stageName] = $stageEnabled
        }
    }

    if (-not $enabled.graph) {
        throw 'The graph stage is mandatory for live equivalence runs and cannot be disabled.'
    }

    $effective = @()
    foreach ($stageName in $stageOrder) {
        if ($enabled[$stageName]) {
            $effective += $stageName
        }
    }

    $arguments = @()
    if ($requested.Count -gt 0) {
        $arguments = @('--stages', ($requested -join ','))
    }

    return [pscustomobject]@{
        Arguments = @($arguments)
        RequestedStages = @($requested)
        EffectiveStages = @($effective)
        UsesImplicitDefaults = ($requested.Count -eq 0)
    }
}

Export-ModuleMember -Function Resolve-LiveEquivalenceScope, Resolve-LiveEquivalenceStages
