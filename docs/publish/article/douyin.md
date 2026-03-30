# 抖音文章发布参数 (Douyin Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `description` | `string` | 否 | 文章描述或摘要 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 封面图 OSS 列表 (`OldCover[]`) | - |
| `headImage` | `Object` | 否 | 文章头图 | - |
| `music` | `Object` | 否 | 平台音乐背景 | - |
| `topics` | `Array` | 否 | 话题标签列表 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 | - |
| `visibleType` | `number` | 否 | 0-公开, 1-私密, 3-好友 | `0` |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "这是一篇抖音文章",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "cover_oss_key", "size": 100, "width": 800, "height": 600 }
          ]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DouyinArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`
