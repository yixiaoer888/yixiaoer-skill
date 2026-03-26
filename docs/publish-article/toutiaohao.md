# 头条号文章发布 (Publish Toutiao Article)

该指令用于通过文章引擎向头条号分发长内容，支持头条特有的首发、广告收益设置及创作者申明。

## DTO 溯源 (Knowledge from TouTiaoHaoArticleForm)
*逻辑来源: `apps/server-api/.../toutiaohao.dto.ts`*

### 核心参数 (Command Arguments)
| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | DTO 映射为 `covers` 数组 |
| `--isFirst` | boolean | 否 | 头条首发 | 是否在头条首发 |
| `--advertisement` | number | 否 | 广告收益 | 2: 投放广告赚收益, 3: 不投放（默认 3） |
| `--declaration` | number | 否 | 创作类型说明 | 1:自行拍摄, 2:取自站外, 3:AI生成, 6:虚构, 7:投资, 8:健康 |
| `--pub_type` | number | 否 | 发布类型 | 0: 草稿, 1: 立即发布 |

## 调用指令示例 (Usage)

### 1. 发布原创首发文章
```bash
node scripts/publish-article.ts \
  --title="2024 AI 趋势报告" \
  --content="<p>正文内容...</p>" \
  --platforms="头条号" \
  --account_ids="tt_123" \
  --cover_url="https://example.com/cover.jpg" \
  --isFirst=true \
  --declaration=1 \
  --advertisement=2
```

### 2. 存为头条草稿
```bash
node scripts/publish-article.ts \
  --title="待完善的草稿" \
  --content="<p>初稿内容</p>" \
  --platforms="头条号" \
  --account_ids="tt_123" \
  --pub_type=0
```

## 逻辑说明
- **多封面支持**: 虽然指令通常传递一个 `--cover_url`，但引擎会自动将其包装为头条要求的 `covers` 数组格式。
- **收益逻辑**: 若未指定 `advertisement`，默认按 DTO 缺省值 `3` 执行（不投放广告）。
