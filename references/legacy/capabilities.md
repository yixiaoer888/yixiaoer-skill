# 未 CLI 化能力索引

本页用于标记当前 `yxer` CLI 尚未覆盖、或仅部分覆盖的能力。它不是兼容入口说明，也不代表仓库内存在其他可执行主流程。

## 使用规则

1. 如果能力已经有正式 `yxer` 命令，优先直接调用现有命令。
2. 如果能力只有“部分覆盖”，优先使用现有命令完成已支持部分，剩余字段仅把对应文档当参考。
3. 如果能力仍“未覆盖”，应优先补成新的 `yxer` 子命令，而不是寻找旧脚本或新增脚本入口。

## 能力状态表

| 能力 | 当前状态 | 推荐/现有命令 | 参考文档 |
| --- | --- | --- | --- |
| 蚁小二内部草稿保存 | 已覆盖 | `yxer draft save <payload.json>` | `docs/save-draft.md` |
| 平台草稿发布 | 部分覆盖 | `yxer publish <type> <platform> <payload.json>` | `docs/save-draft.md` |
| 素材库登记 | 已覆盖 | `yxer material create <payload.json>` | `docs/material-resource.md` |
| 素材上传并登记一体化 | 已覆盖 | `yxer material add --file <文件路径或URL>` | `docs/material-resource.md` |
| 账号数据概览 | 未覆盖 | 建议新增：`yxer account-overviews ...` | `docs/get-account-overviews.md` |
| 作品数据概览 | 未覆盖 | 建议新增：`yxer content-overviews ...` | `docs/get-content-overviews.md` |
| 征文活动 | 未覆盖 | 建议新增：`yxer activities <account_id> ...` | `docs/get-publish-activities.md` |
| 小程序列表 | 未覆盖 | 建议新增：`yxer miniapps <account_id> ...` | `docs/get-miniapps.md` |
| 同步发布应用 | 未覆盖 | 建议新增：`yxer syncapps <account_id>` | `docs/get-sync-apps.md` |
| 热点列表 | 未覆盖 | 建议新增：`yxer hot-events <account_id> ...` | `docs/get-hot-events.md` |
| 群聊列表 | 未覆盖 | 建议新增：`yxer groups <account_id>` | `docs/get-groups.md` |
| 音乐分类 | 未覆盖 | 建议新增：`yxer music-categories <account_id>` | `docs/get-music-categories.md` |
| 游戏挂载 | 未覆盖 | 建议新增：`yxer games <account_id> ...` | `docs/get-games.md` |
| 代理管理 | 未覆盖 | 建议新增：`yxer proxies ...` / `yxer proxy-areas` / `yxer update-account ...` | `docs/proxy-management.md` |

## 推荐迁移顺序

1. `yxer miniapps` / `yxer syncapps`
2. `yxer games` / `yxer hot-events` / `yxer groups`
3. `yxer account-overviews` / `yxer content-overviews`
4. `yxer activities`
5. `yxer proxies` / `yxer proxy-areas` / `yxer update-account`
