# @yixiaoermail/cli

Packaged `yxer` CLI for global npm installation.

## Install

```bash
npm install -g @yixiaoermail/cli
```

The npm package installs a lightweight launcher. During install or first run it downloads the matching platform binary from `https://oss-v2.yixiaoer.cn/yxer/releases/v<version>/`.

## Verify

```bash
yxer --version
```

## Install skill

The npm package includes the `yixiaoer` skill bundle. Sync it into your agent host with:

```bash
yxer skill sync
```

Use `yxer skill sync --global` if your host expects a global skill install.

## Release Packaging

This package expects assets in the versioned TOS/CDN directory named like:

```text
yxer-cli-<version>-windows-amd64.zip
yxer-cli-<version>-windows-arm64.zip
yxer-cli-<version>-darwin-amd64.tar.gz
yxer-cli-<version>-darwin-arm64.tar.gz
yxer-cli-<version>-linux-amd64.tar.gz
yxer-cli-<version>-linux-arm64.tar.gz
checksums.txt
```

To override the versioned download directory, set `YXER_DOWNLOAD_BASE_URL` before installation.
