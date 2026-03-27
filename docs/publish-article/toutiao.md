# 头条号文章发布 (Publish Toutiao Article)

该指令用于通过文章引擎向头条号分发长内容，支持头条首发选项、广告收益设置、内容声明、位置设置及定时发布。

## DTO 溯源 (Knowledge from TouTiaoHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 (支持标准 HTML 标签) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 直连地址，引擎自动上传并映射为 `covers` 数组。 |
| `--isFirst` | boolean | 否 | **头条首发** | 是否在头条首发。`true`/`false`。 |
| `--advertisement` | number | 否 | **广告投放** | `2`: 投放广告赚取收益, `3`: 不投放 (默认 `3`)。 |
| `--declaration` | number | 否 | **创作类型声明** | `1`: 自行拍摄, `2`: 取自站外, `3`: AI生成, `6`: 虚构演绎, `7`: 投资观点, `8`: 健康医疗分享。 |
| `--location` | json | 否 | **地理位置** | 对象格式。通过 [查询地理位置](#) 相关接口获取。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

## 嵌套模型定义 (Nested Models)

### 地理位置 (Location Object)
```json
{
  "id": "loc_123",
  "text": "北京市朝阳区",
  "raw": {} 
}
```

## 调用指令示例 (Usage)

### 1. 立即发布一篇科技文章（声明 AI 生成，开启广告收益）
```bash
node scripts/publish.ts \
  --type=article \
  --title="2024 智能硬件趋势报告" \
  --content="<p>正文 HTML 内容...</p>" \
  --platforms="头条号" \
  --account_ids="tt_acc_001" \
  --cover_url="https://assets.example.com/tt-cover.jpg" \
  --isFirst=true \
  --advertisement=2 \
  --declaration=3
```

### 2. 存为草稿
```bash
node scripts/publish.ts \
  --type=article \
  --title="未完成的草稿" \
  --content="<p>内容...</p>" \
  --platforms="头条号" \
  --account_ids="tt_acc_001" \
  --cover_url="https://assets.example.com/tt-cover.jpg" \
  --pubType=0
```

## 逻辑与规范说明
- **引擎适配**: `publish.ts` 已针对 `type=article` 实现了 `covers` 的自动封装（将 `--cover_url` 转换为 `[{key, width, height}]`）以及 `title/content` 的补全。
- **数据驱动**: AI 在生成发布参数时必须严格遵守 DTO 字段名（如 `isFirst`, `advertisement`），引擎会透传所有自定义参数。
- **布尔值处理**: 脚本支持 `--isFirst=true` 这种格式并自动解析为布尔类型。
- **资源预处理**: 若封面图为外链，引擎会自动调用 `uploadResource` 获取内部 Key。
