# 知乎文章发布参数 (ZhiHu Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 (9-100 字符) | - |
| `content` | `string` | **是** | 文章内容 (9-10000 字符) | - |
| `covers` | `Array` | 否 | 文章封面列表 (`OldCover[]`) | - |
| `topics` | `Array` | 否 | 话题列表 (`Category[]`, 最多3个) | - |
| `declaration` | `number` | 否 | 创作申明: 0-无申明, 1-剧透, 2-医疗建议, 3-虚构创作, 4-理财内容, 5-AI辅助 | 0 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["Zhihu"],
  "publishArgs": {
    "content": "<p>这是知乎文章的正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎深度文章标题解析",
          "content": "<p>这是知乎文章的正文内容...</p>",
          "covers": [
            { "key": "cover_key_1", "size": 1024, "width": 800, "height": 600 }
          ],
          "topics": [
            { "yixiaoerId": "123", "yixiaoerName": "科技", "raw": {} }
          ],
          "declaration": 0,
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构说明

### 3.1 OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

### 3.2 Category (用于话题)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 话题 ID |
| `yixiaoerName` | `string` | **是** | 话题名称 |
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `topics` | `challenges` | [获取话题/挑战](../../get-challenges.md) |
| `covers.key` | `upload` | [资源上传](../../upload-resource.md) |
