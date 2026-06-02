# 技能安装与同步

本项目按“CLI 先安装，Skill 再安装”的方式给 AI agent 提供能力。

## 目标

- `yxer` CLI：唯一执行入口，负责真正调用蚁小二 API
- `SKILL.md`：给 AI agent 的规则、工作流、命令选择说明
- `README.md`：给仓库维护者的入口说明

也就是说，Skill 不负责执行，它只负责告诉 Agent 何时调用什么 `yxer` 命令。

## 推荐安装流程

1. 安装或编译 `yxer` CLI。
2. 运行 `yxer skill show`，拿到当前技能包目录。
3. 用 skills 工具安装技能：

```bash
npx skills add "<repo>/skills/yixiaoer" -y
```

如需全局安装：

```bash
npx skills add "<repo>/skills/yixiaoer" -g -y
```

也可以直接让 CLI 代为同步并写入版本戳：

```bash
yxer skill sync
yxer skill sync --global
```

如需统一执行“检查状态 + 同步 skill + 查看 CLI 更新指引”，可运行：

```bash
yxer update
yxer update --check
```

## 何时需要重新同步

- `yxer --version` 升级后
- 当前 skill 包中的 `SKILL.md` 更新后
- `skills/yixiaoer/references/` 中影响 Agent 行为的文档更新后

## 漂移检查

- `yxer skill show` 会显示当前 `skills.stamp` 状态
- `yxer doctor` 会在 `_notice.skills` 中提示是否需要重新同步

## 设计原则

- 技能负责“让 Agent 会用 CLI”
- CLI 负责“把事情真正做完”
- 不再保留默认 Node 脚本执行入口
