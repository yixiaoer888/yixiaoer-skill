# 📄 网易号 文章 参数 (WangYiHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“网易号”平台发布网易新闻相关的资讯、评论或深度专栏时触发：
- **新闻分发**：利用网易新闻客户端的海量用户基数进行全网同步。
- **原创申明**：申明原创或勾选创作类型，以获取更多分成及权益。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装网易号 Payload 时需遵守：
1. **原创逻辑双重判定**：网易号包含 `type` (勾选原创) 和 `declaration` (创作申明) 两个维度。Agent 应根据用户提示词同时设置这两个字段。
2. **封面格式化**：`covers` 数组中的对象必须包含完整的 `OldCover` 结构，且资源需通过 `upload` 动作获取。
3. **定时发布自律**：支持 `scheduledTime`，设置时应确保账号已开通网易定时发布权益。
4. **内容规范**：网易对文章质量要求较高，Agent 建议标题应兼具信息量与吸引力。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **创作申明**: 1-AI生成, 2-个人原创, 3-取材网络, 4-虚构演绎。 | - |
| `type` | `number` | 否 | **是否勾选原创**: 0-不勾选, 1-勾选。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布网易号个人原创文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["网易号"],
  "publishArgs": {
    "content": "<h1>网易生态下的自媒体生存法则</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_wy_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "网易生态下的自媒体生存法则：深度复盘",
          "content": "<h1>生存法则</h1><p>正文内容...</p>",
          "covers": [{ "key": "wy_cov_1", "size": 102400, "width": 800, "height": 600 }],
          "declaration": 2,
          "type": 1,
          "pubType": 1
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
| **申明与类型冲突** | 勾选了原创但申明选择了“取材网络”。 | 确保 `declaration` 与 `type` 的逻辑含义一致。 |
| **封面获取异常** | 未通过 `upload` 接口上传封面且 `key` 格式错误。 | 重新执行上传动作获取 OSS 内部标识。 |
| **账号权限拦截** | 账号处于扣分期或违规处罚中。 | 检查网易号后台信用分。 |
| **定时任务未触发** | 服务端延迟或设置的时间过于接近当前。 | 建议定时发布时间至少比当前晚 30 分钟。 |

---
> [!TIP]
> **网易全家桶**: 网易号内容可能被推送至有道、网易云音乐等生态位，Agent 建议内容尽量保持高品质，以迎合网易用户群体的偏好。
