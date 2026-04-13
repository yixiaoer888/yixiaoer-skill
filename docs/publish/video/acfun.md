# 📄 AcFun 视频 参数 (AcFun Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“AcFun (A站)”分发视频，且涉及以下特有需求时触发：
- **A 站投稿**：发布带有 A 站特色的视频内容。
- **分类挂载**：设置视频分区 (Category)，这是 A 站投稿的必填项。
- **原创申明**：申明视频为原创内容。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 AcFun 视频 Payload 时需遵守：
1. **分类强制原则**：AcFun 强制要求视频分类。必须通过 `categories` 接口获取并在 `category` 数组中透传 `raw` 数据。
2. **描述格式化**：建议 `description` 使用 HTML 格式，至少包含 `<p>` 标签。
3. **标签限制**：标签最多支持 6 个。
4. **资源引用**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`category`** | `Array` | **是** | 视频分类。使用 `CascadingPlatformDataItem[]` 结构。 | - |
| **`type`** | `number` | **是** | **内容类型**: `1`-原创, `0`-非原创。 | `0` |
| `description` | `string` | 否 | 视频描述 (建议使用 `<p>` 包裹)。 | - |
| `tags` | `string[]` | 否 | 视频标签 (最多 6 个)。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布 AcFun 原创投递
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["AcFun"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "AC_ACC_001",
        "video": { "key": "ac_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "A 站发布测试",
          "description": "<p>这是 AcFun 视频的描述内容</p>",
          "tags": ["生活", "美食"],
          "category": [
            { "yixiaoerId": "cat_001", "yixiaoerName": "生活", "raw": {...} }
          ],
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
| **分类未选择** | `category` 字段缺失或数组为空。 | A 站投稿必须选择分类，请调用 `categories` 接口。 |
| **标签过多** | `tags` 数组元素超过 6 个。 | 移除多余标签，保留核心关键词。 |
| **原创申明冲突** | 勾选了原创但内容被识别为盗版。 | 确保内容确实为原创或改为非原创。 |
| **定时发布失败** | `scheduledTime` 设置不合理。 | 检查时间戳是否正确且符合平台预设规则。 |

---
> [!TIP]
> **二次元氛围**: AcFun 是高粘性的二次元社区，建议视频标题和描述可以更加活泼，增加互动引导。
