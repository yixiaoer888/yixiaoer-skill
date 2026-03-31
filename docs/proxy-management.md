# 代理管理 (Proxy Management)

获取租户下的代理列表，并为特定账号配置代理。云发布时，如果账号没有绑定代理，通常需要先为其设置一个有效的代理以确保发布成功。

## 场景描述 (Usage)

- "查询我团队下所有的可用代理列表。"
- "为某个抖音账号绑定一个新的团队代理。"
- "将某个账号的代理切换为内置代理地区（如：上海）。"

## 1. 查询代理列表 (List Proxies)

获取当前团队下所有已配置的代理。

### 参数定义 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `size` | `number` | 否 | 每页数量，默认 `9999`。 |

### 调用示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"proxies"}'
```

### 返回结果说明 (Response Details)

返回一个包含代理对象的列表。每个代理对象包含：

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `id` | `string` | 代理 ID (用于 `update-account`) |
| `name` | `string` | 代理名称 |
| `proxyIp` | `string` | 代理 IP |
| `proxyPort` | `string` | 代理端口 |
| `enabled` | `boolean` | 是否启用 |
| `accounts` | `string[]` | 已绑定该代理的账号 ID 列表 |

---

## 2. 查询内置代理地区列表 (List Default Proxy Areas)

对于大多数没有团队代理的用户，建议使用系统内置代理。配置前需先查询支持的地区编码列表。

### 参数定义 (Payload Properties)

无。

### 调用示例 (Command)

```bash
node scripts/api.ts --payload='{"action":"proxy-areas"}'
```

### 返回结果说明 (Response Details)

返回一个地区对象数组。每一个对象代表一个可用的代理物理地区：

| 字段 | 类型 | 描述 |
| :--- | :--- | :--- |
| `name` | `string` | 地区显示名称 (如：上海、北京、广东等) |
| `code` | `string` | **地区唯一编码**。此值即为 `update-account` 接口中的 `kuaidailiArea` 参数。 |

---

## 3. 更新账号代理 (Update Account Proxy)

为指定账号设置代理信息。支持绑定**团队代理 (Team Proxy)** 或**内置代理 (Build-in Proxy)**。

### 参数定义 (Payload Properties)

| 字段名 | 类型 | 是否必填 | 描述 |
| :--- | :--- | :--- | :--- |
| `account_id` | `string` | 是 | 平台的账号 ID (可通过 `accounts` 接口获取) |
| `proxyId` | `string` | 否 | 团队代理的 ID (从 `proxies` 动作获取)。若要解绑，传 `null`。 |
| `kuaidailiArea` | `string` | 否 | 内置代理的地区编码 (从 `proxy-areas` 动作获取)。若要解绑，传 `null`。 |

> [!TIP]
> - 设置团队代理请使用 `proxyId`。
> - 设置内置代理请使用 `kuaidailiArea`。
> - 若两者同时传入，系统行为取决于后端实现，通常建议二选一。

### 调用示例 (Command)

```bash
# 为账号设置团队代理
node scripts/api.ts --payload='{"action":"update-account", "account_id":"xxx", "proxyId":"proxy-id-here"}'

# 为账号设置内置代理（上海）
node scripts/api.ts --payload='{"action":"update-account", "account_id":"xxx", "kuaidailiArea":"shanghai"}'
```

## 注意事项
- 请确保环境变量 `YIXIAOER_API_KEY` 已设置。
- 云发布环境下，未设置代理的账号可能会导致任务执行失败。
