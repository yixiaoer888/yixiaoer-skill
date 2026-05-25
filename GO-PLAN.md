# Go + Cobra 重构开发计划

> 项目：yixiaoer-skill  
> 日期：2026-05-22  
> 目标：先搭好可扩展 Go + Cobra CLI 架构，再逐步迁移当前技能能力，后续按需扩展草稿、记录详情等功能。

## 1. 方案结论

使用 Go + Cobra 重构当前技能是可行的，且适合采用渐进式迁移。

当前仓库的核心能力集中在：

| 模块 | 文件 | 作用 |
| --- | --- | --- |
| Skill 规则 | `SKILL.md` | Agent 执行约束、命令说明、发布规则 |
| CLI 入口 | `scripts/yxer.ts` | 当前 `yxer` 命令实现 |
| API 封装 | `scripts/api.ts` | 鉴权、HTTP 请求、上传、低级 API action |
| Schema 校验 | `scripts/validator.ts` | AJV JSON Schema 校验 |
| Schema 文件 | `schemas/platforms/*.schema.json` | 各平台发布字段校验 |
| 工作流文档 | `workflows/*.md` | 图文、视频、文章发布流程 |
| 平台文档 | `docs/publish/**` | 各平台字段说明 |

重构时保留 `schemas/`、`workflows/`、`docs/`，先新增 Go CLI，与 TypeScript 版本并行。Go 版本稳定后，再切换 `yxer` 入口。

## 2. 重构原则

1. CLI 对外接口尽量兼容当前 `yxer <command>`。
2. 当前 TypeScript 版本保留，作为迁移期间兜底。
3. Go 版本优先迁移高频和关键链路：`doctor`、`accounts`、`upload`、`validate`、`publish`。
4. 发布能力不能只做 JSON Schema 校验，必须保留当前硬规则 preflight。
5. API 路径必须以现有 `scripts/api.ts` 和 `scripts/yxer.ts` 为准，不使用未验证的 `/openApi/...` 路径。
6. 文档和工作流先不大改，等 Go CLI 稳定后再同步更新 `SKILL.md`。

## 3. 当前真实 API 路径

迁移时优先复用当前 TypeScript 已验证路径：

| 能力 | 当前路径 |
| --- | --- |
| 账号列表 | `GET /v2/platform/accounts` |
| 发布内容 | `POST /taskSets/v2` |
| 保存蚁小二草稿 | `PUT /taskSets/drafts` |
| 发布记录 | `GET /v2/taskSets` |
| 任务详情 | `GET /v2/taskSets/{task_set_id}/tasks` |
| 获取上传 URL | `GET /storages/{bucket}/upload-url` |
| 分类 | `GET /platform-accounts/{accountId}/categories` |
| 位置 | `GET /platform-accounts/{accountId}/location` |
| 音乐 | `GET /platform-accounts/{accountId}/music` |
| 音乐分类 | `GET /platform-accounts/{accountId}/music/category` |
| 合集 | `GET /platform-accounts/{accountId}/collections` |
| 商品 | `GET /platform-accounts/{accountId}/goods` |
| 话题/挑战 | `GET /platform-accounts/{accountId}/challenges` |
| 群聊 | `GET /platform-accounts/{accountId}/group-chats` |
| 小程序 | `GET /platform-accounts/{accountId}/mini-apps` |
| 同步应用 | `GET /platform-accounts/{accountId}/sync-apps` |
| 游戏 | `GET /platform-accounts/{accountId}/games` |
| 代理列表 | `GET /proxys` |
| 代理地区 | `GET /daili/areas` |
| 更新账号 | `PATCH /platform-accounts/{accountId}` |
| 内容数据 | `GET /contents/overviews` |
| 账号数据 | `GET /platform-accounts/overviews-v2` |
| 素材库登记 | `POST /material` |

默认环境变量：

```bash
YIXIAOER_API_KEY
YIXIAOER_API_URL=https://www.yixiaoer.cn/api
```

## 4. 目标目录结构

```text
yixiaoer-skill/
├── cmd/
│   ├── root.go
│   ├── doctor.go
│   ├── account/
│   │   └── accounts.go
│   ├── publish/
│   │   ├── publish.go
│   │   ├── validate.go
│   │   ├── prepare.go
│   │   └── records.go
│   ├── resource/
│   │   └── upload.go
│   └── query/
│       ├── categories.go
│       ├── locations.go
│       ├── music.go
│       ├── goods.go
│       ├── collections.go
│       └── challenges.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── api/
│   │   ├── client.go
│   │   ├── accounts.go
│   │   ├── upload.go
│   │   ├── publish.go
│   │   ├── records.go
│   │   └── query.go
│   ├── modules/
│   │   ├── publish/
│   │   │   ├── service.go
│   │   │   ├── preflight.go
│   │   │   ├── validate.go
│   │   │   ├── prepare.go
│   │   │   └── types.go
│   │   ├── resource/
│   │   │   ├── service.go
│   │   │   └── types.go
│   │   ├── account/
│   │   │   ├── service.go
│   │   │   └── types.go
│   │   └── query/
│   │       ├── service.go
│   │       └── types.go
│   ├── schema/
│   │   ├── loader.go
│   │   └── validator.go
│   ├── domain/
│   │   ├── account.go
│   │   ├── resource.go
│   │   ├── publish.go
│   │   └── response.go
│   ├── output/
│   │   └── formatter.go
│   └── yxerrors/
│       └── errors.go
├── schemas/
├── workflows/
├── docs/
├── tests/fixtures/payloads/
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── .env.example
```

分层职责：

| 层 | 职责 |
| --- | --- |
| `cmd/` | Cobra 命令、参数解析、help 文案 |
| `internal/modules/*` | 按业务域划分的模块层，发布、资源、账号、查询各自独立 |
| `internal/api/` | HTTP 请求、鉴权、上传、后端路径封装 |
| `internal/schema/` | JSON Schema 加载和校验 |
| `internal/output/` | `--json` 和人类可读输出 |
| `internal/yxerrors/` | 统一错误码、hint、exit code |
| `internal/domain/` | DTO、响应模型、通用类型 |

模块分工建议：

| 模块 | 责任边界 |
| --- | --- |
| `publish` | 发布、validate、preflight、prepare、发布记录相关主链路 |
| `resource` | 资源上传、上传结果处理 |
| `account` | 平台账号查询、在线状态、账号更新 |
| `query` | 分类、位置、音乐、合集、商品、话题、群聊、小程序、同步应用、游戏等查询能力 |

## 5. 发布前必须保留的硬规则

Go 版本必须迁移当前 `scripts/yxer.ts` 里的 preflight 逻辑：

1. `publish type` 只能是 `video`、`image-text`、`article`。
2. `platforms` 不能为空。
3. payload 必须包含非空 `accountForms`。
4. 每个 `accountForms[]` 必须包含 `platformAccountId` 或 `account_id`。
5. 每个 `accountForms[]` 必须包含 `contentPublishForm`。
6. 视频发布必须包含上传后的 `video.key`，并校验 `size`、`width`、`height`。
7. 图文发布必须包含上传后的 `images[].key`，并校验 `size`、`width`、`height`。
8. 文章发布必须包含正文 `contentPublishForm.content`。
9. payload 内禁止出现外部 `http://` 或 `https://` 资源 URL。
10. `location`、`music`、`collection`、`challenge`、`goods`、`group`、`miniapp` 等动态对象必须包含完整 `raw`。
11. 发布前必须查询账号并确认 `status=1`。
12. `publish` 命令必须先执行 schema validate，失败则禁止发布。

## 6. 开发阶段

### Phase 0：Go CLI 骨架

目标：先搭出稳定框架，不迁移复杂业务。

交付物：

- `go mod init`
- `main.go`
- `cmd/root.go`
- `cmd/doctor.go`
- `internal/config`
- `internal/api/client.go`
- `internal/output`
- `internal/yxerrors`
- `.env.example`
- `Makefile`

命令：

```bash
yxer --help
yxer --version
yxer doctor
```

验收标准：

- 能编译出 `yxer` 二进制。
- `yxer doctor` 能检查 API Key、API URL、schema 目录是否存在。
- API Key 缺失时输出结构化错误。

### Phase 1：账号与上传

目标：迁移最基础的在线能力。

交付物：

- `yxer accounts [中文平台名] [--name kw] [--status 1] [--json]`
- `yxer upload <file_path_or_url> [--bucket cloud-publish|material-library]`
- 本地文件上传
- URL 下载后上传
- MIME 推断
- 图片宽高读取
- 文件大小读取

验收命令：

```bash
yxer accounts
yxer accounts 抖音 --status 1
yxer upload ./cover.jpg
yxer upload https://example.com/a.jpg
```

### Phase 2：Schema 与 Preflight

目标：把发布安全网先迁完，再接真实发布。

交付物：

- `yxer validate <中文平台名> <type> <payload.json>`
- `internal/schema/loader.go`
- `internal/schema/validator.go`
- `internal/preflight/publish.go`
- `image-text` 到 `imageText` 的内部 schema 文件名映射
- schema 缺失时 fallback basic validate

验收命令：

```bash
yxer validate 抖音 image-text tests/fixtures/payloads/douyin-video-valid.json
yxer validate 抖音 video tests/fixtures/payloads/douyin-video-valid.json
```

验收标准：

- schema 校验错误能指出字段路径。
- preflight 能阻断外部 URL、缺 key、缺 raw、缺账号字段。

### Phase 3：发布主链路

目标：迁移 `publish`，保证发布前校验完整。

交付物：

- `yxer publish <type> <中文平台名列表> <payload.json> [clientId]`
- 自动 schema validate
- 自动 preflight
- 自动检查账号 `status=1`
- 调用 `POST /taskSets/v2`
- 默认 `publishChannel=cloud`
- 传入 `clientId` 时使用 `publishChannel=local`

验收命令：

```bash
yxer publish image-text 抖音 payload.json
yxer publish video 抖音,快手 payload.json
```

验收标准：

- 校验失败不会调用发布 API。
- 账号不在线不会调用发布 API。
- 发布成功输出 task set 信息。

### Phase 4：查询命令迁移

目标：迁移当前 `yxer.ts` 中剩余查询类命令。

交付物：

- `yxer categories <account_id> [--type video|article]`
- `yxer locations <account_id> [--query kw] [--type 0|1|2|3]`
- `yxer music <account_id> [--query kw]`
- `yxer goods <account_id> [--query kw]`
- `yxer collections <account_id> [--type video|article]`
- `yxer challenges <account_id> [--query kw] [--type video]`
- `yxer records [--platform p] [--limit n] [--status s]`
- `yxer prepare <platform> <type>`

验收标准：

- 每个命令支持 `--json`。
- 输出结构和 TS 版本尽量兼容。

### Phase 5：测试与发布

目标：完成 Go 版本稳定验收。

交付物：

- `internal/api/client_test.go`
- `internal/schema/validator_test.go`
- `internal/preflight/publish_test.go`
- `cmd/*_test.go`
- 使用 `tests/fixtures/payloads/` 做回归测试
- mock API 测试
- 可选 live E2E 测试

验收命令：

```bash
go test ./...
make build
```

真实发布测试默认禁止，只有显式设置时才允许：

```bash
YIXIAOER_ALLOW_REAL_PUBLISH=true
```

## 7. 后续扩展阶段

Go 架构稳定后，再逐步增加以下能力。

### Phase 6：草稿

```bash
yxer draft save
yxer publish --platform-draft
```

注意：用户只说“草稿”时，Agent 仍必须询问是“蚁小二草稿”还是“平台草稿”。

### Phase 7：记录详情与诊断

```bash
yxer records detail --task-set-id xxx
yxer records watch --task-set-id xxx --interval 10
```

失败任务应输出：

- 失败平台
- 失败账号
- 后端错误
- 修复建议
- 下一步查询命令

## 8. 模块化演进路线

为了避免后续能力扩展把项目重新拧成一团，建议从一开始就按业务模块落地：

1. `publish` 作为第一核心模块，先把发布链路、校验链路、预检链路打通。
2. `account` 和 `query` 作为发布依赖模块，提供账号、分类、位置、音乐、合集等基础能力。
3. `resource` 作为独立模块，内部封装上传和资源元数据处理。
4. 未来如果新增其他业务模块，再按同样的目录模式单独加目录，不和发布逻辑混写。
5. 每个模块只暴露 service 层接口，`cmd/` 只负责命令参数到 service 的映射。

推荐约束：

- 一个模块一个目录。
- 一个模块一个 service 入口。
- 一个模块一组 DTO 类型。
- 公共能力尽量沉到 `api/`、`schema/`、`output/`，避免模块之间互相直接调用内部实现。

## 9. 命名与兼容规则

### 发布类型

对外统一使用：

```text
video
image-text
article
```

内部需要 schema 文件名时再映射：

```text
image-text -> imageText
```

### 输出格式

成功：

```json
{
  "success": true,
  "action": "accounts",
  "version": "3.0.0",
  "data": {}
}
```

失败：

```json
{
  "success": false,
  "errorCode": "YIXIAOER_USAGE_ERR",
  "message": "Failed to validate payload",
  "details": [],
  "suggestion": "请检查参数、schema 和 workflows 文档"
}
```

第一阶段优先兼容当前 `success/action/version/data` 格式；后续可以再升级为更标准的 `ok/data/error/next`，但不要在迁移初期破坏 Agent 依赖。

## 10. 里程碑

| 里程碑 | 阶段 | 验收条件 |
| --- | --- | --- |
| M0 | Phase 0 | `yxer doctor`、`yxer --help`、`yxer --version` 可用 |
| M1 | Phase 1 | `accounts`、`upload` 与 TS 版本功能等价 |
| M2 | Phase 2 | `validate` 和 preflight 能阻断非法 payload |
| M3 | Phase 3 | `publish` 能完整走通校验、账号检查、发布 |
| M4 | Phase 4 | 当前 12 个命令全部迁移 |
| M5 | Phase 5 | `go test ./...` 通过，可构建多平台二进制 |

## 11. 最小可用版本范围

第一版 Go CLI 不追求所有未来功能，只要求完成：

1. Cobra 框架稳定。
2. 配置、鉴权、错误输出统一。
3. `doctor` 可诊断。
4. `accounts` 可查账号。
5. `upload` 可上传资源。
6. `validate` 可校验 schema 和 preflight。
7. `publish` 可安全发布。

完成以上 7 项后，再开始团队、素材库、草稿和更多平台增强。
