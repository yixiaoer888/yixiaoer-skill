# 账号、环境与 Skill 同步

适用范围：用户要检查环境、初始化配置、查询账号、查看 skill 安装状态。

## 优先读取

1. [`../yixiaoer-shared.md`](../yixiaoer-shared.md)
2. [`../workflows/account-selection.md`](../workflows/account-selection.md)

## 常用命令

```bash
yxer doctor
yxer config get
yxer config init --api-key <apiKey>
yxer config set-local-client-id <clientId>
yxer accounts list [platform] [--name 关键词] [--status 1] [--json]
yxer skill show
yxer skill sync [--global]
```

## 规则

- 首次使用、环境不明、或失败排查后，先 `yxer doctor`
- `doctor` 失败时先修环境，不继续业务流程
- 涉及本机发布前的 `clientId` 判断时，先查询现有配置和账号列表
- 命中多个账号时，结构化列出候选并让用户确认；只有一个可用账号时可自动选中并说明
