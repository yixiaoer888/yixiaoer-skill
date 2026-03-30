# 企鹅号文章发布参数 (QiEHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 文章封面图列表 (`OldCover[]`) | - |
| `tags` | `Array` | **是** | 文章标签列表 (`string[]`) | - |
| `declaration` | `number` | 否 | 创作声明 (0:暂不声明, 1:AI生成, 2:个人观点, 3:剧情演绎, 7:AI辅助, 8:健康, 9:危险行为) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["QiEHao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "腾讯发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ],
          "tags": ["测试", "腾讯"]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `QiEHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/qiehao.dto.ts`
