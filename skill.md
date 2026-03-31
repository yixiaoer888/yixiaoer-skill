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

## 能力地图 (Capabilities)

所有的能力均通过 [api.ts](./scripts/api.ts) 中的 `callApi` 配合对应的 DTO Payload 实现。

| 能力类型 | 功能描述 | 请求 Endpoint | 相关文档 |
| :--- | :--- | :--- | :--- |
| **内容发布** | 视频、图文、文章全平台发布 | `/taskSets/v2` | [docs/publish/index.md](./docs/publish/index.md) |
| **账号管理** | 查询全平台账号列表与 UID | `/v2/platform/accounts` | [query-accounts.md](./docs/query-accounts.md) |
| **资源管理** | 图片/视频上传至 OSS 获取 Key | `/storages/[bucket]/upload-url` | [upload-resource.md](./docs/upload-resource.md) |
| **发布记录** | 查询历史发布任务与各平台状态 | `/v2/taskSets` | [get-publish-records.md](./docs/get-publish-records.md) |
| **发布详情** | 获取特定任务的详细执行记录 | `/v2/taskSets/[id]/tasks` | [get-publish-details.md](./docs/get-publish-details.md) |
| **团队信息** | 获取租户/团队基本信息与额度 | `/v2/teams/current` | [get-team-info.md](./docs/get-team-info.md) |
| **分类查询** | 获取账号下的分类/合集/话题 | `/web/config-data/category-tasks` | [get-publish-categories.md](./docs/get-publish-categories.md) |
| **活动查询** | 获取平台当前的征文活动列表 | `/web/config-data/activity-tasks` | [get-publish-activities.md](./docs/get-publish-activities.md) |
| **地理位置** | 获取发布可选的 POI 地址列表 | `/web/config-data/location-tasks` | [get-locations.md](./docs/get-locations.md) |
| **音乐素材** | 获取抖音/快手发布可选音乐 | `/web/config-data/music-tasks` | [get-music.md](./docs/get-music.md) |
| **合集管理** | 获取账号已创建的合集列表 | `/web/config-data/collection-tasks` | [get-collections.md](./docs/get-collections.md) |

## 开发指南 (Development Guide)

为了简化 API 的调用与脚本开发，我们提供了通用的 **API 助手模块**：

*   **API 助手模块**: `scripts/api.ts`
*   **开发指引文档**: [API 助手使用指南](./docs/scripts/api-guide.md)

在开发新功能或修改现有脚本时，请务必参考此指引。

---
> [!NOTE]
> 所有的敏感信息应通过环境变量 `YIXIAOER_API_KEY` 注入。
