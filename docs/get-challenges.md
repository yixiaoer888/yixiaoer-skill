# 获取挑战/话题列表 (Get Challenges)

获取平台当前进行的挑战活动或推荐的话题列表。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "challenges",
  "account_id": "YOUR_ACCOUNT_ID",
  "keyword": "搜索词"
}'
```

## 2. 返回数据结构

返回 `Category` 数组。
