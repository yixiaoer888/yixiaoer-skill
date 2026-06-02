# 草稿与素材库

适用范围：用户要保存草稿、登记素材、上传素材到素材库，或先存资源再继续发布。

## 读取顺序

1. 草稿任务：[`../workflows/draft-workflow.md`](../workflows/draft-workflow.md)
2. 素材任务：[`../workflows/material-workflow.md`](../workflows/material-workflow.md)
3. 若随后要直接发布，再切回 [`./publish.md`](./publish.md)

## 常用命令

```bash
yxer draft save <payload.json> [--dry-run]
yxer material add --file <文件路径或URL> [--thumb <缩略图路径或URL>] [--type image|video|file] [--dry-run]
yxer upload --file <file_path> --bucket material-library
yxer upload --url <resource_url> --bucket material-library
yxer material create <payload.json> [--dry-run]
```

## 决策规则

- 用户只说“保存草稿”时，先区分蚁小二草稿和平台草稿，不要自行猜测
- 用户只想把资源放进素材库时，优先 `material add`
- 用户已有上传结果，只差登记素材时，再用 `material create`
- 素材任务不自动进入发布主流程；只有用户明确要“上传后马上发布”时，再切回发布域
