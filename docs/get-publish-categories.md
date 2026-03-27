# 获取发布分类 (Get Publish Categories)

该指令用于获取特定账号在特定分发模态下的分类列表。

## 场景描述 (Usage)

- "帮我列出我在这台电脑上百家号所有可用的文章分类。"
- "我想发个科技文章，看看百家号的分类 ID 是多少。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `--account_id` | `string` | **必填**。目标账号 ID |
| `--type` | `string` | **必填**。`article`, `image-text`, `video` |

## 脚本逻辑 (Backend)

- **脚本路径**: `../scripts/get-publish-categories.ts`
- **实际接口**: `POST https://www.yixiaoer.cn/api/web/config-data/publish-category-tasks`
- **调用示例**: `node scripts/get-publish-categories.ts --account_id=bjh_123 --type=article`

## 输出结果 (Output)

脚本返回标准的 JSON 数组对象。每一个分类对象应包含 `yixiaoerId`, `yixiaoerName` 等。
请确保环境变量 `YIXIAOER_API_KEY` 已设置。
