# 📄 得物 视频 参数 (Dewu Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“得物 (Poizon)”平台分发潮、鞋、穿搭或生活方式内容时触发：
- **种草短视频**：展示好物、穿搭效果或球鞋开箱。
- **潮流动态更新**：同步生活方式相关的精品短视频。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装得物视频 Payload 时需遵守：
1. **申明必填要求**：得物强制要求设置 `declaration`。Agent 应根据内容自动识别是“AI生成”还是“剧情演绎”。
2. **潮流分区选择**：虽然可选分类，但通过 `categories` 接口获取并透传 `raw` 能够极大提升推荐精准度。
3. **内容调性控制**：标题和描述应侧重于“潮流、种草、分享”。
4. **资源引用规范**：必须通过 `upload` 动作产生视频 key 及封面 key。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述。 | - |
| **`declaration`** | `number` | **是** | **创作者申明**: `0`-不添加, `1`-AI生成, `2`-非营销推广, `3`-专业运动, `4`-剧情演绎。 | `0` |
| `category` | `Array` | 否 | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `id`, `text`, `raw` 元数据。

## 4. 执行指令示例 (Command)

```bash
# 发布得物穿搭视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Dewu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DEWU_ACC_01",
        "video": { "key": "dw_v_key", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "夏日清爽运动风穿搭分享",
          "description": "今日份运动风，OOTD 安排！",
          "category": [{ "id": "1", "text": "穿搭", "raw": {...} }],
          "declaration": 2
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
| **强制声明失效** | 未传 `declaration` 字段或数值不在合法范围内。 | 请补全 [3.1](#31-核心表单参数-contentpublishform) 中定义的数值。 |
| **营销过度过滤** | 内容被判定为硬广未报备。 | 建议申明 `declaration: 2` 并确保内容具有分享价值。 |
| **分类加载失败** | `category.raw` 数据格式错误。 | 重新执行 `categories` 接口获取并透传原始对象。 |
| **封面比例不符** | 得物偏好 9:16 的竖向高清屏展示。 | 请调整封面裁剪比例。 |

---
> [!TIP]
> **得物种草力**: 得物用户不仅关注内容更关注其中的实体商品。Agent 建议在 `description` 中明确提及商品型号，利用平台的垂直检索优势获客。
