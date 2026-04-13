# 📄 腾讯微视 视频 参数 (Weishi Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“腾讯微视”分发短视频动态、才艺秀或生活片段时触发：
- **竖屏分发**：将精品短视频同步至微视社区。
- **社交同步**：在腾讯社交生态内进行短内容曝光。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装微视视频 Payload 时需遵守：
1. **简炼双核原则**：微视发布主要聚焦于 `title` 和 `description`。Agent 应确保文案极具社交亲和力。
2. **竖向兼容性**：微视是典型的竖屏平台，Agent 建议用户优先选择 9:16 的视频素材。
3. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。
4. **定时分发**：支持 `scheduledTime`，设置时需满足微视对未来发布时间点的约束要求。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述内容。请多使用 #话题。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布腾讯微视趣味短视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Weishi"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WS_ACC_01",
        "video": { "key": "ws_v_1", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "今日份开心放送！",
          "description": "生活就要皮一点，大家觉得呢？ #搞笑 #生活"
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
| **内容过短或缺失** | `title` 或 `description` 完全为空。 | 请补全必要的文案。 |
| **素材比例报错** | 选用了横屏视频且未适配。 | 建议在 Agent 辅助下先进行画面裁剪或边框填充。 |
| **发布权限拦截** | 账号处于被禁言或异常状态。 | 检查微视 App 内的账号健康度。 |
| **定时任务未触发** | 网络波动或时间戳格式错误。 | 确认 Unix 时间戳为秒级。 |

---
> [!TIP]
> **强社交属性**: 微视与腾讯社交链紧密协同。Agent 建议描述中包含具共鸣力的话题，并利用微视的背景音乐库优化视频氛围以提升完播率。
