# 得物 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| category | object[] | 否 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| declaration | number | 是 | 创作者申明：0-不添加自主声明, 1-内容由AI生成, 2-内容不含营销推广属性, 3-内容涉及专业运动, 4-剧情演绎仅供娱乐 | 0 |

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
  "platforms": ["Dewu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DEWU_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "得物穿搭视频示例",
          "description": "这是一段关于得物穿搭分享的视频描述。",
          "category": [
            {
              "id": "1",
              "text": "穿搭",
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
