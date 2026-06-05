# 获取商品列表 (Get Goods)

此接口用于获取指定媒体账号在平台上关联的商品、团购或小店商品列表，以便在发布视频时进行商业推广。

## 1. 调用指令

```bash
yxer query goods YOUR_ACCOUNT_ID --query 可选搜索词 --json
```

## 2. 请求参数

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| action | string | 是 | 固定为 `goods` |
| account_id | string | 是 | 蚁小二系统内的媒体账号 ID |
| keyword | string | 否 | 搜索商品的关键词 |

## 3. 返回数据结构

返回一个包含 `ShoppingCartItem` 对象的数组及分页信息。发布时必须使用 `yxer query goods` 返回的完整对象，不能只保留 `data.yixiaoerId`、`data.yixiaoerName` 或局部字段。

### ShoppingCartItem 结构说明
| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `sale_title` | `string` | 挂车推广标题 |
| `images` | `string[]` | 顶层商品图片数组 |
| `data` | `object` | 核心商品数据对象 |

`data` 对象中的核心字段如下：

| 字段名 | 类型 | 说明 |
| :--- | :--- | :--- |
| `yixiaoerId` | `string` | (必填) 商品 ID |
| `yixiaoerName` | `string` | (必填) 商品名称 |
| `raw` | `object` | (必填) 平台原始商品对象。如果在获取时该字段存在，发布表单中必须携带并完整透传 |
| `yixiaoerDesc` | `string` | 商品规格说明 |
| `yixiaoerImageUrl` | `string` | 商品图片 URL |
| `price` | `number` | 商品价格（单位：分） |
| `earnPrice` | `number` | 预估佣金（单位：分） |
| `count` | `number` | 剩余库存 |

