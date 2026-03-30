# 车家号视频发布参数 (Chejiahao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `type` | `number` | **是** | 创作类型: 1-原创, 3-首发, 13-原创首发 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Chejiahao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "CHEJIA_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "汽车之家视频",
          "description": "内容描述",
          "type": 13
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `ChejiahaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/chejiahao.dto.ts`
