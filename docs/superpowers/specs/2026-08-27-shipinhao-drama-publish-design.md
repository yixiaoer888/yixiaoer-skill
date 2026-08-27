# 视频号剧集挂载支持设计

## 背景

当前 `yxer` 对视频号只支持 `collection`（合集）。视频号页面中的“剧集”是另一项能力，不能复用合集字段。

已从视频号页面网络请求确认真实契约：

- 剧集查询：`GET /api/platform-accounts/{platformAccountId}/drama-tasks?keyWord=`。
- 发布提交：`POST /api/taskSets/v2`。
- 剧集字段：`publishArgs.accountForms[].contentPublishForm.drama`。
- 剧集对象字段：`yixiaoerId`、`yixiaoerImageUrl`、`yixiaoerName`。

成功发布请求中的剧集对象没有 `raw` 字段，因此不能套用合集、分类等动态对象的 `raw` 强制规则。

查询接口的候选响应按固定规则处理：先解一层 CLI 成功 envelope 的 `data`；之后接受直接数组、完整单对象，或仅含一个候选列表键的对象。列表键按 `items`、`list`、`dataList`、`records`、`results` 识别，并允许继续解一个嵌套 `data`。同一层出现多个候选列表键视为歧义；空数组、标量、未知 envelope、列表键类型错误以及缺少三个字段的候选均在 `publish form choose` 阶段返回结构化使用错误，不能写入 session。完整单对象必须恰好包含 `yixiaoerId`、`yixiaoerImageUrl`、`yixiaoerName` 三个非空字符串；额外字段同样在 choose 阶段拒绝。Schema validation 继续作为直接 payload 和会话被外部修改时的第二道防线。

## 目标与范围

本次为视频号视频增加独立的剧集查询、字段发现、表单选择、校验和发布透传能力：

1. 客户可以通过 `yxer query drama-tasks` 查询目标账号可挂载的剧集。
2. `yxer prepare`、`schema fields`、`schema get` 和 `publish form` 能发现并填写 `drama` 字段。
3. 选择查询结果后，CLI 将完整剧集对象写入 `contentPublishForm.drama`，保留 `yixiaoerId`、`yixiaoerImageUrl`、`yixiaoerName`。
4. `validate`、`publish --dry-run` 和正式 `publish` 使用同一字段结构，最终请求不把剧集转换成 `collection`，也不丢失剧集图片地址。
5. 错误仍遵守 stdout 纯 JSON、结构化错误和现有 CLI 命令协议。

不在范围内：

- 修改视频号后台或 `/taskSets/v2` 服务端。
- 修改 `posting-eligibility` 资格检查接口。
- 修改已有合集查询和 `collection` 字段语义。
- 自动从视频截取封面、自动选择标题或绕过发布前校验。
- 新增未经网络请求确认的剧集字段、接口或平台映射。

## 方案

### 1. 独立查询命令

新增命令：

```text
yxer query drama-tasks <account_id> [--query 关键词] [--keyword 关键词] --json
```

`--keyword` 作为现有查询命令的兼容别名，内部统一使用 `keyWord` 查询参数。`Client.DramaTasks` 必须使用 `url.Values.Set("keyWord", keyword)` 与 `QueryValues`，不能使用会丢弃空值的 `Query`；即使关键词为空，最终 `RawQuery` 也必须包含 `keyWord=`，以匹配页面请求。命令访问：

```text
/platform-accounts/{accountID}/drama-tasks?keyWord={keyword}
```

查询结果只做现有 `queryData` 的 envelope 解包，不改变接口返回的剧集对象，不补造 `raw`，也不把剧集对象压缩成只含 ID 的对象。`publish form choose` 增加 `dataList` 列表识别，并对 `drama` 做候选形状检查：`--value-file` 中必须能提取一个或多个完整剧集对象；未知的响应对象返回结构化使用错误，而不是把整个响应 envelope 当作剧集。这样 `--value-file` 可以直接复用 `yxer query drama-tasks ... --json` 的结果，仍可通过 `yixiaoerId` 或 `--index` 选择候选。

新增 API client、query workflow service 和 Cobra 子命令；更新 CLI 命令参考及关键词文档。

成功输出复用既有 `output.Success` 契约：stdout 仅有一个 JSON 对象，`ok=true`、`action="drama-tasks"`、`version` 沿用 CLI 版本，解包后的查询结果位于 `data`，`_notice` 仅在现有通知机制产生时出现。远端查询失败时 stdout 必须为空，结构化错误 JSON 仅写 stderr，并保留远端 `code/category/retryable`；同时补充指向当前账号的 `nextCommand: yxer query drama-tasks <account_id> --json` 和修复提示。

### 2. 视频号 schema 与字段发现

仅修改 `schemas/platforms/shipinhao.video.schema.json`，增加：

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

必填字段以成功发布请求中实际出现的三个字段为准。`drama` 的输入路径固定为：

```text
publishArgs.accountForms[].contentPublishForm.drama
```

在 `cmd/schema.go` 中把 `drama` 标记为查询型复杂字段，并提供：

```text
yxer query drama-tasks <account_id> [--query 关键词] --json
```

在 `cmd/dynamic_examples.go` 中为 `drama` 单独生成示例。示例只展示三个真实字段，不添加 `raw`，避免 `schema fields` 与真实发布请求不一致。`publish form choose` 会自动通过 schema 生成的 `dynamicFieldExamples` 获得该字段的可写路径。

`publish form choose` 的来源记录按规范路径绑定查询命令。规范路径固定为 `publishArgs.accountForms[N].contentPublishForm.drama`：

- `publish form set` 禁止写入该路径，必须提示改用 `choose`。
- `choose drama --path ...` 只允许与目标账号对应的规范 drama 路径完全一致；`choose` 的其他字段也不得通过 `--path` 写入规范 drama 路径。
- `drama` 只能使用 `yxer query drama-tasks <account_id> ...` 产生的 `sourceCommand`，并继续校验查询账号与目标账号一致。
- verify、review 和 export 必须重新按实际路径校验 `sourceCommand`、目标账号、`ValueHash` 和来源覆盖；只要 payload 中存在 `drama`，就必须存在对应的 query source。外部手改 session 绕过来源记录时验证失败。

剧集来源记录使用完整对象的 `ValueHash`，不生成也不要求 `RawHash`。是否应用这个例外只由实际规范路径决定，不能只看调用方传入的 `Field` 名称；因此把 `field=drama` 写到 collection 路径仍按 collection 规则要求 `RawHash`，把其他 field 写到 drama 路径则直接拒绝。其他动态字段维持原有来源规则。

为使 Agent 在 `schema get` 中看到严格的剧集对象边界，`PropertyView` 需要传播对象节点的 `additionalProperties`（使用可选布尔值以保留显式 `false`，不改变没有声明该约束的旧字段输出）。`schema fields` 继续通过 `dynamicFieldExamples.drama.value` 展示三个可填写字段。

现有 validator 会在任意对象层级放行 `content`、`clientId`、`video` 等 CLI 通用字段。对 `/drama` 路径必须禁用这项通用白名单，使显式 `additionalProperties: false` 真正严格；`raw`、`content`、`clientId` 及其他未声明字段都必须被拒绝。其他字段沿用现有兼容策略，避免扩大本次变更范围。发布 schema 错误只要包含 `drama`，`schemaValidationHint` 应提示重新执行 `yxer query drama-tasks <account_id> --json` 并使用完整三字段对象。

### 3. 发布与预校验

现有标准 payload 已经会把 `contentPublishForm` 作为账号表单的一部分交给发布 body 构造器，因此只要 schema 接受 `drama`，无需新增另一套发布入口。发布流程保持：

```text
query drama-tasks
    ↓
publish form choose / payload 写入 drama
    ↓
schema validation
    ↓
preflight
    ↓
publish --dry-run
    ↓
taskSets/v2
```

`drama` 不加入现有要求 `raw` 的通用动态对象字段列表，也不加入合集归一化逻辑。预校验仍会拒绝未声明字段和占位符；`yixiaoerImageUrl` 按平台元数据地址处理，沿用现有 `shouldIgnoreExternalURLPath` 规则，不会被当成待上传的视频或封面外链。`drama` 选择、verify、review、export 全链路都只比较完整对象的 `ValueHash`；`collection` 等已有字段继续比较 `RawHash`。

发布边界仅保留已有的分类 wire 转换，不能对 `drama` 做 ID/name 重包装。这样 dry-run 请求和实际 POST 请求都保留页面确认过的对象结构。

### 4. 文档与使用方式

新增 `skills/yixiaoer/references/get-drama-tasks.md`，说明接口、命令、对象字段及“剧集不等于合集”。同步更新：

- `skills/yixiaoer/references/platforms/video/shipinhao.md`：增加剧集触发场景、查询命令和 payload 片段。
- `skills/yixiaoer/references/domains/publish.md`、视频发布 workflow：增加剧集查询入口。
- `skills/yixiaoer/references/cli/command-reference.md`、`keyword-reference.md`：增加命令和术语。

视频号平台文档中的通用动态对象说明要拆分：`collection`/合集继续要求完整 `raw`；`drama` 是只含三个字段的独立对象，`yixiaoerImageUrl` 可以是查询返回的 HTTP 地址，不要求 `raw`。

推荐使用方式：

```bash
yxer query drama-tasks <account_id> --query 剧名 --json
yxer publish form choose publish-form.json drama \
  --value-file drama-tasks.json \
  --id <yixiaoerId> \
  --source-command "yxer query drama-tasks <account_id> --query 剧名 --json"
```

最终 payload 中的字段示例：

```json
{
  "contentPublishForm": {
    "formType": "task",
    "createType": 2,
    "pubType": 1,
    "drama": {
      "yixiaoerId": "event/<真实剧集标识>",
      "yixiaoerImageUrl": "<查询结果中的图片地址>",
      "yixiaoerName": "<查询结果中的剧集名称>"
    }
  }
}
```

示例中的尖括号只表示说明，正式 payload 必须使用当前 CLI 查询返回的完整值。

## 错误处理

- 剧集查询失败：沿用 API client 的结构化远端错误，不创建发布任务；错误中保留接口错误和下一步查询命令。
- 查询返回多个候选：`publish form choose` 要求 `--id` 或 `--index`，不能默认第一项。
- 空候选数组：choose 返回 `form choose has no candidates`；标量、未知 envelope、多个列表键或列表键类型错误：choose 返回 `drama query result is not a recognized candidate envelope`；选中的候选缺字段、字段为空或含额外字段：choose 返回 `drama query result does not contain a valid candidate`。以上错误都不能修改 session。
- 直接 payload 或被外部修改的 session 中，剧集对象缺少 `yixiaoerId`、`yixiaoerName`、`yixiaoerImageUrl`，或含 `raw`、`content`、`clientId` 等额外字段：schema validation 失败，返回字段路径和 `yxer query drama-tasks <account_id> --json` 修复提示。
- `drama` 的 `sourceCommand` 不是 `yxer query drama-tasks ...`：在选择阶段返回结构化来源错误；不允许用合集查询或手工 JSON 冒充剧集查询来源。
- 剧集候选响应不是数组、已知列表 envelope 或完整单对象：返回结构化使用错误，不把未知 envelope 写入发布 payload。
- `form set` 或不匹配的 `choose --path` 尝试写入规范 drama 路径：在写 session 前返回结构化使用错误；verify/review/export 对缺失来源和错误来源再次失败。
- payload 使用 `collection`：维持现有合集语义，不自动转成 `drama`。
- payload 使用未声明的 `drama` 子字段：schema validation 失败，防止把猜测字段发到 `/taskSets/v2`。

不新增专用的远端剧集错误码，除非后续真实接口返回证明需要区分查询、选择和发布阶段；现阶段使用现有错误 envelope 保持兼容。

## 数据流与边界

```text
真实账号 ID
    ↓
GET /platform-accounts/{id}/drama-tasks?keyWord=...
    ↓
CLI 返回完整候选对象
    ↓
publish form choose 按 yixiaoerId 选择
    ↓
contentPublishForm.drama（不加 raw、不改名）
    ↓
schema + preflight + dry-run
    ↓
POST /taskSets/v2
```

账号、剧集对象和发布资源均必须来自当前 CLI 查询或上传结果；剧集查询结果不得跨账号复用。`collection` 与 `drama` 在 schema、命令和 payload 中保持两个独立字段。

## 测试与验收

采用 TDD，先写能够失败的测试，再实现最小代码：

1. API client 使用准确的 `/platform-accounts/{id}/drama-tasks` 路径；分别断言空关键词的 `RawQuery` 包含 `keyWord=`、中文关键词可正确解码，并完整保留三个剧集字段。
2. `query` 命令暴露 `drama-tasks`，支持 `--query` 和 `--keyword`；成功 stdout envelope 的 action 固定为 `drama-tasks`，远端失败时 stdout 为空、stderr 为结构化 JSON 且包含当前账号的 `nextCommand`。
3. 视频号 `schema fields` 返回 `drama` 字段、正确路径、三个字段定义和 `drama-tasks` 查询命令；动态示例不包含 `raw`。
4. schema `get/fields` 能展示 `drama` 及其嵌套 `additionalProperties: false` 约束；视频号 schema 接受完整 `drama` 对象，拒绝缺少三个字段、`raw`、普通额外字段以及 validator 通用白名单中的 `content`/`clientId`。
5. `publish form choose` 覆盖直接数组、`items`、`list`、`dataList`、`records`、`results`、嵌套 `data` 和直接完整对象；空数组、标量、未知 envelope、多个列表键、列表键类型错误、不完整/空值/额外字段候选都返回明确结构化错误且 session 不变。
6. `drama` choose→verify/review/export 使用 `ValueHash` 完成 provenance，允许没有 `raw`；错误的 `sourceCommand`、跨账号来源、payload 漂移、缺失 source、`form set` 和不匹配的 `choose --path` 均被拒绝；collection 仍要求 `RawHash`，并覆盖 RawHash 漂移回归。
7. preflight 不因 `drama.yixiaoerImageUrl` 是 HTTP 地址而拒绝 payload，也不要求 `drama.raw`。
8. dry-run、云发布、本机发布及发布 body fallback 都保留 `contentPublishForm.drama` 的三个字段，不转换成 `collection` 或丢弃图片地址。
9. 多候选剧集未指定 ID/index 时会阻止选择；不存在的 ID 返回结构化错误；schema validation 的剧集错误返回 `drama-tasks` 专用修复提示。
10. 既有 `collection` 查询、合集 payload 和其他平台测试保持通过；全量 Go 测试、CLI schema/query 命令测试通过，且工作区 stdout/stderr 输出契约不回归。

验收标准：用网络请求中捕获的剧集对象填入视频号标准 payload，`yxer schema fields 视频号 video` 能发现该字段，`yxer validate` 和 `yxer publish --dry-run` 均通过，并且最终 `/taskSets/v2` 请求在 `contentPublishForm.drama` 中保留同一份查询对象；使用合集字段不会被误报为剧集，反之亦然。
