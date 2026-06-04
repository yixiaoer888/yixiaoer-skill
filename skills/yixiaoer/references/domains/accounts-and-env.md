# 账号、环境与 Skill 同步

适用范围：用户要检查环境、初始化配置、查询账号、查看 skill 安装状态。

**CRITICAL - 用户只要说“先看看环境/账号/同步状态”，就先停在本域，不要提前进入发布、草稿或素材主流程。**

## 优先读取

1. [`../yixiaoer-shared.md`](../yixiaoer-shared.md)
2. [`../workflows/account-selection.md`](../workflows/account-selection.md)

## 常用命令

```bash
yxer doctor
yxer config get
yxer config init --api-key <apiKey>
yxer config set-local-client-id <clientId>
yxer accounts list [platform] [--name 关键词] [--status 1] [--page 1] [--size 20] [--all] [--json]
yxer skill show
yxer skill sync [--global]
```

## 规则

- 首次使用、环境不明、或失败排查后，先 `yxer doctor`
- `doctor` 失败时先修环境，不继续业务流程
- 涉及本机发布前的 `clientId` 判断时，先查询现有配置和账号列表
- `yxer accounts list` 默认查第 `1` 页、每页 `20` 条；传 `--page`、`--size` 可控制单页范围，传 `--all` 时才继续翻页
- 命中多个账号时，结构化列出候选并让用户确认；只有一个可用账号时可自动选中并说明
- 用户明确说“先不要发布，只查环境/账号”时，完成本域后直接停下，不擅自继续发布流程
