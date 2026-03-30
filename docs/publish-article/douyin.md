# 抖音文章发布 (Publish DouYin Article)

该指令用于通过文章引擎向抖音分发图文笔记/文章内容。

## DTO 溯源 (Knowledge from DouyinArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章正文 (HTML) | 不可为空 |
| `--description`| string | 否 | 文章描述 | 简要描述 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--headImage` | object | 否 | **文章头图** | Standard Image Object |
| `--music` | object | 否 | 音乐素材 | 来源于音乐接口 |
| `--topics` | array | 否 | 话题列表 | 最多 5 个 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--visibleType` | number | 否 | **可见性** | `0`:公开 `1`:私密 `3`:好友可见 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="抖音图文文章测试" \
  --content="<p>正文内容展示...</p>" \
  --platforms="抖音" \
  --account_ids="dy_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --visibleType=0
```
