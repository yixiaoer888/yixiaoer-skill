# 简书文章发布 (Publish JianShu Article)

该指令用于通过文章引擎向简书分发随笔内容。

## DTO 溯源 (Knowledge from JianShuArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/jianshu.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="简书随笔标题" \
  --content="<p>创作你的创作...</p>" \
  --platforms="简书" \
  --account_ids="js_acc_001"
```
