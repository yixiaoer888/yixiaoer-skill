# 蜂网视频发布 (Publish Fengwang Video)

该指令用于通过视频引擎向蜂网分发短视频内容。

## DTO 溯源 (Knowledge from FengWangVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/fengwang.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--content` | string | 是 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--tags` | array | 是 | **视频标签** | 字符串数组 |
| `--category` | array | 是 | **视频分类** | `CascadingPlatformDataItem[]` |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="蜂网视频测试" \
  --content="这是蜂网视频的内容说明" \
  --platforms="蜂网" \
  --account_ids="fw_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["资讯","财经"]' \
  --category='[{"id":"10","text":"生活"}]'
```
