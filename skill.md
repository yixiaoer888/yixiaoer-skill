---
name: yixiaoer
version: 1.5.0
description: "蚁小二支持 50 + 平台一键分发、多账号矩阵管理、团队协作与数据统计，适配图文、短视频等内容，覆盖全平台客户端，助力个人与企业高效运营新媒体矩阵。"
author: wangzhengjiao
---

# OpenClaw 龙虾技能 (OpenClaw Skill)

本项目采用 **DTO 驱动 (Interface-Driven)** 的架构模式。所有的发布能力均通过通用的执行脚本与标准化的 DTO 文档协同完成。

## 核心原则 (Core Principles)

1.  **接口即标准**: `scripts/` 下的脚本严格基于后端 API 的 Request DTO 设计。
2.  **文档即指引**: `docs/` 下的 Markdown 文档用于解释对应接口的参数规则、必填项与约束。
3.  **零映射透传**: 鼓励调用者直接构造符合 DTO 要求的 JSON Payload 进行提交。


## 平台支持 (Platform Support)

API 调用时涉及的平台名称必须使用蚁小二定义的中文枚举或 Code。

*   **平台枚举列表**: [docs/platform.md](./docs/platform.md)

## 能力地图 (Capabilities)

| 能力类型 | 功能描述 | 执行脚本 | 文档指引 |
| :--- | :--- | :--- | :--- |
| **内容发布** | 视频、图文、文章全平台发布 | `publish.ts` | [docs/publish/index.md](./docs/publish/index.md) |
| **账号管理** | 查询全平台账号列表与 UID | `query-accounts.ts` | [query-accounts.md](./docs/query-accounts.md) |
| **资源管理** | 图片/视频上传至 OSS 获取 Key | `upload-resource.ts` | [upload-resource.md](./docs/upload-resource.md) |
| **发布记录** | 查询历史发布任务与各平台状态 | `get-publish-records.ts`| [get-publish-records.md](./docs/get-publish-records.md) |
| **发布详情** | 获取特定任务的详细执行记录 | `get-publish-details.ts`| [get-publish-details.md](./docs/get-publish-details.md) |
| **团队信息** | 获取租户/团队基本信息与额度 | `get-team-info.ts` | [get-team-info.md](./docs/get-team-info.md) |
| **分类查询** | 获取账号下的分类/合集/话题 | `get-publish-categories.ts`| [get-publish-categories.md](./docs/get-publish-categories.md) |
| **活动查询** | 获取平台当前的征文活动列表 | `get-publish-activities.ts`| [get-publish-activities.md](./docs/get-publish-activities.md) |
| **地理位置** | 获取发布可选的 POI 地址列表 | `get-locations.ts` | [get-locations.md](./docs/get-locations.md) |
| **音乐素材** | 获取抖音/快手发布可选音乐 | `get-music.ts` | [get-music.md](./docs/get-music.md) |
| **合集管理** | 获取账号已创建的合集列表 | `get-collections.ts` | [get-collections.md](./docs/get-collections.md) |

---
> [!NOTE]
> 所有的敏感信息应通过环境变量 `YIXIAOER_API_KEY` 注入。
