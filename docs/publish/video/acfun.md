# AcFun视频发布参数 (AcFun Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


在本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | 否 | 视频描述 (建议 HTML 格式用 `<p>` 包裹) | - |
| `tags` | `string[]` | 否 | 视频标签 (最多 6 个) | - |
| `category` | `Array` | **是** | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| `type` | `number` | **是** | 内容类型: 1-原创, 0-非原创 | 0 |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["AcFun"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "AC_ACC_ID",
        "video": {
          "key": "video_oss_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "AcFun 视频发布标题",
          "description": "<p>这是 AcFun 视频的描述内容</p>",
          "tags": ["生活", "美食"],
          "category": [
            { "yixiaoerId": "cat_001", "yixiaoerName": "生活", "raw": {} }
          ],
          "type": 1
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
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| `category` | `categories` | [获取发布分类](../../get-publish-categories.md) |
