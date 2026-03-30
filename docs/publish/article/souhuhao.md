# 搜狐号文章发布参数 (SouHuHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 文章封面图列表 (`OldCover[]`) | - |
| `desc` | `string` | **是** | 文章描述 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["SouHuHao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐网发布测试",
          "content": "<p>正文内容...</p>",
          "desc": "这是一篇搜狐文章的描述",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `SouHuHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/souhuhao.dto.ts`
