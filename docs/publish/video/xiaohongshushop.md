# 📄 小红书店铺 视频 参数 (Xiaohongshu Shop Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“小红书店铺 (商家后台/达人选品)”分发带货视频、新品预告或店铺内容时触发：
- **挂载商品分发**：发布带有小红书商品详情或购物车链接的视频。
- **互动种草**：标注位置信息、设置定时发布或精准关联店铺商品。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装小红书店铺视频 Payload 时需遵守：
1. **购物车逻辑先行**：若用户提到“带货”或指定商品，必须填充 `shoppingCart` 数组。Agent 应通过 `goods` 接口获取最新的商品元数据。
2. **位置信息透传**：支持 `location` 挂载。必须通过 `locations` 接口获取并透传完整的 `raw` 固定对象。
3. **内容调性对齐**：小红书侧重于“种草、质感”。Agent 应建议标题和描述具有吸引力。
4. **资源引用规范**：必须通过 `upload` 动作产生的 key 引用视频和封面。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。建议 1-20 字符。 | - |
| **`description`** | `string` | **是** | 视频描述内容。 | - |
| `location` | `Object` | 否 | **地理位置**: 使用 `PlatformDataItem` 结构。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `shoppingCart` | `Array` | 否 | **关联商品列表**: 使用 `ShoppingCartItem[]` 结构。 | - |

### 3.2 复杂结构说明

- **PlatformDataItem**: 必须包含 `id`, `text`, `raw` 对象。
- **ShoppingCartItem**: 包含 `yixiaoerId`, `yixiaoerName`, `price` 等核心信息，需确保与 `goods` 接口返回一致。

## 4. 执行指令示例 (Command)

```bash
# 发布小红书店铺带货视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Xiaohongshushop"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_SHOP_01",
        "video": { "key": "xhs_v_1", "size": 1024000, "width": 1080, "height": 1440, "duration": 30 },
        "contentPublishForm": {
          "formType": "task",
          "title": "这件新品太绝了！商家版首发实测",
          "description": "家人们，这件真的值得冲！细节满满。 #小红书好物 #新品预告",
          "shoppingCart": [
            { "goods_id": "xhs_item_99", "price": 9990 }
          ],
          "location": { "id": "loc_01", "text": "上海·新天地", "raw": {...} }
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
| **商品 ID 验证不通过** | `shoppingCart.goods_id` 与在该店铺上架的 ID 不符。 | 确认商品已完成小红书商家中心报备且状态正常。 |
| **位置偏移** | `location.raw` 中的元数据格式损坏。 | 必须实时从获取位置接口获取并直传。 |
| **画质不符合要求** | 视频分辨率低于小红书商家端的最低标准。 | 建议在 Agent 辅助下进行资源质量核查。 |
| **定时任务限制** | 小红书对店铺端定时发布的权益有专门等级限制。 | 请检查该账号的商家权益包。 |

---
> [!TIP]
> **精致种草场**: 小红书店铺内容应保持高度的滤镜质感。Agent 建议视频内容应多展现使用细节，并利用小红书的社交分发属性，在描述中积极使用自带表情及热门 #标签。
