# yixiaoer-cli

蚁小二 `yxer` CLI 与 AI Skill 配套仓库。

仓库入口参考飞书 CLI 的组织方式来写：

- `README.md` 面向人类用户和维护者，负责安装、快速开始和常用命令。
- `skills/yixiaoer/SKILL.md` 面向 AI agent，负责技能触发规则、工作流和命令选择。
- `references/` 放命令参考、工作流和平台差异说明。

运行时统一通过 `yxer` 执行，不再假设存在旧 Node 脚本入口。

## 安装

### 环境要求

- Go `1.25.0` 或更高
- Node.js 和 `npx`

### 从源码构建

在仓库根目录执行：

```bash
go build -o bin/yxer.exe .
```

也可以直接使用项目里的 `Makefile`：

```bash
make build
```

构建完成后，可执行文件位于：

```text
bin/yxer.exe
```

### 验证安装

```bash
bin/yxer.exe --version
bin/yxer.exe config set-api-key <apiKey>
bin/yxer.exe doctor
```

如果你已经把 `bin/` 加入了 `PATH`，也可以直接执行：

```bash
yxer --version
yxer doctor
```

## 3 分钟开始

### 1. 检查本地环境

```bash
yxer config set-api-key <apiKey>
yxer doctor
```

### 2. 查看当前配置

```bash
yxer config get
```

### 3. 查询可用账号

```bash
yxer accounts list 抖音 --json
```

### 4. 预览一次发布请求

```bash
yxer publish video 抖音 .\payload.json --dry-run
```

如果你走 flags 模式，也可以直接：

```bash
yxer publish video 抖音 \
  --account "视频账号" \
  --title "视频标题" \
  --description "视频描述" \
  --video .\clip.mp4 \
  --cover .\cover.png \
  --dry-run
```

## AI Skill 安装

本项目采用“CLI 先安装，Skill 再安装”的方式，和飞书 CLI 的 skill 使用习惯保持一致。

### 查看当前技能包位置

```bash
yxer skill show
```

### 安装 skill

```bash
npx skills add "<repo>\\skills\\yixiaoer" -y
```

如需全局安装：

```bash
npx skills add "<repo>\\skills\\yixiaoer" -g -y
```

### 使用 CLI 同步 skill

```bash
yxer skill sync
yxer skill sync --global
```

建议在以下场景重新同步：

- `yxer --version` 升级后
- `skills/yixiaoer/SKILL.md` 更新后
- `references/workflows/` 或 `references/cli/` 更新后

## 面向 AI Agent

如果你是把本仓库接给 AI agent、Codex 或其他命令式助手使用，建议按下面的入口顺序读取：

1. `skills/yixiaoer/SKILL.md`
2. `references/workflows/common-rules.md`
3. 对应类型的工作流：
   `references/workflows/publish-video.md`、
   `references/workflows/publish-imageText.md`、
   `references/workflows/publish-article.md`
4. 需要平台差异时，再查 `docs/publish/`

推荐约束：

- 优先调用 `yxer` CLI，不手写大 JSON
- 发布前先查账号，再上传资源
- 复杂对象通过查询命令获取，不手写 `raw`
- 本机发布显式传 `--publish-channel local` 和 `--client-id`

## 常用命令

### 环境和配置

```bash
yxer doctor
yxer config get
yxer config set-api-key <apiKey>
yxer config set-local-client-id <clientId>
```

### 账号和资源

```bash
yxer accounts list [platform] [--name 关键词] [--status 1] [--json]
yxer upload --file <file_path> [--bucket cloud-publish|material-library] [--dry-run]
yxer upload --url <resource_url> [--bucket cloud-publish|material-library] [--dry-run]
```

### 发布和校验

```bash
yxer validate <platform> <type> <payload.json>
yxer publish <type> <platform> <payload.json> [clientId] [--dry-run]
```

推荐的发布类型只有三种：

- `video`
- `imageText`
- `article`

### 查询类能力

```bash
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
yxer records list [--platform P] [--limit N] [--status S] [--json]
```

## 使用说明

### 默认发布通道

- 用户未明确指定时，默认走云发布。
- 用户明确要求“本机发布 / 本地发布 / 客户端发布”时，必须走本机发布。
- 本机发布必须提供 `clientId`，可通过 `yxer config set-local-client-id <clientId>` 预设。

### 推荐执行顺序

发布类任务建议始终按这个顺序执行：

1. `yxer doctor`
2. `yxer accounts list`
3. `yxer upload`
4. 查询分类、位置、音乐等复杂对象
5. `yxer validate`
6. `yxer publish --dry-run`
7. `yxer publish`

### Skill 与 CLI 的分工

- `README.md`：给人看，负责安装和上手。
- `skills/yixiaoer/SKILL.md`：给 agent 看，负责规则和工作流。
- `yxer` CLI：真正执行账号查询、资源上传、校验和发布。

## 输出示例

### 成功输出

```json
{
  "ok": true,
  "action": "doctor",
  "version": "3.0.0",
  "data": {
    "configPath": "C:\\Users\\<user>\\AppData\\Roaming\\yxer\\config.json",
    "apiUrl": "https://www.yixiaoer.cn/api",
    "apiKeyPresent": true
  },
  "_notice": {
    "skills": {
      "current": "3.0.0",
      "target": "3.0.0"
    }
  }
}
```

### 失败输出

```json
{
  "ok": false,
  "version": "3.0.0",
  "error": {
    "type": "validation_error",
    "code": "YIXIAOER_USAGE_ERR",
    "message": "clientId must not be empty",
    "hint": "请传入有效的本机发布 clientId。",
    "retryable": false
  }
}
```

说明：

- 默认输出更适合人读
- 加 `--json` 后适合脚本和 AI agent 继续处理
- 成功输出使用 `ok/action/version/data`
- 失败输出使用 `ok/version/error`

## 常见问题

### 1. `doctor` 提示 skill 版本不同步

先执行：

```bash
yxer skill sync
```

如果需要全局同步：

```bash
yxer skill sync --global
```

### 2. 云发布失败，提示账号代理不存在

- 先检查账号代理配置
- 如需快速绕过云端代理问题，可改用本机发布
- 本机发布前先确认已经配置 `clientId`

### 3. 本机发布失败，提示客户端不在线

- 启动并登录蚁小二客户端
- 重新执行本机发布
- 如果当前不方便保持客户端在线，可改回云发布

### 4. 用户只说“保存草稿”怎么处理

不要默认。需要先区分：

- 蚁小二草稿
- 平台草稿箱

详细排查和错误说明见：

- `docs/troubleshooting-guide.md`
- `docs/execution-standard.md`
- `references/workflows/`

## 目录结构

```text
README.md
cmd/
internal/
schemas/
skills/
  yixiaoer/
    SKILL.md
references/
  cli/
  legacy/
  platforms/
  workflows/
docs/
tests/
scripts/
```

## 文档索引

- 技能入口：`skills/yixiaoer/SKILL.md`
- 命令参考：`references/cli/command-reference.md`
- 技能安装与同步：`references/cli/skill-install.md`
- 工作流正文：`references/workflows/`
- 平台文档：`docs/publish/`
