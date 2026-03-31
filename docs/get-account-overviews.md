# 获取账号数据概览 (Get Account Overviews V2)

获取当前团队下各账号的综合表现概览。主要用于性能监控和核心指标汇总。

## 场景描述 (Usage)

- "帮我汇总一下当前团队下所有抖音账号的粉丝总量和互动量。"
- "我想看某个特定负责人的旗下账号的整体表现。"
- "最近 30 天各个平台的账号数据增长情况。"

## 参数定义 (Parameters)

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | 是 | 固定值：`account-overviews` |
| `platform` | `string` | 是 | 必填，指定查询某个平台。见 [平台定义](./platform.md)。 |
| `name` | `string` | 否 | 按照账号昵称模糊查询。 |
| `group` | `string` | 否 | 按照分组名称查询。 |
| `loginStatus` | `number` | 否 | 账号状态 (0未登录, 1成功, 2过期, 3失败)。 |
| `memberIds[]` | `string[]` | 否 | 负责人成员 ID 数组 (ObjectId)。 |
| `page` | `number` | 否 | 当前页码，默认 `1`。 |
| `size` | `number` | 否 | 每页数量，默认 `10`。 |

## 返回结果说明 (Response Details)

返回包含账号表现信息列表的对象。主要字段如下：

### 复杂对象：data.data[i] (账号概览对象)

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `platformAccountId` | `string` | 媒体账号 ID |
| `platformAvatar` | `string` | 账号头像地址 |
| `platformAccountName` | `string` | 账号昵称 |
| `principalName` | `string` | 当前负责成员名 |
| `overviewData` | `object` | 统计数据总计。包含 `fansTotal` (总粉丝数), `playTotal` (总播放量), `commentsTotal` (总评论量), `likesTotal` (总点赞量), `favoritesTotal` (总收藏量) 等。 |

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"account-overviews","platform":"抖音","page":1,"size":10}'
```

## 注意事项
- 账号数据 V2 提供的是总量的统计概览，相比 `/contents/overviews` 粒度更粗，适合查看整体情况。
- 请确保环境变量 `YIXIAOER_API_KEY` 已设置。
