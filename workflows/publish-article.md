# 文章发布工作流

> 适用范围：百家号文章、头条号文章、公众号文章、知乎文章、搜狐号文章等  
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)

---

## 与图文/视频发布的差异

| 步骤 | 图文发布 | 视频发布 | 文章发布 |
|------|---------|---------|---------|
| Step 2 | 上传图片（多张） | 上传视频 + 上传封面 | 上传封面 + 上传文章内嵌图片（如有） |
| Step 4 | 查阅 `docs/publish/image-text/` | 查阅 `docs/publish/video/` | 查阅 `docs/publish/article/` |
| Step 5 Payload | `images` 数组 + `description` | `video` 对象 + `description` | `content`（HTML正文）+ `title` + 封面 |
| Step 5 必填 | `images` + `description` | `video` + `cover` | `content`（正文HTML）+ `title` + `cover` |
| 分类 | 部分平台需要 | 部分平台需要 | **大多数平台必须要选到叶子节点** |

---

## 执行步骤（严格按顺序，不得跳步）

### Step 1：查询账号（BLOCKING）

```bash
yxer accounts [platform]
```

**执行规则：**（同图文/视频 Step 1）
1. 确认目标账号 `status: 1`（在线）
2. 记录 `platformAccountId`
3. 多个在线账号且用户未指定 → 列出请用户选择

---

### Step 2：上传封面 + 内嵌图片（BLOCKING）

#### 2.1 上传封面（必须）

```bash
yxer upload <封面图片路径或URL>
```

**执行规则：**
1. 文章**必须**有封面图，用户未提供时**必须询问**
2. 从返回中提取 `key` / `size` / `width` / `height` / `format`
3. 填入 `accountForms[0].cover` 和 `accountForms[0].coverKey`

#### 2.2 上传文章内嵌图片（如有）

文章内容（HTML 的 `<img>` 标签）需要引用已上传图片的 key。

```bash
yxer upload <图片路径或URL>
```

**执行规则：**
1. 每上传一张，记录返回的 `key`
2. 在构造 `content`（HTML正文）时，用 `![alt](key)` 或 `<img src="key">` 引用
3. **严禁**在 HTML 中直接写外部 URL

---

### Step 3：查询可选前置数据（按需，用户要求时才查）

```bash
# 需要挂载分类时（文章发布大多数平台需要！）
yxer categories <account_id> --type article

# 需要挂载地理位置时
yxer locations <account_id> [--query "关键词"]

# 需要加热门话题时
yxer challenges <account_id> [--query "关键词"] [--type article]
```

**⚠️ 分类特别提示（文章发布重点）：**

大部分文章平台（百家号、头条号、公众号等）**要求选择分类，且必须选到叶子节点**（最深层级）。

示例（百家号）：
```
❌ 错误：只填 "美食"（一级分类）
✅ 正确：填 "美食" → "家常菜"（二级分类，叶子节点）
```

查询后如有 `children`，**必须**继续引导用户选择子分类，直到无 `children` 为止。

---

### Step 4：查阅平台参数文档

根据目标平台查阅对应文档：

| 平台 | 文档路径 |
|------|---------|
| 百家号文章 | [docs/publish/article/baijiahao.md](./docs/publish/article/baijiahao.md) |
| 头条号文章 | [docs/publish/article/toutiaohao.md](./docs/publish/article/toutiaohao.md) |
| 公众号文章 | [docs/publish/article/weixin.md](./docs/publish/article/weixin.md) |
| 知乎文章 | [docs/publish/article/zhihu.md](./docs/publish/article/zhihu.md) |
| 搜狐号文章 | [docs/publish/article/sohu.md](./docs/publish/article/sohu.md) |

**先读通用索引**：[docs/publish/article/index.md](./docs/publish/article/index.md)

---

### Step 5：构造 Payload

将内容写入临时 JSON 文件（如 `./publish-payload.json`）：

```json
{
  "accountForms": [
    {
      "platformAccountId": "<Step 1 获取的 ID>",
      "cover": {
        "key": "<Step 2.1 获取的封面 key>",
        "size": 512000,
        "width": 1080,
        "height": 720
      },
      "coverKey": "<Step 2.1 获取的封面 key>",
      "contentPublishForm": {
        "formType": "task",
        "title": "<标题>",
        "content": "<正文 HTML，内嵌图片引用上传后的 key>",
        "thumbMediaId": "<Step 2.1 的封面 key>",
        "<可选字段>": "<Step 3 查询结果（如有）>"
      }
    }
  ]
}
```

**注意：**
- `formType` 固定为 `"task"`
- `content` 是**HTML 格式正文**，按各平台要求格式化（查阅 Step 4 的文档）
- 内嵌图片：`<img src="img_xxx">` 其中 `img_xxx` 是 Step 2.2 上传后返回的 `key`
- `thumbMediaId` 大多数平台要求与 `cover.key` 一致
- 分类字段（`category`）必须选到**叶子节点**，见 Step 3 说明

---

### Step 6：本地校验（BLOCKING）

```bash
yxer validate <platform> article ./publish-payload.json
```

校验底层使用 `schemas/platforms/<platform>.article.schema.json` 的 JSON Schema，
对 `contentPublishForm` 字段进行 `additionalProperties: false` 严格校验。

- ✅ 校验通过 → 继续 Step 7
- ❌ 校验失败 → 读取错误信息，修正 Payload 后重新校验

---

### Step 7：确认并发布

**向用户展示发布摘要：**

```
📋 发布确认
平台：百家号
账号：我的百家号（BJH_ACC_001）
标题：家常红烧肉的做法
封面：已上传
分类：美食 → 家常菜
定时：立即发布
```

用户确认后执行：

```bash
yxer publish article <platform> ./publish-payload.json
```

---

### Step 8：清理与报告

1. 删除临时 Payload 文件
2. 向用户报告发布结果（成功/失败 + 任务 ID）

---

## 完整执行示例

**用户输入：**
> "帮我发一篇百家号文章，标题'家常红烧肉的做法'，正文内容在 C:\articles\hongshaorou.txt，封面用 C:\covers\food.jpg"

**Agent 执行流程：**

```
Step 1:
  yxer accounts BaiJiaHao
  → 找到账号 platformAccountId: "BJH_ACC_001", status: 1 ✅

Step 2.1（上传封面）:
  yxer upload C:\covers\food.jpg
  → { "key": "img_cover001", "size": 512000, "width": 1080, "height": 720, "format": "jpg" }

Step 2.2（上传内嵌图片）:
  读取 C:\articles\hongshaorou.txt，发现正文中有3张图需要上传
  yxer upload C:\articles\img1.jpg → key: "img_inner001"
  yxer upload C:\articles\img2.jpg → key: "img_inner002"
  yxer upload C:\articles\img3.jpg → key: "img_inner003"

Step 3（查询分类）:
  yxer categories BJH_ACC_001 --type article
  → 返回分类列表，引导用户选择到叶子节点
  → 用户选择："美食" → "家常菜"
  → 获取完整分类对象（含 raw）

Step 4:
  查阅 docs/publish/article/baijiahao.md
  → 了解 content 格式要求（支持 HTML 标签）

Step 5:
  读取 C:\articles\hongshaorou.txt
  将正文中的图片引用替换为上传后的 key
  构造 Payload → 写入 ./publish-payload.json

Step 6:
  yxer validate BaiJiaHao article ./publish-payload.json
  → ✅ Validation passed!

Step 7:
  展示摘要 → 用户确认 →
  yxer publish article BaiJiaHao ./publish-payload.json
  → { "success": true, "data": { "taskSetId": "xxx" }}

Step 8:
  删除 ./publish-payload.json
  报告：✅ 发布成功！任务 ID: xxx
```

---

## 常见错误与排查

| 错误现象 | 原因 | 处理方式 |
|---------|------|---------|
| validate 报错 "Missing required field: content" | `contentPublishForm.content` 未填 | 检查 Payload 中是否有正文 HTML |
| 发布成功但封面不显示 | `coverKey` 与 `cover.key` 不一致 | 确保两个字段值完全相同 |
| 发布失败提示"分类错误" | 只填了一级分类，未选到叶子节点 | 重新 `yxer categories` 并选到最深层级 |
| HTML 中的图片无法显示 | 直接使用了外部 URL 或本地路径 | 必须通过 `yxer upload` 上传后引用 key |
| 文章格式错乱 | `content` 的 HTML 标签不符合平台要求 | 查阅对应平台文档的 HTML 支持标签列表 |

---

## 文章正文 HTML 格式说明

各平台对 `content` 字段的 HTML 标签支持程度不同，常见规则：

| 平台 | 支持标签 | 注意事项 |
|------|---------|---------|
| 百家号 | `<p>`, `<br>`, `<img>`, `<strong>`, `<em>` | 不支持 `<div>` 自定义样式 |
| 头条号 | `<p>`, `<br>`, `<img>`, `<h1>`-`<h6>` | 标题建议用 `<h2>`-`<h3>` |
| 公众号 | 完整 HTML | 支持 `<section>`, `<style>` 等 |
| 知乎 | Markdown 或简化 HTML | 建议用 Markdown 格式 |

**建议**：Agent 构造 `content` 时，优先使用各平台文档中明确声明的支持标签，避免使用未声明的标签。
