# 腾讯视频发布 (Publish Tencent Video)

该指令用于通过视频引擎向腾讯视频分发长/短视频内容。

## DTO 溯源 (Knowledge from TengXunShiPinVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/tengxunshipin.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--tags` | array | 否 | **视频标签** | 字符串数组 |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--declaration` | number | 是 | **视频原创类型** | `1`:AI生成 `2`:剧情演绎 `3`:取材网络 `4`:个人观点 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="腾讯视频发布测试" \
  --platforms="腾讯视频" \
  --account_ids="txv_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --declaration=4
```
