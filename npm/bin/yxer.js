#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

function resolveBinary() {
  const platform = os.platform();
  const arch = os.arch();
  const key = `${platform}/${arch}`;

  const binaries = {
    "win32/x64": "yxer-windows-amd64.exe",
    "win32/arm64": "yxer-windows-arm64.exe",
    "darwin/x64": "yxer-darwin-amd64",
    "darwin/arm64": "yxer-darwin-arm64",
    "linux/x64": "yxer-linux-amd64",
    "linux/arm64": "yxer-linux-arm64"
  };

  const filename = binaries[key];
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

  const binaryPath = path.join(__dirname, "..", "dist", filename);
  if (!fs.existsSync(binaryPath)) {
    console.error(
      JSON.stringify({
        error: {
          code: "missing_binary",
          message: `Expected packaged binary not found: ${filename}`,
          category: "environment",
          hint: "Reinstall the npm package or rebuild the release artifact.",
          retryable: false
        }
      })
    );
    process.exit(1);
  }

  return binaryPath;
}

const binary = resolveBinary();
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
