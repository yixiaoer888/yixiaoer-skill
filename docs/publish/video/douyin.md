# 📄 抖音 视频 参数 (DouYin Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `contentPublishForm` 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“抖音”发布视频，且涉及以下特有需求时触发：
- **社交互动**：挂载地点、带上 #话题、关联合集。
- **营销/变现**：挂载小程序、插入购物车、团购推广。
- **合规声明**：如声明视频是由 AI 生成。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装抖音 Payload 时需遵守：
1. **动态参数获取**：对于 `location`, `music`, `challenge` 等字段，**严禁手动拼写**。必须先调用对应的 `get-*` action 获取 `yixiaoerId` 和 `raw` 数据。
2. **存入草稿箱引导**：
   - 存为“平台草稿” -> 设置 `pubType: 0`。
   - 存为“私密发布” -> 设置 `visibleType: 1`。
3. **参数完整性**：必须按需求填入所有参数，严禁随意删除原始定义的任何字段（详见下表）。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (1-50字)。建议包含核心关键词。 | - |
| **`description`** | `string` | **是** | 视频描述 (1-500字)。支持 #话题 标签。 | - |
| `pubType` | `number` | 否 | **发布模式**: `1`: 直接发布, `0`: 存为平台草稿。 | `1` |
| `visibleType` | `number` | 否 | **可见性**: `0`: 公开, `1`: 私密, `2`: 好友可见。 | `0` |
| `horizontalCover` | `object` | 否 | 抖音视频横板封面，使用 `OldCover` 结构。 | - |
| `statement` | `number` | 否 | **声明**: `3`-内容由 AI 生成, `4`-可能引人不适, `5`-虚构演绎, `6`-危险行为。 | - |
| `location` | `object` | 否 | **地理位置**: 须包含 `yixiaoerId`, `yixiaoerName`, `raw`。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `allow_save` | `number` | 否 | 保存权限: `0`-不允许, `1`-允许。 | `0` |
| `shoppingCart` | `object[]` | 否 | 购物车列表，使用 `ShoppingCart` 结构。 | - |
| `groupShopping` | `object` | 否 | 团购信息，使用 `ShoppingCart` 结构。 | - |
| `collection` | `object` | 否 | 合集信息，使用 `Category` 结构。 | - |
| `sub_collection` | `object` | 否 | 合集选集，使用 `Category` 结构。 | - |
| `sync_apps` | `object[]` | 否 | 同时发布应用，使用 `Category[]`。 | - |
| `hot_event` | `object` | 否 | 热点事件，使用 `Category` 结构。 | - |
| `challenge` | `object` | 否 | **话题挑战**: 参与热门话题需注入。 | - |
| `mini_app` | `object` | 否 | **挂载小程序**: 须包含小程序 ID 及 `raw`。 | - |
| `music` | `object` | 否 | **背景音乐**: 建议调用 `music` 接口获取。 | - |
| `cooperation_info` | `object` | 否 | 共创信息。 | - |
| `game` | `object` | 否 | 游戏挂载信息，使用 `GameItem` 结构。 | - |

### 3.1 复杂结构说明

- **OldCover**: 包含 `key`, `size`, `width`, `height`。
- **PlatformDataItem / Category / MiniApp / GameItem**: 必须包含 `yixiaoerId`, `yixiaoerName`, `raw`。
- **ShoppingCart**: 包含 `sale_title` 和 `raw`。
- **MusicItem**: 包含 `yixiaoerId`, `yixiaoerName`, `duration`, `raw`。

## 4. 执行指令示例 (Command)

```bash
# 抖音视频发布：全参数结构演示
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Douyin"],
  "publishArgs": {
    "accountForms": [{
      "platformAccountId": "DY_ACC_001",
      "video": {"key": "v_key", "width": 1080, "height": 1920, "size": 1024000},
      "cover": {"key": "c_key", "width": 1080, "height": 1920, "size": 300000},
      "coverKey": "c_key",
      "contentPublishForm": {
        "formType": "task",
        "title": "测试视频",
        "description": "内容描述 #话题",
        "statement": 3
      }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **标题/描述超长** | 字数超过了 50/500。 | 自动截断或提醒用户修改。 |
| **位置挂载失败** | `location.raw` 数据已过期或结构不完整。 | 重新执行 `get-locations` 动作获取最新数据。 |
| **AI 识别拦截** | 未勾选 `statement: 3` 但内容特征明显。 | 建议在 Agent 辅助确认后自动注入该声明。 |
| **参数缺失报错** | 漏掉了必要的 `raw` 数据或 `yixiaoerId`。 | 重新对照表格检查各复杂对象是否完整。 |

---
> [!IMPORTANT]
> **资源引用提示**：抖音对封面图尺寸有严格要求，建议使用与视频比例一致（通常 9:16）的图片，且必须先通过 `upload` 获取 key。
