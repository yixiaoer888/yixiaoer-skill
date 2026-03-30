# 新浪微博图文发布 (Publish Sina Weibo Image-Text)

该指令用于向新浪微博发布图文/动态内容，支持图片上传、描述（含话题/好友/活动）、地理位置及定时发布。

## DTO 溯源 (Knowledge from XinLangWeiBoDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/xinlangweibo.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--content` | string | 是 | **内容描述** | 对应 `description`。支持 HTML 标签（话题、好友等）。最大 1000 字符。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并映射为 `images` 数组。 |
| `--location` | json | 否 | **地理位置** | 来源于地址搜索接口。包含 `id`, `text`, `raw` 等。 |
| `--statement` | json | 否 | **内容声明** | 格式：`{"type": number}`。1: AI生成, 2: 虚构演绎, 3: 转载, 4: 自主创作。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段支持 HTML 增强格式，用于插入话题、@好友 及 活动：

### 1. 插入话题 (Topics)
```html
<p>这是一个关于 <topic text='搞笑' raw='{"yixiaoerId":"...","yixiaoerName":"搞笑","raw":{"id":"...","topic":"搞笑"}}'>#搞笑</topic> 的内容</p>
```

### 2. @好友 (Mention Friends)
```html
<p>今天天气真好 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend></p>
```

### 3. 参与活动 (Activities)
```html
<p>参加活动啦 <activity raw='{"yixiaoerId":"...","raw":{}}'>#活动名称</activity></p>
```

## 调用指令示例 (Usage)

### 1. 发布带话题和图片的微博
```bash
node scripts/publish.ts \
  --type=image-text \
  --content="<p>分享今天的午餐！<topic text='美食' raw='{\"raw\":{\"topic\":\"美食\"}}'>#美食</topic></p>" \
  --platforms="新浪微博" \
  --account_ids="weibo_acc_123" \
  --image_urls="https://img.com/food1.jpg,https://img.com/food2.jpg"
```

### 2. 发布带地理位置和内容声明的微博（定时发布）
```bash
node scripts/publish.ts \
  --type=image-text \
  --content="在西湖边散步..." \
  --platforms="新浪微博" \
  --account_ids="weibo_acc_123" \
  --image_urls="https://img.com/xihu.jpg" \
  --location="{\"id\":\"loc_789\",\"text\":\"杭州西湖\"}" \
  --statement="{\"type\":4}" \
  --scheduledTime=1711850400
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 的值映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **定时发布**: `scheduledTime` 必须是未来的时间戳。
- **图片限制**: 微博图文建议上传 1-9 张图片。
- **HTML 约束**: 内容应该用 `<p>` 标签包裹。支持多个 `<p>` 标签。
