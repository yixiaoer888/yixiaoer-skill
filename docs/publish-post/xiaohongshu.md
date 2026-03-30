# 小红书图文发布 (Publish XiaoHongShu Image-Text)

该指令用于通过图文/动态引擎向小红书发布笔记，支持小红书要求的标题、带有话题标签的 HTML 描述、可见性设置、多图上传、地理位置、背景音乐及合集功能。

## DTO 溯源 (Knowledge from XiaoHongShuDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/xiaohongshu.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--title` | string | 否 | 笔记标题 | 最大 20 字符。对应 `title` 字段。 |
| `--content` | string | 是 | **笔记正文** | 对应 `description`。支持 HTML 增强格式（见下文）。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎自动上传并映射为 `images` 数组。 |
| `--visibleType` | number | 是 | **可见性** | `0`: 公开, `1`: 私密, `3`: 好友可见。 |
| `--location` | json | 否 | **地理位置** | 对象格式。通过 [查询地理位置](#) 接口获取。包含 `id`, `text`, `raw` 字段。 |
| `--music` | json | 否 | **背景音乐** | 对象格式。包含 `id`, `text`, `raw` 字段。 |
| `--scheduledTime` | number | 否 | **定时发布时间** | Unix 时间戳 (秒)。 |
| `--collection` | json | 否 | **合集信息** | 对象格式。如果账号开通了合集功能。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

---

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段（对应 `--content`）支持 HTML 增强格式，用于插入话题标签、@好友及参与活动：

### 1. 插入话题 (Topics)
```html
<p>今天天气真好 <topic text='好心情' raw='{"yixiaoerId":"...","raw":{"topic":"好心情"}}'>#好心情</topic></p>
```

### 2. @好友 (Mention Friends)
```html
<p>记录生活 <friend raw='{"yixiaoerId":"123","raw":{"nick":"张三"}}'>@张三</friend></p>
```

### 3. 参与活动 (Activities)
```html
<p>打卡成功 <activity raw='{"yixiaoerId":"789","raw":{}}'>参与活动</activity></p>
```

---

## 调用指令示例 (Usage)

### 1. 立即发布一篇带话题的小红书笔记
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="小红书" \
  --account_ids="xhs_acc_001" \
  --title="今日穿搭分享" \
  --content="<p>这套衣服真的超级出片！<topic text='OOTD' raw='{\"raw\":{\"topic\":\"OOTD\"}}'>#OOTD</topic></p>" \
  --image_urls="https://example.com/img1.jpg,https://example.com/img2.jpg" \
  --visibleType=0
```

### 2. 发布定时私密笔记
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="小红书" \
  --account_ids="xhs_acc_001" \
  --content="<p>有些话只想对自己说。</p>" \
  --image_urls="https://example.com/private.jpg" \
  --visibleType=1 \
  --scheduledTime=1743382800
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **强制约束**: **`visibleType` 是小红书图文必传字段**。如果不指定，后端逻辑可能会报错。
- **图片设置**: 每个图片会由引擎默认补全 `width: 1200, height: 800`。
- **定时发布**: `--scheduledTime` 应为 Unix 时间戳（秒）。
