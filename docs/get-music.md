# 获取音乐列表 (Get Music)

获取发布内容时可供选择的平台背景音乐列表。

## 核心指令 (Command)

```bash
node scripts/get-music.ts --account_id=ACCOUNT_ID --keyword="周杰伦"
```

## 参数列表 (Properties)

| 参数名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `--account_id` | `string` | 是 | 账号 ID。 |
| `--keyword` | `string` | 是 | 搜索音乐关键词。 |

## 返回结果 (Response)

返回一个 `PlatformDataItem` 数组。可以直接将其中的对象作为 `musice` (或 `music`) 参数传递给发布脚本。

```json
[
  {
    "id": "music_id",
    "text": "稻香",
    "raw": { ... }
  }
]
```
