# OpenClaw 龙虾技能 (OpenClaw Skill)

该技能定义了蚁小二全平台的媒体管理与运营能力。
通过元数据驱动（Skill -> Doc -> Script）模式，将发布过程原子化。

## 技能定义 (Metadata)

- **ID**: `openclaw-skill-core`
- **版本**: `1.3.0`
- **架构模式**: DTO 驱动型文档与共享引擎 (DTO-Driven & Shared Engine)
- **运行环境**: Node.js v18+ (Direct Runtime)

## 配置与安全 (Config & Secrets)

所有的敏感信息应通过**环境变量**注入：

1.  **生产环境**: 在龙虾系统 (OpenClaw) 的环境变量配置中填入 `YIXIAOER_API_KEY`。
2.  **本地开发**: 
    - 运行脚本时，Node.js 20.6+ 可以使用内置标志加载：`node --env-file=.env scripts/xxx.ts`。

## 通用发布指令接口 (Universal Publish CLI)

所有发布任务均通过 `scripts/publish.ts` 执行，支持以下核心通用参数：

| 参数 | 对应 DTO 逻辑 | 说明 |
| :--- | :--- | :--- |
| `--type` | `publishType` | 必填。`article`, `image-text`, `video`, `weixin-gongzhonghao` |
| `--platforms` | `platforms` | 必填。目标平台名称列表（逗号分隔）。见 [平台定义](./docs/platform.md) |
| `--account_ids` | `accountIds` | 必填。目标账号 ID 列表（逗号分隔） |
| `--title` | `title` | 标题（文章、视频必填） |
| `--content` | `content/description`| 内容（文章为 HTML，图文视频为描述文本） |
| `--video_url` | 自动上传并重组 | 视频直连 URL，引擎会自动上传至 OSS 并填入 `videoKey` |
| `--image_urls` | 自动上传并重组 | 图片/封面 URL（逗号分隔），引擎会自动上传并填入 `images` |
| `--pubType` | `isDraft` | `1`: 立即发布 (默认), `0`: 保存草稿 |

## 能力地图 (Capabilities)

本技能通过映射 `docs/` 下的指令文档到 `scripts/` 下的执行脚本实现功能的动态调度。
为确保表单知识准确，**每个平台+发布模态均拥有独立的指令文档**，其字段逻辑源自后端 DTO。

| 能力名称 | 指令文档 (Trigger) | 执行脚本 (Implementation) | 核心功能 |
| :--- | :--- | :--- | :--- |
| **查询账号列表** | [query-accounts.md](./docs/query-accounts.md) ([平台定义](./docs/platform.md)) | [query-accounts.ts](./scripts/query-accounts.ts) | 获取租户下绑定的媒体账号 |
| **查询发布记录** | [get-publish-records.md](./docs/get-publish-records.md) | [get-publish-records.ts](./scripts/get-publish-records.ts) | 获取发布任务的详细记录与状态 |
| **当前团队信息** | [get-team-info.md](./docs/get-team-info.md) | [get-team-info.ts](./scripts/get-team-info.ts) | 获取团队名称、角色、额度信息 |
| **查询发布分类** | [get-publish-categories.md](./docs/get-publish-categories.md) | [get-publish-categories.ts](./scripts/get-publish-categories.ts) | 获取账号下的分类列表（百家号/公众号等） |
| **查询征文活动** | [get-publish-activities.md](./docs/get-publish-activities.md) | [get-publish-activities.ts](./scripts/get-publish-activities.ts) | 获取账号下的可参与活动（百家号等） |
| **上传资源** | [upload-resource.md](./docs/upload-resource.md) | [upload-resource.ts](./scripts/upload-resource.ts) | **基础能力**: 将文件或 URL 直传蚁小二 OSS |
| **文章发布 (通用)** | [查看 docs/publish-article 目录](./docs/publish-article/) | [publish.ts](./scripts/publish.ts) | 支持 20+ 长文平台自动分发 |
| **图文发布 (通用)** | [查看 docs/publish-post 目录](./docs/publish-post/) | [publish.ts](./scripts/publish.ts) | 支持抖音/快手/小红书/微博等动态平台 |
| **视频发布 (通用)** | [查看 docs/publish-video 目录](./docs/publish-video/) | [publish.ts](./scripts/publish.ts) | 支持 30+ 视频平台自动化分发 |

## DTO 知识提取规范 (DTO Extraction Specs)

为确保 AI 助手在生成指令文档时**完整、无遗漏**地提取参数，必须遵循以下 DTO 阅读准则：

### 1. 目标类定位 (Target Identification)
根据业务模态，在 `<platform_id>.dto.ts` 中定位对应的表单类（通常继承自 `Platform*ViewDTO`）：
- **文章发布**: 提取 `*ArticleForm` 类（如 `DouyinArticleForm`）。
- **视频发布**: 提取 `*VideoForm` 类（如 `DouYinVideoForm`）。
- **图文/动态**: 提取 `*DynamicForm` 类（如 `DouYinDynamicForm`）。
- **注意**: 若类名不符合上述规律，请查找所有继承自 `PlatformFormBaseDTO` 的类。

### 2. 知识提取 (Knowledge Extraction)
**必须深度审计以下文件**以获取完整表单知识：
1.  **后端 DTO** (`*.dto.ts`): 
    - **定位 Form 类**: 文章找 `*ArticleForm`，视频找 `*VideoForm`，图文找 `*DynamicForm`。
    - **全量扫描**: 必须提取类中**所有**带有 `@ApiProperty` 的字段。
    - **规则还原**: 详细记录 `@IsNotEmpty` (必填), `@IsIn` (枚举范围), `@IsInt` (类型), `@MaxLength` (长度限制) 等逻辑。
    - **嵌套对象**: 若字段类型是 class (如 `Cover`), 必须点击跳转查看该类的具体字段，并在文档中详细说明该对象的 JSON 结构。
2.  **环境变量**: 确保 `YIXIAOER_API_KEY` 在 `openclaw-skill/.env` 中正确配置。

### 4. 对象与嵌套能力处理 (Complex Objects & Capabilities)
**核心规则**:
1.  **模型说明**: 如果平台参数中是一个对象（如 `category`, `activity`, `location`），在指令文档中必须提供该对象的完整 JSON 模型示例。
2.  **递归能力创建**: 
    - 若该对象的数据依赖于其他接口（例如“获取分类列表”、“获取活动列表”），AI 助手**必须**在 `docs/` 下创建一个专门的能力说明文件（如 `get-bjh-categories.md`）。
    - 同时在 `scripts/` 下增加对应的原子脚本（如 `get-bjh-categories.ts`）以支持数据查询。
    - 在主发布文档中，应在参数说明中明确指向该查询能力。

### 5. 禁止项 (Strict Prohibitions)
- **严禁遗漏**: 不得因为属性是“可选”的就忽略提取。可选属性在文档中应标注为 `[可选]`。
- **严禁简写**: 必须保留 DTO 中定义的完整字段名（如 `scheduledTime` 不得简写为 `time`）。
- **严禁自创**: 所有的参数逻辑必须以后端代码为准，不得凭经验猜测。

## 平台与类型扩展规范 (Extension & Sync Specs)

当新增发布**类型** (`type`) 或特定**平台** (`platform`) 时，AI 助手必须确保以下两个维度的同步变更：

### 1. 文档端：参数完整说明 (Full Parameter Documentation)
- **位置**: 在 `docs/publish-*/` 目录下创建或更新对应的平台文档。
- **要求**: 文档必须详细列出该平台/类型下所有的业务参数。每个参数需包含：字段名、中文含义、类型、必填性、枚举值范围及特殊约束。
- **依据**: 必须严格遵循上文的 **DTO 知识提取规范**，确保 AI 在阅读文档后能生成 100% 合规的 API 请求。

### 2. 代码端：脚本逻辑适配 (Script Adaptation)
- **位置**: `scripts/publish.ts`。
- **要求**: 检查通用发布引擎是否能处理该新增平台/类型的特殊要求。
- **适配逻辑**:
  - **字段转换**: 若前端传入的参数名与后端 DTO 要求的结构不一致（例如文章列表包装、封面对象嵌套），需在 `publish.ts` 的“构建业务表单”或“补全基础字段”环节添加对应的转换逻辑。
  - **模态分支**: 对于全新的 `type`（如“直播间推送”），需在 `publish.ts` 中新增处理分支，确保 `taskBody` 的构造符合后端 API 规范。
- **验证**: 适配完成后，必须确保脚本能正确将 Markdown 指令转换为合规的 JSON Payload。

## 任务执行最佳实践 (Best Practices)

在处理复杂的发布任务时，AI 助手应遵循以下工作流：
1. **获取账号**: 首先调用其“账号查询”能力，确认目标账号 ID。
2. **预处理素材**: 如果存在外部图片/视频，**优先调用“上传资源”能力**，获取各个素材的 Key。避免在发布步骤中进行耗时的实时上传。
3. **最终派发**: 将获取到的各级 Resource Keys 传递给特定的“发布”原子能力进行最终提交。
