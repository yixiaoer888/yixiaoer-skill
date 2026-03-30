# 获取合集列表 (Get Collections)

获取账号已创建的合集列表。主要用于抖音、今日头条等支持合集的平台。

## 核心指令 (Command)

```bash
node scripts/get-collections.ts --account_id=ACCOUNT_ID
```

## 参数列表 (Properties)

| 参数名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `--account_id` | `string` | 是 | 账号 ID。 |

## 返回结果 (Response)

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
