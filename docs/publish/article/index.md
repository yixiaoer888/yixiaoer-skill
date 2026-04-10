# 文章发布 (Article Publish)

> [!CAUTION]
> **阅读规范 (Reading Protocol)**:
> 本文档是 **所有平台** 文章发布的 **唯一入口** 和 **基础 DTO 定义**。
> 在查阅具体的平台文档（如 `weixingongzhonghao.md`）之前，你 **必须** 首先查阅本文档以理解 Payload 的根结构，否则将导致生成的 JSON 无法通过校验。

## 触发场景 (Trigger)
- **意图辨析**：当用户下达长图文、深度文章或多图文消息（仅限微信公众号）分发指令时触发。
- **典型提示词**：
  - “发布这篇公众号文章，内容是...”
  - “帮我同步这篇长文章到知乎和 CSDN”
  - “把这一周的推文定时在周日发布”

## 执行逻辑 (Logic Flow)
1. **内容转换**：确保正文为 HTML 格式。若用户提供的是 Markdown 或纯文本，需先进行转换。
2. **资源补全**：
   - 调用 `upload` action 上传文章封面图。
   - 正文内嵌图片也建议先行上传获取 Key（视频号卡片、公众号卡片同理）。
3. **平台策略分配**：
   - **微信公众号**：必须单独发布，推荐使用 `platformForms` 结构。
   - **通用平台**（知乎、简书等）：使用 `accountForms` 结构进行分发。
4. **参数装配**：注入 `action: "publish"` 及其余 DTO 字段。
5. **指令执行**：调用 `node scripts/api.ts`。

## 1. 数据结构 (Data Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

### 1.1 基础结构 (Base Structure)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `action` | `string` | **是** | 固定值：`publish` | - |
| `publishType` | `string` | **是** | 固定为 `article` | - |
| `platforms` | `string[]` | **是** | 目标平台枚举数组，详见下方平台列表 | - |
| `coverKey` | `string` | **是** | 任务封面资源 Key | - |
| `publishArgs` | `Object` | **是** | 发布参数核心容器 | - |
| `taskSetId` | `string` | 否 | 任务集唯一标识 (草稿发布时必填) | - |
| `desc` | `string` | 否 | 任务描述/摘要 | - |
| `publishChannel` | `string` | 否 | `cloud` (云端) 或 `local` (本机) | `cloud` |
| `clientId` | `string` | 否 | 客户端连接 ID (`local` 发布时必填) | - |
| `isDraft` | `boolean` | 否 | 是否仅保存为 draft (蚁小二草稿箱) | `false` |

### 1.2 草稿模式选取 (Draft Selection)

| 场景 | 蚁小二草稿箱 | 目标平台草稿箱 |
| :--- | :--- | :--- |
| **位置** | `Payload` 根路径 | `accountForms` -> `contentPublishForm` |
| **参数** | `"isDraft": true` | `"pubType": 0` |
| **效果** | 仅保存在蚁小二系统，不发起平台推送 | 执行推送流程，但最终结果为平台端的草稿态 |
| **用户话术** | “存为蚁小二草稿”、“以后再发” | “存到百家号草稿箱”、“推送到知乎草稿” |

### 1.3 发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `content` | `string` | **是** | **文章正文**: HTML 格式字符串 | - |
| `accountForms` | `Array` | **是** | 账号发布表单列表 (定义目标账号) | - |
| `platformForms` | `Object` | 否 | **平台级表单**: 仅限 `微信公众号` 使用。按平台名称组织的共享配置字典 | - |

> [!IMPORTANT]
> **配置架构约束**:
> - **微信公众号专用性**: **微信公众号必须单独发布**。在一个发布请求中，如果包含微信公众号，则不能包含其他任何平台；反之亦然。
> - **platformForms**: **仅限微信公众号使用**。
> - **优先级**: 后端将优先尝试从 `platformForms` 中获取对应平台的配置，若不存在则回退至账号级的 `contentPublishForm`。

### 1.4 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID | - |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 (`key`, `width`, `height`, `size`) | - |
| `contentPublishForm`| `Object` | 否 | **账号级透传配置**: 若未配置 `platformForms` 则从此读取 | `{}` |
| `coverKey` | `string` | 否 | 账号级封面 Key (通常与 `cover.key` 一致) | - |

## 2. 发布示例 (Payload Example)

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["微信公众号"],
  "publishArgs": {
    "content": "<h1>演示文章标题</h1><p>这是一个演示文章的正文内容...</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_art_001",
        "cover": {
          "key": "article_cover_key",
          "width": 900,
          "height": 500,
          "size": 150000
        },
        "coverKey": "article_cover_key"
      }
    ]
  }
}
```

## 3. 支持平台列表 (Support Platforms)

| 平台名称 | 标识符 | 文档链接 |
| :--- | :--- | :--- |
| **微信公众号** | `WeiXinGongZhongHao` | [weixingongzhonghao.md](./weixingongzhonghao.md) |
| **知乎** | `ZhiHu` | [zhihu.md](./zhihu.md) |
| ... | ... | ... |
