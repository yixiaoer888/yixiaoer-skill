# 抖音图文发布 (Publish Douyin Image-Text)

该指令用于通过图文/动态引擎向抖音分发内容，支持抖音要求的标题、描述（含话题/好友）、地理位置、背景音乐及合集功能。

## DTO 溯源 (Knowledge from DouYinDynamicForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `image-text`** | 业务模态识别 |
| `--title` | string | 否 | 抖音标题 | 对应 `title` 字段 |
| `--content` | string | 是 | **内容描述** | 对应 `description`。支持 HTML 标签（话题、好友等）。最大 1000 字符。 |
| `--image_urls` | string | 是 | **图片列表** | 逗号分隔的 URL。引擎会自动上传并映射为 `images` 数组。 |
| `--location` | json | 否 | **地理位置** | 来源于地址搜索接口。包含 `id`, `name` 等。注意与其它挂载项可能互斥。 |
| `--musice` | json | 否 | **背景音乐** | **注意：字段名为 `musice`**。来源于音乐搜索接口。 |
| `--collection` | json | 否 | **合集信息** | 如果账号开通了合集功能，可传入合集对象。 |
| `--sub_collection` | json | 否 | **合集选集** | 合集下的二级分类。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。需大于当前时间且符合平台限制。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

## 描述格式详解 (Description Formatting)

根据 DTO 定义，`description` 字段支持 HTML 增强格式，用于插入话题和 @好友：

### 1. 插入话题 (Topics)
```html
<p>记录美好生活 <topic text='搞笑' raw='{"yixiaoerId":"...","raw":{"topic":"搞笑"}}'>#搞笑</topic></p>
```

### 2. @好友 (Mention Friends)
```html
<p>今天天气真好 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend></p>
```

## 调用指令示例 (Usage)

### 1. 发布带话题和音乐的图文
```bash
node scripts/publish.ts \
  --type=image-text \
  --title="我的周末生活" \
  --content="<p>周末去爬山啦！<topic text='生活' raw='{\"raw\":{\"topic\":\"生活\"}}'>#生活</topic></p>" \
  --platforms="抖音" \
  --account_ids="douyin_acc_123" \
  --image_urls="https://img.com/1.jpg,https://img.com/2.jpg" \
  --musice="{\"id\":\"music_001\",\"name\":\"欢快背景音\"}"
```

### 2. 发布带地理位置的图文（保存为草稿）
```bash
node scripts/publish.ts \
  --type=image-text \
  --content="在武康路打卡中..." \
  --platforms="抖音" \
  --account_ids="douyin_acc_123" \
  --image_urls="https://img.com/wukang.jpg" \
  --location="{\"id\":\"loc_456\",\"name\":\"武康路历史文化名街\"}" \
  --pubType=0
```

## 逻辑与规范说明
- **字段转换**: `publish.ts` 会自动将 `--content` 的值映射到 `description` 字段，并将 `--image_urls` 转换后的资源 Key 注入到 `images` 数组中。
- **拼写注意**: 后端 DTO 中音乐字段名为 `musice`，调用时请务必使用 `--musice` 而非 `--music`。
- **图片约束**: 抖音图文建议上传 3-35 张图片，单张大小限制请参考平台文档。引擎默认设置图片尺寸为 1200x800。
- **互斥规则**: 地理位置、带货地址、小程序挂载等在抖音平台可能存在互斥逻辑，请确保输入合法。
