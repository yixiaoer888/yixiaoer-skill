# 📄 头条号文章发布参数 (TouTiaoHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“今日头条 (头条号)”发布长文章或新闻资讯时触发。支持：
- **收益模式**：开启“投放广告赚收益”或“投放专属广告”。
- **权益申明**：申明“头条首发”获取更高流量权重。
- **地点挂载**：为文章添加地理位置信息（POI）。
- **创作者申明**：申明为自行拍摄、AI 生成或虚构演绎等。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装头条号 Payload 时需遵守：
1. **封面比例约束**：头条号封面数量推荐为 1 张或 3 张，以适配不同端显示效果。
2. **位置信息动态获取**：若用户提到地点，必须通过 `locations` 接口获取带有 `raw` 数据的 POI 对象。
3. **广告与收益策略**：默认情况下，Agent 应倾向于开启收益（`advertisement: 3`），除非用户明确表示“不投放广告”。
4. **原创首发逻辑**：若用户提到“我是原创”或“首发”，除设置创作类型外，应将 `isFirst` 设为 `true`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 50 字)。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式，最多 50000 字符)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。推荐 1 或 3 张。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `advertisement` | `number` | 否 | **投放规则**: `2`-不投放, `3`-投放广告赚收益。 | `3` |
| `isFirst` | `boolean` | 否 | 是否头条首发。设置为 `true` 可提升流量。 | `false` |
| `declaration` | `number` | 否 | **创作申明**: 1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎, 7-投资观点, 8-健康医疗。 | - |
| `location` | `Object` | 否 | 位置对象。包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布头条文章：开启收益并申明首发
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TT_ACC_123",
        "contentPublishForm": {
          "formType": "task",
          "title": "头条号自动化分发实战指南",
          "content": "<h1>实战指南</h1><p>内容正文...</p>",
          "covers": [
            { "key": "c_key_1", "width": 1280, "height": 720, "size": 256000 }
          ],
          "advertisement": 3,
          "isFirst": true,
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
| **标题敏感词拦截** | 头条对标题审查严格，包含过多夸张词汇或符号。 | 修改标题，使其更具新闻性或叙述性。 |
| **定位信息非法** | 传入了经纬度而非 `locations` 接口返回的 POI 对象。 | 必须使用 `action: "locations"` 获取的标准 POI 结构且透传 `raw`。 |
| **首发申明失败** | 账号权重不足或该文章已在其他平台检测到先发。 | 检查该账号是否具备首发权益，或确认文章是否为全网首发。 |
| **封面尺寸报错** | 封面图比例不符或分辨率过低。 | 推荐使用 16:9 或 3:2 比例的高清图。 |

---
> [!TIP]
> **多平台同步**: 在头条发布时，通常会激活“同步到微头条”，这属于蚁小二的自动扩展行为，无需额外参数。
