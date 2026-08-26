# 搜狐号视频分类发布修复 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让 `yxer` 在搜狐号视频发布前查询每个账号的最新分类，向客户展示父子路径，并自动生成服务端可识别的完整分类 payload。

**Architecture:** 在 `internal/api` 增加纯分类树解析、路径展平和搜狐号 payload 解析能力；`publish.Service.Prepare` 在 schema 校验前调用该能力，使 validate、dry-run、正式发布共享同一份已解析 payload。`query categories --paths` 只扩展查询展示，不改变默认 JSON 结果；其他平台保持现有流程。

**Tech Stack:** Go、Cobra、`httptest`、现有 `yxerrors`、JSON schema validator。

---

### Task 1: 建立搜狐号分类路径解析器

**Files:**
- Create: `internal/api/souhuhao_categories.go`
- Test: `internal/api/souhuhao_categories_test.go`

- [ ] **Step 1: Write the failing tests**

覆盖以下纯函数行为：

- 从 canonical `yixiaoerId/yixiaoerName/raw` 树递归生成任意深度路径；
- 将末级 ID 或名称匹配到唯一叶子，并返回根到叶子的 canonical 对象数组；
- ID 优先于名称；同名或身份字段冲突返回可区分错误；
- 生成 `path`、`nodes`、两级路径的 `parentId/childId` 展示字段；
- 每个 canonical 节点的 `raw.id` 存在且与 `yixiaoerId` 一致。

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/api -run 'Test(Souhuhao|CategoryPath)' -count=1`

Expected: FAIL because the resolver types/functions do not exist.

- [ ] **Step 3: Write the minimal implementation**

在新文件中定义内部分类节点/路径结构和以下边界能力：

- 递归读取 `child` 或 `children`；
- 只接受带 `yixiaoerId`、`yixiaoerName`、`raw.id` 的 canonical 节点；
- 按 ID、名称解析请求候选，拒绝不一致身份；
- 返回完整 canonical 路径，不从请求 payload 拼装父节点；
- 生成客户可读的路径摘要。

固定内部数据契约：查询 canonical 节点为 `{yixiaoerId, yixiaoerName, raw}`，其中 `raw` 必须是 map 且 `raw.id` 非空；展示路径为 `{path, nodes, category}`，`nodes` 支持任意深度。同一候选项的 `id`、`raw.id`、名称冲突，以及 payload 父子路径冲突，都返回 `sohuhao_category_invalid`。

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/api -run 'Test(Souhuhao|CategoryPath)' -count=1`

Expected: PASS。

### Task 2: 实现每账号查询和搜狐号 payload 归一化

**Files:**
- Modify: `internal/api/souhuhao_categories.go`
- Test: `internal/api/souhuhao_categories_test.go`

- [ ] **Step 1: Write the failing tests**

使用 `httptest.Server` 模拟 `/platform-accounts/{accountID}/categories`，测试：

- 每个 `platformAccountId` 独立发起分类查询；
- handler 断言请求方法为 `GET`、路径精确为 `/platform-accounts/{accountID}/categories`，并断言查询参数 `publishType=video`；
- 第二次失败任务的旧格式（`id/text` 加平台原始 `raw`）被转换为 `wireItem.raw` 为 canonical 对象，且 `wireItem.raw.raw.id` 存在；
- 只传子分类时补齐父分类；三层路径保留全部节点；
- payload 原 map 不被修改；
- 查询失败、空树、raw 缺少 `id`、分类不存在、名称歧义和父子冲突分别返回稳定错误 code/category/message/hint/nextCommand/retryable；details 必须包含 `accountId`、`requested`、`availablePaths`，查询失败、空树和结构无效还必须包含 `cause`，且 `availablePaths=[]`；
- 失败时不产生可继续发布的 payload。

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/api -run 'TestNormalizeSohuhao|TestResolveSohuhao' -count=1`

Expected: FAIL because the API method and error contracts are not implemented。

- [ ] **Step 3: Write the minimal implementation**

增加 `Client.NormalizeSohuhaoVideoPayload(payload map[string]interface{}) (map[string]interface{}, error)`。该方法通过 receiver 中已有的 HTTP client 调用 `c.Categories(accountID, "video")`，不新增依赖注入接口；它深拷贝传入 map 后返回新 map，原 map 不变：

- 只处理 `publishType=video` 且平台为 `搜狐号` 的 payload；
- 对每个 account form 查询自身账号的分类树，不共享结果、不使用缓存或旧 payload 回退；
- 支持 canonical、正确前端包装、旧 raw 直挂、只含 ID/名称的候选；
- 以查询结果重新组装父子路径和服务端 wire 结构；
- 查询/解析失败返回 `yxerrors.Error`，详情包含 `accountId`、`requested`、`availablePaths` 和可选 `cause`；
- 不修改传入 payload。

错误契约固定为：

| 场景 | code | category | retryable | message |
| --- | --- | --- | --- | --- |
| 查询失败 | `sohuhao_category_query_failed` | `sohuhao_category` | `true` | `搜狐号视频分类查询失败` |
| 空树/结构无效/raw 缺 ID | `sohuhao_category_data_invalid` | `sohuhao_category` | `false` | `搜狐号视频分类数据无效` |
| 分类不存在 | `sohuhao_category_not_found` | `sohuhao_category` | `false` | `搜狐号视频分类不存在` |
| 名称歧义 | `sohuhao_category_ambiguous` | `sohuhao_category` | `false` | `搜狐号视频分类名称不唯一` |
| 候选身份或父子路径冲突 | `sohuhao_category_invalid` | `sohuhao_category` | `false` | `搜狐号视频分类参数无效` |

每个错误都设置非空 `hint`、`nextCommand`，其中 `<accountID>` 必须替换为当前失败 form 的实际账号 ID（精确命令为 `yxer query categories <accountID> --type video --json`），并在 details 中返回 `accountId`、`requested`、`availablePaths`；查询失败额外返回 `cause`，且 `availablePaths` 固定为空数组。测试必须逐场景断言精确 `message`、非空 `hint`、账号对应的 `nextCommand` 以及 details 字段。

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/api -run 'TestNormalizeSohuhao|TestResolveSohuhao' -count=1`

Expected: PASS。

### Task 3: 将归一化接入 publish prepare

**Files:**
- Modify: `internal/workflows/publish/service.go`
- Test: `cmd/publish_test.go`
- Test: `internal/workflows/publish/service_test.go` (create if no focused test exists)

- [ ] **Step 1: Write the failing tests**

新增 CLI/Service 回归测试：

- Sohu 视频的 `Prepare` 在 schema 校验前查询分类并生成父子 wire 结构；
- 用一个同时存在 schema 错误和分类查询错误的 payload 断言分类查询先发生，最终返回分类错误而不是继续提交 schema 后的远端任务；
- 旧格式 payload 不再只通过 schema 后在远端失败；
- `validate` 在分类查询失败、空树、raw 缺 ID、不存在、歧义、身份冲突和父子冲突时均返回结构化 JSON 错误，且不调用 `/taskSets/v2`；
- dry-run 调用分类查询但不调用 `/taskSets/v2`，输出 request 中包含完整父子路径；
- 正式发布在分类查询/解析任一失败时不调用 `/taskSets/v2`；成功时断言发送给 `/taskSets/v2` 的每个 `wireItem.raw.raw.id` 存在，且服务端模拟 `parsePlatformCategory` 后的 `category.raw.id` 等于平台分类 ID；
- 非搜狐视频不调用搜狐分类接口且维持现有分类结构；
- 多账号时分别返回不同分类树，断言每个 form 使用自己 `platformAccountId` 对应的分类，任一账号失败时整体不发布。

对上述每一种失败场景使用 table-driven 测试逐项断言：`/taskSets/v2` 调用次数为 0，stdout 是可解析 JSON，stderr 不污染 JSON，错误的 `message`、`hint`、`nextCommand`、`details.accountId`、`details.requested`、`details.availablePaths` 和必要时的 `details.cause` 均符合契约。最终发布请求测试在 HTTP handler 中先读取每个 wire item 的 `raw` 作为服务端 canonical category，再断言 `wireItem.raw.raw.id == category.raw.id`。

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd ./internal/workflows/publish -run 'Test.*Sohuhao|Test.*Sohu|Test.*DryRun.*Category' -count=1`

Expected: FAIL because `Prepare` 当前未调用搜狐号分类解析器。

- [ ] **Step 3: Write the minimal implementation**

在 `Service.Prepare` 的第一步将输入 payload 做递归深拷贝，避免现有资源元数据解析和后续归一化修改调用方嵌套 map；完成发布类型/平台解析、资源元数据解析后、schema validation 前调用 API 归一化方法：

```go
if platform == "搜狐号" && input.PublishType == "video" {
    resolvedPayload, err = s.rt.Client.NormalizeSohuhaoVideoPayload(resolvedPayload)
    if err != nil { return PreparedPublish{}, err }
}
```

将返回的副本继续交给现有标准归一化、schema、preflight 和 body builder。`validate`、`DryRun`、`Execute` 都经过同一个 `Prepare`；保留现有其他平台分支，确保查询失败在创建 task set 之前返回。

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd ./internal/workflows/publish -run 'Test.*Sohuhao|Test.*Sohu|Test.*DryRun.*Category' -count=1`

Expected: PASS。

### Task 4: 增加父子分类展示视图和文档

**Files:**
- Modify: `cmd/query.go`
- Test: `cmd/query_test.go`
- Modify: `internal/api/souhuhao_categories.go`
- Modify: `skills/yixiaoer/references/get-publish-categories.md`
- Modify: `skills/yixiaoer/references/platforms/video/souhuhao.md`

- [ ] **Step 1: Write the failing tests**

测试 `yxer query categories <account> --type video --paths --json`：

- 返回 JSON envelope；
- `data.categories` 保留原始规范化树；
- `data.paths` 包含 `path`、`nodes`、父子 ID/名称和完整 category 对象；
- 任意深度路径都能展示；
- 默认不带 `--paths` 的输出 shape 不变。
- 使用 `ExecuteWithIO` 验证 stdout 可被完整解析为 JSON，stderr 不包含分类结果或诊断以外的 stdout 数据。

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./cmd -run 'Test.*Categories.*Paths' -count=1`

Expected: FAIL because `--paths` 尚不存在。

- [ ] **Step 3: Write the minimal implementation**

为 categories 命令增加 `--paths`，在 command 层接收 `Client.Categories` 返回的原值：若是含 `dataList` 的 envelope，先以 `dataList` 作为分类树；若是数组则直接使用；使用查询结果递归生成展示数据：`data.categories` 保留原树，`data.paths` 为完整路径对象。新增 envelope fixture 测试，确认 resolver 和 paths 视图都不会把 `dataList` 外壳误当成分类节点。stdout 仍只输出 JSON，不打印纯文本树。更新帮助文本和平台文档，明确搜狐号支持父子级 `category`。

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./cmd -run 'Test.*Categories.*Paths' -count=1`

Expected: PASS。

### Task 5: 运行回归测试和构建验证

**Files:**
- Verify only; do not modify unrelated user files:
  - `internal/media/video.go`
  - `internal/schema/validator_test.go`
  - `cmd/schema_shipinhao_original_test.go`
  - `internal/media/video_test.go`

- [ ] **Step 1: Run focused tests**

Run: `go test ./internal/api ./internal/workflows/publish ./cmd -run 'Test.*Sohuhao|Test.*Sohu|Test.*Categories.*Paths|Test.*Validate.*Category' -count=1`

Expected: PASS。

- [ ] **Step 2: Verify stdout/error contracts**

Run: `go test ./cmd -run 'Test.*JSON|Test.*Stdout|Test.*Sohuhao' -count=1`

Expected: every command response is parseable JSON; all category failures carry exact `message`, non-empty `hint`, account-specific `nextCommand`, `details.accountId/requested/availablePaths` (and query `cause`), and no failure reaches `/taskSets/v2`。

- [ ] **Step 3: Run all Go tests**

Run: `go test ./... -count=1`

Expected: PASS with no test failures。

- [ ] **Step 4: Build the CLI**

Run: `go build ./...`

Expected: exit code 0。

- [ ] **Step 5: Inspect the diff and user changes**

Run: `git -c safe.directory=D:/Projects/yixiaoer-skill diff --stat; git -c safe.directory=D:/Projects/yixiaoer-skill status --short`

Expected: only the planned Sohu files/docs are changed by this task; existing user changes remain untouched and uncommitted。
