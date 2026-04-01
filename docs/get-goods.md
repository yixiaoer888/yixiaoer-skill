# 获取商品列表 (Get Goods)

此接口用于获取指定媒体账号在平台上关联的商品、团购或小店商品列表，以便在发布视频时进行商业推广。

## 1. 调用指令

```bash
node scripts/api.ts --payload='{
  "action": "goods",
  "account_id": "YOUR_ACCOUNT_ID",
  "keyword": "可选搜索词"
}'
```

## 2. 请求参数

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| action | string | 是 | 固定为 `goods` |
| account_id | string | 是 | 蚁小二系统内的媒体账号 ID |
| keyword | string | 否 | 搜索商品的关键词 |

## 3. 返回数据结构

返回一个包含 `ShoppingCartItem` 对象的数组及分页信息。

### ShoppingCartItem 结构说明
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| yixiaoerId | string | 商品 ID |
| yixiaoerName | string | 商品名称 |
| yixiaoerDesc | string | 商品规格说明 |
| yixiaoerImageUrl | string | 商品图片 URL |
| price | number | 商品价格（单位：分） |
| earnPrice | number | 预估佣金（单位：分） |
| count | number | 剩余库存 |
