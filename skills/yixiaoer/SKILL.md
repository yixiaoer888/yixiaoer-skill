---
name: yixiaoer
version: 3.1.0
description: "蚁小二多平台内容分发 CLI 技能。用于图文、视频、文章的账号查询、资源上传、发布前字段准备、payload 校验、云发布/本机发布、草稿保存、素材登记、发布记录排查，以及分类、位置、音乐、合集、商品、话题等动态字段查询。用户提到蚁小二、yxer、多平台分发、内容发布、云发布、本机发布、客户端发布、保存草稿、素材库、发布失败排查时使用。本 skill 只负责指挥 Agent 统一调用 yxer CLI，不假设存在旧 Node 脚本入口。"
metadata:
  requires:
    bins: ["yxer"]
  cliHelp: "yxer --help; yxer doctor; yxer accounts list --help; yxer publish --help; yxer validate --help"
---

# yxer

> `yxer` 是唯一执行入口。发布、查询、上传、校验、草稿、素材、排查都必须走 CLI，不要假设存在其他脚本或隐式 API。

## CRITICAL

- 开始前先读取 [`./references/yixiaoer-shared.md`](./references/yixiaoer-shared.md)。
- 先按任务类型路由，再读取对应节点文件；未进入正确 workflow 前，禁止直接组装 payload 或执行正式发布。
- 发布、草稿、素材、排查都只允许通过 `yxer` CLI 执行；不要假设旧 Node 入口、隐式 API 或手工 DTO 模板。
- 真正发布前，固定顺序是 `doctor -> accounts list -> prepare -> schema get -> validate -> publish --dry-run -> publish`。

## 快速路由

- 发布任务：
  1. [`./references/workflows/common-rules.md`](./references/workflows/common-rules.md)
  2. [`./references/workflows/account-selection.md`](./references/workflows/account-selection.md)
  3. [`./references/workflows/local-vs-cloud.md`](./references/workflows/local-vs-cloud.md)
  4. [`./references/workflows/payload-sourcing.md`](./references/workflows/payload-sourcing.md)
  5. 按类型读取：
     - 图文：[`./references/workflows/publish-imageText.md`](./references/workflows/publish-imageText.md)
     - 视频：[`./references/workflows/publish-video.md`](./references/workflows/publish-video.md)
     - 文章：[`./references/workflows/publish-article.md`](./references/workflows/publish-article.md)
- 草稿任务：
  1. [`./references/workflows/draft-workflow.md`](./references/workflows/draft-workflow.md)
  2. 如草稿 payload 未成型，再补读对应发布 workflow
- 素材库任务：
  1. [`./references/workflows/material-workflow.md`](./references/workflows/material-workflow.md)
- 发布失败排查 / 历史记录：
  1. [`./references/workflows/publish-troubleshooting.md`](./references/workflows/publish-troubleshooting.md)
- 只生成或修 payload：
  1. [`./references/workflows/payload-sourcing.md`](./references/workflows/payload-sourcing.md)
- 只查账号：
  1. [`./references/workflows/account-selection.md`](./references/workflows/account-selection.md)
- 只判断云/本机通道：
  1. [`./references/workflows/local-vs-cloud.md`](./references/workflows/local-vs-cloud.md)

## 通用参考

- [`./references/cli/command-reference.md`](./references/cli/command-reference.md)
- [`./references/cli/skill-install.md`](./references/cli/skill-install.md)
- [`./references/platforms/index.md`](./references/platforms/index.md)

## 快速决策

- 用户要“发布”：进入发布路由，不要直接从命令列表开跑。
- 用户要“保存草稿”：先读 `draft-workflow.md`，先区分蚁小二草稿和平台草稿箱。
- 用户要“素材库”：先读 `material-workflow.md`，优先判断是 `material add` 还是 `material create`。
- 用户要“查账号”或“选账号”：先读 `account-selection.md`。
- 用户要“判断走云还是本机”：先读 `local-vs-cloud.md`。
- 用户要“修 payload”或“解释字段放哪”：先读 `payload-sourcing.md`。
- 用户要“排查发布失败”或“看历史记录”：先读 `publish-troubleshooting.md`。
- 用户要发图、视频、封面：先 `yxer upload`，禁止在 payload 里直接写外部 URL。
- 用户要填分类、位置、音乐、合集、商品、话题：先查询，再回填 CLI 返回的合法结构。
- 用户说“看发布记录 / 排查失败”：用 `yxer records list`，必要时再回到 `prepare` / `schema get` / workflow 核对。

## 阻塞性规则

1. 首次使用、环境不明、或刚经历失败排查后，先运行 `yxer doctor`。
2. 发布前必须先执行 `yxer accounts list`，并确认目标账号 `status=1`。
3. 发布前必须先执行 `yxer prepare <platform> <type>` 和 `yxer schema get <platform> <type>`。
4. payload 只能依据 workflow、schema、prepare 返回值、平台文档和 CLI 实际输出组装；禁止猜字段、猜枚举、猜默认值。
5. 共享字段与账号字段必须保持标准层级，禁止把 `publishArgs` 下字段乱塞到 `contentPublishForm`，也不要反向移动。
6. 图片、视频、封面等资源必须先上传；只能复用 `yxer upload` 返回的真实 `key` 和元数据。
7. `category`、`location`、`music`、`collection`、`challenge`、`goods` 等动态字段必须先查询，不能手写 `raw` 对象。
8. 用户未提供的业务内容，如标题、正文、描述、定时发布时间，除非 workflow 明确允许并要求先确认，否则不要私自生成。
9. 先 `yxer validate`，再 `yxer publish --dry-run`，最后正式 `yxer publish`。
10. `validate`、`publish --dry-run`、正式 `publish` 必须使用同一套发布通道参数。
11. 统一走 `yxer` CLI，不要假设存在旧 Node 入口。

## 标准发布骨架

```bash
yxer doctor
yxer accounts list <platform> --status 1 --json
yxer prepare <platform> <type>
yxer schema get <platform> <type>
# 如需资源，先上传
yxer upload --file <path>
# 如需动态字段，先查询
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
# 只基于以上结果填写 payload
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
yxer publish <type> <platform> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>] --dry-run
yxer publish <type> <platform> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
```

## Payload 纪律

- 不允许从空白文件手工创建 payload；必须先基于 CLI 模板或标准结构生成骨架，再按 `prepare` / `schema get` 结果填值。
- 顶层使用标准 envelope：`action` / `publishType` / `platforms` / `publishChannel` / `publishArgs`。
- `accountForms` 必须放在 `publishArgs.accountForms[]`。
- 平台业务字段默认放在 `publishArgs.accountForms[].contentPublishForm`。
- 上传产物只能复用 `yxer upload` 返回的真实 `key`、`size`、`width`、`height`、`duration`、`format`。
- 文档没有写、schema 没有定义、CLI 没有返回的字段，一律不要出现在 payload 里。
- 校验失败时，先回到 schema / prepare / 平台文档核对结构，再修 payload；不要靠试错继续瞎填。

## 常用命令

```bash
yxer doctor
yxer config get
yxer config init --api-key <apiKey> [--bind-app --account-id <id> | --account-name <name>]
yxer config set-local-client-id <clientId>
yxer accounts list [platform] [--name 关键词] [--status 1] [--json]
yxer upload --file <file_path> [--bucket cloud-publish|material-library] [--dry-run]
yxer upload --url <resource_url> [--bucket cloud-publish|material-library] [--dry-run]
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
yxer records list [--platform P] [--limit N] [--status S] [--json]
yxer skill show
yxer skill sync [--global]
yxer linked-app status
yxer linked-app connect --account-id <id> --account-name <name>
yxer linked-app disconnect
yxer linked-app toggle
yxer prepare <platform> <type>
yxer schema get <platform> <type>
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
yxer publish <type> <platform> <payload.json> [clientId] [--dry-run]
yxer draft save <payload.json> [--dry-run]
yxer material create <payload.json> [--dry-run]
yxer material add --file <文件路径或URL> [--thumb <缩略图路径或URL>] [--type image|video|file] [--dry-run]
```

## 失败处理

- `doctor` 失败：先修环境，再继续。
- `accounts list` 查不到可用账号：先确认平台、账号名和状态，不要继续填 payload。
- `validate` 失败：回到 `prepare` / `schema get` / workflow 查结构，不要靠猜测修字段。
- 云发布遇到代理相关错误：可建议切到本机发布。
- 本机发布提示客户端不在线：让用户启动并登录蚁小二客户端，或改回云发布。

## 安装与同步

```bash
npx skills add "<repo>\\skills\\yixiaoer" -y
npx skills add "<repo>\\skills\\yixiaoer" -g -y
yxer skill show
yxer skill sync
yxer skill sync --global
```

当 `yxer --version` 升级、`SKILL.md` 更新、或 workflow / CLI reference 更新后，应重新同步 skill。
