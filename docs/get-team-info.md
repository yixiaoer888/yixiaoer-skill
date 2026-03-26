# 当前团队信息 (Get Team Info)

查询当前上下文环境所属团队的核心元数据，包括名称、发布额度及当前用户的角色权限。

## 场景描述 (Usage)

- "帮我查一下我所属团队的名称。"
- "看看目前的团队套餐还有多少发布次数。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `teamId` | `string` | **可选**。如未提供，则默认为当前活动团队。 |

## 脚本映射 (Script)

- **脚本路径**: `../scripts/get-team-info.ts`
- **调用示例**: `node get-team-info.ts --teamId=t_999`

## 输出结果 (Output)

脚本需输出标准的 JSON 对象，包含 `teamId`, `name`, `role`, `quota` 等字段。
- **角色限制**: `admin` 或 `manager` 角色可见全部额度。
- **普通成员**: 屏蔽某些敏感的套餐信息。
