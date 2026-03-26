# 微信公众号文章发布 (Publish WeChat Article)

该指令用于通过文章引擎向微信公众号分发长内容，支持公众号特有的原创申明、群发设置等高级功能。

## DTO 溯源 (Knowledge from WxGongZhongHaoArticleForm)
*逻辑来源: `apps/server-api/.../wxgongzhonghao.dto.ts`*

### 核心参数 (Command Arguments)
| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--title` | string | 是 | 文章标题 | 最多 64 字 |
| `--content` | string | 是 | HTML 内容 | 需包含在 `<html><body>` 中 |
| `--cover_url` | string | 是 | 封面图 | DTO 要求 `cover` 不为空 |
| `--digest` | string | 否 | 摘要 | 最多 120 字 |
| `--author` | string | 否 | 作者 | 申明原创时必须提供，最多 8 字 |
| `--original` | boolean | 否 | 原创申明 | 映射为 DTO 的 `createType` (1: 原创, 0: 普通) |
| `--notify` | boolean | 否 | 群发通知 | 映射为 DTO 的 `notifySubscribers` (1: 群发, 0:仅存草稿/不群发) |
| `--contentSourceUrl`| string| 否 | 原文链接 | 对应 DTO 的 `contentSourceUrl` |

## 调用指令示例 (Usage)

### 1. 发布原创文章并群发
```bash
node scripts/publish-article.ts \
  --title="深度解析 AI 代理" \
  --content="<p>正文内容...</p>" \
  --platforms="微信公众号" \
  --account_ids="wx_123" \
  --cover_url="https://example.com/cover.jpg" \
  --author="Antigravity" \
  --original=true \
  --notify=true
```

### 2. 保存为公众号草稿
```bash
node scripts/publish-article.ts \
  --title="草稿测试" \
  --content="<p>内容</p>" \
  --platforms="微信公众号" \
  --account_ids="wx_123" \
  --pub_type=0
```

## 逻辑说明
- **表单构造**: 引擎会根据 DTO 结构将参数包装进 `contentList` 数组。
- **存证依赖**: 调用前确保 `YIXIAOER_API_KEY` 有效，系统会自动完成 `publishContentId` 存证。
