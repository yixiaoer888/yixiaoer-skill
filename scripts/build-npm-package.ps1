param(
    [string]$Version,

    [string]$PackageName = "@yixiaoermail/cli",

    [string]$OutputDir = "out\npm",

    [string]$ReleaseDir = "out\release",

    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

$repoRoot = Split-Path -Parent $PSScriptRoot
$npmTemplateDir = Join-Path $repoRoot "npm"
$skillSourceDir = Join-Path $repoRoot "skills\yixiaoer"
$schemaSourceDir = Join-Path $repoRoot "schemas"
$referencesSourceDir = Join-Path $repoRoot "references"
$stagingRoot = Join-Path $repoRoot "out\npm-staging"
$packageRoot = Join-Path $stagingRoot "package"
$packagedSkillRoot = Join-Path $packageRoot "skills"
$nativeBinDir = Join-Path $packageRoot "bin-native"
$resolvedOutputDir = Join-Path $repoRoot $OutputDir
$resolvedReleaseDir = Join-Path $repoRoot $ReleaseDir
$goCacheDir = Join-Path $repoRoot "out\go-build-cache"
$npmCacheDir = Join-Path $repoRoot "out\npm-cache"
$goVersionSourcePath = Join-Path $repoRoot "internal\domain\response.go"
$skillManifestPath = Join-Path $repoRoot "skills\yixiaoer\SKILL.md"

function Assert-LastExitCode {
    param(
        [string]$CommandName
    )

    if ($LASTEXITCODE -ne 0) {
        throw "$CommandName failed with exit code $LASTEXITCODE"
    }
}

function Restore-Environment {
    param(
        [hashtable]$OriginalEnvironment
    )

    foreach ($name in $OriginalEnvironment.Keys) {
        if ($null -eq $OriginalEnvironment[$name]) {
            Remove-Item -Path "Env:$name" -ErrorAction SilentlyContinue
        } else {
            Set-Item -Path "Env:$name" -Value $OriginalEnvironment[$name]
        }
    }
}

function Get-GoSkillVersion {
    param(
        [string]$Path
    )

    $content = Get-Content -LiteralPath $Path -Raw
    $match = [regex]::Match($content, 'const\s+SkillVersion\s*=\s*"([^"]+)"')
    if (-not $match.Success) {
        throw "SkillVersion constant not found: $Path"
    }

    return $match.Groups[1].Value.Trim()
}

function Get-SkillManifestVersion {
    param(
        [string]$Path
    )

    foreach ($line in Get-Content -LiteralPath $Path) {
        $trimmed = $line.Trim()
        if ($trimmed -like "version:*") {
            return $trimmed.Substring("version:".Length).Trim().Trim('"').Trim("'")
        }
    }

    throw "version field not found in skill manifest: $Path"
}

function New-ArchiveFromBinary {
    param(
        [string]$BinaryPath,
        [string]$ArchivePath,
        [bool]$TargetIsWindows
    )

    $archiveDir = Split-Path -Parent $ArchivePath
    New-Item -ItemType Directory -Path $archiveDir -Force | Out-Null

    if ($TargetIsWindows) {
        if (Test-Path $ArchivePath) {
            Remove-Item -LiteralPath $ArchivePath -Force
        }

        Compress-Archive -LiteralPath $BinaryPath -DestinationPath $ArchivePath -CompressionLevel Optimal
        return
    }

    tar -czf $ArchivePath -C (Split-Path -Parent $BinaryPath) (Split-Path -Leaf $BinaryPath)
    Assert-LastExitCode "tar -czf $ArchivePath"
}

function Get-FileSha256 {
    param(
        [string]$Path
    )

    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Get-ArchiveName {
    param(
        [string]$Version,
        [string]$Platform,
        [string]$Arch
    )

    $extension = if ($Platform -eq "windows") { "zip" } else { "tar.gz" }
    return "yxer-cli-$Version-$Platform-$Arch.$extension"
}

if (-not (Test-Path $npmTemplateDir)) {
    throw "npm template directory not found: $npmTemplateDir"
}
if (-not (Test-Path $skillSourceDir)) {
    throw "skill source directory not found: $skillSourceDir"
}
if (-not (Test-Path $schemaSourceDir)) {
    throw "schema source directory not found: $schemaSourceDir"
}
if (-not (Test-Path $referencesSourceDir)) {
    throw "references source directory not found: $referencesSourceDir"
}

$goVersion = Get-GoSkillVersion -Path $goVersionSourcePath
$skillVersion = Get-SkillManifestVersion -Path $skillManifestPath

$detectedVersions = @(
    @{ Name = "internal/domain/response.go"; Version = $goVersion },
    @{ Name = "skills/yixiaoer/SKILL.md"; Version = $skillVersion }
)

$distinctVersions = $detectedVersions.Version | Sort-Object -Unique
if ($distinctVersions.Count -ne 1) {
    $details = ($detectedVersions | ForEach-Object { "$($_.Name)=$($_.Version)" }) -join ", "
    throw "Version sources are inconsistent: $details"
}

$resolvedVersion = $goVersion
if ($Version) {
    if ($Version -ne $resolvedVersion) {
        throw "Provided version '$Version' does not match internal version '$resolvedVersion'"
    }
} else {
    $Version = $resolvedVersion
}

Write-Host "Using package version $Version"

$trackedEnvironment = @("GOOS", "GOARCH", "GOCACHE", "npm_config_cache")
$originalEnvironment = @{}
foreach ($name in $trackedEnvironment) {
    $originalEnvironment[$name] = [Environment]::GetEnvironmentVariable($name, "Process")
}

try {
    New-Item -ItemType Directory -Path $goCacheDir -Force | Out-Null
    New-Item -ItemType Directory -Path $npmCacheDir -Force | Out-Null
    $env:GOCACHE = $goCacheDir
    $env:npm_config_cache = $npmCacheDir

    if (-not $SkipTests) {
        $hostGoOS = (go env GOHOSTOS).Trim()
        Assert-LastExitCode "go env GOHOSTOS"
        $hostGoArch = (go env GOHOSTARCH).Trim()
        Assert-LastExitCode "go env GOHOSTARCH"
        $env:GOOS = $hostGoOS
        $env:GOARCH = $hostGoArch

        Write-Host "Running go tests for $hostGoOS/$hostGoArch"
        go test ./...
        Assert-LastExitCode "go test ./..."

        Write-Host "Running npm installer tests"
        node --test .\test\npm\ensure-executable.test.js
        Assert-LastExitCode "node --test .\test\npm\ensure-executable.test.js"
    }

    if (Test-Path $stagingRoot) {
        Remove-Item -LiteralPath $stagingRoot -Recurse -Force
    }
    if (Test-Path $resolvedReleaseDir) {
        Remove-Item -LiteralPath $resolvedReleaseDir -Recurse -Force
    }

    New-Item -ItemType Directory -Path $packageRoot -Force | Out-Null
    New-Item -ItemType Directory -Path $packagedSkillRoot -Force | Out-Null
    New-Item -ItemType Directory -Path $nativeBinDir -Force | Out-Null
    New-Item -ItemType Directory -Path $resolvedOutputDir -Force | Out-Null
    New-Item -ItemType Directory -Path $resolvedReleaseDir -Force | Out-Null

    Copy-Item -Path (Join-Path $npmTemplateDir "*") -Destination $packageRoot -Recurse -Force
    Copy-Item -Path $skillSourceDir -Destination $packagedSkillRoot -Recurse -Force
    Copy-Item -Path $schemaSourceDir -Destination $packageRoot -Recurse -Force
    Copy-Item -Path $referencesSourceDir -Destination $packageRoot -Recurse -Force

    $targets = @(
        @{ GOOS = "windows"; GOARCH = "amd64"; BinaryName = "yxer.exe" },
        @{ GOOS = "windows"; GOARCH = "arm64"; BinaryName = "yxer.exe" },
        @{ GOOS = "darwin"; GOARCH = "amd64"; BinaryName = "yxer" },
        @{ GOOS = "darwin"; GOARCH = "arm64"; BinaryName = "yxer" },
        @{ GOOS = "linux"; GOARCH = "amd64"; BinaryName = "yxer" },
        @{ GOOS = "linux"; GOARCH = "arm64"; BinaryName = "yxer" }
    )

    $checksumLines = New-Object System.Collections.Generic.List[string]

    foreach ($target in $targets) {
        $archiveName = Get-ArchiveName -Version $Version -Platform $target.GOOS -Arch $target.GOARCH
        $buildDir = Join-Path $resolvedReleaseDir "$($target.GOOS)-$($target.GOARCH)"
        $binaryPath = Join-Path $buildDir $target.BinaryName
        $archivePath = Join-Path $resolvedReleaseDir $archiveName

        New-Item -ItemType Directory -Path $buildDir -Force | Out-Null

        Write-Host "Building $($target.GOOS)/$($target.GOARCH) -> $binaryPath"
        $env:GOOS = $target.GOOS
        $env:GOARCH = $target.GOARCH
        go build -buildvcs=false -o $binaryPath .
        Assert-LastExitCode "go build ($($target.GOOS)/$($target.GOARCH))"

        New-ArchiveFromBinary -BinaryPath $binaryPath -ArchivePath $archivePath -TargetIsWindows:($target.GOOS -eq "windows")
        $checksumLines.Add("$(Get-FileSha256 -Path $archivePath)  $archiveName")
    }

    $checksumsPath = Join-Path $packageRoot "checksums.txt"
    $checksumLines | Set-Content -LiteralPath $checksumsPath -Encoding ascii
    $checksumLines | Set-Content -LiteralPath (Join-Path $resolvedReleaseDir "checksums.txt") -Encoding ascii

    $packageJsonPath = Join-Path $packageRoot "package.json"
    $packageJson = Get-Content -LiteralPath $packageJsonPath -Raw | ConvertFrom-Json
    $packageJson.version = $Version
    $packageJson.name = $PackageName
    $packageJson | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $packageJsonPath -Encoding utf8

    Write-Host "Packing npm artifact"
    $packOutput = npm pack $packageRoot --pack-destination $resolvedOutputDir
    Assert-LastExitCode "npm pack"

    Write-Host "Generated release archives:"
    Get-ChildItem -LiteralPath $resolvedReleaseDir -File | ForEach-Object { Write-Host $_.Name }

    Write-Host "Generated npm package:"
    $packOutput | ForEach-Object { Write-Host $_ }
} finally {
    Restore-Environment -OriginalEnvironment $originalEnvironment
}
