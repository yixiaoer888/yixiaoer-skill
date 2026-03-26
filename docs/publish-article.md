# 发布文章 (Publish Article Engine)

全平台文章发布统一接口。支持一次性将内容同步分发到多个平台（百家号、企鹅号、头条号等）的多个账号。

## 场景描述 (Usage)

- "帮我把这篇科技文章发布到我的百家号、企鹅号和头条号所有账号。"
- "在多号分发时，将封面设定为 (Key: xxx)，标签为 AI, 趋势。"

## 支持平台 (Supported Platforms)

目前基座支持以下 21 个主流长文章平台的一键分发：
- **核心平台**: 头条号、百家号、企鹅号、网易号、搜狐号、知乎、一点号、大鱼号
- **垂直平台**: CSDN、车家号、易车号、简书、雪球号
- **社交/动态**: 新浪微博、豆瓣、哔哩哔哩、AcFun (A站)、抖音
- **其他**: 爱奇艺、快传号、WiFi万能钥匙

## 参数定义 (Parameters)

### 1. 核心分发参数 (Distribution)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `platforms` | `string[]` | 是 | 目标平台名称列表，如 `["百家号", "企鹅号"]` |
| `account_ids` | `string[]` | 否 | 涉及的所有账号 ID 列表。如不填，则从 `platforms` 中各取第一个可用账号。 |

### 2. 内容参数 (Content)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | 是 | 文章标题 |
| `content` | `string` | 是 | 文章正文（支持 HTML） |
| `cover_key` | `string` | 否 | **推荐**。通过能力“上传资源”得到的 OSS Key |
| `cover_url` | `string` | 否 | 封面图 URL（脚本将自动完成转存） |
| `tags` | `string[]` | 否 | 标签列表，各平台通用 |
| `pub_type` | `number` | 否 | `1`（公开，默认），`0`（草稿） |

### 3. 平台专有参数 (Platform Specifics)
可通过 `--[platform]_[param]` 形式传入：
- `qiehao_declaration`: 企鹅号创作声明 (0-9)
- `toutiao_is_first`: 头条号是否首发 (true/false)

## 任务执行建议 (Guidance)

1. **分步执行**：建议先用 [上传资源](./upload-resource.md) 处理图片，再调用本能力。
2. **账号获取**：建议先用 [查询账号](./query-accounts.md) 确认 ID。

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/publish-article.ts`
- **核心逻辑**:
  1. **配置归一化**: 自动匹配不同平台的 Form DTO。
  2. **内容云端存证**: 调用 `/api/storages/articles` 获取内容 ID。
  3. **任务原子提交**: 聚合所有平台和账号，调用 `/api/taskSets/v2` 发起单一任务集。

## 调用示例 (AI 指令)
"发布文章到百家号和企鹅号。标题：xxx，内容：...，账号 ID：bjh_123, qh_456"
