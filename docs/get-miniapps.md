# 获取小程序列表 (Get MiniApps)

获取账号可用的挂载小程序列表。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "miniapps",
  "account_id": "YOUR_ACCOUNT_ID",
  "keyword": "搜索词"
}'
```

## 2. 返回数据结构

返回 `Category` 数组。
