# CLI 与 Skill 安装卸载说明

本文档覆盖 `yxer` CLI 的下载/安装、升级、卸载，以及 `yixiaoer` skill 的安装和卸载。

## 1. 环境要求

- Go `1.25.0` 或更高
- Node.js
- `npx`

## 2. CLI 下载与安装

当前仓库没有内置“自动下载最新版 CLI”的能力，推荐以下两种获取方式。

### 方式 A：从源码构建

```bash
git clone <repo>
cd yixiaoer-skill
go build -o bin/yxer.exe .
```

构建完成后，可执行文件位于：

```text
bin\yxer.exe
```

### 方式 B：使用发布人员提供的制品包

如果团队已经产出上线包，直接获取其中的 `yxer.exe` 即可，无需本地重新构建。

建议将 `yxer.exe` 放到固定目录，例如：

```text
C:\tools\yxer\yxer.exe
```

如需全局直接执行，可把所在目录加入 `PATH`。

## 3. CLI 安装后初始化

```bash
yxer --version
yxer config init --api-key <apiKey>
yxer doctor
```

如需要同时把蚁小二作为链接应用启用：

```bash
yxer config init --api-key <apiKey>
```

## 4. Skill 安装

先确认 skill 包位置：

```bash
yxer skill show
```

如果 CLI 通过 npm 成品包装好且技能文件随包分发，优先直接同步：

```bash
yxer skill sync
yxer skill sync --global
```

本地安装：

```bash
npx skills add "<repo>\skills\yixiaoer" -y
```

全局安装：

```bash
npx skills add "<repo>\skills\yixiaoer" -g -y
```

也可以直接通过 CLI 同步：

```bash
yxer skill sync
yxer skill sync --global
```

## 5. 升级方式

### 升级 CLI

当前仓库推荐重新构建或重新分发新版二进制：

```bash
git pull
go build -o bin/yxer.exe .
```

### 升级 Skill

CLI 升级后，必须同步 skill：

```bash
yxer skill sync
```

检查状态：

```bash
yxer update --check
yxer doctor
```

## 6. CLI 卸载

### 标准卸载

如果 `yxer.exe` 是单独分发的二进制，删除可执行文件并移除 `PATH` 即可。

例如：

- 删除 `C:\tools\yxer\yxer.exe`
- 从系统环境变量中移除 `C:\tools\yxer`

### 配置清理

如需要彻底清理本地配置，可删除以下目录中的配置文件：

```text
%USERPROFILE%\.yxer\config.json
%USERPROFILE%\.yxer\skills.stamp
```

说明：

- `config.json` 保存 `apiKey`、本机发布 `clientId` 和 linked app 状态
- `skills.stamp` 用于记录本地 skill 同步版本

### 源码构建场景

如果 CLI 是在仓库内构建的，可删除构建产物：

```text
<repo>\bin\yxer.exe
```

## 7. Skill 卸载

### 推荐做法

优先使用你当前 `skills` 工具所提供的卸载命令卸载 `yixiaoer` skill。

由于不同宿主环境里的 `skills` 工具版本可能不完全一致，建议先查看帮助：

```bash
npx skills --help
```

然后按工具实际支持的删除命令执行。

### 手动卸载

如果当前环境没有统一的 skill 删除命令，可以直接删除已安装的 `yixiaoer` skill 目录，并清理版本戳：

```text
删除已安装的 yixiaoer 技能目录
删除 %USERPROFILE%\.yxer\skills.stamp
```

如果你不确定技能安装到了哪里，先执行：

```bash
yxer skill show
```

再根据宿主工具实际安装位置清理。

## 8. 卸载后的验证

CLI 卸载后，以下命令应不可再执行：

```bash
yxer --version
```

Skill 卸载后，建议确认：

- AI 宿主中不再能读取 `yixiaoer` skill
- `%USERPROFILE%\.yxer\skills.stamp` 已删除或不再指向旧版本

## 9. 常见建议

- 升级时优先“先替换 CLI，再同步 skill”
- 卸载时优先“先卸载 skill，再删除 CLI”
- 如果只是停用能力而不是彻底卸载，保留 `config.json` 会更方便后续恢复
