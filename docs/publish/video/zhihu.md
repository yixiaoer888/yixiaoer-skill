# 知乎视频发布参数 (Zhihu)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数 definition

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `topics` | `Array` | 否 | 话题列表 (`Category[]`) | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | 否 | 声明: 0-不声明, 2-AI生成 | - |
| `type` | `number` | 否 | 创作类型: 1-原创, 2-非原创 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Zhihu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎视频标题",
          "description": "视频内容描述",
          "category": [{"id": "1", "name": "科技"}],
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `ZhiHuVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`
