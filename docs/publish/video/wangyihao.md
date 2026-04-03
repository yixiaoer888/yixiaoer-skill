# 网易号 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 是 | 视频标签 | - |
| category | object[] | 否 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| declaration | number | 否 | 创作申明：1-内容由AI生成, 2-个人原创仅供参考, 3-取材网络谨慎辨别, 4-虚构演绎仅供娱乐 | - |
| type | number | 是 | 原创声明：0-不勾选, 1-勾选原创 | 0 |
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
  "platforms": ["Wangyihao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WANGYI_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "网易号视频标题示例",
          "description": "这是网易号平台的视频内容描述。",
          "tags": ["科技", "前沿"],
          "category": [
            {
              "id": "1",
              "text": "科技",
              "raw": {}
            }
          ],
          "type": 1,
          "declaration": 2
        }
      }
    ]
  }
}
```
