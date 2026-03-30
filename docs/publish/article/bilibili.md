# 哔哩哔哩专栏发布参数 (BiLiBiLi Article)

本平台专栏发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 专栏标题 | - |
| `content` | `string` | **是** | 专栏 HTML 正文 | - |
| `covers` | `Array` | **是** | 专栏封面图列表 (`OldCover[]`) | - |
| `tags` | `Array` | 否 | 标签名称列表 (`string[]`) | - |
| `type` | `number` | 否 | 创作类型 (1:自制, 2:转载) | `1` |
| `category` | `Array` | 否 | 分类列表 (`Category[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["BiLiBiLi"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "B站专栏发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ],
          "tags": ["测试", "B站"],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `BiLiBiLiArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/bilibili.dto.ts`
