# 视频发布工作流

> 适用范围：抖音视频、快手视频、B 站视频、视频号视频、微博视频等。
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)。

---

## 推荐入口

仅支持 `payload.json` 模式。发布前先获取表单字段和 schema：

```bash
yxer prepare <platform> video
yxer schema fields <platform> video
yxer schema get <platform> video
```

开始前，先补读：

- [`account-selection.md`](./account-selection.md)
- [`local-vs-cloud.md`](./local-vs-cloud.md)
- [`payload-sourcing.md`](./payload-sourcing.md)
- 涉及话题/标签时：[`../topic-tags.md`](../topic-tags.md)

## 执行顺序

1. 查询账号：`yxer accounts list [platform] [--status 1] [--json]`
2. 获取前置数据：`yxer prepare <platform> video`
3. 先获取字段视图：`yxer schema fields <platform> video`；需要 payload 骨架时再执行 `yxer schema get <platform> video`
4. 上传视频：`yxer upload <视频路径或URL>`
5. 上传封面：`yxer upload <封面路径或URL>`；视频号封面必须使用 `yxer upload <封面路径或URL> --platform 视频号 --usage cover`，超 512KB 时由 CLI 内部压缩
6. 按需查询分类、位置、音乐、合集、话题、商品
7. 根据前置数据、schema 和字段来源纪律填写 `payload.json`
8. 查阅对应平台文档：`../platforms/video/`
9. 执行校验：`yxer validate <platform> video <payload.json>`
10. 正式发布：`yxer publish video <platform> <payload.json>`

## 关键规则

- 发布前必须确认目标账号 `status=1`
- 视频只能有一个，封面必须单独上传
- 用户未提供封面时，必须补问，不要自动截帧
- 可选复杂对象必须通过查询命令取得完整对象后再填入
- 话题/标签必须直接按 `../topic-tags.md` 的目标格式传入；不要依赖 CLI 从 `description` 自动改写
- 视频号原创声明必须按平台文档映射：用户说“勾选原创”“声明原创”或“开启原创”时，在 `publishArgs.accountForms[].contentPublishForm` 写 `createType: 1`；未提及或明确关闭/转载时写 `createType: 2`。不要把 `pubType` 当作原创开关
- 发布前先看 `prepare` 和 `schema fields` 返回的字段；只有要确认完整骨架时再看 `schema get`
- `payload.json` 必须使用统一标准结构：顶层 `publishArgs`，业务字段放在 `publishArgs.accountForms[].contentPublishForm`
- 用户明确要求本机发布时，必须显式传 `--publish-channel local` 和 `--client-id`

## 发布示例

```bash
yxer validate 抖音 video .\payload.json
yxer publish video 抖音 .\payload.json
```

## 本机发布示例

```bash
yxer publish video 抖音 .\payload.json --publish-channel local --client-id <clientId>
```

## 平台文档入口

- 索引：`../platforms/video/index.md`
- 平台细节：`../platforms/video/*.md`
