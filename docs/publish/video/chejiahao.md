# 📄 车家号 视频 参数 (Chejiahao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“车家号”平台发布汽车相关的专业测试视频、新车导购或行业资讯时触发：
- **汽车视频分发**：利用汽车之家庞大的垂直流量。
- **精品发布**：发布具有原创属性的汽车专业视频。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装车家号视频 Payload 时需遵守：
1. **原创属性申明**：车家号对原创和首发有专门的权重逻辑。Agent 应根据用户提示设置 `type`。
2. **资源引用规范**：确保视频资源的 `key` 已通过 `upload` 动作获取。
3. **内容长度校验**：通常建议视频时长符合汽车之家社区分发习惯。
4. **意图确认**：向用户明确当前发布是否涉及“首发”特权。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`description`** | `string` | **是** | 视频描述。 | - |
| **`type`** | `number` | **是** | **创作类型**: `1`-原创, `3`-首发, `13`-原创首发。 | `1` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布车家号原创首发视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Chejiahao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "CHEJIAHAO_ACC_01",
        "video": { "key": "cjh_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "2026 豪华 SUV 深度横测",
          "description": "本视频深入对比了多款热门豪华 SUV。",
          "type": 13
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
| **发布类型不合法** | `type` 传值超出了 1, 3, 13 的范围。 | 修正传参数值。 |
| **标题重复** | 车家号对重复标题内容审核较严。 | 建议在标题中增加具体评测日期或子型号。 |
| **资源不存在** | `video.key` 无效。 | 请重新执行 `upload` 动作。 |
| **定时设置冲突** | `scheduledTime` 与已有的定时任务时间过于接近。 | 调整发布时间点。 |

---
> [!TIP]
> **汽车垂直流量**: 车家号是高价值汽车购买意向人群的聚集地，Agent 建议在 `description` 中多提及具体的配置参数以增加被搜索到的概率。
