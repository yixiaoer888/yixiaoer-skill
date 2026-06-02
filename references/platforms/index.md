# 平台参考索引

平台字段和差异说明已收敛到 `skills/yixiaoer/references/platforms/`。根目录索引只作为仓库维护入口，避免 Agent 在多套路径之间摇摆。

## 通用索引

- 图文：`skills/yixiaoer/references/platforms/imageText/index.md`
- 视频：`skills/yixiaoer/references/platforms/video/index.md`
- 文章：`skills/yixiaoer/references/platforms/article/index.md`

## 使用原则

1. 先走工作流和 `yxer` CLI。
2. 只有在构造平台特有字段时，才继续查具体平台文档。
3. 如果字段来源可以由 `yxer categories`、`yxer locations`、`yxer music` 等命令获得，优先用命令，不要手写对象。

## 维护约束

- Agent 运行时应优先读取 skill 包内部的 `skills/yixiaoer/references/platforms/`
- 根目录文档用于仓库维护和交叉索引，不再作为 skill 主入口
