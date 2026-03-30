# 视频号图文发布 (Publish WeiXinShiPinHao Image-Text)

该指令用于通过图文/动态引擎向微信视频号发布内容，支持视频号要求的标题、带有话题/好友标签的 HTML 描述、图片列表、地理位置、背景音乐及合集功能。

## DTO 溯源 (Knowledge from ShiPingHaoDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/shipinghao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--title` | string | 否 | **内容标题** | 对应 `title`。 |
| `--content` | string | 是 | **图文描述** | 对应 `description`。支持 HTML 增强格式（见下文）。最大 1000 字符。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并映射为 `images` 数组。 |
| `--location` | json | 否 | **地理位置** | 包含 `id`, `text`, `raw` 的 JSON 对象。 |
| `--music` | json | 否 | **背景音乐** | 包含 `id`, `text`, `raw` 的 JSON 对象。 |
| `--collection` | json | 否 | **合集信息** | 对象格式。通过 [查询合集列表](./get-collections.md) 获取。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即发布 (默认 1)。 |

---

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段（对应 `--content`）支持 HTML 增强格式，用于插入话题标签、@好友及参与活动：

### 1. 插入话题 (Topics)
最多支持 5 个话题，每个话题最大 500 字符。
```html
<p>这是一个关于 <topic text='搞笑' raw='{"yixiaoerId":"...","yixiaoerName":"搞笑","raw":{"id":"...","topic":"搞笑"}}'>#搞笑</topic> 的内容</p>
```

### 2. @好友 (Mention Friends)
```html
<p>记录生活 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend></p>
```

### 3. 参与活动 (Activities)
```html
<p>打卡活动 <activity raw='{"yixiaoerId":"...","raw":{}}'>参与活动</activity></p>
```

---

## 调用指令示例 (Usage)

### 1. 立即发布一篇带话题的视频号图文
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="视频号" \
  --account_ids="sph_acc_001" \
  --title="今日穿搭分享" \
  --content="<p>这套衣服真的超级出片！<topic text='OOTD' raw='{\"raw\":{\"topic\":\"OOTD\"}}'>#OOTD</topic></p>" \
  --image_urls="https://example.com/img1.jpg,https://example.com/img2.jpg"
```

### 2. 发布定时草稿
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="视频号" \
  --account_ids="sph_acc_001" \
  --content="<p>有些话只想对自己说。</p>" \
  --image_urls="https://example.com/single.jpg" \
  --pubType=0 \
  --scheduledTime=1743382800
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **强制约束**: 描述内容必须使用 `<p>` 标签包裹。
- **资源上传**: 引擎会自动处理 `--image_urls` 中的图片上传至 OSS 并获取资源 Key。
- **定时发布**: `--scheduledTime` 应为 Unix 时间戳（秒）。
