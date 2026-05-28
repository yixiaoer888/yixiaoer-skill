# 文章发布工作流

> 适用范围：百家号文章、头条号文章、公众号文章、知乎文章、搜狐号文章等。
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)。

---

## 推荐入口

优先使用 flags 模式：

```bash
yxer publish article <platform> \
  --account "<账号名或ID>" \
  --title "<标题>" \
  --content @./article.html \
  --cover ./cover.png
```

只有在需要复杂分类或高级平台字段时，才回退到 `payload.json` 模式。

## 执行顺序

1. 查询账号：`yxer accounts [platform]`
2. 上传封面：`yxer upload <封面路径或URL>`
3. 如正文含图片，先逐张上传并替换引用
4. 按需查询分类、位置、话题
5. 查阅对应平台文档：`docs/publish/article/`
6. 执行校验：`yxer validate <platform> article <payload.json>`
7. 正式发布：`yxer publish article <platform> <payload.json>`

## 关键规则

- 发布前必须确认目标账号 `status=1`
- 文章必须有封面，未提供时必须补问
- 文章正文中的图片不能直接引用外部 URL，必须先上传
- 文章分类通常必须选到叶子节点
- 用户明确要求本机发布时，必须显式传本机发布参数

## 直发示例

```bash
yxer publish article 知乎 \
  --account "知乎账号" \
  --title "文章标题" \
  --content @./article.html \
  --cover ./cover.png
```

## 本机发布示例

```bash
yxer publish article 百家号 \
  --account "文章账号" \
  --title "文章标题" \
  --content @./article.html \
  --cover ./cover.png \
  --publish-channel local \
  --client-id <clientId>
```

## 平台文档入口

- 索引：`docs/publish/article/index.md`
- 平台细节：`docs/publish/article/*.md`
