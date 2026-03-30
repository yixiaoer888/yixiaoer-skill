# CSDN文章发布 (Publish CSDN Article)

该指令用于通过文章引擎向 CSDN 分发技术博客内容。

## DTO 溯源 (Knowledge from CSDNArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/csdn.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML/Markdown 为准) | 不可为空 |
| `--desc` | string | 是 | **文章摘要** | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--tags` | string[] | 是 | 标签 | 字符串数组 |
| `--type` (extra)| number | 是 | **创作类型** | `1`:原创 `2`:转载 `4`:翻译 |
| `--contentSourceUrl`| string | 否 | 原文链接 | 转载时必须填写 |
| `--declaration` | number | 否 | **创作申明** | `0`:无 `1`:AI辅助 `2`:整合 `3`:个人观点 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="React 进阶技巧" \
  --content="<p>深入探讨 React Context...</p>" \
  --desc="本文详细介绍了 React Context 的性能优化方案。" \
  --platforms="CSDN" \
  --account_ids="csdn_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["React", "前端", "JavaScript"]' \
  --type=1
```
