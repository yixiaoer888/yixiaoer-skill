# 📄 搜狐视频 视频 参数 (Souhu Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“搜狐视频 (客户端)”分发原创短视频、娱乐段子或 VLOG 内容时触发：
- **客户端流量同步**：触达搜狐视频独立 App 的高频用户群。
- **互动/合规发布**：申明视频背景（如 AI 数字人生成）。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装搜狐视频 Payload 时需遵守：
1. **自研声明规范**：搜狐视频对 `declaration` 有细化的枚举要求，特别是针对“AI 数字人”等前沿内容，Agent 应引导用户准确设置。
2. **标题为核心**：发布表单中 `title` 为关键项。
3. **标签与描述协同**：虽然非必填，但建议通过 `tags` 和 `description` 增加内容的被检索几率。
4. **资源引用协议**：必须通过 `upload` 动作产生视频 key 及封面 key。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| `description` | `string` | 否 | 视频描述内容。 | - |
| `tags` | `string[]` | 否 | 视频标签。 | - |
| `declaration` | `number` | 否 | **声明**: `0`-无需, `3`-AI 生成, `4`-虚构演绎, `5`-AI 数字人生成。 | `0` |

## 4. 执行指令示例 (Command)

```bash
# 发布搜狐视频 AI 数字人视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Souhushipin"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SOUHU_V_ACC_01",
        "video": { "key": "sh_v_api_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "AI 时代的虚拟偶像生存报告",
          "description": "内容由 AI 数字人自动生成测试。",
          "tags": ["AI", "黑科技", "虚拟"],
          "declaration": 5
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
| **标题缺失** | `title` 字段未提供。 | 搜狐视频强制要求提供视频标题。 |
| **申明不合规** | 传入了不在 0, 3, 4, 5 枚举范围内的数值。 | 修正声明传参。 |
| **视频画质异常** | 低清晰度视频导致分发受限。 | 建议使用 720P 及以上素材，首选 1080P。 |
| **账号令牌过期** | 账号在蚁小二中的授权已失效。 | 引导用户重新扫码登录该搜狐视频账号。 |

---
> [!TIP]
> **多元化分发**: 搜狐视频对脑洞类内容接受度较高。Agent 建议描述内容具有互动性，并配合高清封面以适应 App 的大瀑布流展示逻辑。
