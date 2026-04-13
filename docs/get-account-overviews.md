# 📄 获取账号数据概览 Query 参数 (Get Account Overviews Query)

获取当前团队下各账号的综合表现概览。主要用于账号矩阵的权重分析、整体粉丝量监控及核心指标汇总。

> [!IMPORTANT]
> 执行本查询前，建议先通过 `accounts` 接口获取最新的账号列表，以确保查询维度的准确性。

## 1. 触发场景 (Trigger)

- **意图辨析**：当用户需要宏观查看某个平台下的“账号全景图”，而不是单个作品的数据时。例如：分析团队资产总量、对比部门/负责人之间的账号表现差异。
- **典型提示词**：
  - “汇总一下我旗下所有抖音号的粉丝量”
  - “我的视频号账号最近表现如何？”
  - “查看张三负责的账号的数据看板”

## 2. 交互协议 (Interactive Protocol)

1. **维度确定**：识别查询的目标平台（`platform` 为必填）。
2. **权限过滤**：若提到特定人员，解析并注入其 `memberIds`。
3. **分步确认**：在执行大规模数据汇总前，告知用户将查询的范围（如：“我将为您查询您名下 5 个抖音号的数据概览”）。
4. **看板交付**：提取 `overviewData` 中的总粉丝数、总播放量等关键指标，使用表格或摘要形式反馈给用户。

## 3. 参数定义 (Parameters)

| 字段名 | 类型 | 必填 | 默认值 | 描述 |
| :--- | :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | `account-overviews` | 固定值。 |
| `platform` | `string` | **是** | - | 指定查询某个平台。见 [平台定义](./platform.md)。 |
| `name` | `string` | 否 | - | 按照账号昵称模糊查询。 |
| `group` | `string` | 否 | - | 按照分组名称查询。 |
| `loginStatus` | `number` | 否 | - | 账号状态 (0未登录, 1成功, 2过期, 3失败)。 |
| `memberIds` | `string[]` | 否 | - | 负责人成员 ID 数组 (ObjectId)。 |
| `page` | `number` | 否 | `1` | 当前页码。 |
| `size` | `number` | 否 | `10` | 每页数量。 |

### 3.1 返回结果结构 (Response Details)

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `platformAccountId` | `string` | 媒体账号 ID |
| `platformAvatar` | `string` | 账号头像地址 |
| `platformAccountName` | `string` | 账号昵称 |
| `principalName` | `string` | 当前负责成员名 |
| `overviewData.fansTotal` | `number` | 总粉丝数 |
| `overviewData.playTotal` | `number` | 总播放量 |

## 4. 执行指令示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"account-overviews","platform":"抖音","page":1,"size":10}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **数据长时间未更新** | 平台风控或蚁小二客户端长时间未在线同步。 | 提醒用户确保蚁小二客户端在后台运行，并点击“同步数据”。 |
| **找不到特定分组账号** | 分组名输入有误（需完全匹配）。 | 建议先通过 `accounts` 接口获取所有分组名。 |
| **粉丝量显示为 0** | 新绑定的账号或数据抓取中。 | 告知用户初次绑定需等待 10-20 分钟进行初始化抓取。 |

---
> [!TIP]
> **总量聚合技巧**：账号数据 V2 提供的是总量的统计概览。如果需要查看增量趋势（如“昨天涨了多少粉”），应结合具体的数据趋势接口使用。

