# 📄 视频号 图文 参数 (WeChat Video Account Image-Text)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [图文发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“微信视频号”发布图文笔记、短动态或朋友圈风格内容时触发：
- **精品图文**：发布高质量多图（1-9 张）笔记，支持标题和长描述。
- **互动/挂载**：挂载 POI 位置、引用话题、关联合集或设置背景音乐。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装视频号图文 Payload 时需遵守：
1. **HTML 话题内嵌**：`description` 支持 HTML。话题必须使用 `<topic text='...' raw='...'>#话题</topic>` 格式进行包裹。
2. **三位一体资源**：`images` 必须上传至 OSS，且每个成员需具备完整的 `OldImage` 结构。
3. **复杂对象透传**：对于 `location`, `music`, `collection` 等字段，必须通过相应接口获取并完整透传 `raw` 原始数据，严禁重构对象属性。
4. **发布模式锁定**：确认用户是希望“存为平台草稿 (pubType: 0)”还是“直接推送至视频号 (pubType: 1)”。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`images`** | `Array` | **是** | 图片数组 (1-9 张)。使用 `OldImage[]` 结构。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `title` | `string` | 否 | 笔记标题。建议包含核心关键词。 | - |
| `description` | `string` | 否 | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `location` | `Object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `music` | `Object` | 否 | **背景音乐**: 使用 `MusicItem` 结构。 | - |
| `collection` | `Object` | 否 | **合集信息**: 使用 `Collection` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 复杂结构说明

- **OldImage**: 包含 `key`, `size`, `width`, `height`, `format`。
- **PlatformDataItem / Collection**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。
- **MusicItem**: 包含 `yixiaoerId`, `yixiaoerName`, `duration`, `playUrl`, `artist` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布视频号精品图文：带音乐和话题
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["视频号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号图文实操案例",
          "description": "<p>今日分享一个实操案例 <topic text=\"自媒体\" raw=\"{\\\"id\\\":\\\"wx_01\\\",\\\"name\\\":\\\"自媒体\\\"}\">#自媒体</topic></p>",
          "images": [
            { "key": "sph_img_01", "size": 307200, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "music": { "yixiaoerId": "m_1", "yixiaoerName": "轻快背景音", "duration":60, "playUrl":"...", "raw": {...} },
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
| **话题标签不生效** | `<topic>` 标签中的 `raw` 数据过期或格式错误。 | 必须重新调用 `challenges` 接口获取最新的话题元数据。 |
| **封面/图片加载失败** | `images` 中的 `key` 缺失或关联资源未上传。 | 确认所有图片已通过 `upload` 动作并正确透传 key。 |
| **音乐挂载异常** | 目标音乐 ID 在当前账号下不可用或版权受限。 | 检查该视频号账号在平台上的音乐可用库。 |
| **位置偏移或无效** | `location.raw` 未被视频号服务端识别。 | 建议在 Agent 辅助下重新筛选最匹配的 POI。 |

---
> [!TIP]
> **微信社交联动**: 视频号图文与朋友圈具有直接的社交分发属性。Agent 建议描述内容具有亲和力，通过 HTML 话题锚点引导用户进入更高流量的相关社群。
