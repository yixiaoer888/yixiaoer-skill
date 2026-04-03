# 快手 视频发布

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 快手标题 | - |
| description | string | 否 | 快手描述 | - |
| declaration | number | 否 | 快手视频声明：1-内容为 AI 生成, 2-演绎情节仅供娱乐, 3-个人观点仅供参考 | - |
| location | object | 否 | 快手视频位置，使用 `PlatformDataItem` 结构 | - |
| visibleType | number | 是 | 可见类型：0-公开, 1-私密, 3-好友可见 | 0 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| shopping_cart | object | 否 | 关联商品信息 | - |
| collection | object | 否 | 合集信息，使用 `Category` 结构 | - |
| mini_app | object | 否 | 挂载小程序（与购物车互斥） | - |
| nearby_show | boolean | 否 | 是否同城展示 | `true` |
| allow_same_frame | boolean | 否 | 是否允许同框 | `false` |
| allow_download | boolean | 否 | 是否允许下载 | `false` |
| music | object | 否 | 背景音乐信息，使用 `MusicItem` 结构 | - |

## 2. 复杂对象结构

### Category
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | 是 | 蚁小二ID |
| `yixiaoerName` | `string` | 是 | 蚁小二名称 |
| `raw` | `object` | 是 | 平台原始数据。如果在获取时该字段存在，发布表单中必须携带并完整透传 |

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

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Kuaishou"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "KS_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "快手视频标题",
          "description": "这是大快手的一条视频动态 #生活 #正能量",
          "visibleType": 0,
          "declaration": 3,
          "nearby_show": true,
          "allow_same_frame": false,
          "allow_download": false
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
