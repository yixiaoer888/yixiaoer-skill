# 蚁小二技能 CLI 标准化重构执行文档

本文档用于指导 `yixiaoer-skill` 按照 `cli-main` 的工程结构和 AI Agent 友好标准进行重构。目标不是简单增加文档，而是把当前“Agent 手写大 JSON + 文档约束”的模式，改造成“小白用户可直接使用、Agent 可稳定调用、错误可自动恢复”的 CLI + Skill 体系。

## 1. 重构目标

### 1.1 用户目标

小白用户只需要表达自然语言需求，例如：

```text
帮我把 C:\videos\a.mp4 发布到抖音，标题是 今天的新品介绍
```

Agent 应能稳定完成：

1. 检查环境是否可用。
2. 查询并选择有效账号。
3. 上传视频和封面资源。
4. 自动组装发布参数。
5. 先执行 dry-run 预检。
6. 预检通过后执行发布。
7. 返回任务 ID 和后续查询命令。
8. 失败时给出可执行修复建议。

### 1.2 工程目标

当前仓库从：

```text
SKILL.md + docs/*.md + 零散脚本式历史实现
```

升级为：

```text
README.md
SKILL.md
cmd/
internal/
schemas/
workflows/
docs/
  quickstart.md
  workflows/
  errors.md
  refactor-to-cli-standard.md
tests/
  fixtures/
```

### 1.3 稳定性目标

重构完成后必须满足：

- 所有文档中的命令都能直接运行。
- 所有写操作都支持 `--dry-run`。
- 所有失败都输出结构化 JSON 错误。
- Agent 不需要手写 `publish` 大 JSON 即可完成常见任务。
- 平台名、发布类型、字段名只有一个标准源。
- 每个核心工作流都有 dry-run 测试。

## 2. 参考 `cli-main` 的标准

`cli-main` 中值得迁移的不是具体 Go 实现，而是这些工程原则：

| 标准 | `cli-main` 做法 | 蚁小二重构要求 |
| --- | --- | --- |
| 快速开始 | 人类用户和 AI Agent 分开写 | `README.md` 和 `SKILL.md` 分工明确 |
| 三层调用 | shortcut -> API command -> generic API | 快捷发布命令 -> 资源/账号命令 -> `api` 透传 |
| 参数约束 | flag、枚举、必填校验 | 禁止小白直接手写大 JSON |
| dry-run | 写操作可预览请求 | 发布、素材入库、账号更新必须支持 |
| 结构化错误 | `ok:false,error:{type,code,message,hint}` | 所有错误统一 envelope |
| 版本漂移 | `_notice.update` / `_notice.skills` | 增加 `doctor` 和版本检查 |
| 测试门禁 | unit + dry-run E2E + live E2E | 核心平台必须有 dry-run golden |

## 3. 目标 CLI 使用方式

### 3.1 初始化和诊断

```bash
yxer doctor
yxer config set-api-key
yxer auth check
```

输出示例：

```json
{
  "ok": true,
  "data": {
    "skillVersion": "2.0.0",
    "cliVersion": "2.0.0",
    "apiKey": "configured",
    "apiReachable": true,
    "defaultPublishChannel": "cloud"
  }
}
```

### 3.2 查询账号

```bash
yxer accounts list --platform 抖音
yxer accounts list --platform 抖音 --status valid
yxer accounts invalid
```

Agent 规则：

- 发布前必须查询账号。
- 只能使用登录状态有效的账号。
- 如果有多个匹配账号，必须让用户选择，或按用户指定昵称过滤。

### 3.3 上传资源

```bash
yxer upload --file ./cover.jpg
yxer upload --file ./video.mp4 --type video
yxer upload --url https://example.com/a.jpg --type image
```

CLI 必须自动完成：

- 推断 MIME 类型。
- 读取文件大小。
- 图片读取宽高。
- 视频尽量读取宽高、时长。
- 保证获取上传 URL 和 PUT 上传使用同一个 `contentType`。

### 3.4 视频发布

```bash
yxer publish video \
  --platform 抖音 \
  --account "账号昵称或ID" \
  --video ./video.mp4 \
  --cover ./cover.jpg \
  --title "今天的新品介绍" \
  --description "新品细节展示" \
  --dry-run
```

dry-run 通过后：

```bash
yxer publish video \
  --platform 抖音 \
  --account "账号昵称或ID" \
  --video ./video.mp4 \
  --cover ./cover.jpg \
  --title "今天的新品介绍" \
  --description "新品细节展示"
```

### 3.5 图文发布

```bash
yxer publish image-text \
  --platform 小红书 \
  --account "账号昵称或ID" \
  --image ./1.jpg \
  --image ./2.jpg \
  --title "新品图文" \
  --content "今天上新，欢迎了解" \
  --dry-run
```

### 3.6 文章发布

```bash
yxer publish article \
  --platform 知乎 \
  --account "账号昵称或ID" \
  --title "文章标题" \
  --content @article.md \
  --cover ./cover.jpg \
  --dry-run
```

CLI 必须负责：

- Markdown 转 HTML。
- 封面上传。
- 正文图片处理策略提示或自动上传。

### 3.7 保存草稿

蚁小二草稿：

```bash
yxer draft save \
  --type video \
  --platform 抖音 \
  --account "账号昵称或ID" \
  --video ./video.mp4 \
  --cover ./cover.jpg \
  --title "草稿标题"
```

平台草稿：

```bash
yxer publish video \
  --platform 抖音 \
  --account "账号昵称或ID" \
  --video ./video.mp4 \
  --cover ./cover.jpg \
  --title "草稿标题" \
  --platform-draft
```

如果用户只说“存草稿”，Agent 必须询问：

```text
你希望存为蚁小二草稿，还是推送到平台草稿箱？
```

### 3.8 素材库

```bash
yxer material add --file ./video.mp4 --cover ./cover.jpg --name "新品视频"
```

CLI 内部必须执行两步：

1. `upload` 到 `material-library` bucket。
2. `material` 登记入库。

禁止让用户或 Agent 只执行第一步后误以为入库成功。

### 3.9 发布记录和详情

```bash
yxer records list --status failed
yxer records detail --task-set-id TS123456
yxer records watch --task-set-id TS123456 --interval 10
```

失败任务必须输出：

- 失败平台。
- 失败账号。
- 后端错误。
- 建议修复命令。

## 4. 必须处理的功能清单

### 4.1 第一阶段：可用性修复

必须处理：

- 仓库内历史脚本入口已移除。
- 不再维护 TypeScript CLI 兼容层。
- 文档示例缺少必填字段。
- `imageText` / `image-text` 冲突。
- 平台枚举大小写和中文名冲突。
- 文章正文位置校验错误。
- `save-draft` 被错误要求 `contentPublishForm`。
- 错误输出不够结构化。

交付物：

- 可执行命令 `yxer`
- `yxer doctor`
- 结构化错误输出
- 基础单元测试

### 4.2 第二阶段：核心工作流

必须处理：

- `accounts list`
- `upload`
- `publish video`
- `publish image-text`
- `publish article`
- `draft save`
- `material add`
- `records list/detail`

交付物：

- 每个命令有 `--help`。
- 每个写命令有 `--dry-run`。
- 每个命令有 JSON 输出。
- 每个命令有至少一个 dry-run 测试。

### 4.3 第三阶段：平台增强

必须处理：

- 抖音：话题、位置、音乐、合集、小程序、游戏。
- 小红书：话题、位置、图文/视频差异。
- B 站：完整分类路径。
- 微信公众号：单独发布约束、`platformForms`。
- 快手、视频号、头条号等核心平台。

交付物：

- 平台能力矩阵。
- 平台字段 schema。
- 平台 dry-run golden 测试。

### 4.4 第四阶段：稳定性和恢复能力

必须处理：

- 账号失效自动阻断。
- 本机发布缺少 `clientId` 自动阻断。
- 云发布代理不存在时给出 `proxy-areas` / `update-account` 修复建议。
- 上传签名错误给出 `contentType` 修复建议。
- 分类缺失时提示先查询分类。
- 任务失败时可自动调用 `records detail` 辅助诊断。

交付物：

- `docs/errors.md`
- 错误码表
- 自动恢复建议
- live E2E 测试清单

## 5. 工程结构建议

### 5.1 目录结构

```text
src/
  cli.ts
  commands/
    doctor.ts
    config.ts
    accounts.ts
    upload.ts
    publish.ts
    draft.ts
    material.ts
    records.ts
    api.ts
  core/
    client.ts
    errors.ts
    output.ts
    env.ts
    mime.ts
    media.ts
    validators.ts
    platform-alias.ts
  schemas/
    publish.schema.json
    account.schema.json
    upload.schema.json
    platform.schema.json
  workflows/
    publish-video.ts
    publish-image-text.ts
    publish-article.ts
    material-add.ts
tests/
  unit/
  dryrun/
  e2e/
```

### 5.2 模块职责

| 模块 | 职责 |
| --- | --- |
| `commands/*` | 解析 CLI flags，不直接拼复杂 DTO |
| `core/client.ts` | 统一请求、鉴权、超时、响应解析 |
| `core/errors.ts` | 统一错误码和 hint |
| `core/output.ts` | 统一 `ok/data/error/notice` 输出 |
| `core/media.ts` | 图片/视频元数据读取 |
| `core/platform-alias.ts` | 平台别名归一化 |
| `core/validators.ts` | DTO 预检 |
| `schemas/*` | 唯一字段标准 |
| `workflows/*` | 自动执行多步骤业务链路 |

## 6. 统一输出标准

### 6.1 成功输出

```json
{
  "ok": true,
  "action": "publish.video",
  "version": "2.0.0",
  "data": {
    "taskSetId": "TS123456",
    "platforms": ["抖音"],
    "accounts": ["账号昵称"]
  },
  "next": {
    "check": "yxer records detail --task-set-id TS123456"
  }
}
```

### 6.2 失败输出

```json
{
  "ok": false,
  "error": {
    "type": "validation_error",
    "code": "YIXIAOER_MISSING_COVER",
    "message": "视频发布缺少封面图",
    "hint": "请传入 --cover ./cover.jpg，或先运行 yxer upload --file ./cover.jpg",
    "retryable": true,
    "nextCommand": "yxer publish video --cover ./cover.jpg ..."
  }
}
```

### 6.3 notice 输出

当 CLI 和 Skill 版本不一致时：

```json
{
  "ok": true,
  "data": {},
  "_notice": {
    "skills": {
      "current": "1.6.4",
      "target": "2.0.0",
      "message": "当前技能文档和 CLI 版本不一致",
      "command": "yxer update"
    }
  }
}
```

## 7. 错误码标准

| 错误码 | 类型 | 场景 | 修复建议 |
| --- | --- | --- | --- |
| `YIXIAOER_ENV_NO_API_KEY` | env_error | 未配置 API Key | 运行 `yxer config set-api-key` |
| `YIXIAOER_ENV_NETWORK` | env_error | API 不可达 | 检查网络或 `YIXIAOER_API_URL` |
| `YIXIAOER_VALIDATION_JSON` | validation_error | JSON 格式错误 | 使用 `--payload @file.json` |
| `YIXIAOER_INVALID_PLATFORM` | validation_error | 平台名无法识别 | 运行 `yxer platforms list` |
| `YIXIAOER_INVALID_PUBLISH_TYPE` | validation_error | 发布类型错误 | 使用 `video/image-text/article` |
| `YIXIAOER_ACCOUNT_NOT_FOUND` | validation_error | 找不到账号 | 运行 `yxer accounts list` |
| `YIXIAOER_ACCOUNT_INVALID` | validation_error | 账号登录失效 | 重新登录账号 |
| `YIXIAOER_MISSING_MEDIA` | validation_error | 缺少视频/图片 | 传入 `--video` 或 `--image` |
| `YIXIAOER_MISSING_COVER` | validation_error | 缺少封面 | 传入 `--cover` |
| `YIXIAOER_UPLOAD_SIGNATURE` | upload_error | 上传签名不匹配 | 使用自动 MIME 推断，保证 Content-Type 一致 |
| `YIXIAOER_LOCAL_CLIENT_REQUIRED` | validation_error | 本机发布缺少客户端 | 启动客户端或改用 `--channel cloud` |
| `YIXIAOER_PROXY_REQUIRED` | remote_error | 云发布代理缺失 | 运行 `yxer proxy areas` 后绑定代理 |
| `YIXIAOER_REMOTE_API` | remote_error | 后端返回错误 | 输出后端详情和 log id |

## 8. Schema 统一规则

### 8.1 发布类型

统一使用：

```text
video
image-text
article
```

如果后端需要 `imageText`，只允许在 client 层映射，文档和 CLI 不再暴露两套名字。

### 8.2 平台名

用户输入允许：

```text
抖音
Douyin
DouYin
douyin
```

内部统一为：

```text
抖音
```

后端如果需要 Code，由 client 层转换。

### 8.3 内容字段

统一规则：

- 文章正文：`accountForms[].contentPublishForm.content`
- 图文正文：`publishArgs.content`
- 视频标题/描述：`accountForms[].contentPublishForm.title/description`

禁止再出现脚本校验和文档不一致。

### 8.4 资源字段

统一要求：

- 视频必须有 `video.key`
- 视频必须有 `cover.key`
- 图文必须有 `images[].key`
- 文章封面按平台要求处理
- 所有资源必须来自 `upload` 返回 key

## 9. Skill 文档改造

### 9.1 `SKILL.md` 应缩短

`SKILL.md` 只保留 Agent 必须遵守的操作规则：

```md
1. 首次使用或失败时先运行 yxer doctor。
2. 发布不要手写大 JSON，优先使用 yxer publish 子命令。
3. 写操作先 dry-run。
4. 发布前必须查询账号并确认 status valid。
5. 本地或 URL 资源必须通过 CLI 自动上传。
6. 遇到 ok:false 时优先读取 error.hint 和 error.nextCommand。
7. 用户只说“存草稿”时必须询问草稿类型。
```

### 9.2 详细平台文档放到 references

平台字段文档保留，但不作为 agent 拼 DTO 的第一入口。Agent 优先调用 CLI，只有高级自定义时才查平台 references。

## 10. 测试标准

### 10.1 单元测试

必须覆盖：

- 平台别名归一化。
- 发布类型归一化。
- MIME 推断。
- 错误码生成。
- DTO 校验。
- 草稿歧义处理。

### 10.2 Dry-run 测试

每个核心命令至少一个 dry-run golden：

```bash
yxer publish video --platform 抖音 --account acc --video ./fixtures/a.mp4 --cover ./fixtures/c.jpg --title t --dry-run
yxer publish image-text --platform 小红书 --account acc --image ./fixtures/1.jpg --title t --content c --dry-run
yxer material add --file ./fixtures/a.mp4 --dry-run
```

测试断言：

- 不调用真实 API。
- 输出 `ok:true`。
- 生成的 DTO 字段完整。
- 平台名和发布类型已归一化。

### 10.3 Live E2E 测试

需要真实环境变量时才执行：

```text
YIXIAOER_API_KEY
YIXIAOER_TEST_DOUYIN_ACCOUNT
YIXIAOER_TEST_DRYRUN_ONLY=false
```

覆盖：

- 账号查询。
- 小文件上传。
- 保存蚁小二草稿。
- 查询发布记录。

真实发布应默认跳过，除非显式设置：

```text
YIXIAOER_ALLOW_REAL_PUBLISH=true
```

## 11. 实施里程碑

### Milestone 1：入口可用

预计 1-2 天。

- 巩固 `cmd/` 与 `internal/` 的 Go CLI 结构。
- 新增 `yxer doctor`。
- 补齐统一命令入口与错误输出。
- 清理历史脚本式入口残留。
- 文档命令全部替换为可运行命令。

验收：

```bash
yxer doctor
yxer accounts 抖音 --json
```

### Milestone 2：账号、上传、输出和错误

预计 2-3 天。

- `accounts list`
- `upload`
- 统一输出 envelope。
- 统一错误码。
- MIME 自动推断。
- API Key 检查。

验收：

```bash
yxer accounts list --platform 抖音 --format json
yxer upload --file ./fixtures/cover.jpg --dry-run
```

### Milestone 3：发布工作流

预计 3-5 天。

- `publish video`
- `publish image-text`
- `publish article`
- 自动上传资源。
- 自动查询账号。
- dry-run DTO 预览。

验收：

```bash
yxer publish video --platform 抖音 --account test --video ./fixtures/a.mp4 --cover ./fixtures/c.jpg --title test --dry-run
```

### Milestone 4：草稿、素材库、记录

预计 2-3 天。

- `draft save`
- `publish --platform-draft`
- `material add`
- `records list/detail/watch`

验收：

```bash
yxer material add --file ./fixtures/a.mp4 --cover ./fixtures/c.jpg --dry-run
yxer records list --status failed --dry-run
```

### Milestone 5：平台专项和稳定性

预计持续迭代。

- 抖音增强字段。
- 小红书增强字段。
- B 站分类路径。
- 微信公众号单独发布。
- 云发布代理修复建议。

验收：

- 每个平台有平台能力矩阵。
- 每个平台至少一个 dry-run golden。
- 常见错误都有 `hint` 和 `nextCommand`。

## 12. Agent 使用协议

重构完成后，Agent 必须按以下顺序处理用户发布请求：

1. 如果没有近期诊断结果，运行 `yxer doctor`。
2. 识别内容类型：`video` / `image-text` / `article`。
3. 查询目标平台账号：`yxer accounts list --platform <平台>`。
4. 如果账号不唯一，询问用户选择。
5. 调用对应 `yxer publish ... --dry-run`。
6. 如果 dry-run 失败，按 `error.hint` 修复。
7. dry-run 成功后执行真实发布。
8. 返回 `taskSetId` 和查询命令。

禁止：

- 未 dry-run 直接发布。
- 未上传资源直接把外部 URL 填进发布 DTO。
- 使用失效账号发布。
- 用户草稿意图不明时自行默认。
- 看到结构化错误后忽略 `hint`。

## 13. 小白用户 README 模板

重构完成后 README 首页应类似：

```md
# 蚁小二 CLI

## 3 分钟开始

1. 配置 API Key
   yxer config set-api-key

2. 检查环境
   yxer doctor

3. 查询账号
   yxer accounts list --platform 抖音

4. 预览发布
   yxer publish video --platform 抖音 --account 我的账号 --video ./a.mp4 --cover ./c.jpg --title 测试 --dry-run

5. 正式发布
   yxer publish video --platform 抖音 --account 我的账号 --video ./a.mp4 --cover ./c.jpg --title 测试
```

## 14. 最终验收清单

重构完成必须全部通过：

- [ ] `yxer doctor` 可运行。
- [ ] 所有 README 示例可运行。
- [ ] `SKILL.md` 不再要求直接执行 `.ts` 文件。
- [ ] 所有写命令支持 `--dry-run`。
- [ ] 所有错误输出结构化 JSON。
- [ ] 发布类型只有一套公开命名。
- [ ] 平台名支持别名并内部归一化。
- [ ] 发布前账号状态会校验。
- [ ] 上传自动处理 `contentType`。
- [ ] 素材库命令自动执行 upload + material 两步。
- [ ] 草稿语义明确区分蚁小二草稿和平台草稿。
- [ ] dry-run 测试覆盖视频、图文、文章、素材库。
- [ ] live E2E 默认不会真实发布。

## 15. 推荐优先级

如果资源有限，按这个顺序做：

1. 可执行 CLI 入口。
2. 结构化错误。
3. `doctor`。
4. `accounts list`。
5. `upload`。
6. `publish video --dry-run`。
7. `publish video` 正式执行。
8. 图文和文章。
9. 草稿和素材库。
10. 平台增强字段。

只要前 7 项完成，蚁小二技能的基础成功率就会显著高于当前版本。

