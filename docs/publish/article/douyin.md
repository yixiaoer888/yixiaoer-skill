# 📄 抖音文章发布参数 (Douyin Article)

> [!IMPORTANT]
> **阅读前置规则**: 本文档仅描述 `platformForms` (或账号级 `contentPublishForm`) 内部的平台差异化字段。在开始前，你 **必须** 已经掌握并应用了 [文章发布通用索引](./index.md) 中的根 Payload 结构。

## 1. 触发场景 (Trigger)

当用户指定在“抖音 (Douyin)”平台发布文章或图文内容时触发。支持：
- **图文分享**：发布多图文形式的文章内容。
- **背景音乐**：为文章配置背景音乐 BGM。
- **话题挂载**：添加抖音热门话题或专业话题。
- **可见性控制**：设置公开、私密或好友可见。

## 2. 交互协议 (Interactive Protocol)

Agent 在拼装抖音 Payload 时需遵守：
1. **话题动态检索**：抖音话题必须通过 `categories` 或对应的检索接口获取，并完整透传 `raw` 数据。
2. **音乐合规性**：若用户请求添加音乐，必须调用 `music` 接口获取合法 `MusicItem`。
3. **封面图组约束**：抖音文章建议提供多张封面图（1-9 张），以提升图文轮播效果。
4. **可见性默认为公开**：除非用户明确要求“悄悄发”或“私密”，否则 `visibleType` 应默认为 `0` (公开)。

## 3. 参数定义 (Parameters)

### 3.1 核心表单参数 (contentPublishForm)

| 字段名 | 类型 | 必填 | 描述 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| **`formType`** | `string` | **是** | 固定值: `task` | `task` |
| **`title`** | `string` | **是** | 文章标题 (最多 50 字)。 | - |
| **`content`** | `string` | **是** | 文章正文 (HTML 格式，最多 50000 字符)。 | - |
| **`covers`** | `Array` | **是** | 封面图列表。详见 [3.2 OldCover 定义](#32-oldcover-定义)。 | - |
| **`visibleType`** | `number` | **是** | **可见性**: `0`-公开, `1`-私密, `3`-好友。 | `0` |
| `description` | `string` | 否 | 文章描述或摘要 (最多 200 字)。 | - |
| `headImage` | `Object` | 否 | 文章头图。使用 `OldCover` 结构。 | - |
| `music` | `Object` | 否 | 背景音乐。详见 [3.3 MusicItem 定义](#33-musicitem-定义)。 | - |
| `topics` | `Array` | 否 | 话题标签列表。最多 5 个。 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳，单位: 秒)。 | - |

### 3.2 OldCover 定义
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `key` | `string` | **是** | OSS 资源 Key。 |
| `size` | `number` | **是** | 文件大小 (Bytes)。 |
| `width` | `number` | **是** | 宽度。 |
| `height` | `number` | **是** | 高度。 |

### 3.3 MusicItem 定义 (音乐)
| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 内部 ID。 |
| `yixiaoerName` | `string` | **是** | 歌曲名称。 |
| `duration` | `number` | **是** | 时长 (秒)。 |
| `playUrl` | `string` | **是** | 播放链接。 |
| `raw` | `object` | **是** | 原始数据。必须完整透传。 |

---

## 4. 执行指令示例 (Command)

```bash
# 发布抖音图文文章：带背景音乐与话题
node scripts/api.ts --payload='{
  "action": "publish",
  "publishType": "article",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "DY_ACC_123",
        "contentPublishForm": {
          "formType": "task",
          "title": "记录我的 AI 探索之旅",
          "content": "<h1>探索之旅</h1><p>正文内容...</p>",
          "covers": [
            { "key": "c_key_1", "size": 102400, "width": 1080, "height": 1440 }
          ],
          "music": {
            "yixiaoerId": "m1",
            "yixiaoerName": "好日子",
            "duration": 180,
            "playUrl": "http://...",
            "raw": {...}
          },
          "topics": [
            { "yixiaoerId": "t1", "yixiaoerName": "AI技术", "raw": {...} }
          ],
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
| **音乐版权拦截** | 选用的 BGM 不在抖音当前的曲库或有地域限制。 | 重新调用 `action: "music"` 获取推荐曲目。 |
| **封面尺寸建议** | 抖音针对竖屏优化，横向图片会被拉伸。 | 推荐使用 3:4 或 9:16 的竖版图片作为封面。 |
| **内容疑似搬运** | 抖音内容查重机制检测到高度雷同。 | 增加原创文字段落或修改封面图。 |
| **可见性设置失败** | 账号等级不足或被风控限制。 | 检查账号是否可以在抖音端手动正常更改可见性。 |

---
> [!TIP]
> **流量密码**: 在抖音发布图文时，话题的精准度 (`topics`) 比正文字数更重要。
