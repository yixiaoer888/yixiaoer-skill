# 多多视频发布 (Publish Duo Duo Video)

该指令用于通过视频引擎向拼多多 (多多视频) 分发短视频内容。

## DTO 溯源 (Knowledge from DuoDuoShiPinVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/duoduoshipin.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--content` | string | 否 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--shopping_cart` | object | 否 | 关联商品 | `object` |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --content="多多视频好物分享" \
  --platforms="多多视频" \
  --account_ids="pdd_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg"
```
