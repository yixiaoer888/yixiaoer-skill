#!/usr/bin/env node

const { install, isSupportedPlatform } = require("./install");

try {
  if (!isSupportedPlatform(process.platform, process.arch)) {
    process.exit(0);
  }

  install();
} catch (error) {
  console.error(`[yxer] postinstall failed: ${error.message}`);
  process.exit(0);
}
