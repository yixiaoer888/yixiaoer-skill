---
name: yixiaoer
version: 1.5.0
description: "蚁小二支持 50 + 平台一键分发、多账号矩阵管理、团队协作与数据统计，适配图文、短视频等内容，覆盖全平台客户端，助力个人与企业高效运营新媒体矩阵。"
author: wangzhengjiao
---

# OpenClaw 龙虾技能 (OpenClaw Skill)

本项目采用 **DTO 驱动 (Interface-Driven)** 的架构模式。所有的功能均可通过 `api.ts` 配合标准化的 DTO Payload 进行调用。

## 核心原则 (Core Principles)

1.  **接口即标准**: 所有的功能调用严格基于后端 API 的 Request DTO 设计。
2.  **文档即指引**: `docs/` 下的 Markdown 文档用于解释对应接口的参数规则、必填项与约束。
3.  **零映射透传**: 鼓励调用者使用 `api.ts` 透明地提交符合 DTO 要求的 JSON Payload。

> [!IMPORTANT]
> **严格合规性 (Strict Compliance)**:
> 1. 所有接口调用**必须**严格遵守各文档中定义的**必填字段 (`必填: 是`)** 以及对应的**数据格式要求**（如时间戳、数组嵌套结构等）。
> 2. **资源引用规范**: 所有的**封面图 (cover)**、**图文图片 (images)** 以及 **视频文件 (video)** 必须先通过[资源上传接口](./docs/upload-resource.md)上传至系统并获得唯一的 `key`。禁止填入非系统内的网络 URL 或随意留空，否则会导致发布任务执行失败。


## 平台支持 (Platform Support)

API 调用时涉及的平台名称必须使用蚁小二定义的中文枚举或 Code。

*   **平台枚举列表**: [docs/platform.md](./docs/platform.md)

## 统一执行指令 (Unified Execution)

所有的 API 功能统一通过 [api.ts](./scripts/api.ts) 脚本执行。调用时需通过 `--payload` 参数传入 JSON，且 **`action` 字段为必填项**，用于指定具体功能。

| action 值 | 对应的能力描述 | 相关文档 |
| :--- | :--- | :--- |
| `publish` | 图文、视频、文章分发发布 | [docs/publish/index.md](./docs/publish/index.md) |
| `accounts` | 查询已绑定的账号列表 | [query-accounts.md](./docs/query-accounts.md) |
| `upload` | 上传本地或 URL 图片/视频 | [upload-resource.md](./docs/upload-resource.md) |
| `records` | 查询发布任务概览列表 | [get-publish-records.md](./docs/get-publish-records.md) |
| `details` | 查询特定任务的执行详情 | [get-publish-records.md](./docs/get-publish-records.md) |
| `team-info`| 获取租户与配额信息 | [get-team-info.md](./docs/get-team-info.md) |
| `categories`| 获取账号分类/话题列表 | [get-publish-categories.md](./docs/get-publish-categories.md) |
| `activities`| 获取征文活动列表 | [get-publish-activities.md](./docs/get-publish-activities.md) |
| `locations` | 获取 POI 物理位置列表 | [get-locations.md](./docs/get-locations.md) |
| `music` | 获取抖音/快手可选背景音乐 | [get-music.md](./docs/get-music.md) |
| `collections`| 获取账号已创建的合集列表 | [get-collections.md](./docs/get-collections.md) |

### 调用示例 (Example)

```bash
# 查询账号列表 (action: accounts)
node scripts/api.ts --payload='{"action": "accounts", "platform": "抖音"}'
```

## 开发指南 (Development Guide)

为了简化 API 的调用与脚本开发，我们提供了通用的 **API 助手模块**：

*   **API 助手模块**: `scripts/api.ts`
*   **开发指引文档**: [API 助手使用指南](./docs/scripts/api-guide.md)

在开发新功能或修改现有脚本时，请务必参考此指引。

---
> [!NOTE]
> 所有的敏感信息应通过环境变量 `YIXIAOER_API_KEY` 注入。
