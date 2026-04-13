# 📄 微信公众号文章发布参数 (WeiXinGongZhongHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“微信公众号”发布文章时触发。支持：
- **多图文群发**：一次性推送 1-8 篇文章。
- **原创申明**：申明内容原创并设置转载权限。
- **粉丝定向**：按性别、地区对群发对象进行筛选（仅限服务号/订阅号群发）。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装公众号 Payload 时需遵守：
1. **结构特殊性**：微信公众号推荐使用 `publishArgs.platformForms["微信公众号"]` 容器，以支持多账号合并群发。
2. **原创确认流程**：若用户提到“原创”，必须设置 `type: 1` 并在 `categories` 中注入由 `categories` 接口返回的 `raw` 分类数据。
3. **资源引用规范**：文章封面 `cover` 必须包含 `key` 和 `raw`。正文 `content` 支持内嵌 `<account-card>` (公众号) 和 `<video-card>` (视频号) 标签。
4. **群发限制**：确认用户是希望“存为草稿”还是“直接群发 (notifySubscribers: 1)”。

## 3. 参数定义 (Platform Form Parameters)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`articles`** | `Array` | **是** | 文章列表。支持 1-8 篇，见下表 [3.1 Articles 定义](#31-articles-定义)。 | - |
| **`notifySubscribers`** | `number` | **是** | **推送开关**: `0`-不群发（存入草稿箱）, `1`-直接群发通知粉丝。 | `0` |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。通常与群发开关同步设置。 | `1` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |
| `sex` | `number` | 否 | 粉丝性别筛选: `0`-全部, `1`-男, `2`-女。 | `0` |
| `country`/`province`/`city` | `string` | 否 | 分地区群发配置。 | - |

### 3.1 Articles 定义 (WxGongZhongHaoContentFrom)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`title`** | `string` | **是** | 标题 (最多 64 字)。 |
| **`content`** | `string` | **是** | 正文 (HTML)。支持内嵌视频号卡片等自定义标签。 |
| **`cover`** | `object` | **是** | 封面对象。必须包含 `key` 和 `raw` 数据。 |
| **`type`** | `number` | **是** | 原创标志: `0`-不申明, `1`-申明原创。 |
| `authorName` | `string` | 否 | 作者署名 (最多 8 字)。申明原创时必填。 |
| `digest` | `string` | 否 | 摘要 (最多 120 字)。 |
| `quickRepost` | `number` | 否 | 允许快捷转载: `1`-是, `0`-否。 |
| `categories` | `Array` | 否 | 原创分类数组。须包含 `yixiaoerId`, `yixiaoerName`, `raw`。 |

## 4. 执行指令示例 (Command)

```bash
# 微信公众号群发：单图文原创文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["微信公众号"],
  "publishArgs": {
    "accountForms": [{ "platformAccountId": "WX_001" }],
    "platformForms": {
      "微信公众号": {
        "articles": [{
          "title": "深度：AI 时代的自媒体运营",
          "content": "<h1>正文标题</h1><p>内容...</p>",
          "cover": { "key": "c_key", "raw": {...} },
          "type": 1,
          "authorName": "王小二",
          "categories": [{ "yixiaoerId": "c1", "yixiaoerName": "科技", "raw": {...} }]
        }],
        "notifySubscribers": 1,
        "pubType": 1
      }
    }
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **封面获取失败** | 封面对象缺失 `raw` 透传数据。 | 微信 API 强制要求封面必须包含原始资源 ID（位于 raw 中）。 |
| **原创校验不通过** | 缺少作者或分类信息。 | 申明原创时，确保 `authorName` 和 `categories` 已填入。 |
| **群发额度耗尽** | 订阅号/服务号每天群发字数或次数超限。 | 建议用户将多篇文章合并为图文消息一次性发出。 |
| **正文格式拦截** | 包含非法外链或 script 脚本。 | 移除所有外部 JS，确保链接均为微信认可的域名（如视频号卡片）。 |

---
> [!IMPORTANT]
> **多图文排序原则**：`articles` 数组的顺序即为用户在手机端看到的排序。数组第一个元素为“大图”头条，后续为并列小图。
