# 获取作品数据 (Get Content Overviews)

获取当前团队下各账号发布的作品数据统计信息。

## 场景描述 (Usage)

- "帮我查询一下最近 7 天所有抖音账号的作品播放量。"
- "我想看看某个特定账号发布的所有视频的数据。"
- "导出本月所有小红书账号的作品表现。"

## 参数定义 (Parameters)

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | 是 | 固定值：`content-overviews` |
| `platform` | `string` | 否 | 指定查询某个平台。见 [平台定义](./platform.md)。 |
| `platformAccountId` | `string` | 否 | 媒体账号 ID (ObjectId)。 |
| `publishUserId` | `string` | 否 | 发布人 ID (ObjectId)。 |
| `type` | `string` | 否 | 作品类型：`article` (文章), `video` (视频), `miniVideo` (短视频), `dynamic` (动态)。 |
| `title` | `string` | 否 | 作品标题，支持模糊查询。 |
| `publishStartTime` | `number` | 否 | 发布时间范围开始（Unix 时间戳，毫秒）。 |
| `publishEndTime` | `number` | 否 | 发布时间范围结束（Unix 时间戳，毫秒）。 |
| `page` | `number` | 否 | 当前页码，默认 `1`。 |
| `size` | `number` | 否 | 每页数量，默认 `10`。 |

## 返回结果说明 (Response Details)

返回包含作品数据列表的对象。主要字段如下：

### 复杂对象：data.data (作品对象)

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `platformAccountId` | `string` | 媒体账号 ID |
| `accountName` | `string` | 账号昵称 |
| `publishUserName` | `string` | 发布人昵称 |
| `contentData` | `object` | 原始统计数据。包含 `reCommand` (推荐), `play` (播放), `read` (阅读), `great` (点赞), `comment` (评论), `share` (分享), `collect` (收藏) 等。 |
| `updatedAt` | `number` | 数据更新时间戳 |

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"content-overviews","platform":"抖音","page":1,"size":10}'
```

## 注意事项
- 查询大范围数据时，建议合理设置 `publishStartTime` 和 `publishEndTime` 以提高查询效率。
