# 📄 美拍 视频 参数 (Meipai Video)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [视频发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“美拍”平台发布美妆、穿搭、才艺展示或高颜值短视频时触发：
- **精品短视频分发**：利用美拍的高颜值社区属性进行内容触达。
- **高活跃分享**：同步生活方式相关的精致短片。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装美拍视频 Payload 时需遵守：
1. **分类精准对齐**：美拍对分类（频道）有较强的依赖。必须通过 `categories` 接口获取并在 `category` 数组中透传 `raw` 原始数据。
2. **标题必填原则**：美拍对 `title` 是强校验。Agent 需确保其具备足够的吸引力且不超过字数限制。
3. **资源引用规范**：必须通过 `upload` 动作获取视频和封面的 `key`。
4. **定时发布校验**：支持 `scheduledTime`，请确保时间点处于未来时。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定为 `task`。 | `task` |
| **`title`** | `string` | **是** | 视频标题。 | - |
| **`category`** | `Array` | **是** | **视频分类**: 使用 `CascadingPlatformDataItem[]` 结构。 | - |
| `description` | `string` | 否 | 视频描述内容。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 (单位: 秒)。 | - |

### 3.2 复杂结构说明

- **CascadingPlatformDataItem**: 必须包含 `id`, `text`, `raw` 元数据。

## 4. 执行指令示例 (Command)

```bash
# 发布美拍精品短视频
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "video",
  "platforms": ["Meipai"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "MP_ACC_01",
        "video": { "key": "mp_v_1", "size": 1024000, "width": 1080, "height": 1920, "duration": 15 },
        "contentPublishForm": {
          "formType": "task",
          "title": "夏日清透果汁感妆容教学",
          "description": "今天教大家画一个超简单的夏日妆容！",
          "category": [{ "id": "1", "text": "生活", "raw": {...} }]
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
| **标题缺失** | 未传 `title` 字段。 | 美拍要求标题必填，请补全。 |
| **分类不匹配** | `category.raw` 中的频道数据已过期。 | 重新执行 `categories` 动作以刷新获取。 |
| **视频画质被降权** | 视频清晰度低于 720P。 | 建议使用 1080P 素材以获得平台精品推荐。 |
| **描述包含敏感词** | 美拍社区环境审核较严。 | 检查并替换夸大其词或诱导类的文案。 |

---
> [!TIP]
> **颜值社区属性**: 美拍用户对画面美感要求极高。Agent 建议视频封面应选用色彩明亮的精美截图，并配合带有正向情绪的标题以提升转化。
