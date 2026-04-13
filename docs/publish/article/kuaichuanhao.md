# 📄 快传号 文章 参数 (KuaiChuanHao Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“快传号 (360 旗下)”发布新闻、生活、娱乐类文章时触发：
- **资讯分发**：将内容同步到 360 浏览器等分发渠道。
- **原创保护**：申明原创以获取更高的流量推荐和版权保护。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装快传号 Payload 时需遵守：
1. **原创申明优先**：快传号对原创内容推荐力度较大。若用户提到“原创”，务必将 `type` 设为 `1`。
2. **标签精准化**：`tags` 应包含文章的核心关键词，有助于在 360 搜索中获得更好的精准曝光。
3. **封面规范**：必须提供 `covers` 表单项，且封面需通过 `upload` 动作获取 OSS Key。
4. **定时发布自检**：若启用 `scheduledTime`，请确保时间点晚于当前发布节点。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题。 | - |
| **`content`** | `string` | **是** | 文章内容 (HTML 格式)。 | - |
| **`covers`** | `Array` | **是** | 封面列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`type`** | `number` | **是** | **创作类型**: `0`-不申明, `1`-申明原创。 | `1` |
| **`pubType`** | `number` | **是** | **存储模式**: `0`-草稿, `1`-发布。 | `1` |
| `tags` | `string[]` | 否 | 文章标签。 | - |
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
# 发布快传号原创资讯
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["快传号"],
  "publishArgs": {
    "content": "<h1>实战：跨平台内容分发避坑指南</h1><p>正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_kch_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "实战：跨平台内容分发避坑指南",
          "content": "<h1>实战：跨平台内容分发避坑指南</h1><p>正文内容...</p>",
          "covers": [{ "key": "kch_cov_1", "size": 120000, "width": 800, "height": 600 }],
          "type": 1,
          "pubType": 1,
          "tags": ["运营", "自媒体"]
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
| **创作类型不支持** | 传入了不支持的 `type` 数值。 | 仅限 0 (普通) 或 1 (原创)。 |
| **标题包含敏感词** | 360 审核系统拦截标题关键词。 | 优化标题文案。 |
| **封面上传失败** | `key` 无效或图片尺寸不合规。 | 检查图片是否成功上传 OSS，推荐 16:9 或 3:2 比例。 |
| **草稿状态异常** | `pubType` 设为 0 但在后台找不到草稿。 | 确认 `platformAccountId` 是否正确且在线。 |

---
> [!TIP]
> **流量入口**: 快传号是非原创内容高发区，因此坚持“申明原创”并保持高频更新是提升账号权重的核心手段。
