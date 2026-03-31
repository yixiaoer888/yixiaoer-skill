# 爱奇艺文章发布参数 (AiQiYi Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["爱奇艺"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_aqy_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题"
        }
      }
    ]
  }
}
```

## 4. DTO 参考
- 后端类: `AiQiYiArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/aiqiyi.dto.ts`
