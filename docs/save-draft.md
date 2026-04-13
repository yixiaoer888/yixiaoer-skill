# 📄 草稿管理 Action 参数 (Draft Management Action)

在蚁小二生态中，存在两种“草稿”能力：保存到**蚁小二草稿箱**（系统内部存储）或保存到**目标平台草稿箱**（第三方平台后台）。

> [!IMPORTANT]
> **意图确认原则**：若用户说“存草稿”但未说明位置，Agent **必须询问用户** 明确意图（是存为“蚁小二草稿”还是“平台草稿”），严禁自行猜测。

## 1. 触发场景 (Trigger)

- **蚁小二草稿 $(\text{save-draft})$**：当用户希望先录入内容，待以后手动检查或由他人审核，不希望立即启动任何平台推送流程时。
- **平台草稿 $(\text{publish} + \text{pubType: 0})$**：当内容已比较完善，但希望在平台后台进行最后的 SEO、话题优化或由于平台规则（如视频号）需要移动端二次确认时。
- **典型提示词**：
  - “帮我把这个视频存为蚁小二的草稿”
  - “暂时不发布，先存草稿”
  - “推送到抖音的草稿箱里”

## 2. 交互协议 (Interactive Protocol)

1. **意图深度研判**：根据提示词关键字选取模式。
2. **模式执行链**：
   - **蚁小二草稿**：构造 `action: "save-draft"`。此模式不消耗发布配额。
   - **平台草稿**：构造 `action: "publish"` 并设置 `pubType: 0`。此模式**消耗**发布配额。
3. **分步确认**：告知用户草稿保存的位置（如：“我将把您的文章保存到蚁小二云端草稿箱”）。

## 3. 参数定义 (Parameters)

### 3.1 蚁小二草稿模式 (`action: save-draft`)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | 固定值：`save-draft` |
| `publishType` | `string` | **是** | `video` (视频) 或 `article` (文章) |
| `platforms` | `string[]` | 否 | 记录该草稿拟发布的平台。 |

### 3.2 平台草稿模式 (`action: publish`)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | 固定值：`publish` |
| `contentPublishForm.pubType` | `number` | **是** | **必须设为 `0`**（0: 草稿, 1: 发布）。 |

## 4. 执行指令示例 (Command)

```bash
# 示例：存为蚁小二草稿
node scripts/api.ts --payload='{"action":"save-draft","publishType":"video","desc":"暂存内容"}'

# 示例：存为抖音平台草稿
node scripts/api.ts --payload='{"action":"publish","publishArgs":{"accountForms":[{"platformAccountId":"xxx","contentPublishForm":{"pubType":0,"title":"标题"}}]}}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **存平台草稿却直接发布了** | `pubType` 填错或该平台不支持草稿模式。 | 检查 `pubType` 是否为 `0`。若平台不支持，Agent 应提前告知。 |
| **蚁小二草稿找不到** | 登录了不同的团队或租户环境。 | 请用户核实当前所在的团队 ID 是否正确。 |
| **发布配额被扣除** | 用户选择了“平台草稿”模式。 | 告知用户平台草稿本质上是一次发布过程的预执行。 |

---
> [!TIP]
> **兼容性补丁**：若平台无 `pubType` 字段，可尝试将 `privacy` (或 `status`) 设为 **1 (私密)** 作为替代草稿方案。

