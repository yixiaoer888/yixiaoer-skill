# 📄 获取音乐素材 Query 参数 (Get Music Query)

获取在发布短视频内容时可供挂载的平台背景音乐素材。支持平台：**抖音、快手、视频号**。

> [!IMPORTANT]
> **资源透传原则**：在发布表单（如 `contentPublishForm.music`）中使用音乐时，**必须**完整保留并透传本接口返回的 `raw` 原始对象，严禁仅填写 ID 或名称。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户在录制/编辑视频后，希望添加特定配乐、搜寻热门歌曲或按分类查找场景音乐（如：励志、动感、抒情）。
- **典型提示词**：
  - “帮我找一下适合跳舞的抖音音乐”
  - “搜索周杰伦的《稻香》用于这个视频”
  - “我的视频号发布需要配乐”

## 2. 交互协议 (Interactive Protocol)

1. **权限前置**：必须先获取 `account_id`，因为不同账号在不同平的可选音乐库存在较大差异。
2. **分类优先**：若用户指定了特定分类（如“流行”），Agent 应优先通过 `music-categories` 获取分类 ID 后再调用本接口。
3. **列表交付**：将返回的 `yixiaoerName` (歌名) 和 `artist` (歌手) 展示给用户确认。

## 3. 参数 definition (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `music` | 固定值。 |
| **`account_id`** | `string` | **是** | - | 蚁小二账号 ID (`platformAccountId`)。 |
| `keyword` | `string` | 否 | - | 音乐搜索关键词。 |
| `categoryId` | `string` | 否 | - | 音乐分类 ID (**仅支持抖音**，需配合 `categoryName`)。 |
| `categoryName`| `string` | 否 | - | 音乐分类名称 (**仅支持抖音**)。 |
| `nextPage` | `string` | 否 | - | 分页标识。 |

### 3.1 返回结果结构 (MusicItem)

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | 统一音乐 ID |
| `yixiaoerName` | `string` | 歌曲名称 |
| `artist` | `string` | 歌手/作者名 |
| `playUrl` | `string` | 试听地址 |
| `duration` | `number` | 音乐时长 (秒) |
| `raw` | `object` | **核心透传对象**。 |

## 4. 执行指令示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"music","account_id":"67fb2f1735eeb3cf31db3d65","keyword":"周杰伦"}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **找不到特定歌曲** | 平台版权库未覆盖或该歌曲被下架。 | 建议更换关键词尝试搜索同风格音乐。 |
| **分類搜索无效** | `categoryId` 目前仅在抖音平台有效。 | 核实当前目标平台是否为“抖音”。 |
| **发布时配乐失败** | 未透传 `raw` 字段或 `raw` 数据已过期。 | 确保获取列表后立即执行发布挂载。 |

---
> [!TIP]
> **爆款配乐技巧**：建议优先选择“流行”或“热门”分类下的前 5 首音乐，这些音乐通常具有更高的流量权重。
