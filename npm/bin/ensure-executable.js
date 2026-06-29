const fs = require("node:fs");
const os = require("node:os");

function shouldEnsureExecutable(platform) {
  return platform !== "win32";
}

function ensureExecutable(binaryPath) {
  if (!shouldEnsureExecutable(os.platform())) {
    return;
  }

  const stat = fs.statSync(binaryPath);
  const executeBits = 0o111;
  if ((stat.mode & executeBits) === executeBits) {
    return;
  }

  fs.chmodSync(binaryPath, stat.mode | 0o755);
}

module.exports = {
  ensureExecutable,
  shouldEnsureExecutable
};
