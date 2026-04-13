# 📄 易车号 视频 参数 (Yichehao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“易车号”平台发布汽车相关的专业测试视频、导购指南或行业动态时触发：
- **汽车专业分发**：将精品内容推送至易车 App 及其垂直社区。
- **资讯同步**：同步与汽车生活相关的长短视频内容。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装易车号视频 Payload 时需遵守：
1. **专业化表达**：易车号用户群体垂直且专业。Agent 建议标题和描述应侧重于参数、体验或具体车型。
2. **极简必填核项**：`title` 和 `description` 是系统的核心校验点。
3. **资源引用协议**：必须通过 `upload` 动作获取视频和封面的 `key`。
4. **定时发布自律**：支持 `scheduledTime`，请确保时间点处于未来且符合平台频率限制。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述文本。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布易车号汽车深度评测视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Yichehao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YC_ACC_01",
        "video": { "key": "yc_v_1", "size": 1024000, "width": 1920, "height": 1080, "duration": 300 },
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 款豪华轿车性能首测",
          "description": "深度解析：这款新车在百公里加速与智能辅助驾驶上的表现。"
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
| **标题包含敏感词** | 易车号对标题中的违规宣传语审核较严。 | 修改标题，移除夸大不实的表述。 |
| **描述字数不足**| 描述内容过短。 | 建议补充至少 20 字符的视频背景介绍。 |
| **素材宽高比不符** | 汽车评测偏好 16:9 横屏。 | 确认视频比例是否符合汽车垂类展示逻辑。 |
| **发布频率受限** | 同一账号触发了平台的防骚扰机制。 | 降低发布频次或安排定时发布。 |

---
> [!TIP]
> **汽车垂直流量**: 易车号聚集了高意向购车人群。Agent 建议在 `description` 中明确车型名称及核心参数，利用平台的数据库索引能力提升内容的被发现几率。
