# 微信公众号文章发布参数 (WeiXinGongZhongHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。微信公众号支持单次发布多篇文章（图文消息）。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `notifySubscribers` | `number` | **是** | 是否群发: 0-不群发（仅存为草稿/定时）, 1-群发 | `1` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，秒) | - |
| `articles` | `Array` | **是** | 文章列表 (`WxArticleItem[]`)，至少包含 1 篇 | - |
| `sex` | `number` | 否 | 群发性别: 0-全部, 1-男, 2-女 | `0` |
| `country` | `string` | 否 | 群发国家名称 | - |
| `province` | `string` | 否 | 群发省份名称 | - |
| `city` | `string` | 否 | 群发城市名称 | - |

## 2. 复杂对象结构说明

### WxArticleItem (文章项)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | **是** | 文章标题 (最多 64 字) |
| `content` | `string` | **是** | 文章正文 (HTML 字符串) |
| `cover` | `Object` | **是** | 封面对象 (`OldCover`) |
| `createType` | `number` | **是** | 创作类型: 0-不声明, 1-声明原创 |
| `digest` | `string` | 否 | 文章摘要 |
| `authorName` | `string` | 否 | 作者名称 |
| `categories` | `Array` | 否 | 文章分类列表 (`Category[]`) |
| `quickRepost` | `number` | 否 | 快捷转载: 0-关闭, 1-开启 |
| `contentSourceUrl` | `string` | 否 | 原文链接 |

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key |
| `size` | `number` | **是** | 文件大小 (Bytes) |
| `width` | `number` | **是** | 宽度 |
| `height` | `number` | **是** | 高度 |

### Category (分类)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 分类 ID |
| `yixiaoerName` | `string` | **是** | 分类名称 |
| `raw` | `object` | 否 | 原始分类数据 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `categories` | `categories` | [获取账号分类](../../get-publish-categories.md) |
| `sex/country/city` | - | 请按平台常规中文名称填写 |

## 3. Payload 完整示例
- 后端类: `WxGongZhongHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/wxgongzhonghao.dto.ts`
