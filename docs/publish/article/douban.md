# 📄 豆瓣文章发布参数 (DouBan Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“豆瓣 (Douban)”平台发布日记、文章或专栏内容时触发。支持：
- **日记同步**：同步内容至豆瓣日记。
- **原创勾选**：申明内容为原创以保护版权。
- **标签标记**：为日记添加分类标签（如“影评”、“生活”）。
- **草稿保存**：同步至豆瓣日记草稿箱。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装豆瓣 Payload 时需遵守：
1. **社区氛围适配**：豆瓣对营销内容管控严格，Agent 应对正文内容进行基础合规研判，避免包含大量导流外链。
2. **原创声明原则**：若用户明确提到原创，设置 `createType: 1`。
3. **标签简化**：豆瓣标签不建议过长，通常使用 2-4 个关键词为佳。
4. **格式保持**：豆瓣支持基础 HTML 标签，Agent 应确保 `content` 结构整洁。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式)。 | - |
| **`createType`** | `number` | **是** | **创作类型**: `0`-非原创, `1`-申明原创。 | `1` |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `tags` | `string[]` | 否 | 标签数组（如 `["电影", "读书"]`）。 | - |

---

## 4. 执行指令示例 (Command)

```bash
# 发布豆瓣原创日记
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["豆瓣"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DB_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "今日漫想：AI 与文学的交织",
          "content": "<h1>今日漫想</h1><p>正文内容...</p>",
          "createType": 1,
          "pubType": 1,
          "tags": ["文学", "科技"]
        }
      }
    ]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **内容疑似营销被拦截** | 包含多个外部链接或二维码。 | 移除营销干扰项，回归日记本质。 |
| **标题重复报错** | 豆瓣不允许在短时间内发布相同标题的日记。 | 修改标题内容。 |
| **账号需要验证** | 豆瓣触发了双重身份校验。 | 登录蚁小二客户端，按照提示完成验证。 |
| **HTML 标签丢失** | 使用了豆瓣不兼容的高级 HTML 标签。 | 仅保留 h1, h2, p, img 等基础标签。 |

---
> [!TIP]
> **互动特性**: 豆瓣日记发布后，会自动进入用户的广播流，建议配以文艺范儿的封面图。
