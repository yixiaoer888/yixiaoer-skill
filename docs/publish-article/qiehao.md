# 企鹅号文章发布 (Publish QiEHao Article)

该指令用于通过文章引擎向腾讯企鹅号分发长内容，支持企鹅号要求的内容声明、标签及封面设置。

## DTO 溯源 (Knowledge from QiEHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/qiehao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 (支持标准 HTML 标签) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 直连地址，引擎自动上传并映射为 `covers` 数组。建议比例 3:2。 |
| `--tags` | json | 是 | **文章标签** | 字符串数组。至少需要一个标签。示例：`["科技", "AI"]` |
| `--declaration` | number | 否 | **创作申明** | `0`: 暂不申明, `1`: AI生成, `2`: 个人观点, `3`: 剧情演绎, `7`: AI辅助创作, `8`: 健康医疗, `9`: 危险行为。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (单位：秒)。注意：企鹅号 DTO 中映射为 `prePubTime`。 |

## 调用指令示例 (Usage)

### 1. 立即发布一篇科技类文章（设置多个标签）
```bash
node scripts/publish.ts \
  --type=article \
  --title="企鹅号自动发布测试" \
  --content="<p>这是通过龙虾技能自动发布的企鹅号正文内容...</p>" \
  --platforms="企鹅号" \
  --account_ids="qeh_acc_001" \
  --cover_url="https://assets.example.com/qeh-cover.jpg" \
  --tags="[\"自媒体\", \"技术干货\"]" \
  --declaration=1
```

### 2. 存为草稿
```bash
node scripts/publish.ts \
  --type=article \
  --title="我的草稿文章" \
  --content="<p>草稿内容...</p>" \
  --platforms="企鹅号" \
  --account_ids="qeh_acc_001" \
  --cover_url="https://assets.example.com/cover.jpg" \
  --tags="[\"测试\"]" \
  --pubType=0
```

## 逻辑与规范说明
- **标签要求 (Tags)**: 企鹅号要求文章必须带有标签，且 `--tags` 必须是合法的 JSON 字符串数组格式。
- **引擎适配**: `publish.ts` 已针对 `type=article` 实现了 `covers` 的自动封装。建议封面尺寸符合企鹅号长文规范。
- **内容声明 (Declaration)**: 企鹅号提供了细粒度的内容声明选项（如 AI 辅助、医疗健康等），请根据内容实际属性准确选择。
- **资源预处理**: 若封面图为外链，引擎会自动调用 `uploadResource` 获取内部 Key 并完成 `OldCover` 模型的填充。
