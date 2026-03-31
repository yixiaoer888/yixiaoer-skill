# 查询账号列表 (Query Account List)

获取当前租户下绑定的自媒体平台账号列表及其 UID 信息。

## 场景描述 (Usage)

- "帮我列出我在这台电脑上绑定的所有抖音账号。"
- "我需要在发布前确认视频号的 UID。"
- "查询我旗下的某个分组内的所有账号。"

## 参数定义 (Parameters)

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `platform` | `string` | 否 | 指定展示某个平台。见 [平台定义](./platform.md)。 |
| `name` | `string` | 否 | 按账号昵称模糊查询。 |
| `group` | `string` | 否 | 按分组名称查询。 |
| `page` | `number` | 否 | 当前页码，默认 `1`。 |
| `size` | `number` | 否 | 每页数量，默认 `20`，最大 `1000`。 |
| `platforms` | `string[]` | 否 | 平台批量查询，支持多个。见 [平台定义](./platform.md)。 |
| `platformType` | `number[]` | 否 | 平台类型批量查询。见 [平台定义-平台类型](./platform.md#平台类型枚举-platformtype)。 |
| `loginStatus` | `number` | 否 | 账号登录状态查询。见下方【登录状态枚举】。 |
| `isolation` | `string` | 否 | 是否显示隔离数据，`true`|`false`，默认 `false`。 |
| `parentId` | `string` | 否 | 父级ID，用于查询子账号列表。 |
| `time` | `number` | 否 | 查询此时间戳之后是否有新消息（unix ms）。 |

### 枚举定义 (Enumerations)

#### 登录状态 (LoginStatus)
| 值 | 定义 | 说明 |
| :--- | :--- | :--- |
| `0` | `Never` | 未曾登录/待登录 |
| `1` | `Succesed` | 登录成功 |
| `2` | `Expired` | 登录过期/失效。**提示：在查询时传 2 会同时匹配 2(过期), 3(失败), 4(取消授权)** |
| `3` | `Failed` | 登录失败 |
| `4` | `CancelAuth` | 取消授权 |


## 返回结果说明 (Response Details)

脚本返回标准的 JSON 数组或对象。每一个账号对象 (`PlatformAccountDTO`) 包含：

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `id` | `string` | 平台内部账号 ID |
| `platformAccountName` | `string` | 账号昵称 |
| `platformName` | `string` | 平台名称（如：抖音、百家号） |
| `platformAvatar` | `string` | 账号头像地址 |
| `platformAuthorId` | `string` | 平台方的 UID / 作者 ID |
| `status` | `number` | 登录状态（见上方枚举） |
| `platformType` | `number` | 平台类型（见上方枚举） |
| `parentId` | `string` | 父账号 ID（针对微信视频号等子账号结构） |
| `groups` | `string[]` | 账号所属的分组 ID 列表 |
| `proxyId` | `string` | 绑定的团队代理 ID (如果有) |
| `kuaidailiArea` | `string` | 绑定的内置代理地区编码 (如果有) |
| `favorites` | `Object[]` | 账号收藏夹列表。包含 `id`, `name`, `websiteUrl` |
| `isOperate` | `boolean` | 当前用户是否拥有该账号的运营权限 |
| `isFreeze` | `boolean` | 账号是否被冻结 |
| `createdAt` | `number` | 账号创建时间戳 |

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"accounts","platform":"抖音","name":"昵称","page":1,"size":20}'
```

## 注意事项
- 请确保环境变量 `YIXIAOER_API_KEY` 已设置。
