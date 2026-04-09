# YiXiaoEr Skill 严格执行标准 (Strict Execution Standard)

为了确保 蚁小二 (YiXiaoEr) 自动化技能在复杂 Agent 环境下的高可靠性、可追溯性以及版本一致性，所有调用者及技能维护者必须遵循本标准。

## 1. 版本控制与一致性 (Versioning & Consistency)

### 1.1 语义化版本
- 技能必须在 `SKILL.md` 的 frontmatter 中明确标注 `version`（遵循 SemVer 规范）。
- **执行要求**：Agent 在启动时**必须**校验环境中的技能版本是否符合预期（如 `>= 1.6.0`）。

### 1.2 物理一致性 (Integrity)
- 技能的核心逻辑 (scripts/api.ts) 不得在未修改版本号的情况下随意变动。
- 自动化流水线应建议生成 Hash 值进行对比。

## 2. 错误分类与标准码 (Standardized Error Codes)

为了准确区分“Agent 使用错误”与“后端系统错误”，所有输出的错误 JSON 必须包含 `errorCode` 字段：

| 错误代码 (errorCode) | 分类 | 触发场景 |
| :--- | :--- | :--- |
| `YIXIAOER_USAGE_ERR` | Agent 侧错误 | 参数缺失、JSON 格式错误、平台枚举值非法、资源未上传即发布。 |
| `YIXIAOER_REMOTE_ERR` | 远端 API 错误 | 蚁小二后端返回 4xx/5xx、API Key 失效、配额超限。 |
| `YIXIAOER_AUTH_ERR` | 鉴权错误 | 缺失 `YIXIAOER_API_KEY` 或 Key 无权访问特定接口。 |
| `YIXIAOER_ENV_ERR` | 环境错误 | Node 版本不足、关键依赖缺失、网络不通。 |
| `YIXIAOER_INTEGRITY_ERR` | 完整性错误 | 运行的代码与文档版本不一致。 |

## 3. 资源预处理标准 (Pre-processing Standard)

**核心约束：发布与上传必须物理分离**。

1. **先上传再引用**：严禁在 `publish` Payload 中直接使用外部原始 URL。必须通过 `upload` action 获取系统内部 `key` 后进行替换。
2. **ContentType 签名一致性**：获取预签名 URL 时声明的 `contentType` 必须与执行 PUT 上传时的 `Content-Type` Header **完全一致**。禁止自行猜测或随意留空，否则将导致 `SignatureDoesNotMatch` 错误。
3. **封面图 (Cover)**：视频内容的封面图必须同样经过 `upload` 流程。
4. **状态校验**：在上传统一资源后，建议通过 `material` 接口（若需要）将其同步至素材库以增加发布稳定性。

## 4. 调用序列与检索规范 (Interaction Flow)

Agent 在执行任务时必须遵循以下“分级检索”逻辑：

1. **L1 根检索**：读取 `SKILL.md` 确定 `action` 类型。
2. **L2 索引检索**：在执行 `publish` 前，必须先读取对应内容类型的 `index.md`（如 `docs/publish/article/index.md`）确定 DTO 架构。
3. **L3 平台细化**：根据目标平台（如“抖音”），读取 `docs/publish/video/douyin.md` 对 `platformSettings` 进行深度补全。
4. **禁止跳级**：严禁在未读取 `index.md` 的情况下直接拼装平台参数。

## 5. 输出格式规范 (Output Schema)

脚本输出必须为严格的 JSON。

### 成功响应示例
```json
{
  "success": true,
  "action": "publish",
  "version": "1.6.2",
  "data": {
    "task_set_id": "TS12345678"
  }
}
```

## 6. 错误处理与自检流程 (Error Handling & Self-Diagnosis)

当脚本执行失败并返回 `success: false` 时，Agent 或开发者**必须**按照以下顺序进行“三步自检”：

### 第一步：检查版本一致性 (Check Version)
- 检查当前运行环境中的 `SKILL.md` 版本号。
- 确认 Agent 使用的 API 调用逻辑是否与该版本匹配。
- **排查项**：是否因版本过低导致某些新增字段（如 `publishChannel`）未生效。

### 第二步：检查请求参数规范 (Check Parameters)
- 对比输出的 `details` 信息与 `docs/` 下的 DTO 定义。
- **排查项**：检查必填字段是否遗漏、数据类型（String/Number）是否符合文档、平台枚举名称是否正确。

### 第三步：检查缓存与过期数据 (Check Cache & Stale Data)
- 检查是否存在过期的缓存文件（如旧的 Payload JSON、过时的资源 Key）。
- **资源上传专项排查**：在上传阶段如果出现 `SignatureDoesNotMatch` 错误，优先检查 `contentType` 定义是否与上传时的 Header **完全对齐**。
- **执行逻辑**：如果怀疑使用了缓存导致的异常，必须清理相关临时文件并严格按照最新技能定义的流程重新执行（如重新执行 `upload` 流程）。

---

---
*本标准由 OpenClaw 架构组制定。违反本标准导致的发布失败或数据异常将由调用者承担相应责任。*
