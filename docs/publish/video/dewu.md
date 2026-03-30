# 得物视频发布参数 (Dewu)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | **是** | 声明: 0-无, 1-AI生成, 2-不含营销, 3-专业运动, 4-剧情演绎 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Dewu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DEWU_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "得物视频",
          "description": "内容描述 #潮鞋",
          "declaration": 2
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DeWuVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/dewu.dto.ts`
