# 📄 视频号 视频 参数 (WeiXin ShiPingHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `contentPublishForm` 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“微信视频号”发布视频，且涉及以下需求时触发：
- **位置标记**：在视频中显示打卡位置。
- **商业变现**：挂载视频号带货商品（橱窗商品）。
- **激励参与**：参加视频号官方发起的创作活动。
- **存为草稿**：推送到视频号手机端/助手端的草稿箱。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装视频号 Payload 时需遵守：
1. **标题长度规范**：视频标题上限 **80 字**。
2. **描述格式化**：`description` 支持 HTML。话题和好友提醒应通过特定标签处理（若适用）。
3. **关联性数据校验**：挂载商品 (shoppingCart) 或活动 (activity) 前，必须先通过 `goods` 或 `activities` 接口获取原始 `raw` 数据。
4. **草稿判别原则**：
   - 存为视频号草稿 -> 设置 `pubType: 0` 且 `createType: 1`。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| `title` | `string` | 否 | 视频标题 (最多 80 字)。 | - |
| `short_title` | `string` | 否 | 视频短标题（主要用于搜索索引）。 | - |
| `description` | `string` | 否 | 视频描述，支持 HTML 格式。 | - |
| `horizontalCover` | `object` | 否 | **横板封面**: 使用 `OldCover` 结构。 | - |
| **`createType`** | `number` | **是** | **创建模式**: `1`-草稿，`2`-发布。 | `2` |
| **`pubType`** | `number` | **是** | **发布类型**: `0`-草稿，`1`-发布。 | `1` |
| `location` | `object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `shoppingCart` | `object` | 否 | **关联商品**: 须包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |
| `collection` | `object` | 否 | **合集信息**: 须包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |
| `activity` | `object` | 否 | **活动信息**: 须包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |

### 3.1 复杂结构说明

- **OldCover**: 包含 `key`, `size`, `width`, `height`。
- **PlatformDataItem / Activity / Goods**: 必须包含 `yixiaoerId`, `yixiaoerName` 和 **完整的 `raw` 对象**。

## 4. 执行指令示例 (Command)

```bash
# 视频号发布：带位置和定时功能
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Shipinghao"],
  "publishArgs": {
    "accountForms": [{
      "platformAccountId": "SPH_001",
      "video": {"key": "v_key", "width": 1080, "height": 1920, "size": 1024000},
      "cover": {"key": "c_key", "width": 1080, "height": 1920, "size": 300000},
      "coverKey": "c_key",
      "contentPublishForm": {
        "formType": "task",
        "title": "记录生活碎片",
        "pubType": 1,
        "createType": 2,
        "location": { "yixiaoerId": "gz_001", "yixiaoerName": "广州塔", "raw": {...} }
      }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **禁止发布该类目** | 视频内容/挂载商品不符合视频号当前类目权限。 | 检查账号在视频号后台的准入类目。 |
| **位置标记显示失败** | `location.raw` 数据格式不完整或 POI ID 已失效。 | 重新执行 `locations` 动作获取。 |
| **定时发布失败** | 时间设置过近（通常要求 > 1 小时）或过远。 | 建议设置在当前时间 2 小时后。 |
| **商品链接不可用** | 关联的小程序商品路径已发生变化。 | 引导用户在视频号橱窗重新确认商品后再查。 |

---
> [!IMPORTANT]
> **私域流量联动**：视频号是微信生态的核心环节，发布时建议开启“显示到公众号首页”等开关（若接口支持），以最大化分发效果。
