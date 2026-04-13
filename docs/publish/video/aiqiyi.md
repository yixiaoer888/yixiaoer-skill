# 📄 爱奇艺 视频 参数 (AiQiYi Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“爱奇艺”分发视频内容时触发：
- **影视综分发**：发布高质量的长短视频、独家剪辑或行业动态。
- **内容合规声明**：标注内容是否由 AI 生成、虚构演绎或取材网络。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装爱奇艺视频 Payload 时需遵守：
1. **分类级联必填**：爱奇艺对分类要求严格。必须通过 `categories` 接口获取并在 `category` 数组中透传完整的 `raw` 对象。
2. **字数与格式**：标题限制在 1-50 字符，描述限制在 1-500 字符。Agent 应确保内容不超限。
3. **创作申明判定**：根据视频来源准确设置 `declaration`。例如，识别到 AI 痕迹应设为 `1`。
4. **资源引用**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (1-50 字符)。 | - |
| **`description`** | `string` | **是** | 视频描述 (1-500 字符)。 | - |
| **`category`** | `Array` | **是** | 视频分类。使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`createType`** | `number` | **是** | **内容类型**: `1`-原创, `2`-转载。 | `1` |
| **`pubType`** | `number` | **是** | **发布类型**: `0`-草稿, `1`-直接发布。 | `1` |
| `tags` | `string[]` | 否 | 视频标签 (最多 10 个)。 | - |
| `declaration` | `number` | 否 | **声明**: `0`-无, `1`-AI 生成, `2`-虚构演绎, `3`-取材网络。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。支持父子级级联。

## 4. 执行指令示例 (Command)

```bash
# 发布爱奇艺原创影视视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["爱奇艺"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_aqy_vid_001",
        "video": { "key": "aqy_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 120 },
        "cover": { "key": "aqy_c_key", "size": 102400, "width": 800, "height": 600 },
        "coverKey": "aqy_c_key",
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 影视产业深度复盘",
          "description": "本文深度探讨了影视产业的未来走向。",
          "tags": ["影视", "科技", "解析"],
          "category": [
            { "yixiaoerId": "cat_id_001", "yixiaoerName": "影视", "raw": {...} }
          ],
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
| **分类未匹配** | `category` 数组中的 ID 与当前账号分区不一致。 | 必须调用 `action: "categories"` 重新获取。 |
| **标题或描述过长** | 超出了 50/500 的字符上限。 | 自动截断或引导用户精简文案。 |
| **标签过多报错** | `tags` 成员超过 10 个。 | 请移除次要标签。 |
| **定时发布失效** | `scheduledTime` 参数格式不正确。 | 提供标准 Unix 秒级时间戳。 |

---
> [!TIP]
> **爱艺奇生态**: 爱奇艺对高质量的原创视频有长尾分发优势。Agent 建议标题尽量吸引人，并配以高清封面以提升点击率。
