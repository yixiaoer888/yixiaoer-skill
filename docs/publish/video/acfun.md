# AcFun 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 否 | 视频描述 | - |
| tags | string[] | 否 | 视频标签 | - |
| category | object[] | 是 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| type | number | 是 | 内容类型：1-原创，0-非原创 | 0 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |

## 2. 复杂对象结构

### CascadingPlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| id | string | 是 | 选项ID |
| text | string | 是 | 选项文本 |
| children | object[] | 否 | 子级选项列表 (CascadingPlatformDataItem[]) |
| raw | object | 是 | 平台原始数据 |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["AcFun"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ACFUN_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "AcFun视频标题示例",
          "description": "这是AcFun视频的描述内容",
          "tags": ["生活", "美食"],
          "category": [
            {
              "id": "1",
              "text": "生活",
              "raw": {}
            }
          ],
          "type": 1
        }
      }
    ]
  }
}
```
