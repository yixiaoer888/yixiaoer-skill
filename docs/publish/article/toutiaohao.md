# 头条号文章发布参数 (TouTiaoHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `isFirst` | `boolean` | 否 | 是否头条首发 | `false` |
| `location` | `Object` | 否 | 位置对象 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `advertisement` | `number` | **是** | 广告投放收益: 2-无收益, 3-投放广告赚收益 | `3` |
| `declaration` | `number` | 否 | 创作申明: 1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎, 7-投资观点, 8-健康医疗 | - |

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["头条号"],
  "publishArgs": {
    "content": "<h1>文章标题</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_th_001",
        "coverKey": "article_cover_key",
        "cover": { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这是文章标题",
          "covers": [
            { "key": "article_cover_key", "size": 102400, "width": 800, "height": 600 }
          ],
          "advertisement": 3,
          "isFirst": false
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构

### 3.1 OldCover (封面对象)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `key` | `string` | OSS 资源 Key |
| `size` | `number` | 文件大小 (bytes) |
| `width` | `number` | 宽度 |
| `height` | `number` | 高度 |

### 3.2 PlatformDataItem (位置信息)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `id` | `string` | 位置 ID |
| `text` | `string` | 位置名称 |
| `raw` | `Object` | 原始位置对象 |

## 4. DTO 参考
- 后端类: `TouTiaoHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`
