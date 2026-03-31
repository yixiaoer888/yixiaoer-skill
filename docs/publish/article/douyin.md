# 抖音文章发布参数 (Douyin Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `description` | `string` | 否 | 文章描述或摘要 | - |
| `covers` | `Array` | **是** | 封面图 OSS 列表 (`OldCover[]`) | - |
| `headImage` | `Object` | 否 | 文章头图 (`OldCover`) | - |
| `music` | `Object` | 否 | 平台音乐背景 (`PlatformDataItem`) | - |
| `topics` | `Array` | 否 | 话题标签列表 (`Category[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `visibleType` | `number` | 否 | 可见性: 0-公开, 1-私密, 3-好友 | `0` |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["抖音"],
  "publishArgs": {
    "content": "<h1>抖音文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "coverKey": "cover_oss_key",
        "cover": { "key": "cover_oss_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是一篇抖音文章",
          "covers": [
            { "key": "cover_oss_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "visibleType": 0
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

### 3.2 Category (分类/话题对象)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | ID |
| `yixiaoerName` | `string` | 名称 |
| `raw` | `Object` | 原始对象 |

## 4. DTO 参考
- 后端类: `DouyinArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`
