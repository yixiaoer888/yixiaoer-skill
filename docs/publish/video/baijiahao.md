# 百家号 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 是 | 视频标签 | - |
| declaration | number | 是 | 声明：0-不声明，1-内容由AI生成 | 0 |
| location | object | 否 | 位置信息，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| collection | object | 否 | 合集信息 | - |
| activity | object | 否 | 征文活动信息 | - |

## 2. 复杂对象结构

### PlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| id | string | 是 | ID |
| text | string | 是 | 文本 |
| raw | object | 是 | 平台原始数据 |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Baijiahao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BAIJIAHAO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号视频标题示例",
          "description": "这是关于百家号视频发布的合规描述，包含重要信息。",
          "tags": ["科技", "互联网"],
          "declaration": 0,
          "location": {
            "id": "123",
            "text": "北京市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```
