# 知乎 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | 否 | 视频标题 (1-50 字符) | - |
| `description` | `string` | **是** | 视频描述 (1-500 字符) | - |
| `category` | `Array` | **是** | 视频分类 (`CascadingPlatformDataItem[]`) | - |
| `declaration` | `number` | 否 | 视频创作申明：0-不声明, 2-图片/视频由AI生成 | 0 |
| `createType` | `number` | **是** | 内容类型：1-原创, 2-转载 | 1 |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |
| `scheduledTime` | `number` | 否 | 定时发布时间戳（单位：秒） | - |

## 2. 复杂对象结构

### Category
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 蚁小二ID |
| `yixiaoerName` | `string` | **是** | 蚁小二名称 |
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

### CascadingPlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `id` | `string` | **是** | 选项ID |
| `text` | `string` | **是** | 选项文本 |
| `children` | `Array` | 否 | 子级选项列表 (`CascadingPlatformDataItem[]`) |
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

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
          "category": [
            {
              "id": "1",
              "text": "科技",
              "raw": {}
            }
          ],
          "createType": 1,
          "declaration": 0,
          "pubType": 1
        }

      }
    ]
  }
}
```

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `category` | `categories` | [获取视频分类](../../get-categories.md) |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
