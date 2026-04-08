# 上传到素材库 (Material Library Save)

将已上传到蚁小二 OSS 的资源登记到素材库，生成可在素材管理中复用的视频或图片记录。

## 场景描述 (Usage)

- "把这个已经上传好的视频存进素材库。"
- "我已经拿到资源 key 了，帮我入库到团队素材库分组。"

> [!IMPORTANT]
> 素材库接口 `action: "material"` 不负责原始文件上传。请先调用 [upload-resource.md](./upload-resource.md) 获取 `key`，再把该 `key` 作为 `filePath` 传给当前接口。
> 若目标是素材库，请在上传步骤显式使用 `bucket: "material-library"`。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{
  "action": "material",
  "filePath": "cloud-publish/2026/04/08/demo/video.mp4",
  "thumbPath": "cloud-publish/2026/04/08/demo/cover.jpg",
  "fileName": "演示视频.mp4",
  "width": 1080,
  "height": 1920,
  "type": "video"
}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | 是 | 固定值: `material` |
| `filePath` | `string` | **是** | 已上传到蚁小二 OSS 的资源 `key`，不是外部 URL |
| `fileName` | `string` | **是** | 展示用文件名，建议保留扩展名 |
| `width` | `number` | **是** | 素材宽度 |
| `height` | `number` | **是** | 素材高度 |
| `type` | `string` | **是** | 素材类型。常用值为 `video` 或 `image` |
| `thumbPath` | `string` | 否 | 缩略图资源 `key`。视频建议传封面，图片可不传 |

> [!NOTE]
> 当前龙虾技能的 `material` action 不支持指定素材分组，素材会按后端默认行为进入“全部”分组。

## 推荐链路 (Recommended Flow)

1. 先调用 `action: "upload"` 上传原始文件，获得资源 `key`
2. 再调用 `action: "material"`，把 `key` 写入素材库
3. 后续在发布能力中直接复用素材库记录或对应资源 `key`

上传步骤建议示例：
```bash
node scripts/api.ts --payload='{"action":"upload","url":"https://example.com/video.mp4","bucket":"material-library","contentType":"video/mp4"}'
```

## 输出结果 (Output)

成功时返回素材库记录，通常包含素材 ID、访问地址、资源 key、分组信息和创建时间：

```json
{
  "id": "67f4cb0c1d593b4ad42c8aa1",
  "type": "video",
  "filePath": "https://cdn.example.com/cloud-publish/2026/04/08/demo/video.mp4",
  "fileKey": "cloud-publish/2026/04/08/demo/video.mp4",
  "thumbKey": "cloud-publish/2026/04/08/demo/cover.jpg",
  "fileName": "演示视频.mp4",
  "createdAt": 1775635200000
}
```

> [!TIP]
> 若同一个 `sourceTaskId` 或相同资源已入库，后端可能直接返回已有素材记录，而不是重复创建新素材。
