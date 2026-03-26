# 当前团队信息 (Get Team Info)

查询当前上下文环境所属团队的核心元数据，包括名称、角色权限及最新关联团队信息。

## 场景描述 (Usage)

- "帮我查一下我所属团队的名称。"
- "看看目前团队的信息详情。"

## 参数定义 (Parameters)

该能力无需额外参数，将自动获取当前 Token 关联的最新团队。

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/get-team-info.ts`
- **流程**:
  1. 调用 `GET /api/users/info` 获取用户的基础元数据及关联的 `latestTeamId`。
  2. 使用 `latestTeamId` 调用 `GET /api/teams/:id` 获取详细的团队配置。
- **调用示例**: `node get-team-info.ts`

## 输出结果 (Output)

脚本需输出标准的 JSON 对象，包含 `user` 简要信息以及 `team` 详细对象。
请确保环境变量 `YIXIAOER_API_KEY` 已设置。
