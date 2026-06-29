param(
    [string]$Version,

    [string]$PackageName = "@yixiaoermail/cli",

    [string]$OutputDir = "out\npm",

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
$distDir = Join-Path $packageRoot "dist"
$packagedSkillRoot = Join-Path $packageRoot "skills"
$resolvedOutputDir = Join-Path $repoRoot $OutputDir
$goCacheDir = Join-Path $repoRoot "out\go-build-cache"
$npmCacheDir = Join-Path $repoRoot "out\npm-cache"
$goVersionSourcePath = Join-Path $repoRoot "internal\domain\response.go"
$skillManifestPath = Join-Path $repoRoot "skills\yixiaoer\SKILL.md"
$pluginManifestPath = Join-Path $repoRoot "skills\yixiaoer\plugin.json"

function Assert-LastExitCode {
    param(
        [string]$CommandName
    )

    if ($LASTEXITCODE -ne 0) {
        throw "$CommandName failed with exit code $LASTEXITCODE"
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

function Get-PluginManifestVersion {
    param(
        [string]$Path
    )

    $json = Get-Content -LiteralPath $Path -Raw | ConvertFrom-Json
    if (-not $json.version) {
        throw "version field not found in plugin manifest: $Path"
    }

    return [string]$json.version
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
$pluginVersion = Get-PluginManifestVersion -Path $pluginManifestPath

$detectedVersions = @(
    @{ Name = "internal/domain/response.go"; Version = $goVersion },
    @{ Name = "skills/yixiaoer/SKILL.md"; Version = $skillVersion },
    @{ Name = "skills/yixiaoer/plugin.json"; Version = $pluginVersion }
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

New-Item -ItemType Directory -Path $goCacheDir -Force | Out-Null
New-Item -ItemType Directory -Path $npmCacheDir -Force | Out-Null
$env:GOCACHE = $goCacheDir
$env:npm_config_cache = $npmCacheDir

if (-not $SkipTests) {
    Write-Host "Running go test ./..."
    go test ./...
    Assert-LastExitCode "go test ./..."
}

if (Test-Path $stagingRoot) {
    Remove-Item -LiteralPath $stagingRoot -Recurse -Force
}

New-Item -ItemType Directory -Path $packageRoot -Force | Out-Null
New-Item -ItemType Directory -Path $distDir -Force | Out-Null
New-Item -ItemType Directory -Path $packagedSkillRoot -Force | Out-Null
New-Item -ItemType Directory -Path $resolvedOutputDir -Force | Out-Null

Copy-Item -Path (Join-Path $npmTemplateDir "*") -Destination $packageRoot -Recurse -Force
Copy-Item -Path $skillSourceDir -Destination $packagedSkillRoot -Recurse -Force
Copy-Item -Path $schemaSourceDir -Destination $packageRoot -Recurse -Force
Copy-Item -Path $referencesSourceDir -Destination $packageRoot -Recurse -Force

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Output = "yxer-windows-amd64.exe" },
    @{ GOOS = "windows"; GOARCH = "arm64"; Output = "yxer-windows-arm64.exe" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Output = "yxer-darwin-amd64" },
    @{ GOOS = "darwin"; GOARCH = "arm64"; Output = "yxer-darwin-arm64" },
    @{ GOOS = "linux"; GOARCH = "amd64"; Output = "yxer-linux-amd64" },
    @{ GOOS = "linux"; GOARCH = "arm64"; Output = "yxer-linux-arm64" }
)

foreach ($target in $targets) {
    $outputPath = Join-Path $distDir $target.Output
    Write-Host "Building $($target.GOOS)/$($target.GOARCH) -> $outputPath"

    $env:GOOS = $target.GOOS
    $env:GOARCH = $target.GOARCH
    go build -buildvcs=false -o $outputPath .
    Assert-LastExitCode "go build ($($target.GOOS)/$($target.GOARCH))"
}

$packageJsonPath = Join-Path $packageRoot "package.json"
$packageJson = Get-Content -LiteralPath $packageJsonPath -Raw | ConvertFrom-Json
$packageJson.version = $Version
$packageJson.name = $PackageName
$packageJson | ConvertTo-Json -Depth 10 | Set-Content -LiteralPath $packageJsonPath -Encoding utf8

Write-Host "Packing npm artifact"
$packOutput = npm pack $packageRoot --pack-destination $resolvedOutputDir
Assert-LastExitCode "npm pack"

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
Remove-Item Env:GOCACHE -ErrorAction SilentlyContinue
Remove-Item Env:npm_config_cache -ErrorAction SilentlyContinue

Write-Host "Generated package:"
$packOutput | ForEach-Object { Write-Host $_ }
