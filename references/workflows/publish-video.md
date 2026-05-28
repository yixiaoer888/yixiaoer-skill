# 视频发布工作流

> 适用范围：抖音视频、快手视频、B 站视频、视频号视频、微博视频等。
> 阅读本文档前，请先阅读 [common-rules.md](./common-rules.md)。

---

## 推荐入口

优先使用 flags 模式：

```bash
yxer publish video <platform> \
  --account "<账号名或ID>" \
  --title "<标题>" \
  --description "<描述>" \
  --video ./clip.mp4 \
  --cover ./cover.png
```

只有在需要高级平台字段时，才回退到 `payload.json` 模式。

## 执行顺序

1. 查询账号：`yxer accounts [platform]`
2. 上传视频：`yxer upload <视频路径或URL>`
3. 上传封面：`yxer upload <封面路径或URL>`
4. 按需查询分类、位置、音乐、合集、话题、商品
5. 查阅对应平台文档：`docs/publish/video/`
6. 执行校验：`yxer validate <platform> video <payload.json>`
7. 正式发布：`yxer publish video <platform> <payload.json>`

## 关键规则

- 发布前必须确认目标账号 `status=1`
- 视频只能有一个，封面必须单独上传
- 用户未提供封面时，必须补问，不要自动截帧
- 可选复杂对象必须通过查询命令取得完整对象后再填入
- 用户明确要求本机发布时，必须显式传 `--publish-channel local` 和 `--client-id`

## 直发示例

```bash
yxer publish video 抖音 \
  --account "视频账号" \
  --title "视频标题" \
  --description "视频描述" \
  --video ./clip.mp4 \
  --cover ./cover.png
```

## 本机发布示例

```bash
yxer publish video 抖音 \
  --account "视频账号" \
  --title "视频标题" \
  --description "视频描述" \
  --video ./clip.mp4 \
  --cover ./cover.png \
  --publish-channel local \
  --client-id <clientId>
```

## 平台文档入口

- 索引：`docs/publish/video/index.md`
- 平台细节：`docs/publish/video/*.md`
