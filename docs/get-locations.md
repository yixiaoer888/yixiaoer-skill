# 获取地理位置 (Get Locations)

获取在发布内容时可选的地理位置列表（地址/带货地址）。

## 核心指令 (Command)

```bash
node scripts/get-locations.ts --account_id=ACCOUNT_ID --keyword="深圳" --type=1
```

## 参数列表 (Properties)

| 参数名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `--account_id` | `string` | 是 | 账号 ID。 |
| `--keyword` | `string` | 否 | 搜索关键词。 |
| `--type` | `number` | 否 | 搜索类型。`1`: 搜索位置（默认）；其他类型视平台支持而定。 |

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
