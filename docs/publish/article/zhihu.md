# 知乎文章发布参数 (ZhiHu Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `topics` | `Array` | 否 | 话题列表 (`Category[]`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| scheduledTime | number | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| declaration | number | 否 | 创作申明: 0-无申明, 1-剧透, 2-医疗建议, 3-虚构创作, 4-理财内容, 5-AI辅助 | - |

## 2. 复杂对象结构说明

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

### Category (用于话题)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 话题 ID |
| `yixiaoerName` | `string` | **是** | 话题名称 |
| `raw` | `object` | 否 | 平台原始数据 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `topics` | `challenges` | [获取话题/挑战](../../get-challenges.md) |

## 3. Payload 完整示例
- 后端类: `ZhiHuArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`
