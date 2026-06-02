# 文章发布工作流

> 适用范围：百家号文章、头条号文章、公众号文章、知乎文章、搜狐号文章等。
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)。

---

## 推荐入口

仅支持 `payload.json` 模式。发布前先获取表单字段和 schema：

```bash
yxer prepare <platform> article
yxer schema get <platform> article
```

开始前，先补读：

- [`account-selection.md`](./account-selection.md)
- [`local-vs-cloud.md`](./local-vs-cloud.md)
- [`payload-sourcing.md`](./payload-sourcing.md)

## 执行顺序

1. 查询账号：`yxer accounts list [platform] [--status 1] [--json]`
2. 获取前置数据：`yxer prepare <platform> article`
3. 获取 schema：`yxer schema get <platform> article`
4. 上传封面：`yxer upload <封面路径或URL>`
5. 如正文含图片，先逐张上传并替换引用
6. 按需查询分类、位置、话题
7. 根据前置数据、schema 和字段来源纪律填写 `payload.json`
8. 查阅对应平台文档：`../platforms/article/`
9. 执行校验：`yxer validate <platform> article <payload.json>`
10. 正式发布：`yxer publish article <platform> <payload.json>`

## 关键规则

- 发布前必须确认目标账号 `status=1`
- 文章必须有封面，未提供时必须补问
- 文章正文中的图片不能直接引用外部 URL，必须先上传
- 文章分类通常必须选到叶子节点
- 发布前先看 `prepare` 和 `schema get` 返回的表单字段，再填写 payload
- `payload.json` 必须使用统一标准结构：顶层 `publishArgs`
- 文章正文 `content` 应放在 `publishArgs.content`，与 `accountForms` 同级
- 账号和平台差异字段放在 `publishArgs.accountForms[].contentPublishForm`
- CLI 会在缺失时把 `publishArgs.content` 自动补齐到 `accountForms[].contentPublishForm.content`
- 用户明确要求本机发布时，必须显式传本机发布参数

## 发布示例

```bash
yxer validate 知乎 article .\payload.json
yxer publish article 知乎 .\payload.json
```

## 本机发布示例

```bash
yxer publish article 百家号 .\payload.json --publish-channel local --client-id <clientId>
```

## 平台文档入口

- 索引：`../platforms/article/index.md`
- 平台细节：`../platforms/article/*.md`
