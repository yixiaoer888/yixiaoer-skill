# 📄 获取征文活动 Query 参数 (Get Activities Query)

获取各媒体平台（如百家号、企鹅号、头条号等）当前正在开展的征文任务、激励计划或官方活动。

> [!TIP]
> **流量激励**：参加平台官方征文活动是获取初始推荐量和奖金激励的高效途径。Agent 应对长期更新的账号主动建议查询活动。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户希望参与平台激励计划，或需要获取特定的活动 ID 挂载到其图文/视频作品中。
- **典型提示词**：
  - “最近百家号有什么征文活动？”
  - “帮我找一下适合这个视频参加的活动”

## 2. 交互协议 (Interactive Protocol)

1. **链式建议**：在执行 `categories` 查询后，Agent 可主动提供 `activities` 查询选项，以丰富发布选项。
2. **列表交付**：展示活动标题、奖励说明（若有）及活动截止时间，引导用户选择 1 个进行挂载。
3. **透传原则**：获取活动详情后，将 `id` 或 `raw` 对象正确映射到发布表单。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `activities` | 固定值。 |
| **`account_id`** | `string` | **是** | - | 蚁小二账号 ID (`platformAccountId`)。 |
| `type` | `string` | 否 | `article` | 类型：`video` 或 `article`。 |
| `categoryId` | `string` | 否 | - | 仅查询特定分类下的活动。 |
| `keyword` | `string` | 否 | - | 活动名称搜索关键词。 |

## 4. 执行指令示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"activities","account_id":"67fb2f1735eeb3cf31db3d65","type":"article"}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **找不到特定活动** | 活动已过期，或该账号不满足参与门槛。 | 核实活动起止时间，并检查账号是否符合活动要求（如：分值、权重）。 |
| **挂载后提示非法** | 活动 ID 输入有误，或跨平台使用了 ID。 | 确保 `activityId` 严格来源于对应账号的 `activities` 查询结果。 |

---
> [!IMPORTANT]
> **资源互斥说明**：某些平台可能限制一个作品只能参加一个活动，Agent 应在用户尝试挂载多个活动时予以提醒。
