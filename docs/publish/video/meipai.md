# 美拍 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 否 | 视频描述 | - |
| category | object[] | 否 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
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
  "platforms": ["Meipai"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "MEIPAI_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "美拍视频标题示例",
          "description": "这是美拍平台的一段精彩短视频内容描述。",
          "category": [
            {
              "id": "1",
              "text": "生活",
              "raw": {}
            }
          ]
        }
      }
    ]
  }
}
```
