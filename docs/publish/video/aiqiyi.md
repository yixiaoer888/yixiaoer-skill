# 爱奇艺视频发布 (AiQiYi Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **是** | 视频标题 (长度: 1-50 字符) | - |
| `description` | `string` | **是** | 视频描述 (长度: 1-500 字符) | - |
| `tags` | `string[]` | **否** | 视频标签 (最多 10 个，每个 20 字符) | - |
| `category` | `object[]` | **是** | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| `createType` | `number` | **是** | 内容类型：1-原创，2-转载 | 1 |
| `pubType` | `number` | **是** | 发布类型：0-草稿，1-直接发布 | 1 |
| `declaration` | `number` | **否** | 声明：0-无需申明，1-内容由 AI 生成，2-虚构演绎仅供娱乐，3-取材网络谨慎辨别 | 0 |
| `scheduledTime` | `number` | **否** | 定时任务 (单位: 秒) | - |

## 2. 复杂对象结构

### CascadingPlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `id` | `string` | **是** | 选项 ID |
| `text` | `string` | **是** | 选项文本 |
| `children` | `object[]` | **否** | 子级选项列表 (CascadingPlatformDataItem[]) |
| `raw` | `object` | **是** | 平台原始数据 (如果接口返回，必须原样回传) |

## 3. JSON 示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["爱奇艺"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_aqy_vid_001",
        "video": {
          "key": "v_key_example",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 120
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "爱奇艺精彩视频",
          "description": "这是一段有趣的视频描述，展示了爱奇艺发布流程。",
          "tags": ["影视", "科技", "教程"],
          "category": [
            {
              "id": "cat_id_001",
              "text": "影视",
              "raw": {}
            }
          ],
          "createType": 1,
          "pubType": 1,
          "declaration": 0
        }
      }
    ]
  }
}
```

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |

