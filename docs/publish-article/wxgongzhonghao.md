# 微信公众号文章发布 (Publish WeChat Article)

该指令用于通过文章引擎向微信公众号分发长内容，支持公众号特有的原创申明、群发设置等高级功能。

## DTO 溯源 (Knowledge from WxGongZhongHaoArticleForm)
*逻辑来源: `apps/server-api/.../wxgongzhonghao.dto.ts`*

### 核心参数 (Command Arguments)
| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `weixin-gongzhonghao`**| 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 最多 64 字 |
| `--content` | string | 是 | HTML 内容 | 需包含在 `<html><body>` 中 |
| `--cover_url` | string | 是 | 封面图 | 直连地址，引擎自动上传 |
| `--digest` | string | 否 | 摘要 | 最多 120 字 |
| `--author` | string | 否 | 作者 | 申明原创时必须提供，最多 8 字 |
| `--original` | boolean | 否 | 原创申明 | 是否在公众号申明原创 |
| `--notify` | boolean | 否 | 群发通知 | true: 群发, false: 仅存草稿 (默认 true) |
| `--contentSourceUrl`| string| 否 | 原文链接 | 对应 DTO 的 `contentSourceUrl` |

## 调用指令示例 (Usage)

### 1. 发布原创文章并群发
```bash
node scripts/publish.ts \
  --type=weixin-gongzhonghao \
  --title="深度解析 AI 代理" \
  --content="<p>正文内容...</p>" \
  --platforms="微信公众号" \
  --account_ids="wx_123" \
  --cover_url="https://example.com/cover.jpg" \
  --author="Antigravity" \
  --original=true \
  --notify=true
```

## 逻辑说明
- **文档驱动**: `publish.ts` 引擎支持透传参数。AI 助手根据本指令文档中的“核心参数”传参，引擎会自动将其转换为 `WxGongZhongHaoArticleForm` 要求的 `contentList` 结构。
- **发布通道**: `type=weixin-gongzhonghao` 会自动执行文章预存证 (Storage) 逻辑。
