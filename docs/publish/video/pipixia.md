# 皮皮虾视频发布参数 (Pipixia)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | 否 | 视频描述 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Pipixia"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PIPIXIA_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "description": "皮皮虾视频描述 #话题"
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `PiPiXiaVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/pipixia.dto.ts`
