# 发布百家号文章 (Publish Baijiahao Article)

将文章内容发布到百家号平台，支持自定义分类、封面及内容声明。

## 场景描述 (Usage)

- "帮我把这篇关于 AI 的文章发布到我的百家号。"
- "在百家号上发布这篇草稿，封面使用这张图片。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | 是 | 文章标题 |
| `content` | `string` | 是 | 文章正文（支持 HTML） |
| `cover_url` | `string` | 否 | 封面图 URL |
| `account_ids` | `string[]` | 否 | 目标账号 ID 列表。如果不传，则自动使用第一个百家号账号。 |
| `coverType` | `'single'\|'triple'` | 否 | 封面类型：`single`（单图，默认），`triple`（三图） |
| `declaration` | `number` | 否 | 内容声明：`0`（无），`1`（AI 生成），`2`（原创，默认） |
| `pubType` | `number` | 否 | 发布类型：`1`（公开，默认），`0`（草稿） |

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/publish-baijiahao-article.ts`
- **核心逻辑**:
  1. **资源上传**: 如果提供 `cover_url`，将其上传至蚁小二 OSS。
  2. **内容存证**: 提交 HTML 内容至 `/api/storages/articles` 获取 `publishContentId`。
  3. **任务创建**: 调用 `/api/taskSets/v2` 发起多平台（或单平台）发布任务。
- **调用示例**: `node publish-baijiahao-article.ts --title="我的第一篇文章" --content="<p>内容...</p>"`

## 输出结果 (Output)

- 成功：返回包含 `task_set_id` 的 JSON 对象。
- 失败：返回包含 `error` 详情的 JSON 对象。
