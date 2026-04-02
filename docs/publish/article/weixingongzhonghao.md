# 微信公众号文章发布参数 (WeiXinGongZhongHao Article)

本平台文章发布推荐使用 `platformForms` 模式，以支持多账号共用单条图文包消息。

## 1. 结构概览 (Structure)

对于微信公众号，发布参数应位于 `publishArgs.platformForms["微信公众号"]` 中：

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `articles` | `Array` | **是** | 文章列表 (`WxArticleItem[]`) | - |
| `notifySubscribers` | `number` | **是** | 是否群发: 0-不群发（存为草稿）, 1-群发 | `0` |
| `pubType` | `number` | 否 | 发布方式: 1-正常发布, 2-仅存为草稿 | `1` |
| `sex` | `number` | 否 | 群发性别: 0-全部, 1-男, 2-女 | `0` |
| `country` | `string` | 否 | 群发国家 | - |
| `province` | `string` | 否 | 群发省份 | - |
| `city` | `string` | 否 | 群发城市 | - |

---

## 2. 复杂对象结构说明

### 2.1 WxArticleItem (文章项)

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `title` | `string` | **是** | 文章标题 |
| `content` | `string` | **是** | 正文 (HTML 字符串) |
| `authorName` | `string` | 否 | 作者姓名 |
| `digest` | `string` | 否 | 文章摘要 |
| `cover` | `Object` | **是** | 封面对象 (含 `key` 或 `url`) |
| `type` | `number` | 否 | 原创类型 (0: 否, 1: 是) | `0` |
| `quickRepost` | `number` | 否 | 允许转载: 0-否, 1-是 | `0` |
| `contentSourceUrl` | `string` | 否 | 原文地址链接 | - |
| `videoCardCount` | `number` | 否 | 视频卡片数量 | `0` |

### 2.2 Cover (封面模型)

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | 否 | 上传到 OSS 后的资源 Key |
| `url` | `string` | 否 | 资源的公网访问 URL |

---

## 3. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["微信公众号"],
  "publishArgs": {
    "accountForms": [
      { "platformAccountId": "68ccfed805159d64462300" }
    ],
    "platformForms": {
      "微信公众号": {
        "articles": [
          {
            "title": "以忍道为炬，赴青春之约",
            "content": "<p>正文内容...</p>",
            "authorName": "六道目",
            "digest": "摘要内容...",
            "type": 0,
            "cover": {
                "key": "wxp/test/1775097820677.jpg",
                "url": "https://yixiaoer.cn/wxp/test/1775097820677.jpg"
            },
            "quickRepost": 0,
            "videoCardCount": 0
          }
        ],
        "notifySubscribers": 0,
        "sex": 0,
        "pubType": 1
      }
    }
  }
}
```

## 4. 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `cover.key` | `upload` | [资源上传](../../upload-resource.md) |
| `authorName` | `accounts` | [查询账号列表](../../query-accounts.md) |
