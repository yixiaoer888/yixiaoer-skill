# 视频号图文发布参数 (WeChat Video Channel Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 标题 | - |
| `description` | `string` | 否 | 图文内容描述 | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐信息 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |
| `collection` | `Object` | 否 | 合集信息 | - |
| `pubType` | `number` | 否 | 发布类型: 0-草稿, 1-直接发布 | 1 |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["视频号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号图文标题",
          "description": "这是描述内容",
          "images": [
            { "key": "img_key_1", "size": 1024, "width": 1080, "height": 1080 }
          ],
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `ShiPingHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/shipinghao.dto.ts`
