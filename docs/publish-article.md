# 发布文章与公众号 (Publish Article & WeChat)

该能力是一个统一的长内容分发引擎，支持向 22 个主流长文平台（含微信公众号、头条、百家号等）发布文章。

## 支持平台 (Supported Platforms)
- **主流门户**: 微信公众号, 头条号, 百家号, 企鹅号, 网易号, 搜狐号, 一点号, 大鱼号, 快传号
- **技术/社区**: CSDN, 知乎, 豆瓣, 简书, 雪球号
- **视频/互娱**: 哔哩哔哩, AcFun (A站), 抖音, 爱奇艺
- **汽车/其他**: 车家号, 易车号, WiFi万能钥匙

## 参数定义 (Parameters)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| --title | string | 是 | 文章标题 |
| --content | string | 是 | 文章 HTML 内容 |
| --platforms | string | 是 | 平台名称列表，逗号分隔 (如 "微信公众号,头条号,百家号") |
| --account_ids | string | 是 | 账号 ID 列表，逗号分隔 |
| --cover_url | string | 否 | 封面图直连 URL |
| --tags | string | 否 | 标签列表，逗号分隔 |
| --author | string | 否 | (仅限公众号) 作者名称 |
| --digest | string | 否 | 文章摘要 |
| --original | boolean| 否 | (仅限公众号) 是否申明原创，默认 false |
| --notify | boolean | 否 | (仅限公众号) 是否群发通知，默认 true |
| --pub_type | number | 否 | 0: 草稿, 1: 立即发布。默认 1 |

## 调用指令 (Commands)

### 1. 发布到多个常规平台
```bash
node scripts/publish-article.ts \
  --title="AI 时代的工具演进" \
  --content="<p>长文内容...</p>" \
  --platforms="头条号,百家号,知乎" \
  --account_ids="id1,id2,id3"
```

### 2. 发布到微信公众号
```bash
node scripts/publish-article.ts \
  --title="技术周刊" \
  --content="内容..." \
  --platforms="微信公众号" \
  --account_ids="wx_id_123" \
  --author="张三" \
  --original=true
```

## 注意事项
- **公众号识别**: 当 `platforms` 仅包含 `微信公众号` 时，系统会自动切换为专用发布通道以支持群发设置。
- **存证逻辑**: 脚本会自动在客户端生成 `publishContentId` 并完成内容存证。
