# 获取合集列表 (Get Collections)

获取账号已创建的合集列表。主要用于抖音、今日头条等支持合集的平台。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"collections","account_id":"ACCOUNT_ID"}'
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

## 脚本逻辑 (Backend)

- **脚本路径**: `scripts/api.ts`
- **功能**: 封装蚁小二标准化合集查询接口 (`GET /platform-accounts/{platformAccountId}/collections`)。
- **参数映射**: 将 `account_id` 映射为 URL 路径变量，将 `type` 映射为 `publishType` 查询参数。
