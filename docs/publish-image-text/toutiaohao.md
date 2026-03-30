# 头条号图文发布 (Publish TouTiaoHao Image-Text)

该指令用于通过图文/动态引擎向头条号发布动态（微头条），支持微头条要求的描述（含话题/好友）、多图列表、内容声明及可见性发布。

## DTO 溯源 (Knowledge from TouTiaoHaoDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--content` | string | 是 | **微头条描述** | 对应 `description`。支持 HTML 增强格式（见下文）。最大 1000 字符。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并映射为 `images` 数组。 |
| `--declaration` | number | 否 | **创作类型** | `1`: 自行拍摄, `2`: 取自站外, `3`: AI 生成, `6`: 虚构演绎, `7`: 投资观点, `8`: 健康医疗。 |
| `--pubType` | number | 是 | **发布类型** | `0`: 保存草稿, `1`: 立即发布 (默认 1)。 |

---

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段（对应 `--content`）支持 HTML 增强格式，用于插入话题标签、@好友及参与活动：

### 1. 插入话题 (Topics)
```html
<p>大家觉得这个怎么样？ <topic text='热门讨论' raw='{"yixiaoerId":"...","yixiaoerName":"热门讨论","raw":{"id":"...","topic":"热门讨论"}}'>#热门讨论</topic></p>
```

### 2. @好友 (Mention Friends)
```html
<p>看到了一个很有意思的观点 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend></p>
```

### 3. @活动 (Activities)
```html
<p>参加活动啦 <activity raw='{"yixiaoerId":"...","raw":{}}'>征文活动</activity></p>
```

---

## 调用指令示例 (Usage)

### 1. 立即发布一条简单的微头条
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="头条号" \
  --account_ids="tt_acc_001" \
  --content="<p>今日份的好心情！<topic text='日常' raw='{\"raw\":{\"topic\":\"日常\"}}'>#日常</topic></p>" \
  --image_urls="https://example.com/img1.jpg,https://example.com/img2.jpg"
```

### 2. 发布带原创声明的微头条
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="头条号" \
  --account_ids="tt_acc_001" \
  --content="<p>这是我亲手拍摄的风景照。</p>" \
  --image_urls="https://example.com/scenery.jpg" \
  --declaration=1
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **强制约束**: 头条号图文（微头条）发布时要求必须传入 `pubType`，缺省为 1 (发布)。
- **HTML 格式**: 头条号后端接口对微头条描述的支持和文章支持不同，建议核心文案均包裹在 `<p>` 内。
- **创作声明**: 头条号对特定领域的声明 (如医疗、健康、投资) 要求严格，须正确指定 `declaration` 枚举。
