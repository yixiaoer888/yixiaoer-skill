# 📄 小红书图文发布参数 (Xiaohongshu Image-Text)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在进入细分参数前，你 **必须** 已经阅读并理解了 [图文发布通用索引](./index.md) 中定义的 Payload 根结构。本文档仅描述 `contentPublishForm` 内部的平台差异化字段。

## 1. 触发场景 (Trigger)

当用户明确要在“小红书”发布图文笔记，且涉及以下特有需求时触发：
- **精品笔记**：发布高质量多图笔记（1-9 张）、配置笔记标题。
- **互动/搜索**：挂载 POI 地理位置、设置笔记话题（正文内嵌或专属标签）。
- **可见性控制**：设置公开、私密或仅好友可见。

## 2. 交互协议 (Interactive Protocol)

Agent 在构造小红书图文 Payload 时需遵守：
1. **标题长度红线**：小红书标题严格限制在 **20 字** 以内。Agent 应主动对超长标题进行截断。
2. **正文 HTML 规范**：`description` 支持 HTML。话题应使用 `<topic text='...' raw='...'>#话题</topic>` 格式内嵌。
3. **数据透传要求**：对于 `location`, `collection`, `music` 等，必须完整透传接口返回的 `raw` 原始数据。
4. **可见性提示**：默认为公开 (0)，若涉及隐私内容需提示用户切换为私密 (1)。

## 3. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | **笔记标题** (最多 20 字)。建议包含关键词以提升搜索权重。 | - |
| **`description`** | `string` | **是** | **笔记正文**，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| **`images`** | `Array` | **是** | 图片数组，支持 1-9 张。使用 `OldImage[]` 结构。 | - |
| `location` | `Object` | 否 | **位置信息**: 使用 `PlatformDataItem` 结构。 | - |
| `music` | `Object` | 否 | **音乐信息**: 使用 `MusicItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |
| `collection` | `Object` | 否 | **合集信息**: 使用 `Collection` 结构。 | - |
| **`visibleType`** | `number` | **是** | **可见类型**: `0`-公开, `1`-私密, `3`-好友可见。 | `0` |

### 3.1 复杂对象结构 (Data Schemas)

- **OldImage**: 必须包含 `key`, `size`, `width`, `height`, `format`。
- **PlatformDataItem / Collection / MusicItem**: 所有统一结构必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 小红书图文发布：带话题和可见性设置
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["Xiaohongshu"],
  "publishArgs": {
    "accountForms": [{
      "platformAccountId": "XHS_ACC_100",
      "contentPublishForm": {
        "formType": "task",
        "title": "今日好物分享",
        "description": "<p>真的超级好用！ <topic text=\"好物推荐\">#好物推荐</topic></p>",
        "images": [
          { "key": "xhs_img_1", "size": 204800, "width": 1080, "height": 1440, "format": "jpg" }
        ],
        "visibleType": 0
      }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **标题超出长度** | 字符数超过了 20 字的物理限制。 | 向用户解释小红书规范并提供 20 字内的精简版标题。 |
| **图片数量超限** | `images` 数组中的图片超过了 9 张。 | 提醒用户小红书最多只能发布 9 张图片。 |
| **话题不生效** | `<topic>` 标签中的 `raw` 数据格式不正确。 | 严格按照 `challenges` 接口返回的数据结构填入。 |
| **无法定时发布** | `scheduledTime` 设置的时间点已被过，或账号权限不足。 | 校对时间，并确保该账号在蚁小二中状态正常。 |

---
> [!IMPORTANT]
> **描述内容策略**：小红书是搜索型社区。Agent 在生成 `description` 时，除了遵循 HTML 话题规范外，应建议用户在文中多埋入搜索热词（如“怎么破”、“神器”等）。
