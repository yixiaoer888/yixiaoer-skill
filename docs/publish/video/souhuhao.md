# 搜狐号 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 是 | 视频标签 | - |
| category | object[] | 是 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| declaration | number | 是 | 原创类型信息来源：0-无特别声明, 1-引用申明, 2-自行拍摄, 3-包含AI创作内容, 4-包含虚构创作 | 0 |
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
  "platforms": ["Souhuhao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SOUHUHAO_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐号视频标题示例",
          "description": "这是关于搜狐号视频发布的精彩描述信息。",
          "tags": ["娱乐", "电影"],
          "category": [
            {
              "id": "1",
              "text": "娱乐",
              "raw": {}
            }
          ],
          "declaration": 2
        }
      }
    ]
  }
}
```
