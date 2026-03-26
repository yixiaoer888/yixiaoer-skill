# 发布微信公众号 (Publish WeChat Official Account)

该能力允许用户向微信公众号发布文章，支持完整的表单结构，包括作者、摘要、原文链接、原创申明及推文设置。

## 场景示例 (Scenarios)
- "帮我把这篇关于 AI 的深度长文发布到‘我的技术周刊’公众号，作者署名‘张三’，并申明原创。"
- "将文章 (ID: xxx) 作为草稿保存到微信公众号。"

## 参数定义 (Parameters)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| --title | string | 是 | 文章标题 |
| --content | string | 是 | 文章 HTML 内容 (或纯文本，脚本会包装为 HTML) |
| --account_ids | string | 是 | 媒体账号 ID 列表，逗号分隔 (微信公众号通常一次只能选一个，但脚本支持批量) |
| --author | string | 否 | 作者姓名 |
| --digest | string | 否 | 文章摘要，默认为标题前 50 字 |
| --cover_url | string | 否 | 封面图直连 URL |
| --cover_key | string | 否 | 封面图 Key |
| --content_source_url | string | 否 | 原文阅读链接 |
| --original | boolean | 否 | 是否申明原创，默认 false |
| --notify | boolean | 否 | 是否群发通知粉丝，默认 true |
| --pub_type | number | 否 | 0: 草稿, 1: 立即发布。默认 1 |

## 调用指令 (Commands)
```bash
node scripts/publish-wechat.ts \
  --title="我的第一篇推文" \
  --content="Hello WeChat!" \
  --account_ids="wx_account_id_123" \
  --author="张三" \
  --original=true
```

## 注意事项
- 微信公众号对图片素材有严格限制，请确保封面图已通过 `upload-resource.ts` 预先上传或使用 `cover_url`。
- `notify=true` 会占用公众号的群发额度，请谨慎使用。
