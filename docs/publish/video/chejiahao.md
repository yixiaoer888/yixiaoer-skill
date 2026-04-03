# 车家号 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| type | number | 是 | 创作类型：1-原创, 3-首发, 13-原创首发 | 1 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Chejiahao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "CHEJIAHAO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "车家号视频标题示例",
          "description": "这是关于汽车评测的视频描述内容。",
          "type": 1
        }
      }
    ]
  }
}
```
