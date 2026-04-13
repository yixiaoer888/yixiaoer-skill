# 📄 爱奇艺文章发布参数 (AiQiYi Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“爱奇艺 (iQIYI)”平台分发文章或图文内容时触发。支持：
- **存为草稿**：先同步到爱奇艺后台而不发布。
- **直接发布**：同步并立即公开。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装爱奇艺 Payload 时需遵守：
1. **基础参数锁定**：爱奇艺的文章发布主要透传 `title` 和 `pubType`，正文 `content` 承载于根 `publishArgs` 中。
2. **意图辨析**：明确用户是希望“存为爱奇艺草稿 (pubType: 0)”还是“立即发布 (pubType: 1)”。
3. **资源引用**：虽主要由根结构控制正文，但确保正文中的图片已转换为 key 标识。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 50 字符)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |

---

## 4. 执行指令示例 (Command)

```bash
# 发布爱奇艺文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["爱奇艺"],
  "publishArgs": {
    "content": "<h1>爱奇艺发布演示</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "AQY_ACC_123",
        "contentPublishForm": {
          "formType": "task",
          "title": "爱奇艺媒体号图文分发示例",
          "pubType": 1
        }
      }
    ]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **标题过长** | 标题超过了 50 字符。 | 精简标题。 |
| **账号状态异常** | 爱奇艺账号由于未实名或被封禁导致分发失败。 | 检查账号在蚁小二中的状态。 |
| **内容同步卡顿** | 爱奇艺接口连接超时。 | 建议稍后重试或检查网络。 |

---
> [!TIP]
> **多终端展示**: 在爱奇艺发布的文章通常会在其移动端“号”频道及搜索结果中展示。
