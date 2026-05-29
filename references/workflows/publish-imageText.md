# 图文发布工作流

> 适用范围：抖音图文、小红书笔记、知乎想法、快手图文、微博图文、微信视频号图文等。
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)。

---

## 推荐入口

仅支持 `payload.json` 模式。发布前先获取表单字段和 schema：

```bash
yxer prepare <platform> imageText
yxer schema get <platform> imageText
```

## 执行顺序

1. 查询账号：`yxer accounts [platform]`
2. 获取前置数据：`yxer prepare <platform> imageText`
3. 获取 schema：`yxer schema get <platform> imageText`
4. 逐张上传图片：`yxer upload <文件路径或URL>`
5. 按需查询分类、位置、音乐、合集、话题、商品
6. 根据前置数据与 schema 填写 `payload.json`
7. 查阅对应平台文档：`docs/publish/imageText/`
8. 执行校验：`yxer validate <platform> imageText <payload.json>`
9. 正式发布：`yxer publish imageText <platform> <payload.json>`

## 关键规则

- 发布前必须确认目标账号 `status=1`
- 每张图片都要单独上传，封面默认取第一张图
- 复杂对象必须通过查询命令获取完整对象，不要手写 `raw`
- 发布前先看 `prepare` 和 `schema get` 返回的表单字段，再填写 payload
- 多平台发布时，必须逐个平台分别执行完整流程
- 用户明确要求本机发布时，必须显式切换到 `local`

## 发布示例

```bash
yxer validate 小红书 imageText .\payload.json
yxer publish imageText 小红书 .\payload.json
```

## 本机发布示例

```bash
yxer publish imageText 抖音 .\payload.json --publish-channel local --client-id <clientId>
```

## 平台文档入口

- 索引：`docs/publish/imageText/index.md`
- 平台细节：`docs/publish/imageText/*.md`
