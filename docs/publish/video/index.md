# 视频发布 (Video Publish)

> [!CAUTION]
> **阅读规范 (Reading Protocol)**:
> 本文档是 **所有平台** 视频发布的 **唯一入口** 和 **基础 DTO 定义**。
> 在查阅具体的平台文档（如 `douyin.md`）之前，你 **必须** 首先查阅本文档以理解 Payload 的根结构，否则将导致生成的 JSON 无法通过校验。

## 触发场景 (Trigger)
- **意图辨析**：当用户下达分发视频指令（无论是单平台发布还是多平台矩阵分发）时触发。涵盖从本地视频上传到最终推送的全链路。
- **典型提示词**：
  - “把这个视频发布到全平台”
  - “帮我同步这个短剧到抖音和快手”
  - “我的视频号要更新了，标题是 XXX”
  - “使用本地模式分发这个视频”

## 执行逻辑 (Logic Flow)
1. **资源预处理**：
   - 调用 `upload` action 将本地或 URL 视频及封面图上传至云端，并持有获得的 `key`。
   - 严禁在 `publish` 负载中直接透传原始 URL。
2. **账号与平台选取**：识别目标 `platforms` 列表及具体的 `platformAccountId`（通过 `accounts` 查询）。
3. **参数深度补全**：若涉及分类、地理位置、音乐等动态字段，调用对应 `get-*` 接口获取合法 ID。
4. **Payload 装配**：按照本文档 1.1 - 1.3 节定义的 DTO 结构，组装包含 `action: "publish"` 的完整 JSON。
5. **指令交付**：调用 `node scripts/api.ts --payload='{...}'` 执行发布。
6. **状态跟踪**：记录返回的 `taskSetId`，以便后续通过 `records` 查询进度。

## 1. 数据结构 (Data Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

### 1.1 基础结构 (Base Structure)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `action` | `string` | **是** | 固定值：`publish` | - |
| `publishType` | `string` | **是** | 固定为 `video` | - |
| `platforms` | `string[]` | **是** | 目标平台枚举数组，详见下方平台列表 | - |
| `coverKey` | `string` | **是** | 任务封面资源 Key | - |
| `publishArgs` | `Object` | **是** | 发布参数核心容器 | - |
| `taskSetId` | `string` | 否 | 任务集唯一标识 (草稿发布时必填) | - |
| `desc` | `string` | 否 | 任务描述/摘要 | - |
| `publishChannel` | `string` | 否 | `cloud` (云端) 或 `local` (本机) | `cloud` |
| `clientId` | `string` | 否 | 客户端连接 ID (`local` 发布时必填) | - |
| `isDraft` | `boolean` | 否 | 是否仅保存为草稿 (蚁小二草稿) | `false` |

### 1.2 发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `accountForms` | `Array` | **是** | 账号发布表单列表 | - |

### 1.3 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID | - |
| `video` | `Object` | **是** | **VideoFormItem**: 视频对象 (`key`, `width`, `height`, `size`) | - |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 | - |
| `contentPublishForm`| `Object` | **是** | **透传层**: `{}` | - |
| `coverKey` | `string` | **是** | 账号级封面 Key (必须与 `cover.key` 一致) | - |
| `fps` | `number` | 否 | 视频发布帧率 (海外平台使用) | - |

## 2. 发布示例 (Payload Example)

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["抖音"],
  "coverKey": "video_cover_key",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_vid_003",
        "video": {
          "key": "video_oss_key",
          "width": 1080,
          "height": 1920,
          "size": 52428800
        },
        "coverKey": "video_cover_key",
        "cover": {
          "key": "video_cover_key",
          "width": 1080,
          "height": 1920,
          "size": 307200
        },
        "contentPublishForm": {
          "formType": "task"
        }
      }
    ]
  }
}
```

## 3. 支持平台列表 (Support Platforms)

以下平台支持通过 `publishType: "video"` 进行发布。

| 平台名称 | 标识符 | 文档链接 |
| :--- | :--- | :--- |
| **头条号** | `Toutiaohao` | [toutiaohao.md](./toutiaohao.md) |
| **哔哩哔哩** | `Bilibili` | [bilibili.md](./bilibili.md) |
| **抖音** | `Douyin` | [douyin.md](./douyin.md) |
| **视频号** | `Shipinghao` | [shipinghao.md](./shipinghao.md) |
| ... | ... | ... |

> [!TIP]
> 完整列表请参考 [SKILL.md](../../SKILL.md)。
