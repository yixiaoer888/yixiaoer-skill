# 📄 百家号文章发布参数 (BaiJiaHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 1. 触发场景 (Trigger)

当用户明确要在“百家号”发布文章（长文），且需要配置如“选择分类”、“设置封面”或“参与征文活动”等功能时触发。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装百家号文章 Payload 时需遵守：
1. **标题/正文限制**：标题必须在 2-30 字符内，正文 HTML 必须在 9-10000 字符内。
2. **分类深度要求**：百家号对分类准确度极高。Agent **必须** 引导用户选择至最细层级（若存在二级分类，则必须选中二级分类）。通常提供 1-2 个分类，必须先调用 `categories` 接口获取。
3. **活动关联引导**：若用户提到“参加活动”，必须调用 `activities` 获取最新活动列表并填入 `activity` 字段。
4. **资源合规**：文章封面 `covers` 必须包含 1-3 张图片。

## 3. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (2-30 字符)。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式, 9-10000 字符)。 | - |
| **`covers`** | `Array` | **是** | 文章封面列表 (`OldCover[]`, 1-3 张)。 | - |
| **`category`** | `Array` | **是** | 文章分类列表 (`Category[]`, 1-2 个)。**必须选中二级分类（若存在）**。 | - |
| **`pubType`** | `number` | **是** | 发布类型: `0`-草稿, `1`-直接发布。 | - |
| `declaration` | `number` | 否 | 内容声明: `0`-不声明, `1`-内容由 AI 生成。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |
| `activity` | `Object` | 否 | **征文活动**: 使用 `Activity` 结构。 | - |

### 3.1 复杂对象结构 (Data Schemas)

- **OldCover**: 包含 `key`, `size`, `width`, `height`。
- **Category (分类)**: 必须包含 `yixiaoerId`, `yixiaoerName` 和 **完整的 `raw` 对象**。
  - **重要原则**：如果 `categories` 接口返回的数据中有子分类（Children），Agent 必须引导或自动匹配到对应的二级分类，不得仅选择一级大类。
- **Activity (活动)**: 包含 `yixiaoerId`, `yixiaoerName`。

## 4. 执行指令示例 (Command)

```bash
# 百家号文章发布：带分类和 AI 声明
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["百家号"],
  "publishArgs": {
    "content": "<h1>百家号发布测试</h1><p>正文内容至少需要9个字以上。</p>",
    "accountForms": [{
      "platformAccountId": "BJH_001",
      "contentPublishForm": {
        "formType": "task",
        "title": "今日科技头条",
        "content": "<h1>百家号发布测试</h1><p>正文内容至少需要9个字以上。</p>",
        "covers": [{"key": "c_001", "size": 102400, "width": 800, "height": 600}],
        "category": [{ "yixiaoerId": "cat_001", "yixiaoerName": "文化", "raw": {...} }],
        "pubType": 1,
        "declaration": 1
      }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **标题长度不符** | 标题少于 2 字或多于 30 字。 | 实时检查标题字数并提示用户。 |
| **正文内容过短** | HTML 正文物理字符少于 9 个字。 | 提醒用户增加正文丰富度。 |
| **分类信息错误** | `category.raw` 数据格式不完整。 | 重新执行 `categories` 查询流程获取。 |
| **封面数量不匹配** | `covers` 数组长度不在 1-3 之间。 | 引导用户至少选择 1 张封面图。 |

---
> [!IMPORTANT]
> **分类层级强约束**：百家号对文章分类的准确度及层级非常敏感。
> 1. **严禁漏选**：发布任务前必须确保分类已填入。
> 2. **必须到二级**：若目标分类存在二级细分，**必须** 选择二级分类作为最终值（例如：不能只选“科技”，必须选“科技 > 智能家居”）。这直接影响百度搜索的流量分配。
