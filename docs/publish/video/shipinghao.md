# 视频号视频发布参数 (Shipinghao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 视频标题 | - |
| `short_title` | `string` | 否 | 视频短标题 | - |
| `description` | `string` | 否 | 视频描述 (支持 HTML 及话题) | - |
| `location` | `Object` | 否 | 视频位置 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `type` | `number` | **是** | 原创类型: 1-非原创, 2-原创 | - |
| `shoppingCart` | `object` | 否 | 关联商品 | - |
| `horizontalCover` | `Object` | 否 | 横版封面 (`OldCover`) | - |
| `collection` | `object` | 否 | 合集信息 | - |
| `activity` | `object` | 否 | 活动信息 | - |
| `pubType` | `number` | 否 | 发布类型: 0-草稿, 1-直接发布 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Shipinghao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>视频号描述 #话题</p>",
          "type": 2,
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `ShiPingHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/shipinghao.dto.ts`
