# 爱奇艺视频发布参数 (Aiqiyi)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | 否 | 视频标签 | - |
| `category` | `Array` | **是** | 分类信息 (`CascadingPlatformDataItem[]`) | - |
| `type` | `number` | **是** | 原创类型: 1-原创, 2-转载 | - |
| `declaration` | `number` | 否 | 声明: 0-无需申明, 1-AI生成, 2-虚构演绎, 3-取材网络 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Aiqiyi"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "IQIYI_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "爱奇艺标题",
          "description": "内容描述",
          "category": [{"id": "1", "name": "原创"}],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `AiQiYiVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/aiqiyi.dto.ts`
