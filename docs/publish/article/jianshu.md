# 📄 简书 文章 参数 (JianShu Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“简书”平台发布随笔、短文、技术分享或文学作品时触发：
- **极简发布**：快速将文字内容同步到简书社区。
- **草稿同步**：将内容保存到简书后台以便后续精修。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装简书 Payload 时需遵守：
1. **轻量化原则**：简书对参数要求相对简单，主要聚焦于标题和 HTML 正文。
2. **正文转换**：确保 `content` 为 HTML 格式，简书对 Markdown 渲染支持良好，但 Payload 传输建议使用 HTML。
3. **可见性逻辑**：简书默认发布即公开，若需私密保存，建议设置 `pubType: 0` (草稿)。
4. **资源策略**：简书正文内引用的图片建议通过 `upload` 动作预先转换。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |

## 4. 执行指令示例 (Command)

```bash
# 在简书发布一篇技术随笔
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["简书"],
  "publishArgs": {
    "content": "<h1>我的简书创作日记</h1><p>今日分享内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_js_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "我的简书创作日记",
          "content": "<h1>我的简书创作日记</h1><p>今日分享内容...</p>",
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
| **标题重复报错** | 简书限制发布完全相同标题的内容。 | 稍微修改标题或增加唯一性后缀。 |
| **HTML 标签丢失** | 正文中包含了简书不支持的复杂自定义标签。 | 保持正文结构简单，仅使用基础图文标签。 |
| **发布频率受限** | 账号短时间内发布过多内容。 | 建议设置时间间隔，或分阶段发布。 |
| **草稿同步失败** | 简书后台接口波动。 | 尝试重新执行发布动作。 |

---
> [!TIP]
> **社区调性**: 简书是一个鼓励原创的写作社区。Agent 建议用户在文章末尾增加“版权申明”或引导关注的文字，以提升个人品牌。
