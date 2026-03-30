# 头条号视频发布 (Publish Toutiao Video)

该指令用于通过视频引擎向头条号分发短视频内容。

## DTO 溯源 (Knowledge from TouTiaoHaoVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--content` | string | 是 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--tags` | array | 是 | **视频标签** | 字符串数组，描述视频关键词 |
| `--declaration` | number | 否 | **原创类型** | `1`:自行拍摄 `2`:取自站外 `3`:AI生成 `6`:虚构演绎 `7`:投资观点 `8`:健康医疗 |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="头条视频测试" \
  --content="这是头条视频的描述内容" \
  --platforms="头条号" \
  --account_ids="tt_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["测试","视频"]' \
  --declaration=1
```
