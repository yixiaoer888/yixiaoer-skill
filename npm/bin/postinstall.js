#!/usr/bin/env node

const { ensureExecutable } = require("./ensure-executable");
const { resolveBinaryPath } = require("./resolve-binary");

try {
  const binaryPath = resolveBinaryPath(__dirname, process.platform, process.arch);
  if (!binaryPath) {
    process.exit(0);
  }

  ensureExecutable(binaryPath);
} catch (error) {
  console.error(
    `[yxer] postinstall failed to ensure executable permission: ${error.message}`
  );
  process.exit(0);
}
