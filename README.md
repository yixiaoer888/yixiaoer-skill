# yixiaoer-cli

蚁小二 `yxer` CLI 与 AI Skill 配套仓库。

**🚀 新功能（2026-06-02）：**
- ✨ **智能字段分组** - `schema fields` 现在按必填/可选/复杂字段分组展示
- 🔍 **智能错误分析** - `validate` 失败时自动分析原因并给出修复建议  
- 📚 **5分钟快速开始** - 查看 [`skills/yixiaoer/QUICKSTART.md`](skills/yixiaoer/QUICKSTART.md)
- 💡 **自动查询提示** - 复杂字段自动提示对应的查询命令

详见 [CHANGELOG.md](CHANGELOG.md)

---

仓库入口按当前项目的使用对象分层组织：

- `README.md` 面向人类用户和维护者，负责安装、快速开始和常用命令。
- `skills/yixiaoer/SKILL.md` 面向 AI agent，负责共享规则、能力索引和命令探索。
- `skills/yixiaoer/references/` 是引用主源，放任务分域、工作流、平台文档和命令参考。
- 根目录 `references/` 仅作为打包时的输出目录，内容来源于 `skills/yixiaoer/references/`。

命令树设计标准与演进建议见 [skills/yixiaoer/references/cli/command-design.md](skills/yixiaoer/references/cli/command-design.md)。

运行时统一通过 `yxer` 执行，不再假设存在旧 Node 脚本入口。

## 安装与快速开始

### 环境要求

- Go `1.25.0` 或更高
- Node.js 和 `npx`

### 快速开始（npm 安装，推荐）

如果你已经把 `yxer` 发布到了 npm，推荐优先使用全局安装：

```powershell
npm install -g @yixiaoermail/cli@latest
yxer --version
yxer update
```

npm 包现在采用轻量安装器模式：

- npm 包本身只包含启动器、skill 源文档、schema 和 references 打包输出
- 安装阶段会按当前系统下载匹配的 `yxer` 二进制归档
- 如果 `postinstall` 被跳过，首次运行 `yxer` 时也会自动补装二进制

如需使用私有镜像或自建发布源，可在安装前设置：

```powershell
$env:YXER_DOWNLOAD_BASE_URL = "https://your-release-host/yxer/v3.2.2"
npm install -g @yixiaoermail/cli@latest
```

如需全局同步 skill：

```powershell
yxer skill sync --global
```

升级后建议再执行一次：

```powershell
yxer doctor
```

本次版本修改内容请查看 [CHANGELOG.md](CHANGELOG.md)。

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

## 升级

如果 CLI 是通过 npm 安装的，推荐按下面顺序升级：

```powershell
yxer update
```

如需全局 skill：

```powershell
yxer update --global
```

只检查当前状态而不执行同步，可运行：

```bash
yxer update --check
```

执行“检查状态 + 同步 skill + 查看升级指引”：

```bash
yxer update
```

本次版本更新内容统一记录在 [CHANGELOG.md](CHANGELOG.md)。

## AI Skill 安装

本项目采用“CLI 先安装，Skill 再安装”的方式。

如果要接入 SkillHub / skills 市场，建议使用独立的 bootstrap skill：

- SkillHub 安装入口：`skills/yixiaoer-bootstrap`
- 正式业务 skill：`skills/yixiaoer`
- 推荐职责划分：
  - bootstrap skill 只负责安装 `@yixiaoermail/cli`、执行 `yxer skill sync --global`、引导 `config init`
  - 正式 skill 负责发布、查询、校验、草稿、素材、排障

如果 CLI 是通过 npm 成品包安装的，推荐优先使用：

```bash
yxer skill sync
yxer skill sync --global
```

npm 包会内置 `skills/yixiaoer`，`skill sync` 会直接使用本地随包分发的 skill 源文件，不依赖 GitHub 仓库地址。

### 生成 npm 成品包

如需产出可通过 `npm install -g` 安装的 CLI 成品包，可在仓库根目录执行：

```powershell
.\scripts\build-npm-package.ps1
```

默认会自动读取仓库内部版本，并校验以下版本源保持一致：

- `internal/domain/response.go`
- `skills/yixiaoer/SKILL.md`

如需显式传版本号，也必须与内部版本一致：

```powershell
.\scripts\build-npm-package.ps1 -Version 3.1.1
```

该脚本会：

- 先运行 `go test ./...`
- 交叉编译 `windows/darwin/linux` 的 `amd64/arm64` 二进制
- 生成 release 压缩包与 `checksums.txt` 到 `out\release\`
- 生成轻量 npm 发布包到 `out\npm\`

生成完成后，可用下列命令验证 tarball：

```powershell
npm install -g .\out\npm\<generated-tarball>.tgz
yxer --version
yxer skill sync
```

发布到 GitHub Release 或其他下载源时，需同时上传：

- `yxer-cli-<version>-windows-amd64.zip`
- `yxer-cli-<version>-windows-arm64.zip`
- `yxer-cli-<version>-darwin-amd64.tar.gz`
- `yxer-cli-<version>-darwin-arm64.tar.gz`
- `yxer-cli-<version>-linux-amd64.tar.gz`
- `yxer-cli-<version>-linux-arm64.tar.gz`
- `checksums.txt`

如果仓库已配置 GitHub Actions 发版流，也可以直接通过打 tag 触发自动上传：

```powershell
git tag v3.2.2
git push origin v3.2.2
```

注意：仅本地创建 tag 不会触发远端发版，必须把 tag push 到 GitHub。

按当前仓库的自动发版逻辑，push `v*` tag 后会依次完成：

- 构建并上传 GitHub Release 资产
- 上传 npm tarball 到 Release
- 使用 Release 中的 tarball 自动发布到 npmjs

如需启用 npm 自动发布，需要在 GitHub 仓库 Secrets 中配置 `NPM_TOKEN`。

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
- `skills/yixiaoer/references/workflows/` 或 `skills/yixiaoer/references/cli/` 更新后

## 面向 AI Agent

如果你是把本仓库接给 AI agent、Codex 或其他命令式助手使用，不要只按一个线性发布流程读取；应先读技能入口，再进入对应 domain 文档，再下钻到 workflow 或平台文档。

### 统一入口

1. `skills/yixiaoer/SKILL.md`
2. 再按任务类型继续读取下列 domain 节点

### 任务路由

- 发布任务：
  - `skills/yixiaoer/references/domains/publish.md`
  - 继续进入 `skills/yixiaoer/references/workflows/common-rules.md`、`skills/yixiaoer/references/workflows/account-selection.md`、`skills/yixiaoer/references/workflows/local-vs-cloud.md`、`skills/yixiaoer/references/workflows/payload-sourcing.md`
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
yxer accounts update <account_id> [--proxy-id ID] [--kuaidaili-area CODE] [--remark 文本] [--group ID] [--dry-run]
yxer upload --file <file_path> [--bucket cloud-publish|material-library] [--dry-run]
yxer upload --url <resource_url> [--bucket cloud-publish|material-library] [--dry-run]
yxer material list [--name <file_name>] [--type image|video|file] [--page 1] [--size 100]
yxer material move <material_id> --group-id <group_id> [--dry-run]
yxer material move-by-name <file_name> --group-id <group_id> [--dry-run]
yxer material groups [--page 1] [--size 50]
```

### 发布和校验

```bash
yxer prepare <platform> <type>
yxer schema fields <platform> <type>
yxer schema get <platform> <type>
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>] [--content-file <article.md>]
yxer publish <type> <platform> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>] [--content-file <article.md>] [--dry-run]
yxer publish delete preview <task_set_id>
yxer publish delete from-record <task_set_id> --index N [--dry-run]
yxer publish delete <task_id> [--dry-run]
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
yxer query account-increments --start-date YYYY-MM-DD --end-date YYYY-MM-DD
```

`account-increments` 支持 `--group-id` 按账号分组、`--platform` 按平台和 `--account-name` 按账号名称筛选，同时返回增量数据、`dmMessageStats` 私信收发统计和 `managedAccounts` 管理账号详情。私信总计和逐日趋势均按增量账号列表统计；`data.accounts` 还带每个增量账号的 `dmInCount`、`dmOutCount` 和状态。

说明：

- 查询类能力统一使用 `yxer query ...` 入口，不再保留旧的一层查询命令。

## 使用说明

### 默认发布通道

- 用户未明确指定时，默认走云发布。
- 用户明确要求“本机发布 / 本地发布 / 客户端发布”时，必须走本机发布。
- 本机发布必须提供 `clientId`，可通过 `yxer config set-local-client-id <clientId>` 预设。
- `validate`、`publish --dry-run`、`publish` 会共用同一套发布通道解析规则；若已预设默认 `clientId`，可在本机发布时只传 `--publish-channel local`。
- `publish --dry-run` 不创建发布任务；它返回最终请求预览，并在有 API key 的云发布场景执行账号/代理 preflight。local dry-run 不检测客户端在线状态。

### 推荐任务分流

- 用户要“发布”：先走 `skills/yixiaoer/references/domains/publish.md`，再进具体 workflow
- 用户要“保存草稿”或“素材库”：先走 `skills/yixiaoer/references/domains/draft-and-material.md`
- 用户要“查账号”或“看环境”：先走 `skills/yixiaoer/references/domains/accounts-and-env.md`
- 用户要“修 payload”：先走 `skills/yixiaoer/references/domains/publish.md`，再下钻 `payload-sourcing`
- 用户要“排查发布失败”或“看历史记录”：先走 `skills/yixiaoer/references/domains/troubleshooting.md`

### 推荐发布顺序

已有完整标准 `payload.json` 时，发布类任务按这个快速路径执行：

1. `yxer validate`
2. `yxer publish --dry-run`
3. 用户授权后 `yxer publish`

缺账号、字段、资源或动态对象时，先按这个组装路径补齐：

1. `yxer doctor`
2. `yxer accounts list`
3. `yxer prepare`
4. `yxer schema fields`（必要时再补 `yxer schema get`）
5. `yxer upload`
6. 查询分类、位置、音乐等复杂对象
7. 填写 `payload.json`
8. `yxer validate`
9. `yxer publish --dry-run`
10. 用户授权后 `yxer publish`

上面这 10 步只适用于“从零组装发布 payload”；草稿、素材库、排查等任务不要强行套这条主流程。

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
  "version": "3.1.1",
  "data": {
    "configPath": "C:\\Users\\<user>\\AppData\\Roaming\\yxer\\config.json",
    "apiUrl": "https://www.yixiaoer.cn/api",
    "apiKeyPresent": true
  },
  "_notice": {
    "skills": {
      "current": "3.1.1",
      "target": "3.1.1"
    }
  }
}
```

### 失败输出

```json
{
  "ok": false,
  "version": "3.1.1",
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

- stdout 始终输出 JSON 数据
- stderr 只输出诊断、警告、提示和结构化错误
- 成功输出使用 `ok/action/version/data`
- 失败输出使用 `ok/version/error`

## 常见问题

### 1. `doctor` 提示 skill 版本不同步

先执行：

```bash
yxer update
```

如果用户明确只想单独同步 skill，也可以：

```bash
yxer skill sync
yxer skill sync --global
```

### 2. 云发布失败，提示账号代理不存在

- 先检查账号代理配置
- 如需快速绕过云端代理问题，可改用本机发布
- 本机发布前先确认已经配置 `clientId`
- CLI 默认只返回本机发布 `nextCommand`，不会自动切换通道；`--auto-fallback-local` 只适合用户明确授权自动回退的场景。

### 3. 本机发布失败，提示客户端不在线

- 启动并登录蚁小二客户端
- 重新执行本机发布
- 如果当前不方便保持客户端在线，可改回云发布

### 4. 用户只说“保存草稿”怎么处理

不要默认。需要先区分：

- 蚁小二草稿
- 平台草稿箱

对应入口：

- 蚁小二草稿：`skills/yixiaoer/references/workflows/draft-workflow.md`
- 发布失败排查：`skills/yixiaoer/references/workflows/publish-troubleshooting.md`
- 通道判断：`skills/yixiaoer/references/workflows/local-vs-cloud.md`
- payload 修订：`skills/yixiaoer/references/workflows/payload-sourcing.md`

## 目录结构

```text
README.md
cmd/
internal/
schemas/
skills/
  yixiaoer/
    SKILL.md
      references/  # 打包时从 skills/yixiaoer/references/ 复制
      domains/
      platforms/
references/  # 打包输出，不作为主维护源
tests/
scripts/
```

## 文档索引

- 技能入口：`skills/yixiaoer/SKILL.md`
- 任务分域：`skills/yixiaoer/references/domains/`
- 命令参考：`skills/yixiaoer/references/cli/command-reference.md`
- 安装、升级与同步：`skills/yixiaoer/references/domains/install-and-sync.md`
- 上线流程：`skills/yixiaoer/references/go-live-process.md`
- 关键词文档：`skills/yixiaoer/references/keyword-reference.md`
- 使用流程文档：`skills/yixiaoer/references/usage-workflow.md`
- 工作流正文：`skills/yixiaoer/references/workflows/`
- 平台文档：`skills/yixiaoer/references/platforms/`
- 平台文档维护规范：`skills/yixiaoer/references/platform-doc-maintenance.md`
