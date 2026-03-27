# 查询账号列表 (Query Account List)

获取当前租户下绑定的自媒体平台账号列表及其 UID 信息。

## 场景描述 (Usage)

- "帮我列出我在这台电脑上绑定的所有抖音账号。"
- "我需要在发布前确认视频号的 UID。"
- "查询我旗下的某个分组内的所有账号。"

## 参数定义 (Parameters)

### 请求参数 (Query Parameters)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `platform` | `string` | **可选**。指定展示某个平台。见 [平台定义](./platform.md)。 |
| `name` | `string` | **可选**。按账号昵称模糊查询。 |
| `group` | `string` | **可选**。按分组名称查询。 |
| `page` | `number` | **可选**。当前页码，默认 `1`。 |
| `size` | `number` | **可选**。每页数量，默认 `20`，最大 `1000`。 |
| `platforms` | `string[]` | **可选**。平台批量查询，支持多个。见 [平台定义](./platform.md)。 |
| `platformType` | `number[]` | **可选**。平台类型批量查询。见 [平台定义-平台类型](./platform.md#平台类型枚举-platformtype)。 |
| `loginStatus` | `number` | **可选**。账号登录状态查询。见下方【登录状态枚举】。 |
| `isolation` | `string` | **可选**。是否显示隔离数据，`true`\|`false`，默认 `false`。 |
| `parentId` | `string` | **可选**。父级ID，用于查询子账号列表。 |
| `time` | `number` | **可选**。查询此时间戳之后是否有新消息（unix ms）。 |

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
| `favorites` | `Object[]` | 账号收藏夹列表。包含 `id`, `name`, `websiteUrl` |
| `isOperate` | `boolean` | 当前用户是否拥有该账号的运营权限 |
| `isFreeze` | `boolean` | 账号是否被冻结 |
| `createdAt` | `number` | 账号创建时间戳 |

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/query-accounts.ts`
- **实际接口**: `GET https://www.yixiaoer.cn/api/platform-accounts`
- **调用示例**: 
  - `node query-accounts.ts --platform=抖音 --name=倒霉蛋`
  - `node query-accounts.ts --page=1 --size=50`

## 注意事项

- 请确保环境变量 `YIXIAOER_API_KEY` 已设置，目前该接口优先通过 API Key 进行租户鉴权。
- 复杂数组参数（如 `platforms[]`）在目前的脚本中可能需要通过特定方式传入，建议优先使用单选 `platform` 参数。
