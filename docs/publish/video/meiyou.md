# 美柚视频发布参数 (Meiyou)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Meiyou"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "MEIYOU_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "美柚视频标题"
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `MeiYouVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/meiyou.dto.ts`
