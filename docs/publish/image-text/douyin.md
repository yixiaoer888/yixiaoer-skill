# 📄 抖音图文发布参数 (Douyin Image-Text)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在进入细分参数前，你 **必须** 已经阅读并理解了 [图文发布通用索引](./index.md) 中定义的 Payload 根结构。本文档仅描述 `contentPublishForm` 内部的平台差异化字段。

## 1. 触发场景 (Trigger)

当用户明确要在“抖音”发布图文动态，且涉及以下特有需求时触发：
- **视觉社交**：发布多图动态、设置图文标题。
- **互动增强**：挂载地理位置、关联背景音乐、加入合集。
- **标签/话题**：在描述中使用 `<topic>` 标签对内容进行分类。

## 2. 交互协议 (Interactive Protocol)

Agent 在构造抖音图文 Payload 时需遵守：
1. **内容格式化**：`description` 支持简易 HTML。话题应封装在 `<topic text='...' raw='...'>#话题</topic>` 结构中。
2. **辅助查询必选**：涉及 `location`, `music`, `collection` 等复杂对象时，**严禁手动构造**。必须调用对应接口获取 `raw` 原始数据。
3. **资源合规性**：确保 `images` 数组中每个 `OldImage` 对象的 `key` 均来自有效的 `upload` 动作。

## 3. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 图文标题 (建议简短有力)。 | - |
| `description` | `string` | 否 | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| **`images`** | `Array` | **是** | 图片数组，须包含 1-35 张图片。使用 `OldImage[]` 结构。 | - |
| `location` | `Object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `music` | `Object` | 否 | **背景音乐**: 使用 `MusicItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |
| `collection` | `Object` | 否 | **合集信息**: 使用 `Category` 结构。 | - |
| `sub_collection` | `Object` | 否 | **合集选集信息**: 使用 `Category` 结构。 | - |

### 3.1 复杂对象结构 (Data Schemas)

- **OldImage**: 包含 `key`, `size`, `width`, `height`, `format`。
- **PlatformDataItem / Category**: 必须包含 `yixiaoerId`, `yixiaoerName` 和 **完整的 `raw` 对象**。
- **MusicItem**: 包含 `yixiaoerId`, `yixiaoerName`, `duration`, `playUrl` 和 `raw`。

## 4. 执行指令示例 (Command)

```bash
# 抖音图文发布：带话题和地点的多图任务
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["Douyin"],
  "publishArgs": {
    "accountForms": [{
      "platformAccountId": "DY_ACC_002",
      "contentPublishForm": {
        "formType": "task",
        "title": "记录生活",
        "description": "<p>今日份的好心情 <topic text=\"开心\">#开心</topic></p>",
        "images": [
          { "key": "img_k1", "size": 512000, "width": 1080, "height": 1920, "format": "jpg" }
        ],
        "location": { "yixiaoerId": "loc_001", "yixiaoerName": "上海", "raw": {...} }
      }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **描述文字截断** | `description` 长度超过 1000 字符限制。 | 提醒用户精简描述内容。 |
| **音乐挂载失败** | `MusicItem` 中的 `raw` 数据过期或 ID 不匹配。 | 重新执行 `music` 查询动作。 |
| **图片加载异常** | `key` 对应的 OSS 资源权限受限或已被清理。 | 确认上传动作在 24 小时内执行。 |
| **定时任务失败** | `scheduledTime` 设置在过去或超过平台限制（如 7 天后）。 | 校对当前服务器时间并按平台规则设置。 |

---
> [!IMPORTANT]
> **话题识别规范**：抖音极度依赖话题搜索。Agent 在解析用户意图时，若用户提及“带上话题 XXX”，应自动将其包装为 HTML 的 `<topic>` 标签形式。
