# 哔哩哔哩视频发布参数 (BiLiBiLi Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


在本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **是** | 视频标题 (最多 80 字符) | - |
| `description` | `string` | 否 | 视频描述 (最多 2000 字符) | - |
| `tags` | `string[]` | **是** | 视频标签 (1-10 个) | - |
| `category` | `Array` | **是** | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| `declaration` | `number` | 否 | 创作者申明：0-不申明, 1-AI合成, 2-危险行为, 3-仅供娱乐, 4-引人不适, 5-理性适度消费, 6-个人观点 | 0 |
| `createType` | `number` | **是** | 类型：1-自制, 2-转载 | 1 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒) | - |
| `contentSourceUrl` | `string` | 否 | 原文 URL 链接 (当 `createType` 为 2 时必填) | - |
| `collection` | `object` | 否 | 合集信息，使用 `Collection` 结构 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Bilibili"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BILIBILI_ACC_ID",
        "video": {
          "key": "video_oss_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "coverKey": "video_cover_key",
        "cover": { "key": "video_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "B站视频发布标题",
          "description": "这是关于B站视频发布的详细描述。",
          "tags": ["生活", "摄影"],
          "category": [
            { "yixiaoerId": "cat_001", "yixiaoerName": "生活", "raw": {} }
          ],
          "declaration": 0,
          "createType": 1,
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构说明

### 3.1 CascadingPlatformDataItem (多级分类)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | ID |
| `yixiaoerName` | `string` | **是** | 文本内容 |
| `raw` | `object` | **是** | 平台原始数据 (透传) |

### 3.2 Collection (合集)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 合集 ID |
| `yixiaoerName` | `string` | **是** | 合集名称 |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| `category` | `categories` | [获取发布分类](../../get-publish-categories.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
