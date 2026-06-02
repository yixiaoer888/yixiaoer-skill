# 发布通道判断工作流

> 适用范围：任何会触发 `validate`、`publish --dry-run`、正式 `publish` 的任务。

## 何时读取

- 用户提到“云发布 / 本机发布 / 本地发布 / 客户端发布”
- 当前任务已经进入发布阶段，需要确定通道
- 云发布或本机发布刚刚失败，需要回退

## 决策树

1. 用户明确要求本机/本地/客户端发布：走 `local`。
2. 用户明确要求不要走云端代理：走 `local`。
3. 用户未指定发布通道：默认走 `cloud`。
4. 云发布遇到代理错误且用户接受切换：改走 `local`。
5. 本机发布提示客户端不在线且用户不便保持在线：改走 `cloud`。

## 本机发布门禁

- 必须显式带 `--publish-channel local`
- 必须确认 `clientId`
- `validate`、`publish --dry-run`、正式 `publish` 必须保持同一套通道参数

## clientId 获取顺序

1. payload 中已有 `clientId`
2. 显式 flags：`--client-id <clientId>`
3. 第四个位置参数：`yxer publish <type> <platform> <payload.json> <clientId>`
4. 本地默认配置：`yxer config set-local-client-id <clientId>`

## 推荐命令

```bash
yxer config get
yxer validate 抖音 video .\payload.json --publish-channel local --client-id <clientId>
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId> --dry-run
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId>
```

## 回退策略

- 云发布报“账号代理不存在”：提示检查代理配置，或改走本机发布
- 本机发布报“客户端不在线”或“获取在线设备列表失败”：提示用户启动并登录客户端，或改回云发布

## 严禁行为

- 用户明确要求本机发布时仍默认走 `cloud`
- 本机发布未确认 `clientId`
- `validate` 用云发布、`publish` 却改成本机发布
