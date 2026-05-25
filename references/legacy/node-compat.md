# 旧 Node 入口兼容说明

仓库历史上存在两类 Node 入口：

- `scripts/yxer.ts`
- `scripts/api.ts`

它们现在都不是默认执行路径。

## 当前策略

1. 对外公开能力统一收口到 Go 版 `yxer` CLI。
2. `SKILL.md`、工作流、查询文档、发布文档都应优先引用 `yxer`。
3. 只有当 `yxer` CLI 尚未覆盖某个底层能力时，才允许临时参考旧 Node 脚本的行为或参数。
4. 新功能不得继续扩展旧 Node 主流程；应优先补到 `cmd/`、`internal/` 和 `yxer` CLI。

## 为什么保留

- 便于对照旧实现细节
- 便于迁移期间核对 API 路径和 payload 结构
- 便于在 CLI 尚未补齐时做兼容参考

## 后续方向

- 文档层面逐步去掉 `node scripts/api.ts` 示例
- CLI 能力补齐后，再考虑彻底下线旧脚本
- 仍未完成 CLI 化的能力列表见 `capabilities.md`
