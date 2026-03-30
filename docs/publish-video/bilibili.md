# 哔哩哔哩视频发布 (Publish Bilibili Video)

该指令用于通过视频引擎向 B 小站分发视频内容。

## DTO 溯源 (Knowledge from BilibiliVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/bilibili.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--description` | string | 是 | 视频描述 | 不可为空 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--tags` | array | 是 | **视频标签** | 字符串数组 |
| `--category` | array | 是 | **分区/分类** | `CascadingPlatformDataItem[]` |
| `--declaration` | number | 是 | **创作者声明** | `0`:不申明 `1`:AI合成 `2`:危险行为 `3`:娱乐演绎 `4`:谨慎观看 `5`:理财适度 `6`:个人观点 |
| `--type` | number | 是 | **发布类型** | `1`:自制 `2`:转载 |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--contentSourceUrl`| string | 否 | **原文链接** | 若 `--type=2` (转载) 时必填 |
| `--collection` | object | 否 | 合集信息 | `object` |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="B站视频测试" \
  --description="这是B站视频的内容简介" \
  --platforms="哔哩哔哩" \
  --account_ids="bili_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["动画","自研"]' \
  --category='[{"id":"1","text":"动画"}]' \
  --declaration=0 \
  --type=1
```
