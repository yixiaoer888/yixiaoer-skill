# AcFun文章发布参数 (AcFun Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `desc` | `string` | 否 | 文章描述 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `tags` | `Array` | **是** | 标签列表 (`string[]`) | - |
| `category` | `Array` | **是** | 分类列表 (`Category[]`) | - |
| `type` | `number` | **是** | 创作类型 (1:原创, 0:非原创) | - |
| `contentSourceUrl` | `string` | 否 | 原文链接 (非原创时选填) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["AcFun"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "A站文章发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [ { "key": "cover_key", "size": 100, "width": 800, "height": 600 } ],
          "tags": ["测试", "A站"],
          "category": [ { "yixiaoerId": "cat_id", "yixiaoerName": "生活" } ],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `AcFunArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/acfun.dto.ts`
