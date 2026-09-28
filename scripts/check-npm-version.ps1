param(
    [Parameter(Mandatory = $true)][string]$PackageName,
    [Parameter(Mandatory = $true)][string]$Version,
    [string[]]$NpmAuthArgs = @()
)

$ErrorActionPreference = "Stop"
$originalErrorActionPreference = $ErrorActionPreference
try {
    # Windows PowerShell turns native stderr into error records. A missing npm
    # version writes E404 to stderr, which is expected during this preflight.
    $ErrorActionPreference = "Continue"
    $npmOutput = & npm view "$PackageName@$Version" version --registry "https://registry.npmjs.org" @NpmAuthArgs 2>&1
    $npmExitCode = $LASTEXITCODE
} finally {
    $ErrorActionPreference = $originalErrorActionPreference
}

$diagnostic = ($npmOutput | Out-String).Trim()
if ($npmExitCode -eq 0) {
    throw "$PackageName@$Version already exists on npm and cannot be published again"
}
if ($diagnostic -notmatch '(?i)\bnpm\s+(?:error|ERR!)\s+code\s+E404\b') {
    throw "npm version check failed with exit code $npmExitCode`: $diagnostic"
}
