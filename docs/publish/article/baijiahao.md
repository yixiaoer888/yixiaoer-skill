# 📄 百家号文章发布参数 (BaiJiaHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 1. 触发场景 (Trigger)

当用户明确要在“百家号”发布文章（长文），且需要配置如“选择分类”、“设置封面”或“参与征文活动”等功能时触发。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装百家号文章 Payload 时需遵守：
1. **标题/正文限制**：标题必须在 2-30 字符内，正文 HTML 必须在 9-10000 字符内。
2. **分类必选要求**：百家号必须提供 `category` 数据（1-2 个）。必须先调用 `categories` 接口获取。
3. **活动关联引导**：若用户提到“参加活动”，必须调用 `activities` 获取最新活动列表并填入 `activity` 字段。
4. **资源合规**：文章封面 `covers` 必须包含 1-3 张图片。

## 3. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (2-30 字符)。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式, 9-10000 字符)。 | - |
| **`covers`** | `Array` | **是** | 文章封面列表 (`OldCover[]`, 1-3 张)。 | - |
| **`category`** | `Array` | **是** | 文章分类列表 (`Category[]`, 1-2 个)。 | - |
| **`pubType`** | `number` | **是** | 发布类型: `0`-草稿, `1`-直接发布。 | - |
| `declaration` | `number` | 否 | 内容声明: `0`-不声明, `1`-内容由 AI 生成。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |
| `activity` | `Object` | 否 | **征文活动**: 使用 `Activity` 结构。 | - |

### 3.1 复杂对象结构 (Data Schemas)

- **OldCover**: 包含 `key`, `size`, `width`, `height`。
- **Category (分类)**: 必须包含 `yixiaoerId`, `yixiaoerName` 和 **完整的 `raw` 对象**。
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
> **百度分发流量提示**：百家号对文章分类的准确度非常敏感。Agent 建议用户优先选择二级或三级细分分类以获得更精准的百度搜索流量分配。
