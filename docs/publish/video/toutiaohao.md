# 📄 头条号 视频 参数 (TouTiaoHao Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“今日头条 (头条号)”分发视频，且需要配置合规申明、定时发布或可见性控制时触发：
- **流量激励申明**：申明原创身份（如自行拍摄）以获取更高推荐权。
- **定时分发**：在用户活跃高峰期自动发布视频。
- **隐私管理**：发布公开或私密视频内容。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装头条号视频 Payload 时需遵守：
1. **标题吸引力与长度**：头条号标题上限 80 字符。Agent 应确保标题兼具描述性与吸引力。
2. **标签强校验**：强制要求 1-5 个 `tags` 字符串数组。Agent 应根据视频核心动量自动提取。
3. **合规申明必填**：一点号对内容来源有严格合规要求。Agent 必须根据视频属性准确设置 `declaration`（如 AI 创作设为 `3`）。
4. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题 (1-80 字符)。 | - |
| **`description`** | `string` | **是** | 视频描述内容 (1-400 字符)。 | - |
| **`tags`** | `string[]` | **是** | 视频标签 (1-5 个)。 | - |
| **`visibleType`** | `number` | **是** | **可见性**: `0`-公开, `1`-私密。 | `0` |
| **`pubType`** | `number` | **是** | **发布类型**: `0`-草稿, `1`-直接发布。 | `1` |
| `declaration` | `number` | 否 | **创作者申明**: 1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎, 7-投资观点, 8-健康医疗。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

## 4. 执行指令示例 (Command)

```bash
# 发布头条号原创视频：带申明和定时
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_tt_v_01",
        "video": { "key": "tt_v_key", "size": 10240000, "width": 1920, "height": 1080, "duration": 60 },
        "contentPublishForm": {
          "formType": "task",
          "title": "今日头条 2026 内容创作全景扫描",
          "description": "本文为您解析 2026 年头条号的新算法与分发逻辑。",
          "tags": ["运营", "自媒体", "干货"],
          "declaration": 1,
          "visibleType": 0,
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
| **标签数量错误** | `tags` 成员少于 1 个或多于 5 个。 | 修正标签数组长度。 |
| **标题或描述超限** | 字数超出了 80/400 限制。 | 精简文案或自动截断。 |
| **创作申明不合法** | `declaration` 未传或数值不在 1-8 范围内。 | 准确设置合规来源申明。 |
| **定时任务未触发** | `scheduledTime` 过期或与服务器时钟同步失败。 | 重新校对 Unix 时间戳。 |

---
> [!TIP]
> **算法流量池**: 头条号对互动性强、申明合规且标题抓眼球的高质量内容有极大流量倾斜。Agent 建议配以高清横版封面以适配推荐流。
