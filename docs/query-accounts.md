# 查询账号列表 (Query Account List)

获取当前租户下绑定的自媒体平台账号列表及其 UID 信息。

## 场景描述 (Usage)

- "帮我列出我在这台电脑上绑定的所有抖音账号。"
- "我需要在发布前确认视频号的 UID。"

## 参数定义 (Parameters)

| 参数名 | 类型 | 描述 |
| :--- | :--- | :--- |
| `platform` | `string` | **可选**。如 `DouYin`, `XiaoHongShu` |
| `status` | `string` | **可选**。如 `active`, `expired` |

## 脚本映射 (Script)

- **脚本路径**: `../scripts/query-accounts.ts`
- **调用示例**: `node query-accounts.ts --platform=DouYin`

## 输出结果 (Output)

脚本返回标准的 JSON 数组对象。每一个账号对象应包含 `platform`, `name`, `uid`, `status`。
