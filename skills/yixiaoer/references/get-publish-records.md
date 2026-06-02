# 查询发布记录 (Get Publish Records)

获取历史发布的任务集（TaskSet）概览列表，支持按平台、状态、时间等维度进行筛选。

## 触发场景 (Trigger)
- **意图辨析**：当用户需要确认任务是否发布成功、查看历史任务状态、获取任务 ID 以便查询详情或执行重发/删除操作时触发。
- **典型提示词**：
  - “查看我昨天的发布记录”
  - “帮我查一下抖音发布失败的任务”
  - “列出最近 10 条发布记录”
  - “确认一下任务 ID 为 TS123 的执行状态”

## 参数定义 (Parameters)

### 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `action` | `string` | **是** | 固定值：`records` |
| `platforms` | `string[]` | 否 | 平台过滤。见 [平台定义](./platform.md)。 |
| `status` | `number[]` | 否 | 状态过滤：0-发布中, 1-发布成功, 2-发布失败, 3-已取消, 4-待审核。 |
| `keyword` | `string` | 否 | 标题关键词模糊搜索 (别名: `keyWord`) |
| `start_time` | `number` | 否 | 开始时间戳 (Unix ms) |
| `end_time` | `number` | 否 | 结束时间戳 (Unix ms) |
| `page` | `number` | 否 | 分页，默认 1 |
| `size` | `number` | 否 | 每页数量，默认 20 |

## 执行逻辑 (Logic Flow)
1. **维度提取**：识别用户查询的时间范围、关键词及状态意图。
2. **参数装配**：构造 `action: "records"` 负载，处理 `status` 数组。
3. **指令执行**：调用 `yxer records [--platform P] [--limit N] [--status S] [--json]`。
4. **结果交付**：从 `data.data` 中提取关键字段（如 `task_set_id`, `status`）反馈给用户。若涉及失败任务，引导用户调用 `details` 进一步查询。

## 返回数据说明 (Response Details)

返回标准的任务集列表对象。每个任务集包含基础信息、状态汇总及关联的账号统计。

## 调用指令 (Command)

```bash
yxer records --status 1 --json
```

