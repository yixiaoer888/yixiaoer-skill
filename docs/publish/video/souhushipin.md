# 搜狐视频 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 否 | 视频描述 | - |
| tags | string[] | 否 | 视频标签 | - |
| declaration | number | 否 | 搜狐视频申明：0-无需申明, 3-AI生成, 4-虚构演绎, 5-AI数字人生成 | 0 |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Souhushipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SOUHU_VIDEO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐视频标题示例",
          "description": "这是大搜狐平台的精彩视频分享内容。",
          "tags": ["生活", "见闻"],
          "declaration": 0
        }
      }
    ]
  }
}
```
