# 百家号视频发布参数 (BaiJiaHao Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


在本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **是** | 视频标题 (1-30 字符) | - |
| `description` | `string` | **是** | 视频描述 (1-100 字符) | - |
| `tags` | `string[]` | **是** | 视频标签 (1-6 个) | - |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `declaration` | `number` | 否 | 声明: 0-不声明, 1-内容由 AI 生成 | 0 |
| `location` | `Object` | 否 | 位置信息，使用 `PlatformDataItem` 结构 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒) | - |
| `collection` | `Object` | 否 | 合集信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |
| `activity` | `Object` | 否 | 征文活动信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["百家号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_bjh_video_001",
        "video": {
          "key": "video_oss_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号视频标题",
          "description": "视频精彩描述内容...",
          "tags": ["科技", "未来"],
          "pubType": 1,
          "declaration": 0,
          "location": {
            "yixiaoerId": "loc_001",
            "yixiaoerName": "北京市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构说明

### 3.1 PlatformDataItem / Category
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 统一 ID |
| `yixiaoerName` | `string` | **是** | 显示名称 |
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `activity` | `activities` | [获取征文活动](../../get-publish-activities.md) |
