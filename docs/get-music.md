# 获取音乐素材列表 (Get Music)

获取发布内容时可供选择的平台背景音乐素材。支持平台：**抖音、快手、视频号**。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"music","account_id":"XXX","keyword":"周杰伦"}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | **是** | 蚁小二账号 ID (32位十六进制) |
| `keyword` | `string` | 否 | 搜索音乐关键词 |
| `categoryId` | `string` | 否 | 音乐分类 ID。**需与 categoryName 同时提供，仅支持抖音**，优先级高于 `keyword` |
| `categoryName`| `string` | 否 | 音乐分类名称。**需与 categoryId 同时提供，仅支持抖音** |
| `nextPage` | `string` | 否 | 分页标识，在返回结果中获取 |

## 返回结果 (Response)

返回一个包含音乐素材对象的数组。可以直接将其中的 `raw` 对象作为 `music` 参数传递给发布脚本。

```json
[
  {
    "yixiaoerId": "music_123",
    "yixiaoerName": "稻香",
    "authorName": "周杰伦",
    "playUrl": "http://...",
    "duration": 240,
    "raw": { ... }
  }
]
```

### 复杂对象：MusicItem
- `yixiaoerId`: 情况内部音乐 ID。
- `yixiaoerName`: 歌曲名。
- `authorName`: 歌手/作者名。
- `playUrl`: 试听链接。
- `duration`: 时长（秒）。
- `raw`: 原始平台返回的音乐对象，发布时需完整透传给 `music` 字段。

## 脚本逻辑 (Backend)

- **脚本路径**: `scripts/api.ts`
- **功能**: 封装蚁小二标准音乐素材查询接口 (`GET /platform-accounts/{platformAccountId}/music`)。
- **参数映射**: 将 `account_id` 映射为路径变量，将 `keyword`/`keyWord`, `nextPage`, `categoryId` 映射为查询参数。
