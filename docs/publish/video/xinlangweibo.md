# 新浪微博视频发布参数 (Xinlangweibo)

本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 微博视频标题 | - |
| `description` | `string` | **是** | 微博正文/描述 | - |
| `type` | `number` | 否 | 类型: 1-原创, 2-转载, 3-二次创作 | `1` |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `collection` | `object` | 否 | 合集字段 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "video",
  "platforms": ["Xinlangweibo"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WEIBO_ACC_ID",
        "video": { "key": "v_key", "size": 1024, "width": 720, "height": 1280 },
        "contentPublishForm": {
          "formType": "task",
          "title": "微博视频",
          "description": "这是微博正文 #话题",
          "type": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `XinLangWeiBoVideoForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/xinlangweibo.dto.ts`
