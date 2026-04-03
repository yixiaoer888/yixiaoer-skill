# 大鱼号视频发布参数 (DaYuHao Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


在本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值为 `task` | `task` |
| `title` | `string` | **是** | 视频标题 (最多 50 字符) | - |
| `description` | `string` | **是** | 视频描述 (最多 1000 字符) | - |
| `tags` | `string[]` | **是** | 视频标签 (1-6 个) | - |
| `category` | `Array` | 否 | 视频分类，使用 `CascadingPlatformDataItem[]` | - |
| `createType` | `number` | 否 | 创作类型: 0-非原创, 1-原创 | 0 |
| `declaration` | `number` | 否 | 声明字段: 0-无需申明, 3-虚构演绎, 4-AI 生成 | 0 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["大鱼号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_dy_vid_001",
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
          "title": "大鱼号视频发布标题",
          "description": "这是关于此视频的详细描述内容。",
          "tags": ["生活", "摄影"],
          "category": [
            { "yixiaoerId": "cat_001", "yixiaoerName": "生活", "raw": {} }
          ],
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

### 3.1 CascadingPlatformDataItem (分类对象)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 级联 ID |
| `yixiaoerName` | `string` | **是** | 级联显示的名称 |
| `children` | `Array` | 否 | 子级对象列表 |
| `raw` | `object` | **是** | 平台原始对象 (必须透传) |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| `category` | `categories` | [获取发布分类](../../get-publish-categories.md) |
