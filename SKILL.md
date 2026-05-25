---
name: yixiaoer
version: 3.0.0
description: "蚁小二多平台内容分发技能。统一通过 yxer CLI 执行，技能说明按飞书 Skill 的 SKILL.md + references 架构组织。"
author: wangzhengjiao
---

# 蚁小二技能

用 `yxer` CLI 做唯一执行入口，不再把 `node scripts/*.ts` 作为默认调用路径。

## 何时使用

当用户要做以下事情时使用本技能：

- 多平台内容发布：图文、视频、文章
- 账号查询与账号状态确认
- 发布前资源上传
- 分类、位置、音乐、合集、商品、话题等发布前置查询
- 发布记录查询和问题排查

## 必须遵守的规则

1. 首次使用、环境不明、或执行失败后，先运行 `yxer doctor`。
2. 发布类请求必须先读取对应工作流文档，再执行命令。
3. 发布前必须运行 `yxer accounts`，确认目标账号 `status=1`。
4. 图片和视频资源必须先通过 `yxer upload` 获取 key，禁止在 payload 中直接放外部 URL。
5. `location`、`music`、`collection`、`challenge` 这类复杂对象必须通过查询命令获得，禁止手写 `raw`。
6. 生成 payload 后必须先执行 `yxer validate`，再执行 `yxer publish`。
7. 用户只说“草稿”时，必须先确认是“蚁小二草稿”还是“平台草稿”。
8. 除非 `yxer` CLI 还没覆盖，否则不要把 `node scripts/api.ts` 当作公开主流程。

## 执行入口

优先命令：

```bash
yxer doctor
yxer accounts [platform] [--name 关键词] [--status 1]
yxer upload <file_path_or_url> [--bucket cloud-publish|material-library]
yxer validate <platform> <type> <payload.json>
yxer publish <type> <platform> <payload.json> [clientId]
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
```

## 工作流入口

按场景先读对应 reference：

- 图文发布：`references/workflows/publish-image-text.md`
- 视频发布：`references/workflows/publish-video.md`
- 文章发布：`references/workflows/publish-article.md`
- 通用发布规则：`references/workflows/common-rules.md`

## 查询与能力参考

- CLI 命令与输出约定：`references/cli/command-reference.md`
- 技能目录与扩展规范：`references/skill-architecture.md`
- 平台字段文档入口：`references/platforms/index.md`
- 旧 Node 通道说明：`references/legacy/node-compat.md`
- 迁移中能力索引：`references/legacy/capabilities.md`

## 扩展约定

以后新增能力时，按飞书技能文件架构落盘：

1. `SKILL.md` 只放使用边界、规则、入口和 reference 索引。
2. 工作流放到 `references/workflows/`。
3. 平台细节放到 `references/platforms/`。
4. CLI 规范、兼容策略、迁移说明放到 `references/cli/` 或 `references/legacy/`。
5. 真正可执行能力统一在 `yxer` CLI 中扩展，不再新增新的默认 Node 执行入口。
