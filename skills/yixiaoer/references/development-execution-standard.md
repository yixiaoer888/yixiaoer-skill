# 蚁小二技能后续开发迭代执行标准

本文档用于约束 `yixiaoer-skill` 后续新增功能、平台能力扩展、Bug 修复和文档更新，避免功能开发跑偏、文档与代码不一致、Agent 调用不稳定。

核心原则：

```text
任何新增能力都必须同时交付：命令入口 + schema 校验 + dry-run + 结构化错误 + 文档 + 测试。
```

如果某项能力无法满足以上要求，不允许标记为完成。

## 1. 开发总原则

### 1.1 以用户工作流为中心

新增功能不能只暴露后端接口字段，必须回答：

- 小白用户怎么用？
- Agent 通过什么命令稳定调用？
- 失败时用户/Agent 怎么修？
- 是否需要 dry-run？
- 是否会产生真实发布、写入、删除、账号修改等副作用？

### 1.2 CLI 是第一入口

面向用户和 Agent 的首选入口必须是 CLI 命令，而不是让 Agent 手写大 JSON。

允许存在底层 `api` 透传，但只能作为高级调试和兜底能力。

### 1.3 Schema 是唯一标准源

字段名、类型、枚举、必填项必须来自 schema。

禁止出现：

- 文档写一种字段，脚本校验另一种字段。
- 平台文档和 index 文档枚举不一致。
- 示例能跑但 schema 不允许。
- schema 允许但 CLI 无法生成。

### 1.4 写操作默认支持 dry-run

以下操作必须支持 `--dry-run`：

- 发布内容。
- 保存草稿。
- 上传素材库。
- 更新账号代理或分组。
- 删除、取消、批量修改类操作。

dry-run 必须输出最终请求结构，但不能调用真实写接口。

### 1.5 错误必须可恢复

错误不能只说“失败”，必须提供可执行修复建议。

错误输出标准：

```json
{
  "ok": false,
  "error": {
    "type": "validation_error",
    "code": "YIXIAOER_XXX",
    "message": "发生了什么",
    "hint": "用户或 Agent 下一步应该怎么做",
    "retryable": true,
    "nextCommand": "可选的修复命令"
  }
}
```

## 2. 新增功能准入流程

任何新增功能必须按以下流程执行。

### 2.1 第一步：写功能设计卡

每个功能先在 issue、PR 描述或设计文档中写清楚：

```md
## 功能名称

## 用户场景

## CLI 命令

## 输入参数

## 输出结构

## 后端接口

## 是否写操作

## 是否需要 dry-run

## 错误码

## 测试用例
```

没有功能设计卡，不进入编码。

### 2.2 第二步：更新或新增 schema

先改 schema，再改命令和文档。

必须明确：

- 字段名。
- 类型。
- 是否必填。
- 默认值。
- 枚举值。
- 字段来源，是用户输入、自动查询、上传结果还是后端返回。

### 2.3 第三步：实现 CLI 命令

新增能力必须有清晰命令。

命令命名规范：

```text
yxer <domain> <action>
yxer publish <type>
yxer records <action>
yxer material <action>
yxer accounts <action>
```

不允许把多个无关能力塞进一个巨大的 `--payload`。

### 2.4 第四步：实现 dry-run

写操作必须先实现 dry-run。

dry-run 输出必须包含：

- 归一化后的平台名。
- 归一化后的发布类型。
- 自动选择或匹配的账号。
- 自动上传后会使用的资源结构，测试环境可用 mock key。
- 最终将提交的 DTO。
- 风险提示。

### 2.5 第五步：实现真实执行

真实执行必须复用 dry-run 的构造逻辑。

禁止 dry-run 和真实执行各写一套 DTO 拼装代码。

推荐结构：

```text
buildRequest() -> dry-run 输出
buildRequest() -> execute 提交
```

### 2.6 第六步：补文档

每个新增命令必须补：

- `--help` 文案。
- README 或 quickstart 示例。
- Agent 使用规则。
- 错误处理说明。
- 至少一个可复制运行的示例。

示例必须真实可执行，不能缺少必填字段。

### 2.7 第七步：补测试

必须至少包含：

- 单元测试。
- dry-run golden 测试。
- 错误分支测试。

如果是核心链路，还要补 live E2E，但真实发布默认跳过。

## 3. 平台能力新增标准

新增一个平台或给平台新增字段，必须遵循本节。

### 3.1 平台能力卡

每个平台能力必须维护能力卡：

```md
# 抖音视频发布能力卡

## 支持内容类型
- video
- imageText

## 必填字段
- title
- description
- video
- cover

## 可选增强字段
- location
- music
- challenge
- collection

## 动态字段来源
- location -> yxer locations search
- music -> yxer music search
- challenge -> yxer challenges search

## 草稿能力
- 平台草稿：支持 pubType=0
- 蚁小二草稿：支持 draft save

## 已覆盖测试
- publish video dry-run
- publish video missing cover
- publish video invalid account
```

### 3.2 平台字段必须归一

用户输入允许别名，但内部必须唯一。

示例：

```text
用户输入：抖音 / Douyin / DouYin / douyin
内部标准：抖音
后端映射：由 client 层决定
```

### 3.3 动态字段不能猜

以下字段必须通过查询接口获取，不允许 Agent 凭常识填写：

- 分类。
- 地点。
- 音乐。
- 话题/挑战。
- 合集。
- 商品。
- 小程序。
- 游戏。
- 同步账号。

如果用户没有指定精确值，Agent 应先查询候选项，让用户选择或按明确规则选取。

## 4. 代码结构约束

### 4.1 命令层只做输入解析

`commands/*` 只负责：

- 定义 flags。
- 读取参数。
- 调用 workflow。
- 输出结果。

禁止在命令层直接拼后端 DTO。

### 4.2 Workflow 层负责业务流程

`references/workflows/*` 负责：

- 查询账号。
- 上传资源。
- 补全动态字段。
- 调用 validators。
- 构建请求。
- 执行或 dry-run。

### 4.3 Core 层负责通用能力

`core/*` 负责：

- API client。
- 结构化错误。
- 输出 envelope。
- MIME 推断。
- 媒体元数据。
- 平台别名。
- schema 校验。

### 4.4 文档不得成为逻辑源

文档只解释功能，不承载唯一逻辑。

字段标准必须在 schema 或代码常量中，文档由标准同步更新。

## 5. 错误码新增标准

新增错误码必须写入 `docs/errors.md`。

错误码命名：

```text
YIXIAOER_<DOMAIN>_<REASON>
```

示例：

```text
YIXIAOER_ACCOUNT_INVALID
YIXIAOER_UPLOAD_SIGNATURE
YIXIAOER_PUBLISH_MISSING_COVER
YIXIAOER_PLATFORM_UNSUPPORTED_FIELD
```

每个错误码必须包含：

- 触发场景。
- 用户可读 message。
- Agent 可执行 hint。
- 是否 retryable。
- 推荐 nextCommand。

## 6. 文档标准

### 6.1 README 面向小白用户

README 只放：

- 安装。
- 配置。
- `doctor`。
- 查询账号。
- 发布示例。
- 常见问题入口。

不要堆平台 DTO。

### 6.2 `skills/yixiaoer/SKILL.md` 面向 Agent

`skills/yixiaoer/SKILL.md` 只放 Agent 执行规则：

- 优先 CLI，不手写大 JSON。
- 写操作先 dry-run。
- 发布前查账号。
- 资源必须上传。
- 草稿歧义必须询问。
- 失败读取 `error.hint`。

### 6.3 平台文档面向高级配置

平台文档只描述平台差异，不重复基础 DTO。

每个平台文档必须写：

- 支持的内容类型。
- 必填字段。
- 可选字段。
- 动态字段来源。
- dry-run 示例。
- 不支持项。

## 7. 测试门禁

### 7.1 每次 PR 必跑

```bash
go test ./...
go build -o bin/yxer.exe .
```

如果后续增加其他语言工具链，也必须提供与 Go 主链路等价的校验命令。

### 7.2 Dry-run golden 测试

每个写命令必须有 golden 测试。

测试内容：

- 输入命令。
- 输出 JSON。
- 归一化结果。
- DTO 结构。
- 不触发真实 API。

### 7.3 文档示例测试

文档里的命令必须被自动扫描或人工清单覆盖。

禁止文档示例缺必填字段。

### 7.4 回归测试

修 Bug 时必须补一个能复现 Bug 的测试。

没有回归测试的 Bug 修复不允许合并，除非明确说明无法自动化，并提供人工验收步骤。

## 8. PR 验收清单

每个 PR 合并前必须检查：

- [ ] 是否有功能设计卡。
- [ ] 是否更新 schema。
- [ ] 是否提供 CLI 命令。
- [ ] 写操作是否支持 `--dry-run`。
- [ ] dry-run 和真实执行是否复用同一构造逻辑。
- [ ] 是否有结构化错误和 hint。
- [ ] 是否更新文档。
- [ ] 是否有单元测试。
- [ ] 是否有 dry-run 测试。
- [ ] 是否不会真实发布除非显式开启。
- [ ] 是否没有新增字段命名冲突。
- [ ] 是否没有让 Agent 手写大 JSON 的新入口。

## 9. 发布版本标准

### 9.1 版本号规则

使用 SemVer：

- `patch`：Bug 修复、文档修复、错误提示优化。
- `minor`：新增命令、新增平台能力、兼容性字段。
- `major`：破坏性命令或 schema 变更。

### 9.2 版本同步

以下版本必须同步：

- `skills/yixiaoer/SKILL.md` frontmatter
- CLI 输出版本
- changelog

如果版本不一致，`yxer doctor` 必须提示。

### 9.3 Changelog 标准

每次发布必须记录：

```md
## 2.1.0

### Added
- 新增 `yxer publish imageText`

### Fixed
- 修复 `imageText` 和 `imageText` 命名冲突

### Changed
- 平台名统一归一为中文名

### Migration
- 旧字段 `imageText` 仍兼容，但 CLI 统一使用 `imageText`
```

## 10. 高风险操作标准

以下操作视为高风险：

- 真实发布到平台。
- 批量发布。
- 删除内容。
- 修改账号代理。
- 批量修改账号配置。

高风险操作必须满足：

- 支持 `--dry-run`。
- 支持 `--yes` 或显式确认机制。
- 输出风险说明。
- Agent 不得自动加 `--yes`，必须获得用户明确确认。

## 11. Agent 行为标准

Agent 处理任何任务必须遵循：

1. 优先使用 CLI shortcut。
2. 不清楚环境时先 `yxer doctor`。
3. 发布前先 `accounts list`。
4. 写操作先 `--dry-run`。
5. dry-run 失败时按 `error.hint` 修复。
6. 真实执行前确认高风险意图。
7. 失败后优先读取 `error.nextCommand`。
8. 不要直接编造动态字段。

## 12. 防跑偏机制

为了保证后续开发不跑偏，必须建立以下机制：

### 12.1 PR 模板

新增 `.github/pull_request_template.md`，强制填写：

```md
## 用户场景

## CLI 命令

## Schema 变化

## Dry-run 输出

## 错误码

## 测试

## 文档
```

### 12.2 CI 检查

CI 至少包含：

- build。
- lint。
- unit test。
- dry-run golden。
- skill format check。
- docs command check。

### 12.3 文档示例检查

新增脚本扫描 Markdown 中的 `yxer ...` 命令，确保：

- 命令存在。
- 必填参数完整。
- dry-run 示例可执行。

### 12.4 Schema 漂移检查

新增脚本对比：

- schema 字段。
- CLI 参数。
- 文档示例。
- 平台文档字段。

发现同一概念多个名字时失败，例如：

```text
imageText
imageText
image_text
```

### 12.5 Golden 输出检查

dry-run 输出作为 golden 文件保存。功能变更导致 DTO 变化时，必须显式更新 golden，并在 PR 说明原因。

## 13. 新增功能完成定义

一个功能只有同时满足以下条件才算完成：

- 用户知道怎么用。
- Agent 知道怎么调用。
- CLI 能校验输入。
- dry-run 能看到请求。
- 真实执行复用 dry-run 构造。
- 失败有结构化错误。
- 文档示例能运行。
- 测试能防止回归。

任何只完成后端透传或只补文档的功能，都不能算完成。


