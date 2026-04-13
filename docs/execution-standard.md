# 📄 YiXiaoEr Skill 严格执行标准 (Core Execution Standard)

> [!IMPORTANT]
> 为了确保蚁小二 (YiXiaoEr) 在复杂 Agent 环境下的高可靠性、安全性和极致体验，所有 Skill 开发与 Agent 执行必须严格遵循本标准。

---

## 1. 核心架构设计 (Architecture Philosophy)

本技能生态采用 **SSRN (Skill-Reference-Script-Notify)** 闭环架构：

*   **Skill (大脑)**: `SKILL.md` 是全局调度中心，定义能力边界与加载策略。
*   **Reference (知识)**: `docs/` 下的结构化文档，按“索引 -> 详情”层级组织，严禁 Agent 脱离文档进行“常识假设”。
*   **Script (双手)**: `scripts/` 下的工具，遵循“无依赖、安全透明、自诊断”原则。
*   **Notify (反馈)**: 统一的错误分类码与诊断建议。

---

## 2. 全局交互规范 (Global Interaction Protocol)

> ‼️ **所有 Agent 在执行任务时必须严格遵守，优先级高于局部规则。**

### 2.1 分步确认协议 (Step-by-Step Consent)
任何非简单问答的操作（发布、上传、修改）必须遵循：
1.  **明确需求**：准确识别用户意图（如：发视频还是发图文？云发布还是本机发布？）。
2.  **前置校验**：在收集参数前，必须先确认账号状态（`accounts` 接口）及环境权限。
3.  **征得同意**：告知下一步动作（如：“我将为您查询抖音的分类并上传视频”），**等用户明确同意后**才继续。
4.  **执行前预览**：构造好 Payload 后，列表摘要展示关键参数，确认后再调用脚本。

### 2.2 “只检索，不生成” (Strict Reference)
*   **严禁幻觉**：Agent 严禁凭记忆构造 `platformSettings` 或 DTO。
*   **物理溯源**：所有填写的 JSON 结构必须有对应的 `docs/` 文档来源。
*   **按需加载**：遵循 `Index -> Platform Specific` 的加载策略，不要一次性加载所有平台文档。

### 2.3 零假设原则 (Zero Assumptions)
*   对缺失参数必须逐一追问。
*   禁止推断默认值（除非文档明确标注 `Default: XXX`）。
*   若语义歧义（如“存草稿”指蚁小二还是平台？），必须询问确认。

---

## 3. 文档化标准 (Documentation Standards)

所有技能文档（Markdown）必须遵循以下视觉与结构规范：

### 3.1 命名与标题标准
*   文件标题格式：`📄 [功能名] [文档类型] 参数/索引`。
*   使用 Emoji 增强可读性：‼️ (禁令)、📄 (文档)、✅ (成功)、⚠️ (风险)、🔧 (工具)。

### 3.2 结构五大板块
1.  **Trigger (触发器)**: 场景说明。
2.  **Interactive Protocol (交互协议)**: 该场景特有的确认逻辑。
3.  **Parameters (参数表)**: 严格区分“必填”与“非必填”，包含类型与示例值。
4.  **Command (示例命令)**: 具体的 JSON Payload 或 CLI 示例。
5.  **Troubleshooting (排障)**: 该接口的常见错误 `errorCode` 及修复方案。

---

## 4. 脚本开发标准 (Secure Action Standard)

为了保证脚本在各种开发者环境及 Agent 托管环境下的普适性：

1.  **零依赖 (Zero-Dependency)**: 鼓励使用 Nodejs 原生 `https`、`fs` 或 Python 原生 `urllib`。减少 `npm install`。
2.  **安全透明 (Security)**: 
    *   敏感信息（API Key）禁止硬编码，通过变量传递。
    *   涉及签名的脚本，应支持外部传入签名值而非在内部存储私钥。
3.  **自诊断 (Self-Diagnosis)**: 脚本输出必须包含解析建议。当 API 报错时，脚本应将 `errorCode` 翻译为人类可读的改进动作。

---

## 5. 错误分类与诊断流程

| 错误代码 (errorCode) | 含义 | 处理策略 |
| :--- | :--- | :--- |
| `YIXIAOER_USAGE_ERR` | Agent 参数或 JSON 格式错误 | Agent 重新读文档核对字段必填性。 |
| `YIXIAOER_REMOTE_ERR` | 远端 API 或后端逻辑错误 | 按 [🛡️ 避坑手册](./troubleshooting-guide.md) 排查。 |
| `YIXIAOER_AUTH_ERR` | 鉴权或账号状态失效 | 引导用户通过 `accounts` 更新状态或检查 Key。 |

---

## 6. 加载策略系统 (Loading Strategy)

Agent 必须使用以下策略动态管理上下文，防止上下文污染：

*   **策略 A (发布场景)**: `SKILL.md` (识别 Action) -> `index.md` (识别 DTO) -> `Platform.md` (细化表单)。
*   **策略 B (查询场景)**: `SKILL.md` -> `Target Document`。
*   **策略 C (排障场景)**: `SKILL.md` -> `errorCode Match` -> `Troubleshooting Guide`。

---

*本标准遵循 YiXiaoEr Skill Open Standard v1.8+。违反本协议导致的执行失败将被记录为交互缺陷。*
致。
- **发布通道与代理校验**：
    - 若执行发布任务返回“获取在线设备列表失败”或提示“客户端不在线”，原因为选择了 **本机发布 (publishChannel: 'local')** 但环境中未运行蚁小二客户端。建议保持客户端在线，或修改为 **云发布 (publishChannel: 'cloud')**。
    - 若**云发布**返回“账号代理不存在”，原因为账号绑定的独立代理 (proxyId) 失效或该平台强制需要代理配置。Agent 应先检查账号的 `proxyId` 或 `kuaidailiArea` 是否设置，并提示用户通过 `update-account` 绑定新的代理。
- **执行逻辑**：如果怀疑使用了缓存导致的异常，必须清理相关临时文件并严格按照最新技能定义的流程重新执行（如重新执行 `upload` 流程）。
- **草稿处理逻辑 (Draft Selection)**：当涉及“草稿”时，Agent 必须通过语义判断用户是指 **“蚁小二内部草稿”** (action: save-draft) 还是 **“发布到平台草稿箱”** (action: publish + pubType: 0)。**若语义不明确，Agent 必须询问用户以确认，严禁在未明确用户意图时随意选取一种或默认执行。**
    - **平台草稿兼容性补丁**：若用户指定“发布到平台草稿箱”，Agent 遵循以下执行链：
        1. **pubType 优先**：若平台定义了 `pubType`，使用 `pubType: 0`。
        2. **visibleType 次之**：若无 `pubType` 但有 `visibleType` (或 `status`/`privacy`)，将其设为 **`1` (私密)**。对应的，`0` 表示公开。
        3. **不支持告知**：若均无定义，则判定为不支持，必须提示用户。
### 6.4 分类路径完整性校验 (Category Path Integrity)
- **核心逻辑**：对于支持多级结构的分类（如 B 站视频分区），Agent 在填充表单时**必须**使用完整的 `child` 路径数组。
- **数据获取**：直接透传从 `action: "categories"` 结果中提取的 `child` 字段。
- **执行动作**：Agent 选中目标分类后，应直接将该对象的整个 `child` 数组填入 `contentPublishForm.category` 或对应字段。
- **严禁**：严禁只填入单一的 `yixiaoerId` 或末级分类名。

---

## 7. 账号有效性校验 (Account Validity)
- **校验标准**：Agent 在执行 `publish` 前，其选择的账号 `status` 必须为 `1` (有效)。
- **严禁执行**：若账号 `status` 为 `2` (失效)，Agent **严禁**尝试执行发布任务，应立即提示用户账号需要重新登录或处理。
- **数据来源**：账号状态必须通过 `action: "accounts"` 接口实时获取，不得依赖过期的历史数据。

## 8. 避坑指南：常见失败场景 (Common Failure Scenarios)

| 场景 | 原因分析 | 改进建议 |
| :--- | :--- | :--- |
| **版本不支持该功能** | 技能版本号过低，未包含新增动作或字段。 | 检查 `SKILL.md` 版本，提示用户升级或回退逻辑。 |
| **获取在线设备失败** | 选用了“本机发布”但客户端未开启或未登录。 | 提示用户：**“蚁小二客户端需要保持在线”**，或者建议改为 **“云发布”** (publishChannel: 'cloud')。 |
| **发布失败：表单错误** | Agent 生成的 JSON 结构与 `docs/` 下 aDTO 架构不匹配。 | **强制核对**：Agent 必须校验传入表单是否符合对应平台的文档要求。 |
| **任务卡在发布中/待发布** | 虽然 API 接收了请求，但表单内容不合规导致引擎挂起。 | **规则校验**：检查表单是否严格按照文档规则填入，特别是 `raw` 数据透传。 |
| **发布失败：必填项缺失** | 忽略了平台要求的必填核心字段（如标题、分类）。 | **必填项自检**：检查文档中标记为“必填”的字段是否已全部填入。 |
| **上传失败/签名不匹配** | `contentType` 在获取上传地址与执行 PUT 时不一致。 | **一致性校验**：必须确保 `action: "upload-resource"` 时传入的 `contentType` 与实际上传一致。 |
| **云发布报错：账号代理不存在** | 账号绑定的独立代理 (`proxyId`) 已失效，或云发布环境未配置代理。 | **修复建议**：通过 `update-account` 为账号设置 `kuaidailiArea`（内置代理）或有效的 `proxyId`（独立代理）。 |

---

## 9. 技能维护标准：文档重构规范 (Maintenance Standard)

为了保持技能文档的高度一致性与现代美学，在对技能 Markdown 文档进行新增、修改或重构时，技术人员与 Agent **必须**参照：
- **标准文档**：[🎨 蚁小二技能文档标准化优化指令 (Unified Prompt)](./standardization-prompt.md)

### 核心动作要求：
1. **结构对标**：必须包含 Trigger, Protocol, Parameters, Command, Troubleshooting 五大板块。
2. **视觉统一**：使用 GitHub Alert 展示关键信息，标题符合 `📄 [名] [型] 参数` 规范。
3. **数据完整**：不得删除历史参数，必须透传 `raw` 对象。

---

*本标准由 OpenClaw 架构组制定。违反本标准导致的发布失败或数据异常将由调用者承担相应责任。*
