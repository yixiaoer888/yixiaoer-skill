# 📄 美柚 视频 参数 (Meiyou Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“美柚”平台分发女性健康、育儿分享、生活小技巧或情感短视频时触发：
- **垂直女性流量同步**：触达美柚平台的高活跃女性用户群体。
- **精品发布**：发布符合女性关注偏好的各类短视频内容。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装美柚视频 Payload 时需遵守：
1. **标题先行原则**：美柚为极简发布模式，`title` 是唯一的必填核心表单项。
2. **调性对齐**：Agent 建议标题应贴合女性、生活或母婴等社区核心关注点。
3. **全量引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。
4. **定时自律**：支持 `scheduledTime`，请确保时间点准确且具备必要的发布权益。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布美柚育儿经验分享视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Meiyou"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "MEIYOU_ACC_01",
        "video": { "key": "my_v_1", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "新手妈妈必看的 5 个育儿黑科技"
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
| **无视频标题** | 缺失 `title` 字段。 | 请补全视频标题。 |
| **分类缺失** | 美柚通常由平台根据标题自动分区，无需传参。 | 若发布失败，请检查标题关键字是否明确。 |
| **描述不支持** | 误填了 `description`。 | 美柚视频发布不支持描述字段，请仅保留 title。 |
| **素材尺寸违规** | 平台更偏好 9:16 竖向素材。 | 建议用户转换视频比例。 |

---
> [!TIP]
> **女性群体触达**: 美柚是精细化的女性流量池。Agent 建议标题内容应多使用“贴心”、“实测”、“避坑”等情感或经验类词汇以增加亲密度。
