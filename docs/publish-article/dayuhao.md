# 大鱼号文章发布 (Publish DaYuHao Article)

该指令用于通过文章引擎向大鱼号分发长内容。

## DTO 溯源 (Knowledge from DaYuHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/dayuhao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 (横版) | 引擎自动上传并映射为 `covers` 数组 |
| `--verticalCovers`| array | 是 | **竖版封面** | 对象数组, e.g. `[{"key": "oss_key", ...}]` |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="大鱼号文章标题" \
  --content="<p>文章内容...</p>" \
  --platforms="大鱼号" \
  --account_ids="dyh_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --verticalCovers='[{"key": "xxx", "width": 800, "height": 1200}]'
```
