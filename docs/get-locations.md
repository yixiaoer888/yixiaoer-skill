# 获取地理位置 (Get Locations)

获取在发布内容时可选的地理位置列表（地址/带货地址）。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"locations","account_id":"XXX","keyword":"深圳","type":1}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | **是** | 蚁小二账号 ID (32位十六进制) |
| `keyword` | `string` | 否 | 搜索关键词 |
| `type` | `number` | 否 | 搜索类型：`1`: 搜索位置 (默认) |

## 返回结果 (Response)

返回一个 `PlatformDataItem` 数组。可以直接将其中的对象作为 `location` 参数传递给发布脚本。

```json
[
  {
    "id": "12345",
    "text": "深圳市南山区...",
    "raw": { ... }
  }
]
```
