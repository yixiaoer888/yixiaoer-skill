# 迁移中能力索引

下面这些能力还没有收口成独立的 `yxer` CLI 子命令，因此当前文档只保留迁移提示，不再推荐 `node scripts/api.ts` 直调。

## 当前仍在迁移中的能力

- 账号数据概览：`docs/get-account-overviews.md`
- 作品数据概览：`docs/get-content-overviews.md`
- 征文活动：`docs/get-publish-activities.md`
- 小程序列表：`docs/get-miniapps.md`
- 同步发布应用：`docs/get-sync-apps.md`
- 热点列表：`docs/get-hot-events.md`
- 群聊列表：`docs/get-groups.md`
- 音乐分类：`docs/get-music-categories.md`
- 游戏挂载：`docs/get-games.md`
- 素材库入库（高级登记场景已由 `yxer material create` 覆盖；一步式 `material add` 仍待补齐）：`docs/material-resource.md`
- 代理管理：`docs/proxy-management.md`
- 草稿保存（蚁小二内部草稿已由 `yxer draft save` 覆盖；平台草稿仍走 `yxer publish`）：`docs/save-draft.md`

## 处理原则

1. 这些文档目前只作为字段说明和业务约束参考。
2. 新增或重构时，优先补 `yxer` CLI，不再扩展旧 Node 主流程。
3. 完成 CLI 化后，应从本索引移除，并把文档示例切成真实 `yxer` 命令。

## 推荐迁移顺序

1. `material add`
2. `miniapps` / `syncapps`
3. `games` / `hot-events` / `groups`
4. `account-overviews` / `content-overviews`
5. `proxy-management`
