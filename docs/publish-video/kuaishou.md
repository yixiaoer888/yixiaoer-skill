# 快手视频发布 (Publish Kuaishou Video)

该指令用于通过视频引擎向快手分发短视频内容。

## DTO 溯源 (Knowledge from KuaiShouVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaishou.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 否 | 视频标题 | 不可为空 |
| `--content` | string | 否 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 封面图 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--visibleType` | number | 是 | **可见性** | `0`:公开 `1`:私密 `3`:好友可见 |
| `--declaration` | number | 否 | **视频声明** | `1`:AI生成 `2`:演绎情节 `3`:个人观点 |
| `--location` | object | 否 | 地理位置 | `PlatformDataItem` |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--nearby_show` | boolean| 否 | **是否同城展示** | 默认 `true` |
| `--allow_same_frame`| boolean| 否 | **允许同框** | 默认 `false` |
| `--allow_download` | boolean| 否 | **允许下载** | 默认 `false` |
| `--shopping_cart` | object | 否 | 关联商品 | `object` |
| `--mini_app` | object | 否 | 挂载小程序 | 与购物车互斥 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="快手发布测试" \
  --content="这是快手视频的描述内容" \
  --platforms="快手" \
  --account_ids="ks_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/cover.jpg" \
  --visibleType=0 \
  --nearby_show=true
```
