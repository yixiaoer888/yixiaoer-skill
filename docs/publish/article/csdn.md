# 📄 CSDN 文章发布参数 (CSDN Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“CSDN”平台发布技术博客、文章或教程时触发。支持：
- **创作类型选择**：申明原创、转载或翻译。
- **标签精准分类**：为技术文章添加多个专业标签。
- **AI 申明**：申明内容是否由 AI 辅助生成。
- **草稿保存**：同步到 CSDN 个人中心的草稿箱。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 CSDN Payload 时需遵守：
1. **摘要必填逻辑**：CSDN 强制要求摘要 `desc`。Agent 应从正文中提取前 100-200 字作为默认摘要。
2. **标签规范化**：`tags` 应包含文章涉及的核心技术栈（如 "Java", "Python", "架构"）。建议数量为 3-5 个。
3. **创作申明校验**：若 `createType` 设为 `2` (转载)，必须强制用户提供 `contentSourceUrl`。
4. **封面完整性**：封面 `covers` 必须通过 `upload` 动作产生，且包含完整的 `OldCover` 结构（包含宽、高、大小）。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式)。 | - |
| **`desc`** | `string` | **是** | 文章摘要。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`tags`** | `string[]` | **是** | 文章标签数组（如 `["AI", "编程"]`）。 | - |
| **`createType`** | `number` | **是** | **创作类型**: `1`-原创, `2`-转载, `4`-翻译。 | `1` |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **声明**: `0`-无, `1`-AI 辅助, `2`-网络来源, `3`-个人观点。 | `0` |
| `contentSourceUrl` | `string` | 否 | 原文链接 (当 `createType` 为 2 时必填)。 | - |
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
# 发布 CSDN 技术原创文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["CSDN"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "CSDN_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "深度：2026 年后端架构演进趋势",
          "content": "<h1>趋势分析</h1><p>正文内容...</p>",
          "desc": "本文深度探讨了未来后端架构的演进方向，重点关注 Agentic Workflow 的集成。",
          "covers": [
            { "key": "c_key_1", "size": 150000, "width": 800, "height": 450 }
          ],
          "tags": ["后端", "架构", "AI"],
          "createType": 1,
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
| **摘要字符过短** | `desc` 字段内容少于要求（通常至少 20 字符）。 | 增加摘要长度。 |
| **转载链接格式错误** | `contentSourceUrl` 为空或不是合法的 URL。 | 检查并补全原文链接。 |
| **标签数量超限** | `tags` 数组元素超过 5 个。 | 精简标签数量至核心词。 |
| **封面上传失败** | CSDN 对封面图的大小或比例有特定限制。 | 确保图片小于 5MB，推荐 16:9 比例。 |

---
> [!TIP]
> **开发者专区**: CSDN 文章发布后会自动进入“开发者社区”分发，建议开启评论功能以增加互动权重。
