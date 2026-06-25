# yixiaoer-cli

蚁小二 `yxer` CLI 与 AI Skill 配套仓库。

**🚀 新功能（2026-06-02）：**
- ✨ **智能字段分组** - `schema fields` 现在按必填/可选/复杂字段分组展示
- 🔍 **智能错误分析** - `validate` 失败时自动分析原因并给出修复建议  
- 📚 **5分钟快速开始** - 查看 [`skills/yixiaoer/QUICKSTART.md`](skills/yixiaoer/QUICKSTART.md)
- 💡 **自动查询提示** - 复杂字段自动提示对应的查询命令

详见 [CHANGELOG.md](CHANGELOG.md)

---

仓库入口参考飞书 CLI 的组织方式来写：

- `README.md` 面向人类用户和维护者，负责安装、快速开始和常用命令。
- `skills/yixiaoer/SKILL.md` 面向 AI agent，负责共享规则、能力索引和命令探索。
- `skills/yixiaoer/references/domains/` 放任务分域入口。
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

### 设置为全局命令

安装成功后，建议把 `yxer.exe` 所在目录加入系统 `PATH`，这样后续就可以在任意目录直接执行 `yxer`。

Windows PowerShell 示例：

```powershell
$yxerBin = (Resolve-Path .\bin).Path
[Environment]::SetEnvironmentVariable(
  "Path",
  $env:Path + ";" + $yxerBin,
  "User"
)
```

执行完成后，请重新打开一个新的终端窗口，再运行 `yxer --version` 验证全局命令是否生效。

### 验证安装

```bash
bin/yxer.exe --version
bin/yxer.exe config init --api-key <apiKey>
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
yxer config init --api-key <apiKey>
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

发布前建议先获取表单字段和前置数据，再填写 `payload.json`：

```bash
yxer prepare 抖音 video
yxer schema fields 抖音 video
yxer schema get 抖音 video
yxer validate 抖音 video .\payload.json
yxer publish video 抖音 .\payload.json --dry-run
```

`payload.json` 必须使用统一标准结构，所有平台都一样：

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["抖音"],
  "publishChannel": "cloud",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "<platformAccountId>",
        "contentPublishForm": {
          "formType": "task"
        }
      }
    ]
  }
}
```

说明：

- 顶层必须包含 `publishArgs`
- `accountForms` 必须放在 `publishArgs.accountForms`
- 平台差异字段默认放在 `publishArgs.accountForms[].contentPublishForm`
- 共享资源字段可直接放在 `publishArgs` 下，与 `accountForms` 同级，例如 `video`、`images`、`cover`、`coverKey`、`content`
- 文章正文推荐写在 `publishArgs.content`；CLI 会在缺失时自动补齐到 `publishArgs.accountForms[].contentPublishForm.content`
- 不再支持顶层直接放 `accountForms`
- 不再支持把 `title`、`description`、`visibleType` 等内层字段直接写在 payload 顶层

如需本机发布，校验阶段也建议显式带上发布通道，保证 `validate`、`--dry-run` 和正式发布使用同一套模式解析：

```bash
yxer validate 抖音 video .\payload.json --publish-channel local --client-id <clientId>
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId> --dry-run
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId>
```

## AI Skill 安装

本项目采用“CLI 先安装，Skill 再安装”的方式，和飞书 CLI 的 skill 使用习惯保持一致。

如果 CLI 是通过 npm 成品包安装的，推荐优先使用：

```bash
yxer skill sync
yxer skill sync --global
```

npm 包会内置 `skills/yixiaoer`，`skill sync` 会直接使用本地随包分发的 skill 源文件，不依赖 GitHub 仓库地址。

### 生成 npm 成品包

如需产出可通过 `npm install -g` 安装的 CLI 成品包，可在仓库根目录执行：

```powershell
.\scripts\build-npm-package.ps1 -Version 0.1.0
```

该脚本会：

- 先运行 `go test ./...`
- 交叉编译 `windows/darwin/linux` 的 `amd64/arm64` 二进制
- 生成 npm 发布包到 `out\npm\`

生成完成后，可用下列命令验证 tarball：

```powershell
npm install -g .\out\npm\<generated-tarball>.tgz
yxer --version
yxer skill sync
```

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
- `skills/yixiaoer/references/domains/` 更新后
- `references/workflows/` 或 `references/cli/` 更新后

## 面向 AI Agent

如果你是把本仓库接给 AI agent、Codex 或其他命令式助手使用，不要只按一个线性发布流程读取；应先读技能入口，再进入对应 domain 文档，再下钻到 workflow 或平台文档。

### 统一入口

1. `skills/yixiaoer/SKILL.md`
2. 再按任务类型继续读取下列 domain 节点

### 任务路由

- 发布任务：
  - `skills/yixiaoer/references/domains/publish.md`
  - 继续进入 `references/workflows/common-rules.md`、`account-selection.md`、`local-vs-cloud.md`、`payload-sourcing.md`
  - 再按类型进入 `publish-video.md`、`publish-imageText.md`、`publish-article.md`
- 草稿或素材任务：
  - `skills/yixiaoer/references/domains/draft-and-material.md`
- 发布失败排查 / 历史记录：
  - `skills/yixiaoer/references/domains/troubleshooting.md`
- 只查账号、环境、skill 同步：
  - `skills/yixiaoer/references/domains/accounts-and-env.md`
- 安装、升级、分发：
  - `skills/yixiaoer/references/domains/install-and-sync.md`
- 需要平台差异时，再查：
  - `skills/yixiaoer/references/platforms/`

推荐约束：

- 优先调用 `yxer` CLI，通过 `prepare` / `schema fields` 了解字段；`schema fields` 默认返回更紧凑的扁平路径视图，只有需要完整 payload 骨架时再看 `schema get`
- `payload.json` 必须使用标准 envelope：`action/publishType/platforms/publishArgs`
- 账号选择、通道判断、payload 来源、失败排查都应走对应 workflow，不要混写在一个大 prompt 里
- 复杂对象通过查询命令获取，不手写 `raw`
- 本机发布显式传 `--publish-channel local` 和 `--client-id`

## 常用命令

### 环境和配置

```bash
yxer doctor
yxer config get
yxer config init --api-key <apiKey>
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
yxer prepare <platform> <type>
yxer schema fields <platform> <type>
yxer schema get <platform> <type>
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
yxer publish <type> <platform> <payload.json> [clientId] [--dry-run]
```

推荐的发布类型只有三种：

- `video`
- `imageText`
- `article`

### 查询类能力

```bash
yxer query categories <account_id> [--type video|article]
yxer query locations <account_id> [--query 关键词]
yxer query music <account_id> [--query 关键词]
yxer query goods <account_id> [--query 关键词]
yxer query collections <account_id> [--type video|article]
yxer query challenges <account_id> [--query 关键词] [--type video]
yxer query records [--platform P] [--limit N] [--status S] [--json]
```

说明：

- 推荐优先使用 `yxer query ...` 作为统一查询入口
- 旧的一层命令如 `yxer locations ...`、`yxer records list ...` 仍兼容，可逐步迁移

## 使用说明

### 默认发布通道

- 用户未明确指定时，默认走云发布。
- 用户明确要求“本机发布 / 本地发布 / 客户端发布”时，必须走本机发布。
- 本机发布必须提供 `clientId`，可通过 `yxer config set-local-client-id <clientId>` 预设。
- `validate`、`publish --dry-run`、`publish` 会共用同一套发布通道解析规则；若已预设默认 `clientId`，可在本机发布时只传 `--publish-channel local`。

### 推荐任务分流

- 用户要“发布”：先走 `skills/yixiaoer/references/domains/publish.md`，再进具体 workflow
- 用户要“保存草稿”或“素材库”：先走 `skills/yixiaoer/references/domains/draft-and-material.md`
- 用户要“查账号”或“看环境”：先走 `skills/yixiaoer/references/domains/accounts-and-env.md`
- 用户要“修 payload”：先走 `skills/yixiaoer/references/domains/publish.md`，再下钻 `payload-sourcing`
- 用户要“排查发布失败”或“看历史记录”：先走 `skills/yixiaoer/references/domains/troubleshooting.md`

### 推荐发布顺序

发布类任务建议始终按这个顺序执行：

1. `yxer doctor`
2. `yxer accounts list`
3. `yxer prepare`
4. `yxer schema fields`（必要时再补 `yxer schema get`）
5. `yxer upload`
6. 查询分类、位置、音乐等复杂对象
7. 填写 `payload.json`
8. `yxer validate`
9. `yxer publish --dry-run`
10. `yxer publish`

上面这 10 步只适用于“正式发布链路”；草稿、素材库、排查等任务不要强行套这条主流程。

### Skill 与 CLI 的分工

- `README.md`：给人看，负责安装和上手。
- `skills/yixiaoer/SKILL.md`：给 agent 看，负责共享规则和能力索引。
- `skills/yixiaoer/references/domains/`：给 agent 做任务分流。
- `yxer` CLI：真正执行账号查询、资源上传、校验和发布。

## 输出示例

### 成功输出

```json
{
  "ok": true,
  "action": "doctor",
  "version": "3.1.0",
  "data": {
    "configPath": "C:\\Users\\<user>\\AppData\\Roaming\\yxer\\config.json",
    "apiUrl": "https://www.yixiaoer.cn/api",
    "apiKeyPresent": true
  },
  "_notice": {
    "skills": {
      "current": "3.1.0",
      "target": "3.1.0"
    }
  }
}
```

### 失败输出

```json
{
  "ok": false,
  "version": "3.1.0",
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

对应入口：

- 蚁小二草稿：`references/workflows/draft-workflow.md`
- 发布失败排查：`references/workflows/publish-troubleshooting.md`
- 通道判断：`references/workflows/local-vs-cloud.md`
- payload 修订：`references/workflows/payload-sourcing.md`

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
      domains/
      platforms/
references/
  cli/
  legacy/
  platforms/
  workflows/
tests/
scripts/
```

## 文档索引

- 技能入口：`skills/yixiaoer/SKILL.md`
- 任务分域：`skills/yixiaoer/references/domains/`
- 命令参考：`references/cli/command-reference.md`
- 技能安装与同步：`references/cli/skill-install.md`
- 上线流程：`skills/yixiaoer/references/go-live-process.md`
- CLI/Skill 安装卸载：`skills/yixiaoer/references/cli-install-uninstall.md`
- 关键词文档：`skills/yixiaoer/references/keyword-reference.md`
- 使用流程文档：`skills/yixiaoer/references/usage-workflow.md`
- 工作流正文：`references/workflows/`
- 平台文档：`skills/yixiaoer/references/platforms/`
- 平台文档维护规范：`skills/yixiaoer/references/platform-doc-maintenance.md`
