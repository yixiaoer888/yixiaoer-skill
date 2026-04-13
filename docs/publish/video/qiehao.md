# 📄 企鹅号 视频 参数 (QiEHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“企鹅号 (腾讯内容开放平台)”发布视频资讯、深度解析或短视频动态时触发：
- **全平台分发**：将视频推送至腾讯视频、腾讯新闻、QQ 浏览器、天天快报等渠道。
- **合规申明**：标注 AI 生成、虚构演绎、个人观点或取材网络。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装企鹅号视频 Payload 时需遵守：
1. **分类级联必填**：企鹅号对频道划分极其细致。必须通过 `categories` 接口获取并在 `category` 数组中透传 `raw` 原始对象。
2. **标签规范**：强制要求 1-5 个 `tags`。Agent 应从视频描述中精准提取核心动量。
3. **内容安全声明**：根据视频真实来源设置 `declaration`。例如，若检测到 AI 痕迹应设为 `1`。
4. **资源引用协议**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述内容。 | - |
| **`tags`** | `string[]` | **是** | 视频标签 (1-10 个)。 | - |
| **`category`** | `Array` | **是** | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`pubType`** | `number` | **是** | **发布类型**: `0`-草稿, `1`-直接发布。 | `1` |
| `declaration` | `number` | 否 | **声明**: `0`-无, `1`-AI 生成, `2`-个人观点, `3`-剧情演绎, `4`-取材网络, `5`-旧闻。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。支持多级嵌套。

## 4. 执行指令示例 (Command)

```bash
# 发布企鹅号原创深度解析视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["企鹅号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_qh_v_001",
        "video": { "key": "qh_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 科技趋势深度复盘",
          "description": "探讨 2026 年人工智能对实体经济的影响。",
          "tags": ["科技", "AI", "解析"],
          "category": [{ "yixiaoerId": "cat_id_001", "yixiaoerName": "生活", "raw": {...} }],
          "declaration": 2,
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
| **分类 ID 不合法** | `category.raw` 数据格式错误或已过期。 | 必须实时调用 `categories` 接口并完整透传原始对象。 |
| **标签解析失败** | 标签中包含了腾讯禁止的敏感词。 | 优化并精简标签列表。 |
| **封面上传失败** | `coverKey` 缺失。 | 请确保封面图片已成功通过 `upload` 并获得 key。 |
| **申明与内容不符** | 实际为 AI 合成但未勾选 `declaration: 1`。 | 准确设置合规申明数值。 |

---
> [!TIP]
> **全场景覆盖**: 企鹅号是打通腾讯多社交分发的单一入口，Agent 建议标题多体现“新鲜感”和“专业性”，以融入腾讯新闻等不同端的信息流特征。
