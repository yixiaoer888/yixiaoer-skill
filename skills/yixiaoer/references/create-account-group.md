# 创建账号分组

用于创建账号分组，对应接口 `POST /groups`。

## 推荐命令

```bash
yxer account-group create 核心账号组 --dry-run
yxer account-group create 核心账号组
```

## 请求说明

- 当前 CLI 以 `{ "name": "<分组名>" }` 作为请求体。
- 写操作先执行 `--dry-run`，确认请求体后再正式提交。

## 输出说明

- `--dry-run` 返回预览请求，不产生实际副作用。
- 正式执行返回远端创建结果，字段以 CLI 实际 JSON 输出为准。
