# 企鹅号视频发布参数 (Qiehao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | **是** | 视频标签 | - |
| `category` | `Array` | **是** | 分类信息 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | 否 | 声明: 0-暂不申明, 1-AI生成, 2-个人观点, 3-剧情演绎, 4-取材网络, 5-旧闻 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Qiehao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "QQ_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "视频标题",
          "description": "内容描述",
          "tags": ["腾讯"],
          "category": [{"id": "1", "name": "科技"}],
          "declaration": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `QiEHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/qiehao.dto.ts`
