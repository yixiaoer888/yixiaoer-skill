# 发布通用规则

> 所有发布工作流均需遵守本文档规则。

---

## 智能助理原则

做智能助理，不做表单填写机。**能自动补全的默认值先补全，只在必须决策时才追问用户。**

---

## 默认值自动补全规则

以下字段 Agent **应自动填入**，无需询问用户：

| 字段 | 默认值 | 说明 |
|------|--------|------|
| `formType` | `"task"` | 固定值，无需询问 |
| `publishChannel` | `"cloud"` | 默认云端发布 |
| `images[].key/size/width/height/format` | 从 `yxer upload` 结果自动获取 | 严禁手动编造 |
| `video.key/size/width/height/duration/format` | 从 `yxer upload` 结果自动获取 | 严禁手动编造 |
| `cover` / `coverKey` | 默认使用 `images[0]` 或视频封面 | 用户未单独指定封面时自动使用 |

以下字段 **应先向用户确认**，再填入：

| 字段 | 确认方式 |
|------|---------|
| `title` | 展示标题内容，请用户确认 |
| `description` / `content` | 如用户未提供，从标题自动生成并展示给用户确认 |
| `platformAccountId` | 如用户只有一个在线账号，自动选择并告知用户；多个则列出让用户选 |
| `scheduledTime` | 用户明确要求"定时发布"才询问时间；CLI 中统一传 13 位 Unix 毫秒时间戳，默认立即发布 |

---

## 资源上传规范

### 图片上传

- 支持格式：`jpg` / `png` / `webp`
- 每张图**必须**单独调用 `yxer upload` 获取 key
- 从返回结果中提取 `key` / `size` / `width` / `height` / `format`
- 禁止手动编造这些字段

### 视频上传

- 支持格式：`mp4` / `mov`
- 调用 `yxer upload <视频路径>` 获取 key
- 返回结果额外包含 `duration`（时长，单位秒）
- 视频封面**必须**单独上传，不能用视频文件本身代替封面

### URL 资源

- 直接传 HTTP/HTTPS URL，`yxer upload` 会自动下载后上传
- 本地文件传绝对路径

---

## 复杂对象查询规范

以下字段**严禁手动构造**，必须通过查询命令获取完整 `raw` 对象：

| 字段 | 查询命令 | 返回必需字段 |
|------|---------|-------------|
| `location` | `yxer locations <account_id> [--query 关键词]` | `yixiaoerId`, `yixiaoerName`, `raw` |
| `music` | `yxer music <account_id> [--query 关键词]` | `yixiaoerId`, `yixiaoerName`, `duration`, `playUrl`, `raw` |
| `collection` / `sub_collection` | `yxer collections <account_id> [--type video]` | `yixiaoerId`, `yixiaoerName`, `raw` |
| `challenge` | `yxer challenges <account_id> [--query 关键词]` | `yixiaoerId`, `yixiaoerName`, `raw` |
| `category` | `yxer categories <account_id> [--type video\|article]` | `yixiaoerId`, `yixiaoerName`, `raw`, `children` |

**关键**：查询后，将完整返回对象（含 `raw`）填入 Payload 对应字段，**不要只填 ID 或名称**。

---

## 分类层级规则（全平台通用）

若分类存在层级结构（如：动漫 → 国产动漫），Agent **必须**选择并提交**最深层的叶子节点**（二级或三级分类）。

**严禁**仅提交一级大类分类。

示例（百家号）：
```
❌ 错误：只填 "美食"（一级分类）
✅ 正确：填 "美食" 下的 "家常菜"（二级分类，叶子节点）
```

---

## 错误处理

| 错误场景 | 处理方式 |
|---------|---------|
| `yxer validate` 失败 | 读取错误信息，修正对应字段后重新校验 |
| `yxer publish` 失败 | 读取错误信息，判断是否需要重新 upload 或修正参数 |
| `yxer upload` 失败 | 检查文件路径/URL 是否有效，重试一次 |
| `yxer accounts` 无在线账号 | 告知用户，建议检查账号 cookie 是否过期 |
| 查询命令返回空 | 放宽关键词重试；仍为空则告知用户该账号不支持此功能 |

---

## 严禁行为（BLOCKING）

- ❌ 未确认账号 `status=1`（在线）就构造 Payload
- ❌ 手动编造 `key` / `size` / `width` / `height` / `duration`
- ❌ 手动构造 `location` / `music` / `collection` / `challenge` 的 `raw` 字段
- ❌ 跳过 `yxer validate` 直接执行 `yxer publish`
- ❌ 在 Payload 中直接使用外部 URL 作为图片/视频地址
- ❌ 跳过工作流步骤，自行拼 JSON Payload
- ❌ 用户说"草稿"时不询问类型，自行猜测是蚁小二草稿还是平台草稿
