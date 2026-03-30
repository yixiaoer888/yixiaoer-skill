# 知乎文章发布参数 (ZhiHu Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 文章封面图列表 (`OldCover[]`) | - |
| `topics` | `Array` | 否 | 话题列表 (`Category[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |
| `declaration` | `number` | 否 | 创作声明 (0:无, 1:剧透, 2:医疗建议, 3:虚构, 4:理财, 5:AI辅助) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["ZhiHu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎文章发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ],
          "topics": [
            { "yixiaoerId": "topic_id", "yixiaoerName": "互联网" }
          ]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `ZhiHuArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`
