# 头条号文章发布参数 (TouTiaoHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 (1-50 字符) | - |
| `content` | `string` | **是** | 文章内容 (HTML 字符串，最多 50000 字符) | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`, 1-9 张) | - |
| `isFirst` | `boolean` | 否 | 是否头条首发 | `false` |
| `location` | `Object` | 否 | 位置对象 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `advertisement` | `number` | 否 | 广告投放收益: 2-无收益, 3-投放广告赚收益 | `3` |
| `declaration`| `number` | 否 | 创作类型 1:自行拍摄 2:取自站外 3:AI生成 6:虚构演绎,故事经历 7:投资观点,仅供参考 8:健康医疗分享,仅供参考 | - |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["头条号"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>文章内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_tt_001",
        "cover": {
          "key": "cover_key_001",
          "width": 800,
          "height": 600,
          "size": 150000
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "头条差异化标题",
          "content": "<h1>头条差异化内容</h1><p>具体内容...</p>",
          "covers": [
            {
              "key": "cover_key_001",
              "width": 800,
              "height": 600,
              "size": 150000
            }
          ],
          "advertisement": 3,
          "declaration": 1,
          "pubType": 1
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

### 3.2 PlatformDataItem (位置信息)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 位置 ID |
| `yixiaoerName` | `string` | **是** | 位置名称 |
| `raw` | `object` | **是** | 原始位置对象 (透传) |

## 相关接口

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |
| `covers.key` | `upload` | [资源上传](../../upload-resource.md) |
