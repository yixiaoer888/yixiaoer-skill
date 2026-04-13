# 📄 新浪微博 文章 参数 (XinLangWeiBo Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“新浪微博”发布长文章（头条文章）时触发。典型用途包括：
- **深度长文**：发布超出 140/2000 字限制的图文说明、专栏博客。
- **粉丝订阅**：通过长文章形式为粉丝提供更沉浸的阅读体验。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装新浪微博文章 Payload 时需遵守：
1. **封面校验**：微博长文章强制要求至少一张封面。必须通过 `upload` 动作获取 `key`。
2. **正文渲染**：`content` 必须为 HTML 格式，Agent 应确保图片标签 (`<img>`) 使用的是系统内部 URL 或 `key`。
3. **字数与阈值**：虽然是长文章，但标题仍建议精简在 30 字以内以获得最佳展示。
4. **意图区分**：明确区分“发布短微博 (imageText)”与“发布长文章 (article)”。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布新浪微博长文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["新浪微博"],
  "publishArgs": {
    "content": "<h1>深度：2026 社交媒体趋势</h1><p>内容正文...</p>",
    "accountForms": [
      {
        "platformAccountId": "WB_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 社交媒体趋势报告",
          "content": "<h1>深度趋势报告</h1><p>正文内容...</p>",
          "covers": [
            { "key": "wb_cover_key", "size": 150000, "width": 800, "height": 450 }
          ],
          "pubType": 1
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
| **图片无法显示** | `content` 中的图片链接无效或未通过微博审核。 | 确保所有图片已通过 `upload` 并正确应用在 HTML 标签中。 |
| **标题重复** | 微博限制短时间内发布完全相同标题的文章。 | 稍微调整标题，或增加时间/版本标识。 |
| **定时发布失败** | 时间戳比当前时间早或账号权限不足。 | 检查 Unix 时间戳准确性。 |
| **内容包含敏感词** | 微博审核机制拦截。 | 优化文案，避免违规词汇。 |

---
> [!TIP]
> **传播技巧**: 发布长文章后，系统通常会同步生成一条带有文章卡片的微博。Agent 建议在生成的微博正文中增加引导语，以提升点击率。
