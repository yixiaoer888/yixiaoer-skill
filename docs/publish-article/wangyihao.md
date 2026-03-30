# 网易号文章发布 (Publish WangYiHao Article)

该指令用于通过文章引擎向网易号分发长内容。

## DTO 溯源 (Knowledge from WangYiHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/wangyihao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--declaration` | number | 否 | **创作申明** | `1`:AI生成 `2`:个人原创 `3`:取材网络 `4`:虚构演绎 |
| `--type` (extra)| number | 否 | **原创勾选** | `0`:不勾选 `1`:原创 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

> [!NOTE]
> 引擎在处理 CLI 参数时，若存在多个 `--type` 可能会导致冲突。建议使用时确认参数映射逻辑。若 DTO 中的核心 `type` 为发布模态（article/video/image-text），平台专有参数建议在文档中明确使用。

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="网易号文章标题" \
  --content="<p>文章内容...</p>" \
  --platforms="网易号" \
  --account_ids="wyh_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --declaration=1
```
