# 哔哩哔哩文章发布 (Publish BiLiBiLi Article)

该指令用于通过文章引擎向 B 站分发“专栏”文章。

## DTO 溯源 (Knowledge from BiLiBiLiArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/bilibili.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--tags` | string[] | 否 | 标签 | 字符串数组 |
| `--type` (extra)| number | 否 | **原创类型** | `1`:非原创 `2`:原创 (注意 DTO 为 1/2) |
| `--category` | array | 否 | 分类 | 对象列表 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="B站专栏标题" \
  --content="<p>文章内容...</p>" \
  --platforms="哔哩哔哩" \
  --account_ids="bili_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["科技", "教程"]'
```
