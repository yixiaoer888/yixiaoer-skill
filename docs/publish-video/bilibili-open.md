# 哔哩哔哩-Open 视频发布 (Publish Bilibili-Open Video)

该指令用于通过 B 站开放接口 (Open API) 分发视频内容。其参数结构与通用 [哔哩哔哩视频发布](./bilibili.md) 完全一致。

## DTO 溯源 (Knowledge from BilibiliVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/bilibili.dto.ts`*

### 核心参数 (Command Arguments)

详见：[哔哩哔哩视频发布参数说明](./bilibili.md#核心参数-command-arguments)

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="B站-Open视频测试" \
  --description="这是通过开放接口发布的视频描述" \
  --platforms="哔哩哔哩-Open" \
  --account_ids="bili_open_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["测试","Open"]' \
  --category='[{"id":"1","text":"动画"}]' \
  --declaration=0 \
  --type=1
```
