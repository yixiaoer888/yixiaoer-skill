# 📄 一点号 文章 参数 (YiDianHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“一点号 (一点资讯)”发布新闻、评论、生活指南或娱乐报道时触发：
- **资讯分发**：利用一点资讯的个性化引擎进行内容分发。
- **合规申明**：标注创作来源，如取材网络、AI 生成或虚构情节。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装一点号 Payload 时需遵守：
1. **标题长度限制**：一点号标题要求在 5-64 字符之间。Agent 应检查标题是否在此范围内。
2. **创作声明规范**：必须根据内容来源设置 `declaration`。例如，若检测到 AI 痕迹，应主动设置为 `4` (AI 生成)。
3. **资源上传原则**：封面 `covers` 中的图片 `key` 必须由系统内部 `upload` 动作产生。
4. **意图确认**：向用户明确是否存为“一点号草稿 (pubType: 0)”。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (长度: 5-64 字符)。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `declaration` | `number` | 否 | **创作声明**: 3-取材网络, 4-AI 生成, 5-虚构情节。 | - |
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
# 发布一点号资讯文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["一点号"],
  "publishArgs": {
    "content": "<h1>深度：新一代智能分发系统的实践</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_yd_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "深度：新一代智能分发系统的实践与探索",
          "content": "<h1>实践与探索</h1><p>正文内容...</p>",
          "covers": [{ "key": "yd_cov_1", "size": 102400, "width": 800, "height": 600 }],
          "declaration": 3,
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
| **标题长度超限** | 标题少于 5 字或超过 64 字。 | 调整标题长度以符合一點号规范。 |
| **声明不合法** | `declaration` 填入了不存在的枚举值。 | 仅限 [3.1](#31-核心表单参数-contentpublishform) 中定义的数值。 |
| **封面尺寸错误** | 封面的分辨率过低或比例严重失调。 | 建议使用 16:9 比例的高清原图。 |
| **发布中途失败** | 账号登录状态已过期或令牌失效。 | 检查该一点号账号在蚁小二中的在线状态。 |

---
> [!TIP]
> **一点资讯渠道**: 一点号与小米浏览器、OPPO 浏览器等深度打通，Agent 建议标题多体现“生活节奏”或“社会热点”，有助于提升点击转化。
