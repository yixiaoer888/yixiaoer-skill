# 📄 车家号 文章 参数 (Chejiahao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“车家号”平台发布汽车相关的专业文章、深度评测或行业资讯时触发：
- **汽车评测**：发布包含多尺寸封面的深度好文。
- **创作申明**：申明原创、首发或原创首发以获取更高推荐权重。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装车家号 Payload 时需遵守：
1. **多重封面规范**：车家号通常要求同时提供横版封面 (`covers`) 和竖版封面 (`verticalCovers`)。Agent 应确保两组封面均已通过 `upload` 动作转换。
2. **创作类型匹配**：明确用户的创作意图。`type` 字段必须与用户的“原创/首发”描述严格对应。
3. **内容完整性**：标题和正文 (HTML) 必须同步填入 `contentPublishForm`，确保账号级配置覆盖根结构。
4. **资源标识**：所有图片的 `key` 必须是系统内部 OSS Key。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | **横版封面**列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`verticalCovers`** | `Array` | **是** | **竖版封面**列表。结构同 OldCover。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `type` | `number` | 否 | **创作类型**: `1`-原创, `3`-首发, `13`-原创首发。 | `1` |
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
# 发布车家号“原创首发”文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["车家号"],
  "publishArgs": {
    "content": "<h1>深度评测：2026 款豪华轿车</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_cjh_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "深度评测：2026 款豪华轿车性能解析",
          "content": "<h1>深度评测</h1><p>正文内容...</p>",
          "covers": [{ "key": "cjh_h_1", "size": 150000, "width": 800, "height": 600 }],
          "verticalCovers": [{ "key": "cjh_v_1", "size": 150000, "width": 600, "height": 800 }],
          "type": 13,
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
| **竖版封面缺失** | 未提供 `verticalCovers` 数组。 | 车家号特定分区强制要求竖版封面，请补全。 |
| **首发申明失败** | 账号不具备首发权限或该文章已在其他平台发过。 | 检查账号权益，确保文章为全网首发。 |
| **HTML 渲染异常** | 包含车家号不支持的复杂 HTML 标签或外链。 | 简化 HTML 结构，移除不可用标签。 |
| **封面比例报错** | 宽高比不符合横屏或竖屏的特定标准。 | 建议横版 4:3 或 16:9，竖版 3:4。 |

---
> [!TIP]
> **汽车流量策略**: 车家号对评测类内容有较高的收录优先级，建议 Agent 在 `content` 中尽可能多展示车辆细节图片并保持专业性。
