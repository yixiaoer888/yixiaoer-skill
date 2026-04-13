# 📄 腾讯视频 视频 参数 (Tencent Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“腾讯视频 (含 App 及网页版)”分发影视剪辑、原创短剧、评测解析或生活记录时触发：
- **精品视频同步**：推送至腾讯视频独立 App 的长短视频流。
- **多维度合规**：满足平台对 AI 生成、剧情演绎及取材网络的强制申明要求。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装腾讯视频 Payload 时需遵守：
1. **标题权重核心**：腾讯视频高度依赖 `title` 参数进行语义推荐。Agent 需提供具吸引力且准确的标题。
2. **强制申明规范**：腾讯视频包含细致的合规项。Agent 应根据视频属性准确设置 `declaration` 数值。
3. **资源引用协议**：必须通过 `upload` 动作产生的 key 引用视频和封面。
4. **定时自律要求**：支持 `scheduledTime`，设置时需满足账号的定时发布等权益等级要求。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| `tags` | `string[]` | 否 | 视频标签。建议 1-5 个。 | - |
| `declaration` | `number` | 否 | **申明**: `1`-AI生成, `2`-剧情演绎, `3`-取材网络, `4`-个人观点。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布腾讯视频个人观点类视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Tengxunshipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TX_V_ACC_01",
        "video": { "key": "tx_v_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "实战复盘：腾讯视频流量分发实操手册",
          "tags": ["自媒体", "经验", "干货"],
          "declaration": 4
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
| **无有效标题** | `title` 字段缺失。 | 腾讯视频必须传标题。 |
| **申明数值异常** | `declaration` 使用了非合法枚举值。 | 检查并锁定为 [3.1](#31-核心表单参数-contentpublishform) 中定义的数值。 |
| **标签解析失败** | 解析 `tags` 数组为字符串时出现非法字符。 | 移除数组中的特殊符号。 |
| **视频文件过大** | 素材超过当前账户等级的直接上传体积上限。 | 建议在 Agent 辅助下先进行质量压缩。 |

---
> [!TIP]
> **顶级娱乐生态**: 腾讯视频是打通腾讯生态娱乐入口的关键。Agent 建议视频内容应具有精品感，配合高清横版封面，以适应 App 端的沉浸式浏览逻辑。
