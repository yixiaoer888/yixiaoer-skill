# 获取可挂载游戏列表 (Get Games)

获取发布视频时可挂载推广的游戏列表（主要用于抖音游戏推广）。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "games",
  "account_id": "YOUR_ACCOUNT_ID",
  "keyword": "搜索词"
}'
```

## 2. 返回数据结构

返回 `Category` 数组。
