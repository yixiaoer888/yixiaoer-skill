# 保存草稿 (Draft Management)

在蚁小二生态中，存在两种不同维度的“草稿”概念，开发者需根据业务场景进行区分。

## 概念区分 (Draft Types)

| 类型 | 蚁小二草稿 (YiXiaoEr Draft) | 平台草稿 (Platform Draft) |
| :--- | :--- | :--- |
| **定义** | 仅保存在蚁小二系统数据库中 | 发送到目标平台（如抖音、B站）的草稿箱 |
| **是否触发任务** | **否**。仅做内容存储，不触发机器人 | **是**。机器人会登录平台并执行“存草稿”操作 |
| **调用 Action** | `save-draft` | `publish` |
| **核心参数** | `isDraft: true` | `contentPublishForm.pubType: 0` |
| **主要用途** | 跨端同步编辑、团队预审 | 将内容预推送到平台后台，方便手动二次微调 |

---

## 1. 蚁小二草稿 (YiXiaoEr Draft)

使用 `action: "save-draft"` 将发布任务整体保存至蚁小二后台。

### 调用指令 (Command)

```bash
node scripts/api.ts --payload='{
  "action": "save-draft",
  "publishType": "verticalVideo",
  "platforms": ["抖音", "视频号"],
  "isDraft": true,
  "desc": "这是一个蚁小二草稿",
  "platformAccounts": [
    {
      "platformAccountId": "67fb2f1735eeb3cf31db3d65",
      "videoKey": "v-xxxxxx",
      "coverKey": "c-xxxxxx"
    }
  ]
}'
```

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | 是 | 固定值: `save-draft` |
| `taskSetId` | `string` | 否 | 任务集 ID。传值则为更新草稿，不传则为创建新草稿 |
| `publishType` | `string` | 是 | 发布类型。可选: `verticalVideo`, `article`, `imageText` 等 |
| `platforms` | `string[]` | 是 | 目标平台名称列表 |
| `isDraft` | `boolean` | 是 | 必须为 `true` |
| `desc` | `string` | 否 | 任务集描述 |

---

## 2. 平台草稿 (Platform Draft)

使用 `action: "publish"` 但通过 `pubType` 参数控制机器人执行行为。

### 调用指令 (Command)

```bash
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_vid_003",
        "video": { "key": "v_key" },
        "coverKey": "c_key",
        "contentPublishForm": {
          "pubType": 0,
          "title": "存入抖音草稿箱的内容"
        }
      }
    ]
  }
}'
```

> [!TIP]
> 平台草稿的 `pubType` 字段定义：
> - `0`: 存入平台草稿箱 (Save to platform draft box)
> - `1`: 直接发布 (Publicly publish)
