# 获取音乐分类 (Get Music Categories)

获取在选择音乐素材时可选的分类列表。

## 调用指令 (Command)

```bash
node scripts/api.ts --payload='{"action":"music-category","account_id":"XXX"}'
```

## 参数列表 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | **是** | 蚁小二账号 ID (32位十六进制) |

## 返回结果 (Response)

返回一个包含音乐分类对象的数组。

```json
[
  {
    "yixiaoerId": "123",
    "yixiaoerName": "流行",
    "raw": {
      "id": "123",
      "name": "流行"
    }
  }
]
```

### 复杂对象：CategoryItem
- `yixiaoerId`: 内部分类 ID。
- `yixiaoerName`: 分类名称。
- `raw`: 原始平台返回的分类对象。如果在获取时该字段存在，发布表单中必须携带并完整透传。

## 脚本逻辑 (Backend)

- **脚本路径**: `scripts/api.ts`
- **功能**: 封装蚁小二标准音乐分类查询接口 (`GET /platform-accounts/{platformAccountId}/music/category`)。
- **参数映射**: 将 `account_id` 映射为路径变量。
