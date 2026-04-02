# 获取音乐列表 (Get Music)

获取发布内容时可供选择的平台背景音乐列表。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"music","account_id":"XXX","keyword":"周杰伦"}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `account_id`| `string` | **是** | 蚁小二账号 ID (32位十六进制) |
| `keyword`   | `string` | **是** | 搜索音乐关键词 |

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

## 脚本逻辑 (Backend)

- **脚本路径**: `scripts/api.ts`
- **功能**: 封装蚁小二标准化音乐查询接口 (`GET /platform-accounts/{platformAccountId}/music`)。
- **参数映射**: 将 `account_id` 映射为 URL 路径变量，将 `keyword`, `nextPage`, `categoryId` 映射为查询参数。
