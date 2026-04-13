# 📄 搜狐号 视频 参数 (SouHuHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“搜狐号”平台分发视频动态、行业资讯或生活 VLOG 时触发：
- **搜狐信息流分发**：利用搜狐网及其手机端的庞大用户群进行触达。
- **原创品牌建设**：标注原创身份并提供准确分类以获取收录。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装搜狐号视频 Payload 时需遵守：
1. **分类级联必填**：搜狐号分类对搜索引擎抓取具有极高权值。必须通过 `categories` 接口获取并在 `category` 数组中透传 `raw` 原始数据。
2. **标题与描述红线**：标题需控制在 5-72 字，描述需在 5-200 字。Agent 应校验输入字数。
3. **原创申明判定**：搜狐号强制要求 `declaration` 字段数值。根据内容特征设置（如自行拍摄设为 `2`）。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (5-72 字符)。 | - |
| **`description`** | `string` | **是** | 视频描述 (5-200 字符)。 | - |
| **`tags`** | `string[]` | **是** | 视频标签。 | - |
| **`category`** | `Array` | **是** | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`declaration`** | `number` | **是** | **创作申明**: `0`-无, `1`-引用, `2`-自行拍摄, `3`-AI创作, `4`-虚构。 | `0` |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `id`, `text`, `raw` 元数据。

## 4. 执行指令示例 (Command)

```bash
# 发布搜狐号原创实拍视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["搜狐号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_sh_v_001",
        "video": { "key": "sh_v_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐号：自媒体时代的流量常青树",
          "description": "探讨搜狐号在 2026 年的流量分发逻辑。",
          "tags": ["运营", "搜狐", "经验"],
          "category": [{ "id": "1", "text": "科技", "raw": {...} }],
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
| **标题太短或过长** | 字数不在 5-72 范围内。 | 修改标题内容长度。 |
| **原创申明缺失** | 未传 `declaration` 数值。 | 搜狐号视频通道对声明通常为强校验，请补全。 |
| **分类数据错误** | `category.raw` 缺失或非原始对象。 | 必须实时从接口获取并透传。 |
| **标签解析失败** | 标签含特殊字符或敏感词汇。 | 移除不合规标签。 |

---
> [!TIP]
> **SEO 全渠道分发**: 搜狐内容在搜索引擎中权重极高。Agent 建议标题尽量贴合用户的搜索意图（如“怎么做的”、“最新技巧”），并配以专业化的分类。
