# 📄 图文发布通用索引 (Image-Text Publish Index)

本文档定义了蚁小二所有平台图文（短动态/笔记）发布的 **根 Payload 结构**。在查阅具体平台（如小红书、新浪微博）文档前，**必须**先理解并遵循本规范。

## 1. 触发场景 (Trigger)

当用户下达发布短动态、笔记或朋友圈风格内容的指令时触发：
- **生活分享**：“把这三张风景照发到小红书”。
- **动态同步**：“发一条带图片的微博告知大家活动开始”。
- **多图矩阵**：“将产品宣传图同步到所有自媒体账号的动态/笔记中”。

## 2. 交互协议 (Interactive Protocol)

Agent 在构造发布指令前需遵守：
1. **多图上传原则**：图文发布通常包含多张图片。必须遍历所有图片并调用 `upload` 动作获取对应的 `key`，严禁使用原始 URL。
2. **账号校验原则**：必须先通过 `accounts` 确认 `platformAccountId` 处于 `status: 1`。
3. **内容一致性原则**：根结构的 `publishArgs.content` 通常作为各账号的默认描述，除非在 `contentPublishForm` 中有特殊覆盖。
4. **意图确认**：向用户展示 **[平台] - [图片张数] - [描述片段]** 并确认后执行。

## 3. 根结构参数定义 (Base Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `publish` | 固定值。 |
| **`publishType`** | `string` | **是** | `imageText`| 发布类型。 |
| **`platforms`** | `string[]` | **是** | - | 目标平台标识数组。 |
| `publishArgs` | `object` | **是** | - | 详见下方 [3.1 publishArgs](#31-publishargs-定义)。 |
| `publishChannel` | `string` | 否 | `cloud` | 发布通道。某些移动端平台建议使用 `cloud` 以提高成功率。 |
| `isDraft` | `boolean` | 否 | `false` | 若为 `true`，存入蚁小二内部草稿。 |

### 3.1 publishArgs 定义

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`content`** | `string` | **是** | **核心描述**: 支持换行和表情。 |
| **`accountForms`** | `Array` | **是** | 账号发布表单列表。 |

### 3.2 accountForms 元素定义 (AccountFormItem)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`platformAccountId`** | `string` | **是** | 账号 ID。 |
| **`images`** | `Array` | **是** | **ImageFormItem[]**: 图片对象数组，包含 `key`, `width`, `height`, `size`。 |
| **`cover`** | `object` | **是** | **ImageFormItem**: 主封面（通常是第一张图）。 |
| **`coverKey`** | `string` | **是** | 必须与 `cover.key` 一致。 |
| **`contentPublishForm`**| `object` | **是** | **平台差异层**: 内部字段（如话题、地点）见各平台文档。 |

## 4. 执行指令示例 (Command)

```bash
# 发布单图动态到微博
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["XinLangWeiBo"],
  "coverKey": "img_001",
  "publishArgs": {
    "content": "今日天气不错！",
    "accountForms": [{
      "platformAccountId": "WB_123",
      "images": [{"key": "img_001", "width": 800, "height": 800, "size": 150000}],
      "cover": {"key": "img_001", "width": 800, "height": 800, "size": 150000},
      "coverKey": "img_001",
      "contentPublishForm": { "formType": "task" }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **图片比例不合规** | 部分平台（如小红书）对封面比例有特定要求。 | 建议使用 3:4 或 1:1 比例，并在 Agent 侧根据平台标准预警。 |
| **描述文字超限** | 字数超过了平台限制（微博 140/2000 字等）。 | 检查具体平台文档的长度限制。 |
| **图片 Key 缺失** | `images` 数组中有元素未包含 `key`。 | 确保所有图片都经过了 `upload` 动作处理。 |
| **发布中途失败** | 大多为 `cloud` 发布时的代理问题。 | 请参考《避坑与故障排查手册》中的代理配置部分。 |

---
> [!IMPORTANT]
> **图文首选通道**：对于注重垂直领域的图文平台（如小红书），强烈建议在 `contentPublishForm` 中增加标签 (tags) 和地点 (location) 以获得更多曝光。
