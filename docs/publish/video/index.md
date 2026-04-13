# 📄 视频 发布 参数 (Video Publish Index)

本文档定义了蚁小二所有平台视频发布的 **根 Payload 结构**。在查阅具体平台（如抖音、视频号）文档前，**必须**先理解并遵循本索引定义的 DTO 规范。

## 1. 触发场景 (Trigger)

当用户下达分发视频指令（单账号或多账号矩阵）时触发：
- **全平台同步**：“把这个视频发到我所有账号”。
- **特定平台分发**：“帮我同步这个短剧到抖音和快手”。
- **草稿保存**：“把这段素材先存到蚁小二视频草稿箱”。

## 2. 交互协议 (Interactive Protocol)

Agent 在构造发布指令前需遵守：
1. **资源先行原则**：禁止直接使用外部 URL。必须先调用 `upload` 动作将视频和封面物理上传至 OSS，获取 `key`。
2. **账号校验原则**：必须先通过 `accounts` 确认 `platformAccountId` 的 `status: 1` (在线)。
3. **分步确认原则**：构造好 Payload 后，必须列表展示：**[平台] - [账号昵称] - [标题]**，征得用户同意后再执行。
4. **模式判别原则**：明确区分 `publishChannel: "cloud"` (云端代理发布) 与 `"local"` (本机客户端发布)。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `publish` | 固定值。 |
| **`publishType`** | `string` | **是** | `video` | 发布类型。 |
| **`platforms`** | `string[]` | **是** | - | 目标平台标识数组（如 `["Douyin", "ShiPinHao"]`）。 |
| `publishArgs` | `object` | **是** | - | 详见下方 [3.1 publishArgs](#31-publishargs-定义)。 |
| `publishChannel` | `string` | 否 | `cloud` | 发布通道。`cloud` (云发布) 或 `local` (本机发布)。 |
| `isDraft` | `boolean` | 否 | `false` | 若为 `true`，则仅保存至蚁小二草稿箱，不推送到平台。 |
| `taskSetId` | `string` | 否 | - | 若从草稿箱重新发起发布，需透传此 ID。 |

### 3.1 publishArgs 定义

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`accountForms`** | `Array` | **是** | 账号发布表单列表。每个元素对应一个账号。 |

### 3.2 accountForms 元素定义 (AccountFormItem)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`platformAccountId`** | `string` | **是** | `accounts` 接口返回的 `id`。 |
| **`video`** | `object` | **是** | **VideoFormItem**: 须含 `key`, `width`, `height`, `size` (Bytes), `duration` (秒)。 |
| **`cover`** | `object` | **是** | **ImageFormItem**: 须含 `key`, `width`, `height`, `size`。 |
| **`coverKey`** | `string` | **是** | 必须与 `cover.key` 保持一致，用于快速检索。 |
| **`contentPublishForm`**| `object` | **是** | **核心差异层**: 内部字段见各平台专属文档。 |

## 4. 执行指令示例 (Command)

```bash
# 简单的视频发布指令
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Douyin"],
  "publishArgs": {
    "accountForms": [{
      "platformAccountId": "ACC_001",
      "video": {"key": "v_123", "width": 1080, "height": 1920, "size": 52428800},
      "cover": {"key": "c_123", "width": 1080, "height": 1920, "size": 307200},
      "coverKey": "c_123",
      "contentPublishForm": { "formType": "task", "title": "示例视频" }
    }]
  }
}'
```

---

## 5. 常见问题排查 (Troubleshooting)

| 报错信息 / 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **API 调用无响应** | 脚本未安装或 `YIXIAOER_API_KEY` 缺失。 | 检查环境变量配置。 |
| **JSON 校验失败** | 遗漏了 `publishType` 或 `coverKey`。 | 重新对照 [3. 根结构参数定义](#3-根结构参数定义)。 |
| **获取在线设备失败** | 设置了 `publishChannel: "local"` 但客户端不在线。 | 切换至 `cloud` 模式或启动蚁小二客户端。 |
| **资源非法 (400)** | `video.key` 传入了 HTTP URL。 | 必须通过 `upload-resource` 先获取 OSS key。 |

---
> [!IMPORTANT]
> **发布前置自检**：执行前必须确认 `publishArgs.accountForms` 中的每个 `platformAccountId` 都属于 `platforms` 数组中声明的平台，严禁跨平台错位。
