# 📄 大鱼号 视频 参数 (DaYuHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户明确要在“大鱼号 (阿里大文娱)”分发内容时触发：
- **阿里生态投放**：将视频推送至优酷、UC 浏览器、土豆等渠道。
- **原创深度运营**：标注原创身份并添加 AI 合成声明。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装大鱼号视频 Payload 时需遵守：
1. **标签必填逻辑**：大鱼号强制要求 1-6 个标签 (`tags`)。Agent 应根据视频核心内容进行提取。
2. **分类级联原则**：虽然可选，但建议提供分类以获得更精准的 UC 信息流推荐。必须透传 `raw` 对象。
3. **字数风向标**：标题上限 50 字符，描述上限 1000 字符。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (最多 50 字符)。 | - |
| **`description`** | `string` | **是** | 视频描述 (最多 1000 字符)。 | - |
| **`tags`** | `string[]` | **是** | 视频标签 (1-6 个)。 | - |
| **`pubType`** | `number` | **是** | **发布类型**: `0`-草稿, `1`-直接发布。 | `1` |
| `category` | `Array` | 否 | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| `createType` | `number` | 否 | **创作类型**: `0`-非原创, `1`-原创。 | `0` |
| `declaration` | `number` | 否 | **声明**: `0`-不声明, `3`-虚构演绎, `4`-AI 生成。 | `0` |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `yixiaoerId`, `yixiaoerName` 及其对应的 **`raw`** 对象。

## 4. 执行指令示例 (Command)

```bash
# 发布大鱼号原创视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["大鱼号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_dy_v_01",
        "video": { "key": "dy_v_key", "size": 1024000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "阿里大文娱生态分发全攻略",
          "description": "探讨优酷与 UC 浏览器的流量对接逻辑。",
          "tags": ["运营", "阿里", "流量"],
          "createType": 1,
          "pubType": 1
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
| **标签数量不符** | `tags` 数组不在 1-6 个范围内。 | 请根据规范增删标签。 |
| **原创申明冲突** | `createType` 与实际账号权限不符。 | 检查账号是否通过原创作者认证。 |
| **封面上传失败** | `coverKey` 无效。 | 请重新上传图片至 OSS 获取 key。 |
| **定时任务过期** | `scheduledTime` 晚于当前系统时间。 | 校准 Unix 时间戳。 |

---
> [!TIP]
> **UC 流量爆发**: 大鱼号在 UC 浏览器具有极强的爆发力。Agent 建议标题多采用资讯类格式，并配以抓眼球的封面以适应信息流逻辑。
