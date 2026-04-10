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

返回一个包含音乐分类对象的数组。脚本会自动将多级嵌套的分类**铺平 (Flatten)**，并为每个对象生成 `child` 路径数组。

```json
[
  {
    "yixiaoerId": "123",
    "yixiaoerName": "流行",
    "child": [ { "yixiaoerId": "123", "yixiaoerName": "流行" } ],
    "raw": { "id": "123", "name": "流行" }
  }
]
```

### 复杂对象：CategoryItem
- `yixiaoerId`: 内部分类 ID。
- `yixiaoerName`: 分类名称。
- `child`: **完整路径对象数组**。如果分类有父子级关系，发布表单时通常需要在此处填入整个生成的 `child`。
- `raw`: 原始平台返回的分类对象。

## 脚本逻辑 (Backend)

- **脚本路径**: `scripts/api.ts`
- **功能**: 封装蚁小二标准音乐分类查询接口并自动执行 `flattenTree` 逻辑。
