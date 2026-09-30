# 头条号文章发布参数 (TouTiaoHao Article)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [文章发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

## 触发场景 (Trigger)
- **意图辨析**：用户指定在“今日头条 (头条号)”发布长文章或新闻资讯，且需要配置如“广告收益 (Advertisement)”、“头条首发”、“地点挂载”或“创作者申明”等头条特有商业与分发逻辑时。
- **典型提示词**：
  - “发一篇今日头条的文章，开启广告收益模式”
  - “头条文章发布，声明是我的原创首发内容”
  - “在头条文章里加上我所在的物理位置”
  - “帮我同步这个图文到头条号并存为草稿”

## 执行逻辑 (Logic Flow)
1. **封面适配**：头条号对封面 (`covers`) 数量有严格要求（1 张或 3 张）。
2. **位置检索**：若需精准挂载位置，调用 `locations` 接口获取 DTO。
3. **收益装配**：根据意图注入 `advertisement: 3`（赚取收益）或 `2`（不投放广告）。
4. **参数装配**：构造 `accountForms[i].contentPublishForm`。
5. **指令执行**：先执行 `yxer validate <platform> <type> <payload.json>`，再执行 `yxer publish <type> <platform> <payload.json> [--publish-channel local --client-id <clientId>]`。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 (1-50 字符) | - |
| `content` | `string` | **是** | 文章内容 (HTML 字符串，最多 50000 字符) | - |
| `covers` | `Array` | **是** | 文章封面列表 (`OldCover[]`, 1-9 张) | - |
| `isFirst` | `boolean` | 否 | 是否头条首发 | `false` |
| `location` | `Object` | 否 | 位置对象 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (13 位 Unix 时间戳，单位: 毫秒) | - |
| `advertisement` | `number` | 否 | 广告投放收益: 2-无收益, 3-投放广告赚收益 | `3` |
| `declaration` | `number` | 否 | 创作声明：`0`-不声明，`1`-取材网络，`2`-个人观点，仅供参考，`3`-引用 AI，`6`-虚构演绎/故事经历，`7`-投资观点，仅供参考，`8`-健康医疗分享，仅供参考；`9`-独家精选（需 `right_exclusive_selection`），`10`-首发精选（需 `right_first_selection`） | `0` |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |

普通声明值 `0/1/2/3/6/7/8` 不受精选权益影响；`9/10` 只有账号信息返回对应权益时才填写，否则使用 `0` 或省略。CLI 接收 `declaration` 数字，并会在 `validate`、`publish`（含 `--dry-run`）前查询 `/platform-accounts/{id}/account-info`，拒绝无对应权益的值；服务端会将其转换为 `statement.type` 透传给头条号。`isFirst=true` 表示“头条首发”，与 `declaration=10` 的“首发精选”是两个独立字段。

## 2. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "article",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_tt_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "头条号自动化发布演示",
          "content": "<h1>演示</h1><p>正文内容...</p>",
          "covers": [
            { "key": "cover_key_001", "width": 800, "height": 600, "size": 150000 }
          ],
          "advertisement": 3,
          "isFirst": true,
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. 复杂对象结构说明

### 3.1 OldCover
包含 `key`, `size`, `width`, `height`。

### 3.2 PlatformDataItem (位置信息)
包含 `yixiaoerId`, `yixiaoerName`, `raw` (必须完整透传)。

## 相关接口

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `location`  | `locations` | [获取位置信息](../../get-locations.md) |
| `covers.key`| `upload`    | [资源上传](../../upload-resource.md) |
