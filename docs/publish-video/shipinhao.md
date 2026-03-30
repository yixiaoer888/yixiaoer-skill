# 视频号视频发布 (Publish Shipinghao Video)

该指令用于通过视频引擎向微信视频号分发短视频内容。

## DTO 溯源 (Knowledge from ShiPingHaoVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/shipinghao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 否 | 视频标题 | 不可超过 200 字 |
| `--short_title` | string | 否 | 视频短标题 | 可选 |
| `--content` | string | 否 | 视频描述 | HTML 格式，映射到 `description` |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--type` | number | 是 | **创作类型** | `1`:非原创 `2`:原创 |
| `--horizontalCover`| object | 否 | **横版封面** | `OldCover` 对象 |
| `--location` | object | 否 | 地理位置 | `PlatformDataItem` |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--shoppingCart` | object | 否 | 关联商品 | `object` |
| `--collection` | object | 否 | 合集信息 | `object` |
| `--activity` | object | 否 | 活动信息 | `object` |
| `--pubType` | number | 否 | **发布模式** | `0`:存草稿 `1`:直接发布 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="视频号发布测试" \
  --content="<p>这是视频号的内容描述 #生活</p>" \
  --platforms="视频号" \
  --account_ids="sph_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --type=2 \
  --pubType=1
```
