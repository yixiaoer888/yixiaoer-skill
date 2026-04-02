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

## 2. 复杂对象结构说明

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

### PlatformDataItem (位置信息)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `id` | `string` | **是** | 位置 ID |
| `text` | `string` | **是** | 位置名称 |
| `raw` | `object` | 否 | 原始位置对象 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location` | `locations` | [获取位置信息](../../get-locations.md) |

## 3. Payload 完整示例
- 后端类: `TouTiaoHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`
