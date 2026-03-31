# CSDN文章发布参数 (CSDN Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `desc` | `string` | **是** | 文章描述 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `tags` | `string[]` | **是** | 文章标签 | - |
| `type` | `number` | **是** | 创作类型: 1-原创, 2-转载, 4-翻译 | `1` |
| `contentSourceUrl` | `string` | 否 | 原文链接 (转载时必填) | - |
| `declaration` | `number` | 否 | 声明: 0-无, 1-AI辅助, 2-来源网络, 3-个人观点 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["CSDN"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_csdn_001",
        "coverKey": "article_cover_key",
        "cover": { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题",
          "desc": "这是摘要",
          "covers": [
            { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "tags": ["Java", "Spring"],
          "type": 1,
          "declaration": 0
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构

### 3.1 OldCover (封面对象)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `key` | `string` | OSS 资源 Key |
| `size` | `number` | 文件大小 (bytes) |
| `width` | `number` | 宽度 |
| `height` | `number` | 高度 |

## 4. DTO 参考
- 后端类: `CSDNArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/csdn.dto.ts`
