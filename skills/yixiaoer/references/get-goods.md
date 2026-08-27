# 获取商品列表 (Get Goods)

此接口用于获取指定媒体账号在平台上关联的商品、团购或小店商品列表，以便在发布视频时进行商业推广。

若需要从一个商品链接生成可挂购物车的完整商品对象，使用 `goods-detail`。它解析单个链接，不查询账号的选品车或商品列表。

## 1. 调用指令

```bash
yxer query goods YOUR_ACCOUNT_ID --query 可选搜索词 --json
yxer query goods-detail YOUR_ACCOUNT_ID --url "商品链接" --json
```

## 2. 请求参数

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| action | string | 是 | 固定为 `goods` |
| account_id | string | 是 | 蚁小二系统内的媒体账号 ID |
| keyword | string | 否 | 搜索商品的关键词 |

`goods-detail` 参数：

| 参数名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| account_id | string | 是 | 蚁小二系统内的媒体账号 ID |
| url | string | 是 | 待解析的商品链接 |

## 3. 返回数据结构

返回一个包含 `ShoppingCartItem` 对象的数组及分页信息。对于使用 `ShoppingCartItem` 结构的平台，发布时必须使用 `yxer query goods` 返回的完整对象，不能只保留 `data.yixiaoerId`、`data.yixiaoerName` 或局部字段。

`yxer query goods-detail` 返回相同的完整商品对象结构，可将其结果中的商品对象用于抖音 `shopping_cart`；该命令只解析链接，不执行发布。

多多视频不使用 `ShoppingCartItem` 查询对象。多多视频的 `shopping_cart` 是平台业务对象，`goods_id` 必须由用户手工输入，CLI 固定补充 `source: "pdd"`；不要把 `query goods` 返回的 `yixiaoerId` 当作多多视频的 `goods_id`。

对于支持该接口的平台，发布前可执行 `yxer query entitlements YOUR_ACCOUNT_ID --json` 确认返回的 `shopping_cart` 为 `true`。CLI 在 `shopping_cart` 非空时会在 `validate`、`publish --dry-run` 与正式 `publish` 前执行同一权限检查；多多视频不支持该通用权限接口，因此会跳过该接口，仅依赖 payload/schema 校验及平台侧发布校验。

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

