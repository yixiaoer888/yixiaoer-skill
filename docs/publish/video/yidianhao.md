# 一点号视频发布参数 (Yidianhao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题字段 | - |
| `description` | `string` | **是** | 描述字段 | - |
| `tags` | `string[]` | **是** | 标签字段 | - |
| `category` | `Array` | **是** | 分类信息 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | 否 | 声明: 3-取材网络, 4-AI生成, 5-虚构情节 | - |
| `type` | `number` | **是** | 原创类型: 0-非原创, 1-原创 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Yidianhao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YIDIAN_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "一点号视频",
          "description": "描述内容",
          "tags": ["一点"],
          "category": [{"id": "1", "name": "生活"}],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `YiDianHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/yidianhao.dto.ts`
