# CSDN文章发布参数 (CSDN Article)

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
| `desc` | `string` | **是** | 文章内容摘要 | - |
| `covers` | `OldCover[]` | **是** | 文章封面列表 | - |
| `tags` | `string[]` | **是** | 文章标签 (字符串数组) | - |
| `createType` | `number` | **是** | 创作类型: 1-原创, 2-转载, 4-翻译 | 1 |
| `contentSourceUrl` | `string` | 否 | 原文链接 (当 `createType` 为 2 时必填) | - |
| `declaration` | `number` | **否** | 声明: 0-无, 1-AI 辅助生成, 2-内容来源网络, 3-个人观点仅供参考 | 0 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | **否** | 定时发布时间 (Unix 时间戳，单位: 秒) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["CSDN"],
  "publishArgs": {
    "content": "<h1>CSDN 文章发布演示</h1><p>这是正文内容。</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_csdn_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "如何使用 CSDN 发布 API",
          "content": "<h1>CSDN 文章发布演示</h1><p>这是正文内容。</p>",
          "desc": "本文介绍了 CSDN 发布 API 的具体参数和使用方法。",
          "covers": [
            {
              "key": "csdn_cover_key_001",
              "size": 102400,
              "width": 800,
              "height": 600
            }
          ],
          "tags": ["开发者", "编程", "API"],
          "createType": 1,
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

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `covers.key` | `upload` | [资源上传](../../upload-resource.md) |
