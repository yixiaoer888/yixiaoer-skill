# yxer CLI 安装指南

以下步骤面向 AI Agent，部分步骤需要用户提供本地环境信息或蚁小二凭证配合完成。

仓库地址：

- <https://github.com/yixiaoer888/yixiaoer-skill/tree/feature/yxrcli>

## 环境要求

开始安装之前，请确保环境中已安装：

- Node.js（需包含 `npx`，用于安装 skill）
- Go `1.25.0+`（用于从源码构建 `yxer` CLI）

## 第 1 步 获取源码

```shell
git clone https://github.com/yixiaoer888/yixiaoer-skill.git
cd yixiaoer-skill
git checkout feature/yxrcli
```

如果用户已经在本地提供了该仓库目录，可以直接进入仓库根目录，无需重复克隆。

## 第 2 步 构建 CLI

在仓库根目录执行：

```shell
go build -o bin/yxer.exe .
```

如果构建成功，可执行文件位于：

```text
bin\yxer.exe
```

如果用户已经提供了现成的 `yxer.exe`，也可以跳过源码构建，直接把该可执行文件加入 `PATH` 后继续后续步骤。

## 第 3 步 验证 CLI

```shell
.\bin\yxer.exe --version
```

如果已经加入 `PATH`，也可以直接执行：

```shell
yxer --version
```

## 第 4 步 初始化配置

Agent 运行以下命令，并让用户提供蚁小二 `apiKey`。

```shell
.\bin\yxer.exe config init --api-key <apiKey>
```

如果需要同时绑定链接应用，可让用户补充账号名后执行：

```shell
.\bin\yxer.exe config init --api-key <apiKey> --bind-app --account-name <账号名>
```

## 第 5 步 安装 CLI Skill

Agent 在仓库根目录执行以下命令安装 `yixiaoer` skill：

```shell
npx skills add ".\skills\yixiaoer" -y
```

如需全局安装：

```shell
npx skills add ".\skills\yixiaoer" -g -y
```

也可以直接通过 CLI 同步 skill：

```shell
.\bin\yxer.exe skill sync
```

如需全局同步：

```shell
.\bin\yxer.exe skill sync --global
```

## 第 6 步 验证安装结果

```shell
.\bin\yxer.exe doctor
.\bin\yxer.exe skill show
```

如果已经加入 `PATH`，也可以执行：

```shell
yxer doctor
yxer skill show
```

## 第 7 步 交给 AI Agent 使用

安装完成后，AI Agent 应按以下入口读取能力，而不是直接猜命令：

1. 先读 `skills/yixiaoer/SKILL.md`
2. 再按任务进入 `skills/yixiaoer/references/domains/`
3. 需要具体命令时，再调用 `yxer`

常见任务入口：

- 发布任务：`skills/yixiaoer/references/domains/publish.md`
- 账号与环境：`skills/yixiaoer/references/domains/accounts-and-env.md`
- 安装与同步：`skills/yixiaoer/references/domains/install-and-sync.md`
- 排障：`skills/yixiaoer/references/domains/troubleshooting.md`

## 常用验证命令

```shell
yxer config get
yxer accounts list 抖音 --json
yxer update --check
```

## 补充说明

- 安装顺序固定为“先安装 CLI，再安装 Skill”
- `Skill` 负责告诉 Agent 如何路由任务和调用命令，不负责真正执行发布
- `yxer` 是唯一执行入口
- 当 `yxer` 版本、`SKILL.md` 或 `references/` 文档更新后，建议重新执行 `yxer skill sync`
