# 知乎文章发布参数 (ZhiHu Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `topics` | `Array` | 否 | 话题列表 (`Category[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `declaration`| `number` | 否 | 创作申明: 0-无申明, 1-剧透, 2-医疗建议, 3-虚构创作, 4-理财内容, 5-AI辅助 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["知乎"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_zh_001",
        "coverKey": "article_cover_key",
        "cover": { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题",
          "covers": [
            { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "topics": [
            { "yixiaoerId": "topic_001", "yixiaoerName": "AI", "yixiaoerImageUrl": "", "yixiaoerDesc": "", "viewNum": "0", "raw": {} }
          ]
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

### 3.2 Category (话题/分类对象)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | ID |
| `yixiaoerName` | `string` | 名称 |
| `yixiaoerImageUrl` | `string` | 图片 URL |
| `yixiaoerDesc` | `string` | 描述 |
| `viewNum` | `string` | 浏览量 |
| `raw` | `Object` | 原始对象 |

## 4. DTO 参考
- 后端类: `ZhiHuArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`
