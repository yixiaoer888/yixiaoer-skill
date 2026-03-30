# 搜狐视频发布参数 (Souhushipin)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | 否 | 视频描述 | - |
| `tags` | `string[]` | 否 | 视频标签 | - |
| `declaration` | `number` | 否 | 声明: 0-无需申明, 3-AI生成, 4-虚构演绎, 5-AI数字人 | `0` |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Souhushipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SOHU_VIDEO_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐视频标题",
          "declaration": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `SouHuShiPinVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/souhushipin.dto.ts`
