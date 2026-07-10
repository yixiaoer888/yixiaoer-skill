const test = require("node:test");
const assert = require("node:assert/strict");
const crypto = require("node:crypto");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

const {
  buildDownloadUrl,
  download,
  getExpectedChecksum,
  getTarget,
  parseChecksums,
  verifyChecksum
} = require("../../npm/bin/install");
const { ensureExecutable, shouldEnsureExecutable } = require("../../npm/bin/ensure-executable");
const { getBinaryFilename, resolveBinaryPath } = require("../../npm/bin/resolve-binary");

test("shouldEnsureExecutable skips windows only", () => {
  assert.equal(shouldEnsureExecutable("win32"), false);
  assert.equal(shouldEnsureExecutable("darwin"), true);
  assert.equal(shouldEnsureExecutable("linux"), true);
});

test("getBinaryFilename resolves mac arm64 binary", () => {
  assert.equal(getBinaryFilename("darwin", "arm64"), "yxer");
});

test("resolveBinaryPath returns null when file is missing", () => {
  const baseDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-bin-missing-"));
  assert.equal(resolveBinaryPath(baseDir, "darwin", "arm64"), null);
});

test("getTarget returns archive metadata for windows x64", () => {
  assert.deepEqual(getTarget("win32", "x64"), {
    platform: "win32",
    arch: "x64",
    mappedPlatform: "windows",
    mappedArch: "amd64",
    archiveName: "yxer-cli-0.0.0-windows-amd64.zip",
    binaryName: "yxer.exe"
  });
});

test("buildDownloadUrl trims trailing slash", () => {
  assert.equal(
    buildDownloadUrl("https://example.com/releases/", "yxer-cli-3.2.2-linux-amd64.tar.gz"),
    "https://example.com/releases/yxer-cli-3.2.2-linux-amd64.tar.gz"
  );
});

test("parseChecksums reads sha256 entries", () => {
  const parsed = parseChecksums("abc123  yxer-cli-3.2.2-linux-amd64.tar.gz\nfff999  yxer-cli-3.2.2-windows-amd64.zip\n");
  assert.equal(parsed.get("yxer-cli-3.2.2-linux-amd64.tar.gz"), "abc123");
  assert.equal(parsed.get("yxer-cli-3.2.2-windows-amd64.zip"), "fff999");
});

test("download falls back to PowerShell on Windows when curl fails", () => {
  const calls = [];
  const runner = (command, args, options) => {
    calls.push({ command, args, options });
    if (command === "curl") {
      const error = new Error("curl: (56) Recv failure: Connection was reset");
      throw error;
    }
  };

  assert.doesNotThrow(() => download("https://github.com/example/release.zip", "C:\\tmp\\release.zip", "win32", runner));
  assert.equal(calls.length, 2);
  assert.equal(calls[0].command, "curl");
  assert.equal(calls[1].command, "powershell.exe");
  assert.equal(calls[1].options.env.YXER_URL, "https://github.com/example/release.zip");
  assert.equal(calls[1].options.env.YXER_DEST, "C:\\tmp\\release.zip");
});

test("download preserves curl failure on non-windows platforms", () => {
  const runner = () => {
    throw new Error("curl failed");
  };

  assert.throws(() => download("https://github.com/example/release.tar.gz", "/tmp/release.tar.gz", "linux", runner), {
    message: "curl failed"
  });
});

test("getExpectedChecksum returns matching checksum from file", () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-checksums-"));
  const checksumsPath = path.join(tempDir, "checksums.txt");
  fs.writeFileSync(checksumsPath, "deadbeef  yxer-cli-3.2.2-linux-amd64.tar.gz\n");

  assert.equal(getExpectedChecksum("yxer-cli-3.2.2-linux-amd64.tar.gz", checksumsPath), "deadbeef");
  assert.equal(getExpectedChecksum("missing.tar.gz", checksumsPath), null);
});

test("verifyChecksum accepts matching sha256", () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-hash-"));
  const archivePath = path.join(tempDir, "archive.tgz");
  const content = "checksum-test";
  fs.writeFileSync(archivePath, content);
  const expectedHash = crypto.createHash("sha256").update(content).digest("hex");

  assert.doesNotThrow(() => verifyChecksum(archivePath, expectedHash));
});

test("ensureExecutable adds execute bits on non-windows platforms", { skip: process.platform === "win32" }, () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "yxer-bin-"));
  const binaryPath = path.join(tempDir, "yxer");
  fs.writeFileSync(binaryPath, "echo test\n", { mode: 0o644 });

  ensureExecutable(binaryPath);

  const stat = fs.statSync(binaryPath);
  assert.equal(stat.mode & 0o111, 0o111);
});
