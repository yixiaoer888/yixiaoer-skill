# 头条号视频发布参数 (TouTiaoHao Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 (1-80 字符) | - |
| `description` | `string` | **是** | 视频描述 (1-400 字符) | - |
| `tags` | `string[]` | **是** | 视频标签 (1-5 个) | - |
| `declaration` | `number` | 否 | 创作者申明：1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎故事经历, 7-投资观点仅供参考, 8-健康医疗分享仅供参考 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳（单位：秒） | - |
| `visibleType` | `number` | **是** | 可见性: 0-公开, 1-私密 | - |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_tt_video_001",
        "video": {
          "key": "video_resource_key",
          "size": 10240000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "头条视频发布标题",
          "description": "视频描述内容...",
          "tags": ["生活", "摄影"],
          "declaration": 1,
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
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
