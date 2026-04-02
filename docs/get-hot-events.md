# 获取热点列表 (Get Hot Events)

获取平台实时的搜索热点、活动热贴，用于绑定到视频以提升曝光。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "hot-events",
  "account_id": "YOUR_ACCOUNT_ID"
}'
```

## 2. 返回数据结构

返回 `Category` 数组。
