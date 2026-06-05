# 记录与排查

适用范围：用户要查看发布记录、排查失败、解释为什么 `validate` 或 `publish` 报错。

**CRITICAL - 进入本域后，先定位失败阶段，再决定回退到哪一步；不要把“正式重试发布”当成第一反应。**

## 优先读取

1. [`../workflows/publish-troubleshooting.md`](../workflows/publish-troubleshooting.md)
2. 必要时回读：
   - [`../workflows/common-rules.md`](../workflows/common-rules.md)
   - [`../workflows/local-vs-cloud.md`](../workflows/local-vs-cloud.md)
   - [`./publish.md`](./publish.md)

## 常用命令

```bash
yxer query records [--platform P] [--limit N] [--status S] [--json]
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
yxer publish <type> <platform> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>] --dry-run
yxer doctor
```

## 排查规则

- 先给出错误阶段和错误摘要，再决定回到哪一步修复
- `validate` 失败时，优先回查 `prepare`、`schema get` 和字段来源
- 云发布报代理问题时，可建议改本机发布
- 本机发布报客户端不在线时，提示用户启动并登录客户端，或改回云发布
- 用户只说“解释下这条报错”时，可先停留在本域做解释和修复建议，不擅自触发新的发布写操作
