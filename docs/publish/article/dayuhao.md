# 📄 大鱼号文章发布参数 (DaYuHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“大鱼号 (阿里大鱼)”平台发布文章、资讯或图文内容时触发。支持：
- **封面适配**：支持设置横版封面及竖版封面。
- **创作者申明**：申明内容为虚构演绎或 AI 生成。
- **草稿同步**：同步到大鱼号创作后台草稿箱。
- **定时发布**：预约时间进行内容分发。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装大鱼号 Payload 时需遵守：
1. **标题字数校验**：标题建议控制在 50 字符以内。
2. **封面多样性**：大鱼号支持 `covers` (横版) 和 `verticalCovers` (竖版)。Agent 应检查资源比例并填入对应字段。
3. **AI 申明意识**：若识别到内容由 AI 生成，必须主动设置 `declaration: 4`。
4. **资源依赖**：封面对象必须包含完整的 `OldCover` 结构，且 key 必须通过 `upload` 动作产生。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 50 字符)。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式，最多 50000 字符)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `covers` | `Array` | 否 | 文章横版封面列表。见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| `verticalCovers` | `Array` | 否 | 文章竖版封面列表。见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| `declaration` | `number` | 否 | **声明**: `0`-无, `3`-虚构演绎, `4`-AI生成。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布大鱼号文章：附带封面与 AI 申明
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["大鱼号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DY_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "阿里大鱼平台内容分发策略",
          "content": "<h1>策略报告</h1><p>正文内容...</p>",
          "covers": [
            { "key": "c_key_h", "size": 150000, "width": 800, "height": 450 }
          ],
          "declaration": 4,
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
| **封面比例不符** | 使用了不适合大鱼号的非标准比例图片。 | 推荐横版 16:9，竖版 3:4。 |
| **账号鉴权失败** | 大鱼号 Token 过期或账号被限制。 | 登录蚁小二客户端重新授权账号。 |
| **内容包含外链** | 正文 HTML 中包含非法外部超链接。 | 移除所有非官方认可的外部链接。 |
| **标题重复** | 该账号下已发布过标题完全一致的文章。 | 微调标题内容再尝试发布。 |

---
> [!TIP]
> **全端推流**: 大鱼号发布后会自动分发至 UC 浏览器、优酷、土豆等阿里系多个终端。
