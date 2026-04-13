# 📄 百家号 视频 参数 (BaiJiaHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“百家号”分发视频动态或快讯时触发：
- **百度流量分发**：将视频推送至百度搜索及信息流。
- **互动/营销**：挂载 POI 位置、关联合集或参与正在进行的征文活动。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装百家号视频 Payload 时需遵守：
1. **标签必填逻辑**：百家号强制要求 1-6 个标签 (`tags`)。Agent 应根据标题自动补充关键词。
2. **描述字数红线**：视频描述严格限制在 1-100 字符。
3. **位置透传原则**：若需要 `location`，必须调用相关接口获取并完整透传 `raw` 数据。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (1-30 字符)。 | - |
| **`description`** | `string` | **是** | 视频描述 (1-100 字符)。 | - |
| **`tags`** | `string[]` | **是** | 视频标签 (1-6 个)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **声明**: `0`-不声明, `1`-内容由 AI 生成。 | `0` |
| `location` | `Object` | 否 | **位置信息**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `collection` | `Object` | 否 | **合集信息**: 包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |
| `activity` | `Object` | 否 | **征文活动**: 包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |

### 3.2 复杂结构说明

- **PlatformDataItem / Category**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布百家号视频动态
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["百家号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_bjh_v_001",
        "video": { "key": "bjh_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "百度生态下的视频分发技巧",
          "description": "视频描述：探讨 2026 年视频分发趋势。",
          "tags": ["科技", "分发", "运营"],
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
| **标签缺失或超限** | `tags` 数量不在 1-6 范围内。 | 增加或减少标签。 |
| **描述过长** | 文字超过 100 字符限制。 | 精简描述。 |
| **位置无效** | `location.raw` 数据格式错误。 | 重新执行 `locations` 动作获取最新 POI。 |
| **活动参与失败** | `activity` 数据与当前账号不匹配。 | 确保用户账号具备参加特定征文活动的权限。 |

---
> [!TIP]
> **SEO 强化**: 百家号视频在百度搜索端权值很高，建议 Agent 引导用户在标题中埋入搜索热词（如“怎么”、“技巧”、“避坑”）。
