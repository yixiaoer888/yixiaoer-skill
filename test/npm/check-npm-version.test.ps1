$ErrorActionPreference = "Stop"
$checkScript = Join-Path $PSScriptRoot "..\..\scripts\check-npm-version.ps1"

function Assert-Throws {
    param([scriptblock]$Action, [string]$ExpectedMessage)
    try {
        & $Action
    } catch {
        if ($_.Exception.Message -notlike "*$ExpectedMessage*") {
            throw "Expected error containing '$ExpectedMessage', got '$($_.Exception.Message)'"
        }
        return
    }
    throw "Expected error containing '$ExpectedMessage'"
}

function npm {
    & node -e "console.error('npm error code E404'); process.exit(1)"
}
& $checkScript -PackageName "@yixiaoermail/cli" -Version "3.2.25"

function npm {
    $global:LASTEXITCODE = 0
    "3.2.25"
}
Assert-Throws { & $checkScript -PackageName "@yixiaoermail/cli" -Version "3.2.25" } "already exists on npm"

function npm {
    & node -e "console.error('npm error code E401'); process.exit(1)"
}
Assert-Throws { & $checkScript -PackageName "@yixiaoermail/cli" -Version "3.2.25" } "npm version check failed"

$global:LASTEXITCODE = 0
Write-Host "npm version preflight tests passed"
