# 📄 凤凰网 视频 参数 (Fengwang Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前， you **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“凤凰网 (iFeng)”平台发布严肃新闻、时事评论、社会纪实或高端专题视频时触发：
- **媒体同步**：将权威资讯同步到凤凰新闻客户端。
- **深度专题**：发布带有精准行业分区的高端长短视频。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装凤凰网视频 Payload 时需遵守：
1. **分类强校验**：凤凰网对标签和分类有明确的频道划分要求。必须通过 `categories` 接口获取并透传 `raw` 对象。
2. **标题与标签规范**：标题应保持专业，`tags` 是系统分发的重要索引，不可缺失。
3. **内容正向引导**：凤凰网社区调性偏向严肃，Agent 应引导内容符合平台格调。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数 definition (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述。 | - |
| **`tags`** | `string[]` | **是** | 视频标签。建议 2-5 个核心词。 | - |
| **`category`** | `Array` | **是** | 视频分类。使用 `CascadingPlatformDataItem[]` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `id`, `text`, `raw` 元数据。

## 4. 执行指令示例 (Command)

```bash
# 发布凤凰网社会观察视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Fengwang"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "IFENG_ACC_01",
        "video": { "key": "ifeng_v_1", "size": 2048000, "width": 1920, "height": 1080, "duration": 150 },
        "contentPublishForm": {
          "formType": "task",
          "title": "全球地缘政治格局 2026 深度研判",
          "description": "深度解析：未来十年全球地缘政治的核心动量。",
          "tags": ["时事", "评论", "深度"],
          "category": [
            { "id": "1", "text": "社会", "raw": {...} }
          ]
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
| **分类未命中** | `category.raw` 数据与当前账号频道权限不匹配。 | 重新执行 `categories` 接口以刷新权限缓存。 |
| **标签解析失败** | `tags` 包含敏感、偏激或非法字符。 | 移除不合规标签。 |
| **视频审核延迟** | 凤凰网对内容有较长的人工审核流程。 | 请稍后在控制台查看发布进度。 |
| **描述字数过短** | 描述内容少于 10 个字符。 | 请补充视频内容的背景介绍。 |

---
> [!TIP]
> **高端受众分发**: 凤凰网在精英阶层具有极高覆盖率。Agent 建议视频标题采用专业、克制的表达方式，并配以精细化的频道分类以触达目标读者。
