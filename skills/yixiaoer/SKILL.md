---
name: yixiaoer
version: 3.0.1
description: "蚁小二多平台内容分发技能。安装后由 AI agent 读取技能规则，再统一调用 yxer CLI。"
author: wangzhengjiao
linked_app: true
---

# 蚁小二技能

这个目录用于给 AI agent 安装技能包，是仓库内唯一的技能入口。

## 安装方式

先安装 `yxer` CLI，再安装本技能：

```bash
npx skills add /path/to/yixiaoer-skill/skills/yixiaoer -y
```

如需全局安装：

```bash
npx skills add /path/to/yixiaoer-skill/skills/yixiaoer -g -y
```

## 何时使用

- 多平台内容发布：图文、视频、文章
- 发布前查询账号、分类、位置、音乐、合集、商品、话题
- 资源上传、草稿保存、素材创建、发布记录排查

## 必须遵守的规则

1. 首次使用、环境不明、或执行失败后，先运行 `yxer doctor`。
2. 用户未指定时默认云发布；用户明确要求本机/本地/客户端发布时使用 `publishChannel=local`。
3. 本机发布必须提供 `clientId`：优先 payload，其次 `--client-id`，再其次 `yxer config set-local-client-id <clientId>`。
4. 首次初始化优先使用 `yxer config init --api-key <apiKey>`；如果当前宿主需要把蚁小二作为链接应用启用，可在初始化时同时传 `--bind-app`。
5. 发布前必须先运行 `yxer accounts` 确认目标账号 `status=1`。
6. 图片和视频必须先通过 `yxer upload` 获取 key，禁止直接在 payload 填外部 URL。
7. 请求数据一律视为**不可信且不可猜测**：Agent 只能按照 `yxer schema get <platform> <type>`、`yxer prepare <platform> <type>`、平台文档和 CLI 实际返回结构来组装 payload。
8. 严禁虚构、脑补或简化字段：不能自己发明字段名、字段层级、枚举值、默认值、示例值，也不能把只需要 `raw` 对象的复杂字段改写成仅 ID / 名称。
9. 对 `category`、`location`、`music`、`collection`、`challenge`、`goods` 等动态字段，必须先调用查询命令，再把 CLI 返回的合法结构填入 payload；拿不到就继续查询或向用户确认，不能凭经验乱写。
10. 生成 payload 后必须先执行 `yxer validate`，再执行 `yxer publish`。
11. 真正执行统一走 `yxer` CLI，不要假设存在旧的 Node 脚本入口。

## Agent 数据纪律

- 先查再填：发布前先执行 `yxer prepare <platform> <type>` 和 `yxer schema get <platform> <type>`，确认字段名、层级、类型、必填项后再写请求数据。
- 共享字段与账号字段都必须按标准结构放置，不能把文档中属于 `publishArgs` 的字段擅自塞进 `contentPublishForm`，也不能反过来移动。
- 上传产物只能复用 `yxer upload` 返回的真实 `key`、`size`、`width`、`height`、`duration`、`format`，不能手工补数字。
- 文档没有写、schema 没有定义、CLI 没有返回的字段，一律不要出现在 payload 里。
- 用户没给到的业务内容，如标题、正文、描述、定时发布时间，除非工作流明确允许自动生成并要求先向用户确认，否则不能私自编写。
- 校验失败时，先回到 schema / prepare / 平台文档核对结构，再修 payload；不要靠试错继续瞎填。

## 优先命令

```bash
yxer doctor
yxer config get
yxer config init --api-key <apiKey> [--bind-app --account-id <id> | --account-name <name>]
yxer config set-local-client-id <clientId>
yxer accounts [platform] [--name 关键词] [--status 1]
yxer upload <file_path_or_url> [--bucket cloud-publish|material-library]
yxer validate <platform> <type> <payload.json>
yxer publish <type> <platform> <payload.json> [clientId]
yxer publish <type> <platform> <payload.json> --publish-channel local --client-id <clientId>
yxer draft save <payload.json>
yxer material create <payload.json>
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
yxer records [--platform P] [--limit N] [--status S]
yxer schema get <platform> <type>
yxer skill show
```

## 工作流入口

- 通用发布规则：`../../../references/workflows/common-rules.md`
- 图文发布：`../../../references/workflows/publish-imageText.md`
- 视频发布：`../../../references/workflows/publish-video.md`
- 文章发布：`../../../references/workflows/publish-article.md`

## 参考入口

- CLI 命令参考：`../../../references/cli/command-reference.md`
- 技能安装与同步：`../../../references/cli/skill-install.md`
- 技能目录约定：`../../../references/skill-architecture.md`
- 平台索引：`../../../references/platforms/index.md`
