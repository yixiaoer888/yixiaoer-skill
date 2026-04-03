# 抖音视频发布参数 (DouYin Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


在本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **是** | 抖音视频标题 (1-50 字) | - |
| `description` | `string` | **是** | 抖音视频描述 (1-500 字) | - |
| `horizontalCover` | `object` | 否 | 抖音视频横板封面，使用 `OldCover` 结构 | - |
| `statement` | `number` | 否 | 声明: 3-内容从 AI 生成, 4-可能引人不适, 5-虚构演绎, 6-危险行为 | - |
| `location` | `object` | 否 | 抖音视频位置，使用 `PlatformDataItem` 结构 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒) | - |
| `allow_save` | `number` | 否 | 保存权限: 0-不允许, 1-允许 | 0 |
| `shoppingCart` | `object[]` | 否 | 购物车列表，使用 `ShoppingCart` 结构 | - |
| `groupShopping` | `object` | 否 | 团购信息，使用 `ShoppingCart` 结构 | - |
| `collection` | `object` | 否 | 合集信息，使用 `Category` 结构 | - |
| `sub_collection` | `object` | 否 | 合集选集，使用 `Category` 结构 | - |
| `sync_apps` | `object[]` | 否 | 同时发布应用，使用 `Category[]` | - |
| `hot_event` | `object` | 否 | 热点事件，使用 `Category` 结构 | - |
| `challenge` | `object` | 否 | 挑战/话题，使用 `Category` 结构 | - |
| `mini_app` | `object` | 否 | 挂载小程序 (与购物车互斥)，使用 `MiniApp` 结构 | - |
| `music` | `object` | 否 | 背景音乐信息，使用 `MusicItem` 结构 | - |
| `cooperation_info` | `object` | 否 | 共创信息 | - |
| `game` | `object` | 否 | 游戏挂载信息，使用 `GameItem` 结构 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
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

## 3. 复杂对象结构说明

### 3.1 OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

### 3.2 PlatformDataItem / Category / MiniApp / GameItem
所有统一的基础结构必须包含 `yixiaoerId`, `yixiaoerName`, `raw`。
- `raw`: 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传。

### 3.3 ShoppingCart (购物车/团购)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `sale_title` | `string` | **是** | 推广标题 (最多 10 字) |
| `raw` | `object` | **是** | 平台原始数据 (透传) |
| `brand_switch_value` | `number` | 否 | (团购专用) 0:不推荐, 1:推荐其他品牌, 2:只推荐同品牌 |

### 3.4 MusicItem (音乐)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 内部 ID |
| `yixiaoerName` | `string` | **是** | 歌曲名称 |
| `duration` | `number` | **是** | 时长 (秒) |
| `playUrl` | `string` | **是** | 播放链接 |
| `raw` | `object` | **是** | 原始数据 (透传) |

## 相关接口

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
| `shoppingCart` | `goods` | [获取商品列表](../../get-goods.md) |
| `mini_app` | `miniapps` | [获取小程序列表](../../get-miniapps.md) |
| `challenge` | `challenges` | [获取挑战列表](../../get-challenges.md) |
| `music` | `music` | [获取背景音乐](../../get-music.md) |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
