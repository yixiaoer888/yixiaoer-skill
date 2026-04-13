# 📄 快手 视频 参数 (Kuaishou Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“快手”分发视频，且涉及以下特有需求时触发：
- **发布策略**：设置可见性（公开、私密、好友）、控制同城展示、开启同框/下载权限。
- **互动/营销**：挂载 POI 位置、关联合集、插入购物车或小程序。
- **内容标注**：宣称内容由 AI 生成或属于演绎情节。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装快手视频 Payload 时需遵守：
1. **意图精准识别**：快手支持精细的权限控制（如 `nearby_show`），Agent 应根据用户提示词（如“别发同城”）准确设置布尔值。
2. **资源先行**：确保视频和封面的 `key` 已通过 `upload` 动作获取。
3. **数据完整性**：对于 `location`, `collection`, `mini_app` 等复杂字段，必须通过相关接口获取并在 Payload 中透传完整的 `raw` 数据。
4. **互斥逻辑校验**：购物车 (`shopping_cart`) 与小程序 (`mini_app`) 通常无法同时挂载，Agent 应根据优先级或询问用户进行二选一处理。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`visibleType`** | `number` | **是** | **可见范围**: `0`-公开, `1`-私密, `3`-好友可见。 | `0` |
| `title` | `string` | 否 | 视频标题。 | - |
| `description` | `string` | 否 | 视频描述。 | - |
| `declaration` | `number` | 否 | **视频声明**: `0`-不申明, `1`-AI 生成, `2`-演绎演绎, `3`-个人观点。 | `0` |
| `location` | `Object` | 否 | **位置信息**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `collection` | `Object` | 否 | **合集信息**: 使用 `Category` 结构。 | - |
| `nearby_show` | `boolean` | 否 | **同城展示**: 是否在同城频道显示。 | `true` |
| `allow_same_frame` | `boolean` | 否 | 是否允许同框。 | `false` |
| `allow_download` | `boolean` | 否 | 是否允许下载。 | `false` |
| `shopping_cart` | `Object` | 否 | 关联商品信息。 | - |
| `mini_app` | `Object` | 否 | 挂载小程序。 | - |

### 3.2 复杂结构说明

- **PlatformDataItem / Category**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布快手视频：开启 AI 声明并关闭同城展示
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["快手"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_ks_vid_001",
        "video": { "key": "v_key_ks", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "快手视频实战测试",
          "visibleType": 0,
          "nearby_show": false,
          "declaration": 1
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
| **位置加载失败** | `location.raw` 过期或格式不完整。 | 重新执行 `action: "locations"` 获取最新 POI 信息。 |
| **同城展示配置无效** | `nearby_show` 传值类型错误（需为 boolean）。 | 确保 JSON 中传入的是布尔值而非字符串。 |
| **挂载冲突** | 同时尝试挂载购物车和小程序。 | 快手仅支持单一挂载项，请保留其一。 |
| **可见性设置被重置** | 账号权重或平台规则限制了该内容的可见性。 | 检查账号状态，或尝试使用更平衡的内容策略。 |

---
> [!TIP]
> **老铁情怀**: 快手社区非常看重博主与粉丝的连结，建议开启 `allow_same_frame` 以增加互动和病毒式传播的机会。
