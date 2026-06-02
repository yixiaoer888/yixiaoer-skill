# yxer CLI 命令参考

`yxer` 是本技能唯一执行入口。Agent 和用户都应直接使用它。

## 命令分组

### 环境与版本

```bash
yxer --version
yxer doctor
yxer update [--check] [--global]
```

### 本地配置

```bash
yxer config get
yxer config set-local-client-id <clientId>
```

### Skill 安装与同步

```bash
yxer skill show
yxer skill sync [--global]
```

### 账号与资源

```bash
yxer accounts list [platform] [--name 关键词] [--status 1] [--json]
yxer upload --file <file_path> [--bucket cloud-publish|material-library] [--dry-run]
yxer upload --url <resource_url> [--bucket cloud-publish|material-library] [--dry-run]
```

### 发布与校验

```bash
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
yxer publish <type> <platform> <payload.json> [clientId] [--dry-run]
```

### 草稿与素材库

```bash
yxer draft save <payload.json> [--dry-run]
yxer material create <payload.json> [--dry-run]
yxer material add --file <文件路径或URL> [--thumb <缩略图路径或URL>] [--type image|video|file] [--dry-run]
```

### 查询类能力

```bash
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词] [--type 0|1|2|3]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
yxer records list [--platform P] [--limit N] [--status S] [--json]
yxer prepare <platform> <type>
yxer schema get <platform> <type>
```

## 基本约束

- 发布类型统一使用：`video`、`imageText`、`article`
- 单次 `yxer publish` 只处理一个平台
- `publish` 仅支持 `payload.json` 模式
- 发布前必须先执行 `yxer prepare <platform> <type>` 和 `yxer schema get <platform> <type>`，确认表单字段和 schema 后再填写 payload
- `payload.json` 只支持标准 `publishArgs` 结构，所有平台统一
- CLI 会根据 `publishArgs` 自动补齐最外层 `cover`、`coverKey`、`desc`、`isDraft`、`isAppContent`
- 云发布是默认模式
- 本机发布时必须提供 `clientId`
- `yxer validate`、`yxer publish --dry-run`、`yxer publish` 使用同一套发布通道解析逻辑
- 本机发布可通过三种方式提供 `clientId`：
  - 第四个位置参数：`yxer publish <type> <platform> <payload.json> <clientId>`
  - flags：`yxer publish <type> <platform> <payload.json> --publish-channel local --client-id <clientId>`
  - 预设默认值：`yxer config set-local-client-id <clientId>` 后，再执行 `--publish-channel local`
- 本机发布校验时，推荐在 `validate` 阶段就显式传入 `--publish-channel local`；若未显式传入但 payload 中已写 `publishChannel=local`，CLI 也会尝试从默认配置读取 `clientId`
- `yxer draft save` 只处理蚁小二内部草稿，不等同于平台草稿箱
- `yxer material create` 只做素材登记，前提是资源已经通过 `yxer upload --bucket material-library` 上传
- 推荐优先使用 `yxer material add --file ...`，由 CLI 自动完成上传和素材登记
- 查询类操作可以直接执行
- 发布类操作必须遵守“查账号 -> prepare/schema -> 上传资源 -> 查询复杂对象 -> 填 payload -> validate -> publish”顺序
- 所有请求字段都必须来自 schema、平台文档或 CLI 返回结果；严禁虚构字段、乱猜枚举、手写 `raw` 对象或编造资源元数据

## 快速示例

### 环境检查

```bash
yxer doctor
yxer config get
```

### 查询账号

```bash
yxer accounts list 抖音 --json
yxer accounts list 小红书 --status 1
```

### 上传资源

```bash
yxer upload --file .\cover.jpg --dry-run
yxer upload --file .\video.mp4
yxer upload --url https://example.com/demo.jpg
```

## 推荐发布流程

### 标准 payload 结构

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["小红书"],
  "publishChannel": "cloud",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "<platformAccountId>",
        "contentPublishForm": {
          "formType": "task"
        }
      }
    ]
  }
}
```

约束：

- 顶层必须有 `publishArgs`
- 账号列表必须放在 `publishArgs.accountForms[]`
- 平台业务字段必须放在 `publishArgs.accountForms[].contentPublishForm`
- 不再支持顶层 `accountForms`
- 不再支持直接提交内层业务表单 JSON

### 获取表单字段与 schema

```bash
yxer prepare 小红书 imageText
yxer schema get 小红书 imageText
```

### 校验与预览发布

```bash
yxer validate 小红书 imageText .\payload.json
yxer publish imageText 小红书 .\payload.json --dry-run
```

### 本机发布校验

```bash
yxer validate 抖音 video .\payload.json --publish-channel local --client-id <clientId>
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId> --dry-run
```

### 正式发布

```bash
yxer publish imageText 小红书 .\payload.json
```

### 本机发布

```bash
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId>
```

## 常见工作流入口

- 共享规则：`../yixiaoer-shared.md`
- 通用规则：`../workflows/common-rules.md`
- 图文发布：`../workflows/publish-imageText.md`
- 视频发布：`../workflows/publish-video.md`
- 文章发布：`../workflows/publish-article.md`

## 发布通道约定

- 用户未指定“本机发布 / 本地发布 / 客户端发布”时，Agent 应默认使用云发布。
- 用户明确要求本机发布，或说明要走本机客户端/本机网络时，Agent 必须显式传 `--publish-channel local`，不要只在说明文字里表达。
- 若云发布返回“账号代理不存在”等代理相关错误，可建议切换到本机发布。
- 若本机发布返回“客户端不在线”或“获取在线设备列表失败”，可建议用户启动蚁小二客户端，或改回云发布。

## 输出约定

- 默认输出适合人读
- 加 `--json` 时输出结构化结果，适合 Agent 二次处理
- 成功输出格式：`ok/action/version/data`
- 失败输出格式：`ok/version/error`
- 错误通过统一错误 envelope 输出
- `yxer doctor` 可能返回 `_notice.skills`，提示当前 AI skill 与 CLI 版本不同步
- `yxer update` 当前会同步 AI skill，并给出 CLI 本体更新指引

## 入口约束

- 仓库已移除旧 Node 入口，不再提供脚本兼容通道
- 未完成 CLI 化的能力只保留文档提示，不代表存在其他可执行入口
