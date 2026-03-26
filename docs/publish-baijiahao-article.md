# 发布百家号文章 (Publish Baijiahao Article)

将文章内容发布到百家号平台。建议配合“上传资源”能力完成封面的预处理。

## 场景描述 (Usage)

- "帮我把这篇关于 AI 的文章发布到我的百家号。"
- "使用我刚刚上传好的封面 (Key: xxx) 发布文章。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | 是 | 文章标题 |
| `content` | `string` | 是 | 文章正文（支持 HTML） |
| `cover_key` | `string` | 否 | **推荐**。通过能力“上传资源”得到的 OSS Key |
| `cover_url` | `string` | 否 | 封面图 URL（如传入，脚本将内部尝试自动转存上传） |
| `account_ids` | `string[]` | 否 | 目标账号 ID 列表。 |
| `coverType` | `'single'\|'triple'` | 否 | 封面类型：`single`（默认），`triple` |
| `declaration` | `number` | 否 | 内容声明：`0`（无），`1`（AI 生成），`2`（原创，默认） |
| `pubType` | `number` | 否 | 发布类型：`1`（公开，默认），`0`（草稿） |

## 任务执行建议 (Guidance)

**最佳性能方案**：
1. 先调用 [上传资源](./upload-resource.md) 能力。
2. 将得到的 `key` 填入本能力的 `cover_key` 参数中。

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/publish-baijiahao-article.ts`
- **核心逻辑**:
  1. **资源处理**: 优先使用 `cover_key`；如只有 `cover_url` 则内部同步上传。
  2. **内容存证**: 提交 HTML 内容至 `/api/storages/articles` 获取 `publishContentId`。
  3. **任务创建**: 调用 `/api/taskSets/v2` 发起百家号发布任务。

## 调用示例 (AI 指令)
"帮我发布百家号。标题是：xxx，内容是：xxx，封面 Key 是：cloud-publish/2026/03/xxx.jpg"
