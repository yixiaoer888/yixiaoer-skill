# 快传号文章发布参数 (KuaiChuanHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章内容 (HTML 格式) | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `tags` | `string[]` | 否 | 文章标签 | - |
| `type`| `number` | **是** | 创作类型: 0-不申明, 1-申明原创 | 1 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["快传号"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_kch_001",
        "coverKey": "article_cover_key",
        "cover": { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题",
          "content": "<h1>文章标题</h1><p>正文内容...</p>",
          "covers": [
            { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "type": 1,
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构说明

### 3.1 OldCover
| 字段名 | 类型 | 必填 |说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `covers.key` | `upload` | [资源上传](../../upload-resource.md) |
