# 📄 知乎 视频 参数 (Zhihu Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“知乎”分发硬核科普、专业解析、生活 VLOG 或好物种草视频时触发：
- **专业知识分发**：利用知乎的高知社区属性发布深度视频内容。
- **原创权益护航**：标注原创身份并提供准确的 AI 生成申明。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装知乎视频 Payload 时需遵守：
1. **分类级联必填**：知乎分区严谨。必须通过 `categories` 接口获取并在 `category` 数组中透传完整的 `raw` 原始对象。
2. **描述深度原则**：知乎对内容质量敏感。Agent 建议描述文字 (`description`) 应包含核心观点，而非简单的引导词。
3. **原创与存储模式**：必须明确 `createType` (原创/转载) 和 `pubType` (草稿/发布)。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`description`** | `string` | **是** | 视频描述内容 (1-500 字符)。 | - |
| **`category`** | `Array` | **是** | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`createType`** | `number` | **是** | **内容类型**: `1`-原创, `2`-转载。 | `1` |
| **`pubType`** | `number` | **是** | **模式**: `0`-草稿, `1`-直接发布。 | `1` |
| `title` | `string` | 否 | 视频标题 (1-50 字符)。 | - |
| `declaration` | `number` | 否 | **声明**: `0`-不声明, `2`-AI 生成。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName`, `raw` 原生对象。

## 4. 执行指令示例 (Command)

```bash
# 发布知乎原创深度解析视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Zhihu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_01",
        "video": { "key": "zh_v_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 120 },
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎内容的分发逻辑与博弈",
          "description": "探讨 2026 年知乎视频权重的变化与创作者对策。",
          "category": [{ "yixiaoerId": "cat_id_01", "yixiaoerName": "科技", "raw": {...} }],
          "createType": 1,
          "pubType": 1,
          "declaration": 0
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
| **分类未锁定** | `category.raw` 数据格式不满足知乎 DTO 要求。 | 确保调用 `categories` 接口并完整透传。 |
| **描述太短** | description 内容少于 10 个字符。 | 建议增加视频核心观点总结。 |
| **原创申明错误** | 勾选了转载但未提供正确的来源标识。 | 准确设置 `createType` 数值。 |
| **封面画质不佳** | 封面模糊被知乎系统降权。 | 建议使用 1080P 素材裁剪。 |

---
> [!TIP]
> **硬核知识社区**: 知乎用户偏好逻辑清晰、具备洞见的内容。Agent 建议标题采用“探究式”，并在描述中简述视频解决的问题，以利用知乎的长尾搜索优势。
