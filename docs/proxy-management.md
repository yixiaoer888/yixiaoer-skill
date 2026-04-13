# 📄 代理管理 动作 参数 (Proxy Management Actions)

获取及配置团队下的代理网络资源。在云发布 (Cloud Publish) 环境下，正确的代理配置是规避平台风控、确保发布成功率的核心要素。

> [!IMPORTANT]
> **云发布强制要求**：若使用“云发布”模式且遇到“账号代理不存在”报错，Agent **必须**通过本文档中的 `update-account` 动作协助用户绑定有效代理。

## 1. 触发场景 (Trigger)

- **意图辨析**：用户反馈发布由于 IP 问题失败、需要查看可用网络地区，或新绑定账号后需要配置出口 IP 以保护账号权重。
- **典型提示词**：
  - “我的抖音号发布报错了，是不是 IP 不对？”
  - “帮我查查有哪些上海的代理”
  - “把这个账号绑定到北京的内置代理上”

## 2. 交互协议 (Interactive Protocol)

1. **资源匹配协议**：
   - 若用户需要固定 IP，引导使用 `proxies` (独立代理)。
   - 若用户需要动态地区 IP，引导使用 `proxy-areas` (内置地区代理)。
2. **编码匹配原则**：在 `proxy-areas` 返回结果中：
   - **直辖市**（如北京、上海）：取第一层级的 `code`（如：`11`）。
   - **地级市**（如长沙）：必须进入 `cities` 数组获取对应的精确 `code`（如：`430100`）。
3. **两步执行逻辑**：先查询 (`proxies`/`proxy-areas`) -> 确认选择 -> 执行更新 (`update-account`)。

## 3. 参数定义 (Parameters)

### 3.1 查询代理列表 (`action: proxies`)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | 固定值为 `proxies`。 |

### 3.2 更新账号代理 (`action: update-account`)

| 字段名 | 类型 | 必填 | 描述 |
| :--- | :--- | :--- | :--- |
| **`action`** | `string` | **是** | 固定值为 `update-account`。 |
| **`account_id`** | `string` | **是** | 蚁小二账号 ID (`platformAccountId`)。 |
| `proxyId` | `string` | 否 | 独立代理 ID (从 `proxies` 获取)。 |
| `kuaidailiArea` | `string` | 否 | 内置地区编码 (从 `proxy-areas` 获取)。 |

## 4. 执行指令示例 (Command)

```bash
# 示例：将账号绑定到“上海”内置代理环境
node scripts/api.ts --payload='{"action":"update-account","account_id":"67fb2f1735eeb3cf31db3d65","kuaidailiArea":"31"}'
```

## 5. 常见问题排查 (Troubleshooting)

| 现象 | 可能原因 | 处理建议 |
| :--- | :--- | :--- |
| **账号代理不存在** | 云发布未配置代理，或原有 `proxyId` 已失效。 | 重新执行 `action: "proxy-areas"` 并通过 `update-account` 绑定。 |
| **设置后依然发布失败** | 代理负载过高或该地区 IP 被目标平台暂时封禁。 | 建议更换其他地区的内置代理尝试。 |
| **编码错误** | 使用了非法的 `kuaidailiArea` 字符串。 | 请严格核对 `proxy-areas` 接口返回的数字编码。 |

---
> [!TIP]
> **安全建议**：单一代理下挂载过多账号可能会触发平台“批量操作”预警。建议按业务线将账号均匀分布到不同的代理地区。
