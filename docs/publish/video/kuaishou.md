# 快手视频发布参数 (Kuaishou Video)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [视频发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台视频发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定为 `task` | `task` |
| `title` | `string` | **否** | 视频标题 | - |
| `description` | `string` | **否** | 视频描述 | - |
| `declaration` | `number` | **否** | 视频声明：0-不申明, 1-内容为 AI 生成, 2-演绎情节仅供娱乐, 3-个人观点仅供参考 | 0 |
| `location` | `object` | **否** | 视频位置，使用 `PlatformDataItem` 结构 | - |
| `visibleType` | `number` | **是** | 可见类型：0-公开, 1-私密, 3-好友可见 | 0 |
| `scheduledTime` | `number` | **否** | 定时发布时间戳 (单位: 秒) | - |
| `shopping_cart` | `object` | **否** | 关联商品信息 | - |
| `collection` | `object` | **否** | 合集信息，使用 `Category` 结构 | - |
| `mini_app` | `object` | **否** | 挂载小程序 (与购物车互斥) | - |
| `nearby_show` | `boolean` | **否** | 是否同城展示 | `true` |
| `allow_same_frame` | `boolean` | **否** | 是否允许同框 | `false` |
| `allow_download` | `boolean` | **否** | 是否允许下载 | `false` |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["快手"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_ks_vid_001",
        "video": {
          "key": "v_key_example",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "记录美好生活",
          "description": "这是快手的一条视频动态 #生活 #正能量",
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

## 3. 复杂对象结构说明

### 3.1 PlatformDataItem (基础结构)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 统一 ID |
| `yixiaoerName` | `string` | **是** | 显示名称 |
| `raw` | `object` | **是** | 平台原始数据 (必须完整透传) |

### 3.2 Category
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 蚁小二内部 ID |
| `yixiaoerName` | `string` | **是** | 显示名称 |
| `raw` | `object` | **是** | 平台原始数据 (透传) |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `video.key` | `upload` | [资源上传](../../upload-resource.md) |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |
