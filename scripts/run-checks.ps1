param(
    [switch]$GoOnly,
    [switch]$NpmOnly
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
$workspaceRoot = Join-Path $repoRoot "out\test-workspace-$PID"
$goTestRoot = Join-Path $repoRoot "test\go"
$fixtureRoot = Join-Path $repoRoot "test\fixtures"
$goCacheDir = Join-Path $repoRoot "out\go-test-cache"
$npmCacheDir = Join-Path $repoRoot "out\npm-test-cache"

function Assert-LastExitCode {
    param(
        [string]$CommandName
    )

    if ($LASTEXITCODE -ne 0) {
        throw "$CommandName failed with exit code $LASTEXITCODE"
    }
}

function Copy-RepositoryForTests {
    if (Test-Path $workspaceRoot) {
        Remove-Item -LiteralPath $workspaceRoot -Recurse -Force
    }

    New-Item -ItemType Directory -Path $workspaceRoot -Force | Out-Null

    $excludedRoots = @(
        (Join-Path $repoRoot ".git"),
        (Join-Path $repoRoot ".agents"),
        (Join-Path $repoRoot ".codex"),
        (Join-Path $repoRoot "node_modules"),
        (Join-Path $repoRoot "out"),
        (Join-Path $repoRoot "test")
    )

    Get-ChildItem -LiteralPath $repoRoot -Force | ForEach-Object {
        $sourcePath = $_.FullName
        if ($excludedRoots -contains $sourcePath) {
            return
        }

        Copy-Item -LiteralPath $sourcePath -Destination $workspaceRoot -Recurse -Force
    }

    if (Test-Path $fixtureRoot) {
        $workspaceTestRoot = Join-Path $workspaceRoot "test"
        New-Item -ItemType Directory -Path $workspaceTestRoot -Force | Out-Null
        Copy-Item -LiteralPath $fixtureRoot -Destination $workspaceTestRoot -Recurse -Force
    }
}

function Copy-GoTestsIntoWorkspace {
    if (-not (Test-Path $goTestRoot)) {
        return
    }

    Get-ChildItem -LiteralPath $goTestRoot -Recurse -File | ForEach-Object {
        $relativePath = $_.FullName.Substring($goTestRoot.Length).TrimStart('\', '/')
        $destination = Join-Path $workspaceRoot $relativePath
        New-Item -ItemType Directory -Path (Split-Path -Parent $destination) -Force | Out-Null
        Copy-Item -LiteralPath $_.FullName -Destination $destination -Force
    }
}

if (-not $NpmOnly) {
    New-Item -ItemType Directory -Path $goCacheDir -Force | Out-Null
    $env:GOCACHE = $goCacheDir

    Copy-RepositoryForTests
    Copy-GoTestsIntoWorkspace

    Push-Location $workspaceRoot
    try {
        go test ./...
        Assert-LastExitCode "go test ./..."
    } finally {
        Pop-Location
    }
}

if (-not $GoOnly) {
    New-Item -ItemType Directory -Path $npmCacheDir -Force | Out-Null
    $env:npm_config_cache = $npmCacheDir

    Push-Location (Join-Path $repoRoot "npm")
    try {
        npm test
        Assert-LastExitCode "npm test"
    } finally {
        Pop-Location
    }
}
