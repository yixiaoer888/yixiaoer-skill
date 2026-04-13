# 📄 一点号 视频 参数 (Yidianhao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“一点资讯 (一点号)”平台分发视频动态、资讯快讯或生活见闻时触发：
- **全网信息流投放**：将视频推送至一点资讯及相关合作端。
- **原创属性申明**：标注原创身份并提供准确的合规申明。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装一点号视频 Payload 时需遵守：
1. **分类级联必填**：一点号对内容精准分类有强校验。必须通过 `categories` 接口获取并在 `category` 数组中透传完整的 `raw` 元数据。
2. **原创与申明双核**：必须同时锁定 `type` (0-非原创, 1-原创) 和 `declaration` (申明)。Agent 需确保两者在逻辑上不冲突。
3. **资源引用规范**：确保视频资源的 `key` 已通过 `upload` 动作获取。
4. **意图确认**：若内容包含 AI 创作，Agent 应主动提醒并设置 `declaration: 4`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述内容。 | - |
| **`tags`** | `string[]` | **是** | 视频标签。建议 1-10 个。 | - |
| **`category`** | `Array` | **是** | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`type`** | `number` | **是** | **原创类型**: `0`-非原创, `1`-原创。 | `1` |
| `declaration` | `number` | 否 | **声明**: `3`-素材取自网络, `4`-内容由 AI 生成, `5`-虚构情节。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `id`, `text`, `raw` 原始对象。

## 4. 执行指令示例 (Command)

```bash
# 发布一点号原创 AI 视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Yidianhao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YD_ACC_01",
        "video": { "key": "yd_v_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "2026：人工智能进化的分水岭",
          "description": "探讨 AI 如何重塑我们的日常生活。",
          "tags": ["AI", "未来", "科技"],
          "category": [{ "id": "1", "text": "社会", "raw": {...} }],
          "type": 1,
          "declaration": 4
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
| **分类 ID 未命中** | `category.raw` 数据格式不满足一点号接口规范。 | 必须实时重新获取并完整透传。 |
| **申明与原创冲突** | `type` 为 1 但 `declaration` 描述为网络素材。 | 逻辑自洽检查，修正申明项。 |
| **标签解析失败** | 解析后的 tags 字符串超长。 | 减少标签数量或精简关键词长度。 |
| **封面上传失败** | 缺失 `coverKey`。 | 请执行 `upload` 获取图片 key。 |

---
> [!TIP]
> **精准分发模型**: 一点号具有极强的兴趣图谱匹配能力。Agent 建议分类选择应尽可能下钻，并配合垂直化的描述内容以触达精准读者。
