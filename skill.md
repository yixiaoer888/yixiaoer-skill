---
name: openclaw-skill-core
version: 1.4.0
description: "蚁小二全平台媒体管理与运营能力，通过 DTO 驱动型文档与共享引擎实现发布过程原子化。"
author: wangzhengjiao
---

# OpenClaw 龙虾技能 (OpenClaw Skill)

## 配置与安全 (Config & Secrets)

所有的敏感信息应通过**环境变量**注入：

1.  **生产环境**: 在龙虾系统 (OpenClaw) 的环境变量配置中填入 `YIXIAOER_API_KEY`。
2.  **本地开发**: 
    - 运行脚本时，Node.js 20.6+ 可以使用内置标志加载：`node --env-file=.env scripts/xxx.ts`。

## 通用发布指令接口 (Universal Publish CLI)

所有发布任务均通过 `scripts/publish.ts` 执行。引擎会自动将 CLI 参数转换为后端 DTO 要求的复杂结构。

### 1. 核心映射逻辑 (Core Mapping)

引擎将指令参数映射为 `CloudTaskPushRequest` 结构：

| 指令参数 | DTO 字段 | 映射目标 | 说明 |
| :--- | :--- | :--- | :--- |
| `--type` | `publishType` | Root | `article`, `imageText`, `video` |
| `--platforms` | `platforms` | Root | 目标平台数组 |
| `--title` | `desc` | Root | 任务集描述/标题 |
| `--content` | `content` | `publishArgs` | **关键**: 文章和图文必须在 `publishArgs` 顶层填入原始内容 |
| `--account_ids`| - | `accountForms`| 拆分为多个账号级对象 |
| `--cover_url` | `cover` | `accountForms`| **必填**: 账号级对象必须包含 `cover` (ImageFormItem) |
| `--image_urls`| `images` | `accountForms`| 图文发布时填入账号级 `images` 列表 |
| `--video_url` | `video` | `accountForms`| 视频发布时填入账号级 `video` 对象 |
| `--client_id` | `clientId` | Root | **RPA 必填**: 当 channel=local 时，对应本地客户端连接 ID |
| `--channel` | `publishChannel`| Root | `cloud` (默认), `local` (本机/RPA) |
| `--task_set_id`| `taskSetId` | Root | 指定任务集 ID (用于更新草稿或特定管理) |
| `--is_app` | `isAppContent` | Root | `true` 或 `false` |
| `--media_id` | `mediaId` | `accountForms`| 使用平台素材库 ID (如已上传) |
| `--fps` | `fps` | `accountForms`| 视频发布帧率 (海外平台必填项) |
| `--interval` | `interval` | `intervalConfig`| 间隔发布数值 (数字) |
| `--time_unit` | `timeUnit` | `intervalConfig`| `minute` (5,10,15,20,30), `day` (1,2,3,4,5) |
| `--start_time` | `dailyStartTime`| `intervalConfig`| 按天发布时的起始时间 (HH:MM) |
| `--rotation` | `accountRotation`| `intervalConfig`| 是否开启账号轮询 (`true`/`false`) |
| `[Other]` | `*` | `contentPublishForm`| 所有不在上述核心列表的参数透传到平台 DTO |

### 2. 标准表单项模型 (Standard Form Models)

当 DTO 中引用以下类时，引擎会按照以下默认结构构造数据：

#### ImageFormItem (图片/封面)
```json
{
  "key": "oss_key",
  "path": "http://...",
  "width": 1200, 
  "height": 800,
  "size": 0
}
```

#### VideoFormItem (视频)
```json
{
  "key": "oss_key",
  "path": "http://...",
  "duration": 0,
  "width": 1920,
  "height": 1080,
  "size": 0
}
```
> [!NOTE]
> 引擎会自动处理资源上传。若指令中只提供 URL，引擎会获取 Key 后填入。

## 能力地图 (Capabilities)

本技能通过映射 `docs/` 下的指令文档到 `scripts/` 下的执行脚本实现功能的动态调度。
为确保表单知识准确，**每个平台+发布模态均拥有独立的指令文档**，其字段逻辑源自后端 DTO。

| 能力名称 | 指令文档 (Trigger) | 执行脚本 (Implementation) | 核心功能 |
| :--- | :--- | :--- | :--- |
| **查询账号列表** | [query-accounts.md](./docs/query-accounts.md) ([平台定义](./docs/platform.md)) | [query-accounts.ts](./scripts/query-accounts.ts) | 获取租户下绑定的媒体账号 |
| **查询发布记录** | [get-publish-records.md](./docs/get-publish-records.md) | [get-publish-records.ts](./scripts/get-publish-records.ts) | 获取发布任务的详细记录与状态 |
| **当前团队信息** | [get-team-info.md](./docs/get-team-info.md) | [get-team-info.ts](./scripts/get-team-info.ts) | 获取团队名称、角色、额度信息 |
| **查询发布分类** | [get-publish-categories.md](./docs/get-publish-categories.md) | [get-publish-categories.ts](./scripts/get-publish-categories.ts) | 获取账号下的分类列表（百家号/公众号等） |
| **查询地理位置** | [get-locations.md](./docs/get-locations.md) | [get-locations.ts](./scripts/get-locations.ts) | 获取发布可选的地址/带货地址 |
| **查询视频音乐** | [get-music.md](./docs/get-music.md) | [get-music.ts](./scripts/get-music.ts) | 获取发布可选的音乐素材 |
| **查询合集列表** | [get-collections.md](./docs/get-collections.md) | [get-collections.ts](./scripts/get-collections.ts) | 获取账号已创建的合集 (抖音/头条) |
| **查询征文活动** | [get-publish-activities.md](./docs/get-publish-activities.md) | [get-publish-activities.ts](./scripts/get-publish-activities.ts) | 获取账号下的可参与活动（百家号等） |
| **上传资源** | [upload-resource.md](./docs/upload-resource.md) | [upload-resource.ts](./scripts/upload-resource.ts) | **基础能力**: 将文件或 URL 直传蚁小二 OSS |
| **文章发布 (通用)** | [支持 20+ 文章平台](#文章发布平台支持情况) | [publish.ts](./scripts/publish.ts) | 支持 30+ 长文平台自动分发 |
| **图文发布 (通用)** | [支持 8+ 图文平台](#图文发布平台支持情况) | [publish.ts](./scripts/publish.ts) | 支持抖音/快手/小红书/微博/视频号/百家/头条/知乎等动态平台 |
| **视频发布 (通用)** | [支持 30+ 视频平台](#视频发布平台支持情况) | [publish.ts](./scripts/publish.ts) | 支持抖音/快手/B站/视频号等全平台自动化分发 |

### 视频发布平台支持情况

| 平台名称 | 文档路径 | 平台名称 | 文档路径 |
| :--- | :--- | :--- | :--- |
| **抖音 (Douyin)** | [douyin.md](./docs/publish-video/douyin.md) | **头条号 (Toutiao)** | [toutiaohao.md](./docs/publish-video/toutiaohao.md) |
| **哔哩哔哩 (Bilibili)** | [bilibili.md](./docs/publish-video/bilibili.md) | **哔哩哔哩-Open** | [bilibili-open.md](./docs/publish-video/bilibili-open.md) |
| **百家号 (Baijiahao)** | [baijiahao.md](./docs/publish-video/baijiahao.md) | **小红书 (Xiaohongshu)** | [xiaohongshu.md](./docs/publish-video/xiaohongshu.md) |
| **小红书商家号** | [xiaohongshu-shop.md](./docs/publish-video/xiaohongshu-shop.md) | **快手 (Kuaishou)** | [kuaishou.md](./docs/publish-video/kuaishou.md) |
| **快手-Open** | [kuaishou-open.md](./docs/publish-video/kuaishou-open.md) | **新浪微博 (Weibo)** | [weibo.md](./docs/publish-video/weibo.md) |
| **视频号 (Video Account)** | [shipinhao.md](./docs/publish-video/shipinhao.md) | **知乎 (Zhihu)** | [zhihu.md](./docs/publish-video/zhihu.md) |
| **企鹅号 (Qiehao)** | [qiehao.md](./docs/publish-video/qiehao.md) | **爱奇艺 (iQIYI)** | [aiqiyi.md](./docs/publish-video/aiqiyi.md) |
| **网易号 (Wangyi)** | [wangyihao.md](./docs/publish-video/wangyihao.md) | **一点号 (Yidianhao)** | [yidianhao.md](./docs/publish-video/yidianhao.md) |
| **搜狐号 (Sohuhao)** | [souhuhao.md](./docs/publish-video/souhuhao.md) | **搜狐视频 (Sohu Video)** | [souhushipin.md](./docs/publish-video/souhushipin.md) |
| **腾讯微视 (Weishi)** | [weishi.md](./docs/publish-video/weishi.md) | **皮皮虾 (Pipixia)** | [pipixia.md](./docs/publish-video/pipixia.md) |
| **腾讯视频 (Tencent)** | [v-qq.md](./docs/publish-video/v-qq.md) | **多多视频 (Duo Duo)** | [duoduoshipin.md](./docs/publish-video/duoduoshipin.md) |
| **美拍 (Meipai)** | [meipai.md](./docs/publish-video/meipai.md) | **AcFun** | [acfun.md](./docs/publish-video/acfun.md) |
| **大鱼号 (Dayuhao)** | [dayuhao.md](./docs/publish-video/dayuhao.md) | **车家号 (Chejiahao)** | [chejiahao.md](./docs/publish-video/chejiahao.md) |
| **蜂网 (Fengwang)** | [fengwang.md](./docs/publish-video/fengwang.md) | **得物 (Dewu)** | [dewu.md](./docs/publish-video/dewu.md) |
| **美柚 (Meiyou)** | [meiyou.md](./docs/publish-video/meiyou.md) | **易车号 (Yichehao)** | [yichehao.md](./docs/publish-video/yichehao.md) |

### 文章发布平台支持情况

| 平台名称 | 文档路径 | 平台名称 | 文档路径 |
| :--- | :--- | :--- | :--- |
| **微信公众号** | [wxgongzhonghao.md](./docs/publish-article/wxgongzhonghao.md) | **百家号 (Baijiahao)** | [baijiahao.md](./docs/publish-article/baijiahao.md) |
| **今日头条 (Toutiao)** | [toutiao.md](./docs/publish-article/toutiao.md) | **知乎 (Zhihu)** | [zhihu.md](./docs/publish-article/zhihu.md) |
| **新浪微博 (Weibo)** | [xinlangweibo.md](./docs/publish-article/xinlangweibo.md) | **网易号 (Wangyi)** | [wangyihao.md](./docs/publish-article/wangyihao.md) |
| **大鱼号 (Dayuhao)** | [dayuhao.md](./docs/publish-article/dayuhao.md) | **一点号 (Yidianhao)** | [yidianhao.md](./docs/publish-article/yidianhao.md) |
| **企鹅号 (Qiehao)** | [qiehao.md](./docs/publish-article/qiehao.md) | **搜狐号 (Sohuhao)** | [souhuhao.md](./docs/publish-article/souhuhao.md) |
| **哔哩哔哩 (Bilibili)** | [bilibili.md](./docs/publish-article/bilibili.md) | **CSDN** | [csdn.md](./docs/publish-article/csdn.md) |
| **简书 (Jianshu)** | [jianshu.md](./docs/publish-article/jianshu.md) | **雪球号 (Xueqiu)** | [xueqiuhao.md](./docs/publish-article/xueqiuhao.md) |
| **豆瓣 (Douban)** | [douban.md](./docs/publish-article/douban.md) | **快传号 (Kuaichuan)** | [kuaichuanhao.md](./docs/publish-article/kuaichuanhao.md) |
| **抖音 (Douyin)** | [douyin.md](./docs/publish-article/douyin.md) | **爱奇艺 (iQIYI)** | [aiqiyi.md](./docs/publish-article/aiqiyi.md) |
| **车家号 (Chejiahao)** | [chejiahao.md](./docs/publish-article/chejiahao.md) | **易车号 (Yichehao)** | [yichehao.md](./docs/publish-article/yichehao.md) |
| **WiFi万能钥匙** | [wifiwanneng.md](./docs/publish-article/wifiwanneng.md) | **AcFun** | [acfun.md](./docs/publish-article/acfun.md) |

### 图文发布平台支持情况

| 平台名称 | 文档路径 | 平台名称 | 文档路径 |
| :--- | :--- | :--- | :--- |
| **小红书 (Xiaohongshu)** | [xiaohongshu.md](./docs/publish-image-text/xiaohongshu.md) | **抖音 (Douyin)** | [douyin.md](./docs/publish-image-text/douyin.md) |
| **快手 (Kuaishou)** | [kuaishou.md](./docs/publish-image-text/kuaishou.md) | **视频号 (Video Account)** | [shipinhao.md](./docs/publish-image-text/shipinhao.md) |
| **新浪微博 (Weibo)** | [weibo.md](./docs/publish-image-text/weibo.md) | **百家号 (Baijiahao)** | [baijiahao.md](./docs/publish-image-text/baijiahao.md) |
| **今日头条 (Toutiao)** | [toutiaohao.md](./docs/publish-image-text/toutiaohao.md) | **知乎 (Zhihu)** | [zhihu.md](./docs/publish-image-text/zhihu.md) |

## DTO 知识提取规范 (DTO Extraction Specs)

为确保 AI 助手在生成指令文档时**完整、无遗漏**地提取参数，必须遵循以下 DTO 阅读准则：

### 1. 目标类定位 (Target Identification)
根据业务模态，在 `<platform_id>.dto.ts` 中定位对应的表单类：
- **文章发布 (`article`)**: 提取 `*ArticleForm` 类（如 `DouyinArticleForm`）。
- **视频发布 (`video`)**: 提取 `*VideoForm` 类（如 `DouYinVideoForm`）。
- **图文/动态 (`image-text`)**: 提取 `*DynamicForm` 类（如 `DouYinDynamicForm`）。
- **注意**: 若类名不符合上述规律，请查找所有继承自 `PlatformFormBaseDTO` 的类。

### 2. 核心结构审计 (Core Structure Audit)
**必须深度审计以下文件**以获取完整表单知识：
1.  **后端 DTO** (`*.dto.ts`): 
    - **扫描 ApiProperty**: 必须提取类中**所有**带有 `@ApiProperty` 的字段。
    - **标准对象识别**: 若字段类型引用了 `ImageFormItem` 或 `VideoFormItem`（见上文模型），文档中只需指明为“标准图片/视频对象”即可。
    - **规则还原**: 详细记录 `@IsNotEmpty` (必填), `@IsIn` (枚举范围), `@IsInt` (类型), `@MaxLength` (长度限制) 等逻辑。
    - **嵌套对象**: 若字段类型是 class (且非标准项), 必须点击跳转查看该类的具体字段，并在文档中详细说明该对象的 JSON 结构。

### 3. 上层结构知识 (Upper-Level Mapping)

AI 助手必须确保生成的请求符合以下 `CloudTaskPushRequest` 的嵌套层级：

```json
{
  "desc": "任务集标题",
  "publishType": "article | video | imageText",
  "publishArgs": {
    "content": "原始内容 (针对文章/图文必填)",
    "accountForms": [
      {
        "platformAccountId": "账号ID",
        "cover": { "key": "...", "width": 1200, "height": 800 }, // 必填
        "images": [], // 图文必填
        "video": {},  // 视频必填
        "contentPublishForm": {
          // 此处填入 DTO 提取规范中的平台专有参数
        }
      }
    ]
  }
}
```
> [!IMPORTANT]
> 注意 `content` 的双重位置：原始内容在 `publishArgs.content`；平台特有格式（如 HTML 或带标签的描述）在 `contentPublishForm`。

### 4. 禁止项 (Strict Prohibitions)
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
4. **批量发布策略**: 若任务涉及多个目标账号，**必须优先通过单个请求**（在 `accountForms` 数组中包含多个账号对象）完成派发。除非用户有明确的拆分需求，严禁为每个账号重复调用发布接口。
