# 抖音视频发布 (Publish DouYin Video)

该指令用于通过视频引擎向抖音分发短视频内容。

## DTO 溯源 (Knowledge from DouYinVideoForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `video`** | 业务模态识别 |
| `--title` | string | 是 | 视频标题 | 不可为空 |
| `--content` | string | 是 | 视频描述 | 映射到 DTO 的 `description` 字段 |
| `--video_url` | string | 是 | 视频文件 URL | 引擎自动上传并映射为 `video` 对象 |
| `--cover_url` | string | 是 | 视频竖版封面 URL | 映射到顶层 `cover` 和 `accForm.cover` |
| `--horizontalCover` | object | 否 | **视频横板封面** | 抖音特有字段，Standard Cover Object |
| `--statement` | number | 否 | **创作声明** | `3`:内容由AI生成 `4`:引人不适 `5`:虚构演绎，仅供娱乐 `6`:危险行为，请勿模仿 |
| `--location` | object | 否 | 地理位置 | 来源于位置接口 `PlatformDataItem` |
| `--scheduledTime` | number | 否 | 定时发布时间 | Unix 时间戳 (秒) |
| `--collection` | object | 否 | 合集信息 | `Category` 对象 |
| `--sub_collection` | object | 否 | 合集选集 | `Category` 对象 |
| `--challenge` | object | 否 | 挑战/抖音话题 | `Category` 对象 |
| `--hot_event` | object | 否 | 热点事件 | `Category` 对象 |
| `--sync_apps` | array | 否 | 同时发布到的应用 | `Category[]` |
| `--allow_save` | number | 否 | **允许保存** | `0`:不允许 `1`:允许 (默认0) |
| `--shoppingCart` | array | 否 | 购物车列表 | `ShoppingCartDTO[]` |
| `--groupShopping` | object | 否 | 团购信息 | `GroupShoppingDTO` |
| `--mini_app` | object | 否 | 挂载小程序 | 互斥逻辑: 购物车和小程序只能二选一 |
| `--music` | object | 否 | 背景音乐 | 来源于音乐接口 |
| `--cooperation_info`| object | 否 | 共创信息 | 对于专业创作者的协助字段 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=video \
  --title="我的抖音视频标题" \
  --content="这是视频的精彩描述 #搞笑 #生活" \
  --platforms="抖音" \
  --account_ids="dy_acc_001" \
  --video_url="https://example.com/video.mp4" \
  --cover_url="https://example.com/vertical_cover.jpg" \
  --horizontalCover='{"key":"oss_key","width":1920,"height":1080}' \
  --allow_save=1
```

> [!NOTE]
> 引擎会自动处理 `horizontalCover`：如果不显式提供，默认将使用 `cover_url` 的 Key 并假设为 1920x1080 尺寸。
