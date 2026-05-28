# 图文发布工作流

> 适用范围：抖音图文、小红书笔记、知乎想法、快手图文、微博图文、微信视频号图文等。
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)。

---

## 推荐入口

优先使用 flags 模式：

```bash
yxer publish image-text <platform> \
  --account "<账号名或ID>" \
  --title "<标题>" \
  --description "<正文>" \
  --image ./1.jpg \
  --image ./2.jpg
```

只有在需要高级自定义字段时，才回退到 `payload.json` 模式。

## 执行顺序

1. 查询账号：`yxer accounts [platform]`
2. 逐张上传图片：`yxer upload <文件路径或URL>`
3. 按需查询分类、位置、音乐、合集、话题、商品
4. 查阅对应平台文档：`docs/publish/image-text/`
5. 执行校验：`yxer validate <platform> image-text <payload.json>`
6. 正式发布：`yxer publish image-text <platform> <payload.json>`

## 关键规则

- 发布前必须确认目标账号 `status=1`
- 每张图片都要单独上传，封面默认取第一张图
- 复杂对象必须通过查询命令获取完整对象，不要手写 `raw`
- 多平台发布时，必须逐个平台分别执行完整流程
- 用户明确要求本机发布时，必须显式切换到 `local`

## 直发示例

```bash
yxer publish image-text 小红书 \
  --account "图文账号" \
  --title "图文标题" \
  --description "图文正文" \
  --image ./1.jpg \
  --image ./2.jpg
```

## 本机发布示例

```bash
yxer publish image-text 抖音 \
  --account "图文账号" \
  --title "图文标题" \
  --description "图文正文" \
  --image ./1.jpg \
  --publish-channel local \
  --client-id <clientId>
```

## 平台文档入口

- 索引：`docs/publish/image-text/index.md`
- 平台细节：`docs/publish/image-text/*.md`
