# 百家号文章发布 (Publish Baijiahao Article)

该指令用于通过文章引擎向百家号分发长内容，支持百家号要求的分类选择与内容声明。

## DTO 溯源 (Knowledge from BaiJiaHaoArticleForm)
*逻辑来源: `apps/server-api/.../baijiahao.dto.ts`*

### 核心参数 (Command Arguments)
| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | HTML 内容 | 不可为空 |
| `--cover_url` | string | 是 | 封面图 | DTO 映射为 `covers` 数组 |
| `--category_id` | string | 是 | 百家号分类 ID | 必须提供有效的分类 ID |
| `--category_name` | string | 否 | 百家号分类名称 | 默认为“默认” |
| `--declaration` | number | 是 | 内容声明 | 0: 不声明, 1: 内容由 AI 生成 |
| `--pub_type` | number | 否 | 发布类型 | 0: 草稿, 1: 立即发布 (默认 1) |

## 调用指令示例 (Usage)

### 1. 发布一篇科技类原创文章
```bash
node scripts/publish-article.ts \
  --title="2024 量子计算突破" \
  --content="<p>正文内容...</p>" \
  --platforms="百家号" \
  --account_ids="bjh_123" \
  --cover_url="https://example.com/cover.jpg" \
  --category_id="1" \
  --category_name="科技" \
  --declaration=0
```

### 2. 存为百家号草稿
```bash
node scripts/publish-article.ts \
  --title="待发布的深度好文" \
  --content="<p>初稿内容</p>" \
  --platforms="百家号" \
  --account_ids="bjh_123" \
  --category_id="2" \
  --pub_type=0 \
  --declaration=1
```

## 逻辑说明
- **分类强制性**: 百家号 DTO 要求 `category` 数组不能为空。脚本会自动将 `--category_id` 包装为符合要求的对象。
- **封面处理**: 默认为 `single` 类型封面，由 `baseArticleForm` 统一管理。
- **声明规则**: 必须显示指定 `--declaration`，否则根据 DTO 校验可能会失败。
