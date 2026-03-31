# 知乎 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| topics | object[] | 否 | 知乎话题列表，使用 `Category[]` 结构 | - |
| category | object[] | 是 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| declaration | number | 否 | 视频创作申明：0-不声明, 2-图片/视频由AI生成 | 0 |
| type | number | 否 | 视频类型：1-原创, 2-非原创 | 1 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |

## 2. 复杂对象结构

### Category
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| yixiaoerId | string | 是 | 蚁小二ID |
| yixiaoerName | string | 是 | 蚁小二名称 |
| raw | object | 是 | 平台原始数据 |

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
  "platforms": ["Zhihu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZHIHU_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎精选视频标题",
          "description": "这是关于知乎深度内容的视频分享描述内容。",
          "topics": [
            {
              "yixiaoerId": "123",
              "yixiaoerName": "人工智能",
              "raw": {}
            }
          ],
          "category": [
            {
              "id": "1",
              "text": "科技",
              "raw": {}
            }
          ],
          "type": 1,
          "declaration": 0
        }
      }
    ]
  }
}
```
