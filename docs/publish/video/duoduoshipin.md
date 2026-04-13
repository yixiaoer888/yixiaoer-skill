# 📄 多多视频 视频 参数 (Duoduoshipin Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“多多视频 (拼多多)”分发内容时触发：
- **带货分发**：发布带有商品购物车链接的种草或带货短视频。
- **流量同步**：将个人 IP 或品牌视频同步到拼多多高频流量位。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装多多视频 Payload 时需遵守：
1. **购物车逻辑确认**：多多视频通常与商品挂载相关。若用户提到“带货”或指定商品，必须填充 `shopping_cart` 对象。
2. **极简描述**：多多视频描述不宜过长，重点突出优惠或产品特点。
3. **定时任务**：支持 `scheduledTime`，设置时需确保时间戳的有效性。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| `description` | `string` | 否 | 视频描述内容。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `shopping_cart` | `Object` | 否 | **关联商品**: 包含 `goods_id` (商品 ID) 和 `source` (系统来源, 如 `pdd`)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布多多视频带货内容
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Duoduoshipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PDD_V_ACC_01",
        "video": { "key": "pdd_v_1", "size": 1024000, "width": 1080, "height": 1920, "duration": 30 },
        "contentPublishForm": {
          "formType": "task",
          "description": "这款厨房好物，用了都说好！",
          "shopping_cart": {
            "goods_id": "item_998877",
            "source": "pdd"
          }
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
| **商品 ID 验证失败** | `goods_id` 不存在或已下架。 | 确认拼多多平台该商品的销售状态。 |
| **描述涉及违规** | 包含“最低价”、“第一”等违禁营销词汇。 | 请修改描述，遵守电商合规导向。 |
| **视频画质异常** | 分辨率低于 720P。 | 多多视频建议使用原生手机比例的高清素材。 |
| **关联权限缺失** | 账号未开通多多视频带货权限。 | 检查拼多多后台的账号达人等级。 |

---
> [!TIP]
> **拼多多流量金矿**: 多多视频是转化效果极佳的电商分流位。Agent 建议内容尽量接地气，并在描述中精准利用优惠信息激发用户点击购物车。
