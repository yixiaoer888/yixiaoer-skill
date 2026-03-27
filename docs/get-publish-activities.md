# 获取发布活动 (Get Publish Activities)

该指令用于获取特定账号在特定分发模态下可参与的征文活动列表。

## 场景描述 (Usage)

- "我的百家号账号 bj-123 最近有什么可以参加的文章征文活动吗？"

## 参数定义 (Parameters)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `--account_id` | `string` | **必填**。目标账号 ID |
| `--type` | `string` | **必填**。`article`, `image-text`, `video` |

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/get-publish-activities.ts`
- **实际接口**: `POST https://www.yixiaoer.cn/api/web/config-data/activity-tasks`
- **调用示例**: `node scripts/get-publish-activities.ts --account_id=bjh_123 --type=article`

## 输出结果 (Output)

脚本返回标准的 JSON 数组对象。每一个活动对象应包含 `id`, `name`, `raw` 等。
请确保环境变量 `YIXIAOER_API_KEY` 已设置。
