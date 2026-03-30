# 搜狐号视频发布参数 (Souhuhao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | **是** | 视频标签 | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | **是** | 原创类型: 0-无特别声明, 1-引用申明, 2-自行拍摄, 3-AI创作, 4-虚构创作 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Souhuhao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SOHU_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐号视频",
          "description": "内容正文",
          "tags": ["搜狐"],
          "category": [{"id": "1", "name": "综合"}],
          "declaration": 2
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `SouHuHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/souhuhao.dto.ts`
