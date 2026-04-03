# 抖音图文发布参数 (Douyin Image-Text)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [图文发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 标题 | - |
| `description` | `string` | 否 | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 地址信息 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐信息 (`MusicItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |
| `collection` | `Object` | 否 | 合集信息 (`Category`) | - |
| `sub_collection` | `Object` | 否 | 合集选集信息 (`Category`) | - |

## 2. 复杂对象结构说明

### OldImage
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `width` | `number` | **是** | 图片宽度 |
| `height` | `number` | **是** | 图片高度 |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `key` | `string` | **是** | 资源 Key (通过上传接口获取) |
| `format` | `string` | **是** | 文件格式 (e.g., `jpg`, `png`) |

### Category (用于合集等)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 蚁小二内部 ID |
| `yixiaoerName` | `string` | **是** | 显示名称 |
| `raw` | `object` | 否 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |

### MusicItem (音乐)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 蚁小二端统一音乐 ID |
| `yixiaoerName` | `string` | 是 | 歌曲名称 |
| `duration` | `number` | 是 | 音乐时长（秒） |
| `playUrl` | `string` | 是 | 试听/播放链接 |
| `artist` | `string` | 否 | 歌手/作者名 |
| `raw` | `object` | 否 | 平台原始数据。如果在音乐列表获取时该字段存在，发布表单中必须携带并完整透传 |

### PlatformDataItem (基础结构)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 统一 ID |
| `yixiaoerName` | `string` | 是 | 显示名称 |
| `raw` | `object` | 是 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |



## 3. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DY_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "抖音图文标题",
          "description": "<p>内容描述 <topic text='合拍' raw='{\"id\":\"xxx\",\"name\":\"合拍\"}'>#合拍</topic></p>",
          "images": [
            { "key": "img_1", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" },
            { "key": "img_2", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "location": { "yixiaoerId": "loc_123", "yixiaoerName": "某商位", "raw": {} }
        }
      }
    ]
  }
}
```

## 相关接口

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `music` | `music` | [获取背景音乐](../../get-music.md) |
| `tags` | `challenges` | [获取话题/挑战](../../get-challenges.md) |
| `images.key` | `upload` | [资源上传](../../upload-resource.md) |
