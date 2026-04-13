# 📄 AcFun 文章发布参数 (AcFun Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“AcFun”平台发布长文章时触发。典型提示词包括：
- “发布这篇文章到 AcFun”。
- “AC 在手，天下我有，同步这篇稿件”。
- “帮我把这篇文章申明原创并同步到 AcFun”。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 AcFun Payload 时需遵守：
1. **原创申明逻辑**：若用户提到“原创”，必须将 `type` 设为 `1`。若是转载，需提醒用户提供原文链接 `contentSourceUrl`。
2. **分类必选原则**：AcFun 文章发布强制要求分类 `category`。Agent 必须通过 `categories` 接口获取合法分类数组，并包含 `raw` 透传数据。
3. **封面规范**：必须通过 `covers` 数组提供至少一张封面图，且包含完整的 `OldCover` 结构（key, size, width, height）。
4. **话题限制**：支持添加 `tags` 话题，但建议数量不要超过 1 个，以符合其文章分发习惯。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 50 字符)。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式，最多 50000 字符)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`category`** | `Array` | **是** | 文章分类。详见 [3.3 Category 定义](#33-category-定义)。 | - |
| **`type`** | `number` | **是** | **创作类型**: `0`-不申明, `1`-申明原创。 | `0` |
| `desc` | `string` | 否 | 文章摘要/描述 (最多 200 字)。 | - |
| `tags` | `string[]` | 否 | 话题标签 (最多 1 个)。 | - |
| `contentSourceUrl` | `string` | 否 | 原文链接 (转载/非原创时必填)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

### 3.3 Category 定义 (分类对象)
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 分类 ID。 |
| `yixiaoerName` | `string` | **是** | 分类名称。 |
| `raw` | `Object` | **是** | 原始对象。必须完整透传。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布 AcFun 原创文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["AcFun"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "AC_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "A 站发布测试：深度分析二次元趋势",
          "content": "<h1>深度分析</h1><p>正文内容...</p>",
          "covers": [
            { "key": "c_key_1", "size": 102400, "width": 800, "height": 600 }
          ],
          "category": [
            { "yixiaoerId": "cat_001", "yixiaoerName": "综合", "raw": {...} }
          ],
          "type": 1,
          "tags": ["动漫"]
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
| **缺少封面图** | `covers` 数组为空或对象结构不完整。 | 确保至少提供一张封面且包含完整的 OldCover 字段。 |
| **分类未选择** | `category` 字段缺失。 | AcFun 强制要求分类，请调用 `action: "categories"` 获取并填入。 |
| **标题过长** | 标题超过 50 字符。 | 精简标题内容。 |
| **转载链接缺失** | 设置为非原创但未提供 `contentSourceUrl`。 | 补全原文链接地址。 |

---
> [!TIP]
> **社区规范**: AcFun 社区对内容质量有较高要求，Agent 建议标题尽量活泼，适配二次元社区氛围。
