# 📄 网易号 视频 参数 (WangYiHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“网易号”平台分发网易新闻相关的资讯、评论或原创专栏视频时触发：
- **新闻分发**：触达网易新闻 App 的高价值用户群。
- **原创权益护航**：标注原创类型并配置详细的创作申明。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装网易号视频 Payload 时需遵守：
1. **分类级联必填**：网易号对频道分区有严格要求。必须通过 `categories` 接口获取并在 `category` 数组中透传 `raw` 原始数据。
2. **双重原创申明**：包含 `createType` (创作类型) 和 `declaration` (申明)。Agent 应根据用户提示词同时锁定这两个字段的逻辑一致性。
3. **标签强索引**：必须提供 `tags` 字符串数组作为分发引擎的抓取依据。
4. **资源引用规范**：必须通过 `upload` 动作产生视频 key 及封面 key。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`tags`** | `string[]` | **是** | 视频标签。建议 1-10 个。 | - |
| **`category`** | `Array` | **是** | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **创作申明**: 1-AI生成, 2-个人原创, 3-取材网络, 4-虚构演绎。 | - |
| `createType` | `number` | 否 | **创作类型**: `0`-非原创, `1`-原创。 | `0` |
| `description` | `string` | 否 | 视频描述内容。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布网易号个人原创视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["网易号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_wy_v_02",
        "video": { "key": "wy_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "网易生态视频分发实战策略",
          "description": "探讨网易号长短视频的融合运营技巧。",
          "tags": ["自媒体", "网易", "解析"],
          "category": [{ "yixiaoerId": "cat_id_01", "yixiaoerName": "科技", "raw": {...} }],
          "createType": 1,
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
| **级联分类不匹配** | `category.raw` 数据格式不满足网易接口要求。 | 必须实时调用 `categories` 获取最新原始对象。 |
| **申明与类型冲突** | 勾选了原创 `createType: 1` 但申明为 `declaration: 3` (网络取材)。 | 修正合规逻辑，确保申明的一致性。 |
| **封面上传失败** | `coverKey` 无效。 | 请重新执行 `upload` 动作。 |
| **发布频率过快** | 短时间内对同一网易号连续发布。 | 建议增加时间间隔。 |

---
> [!TIP]
> **网易态度生态**: 网易用户偏好有深度、有态度的内容。Agent 建议视频标题应具备观点性，并在描述中引入用户关心的社会或科技话题。
