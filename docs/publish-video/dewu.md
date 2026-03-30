# 得物视频发布 (Publish Dewu Video)

该指令用于通过视频引擎向得物 (Dewu) 分发短视频内容。

## DTO 溯源 (Knowledge from DeWuVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/dewu.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--content` | string | 是 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--declaration` | number | 是 | **创作声明** | `0`:不添加声明 `1`:AI生成 `2`:不含营销推广 `3`:专业运动 `4`:剧情演绎 |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="得物好物视频发布测试" \
  --content="这是得物视频的内容说明" \
  --platforms="得物" \
  --account_ids="dw_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --declaration=2
```
