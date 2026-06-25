param(
    [Parameter(Mandatory = $true)]
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
$stagingRoot = Join-Path $repoRoot "out\npm-staging"
$packageRoot = Join-Path $stagingRoot "package"
$distDir = Join-Path $packageRoot "dist"
$packagedSkillRoot = Join-Path $packageRoot "skills"
$resolvedOutputDir = Join-Path $repoRoot $OutputDir
$goCacheDir = Join-Path $repoRoot "out\go-build-cache"
$npmCacheDir = Join-Path $repoRoot "out\npm-cache"

function Assert-LastExitCode {
    param(
        [string]$CommandName
    )

    if ($LASTEXITCODE -ne 0) {
        throw "$CommandName failed with exit code $LASTEXITCODE"
    }
}

if (-not (Test-Path $npmTemplateDir)) {
    throw "npm template directory not found: $npmTemplateDir"
}
if (-not (Test-Path $skillSourceDir)) {
    throw "skill source directory not found: $skillSourceDir"
}

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
