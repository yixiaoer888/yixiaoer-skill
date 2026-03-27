# 当前团队信息 (Get Team Info)

查询当前授权 Token 关联的团队核心元数据，包括资源限额、成员数、点数消耗及插件状态。

## 场景描述 (Usage)

- "帮我查一下我的团队 ID 和可用点数。"
- "看看目前团队的账号点数还能加几个账号。"

## 参数定义 (Parameters)

该能力无需额外参数，自动根据当前 Token 关联的 `latestTeamId` 查询。

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/get-team-info.ts`
- **流程**:
  1. 调用 `GET /api/users/info` 获取用户元数据及 `latestTeamId`。
  2. 调用 `GET /api/teams/:id` 获取团队详情。
- **调用示例**: `node scripts/get-team-info.ts`

## 输出结果 (Output)

脚本输出合并后的 JSON 对象。核心字段如下：

### Team 字段 (Team Details)
| 字段名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `id` | `string` | 团队唯一 ID |
| `name` | `string` | 团队名称 |
| `logoUrl` | `string` | 团队 LOGO |
| `isVip` | `boolean` | 是否为付费版本 |
| `accountCapacityLimit` | `number` | 总账号点数上限 |
| `accountCapacity` | `number` | 已使用的账号点数 |
| `memberCountLimit` | `number` | 成员数上限 |
| `memberCount` | `number` | 当前成员数 |
| `capacity` | `number` | 素材云存储总量 (B) |
| `usedCapacity` | `number` | 已使用素材容量 (B) |
| `interestCount` | `number` | 关联的权益包数 |
| `appPublish` | `boolean` | 是否支持移动端发布 |
| `createdAt` | `number` | 团队创建时间戳 |
| `components` | `array` | 激活的插件列表（如 AI 助手、数据大屏等） |

请确保环境变量 `YIXIAOER_API_KEY` 已设置。

