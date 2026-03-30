# 查询发布记录 (Query Publish Records)

该能力允许用户查询任务集（TaskSet）的历史发布记录，包括分页查看、状态筛选、发布方式（本地/云端）和关键词搜索。

## 调用指令 (Command)

```bash
node scripts/get-publish-records.ts --payload='{"page":1,"size":10}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `page` | `number` | 否 | 当前页码，默认 `1` |
| `size` | `number` | 否 | 每页条数，默认 `20` |
| `platforms` | `string[]` | 否 | 按平台过滤，如 `["抖音"]` |
| `publishStatus` | `number` | 否 | 状态过滤 |
| `keywords` | `string` | 否 | 标题描述关键词搜索 |
| `status` | `TaskSetStatusEnum` | 否 | 任务集状态过滤 |
| `start_time` | `number` | 否 | 发布起始时间戳（毫秒） |
| `end_time` | `number` | 否 | 发布截止时间戳（毫秒） |

### 枚举值定义

#### TaskSetStatusEnum (任务集总状态)
- `pending`: 待发布
- `publishing`: 发布中
- `allsuccessful`: 全部发布成功
- `partialsuccessful`: 部分发布成功
- `allfailed`: 全部发布失败
- `cancel`: 已取消

## 调用示例 (Examples)

### 1. 分页查询成功记录
```bash
node scripts/get-publish-records.ts --payload='{"page":1,"size":10,"status":"allsuccessful"}'
```

### 2. 查询特定任务集的子任务执行状态
需要先从记录列表获取 `id`，然后调用 `get-publish-details`：
```bash
node scripts/get-publish-details.ts --payload='{"task_set_id":"TASK_SET_ID"}'
```

## 响应数据模型 (Response JSON)

### 任务集列表 (TaskSet List)
| 字段名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `id` | `string` | 任务集 ID |
| `title` | `string` | 标题 |
| `taskSetStatus` | `string` | 总体执行状态 |
| `publishChannel` | `string` | 发布渠道：`local` (本机), `cloud` (云端) |
| `accountTotal` | `number` | 涉及的总账号数 |
| `failedTotal` | `number` | 失败的账号数 |
| `createdAt` | `number` | 创建时间戳 |

### 子任务详情 (Sub-Task Detail)
| 字段名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `platformName` | `string` | 平台名称。见 [平台定义](./platform.md) |
| `platformAccountName` | `string` | 账号名称 (昵称) |
| `stageStatus` | `StageStatus` | 子任务当前阶段的状态 |
| `stages` | `TaskStages` | 子任务当前所处阶段 |
| `openUrl` | `string` | 发布成功后的内容链接（若有） |
| `errorMessage` | `string` | 失败原因描述 |

#### StageStatus (阶段状态)
- `waiting`: 等待中
- `running`: 执行中
- `success`: 成功
- `failed`: 失败
- `cancel`: 取消

#### TaskStages (任务阶段)
- `upload`: 上传
- `push`: 推送
- `transcoding`: 转码
- `review`: 审核
- `scheduled`: 定时
- `success`: 成功
