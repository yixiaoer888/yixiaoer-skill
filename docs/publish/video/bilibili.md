# 📄 哔哩哔哩 视频 参数 (BiLiBiLi Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `contentPublishForm` 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“哔哩哔哩 (B站)”分发视频，且涉及以下特有需求时触发：
- **社区互动**：设置视频分区 (Category)、打上 #话题 标签、关联合集。
- **创作申明**：申明视频为自制原创或转载，并提供原文链接。
- **内容合规**：标注是否为 AI 合成、包含危险行为或个人观点等。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 B 站 Payload 时需遵守：
1. **级联分类校验**：B 站投稿对分类要求极严。必须包含父子级 DTO。**严禁手动指定 ID**，必须通过 `categories` 接口获取并在 `category` 数组中透传 `raw` 数据。
2. **标签规范化**：`tags` 应包含视频的核心关键词（1-10 个）。若用户未提供，Agent 应根据标题和描述自动提取。
3. **申明判定**：若 `createType` 设为 `2` (转载)，必须强制要求用户提供 `contentSourceUrl`。
4. **资源引用**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (最多 80 字符)。 | - |
| `description` | `string` | 否 | 视频描述 (最多 2000 字符)。 | - |
| **`tags`** | `string[]` | **是** | 视频标签数组 (1-10 个)。 | - |
| **`category`** | `Array` | **是** | 视频分类。使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`createType`** | `number` | **是** | **类型**: `1`-自制, `2`-转载。 | `1` |
| **`pubType`** | `number` | **是** | **发布类型**: `0`-草稿, `1`-直接发布。 | `1` |
| `declaration` | `number` | 否 | **创作者申明**: `0`-无, `1`-AI合成, `2`-危险行为, `3`-仅供娱乐, `4`-引人不适, `5`-理性消费, `6`-个人观点。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `contentSourceUrl` | `string` | 否 | 原文 URL (当 `createType` 为 2 时必填)。 | - |
| `collection` | `Object` | 否 | 合集信息。包含 `yixiaoerId`, `yixiaoerName`。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及 **`raw`** 对象。必须按照“父分类 -> 子分类”的顺序排列。
- **Collection**: 必须包含 `yixiaoerId`, `yixiaoerName`。

## 4. 执行指令示例 (Command)

```bash
# 发布 B 站自制原创视频到生活区
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Bilibili"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BILIBILI_ACC_123",
        "video": { "key": "v_key_001", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "cover": { "key": "c_key_001", "size": 300000, "width": 1920, "height": 1080 },
        "coverKey": "c_key_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "我的 B 站生活 VLOG",
          "tags": ["生活", "VLOG", "摄影"],
          "category": [
            { "yixiaoerId": "cat_parent_id", "yixiaoerName": "生活", "raw": {...} },
            { "yixiaoerId": "cat_child_id", "yixiaoerName": "日常", "raw": {...} }
          ],
          "createType": 1,
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
| **分类未选择或非法** | `category` 数组为空或 ID 不匹配 B 站当前分区。 | 必须调用 `action: "categories"` 并确保包含完整的 `raw` 数据。 |
| **标签数量不符** | 标签少于 1 个或多于 10 个。 | 调整标签数量至合规范围。 |
| **转载链接缺失** | 申明为转载但未填写 `contentSourceUrl`。 | 补全原文链接。 |
| **标题含有关键词** | 包含 B 站违禁词或政治敏感词。 | 修改标题，参考 B 站社区准则。 |

---
> [!TIP]
> **分区匹配技巧**: B 站不同的分区对内容审核尺度不同，建议 Agent 在获取分类时，向用户确认最精准的子分区以获得更好的推荐效果。
