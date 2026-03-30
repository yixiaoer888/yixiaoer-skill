# 抖音图文发布 (Douyin Image-Text)

支持发布抖音图文（动态/图文快传）。该模态在手机端表现为图片轮播形式，支持配乐、话题标签、地理位置及合集挂载。

## 指令定义 (Capability)

- **类型 (`--type`)**: `image-text`
- **平台 (`--platforms`)**: `抖音`
- **核心 DTO**: `DouYinDynamicForm`
- **发布引擎映射**: `publishType` -> `imageText`

## 参数说明 (Parameters)

### 1. 核心参数 (Core)

| 参数 | 类型 | 必填 | 说明 | DTO 校验规则 |
| :--- | :--- | :--- | :--- | :--- |
| `--title` | `string` | 否 | 抖音图文标题。 | `@IsOptional` |
| `--content` | `string` | 是 | **描述内容**。支持 HTML 格式（`<p>`）及特殊标签（见下文）。 | 后端要求 `1000` 字符以内 |
| `--image_urls` | `url[]` | 是 | **图片链接**。至少 1 张，多张用逗号分隔。 | `@IsNotEmpty` (映射为 `images`) |
| `--cover_url` | `url` | 否 | **封面图链接**。若未指定，引擎自动取首图。 | 自动构建 `covers` 数组 |

### 2. 业务参数 (Business)

这些参数将透传至 `contentPublishForm`：

| 参数 | 类型 | 必填 | 说明 | 特殊要求 |
| :--- | :--- | :--- | :--- | :--- |
| `--location` | `object` | 否 | **地理位置**。来源于 `get-publish-categories` 或符合 `PlatformDataItem`。 | 与购物车/小程序互斥 |
| `--musice` | `object` | 否 | **背景音乐**。来源于音乐搜索接口。 | **注意**: 字段名包含 `e` (`musice`) |
| `--scheduledTime`| `number` | 否 | **定时发布时间**。单位：秒级时间戳。 | 需大于当前时间 |
| `--collection` | `object` | 否 | **合集信息**。属于 `Category` 结构。 | 挂载到指定合集 |
| `--sub_collection`| `object` | 否 | **合集选集**。属于 `Category` 结构。 | |

## 高级特性：描述内容标签 (Rich Tags)

抖音描述 (`--content`) 支持在文本中嵌入特殊互动标签，格式如下：

### 1. 话题标签 (`topic`)
```html
<p>探讨 <topic text='搞笑' raw='{"yixiaoerId":"...","yixiaoerName":"搞笑","raw":{...}}'>#搞笑</topic> 的日常 </p>
```

### 2. 艾特好友 (`friend`)
```html
<p>感谢 <friend raw='{"yixiaoerId":"...","raw":{"nick":"张三"}}'>@张三</friend> 的出镜 </p>
```

### 3. 参与活动 (`activity`)
```html
<p>参加 <activity raw='{"yixiaoerId":"...","raw":{}}'>#热门挑战</activity> </p>
```

> [!TIP]
> 引擎会自动为未包裹 `<p>` 标签的普通纯文本添加标签并处理换行。

## 对象模型 (Models)

### PlatformDataItem (location / musice)
```json
{
  "id": "123",
  "text": "地理位置名称/音乐名称",
  "raw": { "platform_data": "..." }
}
```

### Category (collection / sub_collection)
```json
{
  "yixiaoerId": "cid_123",
  "yixiaoerName": "合集名",
  "yixiaoerImageUrl": "http...",
  "yixiaoerDesc": "描述",
  "viewNum": "100",
  "raw": {}
}
```

## 调用示例 (Usage)

```bash
# 基础发布
node scripts/publish.ts \
  --type=image-text \
  --platforms="抖音" \
  --account_ids="acc_douyin_001" \
  --content="记录美好周末 #生活感悟" \
  --image_urls="https://img.com/1.jpg,https://img.com/2.jpg"

# 带地理位置与定时发布
node scripts/publish.ts \
  --type=image-text \
  --platforms="抖音" \
  --account_ids="acc_douyin_001" \
  --content="西湖的风景真美" \
  --image_urls="https://img.com/scenery.jpg" \
  --location='{"id":"loc_001","text":"西湖景区","raw":{}}' \
  --scheduledTime=1735689600
```
