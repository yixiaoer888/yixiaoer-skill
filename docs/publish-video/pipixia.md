# 皮皮虾视频发布 (Publish Pipixia Video)

该指令用于通过视频引擎向皮皮虾分发内容。

## DTO 溯源 (Knowledge from PiPiXiaVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/pipixia.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--content` | string | 否 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --content="这是皮皮虾视频的好笑评论" \
  --platforms="皮皮虾" \
  --account_ids="ppx_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg"
```
