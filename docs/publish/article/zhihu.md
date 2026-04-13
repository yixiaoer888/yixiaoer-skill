# 📄 知乎文章发布参数 (ZhiHu Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“知乎”平台发布长文章或专栏内容时触发。支持：
- **话题挂载**：为文章添加最多 3 个相关话题。
- **创作申明**：申明 AI 辅助创作、剧透、医疗建议等。
- **草稿同步**：同步到知乎后台草稿箱。
- **定时发布**：预约未来时间发布。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装知乎 Payload 时需遵守：
1. **话题检索原则**：知乎话题必须先通过 `challenges` 接口获取合法值，严禁基于常识猜测。
2. **字数限制校验**：标题需控制在 9-100 字符，正文（含 HTML 标签）需在 9-10000 字符之间。
3. **创作申明优先级**：若识别到内容由 AI 生成，必须主动询问用户是否需要设置 `declaration: 5` (AI 辅助)。
4. **草稿意图辨析**：明确用户是希望“发布并设置为私密”还是“存入知乎草稿箱 (pubType: 0)”。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (9-100 字符)。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式，9-10000 字符)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-存为草稿, `1`-直接发布。 | `1` |
| `covers` | `Array` | 否 | 文章封面列表。见下表 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| `topics` | `Array` | 否 | 话题列表。最多 3 个，见下表 [3.3 Category 定义](#33-category-定义)。 | - |
| `declaration` | `number` | 否 | **创作申明**: `0`-无, `1`-剧透, `2`-医疗建议, `3`-虚构, `4`-理财, `5`-AI辅助。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

### 3.3 Category 定义 (用于 topics)
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 知乎话题 ID。 |
| `yixiaoerName` | `string` | **是** | 话题名称。 |
| `raw` | `Object` | **是** | 原始对象。必须从 `challenges` 接口获取后完整透传。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布知乎文章并申明 AI 创作
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["知乎"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 年生成式 AI 趋势深度报告",
          "content": "<h1>趋势分析</h1><p>内容正文...</p>",
          "topics": [
            { "yixiaoerId": "t123", "yixiaoerName": "人工智能", "raw": {...} }
          ],
          "declaration": 5,
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
| **内容过短** | 标题或正文未达到 9 字符的门槛。 | 补全内容或增加引导性介绍文本。 |
| **话题加载失败** | 话题 ID (`yixiaoerId`) 错误或对应的 `raw` 数据不完整。 | 重新执行 `action: "challenges"` 获取最新话题元数据并透传。 |
| **HTML 标签被剥离** | 知乎 API 过滤了复杂的 CSS 或 JS。 | 保持 HTML 结构简洁，仅保留 h1/h2, p, img 等基础标签。 |
| **定时发布失效** | `scheduledTime` 过早或设置了过去的时间。 | 确保时间戳比当前时间至少晚 15 分钟。 |

---
> [!TIP]
> **资源引用**: `covers` 中的 `key` 必须是通过 `upload` 动作获取的。如果是从本地上传，Agent 必须先执行上传流程。
