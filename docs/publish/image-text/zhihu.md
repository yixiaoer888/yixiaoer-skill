# 📄 知乎 图文 参数 (Zhihu Image-Text)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [图文发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“知乎”发布想法、动态或简短图文分析时触发：
- **想法投递**：发布类似朋友圈的碎碎念或动态。
- **互动提及**：在描述中 @好友、关联热门话题或参与知乎官方活动。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装知乎图文 Payload 时需遵守：
1. **HTML 标签闭环**：知乎图文描述高度依赖 HTML。所有正文需由 `<p>` 包裹，并精准嵌套 `<topic>`, `<friend>`, `<activity>` 等自定义标签。
2. **资源转换要求**：知乎最多支持 9 张图片，每张成员必须具备完整的 `OldImage` 结构且 Key 已通过 `upload` 转换。
3. **属性透传原则**：所有的标签（话题、好友、活动）中的 `raw` 属性必须存放完整的原始数据 JSON 序列化字符串。
4. **字数与阈值**：想法（图文）不宜过于冗长，Agent 应建议将长篇大论转为专栏文章。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`description`** | `string` | **是** | 图文内容。支持 HTML (`<p>`, `<topic>`, `<friend>`, `<activity>`)。 | - |
| `title` | `string` | 否 | 图文标题 (1-50 字符)。 | - |
| `images` | `Array` | 否 | 图片数组 (1-9 张)。使用 `OldImage[]` 结构。 | - |

### 3.2 复杂结构说明

- **OldImage**: 包含 `key`, `size`, `width`, `height`, `format`。
- **自定义标签属性**: 
  - **`<topic text='...' raw='...'>#名称</topic>`**: 最多 5 个话题。
  - **`<friend raw='...'>@好友名</friend>`**: 关联艾特好友。
  - **`<activity raw='...'>活动名</activity>`**: 参与特定官方活动。

## 4. 执行指令示例 (Command)

```bash
# 发布知乎想法：带话题和好友提及
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["Zhihu"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "今日学习打卡",
          "description": "<p>今天又学到了新东西 <topic text=\"知识\" raw=\"{\\\"yixiaoerId\\\":\\\"123\\\",\\\"yixiaoerName\\\":\\\"知识\\\",\\\"raw\\\":{}}\">#知识</topic> @<friend raw=\"{\\\"yixiaoerId\\\":\\\"456\\\",\\\"raw\\\":{\\\"nick\\\":\\\"李四\\\"}}\">李四</friend></p>",
          "images": [
            { "key": "zh_img_1", "size": 102400, "width": 1080, "height": 1440, "format": "jpg" }
          ]
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
| **话题不识别** | `topic` 标签中的 `raw` 数据格式不满足 JSON 序列化要求。 | 严格校验 `raw` 内容，确保包含必要的平台内部 ID。 |
| **好友提及失败** | `@好友` 后的 `raw` 属性中 ID 已注销或不合法。 | 请通过 `get-friends` 接口获取准确的好友元数据。 |
| **图片加载异常** | 使用了外部链接而非系统的 `key`。 | 必须先执行 `upload` 动作并引用产生的 OSS Key。 |
| **内容过短或空白** | `description` 字段内容缺失或仅包含空标签。 | 请补全博文描述文字。 |

---
> [!TIP]
> **知乎硬核社区**: 知乎用户更偏好理性的干货分享。Agent 建议在想法中通过话题钩子引导用户参与深度讨论，利用知乎的长尾搜索效应进行获客。
