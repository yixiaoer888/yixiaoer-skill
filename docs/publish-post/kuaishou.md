# 快手图文发布 (Publish Kuaishou Image-Text)

该指令用于通过图文/动态引擎向快手分发内容，支持快手要求的描述（含话题/好友）、图片列表、可见性控制、地理位置、背景音乐及合集功能。

## DTO 溯源 (Knowledge from KuaiShouDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaishou.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--content` | string | 是 | **图文描述** | 对应 `description`。支持 HTML 格式和话题标签。最大 1000 字符，最多 5 个话题。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并映射为 `images` 数组。 |
| `--visibleType` | number | 是 | **可见类型** | `0`: 公开, `1`: 私密, `3`: 好友可见。 |
| `--location` | json | 否 | **地理位置** | 包含 `id`, `text`, `raw` 的 JSON 对象。 |
| `--music` | json | 否 | **背景音乐** | 包含 `id`, `text`, `raw` 的 JSON 对象。 |
| `--collection` | json | 否 | **合集信息** | 包含 `yixiaoerId`, `yixiaoerName`, `yixiaoerImageUrl`, `yixiaoerDesc`, `viewNum`, `raw` 的 JSON 对象。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |
| `--pubType` | number | 否 | 发布状态 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段支持 HTML 增强格式，用于插入话题和 @好友：

### 1. 插入话题 (Topics)
```html
<p>记录美好生活 <topic text='搞笑' raw='{"yixiaoerId":"...","yixiaoerName":"搞笑","raw":{"id":"...","topic":"搞笑"}}'>#搞笑</topic></p>
```

### 2. @好友 (Mention Friends)
```html
<p>和朋友出去玩 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend></p>
```

## 调用指令示例 (Usage)

### 1. 发布带话题和可见性设置的图文
```bash
node scripts/publish.ts \
  --type=image-text \
  --content="<p>快手的友友们好！<topic text='日常' raw='{\"raw\":{\"topic\":\"日常\"}}'>#日常</topic></p>" \
  --platforms="快手" \
  --account_ids="kuaishou_acc_123" \
  --image_urls="https://img.com/1.jpg,https://img.com/2.jpg" \
  --visibleType=0
```

### 2. 发布带地理位置和定时任务的图文
```bash
node scripts/publish.ts \
  --type=image-text \
  --content="今天天气真不错" \
  --platforms="快手" \
  --account_ids="kuaishou_acc_123" \
  --image_urls="https://img.com/landscape.jpg" \
  --location="{\"id\":\"loc_789\",\"text\":\"北京天安门\",\"raw\":{}}" \
  --visibleType=0 \
  --scheduledTime=1735689600
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 的值映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **可见性**: 快手平台必须明确指定 `--visibleType`，通常 `0` 代表公开。
- **定时发布**: `scheduledTime` 如果提供，必须是未来的 Unix 时间戳。
- **图片约束**: 快手图文支持多图发布，引擎默认设置图片尺寸为 1200x800。
- **合集**: 如果账号支持合集，可以使用 `--collection` 传入通过 `get-publish-categories.ts` 获取到的合集对象结构。
