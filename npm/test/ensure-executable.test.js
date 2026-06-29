const test = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

const { ensureExecutable, shouldEnsureExecutable } = require("../bin/ensure-executable");
const { getBinaryFilename, resolveBinaryPath } = require("../bin/resolve-binary");

test("shouldEnsureExecutable skips windows only", () => {
  assert.equal(shouldEnsureExecutable("win32"), false);
  assert.equal(shouldEnsureExecutable("darwin"), true);
  assert.equal(shouldEnsureExecutable("linux"), true);
});

test("getBinaryFilename resolves mac arm64 binary", () => {
  assert.equal(getBinaryFilename("darwin", "arm64"), "yxer-darwin-arm64");
});

test("resolveBinaryPath returns null when file is missing", () => {
  const baseDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-bin-missing-"));
  assert.equal(resolveBinaryPath(baseDir, "darwin", "arm64"), null);
});

test("ensureExecutable adds execute bits on non-windows platforms", { skip: process.platform === "win32" }, () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-bin-"));
  const binaryPath = path.join(tempDir, "yxer-darwin-arm64");
  fs.writeFileSync(binaryPath, "echo test\n", { mode: 0o644 });

  ensureExecutable(binaryPath);

  const stat = fs.statSync(binaryPath);
  assert.equal(stat.mode & 0o111, 0o111);
});
