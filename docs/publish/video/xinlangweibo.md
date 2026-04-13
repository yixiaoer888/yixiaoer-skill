# 📄 新浪微博 视频 参数 (Sina Weibo Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“新浪微博”分发短视频资讯、VLOG 或博文配图视频时触发：
- **实时热点分发**：利用微博的强时效性分发短视频内容。
- **互动/定位发布**：挂载 POI 位置、设置定时任务或标注创作者类型。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装微博视频 Payload 时需遵守：
1. **文案社交化**：微博对 `title` 和 `description` 有较高权值。Agent 建议描述内容包含 #话题，增强博文可读性。
2. **位置透传原则**：若包含 `location`，必须通过 `locations` 接口获取并透传完整的 `raw` 对象。
3. **内容来源申明**：准确设置 `type` (1-原创, 2-转载, 3-二创)，这直接影响视频在发现页的权重。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述内容 (博文正文)。推荐包含 #话题。 | - |
| **`type`** | `number` | **是** | **内容类型**: `1`-原创, `2`-转载, `3`-二创。 | `1` |
| `location` | `Object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `collection` | `Object` | 否 | **合集信息**: 包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |

### 3.2 复杂结构说明

- **PlatformDataItem**: 必须包含 `id`, `text`, `raw` 原始元数据。

## 4. 执行指令示例 (Command)

```bash
# 发布微博原创视频：带定位和话题
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Xinlangweibo"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WB_ACC_01",
        "video": { "key": "wb_v_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 科技大赏实测",
          "description": "今天带大家看看 2026 年最硬核的科技产品！ #科技 #微博视频号",
          "type": 1,
          "location": { "id": "loc_sh_001", "text": "上海·静安", "raw": {...} }
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
| **位置解析失败** | `location.raw` 数据格式不满足微博 API 要求。 | 确保调用 `locations` 接口获取并原样透传。 |
| **标题重复拦截** | 相同微博号在短时间内发布标题高度一致的内容。 | 在标题中增加时效性或差异化字符。 |
| **图片/视频解析失败** | key 引用错误或素材已在 OSS 端失效。 | 重新执行 `upload` 动作确认资源可用性。 |
| **博文触发敏感词** | `description` 包含博文审核红线词汇。 | 请系统性检查并优化文案表达。 |

---
> [!TIP]
> **热搜驱动力**: 微博内容的核心在于传播。Agent 建议标题采用“爆料式”或“提问式”，并在正文描述中多 @相关博主，利用微博的社交关系网闭环流量。
