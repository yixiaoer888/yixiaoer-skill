---
name: yixiaoer
version: 3.1.0
description: "通过 yxer CLI 操作蚁小二多平台内容分发：账号查询、资源上传、发布前准备、payload 校验、云发布/本机发布、草稿保存、素材登记、发布记录排查与技能同步。"
metadata:
  category: "productivity"
  requires:
    bins: ["yxer"]
  cliHelp: "yxer --help; yxer doctor; yxer accounts list --help; yxer publish --help; yxer validate --help"
---

# 蚁小二 Skill

你是 AI Agent，通过 `yxer` CLI 操作蚁小二资源。真正执行一律走 CLI，不要假设存在旧 Node 脚本入口、隐式 API 或手工脚本。

**🚀 首次使用？先读 [`./QUICKSTART.md`](./QUICKSTART.md) - 5 分钟完成首次发布**

**CRITICAL - 开始前必须先读取 [`./references/yixiaoer-shared.md`](./references/yixiaoer-shared.md)，其中包含环境检查、linked-app、发布通道、同步和输出协议。**

## 能力索引

根据用户需求，必须先读取对应业务域文档，再进入具体 workflow 或 reference。不要直接凭记忆拼 payload 或执行正式发布。

- 发布与 payload 修订
  - 入口：[`./references/domains/publish.md`](./references/domains/publish.md)
  - 覆盖视频、图文、文章发布，账号选择，云/本机通道判断，payload 来源纪律，动态字段查询，平台差异文档入口。
- 账号、环境与 linked-app
  - 入口：[`./references/domains/accounts-and-env.md`](./references/domains/accounts-and-env.md)
  - 覆盖 `doctor`、`config`、账号查询、linked-app 连接与技能同步。
- 草稿与素材库
  - 入口：[`./references/domains/draft-and-material.md`](./references/domains/draft-and-material.md)
  - 覆盖蚁小二草稿、平台草稿判断、素材上传、素材登记与“上传后立即发布”的切换路径。
- 发布记录与失败排查
  - 入口：[`./references/domains/troubleshooting.md`](./references/domains/troubleshooting.md)
  - 覆盖 `records list`、校验失败修复、本机/云发布错误分流与回退策略。
- 安装、升级与分发
  - 入口：[`./references/domains/install-and-sync.md`](./references/domains/install-and-sync.md)
  - 覆盖 skill 安装、同步、升级和宿主侧接入说明。

## 命令探索

```bash
yxer --help
yxer doctor
yxer <command> --help
yxer prepare <platform> <type>
yxer schema fields <platform> <type>
yxer schema get <platform> <type>
```

## 全局规则

- 发布、草稿、素材、排查都只允许通过 `yxer` CLI 执行。
- 正式发布前固定顺序是：`doctor -> accounts list -> prepare -> schema fields -> validate -> publish --dry-run -> publish`；只有需要 payload 骨架时再补 `schema get`。
- `prepare`、`schema fields` / `schema get`、workflow、平台文档和 CLI 实际输出，是组装 payload 的唯一依据。
- 图片、视频、封面等资源必须先上传，且只能复用 `yxer upload` 返回的真实字段。
- `category`、`location`、`music`、`collection`、`challenge`、`goods` 等动态字段必须先查询，不能手写对象。
- `validate`、`publish --dry-run`、正式 `publish` 必须使用同一套发布通道参数。
