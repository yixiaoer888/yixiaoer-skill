# 腾讯视频发布参数 (TencentVideo)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `tags` | `string[]` | 否 | 视频标签 | - |
| `declaration` | `number` | **是** | 腾讯视频申明: 1-AI生成, 2-剧情演绎, 3-取材网络, 4-个人观点 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["TencentVideo"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TX_VIDEO_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "腾讯视频标题",
          "declaration": 2
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `TengXunShiPinVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/tengxunshipin.dto.ts`
