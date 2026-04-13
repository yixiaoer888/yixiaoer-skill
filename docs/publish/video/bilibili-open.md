# 📄 哔哩哔哩-Open 视频 参数 (Bilibili-Open Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户通过 B 站开放平台接口（通常为企业级或特殊接入端）分发视频时触发：
- **合规投稿**：发布带有详细创作者申明的内容。
- **级联分区**：设置复杂的 B 站视频分区。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 Bilibili-Open Payload 时需遵守：
1. **级联分区必填**：分区对内容收录至关重要。必须调用 `categories` 接口获取并在 `category` 数组中透传 `raw` 元数据。
2. **原创声明策略**：根据 `type` (1-自制, 2-转载) 准确映射内容。转载时强制需要 `contentSourceUrl`。
3. **内容安全声明**：必须明确 `declaration` 字段数值，以告知审核系统该视频的创作背景（如 AI 合成）。
4. **资源引用**：必须通过 `upload` 动作产生视频 key 及封面 key。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述。 | - |
| **`tags`** | `string[]` | **是** | 视频标签。建议 1-10 个。 | - |
| **`category`** | `Array` | **是** | 视频分类。使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`declaration`** | `number` | **是** | **创作者申明**: `0`-不申明, `1`-AI合成, `2`-危险行为, `3`-仅供娱乐, `4`-引人不适, `5`-个人观点。 | `0` |
| **`type`** | `number` | **是** | **发布类型**: `1`-自制, `2`-转载。 | `1` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |
| `contentSourceUrl` | `string` | 否 | 原文 URL 链接 (当 `type` 为 2 时必填)。 | - |
| `collection` | `Object` | 否 | 合集信息。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `id`, `text`, `raw`。支持 `children` 嵌套列表。

## 4. 执行指令示例 (Command)

```bash
# 通过开放平台发布 B 站自制视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["BilibiliOpen"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BILI_OPEN_ACC_01",
        "video": { "key": "v_key_open", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "B 站 Open 版发布测试",
          "description": "测试视频描述内容。",
          "tags": ["科技", "OpenAPI"],
          "category": [{ "id": "1", "text": "生活", "raw": {...} }],
          "declaration": 0,
          "type": 1
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
| **转载链接缺失** | `type` 为 2 但未提供 `contentSourceUrl`。 | 补全原文链接。 |
| **分区错误 (400)** | `category` 数据结构不完整或使用了非法 ID。 | 调用 `categories` 获取并完整透传 `raw` 字段。 |
| **申明信息缺失** | `declaration` 未传值。 | Open 平台对安全声明通常为强校验，请确保传参。 |
| **标签格式错误** | `tags` 传入了非字符串数组格式。 | 检查并格式化为标准 `string[]`。 |

---
> [!IMPORTANT]
> **OpenAPI 限制**: 此发布通道主要针对具有第三方 API 接入权限的账号，参数对合规性要求高于普通客户端。
