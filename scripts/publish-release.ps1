[CmdletBinding()]
param(
    [string]$Version,
    [string]$PackageName = "@yixiaoermail/cli",
    [string]$DownloadRootUrl = $env:YXER_DOWNLOAD_ROOT_URL,
    [string]$Bucket = $env:TOS_BUCKET,
    [string]$Region = $env:TOS_REGION,
    [string]$S3Endpoint = $env:TOS_S3_ENDPOINT,
    [string]$ObjectPrefix = $env:TOS_OBJECT_PREFIX,
    [int]$MaxConcurrentRequests = 20,
    [string]$MultipartChunkSize = "16MB",
    [switch]$DryRun,
    [switch]$Execute,
    [switch]$SkipBuildTests
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$repoRoot = Split-Path -Parent $PSScriptRoot
$localCredentialDir = Join-Path $repoRoot ".local-release"
$releaseDir = Join-Path $repoRoot "out\release"
$npmDir = Join-Path $repoRoot "out\npm"
$verifyScript = Join-Path $PSScriptRoot "verify-release-artifacts.ps1"
$buildScript = Join-Path $PSScriptRoot "build-npm-package.ps1"
$checkNpmVersionScript = Join-Path $PSScriptRoot "check-npm-version.ps1"
$awsConfigPath = Join-Path ([System.IO.Path]::GetTempPath()) ("yxer-aws-" + [guid]::NewGuid().ToString("N"))
$npmConfigPath = Join-Path ([System.IO.Path]::GetTempPath()) ("yxer-npm-" + [guid]::NewGuid().ToString("N") + ".npmrc")

function Assert-LastExitCode {
    param([string]$CommandName)
    if ($LASTEXITCODE -ne 0) {
        throw "$CommandName failed with exit code $LASTEXITCODE"
    }
}

function Get-EnvOrDefault {
    param([string]$Value, [string]$DefaultValue)
    if ([string]::IsNullOrWhiteSpace($Value)) { return $DefaultValue }
    return $Value.Trim()
}

function Get-LocalCredential {
    param([string]$Value, [string]$FileName, [string]$Placeholder)

    if (-not [string]::IsNullOrWhiteSpace($Value)) { return $Value.Trim() }
    $credentialPath = Join-Path $localCredentialDir $FileName
    if (-not (Test-Path -LiteralPath $credentialPath -PathType Leaf)) { return "" }

    $fileValue = (Get-Content -LiteralPath $credentialPath -Raw).Trim()
    if ($fileValue -eq $Placeholder) { return "" }
    if ($fileValue -match "[`r`n]") {
        throw "Credential file must contain one value on one line: $credentialPath"
    }
    return $fileValue
}

function Get-SourceVersion {
    $responsePath = Join-Path $repoRoot "internal\domain\response.go"
    $skillPath = Join-Path $repoRoot "skills\yixiaoer\SKILL.md"
    $response = Get-Content -LiteralPath $responsePath -Raw
    $match = [regex]::Match($response, '(?:const|var)\s+SkillVersion\s*=\s*"([^"]+)"')
    if (-not $match.Success) { throw "SkillVersion constant not found in $responsePath" }
    $goVersion = $match.Groups[1].Value.Trim()

    $skillVersion = $null
    foreach ($line in Get-Content -LiteralPath $skillPath) {
        if ($line.Trim() -like "version:*") {
            $skillVersion = $line.Trim().Substring("version:".Length).Trim().Trim('"').Trim("'")
            break
        }
    }
    if (-not $skillVersion) { throw "version field not found in $skillPath" }
    if ($goVersion -ne $skillVersion) {
        throw "Version sources differ: response.go=$goVersion, SKILL.md=$skillVersion"
    }
    return $goVersion
}

function Get-FileSha256 {
    param([string]$Path)
    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Invoke-DownloadVerification {
    param([string]$BaseUrl)

    $downloadDir = Join-Path ([System.IO.Path]::GetTempPath()) ("yxer-download-verify-" + [guid]::NewGuid().ToString("N"))
    New-Item -ItemType Directory -Path $downloadDir -Force | Out-Null
    try {
        $checksumsPath = Join-Path $downloadDir "checksums.txt"
        $checksumsUrl = "$BaseUrl/checksums.txt"
        Write-Host "Downloading release checksum list: $checksumsUrl"
        & curl.exe --fail --location --silent --show-error --retry 5 --retry-delay 5 $checksumsUrl --output $checksumsPath
        Assert-LastExitCode "Download release checksums from $checksumsUrl"

        $actualChecksums = (Get-Content -LiteralPath $checksumsPath -Raw).Replace("`r", "").TrimEnd("`n")
        $localChecksums = (Get-Content -LiteralPath (Join-Path $releaseDir "checksums.txt") -Raw).Replace("`r", "").TrimEnd("`n")
        if ($actualChecksums -ne $localChecksums) {
            throw "Downloaded checksums.txt does not match the local release artifact"
        }

        foreach ($line in $actualChecksums -split "`n") {
            if (-not $line.Trim()) { continue }
            if ($line -notmatch '^([0-9a-fA-F]{64})\s{2}(.+)$') {
                throw "Invalid checksum line from download URL: $line"
            }
            $expectedHash = $Matches[1].ToLowerInvariant()
            $fileName = $Matches[2].Trim()
            $downloadPath = Join-Path $downloadDir $fileName
            Write-Host "Downloading release artifact: $fileName"
            & curl.exe --fail --location --silent --show-error --retry 5 --retry-delay 5 "$BaseUrl/$fileName" --output $downloadPath
            Assert-LastExitCode "Download $fileName from $BaseUrl"
            $actualHash = Get-FileSha256 -Path $downloadPath
            if ($actualHash -ne $expectedHash) {
                throw "Download checksum mismatch for ${fileName}: expected $expectedHash but got $actualHash"
            }
        }
        Write-Host "Verified all release archives and checksums from download URL"
    } finally {
        Remove-Item -LiteralPath $downloadDir -Recurse -Force -ErrorAction SilentlyContinue
    }
}

if ($MaxConcurrentRequests -lt 1) { throw "MaxConcurrentRequests must be at least 1" }
if (-not (Get-Command go -ErrorAction SilentlyContinue)) { throw "Go is required on PATH" }
if (-not (Get-Command node -ErrorAction SilentlyContinue)) { throw "Node.js is required on PATH" }
if (-not (Get-Command npm -ErrorAction SilentlyContinue)) { throw "npm is required on PATH" }
if (-not (Get-Command aws -ErrorAction SilentlyContinue)) { throw "AWS CLI is required on PATH" }
if (-not (Get-Command curl.exe -ErrorAction SilentlyContinue)) { throw "curl.exe is required on PATH" }

$sourceVersion = Get-SourceVersion
if (-not $Version) { $Version = $sourceVersion }
if ($Version -ne $sourceVersion) {
    throw "Requested version '$Version' does not match source version '$sourceVersion'"
}
if ($Version -notmatch '^\d+\.\d+\.\d+$') { throw "Version must have the form major.minor.patch" }

$templatePackage = Get-Content -LiteralPath (Join-Path $repoRoot "npm\package.json") -Raw | ConvertFrom-Json
$DownloadRootUrl = Get-EnvOrDefault -Value $DownloadRootUrl -DefaultValue $templatePackage.yxerDownloadRootUrl
$Bucket = Get-EnvOrDefault -Value $Bucket -DefaultValue "yixiaoer-lite-asserts"
$Region = Get-EnvOrDefault -Value $Region -DefaultValue "cn-shanghai"
$S3Endpoint = Get-EnvOrDefault -Value $S3Endpoint -DefaultValue "https://tos-s3-cn-shanghai.volces.com"
$ObjectPrefix = Get-EnvOrDefault -Value $ObjectPrefix -DefaultValue "yxer/releases"
$ObjectPrefix = $ObjectPrefix.Trim('/')
if ([string]::IsNullOrWhiteSpace($Bucket)) { throw "TOS bucket must not be empty" }
if ([string]::IsNullOrWhiteSpace($ObjectPrefix)) { throw "TOS object prefix must not be empty" }

$parsedDownloadRoot = $null
if (-not [System.Uri]::TryCreate($DownloadRootUrl, [System.UriKind]::Absolute, [ref]$parsedDownloadRoot) -or $parsedDownloadRoot.Scheme -ne "https") {
    throw "DownloadRootUrl must be an absolute HTTPS URL"
}
$parsedS3Endpoint = $null
if (-not [System.Uri]::TryCreate($S3Endpoint, [System.UriKind]::Absolute, [ref]$parsedS3Endpoint) -or $parsedS3Endpoint.Scheme -ne "https") {
    throw "S3Endpoint must be an absolute HTTPS URL"
}

$releaseTag = "v$Version"
$destination = "s3://$Bucket/$ObjectPrefix/$releaseTag"
$downloadBaseUrl = "$($DownloadRootUrl.TrimEnd('/'))/$releaseTag"

$awsKey = Get-EnvOrDefault -Value $env:TOS_ACCESS_KEY_ID -DefaultValue $env:AWS_ACCESS_KEY_ID
$awsSecret = Get-EnvOrDefault -Value $env:TOS_SECRET_ACCESS_KEY -DefaultValue $env:AWS_SECRET_ACCESS_KEY
$npmToken = ""
if ($Execute -and -not $DryRun) {
    $awsKey = Get-LocalCredential -Value $awsKey -FileName "TOS_ACCESS_KEY_ID.txt" -Placeholder "PASTE_TOS_ACCESS_KEY_ID_HERE"
    $awsSecret = Get-LocalCredential -Value $awsSecret -FileName "TOS_SECRET_ACCESS_KEY.txt" -Placeholder "PASTE_TOS_SECRET_ACCESS_KEY_HERE"
    $npmTokenEnv = Get-EnvOrDefault -Value $env:NODE_AUTH_TOKEN -DefaultValue $env:NPM_TOKEN
    $npmToken = Get-LocalCredential -Value $npmTokenEnv -FileName "NPM_TOKEN.txt" -Placeholder "PASTE_NPM_TOKEN_HERE"
    if ([string]::IsNullOrWhiteSpace($awsKey) -or [string]::IsNullOrWhiteSpace($awsSecret)) {
        throw "Set TOS_ACCESS_KEY_ID and TOS_SECRET_ACCESS_KEY in the environment, or fill both files in .local-release before executing a release"
    }

    $preflightNodeAuthToken = $env:NODE_AUTH_TOKEN
    try {
        if ($npmToken) {
            $env:NODE_AUTH_TOKEN = $npmToken
            '//registry.npmjs.org/:_authToken=${NODE_AUTH_TOKEN}' | Set-Content -LiteralPath $npmConfigPath -Encoding ascii
            & npm whoami --registry "https://registry.npmjs.org" --userconfig $npmConfigPath
        } else {
            & npm whoami --registry "https://registry.npmjs.org"
        }
        $npmAuthExitCode = $LASTEXITCODE
    } finally {
        $env:NODE_AUTH_TOKEN = $preflightNodeAuthToken
        Remove-Item -LiteralPath $npmConfigPath -Force -ErrorAction SilentlyContinue
    }
    if ($npmAuthExitCode -ne 0) {
        if ($npmToken) {
            throw "The configured npm token was rejected by registry.npmjs.org; replace .local-release/NPM_TOKEN.txt or NODE_AUTH_TOKEN with a valid publishing token"
        }
        throw "npm authentication failed; fill .local-release/NPM_TOKEN.txt, set NODE_AUTH_TOKEN, or run npm login"
    }
}

Write-Host "Release version: $Version"
Write-Host "TOS destination: $destination/"
Write-Host "Public download path: $downloadBaseUrl/"
Write-Host "Building artifacts locally and verifying the npm tarball..."

& $buildScript -Version $Version -PackageName $PackageName -DownloadRootUrl $DownloadRootUrl -SkipTests:$SkipBuildTests
Assert-LastExitCode "Build release artifacts"

$packageFilePrefix = ($PackageName -replace '^@', '') -replace '/', '-'
$tarballPattern = "$packageFilePrefix-$Version.tgz"
$tarballs = @(Get-ChildItem -LiteralPath $npmDir -Filter $tarballPattern -File)
if ($tarballs.Count -ne 1) { throw "Expected exactly one npm tarball in $npmDir, found $($tarballs.Count)" }
$tarballPath = $tarballs[0].FullName
& $verifyScript -ReleaseDir "out\release" -TarballDir "out\npm" -TarballPattern $tarballPattern
Assert-LastExitCode "Verify release artifacts"

$originalAwsConfig = $env:AWS_CONFIG_FILE
$originalAwsKey = $env:AWS_ACCESS_KEY_ID
$originalAwsSecret = $env:AWS_SECRET_ACCESS_KEY
$originalAwsRegion = $env:AWS_DEFAULT_REGION
$originalMetadataDisabled = $env:AWS_EC2_METADATA_DISABLED
$originalNodeAuthToken = $env:NODE_AUTH_TOKEN
try {
    $npmAuthArgs = @()
    if ($Execute -and -not $DryRun -and $npmToken) {
        $env:NODE_AUTH_TOKEN = $npmToken
        '//registry.npmjs.org/:_authToken=${NODE_AUTH_TOKEN}' | Set-Content -LiteralPath $npmConfigPath -Encoding ascii
        $npmAuthArgs = @("--userconfig", $npmConfigPath)
    }

    $env:AWS_DEFAULT_REGION = $Region
    $env:AWS_EC2_METADATA_DISABLED = "true"
    $env:AWS_CONFIG_FILE = $awsConfigPath
    @"
[default]
region = $Region
s3 =
    addressing_style = virtual
    max_concurrent_requests = $MaxConcurrentRequests
    multipart_threshold = 16MB
    multipart_chunksize = $MultipartChunkSize
"@ | Set-Content -LiteralPath $awsConfigPath -Encoding ascii

    $awsBaseArgs = @(
        "s3", "cp", $releaseDir.TrimEnd('\'), "$destination/", "--recursive",
        "--exclude", "*", "--include", "checksums.txt", "--include", "yxer-cli-*.zip",
        "--include", "yxer-cli-*.tar.gz", "--endpoint-url", $S3Endpoint, "--no-progress"
    )

    if ($DryRun -or -not $Execute) {
        Write-Host "DRY RUN: no TOS upload or npm publish will be performed."
        Write-Host "TOS upload preview:"
        $dryRunArgs = $awsBaseArgs + @("--dryrun")
        & aws @dryRunArgs
        Assert-LastExitCode "Preview TOS upload"
        Write-Host "npm publish preview:"
        & npm publish $tarballPath --access public --registry "https://registry.npmjs.org" --dry-run
        Assert-LastExitCode "Preview npm publish"
        Write-Host "Review the generated files, then rerun with -Execute to upload and publish."
        return
    }

    & $checkNpmVersionScript -PackageName $PackageName -Version $Version -NpmAuthArgs $npmAuthArgs

    if ($npmToken) { $env:NODE_AUTH_TOKEN = $originalNodeAuthToken }
    $env:AWS_ACCESS_KEY_ID = $awsKey
    $env:AWS_SECRET_ACCESS_KEY = $awsSecret
    Write-Host "Uploading release archives to TOS..."
    & aws @awsBaseArgs
    Assert-LastExitCode "Upload release archives to TOS"
    $env:AWS_ACCESS_KEY_ID = $originalAwsKey
    $env:AWS_SECRET_ACCESS_KEY = $originalAwsSecret

    Invoke-DownloadVerification -BaseUrl $downloadBaseUrl

    if ($npmToken) { $env:NODE_AUTH_TOKEN = $npmToken }
    Write-Host "Publishing $PackageName@$Version to npm..."
    & npm publish $tarballPath --access public --registry "https://registry.npmjs.org" @npmAuthArgs
    Assert-LastExitCode "Publish npm package"
    Write-Host "Release $Version completed"
} finally {
    Remove-Item -LiteralPath $awsConfigPath -Force -ErrorAction SilentlyContinue
    Remove-Item -LiteralPath $npmConfigPath -Force -ErrorAction SilentlyContinue
    $env:AWS_CONFIG_FILE = $originalAwsConfig
    $env:AWS_ACCESS_KEY_ID = $originalAwsKey
    $env:AWS_SECRET_ACCESS_KEY = $originalAwsSecret
    $env:AWS_DEFAULT_REGION = $originalAwsRegion
    $env:AWS_EC2_METADATA_DISABLED = $originalMetadataDisabled
    $env:NODE_AUTH_TOKEN = $originalNodeAuthToken
}
