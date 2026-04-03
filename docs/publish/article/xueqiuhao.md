# 雪球号文章发布参数 (XueQiuHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 (长度: 9-100) | - |
| `content` | `string` | **是** | 文章内容 (HTML 格式) | - |
| `covers` | `Array` | 否 | 文章封面列表 (`OldCover[]`) | - |
| `visibleType` | `number` | **是** | 可见类型: 0-公开, 1-私密 | 0 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `declaration` | `number` | 否 | 内容申明: 0-不申明, 1-包含 AI 内容 | 0 |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |

## 2. 复杂对象结构说明

### 2.1 OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

## 3. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["雪球号"],
  "publishArgs": {
    "content": "<h1>雪球号财经文章</h1><p>内容正文...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_xq_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "这是财经文章标题示例",
          "content": "<h1>雪球号财经文章</h1><p>内容正文...</p>",
          "visibleType": 0,
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
