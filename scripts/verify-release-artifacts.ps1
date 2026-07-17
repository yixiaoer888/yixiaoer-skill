param(
    [string]$ReleaseDir = "out\release",
    [string]$TarballDir = "",
    [string]$TarballPattern = "*.tgz"
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

function Get-FileSha256 {
    param(
        [string]$Path
    )

    return (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
}

function Get-ChecksumMap {
    param(
        [string]$Path
    )

    if (-not (Test-Path -LiteralPath $Path)) {
        throw "checksums.txt not found: $Path"
    }

    $map = @{}
    foreach ($line in Get-Content -LiteralPath $Path) {
        $trimmed = $line.Trim()
        if (-not $trimmed) {
            continue
        }

        $parts = $trimmed -split "\s{2,}", 2
        if ($parts.Count -ne 2) {
            throw "Invalid checksum line in ${Path}: $line"
        }

        $map[$parts[1]] = $parts[0].ToLowerInvariant()
    }

    return $map
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$resolvedReleaseDir = if ([System.IO.Path]::IsPathRooted($ReleaseDir)) {
    $ReleaseDir
} else {
    Join-Path $repoRoot $ReleaseDir
}

$resolvedTarballDir = if ([string]::IsNullOrWhiteSpace($TarballDir)) {
    $resolvedReleaseDir
} elseif ([System.IO.Path]::IsPathRooted($TarballDir)) {
    $TarballDir
} else {
    Join-Path $repoRoot $TarballDir
}

if (-not (Test-Path -LiteralPath $resolvedReleaseDir)) {
    throw "Release directory not found: $resolvedReleaseDir"
}

$tarballs = Get-ChildItem -LiteralPath $resolvedTarballDir -Filter $TarballPattern -File
if ($tarballs.Count -ne 1) {
    throw "Expected exactly 1 tarball matching '$TarballPattern' in $resolvedTarballDir, found $($tarballs.Count)"
}

$releaseChecksumsPath = Join-Path $resolvedReleaseDir "checksums.txt"
$releaseChecksums = Get-ChecksumMap -Path $releaseChecksumsPath

foreach ($archiveName in $releaseChecksums.Keys) {
    $archivePath = Join-Path $resolvedReleaseDir $archiveName
    if (-not (Test-Path -LiteralPath $archivePath)) {
        throw "Release archive missing: $archiveName"
    }

    $actualHash = Get-FileSha256 -Path $archivePath
    if ($actualHash -ne $releaseChecksums[$archiveName]) {
        throw "Checksum mismatch for ${archiveName}: expected $($releaseChecksums[$archiveName]) but got $actualHash"
    }
}

$tmpDir = Join-Path ([System.IO.Path]::GetTempPath()) ("yxer-release-verify-" + [System.Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

try {
    tar -xzf $tarballs[0].FullName -C $tmpDir
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to extract tarball: $($tarballs[0].FullName)"
    }

    $embeddedChecksumsPath = Join-Path $tmpDir "package\checksums.txt"
    $embeddedChecksums = Get-ChecksumMap -Path $embeddedChecksumsPath

    if ($embeddedChecksums.Count -ne $releaseChecksums.Count) {
        throw "Embedded checksum count ($($embeddedChecksums.Count)) does not match release checksum count ($($releaseChecksums.Count))"
    }

    foreach ($archiveName in $releaseChecksums.Keys) {
        if (-not $embeddedChecksums.ContainsKey($archiveName)) {
            throw "Embedded checksums missing archive entry: $archiveName"
        }

        if ($embeddedChecksums[$archiveName] -ne $releaseChecksums[$archiveName]) {
            throw "Embedded checksum mismatch for ${archiveName}: tarball has $($embeddedChecksums[$archiveName]) but release has $($releaseChecksums[$archiveName])"
        }
    }

    Write-Host "Verified release artifacts and embedded checksums for $($tarballs[0].Name)"
} finally {
    if (Test-Path -LiteralPath $tmpDir) {
        Remove-Item -LiteralPath $tmpDir -Recurse -Force
    }
}
