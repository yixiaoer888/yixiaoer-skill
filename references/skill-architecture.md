# 技能目录架构

本仓库现在按“飞书风格技能文件架构”管理蚁小二技能信息：

```text
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

- `SKILL.md`
  只放技能简介、触发场景、硬性规则、CLI 入口、reference 索引
- `references/workflows/`
  放 Agent 必须先阅读的场景工作流
- `references/cli/`
  放 CLI 用法、命令规范、输出约定
- `references/platforms/`
  放平台索引和平台差异入口
- `references/legacy/`
  放迁移说明和旧入口兼容说明
- `cmd/` + `internal/`
  放真正可执行的 `yxer` CLI 能力

## 扩展规则

新增能力时：

1. 先补 `yxer` CLI 命令。
2. 再补对应 workflow/reference。
3. 最后只在 `SKILL.md` 增加索引，不把长文档继续堆进主技能文件。

## 这样做的目的

- 避免技能主文件越来越长
- 避免 Agent 误走旧 Node 入口
- 让后续平台扩展、命令扩展、工作流扩展都有固定落点
