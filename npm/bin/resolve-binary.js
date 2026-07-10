const fs = require("node:fs");
const os = require("node:os");
const { getBinaryPath, getTarget } = require("./install");

function getBinaryFilename(platform = os.platform(), arch = os.arch()) {
  const target = getTarget(platform, arch);
  return target ? target.binaryName : null;
}

function resolveBinaryPath(baseDir = __dirname, platform = os.platform(), arch = os.arch()) {
  const binaryPath = getBinaryPath(baseDir, platform, arch);
  if (!binaryPath || !fs.existsSync(binaryPath)) {
    return null;
  }

  return binaryPath;
}

module.exports = {
  getBinaryFilename,
  resolveBinaryPath
};
