# CSDN文章发布参数 (CSDN Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `desc` | `string` | **是** | 文章摘要/描述 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `tags` | `Array` | **是** | 文章标签列表 (`string[]`) | - |
| `type` | `number` | **是** | 创作类型 (1:原创, 2:转载, 4:翻译) | - |
| `contentSourceUrl` | `string` | 否 | 原文链接 (转载/翻译时必填) | - |
| `declaration` | `number` | 否 | 声明 (0:无, 1:AI辅助, 2:网络来源, 3:个人观点) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["CSDN"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "CSDN文章发布测试",
          "content": "<p>正文内容...</p>",
          "desc": "这是一篇技术文章的摘要",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ],
          "tags": ["JavaScript", "Node.js"],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `CSDNArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/csdn.dto.ts`
