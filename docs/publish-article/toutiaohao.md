# 头条号文章发布 (Publish Toutiao Article)

该指令用于通过文章引擎向头条号分发长内容，支持头条特有的首发、广告收益设置及创作者申明。

## DTO 溯源 (Knowledge from TouTiaoHaoArticleForm)
*逻辑来源: `apps/server-api/.../toutiaohao.dto.ts`*

### 核心参数 (Command Arguments)
| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 | 不可为空 |
| `--cover_url` | string | 是 | 封面图 |直连地址，引擎自动上传并映射为 `covers` 数组 |
| `--isFirst` | boolean | 否 | 头条首发 | 是否在头条首发 |
| `--advertisement` | number | 否 | 广告收益 | 2: 投放广告赚收益, 3: 不投放（默认 3） |
| `--declaration` | number | 否 | 创作类型 | 1:自行拍摄, 2:取自站外, 3:AI生成, 6:虚构... |
| `--pubType` | number | 否 | 发布类型 | 0: 草稿, 1: 立即发布 (默认 1) |

## 调用指令示例 (Usage)

### 1. 发布原创首发文章
```bash
node scripts/publish.ts \
  --type=article \
  --title="2024 AI 趋势报告" \
  --content="<p>正文内容...</p>" \
  --platforms="头条号" \
  --account_ids="tt_123" \
  --cover_url="https://example.com/cover.jpg" \
  --isFirst=true \
  --declaration=1 \
  --advertisement=2
```

## 逻辑说明
- **文档驱动**: `publish.ts` 引擎不维护任何 DTO 知识。AI 助手必须依据本指令文档中的“核心参数”提供的 DTO 字段名（如 `pubType` 而非 `pub_type`）直接传参。
- **发布通道**: `type=article` 会自动执行文章预存证 (Storage) 逻辑。
