# 车家号视频发布 (Publish Chejiahao Video)

该指令用于通过视频引擎向汽车之家 (车家号) 分发视频内容。

## DTO 溯源 (Knowledge from ChejiahaoVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/chejiahao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--content` | string | 是 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--type` | number | 是 | **创作类型** | `1`:原创 `3`:首发 `13`:原创首发 |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="车家号视频测试" \
  --content="这是车家号视频的内容说明" \
  --platforms="车家号" \
  --account_ids="cjh_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --type=1
```
