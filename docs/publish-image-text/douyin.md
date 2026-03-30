# 抖音图文发布 (Douyin Image-Text Publishing)

该文档定义了在抖音平台发布图文（动态）能力的参数规范。

## 核心指令 (Command)

```bash
node scripts/publish.ts --type=image-text --platforms=抖音 --account_ids=ACCOUNT_ID --title="标题" --description="描述内容" --image_urls="URL1,URL2"
```

## 参数列表 (Properties)

| 参数名 | 类型 | 是否必填 | 说明 | DTO 校验规则 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | `string` | 是 | 固定值: `image-text` | - |
| `--platforms` | `string` | 是 | 固定值: `抖音` | - |
| `--account_ids` | `string` | 是 | 账号 ID，多个请用逗号分隔 | - |
| `--title` | `string` | 否 | 抖音标题 | `title` |
| `--description` | `string` | 否 | 图文描述。支持 HTML 格式和话题标签。**约束**: 必须用 `<p>` 标签包裹段落，最多 5 个话题。如未包裹，引擎会自动处理。 | `description` |
| `--image_urls` | `string` | 是 | 图片 URL 列表，多个请用逗号分隔。 | `images` (OldImage[]) |
| `--cover_url` | `string` | 否 | 封面图 URL。如不填，引擎默认取第一张图。 | `covers` (OldCover[]) |
| `--location` | `object` | 否 | 抖音地址/带货地址。JSON 字符串。来源于获取地址接口。 | `location` (PlatformDataItem) |
| `--musice` | `object` | 否 | **注意**: DTO 拼写为 `musice`。图文发布的音乐。JSON 字符串。 | `musice` (PlatformDataItem) |
| `--scheduledTime`| `number`| 否 | 定时发布时间戳（毫秒）。 | `scheduledTime` |
| `--collection` | `object` | 否 | 合集信息。JSON 字符串。要求抖音账号已创建合集。 | `collection` (Category) |
| `--sub_collection`| `object`| 否 | 合集选集。JSON 字符串。要求抖音账号已创建合集。 | `sub_collection` (Category) |

## 复杂对象结构 (Complex Object Structures)

### PlatformDataItem (用于 location, musice)
```json
{
  "id": "如果接口有返回 ID",
  "text": "显示的文本名称",
  "raw": { "根据不同接口返回的原始数据内容" }
}
```
> [!TIP]
> 建议直接通过查询接口（如 `get-locations.ts`）获取该对象，不要手动构造。

### Category (用于 collection, sub_collection)
```json
{
  "yixiaoerId": "内部 ID",
  "yixiaoerName": "名称",
  "yixiaoerImageUrl": "图标 URL",
  "yixiaoerDesc": "描述",
  "viewNum": "浏览量/统计",
  "raw": { "原始数据" }
}
```

## 依赖查询文档 (Dependencies)

如果需要获取上述复杂对象的数据，请参考以下查询指令：
- **获取地理位置**: [get-locations.md](../get-locations.md) -> `scripts/get-locations.ts`
- **获取音乐列表**: [get-music.md](../get-music.md) -> `scripts/get-music.ts`
- **获取合集列表**: [get-collections.md](../get-collections.md) -> `scripts/get-collections.ts`

## 引擎逻辑说明 (Engine Logic)

1. **描述处理**: 引擎会自动检测 `description`。如果是抖音平台且没有 HTML 标签，会按照换行符自动将其拆分为多个 `<p>` 标签包裹的内容。
2. **话题处理**: 描述中支持 `<topic text='话题名' raw='{...}'>#话题名</topic>` 格式。
3. **封面补全**: 如果未提供 `--cover_url`，引擎会自动将第一张图片作为封面填入 `covers` 数组。
4. **资源映射**: `--image_urls` 会被自动转换为后端要求的 `images: OldImage[]` 数组结构。

## 调用示例 (Example)

```bash
node scripts/publish.ts \
  --type=image-text \
  --platforms=抖音 \
  --account_ids=67c824558e8b233a00000000 \
  --title="生活随笔" \
  --description="今天阳光真好 #生活\n记录一下美好瞬间" \
  --image_urls="https://example.com/1.jpg,https://example.com/2.jpg"
```
