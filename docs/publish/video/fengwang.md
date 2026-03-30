# 蜂网视频发布参数 (Fengwang)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | **是** | 视频标签 | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Fengwang"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "FENGWANG_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "蜂网视频",
          "description": "描述内容",
          "tags": ["蜂网"],
          "category": [{"id": "1", "name": "科技"}]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `FengWangVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/fengwang.dto.ts`
