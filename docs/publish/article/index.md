# 📄 文章发布通用索引 (Article Publish Index)

> [!IMPORTANT]
> **能力定位**: 本文档定义了“文章”类内容发布的根 Payload 结构与核心校验逻辑。在查阅具体的平台文档（如 `weixingongzhonghao.md`, `zhihu.md`）之前，Agent **必须** 首先掌握本文档，以确保根部结构的完整性。

## 1. 触发场景 (Trigger)

当系统分析到用户的意图为“分发长图文内容”时加载。典型提示词包括：
- “发布这篇文章到知乎和 CSDN”。
- “把这段 HTML 发布到我的所有公众号”。
- “同步这篇博文，并设置封面图”。
- “帮我把这个内容存为蚁小二草稿”。

## 2. 交互协议 (Interactive Protocol)

Agent 在构造文章发布 Payload 时，必须遵循以下“三步法”：

1. **资源预处理 (Pre-processing)**：
   - **格式转换**：确保正文 `content` 已转换为标准的 HTML 格式。
   - **封面锁定**：检查封面图。若为外部 URL，**必须** 先调用 `upload` 动作获取系统 `key`，严禁直接透传原始 URL。
2. **发布模式判定 (Mode Selection)**：
   - **微信公众号模式**：若包含“微信公众号”，则 `platforms` 数组中**严禁** 出现其他平台。必须使用 `platformForms` 结构进行精细化控制（详见其专属文档）。
   - **多平台通用模式**：使用 `accountForms` 数组承载多个账号的配置，支持不同账号使用不同的分类或封面。
3. **级联分类处理 (Cascading Categories)**：
   - Agent **必须** 自行装配 `category` 数组，每一级均需包含 `yixiaoerId`, `yixiaoerName` 及 `raw` 对象。
   - **下钻约束 (Leaf Node)**：若平台分类存在二级或多级结构，Agent **必须** 引导选择至叶子节点（最底层分类）。严禁漏选子分类或仅选择一级分类。
4. **预览确认 (Preview & Confirm)**：
   - 在执行前，向用户展示生成的标题、选中的封面缩略图及目标平台列表。**明确获得用户“发布”指令后**方可调用脚本。

## 3. 参数定义 (Parameters)

### 3.1 根 Payload 结构 (Root Structure)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | 固定值：`publish` | `publish` |
| **`publishType`** | `string` | **是** | 固定值：`article` | `article` |
| **`platforms`** | `string[]` | **是** | 目标平台枚举数组。例如 `["知乎", "CSDN"]`。 | - |
| **`publishArgs`** | `Object` | **是** | 核心发布参数，包含正文与账号配置。 | - |
| `isDraft` | `boolean` | 否 | 若为 `true`，则仅作为蚁小二系统的“内部草稿”保存，不执行远端分发。 | `false` |

### 3.2 核心发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`content`** | `string` | **是** | **HTML 格式** 的正文内容。 |
| **`accountForms`** | `Array` | **是** | 账号级发布配置列表。包含账号 ID、封面等信息。 |
| `platformForms` | `Object | 否 | **仅限微信公众号使用**。用于承载多图文、设置同步开关等。 |

### 3.3 账号配置项 (accountForms Item)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`platformAccountId`** | `string` | **是** | 蚁小二平台账号唯一 ID。 |
| **`cover`** | `Object` | **是** | **ImageFormItem**: 主封面对象。包含 `key`, `width`, `height`, `size`。 |
| `contentPublishForm` | `Object` | 否 | **平台私有字段容器**：若目标平台有特殊需求（如 AcFun 的标签），在此填入。 |
| `coverKey` | `string | 否 | 冗余字段。通常与 `cover.key` 保持一致，建议填入。 |

---

## 4. 执行指令示例 (Command)

`ash
# 通用文章发布示例：同步到知乎与 CSDN
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["知乎", "CSDN"],
  "publishArgs": {
    "content": "<h1>深度：AI 驱动的写作变革</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "ACC_ZH_001",
        "cover": {
          "key": "article_cover_key_001",
          "width": 800,
          "height": 600,
          "size": 102400
        },
        "coverKey": "article_cover_key_001"
      }
    ]
  }
}'
`

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **`YIXIAOER_USAGE_ERR`** | 正文包含非法标签或内容字段缺失。 | 检查 `content` 是否为有效的 HTML。确认 `accountForms` 是否为空。 |
| **分类不识别 / ID 非法** | `category` 仅传入了末级 ID 或停留在了一级大类（未到二级叶子节点）。 | **级联下钻要求**：必须提供从根到叶子的完整对象数组。对于百家号等强分类平台，严禁停留在一级分类。 |
| **封面显示破损** | 使用了外部 URL 且未通过 `upload` 动作转换。 | **禁止透传外链**：必须先执行 `upload` 动作获取系统内部 `key` 后再引用。 |
| **多图文发布冲突** | 将“微信公众号”与其他通用平台在同一个 Payload 中混发。 | **物理隔离原则**：微信公众号必须作为一个独立的发布动作执行。 |
| **发布平台不在支持列表** | 目标账号（如小红书）不支持“文章”类型。 | **类型对标**：长文章平台（百家号等）不支持短动态，短动态平台（小红书等）通常不支持长文章。请根据账号能力切换 `publishType`。 |

---
> [!CAUTION]
> **严禁生成代码**：Agent 必须严格按照上述 JSON 结构构造 Payload。严禁直接调用文件操作工具尝试生成 `.ts` 或脚本文件。
