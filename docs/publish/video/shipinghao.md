# 视频号视频发布参数 (WeiXin ShiPingHao Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


在本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | 否 | 视频标题 | - |
| `short_title` | `string` | 否 | 视频短标题 | - |
| `description` | `string` | 否 | 视频描述，支持 HTML 格式和 `@` 好友/话题标签 | - |
| `horizontalCover` | `object` | 否 | 视频横板封面，使用 `OldCover` 结构 | - |
| `createType` | `number` | **是** | 创建类型：1-草稿，2-直接发布 | 2 |
| `pubType` | `number` | **是** | 发布类型：0-草稿，1-直接发布 | 1 |
| `location` | `object` | 否 | 视频位置，使用 `PlatformDataItem` 结构 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒) | - |
| `shoppingCart` | `object` | 否 | 关联商品信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |
| `collection` | `object` | 否 | 合集信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |
| `activity` | `object` | 否 | 活动信息 (`yixiaoerId`, `yixiaoerName`, `raw`) | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Shipinghao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "video": {
          "key": "video_oss_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "coverKey": "video_cover_key",
        "cover": { "key": "video_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号视频标题",
          "description": "<p>这是视频号内容的精彩描述 #生活 #记录</p>",
          "createType": 2,
          "pubType": 1,
          "location": {
            "yixiaoerId": "sph_loc_1",
            "yixiaoerName": "广州市",
            "raw": {}
          }
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

### 3.2 PlatformDataItem (位置/商品/合集/活动)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 统一 ID |
| `yixiaoerName` | `string` | **是** | 显示名称 |
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

## 相关接口

| 目标数据 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `activity` | `activities` | [获取活动列表](../../get-publish-activities.md) |
| `shoppingCart` | `goods` | [获取商品列表](../../get-goods.md) |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
