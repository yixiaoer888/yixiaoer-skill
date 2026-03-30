# 豆瓣文章发布 (Publish DouBan Article)

该指令用于通过文章引擎向豆瓣分发日记/文章内容。

## DTO 溯源 (Knowledge from DouBanArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/douban.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--type` (extra)| number | 是 | **创作类型** | `0`:不申明 `1`:原创 |
| `--tags` | string[] | 否 | 标签 | 字符串数组 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="豆瓣书评/日记" \
  --content="<p>内容...</p>" \
  --platforms="豆瓣" \
  --account_ids="db_acc_001" \
  --type=1 \
  --tags='["文学", "随笔"]'
```
