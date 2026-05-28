# 技能目录架构

本仓库现在按“飞书风格技能文件架构”管理蚁小二技能信息。根目录不再放 `SKILL.md`，真正的技能入口固定放在 `skills/<name>/SKILL.md`：

```text
README.md
skills/
  yixiaoer/
    SKILL.md
references/
  cli/
  legacy/
  platforms/
  workflows/
cmd/
internal/
schemas/
tests/
```

## 分层职责

- `skills/<name>/SKILL.md`
  给外部 AI agent 安装的技能包入口；内容可更精简，但必须明确“执行统一走 yxer CLI”
- `README.md`
  仓库入口，面向人类维护者说明整体结构、安装方式和参考目录
- `references/workflows/`
  放 Agent 必须先阅读的场景工作流
- `references/cli/`
  放 CLI 用法、命令规范、输出约定
- `references/platforms/`
  放平台索引和平台差异入口
- `references/legacy/`
  放迁移说明和未 CLI 化能力索引
- `cmd/` + `internal/`
  放真正可执行的 `yxer` CLI 能力

## 扩展规则

新增能力时：

1. 先补 `yxer` CLI 命令。
2. 再补对应 workflow/reference。
3. 若该能力需要对外给 agent 安装，补 `skills/<name>/SKILL.md`。
4. 如需补充仓库导航或维护说明，更新 `README.md`，不要再新增根级 `SKILL.md`。

## 这样做的目的

- 避免技能主文件越来越长
- 避免 Agent 误走旧 Node 入口
- 让后续平台扩展、命令扩展、工作流扩展都有固定落点
