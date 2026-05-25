# yxer CLI 命令参考

`yxer` 是本技能唯一默认执行入口。Agent 和用户都应优先使用它，而不是直接调用 `node scripts/*.ts`。

## 核心命令

```bash
yxer doctor
yxer accounts [platform] [--name 关键词] [--status 1] [--json]
yxer upload <file_path_or_url> [--bucket cloud-publish|material-library]
yxer validate <platform> <type> <payload.json>
yxer publish <type> <platform> <payload.json> [clientId]
yxer draft save <payload.json>
yxer material create <payload.json>
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词] [--type 0|1|2|3]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
yxer records [--platform P] [--limit N] [--status S] [--json]
yxer prepare <platform> <type>
yxer schema get <platform> <type>
```

## 基本约束

- 发布类型统一使用：`video`、`image-text`、`article`
- 单次 `yxer publish` 只处理一个平台
- 本地发布时传入 `clientId`
- `yxer draft save` 只处理蚁小二内部草稿，不等同于平台草稿箱
- `yxer material create` 只做素材登记，前提是资源已经通过 `yxer upload --bucket material-library` 上传
- 查询类操作可以直接执行
- 发布类操作必须遵守“查账号 -> 上传资源 -> 查询复杂对象 -> validate -> publish”顺序

## 输出约定

- 默认输出适合人读
- 加 `--json` 时输出结构化结果，适合 Agent 二次处理
- 错误通过统一错误 envelope 输出

## 与旧入口的关系

- `scripts/yxer.ts` 和 `scripts/api.ts` 不再作为默认主入口
- 旧入口只作为兼容和补漏通道，详情见 `../legacy/node-compat.md`
