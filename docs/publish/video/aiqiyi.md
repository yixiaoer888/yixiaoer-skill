# 爱奇艺 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 否 | 视频标签 | - |
| category | object[] | 是 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| type | number | 是 | 内容发布类型：0-草稿，1-直接发布 | 1 |
| declaration | number | 否 | 声明：0-无需申明，1-内容由AI生成，2-虚构演绎仅供娱乐，3-取材网络谨慎辨别 | 0 |
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
  "platforms": ["AiQiYi"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "AIQIYI_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "走进爱奇艺视频标题",
          "description": "这是关于爱奇艺视频发布的一个示例描述信息。",
          "tags": ["影视", "花絮"],
          "category": [
            {
              "id": "101",
              "text": "电影",
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
