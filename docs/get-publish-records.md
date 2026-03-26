# 查询发布记录 (Query Publish Records)

该能力允许用户查询任务集（TaskSet）的历史发布记录，包括分页查看、类型筛选和关键词搜索。

## 场景示例 (Scenarios)
- "查看我最近发布的 5 条记录。"
* "查询昨天发布的所有视频任务。"
- "检查 ID 为 `xxx` 的任务发布状态。"

## 参数定义 (Parameters)
| 参数名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| --page | number | 否 | 页码，默认 1 |
| --size | number | 否 | 每页数量，默认 10 |
| --publish_type | string| 否 | 发布类型 (article/video/image-text) |
| --keywords | string | 否 | 标题描述关键词搜索 |
| --id | string | 否 | (仅限详情查询) 任务集 ID |

## 调用指令 (Commands)

### 1. 查询任务集列表
```bash
node scripts/get-publish-records.ts --page=1 --size=10
```

### 2. 查询特定任务详情 (子任务状态)
```bash
node scripts/get-publish-details.ts --id={TASK_SET_ID}
```

## 响应数据模型 (Response JSON)
数据包含任务集的 ID、状态 (taskSetStatus)、发布渠道 (publishChannel)、总账号数 (accountTotal) 等。详情查询则展开为具体平台的发布链接 (openUrl) 和错误信息 (errorMessage)。
