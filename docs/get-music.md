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

返回一个包含音乐素材对象的数组。发布时请将整个对象（使用 `MusicItem` 结构）作为 `music` 参数传递给发布脚本。如果在获取时 `raw` 字段有值，发布表单中必须完整保留并透传。

```json
[
  {
    "yixiaoerId": "music_123",
    "yixiaoerName": "稻香",
    "artist": "周杰伦",
    "playUrl": "http://...",
    "duration": 240,
    "raw": { ... }
  }
]
```

### 复杂对象：MusicItem
- `yixiaoerId`: (必填) 蚁小二端统一音乐 ID。
- `yixiaoerName`: (必填) 歌曲名称。
- `duration`: (必填) 音乐时长（单位：秒）。
- `playUrl`: (必填) 试听/播放链接。
- `artist`: (可选) 歌手/作者名。
- `raw`: (可选) 平台原始数据。如果在音乐列表获取时该字段存在，发布表单中必须携带并完整透传。

## 脚本逻辑 (Backend)

- **脚本路径**: `scripts/api.ts`
- **功能**: 封装蚁小二标准音乐素材查询接口 (`GET /platform-accounts/{platformAccountId}/music`)。
- **参数映射**: 将 `account_id` 映射为路径变量，将 `keyword`/`keyWord`, `nextPage`, `categoryId` 映射为查询参数。
