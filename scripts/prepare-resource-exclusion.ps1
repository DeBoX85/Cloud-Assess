#requires -Version 5.1

<#
.SYNOPSIS
Creates paired AZQR and Cloud Assess resource-exclusion filters from local evidence.

.DESCRIPTION
Reads only retained reports. Selects the one resource associated with an observed
baseline recommendation, verifies that it remains in the newer target inventory,
and writes YAML under the git-ignored artifacts directory. No Azure calls are made.
The filter files contain an unredacted ARM ID and must remain local.
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string] $BaselineTargetJson,

    [Parameter(Mandatory = $true)]
    [string] $CurrentTargetJson,

    [Parameter(Mandatory = $true)]
    [string] $ExpectedSubscriptionId,

    [string] $RecommendationId = '1981f704-97b9-b645-9c57-33f8ded9261a'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Read-Report {
    param([string] $Path)
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) {
        throw "Report not found: $Path"
    }
    return (Get-Content -LiteralPath $Path -Raw | ConvertFrom-Json)
}

function Normalize-Id {
    param([string] $Id)
    return $Id.Trim().TrimEnd('/').ToLowerInvariant()
}

function Quote-Yaml {
    param([string] $Value)
    return "'$($Value.Replace("'", "''"))'"
}

$baseline = Read-Report $BaselineTargetJson
$current = Read-Report $CurrentTargetJson
$expected = $ExpectedSubscriptionId.Trim().ToLowerInvariant()

if ($expected -notmatch '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$') {
    throw 'ExpectedSubscriptionId must be a subscription GUID.'
}

$baselineRows = @($baseline.findings | Where-Object {
    [string]$_.recommendationId -ieq $RecommendationId
})
$candidateIds = @($baselineRows | ForEach-Object {
    Normalize-Id ([string]$_.resourceId)
} | Where-Object { $_ } | Sort-Object -Unique)
if ($candidateIds.Count -ne 1) {
    throw "Expected one distinct baseline resource for recommendation $RecommendationId; found $($candidateIds.Count)."
}

$matching = @($current.resources | Where-Object {
    (Normalize-Id ([string]$_.id)) -eq $candidateIds[0]
})
if ($matching.Count -ne 1) {
    throw "Expected the baseline resource once in current target inventory; found $($matching.Count)."
}

$resourceId = ([string]$matching[0].id).Trim().TrimEnd('/')
$armMatch = [regex]::Match($resourceId, '^/subscriptions/([^/]+)/resourceGroups/[^/]+/providers/[^/]+/[^/]+/[^/]+$', 'IgnoreCase')
if (-not $armMatch.Success) {
    throw 'Selected inventory ID is not an expected ARM resource ID.'
}
if ($armMatch.Groups[1].Value -ine $expected -or [string]$matching[0].subscriptionId -ine $expected) {
    throw 'Selected resource does not belong to the expected subscription.'
}

$scopeStages = @($current.stages | Where-Object { $_.name -eq 'scope' })
if ($scopeStages.Count -ne 1 -or [int]$scopeStages[0].records -ne 1) {
    throw 'Current report must show one resolved scope subscription.'
}

$otherSubscriptions = @($current.resources | Where-Object {
    [string]$_.subscriptionId -ine $expected
})
if ($otherSubscriptions.Count -gt 0) {
    throw 'Current inventory contains resources outside the expected subscription.'
}

$outputRoot = Join-Path (Split-Path -Parent $PSScriptRoot) 'artifacts/equivalence/resource-exclusion-filters'
$outputDir = Join-Path $outputRoot ("{0}_{1}" -f [DateTime]::UtcNow.ToString('yyyyMMdd_HHmmss_fffZ'), [Guid]::NewGuid().ToString('N').Substring(0, 8))
New-Item -ItemType Directory -Path $outputDir -ErrorAction Stop | Out-Null

$azqrPath = Join-Path $outputDir 'azqr.yml'
$targetPath = Join-Path $outputDir 'cloud-assess.yml'
$quotedId = Quote-Yaml $resourceId
$azqrYaml = "azqr:`n  exclude:`n    services:`n      - $quotedId`n"
$targetYaml = "assessment:`n  exclude:`n    resources:`n      - $quotedId`n"
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
[System.IO.File]::WriteAllText($azqrPath, $azqrYaml, $utf8NoBom)
[System.IO.File]::WriteAllText($targetPath, $targetYaml, $utf8NoBom)

$sha256 = [System.Security.Cryptography.SHA256]::Create()
try {
    $digest = $sha256.ComputeHash([System.Text.Encoding]::UTF8.GetBytes((Normalize-Id $resourceId)))
    $idHash = ([BitConverter]::ToString($digest) -replace '-', '').Substring(0, 16).ToLowerInvariant()
}
finally {
    $sha256.Dispose()
}

[pscustomobject]@{
    ResourceIdHash = $idHash
    BaselineFindingRows = $baselineRows.Count
    CurrentInventoryMatches = $matching.Count
    ReferenceFilters = $azqrPath
    TargetFilters = $targetPath
}
