# 📄 获取作品数据概览 Query 参数 (Get Content Overviews Query)

获取当前团队下各账号发布的作品（视频、文章、动态）的数据统计信息。用于分析单件作品的传播效果。

> [!NOTE]
> **数据实时性说明**：作品数据由蚁小二后台及客户端定期更新。若发现刚发布的视频没有播放数据，请引导用户稍等 1-2 小时。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户需要查询、对比或分析特定作品的表现指标（如阅读量、点赞量、分享量）。
- **典型提示词**：
  - “我的视频号作品最近 7 天的播放量如何？”
  - “查询标题中包含‘测评’的视频数据情况”
  - “由于最近流量下滑，帮我找出表现最差的作品”

## 2. 交互协议 (Interactive Protocol)

1. **时间解析**：Agent 应将用户提到的相对时间（如“上周”、“昨天”）自动转换为 Unix 毫秒时间戳。
2. **数据分类**：交付结果时，Agent 应对数据进行优先级排序，优先展示点赞、播放等核心转化指标。
3. **分步建议**：若查询结果中单条作品表现极佳，可建议用户针对该风格进行复刻。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `content-overviews` | 固定值。 |
| `platform` | `string` | 否 | - | 单平台过滤。 |
| `platformAccountId` | `string` | 否 | - | 指定媒体账号 ID。 |
| `publishUserId` | `string` | 否 | - | 按发布人成员 ID 过滤。 |
| `type` | `string` | 否 | - | 类型：`article`, `video`, `miniVideo`, `dynamic`。 |
| `title` | `string` | 否 | - | 标题关键词搜索。 |
| `publishStartTime` | `number` | 否 | - | 开始时间戳 (ms)。 |
| `publishEndTime` | `number` | 否 | - | 结束时间戳 (ms)。 |
| `page` | `number` | 否 | `1` | 页码。 |
| `size` | `number` | 否 | `20` | 每页数量。 |

### 3.1 返回结果结构 (Response Details)

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `accountName` | `string` | 账号昵称 |
| `title` | `string` | 作品标题 |
| `contentData.reCommand` | `number` | 平台推荐量 |
| `contentData.play` | `number` | 播放量 (视频) |
| `contentData.read` | `number` | 阅读量 (文章) |
| `contentData.great` | `number` | 点赞数 |
| `contentData.comment` | `number` | 评论数 |

## 4. 执行指令示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"content-overviews","platform":"抖音","page":1,"size":10}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **查询结果为空** | 时间跨度过小、标题匹配太严或所选平台无此类型作品。 | 建议扩大搜索时间范围或使用更通用的标题关键字。 |
| **数据长时间保持 0** | 尚未触发平台数据回调。 | 若作品已发布超过 24 小时仍为 0，建议检查客户端登录状态。 |
| **数据异常偏高/偏低** | 数据抓取逻辑受平台更新影响产生短暂偏差。 | 点击“刷新数据”或在蚁小二客户端手动同步作品。 |

---
> [!IMPORTANT]
> **维度声明**：作品数据专注于“单内容”的表现，若需要查看账号整体矩阵的演进，请使用 `account-overviews` 接口。
