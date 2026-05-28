# 发布通用规则

> 所有发布工作流均需遵守本文档规则。

---

## 智能助理原则

做智能助理，不做表单填写机。能自动补全的默认值先补全，只在必须决策时才追问用户。

---

## 发布通道判断规则

Agent 在任何 `publish` 之前，都要先判断这次任务是云发布还是本机发布。

### 何时用云发布

- 用户未明确指定发布通道
- 用户只说“帮我发布”，没有强调本机客户端
- 当前环境没有可用 `clientId`

### 何时用本机发布

- 用户明确说“本机发布”“本地发布”“走客户端发布”
- 用户明确表示不要走云端代理，或者希望走当前机器网络环境
- 云发布已因代理问题失败，用户接受改走本机

### 本机发布执行规则

- 必须显式使用 `publishChannel=local`
- 必须提供 `clientId`
- `clientId` 获取优先级：
  1. payload 中已有 `clientId`
  2. 命令 flags：`--client-id <clientId>`
  3. 第四个位置参数：`yxer publish <type> <platform> <payload.json> <clientId>`
  4. 本地配置：`yxer config set-local-client-id <clientId>` 后由 CLI 自动读取

### 发布通道失败后的回退

- 本机发布报“客户端不在线”或“获取在线设备列表失败”：
  - 提示用户启动并登录蚁小二客户端
  - 若用户不方便保持在线，建议改用云发布
- 云发布报“账号代理不存在”：
  - 提示检查账号代理配置
  - 若用户希望立即绕过代理限制，可改用本机发布

---

## 默认值自动补全规则

以下字段 Agent 应自动填入，无需询问用户：

| 字段 | 默认值 | 说明 |
| --- | --- | --- |
| `formType` | `"task"` | 固定值，无需询问 |
| `publishChannel` | `"cloud"` | 仅在用户未指定本机发布时默认使用 |
| `images[].key/size/width/height/format` | 从 `yxer upload` 结果自动获取 | 严禁手动编造 |
| `video.key/size/width/height/duration/format` | 从 `yxer upload` 结果自动获取 | 严禁手动编造 |
| `cover` / `coverKey` | 默认使用 `images[0]` 或视频封面 | 用户未单独指定封面时自动使用 |

以下字段应先向用户确认，再填入：

| 字段 | 确认方式 |
| --- | --- |
| `title` | 展示标题内容，请用户确认 |
| `description` / `content` | 如用户未提供，从标题自动生成并展示给用户确认 |
| `platformAccountId` | 如用户只有一个在线账号，自动选择并告知用户；多个则列出让用户选 |
| `scheduledTime` | 用户明确要求定时发布才询问时间；CLI 中统一传 13 位 Unix 毫秒时间戳，默认立即发布 |
| `publishChannel` | 用户明确提到本机/本地/客户端发布时，必须切换为 `local` |
| `clientId` | 用户要求本机发布且当前配置中没有默认值时，必须补齐 |

---

## 资源上传规范

### 图片上传

- 支持格式：`jpg` / `png` / `webp`
- 每张图必须单独调用 `yxer upload` 获取 key
- 从返回结果中提取 `key` / `size` / `width` / `height` / `format`
- 禁止手动编造这些字段

### 视频上传

- 支持格式：`mp4` / `mov`
- 调用 `yxer upload <视频路径>` 获取 key
- 返回结果额外包含 `duration`
- 视频封面必须单独上传，不能用视频文件本身代替封面

### URL 资源

- 直接传 HTTP/HTTPS URL，`yxer upload` 会自动下载后上传
- 本地文件传绝对路径

---

## 复杂对象查询规范

以下字段严禁手动构造，必须通过查询命令获取完整 `raw` 对象：

| 字段 | 查询命令 | 返回必需字段 |
| --- | --- | --- |
| `location` | `yxer locations <account_id> [--query 关键词]` | `yixiaoerId`, `yixiaoerName`, `raw` |
| `music` | `yxer music <account_id> [--query 关键词]` | `yixiaoerId`, `yixiaoerName`, `duration`, `playUrl`, `raw` |
| `collection` / `sub_collection` | `yxer collections <account_id> [--type video]` | `yixiaoerId`, `yixiaoerName`, `raw` |
| `challenge` | `yxer challenges <account_id> [--query 关键词]` | `yixiaoerId`, `yixiaoerName`, `raw` |
| `category` | `yxer categories <account_id> [--type video\|article]` | `yixiaoerId`, `yixiaoerName`, `raw`, `children` |

查询后，将完整返回对象填入 payload 对应字段，不要只填 ID 或名称。

---

## 分类层级规则

若分类存在层级结构，Agent 必须选择并提交最深层的叶子节点。

示例：

```text
错误：只填 "美食"
正确：填 "美食" 下的 "家常菜"
```

---

## 错误处理

| 错误场景 | 处理方式 |
| --- | --- |
| `yxer validate` 失败 | 读取错误信息，修正对应字段后重新校验 |
| `yxer publish` 失败 | 读取错误信息，判断是否需要重新 upload 或修正参数 |
| `yxer upload` 失败 | 检查文件路径或 URL 是否有效，重试一次 |
| `yxer accounts` 无在线账号 | 告知用户，建议检查账号 cookie 是否过期 |
| 查询命令返回空 | 放宽关键词重试；仍为空则告知用户该账号不支持此功能 |
| 本机发布失败且提示客户端不在线 | 引导用户启动客户端，或建议改用云发布 |
| 云发布失败且提示账号代理不存在 | 提示检查代理，或建议改用本机发布 |

---

## 严禁行为

- 未确认账号 `status=1` 就构造 payload
- 用户明确要求本机发布时，仍然默认走云发布
- 使用本机发布却没有显式提供或确认 `clientId`
- 手动编造 `key` / `size` / `width` / `height` / `duration`
- 手动构造 `location` / `music` / `collection` / `challenge` 的 `raw`
- 跳过 `yxer validate` 直接执行 `yxer publish`
- 在 payload 中直接使用外部 URL 作为图片或视频地址
- 跳过工作流步骤，自行拼大 JSON payload
- 用户说“草稿”时不询问类型，自行猜测是蚁小二草稿还是平台草稿
