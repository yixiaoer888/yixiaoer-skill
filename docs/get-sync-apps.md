# 获取同步发布应用列表 (Get Sync Apps)

获取发布当前作品时，可同时同步到的其他平台或应用列表（如抖音发布时同步到今日头条）。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "syncapps",
  "account_id": "YOUR_ACCOUNT_ID"
}'
```

## 2. 返回数据结构

返回 `Category` 数组。如果在获取时 `raw` 字段有值，发布表单中必须完整保留并透传。
