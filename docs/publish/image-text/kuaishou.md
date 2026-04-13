# 📄 快手 图文 参数 (Kuaishou Image-Text)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [图文发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“快手”发布图文笔记或动态，且涉及以下特有需求时触发：
- **精品动态**：发布多图动态（支持 HTML 格式描述）。
- **互动增强**：挂载 POI 地理位置、内嵌 #话题 标签、关联背景音乐。
- **发布控制**：设置公开、私密或好友可见，或安排定时发布。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装快手图文 Payload 时需遵守：
1. **HTML 话题内嵌**：`description` 支持 HTML 标签。话题必须使用 `<topic text='...' raw='...'>#话题</topic>` 格式。
2. **资源先行上传**：所有 `images` 必须先调用 `upload` 动作获取 `key`。
3. **数据完整性**：对于 `location`, `music`, `collection` 等复杂字段，**严禁手动构造**，必须调用相关接口获取并完整透传 `raw` 对象。
4. **可见性校验**：默认为公开 (0)，若涉及个人隐私，Agent 应提醒用户确认是否需要设置为私密 (1)。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`description`** | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| **`images`** | `Array` | **是** | 图片数组。使用 `OldImage[]` 结构。 | - |
| **`visibleType`** | `number` | **是** | **可见类型**: `0`-公开, `1`-私密, `3`-好友可见。 | `0` |
| `location` | `Object` | 否 | **位置信息**: 使用 `PlatformDataItem` 结构。 | - |
| `music` | `Object` | 否 | **音乐信息**: 使用 `MusicItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |
| `collection` | `Object` | 否 | **合集信息**: 使用 `Category` 结构。 | - |

### 3.2 复杂结构说明

- **OldImage**: 包含 `key`, `size`, `width`, `height`, `format`。
- **PlatformDataItem / Category / Collection**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。
- **MusicItem**: 包含 `yixiaoerId`, `yixiaoerName`, `duration`, `playUrl`, `artist` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布快手图文动态：带位置和话题
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["快手"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_ks_it_001",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>快手动态测试 <topic text=\"快手\" raw=\"{\\\"id\\\":\\\"ks_01\\\",\\\"name\\\":\\\"快手\\\"}\">#快手</topic></p>",
          "images": [
            { "key": "ks_img_001", "size": 204800, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "location": { "yixiaoerId": "poi_001", "yixiaoerName": "快手总部", "raw": {...} },
          "visibleType": 0
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
| **描述内容超限** | 字数超过 1000 字符限制。 | 精简描述内容。 |
| **话题标签不生效** | `<topic>` 标签格式不正确或缺少 `raw` 属性。 | 严格按照规范拼装 HTML。 |
| **封面/图片获取失败** | `images` 数组中的 `key` 无效或资源已过期。 | 重新执行 `upload` 动作获取最新的 `key`。 |
| **位置/音乐加载异常** | `raw` 数据过期或与该账号平台不匹配。 | 重新调用对应的 `get-*` action 获取最新数据。 |

---
> [!TIP]
> **快手流量建议**: 快手社区倾向于生活化和互动性强的内容，建议在 `description` 中多使用表情符号并积极引导评论。
