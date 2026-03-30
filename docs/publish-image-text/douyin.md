# 抖音图文发布 (Douyin Image-Text)

支持发布抖音图文（动态）。

## 指令定义 (Capability)

- **类型**: `image-text`
- **平台**: `抖音`
- **对应 DTO**: `DouYinDynamicForm`

## 参数说明 (Parameters)

| 参数 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `--title` | `string` | 否 | 抖音图文标题 |
| `--content` | `string` | 否 | 图文描述内容。支持简单的 HTML 或话题标签。 |
| `--image_urls` | `url[]` | 是 | 图片链接（逗号分隔）。至少 1 张。 |
| `--location` | `object` | 否 | 地理位置对象。需调用 `get-publish-categories.md` 获取或符合 `PlatformDataItem` 结构。 |
| `--musice` | `object` | 否 | 音乐对象。来源于音乐搜索接口。注意：字段名为 `musice`。 |
| `--scheduledTime` | `number` | 否 | 定时发布时间戳（单位：秒）。 |
| `--collection` | `object` | 否 | 合集信息对象。 |
| `--sub_collection` | `object` | 否 | 合集选集对象。 |

## 对象模型 (Models)

### location / musice (PlatformDataItem)
```json
{
  "id": "123",
  "text": "名称",
  "raw": { ... }
}
```

## 使用示例 (Usage)

```bash
# 发布即时图文
pnpm skill:publish \
  --type=image-text \
  --platforms="抖音" \
  --account_ids="your_account_id" \
  --title="我的周末生活" \
  --content="今天天气真不错 #生活记录" \
  --image_urls="https://example.com/1.jpg,https://example.com/2.jpg"
```
