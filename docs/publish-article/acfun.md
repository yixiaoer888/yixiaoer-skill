# AcFun文章发布 (Publish AcFun Article)

该指令用于通过文章引擎向 AcFun（A站）分发长文内容。

## DTO 溯源 (Knowledge from AcFunArticleForm)
*逻辑来源: `apps/server-api/packages/yxr-open-platform/src/models/platform/acfun.dto.ts`*

### 核心参数 (Command Arguments)

| 参数名 | 类型 | 必填 | 说明 | DTO 校验与逻辑 |
| :--- | :--- | :--- | :--- | :--- |
| `--type` | string | 是 | **必须固定为 `article`** | 业务模态识别 |
| `--title` | string | 是 | 文章标题 | 不可为空 |
| `--content` | string | 是 | 文章内容 (HTML) | 不可为空 |
| `--desc` | string | 否 | **文章摘要/描述** | 简短描述 |
| `--cover_url` | string | 是 | 封面图 | 引擎自动上传并映射为 `covers` 数组 |
| `--tags` | string[] | 是 | 标签 | 字符串数组 |
| `--type` (extra)| number | 是 | **创作类型** | `0`:不申明 `1`:申明原创 |
| `--category` | array | 是 | 文章分类 | 对象列表 |
| `--contentSourceUrl`| string | 否 | 原文链接 | 转载时必须填写 |

## 调用指令示例 (Usage)

```bash
node scripts/publish.ts \
  --type=article \
  --title="A站长文发布测试" \
  --content="<p>文章内容展示...</p>" \
  --desc="测试文章摘要" \
  --platforms="AcFun" \
  --account_ids="ac_acc_001" \
  --cover_url="https://example.com/cover.jpg" \
  --tags='["二次元", "技术"]' \
  --category='[{"id": "1", "text": "文章"}]' \
  --type=1
```
