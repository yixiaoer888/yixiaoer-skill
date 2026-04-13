# 📄 易车号 文章 参数 (YiCheHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“易车号”平台发布汽车资讯、行情点评、新车试驾或技术分析时触发：
- **行业报道**：发布包含横竖版双封面的汽车专业文章。
- **多端分发**：利用易车网全平台渠道进行内容曝光。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装易车号 Payload 时需遵守：
1. **双封面试配原则**：易车号强烈建议同时提供 `covers` (横版) 和 `verticalCovers` (竖版) 封面，以适配手机与 PC 不同端的展示。
2. **话题内嵌逻辑**：支持添加 `topics` 话题。Agent 必须通过 `categories` 接口获取并透传 `raw` 对象数据。
3. **内容权限控制**：通过 `allowForward` 和 `allowAbstract` 控制内容的转发和摘要生成权限。
4. **资源先行原则**：所有图片资源的 `key` 必须通过 `upload` 动作获取。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 28 字符)。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | **横版封面**列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`verticalCovers`** | `Array` | **是** | **竖版封面**列表。结构同 OldCover。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **创作申明**: 0-不申明, 1-个人观点, 2-网络来源, 3-AI生成, 4-引用站内。 | - |
| `allowForward` | `boolean` | 否 | 是否允许转发。 | `false` |
| `allowAbstract` | `boolean` | 否 | 是否允许生成摘要。 | `false` |
| `topics` | `Array` | 否 | 话题列表。使用 `Category[]` 结构。 | - |
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
| `yixiaoerId` | `string` | **是** | 话题 ID。 |
| `yixiaoerName` | `string` | **是** | 话题名称。 |
| `raw` | `Object` | **是** | 原始对象。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布易车号单图文资讯
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["易车号"],
  "publishArgs": {
    "content": "<h1>起底：2026 款豪华越野车性能实测</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_ych_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 款豪华越野车性能实测",
          "content": "<h1>性能实测</h1><p>正文内容...</p>",
          "covers": [{ "key": "ych_h_1", "size": 150000, "width": 800, "height": 600 }],
          "verticalCovers": [{ "key": "ych_v_1", "size": 150000, "width": 600, "height": 800 }],
          "declaration": 1,
          "allowForward": true,
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
| **标题过长** | 标题超过 28 个字符。 | 易车号标题要求简练，请精简至 28 字以内。 |
| **竖版封面不合规** | `verticalCovers` 比例或尺寸错误。 | 建议使用 3:4 比例，确保清晰度。 |
| **话题加载失败** | `topics.raw` 数据格式不正确。 | 必须调用 `categories` 接口并透传完整 `raw` 内容。 |
| **申明与内容不符** | 实际为 AI 生成但未勾选 `declaration: 3`。 | 根据实际创作情况准确设置申明。 |

---
> [!TIP]
> **汽车垂直生态**: 易车号对高质量的图文报告有较好的流量倾斜，建议 Agent 引导用户在标题中包含具体的“汽车品牌+车型”关键字以提升分发精准度。
