const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

function getBinaryFilename(platform = os.platform(), arch = os.arch()) {
  const key = `${platform}/${arch}`;
  const binaries = {
    "win32/x64": "yxer-windows-amd64.exe",
    "win32/arm64": "yxer-windows-arm64.exe",
    "darwin/x64": "yxer-darwin-amd64",
    "darwin/arm64": "yxer-darwin-arm64",
    "linux/x64": "yxer-linux-amd64",
    "linux/arm64": "yxer-linux-arm64"
  };

  return binaries[key] || null;
}

function resolveBinaryPath(baseDir = __dirname, platform = os.platform(), arch = os.arch()) {
  const filename = getBinaryFilename(platform, arch);
  if (!filename) {
    return null;
  }

  const binaryPath = path.join(baseDir, "..", "dist", filename);
  if (!fs.existsSync(binaryPath)) {
    return null;
  }

  return binaryPath;
}

module.exports = {
  getBinaryFilename,
  resolveBinaryPath
};
