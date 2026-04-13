# 📄 素材库保存 动作 参数 (Material Library Save Action)

将已物理上传至蚁小二 OSS 的资源（图片或视频）正式登记到素材库中，生成可在蚁小二全端复用的固定资产记录。

> [!IMPORTANT]
> **串行执行规范**：素材入库是一个两步走的闭环。Agent **必须**先执行 `upload` (建议设置 `bucket: "material-library"`) 获取 Key，然后立即执行本文档定义的 `material` 动作。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户希望将临时的上传资源“资产化”，以便后续在不同任务中多次复用、团队共享，或在蚁小二网页端进行素材管理。
- **典型提示词**：
  - “帮我把这个视频存入我的素材库”
  - “将这个文件入库，起名叫‘周年庆宣传片’”

## 2. 交互协议 (Interactive Protocol)

1. **两步连动提醒**：若用户直接说“上传”，Agent 应询问是否需要同时存入素材库以便日后查看。
2. **元数据嗅探**：Agent 应尽可能从原始文件中提取 `width`, `height` 等元数据，严禁在 `fileName` 中填入随机乱码。
3. **完成反馈**：登记成功后，告知用户素材已在“素材管理”中可见。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `material` | 固定操作码。 |
| **`filePath`** | `string` | **是** | - | 通过 `upload` 动作获取的资源 `key`。 |
| **`fileName`** | `string` | **是** | - | 素材展示名称（建议带上扩展名）。 |
| **`width`** | `number` | **是** | - | 资源宽度 (像素)。 |
| **`height`** | `number` | **是** | - | 资源高度 (像素)。 |
| **`type`** | `string` | **是** | - | 类型：`video` 或 `image`。 |
| `thumbPath` | `string` | 否 | - | 缩略图 `key` (视频建议上传封面图作为缩略图)。 |

## 4. 执行指令示例 (Command)

```bash
# 示例：将一段已上传的周年庆视频登记入库
node scripts/api.ts --payload='{
  "action": "material",
  "filePath": "material-library/v/demo_clip.mp4",
  "fileName": "周年庆正式版.mp4",
  "width": 1080,
  "height": 1920,
  "type": "video"
}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **登记失败：Key 不存在** | `upload` 动作未成功或 `bucket` 选错。 | 确保上传时使用了 `material-library` 桶，并检查 `upload` 的返回结果。 |
| **预览图无法显示** | `thumbPath` 为空或指向了错误的图片资源。 | 如果是视频入库，强烈建议先上传一张封面图并将 Key 作为 `thumbPath` 传入。 |
| **素材名称重合** | 库中已存在同名素材。 | 系统通常支持重名，但建议 Agent 在名称后添加时间戳或版本号以示区分。 |

---
> [!CAUTION]
> **物理存储隔离**：严禁将 `bucket: "cloud-publish"` 下的临时发布 Key 直接用于素材库登记，否则该资源可能会在 7 天后被系统自动清理。
