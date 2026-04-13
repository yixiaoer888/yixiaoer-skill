# 📄 快手-Open 视频 参数 (Kuaishou-Open Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户通过快手开放平台接口分发短视频时触发：
- **短视频同步**：将竖屏短视频同步至快手社区。
- **可见性控制**：设置视频为公开、私密或仅好友可见。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装 Kuaishou-Open Payload 时需遵守：
1. **可见性强校验**：必须明确 `visibleType` 数值。Agent 应根据用户意图（如“发个私密动态”）精准设置字段。
2. **定时发布自省**：支持 `scheduledTime`，但快手对定时任务的最小间隔有要求，Agent 建议时间点应晚于当前 15 分钟。
3. **资源引用规范**：必须通过 `upload` 动作产生视频 key 及封面 key。
4. **极简原则**：虽然支持 `title` 和 `description`，但快手核心推荐逻辑高度依赖标签（通过 HTML 话题或 description 文本体现）。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`visibleType`** | `number` | **是** | **可见类型**: `0`-公开, `1`-私密, `3`-好友可见。 | `0` |
| `title` | `string` | 否 | 视频标题。 | - |
| `description` | `string` | 否 | 视频描述。推荐包含 #标签。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布快手 Open 版公开视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["KuaishouOpen"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "KS_OPEN_ACC_01",
        "video": { "key": "ks_v_1", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "快手开放平台实测案例",
          "description": "内容分享 #短视频 #实测",
          "visibleType": 0
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
| **可见性参数错误** | `visibleType` 传入了非定义数值（如 2）。 | 修正为 0, 1 或 3。 |
| **视频比例拦截** | 视频不符合快手主流的 9:16 或 3:4 比例。 | 强制建议用户使用竖屏素材。 |
| **定时任务失败** | `scheduledTime` 时间戳过期或不足 15 分钟间隔。 | 重新校正发布时间点。 |
| **描述字数过长** | description 超过了快手平台限制。 | 建议控制在 200 字符以内。 |

---
> [!TIP]
> **老铁经济**: 快手强调真实性与互动。Agent 建议在 `description` 中多使用口语化表达，并配以热门话题以融入社区氛围。
