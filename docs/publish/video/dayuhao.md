# 大鱼号视频发布参数 (Dayuhao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | **是** | 视频标签 | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `type` | `number` | **是** | 创作类型: 1-原创, 0-非原创 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Dayuhao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DAYU_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "大鱼号视频",
          "description": "内容描述",
          "tags": ["大鱼"],
          "category": [{"id": "1", "name": "科技"}],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DaYuHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/dayuhao.dto.ts`
