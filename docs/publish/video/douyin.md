# 抖音 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 抖音视频标题 | - |
| description | string | 否 | 抖音视频描述 | - |
| horizontalCover | object | 否 | 抖音视频横板封面，使用 `OldCover` 结构 | - |
| statement | number | 否 | 抖音视频声明：3-内容由AI生成, 4-可能引人不适, 5-虚构演绎仅供娱乐, 6-危险行为请勿模仿 | - |
| location | object | 否 | 抖音视频位置，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| allow_save | number | 否 | 保存权限：0-不允许, 1-允许 | 0 |
| shoppingCart | object[] | 否 | 购物车列表 | - |
| groupShopping | object | 否 | 团购信息 | - |
| collection | object | 否 | 合集信息 | - |
| sub_collection | object | 否 | 合集选集，使用 `Category` 结构 | - |
| sync_apps | object[] | 否 | 同时发布到的应用列表，使用 `Category` 结构数组 | - |
| hot_event | object | 否 | 热点事件，使用 `Category` 结构 | - |
| challenge | object | 否 | 挑战/话题，使用 `Category` 结构 | - |
| mini_app | object | 否 | 挂载小程序（与购物车互斥），使用 `MiniApp` 结构 | - |
| music | object | 否 | 背景音乐信息，使用 `MusicItem` 结构 | - |
| cooperation_info | object | 否 | 共创信息 | - |
| game | object | 否 | 游戏挂载信息，使用 `GameItem` 结构 | - |

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

### Category
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 蚁小二端统一 ID |
| `yixiaoerName` | `string` | 是 | 平台原始名称 |
| `raw` | `object` | 否 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |

### ShoppingCart (购物车/团购)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| sale_title | string | 是 | 推广标题 (最多10个字) |
| `raw` | `object` | 是 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |
| brand_switch_value | number | 否 | 仅团购使用。0:不推荐, 1:推荐其他品牌, 2:只推荐相同品牌 |

### MiniApp (小程序)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 小程序 ID |
| `yixiaoerName` | `string` | 是 | 小程序名称 |
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

### GameItem (游戏)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 游戏 ID |
| `yixiaoerName` | `string` | 是 | 游戏规格/名称 |
| `raw` | `object` | 否 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |



## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Douyin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DOUYIN_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "记录美好生活",
          "description": "这是我在抖音的第一条视频 #美好生活 #见闻",
          "statement": 3,
          "allow_save": 1,
          "location": {
            "yixiaoerId": "123",
            "yixiaoerName": "上海市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```

## 相关接口

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `shoppingCart` / `groupShopping` | `goods` | [获取商品列表](../../get-goods.md) |
| `hot_event` | `hot-events` | [获取热点列表](../../get-hot-events.md) |
| `challenge` | `challenges` | [获取挑战列表](../../get-challenges.md) |
| `mini_app` | `miniapps` | [获取小程序列表](../../get-miniapps.md) |
| `sync_apps` | `syncapps` | [获取同步应用](../../get-sync-apps.md) |
| `music` | `music` | [获取背景音乐](../../get-music.md) |
| `game` | `games` | [获取游戏列表](../../get-games.md) |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
