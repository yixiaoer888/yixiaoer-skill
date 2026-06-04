# 发布与 Payload

适用范围：用户要发布视频、图文、文章，或要修订发布 payload、解释字段归属、确认发布通道。

**CRITICAL - 只要用户意图落在本域，MUST 先读取下方 workflow，再决定是否执行任何 `prepare`、字段查询、`validate` 或 `publish`。**
**BLOCKING REQUIREMENT - 未完成 `doctor`、账号确认、`prepare`、`schema fields`、`validate`、`publish --dry-run` 前，绝对禁止正式 `publish`。**

## 读取顺序

1. [`../workflows/common-rules.md`](../workflows/common-rules.md)
2. [`../workflows/account-selection.md`](../workflows/account-selection.md)
3. [`../workflows/local-vs-cloud.md`](../workflows/local-vs-cloud.md)
4. [`../workflows/payload-sourcing.md`](../workflows/payload-sourcing.md)
5. 按类型继续读取：
   - 图文：[`../workflows/publish-imageText.md`](../workflows/publish-imageText.md)
   - 视频：[`../workflows/publish-video.md`](../workflows/publish-video.md)
   - 文章：[`../workflows/publish-article.md`](../workflows/publish-article.md)

## 意图路由

- 用户只说“帮我发一下”“发个抖音/小红书”“发视频/图文/文章”时，直接进入本域，并继续按类型读取 workflow。
- 用户明确说“先别发，只生成 payload”“帮我修 payload”时，仍进入本域，但最终动作只停留在 payload 修订或 `validate` / `publish --dry-run`。
- 用户明确说“查为什么发失败了”“解释报错”时，先切 [`./troubleshooting.md`](./troubleshooting.md)，不要直接重试正式发布。
- 用户明确说“上传素材后马上发”时，先完成素材流程，再回切本域继续执行发布主流程。

## 平台差异入口

- 总索引：[`../platforms/index.md`](../platforms/index.md)
- 视频平台：[`../platforms/video/index.md`](../platforms/video/index.md)
- 图文平台：[`../platforms/imageText/index.md`](../platforms/imageText/index.md)
- 文章平台：[`../platforms/article/index.md`](../platforms/article/index.md)

只有在 `prepare` / `schema get` 之后，且当前 workflow 无法回答平台差异时，才继续读取具体平台文档。

## 强制门禁

- 未执行 `yxer doctor` 不进入发布流程
- 未确认 `accounts list` 中账号 `status=1` 不继续
- 未执行 `prepare` / `schema fields` 不组装 payload；只有需要 payload 骨架时再补 `schema get`
- 未先 `validate` 与 `publish --dry-run` 不执行正式 `publish`

## 常用命令

```bash
yxer accounts list [platform] [--name 关键词] [--status 1] [--json]
yxer prepare <platform> <type>
yxer schema get <platform> <type>
yxer upload --file <file_path>
yxer upload --url <resource_url>
yxer categories <account_id> [--type video|article]
yxer locations <account_id> [--query 关键词]
yxer music <account_id> [--query 关键词]
yxer goods <account_id> [--query 关键词]
yxer collections <account_id> [--type video|article]
yxer challenges <account_id> [--query 关键词] [--type video]
yxer validate <platform> <type> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
yxer publish <type> <platform> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>] --dry-run
yxer publish <type> <platform> <payload.json> [--publish-channel cloud|local] [--client-id <clientId>]
```

## 决策提示

- 用户只说“帮我发”时，默认云发布；明确说“本机发布”“客户端发布”时切到本机通道。
- 用户要“只生成 payload”时，仍要先走 `prepare` / `schema get` 和字段查询纪律。
- 用户要填分类、位置、音乐、合集、话题、商品时，先查询，再回填完整对象。
