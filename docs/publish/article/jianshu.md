# 简书文章发布参数 (JianShu Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["JianShu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "简书发布测试",
          "content": "<p>正文内容...</p>"
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `JianShuArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/jianshu.dto.ts`
