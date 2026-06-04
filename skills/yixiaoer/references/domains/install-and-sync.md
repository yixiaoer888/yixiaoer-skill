# 安装、升级与同步

适用范围：用户要安装 skill、重新同步、升级 CLI、查看宿主如何接入该技能。

**CRITICAL - 用户意图是安装、升级或同步时，先完成本域动作，不要跳到发布流程里“顺手再发一次”验证。**

## 优先读取

1. [`../cli/skill-install.md`](../cli/skill-install.md)
2. [`../cli-install-uninstall.md`](../cli-install-uninstall.md)

## 常用命令

```bash
yxer skill show
yxer skill sync
yxer skill sync --global
yxer --version
yxer doctor
```

## 安装命令

```bash
npx skills add "<repo>\\skills\\yixiaoer" -y
npx skills add "<repo>\\skills\\yixiaoer" -g -y
```

## 规则

- `yxer --version` 升级后，应提示重新同步 skill
- `SKILL.md` 或 `references/` 中影响 Agent 行为的文档更新后，应提示重新同步
- `yxer doctor` 返回 `_notice.skills` 时，优先执行 `yxer skill sync`
- 用户明确说“只想同步 skill / 看安装方法”时，完成本域说明后直接停下
