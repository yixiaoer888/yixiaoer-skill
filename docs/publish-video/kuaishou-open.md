# 快手-Open 视频发布 (Publish Kuaishou-Open Video)

该指令用于通过快手开放接口 (Open API) 分发视频内容。其参数结构与通用 [快手视频发布](./kuaishou.md) 基本一致。

## DTO 溯源 (Knowledge from KuaiShouVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaishou.dto.ts`*

### 核心参数 (Command Arguments)

详见：[快手视频发布参数说明](./kuaishou.md#核心参数-command-arguments)

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="快手-Open发布测试" \
  --content="这是通过快手开放接口发布的描述内容" \
  --platforms="快手-Open" \
  --account_ids="ks_open_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --visibleType=0
```
