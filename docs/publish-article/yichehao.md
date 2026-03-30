# 易车号文章发布 (Publish YiCheHao Article)

该指令用于通过文章引擎向易车号分发长内容。

## DTO 溯源 (Knowledge from YiCheHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/yichehao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 (横版) | 引擎自动上传并映射为 `covers` 数组 |
| `--verticalCovers`| array | 是 | **文章竖版封面** | 对象数组, e.g. `[{"key": "oss_key", ...}]` |
| `--declaration` | number | 否 | **创作申明** | `0`:无 `1`:个人观点 `2`:网络来源 `3`:AI生成 `4`:引用站内 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--allowForward` | boolean | 否 | 允许转发 | 是否允许转发到其他平台 |
| `--allowAbstract`| boolean | 否 | 允许生成摘要 | 是否允许 AI 生成摘要 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="易车号评测文章" \
  --content="<p>这款车的内饰细节处理...</p>" \
  --platforms="易车号" \
  --account_ids="ych_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --verticalCovers='[{"key": "xxx", "width": 800, "height": 1200}]'
```
