#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const fs = require("node:fs");
const { ensureExecutable } = require("./ensure-executable");
const { install } = require("./install");
const { getBinaryFilename, resolveBinaryPath } = require("./resolve-binary");

function resolveBinary() {
  const platform = process.platform;
  const arch = process.arch;
  const key = `${platform}/${arch}`;
  const filename = getBinaryFilename(platform, arch);
  if (!filename) {
    console.error(
      JSON.stringify({
        error: {
          code: "unsupported_platform",
          message: `No packaged yxer binary for ${key}`,
          category: "environment",
          retryable: false
        }
      })
    );
    process.exit(1);
  }

  const binaryPath = resolveBinaryPath(__dirname, platform, arch);
  if (!binaryPath || !fs.existsSync(binaryPath)) {
    try {
      return install(__dirname);
    } catch (error) {
      console.error(
        JSON.stringify({
          error: {
            code: "missing_binary",
            message: `Expected packaged binary not found: ${filename}`,
            category: "environment",
            hint: `Automatic install failed: ${error.message}`,
            nextCommand: "npm install -g @yixiaoermail/cli@latest",
            retryable: true
          }
        })
      );
      process.exit(1);
    }
  }

  return binaryPath;
}

const binary = resolveBinary();
ensureExecutable(binary);
const result = spawnSync(binary, process.argv.slice(2), { stdio: "inherit" });

if (result.error) {
  console.error(
    JSON.stringify({
      error: {
        code: "spawn_failed",
        message: result.error.message,
        category: "environment",
        hint: "Check that the packaged binary has execute permission and matches the current platform.",
        retryable: false
      }
    })
  );
  process.exit(1);
}

process.exit(result.status === null ? 1 : result.status);
