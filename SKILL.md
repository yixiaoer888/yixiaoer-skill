---
name: yixiaoer
version: 3.0.0
description: "蚁小二多平台内容分发智能助手。支持图文/视频/文章一键发布到50+平台，提供账号管理、资源上传、数据查询等能力。参照飞书CLI架构设计，确保Agent执行稳定可预期。"
author: wangzhengjiao
---

# 蚁小二多平台内容分发

做智能助理，不做表单填写机。能自动补全的默认值先补全，只在必须决策时才追问用户。

## ⚠️ BLOCKING REQUIREMENT（最高优先级，任何场景均不得违反）

1. 【发布类操作】执行前**必须**先读取对应场景的工作流文档（workflows/ 目录），未读取前绝对禁止执行任何发布操作！
2. 所有图片/视频资源**必须**先通过 `yxer upload` 获取 key，严禁在 Payload 中直接使用外部 URL
3. 发布前**必须**通过 `yxer accounts` 确认目标账号 `status=1`（在线）
4. 复杂对象（`location`/`music`/`collection`/`challenge`）**严禁**手动构造，必须通过对应查询命令获取完整 `raw` 对象
5. 生成 Payload 后**必须**先 `yxer validate` 校验通过，再执行 `yxer publish`，禁止跳过校验
6. 【草稿识别】用户提及"草稿"时，必须先询问是"蚁小二草稿"还是"平台草稿"，严禁自行猜测

---

## 核心场景

### 发布内容

⚠️ **BLOCKING**: 根据发布类型，必须先读取对应工作流文档！未读取前不得执行任何发布操作！

| 场景 | 工作流文档 | 触发条件 |
|------|-----------|---------|
| 发布图文（图片+文字动态/笔记） | [workflows/publish-image-text.md](./workflows/publish-image-text.md) | "发图文"、"发笔记"、"发动态"、"配图发" |
| 发布视频 | [workflows/publish-video.md](./workflows/publish-video.md) | "发视频"、"上传视频" |
| 发布文章（长图文） | [workflows/publish-article.md](./workflows/publish-article.md) | "发文章"、"写文章"、"发长文" |

### 查询类操作（直接执行，无需读工作流）

| 操作 | 命令 | 说明 |
|------|------|------|
| 查看账号列表 | `yxer accounts [platform]` | 支持按平台筛选，加 `--name 关键词` 模糊匹配昵称 |
| 查看发布记录 | `yxer records [--platform P] [--limit N]` | 默认最近10条 |
| 查看分类列表 | `yxer categories <account_id> [--type video\|article]` | 发布前获取可选分类 |
| 搜索地理位置 | `yxer locations <account_id> [--query 关键词]` | POI位置，用于图文/视频挂载 |
| 搜索背景音乐 | `yxer music <account_id> [--query 关键词]` | 用于视频/图文挂载音乐 |
| 查看商品列表 | `yxer goods <account_id> [--query 关键词]` | 用于挂车 |
| 查看合集列表 | `yxer collections <account_id> [--type video]` | 用于视频加入合集 |
| 查看话题/挑战 | `yxer challenges <account_id> [--query 关键词]` | 用于视频/图文加话题 |

### 资源上传

```bash
yxer upload <file_path_or_url> [--bucket cloud-publish|material-library]
```

- 自动推断文件类型（图片/视频），无需手动指定 `content_type`
- 图片自动返回 `key`/`size`/`width`/`height`/`format`
- 视频自动返回 `key`/`size`/`width`/`height`/`duration`/`format`
- 支持本地文件路径和 HTTP/HTTPS URL

---

## 命令参考

### 完整命令列表

```bash
yxer <command> [options]

Commands:
  accounts [platform] [--name 关键词] [--status 1] [--json]     查询账号
  upload <path_or_url> [--bucket cloud-publish|material-library]  上传资源
  validate <platform> <type> <payload.json>                       校验Payload
  publish <type> <platforms> <payload.json> [clientId]           发布内容
  categories <account_id> [--type video|article]                  查分类
  locations <account_id> [--query 关键词] [--type 0|1|2|3]      查位置
  music <account_id> [--query 关键词]                             查音乐
  goods <account_id> [--query 关键词]                             查商品
  collections <account_id> [--type video|article]                 查合集
  challenges <account_id> [--query 关键词] [--type video]         查话题
  records [--platform P] [--limit N] [--status S] [--json]      查发布记录
  prepare <platform> <type>                                       一步获取发布前置数据

Options:
  --json     输出JSON格式（默认人类可读格式）
  --debug    显示详细日志
  --help     显示帮助
```

### 底层 API 调用（高级用法，一般不推荐）

```bash
node scripts/api.ts --payload='{"action":"...", ...}'
```

仅在 `yxer` CLI 不支持的场景下使用。大多数场景请优先使用 `yxer` 命令。

---

## JSON Schema 验证参考

`yxer validate <platform> <type> <payload.json>` 使用 JSON Schema (draft-07) 对 payload 进行校验。
Schema 文件位于 `schemas/platforms/` 目录，命名格式：`{platform}.{type}.schema.json`。

**当前已覆盖 56 个 schema（22 文章 + 8 图文 + 26 视频平台）**：

| 类型 | 覆盖平台 |
|------|----------|
| 文章 (22) | acfun, aiqiyi, baijiahao, bilibili, chejiahao, csdn, dayuhao, douban, douyin, jianshu, qiehao, souhuhao, toutiaohao, wangyihao, weixin.account, xinlang, xueqiuhao, yichehao, yidianhao, xinlang.article, weixin.account(公众号) |
| 图文 (8) | douyin, kuaishou, weixin.shipinhao, xhs, xinlang, zhihu, baijiahao, toutiaohao |
| 视频 (26) | acfun, aiqiyi, bilibili, douyin, kuaishou, toutiaohao, xinlang, chejiahao, dayuhao, duoduoshipin, fengwang, meipai, meiyou, qiehao, shipinghao, baijiahao, kuaishou-open, souhushipin, wangyihao, xiaohongshu, xiaohongshushop, yidianhao, yichehao, zhihu, weishi, pipixia, dewu, tengxunshipin |

> [!TIP]
> Agent 构造 payload 后，先用 `yxer validate` 校验，再执行 `yxer publish`。
> Schema 校验失败会明确提示缺少哪些 required 字段或多余的 additionalProperties。

---

## 平台参考文档

各平台的 `contentPublishForm` 差异字段定义（Agent 在构造 Payload 时按需查阅）：

- [图文发布通用索引](./docs/publish/image-text/index.md) — 所有图文平台的根结构定义
- [视频发布通用索引](./docs/publish/video/index.md) — 所有视频平台的根结构定义
- [文章发布通用索引](./docs/publish/article/index.md) — 所有文章平台的根结构定义
- 各平台详细参数：查阅 `docs/publish/<类型>/<平台>.md`

运维与故障排查：
- [蚁小二 Skill 严格执行标准](./docs/execution-standard.md)
- [蚁小二 Skill 避坑与故障排查手册](./docs/troubleshooting-guide.md)

---

## 通用规则（Agent 必须遵守）

### 智能助理原则

- 能自动补全的默认值**直接补全**，不要反复追问用户
- 只在必须用户决策时才提问（选账号、确认发布摘要）
- 用户描述模糊时，根据上下文推断；推断不出再问

### 默认值自动补全规则

| 字段 | 默认值 | 是否需确认 |
|------|--------|------------|
| `formType` | `"task"` | 无需确认，直接填 |
| `publishChannel` | `"cloud"` | 无需确认，直接填 |
| `cover` / `coverKey` | 使用 `images[0]` 的 key | 无需确认，直接填 |
| `title` | 从用户描述中提取 | 展示给用户确认 |
| `description` | 从标题或用户描述自动生成 | 展示给用户确认 |
| `scheduledTime` | 不填（立即发布） | 用户明确要求定时才填 |

### 严禁行为

- ❌ 未确认账号在线（`status=1`）就构造 Payload
- ❌ 手动编造图片/视频的 `key`/`size`/`width`/`height`
- ❌ 手动构造 `location`/`music`/`collection`/`challenge` 的 `raw` 字段
- ❌ 跳过 `yxer validate` 直接 `yxer publish`
- ❌ 在 Payload 中直接使用外部 URL 作为图片/视频地址
- ❌ 跳过工作流文档，自行拼 JSON Payload
- ❌ 用户说"草稿"时不询问类型，自行猜测是蚁小二草稿还是平台草稿

---

> [!TIP]
> **快速开始**: 用户说"帮我发一条抖音图文"→ 读 `workflows/publish-image-text.md` → 按 Step 1→2→3→4→5→6 执行
