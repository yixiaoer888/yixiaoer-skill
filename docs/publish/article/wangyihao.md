# 网易号文章发布参数 (WangYiHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`) | - |
| `declaration` | `number` | 否 | 创作申明: 1-AI生成, 2-个人原创, 3-取材网络, 4-虚构演绎 | - |
| `type`| `number` | 否 | 是否勾选原创: 0-不勾选, 1-勾选 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |

## 2. 复杂对象结构说明

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

## 3. Payload 完整示例
- 后端类: `WangYiHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/wangyihao.dto.ts`
