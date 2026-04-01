# 获取群聊列表 (Get Groups)

此接口用于获取指定媒体账号在平台上已创建或加入的群聊列表，以便在发布视频时进行绑定。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "groups",
  "account_id": "YOUR_ACCOUNT_ID"
}'
```

## 2. 请求参数

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| action | string | 是 | 固定为 `groups` |
| account_id | string | 是 | 蚁小二系统内的媒体账号 ID |

## 3. 返回数据结构

返回一个包含 `Group` 对象的数组。

### Group 结构说明
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| yixiaoerId | string | 群聊 ID |
| yixiaoerName | string | 群聊标题 |
| yixiaoerDesc | string | 群聊描述 |
| yixiaoerImageUrl | string | 群聊头像 URL |
