# 豆瓣日记发布参数 (DouBan Article)

本平台日记发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 日记标题 | - |
| `content` | `string` | **是** | 日记 HTML 正文 | - |
| `type` | `number` | **是** | 创作类型 (0:不声明, 1:声明原创) | - |
| `tags` | `Array` | 否 | 标签列表 (`string[]`) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["DouBan"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "豆瓣日记发布测试",
          "content": "<p>正文内容...</p>",
          "type": 1,
          "tags": ["测试", "豆瓣"]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DouBanArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/douban.dto.ts`
