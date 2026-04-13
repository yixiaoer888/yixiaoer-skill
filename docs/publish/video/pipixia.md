# 📄 皮皮虾 视频 参数 (Pipixia Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“皮皮虾”平台发布幽默段子、脑洞吐槽或生活趣闻短视频时触发：
- **幽默社区分发**：利用皮皮虾独特的社区氛围进行流量曝光。
- **互动同步**：将个人创意内容同步分享至皮友圈。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装皮皮虾视频 Payload 时需遵守：
1. **纯描述驱动模式**：皮皮虾采用极简发布逻辑，没有独立的 `title` 字段，核心意图全在 `description`。
2. **标签内嵌建议**：虽然没有独立标签位，但 Agent 建议在 `description` 末尾手动增加 #搞笑 #神评论 等后缀。
3. **调性适配**：内容建议轻松愉快，Agent 可建议用户在描述中使用皮友常用黑话。
4. **资源引用规范**：必须通过 `upload` 动作产生视频 key 及封面 key。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| `description` | `string` | 否 | 视频描述/文案内容。建议包含热门话题。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布皮皮虾幽默短视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Pipixia"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "PPX_ACC_01",
        "video": { "key": "ppx_v_1", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "description": "皮皮虾，我们走！这年头谁还不是个段子手啊 #搞笑 #生活"
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
| **描述为空** | `description` 字段完全缺失或全为空格。 | 请补充至少一句话的描述内容。 |
| **误填标题字段** | 传入了 `title` 指令。 | 皮皮虾发布仅识别 description，请合并内容。 |
| **视频画质不佳** | 低分辨率视频被平台限制流向。 | 建议使用 720P 及以上比例素材。 |
| **版权拦截** | 包含违规水印或搬运痕迹明显。 | 请确保内容的原创性或经过合法授权。 |

---
> [!TIP]
> **皮友精神**: 皮皮虾是强社区属性平台。Agent 建议描述文字应具备“皮、逗、暖”的特质，甚至可以主动反问读者以激发“神评论”区的讨论热烈度。
