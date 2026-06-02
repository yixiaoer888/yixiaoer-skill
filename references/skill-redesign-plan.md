# yxer Skill 重构方案

## 目标

把 `yxer` skill 从“单入口 + 外挂参考文档”改成“可打包、自包含、强路由、可校验”的结构，减少 Agent 在安装态和运行态的路径漂移、步骤跳过和命令乱猜。

## 当前缺陷

### 1. 技能包不自包含

- 当前安装入口是 `skills/yixiaoer/SKILL.md`
- 真实规则却主要放在仓库根下的 `references/` 和 `docs/publish/`
- 当 skill 被复制到 `.agents/skills/yixiaoer/` 或外部环境时，`../../references/...` 很容易失效

结果：

- Agent 在运行态无法稳定读取 workflow
- `SKILL.md` 看似完整，实际依赖仓库上下文

### 2. 缺少 shared 层

当前把环境检查、账号检查、linked-app、本机客户端、同步安装、常见报错等内容分散在多个 workflow 里，没有类似 `lark-shared` 的统一入口。

结果：

- 相同规则在多个 workflow 重复
- 新增命令后容易漏改
- Agent 不容易形成一致的错误处理策略

### 3. workflow 偏步骤表，不够像决策树

当前 workflow 对以下场景的分支规则还不够强：

- 多账号命中时如何展示候选项
- 素材不完整时先补什么
- `validate` 失败后回退到 `schema`、`prepare`、`upload` 还是动态字段查询
- 只给模糊业务意图时如何自动补默认值、何时必须追问
- dry-run 成功但正式发布失败时如何分层排查

### 4. 命令示例不够统一

同一能力会同时出现位置参数版和 flags 版，容易诱发 Agent 走“CLI 兼容什么就随便拼什么”的路径。

### 5. 缺少完整性校验

当前没有自动检查：

- `SKILL.md` 引用的相对路径是否存在
- `skills/` 安装态和 `.agents/` 运行态是否都可解析
- workflow 引用的 `docs/publish/*` 是否齐全

## 目标目录

目标是把“Agent 执行必须依赖的资料”全部收进 skill 包内。

```text
skills/
  yixiaoer/
    SKILL.md
    references/
      yixiaoer-shared.md
      cli/
        command-reference.md
        skill-install.md
      workflows/
        account-selection.md
        draft-workflow.md
        local-vs-cloud.md
        material-workflow.md
        payload-sourcing.md
        publish-article.md
        publish-imageText.md
        publish-troubleshooting.md
        publish-video.md
      platforms/
        index.md
        article/
          index.md
          *.md
        imageText/
          index.md
          *.md
        video/
          index.md
          *.md
    assets/
      yixiaoer.png
    plugin.json
```

说明：

- `references/yixiaoer-shared.md` 对标飞书的 `lark-shared`
- `platforms/` 应承接现在 `docs/publish/` 中 Agent 真正需要的内容
- 保留仓库根级 `references/` 和 `docs/` 作为维护态文档可以接受，但 skill 不应再依赖它们才能运行

## SKILL.md 目标形态

`SKILL.md` 应变薄，只做四件事：

1. 定义适用范围和唯一执行入口 `yxer`
2. 强制要求先读 `references/yixiaoer-shared.md`
3. 按任务类型路由到对应 workflow
4. 给出少量高优先级门禁规则

不应继续承担：

- 过长的命令大全
- 大段 payload 字段说明
- 平台细节正文
- 安装态说明和维护态说明混写

## yixiaoer-shared 建议内容

建议新增 `references/yixiaoer-shared.md`，统一承接以下内容：

- 首次配置：`yxer doctor`、`yxer config init`、`yxer config get`
- linked-app：`status/connect/disconnect/toggle`
- 本机发布前置：`yxer config set-local-client-id`
- 技能同步：`yxer skill show`、`yxer skill sync`
- 常见环境故障分级
- 输出协议：账号候选、校验错误、dry-run 结果、发布结果如何向用户展示
- 危险动作协议：真正发布前必须先 `validate` + `publish --dry-run`

## workflow 重写建议

### 发布类 workflow

发布类 workflow 建议统一采用下面的分支结构：

1. 判断内容类型：`video | imageText | article`
2. 判断目标平台是否唯一
3. 判断发布通道：`cloud | local`
4. 判断素材是否齐全
5. 判断动态字段是否需要查询
6. 获取 `prepare` 和 `schema`
7. 生成 payload 骨架
8. 回填资源、动态字段和业务字段
9. `validate`
10. `publish --dry-run`
11. 正式 `publish`

每一步都要写清楚：

- 必须命令
- 允许自动补全的字段
- 必须追问用户的场景
- 失败后的回退目标

### 非发布类 workflow

- `account-selection.md`：要写清多账号命中、唯一账号自动选择、账号状态异常时的展示协议
- `local-vs-cloud.md`：要写清默认走云发布、本机发布的触发词、本机失败后的回退
- `payload-sourcing.md`：保留字段来源纪律，但增加“只修局部、不重写全量 payload”的修订流程
- `publish-troubleshooting.md`：按错误类别分桶，而不是只按命令顺序回溯

## 迁移阶段

### Phase 1

- 重写 `skills/yixiaoer/SKILL.md`
- 新增 `yixiaoer-shared.md`
- 保留对根级 `references/` 和 `docs/publish/` 的兼容链接

### Phase 2

- 把 `references/cli`、`references/workflows` 迁入 `skills/yixiaoer/references/`
- 把 `docs/publish/{article,imageText,video}` 迁入 `skills/yixiaoer/references/platforms/`
- 更新所有相对路径

### Phase 3

- 增加 lint/check 脚本，校验 `SKILL.md` 与所有引用文件
- 增加 `.agents` 运行态自检，验证安装后相对路径仍成立

## 自动检查建议

至少补两类检查：

1. 路径完整性检查
   - 遍历 `SKILL.md` 与 workflow 中的 Markdown 链接
   - 校验源目录和 `.agents` 目录下都能解析

2. 平台文档完整性检查
   - 校验 `article/imageText/video` 三类索引存在
   - 校验平台索引中的条目都能找到对应平台文档

## 本次落地范围

本轮先完成两件事：

1. 给出本重构方案文档
2. 将 `skills/yixiaoer/SKILL.md` 改为更薄的入口形态，先收紧职责和路由

真正的目录搬迁和 shared 文档拆分，适合在下一轮继续做。
