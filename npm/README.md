# @yixiaoer/cli

Packaged `yxer` CLI for global npm installation.

## Install

```bash
npm install -g @yixiaoermail/cli
```

The npm package installs a lightweight launcher. During install or first run it downloads the matching platform binary from the release assets.

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

This package expects release assets named like:

```text
yxer-cli-3.2.2-windows-amd64.zip
yxer-cli-3.2.2-windows-arm64.zip
yxer-cli-3.2.2-darwin-amd64.tar.gz
yxer-cli-3.2.2-darwin-arm64.tar.gz
yxer-cli-3.2.2-linux-amd64.tar.gz
yxer-cli-3.2.2-linux-arm64.tar.gz
checksums.txt
```

If you mirror the download host, set `YXER_DOWNLOAD_BASE_URL` before installation.
