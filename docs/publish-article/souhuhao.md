# 搜狐号文章发布 (Publish Souhuhao Article)

该指令用于通过文章引擎向搜狐号分发长内容，支持搜狐号要求的封面选择、任务定时发布及摘要描述。

## DTO 溯源 (Knowledge from SouHuHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/souhuhao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 (支持标准 HTML 标签) | 不可为空 |
| `--desc` | string| 是 | **文章描述/摘要** | 不可为空。将作为文章摘要展示。 |
| `--cover_url` | string | 是 | 封面图 | 直连地址，引擎自动上传并映射为 `covers` 数组。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

## 调用指令示例 (Usage)

### 1. 立即发布一篇科技类文章
```bash
node scripts/publish.ts \
  --type=article \
  --title="搜狐号新平台接入指南" \
  --content="<p>欢迎使用搜狐号发布功能，本文将详细通过 HTML 展示内容...</p>" \
  --desc="关于搜狐号文章发布能力的详细介绍与使用指南。" \
  --platforms="搜狐号" \
  --account_ids="shh_acc_001" \
  --cover_url="https://assets.example.com/cover.jpg"
```

### 2. 存为草稿并设置定时发布 (2026-04-01 10:00:00)
```bash
node scripts/publish.ts \
  --type=article \
  --title="预发布测试内容" \
  --content="<p>这是一篇定时发布的测试文章。</p>" \
  --desc="定时发布功能的验证。" \
  --platforms="搜狐号" \
  --account_ids="shh_acc_001" \
  --cover_url="https://assets.example.com/cover.jpg" \
  --scheduledTime=1775011200 \
  --pubType=0
```

## 逻辑与规范说明
- **摘要描述 (Description)**: 搜狐号强制要求提供文章摘要（`desc` 字段）。如果未提供 `--desc` 参数，系统将尝试截取 `--content` 中的文本内容作为默认摘要。
- **封面处理 (Cover)**: 虽然 DTO 要求的是 `OldCover[]` 数组，但 `publish.ts` 引擎会自动处理 `--cover_url`，将其映射为符合要求的数组结构。
- **定时发布**: `scheduledTime` 必须大于当前时间的 2 小时后（建议），且格式为秒级 Unix 时间戳。
- **引擎适配**: 搜狐号文章不支持分类和标签（由后端在发布后自动分析或不支持），故无需在指令中传入。
