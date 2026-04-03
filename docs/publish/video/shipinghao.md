# 视频号 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 视频标题 | - |
| short_title | string | 否 | 视频短标题 | - |
| description | string | 否 | 视频描述，支持 HTML 格式和 `@` 好友/话题标签 | - |
| horizontalCover | object | 否 | 视频横板封面，使用 `OldCover` 结构 | - |
| type | number | 是 | 视频原创类型：1-非原创，2-原创 | 1 |
| createType | number | 是 | 创建类型：1-草稿，2-直接发布 | 2 |
| pubType | number | 否 | 发布类型：0-草稿，1-直接发布 | 1 |
| `location` | `object` | 否 | 视频位置，使用 `PlatformDataItem` 基础结构 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳（单位：秒） | - |
| `shoppingCart` | `object` | 否 | 关联商品信息，包含 `yixiaoerId`, `yixiaoerName`, `raw` | - |
| `collection` | `object` | 否 | 合集信息，包含 `yixiaoerId`, `yixiaoerName`, `raw` | - |
| `activity` | `object` | 否 | 活动信息，包含 `yixiaoerId`, `yixiaoerName`, `raw` | - |
| `music` | `object` | 否 | 背景音乐信息，使用 `MusicItem` 结构 | - |

## 2. 复杂对象结构

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| width | number | 是 | 宽度 |
| height | number | 是 | 高度 |
| size | number | 是 | 文件大小 (Byte) |
| path | string | 否 | 文件绝对路径 (与 `key` 二选一) |
| key | string | 否 | OSS 对象存储 Key (与 `path` 二选一) |

### PlatformDataItem (基础结构)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 统一 ID |
| `yixiaoerName` | `string` | 是 | 显示名称 |
| `raw` | `object` | 是 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |

### MusicItem (音乐)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 蚁小二端统一音乐 ID |
| `yixiaoerName` | `string` | 是 | 歌曲名称 |
| `duration` | `number` | 是 | 音乐时长（秒） |
| `playUrl` | `string` | 是 | 试听/播放链接 |
| `artist` | `string` | 否 | 歌手/作者名 |
| `raw` | `object` | 否 | 平台原始数据。如果在音乐列表获取时该字段存在，发布表单中必须携带并完整透传 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `activity` | `activities` | [获取活动列表](../../get-publish-activities.md) |
| `shoppingCart` | `goods` | [获取商品列表](../../get-goods.md) |
| `music` | `music` | [获取背景音乐](../../get-music.md) |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Shipinghao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号视频标题",
          "description": "<p>这是视频号内容的精彩描述 #生活 #记录</p>",
          "type": 2,
          "createType": 2,
          "pubType": 1,
          "location": {
            "yixiaoerId": "sph_loc_1",
            "yixiaoerName": "广州市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```
