# 📄 雪球号 文章 参数 (XueQiuHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“雪球号”平台发布财经评论、投资笔记或市场分析时触发：
- **财经分析**：发布具有专业深度的金融、股票相关图文内容。
- **投资社交**：在财经圈内进行硬核内容同步。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装雪球号 Payload 时需遵守：
1. **标题长度规范**：雪球号标题限制在 9-100 字符之间。Agent 应检查标题是否过短。
2. **正文字数校验**：正文内容同样受到平台审核门槛限制，逻辑上应确保持续的高质量输出。
3. **可见性控制**：支持公开 (0) 和私密 (1) 发布。Agent 应根据内容敏感度提示用户。
4. **资源引用**：封面 `covers` 必须通过 `upload` 动作产生。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (长度: 9-100 字符)。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| **`visibleType`** | `number` | **是** | **可见类型**: `0`-公开, `1`-私密。 | `0` |
| `covers` | `Array` | 否 | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| `declaration` | `number` | 否 | **内容申明**: `0`-不申明, `1`-包含 AI 生成内容。 | `0` |
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
# 发布雪球号公开财经文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["雪球号"],
  "publishArgs": {
    "content": "<h1>深度：2026 年全球芯片市场展望</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_xq_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "深度：2026 年全球芯片市场及投资展望",
          "content": "<h1>市场展望</h1><p>正文内容...</p>",
          "visibleType": 0,
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
| **标题长度不符** | 标题少于 9 个字符或多于 100 字符。 | 请优化标题，确保其在合法范围内。 |
| **内容敏感词拦截** | 雪球号对违禁金融词汇审核极严。 | 检查并替换涉及违规诱导、非法荐股等词汇。 |
| **封面获取异常** | 未通过 `upload` 接口上传封面。 | 确保 `covers` 数组中的 `key` 有效。 |
| **发布频率过快** | 短时间内连续发布。 | 建议设置发布时间间隔。 |

---
> [!TIP]
> **球友互动**: 雪球号是高净值人群聚集地。Agent 建议内容尽量保持专业性，并在正文末尾通过问题形式引导读者在评论区交流。
