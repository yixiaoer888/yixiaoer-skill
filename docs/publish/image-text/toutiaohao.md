# 📄 头条号 图文 参数 (TouTiaoHao Image-Text)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [图文发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“今日头条 (头条号)”发布微头条、短动态或多图快讯时触发：
- **微头条发布**：发布类似微博的短内容。
- **生活/资讯分享**：发布 1-9 张图片并配以精炼描述。
- **创作申明**：申明内容来源以符合头条合规标准。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装头条号图文 Payload 时需遵守：
1. **HTML 话题规范**：`description` 支持 HTML。话题必须使用 `<topic text='...' raw='...'>#话题</topic>` 格式进行包裹。
2. **多图上传原则**：支持 1-9 张图片。所有 `images` 必须先调用 `upload` 动作获取系统 `key`。
3. **合规申明判定**：一点号对 AI 内容及医疗、投资等敏感领域有严格声明要求。根据内容特征准确设置 `declaration`。
4. **意图确认**：明确用户是希望“直接发布 (pubType: 1)”还是“存为草稿 (pubType: 0)”。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`description`** | `string` | **是** | 图文描述，支持 HTML 和话题标签 (`<topic>`)。最多 1000 字符。 | - |
| **`images`** | `Array` | **是** | 图片数组。使用 `OldImage[]` 结构。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **创作申明**: 1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎, 7-投资观点, 8-健康医疗。 | - |

### 3.2 OldImage 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |
| `format` | `string` | **是** | 文件格式 (如 `jpg`, `png`)。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布头条号图文动态：带话题和 AI 申明
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_tt_it_001",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>测试微头条发布 <topic text=\"头条\" raw=\"{\\\"yixiaoerId\\\":\\\"123\\\",\\\"yixiaoerName\\\":\\\"头条\\\",\\\"raw\\\":{\\\"id\\\":\\\"tt_topic_01\\\",\\\"topic\\\":\\\"头条\\\"}}\">#头条</topic></p>",
          "images": [
            { "key": "tt_img_01", "size": 1024000, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "pubType": 1,
          "declaration": 3
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
| **话题识别失败** | 标签结构不完整，缺少 `raw` 属性。 | 必须调用 `challenges` 接口获取并透传元数据。 |
| **图片数量超限** | 图片数组长度多于 9 个。 | 移除多余图片，保持在 9 张以内。 |
| **内容违规拦截** | 描述中包含联系方式或其他导流行为。 | 修改博文内容，遵守头条社区合规手册。 |
| **创作申明错误** | 选取的 `declaration` 数值不在合法范围内。 | 仅限 [3.1](#31-核心表单参数-contentpublishform) 中定义的数值。 |

---
> [!TIP]
> **流量密码**: 头条号对互动性强的内容分发较快。Agent 建议在微头条中积极内嵌“怎么看”、“大家觉得呢”等引导性金句，配合热门话题以获得更高流量推送。
