# 快手开放平台 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 快手标题 | - |
| description | string | 否 | 快手描述 | - |
| visibleType | number | 是 | 可见类型：0-公开, 1-私密, 3-好友可见 | 0 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["KuaishouOpen"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "KS_OPEN_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "快手开放平台视频标题",
          "description": "通过快手开放平台发布的精彩内容内容描述。",
          "visibleType": 0
        }
      }
    ]
  }
}
```
