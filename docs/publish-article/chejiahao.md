# 车家号文章发布 (Publish CheJiaHao Article)

该指令用于通过文章引擎向车家号（汽车之家）分发汽车相关内容。

## DTO 溯源 (Knowledge from ChejiahaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/chejiahao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 (横版) | 引擎自动上传并映射为 `covers` 数组 |
| `--verticalCovers`| array | 是 | **文章竖版封面** | 对象数组, e.g. `[{"key": "oss_key", ...}]` |
| `--type` (extra)| number | 否 | **创作类型** | `1`:原创 `3`:首发 `13`:原创首发 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="2026款车型测评" \
  --content="<p>这款车的性能超乎想象...</p>" \
  --platforms="车家号" \
  --account_ids="cjh_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --verticalCovers='[{"key": "xxx", "width": 800, "height": 1200}]' \
  --type=1
```
