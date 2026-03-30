# 获取合集列表 (Get Collections)

获取账号已创建的合集列表。主要用于抖音、今日头条等支持合集的平台。

## 调用指令 (Command)

```bash
node scripts/get-collections.ts --payload='{"account_id":"ACCOUNT_ID"}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | **是** | 蚁小二账号 ID (32位十六进制) |

返回一个 `Category` 数组。可以直接将其中的对象作为 `collection` 参数传递给发布脚本。

```json
[
  {
    "yixiaoerId": "col_123",
    "yixiaoerName": "美食系列",
    "raw": { ... }
  }
]
```
