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
yxer validate <platform> <type> <payload.json>
yxer publish <type> <platform> <payload.json> [clientId] [--dry-run]
yxer publish imageText <platform> --account <账号名或ID> --title <标题> --description <正文> --image <图片路径或URL> [--image ...] [--dry-run]
yxer publish article <platform> --account <账号名或ID> --title <标题> --content @<html文件> [--cover <封面路径或URL>] [--dry-run]
yxer publish video <platform> --account <账号名或ID> --title <标题> --description <描述> --video <视频路径或URL> [--cover <封面路径或URL>] [--dry-run]
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
- 推荐优先使用 `publish` 的 flags 模式，让 CLI 自动解析账号、上传资源并组装 payload
- `payload.json` 模式继续保留，适合高级定制或手工调试
- 云发布是默认模式
- 本机发布时必须提供 `clientId`
- 本机发布可通过三种方式提供 `clientId`：
  - 第四个位置参数：`yxer publish <type> <platform> <payload.json> <clientId>`
  - flags：`yxer publish <type> <platform> <payload.json> --publish-channel local --client-id <clientId>`
  - 预设默认值：`yxer config set-local-client-id <clientId>` 后，再执行 `--publish-channel local`
- `yxer draft save` 只处理蚁小二内部草稿，不等同于平台草稿箱
- `yxer material create` 只做素材登记，前提是资源已经通过 `yxer upload --bucket material-library` 上传
- 推荐优先使用 `yxer material add --file ...`，由 CLI 自动完成上传和素材登记
- 查询类操作可以直接执行
- 发布类操作必须遵守“查账号 -> 上传资源 -> 查询复杂对象 -> validate -> publish”顺序

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

## 推荐发布方式

### 图文

```bash
yxer publish imageText 小红书 \
  --account "图文账号" \
  --title "图文标题" \
  --description "图文正文" \
  --image ./1.jpg \
  --image ./2.jpg \
  --dry-run
```

### 文章

```bash
yxer publish article 知乎 \
  --account "知乎账号" \
  --title "文章标题" \
  --content @./article.html \
  --cover ./cover.png \
  --dry-run
```

### 视频

```bash
yxer publish video 抖音 \
  --account "视频账号" \
  --title "视频标题" \
  --description "视频描述" \
  --video ./clip.mp4 \
  --cover ./cover.png \
  --dry-run
```

### 本机发布

```bash
yxer publish video 抖音 \
  --account "视频账号" \
  --title "视频标题" \
  --description "视频描述" \
  --video ./clip.mp4 \
  --cover ./cover.png \
  --publish-channel local \
  --client-id <clientId>
```

## 常见工作流入口

- 通用规则：`references/workflows/common-rules.md`
- 图文发布：`references/workflows/publish-imageText.md`
- 视频发布：`references/workflows/publish-video.md`
- 文章发布：`references/workflows/publish-article.md`

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
