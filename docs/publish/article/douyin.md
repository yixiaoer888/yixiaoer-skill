# 抖音文章发布参数 (Douyin Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 (最多 50 字) | - |
| `content` | `string` | **是** | 文章正文 (HTML 格式，最多 50000 字符) | - |
| `description` | `string` | 否 | 文章描述或摘要 (最多 200 字) | - |
| `covers` | `Array` | **是** | 封面图 OSS 列表 (`OldCover[]`, 1-9 张) | - |
| `headImage` | `Object` | 否 | 文章头图 (`OldCover`) | - |
| `music` | `Object` | 否 | 平台音乐背景 (`MusicItem`) | - |
| `topics` | `Array` | 否 | 话题标签列表 (`Category[]`, 最多 5 个) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `visibleType` | `number` | **是** | 可见性: 0-公开, 1-私密, 3-好友 | `0` |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["抖音"],
  "publishArgs": {
    "content": "<h1>抖音文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "coverKey": "cover_oss_key",
        "cover": { "key": "cover_oss_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是一篇抖音文章",
          "content": "<h1>抖音文章标题</h1><p>正文内容...</p>",
          "covers": [
            { "key": "cover_oss_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "visibleType": 0
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

### 3.2 Category (分类/话题对象)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | ID |
| `yixiaoerName` | `string` | **是** | 名称 |
| `raw` | `Object` | **是** | 平台原始对象 (透传) |

### 3.3 MusicItem (音乐)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 内部 ID |
| `yixiaoerName` | `string` | **是** | 歌曲名称 |
| `duration` | `number` | **是** | 时长 (秒) |
| `playUrl` | `string` | **是** | 播放链接 |
| `raw` | `object` | **是** | 原始数据 (透传) |

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `covers.key` | `upload` | [资源上传](../../upload-resource.md) |
| `topics` | `categories` | [获取发布分类](../../get-publish-categories.md) |
| `music` | `music` | [获取背景音乐](../../get-music.md) |
