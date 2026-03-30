# 爱奇艺文章发布 (Publish AiQiYi Article)

该指令用于通过文章引擎向爱奇艺分发长内容。

## DTO 溯源 (Knowledge from AiQiYiArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/aiqiyi.dto.ts`*

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
  --title="爱奇艺文章标题" \
  --content="<p>文章内容...</p>" \
  --platforms="爱奇艺" \
  --account_ids="aqy_acc_001"
```
