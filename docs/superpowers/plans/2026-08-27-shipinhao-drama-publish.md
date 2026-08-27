# 视频号剧集挂载 Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 `yxer` 增加视频号剧集查询、Schema 发现、表单来源追踪和发布原样透传能力，同时保持合集语义与 stdout JSON 契约不变。

**Architecture:** 新增独立的 `drama-tasks` 查询入口，查询对象始终保留页面确认的三个字段。视频号 schema 将 `drama` 声明为严格对象；表单会话为该字段使用完整值哈希而非 `raw` 哈希，并将来源命令绑定到目标账号的 `drama-tasks` 查询。现有发布构造器继续透传 `contentPublishForm`，通过命令级 wire 测试锁定云发布、本机发布与回退行为。

**Tech Stack:** Go 1.25、Cobra、标准库 `net/http/httptest`、JSON Schema、Markdown 技能文档。

**Required skills:** @test-driven-development、@yixiaoer、@verification-before-completion。

---

## 文件结构

- `internal/api/query.go`：新增剧集查询客户端方法并保证空 `keyWord` 参数仍发送。
- `internal/api/query_test.go`：锁定接口路径、空/中文关键词编码和三字段响应透传。
- `internal/workflows/query/service.go`：向命令层暴露 `DramaTasks`。
- `cmd/query.go`：注册 `yxer query drama-tasks` 及 `--query`/`--keyword`。
- `cmd/query_test.go`：锁定命令发现和别名行为。
- `schemas/platforms/shipinhao.video.schema.json`：声明严格的 `drama` 对象。
- `internal/schema/document.go`：在 Agent 可见的 `PropertyView` 中传播显式 `additionalProperties`。
- `internal/schema/validator_test.go`：验证完整对象、缺字段、`raw` 和额外字段。
- `cmd/dynamic_examples.go`：生成不含 `raw` 的剧集动态字段示例。
- `cmd/schema.go`：把 `drama` 纳入复杂字段和查询命令提示。
- `cmd/schema_test.go`：锁定字段路径、严格约束和动态示例输出。
- `cmd/publish_form.go`：识别 `dataList`，验证剧集候选与来源命令，并为剧集使用 `ValueHash` provenance。
- `cmd/publish_form_test.go`：覆盖 choose→verify/review/export、未知 envelope、错误来源和合集回归。
- `internal/modules/publish/preflight_test.go`：锁定外链剧集图片不触发资源外链错误且不要求 `raw`。
- `internal/api/publish_test.go`：锁定 `/taskSets/v2` wire body 原样保留剧集对象。
- `cmd/publish_test.go`：锁定云发布、本机发布和自动回退时的最终请求体。
- `skills/yixiaoer/references/get-drama-tasks.md`：新增查询参考。
- `skills/yixiaoer/SKILL.md`、`skills/yixiaoer/references/domains/publish.md`、`skills/yixiaoer/references/workflows/common-rules.md`、`skills/yixiaoer/references/workflows/data-accuracy.md`、`skills/yixiaoer/references/workflows/payload-sourcing.md`、`skills/yixiaoer/references/workflows/publish-video.md`：将剧集加入动态字段来源纪律并注明无 `raw` 例外。
- `skills/yixiaoer/references/platforms/video/shipinhao.md`：区分剧集与合集，补充 payload 结构。
- `skills/yixiaoer/references/cli/command-reference.md`、`skills/yixiaoer/references/keyword-reference.md`：登记新命令和术语。

### Task 1: 增加剧集查询 API 与 CLI 命令

**Files:**
- Modify: `internal/api/query_test.go`
- Modify: `cmd/query_test.go`
- Modify: `internal/api/query.go`
- Modify: `internal/workflows/query/service.go`
- Modify: `cmd/query.go`

- [ ] **Step 1: 写 API 失败测试**

在 `internal/api/query_test.go` 增加测试，要求空关键词仍存在 `keyWord`，中文关键词正确解码，返回对象不丢字段：

```go
func TestDramaTasksAlwaysSendsKeyWordAndPreservesResponse(t *testing.T) {
	tests := []struct {
		name, keyword string
	}{
		{name: "empty", keyword: ""},
		{name: "chinese", keyword: "护妻"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/platform-accounts/acc_1/drama-tasks" {
					t.Fatalf("unexpected path: %s", r.URL.Path)
				}
				if tt.keyword == "" && r.URL.RawQuery != "keyWord=" {
					t.Fatalf("empty keyword must remain on wire, got %q", r.URL.RawQuery)
				}
				values, ok := r.URL.Query()["keyWord"]
				if !ok || len(values) != 1 || values[0] != tt.keyword {
					t.Fatalf("unexpected keyWord query: %#v", r.URL.Query())
				}
				_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": []map[string]interface{}{{
					"yixiaoerId": "event/1", "yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover", "yixiaoerName": "风浪过后护妻安康",
				}}})
			}))
			defer server.Close()

			client := NewClient(config.Config{APIKey: "test-key", APIURL: server.URL})
			result, err := client.DramaTasks("acc_1", tt.keyword)
			if err != nil { t.Fatal(err) }
			item := result.([]interface{})[0].(map[string]interface{})
			if item["yixiaoerId"] != "event/1" || item["yixiaoerImageUrl"] == "" || item["yixiaoerName"] == "" {
				t.Fatalf("unexpected drama task: %#v", item)
			}
		})
	}
}
```

- [ ] **Step 2: 写命令发现失败测试**

在 `cmd/query_test.go` 增加：

```go
func TestDramaTasksKeywordFlagUsesAliasStorage(t *testing.T) {
	assertKeywordFlagUsesAliasStorage(t, "drama-tasks", "护妻")
}

func TestQueryCommandExistsWithDramaTasksSubcommand(t *testing.T) {
	query := newQueryCmd()
	for _, child := range query.Commands() {
		if child.Name() == "drama-tasks" { return }
	}
	t.Fatal("expected query command to expose drama-tasks subcommand")
}

func TestDramaTasksCommandUsesJSONOutputContract(t *testing.T) {
	// 使用 httptest server 和 ExecuteWithIO 执行 query drama-tasks acc_1 --json；
	// 成功时断言 stderr 为空，stdout 是 {ok:true,action:"drama-tasks",version,data}；
	// 远端失败时断言 stdout 为空，stderr 是 {ok:false,error:{code,category,nextCommand}}，
	// 且 nextCommand 等于 yxer query drama-tasks acc_1 --json。
}
```

- [ ] **Step 3: 运行测试并确认 RED**

Run:

```powershell
go test ./internal/api ./cmd -run 'TestDramaTasks|TestQueryCommandExistsWithDramaTasks' -count=1
```

Expected: FAIL/compile failure，提示 `Client.DramaTasks` 或 `newDramaTasksCmd` 尚不存在。

- [ ] **Step 4: 实现最小查询链路**

在 `internal/api/query.go` 增加：

```go
func (c *Client) DramaTasks(accountID, keyword string) (interface{}, error) {
	values := url.Values{}
	values.Set("keyWord", keyword)
	result, err := c.queryData(QueryValues(fmt.Sprintf("/platform-accounts/%s/drama-tasks", accountID), values))
	if err != nil {
		return nil, withDramaTasksQueryRecovery(err, accountID)
	}
	return result, nil
}
```

`withDramaTasksQueryRecovery` 保留已有结构化远端错误字段，并仅补充 hint 与 `yxer query drama-tasks <account_id> --json` nextCommand；不得把诊断文本写入 stdout。

在 `internal/workflows/query/service.go` 增加转发方法：

```go
func (s Service) DramaTasks(accountID, keyword string) (interface{}, error) {
	return s.rt.Client.DramaTasks(accountID, keyword)
}
```

在 `cmd/query.go` 注册并实现：

```go
func newDramaTasksCmd() *cobra.Command {
	var query, keyword string
	cmd := &cobra.Command{
		Use: "drama-tasks <account_id>", Short: "查询视频号剧集", Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runQuery(cmd, "drama-tasks", func(service queryflow.Service) (interface{}, error) {
				return service.DramaTasks(args[0], resolveQueryAlias(query, keyword))
			})
		},
	}
	cmd.Flags().StringVar(&query, "query", "", "search keyword")
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword (alias for --query)")
	return cmd
}
```

- [ ] **Step 5: 运行测试并确认 GREEN**

Run: `go test ./internal/api ./cmd -run 'TestDramaTasks|TestQueryCommandExistsWithDramaTasks' -count=1`

Expected: PASS。

- [ ] **Step 6: 提交查询功能**

```powershell
git add internal/api/query.go internal/api/query_test.go internal/workflows/query/service.go cmd/query.go cmd/query_test.go
git commit -m "feat: add shipinhao drama task query"
```

### Task 2: 增加严格 Schema 与字段发现

**Files:**
- Modify: `internal/schema/validator_test.go`
- Modify: `cmd/schema_test.go`
- Modify: `schemas/platforms/shipinhao.video.schema.json`
- Modify: `internal/schema/document.go`
- Modify: `internal/schema/validator.go`
- Modify: `cmd/dynamic_examples.go`
- Modify: `cmd/schema.go`
- Create: `internal/workflows/publish/error_hints_test.go`
- Modify: `internal/workflows/publish/error_hints.go`
- Modify: `internal/modules/publish/preflight_test.go`
- Modify: `internal/api/publish_test.go`
- Modify: `cmd/publish_test.go`

- [ ] **Step 1: 写 Schema 失败测试**

在 `internal/schema/validator_test.go` 增加一个合法用例和三个非法用例：

```go
func TestShipinhaoVideoDramaSchemaIsStrict(t *testing.T) {
	validator := NewValidator(filepath.Join("..", "..", "schemas"))
	newValid := func() map[string]interface{} {
		return map[string]interface{}{
			"formType": "task", "createType": float64(2), "pubType": float64(1),
			"drama": map[string]interface{}{
				"yixiaoerId": "event/1", "yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover", "yixiaoerName": "风浪过后护妻安康",
			},
		}
	}
	if result := validator.Validate("视频号", "video", newValid()); !result.Valid {
		t.Fatalf("expected valid drama, got %v", result.Errors)
	}
	for name, mutate := range map[string]func(map[string]interface{}){
		"missing image": func(v map[string]interface{}) { delete(v["drama"].(map[string]interface{}), "yixiaoerImageUrl") },
		"raw": func(v map[string]interface{}) { v["drama"].(map[string]interface{})["raw"] = map[string]interface{}{} },
		"extra": func(v map[string]interface{}) { v["drama"].(map[string]interface{})["unknown"] = true },
		"common content": func(v map[string]interface{}) { v["drama"].(map[string]interface{})["content"] = "must be rejected" },
		"common clientId": func(v map[string]interface{}) { v["drama"].(map[string]interface{})["clientId"] = "must be rejected" },
	} {
		t.Run(name, func(t *testing.T) {
			candidate := newValid()
			mutate(candidate)
			if result := validator.Validate("视频号", "video", candidate); result.Valid {
				t.Fatalf("expected strict drama rejection: %#v", candidate)
			}
		})
	}
}
```

同时在 `TestSchemaReturnsValidShipinhaoVideoSchema` 中断言：

```go
drama := schemaDoc.Properties["drama"]
if drama.AdditionalProperties == nil || *drama.AdditionalProperties {
	t.Fatalf("expected drama additionalProperties=false, got %#v", drama)
}
```

- [ ] **Step 2: 写 CLI 字段发现失败测试**

在 `cmd/schema_test.go` 执行 `newSchemaGetCmd` 和 `newSchemaFieldsCmd`，断言：

```go
drama := data["businessFields"].(map[string]interface{})["drama"].(map[string]interface{})
if drama["additionalProperties"] != false { t.Fatalf("unexpected drama schema: %#v", drama) }
example := fieldsData["dynamicFieldExamples"].(map[string]interface{})["drama"].(map[string]interface{})
value := example["value"].(map[string]interface{})
if example["path"] != "publishArgs.accountForms[].contentPublishForm.drama" || example["queryCommand"] != "yxer query drama-tasks <account_id> [--query 关键词] --json" {
	t.Fatalf("unexpected drama example: %#v", example)
}
if _, exists := value["raw"]; exists { t.Fatalf("drama example must not contain raw: %#v", value) }
```

- [ ] **Step 3: 运行测试并确认 RED**

在运行前同时写好三层发布透传回归：

- `internal/modules/publish/preflight_test.go`：无 `raw` 且图片为 HTTP 地址的 `drama` 不产生错误、不被改写。
- `internal/api/publish_test.go`：捕获 `/taskSets/v2`，三字段对象深度相等，且没有新增 `raw`/`collection`。
- `cmd/publish_test.go`：cloud、local 和 `--auto-fallback-local` 捕获的每次请求都保留相同 `drama`。

preflight 测试主体：

```go
func TestPreflightAcceptsShipinhaoDramaWithoutRaw(t *testing.T) {
	payload := validVideoPayload()
	form := payload["publishArgs"].(map[string]interface{})["accountForms"].([]interface{})[0].(map[string]interface{})
	drama := map[string]interface{}{
		"yixiaoerId": "event/1", "yixiaoerImageUrl": "http://wxapp.tc.qq.com/cover", "yixiaoerName": "风浪过后护妻安康",
	}
	form["contentPublishForm"].(map[string]interface{})["drama"] = drama
	result := Preflight("video", []string{"视频号"}, payload)
	if len(result.Errors) != 0 { t.Fatalf("unexpected errors: %v", result.Errors) }
	if _, exists := drama["raw"]; exists { t.Fatal("preflight must not add drama.raw") }
}
```

API 与命令测试共用如下断言；API test 的 handler 对 decode 后的 body 调用它，命令测试对每次捕获的 body 调用它：

```go
func assertDramaPreservedOnWire(t *testing.T, body map[string]interface{}, expected map[string]interface{}) {
	t.Helper()
	args := body["publishArgs"].(map[string]interface{})
	form := args["accountForms"].([]interface{})[0].(map[string]interface{})
	got := form["contentPublishForm"].(map[string]interface{})["drama"]
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("drama changed on wire: got %#v want %#v", got, expected)
	}
	obj := got.(map[string]interface{})
	if len(obj) != 3 || obj["raw"] != nil || form["contentPublishForm"].(map[string]interface{})["collection"] != nil {
		t.Fatalf("unexpected drama wire shape: %#v", obj)
	}
}
```

命令级用例复用现有 `writePublishPayload`、`configureAPIKey`、`useTestAPIBaseURL` 和本机 client 配置 helper；cloud 断言一次请求，local 断言 `publishChannel=local`，fallback 断言两次请求均通过上述 helper。

前两项属于既有通用透传行为的 characterization test，可能立即通过，因此不授权修改生产代码；命令级测试和 schema 测试必须因 `drama` 尚未声明而 RED，作为本任务的实现门禁。

Run:

```powershell
go test ./internal/schema ./internal/modules/publish ./internal/workflows/publish ./internal/api ./cmd -run 'Test.*Drama|TestSchemaReturnsValidShipinhaoVideoSchema' -count=1
```

Expected: FAIL，因为视频号 schema 尚无 `drama`，`PropertyView` 也未暴露嵌套 `additionalProperties`。

- [ ] **Step 4: 实现严格 Schema 与可见约束**

在 `shipinhao.video.schema.json` 添加：

```json
"drama": {
  "type": "object",
  "required": ["yixiaoerId", "yixiaoerImageUrl", "yixiaoerName"],
  "properties": {
    "yixiaoerId": { "type": "string" },
    "yixiaoerImageUrl": { "type": "string" },
    "yixiaoerName": { "type": "string" }
  },
  "additionalProperties": false
}
```

在 `internal/schema/document.go` 的 `PropertyView` 增加：

```go
AdditionalProperties *bool `json:"additionalProperties,omitempty"`
```

并在 `buildPropertyView` 中仅在 schema 显式声明时赋值：

```go
if value, ok := schemaDoc["additionalProperties"].(bool); ok {
	view.AdditionalProperties = &value
}
```

在 `cmd/dynamic_examples.go` 添加专用示例函数：

```go
func addDramaExample(examples map[string]dynamicFieldExample, doc schema.Document) {
	if _, ok := doc.Properties["drama"]; !ok { return }
	examples["drama"] = dynamicFieldExample{
		Field: "drama", Path: "publishArgs.accountForms[].contentPublishForm.drama", Source: "query",
		QueryCommand: "yxer query drama-tasks <account_id> [--query 关键词] --json",
		Note: "剧集对象必须完整来自当前视频号账号查询；只保留三个真实字段，不添加 raw。",
		Value: map[string]interface{}{
			"yixiaoerId": "<from query>", "yixiaoerImageUrl": "<from query>", "yixiaoerName": "<from query>",
		},
	}
}
```

在 `buildDynamicFieldExamples` 调用该函数；在 `cmd/schema.go` 的复杂字段识别、字段类型识别和查询命令映射中加入 `drama`，并确保判断顺序先于可能的通用模式。

在 `internal/schema/validator.go` 将通用字段兼容改为路径感知，只有 `/drama` 禁用白名单：

```go
func isCLICommonOptionalFieldAtPath(path, key string) bool {
	if strings.TrimRight(path, "/") == "/drama" { return false }
	return isCLICommonOptionalField(key)
}
```

`validateObject` 使用该函数；其他路径行为不变。在 `error_hints_test.go` 先写失败测试，要求任一 drama schema 错误的 hint 包含 `yxer query drama-tasks <account_id> --json`；随后在 `schemaValidationHint` 的 switch 首位增加 drama 专用提示。

- [ ] **Step 5: 运行测试并确认 GREEN**

Run: `go test ./internal/schema ./internal/modules/publish ./internal/workflows/publish ./internal/api ./cmd -run 'Test.*Drama|TestSchemaReturnsValidShipinhaoVideoSchema' -count=1`

Expected: PASS。

- [ ] **Step 6: 提交 Schema 功能**

```powershell
git add schemas/platforms/shipinhao.video.schema.json internal/schema/document.go internal/schema/validator.go internal/schema/validator_test.go cmd/dynamic_examples.go cmd/schema.go cmd/schema_test.go internal/workflows/publish/error_hints.go internal/workflows/publish/error_hints_test.go internal/modules/publish/preflight_test.go internal/api/publish_test.go cmd/publish_test.go
git commit -m "feat: expose strict shipinhao drama field"
```

### Task 3: 修复表单候选与 provenance

**Files:**
- Modify: `cmd/publish_form_test.go`
- Modify: `cmd/publish_form.go`

- [ ] **Step 1: 写候选提取失败测试**

用同一个三字段 fixture 建立表驱动测试，成功输入覆盖直接数组、`items`、`list`、`dataList`、`records`、`results`、单层嵌套 `data` 和直接完整对象；均按 `event/1` 选中并断言写入对象恰好三个字段。失败输入分别覆盖：

```go
tests := []struct {
	name, value, wantError string
}{
	{"empty array", `[]`, "form choose has no candidates"},
	{"scalar", `42`, "drama query result is not a recognized candidate envelope"},
	{"unknown envelope", `{"page":1,"unknown":[]}`, "drama query result is not a recognized candidate envelope"},
	{"two list keys", `{"items":[],"dataList":[]}`, "drama query result is not a recognized candidate envelope"},
	{"wrong list type", `{"dataList":{}}`, "drama query result is not a recognized candidate envelope"},
	{"missing image", `[{"yixiaoerId":"event/1","yixiaoerName":"剧名"}]`, "drama query result does not contain a valid candidate"},
	{"empty name", `[{"yixiaoerId":"event/1","yixiaoerImageUrl":"http://cover","yixiaoerName":""}]`, "drama query result does not contain a valid candidate"},
	{"extra raw", `[{"yixiaoerId":"event/1","yixiaoerImageUrl":"http://cover","yixiaoerName":"剧名","raw":{}}]`, "drama query result does not contain a valid candidate"},
}
```

每个失败用例执行前后读取 session 文件并比较 payload/hash，证明错误不会产生本地写入。

- [ ] **Step 2: 写 provenance 失败测试**

新增完整链路测试：

```go
func TestPublishFormDramaProvenanceUsesValueHashWithoutRawHash(t *testing.T) {
	// choose 使用 yxer query drama-tasks acc_1 --query 护妻 --json；
	// 断言 ValueHash/FetchedAt 非空、RawHash 为空；
	// verify、review --dry-run、export --dry-run 均成功。
}

func TestPublishFormChooseDramaRejectsWrongQueryCommand(t *testing.T) {
	// source-command 使用 yxer query collections acc_1 --json；
	// 断言选择阶段返回结构化来源错误。
}
```

另加：

- `form set ...contentPublishForm.drama` 必须拒绝并提示 `choose`。
- `choose drama --path ...collection` 和 `choose collection --path ...drama` 都必须拒绝。
- 外部写入 payload.drama 但没有 source 时，verify/review/export 必须失败。
- 外部篡改已记录 drama 的 `sourceCommand` 后，verify 必须失败。
- 保留现有 `TestPublishFormChooseSelectsCandidateAndRecordsSource`，它必须继续断言普通动态字段 `RawHash` 非空；新增 collection `raw` 漂移后 verify 失败的回归。

- [ ] **Step 3: 运行测试并确认 RED**

Run:

```powershell
go test ./cmd -run 'TestPublishForm.*Drama|TestPublishFormChooseSelectsCandidateAndRecordsSource' -count=1
```

Expected: FAIL，分别暴露 `dataList` 未识别、错误来源未拒绝和 `missing_raw_hash`。

- [ ] **Step 4: 实现字段特定候选校验**

为 `drama` 使用独立的严格候选解析器，其他字段继续走既有宽松选择器：

```go
var dramaListKeys = []string{"items", "list", "dataList", "records", "results"}

func selectPublishFormCandidateForField(value interface{}, field string, index int, id string) (interface{}, []interface{}, error) {
	if strings.TrimSpace(field) != "drama" {
		return selectPublishFormCandidate(value, index, id)
	}
	candidates, err := dramaPublishFormCandidates(value, 0)
	if err != nil { return nil, nil, err }
	return selectPublishFormCandidate(candidates, index, id)
}

func isExactDramaCandidate(value interface{}) bool {
	obj, ok := value.(map[string]interface{})
	if !ok || len(obj) != 3 { return false }
	for _, key := range []string{"yixiaoerId", "yixiaoerImageUrl", "yixiaoerName"} {
		text, ok := obj[key].(string)
		if !ok || strings.TrimSpace(text) == "" { return false }
	}
	return true
}
```

`dramaPublishFormCandidates` 按规格固定顺序解析：数组必须非空且每项通过 `isExactDramaCandidate`；直接 map 若通过该函数则作为单候选；否则统计同层已知列表键，多个键或非数组值返回 envelope 错误；没有列表键时只允许在 depth=0 解一次 `data`；其余输入返回 envelope 错误。

- [ ] **Step 5: 实现来源命令绑定与 ValueHash 例外**

新增：

```go
func isCanonicalDramaFormPath(path string) bool {
	normalized := normalizePublishFormPath(path)
	return strings.HasSuffix(normalized, "accountForms.[].contentPublishForm.drama")
}

func validateDramaSourceCommand(path, command string) error {
	if !isCanonicalDramaFormPath(path) { return nil }
	parts := strings.Fields(command)
	for i := 0; i+1 < len(parts); i++ {
		if parts[i] == "query" && parts[i+1] == "drama-tasks" { return nil }
	}
	return yxerrors.Usage("form choose source command does not match drama", map[string]interface{}{"sourceCommand": command}).
		WithHint("drama 必须来自 yxer query drama-tasks <account_id> ... --json。")
}
```

在 `set` 写 session 前拒绝规范 drama 路径。choose 先算 contract 期望路径；当 field 或实际 path 涉及 drama 时，要求两者规范化后完全相同，再校验命令与账号。`appendPublishFormSource` 仅在 source 的实际 path 是规范 drama 路径时不生成 `RawHash`；`validatePublishFormSource` 按实际 path 重验命令、时间和 `ValueHash`，并只对该路径跳过 `missing_raw_hash`/`raw_hash_mismatch`。

`validatePublishFormProvenance` 还要遍历 payload 中所有 `accountForms[N].contentPublishForm.drama`：每个存在的字段必须至少有一个 path 完全匹配的 query source，否则新增 `missing_query_source`。verify、review、export 已共用该函数，因此三者都会阻止手改 session 的绕过。其他字段逻辑保持原样。

- [ ] **Step 6: 运行测试并确认 GREEN**

Run: `go test ./cmd -run 'TestPublishForm.*Drama|TestPublishFormChooseSelectsCandidateAndRecordsSource' -count=1`

Expected: PASS。

- [ ] **Step 7: 提交表单功能**

```powershell
git add cmd/publish_form.go cmd/publish_form_test.go
git commit -m "fix: support drama form provenance without raw"
```

### Task 4: 复核发布 wire 回归

**Files:**
- Verify: `internal/modules/publish/preflight_test.go`
- Verify: `internal/api/publish_test.go`
- Verify: `cmd/publish_test.go`
- Modify only if Task 2 中已先失败的测试证明必要: `internal/modules/publish/preflight.go`
- Modify only if Task 2 中已先失败的测试证明必要: `internal/api/publish.go`

Task 2 已在 schema 实现之前写入 preflight、API wire 和命令级 cloud/local/fallback 测试。本任务不再新增事后测试，只复核结果并在既有 RED 证明需要时做最小生产修复。

- [ ] **Step 1: 运行剧集发布回归**

Run:

```powershell
go test ./internal/modules/publish ./internal/api ./cmd -run 'Test.*ShipinhaoDrama|Test.*Drama.*Wire|Test.*Drama.*Fallback' -count=1
```

Expected: PASS。若测试仍失败，只做使对象原样透传的最小修复，不把 `drama` 加入 raw/归一化列表；修复后重跑同一命令确认 GREEN。

- [ ] **Step 2: 重跑发布相关回归**

Run:

```powershell
go test ./internal/modules/publish ./internal/api ./cmd -run 'TestPreflight|TestPublish' -count=1
```

Expected: PASS，包括已有 collection/raw 和本机回退测试。

- [ ] **Step 3: 仅在存在额外生产修复时提交**

```powershell
git add internal/modules/publish/preflight.go internal/api/publish.go
git commit -m "fix: preserve shipinhao drama on publish wire"
```

### Task 5: 同步 Agent 技能文档

**Files:**
- Create: `skills/yixiaoer/references/get-drama-tasks.md`
- Modify: `skills/yixiaoer/SKILL.md`
- Modify: `skills/yixiaoer/references/domains/publish.md`
- Modify: `skills/yixiaoer/references/workflows/common-rules.md`
- Modify: `skills/yixiaoer/references/workflows/data-accuracy.md`
- Modify: `skills/yixiaoer/references/workflows/payload-sourcing.md`
- Modify: `skills/yixiaoer/references/workflows/publish-video.md`
- Modify: `skills/yixiaoer/references/platforms/video/shipinhao.md`
- Modify: `skills/yixiaoer/references/cli/command-reference.md`
- Modify: `skills/yixiaoer/references/keyword-reference.md`

- [ ] **Step 1: 新增剧集查询参考**

文档必须包含：

```text
yxer query drama-tasks <account_id> [--query 关键词] --json
GET /platform-accounts/{platformAccountId}/drama-tasks?keyWord=
```

并说明候选对象只含 `yixiaoerId`、`yixiaoerImageUrl`、`yixiaoerName`，不添加 `raw`，不得跨账号复用。

- [ ] **Step 2: 更新公共来源纪律**

把 `drama` 加入动态查询表；将“所有动态字段都必须有 raw”改为“除平台 schema 明确声明的例外外”。唯一新增例外是视频号 `drama`，合集 `collection` 规则不变。

- [ ] **Step 3: 更新视频号平台文档**

增加触发词“挂剧集/关联剧集”，字段表和三字段 payload 示例；把复杂对象章节拆成：

```text
PlatformDataItem（位置/商品/合集/活动）：保留 raw
DramaTask（剧集）：yixiaoerId + yixiaoerImageUrl + yixiaoerName，不含 raw
```

- [ ] **Step 4: 更新命令和关键词索引**

将 `drama-tasks` 登记为查询视频号剧集；明确“合集≠剧集”。

- [ ] **Step 5: 做文档一致性检查**

Run:

```powershell
rg -n "drama-tasks|剧集|drama" skills/yixiaoer schemas/platforms/shipinhao.video.schema.json cmd internal
rg -n "所有.*raw|一律.*raw|必须.*raw" skills/yixiaoer/references
```

Expected: 所有命令、字段路径和三字段结构一致；不存在要求 `drama.raw` 的说明。

- [ ] **Step 6: 提交文档**

```powershell
git add skills/yixiaoer
git commit -m "docs: document shipinhao drama workflow"
```

### Task 6: 全量验证与交付

**Files:**
- Verify only: repository

- [ ] **Step 1: 格式化修改过的 Go 文件**

Run:

```powershell
gofmt -w internal/api/query.go internal/api/query_test.go internal/workflows/query/service.go cmd/query.go cmd/query_test.go internal/schema/document.go internal/schema/validator_test.go cmd/dynamic_examples.go cmd/schema.go cmd/schema_test.go cmd/publish_form.go cmd/publish_form_test.go internal/modules/publish/preflight_test.go internal/api/publish_test.go cmd/publish_test.go
```

- [ ] **Step 2: 运行聚焦测试**

Run:

```powershell
go test ./internal/api ./internal/schema ./internal/modules/publish ./cmd -run 'Drama|DramaTasks|ShipinhaoVideoSchema|PublishFormChooseSelectsCandidateAndRecordsSource' -count=1
```

Expected: PASS。

- [ ] **Step 3: 运行全量测试**

Run: `go test ./... -count=1`

Expected: PASS，stdout 无非 JSON 行为回归。

- [ ] **Step 4: 构建并检查 CLI 发现能力**

Run:

```powershell
go build ./...
go run . query drama-tasks --help
go run . schema fields 视频号 video
go run . schema get 视频号 video
```

Expected: help 中包含 `--query`/`--keyword`；两个 schema 命令返回合法 JSON，包含 `drama`，动态示例无 `raw`。

- [ ] **Step 5: 检查工作区与提交历史**

Run:

```powershell
git status --short
git diff --check
git log --oneline -8
```

Expected: 仅保留用户原有的 `skills-lock.json` 未跟踪文件；实现文件均已提交，`git diff --check` 无错误。
