# 百家号文章发布参数 (BaiJiaHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 文章封面图列表 (`OldCover[]`) | - |
| `category` | `Array` | **是** | 文章分类列表 (`Category[]`) | - |
| `declaration` | `number` | **是** | 创作声明 (0:不声明, 1:内容由AI生成) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |
| `activity` | `Object` | 否 | 征文活动数据 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["BaiJiaHao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "百度发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ],
          "category": [
            { "yixiaoerId": "123", "yixiaoerName": "科技" }
          ],
          "declaration": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `BaiJiaHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`
