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

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["微信公众号"],
  "publishArgs": {
    "content": "<h1>第一篇文章正文</h1>",
    "accountForms": [
      {
        "platformAccountId": "acc_wx_001",
        "cover": { "key": "main_cover_key", "size": 102400, "width": 900, "height": 500 },
        "contentPublishForm": {
          "formType": "task",
          "notifySubscribers": 1,
          "articles": [
            {
              "title": "这是第一篇文章标题",
              "content": "<h1>第一篇文章正文</h1>",
              "digest": "这是文章摘要",
              "cover": { "key": "main_cover_key", "size": 102400, "width": 900, "height": 500 },
              "createType": 1,
              "authorName": "作者名",
              "quickRepost": 1
            }
          ]
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构

### 3.1 WxArticleItem (文章项)
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | **是** | 文章标题 (最多 64 字) |
| `content` | `string` | **是** | 文章正文 (HTML 字符串)。支持嵌套特殊标签，如公众号卡片 (`account-card`)、视频号卡片 (`video-card`) |
| `cover` | `Object` | **是** | 封面对象 (`OldCover`) |
| `createType` | `number` | **是** | 创作类型: 0-不声明, 1-声明原创 |
| `digest` | `string` | 否 | 文章摘要 (最多 129 字) |
| `authorName` | `string` | 否 | 作者名称 (声明原创时必传，最多 8 字) |
| `categories` | `Array` | 否 | 文章分类数组 (`Category[]`)，声明原创时可选 |
| `quickRepost` | `number` | 否 | 快捷转载: 0-关闭, 1-开启 (声明原创时有效) |
| `contentSourceUrl`| `string`| 否 | 原文链接 |
| `quickPrivateMessage` | `number` | 否 | 快捷私信: 0-关闭, 1-开启 (声明原创时有效，与留言互斥) |

### 3.2 Category (分类结构)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | 蚁小二分类 ID |
| `yixiaoerName`| `string` | 分类名称 |
| `raw` | `Object` | 原始平台分类数据 |

### 3.3 OldCover (封面对象)
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `key` | `string` | OSS 资源 Key |
| `size` | `number` | 文件大小 (bytes) |
| `width` | `number` | 宽度 |
| `height` | `number` | 高度 |

## 4. DTO 参考
- 后端类: `WxGongZhongHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/wxgongzhonghao.dto.ts`
