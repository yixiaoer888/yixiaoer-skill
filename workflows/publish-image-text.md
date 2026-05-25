# 图文发布工作流

> 适用范围：抖音图文、小红书笔记、知乎想法、快手图文、微博图文、微信视频号图文  
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)

---

## 执行步骤（严格按顺序，不得跳步）

### Step 1：查询账号（BLOCKING）

```bash
yxer accounts [platform]
```

**执行规则：**
1. 查看返回的账号列表
2. **必须**确认目标账号 `status: 1`（在线）
3. 记录 `platformAccountId`（后续所有步骤必用）
4. 如果用户指定了账号昵称，模糊匹配找到对应账号
5. 如果有多个在线账号且用户未指定，**列出账号列表请用户选择**
6. 如无在线账号，终止流程并告知用户检查账号状态

---

### Step 2：上传图片（BLOCKING）

对**每一张**图片逐一执行：

```bash
yxer upload <文件路径或URL>
```

**返回示例：**
```json
{
  "success": true,
  "key": "img_abc123",
  "size": 512000,
  "width": 1080,
  "height": 1440,
  "format": "jpg"
}
```

**执行规则：**
1. 从返回中提取 `key` / `size` / `width` / `height` / `format`
2. 构建 `images` 数组（每个图片一个对象）
3. `cover` 和 `coverKey` **默认使用第一张图**的 key（用户未单独指定封面时）
4. 上传失败 → 检查路径/URL 是否有效，重试一次；仍失败则终止并报告错误

---

### Step 3：查询可选前置数据（按需，用户要求时才查）

**用户没提就不查，提了才查：**

```bash
# 需要挂载分类时
yxer categories <account_id> [--type video|article]

# 需要挂载地理位置时
yxer locations <account_id> [--query "关键词"]

# 需要挂载背景音乐时
yxer music <account_id> [--query "关键词"]

# 需要加入合集时
yxer collections <account_id> [--type video]

# 需要加热门话题时
yxer challenges <account_id> [--query "关键词"]

# 需要挂商品（带货）时
yxer goods <account_id> [--query "关键词"]
```

**执行规则：**
- 查询后，将**完整返回对象**（含 `raw`）填入 Payload 对应字段
- **严禁**只填 `yixiaoerId` 或 `yixiaoerName`，必须带入完整 `raw`

---

### Step 4：查阅平台参数文档

根据目标平台查阅对应文档，了解 `contentPublishForm` 的差异字段：

```bash
yxer schema <platform> image-text
```

| 平台 | 文档路径 |
|------|---------|
| 抖音图文 | [docs/publish/image-text/douyin.md](./docs/publish/image-text/douyin.md) |
| 小红书笔记 | [docs/publish/image-text/xiaohongshu.md](./docs/publish/image-text/xiaohongshu.md) |
| 知乎想法 | [docs/publish/image-text/zhihu.md](./docs/publish/image-text/zhihu.md) |
| 快手图文 | [docs/publish/image-text/kuaishou.md](./docs/publish/image-text/kuaishou.md) |
| 微博图文 | [docs/publish/image-text/xinlangweibo.md](./docs/publish/image-text/xinlangweibo.md) |

**先读通用索引**：[docs/publish/image-text/index.md](./docs/publish/image-text/index.md)

---

### Step 5：构造 Payload

将内容写入临时 JSON 文件（如 `./publish-payload.json`）：

```json
{
  "accountForms": [
    {
      "platformAccountId": "<Step 1 获取的 ID>",
      "images": [
        {
          "key": "<Step 2 获取>",
          "size": 512000,
          "width": 1080,
          "height": 1440,
          "format": "jpg"
        }
      ],
      "cover": {
        "key": "<第一张图的 key>",
        "size": 512000,
        "width": 1080,
        "height": 1440
      },
      "coverKey": "<第一张图的 key>",
      "contentPublishForm": {
        "formType": "task",
        "title": "<标题>",
        "description": "<正文，按平台文档格式化>",
        "images": [ "<同上，部分平台需要>" ],
        "<可选字段>": "<Step 3 查询结果（如有）>"
      }
    }
  ]
}
```

**注意：**
- `formType` 固定为 `"task"`，无需询问用户
- `description` 格式按各平台要求（抖音支持 HTML `<topic>`，知乎纯文本，小红书支持 emoji）
- 可选字段**只在用户要求或 Step 3 查询后才填入**，不要自作主张添加

---

### Step 6：本地校验（BLOCKING）

```bash
yxer validate <platform> image-text ./publish-payload.json
```

校验底层使用 `schemas/platforms/<platform>.imageText.schema.json` 的 JSON Schema，
对 `contentPublishForm` 字段进行 `additionalProperties: false` 严格校验。

- ✅ 校验通过 → 继续 Step 7
- ❌ 校验失败 → 读取错误信息，修正 Payload 后重新校验，不得跳过

---

### Step 7：确认并发布

**向用户展示发布摘要，确认后再执行：**

```
📋 发布确认
平台：抖音
账号：我的抖音号1（DY_ACC_001）
内容：3张图 + 标题"今日好心情"
描述：今天天气真好 🌞
定时：立即发布
```

用户确认后执行：

```bash
yxer publish image-text <platform> ./publish-payload.json
```

> 多平台发布时，Agent 必须为每个平台分别执行 Step 1 到 Step 7，分别构造 payload、校验并发布；不要使用逗号分隔平台合并成一次发布。

---

### Step 8：清理与报告

1. 删除临时 Payload 文件
2. 向用户报告发布结果（成功/失败 + 任务 ID）

---

## 完整执行示例

**用户输入：**
> "帮我发一条抖音图文，标题'今日好心情'，配上这3张图 C:\photos\1.jpg C:\photos\2.jpg C:\photos\3.jpg"

**Agent 执行流程：**

```
Step 1:
  yxer accounts Douyin
  → 找到账号 platformAccountId: "DY_ACC_001", status: 1 ✅

Step 2:
  yxer upload C:\photos\1.jpg
  → { "key": "img_aaa", "size": 512000, "width": 1080, "height": 1440, "format": "jpg" }
  yxer upload C:\photos\2.jpg
  → { "key": "img_bbb", ... }
  yxer upload C:\photos\3.jpg
  → { "key": "img_ccc", ... }

Step 3:
  用户未提 location/music/collection → 跳过

Step 4:
  查阅 docs/publish/image-text/douyin.md
  → 了解 description 支持 <topic> 标签

Step 5:
  构造 Payload → 写入 ./publish-payload.json

Step 6:
  yxer validate Douyin image-text ./publish-payload.json
  → ✅ Validation passed!

Step 7:
  展示摘要 → 用户确认 →
  yxer publish image-text Douyin ./publish-payload.json
  → { "success": true, "data": { "taskSetId": "xxx" }}

Step 8:
  删除 ./publish-payload.json
  报告：✅ 发布成功！任务 ID: xxx
```

---

## 带可选字段的完整示例

**用户输入：**
> "发抖音图文，标题'打卡上海'，配上图 C:\photos\sh.jpg，位置标记上海外滩"

**Agent 执行差异（Step 3 增加）：**

```
Step 3:
  yxer locations DY_ACC_001 --query "上海外滩"
  → 返回: { "yixiaoerId": "loc_001", "yixiaoerName": "外滩", "raw": { ... } }

Step 5:
  contentPublishForm 中增加:
    "location": {
      "yixiaoerId": "loc_001",
      "yixiaoerName": "外滩",
      "raw": { ... }   ← 完整 raw 对象，严禁只填 ID
    }
```

---

## 常见错误与排查

| 错误现象 | 原因 | 处理方式 |
|---------|------|---------|
| validate 报错 "Missing required field: platformAccountId" | Payload 结构错误，字段位置不对 | 检查 `accountForms[0].platformAccountId` 是否存在 |
| 图片上传成功但发布时报"资源不存在" | `key` 对应的 OSS 资源已过期 | 重新执行 `yxer upload` |
| 发布成功但位置/音乐未显示 | `raw` 对象未完整传入 | 检查 Payload 中对应字段是否包含完整 `raw` |
| validate 通过但发布报错 | 平台业务校验（如话题不存在） | 查阅对应平台文档的"常见问题"章节 |
| 多个账号在线但未指定 | Agent 自行选择了一个 | 应列出账号列表请用户选择 |
