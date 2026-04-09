# 草稿管理 (Draft Management)

在蚁小二生态中，存在两种“草稿”能力，需按目标选择：
- 保存到**蚁小二草稿箱**（`action: publish` + `isDraft: true`）
- 保存到**目标平台草稿箱**（`action: publish` + `pubType: 0`）

## 快速区分 (Draft Types)

| 类型 | 蚁小二草稿 (YiXiaoEr Draft) | 平台草稿 (Platform Draft) |
| :--- | :--- | :--- |
| **定义** | 仅保存在蚁小二系统数据库中 | 发送到目标平台（如抖音、B站）的草稿箱 |
| **是否触发任务** | 否 (仅存储) | 是 (执行推送流程) |
| **调用 Action** | `publish` | `publish` |
| **核心区分参数** | `isDraft: true` | `contentPublishForm.pubType: 0` |
| **主要用途** | 跨端同步编辑、团队预审 | 将内容预推送到平台后台，方便手动二次微调 |

---

## 存为蚁小二草稿 (`isDraft: true`)

当需要将任务暂存到蚁小二系统的草稿列表，而不启动发布流程时使用。

### 调用指令 (Command)

```bash
node scripts/api.ts --payload='{
  "action": "publish",
  "isDraft": true,
  "publishType": "video",
  "platforms": ["抖音", "视频号"],
  "desc": "这是一个蚁小二草稿",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "67fb2f1735eeb3cf31db3d65",
        "video": { "key": "v-xxxxxx" },
        "coverKey": "c-xxxxxx",
        "contentPublishForm": {
          "pubType": 1,
          "title": "这是一个蚁小二草稿"
        }
      }
    ]
  }
}'
```

---

## 存为平台草稿 (`pubType: 0`)

当需要启动发布流程，但最终结果是在第三方平台后台看到草稿时使用。

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
