# 哔哩哔哩视频发布参数 (Bilibili)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题字段 | - |
| `description` | `string` | **是** | 描述字段 | - |
| `tags` | `string[]` | **是** | 标签字段 | - |
| `category` | `Array` | **是** | 分类字段 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | **是** | 声明: 0-不申明, 1-AI合成, 2-危险行为, 3-仅供娱乐, 4-引人不适, 5-适度消费, 6-个人观点 | - |
| `type` | `number` | **是** | 类型: 1-自制, 2-转载 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `contentSourceUrl` | `string` | 否 | 当 `type` 为 2 (转载) 时必填 | - |
| `collection` | `object` | 否 | 合集信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Bilibili"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BILIBILI_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "B站视频标题",
          "description": "B站视频描述",
          "tags": ["B站", "生活"],
          "category": [{"id": "1", "name": "生活"}],
          "declaration": 0,
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `BilibiliVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/bilibili.dto.ts`
