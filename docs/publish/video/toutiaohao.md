# 头条号视频发布参数 (Toutiaohao)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 视频标题 | - |
| `description` | `string` | **是** | 视频描述 | - |
| `tags` | `string[]` | **是** | 视频标签 | - |
| `declaration` | `number` | 否 | 视频原创类型: 1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎, 7-投资观点, 8-健康医疗 | - |
| `scheduledTime` | `number` | 否 | 视频定时发布时间戳 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Toutiaohao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TOUTIAO_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "视频标题",
          "description": "视频描述内容",
          "tags": ["标签1", "标签2"],
          "declaration": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `TouTiaoHaoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`
