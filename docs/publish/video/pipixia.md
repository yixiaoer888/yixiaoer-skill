# 皮皮虾 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| description | string | 否 | 视频描述 | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Pipixia"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PIPIXIA_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "description": "皮皮虾，我们走！ #搞笑 #段子"
        }
      }
    ]
  }
}
```
