# 📄 搜狐号 文章 参数 (SouHuHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“搜狐号”平台发布社会资讯、娱乐八卦、技术教程或企业通稿时触发：
- **搜索引擎优化**：通过搜狐号在搜索引擎中的高权重获取外部流量。
- **摘要发布**：为文章配置精准摘要，提升首页信息流点击率。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装搜狐号 Payload 时需遵守：
1. **摘要强制校验**：搜狐号强制要求摘要 `desc`。若用户未提供，Agent 应从正文首段自动提取 50-100 字作为补充。
2. **封面单/多图逻辑**：搜狐号封面数组 `covers` 必须通过 `upload` 动作获取 Key，且需包含完整的 `OldCover` 元数据。
3. **内容发布时效**：支持 `scheduledTime` 定时发布，Agent 应校验时间戳准确性。
4. **意图确认**：向用户明确是否存为草稿 (pubType: 0)。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式)。 | - |
| **`desc`** | `string` | **是** | 文章摘要或描述。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
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
# 发布搜狐号带摘要的文章
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["搜狐号"],
  "publishArgs": {
    "content": "<h1>搜狐号：自媒体的常青树</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_sh_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "搜狐号：自媒体的常青树及其流量价值",
          "content": "<h1>常青树解析</h1><p>正文内容...</p>",
          "desc": "搜狐号作为老牌自媒体平台，在搜索引擎权重和信息流分发中依然具有核心价值。",
          "covers": [{ "key": "sh_cov_1", "size": 102400, "width": 800, "height": 600 }],
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
| **摘要字数不足** | `desc` 字符少于 10 字。 | 请补全摘要内容，使之更具概括性。 |
| **标题重复拦截** | 相同标题在同一账号下重复投递。 | 修改标题后缀或更正发布状态。 |
| **封面像素不足** | `width/height` 小于平台要求。 | 搜狐号建议封面宽度不小于 640px。 |
| **定时设置无效** | `scheduledTime` 为过去的时间点。 | 重新校对 Unix 时间戳。 |

---
> [!TIP]
> **SEO 强化**: 搜狐号在百度等搜索引擎。Agent 建议标题多包含热点搜索关键词，并在摘要中精准概括，以获得更多搜索点击。
