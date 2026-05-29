---
name: yixiaoer
version: 3.0.0
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
4. 发布前必须先运行 `yxer accounts` 确认目标账号 `status=1`。
5. 图片和视频必须先通过 `yxer upload` 获取 key，禁止直接在 payload 填外部 URL。
6. 生成 payload 后必须先执行 `yxer validate`，再执行 `yxer publish`。
7. 真正执行统一走 `yxer` CLI，不要假设存在旧的 Node 脚本入口。

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

- 通用发布规则：`../../references/workflows/common-rules.md`
- 图文发布：`../../references/workflows/publish-imageText.md`
- 视频发布：`../../references/workflows/publish-video.md`
- 文章发布：`../../references/workflows/publish-article.md`

## 参考入口

- CLI 命令参考：`../../references/cli/command-reference.md`
- 技能安装与同步：`../../references/cli/skill-install.md`
- 技能目录约定：`../../references/skill-architecture.md`
- 平台索引：`../../references/platforms/index.md`
