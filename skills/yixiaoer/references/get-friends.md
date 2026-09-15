# 获取好友/关联对象 (Get Friends)

用于查询可在小红书描述中艾特的好友对象。查询结果是发布时 `<friend>` 标签的唯一数据来源。

## 1. 调用指令

```bash
yxer query friends <account_id> [--red-id <小红书号>] --json
```

`account_id` 必须是小红书账号的蚁小二账号 ID。CLI 会将后端返回的好友列表原样放在 JSON 的 `data` 中；不要根据昵称或其他自然语言自行构造好友对象。

指定 `--red-id` 时，CLI 会按查询结果中好友对象的 `raw.red_id` 做精确筛选；该参数不会改变后端请求，也不会猜测后端筛选字段。

## 2. 在小红书描述中艾特好友

小红书的艾特内容写在 `contentPublishForm.description` 的 HTML 中，不是单独的 `friends` 字段：

```html
<p>记录生活 <friend raw='好友查询结果的完整 JSON 序列化字符串'>@好友名称</friend></p>
```

组装规则：

1. 执行 `yxer query friends <account_id> --json`，从 `data.list` 中选择目标好友。
2. 将选中的完整好友对象序列化为 JSON，写入 `<friend>` 的 `raw` 属性；不能只保留 ID、昵称或内部 `raw` 子对象。
3. 标签正文使用 `@` 加上查询结果中的 `yixiaoerName`，并将最终 HTML 放入 `description`。
4. 用同一份 payload 执行 `yxer validate` 和 `yxer publish --dry-run`；CLI 会保留 `<friend>` 标签，不会把它改写成普通文本。

`raw` 是 HTML 属性中的字符串，因此生成 JSON payload 时需要按 JSON 规则转义其中的双引号。上面的 HTML 仅展示结构，好友对象必须来自本次 CLI 查询。

## 3. 返回对象使用约束

- `data.list` 中的每一项都是一个完整好友对象，发布时整体透传。
- 查询结果为空时不要编造兜底对象，应让用户确认账号登录状态或好友关系。
- 返回多个好友时，应根据用户明确的昵称或稳定 ID 选择，不能默认选择第一项。
