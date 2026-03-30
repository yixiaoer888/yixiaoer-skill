# 知乎图文发布 (Publish ZhiHu Image-Text)

该指令用于通过图文/动态引擎向知乎分发动态（想法），支持知乎要求的标题、带有话题/好友标签的 HTML 描述及多图列表。

## DTO 溯源 (Knowledge from ZhiHuDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--title` | string | 是 | **内容标题** | 对应 `title`。知乎想法通常有简短标题。 |
| `--content` | string | 是 | **内容描述** | 对应 `description`。支持 HTML 增强格式（见下文）。最大 1000 字符。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并映射为 `images` 数组。 |

---

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段（对应 `--content`）支持 HTML 增强格式，用于插入话题标签、@好友及参与活动：

### 1. 插入话题 (Topics)
```html
<p>探讨一个话题 <topic text='人工智能' raw='{"yixiaoerId":"...","raw":{"topic":"人工智能"}}'>#人工智能</topic></p>
```

### 2. @好友 (Mention Friends)
```html
<p>向大佬请教 <friend raw='{"yixiaoerId":"...","raw":{"nick":"大佬名字"}}'>@大佬名字</friend></p>
```

### 3. @活动 (Activities)
```html
<p>参与想法活动 <activity raw='{"yixiaoerId":"...","raw":{}}'>参与活动</activity></p>
```

---

## 调用指令示例 (Usage)

### 1. 发布一条简单的知乎想法
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="知乎" \
  --account_ids="zhihu_acc_001" \
  --title="今日感慨" \
  --content="<p>知乎的朋友们，大家好！<topic text='日常记录' raw='{\"raw\":{\"topic\":\"日常记录\"}}'>#日常记录</topic></p>" \
  --image_urls="https://example.com/img1.jpg,https://example.com/img2.jpg"
```

### 2. 发布带话题的知乎图文
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="知乎" \
  --account_ids="zhihu_acc_001" \
  --title="知乎好物分享" \
  --content="<p>这件产品真的很值得推荐。</p>" \
  --image_urls="https://example.com/product.jpg"
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **强制约束**: 知乎图文发布时，`title`, `description`, `images` 都是必须的字段。
- **资源限制**: 知乎想法单条内容一般支持 1-9 张图片。
