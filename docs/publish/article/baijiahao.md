# 百家号文章发布参数 (BaiJiaHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `category` | `Array` | **是** | 文章分类列表 (`Category[]`) | - |
| `declaration` | `number` | 否 | 内容声明: 0-不声明, 1-内容由 AI 生成 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `activity` | `Object` | 否 | 征文活动数据对象，使用 `Activity` 结构 | - |

## 2. 复杂对象结构说明

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

### Category (分类)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 分类 ID |
| `yixiaoerName` | `string` | **是** | 分类名称 |
| `raw` | `object` | 否 | 原始分类对象 |

### Activity (活动)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 活动 ID |
| `yixiaoerName` | `string` | **是** | 活动名称 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `category` | `categories` | [获取账号分类](../../get-publish-categories.md) |
| `activity` | `activities` | [获取征文活动](../../get-publish-activities.md) |

## 3. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["百家号"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_bjh_001",
        "coverKey": "article_cover_key",
        "cover": { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题",
          "covers": [
            { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "category": [
            { "yixiaoerId": "cat_001", "yixiaoerName": "文化", "yixiaoerImageUrl": "", "yixiaoerDesc": "文化类", "viewNum": "100", "raw": {} }
          ],
          "declaration": 0
        }
      }
    ]
  }
}
```

## 4. DTO 参考
- 后端类: `BaiJiaHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`
