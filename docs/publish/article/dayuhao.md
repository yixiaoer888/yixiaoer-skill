# 大鱼号文章发布参数 (DaYuHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 横版封面图列表 (`OldCover[]`) | - |
| `verticalCovers` | `Array` | **是** | 竖版封面图列表 (`OldCover[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["DaYuHao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "大鱼号发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "h_cover", "size": 100, "width": 800, "height": 600 }
          ],
          "verticalCovers": [
            { "key": "v_cover", "size": 100, "width": 600, "height": 800 }
          ]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DaYuHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/dayuhao.dto.ts`
