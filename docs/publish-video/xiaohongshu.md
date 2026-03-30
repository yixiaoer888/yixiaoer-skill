# 小红书视频发布 (Publish Xiaohongshu Video)

该指令用于通过视频引擎向小红书分发短视频笔记。

## DTO 溯源 (Knowledge from XiaoHongShuVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/xiaohongshu.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 否 | 笔记标题 | 最多 20 个字 |
| `--content` | string | 否 | 笔记描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--visibleType` | number | 是 | **可见性** | `0`:公开 `1`:私密 `3`:好友可见 |
| `--declaration` | number | 否 | **内容类型申明** | `1`:虚构演绎 `2`:AI合成 |
| `--type` | number | 否 | **创作类型** | `1`:原创 `0`:不申明 |
| `--location` | object | 否 | 地理位置 | `PlatformDataItem` |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--collection` | object | 否 | 合集信息 | `object` |
| `--shopping_cart` | array | 否 | 商品信息 | `object[]` |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="小红书视频测试" \
  --content="这是视频笔记的描述 #生活 #日常" \
  --platforms="小红书" \
  --account_ids="xhs_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --visibleType=0 \
  --type=1
```
