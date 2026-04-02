# 视频号图文发布参数 (WeChat Video Account Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 标题 | - |
| `description` | `string` | 否 | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置对象 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐对象 (`MusicItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |
| `collection` | `Object` | 否 | 合集信息，使用 `Collection` 结构 | - |
| `pubType` | `number` | 否 | 发布类型: 0-草稿, 1-直接发布 | 1 |

## 2. 复杂对象结构说明

### OldImage
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `width` | `number` | **是** | 图片宽度 |
| `height` | `number` | **是** | 图片高度 |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `key` | `string` | **是** | 资源 Key (通过上传接口获取) |
| `format` | `string` | **是** | 文件格式 (e.g., `jpg`, `png`) |

### PlatformDataItem (基础结构)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 统一 ID |
| `yixiaoerName` | `string` | 是 | 显示名称 |
| `raw` | `object` | 是 | 平台原始数据 |

### Collection
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 合集 ID |
| `yixiaoerName` | `string` | 是 | 合集名称 |
| `raw` | `object` | 否 | 平台原始数据 |

### MusicItem (音乐)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 蚁小二端统一音乐 ID |
| `yixiaoerName` | `string` | 是 | 歌曲名称 |
| `duration` | `number` | 是 | 音乐时长（秒） |
| `playUrl` | `string` | 是 | 试听/播放链接 |
| `artist` | `string` | 否 | 歌手/作者名 |
| `raw` | `object` | 否 | 平台原始数据，发布时需完整透传 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `music` | `music` | [获取背景音乐](../../get-music.md) |
| `tags` | `challenges` | [获取话题/挑战](../../get-challenges.md) |

## 3. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["视频号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号动态标题",
          "description": "<p>视频号图文内容 <topic text='微信' raw='{\"id\":\"xxx\",\"name\":\"微信\"}'>#微信</topic></p>",
          "images": [
            { "key": "img_sph_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ]
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `ShiPingHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/shipinghao.dto.ts`
