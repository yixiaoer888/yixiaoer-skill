# 哔哩哔哩-Open 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 是 | 视频标签 | - |
| category | object[] | 是 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| declaration | number | 是 | 创作者申明：0-不申明, 1-AI合成, 2-危险行为, 3-仅供娱乐, 4-引人不适, 5-理性适度消费, 6-个人观点 | 0 |
| type | number | 是 | 类型：1-自制, 2-转载 | 1 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| contentSourceUrl | string | 否 | 原文url链接，当 `type` 为 2 (转载) 时必填 | - |
| collection | object | 否 | 合集信息 | - |

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
  "platforms": ["BilibiliOpen"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BILIBILI_OPEN_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "B站Open视频标题示例",
          "description": "这是关于B站Open版接入视频发布的描述内容。",
          "tags": ["科技", "极客"],
          "category": [
            {
              "id": "1",
              "text": "生活",
              "raw": {}
            }
          ],
          "declaration": 0,
          "type": 1
        }
      }
    ]
  }
}
```
