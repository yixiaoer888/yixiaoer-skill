# 📄 新浪微博 图文 参数 (Sina Weibo Image-Text)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [图文发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“新浪微博”发布碎碎念、日常动态或多图快讯时触发：
- **短图文动态**：发布 1-9 张图片配以描述文字。
- **互动/热点**：内嵌 #话题 标签、关联 POI 地理位置。
- **发布预约**：设置未来某个时间点自动发布。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装新浪微博图文 Payload 时需遵守：
1. **HTML 话题规范**：微博博文支持 HTML 映射。热门话题必须使用 `<topic text='...' raw='...'>#话题</topic>` 格式进行包裹。
2. **多图上传验证**：微博最多支持 9 张图片，每张图片必须先通过 `upload` 动作转换为系统 `key`。
3. **字数与格式**：微博正文（`description`）建议控制在 2000 字以内，Agent 应主动识别用户是否需要将超长内容转换为“长文章”。
4. **地理位置透传**：若用户提到地点，必须通过 `locations` 获取带有 `raw` 数据的对象进行闭环透传。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`description`** | `string` | **是** | 微博博文内容。支持 HTML (`<p>`, `<topic>`)。最多 2000 字符。 | - |
| **`images`** | `Array` | **是** | 图片数组 (1-9 张)。使用 `OldImage[]` 结构。 | - |
| `location` | `Object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 复杂结构说明

- **OldImage**: 包含 `key`, `size`, `width`, `height`, `format`。
- **PlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布微博多图动态：带话题和地理位置
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["新浪微博"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WB_ACC_123",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>今日好心情！ <topic text=\"北京\" raw=\"{\\\"id\\\":\\\"bj_001\\\",\\\"name\\\":\\\"北京\\\"}\">#北京</topic></p>",
          "images": [
            { "key": "img_key_1", "size": 102400, "width": 800, "height": 800, "format": "jpg" }
          ],
          "location": { "yixiaoerId": "poi_123", "yixiaoerName": "天安门", "raw": {...} }
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
| **图片上传超限** | 图片数量超过 9 张或单张图片过大。 | 精简图片数量，确保单张小于 5MB。 |
| **话题未被识别** | 话题格式不符合 `<topic>` 规范，或 `raw` 数据有误。 | 检查 HTML 闭合情况及 `raw` 对象完整性。 |
| **位置偏移或非法** | `location.raw` 数据格式已改变。 | 重新执行 `action: "locations"` 获取最新数据。 |
| **发布过于频繁** | 微博对同一账号的短时间连续发布有冷却期。 | 建议设置 `scheduledTime` 间隔发布或稍后再试。 |

---
> [!TIP]
> **微博热搜策略**: 微博是公域流量平台。Agent 在生成 `description` 时，应建议用户关联当前微博热搜榜中的相关话题以提升博文权重。
