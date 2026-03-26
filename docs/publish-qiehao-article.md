# 发布企鹅号文章 (Publish Qiehao Article)

将文章内容发布到企鹅号平台。建议配合“上传资源”能力完成封面的预处理。

## 场景描述 (Usage)

- "帮我把这篇关于科技发展的文章发布到我的企鹅号。"
- "使用封面 (Key: xxx) 发布企鹅号文章，标签设定为 AI, 科技。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | 是 | 文章标题 |
| `content` | `string` | 是 | 文章正文（支持 HTML） |
| `cover_key` | `string` | 否 | **推荐**。通过能力“上传资源”得到的 OSS Key |
| `cover_url` | `string` | 否 | 封面图 URL（如传入，脚本将内部尝试自动转存上传） |
| `tags` | `string[]` | 否 | 标签列表，例如 `["AI", "科技"]` |
| `declaration` | `number` | 否 | 内容声明：`0`（无，默认），`1`（原创），`2`（独家） |
| `pubType` | `number` | 否 | 发布类型：`1`（公开，默认），`0`（草稿） |
| `account_ids` | `string[]` | 否 | 目标账号 ID 列表。如不填则默认使用第一个可用企鹅号。 |

## 任务执行建议 (Guidance)

**最佳性能方案**：
1. 先调用 [上传资源](./upload-resource.md) 能力。
2. 将得到的 `key` 填入本能力的 `cover_key` 参数中。

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/publish-qiehao-article.ts`
- **核心逻辑**:
  1. **资源处理**: 优先使用 `cover_key`；如只有 `cover_url` 则内部同步上传。
  2. **内容存证**: 提交 HTML 内容至 `/api/storages/articles` 获取 `publishContentId`。
  3. **任务创建**: 调用 `/api/taskSets/v2` 发起企鹅号发布任务。

## 调用示例 (AI 指令)
"帮我发布企鹅号文章。标题是：未来科技，内容是：...，标签是：AI, 未来，封面 Key 是：cloud-publish/2026/03/xxx.jpg"
