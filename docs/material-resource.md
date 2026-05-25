# 上传到素材库 (Material Library Save)

将已上传到蚁小二 OSS 的资源登记到素材库，生成可在素材管理中复用的视频或图片记录。

## 触发场景 (Trigger)
- **意图辨析**：当用户希望将已上传的临时资源“持久化”或“资产化”，以便后续多次复用、团队共享或在网页端素材库中进行统一管理时触发。
- **典型提示词**：
  - “把这个视频存进素材库方便以后调用”
  - “我的宣传片已经传好了，帮我入库”
  - “将这个 Key 对应的图片登记到素材管理中”

> [!CAUTION]
> **强制两步流程 (Mandatory Two-Step)**: 
> 当用户意图为“上传到素材库”时，**禁止单元操作**。Agent 必须严格执行：
> 1. **upload** (上传到 OSS 并获取 Key，建议使用 `bucket: "material-library"`)
> 2. **material** (将获得的 Key 登记到库，产出素材 ID)
> **错误演示**: 只执行 `upload` 而不执行 `material` 属于逻辑截断，会导致用户看不见素材。

## 参数定义 (Parameters)

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | **是** | 固定值: `material` |
| `filePath` | `string` | **是** | 已上传到蚁小二 OSS 的资源 `key`（通过 `upload` 获取） |
| `fileName` | `string` | **是** | 展示用文件名，建议保留扩展名 |
| `width` | `number` | **是** | 素材宽度 |
| `height` | `number` | **是** | 素材高度 |
| `type` | `string` | **是** | 素材类型。常用值为 `video` 或 `image` |
| `thumbPath` | `string` | 否 | 缩略图资源 `key`。视频建议传封面，图片可不传 |

## 执行逻辑 (Logic Flow)
1. **链路检查**：确认资源是否已通过 `upload` 动作上传（建议使用 `bucket: "material-library"`）。
2. **元数据提取**：识别文件的真实名称、宽高及类型。
3. **参数装配**：构造 `action: "material"` 及完整元数据。
4. **指令执行**：先调用 `yxer upload --bucket material-library`，再调用 `yxer material create <payload.json>`。
5. **入库反馈**：向用户确认素材 ID 及入库成功状态。

## 推荐链路 (Recommended Flow)

1. 先调用 `action: "upload"` 上传原始文件，获得资源 `key`
2. 再调用 `action: "material"`，把 `key` 写入素材库

## 调用指令 (Command)

```bash
yxer material create material-payload.json
```

推荐先执行：

```bash
yxer upload ./demo.mp4 --bucket material-library
```

再将返回的 `key` 组装进 `material-payload.json`：

```json
{
  "filePath": "uploaded/demo.mp4",
  "fileName": "demo.mp4",
  "width": 1080,
  "height": 1920,
  "type": "video"
}
```

## 注意事项
- **Bucket 匹配**：若目标是素材库，请在上传步骤显式使用 `bucket: "material-library"`。
- **重复校验**：若相同资源已入库，后端可能返回已有记录。
