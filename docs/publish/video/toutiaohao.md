# 头条号 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 是 | 视频标签 | - |
| declaration | number | 否 | 创作者申明：1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎故事经历, 7-投资观点仅供参考, 8-健康医疗分享仅供参考 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |

## 2. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Toutiaohao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TOUTIAO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "头条号视频标题示例",
          "description": "这是头条号平台的精彩视频描述内容。",
          "tags": ["军事", "历史"],
          "declaration": 1
        }
      }
    ]
  }
}
```
