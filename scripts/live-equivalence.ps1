#requires -Version 5.1

<#
.SYNOPSIS
Runs the pinned AZQR reference and Cloud Assess against the same Azure scope,
then produces a semantic equivalence report and an auditable evidence bundle.

.DESCRIPTION
This script is intended for development validation only. It does not modify Azure.

The generated evidence bundle contains unredacted Azure identifiers by design.
The default output location is under artifacts/, which is ignored by Git.

The AZQR reference checkout must be exactly at the pinned reference commit and
must have a clean working tree. The Cloud Assess checkout must also be clean so
its commit SHA fully describes the code being validated.
#>

[CmdletBinding(DefaultParameterSetName = 'Subscription')]
param(
    [Parameter(Mandatory = $true)]
    [string] $ReferenceRepo,

    [Parameter(ParameterSetName = 'Subscription', Mandatory = $true)]
    [string] $SubscriptionId,

    [Parameter(ParameterSetName = 'Subscription')]
    [string] $ResourceGroup,

    [Parameter(ParameterSetName = 'ManagementGroup', Mandatory = $true)]
    [string] $ManagementGroupId,

    [string[]] $Stages,

    [string] $ReferenceFilters,

    [string] $TargetFilters,

    [string] $TargetRepo,

    [string] $OutputRoot
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$PinnedReferenceCommit = '8e4f0577f3615e6c9014c031bcad079f235369cc'
$PinnedAprlCommit = '60eaddda76541f6adbc1c5ffa686829807e55e29'

function Resolve-Directory {
    param(
        [Parameter(Mandatory = $true)]
        [string] $Path,
        [Parameter(Mandatory = $true)]
        [string] $Label
    )

    if (-not (Test-Path -LiteralPath $Path -PathType Container)) {
        throw "$Label does not exist or is not a directory: $Path"
    }
    return (Resolve-Path -LiteralPath $Path).Path
}

function Resolve-OptionalFile {
    param(
        [string] $Path,
        [string] $Label
    )

    if ([string]::IsNullOrWhiteSpace($Path)) {
        return $null
    }
    if (-not (Test-Path -LiteralPath $Path -PathType Leaf)) {
        throw "$Label does not exist or is not a file: $Path"
    }
    return (Resolve-Path -LiteralPath $Path).Path
}

function Invoke-GitText {
    param(
        [string] $Repo,
        [string[]] $Arguments
    )

    Push-Location $Repo
    try {
        $output = & git @Arguments 2>&1
        if ($LASTEXITCODE -ne 0) {
            throw "git $($Arguments -join ' ') failed in $($Repo). Details: $($output -join [Environment]::NewLine)"
        }
        return (($output | Out-String).Trim())
    }
    finally {
        Pop-Location
    }
}

function Assert-CleanRepo {
    param(
        [string] $Repo,
        [string] $Label
    )

    $status = Invoke-GitText -Repo $Repo -Arguments @('status', '--porcelain')
    if (-not [string]::IsNullOrWhiteSpace($status)) {
        throw "$Label has uncommitted or untracked changes. The live equivalence run requires a clean checkout so the recorded commit SHA is authoritative."
    }
}

function Assert-SubmoduleCommit {
    param(
        [string] $Repo,
        [string] $RelativePath,
        [string] $ExpectedCommit,
        [string] $Label
    )

    $treeCommit = Invoke-GitText -Repo $Repo -Arguments @('rev-parse', "HEAD:$RelativePath")
    if ($treeCommit -ne $ExpectedCommit) {
        throw "$Label tree entry is $treeCommit, expected $ExpectedCommit"
    }

    $checkoutPath = Join-Path $Repo $RelativePath
    if (-not (Test-Path -LiteralPath $checkoutPath -PathType Container)) {
        throw "$Label is not initialized. Run: git submodule update --init --recursive"
    }

    try {
        $checkoutCommit = Invoke-GitText -Repo $checkoutPath -Arguments @('rev-parse', 'HEAD')
    }
    catch {
        throw "$Label is not initialized correctly. Run: git submodule update --init --recursive"
    }

    if ($checkoutCommit -ne $ExpectedCommit) {
        throw "$Label checkout is at $checkoutCommit, expected $ExpectedCommit. Run: git submodule update --init --recursive"
    }
}

function Invoke-CapturedNative {
    param(
        [string] $Label,
        [string] $WorkingDirectory,
        [string] $Command,
        [string[]] $Arguments,
        [string] $StdoutPath,
        [string] $StderrPath
    )

    $started = [DateTime]::UtcNow
    Push-Location $WorkingDirectory
    $previousErrorActionPreference = $ErrorActionPreference
    try {
        # Windows PowerShell 5.1 promotes native-process stderr records according to
        # ErrorActionPreference. Go writes normal module-download progress to stderr,
        # so the script-level 'Stop' setting would otherwise abort a successful run
        # before LASTEXITCODE can be evaluated. Capture stderr and judge the native
        # command solely by its process exit code.
        $ErrorActionPreference = 'Continue'
        & $Command @Arguments 1> $StdoutPath 2> $StderrPath
        $exitCode = $LASTEXITCODE
    }
    finally {
        $ErrorActionPreference = $previousErrorActionPreference
        Pop-Location
    }
    $finished = [DateTime]::UtcNow

    return [ordered]@{
        label = $Label
        command = $Command
        arguments = @($Arguments)
        workingDirectory = $WorkingDirectory
        startedUtc = $started.ToString('o')
        finishedUtc = $finished.ToString('o')
        durationSeconds = [Math]::Round(($finished - $started).TotalSeconds, 3)
        exitCode = $exitCode
        stdout = $StdoutPath
        stderr = $StderrPath
    }
}

function Get-OptionalFileEvidence {
    param([string] $Path)

    if ([string]::IsNullOrWhiteSpace($Path)) {
        return $null
    }

    return [ordered]@{
        path = $Path
        sha256 = (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
    }
}

function Write-Metadata {
    param(
        [System.Collections.IDictionary] $Metadata,
        [string] $Path
    )

    $Metadata | ConvertTo-Json -Depth 12 | Set-Content -LiteralPath $Path -Encoding UTF8
}

if ([string]::IsNullOrWhiteSpace($TargetRepo)) {
    $TargetRepo = Join-Path $PSScriptRoot '..'
}

$ReferenceRepo = Resolve-Directory -Path $ReferenceRepo -Label 'Reference repository'
$TargetRepo = Resolve-Directory -Path $TargetRepo -Label 'Cloud Assess repository'
$ReferenceFilters = Resolve-OptionalFile -Path $ReferenceFilters -Label 'Reference filters file'
$TargetFilters = Resolve-OptionalFile -Path $TargetFilters -Label 'Target filters file'

if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    throw 'git is required but was not found in PATH'
}
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw 'Go is required but was not found in PATH'
}

$referenceCommit = Invoke-GitText -Repo $ReferenceRepo -Arguments @('rev-parse', 'HEAD')
if ($referenceCommit -ne $PinnedReferenceCommit) {
    throw "Reference repository is at $referenceCommit, but live equivalence requires pinned commit $PinnedReferenceCommit"
}

Assert-CleanRepo -Repo $ReferenceRepo -Label 'Reference repository'
Assert-CleanRepo -Repo $TargetRepo -Label 'Cloud Assess repository'

Assert-SubmoduleCommit -Repo $ReferenceRepo -RelativePath 'internal/graph/aprl' -ExpectedCommit $PinnedAprlCommit -Label 'Reference APRL submodule'
Assert-SubmoduleCommit -Repo $TargetRepo -RelativePath 'internal/rules/upstream/aprl' -ExpectedCommit $PinnedAprlCommit -Label 'Cloud Assess APRL submodule'

$targetCommit = Invoke-GitText -Repo $TargetRepo -Arguments @('rev-parse', 'HEAD')
$referenceBranch = Invoke-GitText -Repo $ReferenceRepo -Arguments @('rev-parse', '--abbrev-ref', 'HEAD')
$targetBranch = Invoke-GitText -Repo $TargetRepo -Arguments @('rev-parse', '--abbrev-ref', 'HEAD')
$goVersion = ((& go version 2>&1) | Out-String).Trim()

if ([string]::IsNullOrWhiteSpace($OutputRoot)) {
    $OutputRoot = Join-Path $TargetRepo 'artifacts/equivalence'
}
elseif (-not [System.IO.Path]::IsPathRooted($OutputRoot)) {
    $OutputRoot = Join-Path $TargetRepo $OutputRoot
}

$runStamp = [DateTime]::UtcNow.ToString('yyyyMMdd_HHmmssZ')
$runDirectory = Join-Path $OutputRoot $runStamp
New-Item -ItemType Directory -Path $runDirectory -Force | Out-Null
$runDirectory = (Resolve-Path -LiteralPath $runDirectory).Path

$referenceBase = Join-Path $runDirectory 'reference'
$targetBase = Join-Path $runDirectory 'target'
$referenceJson = "$referenceBase.json"
$targetJson = "$targetBase.json"
$equivalenceJson = Join-Path $runDirectory 'equivalence.json'
$metadataPath = Join-Path $runDirectory 'run-metadata.json'

$scopeArguments = @()
$scopeMetadata = [ordered]@{}
if ($PSCmdlet.ParameterSetName -eq 'ManagementGroup') {
    $scopeArguments += @('--management-group-id', $ManagementGroupId)
    $scopeMetadata.managementGroupId = $ManagementGroupId
}
else {
    $scopeArguments += @('--subscription-id', $SubscriptionId)
    $scopeMetadata.subscriptionId = $SubscriptionId
    if (-not [string]::IsNullOrWhiteSpace($ResourceGroup)) {
        $scopeArguments += @('--resource-group', $ResourceGroup)
        $scopeMetadata.resourceGroup = $ResourceGroup
    }
}

$stageArguments = @()
if ($null -ne $Stages -and $Stages.Count -gt 0) {
    $stageArguments = @('--stages', ($Stages -join ','))
}

$referenceArguments = @(
    'run', './cmd/azqr', 'scan'
) + $scopeArguments + $stageArguments + @(
    '--json',
    '--xlsx=false',
    '--mask=false',
    '--output-name', $referenceBase
)

if ($ReferenceFilters) {
    $referenceArguments += @('--filters', $ReferenceFilters)
}

$targetArguments = @(
    'run', './cmd/cloud-assess', 'scan'
) + $scopeArguments + $stageArguments + @(
    '--json',
    '--xlsx=false',
    '--redact-subscription-ids=false',
    '--output-name', $targetBase
)

if ($TargetFilters) {
    $targetArguments += @('--filters', $TargetFilters)
}

$metadata = [ordered]@{
    schemaVersion = '1.0'
    purpose = 'Cloud Assess live source-versus-target equivalence'
    warning = 'This evidence bundle contains unredacted Azure identifiers and must be handled as sensitive assessment data.'
    createdUtc = [DateTime]::UtcNow.ToString('o')
    pinnedReferenceCommit = $PinnedReferenceCommit
    pinnedAprlCommit = $PinnedAprlCommit
    reference = [ordered]@{
        repo = $ReferenceRepo
        branch = $referenceBranch
        commit = $referenceCommit
        filters = Get-OptionalFileEvidence -Path $ReferenceFilters
    }
    target = [ordered]@{
        repo = $TargetRepo
        branch = $targetBranch
        commit = $targetCommit
        filters = Get-OptionalFileEvidence -Path $TargetFilters
    }
    scope = $scopeMetadata
    stages = @($Stages)
    environment = [ordered]@{
        powerShell = $PSVersionTable.PSVersion.ToString()
        go = $goVersion
        azureCloud = $env:AZURE_CLOUD
        azureAuthorityHost = $env:AZURE_AUTHORITY_HOST
        azureResourceManagerEndpoint = $env:AZURE_RESOURCE_MANAGER_ENDPOINT
        azureResourceManagerAudience = $env:AZURE_RESOURCE_MANAGER_AUDIENCE
    }
    artifacts = [ordered]@{
        referenceJson = $referenceJson
        targetJson = $targetJson
        equivalenceJson = $equivalenceJson
    }
    executions = @()
}

Write-Metadata -Metadata $metadata -Path $metadataPath

Write-Host "Evidence directory: $runDirectory"
Write-Host 'Running pinned AZQR reference...'

$referenceInvocation = @{
    Label = 'reference'
    WorkingDirectory = $ReferenceRepo
    Command = 'go'
    Arguments = $referenceArguments
    StdoutPath = (Join-Path $runDirectory 'reference.stdout.log')
    StderrPath = (Join-Path $runDirectory 'reference.stderr.log')
}
$referenceExecution = Invoke-CapturedNative @referenceInvocation

$metadata.executions += $referenceExecution
Write-Metadata -Metadata $metadata -Path $metadataPath

if ($referenceExecution.exitCode -ne 0) {
    throw "Pinned reference scan failed with exit code $($referenceExecution.exitCode). See $($referenceExecution.stderr)"
}
if (-not (Test-Path -LiteralPath $referenceJson -PathType Leaf)) {
    throw "Pinned reference completed without producing expected JSON: $referenceJson"
}

Write-Host 'Running Cloud Assess...'

$targetInvocation = @{
    Label = 'target'
    WorkingDirectory = $TargetRepo
    Command = 'go'
    Arguments = $targetArguments
    StdoutPath = (Join-Path $runDirectory 'target.stdout.log')
    StderrPath = (Join-Path $runDirectory 'target.stderr.log')
}
$targetExecution = Invoke-CapturedNative @targetInvocation

$metadata.executions += $targetExecution
Write-Metadata -Metadata $metadata -Path $metadataPath

if ($targetExecution.exitCode -ne 0) {
    throw "Cloud Assess scan failed with exit code $($targetExecution.exitCode). See $($targetExecution.stderr)"
}
if (-not (Test-Path -LiteralPath $targetJson -PathType Leaf)) {
    throw "Cloud Assess completed without producing expected JSON: $targetJson"
}

Write-Host 'Running semantic equivalence comparison...'

$compareArguments = @(
    'run', './tools/equivalence',
    '--reference', $referenceJson,
    '--target', $targetJson,
    '--output', $equivalenceJson
)

$compareInvocation = @{
    Label = 'equivalence'
    WorkingDirectory = $TargetRepo
    Command = 'go'
    Arguments = $compareArguments
    StdoutPath = (Join-Path $runDirectory 'equivalence.stdout.log')
    StderrPath = (Join-Path $runDirectory 'equivalence.stderr.log')
}
$compareExecution = Invoke-CapturedNative @compareInvocation

$metadata.executions += $compareExecution
$metadata.completedUtc = [DateTime]::UtcNow.ToString('o')
$metadata.equivalenceExitCode = $compareExecution.exitCode
Write-Metadata -Metadata $metadata -Path $metadataPath

switch ($compareExecution.exitCode) {
    0 {
        Write-Host 'Semantic equivalence: PASS'
        Write-Host "Report: $equivalenceJson"
        exit 0
    }
    1 {
        Write-Warning 'Semantic equivalence: DIFFERENCES FOUND'
        Write-Host "Report: $equivalenceJson"
        exit 1
    }
    default {
        throw "Equivalence tool failed with exit code $($compareExecution.exitCode). See $($compareExecution.stderr)"
    }
}
