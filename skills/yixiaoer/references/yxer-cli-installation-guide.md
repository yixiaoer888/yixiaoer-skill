# yxer CLI 安装指南

本文档按照飞书 `feishu-cli-installation-guide.md` 的目录思路整理，提供 `yxer` CLI 与 `yixiaoer` skill 的安装、配置与验证步骤。

## 环境要求

安装 `yxer` CLI 前，请先确认本机具备以下环境：

- Go `1.25.0` 或更高版本
- Node.js
- `npx`

说明：

- Go 用于本地构建 `yxer` 二进制
- Node.js 与 `npx` 用于后续安装或同步 `yixiaoer` skill

## 第一步：安装 yxer CLI

当前仓库没有“自动下载安装最新版二进制”的命令，推荐直接从源码构建。

1. 获取仓库代码并进入项目目录。

```bash
git clone -b feature/yxrcli https://github.com/yixiaoer888/yixiaoer-skill.git
cd yixiaoer-skill
```

2. 在仓库根目录构建 CLI。

```bash
go build -o bin/yxer.exe .
```

3. 构建完成后，确认可执行文件已生成。

```text
bin\yxer.exe
```

如需直接使用 `yxer` 命令而不是 `bin\yxer.exe`，可将 `bin` 目录加入系统 `PATH`。

## 第二步：初始化 CLI 配置

安装完成后，先初始化本地配置。

```bash
yxer config init --api-key <apiKey>
```

如果当前宿主需要同时把蚁小二作为链接应用启用，可在初始化时一并绑定：

```bash
yxer config init --api-key <apiKey> --bind-app --account-name <账号名>
```

如需预设本机发布默认 `clientId`，可执行：

```bash
yxer config set-local-client-id <clientId>
```

初始化完成后，可查看当前配置：

```bash
yxer config get
```

## 第三步：安装 yixiaoer Skill

`yxer` CLI 安装完成后，必须继续安装 `yixiaoer` skill。CLI 负责真正执行，skill 负责约束 AI agent 的工作流与命令选择，两者缺一不可。

在仓库根目录执行：

```bash
npx skills add ".\\skills\\yixiaoer" -y
```

如需全局安装：

```bash
npx skills add ".\\skills\\yixiaoer" -g -y
```

也可以直接通过 CLI 同步：

```bash
yxer skill sync
yxer skill sync --global
```

## 第四步：验证安装是否成功

依次执行以下命令检查 CLI 和 skill 的安装结果：

```bash
yxer --version
yxer doctor
yxer skill show
```

验证通过时，通常应满足以下条件：

- `yxer --version` 能正常输出版本号
- `yxer doctor` 不再提示缺少 `apiKey`
- `yxer skill show` 能正确显示 `skills/yixiaoer` 的目录位置
- `yxer doctor` 的返回结果中，skill 状态不应再提示未安装或未同步

## 下一步

安装成功后，建议按以下顺序开始首次使用：

```bash
yxer doctor
yxer accounts
yxer prepare <platform> <type>
yxer schema get <platform> <type>
```

如果后续需要发布内容，再继续执行：

```bash
yxer upload <file_path_or_url>
yxer validate <platform> <type> <payload.json>
yxer publish <type> <platform> <payload.json>
```
