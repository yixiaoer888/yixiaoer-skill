# 平台参考索引

平台字段和差异说明仍保留在原有 `docs/publish/` 体系中，但技能入口统一从这里进入，避免 Agent 直接在仓库里漫游。

## 通用索引

- 图文：`docs/publish/imageText/index.md`
- 视频：`docs/publish/video/index.md`
- 文章：`docs/publish/article/index.md`

## 使用原则

1. 先走工作流和 `yxer` CLI。
2. 只有在构造平台特有字段时，才继续查具体平台文档。
3. 如果字段来源可以由 `yxer categories`、`yxer locations`、`yxer music` 等命令获得，优先用命令，不要手写对象。

## 后续扩展

后续如果要完全贴近飞书技能目录，可以把 `docs/publish/` 逐步迁移到 `references/platforms/` 下；当前阶段先保留原文档路径，减少一次性搬迁风险。
