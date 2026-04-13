# 📄 企鹅号 文章 参数 (QiEHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“企鹅号 (腾讯内容开放平台)”发布新闻资讯、深度解析或社交动态时触发：
- **全平台分发**：将内容同步到腾讯新闻、天天快报、QQ 浏览器等渠道。
- **合规申明**：标注 AI 生成、个人观点或剧评演绎等创作背景。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装企鹅号 Payload 时需遵守：
1. **标签必填逻辑**：企鹅号强制要求至少提供一个标签 (`tags`)。Agent 应根据文章内容自动提取 2-3 个核心关键词。
2. **合规申明校验**：若识别到内容涉及敏感或 AI 创作，应引导用户正确设置 `declaration`。
3. **资源上传原则**：封面 `covers` 必须通过 `upload` 动作获取 OSS Key，且需符合企鹅号对清晰度的要求。
4. **发布意图确认**：确认是保存为“平台草稿 (pubType: 0)”还是“直接发布 (pubType: 1)”。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`tags`** | `string[]` | **是** | 文章标签 (至少 1 个)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **创作申明**: 0-暂不申明, 1-AI生成, 2-个人观点, 3-剧情演绎, 7-AI辅助, 8-健康医疗, 9-危险行为。 | - |

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
# 发布企鹅号 AI 辅助创作文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["企鹅号"],
  "publishArgs": {
    "content": "<h1>深度：腾讯生态内容分发逻辑</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_qh_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "深度：腾讯生态内容分发逻辑解析",
          "content": "<h1>深度解析</h1><p>正文内容...</p>",
          "covers": [{ "key": "qh_cov_1", "size": 102400, "width": 800, "height": 600 }],
          "tags": ["腾讯", "运营", "流量"],
          "declaration": 1,
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
| **缺少标签** | `tags` 数组为空。 | 企鹅号要求至少 1 个标签，请补充。 |
| **审核未通过** | 标题带有诱导性夸张词汇或正文包含违规链接。 | 遵守腾讯内容审核规范，移除违禁标签。 |
| **封面尺寸错误** | 宽/高比不符合企鹅号 3 图或单图的展示逻辑。 | 建议使用 3:2 或 16:9 比例的高清图。 |
| **申明类型与内容不符** | 实际为 AI 创作但未勾选相应申明。 | 准确勾选 `declaration` 字段以通过机器审核。 |

---
> [!TIP]
> **腾讯生态红利**: 企鹅号是打通腾讯多社交终端的关键，建议在 `tags` 中包含热门关键词以获得 QQ 浏览器等渠道的主动抓取支持。
