# AcFun 视频发布参数 (AcFun)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | 否 | 视频描述 | - |
| `tags` | `string[]` | 否 | 视频标签 | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `type` | `number` | **是** | 内容类型: 1-原创, 0-非原创 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["AcFun"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ACFUN_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "A站视频",
          "description": "内容描述",
          "category": [{"id": "1", "name": "科技"}],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `AcFunVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/acfun.dto.ts`
