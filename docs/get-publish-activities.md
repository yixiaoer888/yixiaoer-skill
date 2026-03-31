# 获取征文活动 (Get Publish Activities)

该指令用于获取特定账号在特定发布类型（视频/文章等）下的征文活动列表。征文活动通常带有特定的投稿要求和奖励。

## 场景描述 (Usage)

- "帮我查一下抖音账号 A 最近有哪些可以参加的征文活动。"
- "我想参加个关于科技的征文，看看我的百家号账号能报哪个。"

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"activities","account_id":"XXX","type":1}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | **是** | 蚁小二账号 ID (32位十六进制) |
| `type`       | `number` | **是** | 发布类型：`1`: 文章，`2`: 视频，`3`: 动态 |
| `categoryId` | `string` | 否   | 特定分类下的活动 |
| `keyWord`    | `string` | 否   | 搜索活动关键字 |

### 枚举值定义

#### ContentTypeEnum (发布类型)
- `video`: 视频
- `article`: 文章
- `dynamic`: 动态 (爬虫使用)
- `imageText`: 图文 (前端使用)

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/get-publish-activities.ts`
- **实际接口**: `GET https://www.yixiaoer.cn/api/v2/platform/accounts/:id/activities` (或通过 OpenPlatform 路由)
- **环境要求**: 需设置 `YIXIAOER_API_KEY` 环境变量。

## 输出结果 (Output)

脚本返回活动列表数组，每个对象包含以下核心字段：

| 字段名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | 活动在本系统的唯一标识 |
| `yixiaoerName` | `string` | 活动名称 |
| `yixiaoerDesc` | `string` | 活动描述 |
| `yixiaoerImageUrl` | `string` | 活动封面图 URL |
| `viewNum` | `string` | 浏览量/参与度指标 |
| `raw` | `object` | 平台原始活动数据（包含第三方平台的 Activity ID 等） |

