# 发布公众号 (WeChat Engine)

微信公众号专用的原子分发接口，支持单图文和多图文的推送。

## 场景描述 (Usage)
- "将这篇深度报道推送到我的微信公众号。"
- "使用封面 (Key: xxx) 推送公众号消息。"

## 参数定义 (Parameters)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | 是 | 文章标题 |
| `content` | `string` | 是 | 文章正文 (HTML) |
| `cover_key` | `string` | 是 | 公众号封面的 OSS Key |
| `author` | `string` | 否 | 作者 |
| `digest` | `string` | 否 | 摘要 |
| `account_ids` | `string[]` | 否 | 目标公众号 ID 列表 |

## 脚本逻辑 (Backend)
- **脚本路径**: `../scripts/publish-wechat.ts`
- **核心逻辑**: 向 `/api/storages/articles` 存证 -> 构建微信专有 DTO -> 提交任务集。
