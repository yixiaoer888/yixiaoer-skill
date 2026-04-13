# 📄 哔哩哔哩 文章 参数 (BiLiBiLi Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“哔哩哔哩 (B站)”发布专栏文章时触发。常见需求包括：
- **专栏投稿**：发布长图文教程、心得或动态。
- **原创保护**：申明原创以获取更多权益。
- **社交标签**：为文章添加最多 12 个技术或生活标签。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 B 站文章 Payload 时需遵守：
1. **多封面支持**：B 站文章支持 1-20 张封面。Agent 应确保 `covers` 数组中的每个对象都包含完整的 `OldCover` 结构。
2. **分类必选与深度要求**：即使是文章，B 站也要求提供准确分类。必须通过 `categories` 接口获取。若存在子分类（分区），**必须** 选中并透传最深层级的 `raw` 对象。严禁停留在父级大类。
3. **正文 HTML 转换**：确保 `content` 已经转换为适合专栏展示的 HTML 格式。
4. **资源预处理**：封面的 `key` 必须通过 `upload` 动作产生。

## 3. 参数 definition (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 80 字符)。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式，最多 100,000 字符)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。支持 1-20 张。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `tags` | `string[]` | 否 | 文章标签 (最多 12 个)。 | - |
| `createType` | `number` | 否 | **原创类型**: `0`-非原创, `1`-原创。 | `1` |
| `category` | `Array` | 否 | 文章分类。使用 `Category[]` 结构。**必须下钻至二级分区**。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

### 3.3 Category 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 分类 ID。 |
| `yixiaoerName` | `string` | **是** | 分类名称。 |
| `raw` | `Object` | **是** | 原始数据对象。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布 B 站原创专栏文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["哔哩哔哩"],
  "publishArgs": {
    "content": "<h1>B 站专栏实战</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_bili_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "深度：B 站自动化分发技术解析",
          "content": "<h1>深度解析</h1><p>内容正文...</p>",
          "covers": [
            { "key": "bili_cover_1", "size": 102400, "width": 800, "height": 600 }
          ],
          "createType": 1,
          "pubType": 1,
          "tags": ["科技", "自动化", "程序人生"]
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
| **标题敏感词拦截** | 标题包含违规或夸大词汇。 | 修改标题，参考 B 站创作守则。 |
| **封面上传失败** | 封面超出 20 张限制或单张超大。 | 检查封面数量和体积。 |
| **HTML 标签解析错误** | 正文中包含 B 站专栏不支持的 HTML 标签。 | 建议仅使用基本的 `h1`, `p`, `img` 等标签。 |
| **标签过多** | `tags` 数组长度超过 12。 | 精简标签至 12 个以内。 |

---
> [!TIP]
> **多图文建议**: B 站专栏非常欢迎图片丰富的技术内容，建议在正文中合理插入与封面一致的高清图片以提升阅读体验。
