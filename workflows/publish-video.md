# 视频发布工作流

> 适用范围：抖音视频、快手视频、B站视频、视频号视频、微博视频等  
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)

---

## 与图文发布的差异

| 步骤 | 图文发布 | 视频发布 |
|------|---------|---------|
| Step 2 | 上传图片（多张） | 上传视频（单个）+ 上传封面（单独） |
| Step 2 返回 | `key`/`size`/`width`/`height` | 额外包含 `duration`（时长，单位秒） |
| Step 4 | 查阅 `docs/publish/image-text/` | 查阅 `docs/publish/video/` |
| Step 5 Payload | `images` 数组 | `video` 对象 + `cover` 必须单独上传 |

---

## 执行步骤（严格按顺序，不得跳步）

### Step 1：查询账号（BLOCKING）

```bash
yxer accounts [platform]
```

**执行规则：**（同图文发布 Step 1）
1. 确认目标账号 `status: 1`（在线）
2. 记录 `platformAccountId`
3. 多个在线账号且用户未指定 → 列出请用户选择

---

### Step 2：上传视频 + 上传封面（BLOCKING）

#### 2.1 上传视频

```bash
yxer upload <视频文件路径或URL>
```

**返回示例：**
```json
{
  "success": true,
  "key": "video_abc123",
  "size": 52428800,
  "width": 1080,
  "height": 1920,
  "duration": 60,
  "format": "mp4"
}
```

**执行规则：**
1. 从返回中提取 `key` / `size` / `width` / `height` / `duration` / `format`
2. 视频**只能有一个**，`accountForms[0].video` 填这个对象
3. 上传失败 → 检查文件格式/URL 是否有效，重试一次

#### 2.2 上传封面（必须单独上传）

```bash
yxer upload <封面图片路径或URL>
```

**执行规则：**
1. 封面**不能**用视频文件本身代替
2. 从返回中获取封面的 `key` / `size` / `width` / `height`
3. 填入 `accountForms[0].cover` 和 `accountForms[0].coverKey`
4. 用户未提供封面时：**告知用户需要封面图**，不要自行用视频帧代替

---

### Step 3：查询可选前置数据（按需，用户要求时才查）

同图文发布 Step 3，命令完全一致：

```bash
# 需要挂载分类时
yxer categories <account_id> --type video

# 需要挂载地理位置时
yxer locations <account_id> [--query "关键词"]

# 需要挂载背景音乐时
yxer music <account_id> [--query "关键词"]

# 需要加入合集时
yxer collections <account_id> --type video

# 需要加热门话题/挑战时
yxer challenges <account_id> [--query "关键词"] [--type video]

# 需要挂商品（带货）时
yxer goods <account_id> [--query "关键词"]
```

**执行规则：** 同图文发布 Step 3

---

### Step 4：查阅平台参数文档

根据目标平台查阅对应文档：

| 平台 | 文档路径 |
|------|---------|
| 抖音视频 | [docs/publish/video/douyin.md](./docs/publish/video/douyin.md) |
| 快手视频 | [docs/publish/video/kuaishou.md](./docs/publish/video/kuaishou.md) |
| B站视频 | [docs/publish/video/bilibili.md](./docs/publish/video/bilibili.md) |
| 微博视频 | [docs/publish/video/xinlangweibo.md](./docs/publish/video/xinlangweibo.md) |
| 视频号 | [docs/publish/video/weixinvideo.md](./docs/publish/video/weixinvideo.md) |

**先读通用索引**：[docs/publish/video/index.md](./docs/publish/video/index.md)

---

### Step 5：构造 Payload

将内容写入临时 JSON 文件（如 `./publish-payload.json`）：

```json
{
  "accountForms": [
    {
      "platformAccountId": "<Step 1 获取的 ID>",
      "video": {
        "key": "<Step 2.1 获取的 key>",
        "size": 52428800,
        "width": 1080,
        "height": 1920,
        "duration": 60,
        "format": "mp4"
      },
      "cover": {
        "key": "<Step 2.2 获取的封面 key>",
        "size": 512000,
        "width": 1080,
        "height": 1920
      },
      "coverKey": "<Step 2.2 获取的封面 key>",
      "contentPublishForm": {
        "formType": "task",
        "title": "<标题>",
        "description": "<描述，按平台文档格式化>",
        "<可选字段>": "<Step 3 查询结果（如有）>"
      }
    }
  ]
}
```

**注意：**
- `formType` 固定为 `"task"`
- `video` 对象是**单个对象**，不是数组
- `cover` **必须**单独上传，不能用 `video.key` 代替
- 可选字段只在用户要求或 Step 3 查询后才填入

---

### Step 6：本地校验（BLOCKING）

```bash
yxer validate <platform> video ./publish-payload.json
```

校验底层使用 `schemas/platforms/<platform>.video.schema.json` 的 JSON Schema，
对 `contentPublishForm` 字段进行 `additionalProperties: false` 严格校验。

- ✅ 校验通过 → 继续 Step 7
- ❌ 校验失败 → 读取错误信息，修正 Payload 后重新校验

---

### Step 7：确认并发布

**向用户展示发布摘要：**

```
📋 发布确认
平台：抖音
账号：我的抖音号1（DY_ACC_001）
视频：60秒，1080x1920，mp4
封面：已上传
标题：今日好心情
定时：立即发布
```

用户确认后执行：

```bash
yxer publish video <platform> ./publish-payload.json
```

---

### Step 8：清理与报告

1. 删除临时 Payload 文件
2. 向用户报告发布结果（成功/失败 + 任务 ID）

---

## 完整执行示例

**用户输入：**
> "帮我发一条抖音视频，标题'今日分享'，视频在 C:\videos\myvideo.mp4，封面用 C:\covers\cover.jpg"

**Agent 执行流程：**

```
Step 1:
  yxer accounts Douyin
  → 找到账号 platformAccountId: "DY_ACC_001", status: 1 ✅

Step 2.1（上传视频）:
  yxer upload C:\videos\myvideo.mp4
  → {
      "key": "video_aaa",
      "size": 52428800,
      "width": 1080,
      "height": 1920,
      "duration": 60,
      "format": "mp4"
    }

Step 2.2（上传封面）:
  yxer upload C:\covers\cover.jpg
  → {
      "key": "img_cover001",
      "size": 512000,
      "width": 1080,
      "height": 1920,
      "format": "jpg"
    }

Step 3:
  用户未提 location/music/collection → 跳过

Step 4:
  查阅 docs/publish/video/douyin.md
  → 了解 description 格式要求

Step 5:
  构造 Payload → 写入 ./publish-payload.json

Step 6:
  yxer validate Douyin video ./publish-payload.json
  → ✅ Validation passed!

Step 7:
  展示摘要 → 用户确认 →
  yxer publish video Douyin ./publish-payload.json
  → { "success": true, "data": { "taskSetId": "xxx" }}

Step 8:
  删除 ./publish-payload.json
  报告：✅ 发布成功！任务 ID: xxx
```

---

## 常见错误与排查

| 错误现象 | 原因 | 处理方式 |
|---------|------|---------|
| validate 报错 "Missing required field: video" | Payload 结构错误，`video` 字段位置不对 | 检查 `accountForms[0].video` 是否存在 |
| 发布成功但封面不显示 | `coverKey` 与 `cover.key` 不一致 | 确保两个字段值完全相同 |
| 视频上传成功但发布时报"格式不支持" | `format` 值不正确 | 检查视频文件实际格式，填入正确后缀 |
| 发布成功但位置/音乐未显示 | `raw` 对象未完整传入 | 检查 Payload 中对应字段是否包含完整 `raw` |
| 视频时长超过平台限制 | 平台对视频时长有限制（通常15s-60s） | 查阅对应平台文档的时长限制 |
