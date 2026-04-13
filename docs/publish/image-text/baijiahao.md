# 📄 百家号 图文 参数 (BaiJiaHao Image-Text)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [图文发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“百家号”发布动态、快讯或多图生活分享时触发：
- **动态分发**：利用百度搜索和信息流的强大分发能力增加曝光。
- **互动话题**：内嵌 #话题 标签、挂载 POI 地理位置以触达精准人群。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装百家号图文 Payload 时需遵守：
1. **标题字数红线**：百家号动态标题建议在 1-20 字符。Agent 应主动校验或从描述中提取精简标题。
2. **HTML 话题内嵌**：`description` 支持 HTML。话题必须使用 `<topic text='...' raw='...'>#话题</topic>` 格式进行包裹。
3. **封面逻辑**：百家号图文强制要求 `cover` 封面对象。必须通过 `upload` 动作获取 Key。
4. **数据透传要求**：对于 `location` 等字段，必须通过相关接口获取并完整透传 `raw` 原始数据。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 标题 (1-20 字符)。 | - |
| **`description`** | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| **`cover`** | `Object` | **是** | 封面图对象。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| **`declaration`** | `number` | **是** | **创作声明**: `0`-不声明, `1`-内容由 AI 生成。 | `0` |
| `location` | `Object` | 否 | **位置信息**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

### 3.3 PlatformDataItem 定义 (用于 location)
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `id` | `string` | **是** | 位置 ID。 |
| `text` | `string` | **是** | 项名称/文本。 |
| `raw` | `Object` | **是** | 原始数据对象。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布百家号图文动态：带话题和 AI 申明
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["百家号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_bjh_it_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号自动化发布实测",
          "description": "<p>内容摘要 <topic text=\"科技\" raw=\"{\\\"yixiaoerId\\\":\\\"123\\\",\\\"yixiaoerName\\\":\\\"科技\\\",\\\"raw\\\":{\\\"id\\\":\\\"bjh_topic_01\\\",\\\"topic\\\":\\\"科技\\\"}}\">#科技</topic></p>",
          "cover": { "key": "bjh_img_1", "size": 102400, "width": 800, "height": 600 },
          "pubType": 1,
          "declaration": 1
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
| **标题超出限制** | 标题字符多于 20 个。 | 请精简标题，突出核心动量。 |
| **话题 ID 无效** | `topic.raw` 中的 ID 与百家号平台不匹配。 | 必须调用 `challenges` 接口获取并透传最新元数据。 |
| **封面获取失败** | `cover.key` 为空或不是 OSS 标识。 | 确保已执行 `upload` 动作并正确引用产生的 key。 |
| **AI 申明漏选** | 内容被 AI 识别系统拦截。 | 建议设置 `declaration: 1` 以符合平台合规要求。 |

---
> [!TIP]
> **百度权重获取**: 百家号动态在百度搜索结果中展现极佳。Agent 建议标题多包含搜索热词，并积极内嵌平台热门话题以获得冷启动流量。
