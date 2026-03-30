# 搜狐视频发布 (Publish Sohu Video)

该指令用于通过视频引擎向搜狐视频分发内容。

## DTO 溯源 (Knowledge from SouHuShiPinVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/souhushipin.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--content` | string | 否 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--tags` | array | 否 | **视频标签** | 字符串数组 |
| `--declaration` | number | 否 | **搜索视频申明** | `0`:无需申明 `3`:AI生成 `4`:虚构演绎 `5`:AI数字人生成 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="搜狐视频测试" \
  --content="这是搜狐视频的精彩片段" \
  --platforms="搜狐视频" \
  --account_ids="shv_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --declaration=3
```
