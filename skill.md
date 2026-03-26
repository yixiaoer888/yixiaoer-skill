# OpenClaw 龙虾技能 (OpenClaw Skill)

该技能定义了蚁小二全平台的媒体管理与运营能力。
通过元数据驱动（Skill -> Doc -> Script）模式，将发布过程原子化。

## 技能定义 (Metadata)

- **ID**: `openclaw-skill-core`
- **版本**: `1.0.4`
- **架构模式**: 文档驱动型自律脚本 (Doc-Driven Scripts)
- **运行环境**: Node.js v18+ (Direct Runtime)

## 配置与安全 (Config & Secrets)

所有的敏感信息应通过**环境变量**注入：

1.  **生产环境**: 在龙虾系统 (OpenClaw) 的环境变量配置中填入 `YIXIAOER_API_KEY`。
2.  **本地开发**: 
    - 运行脚本时，Node.js 20.6+ 可以使用内置标志加载：`node --env-file=.env scripts/xxx.ts`。

## 能力地图 (Capabilities)

本技能通过映射 `docs/` 下的指令文档到 `scripts/` 下的执行脚本实现功能的动态调度。

| 能力名称 | 指令文档 (Trigger) | 执行脚本 (Implementation) | 核心功能 |
| :--- | :--- | :--- | :--- |
| **查询账号列表** | [query-accounts.md](./docs/query-accounts.md) | [query-accounts.ts](./scripts/query-accounts.ts) | 获取租户下绑定的媒体账号 |
| **查询发布记录** | [get-publish-records.md](./docs/get-publish-records.md) | [get-publish-records.ts](./scripts/get-publish-records.ts) | 获取发布任务的详细记录与状态 |
| **当前团队信息** | [get-team-info.md](./docs/get-team-info.md) | [get-team-info.ts](./scripts/get-team-info.ts) | 获取团队名称、角色、额度信息 |
| **上传资源** | [upload-resource.md](./docs/upload-resource.md) | [upload-resource.ts](./scripts/upload-resource.ts) | **核心能力**: 将文件或 URL 直传蚁小二 OSS |
| **发布文章 (Article Engine)** | [publish-article.md](./docs/publish-article.md) | [publish-article.ts](./scripts/publish-article.ts) | **统一引擎**: 支持 22+ 长文平台 (含微信公众号/头条/百家号等) |
| **发布图文 (Post Engine)** | [publish-post.md](./docs/publish-post.md) | [publish-post.ts](./scripts/publish-post.ts) | **统一引擎**: 支持抖音/快手/小红书/微博等动态平台 |
| **发布视频 (Video Engine)** | [publish-video.md](./docs/publish-video.md) | [publish-video.ts](./scripts/publish-video.ts) | **统一引擎**: 支持 30+ 视频平台 (B站/西瓜/视频号/抖音/快手等) |


## 任务执行最佳实践 (Best Practices)

在处理复杂的发布任务时，AI 助手应遵循以下工作流：
1. **获取账号**: 首先调用其“账号查询”能力，确认目标账号 ID。
2. **预处理素材**: 如果存在外部图片/视频，**优先调用“上传资源”能力**，获取各个素材的 Key。避免在发布步骤中进行耗时的实时上传。
3. **最终派发**: 将获取到的各级 Resource Keys 传递给特定的“发布”原子能力进行最终提交。
