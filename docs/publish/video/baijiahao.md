# 百家号视频发布参数 (Baijiahao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题字段 | - |
| `description` | `string` | **是** | 描述字段 | - |
| `tags` | `string[]` | **是** | 标签字段 | - |
| `declaration` | `number` | **是** | 声明: 0-不声明, 1-内容由AI生成 | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `collection` | `object` | 否 | 合集信息 | - |
| `activity` | `object` | 否 | 活动信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Baijiahao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BJH_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号视频标题",
          "description": "内容描述",
          "tags": ["标签1"],
          "declaration": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `BaiJiaHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`
