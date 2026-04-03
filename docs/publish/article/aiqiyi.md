# 爱奇艺文章发布参数 (AiQiYi Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 (长度: 1-50 字符) | - |
| `pubType` | `number` | **是** | 发布类型：0-草稿，1-直接发布 | 1 |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["爱奇艺"],
  "publishArgs": {
    "content": "<h1>这是文章标题</h1><p>这是正文内容，支持 HTML 格式。</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_aqy_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题",
          "pubType": 1
        }
      }
    ]
  }
}
```

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `covers.key` | `upload` | [资源上传](../../upload-resource.md) |

