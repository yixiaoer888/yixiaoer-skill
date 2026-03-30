# 百家号图文发布 (Publish BaiJiaHao Image-Text)

该指令用于通过图文/动态引擎向百家号发布动态内容（微百家号），支持百家号要求的标题、描述（含话题/好友/活动）、封面图片、地理位置及创作声明。

## DTO 溯源 (Knowledge from BaiJiaHaoDynamicForm)
*逻辑来源: `apps/server-api-packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--title` | string | 是 | **动态标题** | 对应 `title`。必填。 |
| `--content` | string | 是 | **图文描述** | 对应 `description`。支持 HTML 增强格式（见下文）。最大 1000 字符。 |
| `--cover_url` | string | 是 | **封面图片** | 对应 `cover` 字段。引擎会自动上传。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并注入到 `accountForms.images`。 |
| `--location` | json | 否 | **地理位置** | 包含 `id`, `text`, `raw` 的 JSON 对象。 |
| `--declaration` | number | 是 | **创作声明** | `0`: 不声明, `1`: 内容由 AI 生成。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |

---

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段（对应 `--content`）支持 HTML 增强格式，用于插入话题标签、@好友及参与活动：

### 1. 插入话题 (Topics)
```html
<p>发现一个好地方 <topic text='旅游热点' raw='{"yixiaoerId":"...","yixiaoerName":"旅游热点","raw":{"id":"...","topic":"旅游热点"}}'>#旅游热点</topic></p>
```

### 2. @好友 (Mention Friends)
```html
<p>和朋友一起 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend></p>
```

### 3. @活动 (Activities)
```html
<p>参与活动 <activity raw='{"yixiaoerId":"...","raw":{}}'>征文活动</activity></p>
```

---

## 调用指令示例 (Usage)

### 1. 立即发布一篇带话题和声明的百家号动态
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="百家号" \
  --account_ids="bjh_acc_001" \
  --title="周末出游记" \
  --content="<p>今天天气真好，出来透透气。<topic text='周末去哪儿' raw='{\"raw\":{\"topic\":\"周末去哪儿\"}}'>#周末去哪儿</topic></p>" \
  --image_urls="https://example.com/img1.jpg,https://example.com/img2.jpg" \
  --cover_url="https://example.com/img1.jpg" \
  --declaration=0
```

### 2. 发布 AI 生成的定时动态
```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms="百家号" \
  --account_ids="bjh_acc_001" \
  --title="AI 绘画分享" \
  --content="<p>这是由 AI 生成的一组图片，感觉效果不错。</p>" \
  --image_urls="https://example.com/ai_img.jpg" \
  --cover_url="https://example.com/ai_img.jpg" \
  --declaration=1 \
  --scheduledTime=1743382800
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 映射到 `description` 字段。
- **强制约束**: `title`, `description`, `cover`, `declaration` 均为百家号图文的必填字段。
- **封面说明**: 百家号动态要求必须指定一张封面图。
- **定时发布**: `--scheduledTime` 应为 Unix 时间戳（秒）。
