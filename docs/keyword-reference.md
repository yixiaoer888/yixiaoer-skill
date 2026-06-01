# yxer 关键词文档

本文档用于统一业务沟通、实施交付和 AI 调用时的关键词口径。

## 1. 产品与组件关键词

| 关键词 | 含义 |
| --- | --- |
| `yxer` | 蚁小二 CLI，可执行入口 |
| `yixiaoer skill` | 给 AI agent 读取的技能包 |
| `CLI` | 命令行工具，当前项目即 `yxer` |
| `skill` | AI 技能规则，不直接执行 API |
| `linked app` | 宿主可识别的链接应用状态 |

## 2. 安装与运维关键词

| 关键词 | 推荐理解 |
| --- | --- |
| 下载 CLI | 获取 `yxer.exe`，来源可以是源码构建或制品包 |
| 安装 CLI | 放置 `yxer.exe` 到本地目录并加入 `PATH` |
| 安装 skill | 执行 `npx skills add "<repo>\\skills\\yixiaoer" -y` |
| 同步 skill | 执行 `yxer skill sync`，让 skill 版本与 CLI 对齐 |
| 升级 CLI | 重新构建或替换新版 `yxer.exe` |
| 卸载 CLI | 删除 `yxer.exe`、移除 `PATH`、按需清理 `.yxer` 配置 |
| 卸载 skill | 删除 AI 宿主中的 `yixiaoer` 技能安装，并清理 `skills.stamp` |

## 3. 发布能力关键词

| 关键词 | 含义 |
| --- | --- |
| `publish` | 正式发布 |
| `validate` | 发布前校验 payload |
| `dry-run` | 只预演，不正式提交 |
| `prepare` | 获取平台发布表单定义 |
| `schema get` | 获取标准 schema |
| `upload` | 上传图片、视频等资源 |

## 4. 内容类型关键词

| 关键词 | 含义 |
| --- | --- |
| `video` | 视频发布 |
| `imageText` | 图文发布 |
| `article` | 文章发布 |

说明：

- 当前推荐发布类型只使用这三种
- 文档、命令、payload 中都应保持这三个标准写法

## 5. 发布通道关键词

| 关键词 | 含义 |
| --- | --- |
| `cloud` | 云发布，默认模式 |
| `local` | 本机发布 / 本地发布 / 客户端发布 |
| `clientId` | 本机发布时必须提供的客户端标识 |

同义词归并建议：

- 本机发布 = 本地发布 = 客户端发布 = `publishChannel=local`
- 云发布 = 线上代理发布 = `publishChannel=cloud`

## 6. Payload 结构关键词

| 关键词 | 含义 |
| --- | --- |
| `publishArgs` | 标准发布请求体主体 |
| `accountForms` | 账号维度表单数组，必须位于 `publishArgs` 下 |
| `contentPublishForm` | 单账号的平台业务表单 |
| `platformAccountId` | 平台账号 ID |
| `platforms` | 本次目标平台列表 |

统一约束：

- 顶层必须包含 `publishArgs`
- `accountForms` 只能放在 `publishArgs.accountForms`
- 平台字段默认放在 `contentPublishForm`

## 7. 查询类关键词

| 关键词 | 含义 |
| --- | --- |
| `accounts list` | 查询账号 |
| `categories` | 查询分类 |
| `locations` | 查询地点 |
| `music` | 查询音乐 |
| `goods` | 查询商品 |
| `collections` | 查询合集 |
| `challenges` | 查询话题 |
| `records list` | 查询发布记录 |

## 8. 使用场景关键词

| 用户说法 | 标准动作 |
| --- | --- |
| 查账号 | `yxer accounts list` |
| 发视频 | `yxer publish video` |
| 发图文 | `yxer publish imageText` |
| 发文章 | `yxer publish article` |
| 先试一下 | `--dry-run` |
| 本机发 | `--publish-channel local` |
| 云端发 | 默认 `cloud` |
| 看字段 | `yxer prepare` + `yxer schema get` |

## 9. 推荐培训口径

对外说明时，建议统一这样表述：

- `yxer` 是实际执行工具
- `yixiaoer skill` 是给 AI 读的规则包
- 发布前必须先 `prepare`、`schema get`、`validate`
- 资源必须先上传，不能直接在 payload 里乱填外链
