# 百家号文章发布 (Publish Baijiahao Article)

该指令用于通过文章引擎向百家号分发长内容，支持百家号要求的分类选择、内容声明、征文活动参与及定时发布。

## DTO 溯源 (Knowledge from BaiJiaHaoArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 (支持标准 HTML 标签) | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | 直连地址，引擎自动上传并映射为 `covers` 数组。建议比例 3:2。 |
| `--category` | json | 是 | **文章分类** | 必须传入两级分类的路径数组。格式为 `[一级分类, 二级分类]`。务必遵循 [查询发布分类](../get-publish-categories.md) 结果中的父子关系。 |
| `--declaration` | number | 是 | **内容声明** | `0`: 不声明, `1`: 内容由 AI 生成。 |
| `--scheduledTime`| number | 否 | 定时发布时间 | Unix 时间戳 (秒)。 |
| `--activity` | json | 否 | **征文活动** | 对象格式。参与特定活动以获得流量扶持。通过 [查询征文活动](../get-publish-activities.md) 获取。 |
| `--pubType` | number | 否 | 发布类型 | `0`: 保存草稿, `1`: 立即提交 (默认 1)。 |

## 嵌套模型定义 (Nested Models)

### 文章分类 (Category Object)
百家号要求分类以两级数组形式传入：第一项为一级分类，第二项为二级分类。
```json
[
  {
    "yixiaoerId": "1",
    "yixiaoerName": "科技",
    "raw": {} 
  },
  {
    "yixiaoerId": "101",
    "yixiaoerName": "人工智能",
    "raw": {}
  }
]
```

### 征文活动 (Activity Object)
```json
{
  "id": "activity_123",
  "name": "2024 AI 创作者计划",
  "raw": {} // 平台原始活动数据
}
```

## 调用指令示例 (Usage)

### 1. 立即发布一篇科技类文章（声明 AI 生成）
```bash
node scripts/publish.ts \
  --type=article \
  --title="大模型时代的端侧 AI 突破" \
  --content="<p>正文 HTML 内容...</p>" \
  --platforms="百家号" \
  --account_ids="bjh_acc_001" \
  --cover_url="https://assets.example.com/cover.jpg" \
  --category="[{\"yixiaoerId\": \"1\", \"yixiaoerName\": \"科技\"}, {\"yixiaoerId\": \"101\", \"yixiaoerName\": \"人工智能\"}]" \
  --declaration=1
```

### 2. 存为草稿并参与征文活动
```bash
node scripts/publish.ts \
  --type=article \
  --title="我的征文投稿" \
  --content="<p>投稿内容...</p>" \
  --platforms="百家号" \
  --account_ids="bjh_acc_001" \
  --cover_url="https://assets.example.com/cover.jpg" \
  --category="[{\"yixiaoerId\": \"1\", \"yixiaoerName\": \"科技\"}, {\"yixiaoerId\": \"101\", \"yixiaoerName\": \"人工智能\"}]" \
  --declaration=0 \
  --activity="{\"id\": \"act_789\", \"name\": \"春季征文\"}" \
  --pubType=0
```

## 逻辑与规范说明
- **分类选择 (Category Selection)**: 百家号分类具有严格的两级层级关系。在调用 `get-publish-categories.ts` 时，必须从返回的树形结构中选择一个具体的二级分类（存在 `child` 数组中的项），并构造包含其完整路径（[一级, 二级]）的数组传给 `--category`。
- **引擎适配**: `publish.ts` 已针对 `type=article` 实现了 `covers` 的自动封装（将 `--cover_url` 转换为 `[{key, width, height}]`）以及 `title/content` 的补全。
- **数据驱动**: AI 在生成发布参数时必须严格遵守 DTO 字段名（如 `scheduledTime`），引擎会透传所有自定义参数。
- **资源预处理**: 若封面图为外链，引擎会自动调用 `uploadResource` 获取内部 Key。
