# 快传号文章发布 (Publish KuaiChuanHao Article)

该指令用于通过文章引擎向快传号分发长内容。

## DTO 溯源 (Knowledge from KuaiChuanHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaichuanhao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--tags` | string[] | 否 | 标签 | 字符串数组 |
| `--type` (extra)| number | 是 | **创作类型** | `0`:不申明 `1`:原创 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="快传号文章标题" \
  --content="<p>文章内容...</p>" \
  --platforms="快传号" \
  --account_ids="kch_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --type=1
```
