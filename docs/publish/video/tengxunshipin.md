# 腾讯视频 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| tags | string[] | 否 | 视频标签 | - |
| declaration | number | 否 | 腾讯视频申明：1-内容由AI生成, 2-剧情演绎仅供娱乐, 3-取材网络谨慎甄别, 4-个人观点仅供参考 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Tengxunshipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TX_VIDEO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "腾讯视频标题示例",
          "tags": ["影视", "评论"],
          "declaration": 4
        }
      }
    ]
  }
}
```
