# OpenClaw 龙虾技能 (OpenClaw Skill)

该技能定义了蚁小二全平台的媒体管理与运营能力。
通过元数据驱动（Skill -> Doc -> Script）的模式，为 AI 提供可执行的工具集。

## 技能定义 (Metadata)

- **ID**: `openclaw-skill-core`
- **版本**: `1.0.0`
- **架构模式**: 文档驱动型自律脚本 (Doc-Driven Scripts)
- **运行环境**: Node.js v18+ (Direct Runtime)

## 配置与安全 (Config & Secrets)

为了保持技能的“无包设计 (Package-Free)”，所有的敏感信息应通过**环境变量**注入：

1.  **生产环境**: 在龙虾系统 (OpenClaw) 的环境变量配置中填入 `YIXIAOER_API_KEY`。
2.  **本地开发**: 
    - 建议在根目录创建 `.env` 文件（请务必将其加入 `.gitignore`）。
    - 运行脚本时，Node.js 20.6+ 可以使用内置标志加载：`node --env-file=.env scripts/query-accounts.ts`。
    - 或者在运行前手动导出：`$env:YIXIAOER_API_KEY="xxx"; node scripts/query-accounts.ts` (PowerShell)。

## 能力地图 (Capabilities)

本技能通过映射 `docs/` 下的指令文档到 `scripts/` 下的执行脚本实现功能的动态调度。

| 能力名称 | 指令文档 (Trigger) | 执行脚本 (Implementation) | 核心功能 |
| :--- | :--- | :--- | :--- |
| **查询账号列表** | [query-accounts.md](./docs/query-accounts.md) | [query-accounts.ts](./scripts/query-accounts.ts) | 获取租户下绑定的媒体账号 |
| **当前团队信息** | [get-team-info.md](./docs/get-team-info.md) | [get-team-info.ts](./scripts/get-team-info.ts) | 获取团队名称、角色、额度信息 |

## 运行规范 (Protocol)

1.  **加载**: 系统分析 `skill.md` 确定可用功能列表。
2.  **理解**: 当用户发出需求时，系统读取对应的 `docs/*.md` 确认输入参数与预期输出。
3.  **执行**: 调用 `scripts/*.ts` 并通过命令行参数或环境变量传递经过处理的输入。
4.  **返回**: 脚本输出标准的 JSON 数据，由系统解析并反馈。
