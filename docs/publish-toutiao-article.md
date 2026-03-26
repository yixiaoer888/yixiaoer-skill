# 头条号文章发布 (Toutiao Article Publish)

该能力允许用户通过命令行将文章发布到“今日头条/头条号”平台。

## 场景描述
发布一篇包含标题、正文、封面的文章到指定的头条号账户。支持原创声明、广告投放等高级设置。

## 基础用法
```bash
node scripts/publish-toutiao-article.ts --title="文章标题" --content="<p>文章HTML内容</p>" --account_ids="account_id_1,account_id_2"
```

## 参数列表

| 参数 | 必填 | 说明 | 默认值 |
| :--- | :---: | :--- | :--- |
| `--title` | 是 | 文章标题 | - |
| `--content` | 是 | 文章 HTML 正文内容 | - |
| `--cover_key` | 否 | 封面图存储 Key (推荐使用 `upload-resource.ts` 获取) | - |
| `--cover_url` | 否 | 封面图外部 URL (若无 Key 则尝试同步上传) | - |
| `--is_first` | 否 | 是否原创 (true/false) | `false` |
| `--advertisement` | 否 | 广告投放 (2: 不投放, 3: 投放) | `2` |
| `--declaration` | 否 | 声明字段 (0: 无需声明) | `0` |
| `--pub_type` | 否 | 发布类型 (0: 草稿, 1: 公开) | `1` |
| `--account_ids` | 否 | 指定账号 ID 列表，逗号分隔 (若不指定则默认发布到该平台下第一个账号) | - |

## 发布流程
1. **上传资源**: 建议先使用 `upload-resource.ts` 上传封面图片，获取 `key`。
2. **执行发布**: 运行本脚本。
3. **内容存证**: 脚本会自动将内容存入“文章库”并获取 `publishContentId`。
4. **发起任务**: 调用 V2 发布接口。

## 示例
```bash
node scripts/publish-toutiao-article.ts \
  --title="今日头条测试文章" \
  --content="<p>这是一篇通过 OpenClaw 发布的测试文章。</p>" \
  --cover_key="material-library/2026/03/xxx.jpg" \
  --is_first=true \
  --advertisement=2
```
