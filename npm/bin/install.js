#!/usr/bin/env node

const crypto = require("node:crypto");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const { execFileSync } = require("node:child_process");
const { ensureExecutable } = require("./ensure-executable");

const PACKAGE_NAME = "yxer";
const ARCHIVE_PREFIX = "yxer-cli";
const VERSION = require("../package.json").version.replace(/-.*$/, "");
const DEFAULT_REPO = "yixiaoer888/yixiaoer-skill";
const DEFAULT_BASE_URL = `https://github.com/${DEFAULT_REPO}/releases/download/v${VERSION}`;
const ALLOWED_HOSTS = new Set([
  "github.com",
  "objects.githubusercontent.com",
  "release-assets.githubusercontent.com",
  "registry.npmmirror.com"
]);

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows"
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64"
};

function getTarget(platform = process.platform, arch = process.arch) {
  const mappedPlatform = PLATFORM_MAP[platform];
  const mappedArch = ARCH_MAP[arch];
  if (!mappedPlatform || !mappedArch) {
    return null;
  }

  const ext = platform === "win32" ? "zip" : "tar.gz";
  const binaryExt = platform === "win32" ? ".exe" : "";
  const baseName = `${ARCHIVE_PREFIX}-${VERSION}-${mappedPlatform}-${mappedArch}`;

  return {
    platform,
    arch,
    mappedPlatform,
    mappedArch,
    archiveName: `${baseName}.${ext}`,
    binaryName: `${PACKAGE_NAME}${binaryExt}`
  };
}

function isSupportedPlatform(platform = process.platform, arch = process.arch) {
  return getTarget(platform, arch) !== null;
}

function getBinDir(baseDir = __dirname) {
  return path.join(baseDir, "..", "bin-native");
}

function getChecksumsPath(baseDir = __dirname) {
  return path.join(baseDir, "..", "checksums.txt");
}

function getBinaryPath(baseDir = __dirname, platform = process.platform, arch = process.arch) {
  const target = getTarget(platform, arch);
  if (!target) {
    return null;
  }

  return path.join(getBinDir(baseDir), target.binaryName);
}

function assertAllowedHost(url) {
  const hostname = new URL(url).hostname;
  if (!ALLOWED_HOSTS.has(hostname)) {
    throw new Error(`Download host not allowed: ${hostname}`);
  }
}

function buildDownloadUrl(baseUrl, archiveName) {
  return `${baseUrl.replace(/\/+$/, "")}/${archiveName}`;
}

function downloadWithCurl(url, destPath, execFile = execFileSync) {
  assertAllowedHost(url);

  const args = [
    "--fail",
    "--location",
    "--silent",
    "--show-error",
    "--connect-timeout",
    "10",
    "--max-time",
    "120",
    "--max-redirs",
    "3",
    "--output",
    destPath
  ];

  if (process.platform === "win32") {
    args.unshift("--ssl-revoke-best-effort");
  }

  args.push(url);
  execFile("curl", args, { stdio: ["ignore", "ignore", "pipe"] });
}

function downloadWithPowerShell(url, destPath, execFile = execFileSync) {
  assertAllowedHost(url);

  const psCommand =
    "$ProgressPreference='SilentlyContinue';" +
    "$ErrorActionPreference='Stop';" +
    "[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12;" +
    "Invoke-WebRequest -UseBasicParsing -Uri $env:YXER_URL -OutFile $env:YXER_DEST";

  execFile(
    "powershell.exe",
    ["-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand],
    {
      stdio: ["ignore", "ignore", "pipe"],
      env: {
        ...process.env,
        YXER_URL: url,
        YXER_DEST: destPath
      }
    }
  );
}

function download(url, destPath, platform = process.platform, execFile = execFileSync) {
  try {
    downloadWithCurl(url, destPath, execFile);
  } catch (error) {
    if (platform !== "win32") {
      throw error;
    }

    try {
      downloadWithPowerShell(url, destPath, execFile);
    } catch (fallbackError) {
      fallbackError.message =
        `${error.message}; fallback via PowerShell failed: ${fallbackError.message}`;
      throw fallbackError;
    }
  }
}

function parseChecksums(content) {
  const result = new Map();
  for (const line of content.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed) {
      continue;
    }

    const separator = trimmed.indexOf("  ");
    if (separator === -1) {
      continue;
    }

    result.set(trimmed.slice(separator + 2), trimmed.slice(0, separator));
  }

  return result;
}

function getExpectedChecksum(archiveName, checksumsPath = getChecksumsPath()) {
  if (!fs.existsSync(checksumsPath)) {
    return null;
  }

  const checksums = parseChecksums(fs.readFileSync(checksumsPath, "utf8"));
  return checksums.get(archiveName) || null;
}

function verifyChecksum(archivePath, expectedHash) {
  if (!expectedHash) {
    return;
  }

  const hash = crypto.createHash("sha256");
  const fd = fs.openSync(archivePath, "r");

  try {
    const buffer = Buffer.alloc(64 * 1024);
    let bytesRead = 0;
    while ((bytesRead = fs.readSync(fd, buffer, 0, buffer.length, null)) > 0) {
      hash.update(buffer.subarray(0, bytesRead));
    }
  } finally {
    fs.closeSync(fd);
  }

  const actualHash = hash.digest("hex");
  if (actualHash.toLowerCase() !== expectedHash.toLowerCase()) {
    throw new Error(
      `Checksum mismatch for ${path.basename(archivePath)}: expected ${expectedHash} but got ${actualHash}`
    );
  }
}

function extractArchive(archivePath, destDir) {
  if (process.platform === "win32") {
    const psCommand =
      "$ErrorActionPreference='Stop';" +
      "Expand-Archive -LiteralPath $env:YXER_ARCHIVE -DestinationPath $env:YXER_DEST -Force";
    execFileSync(
      "powershell.exe",
      ["-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", psCommand],
      {
        stdio: ["ignore", "inherit", "inherit"],
        env: {
          ...process.env,
          YXER_ARCHIVE: archivePath,
          YXER_DEST: destDir
        }
      }
    );
    return;
  }

  execFileSync("tar", ["-xzf", archivePath, "-C", destDir], { stdio: "ignore" });
}

function resolveBaseUrl() {
  const customUrl = (process.env.YXER_DOWNLOAD_BASE_URL || "").trim();
  if (!customUrl) {
    return DEFAULT_BASE_URL;
  }

  ALLOWED_HOSTS.add(new URL(customUrl).hostname);
  return customUrl;
}

function install(baseDir = __dirname) {
  const target = getTarget();
  if (!target) {
    throw new Error(`Unsupported platform: ${process.platform}/${process.arch}`);
  }

  const binDir = getBinDir(baseDir);
  const binaryPath = path.join(binDir, target.binaryName);
  fs.mkdirSync(binDir, { recursive: true });

  if (fs.existsSync(binaryPath)) {
    ensureExecutable(binaryPath);
    return binaryPath;
  }

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-install-"));
  const archivePath = path.join(tmpDir, target.archiveName);

  try {
    const downloadUrl = buildDownloadUrl(resolveBaseUrl(), target.archiveName);
    download(downloadUrl, archivePath);
    verifyChecksum(archivePath, getExpectedChecksum(target.archiveName, getChecksumsPath(baseDir)));
    extractArchive(archivePath, tmpDir);

    const extractedBinary = path.join(tmpDir, target.binaryName);
    if (!fs.existsSync(extractedBinary)) {
      throw new Error(`Extracted binary not found: ${target.binaryName}`);
    }

    fs.copyFileSync(extractedBinary, binaryPath);
    ensureExecutable(binaryPath);
    return binaryPath;
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

if (require.main === module) {
  try {
    install();
  } catch (error) {
    console.error(`Failed to install yxer: ${error.message}`);
    process.exit(1);
  }
}

module.exports = {
  ALLOWED_HOSTS,
  buildDownloadUrl,
  download,
  downloadWithCurl,
  downloadWithPowerShell,
  getBinaryPath,
  getBinDir,
  getExpectedChecksum,
  getTarget,
  install,
  isSupportedPlatform,
  parseChecksums,
  verifyChecksum
};
