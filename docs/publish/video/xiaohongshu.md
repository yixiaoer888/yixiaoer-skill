# 📄 小红书视频发布参数 (XiaoHongShu Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `contentPublishForm` 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“小红书”发布视频笔记，且涉及以下需求时触发：
- **内容标注**：虚构演绎申明、AI 合成内容申明。
- **社交关联**：挂载地理位置、加入合集、关联群聊或直播预告。
- **商业化**：挂载小红书内部店铺商品。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装小红书 Payload 时需遵守：
1. **标题字数敏感**：小红书标题上限仅 **20 字**，请务必精简。
2. **描述话题规范**：正文描述上限 1000 字，建议包含 **#话题**。
3. **可见性状态确认**：默认为公开 (visibleType: 0)，可按需设为私密 (1) 或好友可见 (3)。
4. **动态数据透传**：对于位置、合集、商品等，必须先行调用查询接口获取 `raw` 对象且完整透传。

## 3. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| `title` | `string` | 否 | 视频标题 (**最高 20 字**)。建议吸引眼球。 | - |
| `description` | `string` | 否 | 视频描述 (最高 1000 字)。支持 Emoji。 | - |
| `declaration` | `number` | 否 | **内容类型申明**: `1`-虚构演绎, `2`-笔记含 AI 合成内容。 | - |
| `createType` | `number` | 否 | **创作类型**: `1`-原创, `0`-不申明。 | `0` |
| `location` | `object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| **`visibleType`** | `number` | **是** | **可见类型**: `0`-公开, `1`-私密, `3`-好友可见。 | `0` |
| `collection` | `object` | 否 | **合集信息**: 使用 `Collection` 结构。 | - |
| `group` | `object` | 否 | **群聊信息**: 使用 `Group` 结构。 | - |
| `bind_live_info` | `object` | 否 | **直播预告信息**: 使用 `LiveInfo` 结构。 | - |
| `shopping_cart` | `object[]` | 否 | **关联商品信息**: 使用 `ShoppingCartItem` 结构数组。 | - |

### 3.1 复杂结构补充 (Complex Structures)

- **PlatformDataItem / Collection / Group / LiveInfo / ShoppingCartItem**: 所有复杂对象必须包含 `yixiaoerId`, `yixiaoerName` 和 **完整的 `raw` 对象**。

## 4. 执行指令示例 (Command)

```bash
# 小红书视频发布：全字段适配演示
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Xiaohongshu"],
  "publishArgs": {
    "accountForms": [{
      "platformAccountId": "XHS_001",
      "video": {"key": "v_key", "width": 1080, "height": 1440, "size": 1024000},
      "cover": {"key": "c_key", "width": 1080, "height": 1440, "size": 300000},
      "coverKey": "c_key",
      "contentPublishForm": {
        "formType": "task",
        "title": "测试标题",
        "description": "内容描述 #话题",
        "visibleType": 0,
        "createType": 1
      }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **标题太长 (20字限)** | 用户提供的标题超过了平台物理限制。 | 提示用户“小红书标题不能超过20字”并展示截断版本。 |
| **视频比例拦截** | 小红书对非 3:4 或 9:16 的视频可能有黑边。 | 建议在 Agent 辅助确认时提示视频比例。 |
| **复杂参数缺失** | `raw` 对象为空或 ID 不匹配。 | 确保先执行了对应数据的查询操作（如 `goods`, `locations`）。 |
| **发布频率受限** | 在短时间内连续发布大量内容触发风控。 | 建议增加任务执行间隔。 |

---
> [!IMPORTANT]
> **资源引用提示**：小红书是一个极度看重封面的平台。建议 `coverKey` 关联的图片具有高清晰度，并确保经过 `upload` 动作产生 key。
