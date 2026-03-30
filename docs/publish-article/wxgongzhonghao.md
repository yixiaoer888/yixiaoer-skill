# 微信公众号文章发布 (Publish WeChat Article)

该指令用于通过文章引擎向微信公众号分发长内容，支持公众号特有的原创申明、群发设置、定时发布以及多图文属性。

## DTO 溯源 (Knowledge from WxGongZhongHaoArticleForm)
*逻辑来源: `apps/server-api/.../wxgongzhonghao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **建议使用 `article`** (兼容 `weixin-gongzhonghao`) | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 最多 64 字 |
| `--content` | string | 是 | HTML 内容 | 支持自定义标签如 `account-card`, `video-card` |
| `--cover_url` | string | 是 | 文章封面图 URL | 引擎自动上传并转为 `cover` 对象 |
| `--digest` | string | 否 | 文章摘要 | 最多 129 字 |
| `--author` | string | 否 | 作者名称 | 申明原创时必填，最多 8 字 |
| `--original` | boolean | 否 | 是否申明原创 | 对应 `createType` (true -> 1, false -> 0) |
| `--notify` | boolean | 否 | 是否群发通知 | true: 群发 (1), false: 仅存草稿 (0)。默认 true |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳（秒） |
| `--contentSourceUrl`| string| 否 | 原文链接 | 对应 DTO 的 `contentSourceUrl` |
| `--quickRepost` | number | 否 | 快捷转载 | 原创时可用。0: 关闭, 1: 开启 (默认 1) |
| `--quickPrivateMessage` | number | 否 | 快捷私信 | 原创时可用。0: 关闭, 1: 开启 (默认 1) |
| `--categories` | array | 否 | 文章分类 | JSON 数组，申明原创时必填。见下面的分类查询 |
| `--sex` | number | 否 | 群发性别过滤 | 0: 全部, 1: 男, 2: 女 |
| `--country` | string | 否 | 群发地区-国家 | 默认全部 |
| `--province` | string | 否 | 群发地区-省份 | |
| `--city` | string | 否 | 群发地区-城市 | |
| `--articles` | json | 否 | **文章列表 (核心)** | 公众号专用的多图文结构。若不传，引擎会根据 title/content 自动构造单图文。详见下文。 |

## 嵌套模型定义 (Nested Models)

### 文章内容项 (WxGongZhongHaoContentFrom)
这是 `articles` 数组中的单个元素模型。一个账号发布任务通常包含 1 到 8 个此类对象。

```json
{
  "title": "文章标题",
  "content": "<p>正文内容...</p>",
  "digest": "文章摘要",
  "cover": {
    "key": "oss_key_123",
    "width": 900,
    "height": 383,
    "size": 0
  },
  "createType": 1, // 0: 不申明, 1: 申明原创
  "authorName": "作者名",
  "quickRepost": 1, // 0: 关闭, 1: 开启
  "categories": [
    {
      "yixiaoerId": "123",
      "yixiaoerName": "分类名",
      "raw": {}
    }
  ],
  "contentSourceUrl": "原文链接",
  "quickPrivateMessage": 1 // 0: 关闭, 1: 开启
}
```

## 调用指令示例 (Usage)

### 1. 发布原创文章并群发（最常用）
```bash
node scripts/publish.ts \
  --type=article \
  --title="深度解析 AI 代理的未来" \
  --content="<p>正文内容...</p>" \
  --platforms="微信公众号" \
  --account_ids="wx_abc_123" \
  --cover_url="https://example.com/cover.jpg" \
  --author="Antigravity" \
  --original=true \
  --notify=true
```

### 2. 定时发布
```bash
node scripts/publish.ts \
  --type=article \
  --title="定时任务测试" \
  --content="<p>内容...</p>" \
  --platforms="微信公众号" \
  --account_ids="wx_abc_123" \
  --cover_url="https://example.com/cover.jpg" \
  --scheduledTime=1711536000
```

## 逻辑映射说明 (Engine Logic)
- **参数桥接**: `publish.ts` 引擎会自动识别 `platforms` 中的“微信公众号”，并将单篇参数（title, content, digest等）封装进其要求的 `articles` 数组中。
- **模式归一**: 即使传入 `--type=weixin-gongzhonghao`，引擎也会在接口层将其归一化为标准的 `article` 发布模态。
- **资源处理**: `cover_url` 会被 `uploadResource` 转换为 `cover` 对象（含 key, width, height），这是公众号发布的核心必填项。
- **默认值**: 引擎默认会开启 `quickRepost` 和 `quickPrivateMessage` 以优化原创体验。
